package tui

import (
	"errors"
	"strings"
	"testing"

	"deutsch-tui/internal/app"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestSettingsConfigPersistenceIsDeferredOrderedAndSnapshotted(t *testing.T) {
	var calls int
	var gotSpeed int
	var gotFront string
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		AIProviderName: "offline",
		AITemplates: map[string]map[string]string{
			"default": {"front": "before"},
		},
		RevealSpeed: 5,
		OnConfigChange: func(_ string, templates map[string]map[string]string, _ bool, _ bool, speed int) error {
			calls++
			gotSpeed = speed
			gotFront = templates["default"]["front"]
			return nil
		},
	})
	model.activeView = ViewSettings
	model.settingsCursor = settingsRevealSpeedItem

	staleCmd, handled := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "+"})
	if !handled || staleCmd == nil {
		t.Fatal("reveal-speed update should return a persistence command")
	}
	latestCmd, _ := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "+"})
	if calls != 0 {
		t.Fatalf("config callbacks during input handling = %d, want 0", calls)
	}
	model.aiTemplates["default"]["front"] = "after dispatch"

	executeCmd(staleCmd)
	executeCmd(latestCmd)
	if calls != 1 {
		t.Fatalf("config callbacks = %d, want only latest request", calls)
	}
	if gotSpeed != 7 || gotFront != "before" {
		t.Fatalf("persisted snapshot = speed %d, front %q; want 7 and before", gotSpeed, gotFront)
	}
}

func TestSettingsSecretPersistenceIsDeferredAndSnapshotted(t *testing.T) {
	var secretCalls int
	var configCalls int
	var gotSecrets app.Secrets
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		AIProviderName: "offline",
		OnSecretsChange: func(secrets app.Secrets) error {
			secretCalls++
			gotSecrets = secrets
			return nil
		},
		OnConfigChange: func(_ string, _ map[string]map[string]string, _ bool, _ bool, _ int) error {
			configCalls++
			return nil
		},
	})
	model.activeView = ViewSettings
	model.editingSecretProvider = "openai"
	model.editingSecretKey = "api_key"
	model.aiSecrets.OpenAI.APIKey = "saved-token"

	cmd, handled := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled || cmd == nil {
		t.Fatal("secret commit should return persistence commands")
	}
	if secretCalls != 0 || configCalls != 0 {
		t.Fatalf("callbacks during input handling = secrets %d, config %d; want 0", secretCalls, configCalls)
	}
	model.aiSecrets.OpenAI.APIKey = "later-token"
	executeCmd(cmd)
	if secretCalls != 1 || configCalls != 1 {
		t.Fatalf("callbacks after command = secrets %d, config %d; want 1 each", secretCalls, configCalls)
	}
	if gotSecrets.OpenAI.APIKey != "saved-token" {
		t.Fatalf("persisted API key = %q, want command snapshot", gotSecrets.OpenAI.APIKey)
	}
}

func TestSettingsPersistenceErrorsAreSurfaced(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
			AIProviderName: "offline",
			OnConfigChange: func(_ string, _ map[string]map[string]string, _ bool, _ bool, _ int) error {
				return errors.New("read-only filesystem")
			},
		})
		model.activeView = ViewSettings
		model.settingsCursor = settingsAutoPlayItem
		cmd, _ := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
		for _, msg := range executeCmd(cmd) {
			model.Update(msg)
		}
		if !model.isErrorStatus || !strings.Contains(model.status, "save settings") {
			t.Fatalf("status = %q (error=%v), want config persistence error", model.status, model.isErrorStatus)
		}
	})

	t.Run("secrets", func(t *testing.T) {
		model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
			AIProviderName: "openai",
			OnSecretsChange: func(app.Secrets) error {
				return errors.New("permission denied")
			},
		})
		model.activeView = ViewSettings
		model.editingSecretProvider = "openai"
		model.editingSecretKey = "api_key"
		cmd, _ := (settingsScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
		for _, msg := range executeCmd(cmd) {
			model.Update(msg)
		}
		if !model.isErrorStatus || !strings.Contains(model.status, "save AI credentials") {
			t.Fatalf("status = %q (error=%v), want secrets persistence error", model.status, model.isErrorStatus)
		}
	})
}

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
