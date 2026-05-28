package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) updateKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// 0. High-priority learning mode trapping
	if m.typingMode || (m.activeView == ViewCram && m.cramActive) {
		if key != "tab" && key != "shift+tab" && key != "ctrl+c" && key != "q" {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
	}

	// 0.05 Fix-proposal preview trapping: when the AI has returned a
	// correction, only y/n/esc are meaningful until the user decides.
	if m.fixProposal != nil {
		switch key {
		case "y", "Y", "enter":
			return m, m.applyFixProposal()
		case "n", "N", "esc":
			m.discardFixProposal()
			return m, nil
		}
		return m, nil
	}

	// 0.1 Deletion confirmation trapping
	if m.confirmingDelete {
		switch key {
		case "y", "Y", "enter":
			m.confirmingDelete = false
			if m.deleteAction != nil {
				return m, m.deleteAction()
			}
		case "n", "N", "esc":
			m.confirmingDelete = false
			m.deleteAction = nil
			m.status = "Deletion cancelled"
			return m, nil
		}
		return m, nil // Trap everything else during confirmation
	}

	// 1. Global critical keys
	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+d":
		if m.activeView == ViewDebug {
			m.activeView = ViewDashboard
		} else {
			m.activeView = ViewDebug
		}
		return m, nil
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
			msg := "Help overlay closed."
			if m.showHelp {
				msg = "Help overlay shown. Press ? to close."
			}
			m.status = msg
			return m, nil
		}
	}

	// 2. Text input trapping
	if m.textInputActive() {
		m.logger.Debug("Text input active, trapping key: %s", key)
		// Only allow certain keys to escape trapping
		if key != "tab" && key != "shift+tab" {
			// If it's a view-specific key for the active editing view, handle it
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				m.logger.Debug("Key %s handled by active view despite text input", key)
				return m, cmd
			}
			m.logger.Debug("Key %s trapped by text input mode", key)
			return m, nil // Trap everything else
		}
	}

	// 3. Global navigation
	switch key {
	case "tab", "right", "s":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() {
			return m, m.nextViewCmd()
		}
	case "shift+tab", "left", "w":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() {
			return m, m.previousViewCmd()
		}
	case "[":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() || m.activeView == ViewBrowser {
			if m.activeView == ViewBrowser {
				m.previousDeck()
				return m, m.reloadBrowserForSelectedDeck()
			}
			m.previousDeck()
			return m, nil
		}
	case "]":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() || m.activeView == ViewBrowser {
			if m.activeView == ViewBrowser {
				m.nextDeck()
				return m, m.reloadBrowserForSelectedDeck()
			}
			m.nextDeck()
			return m, nil
		}
	}

	// 3. View-specific keys
	if cmd, handled := m.updateActiveViewKey(msg); handled {
		return m, cmd
	}

	if cmd, handled := m.updateNumberKey(msg); handled {
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
				m.showHint = false
			}
			return m, nil
		}
	case "down", "j":
		if m.activeView == ViewReview {
			if m.cursor < len(m.dueCards)-1 {
				m.cursor++
				m.resetMCQState()
				m.clearReviewHistory()
				m.showHint = false
			}
			return m, nil
		}
	case "p":
		if m.activeView == ViewReview && len(m.dueCards) > 0 {
			return m, m.playCardAudio(m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)])
		}
		if m.activeView == ViewCram && m.cramActive && len(m.cramCards) > 0 {
			return m, m.playCardAudio(m.cramCards[clampInt(m.cramCursor, 0, len(m.cramCards)-1)])
		}
	}

	return m, nil
}

func (m *Model) textInputActive() bool {
	return (m.activeView == ViewImport && m.editingImportPath) ||
		(m.activeView == ViewSettings && m.editingTemplate) ||
		(m.activeView == ViewSettings && m.editingSecretKey != "") ||
		(m.activeView == ViewImport && m.editingExportTag) ||
		m.taggingCards ||
		m.searchingBrowser ||
		m.searchingTags ||
		m.searchingDecks ||
		m.searchingAI ||
		m.drafting // AI drafting also uses input
}

