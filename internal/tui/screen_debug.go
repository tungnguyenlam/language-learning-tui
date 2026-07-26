package tui

import (
	tea "charm.land/bubbletea/v2"
)

// debugScreen wraps the Debug view to satisfy the screen interface.
type debugScreen struct{}

func (debugScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderDebug(layout.X, layout.Y)
}

func (debugScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if msg.String() == "ctrl+d" {
		return m.updateView(ViewDashboard), true
	}
	return nil, false
}
