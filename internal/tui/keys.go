package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) updateKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Intercept keys when the spotlight dictionary overlay is active
	if m.dictionaryOverlayActive {
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if cmd, handled := m.updateDictionaryOverlayKey(msg); handled {
			return m, cmd
		}
	}

	// Toggle dictionary overlay on '='
	if key == "=" && !m.textInputActive() {
		return m, m.openDictionaryOverlay()
	}

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
		if m.activeView == ViewDictionary {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
		if m.activeView == ViewDebug {
			m.activeView = ViewDashboard
		} else {
			m.activeView = ViewDebug
		}
		return m, nil
	case "q":
		// "q" is a legitimate character in German answers (Qualität, Quelle),
		// so let an actively-typing trainer consume it before quitting.
		if m.trainerInputActive() {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
		if !m.textInputActive() {
			if m.activeView == ViewCram && m.cramActive {
				m.cramActive = false
				return m, nil
			}
			return m, tea.Quit
		}
	case "?":
		if m.trainerInputActive() {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
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
		allowNumbers := (m.activeView == ViewDictionary || m.dictionaryOverlayActive) && m.dictionarySearch == ""
		isNumber := len(key) == 1 && unicode.IsDigit([]rune(key)[0])

		if key != "tab" && key != "shift+tab" && key != "left" && key != "right" && key != "up" && key != "down" && key != "j" && key != "k" && !(allowNumbers && isNumber) {
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
	case "K":
		if !m.textInputActive() {
			return m, tea.Sequence(m.updateView(ViewPractice), func() tea.Msg {
				return m.enterPracticeMode(PracticeSubViewConjugation)()
			})
		}
	case "tab", "right", "s":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() || (m.activeView == ViewDictionary && (key == "tab" || key == "right")) {
			return m, m.nextViewCmd()
		}
	case "shift+tab", "left", "w":
		if cmd, handled := m.updateActiveViewKey(msg); handled {
			return m, cmd
		}
		if !m.textInputActive() || (m.activeView == ViewDictionary && (key == "shift+tab" || key == "left")) {
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
			if m.activeView == ViewReview {
				m.cursor = 0
				m.resetReviewState()
				m.status = fmt.Sprintf("Deck: %s", m.deckLabel())
			}
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
			if m.activeView == ViewReview {
				m.cursor = 0
				m.resetReviewState()
				m.status = fmt.Sprintf("Deck: %s", m.deckLabel())
			}
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
				m.cardTransitioning = true
				m.cardTransitionProgress = 0
				m.cardTransitionFrame = 0
				m.cardTransitionDir = -1
				m.cursor--
				m.resetReviewState()
				return m, m.tickCardTransition()
			}
			return m, nil
		}
	case "down", "j":
		if m.activeView == ViewReview {
			if m.cursor < len(m.dueCards)-1 {
				m.cardTransitioning = true
				m.cardTransitionProgress = 0
				m.cardTransitionFrame = 0
				m.cardTransitionDir = 1
				m.cursor++
				m.resetReviewState()
				return m, m.tickCardTransition()
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
	return (m.activeView == ViewImport && m.importScreen.editingImportPath) ||
		(m.activeView == ViewSettings && m.editingTemplate) ||
		(m.activeView == ViewSettings && m.editingSecretKey != "") ||
		(m.activeView == ViewImport && m.importScreen.editingExportTag) ||
		(m.activeView == ViewAnkiWeb && m.ankiWebScreen.editingQuery) ||
		m.taggingCards ||
		m.searchingBrowser ||
		m.searchingTags ||
		m.searchingDecks ||
		m.searchingAI ||
		(m.activeView == ViewDictionary && !m.dictionaryFocusResults) ||
		m.dictionaryOverlayActive ||
		m.drafting // AI drafting also uses input
}

// trainerInputActive reports whether a generic text-input trainer is waiting
// for the learner to type an answer.
//
// Trainers deliberately stay out of textInputActive(): that predicate also
// disables Tab/arrow view switching, which trainers must keep. This narrower
// check exists so single-letter global shortcuts ("q", "?") are typed into the
// answer instead of quitting the app or opening the help overlay mid-word.
func (m *Model) trainerInputActive() bool {
	if m.activeView != ViewPractice || !isGenericTrainer(m.practiceSubView) {
		return false
	}
	st, ok := m.trainers[m.practiceSubView]
	return ok && len(st.items) > 0 && !st.revealed
}

func (m *Model) updateNumberKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	if m.textInputActive() {
		// Exception: allow view switching from Dictionary if search is empty
		if (m.activeView == ViewDictionary || m.dictionaryOverlayActive) && m.dictionarySearch == "" {
			// fall through
		} else {
			return nil, false
		}
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
	if s, ok := m.screens[m.activeView]; ok {
		return s.HandleKey(m, msg)
	}
	return nil, false
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
		case RevealFlipping:
			m.flipProgress = 100
			m.revealState = RevealRevealing
			m.revealProgress = 0
			return m.loadReviewPredictions(card.ID), true
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
	case "u", "z", "ctrl+z":
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
	case "G":
		if len(m.dueCards) > 0 {
			m.showGrammarHint = !m.showGrammarHint
			if m.showGrammarHint {
				card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
				tip := content.GetRelevantGrammarTip(card.Prompt)
				m.grammarHint = &tip
				m.status = "Grammar hint shown"
			} else {
				m.grammarHint = nil
				m.status = "Grammar hint hidden"
			}
			return nil, true
		}
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
		return m.lookupReviewCardInDictionary(), true
	case "B":
		return m.toggleBookmarkFilter(), true
	case "x":
		return m.suspendCard(), true
	case "!", "F":
		return m.reportCardWrong(), true
	case "delete", "backspace":
		if len(m.dueCards) > 0 {
			return m.deleteReviewCard(), true
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
	case "pgup":
		m.scrollStats(-10)
		return nil, true
	case "pgdown":
		m.scrollStats(10)
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
			m.recordDeckSearch(m.deckFilter)
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
		case "ctrl+x":
			m.deckSearchHistory = nil
			m.saveDeckHistory()
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
	case " ":
		if len(filtered) > 0 {
			id := filtered[m.deckCursor].ID
			m.deckSelected[id] = !m.deckSelected[id]
		}
		return nil, true
	case "x", "m":
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
		m.applyDeckFilter()
		return nil, true
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
				m.onConfigChange(m.theme, m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization, m.revealSpeed)
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
				m.onConfigChange(m.theme, m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization, m.revealSpeed)
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
		if m.settingsCursor < 14 {
			m.settingsCursor++
		}
		return nil, true
	case "g":
		m.settingsCursor = 0
		return nil, true
	case "G":
		m.settingsCursor = 14
		return nil, true
	case "c":
		return m.cycleTheme(), true
	case "enter":
		return m.handleSettingsEnter(), true
	case "+":
		if m.settingsCursor == 14 {
			return m.setRevealSpeed(m.revealSpeed + 1), true
		}
		return m.setDailyGoal(m.stats.DailyGoal + 1), true
	case "-":
		if m.settingsCursor == 14 {
			return m.setRevealSpeed(m.revealSpeed - 1), true
		}
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
			m.browserSearchTimerID++
			id := m.browserSearchTimerID
			return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
				return debounceSearchMsg{id: id, view: ViewBrowser}
			}), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserTag += ch
			m.browserSearchTimerID++
			id := m.browserSearchTimerID
			return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
				return debounceSearchMsg{id: id, view: ViewBrowser}
			}), true
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
			m.browserSearchTimerID++
			id := m.browserSearchTimerID
			return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
				return debounceSearchMsg{id: id, view: ViewBrowser}
			}), true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.browserSearch += ch
			m.browserSearchTimerID++
			id := m.browserSearchTimerID
			return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
				return debounceSearchMsg{id: id, view: ViewBrowser}
			}), true
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
	case "m":
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

