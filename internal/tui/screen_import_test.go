package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

var importActionLabels = []string{
	"[i] Import TSV", "[I] Import APKG", "[A] Browse AnkiWeb",
	"[S] Seed Standard", "[x] Export TSV", "[X] Export APKG",
	"[B] Backup", "[U] Restore", "[R] Reset DB",
}

var importActionIDs = []string{"import-tsv", "import-apkg", "browse-ankiweb",
	"seed-std", "export-tsv", "export-apkg", "backup-progress", "restore-progress", "reset-db"}

// renderImport renders the Import screen at a terminal width and returns its
// content lines alongside the layout Render actually derived, since Render
// recomputes the layout from the panel chrome rather than using the one passed.
func renderImport(t *testing.T, termWidth int) (*Model, []string, viewportLayout) {
	t.Helper()

	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewImport
	m.width, m.height = termWidth, 40

	view := m.importScreen.Render(m, viewportLayout{X: 0, Y: 0})

	width, height := m.activePanelSize()
	layout := contentLayoutForStyle(panelStyle.Width(width).Height(height), 0, 0)
	return m, strings.Split(view, "\n"), layout
}

// The action row grew past the panel width once AnkiWeb was added. Letting the
// panel hard-wrap it split "[x] Export TSV" across two lines and left the
// button's hitbox pointing at cells it no longer occupies, so the row now wraps
// on whole buttons.
func TestImportActionButtonsWrapWhole(t *testing.T) {
	for _, termWidth := range []int{80, 100, 120, 160} {
		_, lines, layout := renderImport(t, termWidth)
		view := strings.Join(lines, "\n")

		for _, label := range importActionLabels {
			if !strings.Contains(view, label) {
				t.Errorf("width %d: %q was split across lines:\n%s", termWidth, label, view)
			}
		}

		// Only the button rows are checked: the footer help line has always
		// relied on the panel to wrap it.
		for _, line := range lines {
			isButtonRow := false
			for _, label := range importActionLabels {
				isButtonRow = isButtonRow || strings.Contains(line, label)
			}
			if !isButtonRow {
				continue
			}
			if got := lipgloss.Width(line); got > layout.Width {
				t.Errorf("width %d: button row overflows the panel by %d cells: %q",
					termWidth, got-layout.Width, line)
			}
		}
	}
}

// Hitboxes are computed alongside the wrapping, so a wrapped row must still
// place every button over the cells it was actually drawn on.
func TestImportActionHitboxesFollowTheWrap(t *testing.T) {
	// Narrow enough that the row must wrap.
	m, lines, layout := renderImport(t, 90)

	for _, id := range importActionIDs {
		var box Hitbox
		for _, h := range m.hitboxes {
			if h.ID == id {
				box = h
			}
		}
		if box.ID == "" {
			t.Fatalf("no hitbox registered for %q", id)
		}

		row := box.Y - layout.Y
		if row < 0 || row >= len(lines) {
			t.Fatalf("%s: hitbox row %d is outside the rendered view", id, row)
		}
		col := box.X - layout.X
		if col < 0 || col+box.Width > lipgloss.Width(lines[row]) {
			t.Errorf("%s: hitbox spans cols %d..%d of a %d-wide row: %q",
				id, col, col+box.Width, lipgloss.Width(lines[row]), lines[row])
		}
	}

	// Two buttons must never claim the same cell.
	for i, a := range importActionIDs {
		for _, bID := range importActionIDs[i+1:] {
			ha, hb := findHitbox(m, a), findHitbox(m, bID)
			if ha.Y == hb.Y && ha.X < hb.X+hb.Width && hb.X < ha.X+ha.Width {
				t.Errorf("%s and %s overlap on row %d", a, bID, ha.Y)
			}
		}
	}
}

func findHitbox(m *Model, id string) Hitbox {
	var found Hitbox
	for _, h := range m.hitboxes {
		if h.ID == id {
			found = h
		}
	}
	return found
}