func (m *Model) updateNumberKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	if m.activeView == ViewCram {
		return nil, false
	}
	if m.textInputActive() {
		return nil, false
	}
	switch key {
	case "0":
		m.activeView = ViewPractice
		return m.updateView(ViewPractice), true
	case "1":
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
	case ViewDashboard:
		return m.updateDashboardKey(msg)
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
	case ViewPractice:
		return m.updatePracticeKey(msg)
	case ViewSessionSummary:
		return m.updateSessionSummaryKey(msg)
	case ViewDecks:
		return m.updateDecksKey(msg)
	case ViewStatistics:
		return m.updateStatisticsKey(msg)
	default:
		return nil, false
	}
}

func (m *Model) updateDashboardKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
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
	}
	return nil, false
}

func (m *Model) updateReviewKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	// Handle typing mode input
	if m.typingMode {
		switch key {
		case "enter", "space", "\r", "\n", " ":
			if !m.typingChecked {
				m.typingChecked = true
				if len(m.dueCards) > 0 {
					card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
					targetAnswer := card.Answer
					if card.Kind == core.CardKindCloze && len(card.Choices) > 0 {
						targetAnswer = card.Choices[0]
					}
					m.typingCorrect = m.normalizeAnswer(m.typedAnswer) == m.normalizeAnswer(targetAnswer)
					m.revealState = RevealRevealed
					m.revealProgress = 100
					return m.loadReviewPredictions(card.ID), true
				}
				m.revealState = RevealRevealed
				m.revealProgress = 100
				return nil, true
			}
		case "esc":
			m.typingMode = false
			m.typedAnswer = ""
			m.typingChecked = false
			m.status = "Typing mode off"
			return nil, true
		case "backspace":
			if !m.typingChecked {
				if len(m.typedAnswer) > 0 {
					m.typedAnswer = trimLastRune(m.typedAnswer)
				}
				return nil, true
			}
		}
		if !m.typingChecked && utf8.RuneCountInString(key) == 1 {
			r, _ := utf8.DecodeRuneInString(key)
			if unicode.IsPrint(r) {
				m.typedAnswer += key
				return nil, true
			}
		}
		// If typingChecked, allow other keys (like grades) to fall through
		if !m.typingChecked {
			return nil, true
		}
	}

	if len(m.dueCards) > 0 {
		card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
		if card.Kind == core.CardKindMCQ && !m.mcqAnswered {
			switch key {
			case "1", "2", "3", "4":
				m.selectMCQChoice(key)
				m.revealState = RevealRevealed
				m.revealProgress = 100
				return m.loadReviewPredictions(card.ID), true
			}
		}

		// Allow 1-4 for grading when revealed
		if m.revealState == RevealRevealed {
			switch key {
			case "1":
				return m.gradeCard(core.GradeAgain), true
			case "2":
				return m.gradeCard(core.GradeHard), true
			case "3":
				return m.gradeCard(core.GradeGood), true
			case "4":
				return m.gradeCard(core.GradeEasy), true
			}
		}
	}

	switch key {
	case "enter", "space", "\r", "\n", " ":
		if len(m.dueCards) == 0 {
			m.status = "No cards due"
			return nil, true
		}
		card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
		switch m.revealState {
		case RevealIdle:
			return tea.Batch(m.startRevealAnimation(card), m.loadReviewPredictions(card.ID)), true
		case RevealRevealing:
			m.revealProgress += 15
			if m.revealProgress >= 100 {
				m.revealProgress = 100
				m.revealState = RevealRevealed
			}
			return m.loadReviewPredictions(card.ID), true
		case RevealRevealed:
			m.status = "Grade with a=again h=hard g=good e=easy"
			return nil, true
		}
	case "c":
		if len(m.dueCards) == 0 {
			m.activeView = ViewCram
			m.status = "Custom Study (Cram Mode)"
			return m.loadCramCards(), true
		}
	case "a":
		if m.revealState == RevealRevealed {
			return m.gradeCard(core.GradeAgain), true
		}
	case "h":
		if m.revealState == RevealRevealed {
			return m.gradeCard(core.GradeHard), true
		}
		if len(m.dueCards) > 0 {
			m.showHint = !m.showHint
			if m.showHint {
				m.status = "Hint shown"
			} else {
				m.status = "Hint hidden"
			}
			return nil, true
		}
	case "g":
		if m.revealState == RevealRevealed {
			return m.gradeCard(core.GradeGood), true
		}
	case "e":
		if m.revealState == RevealRevealed {
			return m.gradeCard(core.GradeEasy), true
		}
	case "u":
		return m.undoLastReview(), true
	case "r":
		return m.toggleReviewHistory(), true
	case "i":
		if len(m.dueCards) > 0 {
			m.showCardInfo = !m.showCardInfo
			if m.showCardInfo {
				m.status = "Card info shown"
			} else {
				m.status = "Card info hidden"
			}
			return nil, true
		}
	case "f":
		m.focusMode = !m.focusMode
		if m.focusMode {
			m.status = "Focus mode enabled"
		} else {
			m.status = "Focus mode disabled"
		}
		return nil, true
	case "H":
		return m.explainCard(), true
	case "p":
		if len(m.dueCards) > 0 {
			return m.playCardAudio(m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]), true
		}
	case "t":
		if !m.typingMode && m.revealState == RevealIdle {
			m.typingMode = true
			m.typedAnswer = ""
			m.typingChecked = false
			m.status = "Type your answer, press Enter to check"
			return nil, true
		}
	case "b":
		return m.toggleBookmark(), true
	case "d":
		if len(m.dueCards) > 0 {
			card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
			word := strings.Split(card.Prompt, "\n")[0] // Use first line in case of multiline
			return m.openDictionary(word), true
		}
	case "B":
		return m.toggleBookmarkFilter(), true
	case "x":
		return m.suspendCard(), true
	case "!", "F":
		return m.reportCardWrong(), true
	case "delete", "backspace":
		if len(m.dueCards) > 0 {
			card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
			return func() tea.Msg {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				if err := m.repo.DeleteCard(ctx, card.ID); err != nil {
					return err
				}
				return m.loadDueCards()
			}, true
		}
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
	case "x":
		return m.exportStatsCSV(), true
	}
	return nil, false
}

