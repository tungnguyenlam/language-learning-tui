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
	if m.searchingDecks {
		switch msg.String() {
		case "enter", "\r", "\n":
			cmd := m.recordDeckSearch(m.deckFilter)
			m.searchingDecks = false
			m.applyDeckFilter()
			filtered := m.filteredDecks()
			if len(filtered) > 0 {
				m.selectDeckByID(filtered[0].ID)
			}
			return cmd, true
		case "esc", "\x1b":
			m.searchingDecks = false
			m.applyDeckFilter()
			return nil, true
		case "ctrl+x":
			m.deckSearchHistory = nil
			return m.saveDeckHistory(), true
		case "backspace":
			if len(m.deckFilter) > 0 {
				m.deckFilter = trimLastRune(m.deckFilter)
				m.deckCursor = 0
			}
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.deckFilter += ch
			m.deckCursor = 0
			return nil, true
		}
		return nil, true
	}

	if m.editingDeckLimits {
		filtered := m.filteredDecks()
		if len(filtered) == 0 {
			m.editingDeckLimits = false
			return nil, false
		}
		m.deckCursor = clampInt(m.deckCursor, 0, len(filtered)-1)
		deck := filtered[m.deckCursor]
		switch msg.String() {
		case "esc", "enter", "\r", "\n":
			m.editingDeckLimits = false
			return nil, true
		case "tab", "right", "l":
			m.limitCursor = (m.limitCursor + 1) % 2
			return nil, true
		case "left", "h":
			m.limitCursor = (m.limitCursor - 1 + 2) % 2
			return nil, true
		case "+", "=":
			if m.limitCursor == 0 {
				deck.NewCardsPerDay += 5
			} else {
				deck.ReviewLimitPerDay += 20
			}
			return m.setDeckLimits(deck.ID, deck.NewCardsPerDay, deck.ReviewLimitPerDay), true
		case "-", "_":
			if m.limitCursor == 0 {
				deck.NewCardsPerDay = maxInt(0, deck.NewCardsPerDay-5)
			} else {
				deck.ReviewLimitPerDay = maxInt(0, deck.ReviewLimitPerDay-20)
			}
			return m.setDeckLimits(deck.ID, deck.NewCardsPerDay, deck.ReviewLimitPerDay), true
		}
		return nil, true
	}

	filtered := m.filteredDecks()
	if len(filtered) > 0 {
		m.deckCursor = clampInt(m.deckCursor, 0, len(filtered)-1)
	}

	switch msg.String() {
	case "L":
		if len(filtered) > 0 {
			m.editingDeckLimits = true
			m.limitCursor = 0
			return nil, true
		}
		return nil, false
	case "up", "k":
		m.moveDeckCursor(-1)
		return nil, true
	case "down", "j":
		m.moveDeckCursor(1)
		return nil, true
	case "g":
		if len(filtered) > 0 {
			m.deckCursor = 0
		}
		return nil, true
	case "G":
		if len(filtered) > 0 {
			m.deckCursor = len(filtered) - 1
		}
		return nil, true
	case " ", "space", "x", "m":
		if len(filtered) > 0 {
			id := filtered[m.deckCursor].ID
			if m.deckSelected[id] {
				delete(m.deckSelected, id)
			} else {
				m.deckSelected[id] = true
			}
		}
		return nil, true
	case "backspace", "delete":
		return m.handleDeckDelete(), true
	case "M":
		return m.handleDeckMerge(), true
	case "/":
		m.searchingDecks = true
		m.deckFilter = ""
		m.deckCursor = 0
		return nil, true
	case "e", "ctrl+e":
		if len(filtered) > 0 {
			deckID := filtered[m.deckCursor].ID
			return m.exportDeckTSVCmd(deckID), true
		}
		return nil, true
	case "v":
		if len(filtered) > 0 {
			m.selectDeckByID(filtered[m.deckCursor].ID)
			m.activeView = ViewStatistics
			return m.loadStatistics(), true
		}
		return nil, true
	case "c":
		if len(filtered) > 0 {
			m.selectDeckByID(filtered[m.deckCursor].ID)
			m.cramType = "all"
			m.cramCursor = 0
			m.cramReviewed = 0
			m.cramCorrect = 0
			m.activeView = ViewCram
			return m.loadCramCards(), true
		}
		return nil, true
	case "enter", "\r", "\n":
		if len(filtered) > 0 {
			m.selectDeckByID(filtered[m.deckCursor].ID)
			m.activeView = ViewDashboard
		}
		return nil, true
	case "esc":
		if len(m.deckSelected) > 0 {
			m.deckSelected = make(map[string]bool)
			return nil, true
		}
		m.deckFilter = ""
		m.deckCursor = 0
		m.applyDeckFilter()
		return nil, true
	}

	return nil, false
}
