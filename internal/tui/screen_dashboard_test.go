package tui

import (
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestDashboardScreenRecentDeckShortcutSelectsDeckAndStartsReview(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDashboard
	model.decks = []core.Deck{
		{ID: "first", Name: "First"},
		{ID: "recent", Name: "Recent"},
	}
	model.recentDecks = []string{"recent"}

	_, handled := (dashboardScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "!"})
	if !handled {
		t.Fatal("recent-deck shortcut was not handled")
	}
	if model.deck.ID != "recent" {
		t.Fatalf("selected deck = %q, want recent", model.deck.ID)
	}
	if model.activeView != ViewReview {
		t.Fatalf("active view = %s, want %s", model.activeView, ViewReview)
	}
}

func TestDashboardScreenConsumesUnavailableRecentDeckShortcut(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})

	cmd, handled := (dashboardScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "#"})
	if !handled {
		t.Fatal("unavailable recent-deck shortcut should still be consumed")
	}
	if cmd != nil {
		t.Fatal("unavailable recent-deck shortcut returned a command")
	}
}
