// Package ankiweb browses and downloads decks from AnkiWeb's public shared
// deck library.
//
// AnkiWeb publishes no API. The endpoints used here are the ones its own web
// client calls, and their responses are protobuf with no published schema, so
// the field numbers below were read off live responses and are decoded
// defensively: a changed or added field yields an incomplete result rather than
// an error, and every call has a timeout.
//
// This is the only part of the app that talks to the network, and it only does
// so when the user explicitly searches or downloads. Nothing is sent about the
// user, and no credentials are involved — the shared library is public.
package ankiweb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BaseURL is AnkiWeb's origin. It is a variable so tests can point the client
// at a local server.
var BaseURL = "https://ankiweb.net"

// MaxDownloadBytes caps a deck download. Shared decks with audio run to
// hundreds of megabytes; without a cap a mis-click could fill the disk.
const MaxDownloadBytes = 512 << 20

// Deck is one entry of the shared deck library.
type Deck struct {
	ID         int64
	Title      string
	ThumbsUp   int
	ThumbsDown int
	Updated    time.Time
	Notes      int
	Cards      int
}

// Details is the extra information the deck page shows, fetched on demand.
type Details struct {
	Title       string
	Tags        string
	Description string
	SizeBytes   int64
	Notes       int
	Cards       int

	// downloadKey is a short-lived token AnkiWeb requires on the download URL.
	downloadKey string
}

// Client talks to AnkiWeb. The zero value is not usable; call New.
type Client struct {
	http *http.Client
}

// New returns a client with sensible timeouts.
func New() *Client {
	return &Client{http: &http.Client{Timeout: 60 * time.Second}}
}

// NewWithHTTPClient returns a client using the given HTTP client.
func NewWithHTTPClient(h *http.Client) *Client {
	return &Client{http: h}
}

// ErrNotAvailable reports that AnkiWeb could not be reached or answered with
// something this client does not understand. Callers should treat the shared
// deck browser as unavailable rather than failing the app.
var ErrNotAvailable = errors.New("ankiweb is unavailable")

// Search returns shared decks matching a query, ranked by AnkiWeb's own
// ordering (most popular first).
func (c *Client) Search(ctx context.Context, query string) ([]Deck, error) {
	body, err := c.get(ctx, "/svc/shared/list-decks?search="+url.QueryEscape(query))
	if err != nil {
		return nil, err
	}

	fields, err := pbFields(body)
	if err != nil && len(fields) == 0 {
		return nil, fmt.Errorf("%w: unreadable search response", ErrNotAvailable)
	}

	var decks []Deck
	for _, f := range fields {
		if f.Num != 1 || f.Wire != wireBytes {
			continue
		}
		entry, err := pbFields(f.Bytes)
		if err != nil && len(entry) == 0 {
			continue
		}
		deck := Deck{
			ID:         int64(pbUint(entry, 1)),
			Title:      strings.TrimSpace(pbString(entry, 2)),
			ThumbsUp:   int(pbUint(entry, 3)),
			ThumbsDown: int(pbUint(entry, 4)),
			Notes:      int(pbUint(entry, 6)),
			Cards:      int(pbUint(entry, 7)),
		}
		if ts := pbUint(entry, 5); ts > 0 {
			deck.Updated = time.Unix(int64(ts), 0)
		}
		if deck.ID == 0 || deck.Title == "" {
			continue
		}
		decks = append(decks, deck)
	}
	return decks, nil
}

// Info fetches a shared deck's description and the token needed to download it.
func (c *Client) Info(ctx context.Context, id int64) (Details, error) {
	body, err := c.get(ctx, fmt.Sprintf("/svc/shared/item-info?sharedId=%d", id))
	if err != nil {
		return Details{}, err
	}

	top, err := pbFields(body)
	if err != nil && len(top) == 0 {
		return Details{}, fmt.Errorf("%w: unreadable deck info", ErrNotAvailable)
	}
	inner, ok := pbSubMessage(top, 1)
	if !ok {
		return Details{}, fmt.Errorf("%w: deck info has no payload", ErrNotAvailable)
	}
	fields, err := pbFields(inner)
	if err != nil && len(fields) == 0 {
		return Details{}, fmt.Errorf("%w: unreadable deck info", ErrNotAvailable)
	}

	details := Details{
		Title:       strings.TrimSpace(pbString(fields, 5)),
		Tags:        strings.TrimSpace(pbString(fields, 6)),
		SizeBytes:   int64(pbUint(fields, 7)),
		Description: strings.TrimSpace(pbString(fields, 9)),
	}
	if stats, ok := pbSubMessage(fields, 10); ok {
		if sf, err := pbFields(stats); err == nil || len(sf) > 0 {
			details.Notes = int(pbUint(sf, 1))
			details.Cards = int(pbUint(sf, 2))
			details.downloadKey = pbString(sf, 5)
		}
	}
	return details, nil
}

