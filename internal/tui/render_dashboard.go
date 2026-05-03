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
		Render(fmt.Sprintf("Active Deck: %s\nDeck: %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Render(m.deckLabel()),
			lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Render(m.deckLabel())))
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)

	reviewQueue := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(0, 1).
		Render(statsStyle.Render("Review Queue") + "\n" +
			fmt.Sprintf("  Due cards:   %d\n", len(m.dueCards)) +
			fmt.Sprintf("  Bookmarked:  %d (%d due)", m.stats.BookmarkedCards, m.stats.BookmarkedDue))
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-review", View: ViewDashboard, X: layout.X, Y: layout.Y + 4, Width: 26, Height: 5})

	collectionStats := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(0, 1).
		Render(statsStyle.Render("Collection") + "\n" +
			fmt.Sprintf("  Decks:       %d (%d active)\n", m.stats.TotalDecks, m.stats.ActiveDecks) +
			fmt.Sprintf("  Leech:       %d\n", m.stats.LeechCards) +
			fmt.Sprintf("  Suspended:   %d", m.stats.SuspendedCards))
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-collection", View: ViewDashboard, X: layout.X + 28, Y: layout.Y + 4, Width: 26, Height: 5})

	goalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	progressBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(0, 1).
		Render(goalStyle.Render("Today's Progress") + "\n" +
			fmt.Sprintf("  Reviews:     %d/%d\n", m.stats.ReviewsToday, m.stats.DailyGoal) +
			fmt.Sprintf("  Streak:      %d days%s", m.stats.CurrentStreak, streakIndicator))
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-progress", View: ViewDashboard, X: layout.X, Y: layout.Y + 9, Width: 26, Height: 5})

	digestStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	message := "All caught up!"
	if len(m.dueCards) > 0 {
		message = fmt.Sprintf("%d due today.", len(m.dueCards))
	}
	if m.stats.CurrentStreak == 0 && len(m.dueCards) > 0 {
		message = fmt.Sprintf("%d cards waiting.", len(m.dueCards))
	}

	dailyDigestBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("213")).
		Padding(0, 1).
		Render(digestStyle.Render("Daily Digest") + "\n" +
			fmt.Sprintf("  %s\n", message) +
			fmt.Sprintf("  M:%d Y:%d N:%d", m.stats.MatureCards, m.stats.YoungCards, m.stats.NewCards))
	m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-digest", View: ViewDashboard, X: layout.X + 28, Y: layout.Y + 9, Width: 26, Height: 5})

	var db strings.Builder
	db.WriteString(titleStyle.Render("DASHBOARD") + "\n")
	db.WriteString(headerBox + "\n")

	db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, reviewQueue, " ", collectionStats) + "\n")
	db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, progressBox, " ", dailyDigestBox) + "\n")

	// Only show grammar tip if there is enough vertical space
	if layout.Height > 20 {
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

	db.WriteString(mutedStyle.Render("Use [ and ] to switch decks.\nUse Review (3) to start studying."))

	return db.String()
}
