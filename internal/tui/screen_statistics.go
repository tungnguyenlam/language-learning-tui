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
	switch msg.String() {
	case "up", "k":
		m.scrollStats(-1)
		return nil, true
	case "down", "j":
		m.scrollStats(1)
		return nil, true
	case "pgup":
		m.scrollStats(-10)
		return nil, true
	case "pgdown":
		m.scrollStats(10)
		return nil, true
	case "x":
		return m.exportStatsCSV(), true
	}
	return nil, false
}
