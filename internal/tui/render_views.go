package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderHelp() string {
	return "Keyboard Shortcuts\n\nGlobal:\n  1-9          Switch to view\n  Tab/arrows   Cycle views\n  w/s          Previous/next view\n  ?            Toggle this help\n  q/Ctrl+c     Quit\n\nDashboard/Decks:\n  [ ]          Previous/next deck\n  /            Search decks\n  L            Edit per-deck limits\n  h/l or <-/-> Switch limit field (while editing)\n  +/-          Adjust selected limit (while editing)\n  Enter        Select deck (Decks view)\n\nReview:\n  Space/Enter  Reveal answer\n  a/h/g/e      Grade Again/Hard/Good/Easy\n  b            Toggle bookmark\n  B            Toggle bookmarked-only mode\n  x            Suspend card\n  u            Undo last review\n  r            Toggle card review history\n  p            Play audio\n  1-4          Select MCQ choice\n\nStatistics:\n  j/k          Scroll stats\n  x            Export deck stats to CSV\n\nBrowser:\n  j/k          Navigate cards\n  /            Search cards\n  m            Select card\n  b/B          Bookmark (single/bulk)\n  x/X          Suspend (single/bulk)\n  t/T          Toggle kind / edit tags\n  C            Cleanup unused tags\n  Enter        Toggle card review history\n  Backspace    Delete selected card(s)\n\nAI:\n  /            Edit topic\n  Enter        Generate or approve selected draft\n  [ / ]        Previous/next template\n  a / A        Approve one/all drafts\n  d / D        Discard one/all drafts\n\nCram:\n  j/k          Navigate cards\n  Enter        Start cram review\n  p            Play audio (in review)\n  1-5          Filter: bookmarked/suspended/leech/flagged/all\n\nImport:\n  j/k          Select field\n  Enter        Start/stop editing path\n  i/I          Import TSV/APKG\n  x/X          Export TSV/APKG\n\nSettings:\n  j/k          Navigate options\n  [ / ]        Cycle AI template set\n  +/-          Adjust daily goal\n  Enter        Toggle AI provider / edit template"
}

func (m *Model) renderDecks(layout viewportLayout) string {
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
	if layout.Height > 20 {
		maxVisible = (layout.Height - 15) / 2
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

		b.WriteString(style.Render(label) + "\n")

		// Show limits only when editing this deck
		if i == m.deckCursor && m.editingDeckLimits {
			newStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
			reviewStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
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
				lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("Backspace"),
				lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("M")))
		} else {
			keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
			b.WriteString(fmt.Sprintf("\nPress %s to select deck. Press %s to search. %s to view stats. %s to multi-select. %s to clear.",
				keyStyle.Render("enter"), keyStyle.Render("/"), keyStyle.Render("v"), keyStyle.Render("x"), keyStyle.Render("Esc")))
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

	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	if strings.TrimSpace(m.importPath) == "" {
		b.WriteString(warningStyle.Render("Import file is empty; set a path before importing.") + "\n")
	}
	if strings.TrimSpace(m.exportPath) == "" {
		b.WriteString(warningStyle.Render("Export file is empty; set a path before exporting.") + "\n")
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
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(fmt.Sprintf("\n\nCurrent Deck: %s\nUse %s/%s to navigate, %s/%s to edit, %s / %s to change deck, or click buttons.\n",
			m.deckLabel(),
			keyStyle.Render("j"), keyStyle.Render("k"),
			keyStyle.Render("Enter"), keyStyle.Render("t"),
			keyStyle.Render("["), keyStyle.Render("]")))
	}

	return b.String()
}
