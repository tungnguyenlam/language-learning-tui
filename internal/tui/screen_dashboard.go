package tui

import (
	tea "charm.land/bubbletea/v2"
)

// dashboardScreen wraps the Dashboard view to satisfy the screen interface.
type dashboardScreen struct{}

func (dashboardScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderDashboard(layout)
}

func (dashboardScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateDashboardKey(msg)
}
