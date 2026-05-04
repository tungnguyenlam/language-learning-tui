package tui

import (
	"strconv"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) updateKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// 1. Global critical keys
	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if !m.textInputActive() {
			if m.activeView == ViewCram && m.cramActive {
				m.cramActive = false
				return m, nil
			}
			return m, tea.Quit
		}
	case "?":
		if !m.textInputActive() {
			m.showHelp = !m.showHelp
			if m.showHelp {
				m.status = "Help overlay shown. Press ? to close."
			} else {
				m.status = "Help overlay closed."
			}
			return m, nil
		}
	}

	// 2. Text input trapping
	if m.textInputActive() {
		// Only allow certain keys to escape trapping
		if key != "enter" && key != "esc" && key != "tab" && key != "shift+tab" {
			// If it's a view-specific key for the active editing view, handle it
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
			return m, nil // Trap everything else
		}
	}

	// 3. Global navigation
	switch key {
	case "tab", "right", "s":
		if !m.textInputActive() {
			return m, m.nextViewCmd()
		}
	case "shift+tab", "left", "w":
		if !m.textInputActive() {
			return m, m.previousViewCmd()
		}
	case "[":
		if !m.textInputActive() || m.activeView == ViewBrowser {
			if m.activeView == ViewBrowser {
				m.previousDeck()
				return m, m.reloadBrowserForSelectedDeck()
			}
			m.previousDeck()
			return m, nil
		}
	case "]":
		if !m.textInputActive() || m.activeView == ViewBrowser {
			if m.activeView == ViewBrowser {
				m.nextDeck()
				return m, m.reloadBrowserForSelectedDeck()
			}
			m.nextDeck()
			return m, nil
		}
	}

	if cmd, handled := m.updateNumberKey(msg); handled {
		return m, cmd
	}

	// 4. View-specific keys
	if cmd, handled := m.updateActiveViewKey(msg); handled {
		return m, cmd
	}

	// 5. Shared navigation keys (if not handled by view)
	switch key {
	case "up", "k":
		if m.activeView == ViewReview {
			if m.cursor > 0 {
				m.cursor--
				m.resetMCQState()
				m.clearReviewHistory()
			}
			return m, nil
		}
	case "down", "j":
		if m.activeView == ViewReview {
			if m.cursor < len(m.dueCards)-1 {
				m.cursor++
				m.resetMCQState()
				m.clearReviewHistory()
			}
			return m, nil
		}
	case "p":
		if m.activeView == ViewReview && len(m.dueCards) > 0 {
			return m, m.playAudio(m.dueCards[m.cursor].Audio)
		}
		if m.activeView == ViewCram && m.cramActive && len(m.cramCards) > 0 {
			return m, m.playAudio(m.cramCards[m.cramCursor].Audio)
		}
	}

	return m, nil
}

func (m *Model) textInputActive() bool {
	return (m.activeView == ViewImport && m.editingImportPath) ||
		(m.activeView == ViewSettings && m.editingTemplate)
}

func (m *Model) updateNumberKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	if m.activeView == ViewCram {
		return nil, false
	}
	if m.textInputActive() {
		return nil, false
	}
	if m.activeView == ViewReview && len(m.dueCards) > 0 && m.dueCards[m.cursor].Kind == core.CardKindMCQ && m.revealState == RevealRevealed && !m.mcqAnswered {
		switch key {
		case "1", "2", "3", "4":
			m.selectMCQChoice(key)
			return nil, true
		}
	}

	switch key {
	case "0", "1":
		m.activeView = ViewDashboard
		return m.updateView(ViewDashboard), true
	case "2":
		m.activeView = ViewDecks
		return m.updateView(ViewDecks), true
	case "3":
		m.activeView = ViewReview
		return m.updateView(ViewReview), true
	case "4":
		m.activeView = ViewStatistics
		return m.updateView(ViewStatistics), true
	case "5":
		m.activeView = ViewImport
		return m.updateView(ViewImport), true
	case "6":
		m.activeView = ViewAI
		return m.updateView(ViewAI), true
	case "7":
		m.activeView = ViewSettings
		return m.updateView(ViewSettings), true
	case "8":
		m.activeView = ViewBrowser
		return m.updateView(ViewBrowser), true
	case "9":
		m.activeView = ViewCram
		return m.updateView(ViewCram), true
	default:
		return nil, false
	}
}

