package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/ai"

	tea "charm.land/bubbletea/v2"
)

const (
	settingsAIProviderItem      = 0
	settingsFrontTemplateItem   = 1
	settingsBackTemplateItem    = 2
	settingsExampleTemplateItem = 3
	settingsDailyGoalItem       = 4
	settingsAutoPlayItem        = 5
	settingsStrictNormItem      = 6
	settingsOpenAIKeyItem       = 7
	settingsOpenAIModelItem     = 8
	settingsOpenAIURLItem       = 9
	settingsAnthropicKeyItem    = 10
	settingsAnthropicModelItem  = 11
	settingsAnthropicURLItem    = 12
	settingsOllamaModelItem     = 13
	settingsOllamaURLItem       = 14
	settingsRevealSpeedItem     = 15
	settingsLastItem            = 15
)

func settingsCredAt(cursor int) (provider, key string, ok bool) {
	switch cursor {
	case settingsOpenAIKeyItem:
		return "openai", "api_key", true
	case settingsOpenAIModelItem:
		return "openai", "model", true
	case settingsOpenAIURLItem:
		return "openai", "base_url", true
	case settingsAnthropicKeyItem:
		return "anthropic", "api_key", true
	case settingsAnthropicModelItem:
		return "anthropic", "model", true
	case settingsAnthropicURLItem:
		return "anthropic", "base_url", true
	case settingsOllamaModelItem:
		return "ollama", "model", true
	case settingsOllamaURLItem:
		return "ollama", "base_url", true
	}
	return "", "", false
}

// settingsScreen wraps the Settings view to satisfy the screen interface.
type settingsScreen struct{}

func (settingsScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderSettings(layout.X, layout.Y)
}

func (settingsScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.editingSecretKey != "" {
		return m.handleSecretEditKey(msg)
	}
	if m.editingTemplate {
		activeSet := m.currentAITemplateSet()
		if activeSet == "" {
			m.editingTemplate = false
			m.originalTemplateValue = ""
			return nil, true
		}
		switch msg.String() {
		case "enter", "\r", "\n":
			m.editingTemplate = false
			m.originalTemplateValue = ""
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{
					Templates: m.aiTemplates,
					ActiveSet: activeSet,
				}
			}
			return m.persistConfig(), true
		case "esc":
			// Restore original value on cancel
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] = m.originalTemplateValue
			m.editingTemplate = false
			m.originalTemplateValue = ""
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{
					Templates: m.aiTemplates,
					ActiveSet: activeSet,
				}
			}
			return m.persistConfig(), true
		case "backspace":
			templateKey := m.templateKeyAtCursor()
			val := m.aiTemplates[activeSet][templateKey]
			if len(val) > 0 {
				m.aiTemplates[activeSet][templateKey] = trimLastRune(val)
			}
			return nil, true
		case "ctrl+u":
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] += ch
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return nil, true
	case "down", "j":
		if m.settingsCursor < settingsLastItem {
			m.settingsCursor++
		}
		return nil, true
	case "g":
		m.settingsCursor = 0
		return nil, true
	case "G":
		m.settingsCursor = settingsLastItem
		return nil, true
	case "enter":
		return m.handleSettingsEnter(), true
	case "+":
		switch m.settingsCursor {
		case settingsRevealSpeedItem:
			return m.setRevealSpeed(m.revealSpeed + 1), true
		case settingsDailyGoalItem:
			return m.setDailyGoal(m.stats.DailyGoal + 1), true
		}
		return nil, true
	case "-":
		switch m.settingsCursor {
		case settingsRevealSpeedItem:
			return m.setRevealSpeed(m.revealSpeed - 1), true
		case settingsDailyGoalItem:
			return m.setDailyGoal(m.stats.DailyGoal - 1), true
		}
		return nil, true
	case "[":
		m.previousAITemplate()
		return nil, true
	case "]":
		m.nextAITemplate()
		return nil, true
	}
	return nil, false
}

// handleSecretEditKey processes keys while the user is typing an API key,
// model name, or base URL into Settings. Like template editing, Enter
// commits the value and triggers a save; Esc reverts to the value we
// stashed in originalSecretValue. Backspace and ctrl+u edit the buffer.
func (m *Model) handleSecretEditKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	provider := m.editingSecretProvider
	key := m.editingSecretKey
	if provider == "" || key == "" {
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		return nil, true
	}

	commit := func() tea.Cmd {
		if m.aiProviderName == "disabled" || m.aiProviderName == "offline" || m.aiProviderName == "template" {
			if provider == "openai" && strings.TrimSpace(m.aiSecrets.OpenAI.APIKey) != "" {
				m.aiProviderName = "openai"
			} else if provider == "anthropic" && strings.TrimSpace(m.aiSecrets.Anthropic.APIKey) != "" {
				m.aiProviderName = "anthropic"
			}
		}

		// Rebuild the provider so the new credentials take effect immediately.
		m.aiProvider = buildProvider(m.aiProviderName, m.aiSecrets, m.aiTemplates, m.currentAITemplateSet())
		return tea.Batch(m.persistSecrets(), m.persistConfig())
	}

	switch msg.String() {
	case "enter", "\r", "\n":
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		m.originalSecretValue = ""
		cmd := commit()
		m.status = fmt.Sprintf("Saved %s %s", provider, key)
		return cmd, true
	case "esc":
		m.setCredValue(provider, key, m.originalSecretValue)
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		m.originalSecretValue = ""
		cmd := commit()
		m.status = "Edit cancelled"
		return cmd, true
	case "backspace":
		val := m.getCredValue(provider, key)
		if len(val) > 0 {
			m.setCredValue(provider, key, trimLastRune(val))
		}
		return nil, true
	case "ctrl+u":
		m.setCredValue(provider, key, "")
		return nil, true
	}
	if ch, ok := singlePrintableInput(msg.String()); ok {
		m.setCredValue(provider, key, m.getCredValue(provider, key)+ch)
		return nil, true
	}
	return nil, true
}
