package tui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderHelp() string {
	return "Keyboard Shortcuts\n\nGlobal:\n  1-9          Switch to view\n  Tab/arrows   Cycle views\n  w/s          Previous/next view\n  ?            Toggle this help\n  q/Ctrl+c     Quit\n\nDashboard/Decks:\n  [ ]          Previous/next deck\n  Enter        Select deck (Decks view)\n\nReview:\n  Space/Enter  Reveal answer\n  a/h/g/e      Grade Again/Hard/Good/Easy\n  b            Toggle bookmark\n  B            Toggle bookmarked-only mode\n  x            Suspend card\n  u            Undo last review\n  r            Toggle card review history\n  p            Play audio\n  1-4          Select MCQ choice\n\nBrowser:\n  j/k          Navigate cards\n  Enter        Toggle card review history\n  Type         Search cards\n  Backspace    Delete search char\n\nCram:\n  j/k          Navigate cards\n  Enter        Start cram review\n  p            Play audio (in review)\n  1-5          Filter: bookmarked/suspended/leech/flagged/all\n\nImport:\n  j/k          Select field\n  Enter        Start/stop editing path\n  i/I          Import TSV/APKG\n  x/X          Export TSV/APKG\n\nSettings:\n  j/k          Navigate options\n  +/-          Adjust daily goal\n  Enter        Toggle AI provider / edit template"
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

		newStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
		dueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
		totalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159"))

		counts := fmt.Sprintf("%s new, %s due, %s total",
			newStyle.Render(strconv.Itoa(deck.NewCards)),
			dueStyle.Render(strconv.Itoa(deck.DueCards)),
			totalStyle.Render(strconv.Itoa(deck.TotalCards)))

		label := fmt.Sprintf("%s%s (%s, today %d, %.0f%% success)",
			prefix, deck.Name, counts, deck.ReviewsToday, deck.SuccessRate*100)
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
		b.WriteString("\nPress enter to select deck. Press / to search. Esc to clear filter.")
	}

	if len(filteredDecks) > maxVisible {
		b.WriteString(fmt.Sprintf(" (Showing %d-%d of %d)", start+1, end, len(filteredDecks)))
	}
	return b.String()
}

func (m *Model) renderImport(x, y int) string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Import / Export") + "\n\n")

	importPathLabel := "Import file: " + m.importPath
	exportPathLabel := "Export file: " + m.exportPath

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	editStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("62"))

	// Import Path
	importLabel := importPathLabel
	if m.editingImportPath && m.importCursor == 0 {
		importLabel = "Import file: " + editStyle.Render(m.importPath+"_")
	} else if m.importCursor == 0 {
		importLabel = "> " + style.Render(importPathLabel)
	} else {
		importLabel = "  " + importPathLabel
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-0",
		View:   ViewImport,
		X:      x,
		Y:      y + 2,
		Width:  lipgloss.Width(importPathLabel) + 2,
		Height: 1,
	})
	b.WriteString(importLabel + "\n")

	// Export Path
	exportLabel := exportPathLabel
	if m.editingImportPath && m.importCursor == 1 {
		exportLabel = "Export file: " + editStyle.Render(m.exportPath+"_")
	} else if m.importCursor == 1 {
		exportLabel = "> " + style.Render(exportPathLabel)
	} else {
		exportLabel = "  " + exportPathLabel
	}
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "import-path-1",
		View:   ViewImport,
		X:      x,
		Y:      y + 3,
		Width:  lipgloss.Width(exportPathLabel) + 2,
		Height: 3,
	})
	b.WriteString(exportLabel + "\n\n")

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

	currentX := x
	for _, a := range actions {
		item := fmt.Sprintf("[%s] %s", a.key, a.label)
		b.WriteString(btnStyle.Render(item))
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     a.id,
			View:   ViewImport,
			X:      currentX,
			Y:      y + 7,
			Width:  lipgloss.Width(item) + 2,
			Height: 1,
		})
		currentX += lipgloss.Width(item) + 3
	}

	b.WriteString("\n\nCurrent Deck: " + m.deckLabel() + "\n\n")
	if m.editingImportPath {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("EDITING - Enter to save, Esc to cancel.") + "\n\n")
	}
	b.WriteString("Use j/k to select path, Enter to edit, or click buttons above.\n")

	return b.String()
}
