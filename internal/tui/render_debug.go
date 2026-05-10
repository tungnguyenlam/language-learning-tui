package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderDebug(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	b.WriteString(titleStyle.Render("DEBUG LOG") + "\n\n")

	b.WriteString(fmt.Sprintf("Active View: %s\n", m.activeView))
	b.WriteString(fmt.Sprintf("Deck: %s (%s)\n", m.deck.Name, m.deck.ID))
	b.WriteString(fmt.Sprintf("Due Cards: %d\n", len(m.dueCards)))
	b.WriteString(fmt.Sprintf("Browser Cards: %d\n", len(m.browserCards)))
	b.WriteString(fmt.Sprintf("Recent Decks: %v\n", m.recentDecks))
	b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("--- LOG CONTENT ---") + "\n")

	// Since we don't have a direct way to read the log file easily without blocking or complexity,
	// let's just show some internal state for now.
	// In a real scenario, we'd have a buffered logger in memory.

	b.WriteString("State Summary:\n")
	b.WriteString(fmt.Sprintf("- textInputActive: %v\n", m.textInputActive()))
	b.WriteString(fmt.Sprintf("- confirmingDelete: %v\n", m.confirmingDelete))
	b.WriteString(fmt.Sprintf("- gradingInProgress: %v\n", m.gradingInProgress))
	b.WriteString(fmt.Sprintf("- drafting: %v\n", m.drafting))

	b.WriteString("\nPress Ctrl+D to exit debug view.")

	return style.Render(b.String())
}