func (m *Model) updateCramKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.cramActive {
		switch msg.String() {
		case "enter", "space", "\r", "\n", " ":
			if m.cramRevealed {
				m.cramRevealed = false
				m.revealState = RevealIdle
				m.revealProgress = 0
				m.flipProgress = 0
				m.flipFrame = 0
			} else {
				switch m.revealState {
				case RevealIdle:
					if len(m.cramCards) > 0 {
						card := m.cramCards[clampInt(m.cramCursor, 0, len(m.cramCards)-1)]
						return m.startRevealAnimation(card), true
					}
				case RevealFlipping:
					m.flipProgress = 100
					m.revealState = RevealRevealing
					m.revealProgress = 0
				case RevealRevealing:
					m.revealProgress += 15
					if m.revealProgress >= 100 {
						m.revealProgress = 100
						m.revealState = RevealRevealed
						m.cramRevealed = true
					}
				case RevealRevealed:
					m.cramRevealed = true
				}
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
			m.cramRevealed = false
			m.revealState = RevealIdle
			m.revealProgress = 0
			m.flipProgress = 0
			m.flipFrame = 0
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
		case "esc", "escape":
			m.searchingAI = false
			m.aiInput = ""
			m.draftSource = ""
			return nil, true
		case "backspace":
			if len(m.aiInput) > 0 {
				m.aiInput = trimLastRune(m.aiInput)
				if len(m.aiInput) == 0 {
					m.draftSource = ""
				}
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
	case "esc", "escape":
		m.aiInput = ""
		m.draftSource = ""
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
		if m.importScreen.editingImportPath {
			if m.importScreen.importCursor == 0 {
				m.importPath += text
			} else {
				m.exportPath += text
			}
			return nil
		}
		if m.importScreen.editingExportTag {
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
		if m.aiProviderName == "disabled" || m.aiProviderName == "offline" || m.aiProviderName == "template" {
			if provider == "openai" && strings.TrimSpace(m.aiSecrets.OpenAI.APIKey) != "" {
				m.aiProviderName = "openai"
			} else if provider == "anthropic" && strings.TrimSpace(m.aiSecrets.Anthropic.APIKey) != "" {
				m.aiProviderName = "anthropic"
			}
		}

		// Rebuild the provider so the new credentials take effect immediately.
		m.aiProvider = buildProvider(m.aiProviderName, m.aiSecrets, m.aiTemplates, m.currentAITemplateSet())
		if m.onSecretsChange != nil {
			m.onSecretsChange(m.aiSecrets)
		}
		if m.onConfigChange != nil {
			m.onConfigChange(m.theme, m.aiProviderName, m.dictionaryProvider, m.aiTemplates, m.autoPlayAudio, m.strictNormalization, m.revealSpeed)
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

func (m *Model) enterPracticeModeByIndex(idx int) tea.Cmd {
	modes := []PracticeSubView{
		PracticeSubViewGender,
		PracticeSubViewConjugation,
		PracticeSubViewCase,
		PracticeSubViewAdjective,
		PracticeSubViewPreposition,
		PracticeSubViewPlural,
		PracticeSubViewSeparable,
		PracticeSubViewNumbers,
		PracticeSubViewConjunctions,
		PracticeSubViewKonjunktiv,
		PracticeSubViewPassive,
		PracticeSubViewRelative,
	}
	if idx >= 0 && idx < len(modes) {
		return m.enterPracticeMode(modes[idx])
	}
	return nil
}

func (m *Model) updatePracticeKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	switch m.practiceSubView {
	case PracticeSubViewHub:
		switch key {
		case "up", "k":
			if m.practiceHubCursor > 0 {
				m.practiceHubCursor--
			} else {
				m.practiceHubCursor = 11
			}
			return nil, true
		case "down", "j":
			if m.practiceHubCursor < 11 {
				m.practiceHubCursor++
			} else {
				m.practiceHubCursor = 0
			}
			return nil, true
		case "enter":
			return m.enterPracticeModeByIndex(m.practiceHubCursor), true
		case "1":
			m.practiceHubCursor = 0
			return m.enterPracticeMode(PracticeSubViewGender), true
		case "2":
			m.practiceHubCursor = 1
			return m.enterPracticeMode(PracticeSubViewConjugation), true
		case "3":
			m.practiceHubCursor = 2
			return m.enterPracticeMode(PracticeSubViewCase), true
		case "4":
			m.practiceHubCursor = 3
			return m.enterPracticeMode(PracticeSubViewAdjective), true
		case "5":
			m.practiceHubCursor = 4
			return m.enterPracticeMode(PracticeSubViewPreposition), true
		case "6":
			m.practiceHubCursor = 5
			return m.enterPracticeMode(PracticeSubViewPlural), true
		case "7":
			m.practiceHubCursor = 6
			return m.enterPracticeMode(PracticeSubViewSeparable), true
		case "8":
			m.practiceHubCursor = 7
			return m.enterPracticeMode(PracticeSubViewNumbers), true
		case "9":
			m.practiceHubCursor = 8
			return m.enterPracticeMode(PracticeSubViewConjunctions), true
		case "0":
			m.practiceHubCursor = 9
			return m.enterPracticeMode(PracticeSubViewKonjunktiv), true
		case "-":
			m.practiceHubCursor = 10
			return m.enterPracticeMode(PracticeSubViewPassive), true
		case "=":
			m.practiceHubCursor = 11
			return m.enterPracticeMode(PracticeSubViewRelative), true
		case "r":
			m.practiceCorrect, m.practiceTotal = 0, 0
			for _, st := range m.trainers {
				st.correct, st.total = 0, 0
			}
			m.status = "Reset all practice session scores"
			return nil, true
		case "q", "esc":
			return m.updateView(ViewDashboard), true
		}
		return nil, false

	case PracticeSubViewGender:
		if len(m.practiceItems) == 0 {
			return nil, false
		}
		if key == "esc" {
			m.practiceSubView = PracticeSubViewHub
			return nil, true
		}
		if m.practiceRevealed {
			m.advanceGenderItem()
			return nil, true
		}
		switch key {
		case "1", "d", "m": // Masculine
			m.practiceTotal++
			m.practiceRevealed = true
			if m.practiceItems[m.practiceIndex].Article == "der" {
				m.practiceCorrect++
				m.practiceLastResult = true
			} else {
				m.practiceLastResult = false
			}
			return nil, true
		case "2", "i", "f": // Feminine
			m.practiceTotal++
			m.practiceRevealed = true
			if m.practiceItems[m.practiceIndex].Article == "die" {
				m.practiceCorrect++
				m.practiceLastResult = true
			} else {
				m.practiceLastResult = false
			}
			return nil, true
		case "3", "a", "n": // Neuter
			m.practiceTotal++
			m.practiceRevealed = true
			if m.practiceItems[m.practiceIndex].Article == "das" {
				m.practiceCorrect++
				m.practiceLastResult = true
			} else {
				m.practiceLastResult = false
			}
			return nil, true
		}

	case PracticeSubViewConjugation, PracticeSubViewCase, PracticeSubViewAdjective,
		PracticeSubViewPreposition, PracticeSubViewPlural, PracticeSubViewSeparable,
		PracticeSubViewNumbers, PracticeSubViewConjunctions, PracticeSubViewKonjunktiv,
		PracticeSubViewPassive, PracticeSubViewRelative:
		return m.updateTrainerKey(m.practiceSubView, msg)
	}

	return nil, false
}

func (m *Model) updateDictionaryKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	switch key {
	case "up", "k":
		if m.dictionaryFocusResults {
			if m.dictionaryCursor == 0 {
				m.dictionaryFocusResults = false
				return nil, true
			}
			if m.dictionaryDetailView {
				if m.dictionaryDetailScroll > 0 {
					m.dictionaryDetailScroll--
				}
				return nil, true
			}
			if m.dictionaryCursor > 0 {
				m.dictionaryCursor--
				m.dictionaryDetailScroll = 0
			}
			return nil, true
		}
		if len(m.dictionarySearchHistory) > 0 {
			return m.cycleDictionaryHistory(-1), true
		}
		return nil, false
	case "down", "j":
		if !m.dictionaryFocusResults {
			if len(m.dictionaryResults) > 0 {
				m.dictionaryFocusResults = true
				return nil, true
			}
			if len(m.dictionarySearchHistory) > 0 {
				return m.cycleDictionaryHistory(1), true
			}
			return nil, false
		}
		if m.dictionaryDetailView {
			maxScroll := maxInt(0, m.dictionaryDetailTotalLines-dictionaryVisibleRows(m.activeViewContentLayout()))
			if m.dictionaryDetailScroll < maxScroll {
				m.dictionaryDetailScroll++
			}
			return nil, true
		}
		if m.dictionaryCursor < len(m.dictionaryResults)-1 {
			m.dictionaryCursor++
			m.dictionaryDetailScroll = 0
		}
		return nil, true
	case "shift+up":
		if m.dictionaryDetailScroll > 0 {
			m.dictionaryDetailScroll--
		}
		return nil, true
	case "shift+down":
		maxScroll := maxInt(0, m.dictionaryDetailTotalLines-dictionaryVisibleRows(m.activeViewContentLayout()))
		if m.dictionaryDetailScroll < maxScroll {
			m.dictionaryDetailScroll++
		}
		return nil, true
	case "pgdown":
		if m.dictionaryDetailView {
			maxScroll := maxInt(0, m.dictionaryDetailTotalLines-dictionaryVisibleRows(m.activeViewContentLayout()))
			m.dictionaryDetailScroll += 10
			if m.dictionaryDetailScroll > maxScroll {
				m.dictionaryDetailScroll = maxScroll
			}
			return nil, true
		}
		m.dictionaryCursor += 10
		if m.dictionaryCursor >= len(m.dictionaryResults) {
			m.dictionaryCursor = len(m.dictionaryResults) - 1
		}
		if m.dictionaryCursor < 0 {
			m.dictionaryCursor = 0
		}
		m.dictionaryDetailScroll = 0
		return nil, true
	case "pgup":
		if m.dictionaryDetailView {
			m.dictionaryDetailScroll -= 10
			if m.dictionaryDetailScroll < 0 {
				m.dictionaryDetailScroll = 0
			}
			return nil, true
		}
		m.dictionaryCursor -= 10
		if m.dictionaryCursor < 0 {
			m.dictionaryCursor = 0
		}
		m.dictionaryDetailScroll = 0
		return nil, true
	case "ctrl+d":
		if len(m.dictionaryResults) > 0 {
			m.dictionaryDetailView = !m.dictionaryDetailView
			m.dictionaryDetailScroll = 0
		}
		return nil, true
	case "ctrl+p":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.playDictionaryAudio(entry.Word), true
		}
		return nil, true
	case "ctrl+a":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			return m.addDictionaryEntryCmd(entry), true
		}
		return nil, true
	case "ctrl+f":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]
			m.browserSearch = entry.Word
			return tea.Batch(m.updateView(ViewBrowser), m.loadBrowserCards()), true
		}
		return nil, true
	case "ctrl+u":
		m.resetDictionarySearchState()
		return nil, true
	case "ctrl+x":
		m.dictionarySearchHistory = nil
		m.saveDictionaryHistory()
		return nil, true
	case "enter":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.recordDictionarySearch(m.dictionarySearch)
			entry := m.dictionaryResults[m.dictionaryCursor]

			m.aiInput = entry.Word + " - " + entry.Translation

			var source strings.Builder
			source.WriteString(fmt.Sprintf("Word: %s\nTranslation: %s\n", entry.Word, entry.Translation))
			if entry.WordClass != "" {
				source.WriteString(fmt.Sprintf("Class: %s\n", entry.WordClass))
			}
			if entry.Gender != "" {
				source.WriteString(fmt.Sprintf("Gender: %s\n", entry.Gender))
			}
			if entry.Forms != "" {
				source.WriteString(fmt.Sprintf("Forms: %s\n", entry.Forms))
			}
			if len(entry.Examples) > 0 {
				source.WriteString("Examples:\n")
				for _, ex := range entry.Examples {
					source.WriteString(fmt.Sprintf("- %s\n", ex))
				}
			}

			m.draftSource = source.String()

			m.updateView(ViewAI)
			m.status = "Drafting flashcard from dictionary entry"
			return m.startDrafting(), true
		}
		return nil, true
	case "backspace", "\x7f", "\x08":
		if m.dictionaryDetailView {
			return nil, true // Swallowed while viewing details
		}
		if len(m.dictionarySearch) > 0 {
			m.dictionarySearch = trimLastRune(m.dictionarySearch)
			m.dictionarySearchTimerID++
			id := m.dictionarySearchTimerID
			return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
				return debounceSearchMsg{id: id, view: ViewDictionary}
			}), true
		}
		return nil, true
	case "esc":
		if m.dictionaryDetailView {
			m.dictionaryDetailView = false
			return nil, true
		}
		m.resetDictionarySearchState()
		destView := ViewDashboard
		if m.dictionaryPreviousView != "" && m.dictionaryPreviousView != ViewDictionary {
			destView = m.dictionaryPreviousView
		}
		return m.updateView(destView), true
	case "tab":
		if len(m.dictionaryResults) > 0 {
			m.dictionaryFocusResults = !m.dictionaryFocusResults
			return nil, true
		}
		return nil, false
	case "shift+tab", "left", "right":
		return nil, false
	}

	if ch, ok := singlePrintableInput(key); ok {
		// If search is empty and it's a number key, don't handle it here to allow global navigation.
		// Standard German/English searches rarely start with a single digit, and nav parity is more important.
		if m.dictionarySearch == "" && unicode.IsDigit([]rune(ch)[0]) {
			return nil, false
		}

		m.dictionarySearch += ch
		m.dictionarySearchTimerID++
		id := m.dictionarySearchTimerID
		return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
			return debounceSearchMsg{id: id, view: ViewDictionary}
		}), true
	}
	if key == "space" {
		m.dictionarySearch += " "
		m.dictionarySearchTimerID++
		id := m.dictionarySearchTimerID
		return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
			return debounceSearchMsg{id: id, view: ViewDictionary}
		}), true
	}

	return nil, false
}

func (m *Model) updateDictionaryOverlayKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	// Toggle off on esc (if details not open) or =
	if key == "=" || (key == "esc" && !m.dictionaryDetailView) {
		m.closeDictionaryOverlay()
		return nil, true
	}

	// Intercept keys that navigate away or perform overlay-closing actions
	switch key {
	case "enter":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.dictionaryOverlayActive = false
		}
	case "ctrl+f":
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			m.dictionaryOverlayActive = false
		}
	}

	// Delegate all other key events to the standard dictionary view key handler
	return m.updateDictionaryKey(msg)
}
