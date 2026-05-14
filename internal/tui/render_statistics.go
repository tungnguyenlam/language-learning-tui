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
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159"))
	title := "Statistics: " + m.deckLabel()
	content.WriteString(titleStyle.Render(title) + "\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)

	// --- Column 1: Collection & Review Stats ---
	col1 := strings.Builder{}
	col1.WriteString(titleStyle.Copy().Underline(true).Render("Collection") + "\n")
	col1.WriteString(fmt.Sprintf("%s %d\n", labelStyle.Render("Total Cards:"), m.stats.TotalCards))
	col1.WriteString(fmt.Sprintf("%s %d (%d active)\n", labelStyle.Render("Total Decks:"), m.stats.TotalDecks, m.stats.ActiveDecks))
	col1.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("New:"), valueStyle.Render(fmt.Sprintf("%d", m.stats.NewCards))))
	col1.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Young:"), valueStyle.Render(fmt.Sprintf("%d", m.stats.YoungCards))))
	col1.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Mature:"), valueStyle.Render(fmt.Sprintf("%d", m.stats.MatureCards))))
	col1.WriteString(fmt.Sprintf("  %s %d (%d due)\n", labelStyle.Render("Bookmarked:"), m.stats.BookmarkedCards, m.stats.BookmarkedDue))
	col1.WriteString(fmt.Sprintf("  %s %d\n", labelStyle.Render("Leech:"), m.stats.LeechCards))
	col1.WriteString(fmt.Sprintf("  %s %d\n", labelStyle.Render("Suspended:"), m.stats.SuspendedCards))

	// --- Column 2: Today & Performance ---
	col2 := strings.Builder{}
	col2.WriteString(titleStyle.Copy().Underline(true).Render("Today's Performance") + "\n")
	col2.WriteString(fmt.Sprintf("%s %d\n", labelStyle.Render("Total Reviews:"), m.stats.TotalReviews))
	col2.WriteString(fmt.Sprintf("%s %d/%d\n", labelStyle.Render("Reviews Today:"), m.stats.ReviewsToday, m.stats.DailyGoal))

	percentage := 0.0
	if m.stats.DailyGoal > 0 {
		percentage = float64(m.stats.ReviewsToday) / float64(m.stats.DailyGoal)
	}
	barColor := "196"
	if percentage >= 1.0 {
		barColor = "46"
	} else if percentage >= 0.5 {
		barColor = "226"
	}
	bar := progressBar(25, percentage, barColor, "238")
	col2.WriteString(fmt.Sprintf("  %s %s\n", bar, valueStyle.Render(fmt.Sprintf("%.0f%%", percentage*100))))

	streakIndicator := ""
	switch {
	case m.stats.CurrentStreak >= 100:
		streakIndicator = " 🏆🔥🔥"
	case m.stats.CurrentStreak >= 30:
		streakIndicator = " 🔥🔥"
	case m.stats.CurrentStreak >= 14:
		streakIndicator = " 🔥✨"
	case m.stats.CurrentStreak >= 7:
		streakIndicator = " 🔥"
	case m.stats.CurrentStreak > 0:
		streakIndicator = " ⚡"
	}
	col2.WriteString(fmt.Sprintf("%s %s%s\n", labelStyle.Render("Current Streak:"), valueStyle.Render(fmt.Sprintf("%d days", m.stats.CurrentStreak)), streakIndicator))

	successColor := "196"
	if m.stats.SuccessRate >= 0.85 {
		successColor = "46"
	} else if m.stats.SuccessRate >= 0.70 {
		successColor = "226"
	}
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(successColor)).Bold(true)
	col2.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Success Rate:"), successStyle.Render(fmt.Sprintf("%.1f%%", m.stats.SuccessRate*100))))

	retentionColor := "196"
	totalMature := m.stats.MatureCards + m.stats.YoungCards
	retentionRate := 0.0
	if totalMature > 0 {
		retentionRate = float64(m.stats.MatureCards) / float64(totalMature)
	}
	if retentionRate >= 0.80 {
		retentionColor = "46"
	} else if retentionRate >= 0.60 {
		retentionColor = "226"
	}
	retentionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(retentionColor)).Bold(true)
	col2.WriteString(fmt.Sprintf("%s %s (Mature/Total: %d/%d)\n", labelStyle.Render("Retention:"), retentionStyle.Render(fmt.Sprintf("%.1f%%", retentionRate*100)), m.stats.MatureCards, totalMature))

	remainingGoal := maxInt(0, m.stats.DailyGoal-m.stats.ReviewsToday)
	availableNow := len(m.dueCards)
	forecastColor := "46"
	forecastText := "On track"
	if remainingGoal > 0 && availableNow == 0 {
		forecastColor = "81"
		forecastText = "No due cards"
	} else if remainingGoal > availableNow {
		forecastColor = "226"
		forecastText = "Need new cards"
	}
	forecastStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(forecastColor)).Bold(true)
	col2.WriteString(fmt.Sprintf("%s %s (%d left)\n\n",
		labelStyle.Render("Forecast:"),
		forecastStyle.Render(forecastText),
		remainingGoal))

	// Combine columns
	col1Str := col1.String()
	col2Str := col2.String()
	maxWidth := layout.Width
	colWidth := (maxWidth - 4) / 2

	joinedCols := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colWidth).Render(col1Str),
		lipgloss.NewStyle().PaddingLeft(4).Render(col2Str),
	)
	content.WriteString(joinedCols + "\n")

	// --- Success Rate per Deck ---
	if len(m.decks) > 1 {
		deckSuccess := strings.Builder{}
		deckSuccess.WriteString(titleStyle.Copy().Underline(true).Render("Success Rate per Deck") + "\n")
		hasData := false
		for _, d := range m.decks {
			if d.ID == "" || d.SuccessRate <= 0 {
				continue
			}
			hasData = true
			successColor := "196"
			if d.SuccessRate >= 0.85 {
				successColor = "46"
			} else if d.SuccessRate >= 0.70 {
				successColor = "226"
			}
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(successColor)).Bold(true)
			bar := progressBar(20, d.SuccessRate, successColor, "238")
			deckSuccess.WriteString(fmt.Sprintf("  %-20s %s %s\n", truncateLine(d.Name, 20), bar, successStyle.Render(fmt.Sprintf("%.1f%%", d.SuccessRate*100))))
		}
		if hasData {
			content.WriteString(deckSuccess.String() + "\n")
		}
	}

	// --- Session Stats ---
	content.WriteString(titleStyle.Copy().Underline(true).Render("Session Stats:") + "\n")
	if m.sessionReviewed > 0 {
		accuracy := float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
		var accuracyColor string
		if accuracy >= 80 {
			accuracyColor = "46"
		} else if accuracy >= 60 {
			accuracyColor = "226"
		} else {
			accuracyColor = "197"
		}
		accStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(accuracyColor)).Bold(true)

		duration := time.Since(m.sessionStartTime).Round(time.Second)
		speed := 0.0
		if duration.Seconds() > 0 {
			speed = float64(m.sessionReviewed) / (duration.Minutes())
		}

		content.WriteString(fmt.Sprintf("  %s %d | %s %d | %s %s\n",
			labelStyle.Render("Reviewed:"), m.sessionReviewed,
			labelStyle.Render("Correct:"), m.sessionCorrect,
			labelStyle.Render("Accuracy:"), accStyle.Render(fmt.Sprintf("%.1f%%", accuracy))))
		content.WriteString(fmt.Sprintf("  %s %s | %s %s cards/min\n\n",
			labelStyle.Render("Duration:"), valueStyle.Render(duration.String()),
			labelStyle.Render("Speed:"), valueStyle.Render(fmt.Sprintf("%.1f", speed))))
	} else {
		content.WriteString("  (no reviews yet this session)\n\n")
	}

	// --- Grades Chart ---
	content.WriteString(titleStyle.Copy().Underline(true).Render("Reviews by Grade") + "\n")
	grades := []core.ReviewGrade{core.GradeAgain, core.GradeHard, core.GradeGood, core.GradeEasy}
	gradeColors := map[core.ReviewGrade]string{
		core.GradeAgain: "196",
		core.GradeHard:  "208",
		core.GradeGood:  "46",
		core.GradeEasy:  "81",
	}
	gradeIcons := map[core.ReviewGrade]string{
		core.GradeAgain: "✗",
		core.GradeHard:  "~",
		core.GradeGood:  "✓",
		core.GradeEasy:  "★",
	}

	maxGradeCount := 0
	for _, g := range grades {
		if m.stats.Grades[g] > maxGradeCount {
			maxGradeCount = m.stats.Grades[g]
		}
	}

	for _, g := range grades {
		count := m.stats.Grades[g]
		gradeLabel := fmt.Sprintf("%-5s", g)
		gPercentage := 0.0
		if maxGradeCount > 0 {
			gPercentage = float64(count) / float64(maxGradeCount)
		}
		gBar := progressBar(20, gPercentage, gradeColors[g], "238")
		content.WriteString(fmt.Sprintf("  %s %s %d  %s\n", labelStyle.Render(gradeLabel), gBar, count, lipgloss.NewStyle().Foreground(lipgloss.Color(gradeColors[g])).Render(gradeIcons[g])))
	}

	// --- 7-Day Activity Chart ---
	content.WriteString("\n" + titleStyle.Copy().Underline(true).Render("Recent Activity (Last 7 Days)") + "\n")
	now := time.Now().Local()
	maxReviews := 0
	last7Days := make([]struct {
		date  string
		count int
	}, 7)

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -6+i).Format("2006-01-02")
		count := m.reviewsPerDay[date]
		last7Days[i] = struct {
			date  string
			count int
		}{date, count}
		if count > maxReviews {
			maxReviews = count
		}
	}

	for i := 0; i < 7; i++ {
		d := last7Days[i]
		dateObj, _ := time.Parse("2006-01-02", d.date)
		dayName := dateObj.Format("Mon")

		percentage := 0.0
		if maxReviews > 0 {
			percentage = float64(d.count) / float64(maxReviews)
		}

		barColor := "240"
		if d.count > 0 {
			barColor = "34"
			if d.count >= m.stats.DailyGoal && m.stats.DailyGoal > 0 {
				barColor = "46"
			}
		}

		bar := progressBar(20, percentage, barColor, "238")
		content.WriteString(fmt.Sprintf("  %s %s %d\n", labelStyle.Render(dayName), bar, d.count))
	}

	// --- Maturity Distribution ---
	content.WriteString("\n" + titleStyle.Copy().Underline(true).Render("Maturity Distribution") + "\n")
	totalKnown := m.stats.NewCards + m.stats.YoungCards + m.stats.MatureCards
	if totalKnown > 0 {
		newPct := float64(m.stats.NewCards) / float64(totalKnown)
		youngPct := float64(m.stats.YoungCards) / float64(totalKnown)
		maturePct := float64(m.stats.MatureCards) / float64(totalKnown)

		content.WriteString(fmt.Sprintf("  %-8s %s %d (%.1f%%)\n", labelStyle.Render("New:"), progressBar(20, newPct, "81", "238"), m.stats.NewCards, newPct*100))
		content.WriteString(fmt.Sprintf("  %-8s %s %d (%.1f%%)\n", labelStyle.Render("Young:"), progressBar(20, youngPct, "220", "238"), m.stats.YoungCards, youngPct*100))
		content.WriteString(fmt.Sprintf("  %-8s %s %d (%.1f%%)\n", labelStyle.Render("Mature:"), progressBar(20, maturePct, "46", "238"), m.stats.MatureCards, maturePct*100))
	} else {
		content.WriteString("  (no cards in collection)\n")
	}

	// Review Heatmap (last 3 months)
	content.WriteString("\n" + titleStyle.Copy().Underline(true).Render("Review Heatmap (last 3 months)") + "\n")
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
	padWidth := layout.Width - 2
	thumbStart, thumbHeight := scrollbarThumb(totalLines, maxVisible, m.statsScroll)

	var visibleContent strings.Builder
	for i := m.statsScroll; i < m.statsScroll+maxVisible && i < totalLines; i++ {
		line := padLine(lines[i], padWidth)
		if totalLines > maxVisible {
			scrollbarChar := "│"
			currentPos := i - m.statsScroll
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			line = line + " " + scrollbarChar

			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("stats-scroll-%d", currentPos),
				View:   ViewStatistics,
				X:      layout.X + padWidth + 1,
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
		footerText := fmt.Sprintf("\nUse %s/%s or Mouse Wheel to scroll. Lines %d-%d of %d.",
			keyStyle.Render("j"), keyStyle.Render("k"), m.statsScroll+1, minInt(m.statsScroll+maxVisible, totalLines), totalLines)
		footer = lipgloss.NewStyle().Foreground(colorMuted).Render(footerText)
	}

	return visibleContent.String() + footer
}
