package tui

import (
	"testing"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestSettingsScreenPlusMinusOnlyAdjustFocusedRow(t *testing.T) {
	model := NewModel(&mockRepo{dailyGoal: 10}, &mockScheduler{})
	model.stats.DailyGoal = 10
	model.revealSpeed = 5
	model.activeView = ViewSettings
	model.settingsCursor = 0

	cmd, handled := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "+"})
	if !handled {
		t.Fatal("plus was not handled on Settings")
	}
	if cmd != nil {
		t.Fatal("plus on the AI provider row started a command")
	}
	if model.stats.DailyGoal != 10 {
		t.Fatalf("daily goal changed from unfocused row: %d", model.stats.DailyGoal)
	}

	model.settingsCursor = settingsRevealSpeedItem
	_, handled = (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "+"})
	if !handled {
		t.Fatal("plus was not handled on Reveal Speed")
	}
	if model.revealSpeed != 6 {
		t.Fatalf("reveal speed = %d, want 6", model.revealSpeed)
	}
}

func TestAIScreenSlashEntersTopicSearch(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewAI

	_, handled := (aiScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "/"})
	if !handled {
		t.Fatal("slash was not handled on AI")
	}
	if !model.searchingAI {
		t.Fatal("slash should enter AI topic search")
	}
}

func TestDictionaryDetailSpaceScrolls(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDictionary
	model.width = 100
	model.height = 30
	model.dictionaryDetailView = true
	model.dictionaryFocusResults = true
	model.dictionaryDetailTotalLines = 80
	model.dictionaryDetailScroll = 0
	model.dictionaryResults = []core.DictionaryEntry{{ID: "1", Word: "Haus", Translation: "house"}}
	model.dictionaryCursor = 0

	_, handled := (dictionaryScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "space"})
	if !handled {
		t.Fatal("space was not handled in dictionary detail")
	}
	if model.dictionaryDetailScroll != 5 {
		t.Fatalf("detail scroll = %d, want 5", model.dictionaryDetailScroll)
	}
}
