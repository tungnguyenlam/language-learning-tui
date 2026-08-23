package tui

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"deutsch-tui/internal/ankiweb"
	"deutsch-tui/internal/content"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ankiWebScreen browses AnkiWeb's public shared deck library and imports a deck
// straight into the collection.
//
// It is the only view that reaches the network, and only when the user searches
// or downloads. It is reachable with "A" from Import rather than living in the
// tab cycle, so the offline-first flow is never routed through it.
type ankiWebScreen struct {
	query        string
	editingQuery bool

	results []ankiweb.Deck
	cursor  int
	scroll  int

	// details is the expanded description for the deck under the cursor, and
	// carries the short-lived token a download needs.
	details      *ankiweb.Details
	detailsForID int64

	// Monotonic request ids so out-of-order network replies cannot clobber
	// a newer search or a different deck's details.
	searchID int
	infoID   int

	busy    string // non-empty while a request is in flight
	lastErr string
}

// --- messages ---------------------------------------------------------------

type ankiWebResultsMsg struct {
	id      int
	query   string
	results []ankiweb.Deck
}

type ankiWebInfoMsg struct {
	reqID   int
	id      int64
	details ankiweb.Details
}

type ankiWebErrorMsg struct {
	action string
	err    error
	id     int // searchID or infoID for the failed request; 0 = unscoped
}

// --- key handling -----------------------------------------------------------

func (s *ankiWebScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	if s.editingQuery {
		switch key {
		case "enter", "\r", "\n":
			s.editingQuery = false
			return s.search(m), true
		case "esc":
			s.editingQuery = false
			return nil, true
		case "backspace":
			if len(s.query) > 0 {
				s.query = trimLastRune(s.query)
			}
			return nil, true
		case "ctrl+u":
			s.query = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(key); ok {
			s.query += ch
			return nil, true
		}
		return nil, true
	}

	switch key {
	case "esc":
		return m.updateView(ViewImport), true
	case "/":
		s.editingQuery = true
		s.lastErr = ""
		return nil, true
	case "up", "k":
		if s.cursor > 0 {
			s.cursor--
			s.details = nil
		}
		return nil, true
	case "down", "j":
		if s.cursor < len(s.results)-1 {
			s.cursor++
			s.details = nil
		}
		return nil, true
	case "enter", "\r", "\n":
		return s.loadDetails(m), true
	case "d", "D":
		return s.download(m), true
	case "r":
		if strings.TrimSpace(s.query) != "" {
			return s.search(m), true
		}
		return nil, true
	}
	return nil, false
}

func (s *ankiWebScreen) selected() (ankiweb.Deck, bool) {
	if s.cursor < 0 || s.cursor >= len(s.results) {
		return ankiweb.Deck{}, false
	}
	return s.results[s.cursor], true
}

// --- commands ---------------------------------------------------------------

func (s *ankiWebScreen) search(m *Model) tea.Cmd {
	query := strings.TrimSpace(s.query)
	if query == "" {
		m.status = "Enter a search term first"
		return nil
	}

	s.searchID++
	reqID := s.searchID
	s.busy = "Searching AnkiWeb…"
	s.lastErr = ""
	m.status = s.busy

	client := m.ankiWebClient()
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		results, err := client.Search(ctx, query)
		if err != nil {
			return ankiWebErrorMsg{action: "search", err: err, id: reqID}
		}
		return ankiWebResultsMsg{id: reqID, query: query, results: results}
	}
}

func (s *ankiWebScreen) loadDetails(m *Model) tea.Cmd {
	deck, ok := s.selected()
	if !ok {
		return nil
	}
	if s.details != nil && s.detailsForID == deck.ID {
		return nil
	}

	s.infoID++
	reqID := s.infoID
	deckID := deck.ID
	s.busy = "Loading deck details…"
	s.lastErr = ""
	m.status = s.busy

	client := m.ankiWebClient()
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		details, err := client.Info(ctx, deckID)
		if err != nil {
			return ankiWebErrorMsg{action: "info", err: err, id: reqID}
		}
		return ankiWebInfoMsg{reqID: reqID, id: deckID, details: details}
	}
}

