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
		Padding(0, 2)
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
		Padding(0, 1)

	var content strings.Builder
	type lineInfo struct {
		itemIdx int
		kind    string // "item", "goal", "other"
	}
	var lineMeta = make(map[int]lineInfo)

	addContent := func(s string, info *lineInfo) {
		lines := strings.Split(s, "\n")
		for _, l := range lines {
			if info != nil {
				lineMeta[strings.Count(content.String(), "\n")] = *info
			}
			content.WriteString(l + "\n")
		}
	}

	addContent(sectionStyle.Render("AI CONFIGURATION"), nil)
	addContent(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("  Template Set: %s", activeSet)), nil)
	addContent(mutedStyle.Render("  Provider cycle: disabled -> offline -> template -> openai -> anthropic."), nil)

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
		addContent(itemStyle.Render(prefix+opt), &lineInfo{itemIdx: i, kind: "item"})
	}

	addContent(sectionStyle.Render("STUDY PREFERENCES"), nil)

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
	addContent(goalLine.String(), &lineInfo{itemIdx: goalIdx, kind: "goal"})

	audioIdx := 5
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if audioIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	addContent(itemStyle.Render(fmt.Sprintf("%sAuto-play audio: %s", prefix, autoPlayStatus)), &lineInfo{itemIdx: audioIdx, kind: "item"})

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
	addContent(itemStyle.Render(fmt.Sprintf("%sStrict Normalization (ss vs ß): %s", prefix, strictStatus)), &lineInfo{itemIdx: strictIdx, kind: "item"})

	addContent("", nil)
	addContent(sectionStyle.Render("API CREDENTIALS"), nil)
	credProvider := m.credProviderName()
	disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	if credProvider == "" {
		addContent(disabledStyle.Render("  (select openai or anthropic above to edit credentials)"), nil)
	} else {
		addContent(mutedStyle.Render(fmt.Sprintf("  Provider: %s. Press Enter on a row to edit.", credProvider)), nil)
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
		addContent(itemStyle.Render(fmt.Sprintf("%s%s%s", prefix, row.label, displayValue)), &lineInfo{itemIdx: row.idx, kind: "item"})
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	addContent("", nil)
	addContent(fmt.Sprintf("Color Theme: %s (Press %s to cycle)", m.theme, keyStyle.Render("c")), nil)

	if !m.editingTemplate && m.editingSecretKey == "" {
		keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		addContent("", nil)
		addContent(fmt.Sprintf("Use %s/%s to move, %s goal, %s theme, Enter edit/toggle, %s/%s templates.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("+/-"), keyStyle.Render("c"), keyStyle.Render("["), keyStyle.Render("]")), nil)
	}

	contentStr := content.String()
	allLines := strings.Split(strings.TrimRight(contentStr, "\n"), "\n")
	m.settingsTotalLines = len(allLines)

	headerLines := strings.Count(b.String(), "\n")
	availableHeight := layout.Height - headerLines
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Auto-scroll to cursor
	// Find line index for current cursor
	cursorLine := 0
	for lineIdx, info := range lineMeta {
		if info.itemIdx == m.settingsCursor {
			cursorLine = lineIdx
			break
		}
	}
	if cursorLine < m.settingsScroll {
		m.settingsScroll = cursorLine
	} else if cursorLine >= m.settingsScroll+availableHeight {
		m.settingsScroll = cursorLine - availableHeight + 1
	}
	m.settingsScroll = clampInt(m.settingsScroll, 0, maxInt(0, m.settingsTotalLines-availableHeight))

	listView := m.renderScrollable(layout.WithHeight(availableHeight), contentStr, m.settingsScroll, scrollableOptions{
		hitboxPrefix: "settings",
		view:         ViewSettings,
		onLine: func(lineIndex int, rY int, content string) {
			info, ok := lineMeta[lineIndex]
			if !ok {
				return
			}
			if info.kind == "goal" {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("settings-%d", info.itemIdx),
					View:   ViewSettings,
					X:      layout.X,
					Y:      rY,
					Width:  lipgloss.Width(goalLabel),
					Height: 1,
				})
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "settings-goal-minus",
					View:   ViewSettings,
					X:      layout.X + lipgloss.Width(goalLabel),
					Y:      rY,
					Width:  lipgloss.Width(minusBtn),
					Height: 1,
				})
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "settings-goal-plus",
					View:   ViewSettings,
					X:      layout.X + lipgloss.Width(goalLabel) + lipgloss.Width(minusBtn),
					Y:      rY,
					Width:  lipgloss.Width(plusBtn),
					Height: 1,
				})
			} else if info.kind == "item" {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("settings-%d", info.itemIdx),
					View:   ViewSettings,
					X:      layout.X,
					Y:      rY,
					Width:  layout.Width - 2,
					Height: 1,
				})
			}
		},
	})

	b.WriteString(listView)
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
