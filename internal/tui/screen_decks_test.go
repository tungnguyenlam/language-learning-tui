package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestDeckSearchHistoryPersistenceRunsAsCommand(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDecks
	model.searchingDecks = true
	model.deckFilter = "German A1"

	cmd, handled := (decksScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled {
		t.Fatal("deck search submission was not handled")
	}
	if repo.setSettingCalls != 0 {
		t.Fatalf("SetSetting calls during input handling = %d, want 0", repo.setSettingCalls)
	}
	if cmd == nil {
		t.Fatal("deck search submission should return a persistence command")
	}

	executeCmd(cmd)
	if repo.setSettingCalls != 1 {
		t.Fatalf("SetSetting calls after command = %d, want 1", repo.setSettingCalls)
	}
	saved, err := repo.GetSetting(context.Background(), "deck_search_history")
	if err != nil {
		t.Fatalf("get saved history: %v", err)
	}
	if !strings.Contains(saved, "German A1") {
		t.Fatalf("saved history = %q, want submitted query", saved)
	}
}

func TestDeckSearchHistoryPersistenceKeepsLatestSnapshot(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})

	staleCmd := model.recordDeckSearch("German A1")
	latestCmd := model.recordDeckSearch("German A2")
	executeCmd(staleCmd)
	executeCmd(latestCmd)

	if repo.setSettingCalls != 1 {
		t.Fatalf("SetSetting calls = %d, want only latest request", repo.setSettingCalls)
	}
	saved, err := repo.GetSetting(context.Background(), "deck_search_history")
	if err != nil {
		t.Fatalf("get saved history: %v", err)
	}
	if !strings.Contains(saved, "German A1") || !strings.Contains(saved, "German A2") {
		t.Fatalf("saved history = %q, want latest complete snapshot", saved)
	}
}

func TestDeckSearchHistoryPersistenceErrorIsSurfaced(t *testing.T) {
	repo := &mockRepo{errSetSetting: errors.New("disk full")}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDecks
	model.searchingDecks = true
	model.deckFilter = "German A1"

	cmd, _ := (decksScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
	msgs := executeCmd(cmd)
	if len(msgs) != 1 {
		t.Fatalf("persistence messages = %d, want 1", len(msgs))
	}
	model.Update(msgs[0])
	if !model.isErrorStatus || !strings.Contains(model.status, "save deck search history") {
		t.Fatalf("status = %q (error=%v), want contextual persistence error", model.status, model.isErrorStatus)
	}
}

func TestDeckSearchHistoryMouseClearReturnsPersistenceCommand(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDecks
	model.width, model.height = 100, 40
	model.searchingDecks = true
	model.deckSearchHistory = []string{"German A1"}
	_ = model.renderDecks(model.activeViewContentLayout())

	var clear Hitbox
	for _, hit := range model.hitboxes {
		if hit.ID == "deck-history-clear" {
			clear = hit
			break
		}
	}
	if clear.Action == nil {
		t.Fatal("deck history clear hitbox not registered")
	}
	cmd := clear.Action()
	if repo.setSettingCalls != 0 {
		t.Fatalf("SetSetting calls during hitbox action = %d, want 0", repo.setSettingCalls)
	}
	if len(model.deckSearchHistory) != 0 {
		t.Fatalf("history after clear = %#v, want empty", model.deckSearchHistory)
	}
	executeCmd(cmd)
	if repo.setSettingCalls != 1 {
		t.Fatalf("SetSetting calls after clear command = %d, want 1", repo.setSettingCalls)
	}
}

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
