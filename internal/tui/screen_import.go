package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// importScreen is the Import/Export view. Its view-local UI state (the row
// cursor and the two inline-edit flags) lives here. The export parameters it
// edits — importPath, exportPath, exportDeckID, exportTag, exportFilter — stay
// on Model because the import/export IO commands (importTSV, exportTSV, …) and
// the statistics CSV export consume them.
type importScreen struct {
	importCursor      int // 0: import path, 1: export path, 2: export deck, 3: export tag, 4: export filter
	editingImportPath bool
	editingExportTag  bool
}

// resetCursor clears the transient UI state; called on entering the view.
func (s *importScreen) resetCursor() {
	s.importCursor = 0
	s.editingImportPath = false
	s.editingExportTag = false
}

func (s *importScreen) HandleKey(m *Model, msg tea.KeyPressMsg) (tea.Cmd, bool) {
	if s.editingImportPath {
		switch msg.String() {
		case "enter", "\r", "\n", "esc":
			s.editingImportPath = false
			return nil, true
		case "backspace":
			if s.importCursor == 0 {
				if len(m.importPath) > 0 {
					m.importPath = trimLastRune(m.importPath)
				}
			} else {
				if len(m.exportPath) > 0 {
					m.exportPath = trimLastRune(m.exportPath)
				}
			}
			return nil, true
		case "ctrl+u":
			if s.importCursor == 0 {
				m.importPath = ""
			} else {
				m.exportPath = ""
			}
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			if s.importCursor == 0 {
				m.importPath += ch
			} else {
				m.exportPath += ch
			}
			return nil, true
		}
		return nil, true
	}

	if s.editingExportTag {
		switch msg.String() {
		case "enter", "\r", "\n", "esc":
			s.editingExportTag = false
			return nil, true
		case "backspace":
			if len(m.exportTag) > 0 {
				m.exportTag = trimLastRune(m.exportTag)
			}
			return nil, true
		case "ctrl+u":
			m.exportTag = ""
			return nil, true
		}
		if ch, ok := singlePrintableInput(msg.String()); ok {
			m.exportTag += ch
			return nil, true
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if s.importCursor > 0 {
			s.importCursor--
		}
		return nil, true
	case "down", "j":
		if s.importCursor < 4 {
			s.importCursor++
		}
		return nil, true
	case "[":
		if s.importCursor == 2 {
			m.previousExportDeck()
			return nil, true
		} else if s.importCursor == 4 {
			m.cycleExportFilter(false)
			return nil, true
		}
	case "]":
		if s.importCursor == 2 {
			m.nextExportDeck()
			return nil, true
		} else if s.importCursor == 4 {
			m.cycleExportFilter(true)
			return nil, true
		}
	case "enter", "\r", "\n":
		if s.importCursor < 2 {
			s.editingImportPath = true
		} else if s.importCursor == 3 {
			s.editingExportTag = true
		}
		return nil, true
	case "t":
		s.importCursor = 3
		s.editingExportTag = true
		return nil, true
	case "i":
		return m.importTSV(), true
	case "I":
		return m.importAPKG(), true
	case "A":
		return m.updateView(ViewAnkiWeb), true
	case "S", "s":
		return m.seedStandardContent(), true
	case "R":
		return m.handleResetDatabase(), true
	case "B":
		return m.handleBackupProgress(), true
	case "U":
		return m.handleRestoreProgress(), true
	case "x":
		return m.exportTSV(), true
	case "X":
		return m.exportAPKG(), true
	}
	return nil, false
}

func (s *importScreen) Render(m *Model, layout viewportLayout) string {
	x, y := layout.X, layout.Y
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout = contentLayoutForStyle(style, x, y)

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
	if s.editingImportPath && s.importCursor == 0 {
		importLabel = "Import file: " + editStyle.Render(m.importPath+"_")
	} else if s.importCursor == 0 {
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
	if s.editingImportPath && s.importCursor == 1 {
		exportLabel = "Export file: " + editStyle.Render(m.exportPath+"_")
	} else if s.importCursor == 1 {
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
	if s.importCursor == 2 {
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
	if s.editingExportTag {
		b.WriteString("  Export Filter: " + editStyle.Render(m.exportTag+"_") + "\n\n")
	} else if s.importCursor == 3 {
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
	if s.importCursor == 4 {
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

	backupLabel := "Progress backup: none yet"
	if m.lastBackupPath != "" {
		backupLabel = "Progress backup: " + filepath.Base(m.lastBackupPath)
	}
	b.WriteString(mutedStyle.Render(backupLabel) + "\n")
	b.WriteString(mutedStyle.Render("Backups exclude the dictionary. Restore replaces decks and review history.") + "\n")

	b.WriteString("Actions:\n")
	actions := []struct {
		id    string
		label string
		key   string
	}{
		{"import-tsv", "Import TSV", "i"},
		{"import-apkg", "Import APKG", "I"},
		{"browse-ankiweb", "Browse AnkiWeb", "A"},
		{"seed-std", "Seed Standard", "S"},
		{"export-tsv", "Export TSV", "x"},
		{"export-apkg", "Export APKG", "X"},
		{"backup-progress", "Backup", "B"},
		{"restore-progress", "Restore", "U"},
		{"reset-db", "Reset DB", "R"},
	}

	// Wrap on whole buttons. Letting the panel hard-wrap the row instead splits
	// a label across lines, which both reads badly and leaves the button's
	// hitbox pointing at the wrong cells.
	rowY = layout.Y + strings.Count(b.String(), "\n")
	currentX := layout.X
	for i, a := range actions {
		item := fmt.Sprintf("[%s] %s", a.key, a.label)
		width := lipgloss.Width(item) + 2
		// btnStyle carries a right margin, so the gap between buttons is
		// already there; only the row break has to be written. The +1 keeps
		// that trailing margin inside the panel too — otherwise it alone wraps
		// and leaves a stray blank row.
		if i > 0 && currentX+width+1 > layout.X+layout.Width {
			b.WriteString("\n")
			rowY++
			currentX = layout.X
		}
		b.WriteString(btnStyle.Render(item))
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     a.id,
			View:   ViewImport,
			X:      currentX,
			Y:      rowY,
			Width:  width,
			Height: 1,
		})
		currentX += width + 1
	}

	if s.editingImportPath || s.editingExportTag {
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