// Download writes a shared deck's `.apkg` to w. It needs the Details returned
// by Info, whose download token AnkiWeb expires after a short while.
//
// progress, when non-nil, is called with the number of bytes written so far.
func (c *Client) Download(ctx context.Context, id int64, details Details, w io.Writer, progress func(int64)) (int64, error) {
	if details.downloadKey == "" {
		return 0, fmt.Errorf("%w: no download token; refresh the deck info", ErrNotAvailable)
	}

	req, err := c.request(ctx, fmt.Sprintf("/svc/shared/download-deck/%d?t=%s",
		id, url.QueryEscape(details.downloadKey)))
	if err != nil {
		return 0, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrNotAvailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, downloadError(resp)
	}

	written, err := copyWithProgress(w, io.LimitReader(resp.Body, MaxDownloadBytes+1), progress)
	if err != nil {
		return written, fmt.Errorf("downloading deck: %w", err)
	}
	if written > MaxDownloadBytes {
		return written, fmt.Errorf("deck is larger than the %d MB download limit", MaxDownloadBytes>>20)
	}
	return written, nil
}

// ErrAnonymousLimit reports that AnkiWeb refused the request because it allows
// only a handful of searches and downloads per address before asking for an
// account. It is a normal outcome rather than a fault, and the deck can still
// be fetched by hand, so callers should say so rather than showing a failure.
var ErrAnonymousLimit = errors.New("AnkiWeb limits anonymous use")

// ErrLinkExpired reports that a download token has gone stale. Fetching the
// deck info again produces a fresh one.
var ErrLinkExpired = errors.New("download link expired")

// classifyError maps AnkiWeb's plain-text error bodies onto the sentinels
// above so the UI can distinguish "try later" from "something is broken".
func classifyError(resp *http.Response, body string) error {
	text := strings.TrimSpace(body)
	lower := strings.ToLower(text)

	switch {
	case strings.Contains(lower, "log in"), resp.StatusCode == http.StatusTooManyRequests:
		return ErrAnonymousLimit
	case strings.Contains(lower, "expired"), strings.Contains(lower, "link is invalid"):
		return ErrLinkExpired
	}
	if text == "" {
		text = resp.Status
	}
	return fmt.Errorf("%w: %s", ErrNotAvailable, text)
}

func downloadError(resp *http.Response) error {
	message, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	return classifyError(resp, string(message))
}

// DeckURL is the page a user can open to download a deck by hand, which is the
// fallback when ErrDownloadLimit is returned.
func DeckURL(id int64) string {
	return fmt.Sprintf("%s/shared/info/%d", BaseURL, id)
}

// copyWithProgress copies src to dst, reporting cumulative bytes written.
func copyWithProgress(dst io.Writer, src io.Reader, progress func(int64)) (int64, error) {
	buf := make([]byte, 256<<10)
	var total int64
	var lastReport time.Time
	for {
		n, err := src.Read(buf)
		if n > 0 {
			written, werr := dst.Write(buf[:n])
			total += int64(written)
			if werr != nil {
				return total, werr
			}
			// Reporting every chunk would flood a TUI with redraws.
			if progress != nil && time.Since(lastReport) > 200*time.Millisecond {
				lastReport = time.Now()
				progress(total)
			}
		}
		if err == io.EOF {
			if progress != nil {
				progress(total)
			}
			return total, nil
		}
		if err != nil {
			return total, err
		}
	}
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	req, err := c.request(ctx, path)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotAvailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return nil, classifyError(resp, string(message))
	}
	// Responses are small protobuf messages; the cap guards against a redirect
	// to something unexpected.
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}

func (c *Client) request(ctx context.Context, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("User-Agent", "deutsch-tui")
	return req, nil
}
