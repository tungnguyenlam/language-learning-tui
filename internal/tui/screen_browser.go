package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// browserScreen wraps the Browser view to satisfy the screen interface.
type browserScreen struct{}

func (browserScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderBrowserAt(layout)
}

func (browserScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.taggingCards {
		switch msg.String() {
		case "enter", "\r", "\n":
			return m.handleTagInput(), true
		case "esc":
			m.taggingCards = false
			m.tagInput = ""
			return nil, true
		case "backspace":
			if len(m.tagInput) > 0 {
				m.tagInput = trimLastRune(m.tagInput)
			}
			return nil, true
		case "ctrl+u":
			m.tagInput = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.tagInput += ch
			return nil, true
		}
		return nil, true
	}

	if m.searchingTags {
		switch msg.String() {
		case "enter", "\r", "\n":
			m.searchingTags = false
			return m.loadBrowserCards(), true
		case "esc":
			m.searchingTags = false
			m.browserTag = ""
			return m.loadBrowserCards(), true
		case "backspace":
			if len(m.browserTag) > 0 {
				m.browserTag = trimLastRune(m.browserTag)
			}
			return m.debounceSearch(ViewBrowser), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserTag += ch
			return m.debounceSearch(ViewBrowser), true
		}

		return nil, true
	}

	if m.searchingBrowser {
		switch msg.String() {
		case "enter", "\r", "\n":
			m.searchingBrowser = false
			if m.browserSearch != "" {
				// Add to history if not already most recent
				if len(m.browserSearchHistory) == 0 || m.browserSearchHistory[len(m.browserSearchHistory)-1] != m.browserSearch {
					m.browserSearchHistory = append(m.browserSearchHistory, m.browserSearch)
					if len(m.browserSearchHistory) > 5 {
						m.browserSearchHistory = m.browserSearchHistory[1:]
					}
				}
			}
			return m.loadBrowserCards(), true
		case "esc":
			m.searchingBrowser = false
			m.browserSearch = ""
			return m.loadBrowserCards(), true
		case "backspace":
			if len(m.browserSearch) > 0 {
				m.browserSearch = trimLastRune(m.browserSearch)
			}
			return m.debounceSearch(ViewBrowser), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserSearch += ch
			return m.debounceSearch(ViewBrowser), true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		m.moveBrowserCursor(-1)
		return nil, true
	case "down", "j":
		m.moveBrowserCursor(1)
		return nil, true
	case "g":
		if len(m.browserCards) > 0 {
			m.browserCursor = 0
		}
		return nil, true
	case "G":
		if len(m.browserCards) > 0 {
			m.browserCursor = len(m.browserCards) - 1
		}
		return nil, true
	case "/":
		m.searchingBrowser = true
		m.browserSearch = ""
		return nil, true
	case "#":
		m.searchingTags = true
		m.browserTag = ""
		return nil, true
	case "m", " ", "space":
		if len(m.browserCards) > 0 {
			cardID := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].ID
			m.browserSelected[cardID] = !m.browserSelected[cardID]
		}
		return nil, true
	case "d":
		return m.lookupBrowserCardInDictionary(), true
	case "b":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserBookmark(true), true
		}
		return m.toggleBrowserBookmark(), true
	case "B":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserBookmark(false), true
		}
		return nil, false
	case "t":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserToggleKind(), true
		}
		return m.toggleCardKind(), true
	case "x":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserSuspend(true), true
		}
		return m.toggleBrowserSuspension(), true
	case "X":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserSuspend(false), true
		}
		return nil, false
	case "T":
		m.taggingCards = true
		m.tagInput = ""
		if len(m.browserCards) > 0 {
			if len(m.getSelectedCardIDs()) == 0 {
				// Pre-fill from current card
				m.tagInput = strings.Join(m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].Tags, " ")
			}
		}
		return nil, true
	case "C":
		return m.cleanupBrowserTags(), true
	case "p":
		if len(m.browserCards) > 0 {
			return m.playCardAudio(m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]), true
		}
		return nil, true
	case "a":
		if len(m.browserCards) > 0 {
			allSelected := true
			for _, card := range m.browserCards {
				if !m.browserSelected[card.ID] {
					allSelected = false
					break
				}
			}
			for _, card := range m.browserCards {
				m.browserSelected[card.ID] = !allSelected
			}
		}
		return nil, true
	case "enter", "\r", "\n":
		if len(m.browserCards) > 0 {
			cardID := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].ID
			if m.showReviewHistory && m.reviewHistoryCard == cardID {
				m.showReviewHistory = false
				return nil, true
			}
			m.reviewHistoryCard = cardID
			return m.loadReviewHistory(cardID), true
		}
		return nil, true
	case "backspace", "delete":
		if len(m.getSelectedCardIDs()) > 0 {
			return m.bulkBrowserDelete(), true
		}
		return m.deleteSelectedCard(), true
	case "esc":
		if len(m.getSelectedCardIDs()) > 0 {
			m.browserSelected = make(map[string]bool)
			return nil, true
		}
		if m.browserSearch != "" || m.browserTag != "" {
			m.browserSearch = ""
			m.browserTag = ""
			return m.loadBrowserCards(), true
		}
	}
	return nil, false
}
