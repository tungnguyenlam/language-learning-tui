package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderPrepositionTrainer(layout viewportLayout) string {
	if len(m.prepItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No preposition exercises found.\nTry adding more grammar content!")
	}

	item := m.prepItems[m.prepIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" PREPOSITION TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.prepTotal > 0 {
		accuracy = float64(m.prepCorrect) / float64(m.prepTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.prepCorrect, m.prepTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Sentence
	sentence := strings.Replace(item.Sentence, "{{...}}", "_____", 1)
	sentenceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(sentence)) + "\n\n")

	// Input/Answer area
	if !m.prepRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Fill in the blank (enter the preposition or article):") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.prepInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.prepLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.prepLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.prepInput) + "\n")
		}

		answer := strings.Replace(item.Sentence, "{{...}}", item.Answer, 1)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(answer)) + "\n")

		if item.Context != "" {
			b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Grammar Context:") + "\n")
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Context)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "preposition-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.prepRevealed = false
				m.prepInput = ""
				m.prepIndex = (m.prepIndex + 1) % len(m.prepItems)
				return nil
			},
		})
	}

	return b.String()
}
