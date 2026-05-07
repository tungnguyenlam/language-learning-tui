package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderReview(x, y int) string {
	if len(m.dueCards) == 0 {
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Render("Review")
		message := "No cards due."
		if m.bookmarkFilter {
			title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Render("Review (Bookmarked)")
			message = "No bookmarked cards due."
		}
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2).
			Width(46).
			Render(message + "\n\nUse [ and ] to switch decks or B to toggle the bookmark filter.")
		return title + "\n\n" + emptyBox
	}
	cursor := clampInt(m.cursor, 0, len(m.dueCards)-1)
	card := m.dueCards[cursor]
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

	gradeAgainText := "Again"
	gradeHardText := "Hard"
	gradeGoodText := "Good"
	gradeEasyText := "Easy"

	if m.reviewPredictions != nil {
		gradeAgainText = fmt.Sprintf("Again (%s)", formatDuration(m.reviewPredictions[core.GradeAgain]))
		gradeHardText = fmt.Sprintf("Hard (%s)", formatDuration(m.reviewPredictions[core.GradeHard]))
		gradeGoodText = fmt.Sprintf("Good (%s)", formatDuration(m.reviewPredictions[core.GradeGood]))
		gradeEasyText = fmt.Sprintf("Easy (%s)", formatDuration(m.reviewPredictions[core.GradeEasy]))
	}

	gradeAgain := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(gradeAgainText)
	gradeHard := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render(gradeHardText)
	gradeGood := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(gradeGoodText)
	gradeEasy := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render(gradeEasyText)

	gaW := lipgloss.Width(gradeAgain)
	ghW := lipgloss.Width(gradeHard)
	ggW := lipgloss.Width(gradeGood)
	geW := lipgloss.Width(gradeEasy)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	header := fmt.Sprintf("%s%s %d/%d", titleStyle.Render("Review"), filterBanner, cursor+1, len(m.dueCards))

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

	width, _ := m.activePanelSize()
	cardWidth := maxInt(30, width-6)
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 2).
		Width(cardWidth)

	var answer string
	answerYOffset := 0

	extraDisplay := ""
	if card.Extra != "" {
		extraStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
		extraDisplay = "\n" + extraStyle.Render("💡 "+card.Extra)
	}

	headerSection := fmt.Sprintf("%s\n%s | %s\n%s%s%s", header, bookmark, keys, leech, suspended, audioIndicator)
	headerLines := strings.Count(headerSection, "\n") + 1

	cardX := x + cardStyle.GetMarginLeft() + cardStyle.GetBorderLeftSize() + cardStyle.GetPaddingLeft()
	cardY := y + headerLines + 1 + cardStyle.GetMarginTop() + cardStyle.GetBorderTopSize() + cardStyle.GetPaddingTop()

	if card.Kind == core.CardKindMCQ && len(card.Choices) > 0 {
		if m.revealState == RevealRevealed {
			mcqChoices := renderMCQChoices(card.Choices, m.mcqChoice)
			if m.mcqAnswered {
				feedback := "Incorrect"
				feedbackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
				if m.mcqCorrect {
					feedback = "Correct"
					feedbackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
				}
				answer = fmt.Sprintf("%s: %s%s\n\n%s\n\nGrade: a %s | h %s | g %s | e %s", feedbackStyle.Render(feedback), card.Answer, extraDisplay, mcqChoices, gradeAgain, gradeHard, gradeGood, gradeEasy)
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n%s: %s\n\n%s\n\nGrade: ", promptDisplay, mature, feedback, card.Answer, mcqChoices), "\n")

				labelAgain := "Grade: a "
				labelHard := fmt.Sprintf("Grade: a %s | h ", gradeAgainText)
				labelGood := fmt.Sprintf("Grade: a %s | h %s | g ", gradeAgainText, gradeHardText)
				labelEasy := fmt.Sprintf("Grade: a %s | h %s | g %s | e ", gradeAgainText, gradeHardText, gradeGoodText)

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			} else {
				answer = fmt.Sprintf("1-4 select answer%s\n\n%s\n\nGrade: a %s | h %s | g %s | e %s", extraDisplay, mcqChoices, gradeAgain, gradeHard, gradeGood, gradeEasy)
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n1-4 select answer\n\n%s\n\nGrade: ", promptDisplay, mature, mcqChoices), "\n")

				labelAgain := "Grade: a "
				labelHard := fmt.Sprintf("Grade: a %s | h ", gradeAgainText)
				labelGood := fmt.Sprintf("Grade: a %s | h %s | g ", gradeAgainText, gradeHardText)
				labelEasy := fmt.Sprintf("Grade: a %s | h %s | g %s | e ", gradeAgainText, gradeHardText, gradeGoodText)

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			}
		} else if m.revealState == RevealRevealing {
			answer = fmt.Sprintf("1-4 select answer\n\n%s", renderMCQChoices(card.Choices, m.mcqChoice))
		} else {
			answer = "Press space or enter to reveal choices."
		}
	} else if m.revealState == RevealRevealed {
		answer = fmt.Sprintf("%s%s\n\nGrade: a %s | h %s | g %s | e %s", card.Answer, extraDisplay, gradeAgain, gradeHard, gradeGood, gradeEasy)
		answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n%s\n\nGrade: ", promptDisplay, mature, card.Answer), "\n")

		labelAgain := "Grade: a "
		labelHard := fmt.Sprintf("Grade: a %s | h ", gradeAgainText)
		labelGood := fmt.Sprintf("Grade: a %s | h %s | g ", gradeAgainText, gradeHardText)
		labelEasy := fmt.Sprintf("Grade: a %s | h %s | g %s | e ", gradeAgainText, gradeHardText, gradeGoodText)

		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
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
		answer = string(revealedRunes) + "▌" + strings.Repeat("▌", maxInt(0, remainingBlocks-1))
	} else {
		answer = "Press space or enter to reveal."
	}

	view := fmt.Sprintf("%s\n\n%s", headerSection, cardStyle.Render(promptDisplay+mature+"\n\n"+answer))
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
