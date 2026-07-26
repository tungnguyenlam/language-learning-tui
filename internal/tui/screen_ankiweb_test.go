package tui

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"deutsch-tui/internal/ankiweb"

	tea "charm.land/bubbletea/v2"
)

// stubAnkiWeb stands in for the network client so these tests never leave the
// machine.
type stubAnkiWeb struct {
	decks     []ankiweb.Deck
	details   ankiweb.Details
	searchErr error
	infoErr   error
	dlErr     error
	payload   []byte

	searched string
	infoFor  int64
}

func (s *stubAnkiWeb) Search(ctx context.Context, query string) ([]ankiweb.Deck, error) {
	s.searched = query
	return s.decks, s.searchErr
}

func (s *stubAnkiWeb) Info(ctx context.Context, id int64) (ankiweb.Details, error) {
	s.infoFor = id
	return s.details, s.infoErr
}

func (s *stubAnkiWeb) Download(ctx context.Context, id int64, d ankiweb.Details, w io.Writer, p func(int64)) (int64, error) {
	if s.dlErr != nil {
		return 0, s.dlErr
	}
	n, err := w.Write(s.payload)
	return int64(n), err
}

func ankiWebModel(t *testing.T, stub *stubAnkiWeb) (*Model, *ankiWebScreen) {
	t.Helper()
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.ankiWeb = stub
	m.activeView = ViewAnkiWeb
	m.width, m.height = 120, 40
	return m, m.ankiWebScreen
}

func sampleDecks() []ankiweb.Deck {
	return []ankiweb.Deck{
		{ID: 74275356, Title: "German 360 - A1", ThumbsUp: 18, ThumbsDown: 1,
			Notes: 1199, Cards: 1191, Updated: time.Unix(1750000000, 0)},
		{ID: 293204297, Title: "Goethe Institute A1 Wordlist", ThumbsUp: 335, Cards: 800},
	}
}

// The browser must not reach the network until the user asks it to.
func TestAnkiWebSearchRequiresExplicitAction(t *testing.T) {
	stub := &stubAnkiWeb{decks: sampleDecks()}
	m, s := ankiWebModel(t, stub)

	m.updateActiveViewKey(tea.KeyPressMsg{Code: 'j'})
	m.updateActiveViewKey(tea.KeyPressMsg{Code: 'k'})
	if stub.searched != "" {
		t.Fatal("navigating must not trigger a network request")
	}

	// "/" opens the query editor; typing goes into the query, not shortcuts.
	if _, handled := m.updateActiveViewKey(tea.KeyPressMsg{Code: '/'}); !handled {
		t.Fatal("expected / to start editing the query")
	}
	if !s.editingQuery || !m.textInputActive() {
		t.Fatal("expected the query editor to count as active text input")
	}
	for _, r := range "german a1" {
		m.updateActiveViewKey(tea.KeyPressMsg{Code: r})
	}
	if s.query != "german a1" {
		t.Fatalf("query not captured: %q", s.query)
	}
	if stub.searched != "" {
		t.Fatal("typing must not trigger a search on every keystroke")
	}

	cmd, _ := m.updateActiveViewKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected Enter to issue a search command")
	}
	msgs := executeCmd(cmd)
	if stub.searched != "german a1" {
		t.Fatalf("expected the query to be searched, got %q", stub.searched)
	}

	for _, msg := range msgs {
		m.Update(msg)
	}
	if len(s.results) != 2 || s.results[0].Title != "German 360 - A1" {
		t.Fatalf("results not applied: %+v", s.results)
	}
}

func TestAnkiWebRendersResultsAndDetails(t *testing.T) {
	stub := &stubAnkiWeb{
		decks: sampleDecks(),
		details: ankiweb.Details{
			Title: "German 360 - A1", Tags: "German A1 Beginner",
			SizeBytes: 19805380, Notes: 1199, Cards: 1191,
			Description: "<p>This deck contains 1202 flashcards of <b>essential</b> German words.</p>",
		},
	}
	m, s := ankiWebModel(t, stub)
	s.query = "german a1"
	s.results = stub.decks

	layout := viewportLayout{X: 0, Y: 0, Width: 100, Height: 30}
	view := s.Render(m, layout)

	for _, want := range []string{"ANKIWEB SHARED DECKS", "German 360 - A1", "1191 cards", "+18/-1"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected %q in the result list, got:\n%s", want, view)
		}
	}

	// Enter loads details for the selected deck.
	cmd, _ := m.updateActiveViewKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	for _, msg := range executeCmd(cmd) {
		m.Update(msg)
	}
	if stub.infoFor != 74275356 {
		t.Fatalf("expected details for the selected deck, got %d", stub.infoFor)
	}

	view = s.Render(m, layout)
	if !strings.Contains(view, "19 MB") {
		t.Errorf("expected a human-readable download size, got:\n%s", view)
	}
	// Descriptions are HTML written by strangers; markup must not reach the UI.
	if strings.Contains(view, "<p>") || strings.Contains(view, "<b>") {
		t.Errorf("raw HTML leaked into the description:\n%s", view)
	}
	if !strings.Contains(view, "essential") {
		t.Errorf("expected the description text, got:\n%s", view)
	}
}

