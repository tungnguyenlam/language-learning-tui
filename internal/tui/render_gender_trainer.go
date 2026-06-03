package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderGenderTrainer(layout viewportLayout) string {
	if len(m.practiceItems) == 0 {
		return lipgloss.NewStyle().
			Width(layout.Width).
			Height(layout.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No nouns found for practice.\nTry adding some nouns to your decks!")
	}

	item := m.practiceItems[m.practiceIndex]

	var b strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" GENDER TRAINER ")) + "\n\n")

	// Score
	accuracy := 0.0
	if m.practiceTotal > 0 {
		accuracy = float64(m.practiceCorrect) / float64(m.practiceTotal) * 100
	}
	scoreStr := fmt.Sprintf("Score: %d/%d (%.0f%%)", m.practiceCorrect, m.practiceTotal, accuracy)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(scoreStr)) + "\n\n")

	// Word
	wordStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, wordStyle.Render(item.Word)) + "\n\n")

	// Options
	if !m.practiceRevealed {
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Which article?") + "\n\n")

		options := []string{
			keyStyle.Render("[1/d/m]") + " der",
			keyStyle.Render("[2/i/f]") + " die",
			keyStyle.Render("[3/a/n]") + " das",
		}
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, strings.Join(options, "   ")) + "\n")

		startX := layout.X + (layout.Width-39)/2
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "gender-opt-der",
			View:   ViewPractice,
			X:      startX,
			Y:      layout.Y + 11,
			Width:  11,
			Height: 1,
			Action: func() tea.Cmd {
				m.practiceTotal++
				m.practiceRevealed = true
				if m.practiceItems[m.practiceIndex].Article == "der" {
					m.practiceCorrect++
					m.practiceLastResult = true
				} else {
					m.practiceLastResult = false
				}
				return nil
			},
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "gender-opt-die",
			View:   ViewPractice,
			X:      startX + 14,
			Y:      layout.Y + 11,
			Width:  11,
			Height: 1,
			Action: func() tea.Cmd {
				m.practiceTotal++
				m.practiceRevealed = true
				if m.practiceItems[m.practiceIndex].Article == "die" {
					m.practiceCorrect++
					m.practiceLastResult = true
				} else {
					m.practiceLastResult = false
				}
				return nil
			},
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "gender-opt-das",
			View:   ViewPractice,
			X:      startX + 28,
			Y:      layout.Y + 11,
			Width:  11,
			Height: 1,
			Action: func() tea.Cmd {
				m.practiceTotal++
				m.practiceRevealed = true
				if m.practiceItems[m.practiceIndex].Article == "das" {
					m.practiceCorrect++
					m.practiceLastResult = true
				} else {
					m.practiceLastResult = false
				}
				return nil
			},
		})
	} else {
		var resultStyle lipgloss.Style
		resultText := ""
		if m.practiceLastResult {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
			resultText = "CORRECT!"
		} else {
			resultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
			resultText = "INCORRECT"
		}

		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, resultStyle.Render(resultText)) + "\n\n")

		revealedWord := item.Article + " " + item.Word
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, wordStyle.Render(revealedWord)) + "\n")

		if item.Meaning != "" {
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render(item.Meaning)) + "\n")
		}

		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press any key for next noun") + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "gender-next",
			View:   ViewPractice,
			X:      layout.X,
			Y:      layout.Y,
			Width:  layout.Width,
			Height: layout.Height,
			Action: func() tea.Cmd {
				m.practiceRevealed = false
				m.practiceIndex = (m.practiceIndex + 1) % len(m.practiceItems)
				return nil
			},
		})
	}

	return b.String()
}
