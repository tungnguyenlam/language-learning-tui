package tui

import (
	"unicode"
	"unicode/utf8"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

// reviewScreen wraps the Review view to satisfy the screen interface.
type reviewScreen struct{}

func (reviewScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderReview(layout.X, layout.Y)
}

func (reviewScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	if m.typingMode {
		switch key {
		case "enter", "space", "\r", "\n", " ":
			if !m.typingChecked {
				m.typingChecked = true
				if len(m.dueCards) > 0 {
					card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
					targetAnswer := card.Answer
					if card.Kind == core.CardKindCloze {
						targetAnswer = clozeAnswerText(card)
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
	case "up", "k":
		return m.moveReviewCursor(-1), true
	case "down", "j":
		return m.moveReviewCursor(1), true
	}
	return nil, false
}

func (m *Model) moveReviewCursor(delta int) tea.Cmd {
	next := m.cursor + delta
	if next < 0 || next >= len(m.dueCards) {
		return nil
	}
	m.cardTransitioning = true
	m.cardTransitionProgress = 0
	m.cardTransitionFrame = 0
	if delta < 0 {
		m.cardTransitionDir = -1
	} else {
		m.cardTransitionDir = 1
	}
	m.cursor = next
	m.resetReviewState()
	return m.tickCardTransition()
}
