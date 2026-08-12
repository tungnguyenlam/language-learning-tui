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
	key := msg.String()
	switch key {
	case "!", "@", "#":
		indices := map[string]int{"!": 0, "@": 1, "#": 2}
		idx := indices[key]
		if idx < len(m.recentDecks) {
			m.selectDeckByID(m.recentDecks[idx])
			return m.updateView(ViewReview), true
		}
		return nil, true
	case "g":
		return m.searchGrammarTipInBrowser(), true
	case "G":
		return m.lookupGrammarTipInDictionary(), true
	case "v":
		return m.practiceVerbOfTheDay(), true
	case "V":
		return m.lookupVerbOfTheDayInDictionary(), true
	case "w":
		return m.addWordOfTheDayToCollection(), true
	case "W":
		return m.lookupWordOfTheDayInDictionary(), true
	case "/":
		return m.updateView(ViewDictionary), true
	}
	return nil, false
}
