package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderHelp(layout viewportLayout) string {
	ctx := NewRenderContext(m, layout, m.activeView)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Underline(true)
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	colStyle := lipgloss.NewStyle().PaddingRight(4).Width(maxInt(32, layout.Width/4))

	ctx.WriteLine(titleStyle.Render("Keyboard Shortcuts"))
	ctx.NewLine()

	// Simplified: just render it as before, but maybe we can make some headers clickable to switch views
	global := sectionStyle.Render("Global:") + "\n" +
		"  1-9, 0   Switch to view\n" +
		"  Tab/arr  Cycle views\n" +
		"  w/s      Prev/next view\n" +
		"  =        Dictionary lookup\n" +
		"  ?        Toggle help\n" +
		"  q/Ctrl+c Quit"

	dash := sectionStyle.Render("Dashboard/Decks:") + "\n" +
		"  [ ]      Prev/next deck\n" +
		"  /        Search decks\n" +
		"  L        Edit deck limits\n" +
		"  +/-      Adjust limits\n" +
		"  !/@/#    Recent decks\n" +
		"  Enter    Select deck"

	review := sectionStyle.Render("Review:") + "\n" +
		"  Spc/Ent  Reveal answer\n" +
		"  a/h/g/e  Again/Hard/Good/Easy\n" +
		"  1-4      Grade/MCQ choice\n" +
		"  b / B    Bookmark / Filter\n" +
		"  u / r    Undo / History\n" +
		"  f / i    Focus / Info\n" +
		"  p / d    Play audio / Dict\n" +
		"  t / x    Type / Suspend\n" +
		"  h / H    Hint / AI Explain\n" +
		"  F        Fix card"

	browser := sectionStyle.Render("Browser:") + "\n" +
		"  j/k      Navigate\n" +
		"  / / #    Search/Tag filter\n" +
		"  m / a    Select / Select all\n" +
		"  b / B    Bookmark / Unmark\n" +
		"  x / X    Suspend / Unsuspend\n" +
		"  t / T    Toggle kind/Tags\n" +
		"  p / d    Play audio / Dict\n" +
		"  Ent/Del  History / Delete\n" +
		"  C        Cleanup tags"

	dictionary := sectionStyle.Render("Dictionary:") + "\n" +
		"  j/k      Navigate\n" +
		"  PgUp/Dn  Fast scroll\n" +
		"  Enter    Draft AI card\n" +
		"  ctrl+a   Quick add to deck\n" +
		"  ctrl+f   Find in Decks\n" +
		"  ctrl+p   Play audio\n" +
		"  Esc      Back"

	other := sectionStyle.Render("Other:") + "\n" +
		"  Stats    j/k scroll, x exp\n" +
		"  AI       / topic, a/d draft\n" +
		"  Cram     Ent start, 1-5 filter\n" +
		"  Practice 1-9/0/-/= pick, r reset\n" +
		"  Import   i/I import, x/X exp\n" +
		"  AnkiWeb  A browse shared\n" +
		"  Settings j/k nav, +/- goal"

	col1 := colStyle.Render(global + "\n\n" + dash)
	col2 := colStyle.Render(review)
	col3 := colStyle.Render(browser + "\n\n" + dictionary)
	col4 := colStyle.Render(other)

	helpContent := lipgloss.JoinHorizontal(lipgloss.Top, col1, col2, col3, col4)

	return titleStyle.Render("Keyboard Shortcuts") + "\n\n" + helpContent
}

