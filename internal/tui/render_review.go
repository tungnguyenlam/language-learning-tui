package tui

import (
	"fmt"
	"strings"
	"time"

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

		// Enhanced empty state with keyboard guidance
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		shortcutGuide := fmt.Sprintf("Use %s and %s to switch decks or %s to toggle the bookmark filter.\nPress %s for Custom Study (Cram Mode).",
			keyStyle.Render("["), keyStyle.Render("]"), keyStyle.Render("B"), keyStyle.Render("c"))

		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2).
			Width(60).
			Align(lipgloss.Center).
			Render(message + "\n\n" + shortcutGuide)
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

	// Enhanced keyboard shortcut display with visual highlighting
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	keys := fmt.Sprintf("%s toggle | %s suspend | %s filter | %s undo | %s history | %s audio | %s type",
		keyStyle.Render("b"), keyStyle.Render("x"), keyStyle.Render("B"), keyStyle.Render("u"), keyStyle.Render("r"), keyStyle.Render("p"), keyStyle.Render("t"))
	if m.bookmarkFilter {
		keys = fmt.Sprintf("%s toggle | %s suspend | %s all cards | %s undo | %s history | %s audio",
			keyStyle.Render("b"), keyStyle.Render("x"), keyStyle.Render("B"), keyStyle.Render("u"), keyStyle.Render("r"), keyStyle.Render("p"))
	}
	audioIndicator := ""
	if card.Audio != "" {
		audioIndicator = " [Audio]"
	}
	deckMeta := fmt.Sprintf("Deck: %s", m.deckNameByID(card.DeckID))
	if len(card.Tags) > 0 {
		deckMeta += " | Tags: #" + strings.Join(card.Tags, " #")
	}
	cardTypeLabel := "Flashcard"
	switch card.Kind {
	case core.CardKindMCQ:
		cardTypeLabel = "Multiple Choice"
	case core.CardKindCloze:
		cardTypeLabel = "Cloze"
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

	// Session progress bar and timer
	sessionTotal := m.sessionReviewed + len(m.dueCards)
	sessionPercentage := 0.0
	if sessionTotal > 0 {
		sessionPercentage = float64(m.sessionReviewed) / float64(sessionTotal)
	}
	sessionBar := progressBar(15, sessionPercentage, "46", "238")

	duration := time.Since(m.sessionStartTime)
	timerStr := fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)

	sessionProgress := fmt.Sprintf(" | Session: %s %d%% | ⏱ %s", sessionBar, int(sessionPercentage*100), timerStr)

	mature := ""
	if card.Mature {
		mature = " ✨"
	}

	// Determine card state/difficulty for badge
	stateBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Padding(0, 1)
	cardBorderColor := "62"
	if card.Mature {
		stateBadge = stateBadge.Background(lipgloss.Color("34")).SetString("MATURE")
		cardBorderColor = "34"
	} else if card.Reviews > 0 {
		stateBadge = stateBadge.Background(lipgloss.Color("39")).SetString("LEARNING")
		cardBorderColor = "39"
	} else {
		stateBadge = stateBadge.Background(lipgloss.Color("208")).SetString("NEW")
		cardBorderColor = "208"
	}

	headerSection := fmt.Sprintf("%s%s %s\n%s | Type: %s\n%s | %s\n%s%s%s", header, sessionProgress, stateBadge.Render(), deckMeta, cardTypeLabel, bookmark, keys, leech, suspended, audioIndicator)
	if m.focusMode {
		headerSection = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render("Focus Mode Active (f to exit)")
	}
	headerLines := strings.Count(headerSection, "\n") + 1

	promptDisplay := card.Prompt
	if card.Kind == core.CardKindCloze {
		// Highlight all [...] or [hint] in cloze prompt
		bracketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		promptDisplay = strings.ReplaceAll(promptDisplay, "[", bracketStyle.Render("["))
		promptDisplay = strings.ReplaceAll(promptDisplay, "]", bracketStyle.Render("]"))
	}

	// Enhanced prompt styling for better visibility
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	promptDisplay = promptStyle.Render(promptDisplay)

	width, _ := m.activePanelSize()
	cardWidth := maxInt(30, width-6)

	// Enhanced card styling with better visual hierarchy
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cardBorderColor)).
		Padding(1, 2).
		Width(cardWidth).
		Align(lipgloss.Center)

	// Special styling for revealed cards
	if m.revealState == RevealRevealed {
		cardStyle = cardStyle.
			BorderForeground(lipgloss.Color("81")).
			Bold(true)
	}

	var answer string
	answerYOffset := 0

	// Calculate card position for hitboxes
	cardX := x + cardStyle.GetMarginLeft() + cardStyle.GetBorderLeftSize() + cardStyle.GetPaddingLeft()
	cardY := y + headerLines + 1 + cardStyle.GetMarginTop() + cardStyle.GetBorderTopSize() + cardStyle.GetPaddingTop()

	// Enhanced answer styling
	answerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	extraDisplay := ""
	if card.Extra != "" {
		extraStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("159")).
			Italic(true)
		extraDisplay = "\n\n" + extraStyle.Render("💡 CONTEXT: "+card.Extra)
	}

	// Typing mode display
	if m.typingMode && m.revealState != RevealRevealing {
		typingBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1).
			Width(maxInt(30, width-6))

		inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		inputDisplay := m.typedAnswer + "_"

		if m.typingChecked {
			if m.typingCorrect {
				correctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
				inputDisplay = correctStyle.Render("✓ " + m.typedAnswer)
				typingBoxStyle = typingBoxStyle.BorderForeground(lipgloss.Color("46"))
			} else {
				wrongStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
				inputDisplay = wrongStyle.Render("✗ " + m.typedAnswer)
				typingBoxStyle = typingBoxStyle.BorderForeground(lipgloss.Color("196"))
			}
			targetAnswer := card.Answer
			if card.Kind == core.CardKindCloze && len(card.Choices) > 0 {
				targetAnswer = card.Choices[0]
			}
			typingContent := fmt.Sprintf("Your answer: %s\nCorrect: %s%s\n\nGrade: %s %s | %s %s | %s %s | %s %s",
				inputDisplay, answerStyle.Render(targetAnswer), extraDisplay,
				keyStyle.Render("a"), gradeAgain, keyStyle.Render("h"), gradeHard, keyStyle.Render("g"), gradeGood, keyStyle.Render("e"), gradeEasy)
			answer = typingBoxStyle.Render(typingContent)
			answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\nYour answer:", promptDisplay, mature), "\n") + 2

			labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
			labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"))
			labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"))
			labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"), gradeGoodText, keyStyle.Render("e"))

			m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
		} else {
			typingContent := fmt.Sprintf("Type your answer:\n\n%s", inputStyle.Render(inputDisplay))
			answer = typingBoxStyle.Render(typingContent)
		}
	} else if card.Kind == core.CardKindMCQ && len(card.Choices) > 0 {
		if m.revealState == RevealRevealed {
			mcqChoices := renderMCQChoices(card.Choices, m.mcqChoice)
			if m.mcqAnswered {
				feedback := "Incorrect"
				feedbackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
				if m.mcqCorrect {
					feedback = "Correct"
					feedbackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
				}
				answer = fmt.Sprintf("%s: %s%s\n\n%s\n\nGrade: %s %s | %s %s | %s %s | %s %s", feedbackStyle.Render(feedback), answerStyle.Render(card.Answer), extraDisplay, mcqChoices,
					keyStyle.Render("a"), gradeAgain, keyStyle.Render("h"), gradeHard, keyStyle.Render("g"), gradeGood, keyStyle.Render("e"), gradeEasy)
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n%s: %s\n\n%s\n\nGrade: ", promptDisplay, mature, feedback, card.Answer, mcqChoices), "\n")

				labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
				labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"))
				labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"))
				labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"), gradeGoodText, keyStyle.Render("e"))

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			} else {
				answer = fmt.Sprintf("1-4 select answer%s\n\n%s\n\nGrade: %s %s | %s %s | %s %s | %s %s", extraDisplay, mcqChoices,
					keyStyle.Render("a"), gradeAgain, keyStyle.Render("h"), gradeHard, keyStyle.Render("g"), gradeGood, keyStyle.Render("e"), gradeEasy)
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n1-4 select answer\n\n%s\n\nGrade: ", promptDisplay, mature, mcqChoices), "\n")

				labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
				labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"))
				labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"))
				labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"), gradeGoodText, keyStyle.Render("e"))

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			}
		} else if m.revealState == RevealRevealing {
			answer = fmt.Sprintf("1-4 select answer\n\n%s", renderMCQChoices(card.Choices, m.mcqChoice))
		} else {
			answer = "Press Space or Enter to reveal choices."
		}
	} else if m.revealState == RevealRevealed {
		answer = fmt.Sprintf("%s%s\n\nGrade: %s %s | %s %s | %s %s | %s %s", answerStyle.Render(card.Answer), extraDisplay,
			keyStyle.Render("a"), gradeAgain, keyStyle.Render("h"), gradeHard, keyStyle.Render("g"), gradeGood, keyStyle.Render("e"), gradeEasy)
		answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n\n%s\n\nGrade: ", promptDisplay, mature, card.Answer), "\n")

		labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
		labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"))
		labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"))
		labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"), gradeGoodText, keyStyle.Render("e"))

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
	choiceStyle := lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	for i, choice := range choices {
		prefix := "  "
		mark := "○ "
		if i == selected {
			mark = selectedStyle.Render("● ")
		} else {
			mark = choiceStyle.Foreground(lipgloss.Color("240")).Render(mark)
		}
		b.WriteString(fmt.Sprintf("%s%d: %s%s\n", prefix, i+1, mark, choice))
	}
	return strings.TrimRight(b.String(), "\n")
}
