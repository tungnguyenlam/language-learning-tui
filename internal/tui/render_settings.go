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

	var allLines []string
	var itemLines = make(map[int]int) // cursor index -> line index

	addLine := func(s string) {
		lines := strings.Split(s, "\n")
		for _, l := range lines {
			allLines = append(allLines, l)
		}
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("236")).
		Padding(0, 2)
	addLine(titleStyle.Render("⚙ SETTINGS"))
	addLine("")

	if m.editingTemplate {
		addLine(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("EDITING - Enter to save, Esc to cancel."))
	} else if strings.HasPrefix(m.status, "Daily goal set") {
		addLine(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(m.status))
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
		Padding(0, 1)

	addLine(sectionStyle.Render("AI CONFIGURATION"))
	addLine(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("  Template Set: %s", activeSet)))
	addLine(mutedStyle.Render("  Provider cycle: disabled -> offline -> template -> openai -> anthropic."))

	setMap := m.aiTemplates[activeSet]
	aiOptions := []string{
		fmt.Sprintf("AI Provider:    %s", aiProviderName),
		fmt.Sprintf("Front Template: %s", strings.ReplaceAll(setMap["front"], "\n", "\\n")),
		fmt.Sprintf("Back Template:  %s", strings.ReplaceAll(setMap["back"], "\n", "\\n")),
		fmt.Sprintf("Example Tmpl:   %s", strings.ReplaceAll(setMap["example"], "\n", "\\n")),
	}

	for i, opt := range aiOptions {
		prefix := "  "
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
		itemLines[i] = len(allLines)
		addLine(itemStyle.Render(item))
	}

	addLine(sectionStyle.Render("STUDY PREFERENCES"))

	goalIdx := 4
	prefix := "  "
	itemStyle := lipgloss.NewStyle()
	if goalIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	goalLabel := fmt.Sprintf("%sDaily Goal: %d cards ", prefix, m.stats.DailyGoal)

	btnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	disabledBtnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)
	minusBtn := "[-] "
	plusBtn := "[+] "

	var goalLine strings.Builder
	goalLine.WriteString(itemStyle.Render(goalLabel))
	if m.stats.DailyGoal <= 1 {
		goalLine.WriteString(disabledBtnStyle.Render(minusBtn))
	} else {
		goalLine.WriteString(btnStyle.Render(minusBtn))
	}
	goalLine.WriteString(btnStyle.Render(plusBtn))

	itemLines[goalIdx] = len(allLines)
	addLine(goalLine.String())

	audioIdx := 5
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if audioIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	audioItem := fmt.Sprintf("%sAuto-play audio: %s", prefix, autoPlayStatus)
	itemLines[audioIdx] = len(allLines)
	addLine(itemStyle.Render(audioItem))

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
	itemLines[strictIdx] = len(allLines)
	addLine(itemStyle.Render(strictItem))

	addLine("")
	addLine(sectionStyle.Render("API CREDENTIALS"))
	credProvider := m.credProviderName()
	disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	if credProvider == "" {
		addLine(disabledStyle.Render("  (select openai or anthropic above to edit credentials)"))
	} else {
		addLine(mutedStyle.Render(fmt.Sprintf("  Provider: %s. Press Enter on a row to edit.", credProvider)))
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
			raw := m.getCredValue(m.editingSecretProvider, m.editingSecretKey)
			if m.editingSecretKey == "api_key" {
				displayValue = fmt.Sprintf("•%d chars (Enter to save, Esc to cancel)", len(raw))
			} else {
				displayValue = raw + "▌"
			}
		}
		item := fmt.Sprintf("%s%s%s", prefix, row.label, displayValue)
		itemLines[row.idx] = len(allLines)
		addLine(itemStyle.Render(item))
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	addLine("")
	addLine(fmt.Sprintf("Color Theme: %s (Press %s to cycle)", m.theme, keyStyle.Render("c")))

	if !m.editingTemplate && m.editingSecretKey == "" {
		keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		addLine("")
		addLine(fmt.Sprintf("Use %s/%s to move, %s goal, %s theme, Enter edit/toggle, %s/%s templates.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("+/-"), keyStyle.Render("c"), keyStyle.Render("["), keyStyle.Render("]")))
	}

	// Scrolling logic
	totalLines := len(allLines)
	m.settingsTotalLines = totalLines
	maxVisible := layout.Height
	if totalLines > maxVisible {
		// Ensure cursor is visible
		cursorLine := itemLines[m.settingsCursor]
		if cursorLine < m.settingsScroll {
			m.settingsScroll = cursorLine
		} else if cursorLine >= m.settingsScroll+maxVisible {
			m.settingsScroll = cursorLine - maxVisible + 1
		}
	} else {
		m.settingsScroll = 0
	}

	m.settingsScroll = clampInt(m.settingsScroll, 0, maxInt(0, totalLines-maxVisible))

	var b strings.Builder
	padWidth := layout.Width - 2
	thumbStart, thumbHeight := scrollbarThumb(totalLines, maxVisible, m.settingsScroll)

	for i := m.settingsScroll; i < m.settingsScroll+maxVisible && i < totalLines; i++ {
		line := allLines[i]
		displayLine := truncateLine(line, padWidth)
		lineWidth := lipgloss.Width(displayLine)

		if lineWidth < padWidth {
			displayLine += strings.Repeat(" ", padWidth-lineWidth)
		}

		if totalLines > maxVisible {
			scrollbarChar := "│"
			currentPos := i - m.settingsScroll
			if currentPos >= thumbStart && currentPos < thumbStart+thumbHeight {
				scrollbarChar = "█"
			}
			displayLine += " " + scrollbarChar

			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("settings-scroll-%d", currentPos),
				View:   ViewSettings,
				X:      layout.X + padWidth + 1,
				Y:      layout.Y + currentPos,
				Width:  1,
				Height: 1,
			})
		}

		b.WriteString(displayLine + "\n")

		// Re-add hitboxes for selectable items
		for itemIdx, lineIdx := range itemLines {
			if i == lineIdx {
				// Special case for goal buttons
				if itemIdx == 4 {
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("settings-%d", itemIdx),
						View:   ViewSettings,
						X:      layout.X,
						Y:      layout.Y + (i - m.settingsScroll),
						Width:  lipgloss.Width(goalLabel),
						Height: 1,
					})
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     "settings-goal-minus",
						View:   ViewSettings,
						X:      layout.X + lipgloss.Width(goalLabel),
						Y:      layout.Y + (i - m.settingsScroll),
						Width:  lipgloss.Width(minusBtn),
						Height: 1,
					})
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     "settings-goal-plus",
						View:   ViewSettings,
						X:      layout.X + lipgloss.Width(goalLabel) + lipgloss.Width(minusBtn),
						Y:      layout.Y + (i - m.settingsScroll),
						Width:  lipgloss.Width(plusBtn),
						Height: 1,
					})
				} else {
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("settings-%d", itemIdx),
						View:   ViewSettings,
						X:      layout.X,
						Y:      layout.Y + (i - m.settingsScroll),
						Width:  padWidth,
						Height: 1,
					})
				}
			}
		}
	}

	return b.String()
}

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
