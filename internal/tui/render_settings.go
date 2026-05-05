package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderSettings(x, y int) string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Settings") + "\n\n")

	autoPlayStatus := "off"
	if m.autoPlayAudio {
		autoPlayStatus = "on"
	}

	aiProviderName := m.aiProviderName
	activeSet := "none"
	if len(m.aiTemplateSets) > 0 {
		activeSet = m.aiTemplateSets[m.aiTemplateIndex]
	}

	sectionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true).MarginTop(1)
	b.WriteString(sectionStyle.Render("AI CONFIGURATION") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("  Template Set: %s", activeSet)) + "\n")

	setMap := m.aiTemplates[activeSet]
	aiOptions := []string{
		fmt.Sprintf("AI Provider: %s", aiProviderName),
		fmt.Sprintf("Front Template: %s", setMap["front"]),
		fmt.Sprintf("Back Template: %s", setMap["back"]),
		fmt.Sprintf("Example Template: %s", setMap["example"]),
	}
	for i, opt := range aiOptions {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == m.settingsCursor {
			prefix = "> "
			if m.editingTemplate {
				style = style.Bold(true).Background(lipgloss.Color("62"))
			} else {
				style = style.Bold(true).Foreground(lipgloss.Color("212"))
			}
		}
		item := prefix + opt
		rowY := y + strings.Count(b.String(), "\n")
		b.WriteString(style.Render(item) + "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("settings-%d", i),
			View:   ViewSettings,
			X:      x,
			Y:      rowY,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
	}

	b.WriteString(sectionStyle.Render("STUDY PREFERENCES") + "\n")

	goalIdx := 4
	prefix := "  "
	style := lipgloss.NewStyle()
	if goalIdx == m.settingsCursor {
		prefix = "> "
		style = style.Bold(true).Foreground(lipgloss.Color("212"))
	}
	goalLabel := fmt.Sprintf("%sDaily Goal: %d ", prefix, m.stats.DailyGoal)
	rowY := y + strings.Count(b.String(), "\n")
	b.WriteString(style.Render(goalLabel))

	btnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	minusBtn := "[-] "
	plusBtn := "[+] "
	b.WriteString(btnStyle.Render(minusBtn))
	b.WriteString(btnStyle.Render(plusBtn))
	b.WriteString("\n")

	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     fmt.Sprintf("settings-%d", goalIdx),
		View:   ViewSettings,
		X:      x,
		Y:      rowY,
		Width:  lipgloss.Width(goalLabel),
		Height: 1,
	})
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "settings-goal-minus",
		View:   ViewSettings,
		X:      x + lipgloss.Width(goalLabel),
		Y:      rowY,
		Width:  lipgloss.Width(minusBtn),
		Height: 1,
	})
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "settings-goal-plus",
		View:   ViewSettings,
		X:      x + lipgloss.Width(goalLabel) + lipgloss.Width(minusBtn),
		Y:      rowY,
		Width:  lipgloss.Width(plusBtn),
		Height: 1,
	})

	audioIdx := 5
	prefix = "  "
	style = lipgloss.NewStyle()
	if audioIdx == m.settingsCursor {
		prefix = "> "
		style = style.Bold(true).Foreground(lipgloss.Color("212"))
	}
	audioItem := fmt.Sprintf("%sAuto-play audio: %s", prefix, autoPlayStatus)
	rowY = y + strings.Count(b.String(), "\n")
	b.WriteString(style.Render(audioItem) + "\n")
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     fmt.Sprintf("settings-%d", audioIdx),
		View:   ViewSettings,
		X:      x,
		Y:      rowY,
		Width:  lipgloss.Width(audioItem),
		Height: 1,
	})

	if m.editingTemplate {
		b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("EDITING") + " - Enter to save, Esc to cancel.")
	} else {
		b.WriteString("\n" + mutedStyle.Render("Use j/k to move, +/- to adjust goal, Enter to toggle/edit."))
	}
	return b.String()
}
