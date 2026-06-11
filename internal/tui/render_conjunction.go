package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderConjunctionTrainer(layout viewportLayout) string {
	if len(m.conjItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No conjunction exercises found.")
	}

	item := m.conjItems[m.conjIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" CONJUNCTIONS & WORD ORDER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.conjTotal > 0 {
		accuracy = float64(m.conjCorrect) / float64(m.conjTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.conjCorrect, m.conjTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Sentence
	sentenceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(item.Sentence)) + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("Meaning: "+item.Meaning)) + "\n\n")

	// Input/Answer area
	if !m.conjRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Enter the missing word:") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.conjInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n\n")

		// Display hint if requested
		if m.practiceShowHint && item.Hint != "" {
			hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Italic(true)
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, hintStyle.Render("Hint: "+item.Hint)) + "\n")
		} else {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("Press 'h' for a hint")) + "\n")
		}
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.conjLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.conjLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.conjInput) + "\n")
		}

		answerBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center)

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, answerBox.Render(item.Answer)) + "\n\n")

		// Display pedagogical explanation
		explanationStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Italic(true)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, explanationStyle.Render(item.Explanation)) + "\n")

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "conjunction-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.conjRevealed = false
				m.practiceShowHint = false
				m.conjInput = ""
				m.conjIndex = (m.conjIndex + 1) % len(m.conjItems)
				return nil
			},
		})
	}

	return b.String()
}
