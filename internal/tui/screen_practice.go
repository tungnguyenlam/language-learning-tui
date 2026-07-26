package tui

import (
	tea "charm.land/bubbletea/v2"
)

// practiceScreen wraps the Practice view to satisfy the screen interface.
type practiceScreen struct{}

func (practiceScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderPractice(layout)
}

func (practiceScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updatePracticeKey(msg)
}
