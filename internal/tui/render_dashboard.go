package tui

import (
	"fmt"
	"image/color"
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
		if m.lastSessionReviewed > 0 {
			accuracy = float64(m.lastSessionCorrect) / float64(m.lastSessionReviewed) * 100
		}
		speedStr := ""
		if m.lastSessionDuration.Minutes() > 0 {
			speed := float64(m.lastSessionReviewed) / m.lastSessionDuration.Minutes()
			speedStr = fmt.Sprintf(", %.1f cards/min", speed)
		}
		summaryStyle := lipgloss.NewStyle().Foreground(colorGreen).Italic(true)
		db.WriteString(summaryStyle.Render(fmt.Sprintf("Last Session: %d cards, %.1f%% accuracy%s", m.lastSessionReviewed, accuracy, speedStr)) + "\n")
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
	if remainingHeight >= 3 {
		actions := []struct {
			id       string
			label    string
			key      string
			category string
			color    color.Color
		}{
			{"nav-review", "Review", "3", "study", colorGreen},
			{"nav-practice", "Practice", "0", "study", colorGreen},
			{"nav-cram", "Cram", "9", "study", colorGreen},
			{"nav-dictionary", "Dictionary", "/", "tools", colorBlue},
			{"nav-decks", "Decks", "2", "tools", colorBlue},
			{"nav-browser", "Browser", "8", "tools", colorBlue},
			{"nav-statistics", "Stats", "4", "progress", colorGold},
			{"nav-ai", "AI Draft", "6", "progress", colorGold},
			{"nav-import", "Import", "5", "content", colorPurple},
			{"nav-settings", "Settings", "7", "content", colorPurple},
		}

		actionsBoxWidth := layout.Width - 2
		contentWidth := actionsBoxWidth - 6 // border(2) + padding(4)
		if contentWidth < 10 {
			contentWidth = 10
		}

		prefix := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("Quick Actions") + "  "
		prefixWidth := lipgloss.Width("Quick Actions  ")

		var lines []string
		var currentLine string
		var currentLineWidth int

		if isNarrow {
			lines = append(lines, lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("Quick Actions"))
			currentLine = ""
			currentLineWidth = 0
		} else {
			currentLine = prefix
			currentLineWidth = prefixWidth
		}

		type relativeHitbox struct {
			id    string
			relX  int
			relY  int
			width int
		}
		var relHitboxes []relativeHitbox

		lastCategory := ""
		for _, action := range actions {
			if action.category != lastCategory && lastCategory != "" {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = ""
				currentLineWidth = 0
			}
			lastCategory = action.category

			keyStr := lipgloss.NewStyle().Foreground(action.color).Bold(true).Render("[" + action.key + "]")
			labelStr := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(action.label)
			item := fmt.Sprintf("%s %s  ", keyStr, labelStr)
			itemWidth := lipgloss.Width(item)

			if currentLineWidth+itemWidth > contentWidth {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = ""
				currentLineWidth = 0
			}

			relHitboxes = append(relHitboxes, relativeHitbox{
				id:    action.id,
				relX:  currentLineWidth,
				relY:  len(lines),
				width: itemWidth - 2, // exclude trailing spaces
			})

			currentLine += item
			currentLineWidth += itemWidth
		}

		if currentLine != "" {
			lines = append(lines, currentLine)
		}

		boxContent := strings.Join(lines, "\n")
		actionsBox := dashActionsStyle.Width(actionsBoxWidth).Render(boxContent)

		startY := strings.Count(db.String(), "\n")
		db.WriteString(actionsBox + "\n")

		for _, rh := range relHitboxes {
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     rh.id,
				View:   ViewDashboard,
				X:      layout.X + rh.relX + 2,          // border(1) + padding(1)
				Y:      layout.Y + startY + rh.relY + 1, // border(1)
				Width:  rh.width,
				Height: 1,
			})
		}
		remainingHeight -= lipgloss.Height(actionsBox)
	}

	if remainingHeight >= 3 {
		tip := content.GetDailyGrammarTip()
		verb := content.GetVerbOfTheDay()
		word := content.GetWordOfTheDay()

		tipLabelStyle := lipgloss.NewStyle().Foreground(colorAITitle).Bold(true)
		verbLabelStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)
		wordLabelStyle := lipgloss.NewStyle().Foreground(colorGold).Bold(true)

		exampleText := ""
		if tip.Example != "" && remainingHeight >= 5 {
			exampleText = fmt.Sprintf("\n  Example: %s", lipgloss.NewStyle().Italic(true).Foreground(colorBlue).Render(tip.Example))
		}

		tipBox := dashTipStyle.
			Width(boxWidth).
			Render(tipLabelStyle.Render("Grammar Tip: [g/G] "+tip.Title) + "\n" +
				fmt.Sprintf("  %s", tip.Tip) + exampleText)

		verbHeader := verbLabelStyle.Render("Verb: [v/V] "+verb.German) +
			lipgloss.NewStyle().Foreground(colorMuted).Render(" — "+verb.English)
		verbBox := dashVerbStyle.
			Width(boxWidth).
			Render(verbHeader + "\n" +
				fmt.Sprintf("  ich %-8s wir %-8s\n  du  %-8s ihr %-8s\n  er/sie/es %-4s sie/Sie %-5s",
					verb.Ich, verb.Wir, verb.Du, verb.Ihr, verb.ErSieEs, verb.SieSie))

		wordHeader := wordLabelStyle.Render("Word: [w/W] "+word.German) +
			lipgloss.NewStyle().Foreground(colorMuted).Render(" — "+word.English)
		wordContent := ""
		if word.Plural != "" {
			wordContent += fmt.Sprintf("  Plural: %s\n", word.Plural)
		}
		if word.Example != "" {
			wordContent += fmt.Sprintf("  %s", lipgloss.NewStyle().Italic(true).Foreground(colorBlue).Render(word.Example))
		}
		wordBox := dashWordStyle.
			Width(boxWidth).
			Render(wordHeader + "\n" +
				wordContent)

		tipY := strings.Count(db.String(), "\n")
		if isNarrow {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-tip", View: ViewDashboard, X: layout.X, Y: layout.Y + tipY, Width: boxWidth, Height: lipgloss.Height(tipBox)})
			db.WriteString(tipBox + "\n")
			verbY := strings.Count(db.String(), "\n")
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-verb", View: ViewDashboard, X: layout.X, Y: layout.Y + verbY, Width: boxWidth, Height: lipgloss.Height(verbBox)})
			db.WriteString(verbBox + "\n")
			wordY := strings.Count(db.String(), "\n")
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-word", View: ViewDashboard, X: layout.X, Y: layout.Y + wordY, Width: boxWidth, Height: lipgloss.Height(wordBox)})
			db.WriteString(wordBox + "\n")
		} else if layout.Width > 110 && remainingHeight >= 5 {
			// Three column row for wide terminals
			thirdWidth := (layout.Width - 4) / 3
			tipBox = dashTipStyle.Width(thirdWidth).Render(tipLabelStyle.Render("Grammar: [g/G] "+tip.Title) + "\n" + tip.Tip)
			verbBox = dashVerbStyle.Width(thirdWidth).Render(verbHeader + "\n" + fmt.Sprintf("  ich %-8s wir %-8s\n  du  %-8s ihr %-8s", verb.Ich, verb.Wir, verb.Du, verb.Ihr))
			wordBox = dashWordStyle.Width(thirdWidth).Render(wordHeader + "\n" + wordContent)

			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-tip", View: ViewDashboard, X: layout.X, Y: layout.Y + tipY, Width: thirdWidth, Height: lipgloss.Height(tipBox)})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-verb", View: ViewDashboard, X: layout.X + thirdWidth + 1, Y: layout.Y + tipY, Width: thirdWidth, Height: lipgloss.Height(verbBox)})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-word", View: ViewDashboard, X: layout.X + 2*thirdWidth + 2, Y: layout.Y + tipY, Width: thirdWidth, Height: lipgloss.Height(wordBox)})
			db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tipBox, " ", verbBox, " ", wordBox) + "\n")
		} else if remainingHeight >= 7 {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-tip", View: ViewDashboard, X: layout.X, Y: layout.Y + tipY, Width: boxWidth, Height: lipgloss.Height(tipBox)})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-verb", View: ViewDashboard, X: layout.X + boxWidth + 1, Y: layout.Y + tipY, Width: boxWidth, Height: lipgloss.Height(verbBox)})
			db.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tipBox, " ", verbBox) + "\n")
			wordY := strings.Count(db.String(), "\n")
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-word", View: ViewDashboard, X: layout.X, Y: layout.Y + wordY, Width: boxWidth, Height: lipgloss.Height(wordBox)})
			db.WriteString(wordBox + "\n")
		} else {
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-tip", View: ViewDashboard, X: layout.X, Y: layout.Y + tipY, Width: boxWidth, Height: lipgloss.Height(tipBox)})
			m.hitboxes = append(m.hitboxes, Hitbox{ID: "dash-verb", View: ViewDashboard, X: layout.X + boxWidth + 1, Y: layout.Y + tipY, Width: boxWidth, Height: lipgloss.Height(verbBox)})
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