// download fetches the selected deck and imports it in one step. The download
// token expires quickly, so details are always refreshed first.
func (s *ankiWebScreen) download(m *Model) tea.Cmd {
	deck, ok := s.selected()
	if !ok {
		m.status = "No deck selected"
		return nil
	}

	s.busy = fmt.Sprintf("Downloading %q…", truncateLine(deck.Title, 40))
	s.lastErr = ""
	m.status = s.busy

	client := m.ankiWebClient()
	repo := m.repo
	return func() tea.Msg {
		// Shared decks with audio run to hundreds of megabytes, so allow a
		// generous window while still bounding the request.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		details, err := client.Info(ctx, deck.ID)
		if err != nil {
			return ankiWebErrorMsg{action: "info", err: err}
		}

		file, err := os.CreateTemp("", "ankiweb-*.apkg")
		if err != nil {
			return ankiWebErrorMsg{action: "download", err: err}
		}
		defer os.Remove(file.Name())
		defer file.Close()

		if _, err := client.Download(ctx, deck.ID, details, file, nil); err != nil {
			return ankiWebErrorMsg{action: "download", err: err}
		}

		notes, err := content.ImportAnkiAPKGFromFile(file.Name())
		if err != nil {
			return ankiWebErrorMsg{action: "import", err: err}
		}
		if len(notes) == 0 {
			return ankiWebErrorMsg{action: "import", err: errors.New("the package contained no notes")}
		}

		// Fall back to the AnkiWeb title when the package has no deck name.
		for i := range notes {
			if strings.TrimSpace(notes[i].DeckID) == "" {
				notes[i].DeckID = deck.Title
			}
		}

		for _, d := range content.DecksFromNotes(notes) {
			d.Description = "Imported from AnkiWeb: " + deck.Title
			if err := repo.UpsertDeck(ctx, d); err != nil {
				return ankiWebErrorMsg{action: "import", err: err}
			}
		}

		decks, err := repo.Decks(ctx)
		if err != nil {
			return ankiWebErrorMsg{action: "import", err: err}
		}
		cards, err := repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return ankiWebErrorMsg{action: "import", err: err}
		}
		return importDoneMsg{decks: decks, cards: cards, count: len(notes), path: "AnkiWeb: " + deck.Title}
	}
}

// handleAnkiWebMsg applies the browser's async results. It lives on Model
// because Update owns the message switch.
func (m *Model) handleAnkiWebMsg(msg tea.Msg) (bool, tea.Cmd) {
	s := m.ankiWebScreen
	switch msg := msg.(type) {
	case ankiWebResultsMsg:
		if msg.id != s.searchID {
			return true, nil
		}
		s.busy = ""
		s.results = msg.results
		s.cursor = 0
		s.scroll = 0
		s.details = nil
		if len(msg.results) == 0 {
			m.status = fmt.Sprintf("No AnkiWeb decks found for %q", msg.query)
		} else {
			m.status = fmt.Sprintf("Found %d AnkiWeb decks for %q", len(msg.results), msg.query)
		}
		return true, nil

	case ankiWebInfoMsg:
		if msg.reqID != s.infoID {
			return true, nil
		}
		if deck, ok := s.selected(); !ok || deck.ID != msg.id {
			return true, nil
		}
		s.busy = ""
		details := msg.details
		s.details = &details
		s.detailsForID = msg.id
		m.status = "Press d to download and import this deck"
		return true, nil

	case ankiWebErrorMsg:
		switch msg.action {
		case "search":
			if msg.id != 0 && msg.id != s.searchID {
				return true, nil
			}
		case "info":
			if msg.id != 0 && msg.id != s.infoID {
				return true, nil
			}
		}
		s.busy = ""
		s.lastErr = ankiWebErrorText(msg)
		m.status = s.lastErr
		m.isErrorStatus = true
		return true, nil
	}
	return false, nil
}