func (m *Model) updateDecksKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.searchingDecks {
		switch msg.String() {
		case "enter", "\r", "\n":
			m.searchingDecks = false
			m.applyDeckFilter()
			filtered := m.filteredDecks()
			if len(filtered) > 0 {
				m.selectDeckByID(filtered[0].ID)
			}
			return nil, true
		case "esc", "\x1b":
			m.searchingDecks = false
			m.applyDeckFilter()
			return nil, true
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
	if m.deckCursor >= len(filtered) && len(filtered) > 0 {
		m.deckCursor = len(filtered) - 1
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
	case " ":
		if len(filtered) > 0 {
			id := filtered[m.deckCursor].ID
			m.deckSelected[id] = !m.deckSelected[id]
		}
		return nil, true
	case "x":
		if len(filtered) > 0 {
			id := filtered[m.deckCursor].ID
			m.deckSelected[id] = !m.deckSelected[id]
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
		return nil, true
	}

	return nil, false
}

func (m *Model) updateImportKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.editingImportPath {
		switch msg.String() {
		case "enter", "\r", "\n", "esc":
			m.editingImportPath = false
			return nil, true
		case "backspace":
			if m.importCursor == 0 {
				if len(m.importPath) > 0 {
					m.importPath = trimLastRune(m.importPath)
				}
			} else {
				if len(m.exportPath) > 0 {
					m.exportPath = trimLastRune(m.exportPath)
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
		if ch, ok := singlePrintableInput(msg.String()); ok {
			if m.importCursor == 0 {
				m.importPath += ch
			} else {
				m.exportPath += ch
			}
			return nil, true
		}
		return nil, true
	}

	if m.editingExportTag {
		switch msg.String() {
		case "enter", "\r", "\n", "esc":
			m.editingExportTag = false
			return nil, true
		case "backspace":
			if len(m.exportTag) > 0 {
				m.exportTag = trimLastRune(m.exportTag)
			}
			return nil, true
		case "ctrl+u":
			m.exportTag = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.exportTag += ch
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
		if m.importCursor < 4 {
			m.importCursor++
		}
		return nil, true
	case "[":
		if m.importCursor == 2 {
			m.previousExportDeck()
			return nil, true
		} else if m.importCursor == 4 {
			m.cycleExportFilter(false)
			return nil, true
		}
	case "]":
		if m.importCursor == 2 {
			m.nextExportDeck()
			return nil, true
		} else if m.importCursor == 4 {
			m.cycleExportFilter(true)
			return nil, true
		}
	case "enter", "\r", "\n":
		if m.importCursor < 2 {
			m.editingImportPath = true
		} else if m.importCursor == 3 {
			m.editingExportTag = true
		}
		return nil, true
	case "t":
		m.importCursor = 3
		m.editingExportTag = true
		return nil, true
	case "i":
		return m.importTSV(), true
	case "I":
		return m.importAPKG(), true
	case "S":
		return m.seedStandardContent(), true
	case "R":
		return m.handleResetDatabase(), true
	case "x":
		return m.exportTSV(), true
	case "X":
		return m.exportAPKG(), true
	}
	return nil, false
}

func (m *Model) updateSettingsKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.editingSecretKey != "" {
		return m.handleSecretEditKey(msg)
	}
	if m.editingTemplate {
		activeSet := m.currentAITemplateSet()
		if activeSet == "" {
			m.editingTemplate = false
			m.originalTemplateValue = ""
			return nil, true
		}
		switch msg.String() {
		case "enter", "\r", "\n":
			m.editingTemplate = false
			m.originalTemplateValue = ""
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{
					Templates: m.aiTemplates,
					ActiveSet: activeSet,
				}
			}
			if m.onConfigChange != nil {
				m.onConfigChange(m.theme, m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization)
			}
			return nil, true
		case "esc":
			// Restore original value on cancel
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] = m.originalTemplateValue
			m.editingTemplate = false
			m.originalTemplateValue = ""
			if m.aiProviderName == "template" {
				m.aiProvider = ai.TemplateProvider{
					Templates: m.aiTemplates,
					ActiveSet: activeSet,
				}
			}
			if m.onConfigChange != nil {
				m.onConfigChange(m.theme, m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization)
			}
			return nil, true
		case "backspace":
			templateKey := m.templateKeyAtCursor()
			val := m.aiTemplates[activeSet][templateKey]
			if len(val) > 0 {
				m.aiTemplates[activeSet][templateKey] = trimLastRune(val)
			}
			return nil, true
		case "ctrl+u":
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			templateKey := m.templateKeyAtCursor()
			m.aiTemplates[activeSet][templateKey] += ch
			return nil, true
		}
		return nil, true
	}

	// Secret-edit fall-through: handled at the top of updateSettingsKey.
	switch msg.String() {
	case "up", "k":
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return nil, true
	case "down", "j":
		if m.settingsCursor < 10 {
			m.settingsCursor++
		}
		return nil, true
	case "c":
		return m.cycleTheme(), true
	case "enter":
		return m.handleSettingsEnter(), true
	case "+":
		return m.setDailyGoal(m.stats.DailyGoal + 1), true
	case "-":
		return m.setDailyGoal(m.stats.DailyGoal - 1), true
	case "[":
		m.previousAITemplate()
		return nil, true
	case "]":
		m.nextAITemplate()
		return nil, true
	}
	return nil, false
}

func (m *Model) updateBrowserKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
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
			return m.loadBrowserCards(), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserTag += ch
			return m.loadBrowserCards(), true
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
			return m.loadBrowserCards(), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserSearch += ch
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
	case "#":
		m.searchingTags = true
		m.browserTag = ""
		return nil, true
	case "m":
		if len(m.browserCards) > 0 {
			cardID := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)].ID
			m.browserSelected[cardID] = !m.browserSelected[cardID]
		}
		return nil, true
	case "d":
		if len(m.browserCards) > 0 {
			card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
			word := strings.Split(card.Prompt, "\n")[0]
			return m.openDictionary(word), true
		}
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
				return m.playCardAudio(m.cramCards[clampInt(m.cramCursor, 0, len(m.cramCards)-1)]), true
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
	if m.searchingAI {
		switch msg.String() {
		case "enter", "\r", "\n":
			m.searchingAI = false
			return m.startDrafting(), true
		case "esc":
			m.searchingAI = false
			return nil, true
		case "backspace":
			if len(m.aiInput) > 0 {
				m.aiInput = trimLastRune(m.aiInput)
			}
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.aiInput += ch
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "/":
		m.searchingAI = true
		return nil, true
	case "esc":
		m.aiInput = ""
		return nil, true
	case "[":
		m.previousAITemplate()
		return nil, true
	case "]":
		m.nextAITemplate()
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
			return nil, true
		}
	}
	return nil, false
}

