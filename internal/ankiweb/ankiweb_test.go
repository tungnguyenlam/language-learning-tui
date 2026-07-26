package ankiweb

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- protobuf builders, so the fixtures are readable ------------------------

func pbTag(num, wire int) []byte {
	return binary.AppendUvarint(nil, uint64(num)<<3|uint64(wire))
}

func pbVarint(num int, v uint64) []byte {
	return append(pbTag(num, wireVarint), binary.AppendUvarint(nil, v)...)
}

func pbBytes(num int, b []byte) []byte {
	out := append(pbTag(num, wireBytes), binary.AppendUvarint(nil, uint64(len(b)))...)
	return append(out, b...)
}

func pbStr(num int, s string) []byte { return pbBytes(num, []byte(s)) }

// serve starts a stub AnkiWeb and points BaseURL at it for the test's duration.
func serve(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	old := BaseURL
	BaseURL = srv.URL
	t.Cleanup(func() { BaseURL = old })

	return NewWithHTTPClient(srv.Client())
}

// --- tests ------------------------------------------------------------------

func TestSearchDecodesDeckList(t *testing.T) {
	// Field numbers match live /svc/shared/list-decks responses.
	entry := func(id uint64, title string, up, down, updated, notes, cards uint64) []byte {
		var b []byte
		b = append(b, pbVarint(1, id)...)
		b = append(b, pbStr(2, title)...)
		b = append(b, pbVarint(3, up)...)
		b = append(b, pbVarint(4, down)...)
		b = append(b, pbVarint(5, updated)...)
		b = append(b, pbVarint(6, notes)...)
		b = append(b, pbVarint(7, cards)...)
		return b
	}

	var body []byte
	body = append(body, pbBytes(1, entry(74275356, "German 360 - A1", 18, 1, 1750000000, 1199, 1191))...)
	body = append(body, pbBytes(1, entry(293204297, "Goethe A1 Wordlist", 335, 9, 1640954352, 926, 800))...)

	var gotQuery string
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("search")
		w.Write(body)
	})

	decks, err := c.Search(context.Background(), "german a1")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if gotQuery != "german a1" {
		t.Errorf("expected the query to be sent, got %q", gotQuery)
	}
	if len(decks) != 2 {
		t.Fatalf("expected 2 decks, got %d", len(decks))
	}

	first := decks[0]
	if first.ID != 74275356 || first.Title != "German 360 - A1" {
		t.Errorf("unexpected first deck: %+v", first)
	}
	if first.ThumbsUp != 18 || first.ThumbsDown != 1 {
		t.Errorf("ratings not decoded: %+v", first)
	}
	if first.Notes != 1199 || first.Cards != 1191 {
		t.Errorf("counts not decoded: %+v", first)
	}
	if !first.Updated.Equal(time.Unix(1750000000, 0)) {
		t.Errorf("timestamp not decoded: %v", first.Updated)
	}
}

// AnkiWeb publishes no schema, so an added or reordered field must degrade
// rather than break the browser.
func TestSearchToleratesUnknownFields(t *testing.T) {
	var entry []byte
	entry = append(entry, pbVarint(1, 42)...)
	entry = append(entry, pbStr(2, "Ein Deck")...)
	entry = append(entry, pbStr(99, "a field this client has never seen")...)
	entry = append(entry, pbVarint(100, 12345)...)

	body := pbBytes(1, entry)
	// A trailing entry with no title must be dropped, not surfaced blank.
	body = append(body, pbBytes(1, pbVarint(1, 43))...)

	c := serve(t, func(w http.ResponseWriter, r *http.Request) { w.Write(body) })

	decks, err := c.Search(context.Background(), "x")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(decks) != 1 || decks[0].Title != "Ein Deck" {
		t.Fatalf("expected 1 usable deck, got %+v", decks)
	}
}

