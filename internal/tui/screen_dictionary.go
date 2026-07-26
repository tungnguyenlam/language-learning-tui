package tui

import (
	tea "charm.land/bubbletea/v2"
)

// dictionaryScreen wraps the Dictionary view to satisfy the screen interface.
type dictionaryScreen struct{}

func (dictionaryScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderDictionary(layout)
}

func (dictionaryScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	return m.updateDictionaryKey(msg)
}
