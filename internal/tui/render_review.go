package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

// renderGrammarHint builds a compact 2-4 line "how is this word used?" block
// shown under the revealed answer. It analyses the card's German side and,
// based on the detected word kind, surfaces gender + case shape (nouns),
// conjugations (verbs), or comparison forms (adjectives), plus a curated or
// auto-generated example sentence.
func renderGrammarHint(card core.Card) string {
	info := content.AnalyzeCard(card.Prompt, card.Answer)
	if info.Kind == content.KindUnknown || info.Display == "" {
		return ""
	}

	kindColor := lipgloss.Color("99")
	switch info.Kind {
	case content.KindNoun:
		switch info.Gender {
		case "masculine":
			kindColor = lipgloss.Color("39") // blue der
		case "feminine":
			kindColor = lipgloss.Color("197") // red die
		case "neuter":
			kindColor = lipgloss.Color("46") // green das
		}
	case content.KindVerb:
		kindColor = lipgloss.Color("214")
	case content.KindAdjective:
		kindColor = lipgloss.Color("117")
	case content.KindPhrase:
		kindColor = lipgloss.Color("141")
	}

	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")).
		Background(kindColor).
		Bold(true).
		Padding(0, 1).
		Render(info.Kind.String())

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	formStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Italic(true)
	exampleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Italic(true)
	noteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)

	var lines []string
	header := badge + " " + lipgloss.NewStyle().Bold(true).Render(info.Display)
	if info.Gender != "" {
		header += "  " + labelStyle.Render("("+info.Gender+")")
	}
	lines = append(lines, header)

	if len(info.Forms) > 0 {
		lines = append(lines, labelStyle.Render("Forms:")+"   "+formStyle.Render(strings.Join(info.Forms, " · ")))
	}
	if info.Example != "" {
		lines = append(lines, labelStyle.Render("Example:")+" "+exampleStyle.Render(info.Example))
	}
	if info.Note != "" {
		lines = append(lines, noteStyle.Render(info.Note))
	}

	return strings.Join(lines, "\n")
}

