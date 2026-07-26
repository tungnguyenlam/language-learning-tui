package tui

import (
	tea "charm.land/bubbletea/v2"
)

// decksScreen wraps the Decks view to satisfy the screen interface.
type decksScreen struct{}

func (decksScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderDecks(layout)
}

func (decksScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateDecksKey(msg)
}
