package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderAdjectiveTrainer(layout viewportLayout) string {
	if len(m.adjItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No adjective exercises found.\nTry adding more grammar content!")
	}

	item := m.adjItems[m.adjIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" ADJECTIVE ENDING TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.adjTotal > 0 {
		accuracy = float64(m.adjCorrect) / float64(m.adjTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.adjCorrect, m.adjTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Sentence
	sentence := strings.Replace(item.Sentence, "{{...}}", "_____", 1)
	sentenceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(sentence)) + "\n\n")

	// Input/Answer area
	if !m.adjRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Fill in the blank (enter the correct adjective ending):") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.adjInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.adjLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.adjLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.adjInput) + "\n")
		}

		answer := strings.Replace(item.Sentence, "{{...}}", item.Answer, 1)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(answer)) + "\n")

		if item.Context != "" {
			b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Grammar Context:") + "\n")
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Context)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")
	}

	return b.String()
}
