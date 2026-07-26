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
	return m.updateAIKey(msg)
}
