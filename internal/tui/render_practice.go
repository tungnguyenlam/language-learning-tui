package tui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderPractice(layout viewportLayout) string {
	switch m.practiceSubView {
	case PracticeSubViewHub:
		return m.renderPracticeHub(layout)
	case PracticeSubViewGender:
		return fillViewportContent(m.renderGenderTrainer(layout), layout)
	}
	if isGenericTrainer(m.practiceSubView) {
		return fillViewportContent(m.renderTrainer(m.practiceSubView, layout), layout)
	}
	return "Unknown Practice View"
}

type practiceHubMode struct {
	id    string
	key   string
	label string
	desc  string
	sub   PracticeSubView
	color color.Color
	icon  string
}

var allPracticeHubModes = []practiceHubMode{
	{"practice-gender", "1", "Gender Trainer", "Practice der/die/das noun genders", PracticeSubViewGender, colorBlue, "🚻"},
	{"practice-conjugation", "2", "Conjugation Trainer", "Practice German verb forms", PracticeSubViewConjugation, colorPink, "🔄"},
	{"practice-case", "3", "Case Ending Trainer", "Practice Nom/Acc/Dat/Gen endings", PracticeSubViewCase, colorGold, "📐"},
	{"practice-adjective", "4", "Adjective Ending Trainer", "Practice weak/strong/mixed endings", PracticeSubViewAdjective, colorYellow, "🎨"},
	{"practice-preposition", "5", "Preposition Trainer", "Practice two-way prepositions & cases", PracticeSubViewPreposition, colorPurple, "📍"},
	{"practice-plural", "6", "Plural Trainer", "Practice German noun plural forms", PracticeSubViewPlural, colorSecondary, "👥"},
	{"practice-separable", "7", "Separable Verb Trainer", "Practice verb prefixes and word order", PracticeSubViewSeparable, colorOrange, "S/"},
	{"practice-numbers", "8", "Numbers & Time", "Practice German numbers and time", PracticeSubViewNumbers, colorCyan, "🔢"},
	{"practice-conjunctions", "9", "Conjunctions & Word Order", "Practice conjunctions & sentence structure", PracticeSubViewConjunctions, colorGreen, "🔗"},
	{"practice-konjunktiv", "0", "Konjunktiv II Trainer", "Practice subjunctive II (würde, wäre, hätte)", PracticeSubViewKonjunktiv, colorAITitle, "Kj"},
	{"practice-passive", "-", "Passive Voice Trainer", "Practice Vorgangspassiv & Zustandspassiv", PracticeSubViewPassive, colorPurple, "Pv"},
	{"practice-relative", "=", "Relative Clauses Trainer", "Practice Relativpronomen & Relativsätze", PracticeSubViewRelative, colorGreen, "Rl"},
}

func (m *Model) visiblePracticeHubModes() []practiceHubMode {
	if m.practiceFilter == "" {
		return allPracticeHubModes
	}
	q := strings.ToLower(m.practiceFilter)
	var visible []practiceHubMode
	for _, mode := range allPracticeHubModes {
		if strings.Contains(strings.ToLower(mode.label), q) || strings.Contains(strings.ToLower(mode.desc), q) || strings.EqualFold(mode.key, q) {
			visible = append(visible, mode)
		}
	}
	return visible
}

func (m *Model) renderPracticeHub(layout viewportLayout) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	header := lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" PRACTICE HUB ")) + "\n\n"
	headerHeight := 2

	displayModes := m.visiblePracticeHubModes()

	if len(displayModes) > 0 {
		m.practiceHubCursor = clampInt(m.practiceHubCursor, 0, len(displayModes)-1)
	} else {
		m.practiceHubCursor = 0
	}

	var filterBar string
	if m.practiceFilterFocus || m.practiceFilter != "" {
		cursorStr := ""
		if m.practiceFilterFocus {
			cursorStr = "█"
		}
		filterText := fmt.Sprintf("🔍 Filter: %s%s (%d/%d trainers)", m.practiceFilter, cursorStr, len(displayModes), len(allPracticeHubModes))
		filterBar = lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Render(filterText)) + "\n\n"
		headerHeight += 2
	}

	spacing := 5
	listHeight := maxInt(1, layout.Height-headerHeight)

	getItemCount := func(sub PracticeSubView) int {
		if sub == PracticeSubViewGender {
			return len(m.practiceItems)
		}
		if st, ok := m.trainers[sub]; ok {
			return len(st.items)
		}
		return 0
	}

	getScoreStr := func(sub PracticeSubView) string {
		var correct, total int
		if sub == PracticeSubViewGender {
			correct, total = m.practiceCorrect, m.practiceTotal
		} else if st, ok := m.trainers[sub]; ok {
			correct, total = st.correct, st.total
		}
		if total > 0 {
			pct := float64(correct) / float64(total) * 100
			return fmt.Sprintf(" • %d/%d (%.0f%%)", correct, total, pct)
		}
		return ""
	}

	btnWidth := 46
	if layout.Width >= 56 {
		btnWidth = 52
	} else if layout.Width < 50 {
		btnWidth = layout.Width - 4
	}
	if btnWidth < 20 {
		btnWidth = 20
	}

	itemStartLine := m.practiceHubCursor * spacing
	totalContentLines := (len(displayModes) * spacing) + 3
	m.practiceScroll = AutoScroll(itemStartLine, m.practiceScroll, listHeight, totalContentLines)

	var b strings.Builder
	if len(displayModes) == 0 {
		b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("No practice trainers found matching query")) + "\n")
	} else {
		for i, mode := range displayModes {
			borderCol := lipgloss.Color("240") // Default border
			prefix := "  "
			if i == m.practiceHubCursor {
				borderCol = colorPink
				prefix = "▶ "
			}

			btnStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0, 2).
				Width(btnWidth).
				BorderForeground(borderCol)

			count := getItemCount(mode.sub)
			countStr := ""
			if count > 0 {
				countStr = fmt.Sprintf(" (%d)", count)
			}

			scoreStr := getScoreStr(mode.sub)

			titlePlain := truncateLine(
				fmt.Sprintf("%s[%s] %s %s%s", prefix, mode.key, mode.icon, mode.label, countStr),
				btnWidth-4,
			)
			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(mode.color)
			if i == m.practiceHubCursor {
				titleStyle = lipgloss.NewStyle().Bold(true).Foreground(colorPink)
			}
			descText := mode.desc
			if scoreStr != "" {
				descText = strings.TrimSpace(scoreStr) + " — " + mode.desc
			}
			descRendered := mutedStyle.Render(truncateLine(descText, btnWidth-4))
			content := titleStyle.Render(titlePlain) + "\n" + descRendered

			btn := btnStyle.Render(content)
			b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, btn) + "\n")

			btnY := headerHeight + (i * spacing) - m.practiceScroll
			if btnY >= headerHeight && btnY < layout.Height {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     mode.id,
					View:   ViewPractice,
					X:      layout.X + (layout.Width-btnWidth)/2,
					Y:      layout.Y + btnY,
					Width:  btnWidth,
					Height: minInt(spacing, layout.Height-btnY),
				})
			}
		}
	}

	b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press a key to select a trainer") + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("/ Filter  •  r Reset scores  •  Esc Dashboard")))

	listLayout := viewportLayout{
		X:      layout.X,
		Y:      layout.Y + headerHeight,
		Width:  layout.Width,
		Height: listHeight,
	}

	return header + filterBar + AutoScrollViewport(b.String(), listLayout, &m.practiceScroll, "practice", ViewPractice, m)
}
