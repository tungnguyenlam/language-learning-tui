package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderConjugation(layout viewportLayout) string {
	if len(m.conjugationItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No verbs found for practice.\nTry adding some verbs to your decks!")
	}

	item := m.conjugationItems[m.conjugationIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" VERB CONJUGATION TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.conjugationTotal > 0 {
		accuracy = float64(m.conjugationCorrect) / float64(m.conjugationTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.conjugationCorrect, m.conjugationTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Verb
	verbStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, verbStyle.Render(item.German)) + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("("+item.English+")")) + "\n\n")

	// Person prompt
	persons := []string{"ich", "du", "er/sie/es", "wir", "ihr", "sie/Sie"}
	person := persons[m.conjugationPerson]

	promptStyle := lipgloss.NewStyle().Bold(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, promptStyle.Render(fmt.Sprintf("Conjugate for: %s", person))) + "\n\n")

	// Input/Answer area
	if !m.conjugationRevealed {
		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.conjugationInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.conjugationLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n")

		answerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center)

		if !m.conjugationLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.conjugationInput) + "\n")
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, answerBox.Render(m.conjugationAnswer)) + "\n")

		if item.Example != "" {
			b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Example:") + "\n")
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Example)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next verb") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "conjugation-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.conjugationRevealed = false
				m.conjugationInput = ""
				m.conjugationIndex = (m.conjugationIndex + 1) % len(m.conjugationItems)
				m.conjugationPerson = int(time.Now().UnixNano() % 6)
				return nil
			},
		})
	}

	return b.String()
}
