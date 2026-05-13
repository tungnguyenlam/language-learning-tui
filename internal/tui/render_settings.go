package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/app"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderSettings(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	var b strings.Builder
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("236")).
		Padding(0, 2).
		MarginBottom(1)
	b.WriteString(titleStyle.Render("⚙ SETTINGS") + "\n\n")
	if m.editingTemplate {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("EDITING - Enter to save, Esc to cancel.") + "\n")
	} else if strings.HasPrefix(m.status, "Daily goal set") {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(m.status) + "\n")
	}

	autoPlayStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("off")
	if m.autoPlayAudio {
		autoPlayStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("on")
	}

	aiProviderName := m.aiProviderName
	activeSet := "none"
	if set := m.currentAITemplateSet(); set != "" {
		activeSet = set
	}

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("62")).
		Bold(true).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)
	b.WriteString(sectionStyle.Render("AI CONFIGURATION") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("  Template Set: %s", activeSet)) + "\n\n")
	b.WriteString(mutedStyle.Render("  Provider cycle: disabled -> offline -> template -> openai -> anthropic.") + "\n")

	setMap := m.aiTemplates[activeSet]
	aiOptions := []string{
		fmt.Sprintf("AI Provider:    %s", aiProviderName),
		fmt.Sprintf("Front Template: %s", strings.ReplaceAll(setMap["front"], "\n", "\\n")),
		fmt.Sprintf("Back Template:  %s", strings.ReplaceAll(setMap["back"], "\n", "\\n")),
		fmt.Sprintf("Example Tmpl:   %s", strings.ReplaceAll(setMap["example"], "\n", "\\n")),
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
	disabledBtnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)
	minusBtn := "[-] "
	plusBtn := "[+] "

	if m.stats.DailyGoal <= 1 {
		b.WriteString(disabledBtnStyle.Render(minusBtn))
	} else {
		b.WriteString(btnStyle.Render(minusBtn))
	}

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

	strictIdx := 6
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if strictIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	strictStatus := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("off")
	if m.strictNormalization {
		strictStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("on")
	}
	strictItem := fmt.Sprintf("%sStrict Normalization (ss vs ß): %s", prefix, strictStatus)
	rowY = layout.Y + strings.Count(b.String(), "\n")
	b.WriteString(itemStyle.Render(strictItem) + "\n")
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     fmt.Sprintf("settings-%d", strictIdx),
		View:   ViewSettings,
		X:      layout.X,
		Y:      rowY,
		Width:  lipgloss.Width(strictItem),
		Height: 1,
	})

	b.WriteString("\n" + sectionStyle.Render("API CREDENTIALS") + "\n")
	credProvider := m.credProviderName()
	disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	if credProvider == "" {
		b.WriteString(disabledStyle.Render("  (select openai or anthropic above to edit credentials)") + "\n")
	} else {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("  Provider: %s. Press Enter on a row to edit.", credProvider)) + "\n")
	}

	creds := m.credsFor(credProvider)
	credRows := []struct {
		idx   int
		label string
		value string
	}{
		{7, "API Key:    ", app.MaskAPIKey(creds.APIKey)},
		{8, "Model:      ", credValueOrDefault(creds.Model, credProvider)},
		{9, "Base URL:   ", credValueOrDefault(creds.BaseURL, credProvider+"-url")},
	}
	for _, row := range credRows {
		prefix = "  "
		itemStyle = lipgloss.NewStyle()
		if row.idx == m.settingsCursor {
			prefix = "> "
			if m.editingSecretKey != "" && credKeyForCursor(row.idx) == m.editingSecretKey {
				itemStyle = itemStyle.Bold(true).Background(lipgloss.Color("62"))
			} else if credProvider == "" {
				itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("240"))
			} else {
				itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
			}
		} else if credProvider == "" {
			itemStyle = itemStyle.Foreground(lipgloss.Color("240"))
		}
		displayValue := row.value
		if m.editingSecretKey != "" && credKeyForCursor(row.idx) == m.editingSecretKey {
			// Show what the user is currently typing.
			raw := m.getCredValue(m.editingSecretProvider, m.editingSecretKey)
			if m.editingSecretKey == "api_key" {
				// Mask while typing too — show count of characters.
				displayValue = fmt.Sprintf("•%d chars (Enter to save, Esc to cancel)", len(raw))
			} else {
				displayValue = raw + "▌"
			}
		}
		item := fmt.Sprintf("%s%s%s", prefix, row.label, displayValue)
		rowY = layout.Y + strings.Count(b.String(), "\n")
		b.WriteString(itemStyle.Render(item) + "\n")
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("settings-%d", row.idx),
			View:   ViewSettings,
			X:      layout.X,
			Y:      rowY,
			Width:  lipgloss.Width(item),
			Height: 1,
		})
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	b.WriteString(fmt.Sprintf("\nColor Theme: %s (Press %s to cycle)\n", m.theme, keyStyle.Render("c")))

	if !m.editingTemplate && m.editingSecretKey == "" {
		keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		b.WriteString(fmt.Sprintf("\nUse %s/%s to move, %s goal, %s theme, Enter edit/toggle, %s/%s templates.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("+/-"), keyStyle.Render("c"), keyStyle.Render("["), keyStyle.Render("]")))
	}
	return b.String()
}

// credValueOrDefault returns value if non-empty, else a short hint of the
// provider's default. We avoid hardcoding the actual URLs/models here so
// the source of truth stays in internal/ai.
func credValueOrDefault(value, kind string) string {
	v := strings.TrimSpace(value)
	if v != "" {
		return v
	}
	switch kind {
	case "openai":
		return "(default: gpt-4o-mini)"
	case "anthropic":
		return "(default: claude-haiku-4-5)"
	case "openai-url":
		return "(default: api.openai.com/v1)"
	case "anthropic-url":
		return "(default: api.anthropic.com)"
	}
	return "(not set)"
}
