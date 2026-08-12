package tui

import (
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestDecksScreenEscapeClearsFilterAfterDeselect(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDecks
	model.decks = []core.Deck{{ID: "deck-1", Name: "German A1"}}
	screen := decksScreen{}

	for range 2 {
		if _, handled := screen.HandleKey(model, tea.KeyPressMsg{Text: "m"}); !handled {
			t.Fatal("deck selection shortcut was not handled")
		}
	}
	if len(model.deckSelected) != 0 {
		t.Fatalf("selection map contains deselected entries: %#v", model.deckSelected)
	}

	model.deckFilter = "German"
	if _, handled := screen.HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEsc}); !handled {
		t.Fatal("escape was not handled")
	}
	if model.deckFilter != "" {
		t.Fatalf("deck filter = %q after escape, want empty", model.deckFilter)
	}
}
