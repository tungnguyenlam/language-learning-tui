package tui

import (
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderStatistics(x, y int) string {
	width, height := m.activePanelSize()
	layout := contentLayoutForStyle(panelStyle.Width(width).Height(height), x, y)
	return m.renderStatisticsAt(layout)
}

func (m *Model) renderStatisticsAt(layout viewportLayout) string {
	var content strings.Builder
	title := "Statistics: " + m.deckLabel()
	content.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Render(title) + "\n\n")

	content.WriteString(fmt.Sprintf("Total Cards:   %d\n", m.stats.TotalCards))
	content.WriteString(fmt.Sprintf("Total Decks:   %d (%d active)\n", m.stats.TotalDecks, m.stats.ActiveDecks))
	content.WriteString(fmt.Sprintf("  New:         %d\n", m.stats.NewCards))
	content.WriteString(fmt.Sprintf("  Young:       %d\n", m.stats.YoungCards))
	content.WriteString(fmt.Sprintf("  Mature:      %d\n\n", m.stats.MatureCards))
	content.WriteString(fmt.Sprintf("  Bookmarked:  %d (%d due)\n", m.stats.BookmarkedCards, m.stats.BookmarkedDue))
	content.WriteString(fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards))
	content.WriteString(fmt.Sprintf("  Suspended:   %d\n\n", m.stats.SuspendedCards))

	content.WriteString(fmt.Sprintf("Total Reviews: %d\n", m.stats.TotalReviews))
	content.WriteString(fmt.Sprintf("Reviews Today: %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal))
	// Colored progress bar for daily goal
	percentage := 0.0
	if m.stats.DailyGoal > 0 {
		percentage = float64(m.stats.ReviewsToday) / float64(m.stats.DailyGoal)
	}
	// Color: green if complete, yellow if halfway, red otherwise
	barColor := "196" // red
	if percentage >= 1.0 {
		barColor = "46" // green
	} else if percentage >= 0.5 {
		barColor = "226" // yellow
	}
	bar := progressBar(30, percentage, barColor, "238")
	content.WriteString(fmt.Sprintf("  %s %.0f%%\n", bar, percentage*100))
	// Streak with fire emoji if > 0
	streakIndicator := ""
	if m.stats.CurrentStreak > 0 {
		streakIndicator = " 🔥"
	}
	content.WriteString(fmt.Sprintf("Current Streak: %d days%s\n", m.stats.CurrentStreak, streakIndicator))
	content.WriteString(fmt.Sprintf("Success Rate:  %.1f%%\n\n", m.stats.SuccessRate*100))

	content.WriteString("Session Stats:\n")
	content.WriteString(fmt.Sprintf("  Reviewed:    %d\n", m.sessionReviewed))
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		var accuracyColor string
		if accuracy >= 80 {
			accuracyColor = "46" // green
		} else if accuracy >= 60 {
			accuracyColor = "226" // yellow
		} else {
			accuracyColor = "197" // red
		}
		content.WriteString(fmt.Sprintf("  Correct:     %d\n", m.sessionCorrect))
		accStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(accuracyColor))
		content.WriteString(fmt.Sprintf("  Accuracy:    %s\n\n", accStyle.Render(fmt.Sprintf("%.1f%%", accuracy))))
	} else {
		content.WriteString("  (no reviews yet)\n\n")
	}

	content.WriteString("Reviews by Grade:\n")
	grades := []core.ReviewGrade{core.GradeAgain, core.GradeHard, core.GradeGood, core.GradeEasy}
	for _, g := range grades {
		count := m.stats.Grades[g]
		content.WriteString(fmt.Sprintf("  %s: %d\n", g, count))
	}

	// Review Heatmap (last 3 months)
	content.WriteString("\nReview Heatmap (last 3 months):\n")
	if len(m.reviewsPerDay) == 0 {
		content.WriteString("  (no review data yet)\n")
	} else {
		now := time.Now().UTC()
		// Start 13 weeks ago
		startDate := now.AddDate(0, 0, -13*7)
		// Align to Sunday (0)
		for startDate.Weekday() != time.Sunday {
			startDate = startDate.AddDate(0, 0, -1)
		}

		days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
		for row := 0; row < 7; row++ {
			content.WriteString(fmt.Sprintf("%s ", days[row]))
			for col := 0; col < 14; col++ {
				day := startDate.AddDate(0, 0, col*7+row)
				if day.After(now) {
					content.WriteString("  ")
					continue
				}
				dayStr := day.Format("2006-01-02")
				count := m.reviewsPerDay[dayStr]

				char := "░" // 0 reviews
				color := "236"
				if count > 0 {
					if count >= 20 {
						char = "█"
						color = "46" // bright green
					} else if count >= 10 {
						char = "▓"
						color = "40" // green
					} else if count >= 5 {
						char = "▒"
						color = "34" // dark green
					} else {
						char = "░"
						color = "22" // very dark green
					}
				}
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				content.WriteString(style.Render(char) + " ")
			}
			content.WriteString("\n")
		}
		content.WriteString("    ")
		for col := 0; col < 14; col++ {
			day := startDate.AddDate(0, 0, col*7)
			if day.Day() <= 7 { // New month
				content.WriteString(day.Format("Jan")[:1] + " ")
			} else {
				content.WriteString("  ")
			}
		}
		content.WriteString("\n")
	}

	lines := strings.Split(content.String(), "\n")
	totalLines := len(lines)
	m.statsTotalLines = totalLines
	maxVisible := m.statisticsVisibleLines(layout.Height)

	m.statsScroll = clampInt(m.statsScroll, 0, maxInt(0, totalLines-maxVisible))
	lineWidth := scrollbarLineWidth(layout.Width)
	thumbStart, thumbHeight := scrollbarThumb(totalLines, maxVisible, m.statsScroll)

	var visibleContent strings.Builder
	for i := m.statsScroll; i < m.statsScroll+maxVisible && i < totalLines; i++ {
		line := padLine(lines[i], lineWidth)
		if totalLines > maxVisible {
			scrollbarChar := "│"
			currentPos := i - m.statsScroll
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			line = line + " " + scrollbarChar

			// Register hitbox for this line of the scrollbar
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("stats-scroll-%d", currentPos),
				View:   ViewStatistics,
				X:      layout.X + lineWidth + 1,
				Y:      layout.Y + currentPos,
				Width:  1,
				Height: 1,
			})
		}
		visibleContent.WriteString(line + "\n")
	}

	footer := ""
	if totalLines > maxVisible {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		footer = fmt.Sprintf("\n%s", mutedStyle.Render(fmt.Sprintf("Use %s/%s or Mouse Wheel to scroll. Lines %d-%d of %d.",
			keyStyle.Render("j"), keyStyle.Render("k"), m.statsScroll+1, minInt(m.statsScroll+maxVisible, totalLines), totalLines)))
	}

	return visibleContent.String() + footer
}