// ankiWebErrorText turns a failure into something a learner can act on. The
// anonymous-use limit is the common case and is not a fault, so it points at
// the manual route instead of reading like a breakage.
func ankiWebErrorText(msg ankiWebErrorMsg) string {
	switch {
	case errors.Is(msg.err, ankiweb.ErrAnonymousLimit):
		if msg.action == "search" {
			return "AnkiWeb limits anonymous searches. Try again later, or browse ankiweb.net and import the .apkg with I."
		}
		return "AnkiWeb limits anonymous downloads. Download the .apkg from ankiweb.net, then import it with I."
	case errors.Is(msg.err, ankiweb.ErrLinkExpired):
		return "That download link expired. Press Enter to reload the deck, then d again."
	case errors.Is(msg.err, ankiweb.ErrNotAvailable):
		return "AnkiWeb is unreachable. Check your connection, or import a downloaded .apkg with I."
	case errors.Is(msg.err, context.DeadlineExceeded):
		return "AnkiWeb timed out. Try again, or import a downloaded .apkg with I."
	}
	return fmt.Sprintf("AnkiWeb %s failed: %v", msg.action, msg.err)
}

// --- rendering --------------------------------------------------------------

func (s *ankiWebScreen) Render(m *Model, layout viewportLayout) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout = contentLayoutForStyle(style, layout.X, layout.Y)

	var b strings.Builder
	var lineY int
	write := func(s string) {
		b.WriteString(s)
		lineY += strings.Count(s, "\n")
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Render("ANKIWEB SHARED DECKS")
	write(title + "  " + mutedStyle.Render("public library — nothing is uploaded") + "\n\n")

	searchLabel := "Search: "
	if s.editingQuery {
		write(searchLabel + editStyle.Render(s.query+"_") + "\n\n")
	} else if s.query == "" {
		write(searchLabel + mutedStyle.Render("(press / to search, e.g. \"german a1\")") + "\n\n")
	} else {
		write(searchLabel + s.query + "\n\n")
	}

	if s.busy != "" {
		write(infoStyle.Render(s.busy) + "\n\n")
	} else if s.lastErr != "" {
		write(warnStyle.Render(wrapToWidth(s.lastErr, layout.Width)) + "\n\n")
	}

	if len(s.results) == 0 {
		if s.busy == "" && s.lastErr == "" {
			b.WriteString(mutedStyle.Render(
				"Browse decks shared by the Anki community and import one directly.\n"+
					"Downloaded decks are converted to this app's cards; your progress stays local.") + "\n")
		}
		b.WriteString("\n" + s.helpLine())
		return b.String()
	}

	// Reserve room for the chrome around the list: the detail block only costs
	// lines when a deck is expanded.
	chrome := 8
	showDetails := false
	if deck, ok := s.selected(); ok && s.details != nil && s.detailsForID == deck.ID {
		showDetails = true
		chrome = 15
	}

	// Keep the cursor inside the visible window.
	rows := maxInt(3, layout.Height-chrome)
	if s.cursor < s.scroll {
		s.scroll = s.cursor
	}
	if s.cursor >= s.scroll+rows {
		s.scroll = s.cursor - rows + 1
	}

	// Right-align the metadata into a column so the list scans vertically even
	// when titles contain wide or right-to-left characters.
	metaFor := func(deck ankiweb.Deck) string {
		size := fmt.Sprintf("%d cards", deck.Cards)
		if deck.Cards == 0 {
			size = fmt.Sprintf("%d notes", deck.Notes)
		}
		rating := fmt.Sprintf("+%d", deck.ThumbsUp)
		if deck.ThumbsDown > 0 {
			rating += fmt.Sprintf("/-%d", deck.ThumbsDown)
		}
		updated := ""
		if !deck.Updated.IsZero() {
			updated = deck.Updated.Format("2006-01")
		}
		return fmt.Sprintf("%12s  %8s  %7s", size, rating, updated)
	}

	metaWidth := 0
	for i := s.scroll; i < len(s.results) && i < s.scroll+rows; i++ {
		metaWidth = maxInt(metaWidth, lipgloss.Width(metaFor(s.results[i])))
	}
	nameWidth := maxInt(12, layout.Width-metaWidth-4)

	var resultLines []string
	for i := s.scroll; i < len(s.results) && i < s.scroll+rows; i++ {
		deck := s.results[i]
		marker := "  "
		nameStyle := lipgloss.NewStyle()
		if i == s.cursor {
			marker = "> "
			nameStyle = nameStyle.Bold(true).Foreground(colorPink)
		}

		name := truncateLine(deck.Title, nameWidth)
		pad := strings.Repeat(" ", maxInt(0, nameWidth-lipgloss.Width(name)))
		resultLines = append(resultLines, marker+nameStyle.Render(name)+pad+mutedStyle.Render(metaFor(deck)))
	}

	startY := layout.Y + lineY
	for i := s.scroll; i < len(s.results) && i < s.scroll+rows; i++ {
		idx := i
		rowY := startY + (i - s.scroll)
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("ankiweb-result-%d", idx),
			View:   ViewAnkiWeb,
			X:      layout.X,
			Y:      rowY,
			Width:  layout.Width,
			Height: 1,
			Action: func() tea.Cmd {
				s.cursor = idx
				s.details = nil
				return s.loadDetails(m)
			},
		})
	}

	resultLinesWithScroll := renderScrollbarColumn(resultLines, rows, len(s.results), s.scroll)
	for _, l := range resultLinesWithScroll {
		b.WriteString(l + "\n")
	}

	if len(s.results) > rows {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("  (%d of %d)", s.cursor+1, len(s.results))) + "\n")
	}
	b.WriteString("\n")

	if showDetails {
		b.WriteString(s.renderDetails(*s.details, layout.Width))
	}

	b.WriteString(s.helpLine())
	return b.String()
}

