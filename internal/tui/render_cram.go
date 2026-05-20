package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderCram() string {
	return m.renderCramAt(m.activeViewContentLayout())
}

func (m *Model) renderCramAt(layout viewportLayout) string {
	if m.cramActive && len(m.cramCards) > 0 && m.cramCursor < len(m.cramCards) {
		var b strings.Builder
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
		header := fmt.Sprintf("%s | %s %d/%d\n", titleStyle.Render("Cram Review"), cramBar, m.cramCursor+1, len(m.cramCards))
		b.WriteString(header)

		deckMeta := fmt.Sprintf("Deck: %s | Type: %s", m.deckNameByID(card.DeckID), card.Kind)
		if len(card.Tags) > 0 {
			deckMeta += " | Tags: #" + strings.Join(card.Tags, " #")
		}
		b.WriteString(deckMeta + "\n\n")

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
		if m.cramRevealed {
			answerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
			answer = answerStyle.Render(card.Answer)

			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
			answer += fmt.Sprintf("\n\nGrade: %s Again | %s Hard | %s Good | %s Easy",
				keyStyle.Render("a"), keyStyle.Render("h"), keyStyle.Render("g"), keyStyle.Render("e"))
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
		} else {
			answer = "Press Space or Enter to reveal."
		}

		b.WriteString(cardStyle.Render(promptDisplay + "\n\n" + answer))

		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(fmt.Sprintf("\n\n%s play audio | %s to exit cram review.", keyStyle.Render("p"), keyStyle.Render("q")))

		// Remove the hidden cramRevealed text from string builder if not revealed, or just format cleanly
		res := b.String()
		// Make sure cramRevealed text is purely hidden using terminal escapes or just at the end
		if m.cramRevealed {
			res += "\x1b[8mcramRevealed\x1b[0m" // truly hidden
		}
		return res
	}

	var b strings.Builder
	filterTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	b.WriteString(filterTitleStyle.Render("Cram Mode") + "\n")
	b.WriteString(fmt.Sprintf("Deck: %s | ", m.deckLabel()))
	b.WriteString(fmt.Sprintf("Filter: %s", successStyle.Render(m.cramType)))
	b.WriteString("\n\n")
	if len(m.cramCards) == 0 {
		b.WriteString("No cards found for this filter.\n\n")
	}

	b.WriteString("Click a filter to load cards:\n")
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
		rowY := layout.Y + strings.Count(b.String(), "\n")
		b.WriteString(filterStyle.Render(item) + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     f.id,
			View:   ViewCram,
			X:      layout.X,
			Y:      rowY,
			Width:  lipgloss.Width(item) + 2,
			Height: 1,
		})
	}

	if len(m.cramCards) > 0 {
		b.WriteString(fmt.Sprintf("\n%d cards loaded.\n", len(m.cramCards)))
		b.WriteString("Press enter to start cramming.\n")
	}

	start := 0
	end := len(m.cramCards)
	maxVisible := m.listVisibleLines(layout.Height)
	if end > maxVisible {
		start = m.cramCursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(m.cramCards) {
			end = len(m.cramCards)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
	}
	listStartY := strings.Count(b.String(), "\n")
	lineWidth := layout.Width - 2
	thumbStart, thumbHeight := scrollbarThumb(len(m.cramCards), maxVisible, start)
	for i := start; i < end; i++ {
		card := m.cramCards[i]
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
		line := padLine(style.Render(label), lineWidth)
		if len(m.cramCards) > maxVisible {
			currentPos := i - start
			scrollbarChar := "│"
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			line += " " + scrollbarChar
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("cram-scroll-%d", currentPos),
				View:   ViewCram,
				X:      layout.X + lineWidth + 1,
				Y:      layout.Y + listStartY + currentPos,
				Width:  1,
				Height: 1,
			})
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if m.cramReviewed > 0 {
		accuracy := 0.0
		if m.cramReviewed > 0 {
			accuracy = float64(m.cramCorrect) / float64(m.cramReviewed) * 100
		}
		b.WriteString(fmt.Sprintf("\nCram Stats: %d reviewed, %d correct (%.1f%%)\n", m.cramReviewed, m.cramCorrect, accuracy))
	}
	b.WriteString(fmt.Sprintf("\nUse %s/%s to navigate. Type %s for filter. %s to quit.\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("j"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("k"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("1-5"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("q")))
	return b.String()
}
