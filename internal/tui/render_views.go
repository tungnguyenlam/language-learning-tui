package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderHelp() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).Underline(true)
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	colStyle := lipgloss.NewStyle().PaddingRight(4).Width(35)

	global := sectionStyle.Render("Global:") + "\n" +
		"  1-9          Switch to view\n" +
		"  Tab/arrows   Cycle views\n" +
		"  w/s          Previous/next view\n" +
		"  ?            Toggle this help\n" +
		"  q/Ctrl+c     Quit"

	dash := sectionStyle.Render("Dashboard/Decks:") + "\n" +
		"  [ ]          Previous/next deck\n" +
		"  /            Search decks\n" +
		"  L            Edit deck limits\n" +
		"  +/-          Adjust limits\n" +
		"  Enter        Select deck"

	review := sectionStyle.Render("Review:") + "\n" +
		"  Space/Enter  Reveal answer\n" +
		"  a/h/g/e      Again/Hard/Good/Easy\n" +
		"  1-4          Grade or MCQ choice\n" +
		"  b / B        Bookmark / Filter\n" +
		"  u / r        Undo / History\n" +
		"  f / i        Focus / Info\n" +
		"  p            Play audio"

	browser := sectionStyle.Render("Browser:") + "\n" +
		"  j/k          Navigate\n" +
		"  / / #        Search / Tag filter\n" +
		"  m            Select card\n" +
		"  b / x        Bookmark / Suspend\n" +
		"  t / T        Toggle kind / Tags\n" +
		"  C            Cleanup tags\n" +
		"  Backspace    Delete card"

	other := sectionStyle.Render("Other:") + "\n" +
		"  Statistics   j/k scroll, x export\n" +
		"  AI           / topic, a/d drafts\n" +
		"  Cram         Enter start, 1-5 filter\n" +
		"  Import       i/I import, x/X export\n" +
		"  Settings     j/k nav, +/- goal"

	col1 := colStyle.Render(global + "\n\n" + dash)
	col2 := colStyle.Render(review)
	col3 := colStyle.Render(browser)
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

	start := 0
	end := len(filteredDecks)
	maxVisible := 10
	if layout.Height > 25 {
		maxVisible = (layout.Height - 10) / 2
	}
	if maxVisible < 5 {
		maxVisible = 5
	}

	if end > maxVisible {
		start = m.deckCursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(filteredDecks) {
			end = len(filteredDecks)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
	}

	listStartY := strings.Count(b.String(), "\n")
	lineWidth := scrollbarLineWidth(layout.Width)
	thumbStart, thumbHeight := scrollbarThumb(len(filteredDecks), maxVisible, start)

	for i := start; i < end; i++ {
		deck := filteredDecks[i]
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
		counts := fmt.Sprintf("%s new, %s due, %s total", newStyled, dueStyled, totalStyled)

		// Add a mini progress bar for deck completion
		deckPercentage := 0.0
		if deck.TotalCards > 0 {
			deckPercentage = float64(deck.TotalCards-deck.DueCards) / float64(deck.TotalCards)
		}
		miniBar := progressBar(5, deckPercentage, "46", "238")

		deckNamePadded := deck.Name
		if len(deckNamePadded) < 40 {
			deckNamePadded += strings.Repeat(" ", 40-len(deckNamePadded))
		} else if len(deckNamePadded) > 40 {
			deckNamePadded = deckNamePadded[:37] + "..."
		}

		label := fmt.Sprintf("%s%s%s %s  %s | today %d, %.0f%% success",
			prefix, selectMark, deckNamePadded, miniBar, counts, deck.ReviewsToday, deck.SuccessRate*100)

		studyBtn := " [Study]"
		cramBtn := " [Cram]"
		currBtnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if i == m.deckCursor {
			currBtnStyle = currBtnStyle.Foreground(colorBlue)
		}

		rowY := layout.Y + strings.Count(b.String(), "\n")
		line := padLine(style.Render(label)+currBtnStyle.Render(studyBtn)+currBtnStyle.Render(cramBtn), lineWidth)

		if len(filteredDecks) > maxVisible {
			currentPos := i - start
			scrollbarChar := "│"
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			line += " " + scrollbarChar

			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("deck-scroll-%d", currentPos),
				View:   ViewDecks,
				X:      layout.X + lineWidth + 1,
				Y:      layout.Y + listStartY + currentPos,
				Width:  1,
				Height: 1,
			})
		}
		b.WriteString(line + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("deck-select-%d", i),
			View:   ViewDecks,
			X:      layout.X,
			Y:      rowY,
			Width:  lipgloss.Width(label),
			Height: 1,
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("deck-study-%d", i),
			View:   ViewDecks,
			X:      layout.X + lipgloss.Width(label),
			Y:      rowY,
			Width:  lipgloss.Width(studyBtn),
			Height: 1,
		})
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("deck-cram-%d", i),
			View:   ViewDecks,
			X:      layout.X + lipgloss.Width(label) + lipgloss.Width(studyBtn),
			Y:      rowY,
			Width:  lipgloss.Width(cramBtn),
			Height: 1,
		})

		// Show limits only when editing this deck
		if i == m.deckCursor && m.editingDeckLimits {
			newStyle := lipgloss.NewStyle().Foreground(colorBlue)
			reviewStyle := lipgloss.NewStyle().Foreground(colorBlue)
			if m.limitCursor == 0 {
				newStyle = newStyle.Bold(true).Underline(true)
			} else {
				reviewStyle = reviewStyle.Bold(true).Underline(true)
			}
			limitsLabel := fmt.Sprintf("     Limits: New %s, Review %s (h/l switch, +/- adjust)",
				newStyle.Render(fmt.Sprintf("%d", deck.NewCardsPerDay)),
				reviewStyle.Render(fmt.Sprintf("%d", deck.ReviewLimitPerDay)))
			b.WriteString(limitsLabel + "\n")
		}

		if deck.Description != "" {
			b.WriteString(fmt.Sprintf("     %s\n", mutedStyle.Render(deck.Description)))
		}
		// Show tags if they exist
		if len(deck.Tags) > 0 {
			tags := strings.Join(deck.Tags, ", ")
			b.WriteString(fmt.Sprintf("     Tags: %s\n", mutedStyle.Render(tags)))
		}
	}

	if m.searchingDecks {
		b.WriteString("\nPress Enter or Esc to finish searching.")
	} else {
		selectedCount := 0
		for _, s := range m.deckSelected {
			if s {
				selectedCount++
			}
		}
		if selectedCount > 0 {
			b.WriteString(fmt.Sprintf("\n%d decks selected. Press %s to delete, %s to merge into current.", selectedCount,
				keyStyle.Render("Backspace"),
				keyStyle.Render("M")))
		} else {
			b.WriteString(fmt.Sprintf("\n%s select | %s search | %s stats | %s multi-select | %s cram deck | %s clear\n",
				keyStyle.Render("enter"), keyStyle.Render("/"), keyStyle.Render("v"), keyStyle.Render("x"), keyStyle.Render("c"), keyStyle.Render("Esc")))
			b.WriteString(mutedStyle.Render("Press enter to select deck."))
		}
	}

	if len(filteredDecks) > maxVisible {
		b.WriteString(fmt.Sprintf(" (Showing %d-%d of %d, Use Mouse Wheel to scroll)", start+1, end, len(filteredDecks)))
	}
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
		ID:     "import-path-1",
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
