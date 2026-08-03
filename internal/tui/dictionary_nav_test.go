package tui

import (
	"context"
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestDictionaryNavigation(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40

	// 1. Full Dictionary tab navigation when empty
	m.activeView = ViewDictionary
	m.dictionarySearch = ""
	m.dictionaryFocusResults = false

	// Pressing '1' should switch to Dashboard
	msg := tea.KeyPressMsg{Code: '1'}
	_, cmd := m.Update(msg)
	if m.activeView != ViewDashboard {
		t.Errorf("expected activeView to be Dashboard, got %s", m.activeView)
	}
	if cmd == nil {
		t.Errorf("expected a command for view update")
	}

	// 2. Dictionary Overlay navigation when empty
	m = NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40
	m.activeView = ViewReview
	m.dictionaryOverlayActive = true
	m.dictionarySearch = ""

	// Pressing '1' should switch to Dashboard AND close overlay
	msg = tea.KeyPressMsg{Code: '1'}
	_, cmd = m.Update(msg)
	if m.activeView != ViewDashboard {
		t.Errorf("expected activeView to be Dashboard from overlay, got %s", m.activeView)
	}
	if m.dictionaryOverlayActive {
		t.Errorf("expected dictionary overlay to be closed")
	}

	// 3. Dictionary search should still work for numbers if not empty
	m = NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40
	m.activeView = ViewDictionary
	m.dictionarySearch = "A"
	m.dictionaryFocusResults = false

	// Pressing '1' should add to search, NOT switch view
	msg = tea.KeyPressMsg{Code: '1'}
	m.Update(msg)
	if m.activeView != ViewDictionary {
		t.Errorf("expected activeView to stay Dictionary, got %s", m.activeView)
	}
	if m.dictionarySearch != "A1" {
		t.Errorf("expected dictionarySearch to be A1, got %s", m.dictionarySearch)
	}
}

func TestCleanLookupQueryAndCardFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<b>das Haus</b> (n.)", "das Haus"},
		{"{{c1::trinken::verb}} nach Hause", "trinken"},
		{"• sprechen (spricht, sprach, hat gesprochen)", "sprechen"},
		{"gefallen [Dat.]", "gefallen"},
		{"der Apfel / die Äpfel", "der Apfel"},
		{"Was bedeutet Haus?", "Was bedeutet Haus"},
	}

	for _, tt := range tests {
		got := cleanLookupQuery(tt.input)
		if got != tt.expected {
			t.Errorf("cleanLookupQuery(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}

	// Test formatDictionaryCardFront
	if got := formatDictionaryCardFront("Haus", "n"); got != "das Haus" {
		t.Errorf("expected 'das Haus', got %q", got)
	}
	if got := formatDictionaryCardFront("Apfel", "m"); got != "der Apfel" {
		t.Errorf("expected 'der Apfel', got %q", got)
	}
	if got := formatDictionaryCardFront("Katze", "f"); got != "die Katze" {
		t.Errorf("expected 'die Katze', got %q", got)
	}
	if got := formatDictionaryCardFront("der Tisch", "m"); got != "der Tisch" {
		t.Errorf("expected 'der Tisch', got %q", got)
	}
}

func TestDictionaryClozeUTF8Replacement(t *testing.T) {
	repo := &mockRepo{}
	m := NewModel(repo, &mockScheduler{})

	entry := core.DictionaryEntry{
		ID:          "test-utf8",
		Word:        "groß",
		Translation: "big",
		Examples:    []string{"Ein SEHR GROßER Baum steht dort."},
	}

	executeCmd(m.addDictionaryClozeEntryCmd(entry))

	note, err := repo.GetNote(context.Background(), "dict-cloze-test-utf8")
	if err != nil {
		t.Fatalf("GetNote error: %v", err)
	}
	expectedFront := "Ein SEHR {{c1::GROß}}ER Baum steht dort."
	if note.Front != expectedFront {
		t.Errorf("cloze front = %q, want %q", note.Front, expectedFront)
	}
}
