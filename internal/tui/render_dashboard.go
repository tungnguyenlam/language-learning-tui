package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/content"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderDashboard(layout viewportLayout) string {
	streakIndicator := ""
	if m.stats.CurrentStreak > 0 {
		streakIndicator = " 🔥"
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(0, 1).
		Width(maxInt(40, m.width-60)).
		Render(fmt.Sprintf("Active Deck: %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Render(m.deckLabel())))
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)

	reviewQueue := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(0, 1).
		Width(maxInt(25, (layout.Width-2)/2)).
		Render(statsStyle.Render("Review Queue") + "\n" +
			fmt.Sprintf("  Due cards:   %d\n", len(m.dueCards)) +
			fmt.Sprintf("  Bookmarked:  %d (%d due)", m.stats.BookmarkedCards, m.stats.BookmarkedDue))

	collectionStats := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(0, 1).
		Width(maxInt(25, (layout.Width-2)/2)).
		Render(statsStyle.Render("Collection") + "\n" +
			fmt.Sprintf("  Decks:       %d (%d active)\n", m.stats.TotalDecks, m.stats.ActiveDecks) +
			fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards) +
			fmt.Sprintf("  Suspended:   %d", m.stats.SuspendedCards))

	goalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	percentage := 0.0
	if m.stats.DailyGoal > 0 {
		percentage = float64(m.stats.ReviewsToday) / float64(m.stats.DailyGoal)
	}
	bar := progressBar(maxInt(10, layout.Width/2-10), percentage, "46", "238")

	progressBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(0, 1).
		Width(maxInt(25, (layout.Width-2)/2)).
		Render(goalStyle.Render("Today's Progress") + "\n" +
			fmt.Sprintf("  Reviews:     %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal) +
			"  " + bar + "\n" +
			fmt.Sprintf("  Streak:      %d days%s", m.stats.CurrentStreak, streakIndicator))

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
	if m.stats.CurrentStreak >= 30 {
		motivation = "  Incredible dedication! 30+ day streak!"
	} else if m.stats.CurrentStreak >= 14 {
		motivation = "  Two weeks strong! Keep it up!"
	} else if m.stats.CurrentStreak >= 7 {
		motivation = "  One week streak! You're building a habit!"
	} else if m.stats.CurrentStreak >= 3 {
		motivation = "  Nice streak! Consistency is key!"
	} else if m.sessionReviewed > 0 {
		motivation = fmt.Sprintf("  Great session! %d cards reviewed today.", m.sessionReviewed)
	}

	dailyDigestBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("213")).
		Padding(0, 1).
		Width(maxInt(25, (layout.Width-2)/2)).
		Render(digestStyle.Render("Daily Digest") + "\n" +
			fmt.Sprintf("  %s", message) + nextPreview + "\n" +
			fmt.Sprintf("  M:%d Y:%d N:%d", m.stats.MatureCards, m.stats.YoungCards, m.stats.NewCards) +
			motivation)

	quickActionsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	quickActionsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(0, 1).
		Width(layout.Width - 4).
		Render(quickActionsStyle.Render("Quick Actions") + "\n" +
			fmt.Sprintf("  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s",
				keyStyle.Render("[3]"), "Review",
				keyStyle.Render("[9]"), "Cram",
				keyStyle.Render("[8]"), "Browser",
				keyStyle.Render("[4]"), "Stats",
				keyStyle.Render("[5]"), "Import",
				keyStyle.Render("[6]"), "AI Draft"))

	var db strings.Builder
	db.WriteString(titleStyle.Render("DASHBOARD") + "\n")

	if m.lastSessionReviewed > 0 {
		accuracy := 0.0
		if m.lastSessionReviewed > 0 {
			accuracy = float64(m.lastSessionCorrect) / float64(m.lastSessionReviewed) * 100
		}
		summaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Italic(true)
		db.WriteString(summaryStyle.Render(fmt.Sprintf("Last Session: %d cards, %.1f%% accuracy", m.lastSessionReviewed, accuracy)) + "\n")
	}

	db.WriteString(headerBox + "\n")

	queueY := strings.Count(db.String(), "\n")
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-review", View: ViewDashboard, X: layout.X, Y: layout.Y + queueY, Width: lipgloss.Width(reviewQueue), Height: lipgloss.Height(reviewQueue)})
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-collection", View: ViewDashboard, X: layout.X + lipgloss.Width(reviewQueue) + 1, Y: layout.Y + queueY, Width: lipgloss.Width(collectionStats), Height: lipgloss.Height(collectionStats)})
	db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, reviewQueue, " ", collectionStats) + "\n")

	progressY := strings.Count(db.String(), "\n")
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-progress", View: ViewDashboard, X: layout.X, Y: layout.Y + progressY, Width: lipgloss.Width(progressBox), Height: lipgloss.Height(progressBox)})
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-digest", View: ViewDashboard, X: layout.X + lipgloss.Width(progressBox) + 1, Y: layout.Y + progressY, Width: lipgloss.Width(dailyDigestBox), Height: lipgloss.Height(dailyDigestBox)})
	db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, progressBox, " ", dailyDigestBox) + "\n")

	if layout.Height > 24 {
		db.WriteString(quickActionsBox + "\n")
	}

	// Only show grammar tip if there is enough vertical space
	if layout.Height > 28 {
		tip := content.GetDailyGrammarTip()
		tipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
		tipBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1).
			Width(layout.Width - 4). // Use layout width
			Render(tipStyle.Render("Grammar Tip: "+tip.Title) + "\n" +
				fmt.Sprintf("  %s", tip.Tip))
		db.WriteString(tipBox + "\n")
	}

	db.WriteString(mutedStyle.Render(fmt.Sprintf("Use %s and %s to switch decks.\nUse %s (%s) to start studying.\nPress %s for help.",
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("["),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("]"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("Review"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("3"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("?"))))

	return db.String()
}
