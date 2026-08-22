package tui

import (
	"strconv"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

// cramScreen wraps the Cram view to satisfy the screen interface.
type cramScreen struct{}

func (cramScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderCramAt(layout)
}

func (cramScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if m.cramActive {
		switch msg.String() {
		case "enter", "space", "\r", "\n", " ":
			if m.cramRevealed {
				return m.gradeCramCard(core.GradeGood), true
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
		case "d":
			return m.lookupCramCardInDictionary(), true
		case "a", "1":
			if !m.cramRevealed {
				return nil, true
			}
			return m.gradeCramCard(core.GradeAgain), true
		case "h", "2":
			if !m.cramRevealed {
				return nil, true
			}
			return m.gradeCramCard(core.GradeHard), true
		case "g", "3":
			if !m.cramRevealed {
				return nil, true
			}
			return m.gradeCramCard(core.GradeGood), true
		case "e", "4":
			if !m.cramRevealed {
				return nil, true
			}
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