func (m *Model) updateSessionSummaryKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	// Any key returns to dashboard and resets session stats
	m.lastSessionReviewed = m.sessionReviewed
	m.lastSessionCorrect = m.sessionCorrect
	m.sessionReviewed = 0
	m.sessionCorrect = 0
	return m.updateView(ViewDashboard), true
}

func (m *Model) handlePaste(text string) tea.Cmd {
	// Clean text: strip newlines/tabs which break single-line inputs
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\t", " ")

	if m.activeView == ViewSettings {
		if m.editingSecretKey != "" {
			provider := m.editingSecretProvider
			key := m.editingSecretKey
			m.setCredValue(provider, key, m.getCredValue(provider, key)+text)
			return nil
		}
		if m.editingTemplate {
			activeSet := m.currentAITemplateSet()
			if activeSet != "" {
				templateKey := m.templateKeyAtCursor()
				m.aiTemplates[activeSet][templateKey] += text
				return nil
			}
		}
	}

	if m.activeView == ViewAI && m.searchingAI {
		m.aiInput += text
		return nil
	}

	if m.activeView == ViewBrowser {
		if m.searchingBrowser {
			m.browserSearch += text
			return m.loadBrowserCards()
		}
		if m.searchingTags {
			m.browserTag += text
			return m.loadBrowserCards()
		}
		if m.taggingCards {
			m.tagInput += text
			return nil
		}
	}

	if m.activeView == ViewImport {
		if m.editingImportPath {
			if m.importCursor == 0 {
				m.importPath += text
			} else {
				m.exportPath += text
			}
			return nil
		}
		if m.editingExportTag {
			m.exportTag += text
			return nil
		}
	}

	if m.activeView == ViewDecks && m.searchingDecks {
		m.deckFilter += text
		m.deckCursor = 0
		return nil
	}

	return nil
}

