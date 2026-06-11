package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	tea "charm.land/bubbletea/v2"
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
	textRunes := []rune(text)
	lowerTextRunes := []rune(strings.ToLower(text))
	lowerQueryRunes := []rune(strings.ToLower(query))
	if len(lowerQueryRunes) == 0 || len(lowerQueryRunes) > len(lowerTextRunes) {
		return text
	}

	var result strings.Builder
	pos := 0
	for pos < len(lowerTextRunes) {
		matchStart := -1
		for i := pos; i <= len(lowerTextRunes)-len(lowerQueryRunes); i++ {
			matched := true
			for j := range lowerQueryRunes {
				if lowerTextRunes[i+j] != lowerQueryRunes[j] {
					matched = false
					break
				}
			}
			if matched {
				matchStart = i
				break
			}
		}
		if matchStart == -1 {
			result.WriteString(string(textRunes[pos:]))
			break
		}

		result.WriteString(string(textRunes[pos:matchStart]))
		matchEnd := matchStart + len(lowerQueryRunes)
		result.WriteString(style.Render(string(textRunes[matchStart:matchEnd])))
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

func dictionaryVisibleRows(layout viewportLayout) int {
	rows := layout.Height - 10
	if rows < 1 {
		return 1
	}
	return rows
}

func (m *Model) renderDictionary(layout viewportLayout) string {
	var b strings.Builder
	if len(m.dictionaryResults) > 0 {
		countSuffix := fmt.Sprintf(" (%d results)", len(m.dictionaryResults))
		if len(m.dictionaryResults) >= 50 {
			countSuffix = " (50+ results)"
		}
		b.WriteString(titleStyle.Render("Dictionary" + countSuffix))
	} else {
		b.WriteString(titleStyle.Render("Dictionary"))
	}
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
		b.WriteString(searchBar.Render("🔍 " + searchText))
	} else {
		// Render with an interactive [x] clear button on the right
		contentLength := 3 + len([]rune(m.dictionarySearch)) + 1 // "🔍 " + query + cursor
		spaces := searchBarWidth - contentLength - 3
		clearText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true).Render("[x]")
		var clearBtn string
		if spaces > 0 {
			clearBtn = strings.Repeat(" ", spaces) + clearText
		} else {
			clearBtn = " " + clearText
		}
		b.WriteString(searchBar.Render("🔍 " + searchText + editStyle.Render("█") + clearBtn))

		// Register hitbox for clear button [x]
		// Y coordinate is layout.Y + 3 (0: title, 1: empty line, 2: border top, 3: content)
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "dict-search-clear",
			View:   ViewDictionary,
			X:      layout.X + 2 + searchBarWidth - 3,
			Y:      layout.Y + 3,
			Width:  3,
			Height: 1,
			Action: func() tea.Cmd {
				m.dictionarySearch = ""
				m.dictionaryResults = nil
				m.dictionaryCursor = 0
				m.dictionaryScroll = 0
				return nil
			},
		})
	}
	b.WriteString("\n\n")

	// Results
	if len(m.dictionaryResults) == 0 {
		if m.dictionarySearch != "" {
			b.WriteString(mutedStyle.Render("No results found. Hint: Seed standard content in Import view [5] by pressing 'S' to populate the local dictionary."))
		} else {
			b.WriteString(mutedStyle.Render("Search dictionary (local dict.cc). Use 'S' in Import view to seed standard content."))
			if len(m.dictionarySearchHistory) > 0 {
				clearHistoryText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("[Clear]")
				b.WriteString("\n\n" + boldStyle.Render("Recent Searches:") + "  " + clearHistoryText + " " + mutedStyle.Render("(ctrl+x)") + "\n")

				// Hitbox for "[Clear]" button
				clearHistoryY := strings.Count(b.String(), "\n") - 1
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "dict-history-clear",
					View:   ViewDictionary,
					X:      layout.X + 18,
					Y:      layout.Y + clearHistoryY,
					Width:  7,
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionarySearchHistory = nil
						return nil
					},
				})

				for i, q := range m.dictionarySearchHistory {
					// Retrieve the current line count to calculate Y coordinate dynamically
					lineY := strings.Count(b.String(), "\n")
					b.WriteString(fmt.Sprintf("  • %s\n", q))
					// Save local variable q for the closure
					queryText := q
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-history-%d", i),
						View:   ViewDictionary,
						X:      layout.X + 4,
						Y:      layout.Y + lineY,
						Width:  len([]rune(queryText)),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = queryText
							return m.searchDictionary()
						},
					})
				}
			}
		}
		return b.String()
	}

	maxResults := dictionaryVisibleRows(layout)

	// Adjust scroll for list
	if m.dictionaryCursor < m.dictionaryScroll {
		m.dictionaryScroll = m.dictionaryCursor
	}
	if m.dictionaryCursor >= m.dictionaryScroll+maxResults {
		m.dictionaryScroll = m.dictionaryCursor - maxResults + 1
	}

	// Single column detail view
	if m.dictionaryDetailView && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
		res := m.dictionaryResults[m.dictionaryCursor]
		detailWidth := layout.Width

		var detailBuilder strings.Builder
		detailBuilder.WriteString(titleStyle.Render(res.Word) + "\n")

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

		detailBuilder.WriteString(mutedStyle.Render(strings.Repeat("─", maxInt(10, detailWidth-6))) + "\n\n")

		if res.Translation != "" {
			detailBuilder.WriteString(boldStyle.Render("Translations:") + "\n")
			translations := strings.Split(res.Translation, ";")
			for _, t := range translations {
				trimmed := strings.TrimSpace(t)
				highlightedT := highlightQuery(trimmed, m.dictionarySearch, dictHighlightStyle)
				detailBuilder.WriteString("  " + highlightedT + "\n")
			}
			detailBuilder.WriteString("\n")
		}

		if res.Forms != "" {
			detailBuilder.WriteString(boldStyle.Render("Word Forms:") + "\n")
			highlightedForms := highlightQuery(res.Forms, m.dictionarySearch, dictHighlightStyle)
			detailBuilder.WriteString("  " + highlightedForms + "\n\n")
		}

		if len(res.Examples) > 0 {
			detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
			for _, ex := range res.Examples {
				highlightedEx := highlightQuery(ex, m.dictionarySearch, dictHighlightStyle)
				detailBuilder.WriteString("  • " + highlightedEx + "\n")
			}
		}

		detailLines := strings.Split(detailBuilder.String(), "\n")
		if len(detailLines) > 0 && detailLines[len(detailLines)-1] == "" {
			detailLines = detailLines[:len(detailLines)-1]
		}
		m.dictionaryDetailTotalLines = len(detailLines)

		maxResultsDetail := dictionaryVisibleRows(layout)
		if m.dictionaryDetailScroll > m.dictionaryDetailTotalLines-maxResultsDetail {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailTotalLines-maxResultsDetail)
		}

		detailContentWidth := detailWidth - 5
		if detailContentWidth < 10 {
			detailContentWidth = 10
		}

		var visibleDetailBuilder strings.Builder
		for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResultsDetail && i < len(detailLines); i++ {
			visibleDetailBuilder.WriteString(padString(detailLines[i], detailContentWidth) + "\n")
		}
		for i := len(detailLines) - m.dictionaryDetailScroll; i < maxResultsDetail; i++ {
			visibleDetailBuilder.WriteString(strings.Repeat(" ", detailContentWidth) + "\n")
		}

		detailPanelContent := visibleDetailBuilder.String()
		if m.dictionaryDetailTotalLines > maxResultsDetail {
			var sb strings.Builder
			thumbStart, thumbHeight := scrollbarThumb(m.dictionaryDetailTotalLines, maxResultsDetail, m.dictionaryDetailScroll)
			lines := strings.Split(detailPanelContent, "\n")
			for i := 0; i < maxResultsDetail && i < len(lines); i++ {
				char := "│"
				if i >= thumbStart && i < thumbStart+thumbHeight {
					char = "┃"
				}
				sb.WriteString(lines[i] + lipgloss.NewStyle().Foreground(colorPanel).Render(char) + "\n")
			}
			detailPanelContent = sb.String()
		} else {
			var sb strings.Builder
			lines := strings.Split(detailPanelContent, "\n")
			for i := 0; i < maxResultsDetail && i < len(lines); i++ {
				sb.WriteString(lines[i] + " \n")
			}
			detailPanelContent = sb.String()
		}

		b.WriteString(detailPanelContent)
		b.WriteString("\n" + mutedStyle.Render("Press esc/ctrl+d to return to list | Enter to draft | ctrl+a to add | ctrl+p to play"))
		return b.String()
	}

	// Two-column layout if wide enough
	if layout.Width > 80 {
		listWidth := maxInt(30, minInt(50, layout.Width*4/10))
		detailWidth := layout.Width - listWidth - 3

		var listBuilder strings.Builder
		contentWidth := listWidth - 2
		listStartLine := strings.Count(b.String(), "\n")
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

			// Click hitbox to select item
			idx := i
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("dict-result-%d", idx),
				View:   ViewDictionary,
				X:      layout.X,
				Y:      layout.Y + listStartLine + (idx - m.dictionaryScroll),
				Width:  listWidth,
				Height: 1,
				Action: func() tea.Cmd {
					m.dictionaryCursor = idx
					m.dictionaryDetailScroll = 0
					return nil
				},
			})
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
					trimmed := strings.TrimSpace(t)
					highlightedT := highlightQuery(trimmed, m.dictionarySearch, dictHighlightStyle)
					detailBuilder.WriteString("  " + highlightedT + "\n")
				}
				detailBuilder.WriteString("\n")
			}

			// Word Forms
			if res.Forms != "" {
				detailBuilder.WriteString(boldStyle.Render("Word Forms:") + "\n")
				highlightedForms := highlightQuery(res.Forms, m.dictionarySearch, dictHighlightStyle)
				detailBuilder.WriteString("  " + highlightedForms + "\n\n")
			}

			// Examples
			if len(res.Examples) > 0 {
				detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
				for _, ex := range res.Examples {
					highlightedEx := highlightQuery(ex, m.dictionarySearch, dictHighlightStyle)
					detailBuilder.WriteString("  • " + highlightedEx + "\n")
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
		listStartLine := strings.Count(b.String(), "\n")
		contentWidth := layout.Width - 4
		if contentWidth < 10 {
			contentWidth = 10
		}

		var listBuilder strings.Builder
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}

			plainLine := fmt.Sprintf("%s - %s", res.Word, res.Translation)
			if res.WordClass != "" {
				plainLine += " [" + res.WordClass + "]"
			}
			if res.Gender != "" {
				plainLine += " {" + res.Gender + "}"
			}

			padded := padString(plainLine, contentWidth)
			highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
			line := prefix + highlighted

			if i == m.dictionaryCursor {
				listBuilder.WriteString(editStyle.Render(line) + "\n")
			} else {
				listBuilder.WriteString(line + "\n")
			}

			// Hitbox for click-to-select and click-again-to-details
			idx := i
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("dict-result-%d", idx),
				View:   ViewDictionary,
				X:      layout.X,
				Y:      layout.Y + listStartLine + (idx - m.dictionaryScroll),
				Width:  layout.Width,
				Height: 1,
				Action: func() tea.Cmd {
					if m.dictionaryCursor == idx {
						m.dictionaryDetailView = true
						m.dictionaryDetailScroll = 0
					} else {
						m.dictionaryCursor = idx
						m.dictionaryDetailScroll = 0
					}
					return nil
				},
			})
		}

		// Fill remaining space
		for i := len(m.dictionaryResults) - m.dictionaryScroll; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", layout.Width) + "\n")
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
			var sb strings.Builder
			lines := strings.Split(listWithScroll, "\n")
			for i := 0; i < maxResults && i < len(lines); i++ {
				sb.WriteString(lines[i] + " \n")
			}
			listWithScroll = sb.String()
		}
		b.WriteString(listWithScroll)

		b.WriteString("\n" + mutedStyle.Render("Press ctrl+d/click selected to view details | Enter to draft | ctrl+a to add | ctrl+p to play"))
	}

	return b.String()
}