func (m *Model) renderReview(x, y int) string {
	if len(m.dueCards) == 0 {
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Render("Review")
		message := "No cards due."
		if m.bookmarkFilter {
			title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Render("Review (Bookmarked)")
			message = "No bookmarked cards due."
		}

		deckLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
		deckLine := fmt.Sprintf("Current Deck: %s", deckLabelStyle.Render(m.deckLabel()))

		// Enhanced empty state with keyboard guidance
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		shortcutGuide := fmt.Sprintf("Use %s and %s to switch decks or %s to toggle the bookmark filter.\nPress %s for Custom Study (Cram Mode).",
			keyStyle.Render("["), keyStyle.Render("]"), keyStyle.Render("B"), keyStyle.Render("c"))

		// Rotating motivational tips based on day
		tips := []string{
			"Tip: Consistent daily practice beats cramming!",
			"Tip: Try reviewing cards out loud for better retention.",
			"Tip: Use new words in sentences to strengthen memory.",
			"Tip: Review before sleep helps consolidate learning.",
			"Tip: Focus on understanding, not just memorizing.",
		}
		dayTip := tips[int(time.Now().Weekday())%len(tips)]
		tipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Italic(true)

		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2).
			Width(60).
			Align(lipgloss.Center).
			Render(deckLine + "\n" + message + "\n\n" + shortcutGuide + "\n\n" + tipStyle.Render(dayTip))
		return title + "\n\n" + emptyBox
	}
	cursor := clampInt(m.cursor, 0, len(m.dueCards)-1)
	card := m.dueCards[cursor]
	leechBadge := ""
	if card.Leech {
		leechBadge = " " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("196")).
			Bold(true).
			Padding(0, 1).
			Render("LEECH")
	}
	filterBanner := ""
	if m.bookmarkFilter {
		filterBanner = " (Bookmarked)"
	}

	// Enhanced keyboard shortcut display with visual highlighting
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
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
	if m.cardTransitioning {
		spinFrames := []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
		spin := spinFrames[m.cardTransitionFrame%len(spinFrames)]
		dirStr := "↑"
		if m.cardTransitionDir > 0 {
			dirStr = "↓"
		}
		transStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
		header += " " + transStyle.Render(fmt.Sprintf("%s %s", spin, dirStr))
	}

	// Session progress bar and timer
	sessionTotal := m.sessionReviewed + len(m.dueCards)
	sessionPercentage := 0.0
	if sessionTotal > 0 {
		sessionPercentage = float64(m.sessionReviewed) / float64(sessionTotal)
	}
	sessionBar := progressBar(15, sessionPercentage, lipgloss.Color("46"), lipgloss.Color("238"))

	duration := time.Since(m.sessionStartTime)
	if m.testMode {
		duration = 0 // Freeze timer in test mode
	}
	timerStr := fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)

	cardsPerMin := 0.0
	if duration.Minutes() > 0 {
		cardsPerMin = float64(m.sessionReviewed) / duration.Minutes()
	}

	etaStr := ""
	if cardsPerMin > 0 && len(m.dueCards) > 0 {
		remaining := len(m.dueCards)
		etaMin := float64(remaining) / cardsPerMin
		if etaMin < 1 {
			etaStr = fmt.Sprintf(" | ETA <1m")
		} else if etaMin < 60 {
			etaStr = fmt.Sprintf(" | ETA ~%dm", int(etaMin))
		} else {
			etaStr = fmt.Sprintf(" | ETA ~%dh%dm", int(etaMin)/60, int(etaMin)%60)
		}
	}

	sessionProgress := fmt.Sprintf(" | Session: %s %d%% | ⏱ %s | %.1f/min%s", sessionBar, int(sessionPercentage*100), timerStr, cardsPerMin, etaStr)

	mature := ""
	if card.Mature {
		mature = " ✨"
	}

	// Determine card state/difficulty for badge
	stateBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Padding(0, 1).Bold(true)
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

	var headerSection string
	if m.focusMode {
		headerSection = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render("Focus Mode Active (f to exit)")
	} else {
		bookmarkText := "Bookmark: off"
		if card.Bookmarked {
			bookmarkText = "Bookmark: on"
		}

		suspendText := "Suspend"
		var suspendStyle lipgloss.Style
		if card.Suspended {
			suspendText = "SUSPENDED"
			suspendStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		} else {
			suspendStyle = mutedStyle
		}

		audioText := ""
		var audioStyle lipgloss.Style
		if card.Audio != "" {
			audioText = "[Audio]"
			audioStyle = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
		} else if m.ttsAvailable() {
			audioText = "[TTS]"
			audioStyle = lipgloss.NewStyle().Foreground(colorCyan)
		}

		// Register hitboxes for Line 3 (Y = y + 2)
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "review-bookmark",
			View:   ViewReview,
			X:      x,
			Y:      y + 2,
			Width:  len(bookmarkText),
			Height: 1,
		})

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "review-suspend",
			View:   ViewReview,
			X:      x + len(bookmarkText) + 3,
			Y:      y + 2,
			Width:  len(suspendText),
			Height: 1,
		})

		if audioText != "" {
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     "review-audio",
				View:   ViewReview,
				X:      x + len(bookmarkText) + 3 + len(suspendText) + 3,
				Y:      y + 2,
				Width:  len(audioText),
				Height: 1,
			})
		}

		metaLine := fmt.Sprintf("%s | %s", bookmarkText, suspendStyle.Render(suspendText))
		if audioText != "" {
			metaLine += fmt.Sprintf(" | %s", audioStyle.Render(audioText))
		}

		headerSection = fmt.Sprintf("%s%s %s%s\n%s | Type: %s\n%s", header, sessionProgress, stateBadge.Render(), leechBadge, deckMeta, cardTypeLabel, metaLine)
	}
	headerLines := strings.Count(headerSection, "\n") + 1

	promptDisplay := card.Prompt
	if card.Kind == core.CardKindCloze {
		// Highlight the entire [...] or [hint] block
		bracketStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("81")).
			Bold(true)

		re := regexp.MustCompile(`\[[^\]]+\]`)
		promptDisplay = re.ReplaceAllStringFunc(promptDisplay, func(s string) string {
			return bracketStyle.Render(s)
		})
	}

	hintDisplay := ""
	if m.showHint {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Italic(true)
		val := card.Hint
		if val == "" {
			val = "(no hint available)"
		}
		hintDisplay = "\n" + hintStyle.Render("Hint: "+val)
	}

	// Enhanced prompt styling for better visibility
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	promptDisplay = promptStyle.Render(promptDisplay + hintDisplay)

	width, height := m.activePanelSize()
	cardWidth := maxInt(30, width-6)
	if cardWidth > 80 {
		cardWidth = 80
	}

	// Enhanced card styling with better visual hierarchy
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cardBorderColor)).
		Width(cardWidth)

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
	if m.revealState == RevealRevealed {
		if hint := renderGrammarHint(card); hint != "" {
			extraDisplay = "\n\n" + hint + extraDisplay
		}
	}

	if m.revealState == RevealFlipping {
		answer = renderFlipAnimation(m, card, cardWidth)
	} else if m.revealState == RevealRevealing {
		// Show gradual reveal animation with blocks for all card types
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
	} else if m.typingMode {
		typingBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1).
			Width(maxInt(30, width-6))

		inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		inputDisplay := m.typedAnswer + "_"

		if m.typingChecked {
			targetAnswer := card.Answer
			if card.Kind == core.CardKindCloze && len(card.Choices) > 0 {
				targetAnswer = card.Choices[0]
			}
			if m.typingCorrect {
				correctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
				inputDisplay = correctStyle.Render("✓ ") + renderTypingDiff(m.typedAnswer, targetAnswer)
				typingBoxStyle = typingBoxStyle.BorderForeground(lipgloss.Color("46"))
			} else {
				wrongStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
				inputDisplay = wrongStyle.Render("✗ ") + renderTypingDiff(m.typedAnswer, targetAnswer)
				typingBoxStyle = typingBoxStyle.BorderForeground(lipgloss.Color("196"))
			}
			typingContent := fmt.Sprintf("Your answer: %s\nCorrect: %s%s\n\nGrade: %s %s | %s %s | %s %s | %s %s",
				inputDisplay, answerStyle.Render(targetAnswer), extraDisplay,
				keyStyle.Render("a"), gradeAgain+" "+keyStyle.Render("(1)"), keyStyle.Render("h"), gradeHard+" "+keyStyle.Render("(2)"), keyStyle.Render("g"), gradeGood+" "+keyStyle.Render("(3)"), keyStyle.Render("e"), gradeEasy+" "+keyStyle.Render("(4)"))
			answer = typingBoxStyle.Render(typingContent)
			answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\nYour answer: %s\nCorrect: %s%s\n\nGrade: ", promptDisplay, mature, inputDisplay, targetAnswer, extraDisplay), "\n")

			labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
			labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"))
			labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"))
			labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"), gradeGood+" (3)", keyStyle.Render("e"))

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
					keyStyle.Render("a"), gradeAgain+" "+keyStyle.Render("(1)"), keyStyle.Render("h"), gradeHard+" "+keyStyle.Render("(2)"), keyStyle.Render("g"), gradeGood+" "+keyStyle.Render("(3)"), keyStyle.Render("e"), gradeEasy+" "+keyStyle.Render("(4)"))
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n%s: %s%s\n\n%s\n\nGrade: ", promptDisplay, mature, feedback, card.Answer, extraDisplay, mcqChoices), "\n")

				labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
				labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"))
				labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"))
				labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"), gradeGood+" (3)", keyStyle.Render("e"))

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			} else {
				answer = fmt.Sprintf("1-4 select answer%s\n\n%s\n\nGrade: %s %s | %s %s | %s %s | %s %s", extraDisplay, mcqChoices,
					keyStyle.Render("a"), gradeAgain, keyStyle.Render("h"), gradeHard, keyStyle.Render("g"), gradeGood, keyStyle.Render("e"), gradeEasy)
				answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n1-4 select answer%s\n\n%s\n\nGrade: ", promptDisplay, mature, extraDisplay, mcqChoices), "\n")

				labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
				labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"))
				labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"))
				labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s %s", keyStyle.Render("a"), gradeAgainText, keyStyle.Render("h"), gradeHardText, keyStyle.Render("g"), gradeGoodText, keyStyle.Render("e"), gradeEasyText)

				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
				m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
			}
		} else {
			answer = "Press Space or Enter to reveal choices."
		}
	} else if m.revealState == RevealRevealed {
		answerDisplay := answerStyle.Render(card.Answer)
		if card.Kind == core.CardKindCloze && len(card.Choices) > 0 {
			clozeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("46")).
				Bold(true)

			re := regexp.MustCompile(`\[[^\]]+\]`)
			answerDisplay = re.ReplaceAllStringFunc(card.Prompt, func(s string) string {
				return clozeStyle.Render(card.Choices[0])
			})
			answerDisplay = answerStyle.Render(answerDisplay)
		}

		answer = fmt.Sprintf("%s%s\n\nGrade: %s %s | %s %s | %s %s | %s %s", answerDisplay, extraDisplay,
			keyStyle.Render("a"), gradeAgain+" "+keyStyle.Render("(1)"), keyStyle.Render("h"), gradeHard+" "+keyStyle.Render("(2)"), keyStyle.Render("g"), gradeGood+" "+keyStyle.Render("(3)"), keyStyle.Render("e"), gradeEasy+" "+keyStyle.Render("(4)"))
		answerYOffset = cardY + strings.Count(fmt.Sprintf("%s%s\n%s%s\n\nGrade: ", promptDisplay, mature, answerDisplay, extraDisplay), "\n")

		labelAgain := fmt.Sprintf("Grade: %s ", keyStyle.Render("a"))
		labelHard := fmt.Sprintf("Grade: %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"))
		labelGood := fmt.Sprintf("Grade: %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"))
		labelEasy := fmt.Sprintf("Grade: %s %s | %s %s | %s %s | %s ", keyStyle.Render("a"), gradeAgain+" (1)", keyStyle.Render("h"), gradeHard+" (2)", keyStyle.Render("g"), gradeGood+" (3)", keyStyle.Render("e"))

		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-again", View: ViewReview, X: cardX + lipgloss.Width(labelAgain), Y: answerYOffset, Width: gaW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-hard", View: ViewReview, X: cardX + lipgloss.Width(labelHard), Y: answerYOffset, Width: ghW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-good", View: ViewReview, X: cardX + lipgloss.Width(labelGood), Y: answerYOffset, Width: ggW, Height: 1})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "grade-easy", View: ViewReview, X: cardX + lipgloss.Width(labelEasy), Y: answerYOffset, Width: geW, Height: 1})
	} else {
		answer = "Press space or enter to reveal."
	}

	view := fmt.Sprintf("%s\n\n%s", headerSection, cardStyle.Render(promptDisplay+mature+"\n"+answer))
	if m.focusMode {
		// Use full width/height in focus mode
		content := promptDisplay + mature + "\n" + answer
		focusHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render("Focus Mode Active (f to exit)")
		view = cardStyle.
			Width(width-4).
			Height(height-4).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content)
		centered := lipgloss.Place(m.width, m.height-10, lipgloss.Center, lipgloss.Center, view)
		return focusHeader + "\n" + centered
	}

	if m.showCardInfo {
		view += "\n\n" + m.renderCardInfo(card)
	}
	if m.showReviewHistory && m.reviewHistoryCard == card.ID {
		view += "\n\n" + m.renderReviewHistory(card.Prompt)
	}
	if m.fixingCard {
		spinFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spin := spinFrames[m.spinnerFrame%len(spinFrames)]
		view += "\n\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorYellow).
			Padding(0, 1).
			Render(infoStyle.Bold(true).Render("Reporting card to AI…")+"  "+spin)
	}
	if m.explainingCard {
		spinFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spin := spinFrames[m.spinnerFrame%len(spinFrames)]
		view += "\n\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1).
			Render(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("Asking AI tutor for explanation…")+"  "+spin)
	}
	if m.showGrammarHint && m.grammarHint != nil {
		hintBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("212")). // Pink border
			Padding(0, 1).
			Width(cardWidth)

		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render("Grammar Hint: " + m.grammarHint.Title)
		content := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(m.grammarHint.Tip)
		example := lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Italic(true).Render(m.grammarHint.Example)

		view += "\n\n" + hintBoxStyle.Render(title+"\n\n"+content+"\n\n"+example)
	}
	if m.explanation != "" {
		explanationStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1).
			Width(cardWidth)

		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")).Render("AI Tutor Explanation")
		view += "\n\n" + explanationStyle.Render(title+"\n\n"+m.explanation)
	}
	if m.fixProposal != nil && m.fixOldNote != nil {
		view += "\n\n" + m.renderFixProposal()
	}

	// Add Interactive Footer
	var hintAction string
	if m.revealState == RevealRevealed {
		hintAction = fmt.Sprintf("%s/%s/%s/%s grade", keyStyle.Render("a"), keyStyle.Render("h"), keyStyle.Render("g"), keyStyle.Render("e"))
	} else {
		hintAction = fmt.Sprintf("%s hint", keyStyle.Render("h"))
	}

	footerActions := fmt.Sprintf("%s | %s grammar | %s explain | %s bookmark | %s suspend | %s undo | %s history | %s info | %s focus | %s audio",
		hintAction, keyStyle.Render("G"), keyStyle.Render("H"), keyStyle.Render("b"), keyStyle.Render("x"),
		keyStyle.Render("u/z"), keyStyle.Render("r"), keyStyle.Render("i"),
		keyStyle.Render("f"), keyStyle.Render("p"))
	if m.bookmarkFilter {
		footerActions += " | " + keyStyle.Render("B") + " all"
	} else {
		footerActions += " | " + keyStyle.Render("B") + " bookmarked"
	}

	view += "\n\n" + lipgloss.NewStyle().
		Foreground(colorMuted).
		Render(footerActions)

	return view
}

