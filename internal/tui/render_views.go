package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderHelp() string {
	return "Keyboard Shortcuts\n\nGlobal:\n  1-9          Switch to view\n  Tab/arrows   Cycle views\n  w/s          Previous/next view\n  ?            Toggle this help\n  q/Ctrl+c     Quit\n\nDashboard/Decks:\n  [ ]          Previous/next deck\n  Enter        Select deck (Decks view)\n\nReview:\n  Space/Enter  Reveal answer\n  a/h/g/e      Grade Again/Hard/Good/Easy\n  b            Toggle bookmark\n  B            Toggle bookmarked-only mode\n  x            Suspend card\n  u            Undo last review\n  r            Toggle card review history\n  p            Play audio\n  1-4          Select MCQ choice\n\nStatistics:\n  j/k          Scroll stats\n  x            Export deck stats to CSV\n\nBrowser:\n  j/k          Navigate cards\n  Enter        Toggle card review history\n  Type         Search cards\n  C            Cleanup unused tags\n  Backspace    Delete search char\n\nCram:\n  j/k          Navigate cards\n  Enter        Start cram review\n  p            Play audio (in review)\n  1-5          Filter: bookmarked/suspended/leech/flagged/all\n\nImport:\n  j/k          Select field\n  Enter        Start/stop editing path\n  i/I          Import TSV/APKG\n  x/X          Export TSV/APKG\n\nSettings:\n  j/k          Navigate options\n  +/-          Adjust daily goal\n  Enter        Toggle AI provider / edit template"
}

func (m *Model) renderDecks(x, y int) string {
	var b strings.Builder
	b.WriteString("Decks\n\n")

	// Show filter if active or searching
	if m.searchingDecks || m.deckFilter != "" {
		if m.searchingDecks {
			b.WriteString(fmt.Sprintf("Search: %s_\n\n", m.deckFilter))
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
	maxVisible := 5
	if m.height > 20 {
		maxVisible = (m.height - 15) / 3 // Each deck takes about 3 lines
	}
	if maxVisible < 3 {
		maxVisible = 3
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

	for i := start; i < end; i++ {
		deck := filteredDecks[i]
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.deckCursor {
			prefix = "> "
			style = style.Bold(true).Foreground(lipgloss.Color("212"))
		}

		selectMark := "[ ] "
		if m.deckSelected[deck.ID] {
			selectMark = "[x] "
		}

		counts := fmt.Sprintf("%d new, %d due, %d total", deck.NewCards, deck.DueCards, deck.TotalCards)

		label := fmt.Sprintf("%s%s%s (%s, today %d, %.0f%% success)",
			prefix, selectMark, deck.Name, counts, deck.ReviewsToday, deck.SuccessRate*100)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
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
			b.WriteString(fmt.Sprintf("\n%d decks selected. Press Backspace to delete, M to merge into current.", selectedCount))
		} else {
			b.WriteString("\nPress enter to select deck. Press / to search. v to view stats. x to multi-select. Esc to clear.")
		}
	}

	if len(filteredDecks) > maxVisible {
		b.WriteString(fmt.Sprintf(" (Showing %d-%d of %d)", start+1, end, len(filteredDecks)))
	}
	return b.String()
}

func (m *Model) renderImport(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Import / Export") + "\n\n")

	importPathLabel := "Import file: " + m.importPath
	exportPathLabel := "Export file: " + m.exportPath

	btnActiveStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	editStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("62"))

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
		for _, d := range m.decks {
			if d.ID == m.exportDeckID {
				exportDeckLabel += d.Name
				break
			}
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

	// Export Tag Filter
	exportTagLabel := "Export Tag: " + m.exportTag
	if m.exportTag == "" {
		exportTagLabel += "(None)"
	}
	rowY = layout.Y + strings.Count(b.String(), "\n")
	if m.editingExportTag {
		b.WriteString("  Export Tag: " + editStyle.Render(m.exportTag+"_") + "\n\n")
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

	b.WriteString("Actions:\n")
	actions := []struct {
		id    string
		label string
		key   string
	}{
		{"import-tsv", "Import TSV", "i"},
		{"import-apkg", "Import APKG", "I"},
		{"export-tsv", "Export TSV", "x"},
		{"export-apkg", "Export APKG", "X"},
	}

	btnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("62")).
		Padding(0, 1).
		MarginRight(1)

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
		b.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("EDITING - Enter to save, Esc to cancel.") + "\n\n")
	} else {
		b.WriteString("\n\nCurrent Deck: " + m.deckLabel() + "\n")
		b.WriteString("Use j/k to navigate, Enter/t to edit, [ / ] to change deck, or click buttons.\n")
	}

	return b.String()
}
