package tui

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"deutsch-tui/internal/core"
)

func TestHighlightQuery(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)

	tests := []struct {
		text     string
		query    string
		expected string
	}{
		{
			text:     "Deutsch",
			query:    "de",
			expected: style.Render("De") + "utsch",
		},
		{
			text:     "Apfelkuchen",
			query:    "kuch",
			expected: "Apfel" + style.Render("kuch") + "en",
		},
		{
			text:     "banana",
			query:    "an",
			expected: "b" + style.Render("an") + style.Render("an") + "a",
		},
		{
			text:     "banana",
			query:    "",
			expected: "banana",
		},
		{
			text:     "Äpfel und Apfel",
			query:    "äpf",
			expected: style.Render("Äpf") + "el und Apfel",
		},
		{
			text:     "Käsekuchen",
			query:    "KÄSE",
			expected: style.Render("Käse") + "kuchen",
		},
		{
			text:     "Straße",
			query:    "ß",
			expected: "Stra" + style.Render("ß") + "e",
		},
	}

	for _, tc := range tests {
		got := highlightQuery(tc.text, tc.query, style)
		if got != tc.expected {
			t.Errorf("highlightQuery(%q, %q) = %q, expected %q", tc.text, tc.query, got, tc.expected)
		}
	}
}

func TestRenderDictionary(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryResults = []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Apfel",
			Translation: "Apple",
			WordClass:   "Noun",
			Gender:      "r",
			Forms:       "Apfel, Äpfel",
			Examples:    []string{"Ein roter Apfel."},
		},
		{
			ID:          "2",
			Word:        "Birne",
			Translation: "Pear",
		},
	}
	m.dictionarySearch = "Ap"
	m.dictionaryCursor = 0

	// Test wide layout (two columns)
	m.width = 100
	m.height = 40
	layout := m.activeViewContentLayout()
	view := m.renderDictionary(layout)
	plainView := stripANSI(view)

	if !strings.Contains(plainView, "Apfel") {
		t.Errorf("expected view to contain Apfel")
	}
	if !strings.Contains(plainView, "Apple") {
		t.Errorf("expected view to contain Apple")
	}
	if !strings.Contains(plainView, "Forms:") {
		t.Errorf("expected view to contain detail panel with Forms. Plain view:\n%s", plainView)
	}
	if !strings.Contains(plainView, "Äpfel") {
		t.Errorf("expected view to contain Äpfel")
	}

	// Test compact layout (single column)
	m.width = 60
	layout = m.activeViewContentLayout()
	view = m.renderDictionary(layout)
	plainView = stripANSI(view)

	if strings.Contains(plainView, "Forms:") {
		t.Errorf("expected compact view NOT to contain detail panel")
	}
}

func TestDictionaryDetailScrollClampsToVisibleRows(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.breakpoint = BreakpointWide
	m.dictionaryDetailTotalLines = 40

	maxScroll := maxInt(0, m.dictionaryDetailTotalLines-dictionaryVisibleRows(m.activeViewContentLayout()))
	for i := 0; i < 100; i++ {
		cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModShift})
		if cmd != nil {
			t.Fatal("shift+down should not return a command")
		}
		if !handled {
			t.Fatal("shift+down should be handled")
		}
	}

	if m.dictionaryDetailScroll != maxScroll {
		t.Fatalf("dictionaryDetailScroll = %d, want %d", m.dictionaryDetailScroll, maxScroll)
	}
}

func TestDictionarySearchHistory(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.recordDictionarySearch("Apfel")
	m.recordDictionarySearch("Birne")
	m.recordDictionarySearch("Apfel") // Duplicate, should be moved to the end (most recent)

	if len(m.dictionarySearchHistory) != 2 {
		t.Fatalf("expected search history to have 2 items, got %d", len(m.dictionarySearchHistory))
	}
	if m.dictionarySearchHistory[0] != "Birne" || m.dictionarySearchHistory[1] != "Apfel" {
		t.Errorf("unexpected history order: %v", m.dictionarySearchHistory)
	}

	// Navigate away to test recording search on view transition
	m.activeView = ViewDictionary
	m.dictionarySearch = "Banane"
	m.updateView(ViewDashboard)

	if len(m.dictionarySearchHistory) != 3 {
		t.Fatalf("expected search history to have 3 items, got %d", len(m.dictionarySearchHistory))
	}
	if m.dictionarySearchHistory[2] != "Banane" {
		t.Errorf("expected Banane at end of search history, got %v", m.dictionarySearchHistory)
	}
}

type captureRepo struct {
	mockRepo
	upsertedNote *core.Note
}

func (r *captureRepo) UpsertNote(ctx context.Context, note core.Note) error {
	r.upsertedNote = &note
	return nil
}

func (r *captureRepo) UpsertDeck(ctx context.Context, deck core.Deck) error {
	return nil
}

func (r *captureRepo) GetDeck(ctx context.Context, id string) (core.Deck, error) {
	return core.Deck{ID: id, Name: "Dictionary"}, nil
}

func TestDictionaryQuickAddCardsGenerated(t *testing.T) {
	repo := &captureRepo{}
	m := NewModel(repo, &mockScheduler{})

	entry := core.DictionaryEntry{
		ID:          "test-1",
		Word:        "Auto",
		Translation: "Car",
		Gender:      "n",
	}

	cmd := m.addDictionaryEntryCmd(entry)
	if cmd == nil {
		t.Fatal("expected command from addDictionaryEntryCmd")
	}

	msg := cmd()
	status, ok := msg.(statusMsg)
	if !ok {
		t.Fatalf("expected statusMsg, got %T: %v", msg, msg)
	}

	if !strings.Contains(status.text, "Added") {
		t.Errorf("expected status text to say Added, got %q", status.text)
	}

	if repo.upsertedNote == nil {
		t.Fatal("expected note to be upserted")
	}

	note := repo.upsertedNote
	if note.Front != "Auto" || note.Back != "Car" {
		t.Errorf("unexpected note content: front=%q, back=%q", note.Front, note.Back)
	}

	if len(note.Cards) == 0 {
		t.Fatal("expected cards to be generated for the quick-added dictionary entry")
	}

	card := note.Cards[0]
	if card.Prompt != "Auto" || card.Answer != "Car" {
		t.Errorf("unexpected card content: prompt=%q, answer=%q", card.Prompt, card.Answer)
	}
}

func TestDictionarySearchClearHitbox(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.dictionarySearch = "Apfel"
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Apfel", Translation: "Apple"},
	}

	layout := m.activeViewContentLayout()
	_ = m.renderDictionary(layout)

	// Check if clear button hitbox is registered
	var foundClearHitbox bool
	for _, hb := range m.hitboxes {
		if hb.ID == "dict-search-clear" {
			foundClearHitbox = true
			if hb.Action == nil {
				t.Fatal("expected clear hitbox action to be set")
			}
			hb.Action() // Execute clear
			break
		}
	}

	if !foundClearHitbox {
		t.Fatal("expected dictionary clear button hitbox to be registered")
	}

	if m.dictionarySearch != "" {
		t.Errorf("expected search query to be cleared, got %q", m.dictionarySearch)
	}
	if len(m.dictionaryResults) != 0 {
		t.Errorf("expected results to be cleared, got %d items", len(m.dictionaryResults))
	}
}
