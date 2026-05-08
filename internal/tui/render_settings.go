package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderSettings(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("159")).MarginBottom(1)
	b.WriteString(titleStyle.Render("Settings") + "\n\n")

	autoPlayStatus := "off"
	if m.autoPlayAudio {
		autoPlayStatus = "on"
	}

	aiProviderName := m.aiProviderName
	activeSet := "none"
	if set := m.currentAITemplateSet(); set != "" {
		activeSet = set
	}

	sectionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true).MarginTop(1)
	b.WriteString(sectionStyle.Render("AI CONFIGURATION") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("  Template Set: %s", activeSet)) + "\n\n")

	setMap := m.aiTemplates[activeSet]
	aiOptions := []string{
		fmt.Sprintf("AI Provider:    %s", aiProviderName),
		fmt.Sprintf("Front Template: %s", setMap["front"]),
		fmt.Sprintf("Back Template:  %s", setMap["back"]),
		fmt.Sprintf("Example Tmpl:   %s", setMap["example"]),
	}
	var rowY int
	var prefix string
	for i, opt := range aiOptions {
		prefix = "  "
		itemStyle := lipgloss.NewStyle()
		if i == m.settingsCursor {
			prefix = "> "
			if m.editingTemplate {
				itemStyle = itemStyle.Bold(true).Background(lipgloss.Color("62"))
			} else {
				itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
			}
		}
		item := prefix + opt
		rowY = layout.Y + strings.Count(b.String(), "\n")
		b.WriteString(itemStyle.Render(item) + "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("settings-%d", i),
			View:   ViewSettings,
			X:      layout.X,
			Y:      rowY,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
	}

	b.WriteString("\n" + sectionStyle.Render("STUDY PREFERENCES") + "\n")

	goalIdx := 4
	prefix = "  "
	itemStyle := lipgloss.NewStyle()
	if goalIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	goalLabel := fmt.Sprintf("%sDaily Goal: %d cards ", prefix, m.stats.DailyGoal)
	rowY = layout.Y + strings.Count(b.String(), "\n")
	b.WriteString(itemStyle.Render(goalLabel))

	btnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	minusBtn := "[-] "
	plusBtn := "[+] "
	b.WriteString(btnStyle.Render(minusBtn))
	b.WriteString(btnStyle.Render(plusBtn))
	b.WriteString("\n")

	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     fmt.Sprintf("settings-%d", goalIdx),
		View:   ViewSettings,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(goalLabel),
		Height: 1,
	})
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "settings-goal-minus",
		View:   ViewSettings,
		X:      layout.X + lipgloss.Width(goalLabel),
		Y:      rowY,
		Width:  lipgloss.Width(minusBtn),
		Height: 1,
	})
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "settings-goal-plus",
		View:   ViewSettings,
		X:      layout.X + lipgloss.Width(goalLabel) + lipgloss.Width(minusBtn),
		Y:      rowY,
		Width:  lipgloss.Width(plusBtn),
		Height: 1,
	})

	audioIdx := 5
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if audioIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	audioItem := fmt.Sprintf("%sAuto-play audio: %s", prefix, autoPlayStatus)
	rowY = layout.Y + strings.Count(b.String(), "\n")
	b.WriteString(itemStyle.Render(audioItem) + "\n")
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     fmt.Sprintf("settings-%d", audioIdx),
		View:   ViewSettings,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(audioItem),
		Height: 1,
	})

	if m.editingTemplate {
		b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("EDITING") + " - Enter to save, Esc to cancel.")
	} else {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(fmt.Sprintf("\nUse %s/%s to move, %s goal, Enter edit/toggle, %s/%s templates.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("+/-"), keyStyle.Render("["), keyStyle.Render("]")))
	}
	return b.String()
}
