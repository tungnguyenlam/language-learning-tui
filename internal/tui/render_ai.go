package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderAI(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	ctx := NewRenderContext(m, layout, ViewAI)

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
	ctx.WriteLine(currTitleStyle.Render("AI Drafts") + spinner)
	if len(m.drafts) > 0 {
		countStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)
		ctx.Write(" " + countStyle.Render(fmt.Sprintf("[%d pending]", len(m.drafts))))
	}
	ctx.NewLine()
	ctx.NewLine()
	ctx.WriteLine(fmt.Sprintf("Deck: %s\nTemplate: %s (use [ / ])", m.deckLabel(), templateName))
	displayText := m.aiInput + "_"
	if m.aiInput == "" && !m.searchingAI {
		displayText = mutedStyle.Render("(e.g., business email, doctor visit, apartment viewing)") + "_"
	}
	ctx.WriteLine(searchStyle.Render(fmt.Sprintf("%s: %s", searchLabel, displayText)))
	ctx.WriteLine(fmt.Sprintf("\nPress %s to edit topic | %s generate | %s approve | %s discard | %s clear",
		keyStyle.Render("/"), keyStyle.Render("Enter"), keyStyle.Render("a"), keyStyle.Render("d"), keyStyle.Render("esc")))
	if m.aiProvider == nil {
		ctx.WriteLine(warnStyle.Render("AI provider is disabled. Enable Offline, Template, OpenAI, or Anthropic in Settings to generate drafts."))
	}

	ctx.WriteLine(mutedStyle.Render("Tip: include level and use case, e.g. B1 workplace small talk with 2 examples."))

	// Suggested Topics Section
	if (m.aiInput == "" || m.aiInput == "der Kaffee") && len(m.drafts) == 0 {
		ctx.WriteLine("\n" + infoStyle.Bold(true).Render("Click a topic or type your own, then press Enter:"))

		suggestions := []struct {
			topic string
			desc  string
		}{
			{"A1 survival", "basic phrases"},
			{"B1 doctor visit", "medical vocabulary"},
			{"B1 apartment viewing", "housing search"},
			{"B1 emergency & safety", "first aid & help"},
			{"B1 bureaucracy appointment", "forms & offices"},
			{"B2 digital privacy", "data protection"},
			{"A2 pharmacy visit", "medicine basics"},
			{"C1 project retrospective", "work reflection"},
			{"C1 business email", "formal writing"},
			{"B2 job interview", "professional talk"},
			{"B2 urban mobility", "transport & city"},
			{"B2 news debate", "current topics"},
			{"C1 academic argument", "university German"},
			{"travel phrases", "tourist essentials"},
			{"weather small talk", "daily conversation"},
			{"B1 at the restaurant", "ordering food"},
			{"B2 complaining politely", "diplomatic German"},
			{"B1 family & friends", "relationships"},
			{"B2 science & tech", "technology vocabulary"},
			{"B1 daily routine", "everyday activities"},
			{"C2 legal German", "juridical terms"},
			{"A1 colors & shapes", "descriptive words"},
			{"B1 sports & fitness", "athletic vocabulary"},
			{"B2 climate change", "sustainability"},
			{"B1 small talk", "casual conversation"},
			{"B1 weather & seasons", "daily German"},
			{"A2 shopping & clothing", "at the store"},
			{"B2 media & news", "journalism terms"},
			{"C1 philosophy & ethics", "abstract concepts"},
			{"A1 greetings formal", "Sie vs du"},
			{"B1 making reservations", "restaurant & hotel"},
			{"B2 discussing art", "museum & cinema"},
			{"B2 job application", "cover letter & CV"},
			{"C1 scientific paper", "research writing"},
			{"B2 environmental issues", "green topics"},
			{"A2 public transport", "train & bus tickets"},
			{"B1 insurance claims", "legal German"},
			{"grammar breakdown", "explain German grammar"},
			{"A1 city directions", "navigating German cities"},
			{"C1 psychology terms", "mind & emotions"},
			{"B2 travel adventure", "trips & flights"},
			{"B2 hiking & camping", "outdoor German"},
			{"B1 hotel & accommodation", "booking vocabulary"},
			{"B2 restaurant phrases", "dining out"},
			{"C1 philosophical concepts", "abstract thinking"},
			{"A2 giving directions", "navigation"},
			{"B2 car & driving", "vehicle vocabulary"},
			{"B1 phone & internet", "communication tech"},
			{"A2 doctor & health", "medical German"},
			{"sentence analysis", "structure & grammar"},
		}
		visibleSuggestions := len(suggestions)
		if layout.Height < 42 {
			visibleSuggestions = minInt(24, visibleSuggestions)
		} else if layout.Height < 55 {
			visibleSuggestions = minInt(32, visibleSuggestions)
		}

		suggestionStyle := lipgloss.NewStyle().
			Foreground(colorCyan).
			Underline(true).
			MarginRight(1)
		descStyle := lipgloss.NewStyle().Foreground(colorMuted)

		maxLineWidth := maxInt(24, layout.Width-2)

		for i, s := range suggestions[:visibleSuggestions] {
			separator := ""
			if ctx.currX > layout.X {
				separator = "  "
			}
			label := s.topic
			labelWidth := lipgloss.Width(separator) + lipgloss.Width(label)

			if ctx.currX > layout.X && ctx.currX-layout.X+labelWidth > maxLineWidth {
				ctx.NewLine()
				separator = ""
				labelWidth = lipgloss.Width(label)
			}

			ctx.Write(separator)
			topic := s.topic
			ctx.RegisterHitboxWithAction("ai-topic-"+s.topic, lipgloss.Width(s.topic), 1, func() tea.Cmd {
				m.aiInput = topic
				m.drafts = nil
				m.draftCursor = 0
				return m.startDrafting()
			})
			ctx.Write(suggestionStyle.Render(s.topic))

			if layout.Width >= 118 && i < 8 {
				desc := descStyle.Render(" (" + s.desc + ")")
				if ctx.currX-layout.X+lipgloss.Width(desc) <= maxLineWidth {
					ctx.Write(desc)
				}
			}
		}
		ctx.NewLine()
		if visibleSuggestions < len(suggestions) {
			ctx.WriteLine(mutedStyle.Render(fmt.Sprintf("Showing %d of %d suggestions. Type any topic for a custom draft.", visibleSuggestions, len(suggestions))))
		}
		ctx.WriteLine(mutedStyle.Render("Tip: be specific (level, situation, examples) for better results."))

		return ctx.String()
	}

	if len(m.drafts) == 0 {
		if m.drafting {
			ctx.WriteLine("\n" + infoStyle.Bold(true).Render("AI is crafting flashcards...") + spinner)
		} else {
			emptyBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted).
				Padding(0, 1).
				Render(infoStyle.Bold(true).Render("✨ Ready to create new flashcards!") + "\n" +
					mutedStyle.Render("No drafts yet. Type a topic and press Enter to generate.") + "\n" +
					mutedStyle.Render("Tip: include level and use case, e.g. B1 workplace small talk."))
			ctx.WriteLine("\n" + emptyBox)
		}
		return ctx.String()
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
	ctx.NewLine()
	listY := ctx.currY + listBoxStyle.GetBorderTopSize() + listBoxStyle.GetPaddingTop()
	for i := start; i < end; i++ {
		draftIdx := i
		prefix := "  "
		if i == m.draftCursor {
			prefix = "> "
		}
		item := fmt.Sprintf("%s%s -> %s", prefix, m.drafts[i].Note.Front, truncateLine(m.drafts[i].Note.Back, 40))
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
			Action: func() tea.Cmd {
				m.draftCursor = draftIdx
				return m.approveDraft()
			},
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("draft-discard-%d", i),
			View:   ViewAI,
			X:      discardX,
			Y:      rowY,
			Width:  lipgloss.Width(discardBtn),
			Height: 1,
			Action: func() tea.Cmd {
				m.draftCursor = draftIdx
				return m.discardDraft()
			},
		})
	}

	ctx.WriteLine(listBox)

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
		ctx.NewLine()
		ctx.WriteLine(previewBox)
	}

	if len(m.drafts) > maxVisible {
		ctx.WriteLine(fmt.Sprintf("\n(Showing %d-%d of %d)", start+1, end, len(m.drafts)))
	}
	return ctx.String()
}