// handleSecretEditKey processes keys while the user is typing an API key,
// model name, or base URL into Settings. Like template editing, Enter
// commits the value and triggers a save; Esc reverts to the value we
// stashed in originalSecretValue. Backspace and ctrl+u edit the buffer.
func (m *Model) handleSecretEditKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	provider := m.editingSecretProvider
	key := m.editingSecretKey
	if provider == "" || key == "" {
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		return nil, true
	}

	commit := func() {
		// Rebuild the provider so the new credentials take effect immediately.
		m.aiProvider = buildProvider(m.aiProviderName, m.aiSecrets, m.aiTemplates, m.currentAITemplateSet())
		if m.onSecretsChange != nil {
			m.onSecretsChange(m.aiSecrets)
		}
	}

	switch msg.String() {
	case "enter", "\r", "\n":
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		m.originalSecretValue = ""
		commit()
		m.status = fmt.Sprintf("Saved %s %s", provider, key)
		return nil, true
	case "esc":
		m.setCredValue(provider, key, m.originalSecretValue)
		m.editingSecretKey = ""
		m.editingSecretProvider = ""
		m.originalSecretValue = ""
		commit()
		m.status = "Edit cancelled"
		return nil, true
	case "backspace":
		val := m.getCredValue(provider, key)
		if len(val) > 0 {
			m.setCredValue(provider, key, trimLastRune(val))
		}
		return nil, true
	case "ctrl+u":
		m.setCredValue(provider, key, "")
		return nil, true
	}
	if ch, ok := singlePrintableInput(msg.String()); ok {
		m.setCredValue(provider, key, m.getCredValue(provider, key)+ch)
		return nil, true
	}
	return nil, true
}

func (m *Model) updatePracticeKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if len(m.practiceItems) == 0 {
		return nil, false
	}

	if m.practiceRevealed {
		// Any key to next noun
		m.practiceRevealed = false
		m.practiceIndex = (m.practiceIndex + 1) % len(m.practiceItems)
		// Reshuffle or repeat? For now just cycle.
		return nil, true
	}

	key := msg.String()
	item := m.practiceItems[m.practiceIndex]
	var choice string

	switch key {
	case "1", "d", "m":
		choice = "der"
	case "2", "i", "f":
		choice = "die"
	case "3", "a", "n":
		choice = "das"
	case "q", "esc":
		return m.updateView(ViewDashboard), true
	default:
		return nil, false
	}

	m.practiceTotal++
	m.practiceRevealed = true
	if choice == item.Article {
		m.practiceCorrect++
		m.practiceLastResult = true
	} else {
		m.practiceLastResult = false
	}

	return nil, true
}
