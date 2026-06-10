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

func padString(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
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
		listWidth := maxInt(30, minInt(50, layout.Width*4/10))
		detailWidth := layout.Width - listWidth - 3

		var listBuilder strings.Builder
		contentWidth := listWidth - 2
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}

			// Format the word text for list item (e.g. "Word {m}" or "Word [verb]")
			wordText := res.Word
			if res.Gender != "" {
				wordText += " {" + res.Gender + "}"
			} else if res.WordClass != "" {
				wordText += " [" + res.WordClass + "]"
			}

			padded := padString(wordText, contentWidth)
			highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
			line := prefix + highlighted

			if i == m.dictionaryCursor {
				listBuilder.WriteString(editStyle.Render(line) + "\n")
			} else {
				listBuilder.WriteString(line + "\n")
			}
		}

		// Fill remaining space in list
		for i := len(m.dictionaryResults) - m.dictionaryScroll; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", listWidth) + "\n")
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
		} else {
			// Append an empty vertical line for alignment when no scrollbar is shown
			var sb strings.Builder
			lines := strings.Split(listBuilder.String(), "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				sb.WriteString(lines[i] + " \n")
			}
			listWithScroll = sb.String()
		}

		var detailBuilder strings.Builder
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			res := m.dictionaryResults[m.dictionaryCursor]
			
			// Modern styled header for detail view
			detailBuilder.WriteString(titleStyle.Render(res.Word) + "\n")
			
			// Part of speech / Gender tags
			meta := ""
			if res.WordClass != "" {
				meta += mutedStyle.Render("["+strings.ToUpper(res.WordClass)+"]") + " "
			}
			if res.Gender != "" {
				meta += renderGender(res.Gender) + " "
			}
			if meta != "" {
				detailBuilder.WriteString(meta + "\n")
			}
			
			detailBuilder.WriteString(mutedStyle.Render(strings.Repeat("─", detailWidth-6)) + "\n\n")

			// Translation Section
			if res.Translation != "" {
				detailBuilder.WriteString(boldStyle.Render("Translations:") + "\n")
				translations := strings.Split(res.Translation, ";")
				for _, t := range translations {
					detailBuilder.WriteString("  " + strings.TrimSpace(t) + "\n")
				}
				detailBuilder.WriteString("\n")
			}
			
			// Word Forms
			if res.Forms != "" {
				detailBuilder.WriteString(boldStyle.Render("Word Forms:") + "\n")
				detailBuilder.WriteString("  " + res.Forms + "\n\n")
			}
			
			// Examples
			if len(res.Examples) > 0 {
				detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
				for _, ex := range res.Examples {
					detailBuilder.WriteString("  • " + ex + "\n")
				}
			}
		}

		detailLines := strings.Split(detailBuilder.String(), "\n")
		// Clean up trailing empty line from splitting if present
		if len(detailLines) > 0 && detailLines[len(detailLines)-1] == "" {
			detailLines = detailLines[:len(detailLines)-1]
		}
		m.dictionaryDetailTotalLines = len(detailLines)

		// Adjust detail scroll
		if m.dictionaryDetailScroll > m.dictionaryDetailTotalLines-maxResults {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailTotalLines-maxResults)
		}

		// Detail content width (accounting for border and padding)
		detailContentWidth := detailWidth - 5
		if detailContentWidth < 10 {
			detailContentWidth = 10
		}

		var visibleDetailBuilder strings.Builder
		for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResults && i < len(detailLines); i++ {
			// Pad detail lines so scrollbar aligns perfectly on the right
			visibleDetailBuilder.WriteString(padString(detailLines[i], detailContentWidth) + "\n")
		}
		// Fill remaining vertical space in detail panel to match list height
		for i := len(detailLines) - m.dictionaryDetailScroll; i < maxResults; i++ {
			visibleDetailBuilder.WriteString(strings.Repeat(" ", detailContentWidth) + "\n")
		}

		detailPanelContent := visibleDetailBuilder.String()
		if m.dictionaryDetailTotalLines > maxResults {
			var sb strings.Builder
			thumbStart, thumbHeight := scrollbarThumb(m.dictionaryDetailTotalLines, maxResults, m.dictionaryDetailScroll)
			lines := strings.Split(detailPanelContent, "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				char := "│"
				if i >= thumbStart && i < thumbStart+thumbHeight {
					char = "┃"
				}
				sb.WriteString(lines[i] + lipgloss.NewStyle().Foreground(colorPanel).Render(char) + "\n")
			}
			detailPanelContent = sb.String()
		} else {
			// Append an empty vertical line for alignment when no scrollbar is shown
			var sb strings.Builder
			lines := strings.Split(detailPanelContent, "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				sb.WriteString(lines[i] + " \n")
			}
			detailPanelContent = sb.String()
		}

		detailPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(colorPanel).
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