func (m *Model) updateActiveViewKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch m.activeView {
	case ViewReview:
		return m.updateReviewKey(msg)
	case ViewAI:
		return m.updateAIKey(msg)
	case ViewImport:
		return m.updateImportKey(msg)
	case ViewBrowser:
		return m.updateBrowserKey(msg)
	case ViewSettings:
		return m.updateSettingsKey(msg)
	case ViewCram:
		return m.updateCramKey(msg)
	case ViewDecks:
		return m.updateDecksKey(msg)
	case ViewStatistics:
		return m.updateStatisticsKey(msg)
	default:
		return nil, false
	}
}

func (m *Model) updateReviewKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "enter", "space", "\r", "\n", " ":
		if len(m.dueCards) == 0 {
			return nil, false
		}
		card := m.dueCards[m.cursor]
		if m.revealState == RevealRevealed {
			m.revealState = RevealIdle
			m.revealProgress = 0
			m.reviewPredictions = nil
		} else {
			m.revealState = RevealRevealing
			m.revealProgress = 0
		}
		m.mcqChoice = -1
		m.mcqAnswered = false
		return tea.Batch(m.startRevealAnimation(card.Audio), m.loadReviewPredictions(card.ID)), true

	case "b":
		return m.toggleBookmark(), true
	case "B":
		return m.toggleBookmarkFilter(), true
	case "x":
		return m.suspendCard(), true
	case "u":
		return m.undoLastReview(), true
	case "r":
		return m.toggleReviewHistory(), true
	case "a":
		return m.gradeCard(core.GradeAgain), true
	case "h":
		return m.gradeCard(core.GradeHard), true
	case "g":
		return m.gradeCard(core.GradeGood), true
	case "e":
		return m.gradeCard(core.GradeEasy), true
	}
	return nil, false
}

func (m *Model) updateStatisticsKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "up", "k":
		m.scrollStats(-1)
		return nil, true
	case "down", "j":
		m.scrollStats(1)
		return nil, true
	}
	return nil, false
}

func (m *Model) updateDecksKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	filtered := m.filteredDecks()
	switch msg.String() {
	case "up", "k":
		if m.deckCursor > 0 {
			m.deckCursor--
		}
		return nil, true
	case "down", "j":
		if m.deckCursor < len(filtered)-1 {
			m.deckCursor++
		}
		return nil, true
	case "enter", "\r", "\n":
		if len(filtered) > 0 {
			m.selectDeckByID(filtered[m.deckCursor].ID)
			m.activeView = ViewDashboard
		}
		return nil, true
	case "backspace":
		if len(m.deckFilter) > 0 {
			m.deckFilter = m.deckFilter[:len(m.deckFilter)-1]
			m.deckCursor = 0
		}
		return nil, true
	case "esc":
		m.deckFilter = ""
		m.deckCursor = 0
		return nil, true
	}

	if len(msg.String()) == 1 {
		m.deckFilter += msg.String()
		m.deckCursor = 0
		return nil, true
	}

	return nil, false
}

func (m *Model) updateImportKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.editingImportPath {
		switch msg.String() {
		case "enter", "esc":
			m.editingImportPath = false
			return nil, true
		case "backspace":
			if m.importCursor == 0 {
				if len(m.importPath) > 0 {
					m.importPath = m.importPath[:len(m.importPath)-1]
				}
			} else {
				if len(m.exportPath) > 0 {
					m.exportPath = m.exportPath[:len(m.exportPath)-1]
				}
			}
			return nil, true
		case "ctrl+u":
			if m.importCursor == 0 {
				m.importPath = ""
			} else {
				m.exportPath = ""
			}
			return nil, true
		}
		if len(msg.String()) == 1 {
			if m.importCursor == 0 {
				m.importPath += msg.String()
			} else {
				m.exportPath += msg.String()
			}
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.importCursor > 0 {
			m.importCursor--
		}
		return nil, true
	case "down", "j":
		if m.importCursor < 1 {
			m.importCursor++
		}
		return nil, true
	case "enter", "\r", "\n":
		m.editingImportPath = true
		return nil, true
	case "i":
		return m.importTSV(), true
	case "I":
		return m.importAPKG(), true
	case "x":
		return m.exportTSV(), true
	case "X":
		return m.exportAPKG(), true
	}
	return nil, false
}

