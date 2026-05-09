package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderSessionSummary(layout viewportLayout) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)

	duration := time.Since(m.sessionStartTime)
	if m.sessionReviewed == 0 {
		duration = 0
	}

	accuracy := 0.0
	if m.sessionReviewed > 0 {
		accuracy = float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("SESSION SUMMARY") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("No cards due") + "\n\n")

	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(1, 2).
		Width(maxInt(40, layout.Width-10))

	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))

	summaryContent := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s",
		statsStyle.Render("Cards Reviewed:"), valStyle.Render(fmt.Sprintf("%d", m.sessionReviewed)),
		statsStyle.Render("Correct Answers:"), valStyle.Render(fmt.Sprintf("%d", m.sessionCorrect)),
		statsStyle.Render("Accuracy Score:"), valStyle.Render(fmt.Sprintf("%.1f%%", accuracy)),
		statsStyle.Render("Session Time:  "), valStyle.Render(duration.Round(time.Second).String()),
	)

	b.WriteString(mainBox.Render(summaryContent) + "\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Press any key or Enter to return to Dashboard.") + "\n")

	// Motivational message
	motivation := "Keep it up!"
	if accuracy >= 90 {
		motivation = "Excellent work! You're mastering this!"
	} else if accuracy >= 75 {
		motivation = "Good job! You're making solid progress!"
	}

	b.WriteString("\n" + lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("213")).Render(motivation))

	return b.String()
}
