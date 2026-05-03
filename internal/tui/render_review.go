package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderReview(x, y int) string {
	if len(m.dueCards) == 0 {
		if m.bookmarkFilter {
			return "Review (Bookmarked)\n\nNo bookmarked cards due."
		}
		return "Review\n\nNo cards due."
	}
	card := m.dueCards[m.cursor]
	bookmark := "Bookmark: off"
	if card.Bookmarked {
		bookmark = "Bookmark: on"
	}
	leech := ""
	if card.Leech {
		leech = " | LEECH"
	}
	suspended := ""
	if card.Suspended {
		suspended = " | SUSPENDED"
	}
	filterBanner := ""
	if m.bookmarkFilter {
		filterBanner = " (Bookmarked)"
	}
	keys := "b toggle | x suspend | B filter | u undo | r history | p audio"
	if m.bookmarkFilter {
		keys = "b toggle | x suspend | B all cards | u undo | r history | p audio"
	}
	audioIndicator := ""
	if card.Audio != "" {
		audioIndicator = " [Audio]"
	}

	gradeAgain := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Again")
	gradeHard := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("Hard")
	gradeGood := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("Good")
	gradeEasy := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Easy")

	var answer string
	if card.Kind == core.CardKindMCQ && len(card.Choices) > 0 {
		if m.revealState == RevealRevealed {
			if m.mcqAnswered {
				feedback := "Incorrect"
				if m.mcqCorrect {
					feedback = "Correct"
				}
				answer = fmt.Sprintf("%s: %s\n\n%s\n\nGrade: a %s | h %s | g %s | e %s", feedback, card.Answer, renderMCQChoices(card.Choices, m.mcqChoice), gradeAgain, gradeHard, gradeGood, gradeEasy)
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 10, Width: 5, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 10, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 10, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 10, Width: 4, Height: 1})
			} else {
				answer = fmt.Sprintf("1-4 select answer\n\n%s\n\nGrade: a %s | h %s | g %s | e %s", renderMCQChoices(card.Choices, m.mcqChoice), gradeAgain, gradeHard, gradeGood, gradeEasy)
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 8, Width: 5, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 8, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 8, Width: 4, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 8, Width: 4, Height: 1})
			}
		} else if m.revealState == RevealRevealing {
			// Show gradual reveal animation with blocks
			progress := int(m.revealProgress)
			if progress > 100 {
				progress = 100
			}
			// Calculate how many characters to reveal
			fullText := card.Answer
			numChars := len([]rune(fullText))
			revealedChars := (numChars * progress) / 100
			if revealedChars < 0 {
				revealedChars = 0
			}
			if revealedChars > numChars {
				revealedChars = numChars
			}
			revealedRunes := []rune(fullText)[:revealedChars]
			remainingBlocks := numChars - revealedChars
			animationText := string(revealedRunes) + strings.Repeat("▌", remainingBlocks-1)
			answer = fmt.Sprintf("1-4 select answer\n\n%s", renderMCQChoices(card.Choices, m.mcqChoice))
			// Replace the answer in choices display
			answer = strings.Replace(answer, fullText, animationText, 1)
		} else {
			answer = "Press space or enter to reveal choices."
		}
	} else if m.revealState == RevealRevealed {
		answer = fmt.Sprintf("%s\n\nGrade: a %s | h %s | g %s | e %s", card.Answer, gradeAgain, gradeHard, gradeGood, gradeEasy)
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: x + 7, Y: y + 6, Width: 5, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: x + 17, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: x + 26, Y: y + 6, Width: 4, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: x + 35, Y: y + 6, Width: 4, Height: 1})
	} else if m.revealState == RevealRevealing {
		// Show gradual reveal animation with blocks
		progress := int(m.revealProgress)
		if progress > 100 {
			progress = 100
		}
		// Calculate how many characters to reveal
		fullText := card.Answer
		numChars := len([]rune(fullText))
		revealedChars := (numChars * progress) / 100
		if revealedChars < 0 {
			revealedChars = 0
		}
		if revealedChars > numChars {
			revealedChars = numChars
		}
		revealedRunes := []rune(fullText)[:revealedChars]
		remainingBlocks := numChars - revealedChars
		answer = string(revealedRunes) + "▌" + strings.Repeat("▌", remainingBlocks-1)
	} else {
		answer = "Press space or enter to reveal."
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	header := fmt.Sprintf("%s%s %d/%d", titleStyle.Render("Review"), filterBanner, m.cursor+1, len(m.dueCards))

	mature := ""
	if card.Mature {
		mature = " ✨"
	}

	promptDisplay := card.Prompt
	if card.Kind == core.CardKindCloze {
		// Highlight the [...] or [hint] in cloze prompt
		promptDisplay = strings.Replace(promptDisplay, "[", lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("["), 1)
		promptDisplay = strings.Replace(promptDisplay, "]", lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("]"), 1)
	}

	view := fmt.Sprintf("%s\n%s | %s\n%s%s%s\n\n%s%s\n\n%s", header, bookmark, keys, leech, suspended, audioIndicator, promptDisplay, mature, answer)
	if m.showReviewHistory && m.reviewHistoryCard == card.ID {
		view += "\n\n" + m.renderReviewHistory(card.Prompt)
	}
	return view
}

func renderMCQChoices(choices []string, selected int) string {
	var b strings.Builder
	for i, choice := range choices {
		prefix := "  "
		mark := " "
		if i == selected {
			mark = ">"
		}
		b.WriteString(fmt.Sprintf("%s%d: %s%s\n", prefix, i+1, mark, choice))
	}
	return strings.TrimRight(b.String(), "\n")
}
