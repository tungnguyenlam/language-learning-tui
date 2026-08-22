package tui

import (
	tea "charm.land/bubbletea/v2"
)

// aiScreen wraps the AI view to satisfy the screen interface.
type aiScreen struct{}

func (aiScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderAI(layout.X, layout.Y)
}

func (aiScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.searchingAI {
		switch msg.String() {
		case "enter", "\r", "\n":
			m.searchingAI = false
			return m.startDrafting(), true
		case "esc", "escape":
			m.searchingAI = false
			m.aiInput = ""
			m.draftSource = ""
			return nil, true
		case "backspace":
			if len(m.aiInput) > 0 {
				m.aiInput = trimLastRune(m.aiInput)
				if len(m.aiInput) == 0 {
					m.draftSource = ""
				}
			}
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.aiInput += ch
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "/":
		m.searchingAI = true
		return nil, true
	case "esc", "escape":
		if m.drafting {
			m.drafting = false
			m.draftCancelled = true
			m.status = "AI drafting cancelled"
			return nil, true
		}
		if m.explanation != "" || m.explainError != "" || m.explainingCard {
			m.explanation = ""
			m.explainError = ""
			m.explainingCard = false
			m.explainCardID = ""
			m.status = "Explanation dismissed"
			return nil, true
		}
		m.aiInput = ""
		m.draftSource = ""
		return nil, true
	case "H", "h":
		if m.explanation != "" || m.explainError != "" {
			m.explanation = ""
			m.explainError = ""
			m.explainingCard = false
			m.explainCardID = ""
			m.status = "Explanation dismissed"
			return nil, true
		}
	case "[":
		m.previousAITemplate()
		return nil, true
	case "]":
		m.nextAITemplate()
		return nil, true
	case "up", "k":
		if m.draftCursor > 0 {
			m.draftCursor--
		}
		return nil, true
	case "down", "j":
		if m.draftCursor < len(m.drafts)-1 {
			m.draftCursor++
		}
		return nil, true
	case "A", "ctrl+a":
		if len(m.drafts) > 0 {
			return m.approveAllDrafts(), true
		}
	case "a":
		if len(m.drafts) > 0 {
			return m.approveDraft(), true
		}
	case "D", "ctrl+d":
		if len(m.drafts) > 0 {
			return m.discardAllDrafts(), true
		}
	case "d":
		if len(m.drafts) > 0 {
			m.discardDraft()
			return nil, true
		}
	case "enter", "\r", "\n":
		if len(m.drafts) > 0 {
			return m.approveDraft(), true
		}
		return m.startDrafting(), true
	case "backspace":
		if len(m.drafts) > 0 {
			m.discardDraft()
			return nil, true
		}
	}
	return nil, false
}