// renderFixProposal shows the AI's proposed correction next to the
// current values so the user can decide whether to apply it.
func (m *Model) renderFixProposal() string {
	old := *m.fixOldNote
	fix := *m.fixProposal

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
	labelStyle := lipgloss.NewStyle().Foreground(colorMuted)
	oldStyle := lipgloss.NewStyle().Foreground(colorOrange)
	newStyle := lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	reasonStyle := lipgloss.NewStyle().Foreground(colorCyan).Italic(true)

	row := func(label, oldVal, newVal string) string {
		if strings.TrimSpace(oldVal) == "" && strings.TrimSpace(newVal) == "" {
			return ""
		}
		changed := strings.TrimSpace(oldVal) != strings.TrimSpace(newVal)
		newRendered := newStyle.Render(newVal)
		if !changed {
			newRendered = labelStyle.Render(newVal + "  (unchanged)")
		}
		return labelStyle.Render(label) + "\n" +
			"  " + labelStyle.Render("old: ") + oldStyle.Render(oldVal) + "\n" +
			"  " + labelStyle.Render("new: ") + newRendered + "\n"
	}

	oldExample := ""
	if len(old.Examples) > 0 {
		oldExample = old.Examples[0]
	}

	body := titleStyle.Render("AI-proposed fix") + "\n"
	if fix.Reason != "" {
		body += reasonStyle.Render("Reason: "+fix.Reason) + "\n\n"
	}
	body += row("Front:", old.Front, fix.Front)
	body += row("Back:", old.Back, fix.Back)
	body += row("Extra:", old.Extra, fix.Extra)
	body += row("Example:", oldExample, fix.Example)
	body += "\n" + labelStyle.Render("Press ") + keyStyle.Render("y") +
		labelStyle.Render(" to apply, ") + keyStyle.Render("n") +
		labelStyle.Render(" to discard.")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorYellow).
		Padding(0, 1).
		Render(body)
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

