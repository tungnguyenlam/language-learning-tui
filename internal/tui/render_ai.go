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
		// Enhanced spinner with more visual appeal
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinner = " " + infoStyle.Render(frames[m.spinnerFrame%len(frames)])
	}
	templateName := "None"
	if set := m.currentAITemplateSet(); set != "" {
		templateName = set
	}

	searchBorderColor := colorPanel
	searchLabel := "Topic"
	if m.searchingAI {
		searchBorderColor = colorBlue
	}
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(searchBorderColor).
		Padding(0, 1).
		Width(maxInt(30, width-20))

	currTitleStyle := dashTitleStyle
	if m.drafting {
		currTitleStyle = currTitleStyle.Foreground(colorBlue)
	}
	b.WriteString(currTitleStyle.Render("AI Drafts") + spinner + "\n\n")
	b.WriteString(fmt.Sprintf("Deck: %s\nTemplate: %s (use [ / ])\n", m.deckLabel(), templateName))
	displayText := m.aiInput + "_"
	if m.aiInput == "" && !m.searchingAI {
		displayText = mutedStyle.Render("(e.g., business email, doctor visit, apartment viewing)") + "_"
	}
	b.WriteString(searchStyle.Render(fmt.Sprintf("%s: %s", searchLabel, displayText)) + "\n")
	b.WriteString(fmt.Sprintf("\nPress %s to edit topic | %s generate | %s approve | %s discard | %s clear\n",
		keyStyle.Render("/"), keyStyle.Render("Enter"), keyStyle.Render("a"), keyStyle.Render("d"), keyStyle.Render("esc")))
	if m.aiProvider == nil {
		b.WriteString(warnStyle.Render("AI provider is disabled. Enable Offline or Template in Settings to generate drafts.") + "\n")
	}

	b.WriteString(mutedStyle.Render("Tip: include level and use case, e.g. B1 workplace small talk with 2 examples.") + "\n")

	// Suggested Topics Section
	if (m.aiInput == "" || m.aiInput == "der Kaffee") && len(m.drafts) == 0 {
		b.WriteString("\n" + infoStyle.Bold(true).Render("Suggested Topics:") + "\n")
		suggestions := []string{"A1 survival", "B1 doctor visit", "B1 apartment viewing", "B2 job interview", "B2 urban mobility", "B2 news debate", "C1 business email", "C1 academic argument", "travel phrases", "weather small talk"}

		suggestionStyle := lipgloss.NewStyle().
			Foreground(colorCyan).
			Underline(true).
			MarginRight(2)

		lineY := layout.Y + strings.Count(b.String(), "\n")
		currentX := layout.X

		for i, s := range suggestions {
			if i > 0 && i%4 == 0 {
				b.WriteString("\n")
				lineY++
				currentX = layout.X
			}
			b.WriteString(suggestionStyle.Render(s))
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     "ai-topic-" + s,
				View:   ViewAI,
				X:      currentX,
				Y:      lineY,
				Width:  len(s),
				Height: 1,
			})
			currentX += len(s) + 2
		}
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("Suggested levels: A1 survival, B1 errands, B2 opinions, C1 formal writing.") + "\n")
	}

	if len(m.drafts) == 0 {
		if m.drafting {
			b.WriteString("\n" + infoStyle.Bold(true).Render("AI is crafting flashcards...") + spinner)
		} else {
			emptyBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted).
				Padding(0, 1).
				Render(infoStyle.Bold(true).Render("✨ Ready to create new flashcards!") + "\n" +
					mutedStyle.Render("No drafts yet. Type a topic and press Enter to generate.") + "\n" +
					mutedStyle.Render("Tip: include level and use case, e.g. B1 workplace small talk."))
			b.WriteString("\n" + emptyBox + "\n")
		}
		return b.String()
	}

	start := 0
	end := len(m.drafts)
	maxVisible := 10
	if layout.Height > 25 {
		maxVisible = layout.Height - 18
	}
	if maxVisible < 3 {
		maxVisible = 3
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

	var listBuilder strings.Builder
	for i := start; i < end; i++ {
		prefix := "  "
		currStyle := lipgloss.NewStyle()
		if i == m.draftCursor {
			prefix = "> "
			currStyle = currStyle.Bold(true).Foreground(colorPink)
		}
		draft := m.drafts[i]
		item := fmt.Sprintf("%s%s -> %s", prefix, draft.Note.Front, truncateLine(draft.Note.Back, 40))

		listBuilder.WriteString(currStyle.Render(item))

		// Add interactive buttons
		currBtnStyle := lipgloss.NewStyle().Foreground(colorPanel)
		approveBtn := " [Approve]"
		discardBtn := " [Discard]"

		if i == m.draftCursor {
			currBtnStyle = currBtnStyle.Foreground(colorBlue)
		}

		listBuilder.WriteString(currBtnStyle.Render(approveBtn))
		listBuilder.WriteString(currBtnStyle.Render(discardBtn))
		listBuilder.WriteString("\n")
	}

	listBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(maxInt(40, width-10))

	listBox := listBoxStyle.Render(listBuilder.String())

	// Re-calculate hitboxes for the list items
	listY := layout.Y + strings.Count(b.String(), "\n") + listBoxStyle.GetBorderTopSize() + listBoxStyle.GetPaddingTop()
	for i := start; i < end; i++ {
		draft := m.drafts[i]
		prefix := "  "
		if i == m.draftCursor {
			prefix = "> "
		}
		item := fmt.Sprintf("%s%s -> %s", prefix, draft.Note.Front, truncateLine(draft.Note.Back, 40))
		itemWidth := lipgloss.Width(item)

		approveBtn := " [Approve]"
		discardBtn := " [Discard]"

		approveX := layout.X + listBoxStyle.GetBorderLeftSize() + listBoxStyle.GetPaddingLeft() + itemWidth
		discardX := approveX + lipgloss.Width(approveBtn)
		rowY := listY + (i - start)

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

	b.WriteString("\n" + listBox + "\n")

	// Show detailed preview for selected draft
	if len(m.drafts) > 0 && m.draftCursor < len(m.drafts) {
		selected := m.drafts[m.draftCursor]
		previewTitleLabelStyle := lipgloss.NewStyle().Bold(true).Foreground(colorAccent)

		previewContent := previewTitleLabelStyle.Render("Preview:") + "\n"
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
			BorderForeground(colorPanel).
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
