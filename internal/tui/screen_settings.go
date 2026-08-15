package tui

import (
	tea "charm.land/bubbletea/v2"
)

const (
	settingsDailyGoalItem   = 5
	settingsLastItem        = 16
	settingsRevealSpeedItem = 16
)

// settingsScreen wraps the Settings view to satisfy the screen interface.
type settingsScreen struct{}

func (settingsScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderSettings(layout.X, layout.Y)
}

func (settingsScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateSettingsKey(msg)
}