func TestInfoDecodesDetailsAndDownloadKey(t *testing.T) {
	stats := append(pbVarint(1, 4214), pbVarint(2, 4222)...)
	stats = append(stats, pbStr(5, "token-abc")...)

	var inner []byte
	inner = append(inner, pbStr(5, "Deutsch: 4000 German Words")...)
	inner = append(inner, pbStr(6, "German Vocabulary")...)
	inner = append(inner, pbVarint(7, 48168450)...)
	inner = append(inner, pbStr(9, "Over 4000 entries.")...)
	inner = append(inner, pbBytes(10, stats)...)

	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(pbBytes(1, inner))
	})

	details, err := c.Info(context.Background(), 653061995)
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if details.Title != "Deutsch: 4000 German Words" || details.Tags != "German Vocabulary" {
		t.Errorf("unexpected details: %+v", details)
	}
	if details.SizeBytes != 48168450 || details.Notes != 4214 || details.Cards != 4222 {
		t.Errorf("stats not decoded: %+v", details)
	}
	if details.downloadKey != "token-abc" {
		t.Errorf("download key not decoded: %q", details.downloadKey)
	}
}

func TestDownloadSendsTokenAndReportsProgress(t *testing.T) {
	payload := bytes.Repeat([]byte("apkg"), 5000)

	var gotToken, gotPath string
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotToken = r.URL.Path, r.URL.Query().Get("t")
		w.Write(payload)
	})

	var buf bytes.Buffer
	var lastProgress int64
	n, err := c.Download(context.Background(), 123, Details{downloadKey: "tok en/+"}, &buf,
		func(sofar int64) { lastProgress = sofar })
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	if gotPath != "/svc/shared/download-deck/123" {
		t.Errorf("unexpected path %q", gotPath)
	}
	if gotToken != "tok en/+" {
		t.Errorf("token not round-tripped through the query string: %q", gotToken)
	}
	if n != int64(len(payload)) || !bytes.Equal(buf.Bytes(), payload) {
		t.Errorf("payload not written intact: got %d of %d bytes", n, len(payload))
	}
	if lastProgress != n {
		t.Errorf("expected a final progress report of %d, got %d", n, lastProgress)
	}
}

func TestDownloadWithoutTokenFails(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("download should not have been attempted without a token")
	})
	if _, err := c.Download(context.Background(), 1, Details{}, &bytes.Buffer{}, nil); err == nil {
		t.Fatal("expected an error when the download token is missing")
	}
}

// AnkiWeb allows a few anonymous downloads per address, then asks for an
// account. That is a normal outcome the UI reports differently from a fault.
func TestAnonymousLimitIsDistinguished(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please log in to download more decks."))
	})

	_, err := c.Download(context.Background(), 1, Details{downloadKey: "t"}, &bytes.Buffer{}, nil)
	if !errors.Is(err, ErrAnonymousLimit) {
		t.Fatalf("expected ErrAnonymousLimit, got %v", err)
	}
}

func TestExpiredLinkIsDistinguished(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Download link is invalid or expired. Please reload the previous page."))
	})

	_, err := c.Download(context.Background(), 1, Details{downloadKey: "t"}, &bytes.Buffer{}, nil)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected an expiry error, got %v", err)
	}
	if errors.Is(err, ErrAnonymousLimit) {
		t.Fatal("an expired link is not a download limit")
	}
}

func TestSearchReportsUnavailableService(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	_, err := c.Search(context.Background(), "german")
	if !errors.Is(err, ErrNotAvailable) {
		t.Fatalf("expected ErrNotAvailable, got %v", err)
	}
}

func TestSearchRespectsContextCancellation(t *testing.T) {
	c := serve(t, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err := c.Search(ctx, "german"); err == nil {
		t.Fatal("expected the request to be cancelled")
	}
}

func TestDeckURL(t *testing.T) {
	old := BaseURL
	BaseURL = "https://ankiweb.net"
	defer func() { BaseURL = old }()

	if got := DeckURL(12345); got != "https://ankiweb.net/shared/info/12345" {
		t.Fatalf("unexpected deck URL %q", got)
	}
}

func TestProtobufReaderHandlesTruncation(t *testing.T) {
	full := pbStr(2, "hello")
	if _, err := pbFields(full[:len(full)-2]); err == nil {
		t.Fatal("expected truncated input to be reported")
	}
	// Fields decoded before the truncation are still returned.
	partial := append(pbVarint(1, 7), full[:len(full)-2]...)
	fields, _ := pbFields(partial)
	if pbUint(fields, 1) != 7 {
		t.Fatalf("expected the intact prefix to survive, got %+v", fields)
	}
}
