package tui

import (
	tea "charm.land/bubbletea/v2"
)

// practiceScreen wraps the Practice view to satisfy the screen interface.
type practiceScreen struct{}

func (practiceScreen) Render(m *Model, layout viewportLayout) string {
	return m.renderPractice(layout)
}

func (practiceScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	key := msg.String()

	switch m.practiceSubView {
	case PracticeSubViewHub:
		if m.practiceFilterFocus {
			switch key {
			case "esc":
				m.practiceFilterFocus = false
				m.practiceFilter = ""
				return nil, true
			case "backspace":
				if len(m.practiceFilter) > 0 {
					runes := []rune(m.practiceFilter)
					m.practiceFilter = string(runes[:len(runes)-1])
				}
				if len(m.practiceFilter) == 0 {
					m.practiceFilterFocus = false
				}
				return nil, true
			case "enter", "down", "up", "tab":
				m.practiceFilterFocus = false
				return nil, true
			}
			if ch, ok := singlePrintableInput(key); ok {
				m.practiceFilter += ch
				m.practiceHubCursor = 0
				return nil, true
			}
			return nil, true
		}

		if key == "/" {
			m.practiceFilterFocus = true
			return nil, true
		}

		visible := m.visiblePracticeHubModes()
		n := len(visible)
		switch key {
		case "up", "k":
			if n > 0 {
				m.practiceHubCursor = (m.practiceHubCursor - 1 + n) % n
			}
			return nil, true
		case "down", "j":
			if n > 0 {
				m.practiceHubCursor = (m.practiceHubCursor + 1) % n
			}
			return nil, true
		case "enter":
			if n == 0 {
				return nil, true
			}
			m.practiceHubCursor = clampInt(m.practiceHubCursor, 0, n-1)
			return m.enterPracticeMode(visible[m.practiceHubCursor].sub), true
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
			if m.practiceFilter != "" {
				m.practiceFilter = ""
				return nil, true
			}
			return m.updateView(ViewDashboard), true
		}
		return nil, false

	case PracticeSubViewGender:
		if key == "esc" {
			m.practiceSubView = PracticeSubViewHub
			return nil, true
		}
		if len(m.practiceItems) == 0 || m.practiceIndex < 0 || m.practiceIndex >= len(m.practiceItems) {
			return nil, false
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
