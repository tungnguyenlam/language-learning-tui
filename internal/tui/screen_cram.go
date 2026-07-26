package tui

import (
	tea "charm.land/bubbletea/v2"
)

// cramScreen wraps the Cram view to satisfy the screen interface.
type cramScreen struct{}

func (cramScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderCramAt(layout)
}

func (cramScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateCramKey(msg)
}
