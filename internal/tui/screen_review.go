package tui

import (
	tea "charm.land/bubbletea/v2"
)

// reviewScreen wraps the Review view to satisfy the screen interface.
type reviewScreen struct{}

func (reviewScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderReview(layout.X, layout.Y)
}

func (reviewScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateReviewKey(msg)
}
