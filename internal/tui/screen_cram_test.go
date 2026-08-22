package tui

import (
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestCramActiveSessionIgnoresGlobalNavigation(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.decks = []core.Deck{
		{ID: "deck-a", Name: "Deck A"},
		{ID: "deck-b", Name: "Deck B"},
	}
	model.deck = model.decks[0]
	model.activeView = ViewCram
	model.cramActive = true
	model.cramRevealed = false
	model.cramCards = []core.Card{
		{ID: "c1", DeckID: "deck-a", Prompt: "Haus", Answer: "house"},
	}

	for _, key := range []tea.KeyPressMsg{
		{Text: "s", Code: 's'},
		{Text: "w", Code: 'w'},
		{Text: "[", Code: '['},
		{Text: "]", Code: ']'},
	} {
		updated, _ := model.Update(key)
		model = updated.(*Model)
		if model.activeView != ViewCram {
			t.Fatalf("key %q left Cram, active view = %s", key.String(), model.activeView)
		}
		if !model.cramActive {
			t.Fatalf("key %q ended the Cram session", key.String())
		}
		if model.deck.ID != "deck-a" {
			t.Fatalf("key %q changed deck to %q", key.String(), model.deck.ID)
		}
	}
}

func TestCramScreenStartsSessionOnEnter(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewCram
	model.cramCards = []core.Card{{ID: "c1", Prompt: "Haus", Answer: "house"}}

	_, handled := (cramScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled {
		t.Fatal("enter was not handled on Cram setup")
	}
	if !model.cramActive {
		t.Fatal("enter should start a Cram session")
	}
}
