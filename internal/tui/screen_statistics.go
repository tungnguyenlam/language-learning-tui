package tui

import (
	tea "charm.land/bubbletea/v2"
)

// statisticsScreen wraps the Statistics view to satisfy the screen interface.
type statisticsScreen struct{}

func (statisticsScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderStatisticsAt(layout)
}

func (statisticsScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateStatisticsKey(msg)
}
