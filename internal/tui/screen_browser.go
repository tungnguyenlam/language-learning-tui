package tui

import (
	tea "charm.land/bubbletea/v2"
)

// browserScreen wraps the Browser view to satisfy the screen interface.
type browserScreen struct{}

func (browserScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderBrowserAt(layout)
}

func (browserScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateBrowserKey(msg)
}
