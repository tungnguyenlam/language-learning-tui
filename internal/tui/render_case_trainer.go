package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderCaseTrainer(layout viewportLayout) string {
	if len(m.caseItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No case exercises found.\nTry adding more grammar content!")
	}

	item := m.caseItems[m.caseIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" CASE ENDING TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.caseTotal > 0 {
		accuracy = float64(m.caseCorrect) / float64(m.caseTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.caseCorrect, m.caseTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Sentence
	sentence := strings.Replace(item.Sentence, "{{...}}", "_____", 1)
	sentenceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(sentence)) + "\n\n")

	// Input/Answer area
	if !m.caseRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Fill in the blank:") + "\n\n")

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(30).
			Align(lipgloss.Center).
			Render(m.caseInput + "▌")
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, inputBox) + "\n")
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.caseLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		if !m.caseLastResult {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "You typed: "+m.caseInput) + "\n")
		}

		answer := strings.Replace(item.Sentence, "{{...}}", item.Answer, 1)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, sentenceStyle.Render(answer)) + "\n")

		if item.Context != "" {
			b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Grammar Context:") + "\n")
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Context)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next exercise") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "case-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.caseRevealed = false
				m.caseInput = ""
				m.caseIndex = (m.caseIndex + 1) % len(m.caseItems)
				return nil
			},
		})
	}

	return b.String()
}
