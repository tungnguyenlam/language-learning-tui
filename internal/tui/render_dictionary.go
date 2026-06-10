package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

var dictHighlightStyle = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)

var (
	genderMascStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	genderFemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	genderNeutStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)

func renderGender(gender string) string {
	if gender == "" {
		return ""
	}
	lower := strings.ToLower(strings.TrimSpace(gender))
	var style lipgloss.Style
	if strings.HasPrefix(lower, "m") || strings.HasPrefix(lower, "der") {
		style = genderMascStyle
	} else if strings.HasPrefix(lower, "f") || strings.HasPrefix(lower, "die") {
		style = genderFemStyle
	} else if strings.HasPrefix(lower, "n") || strings.HasPrefix(lower, "das") {
		style = genderNeutStyle
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	}
	return style.Render("{" + gender + "}")
}

func highlightQuery(text, query string, style lipgloss.Style) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return text
	}
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	var result strings.Builder
	pos := 0
	for {
		idx := strings.Index(lowerText[pos:], lowerQuery)
		if idx == -1 {
			result.WriteString(text[pos:])
			break
		}
		// Write text before match
		result.WriteString(text[pos : pos+idx])
		// Write matched text styled
		matchStart := pos + idx
		matchEnd := matchStart + len(query)
		matchedText := text[matchStart:matchEnd]
		result.WriteString(style.Render(matchedText))

		pos = matchEnd
	}
	return result.String()
}

