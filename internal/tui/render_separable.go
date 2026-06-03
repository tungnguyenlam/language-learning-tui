package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderSeparableTrainer(layout viewportLayout) string {
	if len(m.separableItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No separable verb exercises found.")
	}

	item := m.separableItems[m.separableIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" SEPARABLE VERB TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.separableTotal > 0 {
		accuracy = float64(m.separableCorrect) / float64(m.separableTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.separableCorrect, m.separableTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Sentence
	sentenceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(item.Sentence)) + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("Verb: "+item.Verb+" ("+item.Meaning+")")) + "\n\n")

	// Input/Answer area
	if !m.separableRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Enter the missing prefix:") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.separableInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.separableLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.separableLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.separableInput) + "\n")
		}

		answerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center)

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, answerBox.Render(item.Answer)) + "\n")

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "separable-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.separableRevealed = false
				m.separableInput = ""
				m.separableIndex = (m.separableIndex + 1) % len(m.separableItems)
				return nil
			},
		})
	}

	return b.String()
}