func (m *Model) renderDecks(layout viewportLayout) string {
	var b strings.Builder
	b.WriteString(dashTitleStyle.Render("DECK LIST") + "\n\n")

	// Show filter if active or searching
	if m.searchingDecks || m.deckFilter != "" {
		searchBarWidth := layout.Width - 4
		if searchBarWidth < 20 {
			searchBarWidth = 20
		}
		searchBar := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(searchBarWidth)

		if m.searchingDecks {
			displaySearch := m.deckFilter
			clearText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true).Render("[x]")

			// Calculate how many spaces to put between search and clear button
			contentLength := 8 + len([]rune(displaySearch)) + 1 // "Search: " + query + "_"
			spaces := searchBarWidth - contentLength - 3
			var clearBtn string
			if spaces > 0 {
				clearBtn = strings.Repeat(" ", spaces) + clearText
			} else {
				clearBtn = " " + clearText
			}

			searchBarText := searchBar.Render("Search: " + displaySearch + "_" + clearBtn)

			// Register hitbox for clear button [x]
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     "deck-search-clear",
				View:   ViewDecks,
				X:      layout.X + 2 + searchBarWidth - 3,
				Y:      layout.Y + strings.Count(b.String(), "\n") + 1, // +1 for the top border
				Width:  3,
				Height: 1,
				Action: func() tea.Cmd {
					m.deckFilter = ""
					return nil
				},
			})
			b.WriteString(searchBarText + "\n\n")

			if m.deckFilter == "" && len(m.deckSearchHistory) > 0 {
				clearHistoryText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("[Clear]")
				b.WriteString(boldStyle.Render("Recent Searches:") + "  " + clearHistoryText + " " + mutedStyle.Render("(ctrl+x)") + "\n")

				// Hitbox for "[Clear]" button
				clearHistoryY := layout.Y + strings.Count(b.String(), "\n") - 1
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "deck-history-clear",
					View:   ViewDecks,
					X:      layout.X + 18,
					Y:      clearHistoryY,
					Width:  7,
					Height: 1,
					Action: func() tea.Cmd {
						m.deckSearchHistory = nil
						m.saveDeckHistory()
						return nil
					},
				})

				for i := len(m.deckSearchHistory) - 1; i >= 0; i-- {
					q := m.deckSearchHistory[i]
					lineY := layout.Y + strings.Count(b.String(), "\n")
					b.WriteString(fmt.Sprintf("  • %s\n", q))
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("deck-history-%d", i),
						View:   ViewDecks,
						X:      layout.X + 4,
						Y:      lineY,
						Width:  lipgloss.Width(q),
						Height: 1,
						Action: func() tea.Cmd {
							m.deckFilter = q
							m.applyDeckFilter()
							return nil
						},
					})
				}
				b.WriteString("\n")
			}
		} else {
			b.WriteString(searchBar.Render(fmt.Sprintf("Filter: %s (Press / to edit)", m.deckFilter)) + "\n\n")
		}
	}

	filteredDecks := m.filteredDecks()
	if len(filteredDecks) == 0 {
		if m.deckFilter != "" {
			b.WriteString("No decks match search. Press Esc to clear.\n")
		} else {
			b.WriteString("No decks found. Use Import to add notes.\n")
		}
		if m.searchingDecks {
			b.WriteString("\nPress Enter or Esc to finish.")
		} else {
			b.WriteString("\nPress Esc to clear filter.")
		}
		return b.String()
	}

	// Constants for width calculation
	lineWidth := layout.Width - 2
	studyBtnWidth := 0
	cramBtnWidth := 0
	if layout.Width >= 100 {
		studyBtnWidth = 8
		cramBtnWidth = 7
	} else if layout.Width >= 90 {
		studyBtnWidth = 8
	}

	countsWidth := 16
	if layout.Width >= 100 {
		countsWidth = 32
	}

	statsWidth := 0
	if layout.Width >= 85 {
		statsWidth = 26
	} else if layout.Width >= 70 {
		statsWidth = 12
	}

	prefixWidth := 2
	selectMarkWidth := 4
	miniBarWidth := 5

	otherWidth := prefixWidth + selectMarkWidth + 1 + miniBarWidth + 2 + countsWidth + statsWidth + studyBtnWidth + cramBtnWidth
	nameWidth := lineWidth - otherWidth

	if nameWidth < 15 {
		// Drop stats first
		if statsWidth > 0 {
			statsWidth = 0
		}
		// Drop buttons
		if studyBtnWidth > 0 {
			studyBtnWidth = 0
		}
		if cramBtnWidth > 0 {
			cramBtnWidth = 0
		}

		otherWidth = prefixWidth + selectMarkWidth + 1 + miniBarWidth + 2 + countsWidth + statsWidth + studyBtnWidth + cramBtnWidth
		nameWidth = lineWidth - otherWidth

		if nameWidth < 15 && countsWidth > 12 {
			countsWidth = 12
			otherWidth = prefixWidth + selectMarkWidth + 1 + miniBarWidth + 2 + countsWidth + statsWidth + studyBtnWidth + cramBtnWidth
			nameWidth = lineWidth - otherWidth
		}

		// If still too narrow, drop counts
		if nameWidth < 15 && countsWidth > 0 {
			countsWidth = 0
			otherWidth = prefixWidth + selectMarkWidth + 1 + miniBarWidth + 2 + countsWidth + statsWidth + studyBtnWidth + cramBtnWidth
			nameWidth = lineWidth - otherWidth
		}

		// If still too narrow, drop miniBar
		if nameWidth < 15 && miniBarWidth > 0 {
			miniBarWidth = 0
			otherWidth = prefixWidth + selectMarkWidth + 1 + miniBarWidth + 2 + countsWidth + statsWidth + studyBtnWidth + cramBtnWidth
			nameWidth = lineWidth - otherWidth
		}

		if nameWidth < 10 {
			nameWidth = 10
		}
	}

	type lineInfo struct {
		deckIdx int
		kind    string // "main", "limits", "desc", "tags"
		label   string // raw label before buttons
		study   string
		cram    string
	}

	var content strings.Builder
	var lines []lineInfo

	for i, deck := range filteredDecks {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.deckCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(colorPink)
		}

		selectMark := "[ ] "
		if m.deckSelected[deck.ID] {
			selectMark = "[x] "
		}

		newStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render(fmt.Sprintf("%d", deck.NewCards))
		dueStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(fmt.Sprintf("%d", deck.DueCards))
		totalStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render(fmt.Sprintf("%d", deck.TotalCards))

		deckPercentage := 0.0
		if deck.TotalCards > 0 {
			deckPercentage = float64(deck.TotalCards-deck.DueCards) / float64(deck.TotalCards)
		}

		miniBar := ""
		if miniBarWidth > 0 {
			miniBar = progressBar(5, deckPercentage, lipgloss.Color("46"), lipgloss.Color("238")) + " "
		}

		studyBtn := ""
		cramBtn := ""
		if studyBtnWidth > 0 {
			studyBtn = " [Study]"
		}
		if cramBtnWidth > 0 {
			cramBtn = " [Cram]"
		}

		var counts string
		if countsWidth > 0 {
			if layout.Width >= 115 {
				counts = fmt.Sprintf("%s new, %s due, %s total", newStyled, dueStyled, totalStyled)
			} else {
				counts = fmt.Sprintf("%sN %sD %sT", newStyled, dueStyled, totalStyled)
			}
			counts = padLine(counts, countsWidth) + " "
		}

		statsStr := ""
		if statsWidth > 0 {
			if layout.Width >= 85 {
				statsStr = fmt.Sprintf(" | today %d, %.0f%% success", deck.ReviewsToday, deck.SuccessRate*100)
			} else if layout.Width >= 70 {
				statsStr = fmt.Sprintf(" | today %d", deck.ReviewsToday)
			}
			statsStr = padLine(statsStr, statsWidth)
		}

		nameText := truncateLine(deck.Name, nameWidth)
		resDeckName := nameText
		/*
			if m.deckFilter != "" {
				resDeckName = highlightQuery(resDeckName, m.deckFilter, highlightStyle)
			}
		*/
		resDeckName = padLine(resDeckName, nameWidth)

		label := fmt.Sprintf("%s%s%s %s%s%s",
			prefix, selectMark, resDeckName, miniBar, counts, statsStr)

		currBtnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if i == m.deckCursor {
			currBtnStyle = currBtnStyle.Foreground(colorBlue)
		}

		// Main line
		lines = append(lines, lineInfo{deckIdx: i, kind: "main", label: label, study: currBtnStyle.Render(studyBtn), cram: currBtnStyle.Render(cramBtn)})
		content.WriteString(label + currBtnStyle.Render(studyBtn) + currBtnStyle.Render(cramBtn) + "\n")

		if i == m.deckCursor && m.editingDeckLimits {
			limitsLabel := fmt.Sprintf("     Limits: New %d, Review %d (h/l switch, +/- adjust)",
				deck.NewCardsPerDay, deck.ReviewLimitPerDay)
			lines = append(lines, lineInfo{deckIdx: i, kind: "limits"})
			content.WriteString(limitsLabel + "\n")
		}

		if deck.Description != "" {
			desc := truncateLine(deck.Description, layout.Width-10)
			lines = append(lines, lineInfo{deckIdx: i, kind: "desc"})
			content.WriteString(fmt.Sprintf("     %s\n", mutedStyle.Render(desc)))
		}
		if len(deck.Tags) > 0 {
			tags := strings.Join(deck.Tags, ", ")
			resTags := truncateLine(tags, layout.Width-15)
			/*
				if m.deckFilter != "" {
					resTags = highlightQuery(resTags, m.deckFilter, highlightStyle)
				}
			*/
			resTagsStyled := mutedStyle.Render(resTags)
			lines = append(lines, lineInfo{deckIdx: i, kind: "tags"})
			content.WriteString(fmt.Sprintf("     Tags: %s\n", resTagsStyled))
		}
		lines = append(lines, lineInfo{deckIdx: i, kind: "spacer"})
		content.WriteString("\n")
	}

	availableHeight := layout.Height - strings.Count(b.String(), "\n") - 3
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Auto-scroll to cursor
	firstLineOfCursor := -1
	for idx, li := range lines {
		if li.deckIdx == m.deckCursor {
			firstLineOfCursor = idx
			break
		}
	}
	if firstLineOfCursor != -1 {
		if firstLineOfCursor < m.deckScroll {
			m.deckScroll = firstLineOfCursor
		} else if firstLineOfCursor >= m.deckScroll+availableHeight-2 {
			m.deckScroll = firstLineOfCursor - availableHeight + 3
		}
	}
	m.deckTotalLines = len(lines)
	m.deckScroll = clampInt(m.deckScroll, 0, maxInt(0, m.deckTotalLines-availableHeight))

	footer := ""
	if m.searchingDecks {
		footer = "Press enter or Esc to finish searching."
	} else {
		selectedCount := 0
		for _, s := range m.deckSelected {
			if s {
				selectedCount++
			}
		}
		if selectedCount > 0 {
			footer = fmt.Sprintf("%d decks selected. Press %s to delete, %s to merge into current.", selectedCount,
				keyStyle.Render("Backspace"), keyStyle.Render("M"))
		} else {
			footer = fmt.Sprintf("%s select | %s search | %s stats | %s/%s multi-select | %s clear\nPress enter to select deck.",
				keyStyle.Render("enter"), keyStyle.Render("/"), keyStyle.Render("v"), keyStyle.Render("m"), keyStyle.Render("x"), keyStyle.Render("Esc"))
		}
	}

	listView := m.RenderList(layout.WithHeight(availableHeight).WithY(layout.Y+strings.Count(b.String(), "\n")), content.String(), ListOptions{
		HitboxPrefix: "deck",
		View:         ViewDecks,
		Footer:       footer,
		ScrollOffset: &m.deckScroll,
		TotalLines:   &m.deckTotalLines,
		OnLine: func(lineIndex int, lineCtx *RenderContext, content string) {
			if lineIndex < 0 || lineIndex >= len(lines) {
				return
			}
			li := lines[lineIndex]
			if li.kind == "main" {
				deck := filteredDecks[li.deckIdx]
				lineCtx.RegisterHitboxWithAction(fmt.Sprintf("deck-select-%d", li.deckIdx), lipgloss.Width(li.label), 1, func() tea.Cmd {
					m.selectDeckByID(deck.ID)
					return m.updateView(ViewDashboard)
				})
				if li.study != "" {
					lineCtx.RegisterHitboxAtWithAction(fmt.Sprintf("deck-study-%d", li.deckIdx), lipgloss.Width(li.label), 0, lipgloss.Width(li.study), 1, func() tea.Cmd {
						m.selectDeckByID(deck.ID)
						return m.updateView(ViewReview)
					})
				}
				if li.cram != "" {
					lineCtx.RegisterHitboxAtWithAction(fmt.Sprintf("deck-cram-%d", li.deckIdx), lipgloss.Width(li.label)+lipgloss.Width(li.study), 0, lipgloss.Width(li.cram), 1, func() tea.Cmd {
						m.selectDeckByID(deck.ID)
						m.cramType = "flagged"
						return tea.Sequence(
							m.updateView(ViewCram),
							m.loadCramCards(),
						)
					})
				}
			}
		},
	})

	b.WriteString(listView)
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

