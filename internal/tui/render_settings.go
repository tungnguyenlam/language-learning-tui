package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/app"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) renderSettings(x, y int) string {
	width, height := m.activePanelSize()
	style := panelStyle.Width(width).Height(height)
	layout := contentLayoutForStyle(style, x, y)

	ctx := NewRenderContext(m, layout, ViewSettings)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("236")).
		Padding(0, 2)
	ctx.WriteLine(titleStyle.Render("⚙ SETTINGS"))
	ctx.NewLine()

	if m.editingTemplate {
		ctx.WriteLine(lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true).Render("EDITING - Enter to save, Esc to cancel."))
	} else if strings.HasPrefix(m.status, "Daily goal set") {
		ctx.WriteLine(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(m.status))
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
	addContent(mutedStyle.Render("  Provider cycle: disabled -> offline -> template -> openai -> anthropic -> ollama."), nil)

	setMap := m.aiTemplates[activeSet]
	aiOptions := []string{
		fmt.Sprintf("AI Provider:    %s", aiProviderName),
		fmt.Sprintf("Dictionary:     %s", m.dictionaryProvider),
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

	goalIdx := 5
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

	audioIdx := 6
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if audioIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	addContent(itemStyle.Render(fmt.Sprintf("%sAuto-play audio: %s", prefix, autoPlayStatus)), &lineInfo{itemIdx: audioIdx, kind: "item"})

	strictIdx := 7
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
	addContent(sectionStyle.Render("DICTIONARY DATA"), nil)
	if m.dictCount > 0 {
		addContent(fmt.Sprintf("  Loaded Entries: %s", lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(fmt.Sprintf("%d", m.dictCount))), nil)
	} else {
		addContent(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("  No dictionary data loaded."), nil)
		addContent(mutedStyle.Render("  To enable offline local dictionary search, download the German-English dataset"), nil)
		addContent(mutedStyle.Render("  from dict.cc, zip it, and place it in the 'local_dict_files/' folder."), nil)
		addContent(mutedStyle.Render("  The app will automatically import it on the next startup."), nil)
	}

	addContent("", nil)
	addContent(sectionStyle.Render("API CREDENTIALS"), nil)
	addContent(mutedStyle.Render("  Press Enter on a row to edit credentials."), nil)

	credRows := []struct {
		idx      int
		provider string
		label    string
		value    string
	}{
		{8, "openai", "OpenAI Key:  ", app.MaskAPIKey(m.aiSecrets.OpenAI.APIKey)},
		{9, "openai", "OpenAI Model:", credValueOrDefault(m.aiSecrets.OpenAI.Model, "openai")},
		{10, "openai", "OpenAI URL:  ", credValueOrDefault(m.aiSecrets.OpenAI.BaseURL, "openai-url")},
		{11, "anthropic", "Anthropic Key:", app.MaskAPIKey(m.aiSecrets.Anthropic.APIKey)},
		{12, "anthropic", "Anthropic Mod:", credValueOrDefault(m.aiSecrets.Anthropic.Model, "anthropic")},
		{13, "anthropic", "Anthropic URL:", credValueOrDefault(m.aiSecrets.Anthropic.BaseURL, "anthropic-url")},
		{14, "ollama", "Ollama Model: ", credValueOrDefault(m.aiSecrets.Ollama.Model, "ollama")},
		{15, "ollama", "Ollama URL:   ", credValueOrDefault(m.aiSecrets.Ollama.BaseURL, "ollama-url")},
	}
	for _, row := range credRows {
		prefix = "  "
		itemStyle := lipgloss.NewStyle()
		if row.idx == m.settingsCursor {
			prefix = "> "
			if m.editingSecretKey != "" && row.provider == m.editingSecretProvider {
				itemStyle = itemStyle.Bold(true).Background(lipgloss.Color("62"))
			} else {
				itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
			}
		}
		displayValue := truncateLine(row.value, maxInt(20, layout.Width-lipgloss.Width(prefix+row.label)-4))
		if m.editingSecretKey != "" && row.idx == m.settingsCursor {
			raw := m.getCredValue(m.editingSecretProvider, m.editingSecretKey)
			if m.editingSecretKey == "api_key" {
				displayValue = fmt.Sprintf("•%d chars (Enter to save, Esc to cancel)", len(raw))
			} else {
				displayValue = raw + "▌"
			}
		}
		addContent(itemStyle.Render(fmt.Sprintf("%s%s%s", prefix, row.label, displayValue)), &lineInfo{itemIdx: row.idx, kind: "item"})
	}

	revealIdx := 16
	prefix = "  "
	itemStyle = lipgloss.NewStyle()
	if revealIdx == m.settingsCursor {
		prefix = "> "
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("212"))
	}
	revealValue := fmt.Sprintf("%d", m.revealSpeed)
	if m.revealSpeed == 0 {
		revealValue = "Instant"
	}
	revealLabel := fmt.Sprintf("%sReveal Speed: %s ", prefix, revealValue)
	var revealLine strings.Builder
	revealLine.WriteString(itemStyle.Render(revealLabel))
	if m.revealSpeed <= 0 {
		revealLine.WriteString(disabledBtnStyle.Render(minusBtn))
	} else {
		revealLine.WriteString(btnStyle.Render(minusBtn))
	}
	if m.revealSpeed >= 10 {
		revealLine.WriteString(disabledBtnStyle.Render(plusBtn))
	} else {
		revealLine.WriteString(btnStyle.Render(plusBtn))
	}
	addContent(revealLine.String(), &lineInfo{itemIdx: revealIdx, kind: "reveal"})

	if !m.editingTemplate && m.editingSecretKey == "" {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
		addContent("", nil)
		addContent(fmt.Sprintf("Use %s/%s to move, %s on goal or speed, Enter edit/toggle, %s/%s templates.",
			keyStyle.Render("j"), keyStyle.Render("k"), keyStyle.Render("+/-"), keyStyle.Render("["), keyStyle.Render("]")), nil)
	}

	contentStr := content.String()
	allLines := strings.Split(strings.TrimRight(contentStr, "\n"), "\n")
	m.settingsTotalLines = len(allLines)

	availableHeight := layout.Height - (ctx.currY - layout.Y)
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
	m.settingsScroll = AutoScroll(cursorLine, m.settingsScroll, availableHeight, m.settingsTotalLines)

	listView := m.RenderList(layout.WithHeight(availableHeight).WithY(ctx.currY), contentStr, ListOptions{
		HitboxPrefix: "settings",
		View:         ViewSettings,
		ScrollOffset: &m.settingsScroll,
		TotalLines:   &m.settingsTotalLines,
		OnLine: func(lineIndex int, lineCtx *RenderContext, content string) {
			info, ok := lineMeta[lineIndex]
			if !ok {
				return
			}
			if info.kind == "goal" {
				lineCtx.RegisterHitboxWithAction(fmt.Sprintf("settings-%d", info.itemIdx), lipgloss.Width(goalLabel), 1, func() tea.Cmd {
					m.settingsCursor = info.itemIdx
					return m.handleSettingsEnter()
				})
				lineCtx.RegisterHitboxAtWithAction("settings-goal-minus", lipgloss.Width(goalLabel), 0, lipgloss.Width(minusBtn), 1, func() tea.Cmd {
					return m.setDailyGoal(m.stats.DailyGoal - 1)
				})
				lineCtx.RegisterHitboxAtWithAction("settings-goal-plus", lipgloss.Width(goalLabel)+lipgloss.Width(minusBtn), 0, lipgloss.Width(plusBtn), 1, func() tea.Cmd {
					return m.setDailyGoal(m.stats.DailyGoal + 1)
				})
			} else if info.kind == "reveal" {
				lineCtx.RegisterHitboxWithAction(fmt.Sprintf("settings-%d", info.itemIdx), lipgloss.Width(revealLabel), 1, func() tea.Cmd {
					m.settingsCursor = info.itemIdx
					return m.handleSettingsEnter()
				})
				lineCtx.RegisterHitboxAtWithAction("settings-reveal-minus", lipgloss.Width(revealLabel), 0, lipgloss.Width(minusBtn), 1, func() tea.Cmd {
					return m.setRevealSpeed(m.revealSpeed - 1)
				})
				lineCtx.RegisterHitboxAtWithAction("settings-reveal-plus", lipgloss.Width(revealLabel)+lipgloss.Width(minusBtn), 0, lipgloss.Width(plusBtn), 1, func() tea.Cmd {
					return m.setRevealSpeed(m.revealSpeed + 1)
				})
			} else if info.kind == "item" {
				lineCtx.RegisterHitboxWithAction(fmt.Sprintf("settings-%d", info.itemIdx), layout.Width-2, 1, func() tea.Cmd {
					m.settingsCursor = info.itemIdx
					return m.handleSettingsEnter()
				})
			}
		},
	})

	ctx.Write(listView)
	return ctx.String()
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
		return "(default: claude-3-5-haiku-latest)"
	case "openai-url":
		return "(default: api.openai.com/v1)"
	case "anthropic-url":
		return "(default: api.anthropic.com)"
	case "ollama":
		return "(default: llama3.1)"
	case "ollama-url":
		return "(default: http://localhost:11434)"
	}
	return "(not set)"
}
