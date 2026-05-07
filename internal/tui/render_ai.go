package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderAI(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	var b strings.Builder
	spinner := ""
	if m.drafting {
		frames := []string{"-", "\\", "|", "/"}
		spinner = " " + frames[m.spinnerFrame%len(frames)]
	}
	templateName := "None"
	if set := m.currentAITemplateSet(); set != "" {
		templateName = set
	}

	searchBorderColor := "62"
	searchLabel := "Topic"
	if m.searchingAI {
		searchBorderColor = "81"
	}
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(searchBorderColor)).
		Padding(0, 1).
		Width(maxInt(30, width-20))

	b.WriteString("AI Drafts" + spinner + "\n\n")
	b.WriteString(fmt.Sprintf("Deck: %s\nTemplate: %s (use [ / ])\n", m.deckLabel(), templateName))
	displayText := m.aiInput + "_"
	if m.aiInput == "" && !m.searchingAI {
		placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		displayText = placeholderStyle.Render("(e.g., business email, doctor visit, apartment viewing)") + "_"
	}
	b.WriteString(searchStyle.Render(fmt.Sprintf("%s: %s", searchLabel, displayText)) + "\n")
	b.WriteString("\nPress / to edit topic | Enter generate | a approve | d discard | esc clear\n")
	b.WriteString(mutedStyle.Render("Tip: include level and use case, e.g. B1 workplace small talk.") + "\n")
	b.WriteString(mutedStyle.Render("Suggested topics: business email, doctor visit, apartment viewing.") + "\n")

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
	if layout.Height > 20 {
		maxVisible = layout.Height - 15
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
		item := fmt.Sprintf("%s%s -> %s", prefix, draft.Note.Front, truncateLine(draft.Note.Back, 40))

		rowY := layout.Y + strings.Count(b.String(), "\n")
		b.WriteString(style.Render(item))

		// Add interactive buttons
		btnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		approveBtn := " [Approve]"
		discardBtn := " [Discard]"

		if i == m.draftCursor {
			btnStyle = btnStyle.Foreground(lipgloss.Color("81"))
		}

		itemWidth := lipgloss.Width(item)
		approveX := layout.X + itemWidth
		discardX := approveX + lipgloss.Width(approveBtn)

		b.WriteString(btnStyle.Render(approveBtn))
		b.WriteString(btnStyle.Render(discardBtn))
		b.WriteString("\n")

		// Register hitboxes
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

		previewContent := previewTitleStyle.Render("Preview:") + "\n"
		previewContent += fmt.Sprintf("  Front:    %s\n", selected.Note.Front)
		previewContent += fmt.Sprintf("  Back:     %s\n", selected.Note.Back)
		previewContent += fmt.Sprintf("  Extra:    %s\n", selected.Note.Extra)
		if len(selected.Note.Tags) > 0 {
			previewContent += fmt.Sprintf("  Tags:     %s\n", strings.Join(selected.Note.Tags, ", "))
		}
		if len(selected.Note.Examples) > 0 {
			previewContent += "  Examples:\n"
			for _, ex := range selected.Note.Examples {
				previewContent += fmt.Sprintf("    - %s\n", ex)
			}
		}

		previewBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(maxInt(40, width-10)).
			Render(previewContent)
		b.WriteString("\n" + previewBox + "\n")
	}

	if len(m.drafts) > maxVisible {
		b.WriteString(fmt.Sprintf("\n(Showing %d-%d of %d)", start+1, end, len(m.drafts)))
	}
	return b.String()
}
