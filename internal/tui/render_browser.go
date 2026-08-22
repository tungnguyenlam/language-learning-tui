package tui

import (
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderBrowserAt(layout viewportLayout) string {
	ctx := NewRenderContext(m, layout, ViewBrowser)

	numSelected := 0
	for _, s := range m.browserSelected {
		if s {
			numSelected++
		}
	}

	titleText := "Card Browser"
	if numSelected > 0 {
		titleText = fmt.Sprintf("Card Browser (%d selected)", numSelected)
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	ctx.WriteLine(titleStyle.Render(titleText))

	searchBorderColor := "62"
	searchLabel := "Search"
	if m.searchingBrowser {
		searchBorderColor = "81"
		searchLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("SEARCHING")
	}
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(searchBorderColor)).
		Padding(0, 1).
		Width(maxInt(30, layout.Width-60))

	if m.taggingCards {
		searchStyle = searchStyle.BorderForeground(lipgloss.Color("212"))
		searchLabel = "TAGS"
		ctx.WriteLine(searchStyle.Render(fmt.Sprintf("%s: %s_", searchLabel, m.tagInput)) + "\n")
	} else if m.searchingTags {
		searchStyle = searchStyle.BorderForeground(lipgloss.Color("46"))
		searchLabel = "FILTER BY TAG"
		ctx.WriteLine(searchStyle.Render(fmt.Sprintf("%s: %s_", searchLabel, m.browserTag)) + "\n")
	} else {
		if m.searchingBrowser && len(m.browserSearchHistory) > 0 {
			historyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
			historyStr := "History: " + strings.Join(m.browserSearchHistory, " | ")
			ctx.WriteLine(historyStyle.Render(historyStr))
		}
		filterText := ""
		if m.browserTag != "" {
			filterText = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(" [Tag: " + m.browserTag + "]")
		}
		ctx.WriteLine(searchStyle.Render(fmt.Sprintf("%s: %s_", searchLabel, m.browserSearch)) + filterText + "\n")
	}

	if len(m.browserCards) == 0 {
		ctx.WriteLine(fmt.Sprintf("No cards found in %s.\n", m.deckLabel()))
		if m.browserSearch != "" || m.browserTag != "" {
			ctx.WriteLine(mutedStyle.Render("Press Esc to clear search/tag filters, or use [ / ] to change deck filter."))
		} else {
			ctx.WriteLine(mutedStyle.Render("Press / to search, # to filter by tag, or use [ / ] to change deck filter."))
		}
		return ctx.String()
	}

	totalLines := len(m.browserCards)
	availableHeight := m.listVisibleLines(layout.Height) - (ctx.currY - layout.Y)
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Auto-scroll to cursor
	m.browserScroll = AutoScroll(m.browserCursor, m.browserScroll, availableHeight, totalLines)

	var content strings.Builder
	lineWidth := layout.Width - 2

	endIdx := minInt(totalLines, m.browserScroll+availableHeight)
	for i := m.browserScroll; i < endIdx; i++ {
		card := m.browserCards[i]

		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.browserCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}

		selected := "[ ] "
		if m.browserSelected[card.ID] {
			selected = "[x] "
			if i != m.browserCursor {
				style = style.Foreground(lipgloss.Color("212"))
			}
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
		interval := ""
		if card.Interval > 0 {
			interval = " (" + formatReviewInterval(card.Interval) + ")"
		}
		if card.Mature {
			mature = " ✨"
		}

		stateIcon := mutedStyle.Render("○") // New
		if card.Mature {
			stateIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("●") // Mature
		} else if card.Reviews > 0 {
			stateIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("◐") // Learning
		}

		plainTags := ""
		if len(card.Tags) > 0 {
			plainTags = " #" + strings.Join(card.Tags, " #")
		}

		otherWidth := lipgloss.Width(prefix + selected + stateIcon + " [" + kind + "] " + interval + mature + bookmark + leech + suspended)
		availWidth := lineWidth - otherWidth
		if availWidth < 10 {
			availWidth = 10
		}

		maxTagsWidth := availWidth * 4 / 10
		if maxTagsWidth > 35 {
			maxTagsWidth = 35
		}
		if maxTagsWidth < 5 && len(card.Tags) > 0 {
			maxTagsWidth = 5
		}

		truncatedTags := ""
		if plainTags != "" {
			truncatedTags = truncateLine(plainTags, maxTagsWidth)
		}

		promptWidth := availWidth - lipgloss.Width(truncatedTags)
		if promptWidth < 5 {
			promptWidth = 5
		}
		truncatedPrompt := truncateLine(card.Prompt, promptWidth)
		highlightStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)
		highlightedPrompt := highlightQuery(truncatedPrompt, m.browserSearch, highlightStyle)

		styledTags := ""
		if truncatedTags != "" {
			styledTags = " " + mutedStyle.Render(strings.TrimSpace(truncatedTags))
		}

		label := fmt.Sprintf("%s%s%s [%s] %s%s%s%s%s%s%s", prefix, selected, stateIcon, kind, highlightedPrompt, interval, mature, bookmark, leech, suspended, styledTags)
		content.WriteString(style.Render(label) + "\n")
	}

	contentStr := content.String()

	listView := m.RenderList(layout.WithHeight(availableHeight).WithY(ctx.currY), contentStr, ListOptions{
		HitboxPrefix: "browser",
		View:         ViewBrowser,
		ScrollOffset: &m.browserScroll,
		TotalLines:   &totalLines,
	})
	ctx.Write(listView)

	if m.showReviewHistory && m.browserCursor >= 0 && m.browserCursor < len(m.browserCards) && m.reviewHistoryCard == m.browserCards[m.browserCursor].ID {
		ctx.NewLine()
		ctx.WriteLine(m.renderReviewHistory(m.browserCards[m.browserCursor].Prompt))
		ctx.NewLine()
	} else if layout.Height > 25 && len(m.browserCards) > 0 && m.browserCursor >= 0 && m.browserCursor < len(m.browserCards) {
		selected := m.browserCards[m.browserCursor]
		extra := selected.Extra
		if extra == "" {
			extra = "(none)"
		}
		tags := strings.Join(selected.Tags, ", ")
		if tags == "" {
			tags = "(none)"
		}

		kind := "FC"
		if selected.Kind == core.CardKindMCQ {
			kind = "MCQ"
		}

		state := "NEW"
		stateColor := "208"
		if selected.Mature {
			state = "MATURE"
			stateColor = "34"
		} else if selected.Reviews > 0 {
			state = "LEARNING"
			stateColor = "39"
		}

		stateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(stateColor)).Bold(true)
		intervalStr := formatReviewInterval(selected.Interval)
		if selected.Interval <= 0 {
			intervalStr = "none"
		}

		lastReviewedStr := "never"
		if !selected.LastReviewed.IsZero() && selected.LastReviewed.Year() > 1970 {
			lastReviewedStr = formatDuration(time.Since(selected.LastReviewed)) + " ago"
		}

		previewWidth := maxInt(35, layout.Width-10)
		previewContent := fmt.Sprintf("%s\n\n%s: %s  |  %s: %s  |  %s: %s\n\n%s: %s\n%s: %s\n\n%s: %s  |  %s: %s  |  %s: %s\n%s: %s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Card Preview:"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Deck"), truncateLine(m.deckNameByID(selected.DeckID), previewWidth/2-10),
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Kind"), kind,
			stateStyle.Render("State"), state,
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Front"), truncateLine(selected.Prompt, previewWidth/2-6),
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Back"), truncateLine(selected.Answer, previewWidth/2-6),
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Reviews"), fmt.Sprintf("%d", selected.Reviews),
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Interval"), intervalStr,
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Tags"), tags,
			lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Last Reviewed"), lastReviewedStr)

		if hint := renderGrammarHint(selected); hint != "" {
			previewContent += "\n\n" + hint
		}

		previewBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(previewWidth).
			Render(previewContent)
		ctx.NewLine()
		ctx.WriteLine(previewBox)
		ctx.NewLine()
	}

	if numSelected > 0 {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		ctx.WriteLine(lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Render(fmt.Sprintf("\n%d cards selected. Bulk actions: %s (bookmark/un) | %s (suspend/un) | %s (type) | %s (tags) | %s (delete) | %s (select/deselect all) | %s (clear selection)", numSelected,
			keyStyle.Render("b/B"), keyStyle.Render("x/X"), keyStyle.Render("t"), keyStyle.Render("T"), keyStyle.Render("Del"), keyStyle.Render("a"), keyStyle.Render("esc"))))
	} else {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		ctx.WriteLine(fmt.Sprintf("\nUse %s/%s to navigate, %s to select, %s to select all, %s to play audio, %s to toggle kind, %s for tags, %s to cleanup tags, %s to search, %s for history, %s to delete.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("m"), keyStyle.Render("a"), keyStyle.Render("p"), keyStyle.Render("t"), keyStyle.Render("T"), keyStyle.Render("C"), keyStyle.Render("/"), keyStyle.Render("enter"), keyStyle.Render("Backspace")))
	}
	return ctx.String()
}
