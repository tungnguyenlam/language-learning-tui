package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderAI(x, y int) string {
	var b strings.Builder
	spinner := ""
	if m.drafting {
		frames := []string{"-", "\\", "|", "/"}
		spinner = " " + frames[m.spinnerFrame%len(frames)]
	}
	templateName := "None"
	if len(m.aiTemplateSets) > 0 {
		templateName = m.aiTemplateSets[m.aiTemplateIndex]
	}

	fmt.Fprintf(&b, "AI Drafts%s\n\nDeck: %s\nTemplate: %s (use [ / ])\nTopic: %s\n\nEnter generate | a approve | d discard | esc clear\n", spinner, m.deckLabel(), templateName, m.aiInput)

	if len(m.drafts) == 0 {
		if m.drafting {
			b.WriteString("\nDrafting in progress...")
		} else {
			b.WriteString("\nNo drafts yet.")
		}
		return b.String()
	}

	start := 0
	end := len(m.drafts)
	maxVisible := 10
	if m.height > 20 {
		maxVisible = m.height - 15
	}
	if maxVisible < 5 {
		maxVisible = 5
	}
	if end > maxVisible {
		start = m.draftCursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(m.drafts) {
			end = len(m.drafts)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.draftCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}
		draft := m.drafts[i]
		item := fmt.Sprintf("%s%s -> %s", prefix, draft.Note.Front, draft.Note.Back)
		b.WriteString(style.Render(item))

		// Add interactive buttons
		btnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		approveBtn := " [Approve]"
		discardBtn := " [Discard]"

		if i == m.draftCursor {
			btnStyle = btnStyle.Foreground(lipgloss.Color("81"))
		}

		b.WriteString(btnStyle.Render(approveBtn))
		b.WriteString(btnStyle.Render(discardBtn))
		b.WriteString("\n")

		// Register hitboxes
		rowY := y + 6 + (i - start)
		approveX := x + lipgloss.Width(item)
		discardX := approveX + lipgloss.Width(approveBtn)

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("draft-approve-%d", i),
			View:   ViewAI,
			X:      approveX,
			Y:      rowY,
			Width:  lipgloss.Width(approveBtn),
			Height: 1,
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("draft-discard-%d", i),
			View:   ViewAI,
			X:      discardX,
			Y:      rowY,
			Width:  lipgloss.Width(discardBtn),
			Height: 1,
		})
	}

	// Show detailed preview for selected draft
	if len(m.drafts) > 0 && m.draftCursor < len(m.drafts) {
		selected := m.drafts[m.draftCursor]
		previewTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
		b.WriteString("\n" + previewTitleStyle.Render("Preview:") + "\n")
		b.WriteString(fmt.Sprintf("  Extra:    %s\n", selected.Note.Extra))
		if len(selected.Note.Tags) > 0 {
			b.WriteString(fmt.Sprintf("  Tags:     %s\n", strings.Join(selected.Note.Tags, ", ")))
		}
		if len(selected.Note.Examples) > 0 {
			b.WriteString("  Examples:\n")
			for _, ex := range selected.Note.Examples {
				b.WriteString(fmt.Sprintf("    - %s\n", ex))
			}
		}
	}

	if len(m.drafts) > maxVisible {
		b.WriteString(fmt.Sprintf("\n(Showing %d-%d of %d)", start+1, end, len(m.drafts)))
	}
	return b.String()
}
