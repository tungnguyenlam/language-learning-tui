package tui

import (
	"fmt"
	"strings"
	"unicode"

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

	// Help is a modal overlay. Keep navigation and editing keys from reaching
	// the view underneath it; only scrolling, closing, and quitting are useful
	// until the learner dismisses the shortcut reference.
	if m.showHelp {
		scrollHelp := func(delta int) {
			maxScroll := maxInt(0, m.helpTotalLines-m.helpViewportLines)
			m.helpScroll = clampInt(m.helpScroll+delta, 0, maxScroll)
		}
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "?":
			m.showHelp = false
			m.helpScroll = 0
			m.status = "Help overlay closed."
			return m, nil
		case "up", "k":
			scrollHelp(-1)
			return m, nil
		case "down", "j":
			scrollHelp(1)
			return m, nil
		case "pgup":
			scrollHelp(-maxInt(1, m.helpViewportLines))
			return m, nil
		case "pgdown", "space":
			scrollHelp(maxInt(1, m.helpViewportLines))
			return m, nil
		default:
			return m, nil
		}
	}

	// Toggle dictionary overlay on '='.
	// Practice Hub reserves '=' for the Relative Clauses trainer (#12).
	// Inside a text-input trainer (including Relative), '=' must be typed /
	// used to advance after reveal — not open Spotlight.
	if key == "=" && !m.textInputActive() && !m.practiceBlocksGlobalShortcut() {
		if !(m.activeView == ViewPractice && m.practiceSubView == PracticeSubViewHub) {
			return m, m.openDictionaryOverlay()
		}
	}

	// 0. High-priority learning mode trapping
	// Review typing mode must receive 'q' (Qualität, Quelle). Cram uses 'q'
	// to exit the active session, so only Cram exempts it from the trap.
	if m.typingMode {
		if key != "tab" && key != "shift+tab" && key != "ctrl+c" {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
	}
	if m.activeView == ViewCram && m.cramActive {
		if key != "tab" && key != "shift+tab" && key != "ctrl+c" && key != "q" {
			// Consume every other key, even if Cram does not handle it.
			// Falling through let s/w switch views and [/] change decks
			// while a session was still live.
			cmd, _ := m.updateActiveViewKey(msg)
			return m, cmd
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
		// and after reveal "any key" advances — don't quit mid-session.
		if m.practiceBlocksGlobalShortcut() {
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
		if m.practiceBlocksGlobalShortcut() {
			if cmd, handled := m.updateActiveViewKey(msg); handled {
				return m, cmd
			}
		}
		if !m.textInputActive() {
			m.showHelp = !m.showHelp
			m.helpScroll = 0
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
			return m, m.loadStatistics()
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
			return m, m.loadStatistics()
		}
	}

	// 3. View-specific keys
	if cmd, handled := m.updateActiveViewKey(msg); handled {
		return m, cmd
	}

	if cmd, handled := m.updateNumberKey(msg); handled {
		return m, cmd
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
		// Practice Hub `/` filter is a real text field — without this, `q`
		// quits the app and digits jump to other views mid-filter.
		(m.activeView == ViewPractice && m.practiceSubView == PracticeSubViewHub && m.practiceFilterFocus) ||
		m.drafting // AI drafting also uses input
}

// trainerInputActive reports whether a generic text-input trainer is waiting
// for the learner to type an answer (not yet revealed).
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

// practiceBlocksGlobalShortcut reports whether a practice trainer should
// consume single-letter global shortcuts ("q", "?", "=") instead of letting
// them quit, open help, or open the dictionary overlay.
//
// Covers: typing into a generic trainer, the post-reveal "press any key"
// advance step for generic trainers, and the Gender trainer's advance step.
func (m *Model) practiceBlocksGlobalShortcut() bool {
	if m.activeView != ViewPractice {
		return false
	}
	if m.practiceSubView == PracticeSubViewGender {
		return len(m.practiceItems) > 0 && m.practiceRevealed
	}
	if !isGenericTrainer(m.practiceSubView) {
		return false
	}
	st, ok := m.trainers[m.practiceSubView]
	return ok && len(st.items) > 0
}

func (m *Model) updateNumberKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()
	if (m.activeView == ViewDictionary || m.dictionaryOverlayActive) && (m.dictionaryFocusResults || m.dictionaryDetailView) {
		return nil, false
	}
	// Inside a practice sub-trainer, number keys belong to the trainer
	// (der/die/das, typed answers). Don't jump to global views mid-exercise.
	if m.activeView == ViewPractice && m.practiceSubView != PracticeSubViewHub {
		return nil, false
	}
	// Active Cram uses a/h/g/e for grading; digit view-shortcuts must not
	// abandon the session. Review still allows 0-9 view switching (E2E and
	// day-to-day nav rely on it); grading digits are handled earlier when revealed.
	if m.activeView == ViewCram && m.cramActive {
		return nil, false
	}
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

func (m *Model) handlePaste(text string) tea.Cmd {
	// Clean text: strip newlines/tabs which break single-line inputs
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\t", " ")
	if text == "" {
		return nil
	}

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

	if m.activeView == ViewPractice && m.practiceSubView == PracticeSubViewHub && m.practiceFilterFocus {
		m.practiceFilter += text
		m.practiceHubCursor = 0
		return nil
	}

	if m.activeView == ViewAnkiWeb && m.ankiWebScreen != nil && m.ankiWebScreen.editingQuery {
		m.ankiWebScreen.query += text
		return nil
	}

	// Dictionary search (full tab or Spotlight overlay) — same debounce as typing.
	if !m.dictionaryDetailView && (m.dictionaryOverlayActive || m.activeView == ViewDictionary) {
		m.dictionarySearch += text
		return m.debounceSearch(ViewDictionary)
	}

	// Review typing mode: paste into the answer buffer before check.
	if m.activeView == ViewReview && m.typingMode && !m.typingChecked {
		m.typedAnswer += text
		return nil
	}

	// Practice text-input trainers: paste into the answer while waiting.
	if m.activeView == ViewPractice && isGenericTrainer(m.practiceSubView) {
		st := m.trainerStateFor(m.practiceSubView)
		if len(st.items) > 0 && !st.revealed {
			st.input += text
		}
		return nil
	}

	return nil
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

	// Delegate remaining keys to the Dictionary screen contract.
	return (dictionaryScreen{}).HandleKey(m, msg)
}