func (m *Model) renderNav(x, y int) string {
	labels := []struct {
		id   string
		view View
		text string
	}{
		{"nav-dashboard", ViewDashboard, "Dashboard"},
		{"nav-decks", ViewDecks, "Decks"},
		{"nav-review", ViewReview, "Review"},
		{"nav-statistics", ViewStatistics, "Statistics"},
		{"nav-import", ViewImport, "Import"},
		{"nav-ai", ViewAI, "AI Drafts"},
		{"nav-settings", ViewSettings, "Settings"},
		{"nav-browser", ViewBrowser, "Browser"},
		{"nav-cram", ViewCram, "Cram"},
		{"nav-practice", ViewPractice, "Practice [0]"},
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render("deutsch-tui") + "\n\n")
	for i, l := range labels {
		style := navStyle
		if m.activeView == l.view {
			style = navActiveStyle
		}
		item := style.Render(l.text)
		b.WriteString(item + "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     l.id,
			View:   l.view,
			X:      x,
			Y:      y + 2 + i,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
	}
	return b.String()
}

func (m *Model) renderTabs(x, y int) string {
	tabs := []struct {
		id   string
		view View
		text string
	}{
		{"tab-dashboard", ViewDashboard, "Dashboard"},
		{"tab-decks", ViewDecks, "Decks"},
		{"tab-review", ViewReview, "Review"},
		{"tab-statistics", ViewStatistics, "Statistics"},
		{"tab-import", ViewImport, "Import"},
		{"tab-ai", ViewAI, "AI"},
		{"tab-settings", ViewSettings, "Settings"},
		{"tab-browser", ViewBrowser, "Browser"},
		{"tab-cram", ViewCram, "Cram"},
		{"tab-practice", ViewPractice, "Practice"},
	}

	var renderedTabs []string
	currentX := x
	for _, t := range tabs {
		style := tabStyle
		if m.activeView == t.view {
			style = tabActiveStyle
		}
		item := style.Render(t.text)
		renderedTabs = append(renderedTabs, item)
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     t.id,
			View:   t.view,
			X:      currentX,
			Y:      y,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
		currentX += lipgloss.Width(item)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
