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
		"  1-9      Switch to view\n" +
		"  Tab/arr  Cycle views\n" +
		"  w/s      Prev/next view\n" +
		"  ?        Toggle help\n" +
		"  q/Ctrl+c Quit"

	dash := sectionStyle.Render("Dashboard/Decks:") + "\n" +
		"  [ ]      Prev/next deck\n" +
		"  /        Search decks\n" +
		"  L        Edit deck limits\n" +
		"  +/-      Adjust limits\n" +
		"  !/@/#    Recent decks\n" +
		"  Enter    Select deck\n" +
		"  0        Practice"

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
		"  Practice 1-3/d-a/m-n select\n" +
		"  Import   i/I import, x/X exp\n" +
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
		if m.searchingDecks {
			b.WriteString(fmt.Sprintf("Search: %s_ (name or tag)\n\n", m.deckFilter))
		} else {
			b.WriteString(fmt.Sprintf("Filter: %s (Press / to edit)\n\n", m.deckFilter))
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
			miniBar = progressBar(5, deckPercentage, "46", "238") + " "
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

		highlightStyle := lipgloss.NewStyle().Foreground(colorPink).Bold(true)
		nameText := truncateLine(deck.Name, nameWidth)
		nameTextPadded := padLine(nameText, nameWidth)

		resDeckName := nameTextPadded
		if m.deckFilter != "" {
			resDeckName = highlightMatch(resDeckName, m.deckFilter, highlightStyle)
		}

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
			resTags := mutedStyle.Render(truncateLine(tags, layout.Width-15))
			if m.deckFilter != "" {
				resTags = highlightMatch(resTags, m.deckFilter, highlightStyle)
			}
			lines = append(lines, lineInfo{deckIdx: i, kind: "tags"})
			content.WriteString(fmt.Sprintf("     Tags: %s\n", resTags))
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
		footer = "Press Enter or Esc to finish searching."
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
			footer = fmt.Sprintf("%s select | %s search | %s stats | %s multi-select | %s clear\nPress enter to select deck.",
				keyStyle.Render("enter"), keyStyle.Render("/"), keyStyle.Render("v"), keyStyle.Render("x"), keyStyle.Render("Esc"))
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

func (m *Model) renderImport(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	var b strings.Builder
	titleLabelStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan).MarginBottom(1)
	b.WriteString(titleLabelStyle.Render("Import / Export") + "\n\n")

	importPathDisplay := m.importPath
	exportPathDisplay := m.exportPath
	if len(importPathDisplay) > maxInt(20, width-25) {
		importPathDisplay = "..." + importPathDisplay[len(importPathDisplay)-maxInt(20, width-25):]
	}
	if len(exportPathDisplay) > maxInt(20, width-25) {
		exportPathDisplay = "..." + exportPathDisplay[len(exportPathDisplay)-maxInt(20, width-25):]
	}

	importPathLabel := "Import file: " + importPathDisplay
	exportPathLabel := "Export file: " + exportPathDisplay

	btnActiveStyle := lipgloss.NewStyle().Bold(true).Foreground(colorPink)

	// Import Path
	importLabel := importPathLabel
	rowY := layout.Y + strings.Count(b.String(), "\n")
	if m.editingImportPath && m.importCursor == 0 {
		importLabel = "Import file: " + editStyle.Render(m.importPath+"_")
	} else if m.importCursor == 0 {
		importLabel = "> " + btnActiveStyle.Render(importPathLabel)
	} else {
		importLabel = "  " + importPathLabel
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-0",
		View:   ViewImport,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(importPathLabel) + 2,
		Height: 1,
	})
	b.WriteString(importLabel + "\n")

	// Export Path
	exportLabel := exportPathLabel
	rowY = layout.Y + strings.Count(b.String(), "\n")
	if m.editingImportPath && m.importCursor == 1 {
		exportLabel = "Export file: " + editStyle.Render(m.exportPath+"_")
	} else if m.importCursor == 1 {
		exportLabel = "> " + btnActiveStyle.Render(exportPathLabel)
	} else {
		exportLabel = "  " + exportPathLabel
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "export-path-1",
		View:   ViewImport,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(exportPathLabel) + 2,
		Height: 1,
	})
	b.WriteString(exportLabel + "\n\n")

	// Export Deck Filter
	exportDeckLabel := "Export Deck: "
	if m.exportDeckID == "" {
		exportDeckLabel += "All Decks"
	} else {
		found := false
		for _, d := range m.decks {
			if d.ID == m.exportDeckID {
				exportDeckLabel += d.Name
				found = true
				break
			}
		}
		if !found {
			exportDeckLabel += "Unknown Deck"
		}
	}
	rowY = layout.Y + strings.Count(b.String(), "\n")
	if m.importCursor == 2 {
		b.WriteString("> " + btnActiveStyle.Render(exportDeckLabel) + " (Use [ / ] to change)\n")
	} else {
		b.WriteString("  " + exportDeckLabel + "\n")
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-2",
		View:   ViewImport,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(exportDeckLabel) + 2,
		Height: 1,
	})

	// Export Filter
	exportTagLabel := "Export Filter: " + m.exportTag
	if m.exportTag == "" {
		exportTagLabel += "(None)"
	}
	rowY = layout.Y + strings.Count(b.String(), "\n")
	if m.editingExportTag {
		b.WriteString("  Export Filter: " + editStyle.Render(m.exportTag+"_") + "\n\n")
	} else if m.importCursor == 3 {
		b.WriteString("> " + btnActiveStyle.Render(exportTagLabel) + " (Press t to edit)\n\n")
	} else {
		b.WriteString("  " + exportTagLabel + "\n\n")
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-3",
		View:   ViewImport,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(exportTagLabel) + 2,
		Height: 1,
	})

	// Export Status Filter
	exportStatusLabel := "Export Status: " + m.exportFilter
	rowY = layout.Y + strings.Count(b.String(), "\n")
	if m.importCursor == 4 {
		b.WriteString("> " + btnActiveStyle.Render(exportStatusLabel) + " (Use [ / ] to change)\n\n")
	} else {
		b.WriteString("  " + exportStatusLabel + "\n\n")
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-4",
		View:   ViewImport,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(exportStatusLabel) + 2,
		Height: 1,
	})

	if strings.TrimSpace(m.importPath) == "" {
		b.WriteString(warnStyle.Render("Import file is empty; set a path before importing.") + "\n")
	}
	if strings.TrimSpace(m.exportPath) == "" {
		b.WriteString(warnStyle.Render("Export file is empty; set a path before exporting.") + "\n")
	}

	b.WriteString("Actions:\n")
	actions := []struct {
		id    string
		label string
		key   string
	}{
		{"import-tsv", "Import TSV", "i"},
		{"import-apkg", "Import APKG", "I"},
		{"seed-std", "Seed Standard", "S"},
		{"export-tsv", "Export TSV", "x"},
		{"export-apkg", "Export APKG", "X"},
		{"reset-db", "Reset DB", "R"},
	}

	rowY = layout.Y + strings.Count(b.String(), "\n")
	currentX := layout.X
	for _, a := range actions {
		item := fmt.Sprintf("[%s] %s", a.key, a.label)
		b.WriteString(btnStyle.Render(item))
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     a.id,
			View:   ViewImport,
			X:      currentX,
			Y:      rowY,
			Width:  lipgloss.Width(item) + 2,
			Height: 1,
		})
		currentX += lipgloss.Width(item) + 3
	}

	if m.editingImportPath || m.editingExportTag {
		b.WriteString("\n\n" + infoStyle.Render("EDITING - Enter to save, Esc to cancel.") + "\n\n")
	} else {
		b.WriteString(fmt.Sprintf("\n\nCurrent Deck: %s\nStatus filters apply to TSV and APKG exports.\nUse %s/%s to navigate, %s/%s to edit, %s / %s to change deck, or click buttons.\n",
			m.deckLabel(),
			keyStyle.Render("j"), keyStyle.Render("k"),
			keyStyle.Render("Enter"), keyStyle.Render("t"),
			keyStyle.Render("["), keyStyle.Render("]")))
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
		currentX += lipgloss.Width(item) + 1
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