func (m *Model) updateSettingsKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.editingTemplate {
		switch msg.String() {
		case "enter", "esc":
			m.editingTemplate = false
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{Templates: m.aiTemplates}
			}
			if m.onConfigChange != nil {
				m.onConfigChange(m.aiProviderName, m.aiTemplates, m.autoPlayAudio)
			}
			return nil, true
		case "backspace":
			templateKey := m.templateKeyAtCursor()
			if len(m.aiTemplates[templateKey]) > 0 {
				m.aiTemplates[templateKey] = m.aiTemplates[templateKey][:len(m.aiTemplates[templateKey])-1]
			}
			return nil, true
		}
		if len(msg.String()) == 1 {
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[templateKey] += msg.String()
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return nil, true
	case "down", "j":
		if m.settingsCursor < 5 {
			m.settingsCursor++
		}
		return nil, true
	case "enter", "\r", "\n":
		return m.handleSettingsEnter(), true
	case "+":
		return m.setDailyGoal(m.stats.DailyGoal + 1), true
	case "-":
		return m.setDailyGoal(m.stats.DailyGoal - 1), true
	}
	return nil, false
}

func (m *Model) updateBrowserKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.searchingBrowser {
		switch msg.String() {
		case "enter", "esc":
			m.searchingBrowser = false
			return m.loadBrowserCards(), true
		case "backspace":
			if len(m.browserSearch) > 0 {
				m.browserSearch = m.browserSearch[:len(m.browserSearch)-1]
			}
			return m.loadBrowserCards(), true
		}
		if len(msg.String()) == 1 {
			m.browserSearch += msg.String()
			return m.loadBrowserCards(), true
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
	case "/":
		m.searchingBrowser = true
		m.browserSearch = ""
		return nil, true
	case "m":
		if len(m.browserCards) > 0 {
			cardID := m.browserCards[m.browserCursor].ID
			m.browserSelected[cardID] = !m.browserSelected[cardID]
		}
		return nil, true
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
	case "enter", "\r", "\n":
		if len(m.browserCards) > 0 {
			cardID := m.browserCards[m.browserCursor].ID
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
	}
	return nil, false
}

func (m *Model) updateCramKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.cramActive {
		switch msg.String() {
		case "enter", "space", "\r", "\n", " ":
			if m.cramRevealed {
				m.cramRevealed = false
				m.revealState = RevealIdle
				m.revealProgress = 0
			} else {
				m.cramRevealed = true
				m.revealState = RevealRevealed
				m.revealProgress = 100
			}
			return nil, true

		case "p":
			if len(m.cramCards) > 0 {
				return m.playAudio(m.cramCards[m.cramCursor].Audio), true
			}
		case "a":
			return m.gradeCramCard(core.GradeAgain), true
		case "h":
			return m.gradeCramCard(core.GradeHard), true
		case "g":
			return m.gradeCramCard(core.GradeGood), true
		case "e":
			return m.gradeCramCard(core.GradeEasy), true
		case "q", "esc":
			m.cramActive = false
			return nil, true
		}
		return nil, false
	}

	switch msg.String() {
	case "up", "k":
		m.moveCramCursor(-1)
		return nil, true
	case "down", "j":
		m.moveCramCursor(1)
		return nil, true
	case "1", "2", "3", "4", "5":
		idx, _ := strconv.Atoi(msg.String())
		return m.setCramFilter(idx), true
	case "enter", "\r", "\n":
		if len(m.cramCards) > 0 {
			m.cramActive = true
			m.cramRevealed = false
			m.revealState = RevealIdle
			m.revealProgress = 0
		}
		return nil, true
	}
	return nil, false
}

func (m *Model) updateAIKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "esc":
		m.aiInput = ""
		return nil, true
	case "up", "k":
		if m.draftCursor > 0 {
			m.draftCursor--
		}
		return nil, true
	case "down", "j":
		if m.draftCursor < len(m.drafts)-1 {
			m.draftCursor++
		}
		return nil, true
	case "A":
		return m.approveAllDrafts(), true
	case "a":
		if len(m.drafts) > 0 {
			return m.approveDraft(), true
		}
	case "D":
		if len(m.drafts) > 0 {
			m.drafts = nil
			m.draftCursor = 0
			m.status = "All drafts discarded"
			return nil, true
		}
	case "d":
		if len(m.drafts) > 0 {
			m.discardDraft()
			return nil, true
		}
	case "enter", "\r", "\n":
		if len(m.drafts) > 0 {
			return m.approveDraft(), true
		}
		return m.startDrafting(), true
	case "backspace":
		if len(m.drafts) > 0 {
			m.discardDraft()
		} else if len(m.aiInput) > 0 {
			m.aiInput = m.aiInput[:len(m.aiInput)-1]
		}
		return nil, true
	}
	if len(msg.String()) == 1 && msg.String() >= " " && msg.String() <= "~" {
		m.aiInput += msg.String()
		return nil, true
	}
	return nil, false
}