func (m *Model) renderDictionary(layout viewportLayout) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Dictionary"))
	b.WriteString("\n\n")

	// Search input
	searchBarWidth := layout.Width - 10
	if searchBarWidth < 20 {
		searchBarWidth = 20
	}
	searchBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(searchBarWidth)

	searchText := m.dictionarySearch
	if searchText == "" {
		searchText = mutedStyle.Render("Search German or English...")
	} else {
		searchText += editStyle.Render("█") // Cursor
	}

	b.WriteString(searchBar.Render("🔍 " + searchText))
	b.WriteString("\n\n")

	// Results
	if len(m.dictionaryResults) == 0 {
		if m.dictionarySearch != "" {
			b.WriteString(mutedStyle.Render("No results found. Hint: Seed standard content in Import view [5] by pressing 'S' to populate the local dictionary."))
		} else {
			b.WriteString(mutedStyle.Render("Search dictionary (local dict.cc). Use 'S' in Import view to seed standard content."))
		}
		return b.String()
	}

	maxResults := layout.Height - 10
	if maxResults < 1 {
		maxResults = 1
	}

	// Adjust scroll for list
	if m.dictionaryCursor < m.dictionaryScroll {
		m.dictionaryScroll = m.dictionaryCursor
	}
	if m.dictionaryCursor >= m.dictionaryScroll+maxResults {
		m.dictionaryScroll = m.dictionaryCursor - maxResults + 1
	}

	// Two-column layout if wide enough
	if layout.Width > 80 {
		listWidth := layout.Width / 2
		detailWidth := layout.Width - listWidth - 6

		var listBuilder strings.Builder
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}
			highlightedWord := highlightQuery(res.Word, m.dictionarySearch, dictHighlightStyle)
			line := fmt.Sprintf("%s%s", prefix, highlightedWord)
			if i == m.dictionaryCursor {
				listBuilder.WriteString(editStyle.Render(lipgloss.NewStyle().Width(listWidth-2).MaxHeight(1).Render(line)) + "\n")
			} else {
				listBuilder.WriteString(lipgloss.NewStyle().Width(listWidth-2).MaxHeight(1).Render(line) + "\n")
			}
		}

		// Fill remaining space in list
		for i := len(m.dictionaryResults) - m.dictionaryScroll; i < maxResults; i++ {
			listBuilder.WriteString("\n")
		}

		// List scrollbar
		listWithScroll := listBuilder.String()
		if len(m.dictionaryResults) > maxResults {
			var sb strings.Builder
			thumbStart, thumbHeight := scrollbarThumb(len(m.dictionaryResults), maxResults, m.dictionaryScroll)
			lines := strings.Split(listBuilder.String(), "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				char := "│"
				if i >= thumbStart && i < thumbStart+thumbHeight {
					char = "┃"
				}
				sb.WriteString(lines[i] + lipgloss.NewStyle().Foreground(colorPanel).Render(char) + "\n")
			}
			listWithScroll = sb.String()
		}

		var detailBuilder strings.Builder
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			res := m.dictionaryResults[m.dictionaryCursor]
			detailBuilder.WriteString(titleStyle.Render(res.Word) + "\n")
			if res.Translation != "" {
				translations := strings.Split(res.Translation, ";")
				for _, t := range translations {
					detailBuilder.WriteString(strings.TrimSpace(t) + "\n")
				}
			}
			if res.WordClass != "" || res.Gender != "" {
				meta := ""
				if res.WordClass != "" {
					meta += mutedStyle.Render("["+res.WordClass+"]") + " "
				}
				if res.Gender != "" {
					meta += renderGender(res.Gender) + " "
				}
				detailBuilder.WriteString(meta + "\n")
			}
			if res.Forms != "" {
				detailBuilder.WriteString("\n" + boldStyle.Render("Forms:") + "\n")
				detailBuilder.WriteString(res.Forms + "\n")
			}
			if len(res.Examples) > 0 {
				detailBuilder.WriteString("\n" + boldStyle.Render("Examples:") + "\n")
				for _, ex := range res.Examples {
					detailBuilder.WriteString("• " + ex + "\n")
				}
			}
		}

		detailLines := strings.Split(detailBuilder.String(), "\n")
		m.dictionaryDetailTotalLines = len(detailLines)

		// Adjust detail scroll
		if m.dictionaryDetailScroll > m.dictionaryDetailTotalLines-maxResults {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailTotalLines-maxResults)
		}

		var visibleDetail strings.Builder
		for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResults && i < len(detailLines); i++ {
			visibleDetail.WriteString(detailLines[i] + "\n")
		}

		// Detail scrollbar
		detailPanelContent := visibleDetail.String()
		if m.dictionaryDetailTotalLines > maxResults {
			var sb strings.Builder
			thumbStart, thumbHeight := scrollbarThumb(m.dictionaryDetailTotalLines, maxResults, m.dictionaryDetailScroll)
			lines := strings.Split(detailPanelContent, "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				char := "│"
				if i >= thumbStart && i < thumbStart+thumbHeight {
					char = "┃"
				}
				sb.WriteString(lipgloss.NewStyle().Foreground(colorPanel).Render(char) + lines[i] + "\n")
			}
			detailPanelContent = sb.String()
		}

		detailPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			Padding(0, 2).
			Width(detailWidth).
			Height(maxResults).
			Render(detailPanelContent)

		joined := lipgloss.JoinHorizontal(lipgloss.Top, listWithScroll, detailPanel)
		b.WriteString(joined)
	} else {
		// Single column layout
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}
			highlightedWord := highlightQuery(res.Word, m.dictionarySearch, dictHighlightStyle)
			highlightedTranslation := highlightQuery(res.Translation, m.dictionarySearch, dictHighlightStyle)
			line := fmt.Sprintf("%s%s - %s", prefix, highlightedWord, highlightedTranslation)
			if res.WordClass != "" {
				line += " " + mutedStyle.Render("["+res.WordClass+"]")
			}
			if res.Gender != "" {
				line += " " + renderGender(res.Gender)
			}
			if i == m.dictionaryCursor {
				b.WriteString(editStyle.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		if len(m.dictionaryResults) > m.dictionaryScroll+maxResults {
			b.WriteString(mutedStyle.Render(fmt.Sprintf("+ %d more...", len(m.dictionaryResults)-(m.dictionaryScroll+maxResults))))
		}
	}

	return b.String()
}