// AnkiWeb rate-limits anonymous downloads. That is expected, and the UI must
// point at the manual route rather than showing a raw failure.
func TestAnkiWebDownloadLimitExplainsTheFallback(t *testing.T) {
	stub := &stubAnkiWeb{decks: sampleDecks(), dlErr: ankiweb.ErrAnonymousLimit}
	m, s := ankiWebModel(t, stub)
	s.results = stub.decks

	cmd, _ := m.updateActiveViewKey(tea.KeyPressMsg{Code: 'd'})
	for _, msg := range executeCmd(cmd) {
		m.Update(msg)
	}

	if !strings.Contains(s.lastErr, "limits anonymous downloads") {
		t.Errorf("expected the limit to be named, got %q", s.lastErr)
	}
	if !strings.Contains(s.lastErr, "import it with I") {
		t.Errorf("expected the manual fallback to be offered, got %q", s.lastErr)
	}
	if !m.isErrorStatus {
		t.Error("expected the status line to be marked as an error")
	}

	if view := s.Render(m, viewportLayout{Width: 100, Height: 30}); !strings.Contains(view, "limits anonymous downloads") {
		t.Errorf("expected the error to be visible in the view:\n%s", view)
	}
}

func TestAnkiWebUnreachableServiceIsExplained(t *testing.T) {
	stub := &stubAnkiWeb{searchErr: ankiweb.ErrNotAvailable}
	m, s := ankiWebModel(t, stub)
	s.query = "german"

	cmd := s.search(m)
	for _, msg := range executeCmd(cmd) {
		m.Update(msg)
	}
	if !strings.Contains(s.lastErr, "unreachable") {
		t.Errorf("expected an offline-friendly message, got %q", s.lastErr)
	}
	if s.busy != "" {
		t.Error("expected the busy indicator to clear after a failure")
	}
}

func TestAnkiWebDownloadRejectsUnreadablePackage(t *testing.T) {
	stub := &stubAnkiWeb{decks: sampleDecks(), payload: []byte("not an apkg")}
	m, s := ankiWebModel(t, stub)
	s.results = stub.decks

	cmd, _ := m.updateActiveViewKey(tea.KeyPressMsg{Code: 'd'})
	for _, msg := range executeCmd(cmd) {
		m.Update(msg)
	}
	if s.lastErr == "" {
		t.Fatal("expected a corrupt download to be reported")
	}
	if !strings.Contains(strings.ToLower(s.lastErr), "import") {
		t.Errorf("expected the failing step to be named, got %q", s.lastErr)
	}
}

func TestAnkiWebEscapeReturnsToImport(t *testing.T) {
	m, _ := ankiWebModel(t, &stubAnkiWeb{})

	if _, handled := m.updateActiveViewKey(tea.KeyPressMsg{Code: tea.KeyEscape}); !handled {
		t.Fatal("expected Esc to be handled")
	}
	if m.activeView != ViewImport {
		t.Fatalf("expected Esc to return to Import, got %q", m.activeView)
	}
}

// The browser is reachable from Import and is deliberately not in the tab
// cycle, so the offline flow is never routed through a network view.
func TestAnkiWebIsReachableFromImportButNotInTabCycle(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewImport

	if _, handled := m.updateActiveViewKey(tea.KeyPressMsg{Code: 'A'}); !handled {
		t.Fatal("expected A to open the AnkiWeb browser from Import")
	}
	if m.activeView != ViewAnkiWeb {
		t.Fatalf("expected ViewAnkiWeb, got %q", m.activeView)
	}

	seen := map[View]bool{}
	m.activeView = ViewDashboard
	for i := 0; i < 20; i++ {
		executeCmd(m.nextViewCmd())
		seen[m.activeView] = true
	}
	if seen[ViewAnkiWeb] {
		t.Error("the AnkiWeb browser must not be part of the tab cycle")
	}
}

func TestAnkiWebEmptyStateExplainsItself(t *testing.T) {
	m, s := ankiWebModel(t, &stubAnkiWeb{})
	view := s.Render(m, viewportLayout{Width: 100, Height: 30})

	if !strings.Contains(view, "press / to search") {
		t.Errorf("expected the empty state to say how to search:\n%s", view)
	}
	if !strings.Contains(view, "nothing is uploaded") {
		t.Errorf("expected the empty state to state what leaves the machine:\n%s", view)
	}
}

func TestAnkiWebErrorTextCoversEachFailure(t *testing.T) {
	for _, tc := range []struct {
		name   string
		action string
		err    error
		want   string
	}{
		// The limit applies to searches and downloads alike, but only the
		// download has a useful manual fallback, so the wording differs.
		{"search limit", "search", ankiweb.ErrAnonymousLimit, "limits anonymous searches"},
		{"download limit", "download", ankiweb.ErrAnonymousLimit, "limits anonymous downloads"},
		{"expired link", "download", ankiweb.ErrLinkExpired, "expired"},
		{"offline", "search", ankiweb.ErrNotAvailable, "unreachable"},
		{"timeout", "search", context.DeadlineExceeded, "timed out"},
		{"other", "search", errors.New("boom"), "boom"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ankiWebErrorText(ankiWebErrorMsg{action: tc.action, err: tc.err})
			if !strings.Contains(got, tc.want) {
				t.Errorf("want %q in %q", tc.want, got)
			}
		})
	}
}

func TestHumanBytes(t *testing.T) {
	for _, tc := range []struct {
		in   int64
		want string
	}{
		{512, "512 B"},
		{2048, "2 KB"},
		{19805380, "19 MB"},
		{2 << 30, "2.0 GB"},
	} {
		if got := humanBytes(tc.in); got != tc.want {
			t.Errorf("humanBytes(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
