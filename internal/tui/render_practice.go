package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) renderPractice(layout viewportLayout) string {
	switch m.practiceSubView {
	case PracticeSubViewHub:
		return m.renderPracticeHub(layout)
	case PracticeSubViewGender:
		return m.renderGenderTrainer(layout)
	case PracticeSubViewConjugation:
		return m.renderConjugation(layout)
	case PracticeSubViewCase:
		return m.renderCaseTrainer(layout)
	case PracticeSubViewAdjective:
		return m.renderAdjectiveTrainer(layout)
	case PracticeSubViewPreposition:
		return m.renderPrepositionTrainer(layout)
	case PracticeSubViewPlural:
		return m.renderPluralTrainer(layout)
	case PracticeSubViewSeparable:
		return m.renderSeparableTrainer(layout)
	case PracticeSubViewNumbers:
		return m.renderNumberTrainer(layout)
	case PracticeSubViewConjunctions:
		return m.renderConjunctionTrainer(layout)
	}
	return "Unknown Practice View"
}

func (m *Model) renderPracticeHub(layout viewportLayout) string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, titleStyle.Render(" PRACTICE HUB ")) + "\n\n")

	modes := []struct {
		id    string
		key   string
		label string
		desc  string
		sub   PracticeSubView
	}{
		{"practice-gender", "1", "Gender Trainer", "Practice der/die/das noun genders", PracticeSubViewGender},
		{"practice-conjugation", "2", "Conjugation Trainer", "Practice German verb forms", PracticeSubViewConjugation},
		{"practice-case", "3", "Case Ending Trainer", "Practice Nom/Acc/Dat/Gen endings", PracticeSubViewCase},
		{"practice-adjective", "4", "Adjective Ending Trainer", "Practice Nom/Acc/Dat/Gen endings", PracticeSubViewAdjective},
		{"practice-preposition", "5", "Preposition Trainer", "Practice two-way prepositions & cases", PracticeSubViewPreposition},
		{"practice-plural", "6", "Plural Trainer", "Practice German noun plural forms", PracticeSubViewPlural},
		{"practice-separable", "7", "Separable Verb Trainer", "Practice verb prefixes and word order", PracticeSubViewSeparable},
		{"practice-numbers", "8", "Numbers & Time", "Practice German numbers and time", PracticeSubViewNumbers},
		{"practice-conjunctions", "9", "Conjunctions & Word Order", "Practice German conjunctions & sentence structure", PracticeSubViewConjunctions},
	}

	spacing := 3
	if layout.Height >= 36 {
		spacing = 4
	}

	getItemCount := func(sub PracticeSubView) int {
		switch sub {
		case PracticeSubViewGender:
			return len(m.practiceItems)
		case PracticeSubViewConjugation:
			return len(m.conjugationItems)
		case PracticeSubViewCase:
			return len(m.caseItems)
		case PracticeSubViewAdjective:
			return len(m.adjItems)
		case PracticeSubViewPreposition:
			return len(m.prepItems)
		case PracticeSubViewPlural:
			return len(m.pluralItems)
		case PracticeSubViewSeparable:
			return len(m.separableItems)
		case PracticeSubViewNumbers:
			return len(m.numberItems)
		case PracticeSubViewConjunctions:
			return len(m.conjItems)
		}
		return 0
	}

	getScoreStr := func(sub PracticeSubView) string {
		var correct, total int
		switch sub {
		case PracticeSubViewGender:
			correct, total = m.practiceCorrect, m.practiceTotal
		case PracticeSubViewConjugation:
			correct, total = m.conjugationCorrect, m.conjugationTotal
		case PracticeSubViewCase:
			correct, total = m.caseCorrect, m.caseTotal
		case PracticeSubViewAdjective:
			correct, total = m.adjCorrect, m.adjTotal
		case PracticeSubViewPreposition:
			correct, total = m.prepCorrect, m.prepTotal
		case PracticeSubViewPlural:
			correct, total = m.pluralCorrect, m.pluralTotal
		case PracticeSubViewSeparable:
			correct, total = m.separableCorrect, m.separableTotal
		case PracticeSubViewNumbers:
			correct, total = m.numberCorrect, m.numberTotal
		case PracticeSubViewConjunctions:
			correct, total = m.conjCorrect, m.conjTotal
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

	for i, mode := range modes {
		btnStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 2).
			Width(btnWidth).
			BorderForeground(lipgloss.Color("81"))

		count := getItemCount(mode.sub)
		countStr := ""
		if count > 0 {
			countStr = fmt.Sprintf(" (%d)", count)
		}

		scoreStr := getScoreStr(mode.sub)
		scoreRendered := ""
		if scoreStr != "" {
			scoreRendered = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(scoreStr)
		}

		keyHint := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render("[" + mode.key + "]")
		content := fmt.Sprintf("%s %s%s%s\n%s", keyHint, mode.label, countStr, scoreRendered, mutedStyle.Render(mode.desc))
		if spacing == 3 {
			content = fmt.Sprintf("%s %s%s%s - %s", keyHint, mode.label, countStr, scoreRendered, mutedStyle.Render(mode.desc))
			btnStyle = btnStyle.Padding(0, 1)
		}

		btn := btnStyle.Render(content)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, btn) + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     mode.id,
			View:   ViewPractice,
			X:      layout.X + (layout.Width-btnWidth)/2,
			Y:      layout.Y + 3 + (i * spacing),
			Width:  btnWidth,
			Height: spacing - 1,
		})
	}

	b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press a key to select a trainer") + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("Press Esc to return to Dashboard")))

	return b.String()
}