func (s *ankiWebScreen) renderDetails(details ankiweb.Details, width int) string {
	var b strings.Builder
	if details.Tags != "" {
		b.WriteString(mutedStyle.Render(truncateLine("Tags: "+details.Tags, width)) + "\n")
	}
	if details.SizeBytes > 0 {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("Download size: %s  •  %d notes, %d cards",
			humanBytes(details.SizeBytes), details.Notes, details.Cards)) + "\n")
	}
	if description := content.PlainTextFromHTML(details.Description); description != "" {
		// Descriptions run long; a few lines is enough to judge a deck by.
		lines := strings.Split(wrapToWidth(description, width), "\n")
		if len(lines) > 4 {
			lines = append(lines[:4], "…")
		}
		b.WriteString(strings.Join(lines, "\n") + "\n")
	}
	return b.String() + "\n"
}

func (s *ankiWebScreen) helpLine() string {
	if s.editingQuery {
		return infoStyle.Render("EDITING — Enter to search, Esc to cancel.")
	}
	return fmt.Sprintf("%s search  %s/%s move  %s details  %s download & import  %s back",
		keyStyle.Render("/"),
		keyStyle.Render("j"), keyStyle.Render("k"),
		keyStyle.Render("Enter"),
		keyStyle.Render("d"),
		keyStyle.Render("Esc"))
}

func humanBytes(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(n)/float64(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.0f MB", float64(n)/float64(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(n)/float64(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}

// wrapToWidth is a plain greedy word wrap for untrusted remote text.
func wrapToWidth(s string, width int) string {
	if width < 10 {
		width = 10
	}
	var out []string
	for _, paragraph := range strings.Split(s, "\n") {
		line := ""
		for _, word := range strings.Fields(paragraph) {
			switch {
			case line == "":
				line = word
			case len(line)+1+len(word) <= width:
				line += " " + word
			default:
				out = append(out, line)
				line = word
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// ankiWebClient returns the shared client, created on first use so the app
// never opens a network stack unless the browser is actually opened.
func (m *Model) ankiWebClient() ankiWebSearcher {
	if m.ankiWeb == nil {
		m.ankiWeb = ankiweb.New()
	}
	return m.ankiWeb
}

// ankiWebSearcher is the slice of the AnkiWeb client the TUI uses, so tests can
// substitute a stub without touching the network.
type ankiWebSearcher interface {
	Search(ctx context.Context, query string) ([]ankiweb.Deck, error)
	Info(ctx context.Context, id int64) (ankiweb.Details, error)
	Download(ctx context.Context, id int64, details ankiweb.Details, w io.Writer, progress func(int64)) (int64, error)
}
