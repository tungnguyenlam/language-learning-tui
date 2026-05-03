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
	if m.cramActive {
		var b strings.Builder
		card := m.cramCards[m.cramCursor]
		audioIndicator := ""
		if card.Audio != "" {
			audioIndicator = " [Audio]"
		}
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
		b.WriteString(titleStyle.Render("Cram Review") + "\n\n")
		b.WriteString(fmt.Sprintf("Prompt: %s%s\n\n", card.Prompt, audioIndicator))
		if m.cramRevealed {
			b.WriteString(fmt.Sprintf("Answer: %s\n\n", card.Answer))
			b.WriteString("Grade: a Again | h Hard | g Good | e Easy\n")
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
			b.WriteString(fmt.Sprintf("Answer: %s\n\n", animationText))
			b.WriteString("Press space or enter to finish reveal.\n")
		} else {
			b.WriteString("Press space or enter to reveal.\n")
		}
		b.WriteString("\np play audio | q to exit cram review.")
		return b.String()
	}

	var b strings.Builder
	b.WriteString("Cram Mode\n\n")
	b.WriteString(fmt.Sprintf("Filter: %s\n\n", m.cramType))
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
	lineWidth := scrollbarLineWidth(layout.Width)
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
		label := fmt.Sprintf("%s[%s] %s%s%s%s%s", prefix, kind, card.Prompt, mature, bookmark, leech, suspended)
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
	b.WriteString("\nUse j/k to navigate. Type 1-5 for filter. q to quit.\n")
	return b.String()
}
