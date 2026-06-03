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
	}

	spacing := 3
	if layout.Height >= 32 {
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
		}
		return 0
	}

	for i, mode := range modes {
		btnStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 2).
			Width(40).
			BorderForeground(lipgloss.Color("81"))

		count := getItemCount(mode.sub)
		countStr := ""
		if count > 0 {
			countStr = fmt.Sprintf(" (%d)", count)
		}

		keyHint := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Render("[" + mode.key + "]")
		content := fmt.Sprintf("%s %s%s\n%s", keyHint, mode.label, countStr, mutedStyle.Render(mode.desc))
		if spacing == 3 {
			content = fmt.Sprintf("%s %s%s - %s", keyHint, mode.label, countStr, mutedStyle.Render(mode.desc))
			btnStyle = btnStyle.Padding(0, 1)
		}

		btn := btnStyle.Render(content)
		b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, btn) + "\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     mode.id,
			View:   ViewPractice,
			X:      layout.X + (layout.Width-40)/2,
			Y:      layout.Y + 3 + (i * spacing),
			Width:  40,
			Height: spacing - 1,
		})
	}

	b.WriteString("\n" + lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, "Press a key to select a trainer") + "\n")
	b.WriteString(lipgloss.PlaceHorizontal(layout.Width, lipgloss.Center, mutedStyle.Render("Press Esc to return to Dashboard")))

	return b.String()
}
