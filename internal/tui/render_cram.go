package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderCram() string {
	return m.renderCramAt(m.activeViewContentLayout())
}

func (m *Model) renderCramAt(layout viewportLayout) string {
	ctx := NewRenderContext(m, layout, ViewCram)

	if m.cramActive && len(m.cramCards) > 0 && m.cramCursor < len(m.cramCards) {
		card := m.cramCards[m.cramCursor]
		audioIndicator := ""
		if card.Audio != "" {
			audioIndicator = " [Audio]"
		} else if m.ttsAvailable() {
			audioIndicator = " [TTS]"
		}

		// Session progress bar
		cramPercentage := 0.0
		if len(m.cramCards) > 0 {
			cramPercentage = float64(m.cramCursor) / float64(len(m.cramCards))
		}
		cramBar := progressBar(15, cramPercentage, "81", "238")

		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
		header := fmt.Sprintf("%s | %s %d/%d", titleStyle.Render("Cram Review"), cramBar, m.cramCursor+1, len(m.cramCards))
		ctx.WriteLine(header)

		deckMeta := fmt.Sprintf("Deck: %s | Type: %s", m.deckNameByID(card.DeckID), card.Kind)
		if len(card.Tags) > 0 {
			deckMeta += " | Tags: #" + strings.Join(card.Tags, " #")
		}
		ctx.WriteLine(deckMeta)
		ctx.NewLine()

		width := layout.Width
		cardWidth := maxInt(30, width-6)

		cardBorderColor := "62"
		if m.cramRevealed {
			cardBorderColor = "81"
		}

		cardStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(cardBorderColor)).
			Padding(1, 2).
			Width(cardWidth).
			Align(lipgloss.Center)

		promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
		promptDisplay := promptStyle.Render(card.Prompt) + audioIndicator

		var answer string
		if m.revealState == RevealFlipping {
			answer = renderFlipAnimation(m, card, cardWidth)
		} else if m.revealState == RevealRevealing {
			// Show gradual reveal animation with blocks
			progress := int(m.revealProgress)
			if progress > 100 {
				progress = 100
			}
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
			animationText := string(revealedRunes)
			if remainingBlocks > 0 {
				animationText += strings.Repeat("▌", remainingBlocks)
			}
			answer = animationText + "\n\nPress Space or Enter to finish reveal."
		} else if m.cramRevealed {
			answerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)

			extraDisplay := ""
			if card.Extra != "" {
				extraStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Italic(true)
				extraDisplay = "\n\n" + extraStyle.Render("💡 CONTEXT: "+card.Extra)
			}
			if hint := renderGrammarHint(card); hint != "" {
				extraDisplay = "\n\n" + hint + extraDisplay
			}

			answer = answerStyle.Render(card.Answer) + extraDisplay

			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
			answer += fmt.Sprintf("\n\nGrade: %s Again | %s Hard | %s Good | %s Easy",
				keyStyle.Render("a"), keyStyle.Render("h"), keyStyle.Render("g"), keyStyle.Render("e"))
		} else {
			answer = "Press Space or Enter to reveal."
		}

		ctx.WriteLine(cardStyle.Render(promptDisplay + "\n\n" + answer))

		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		ctx.WriteLine(fmt.Sprintf("\n%s play audio | %s to exit cram review.", keyStyle.Render("p"), keyStyle.Render("q")))

		return ctx.String()
	}

	filterTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	ctx.WriteLine(filterTitleStyle.Render("Cram Mode"))
	ctx.WriteLine(fmt.Sprintf("Deck: %s | Filter: %s", m.deckLabel(), successStyle.Render(m.cramType)))
	ctx.NewLine()
	if len(m.cramCards) == 0 {
		ctx.WriteLine("No cards found for this filter.\n")
	}

	ctx.WriteLine("Click a filter to load cards:")
	filters := []struct {
		id    string
		label string
	}{
		{"cram-filter-1", "Bookmarked"},
		{"cram-filter-2", "Suspended"},
		{"cram-filter-3", "Leeches"},
		{"cram-filter-4", "All flagged"},
		{"cram-filter-5", "All cards"},
	}

	filterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("62")).
		Padding(0, 1).
		MarginBottom(1)

	for i, f := range filters {
		idx := i + 1
		item := fmt.Sprintf("  %d: %s", idx, f.label)
		ctx.RegisterHitboxWithAction(f.id, lipgloss.Width(item)+2, 1, func() tea.Cmd {
			return m.setCramFilter(idx)
		})
		ctx.WriteLine(filterStyle.Render(item))
	}

	if len(m.cramCards) > 0 {
		ctx.WriteLine(fmt.Sprintf("\n%d cards loaded.", len(m.cramCards)))
		ctx.WriteLine("Press enter to start cramming.")
	}

	var content strings.Builder
	lineWidth := layout.Width - 2
	for i, card := range m.cramCards {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.cramCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		kind := "FC"
		if card.Kind == core.CardKindMCQ {
			kind = "MCQ"
		}
		bookmark := ""
		if card.Bookmarked {
			bookmark = " [B]"
		}
		leech := ""
		if card.Leech {
			leech = " [L]"
		}
		suspended := ""
		if card.Suspended {
			suspended = " [S]"
		}
		mature := ""
		if card.Mature {
			mature = " ⭐"
		}
		plainDeckName := ""
		if m.deck.ID == "" {
			plainDeckName = " (" + truncateLine(m.deckNameByID(card.DeckID), 18) + ")"
		}

		otherWidth := lipgloss.Width(prefix + "[" + kind + "] " + plainDeckName + mature + bookmark + leech + suspended)
		availWidth := lineWidth - otherWidth
		if availWidth < 10 {
			availWidth = 10
		}

		truncatedPrompt := truncateLine(card.Prompt, availWidth)

		styledDeckName := ""
		if plainDeckName != "" {
			styledDeckName = " " + mutedStyle.Render(strings.TrimSpace(plainDeckName))
		}

		label := fmt.Sprintf("%s[%s] %s%s%s%s%s%s", prefix, kind, truncatedPrompt, styledDeckName, mature, bookmark, leech, suspended)
		content.WriteString(style.Render(label) + "\n")
	}

	contentStr := content.String()
	totalLines := len(m.cramCards)

	availableHeight := m.listVisibleLines(layout.Height) - (ctx.currY - layout.Y)
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Auto-scroll
	m.cramScroll = AutoScroll(m.cramCursor, m.cramScroll, availableHeight, totalLines)

	listView := m.RenderList(layout.WithHeight(availableHeight).WithY(ctx.currY), contentStr, ListOptions{
		HitboxPrefix: "cram",
		View:         ViewCram,
		ScrollOffset: &m.cramScroll,
	})

	ctx.Write(listView)

	if m.cramReviewed > 0 {
		accuracy := 0.0
		if m.cramReviewed > 0 {
			accuracy = float64(m.cramCorrect) / float64(m.cramReviewed) * 100
		}
		ctx.WriteLine(fmt.Sprintf("\nCram Stats: %d reviewed, %d correct (%.1f%%)", m.cramReviewed, m.cramCorrect, accuracy))
	}
	ctx.WriteLine(fmt.Sprintf("\nUse %s/%s to navigate. Type %s for filter. %s to quit.",
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("j"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("k"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("1-5"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("q")))

	return ctx.String()
}
