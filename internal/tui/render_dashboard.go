package tui

import (
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/content"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderDashboard(layout viewportLayout) string {
	isNarrow := layout.Width < 55
	boxWidth := maxInt(25, (layout.Width-2)/2)
	if isNarrow {
		boxWidth = layout.Width - 2
	}

	streakIndicator := ""
	switch {
	case m.stats.CurrentStreak >= 100:
		streakIndicator = lipgloss.NewStyle().Foreground(colorGold).Bold(true).Render(" 🏆🔥🔥")
	case m.stats.CurrentStreak >= 30:
		streakIndicator = lipgloss.NewStyle().Foreground(colorGold).Render(" 🔥🔥")
	case m.stats.CurrentStreak >= 14:
		streakIndicator = lipgloss.NewStyle().Foreground(colorOrange).Render(" 🔥✨")
	case m.stats.CurrentStreak >= 7:
		streakIndicator = lipgloss.NewStyle().Foreground(colorOrange).Render(" 🔥")
	case m.stats.CurrentStreak > 0:
		streakIndicator = lipgloss.NewStyle().Foreground(colorYellow).Render(" ⚡")
	}

	reviewQueue := dashReviewStyle.
		Width(boxWidth).
		Render(dashStatsStyle.Render("Review Queue") + "\n" +
			fmt.Sprintf("  Due cards:   %d\n", len(m.dueCards)) +
			fmt.Sprintf("  Next 24h:    %d\n", m.stats.Next24hDue) +
			fmt.Sprintf("  Bookmarked:  %d (%d due)", m.stats.BookmarkedCards, m.stats.BookmarkedDue))

	collectionStats := dashCollectionStyle.
		Width(boxWidth).
		Render(dashStatsStyle.Render("Collection") + "\n" +
			fmt.Sprintf("  Decks:       %d (%d active)\n", m.stats.TotalDecks, m.stats.ActiveDecks) +
			fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards) +
			fmt.Sprintf("  Suspended:   %d", m.stats.SuspendedCards))

	totalKnown := m.stats.NewCards + m.stats.YoungCards + m.stats.MatureCards
	mixWidth := maxInt(10, layout.Width/2-16)
	if isNarrow {
		mixWidth = maxInt(10, layout.Width-20)
	}
	newPct, youngPct, maturePct := 0.0, 0.0, 0.0
	if totalKnown > 0 {
		newPct = float64(m.stats.NewCards) / float64(totalKnown)
		youngPct = float64(m.stats.YoungCards) / float64(totalKnown)
		maturePct = float64(m.stats.MatureCards) / float64(totalKnown)
	}
	dueMixBox := dashMixStyle.
		Width(boxWidth).
		Render(lipgloss.NewStyle().Foreground(colorSecondary).Bold(true).Render("Card Mix") + "\n" +
			fmt.Sprintf("  New    %3d %s\n", m.stats.NewCards, progressBar(mixWidth, newPct, "81", "238")) +
			fmt.Sprintf("  Young  %3d %s\n", m.stats.YoungCards, progressBar(mixWidth, youngPct, "220", "238")) +
			fmt.Sprintf("  Mature %3d %s", m.stats.MatureCards, progressBar(mixWidth, maturePct, "46", "238")))

	percentage := 0.0
	barColor := "46" // Green by default
	if m.stats.DailyGoal > 0 {
		percentage = float64(m.stats.ReviewsToday) / float64(m.stats.DailyGoal)
		if percentage >= 1.0 {
			barColor = "220" // Gold when goal is met
		}
	}
	goalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor)).Bold(true)
	bar := progressBar(maxInt(10, layout.Width/2-10), percentage, barColor, "238")
	if isNarrow {
		bar = progressBar(maxInt(10, layout.Width-14), percentage, barColor, "238")
	}

	progressBox := dashProgressStyle.
		BorderForeground(lipgloss.Color(barColor)).
		Width(boxWidth).
		Render(goalStyle.Render("Today's Progress") + "\n" +
			fmt.Sprintf("  Reviews:     %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal) +
			"  " + bar + "\n" +
			fmt.Sprintf("  Streak:      %d days%s", m.stats.CurrentStreak, streakIndicator))

	// Recent Activity Sparkline
	recentData := make([]int, 14)
	now := time.Now()
	for i := 0; i < 14; i++ {
		date := now.AddDate(0, 0, -(13 - i)).Format("2006-01-02")
		recentData[i] = m.reviewsPerDay[date]
	}
	sparkWidth := maxInt(10, layout.Width/2-10)
	if isNarrow {
		sparkWidth = maxInt(10, layout.Width-10)
	}
	spark := sparkline(recentData, sparkWidth)
	activityBox := dashActivityStyle.
		Width(boxWidth).
		Render(lipgloss.NewStyle().Foreground(colorAITitle).Bold(true).Render("Recent Activity") + "\n" +
			"  " + spark + "\n" +
			"  Last 14 days review trend")

	digestStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	message := "All caught up!"
	nextPreview := ""
	if len(m.dueCards) > 0 {
		message = fmt.Sprintf("%d due today.", len(m.dueCards))
		nextCard := m.dueCards[0]
		nextPreview = "\n  Next: " + truncateLine(nextCard.Prompt, 15)
	}
	if m.stats.CurrentStreak == 0 && len(m.dueCards) > 0 {
		message = fmt.Sprintf("%d cards waiting.", len(m.dueCards))
	}

	motivation := ""
	if m.stats.DailyGoal > 0 && m.stats.ReviewsToday >= m.stats.DailyGoal {
		motivation = "  " + successStyle.Render("Daily Goal Met! Great work!")
	} else if m.stats.CurrentStreak >= 30 {
		motivation = "  " + successStyle.Render("30+ day streak! Incredible!")
	} else if m.stats.CurrentStreak >= 7 {
		motivation = "  📈 One week streak! Keep it up!"
	} else if m.sessionReviewed > 0 {
		motivation = fmt.Sprintf("  ✅ Great session! %d cards today.", m.sessionReviewed)
	} else {
		motivation = "  💪 Ready for your daily German practice?"
	}

	dailyDigestBox := dashDigestStyle.
		Width(boxWidth).
		Render(digestStyle.Render("Daily Digest") + "\n" +
			fmt.Sprintf("  %s", message) + nextPreview + "\n" +
			fmt.Sprintf("  M:%d Y:%d N:%d", m.stats.MatureCards, m.stats.YoungCards, m.stats.NewCards) +
			motivation)

	var db strings.Builder
	now = time.Now()
	dateStr := now.Format("Monday, 02. January 2006")
	progressText := fmt.Sprintf("Daily Goal: %d/%d reviews", m.stats.ReviewsToday, m.stats.DailyGoal)
	goalMetBadge := ""
	if m.stats.DailyGoal > 0 && m.stats.ReviewsToday >= m.stats.DailyGoal {
		goalMetBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("220")). // Gold
			Bold(true).
			Padding(0, 1).
			MarginLeft(1).
			Render("GOAL MET 🏆")
	}
	db.WriteString(dashTitleStyle.Render("DASHBOARD - WILLKOMMEN!") + goalMetBadge + " | " + lipgloss.NewStyle().Foreground(colorGreen).Render(progressText) + " | " + mutedStyle.Render(dateStr) + "\n")

	if m.lastSessionReviewed > 0 {
		accuracy := 0.0
		if m.lastSessionCorrect > 0 {
			accuracy = float64(m.lastSessionCorrect) / float64(m.lastSessionReviewed) * 100
		}
		summaryStyle := lipgloss.NewStyle().Foreground(colorGreen).Italic(true)
		db.WriteString(summaryStyle.Render(fmt.Sprintf("Last Session: %d cards, %.1f%% accuracy", m.lastSessionReviewed, accuracy)) + "\n")
	}

	activeDeckText := lipgloss.NewStyle().Foreground(colorCyan).Render("Active Deck: " + m.deckLabel())
	db.WriteString(activeDeckText + "\n")

	queueY := strings.Count(db.String(), "\n")
	if isNarrow {
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-review", View: ViewDashboard, X: layout.X, Y: layout.Y + queueY, Width: lipgloss.Width(reviewQueue), Height: lipgloss.Height(reviewQueue)})
		db.WriteString(reviewQueue + "\n")
		queueY = strings.Count(db.String(), "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-collection", View: ViewDashboard, X: layout.X, Y: layout.Y + queueY, Width: lipgloss.Width(collectionStats), Height: lipgloss.Height(collectionStats)})
		db.WriteString(collectionStats + "\n")
	} else {
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-review", View: ViewDashboard, X: layout.X, Y: layout.Y + queueY, Width: lipgloss.Width(reviewQueue), Height: lipgloss.Height(reviewQueue)})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-collection", View: ViewDashboard, X: layout.X + lipgloss.Width(reviewQueue) + 1, Y: layout.Y + queueY, Width: lipgloss.Width(collectionStats), Height: lipgloss.Height(collectionStats)})
		db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, reviewQueue, " ", collectionStats) + "\n")
	}

	progressY := strings.Count(db.String(), "\n")
	if isNarrow {
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-progress", View: ViewDashboard, X: layout.X, Y: layout.Y + progressY, Width: lipgloss.Width(progressBox), Height: lipgloss.Height(progressBox)})
		db.WriteString(progressBox + "\n")
		progressY = strings.Count(db.String(), "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-digest", View: ViewDashboard, X: layout.X, Y: layout.Y + progressY, Width: lipgloss.Width(dailyDigestBox), Height: lipgloss.Height(dailyDigestBox)})
		db.WriteString(dailyDigestBox + "\n")
	} else {
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-progress", View: ViewDashboard, X: layout.X, Y: layout.Y + progressY, Width: lipgloss.Width(progressBox), Height: lipgloss.Height(progressBox)})
		m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-digest", View: ViewDashboard, X: layout.X + lipgloss.Width(progressBox) + 1, Y: layout.Y + progressY, Width: lipgloss.Width(dailyDigestBox), Height: lipgloss.Height(dailyDigestBox)})
		db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, progressBox, " ", dailyDigestBox) + "\n")
	}

	activityY := strings.Count(db.String(), "\n")
	recentDecksBox := ""
	if layout.Height > 22 && len(m.recentDecks) > 0 {
		recentDecksContent := lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render("Recently Studied") + "\n"
		for i, deckID := range m.recentDecks {
			if i >= 3 {
				break
			}
			name := deckID
			for _, d := range m.decks {
				if d.ID == deckID {
					name = d.Name
					break
				}
			}

			hint := ""
			switch i {
			case 0:
				hint = keyStyle.Render("!") + " "
			case 1:
				hint = keyStyle.Render("@") + " "
			case 2:
				hint = keyStyle.Render("#") + " "
			}

			deckName := truncateLine(name, maxInt(20, boxWidth-6))
			recentDecksContent += fmt.Sprintf("  %s• %s\n", hint, deckName)

			// Note: hitboxes for recent decks might be tricky if we don't know the final Y.
			// For now, we'll skip precise hitboxes in narrow mode for these, or recalculate.
		}
		recentDecksBox = dashRecentStyle.
			Width(boxWidth).
			Render(recentDecksContent)
	}

	activityY = strings.Count(db.String(), "\n")
	if recentDecksBox != "" {
		if isNarrow {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-activity", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(activityBox), Height: lipgloss.Height(activityBox)})
			db.WriteString(activityBox + "\n")
			activityY = strings.Count(db.String(), "\n")
			// Individual deck hitboxes
			for i := 0; i < len(m.recentDecks) && i < 3; i++ {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("dash-recent-%d", i),
					View:   ViewDashboard,
					X:      layout.X + 2,
					Y:      layout.Y + activityY + 2 + i,
					Width:  boxWidth - 4,
					Height: 1,
				})
			}
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-recent-decks", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(recentDecksBox), Height: lipgloss.Height(recentDecksBox)})
			db.WriteString(recentDecksBox + "\n")
		} else {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-activity", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(activityBox), Height: lipgloss.Height(activityBox)})
			recentX := layout.X + lipgloss.Width(activityBox) + 1
			// Individual deck hitboxes
			for i := 0; i < len(m.recentDecks) && i < 3; i++ {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("dash-recent-%d", i),
					View:   ViewDashboard,
					X:      recentX + 2,
					Y:      layout.Y + activityY + 2 + i,
					Width:  boxWidth - 4,
					Height: 1,
				})
			}
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-recent-decks", View: ViewDashboard, X: recentX, Y: layout.Y + activityY, Width: lipgloss.Width(recentDecksBox), Height: lipgloss.Height(recentDecksBox)})
			db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, activityBox, " ", recentDecksBox) + "\n")
		}
	} else if layout.Height > 26 {
		if isNarrow {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-activity", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(activityBox), Height: lipgloss.Height(activityBox)})
			db.WriteString(activityBox + "\n")
			activityY = strings.Count(db.String(), "\n")
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-card-mix", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(dueMixBox), Height: lipgloss.Height(dueMixBox)})
			db.WriteString(dueMixBox + "\n")
		} else {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-activity", View: ViewDashboard, X: layout.X, Y: layout.Y + activityY, Width: lipgloss.Width(activityBox), Height: lipgloss.Height(activityBox)})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-card-mix", View: ViewDashboard, X: layout.X + lipgloss.Width(activityBox) + 1, Y: layout.Y + activityY, Width: lipgloss.Width(dueMixBox), Height: lipgloss.Height(dueMixBox)})
			db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, activityBox, " ", dueMixBox) + "\n")
		}
	} else {
		// Even more compact row
		activityBox = dashActivityStyle.
			Width(layout.Width - 4).
			Render(lipgloss.NewStyle().Foreground(colorAITitle).Bold(true).Render("Recent Activity") + "  " + spark)
		db.WriteString(activityBox + "\n")
	}

	usedHeight := strings.Count(db.String(), "\n") + 2
	remainingHeight := layout.Height - usedHeight

	// Compact Quick Actions
	if remainingHeight >= 2 {
		label := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("  Quick Actions:")
		db.WriteString(label + " ")

		currentY := strings.Count(db.String(), "\n")
		currentX := lipgloss.Width(label) + 1

		actions := []struct {
			id    string
			label string
			key   string
		}{
			{"nav-review", "Review", "3"},
			{"nav-practice", "Practice", "0"},
			{"nav-cram", "Cram", "9"},
			{"nav-decks", "Decks", "2"},
			{"nav-browser", "Browser", "8"},
			{"nav-statistics", "Stats", "4"},
			{"nav-import", "Import", "5"},
			{"nav-ai", "AI Draft", "6"},
			{"nav-settings", "Settings", "7"},
		}

		for _, action := range actions {
			keyStr := keyStyle.Render("[" + action.key + "]")
			item := fmt.Sprintf("%s %s  ", keyStr, action.label)
			itemWidth := lipgloss.Width(item)

			if currentX+itemWidth > layout.Width-2 {
				db.WriteString("\n  ")
				currentY++
				currentX = 2
				remainingHeight--
			}

			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     action.id,
				View:   ViewDashboard,
				X:      layout.X + currentX,
				Y:      layout.Y + currentY,
				Width:  itemWidth - 2, // -2 for trailing spaces
				Height: 1,
			})
			db.WriteString(item)
			currentX += itemWidth
		}
		db.WriteString("\n")
		remainingHeight -= 1
	}

	if remainingHeight >= 3 {
		tip := content.GetDailyGrammarTip()
		verb := content.GetVerbOfTheDay()

		tipLabelStyle := lipgloss.NewStyle().Foreground(colorAITitle).Bold(true)
		verbLabelStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)

		exampleText := ""
		if tip.Example != "" && remainingHeight >= 5 {
			exampleText = fmt.Sprintf("\n  Example: %s", lipgloss.NewStyle().Italic(true).Foreground(colorBlue).Render(tip.Example))
		}

		tipBox := dashTipStyle.
			Width(boxWidth).
			Render(tipLabelStyle.Render("Grammar Tip: "+tip.Title) + "\n" +
				fmt.Sprintf("  %s", tip.Tip) + exampleText)

		verbHeader := verbLabelStyle.Render("Verb: "+verb.German) +
			lipgloss.NewStyle().Foreground(colorMuted).Render(" — "+verb.English)
		verbBox := dashVerbStyle.
			Width(boxWidth).
			Render(verbHeader + "\n" +
				fmt.Sprintf("  ich %-8s wir %-8s\n  du  %-8s ihr %-8s\n  er/sie/es %-4s sie/Sie %-5s",
					verb.Ich, verb.Wir, verb.Du, verb.Ihr, verb.ErSieEs, verb.SieSie))

		if isNarrow {
			db.WriteString(tipBox + "\n" + verbBox + "\n")
		} else {
			db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tipBox, " ", verbBox) + "\n")
		}
	}

	db.WriteString(mutedStyle.Render(fmt.Sprintf("Use %s and %s to switch decks.\nUse %s (%s) to start studying.\nPress %s for help.",
		lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("["),
		lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("]"),
		lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("Review"),
		lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("3"),
		lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("?"))))

	return db.String()
}
