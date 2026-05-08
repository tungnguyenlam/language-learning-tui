package tui

import (
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderBrowser() string {
	return m.renderBrowserAt(m.activeViewContentLayout())
}

func (m *Model) renderBrowserAt(layout viewportLayout) string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Card Browser") + "\n")

	searchBorderColor := "62"
	searchLabel := "Search"
	if m.searchingBrowser {
		searchBorderColor = "81"
		searchLabel = "SEARCHING"
	}
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(searchBorderColor)).
		Padding(0, 1).
		Width(maxInt(30, layout.Width-60))

	if m.taggingCards {
		searchStyle = searchStyle.BorderForeground(lipgloss.Color("212"))
		searchLabel = "TAGS"
		b.WriteString(searchStyle.Render(fmt.Sprintf("%s: %s_", searchLabel, m.tagInput)) + "\n\n")
	} else {
		b.WriteString(searchStyle.Render(fmt.Sprintf("%s: %s_", searchLabel, m.browserSearch)) + "\n\n")
	}

	if len(m.browserCards) == 0 {
		b.WriteString("No cards found. Press / to search.\n\n")
		b.WriteString(mutedStyle.Render("Use left/right/[ to change deck filter.\n"))
		return b.String()
	}
	start := 0
	end := len(m.browserCards)
	maxVisible := m.listVisibleLines(layout.Height)
	if end > maxVisible {
		start = m.browserCursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(m.browserCards) {
			end = len(m.browserCards)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
	}
	listStartY := strings.Count(b.String(), "\n")
	lineWidth := scrollbarLineWidth(layout.Width)
	thumbStart, thumbHeight := scrollbarThumb(len(m.browserCards), maxVisible, start)
	for i := start; i < end; i++ {
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
		if card.Mature {
			mature = " ✨"
		}
		tags := ""
		if len(card.Tags) > 0 {
			tags = " " + mutedStyle.Render("#"+strings.Join(card.Tags, " #"))
		}
		label := fmt.Sprintf("%s%s[%s] %s%s%s%s%s%s", prefix, selected, kind, card.Prompt, mature, bookmark, leech, suspended, tags)
		line := padLine(style.Render(label), lineWidth)
		if len(m.browserCards) > maxVisible {
			currentPos := i - start
			scrollbarChar := "│"
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			line += " " + scrollbarChar
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("browser-scroll-%d", currentPos),
				View:   ViewBrowser,
				X:      layout.X + lineWidth + 1,
				Y:      layout.Y + listStartY + currentPos,
				Width:  1,
				Height: 1,
			})
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if m.showReviewHistory && m.browserCursor < len(m.browserCards) && m.reviewHistoryCard == m.browserCards[m.browserCursor].ID {
		b.WriteString("\n")
		b.WriteString(m.renderReviewHistory(m.browserCards[m.browserCursor].Prompt))
		b.WriteString("\n")
	}

	numSelected := 0
	for _, s := range m.browserSelected {
		if s {
			numSelected++
		}
	}

	if numSelected > 0 {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Render(fmt.Sprintf("\n%d cards selected. Bulk actions: %s (bookmark) | %s (suspend) | %s (type) | %s (tags) | %s (delete) | %s (clear selection)\n", numSelected,
			keyStyle.Render("B"), keyStyle.Render("X"), keyStyle.Render("t"), keyStyle.Render("T"), keyStyle.Render("Del"), keyStyle.Render("esc"))))
	} else {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(fmt.Sprintf("\nUse %s/%s to navigate, %s to select, %s to toggle kind, %s for tags, %s to cleanup tags, type to search, %s for history.\n",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("m"), keyStyle.Render("t"), keyStyle.Render("T"), keyStyle.Render("C"), keyStyle.Render("Enter")))
	}
	return b.String()
}

func (m *Model) renderReviewHistory(label string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Review History: %s\n", label))
	if len(m.reviewHistory) == 0 {
		b.WriteString("  No reviews yet.")
		return b.String()
	}
	for i, log := range m.reviewHistory {
		fmt.Fprintf(&b, "  %d. %s at %s -> next %s (%s, reviews %d, lapses %d)\n",
			i+1,
			log.Grade,
			log.Reviewed.Local().Format("Jan 02 15:04"),
			log.Due.Local().Format("Jan 02"),
			formatReviewInterval(log.Interval),
			log.Reviews,
			log.Lapses,
		)
	}
	return strings.TrimRight(b.String(), "\n")
}

func formatReviewInterval(interval time.Duration) string {
	if interval <= 0 {
		return "same day"
	}
	hours := int(interval.Hours())
	if hours < 24 {
		if hours <= 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := hours / 24
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