func (m *Model) renderSpotlightDictionary() string {
	boxWidth := 86
	if m.width < 92 {
		boxWidth = m.width - 6
	}
	if boxWidth < 30 {
		boxWidth = 30
	}

	boxHeight := 16
	if m.height < 22 {
		boxHeight = m.height - 6
	}
	if boxHeight < 8 {
		boxHeight = 8
	}

	startX := (m.width - boxWidth) / 2
	startY := (m.height - boxHeight) / 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")). // Vibrant pink/magenta spotlight border
		Background(lipgloss.Color("233")).       // Deep dark backdrop
		Padding(0, 1).
		Width(boxWidth).
		Height(boxHeight)

	var b strings.Builder

	title := " 🔍 SPOTLIGHT DICTIONARY "
	if m.dictionarySearch != "" && len(m.dictionaryResults) > 0 {
		if len(m.dictionaryResults) >= 50 {
			title = " 🔍 SPOTLIGHT DICTIONARY (50+ results) "
		} else if len(m.dictionaryResults) == 1 {
			title = " 🔍 SPOTLIGHT DICTIONARY (1 result) "
		} else {
			title = fmt.Sprintf(" 🔍 SPOTLIGHT DICTIONARY (%d results) ", len(m.dictionaryResults))
		}
	}
	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(title)
	closeHint := mutedStyle.Render("Press = or Esc to close")
	b.WriteString(titleStr + "  " + closeHint + "\n\n")

	searchBarWidth := boxWidth - 6
	if searchBarWidth < 10 {
		searchBarWidth = 10
	}

	searchBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(searchBarWidth)

	searchText := m.dictionarySearch
	if searchText == "" {
		searchText = mutedStyle.Render("Search German or English...")
		b.WriteString(searchBar.Render("🔍 "+searchText) + "\n\n")
	} else {
		contentLength := 3 + len([]rune(m.dictionarySearch)) + 1
		spaces := searchBarWidth - contentLength - 3
		clearText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true).Render("[x]")
		var clearBtn string
		if spaces > 0 {
			clearBtn = strings.Repeat(" ", spaces) + clearText
		} else {
			clearBtn = " " + clearText
		}
		b.WriteString(searchBar.Render("🔍 "+searchText+editStyle.Render("█")+clearBtn) + "\n\n")

		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     "dict-overlay-search-clear",
			View:   m.activeView,
			X:      startX + searchBarWidth - 2,
			Y:      startY + 3,
			Width:  3,
			Height: 1,
			Action: func() tea.Cmd {
				m.resetDictionarySearchState()
				return nil
			},
		})
	}

	usedLines := 7
	interiorHeight := boxHeight - usedLines - 2
	if interiorHeight < 1 {
		interiorHeight = 1
	}

	interiorWidth := boxWidth - 4
	if interiorWidth < 10 {
		interiorWidth = 10
	}

	if len(m.dictionaryResults) == 0 {
		if m.dictionarySearch != "" {
			b.WriteString(mutedStyle.Render("No results found."))
		} else {
			b.WriteString(mutedStyle.Render("Type to search local dict.cc dictionary."))
			if len(m.dictionarySearchHistory) > 0 {
				clearHistoryText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("[Clear]")
				b.WriteString("\n\n" + boldStyle.Render("Recent Searches:") + "  " + clearHistoryText + " " + mutedStyle.Render("(ctrl+x)") + "\n")

				clearHistoryLineY := strings.Count(b.String(), "\n") - 1
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "dict-overlay-history-clear",
					View:   m.activeView,
					X:      startX + 20,
					Y:      startY + clearHistoryLineY + 2,
					Width:  7,
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionarySearchHistory = nil
						return nil
					},
				})

				for i, q := range m.dictionarySearchHistory {
					lineY := strings.Count(b.String(), "\n")
					b.WriteString(fmt.Sprintf("  • %s\n", q))
					queryText := q
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-overlay-history-%d", i),
						View:   m.activeView,
						X:      startX + 4,
						Y:      startY + lineY + 2,
						Width:  len([]rune(queryText)),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = queryText
							m.dictionaryResults = nil
							m.dictionaryCursor = 0
							m.dictionaryScroll = 0
							m.dictionaryDetailScroll = 0
							m.dictionaryDetailTotalLines = 0
							m.dictionaryDetailView = false
							return m.searchDictionary()
						},
					})
				}
			}
		}
		return boxStyle.Render(b.String())
	}

	maxResults := interiorHeight
	if m.dictionaryCursor < m.dictionaryScroll {
		m.dictionaryScroll = m.dictionaryCursor
	}
	if m.dictionaryCursor >= m.dictionaryScroll+maxResults {
		m.dictionaryScroll = m.dictionaryCursor - maxResults + 1
	}

	if m.dictionaryDetailView && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
		res := m.dictionaryResults[m.dictionaryCursor]

		var detailBuilder strings.Builder
		detailBuilder.WriteString(titleStyle.Render(res.Word) + "\n")
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
		detailBuilder.WriteString(mutedStyle.Render(strings.Repeat("─", maxInt(10, interiorWidth-6))) + "\n")

		if res.Translation != "" {
			detailBuilder.WriteString(boldStyle.Render("Translations:") + "\n")
			translations := strings.Split(res.Translation, ";")
			for _, t := range translations {
				trimmed := strings.TrimSpace(t)
				highlightedT := highlightQuery(trimmed, m.dictionarySearch, dictHighlightStyle)
				detailBuilder.WriteString("  " + highlightedT + "\n")
			}
		}
		if res.Forms != "" {
			detailBuilder.WriteString(boldStyle.Render("Forms: ") + highlightQuery(res.Forms, m.dictionarySearch, dictHighlightStyle) + "\n")
		}
		if len(res.Examples) > 0 {
			detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
			for _, ex := range res.Examples {
				detailBuilder.WriteString("  • " + highlightQuery(ex, m.dictionarySearch, dictHighlightStyle) + "\n")
			}
		}

		detailLines := strings.Split(detailBuilder.String(), "\n")
		m.dictionaryDetailTotalLines = len(detailLines)
		if m.dictionaryDetailScroll > m.dictionaryDetailTotalLines-maxResults {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailTotalLines-maxResults)
		}

		var visibleDetailBuilder strings.Builder
		for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResults && i < len(detailLines); i++ {
			visibleDetailBuilder.WriteString(padString(detailLines[i], interiorWidth) + "\n")
		}

		b.WriteString(visibleDetailBuilder.String())
		return boxStyle.Render(b.String())
	}

	if interiorWidth > 70 {
		listWidth := maxInt(25, minInt(40, interiorWidth*4/10))
		detailWidth := interiorWidth - listWidth - 3

		var listBuilder strings.Builder
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}
			wordText := res.Word
			if res.Gender != "" {
				wordText += " {" + res.Gender + "}"
			} else if res.WordClass != "" {
				wordText += " [" + res.WordClass + "]"
			}
			padded := padString(wordText, listWidth-2)
			highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
			line := prefix + highlighted
			if i == m.dictionaryCursor {
				listBuilder.WriteString(editStyle.Render(line) + "\n")
			} else {
				listBuilder.WriteString(line + "\n")
			}

			idx := i
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("dict-overlay-result-%d", idx),
				View:   m.activeView,
				X:      startX + 2,
				Y:      startY + usedLines + (idx - m.dictionaryScroll),
				Width:  listWidth,
				Height: 1,
				Action: func() tea.Cmd {
					m.dictionaryCursor = idx
					m.dictionaryDetailScroll = 0
					return nil
				},
			})
		}
		for i := len(m.dictionaryResults) - m.dictionaryScroll; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", listWidth) + "\n")
		}

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
			detailBuilder.WriteString(titleStyle.Render(res.Word) + "\n")
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
			detailBuilder.WriteString(mutedStyle.Render(strings.Repeat("─", detailWidth-4)) + "\n")

			if res.Translation != "" {
				translations := strings.Split(res.Translation, ";")
				for _, t := range translations {
					detailBuilder.WriteString("  " + highlightQuery(strings.TrimSpace(t), m.dictionarySearch, dictHighlightStyle) + "\n")
				}
			}
		}

		detailLines := strings.Split(detailBuilder.String(), "\n")
		m.dictionaryDetailTotalLines = len(detailLines)
		if m.dictionaryDetailScroll > m.dictionaryDetailTotalLines-maxResults {
			m.dictionaryDetailScroll = maxInt(0, m.dictionaryDetailTotalLines-maxResults)
		}

		var visibleDetailBuilder strings.Builder
		for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResults && i < len(detailLines); i++ {
			visibleDetailBuilder.WriteString(padString(detailLines[i], detailWidth-2) + "\n")
		}
		for i := len(detailLines) - m.dictionaryDetailScroll; i < maxResults; i++ {
			visibleDetailBuilder.WriteString(strings.Repeat(" ", detailWidth-2) + "\n")
		}

		detailPanel := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(colorPanel).
			Padding(0, 1).
			Width(detailWidth).
			Height(maxResults).
			Render(visibleDetailBuilder.String())

		joined := lipgloss.JoinHorizontal(lipgloss.Top, listWithScroll, detailPanel)
		b.WriteString(joined)
	} else {
		var listBuilder strings.Builder
		for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
			res := m.dictionaryResults[i]
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}
			plainLine := fmt.Sprintf("%s - %s", res.Word, res.Translation)
			padded := padString(plainLine, interiorWidth-2)
			highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
			line := prefix + highlighted
			if i == m.dictionaryCursor {
				listBuilder.WriteString(editStyle.Render(line) + "\n")
			} else {
				listBuilder.WriteString(line + "\n")
			}

			idx := i
			m.hitboxes = append(m.hitboxes, Hitbox{
				ID:     fmt.Sprintf("dict-overlay-result-%d", idx),
				View:   m.activeView,
				X:      startX + 2,
				Y:      startY + usedLines + (idx - m.dictionaryScroll),
				Width:  interiorWidth,
				Height: 1,
				Action: func() tea.Cmd {
					if m.dictionaryCursor == idx {
						m.dictionaryDetailView = true
						m.dictionaryDetailScroll = 0
					} else {
						m.dictionaryCursor = idx
						m.dictionaryDetailScroll = 0
					}
					return nil
				},
			})
		}
		for i := len(m.dictionaryResults) - m.dictionaryScroll; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", interiorWidth) + "\n")
		}
		b.WriteString(listBuilder.String())
	}

	return boxStyle.Render(b.String())
}