func renderFlipAnimation(m *Model, card core.Card, cardWidth int) string {
	progress := int(m.flipProgress)
	frame := m.flipFrame

	frontStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("159"))

	backStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	spinFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spin := spinFrames[frame%len(spinFrames)]

	if progress < 50 {
		shrinkProgress := progress * 2
		visibleWidth := maxInt(1, (cardWidth-4)*(100-shrinkProgress)/100)
		frontContent := frontStyle.Render(card.Prompt)
		if visibleWidth < lipgloss.Width(frontContent) {
			runes := []rune(card.Prompt)
			visibleRunes := runes[:maxInt(1, len(runes)*(100-shrinkProgress)/100)]
			frontContent = frontStyle.Render(string(visibleRunes))
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Width(cardWidth).
			Align(lipgloss.Center).
			Render(frontContent + " " + spin)
	} else {
		expandProgress := (progress - 50) * 2
		fullText := card.Answer
		runes := []rune(fullText)
		visibleRunes := runes[:maxInt(1, len(runes)*expandProgress/100)]
		backContent := backStyle.Render(string(visibleRunes))
		if expandProgress < 100 {
			remaining := len(runes) - len(visibleRunes)
			backContent += " " + strings.Repeat("▌", maxInt(1, remaining))
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Width(cardWidth).
			Align(lipgloss.Center).
			Render(backContent + " " + spin)
	}
}

func (m *Model) renderCardInfo(card core.Card) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))

	info := fmt.Sprintf("%s\n\n", titleStyle.Render("Card Statistics"))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("ID:"), valueStyle.Render(card.ID))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Deck:"), valueStyle.Render(m.deckNameByID(card.DeckID)))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Ease:"), valueStyle.Render(fmt.Sprintf("%.0f%%", card.Ease*100)))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Interval:"), valueStyle.Render(formatReviewInterval(card.Interval)))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Reviews:"), valueStyle.Render(fmt.Sprintf("%d", card.Reviews)))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Lapses:"), valueStyle.Render(fmt.Sprintf("%d", card.Lapses)))
	info += fmt.Sprintf("%s %s\n", labelStyle.Render("Due:"), valueStyle.Render(card.Due.Local().Format("Jan 02, 2006")))

	if len(card.Tags) > 0 {
		info += fmt.Sprintf("%s #%s\n", labelStyle.Render("Tags:"), valueStyle.Render(strings.Join(card.Tags, " #")))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("159")).
		Padding(0, 1).
		Render(info)
}
