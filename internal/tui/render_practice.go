package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderPractice(layout viewportLayout) string {
	if len(m.practiceItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No nouns with articles found in your decks.\nAdd some cards with 'der', 'die', or 'das' to start practicing!")
	}

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan).MarginBottom(1)
	b.WriteString(titleStyle.Render("GENDER TRAINER") + "\n\n")

	scoreStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	accuracy := 0.0
	if m.practiceTotal > 0 {
		accuracy = float64(m.practiceCorrect) / float64(m.practiceTotal) * 100
	}
	b.WriteString(scoreStyle.Render(fmt.Sprintf("Score: %d/%d (%.0f%%)", m.practiceCorrect, m.practiceTotal, accuracy)) + "\n\n")

	item := m.practiceItems[m.practiceIndex]

	wordStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Padding(1, 4).Border(lipgloss.RoundedBorder())
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, wordStyle.Render(item.Word)) + "\n\n")

	if !m.practiceRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Which article?") + "\n\n")

		btnStyle := lipgloss.NewStyle().Padding(0, 2).Margin(0, 1).Background(lipgloss.Color("235"))
		derBtn := btnStyle.Render("der [1/d/m]")
		dieBtn := btnStyle.Render("die [2/i/f]")
		dasBtn := btnStyle.Render("das [3/a/n]")

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, derBtn, dieBtn, dasBtn)) + "\n")
	} else {
		resultStyle := lipgloss.NewStyle().Bold(true).Padding(0, 2).MarginBottom(1)
		if m.practiceLastResult {
			resultStyle = resultStyle.Foreground(lipgloss.Color("46")).SetString("CORRECT")
		} else {
			resultStyle = resultStyle.Foreground(lipgloss.Color("196")).SetString("INCORRECT")
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.String()) + "\n")

		fullWord := fmt.Sprintf("%s %s", item.Article, item.Word)
		fullWordStyle := lipgloss.NewStyle().Bold(true).Foreground(colorPink)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, fullWordStyle.Render(fullWord)) + "\n")

		if item.Meaning != "" {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Meaning)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next noun") + "\n")
	}

	return b.String()
}
