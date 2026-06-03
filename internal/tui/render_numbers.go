package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderNumberTrainer(layout viewportLayout) string {
	if len(m.numberItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No number exercises found.")
	}

	item := m.numberItems[m.numberIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" NUMBER & TIME TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.numberTotal > 0 {
		accuracy = float64(m.numberCorrect) / float64(m.numberTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.numberCorrect, m.numberTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Question
	questionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, questionStyle.Render(item.Question)) + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Help)) + "\n\n")

	// Input/Answer area
	if !m.numberRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Enter the German translation:") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(40).
			Align(lipgloss.Center).
			Render(m.numberInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.numberLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.numberLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.numberInput) + "\n")
		}

		answerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(40).
			Align(lipgloss.Center)

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, answerBox.Render(item.Answer)) + "\n")

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "number-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.numberRevealed = false
				m.numberInput = ""
				m.numberIndex = (m.numberIndex + 1) % len(m.numberItems)
				return nil
			},
		})
	}

	return b.String()
}
