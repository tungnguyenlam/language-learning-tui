package tui

import (
	"fmt"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// summaryScreen is the post-session statistics screen shown after a review
// session ends. It carries no view-local state: the session counters it reads
// (sessionReviewed/sessionCorrect/sessionStartTime/sessionGrades) are updated
// during review, so they remain shared state on Model.
type summaryScreen struct{}

// HandleKey returns to the dashboard on any key, snapshotting the session
// totals into the "last session" fields and resetting the live counters.
func (summaryScreen) HandleKey(m *Model, _ tea.KeyPressMsg) (tea.Cmd, bool) {
	m.lastSessionReviewed = m.sessionReviewed
	m.lastSessionCorrect = m.sessionCorrect
	m.lastSessionDuration = time.Since(m.sessionStartTime)
	m.sessionReviewed = 0
	m.sessionCorrect = 0
	return m.updateView(ViewDashboard), true
}

func (summaryScreen) Render(m *Model, layout viewportLayout) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)

	duration := time.Since(m.sessionStartTime)
	if m.sessionReviewed == 0 {
		duration = 0
	}

	accuracy := 0.0
	if m.sessionReviewed > 0 {
		accuracy = float64(m.sessionCorrect) / float64(m.sessionReviewed) * 100
	}

	cardsPerMin := 0.0
	if duration.Minutes() > 0 {
		cardsPerMin = float64(m.sessionReviewed) / duration.Minutes()
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("SESSION SUMMARY") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Session complete!") + "\n\n")

	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(1, 2).
		Width(maxInt(40, layout.Width-10))

	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159")).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))

	barWidth := 20
	filled := int((accuracy / 100) * float64(barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}
	barColor := "214"
	if accuracy >= 80 {
		barColor = "46"
	} else if accuracy < 50 && accuracy > 0 {
		barColor = "196"
	}
	accBar := lipgloss.NewStyle().Foreground(lipgloss.Color(barColor)).Render(strings.Repeat("█", filled))
	accEmpty := lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(strings.Repeat("░", barWidth-filled))
	accProgressBar := accBar + accEmpty

	summaryContent := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n\n%s\n  %s",
		statsStyle.Render("Cards Reviewed:"), valStyle.Render(fmt.Sprintf("%d", m.sessionReviewed)),
		statsStyle.Render("Correct Answers:"), valStyle.Render(fmt.Sprintf("%d", m.sessionCorrect)),
		statsStyle.Render("Session Time:  "), valStyle.Render(duration.Round(time.Second).String()),
		statsStyle.Render("Accuracy Score:"), valStyle.Render(fmt.Sprintf("%.1f%%", accuracy)),
		statsStyle.Render("Efficiency:    "), valStyle.Render(fmt.Sprintf("%.1f cards/min", cardsPerMin)),
		statsStyle.Render("Performance:"), accProgressBar,
	)

	b.WriteString(mainBox.Render(summaryContent) + "\n\n")

	if m.sessionGrades != nil && len(m.sessionGrades) > 0 && m.sessionReviewed > 0 {
		gradeOrder := []core.ReviewGrade{core.GradeAgain, core.GradeHard, core.GradeGood, core.GradeEasy}
		gradeColors := map[core.ReviewGrade]string{
			core.GradeAgain: "196",
			core.GradeHard:  "208",
			core.GradeGood:  "46",
			core.GradeEasy:  "81",
		}
		gradeIcons := map[core.ReviewGrade]string{
			core.GradeAgain: "✗",
			core.GradeHard:  "~",
			core.GradeGood:  "✓",
			core.GradeEasy:  "★",
		}

		distributionBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2).
			Width(maxInt(40, layout.Width-10))

		distContent := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Render("Grade Distribution") + "\n"
		maxCount := 0
		for _, g := range gradeOrder {
			if m.sessionGrades[g] > maxCount {
				maxCount = m.sessionGrades[g]
			}
		}
		for _, g := range gradeOrder {
			count := m.sessionGrades[g]
			pct := 0.0
			if maxCount > 0 {
				pct = float64(count) / float64(maxCount)
			}
			bar := progressBar(15, pct, lipgloss.Color(gradeColors[g]), lipgloss.Color("238"))
			distContent += fmt.Sprintf("  %s %s %-5s %d\n",
				lipgloss.NewStyle().Foreground(lipgloss.Color(gradeColors[g])).Render(gradeIcons[g]),
				bar,
				lipgloss.NewStyle().Foreground(lipgloss.Color(gradeColors[g])).Bold(true).Render(string(g)),
				count)
		}
		b.WriteString(distributionBox.Render(distContent) + "\n\n")
	}

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Press any key or Enter to return to Dashboard.") + "\n")

	motivation := "Keep it up!"
	motivationColor := "213"
	if accuracy >= 90 {
		motivation = "Excellent work! You're mastering this!"
		motivationColor = "46"
	} else if accuracy >= 75 {
		motivation = "Good job! You're making solid progress!"
		motivationColor = "226"
	} else if accuracy >= 50 {
		motivation = "Nice effort! Keep practicing to improve!"
		motivationColor = "214"
	} else if m.sessionReviewed > 0 {
		motivation = "Don't give up! Review helps strengthen memory."
		motivationColor = "196"
	}

	b.WriteString("\n" + lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(motivationColor)).Render(motivation))

	return b.String()
}
