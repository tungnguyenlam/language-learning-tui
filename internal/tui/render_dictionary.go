package tui

import (
	"fmt"
	"strings"

	"deutsch-tui/internal/content"

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

func padString(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return truncateLine(s, width)
	}
	return s + strings.Repeat(" ", width-w)
}

func dictionaryVisibleRows(layout viewportLayout) int {
	rows := layout.Height - 12
	if rows < 1 {
		return 1
	}
	return rows
}

func renderFilterPillsRow(m *Model, startX, startY int, targetView View) string {
	pills := []struct {
		Name string
		Tag  string
	}{
		{"All", ""},
		{"★ Starred", ":starred"},
		{"DE", "de:"},
		{"EN", "en:"},
		{"Verb", ":verb"},
		{"Noun", ":noun"},
		{"Adj", ":adj"},
		{"Adv", ":adv"},
		{"Der", ":m"},
		{"Die", ":f"},
		{"Das", ":n"},
		{"Pl", ":pl"},
	}

	activeStyle := lipgloss.NewStyle().Bold(true).Foreground(colorYellow).Background(lipgloss.Color("236")).Padding(0, 1)
	inactiveStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Padding(0, 1)

	var parts []string
	currentX := startX + 9
	for i, p := range pills {
		active := isFilterActive(m.dictionarySearch, p.Tag)
		label := p.Name
		var rendered string
		if active {
			rendered = activeStyle.Render(label)
		} else {
			rendered = inactiveStyle.Render(label)
		}
		pillWidth := lipgloss.Width(rendered)

		tag := p.Tag
		m.hitboxes = append(m.hitboxes, Hitbox{
			ID:     fmt.Sprintf("dict-filter-pill-%d", i),
			View:   targetView,
			X:      currentX,
			Y:      startY,
			Width:  pillWidth,
			Height: 1,
			Action: func() tea.Cmd {
				if tag == "" {
					m.dictionarySearch = clearFilterTags(m.dictionarySearch)
				} else {
					m.dictionarySearch = toggleFilterTag(m.dictionarySearch, tag)
				}
				m.dictionaryCursor = 0
				m.dictionaryScroll = 0
				return m.searchDictionary()
			},
		})
		currentX += pillWidth + 1
		parts = append(parts, rendered)
	}

	targetDeckName := "Dictionary"
	if m.dictionaryTargetDeckID != "" {
		for _, d := range m.decks {
			if d.ID == m.dictionaryTargetDeckID {
				targetDeckName = d.Name
				break
			}
		}
		if m.dictionaryTargetDeckID == "dictionary" {
			targetDeckName = "Dictionary"
		}
	} else if m.deck.ID != "" && m.deck.ID != "all" {
		targetDeckName = m.deck.Name
	}
	deckPillStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Background(lipgloss.Color("235")).Padding(0, 1)
	renderedDeck := deckPillStyle.Render("Target: [" + targetDeckName + "]")
	deckWidth := lipgloss.Width(renderedDeck)
	m.hitboxes = append(m.hitboxes, Hitbox{
		ID:     "dict-target-deck-pill",
		View:   targetView,
		X:      currentX + 1,
		Y:      startY,
		Width:  deckWidth,
		Height: 1,
		Action: func() tea.Cmd {
			return m.cycleDictionaryTargetDeck()
		},
	})
	parts = append(parts, renderedDeck)

	return mutedStyle.Render("Filters: ") + strings.Join(parts, " ")
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
		// Clip query if it's too long for the search bar width
		availableWidth := searchBarWidth - 8
		displaySearch := searchText
		if len([]rune(searchText)) > availableWidth && availableWidth > 5 {
			displaySearch = "..." + string([]rune(searchText)[len([]rune(searchText))-availableWidth+3:])
		}
		// Render with an interactive [x] clear button on the right
		contentLength := 3 + len([]rune(displaySearch)) + 1 // "🔍 " + query + cursor
		spaces := searchBarWidth - contentLength - 3
		clearText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true).Render("[x]")
		var clearBtn string
		if spaces > 0 {
			clearBtn = strings.Repeat(" ", spaces) + clearText
		} else {
			clearBtn = " " + clearText
		}
		b.WriteString(searchBar.Render("🔍 " + displaySearch + editStyle.Render("█") + clearBtn))

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

	// Render interactive filter pills row
	filterLineY := strings.Count(b.String(), "\n")
	b.WriteString(renderFilterPillsRow(m, layout.X, layout.Y+filterLineY, ViewDictionary))
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
						m.saveDictionaryHistory()
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

			if len(m.dictionaryRecentlyViewed) > 0 {
				b.WriteString("\n" + boldStyle.Render("Recently Inspected Words:") + "\n")
				for i, e := range m.dictionaryRecentlyViewed {
					lineY := strings.Count(b.String(), "\n")
					wordStr := e.Word
					if e.Gender != "" {
						wordStr += " {" + e.Gender + "}"
					}
					if e.Translation != "" {
						wordStr += " - " + e.Translation
					}
					b.WriteString(fmt.Sprintf("  • %s\n", wordStr))
					entry := e
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-recent-view-%d", i),
						View:   ViewDictionary,
						X:      layout.X + 4,
						Y:      layout.Y + lineY,
						Width:  len([]rune(wordStr)),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = entry.Word
							return m.searchDictionary()
						},
					})
				}
			}

			if len(m.dictionaryDiscoverEntries) > 0 {
				discoverIcon := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("✦")
				b.WriteString("\n" + discoverIcon + " " + boldStyle.Render("Discover:") + " " + mutedStyle.Render("random words to explore") + "\n")
				for i, e := range m.dictionaryDiscoverEntries {
					lineY := strings.Count(b.String(), "\n")
					wordStr := e.Word
					if e.Gender != "" {
						wordStr += " " + renderGender(e.Gender)
					}
					if e.Translation != "" {
						trans := e.Translation
						if idx := strings.Index(trans, ";"); idx > 0 {
							trans = strings.TrimSpace(trans[:idx])
						}
						wordStr += " — " + mutedStyle.Render(trans)
					}
					b.WriteString(fmt.Sprintf("  • %s\n", wordStr))
					entry := e
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-discover-%d", i),
						View:   ViewDictionary,
						X:      layout.X + 4,
						Y:      layout.Y + lineY,
						Width:  lipgloss.Width(wordStr),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = entry.Word
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
		titleText := res.Word
		if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
			titleText = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ") + titleText
		}
		detailBuilder.WriteString(titleStyle.Render(titleText) + "\n")

		meta := ""
		if res.WordClass != "" {
			meta += mutedStyle.Render("["+strings.ToUpper(res.WordClass)+"]") + " "
		}
		if res.Gender != "" {
			meta += renderGender(res.Gender) + " "
		}
		if len(res.Tags) > 0 {
			tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("236")).Padding(0, 1)
			for _, tg := range res.Tags {
				cleanTag := strings.Trim(tg, "[] ")
				if cleanTag != "" {
					meta += tagStyle.Render("["+strings.ToUpper(cleanTag)+"]") + " "
				}
			}
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
			detailBuilder.WriteString(formatInflectionTable(res.WordClass, res.Gender, res.Forms, m.dictionarySearch))
		}

		if compoundParts := content.DecomposeCompound(res.Word, nil); len(compoundParts) == 2 {
			lineY := strings.Count(detailBuilder.String(), "\n")
			detailBuilder.WriteString(boldStyle.Render("Compound Breakdown:") + "\n")
			part1 := compoundParts[0].Word
			part2 := compoundParts[1].Word
			detailBuilder.WriteString(fmt.Sprintf("  • %s + %s\n\n", part1, part2))

			maxRows := dictionaryVisibleRows(layout)
			if lineY+1 >= m.dictionaryDetailScroll && lineY+1 < m.dictionaryDetailScroll+maxRows {
				screenY := layout.Y + strings.Count(b.String(), "\n") + (lineY + 1 - m.dictionaryDetailScroll)
				p1Word := part1
				p2Word := part2
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "dict-compound-sc-0",
					View:   ViewDictionary,
					X:      layout.X + 4,
					Y:      screenY,
					Width:  lipgloss.Width(p1Word),
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionarySearch = p1Word
						return m.searchDictionary()
					},
				})
				p2X := layout.X + 4 + lipgloss.Width(p1Word) + 3
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "dict-compound-sc-1",
					View:   ViewDictionary,
					X:      p2X,
					Y:      screenY,
					Width:  lipgloss.Width(p2Word),
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionarySearch = p2Word
						return m.searchDictionary()
					},
				})
			}
		}

		if len(res.Examples) > 0 {
			detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
			for _, ex := range res.Examples {
				highlightedEx := highlightQuery(ex, m.dictionarySearch, dictHighlightStyle)
				detailBuilder.WriteString("  • " + highlightedEx + "\n")
			}
		}

		if len(m.dictionaryRelatedEntries) > 0 {
			if len(res.Examples) > 0 || res.Forms != "" || res.Translation != "" {
				detailBuilder.WriteString("\n")
			}
			detailBuilder.WriteString(boldStyle.Render("Related Words:") + "\n")
			for i, related := range m.dictionaryRelatedEntries {
				lineY := strings.Count(detailBuilder.String(), "\n")
				line := "  • " + related.Word
				if related.Gender != "" {
					line += " {" + related.Gender + "}"
				}
				if related.Translation != "" {
					line += " - " + related.Translation
				}
				detailBuilder.WriteString(line + "\n")

				// Register hitbox if visible
				maxRows := dictionaryVisibleRows(layout)
				if lineY >= m.dictionaryDetailScroll && lineY < m.dictionaryDetailScroll+maxRows {
					screenY := layout.Y + strings.Count(b.String(), "\n") + (lineY - m.dictionaryDetailScroll)
					relEntry := related
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-related-sc-%d", i),
						View:   ViewDictionary,
						X:      layout.X,
						Y:      screenY,
						Width:  lipgloss.Width(line),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = relEntry.Word
							return m.searchDictionary()
						},
					})
				}
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
		writtenDetailLines := maxInt(0, minInt(len(detailLines)-m.dictionaryDetailScroll, maxResultsDetail))
		for i := writtenDetailLines; i < maxResultsDetail; i++ {
			visibleDetailBuilder.WriteString(strings.Repeat(" ", detailContentWidth) + "\n")
		}

		renderedDetailLines := strings.Split(strings.TrimSuffix(visibleDetailBuilder.String(), "\n"), "\n")
		detailLinesWithScroll := renderScrollbarColumn(renderedDetailLines, maxResultsDetail, m.dictionaryDetailTotalLines, m.dictionaryDetailScroll)
		if m.dictionaryDetailTotalLines > maxResultsDetail {
			detailStartLine := strings.Count(b.String(), "\n")
			for i := 0; i < maxResultsDetail; i++ {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("dict-detail-scroll-%d", i),
					View:   ViewDictionary,
					X:      layout.X + layout.Width - 1,
					Y:      layout.Y + detailStartLine + i,
					Width:  1,
					Height: 1,
				})
			}
		}
		detailPanelContent := strings.Join(detailLinesWithScroll, "\n") + "\n"

		b.WriteString(detailPanelContent)
		b.WriteString("\n" + mutedStyle.Render("Press esc/ctrl+d to return to list | Enter to draft | ctrl+a to add | ctrl+e to explain | ctrl+p to play"))
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

			starStr := ""
			if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
				starStr = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ")
			}

			// Format the word text for list item (e.g. "Word {m} - house")
			wordText := res.Word
			if res.Gender != "" {
				wordText += " {" + res.Gender + "}"
			} else if res.WordClass != "" {
				wordText += " [" + res.WordClass + "]"
			}
			if res.Translation != "" {
				wordText = fmt.Sprintf("%s - %s", wordText, res.Translation)
			}

			padded := padString(wordText, maxInt(5, contentWidth-lipgloss.Width(starStr)))
			highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
			line := prefix + starStr + highlighted

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
					m.inspectDictionaryCursor()
					if m.dictionaryDetailVisible() && idx >= 0 && idx < len(m.dictionaryResults) {
						return m.findRelatedEntries(m.dictionaryResults[idx].Word)
					}
					return nil
				},
			})
		}

		// Fill remaining space in list
		writtenListLines := maxInt(0, minInt(len(m.dictionaryResults)-m.dictionaryScroll, maxResults))
		for i := writtenListLines; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", listWidth) + "\n")
		}

		// List scrollbar
		listLines := strings.Split(strings.TrimSuffix(listBuilder.String(), "\n"), "\n")
		listLinesWithScroll := renderScrollbarColumn(listLines, maxResults, len(m.dictionaryResults), m.dictionaryScroll)
		if len(m.dictionaryResults) > maxResults {
			for i := 0; i < maxResults; i++ {
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     fmt.Sprintf("dict-scroll-%d", i),
					View:   ViewDictionary,
					X:      layout.X + listWidth - 1,
					Y:      layout.Y + listStartLine + i,
					Width:  1,
					Height: 1,
				})
			}
		}
		listWithScroll := strings.Join(listLinesWithScroll, "\n") + "\n"

		var detailBuilder strings.Builder
		if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			res := m.dictionaryResults[m.dictionaryCursor]

			// Modern styled header for detail view
			titleText := res.Word
			if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
				titleText = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ") + titleText
			}
			detailBuilder.WriteString(titleStyle.Render(titleText) + "\n")

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
				detailBuilder.WriteString(formatInflectionTable(res.WordClass, res.Gender, res.Forms, m.dictionarySearch))
			}

			if compoundParts := content.DecomposeCompound(res.Word, nil); len(compoundParts) == 2 {
				lineY := strings.Count(detailBuilder.String(), "\n")
				detailBuilder.WriteString(boldStyle.Render("Compound Breakdown:") + "\n")
				part1 := compoundParts[0].Word
				part2 := compoundParts[1].Word
				detailBuilder.WriteString(fmt.Sprintf("  • %s + %s\n\n", part1, part2))

				if lineY+1 >= m.dictionaryDetailScroll && lineY+1 < m.dictionaryDetailScroll+maxResults {
					screenY := layout.Y + listStartLine + (lineY + 1 - m.dictionaryDetailScroll)
					p1Word := part1
					p2Word := part2
					detailLeftX := layout.X + listWidth + 3
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     "dict-compound-tc-0",
						View:   ViewDictionary,
						X:      detailLeftX + 4,
						Y:      screenY,
						Width:  lipgloss.Width(p1Word),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = p1Word
							return m.searchDictionary()
						},
					})
					p2X := detailLeftX + 4 + lipgloss.Width(p1Word) + 3
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     "dict-compound-tc-1",
						View:   ViewDictionary,
						X:      p2X,
						Y:      screenY,
						Width:  lipgloss.Width(p2Word),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = p2Word
							return m.searchDictionary()
						},
					})
				}
			}

			// Examples
			if len(res.Examples) > 0 {
				detailBuilder.WriteString(boldStyle.Render("Examples:") + "\n")
				for _, ex := range res.Examples {
					highlightedEx := highlightQuery(ex, m.dictionarySearch, dictHighlightStyle)
					detailBuilder.WriteString("  • " + highlightedEx + "\n")
				}
			}

			if len(m.dictionaryRelatedEntries) > 0 {
				if len(res.Examples) > 0 || res.Forms != "" || res.Translation != "" {
					detailBuilder.WriteString("\n")
				}
				detailBuilder.WriteString(boldStyle.Render("Related Words:") + "\n")
				for i, related := range m.dictionaryRelatedEntries {
					lineY := strings.Count(detailBuilder.String(), "\n")
					line := "  • " + related.Word
					if related.Gender != "" {
						line += " {" + related.Gender + "}"
					}
					if related.Translation != "" {
						line += " - " + related.Translation
					}
					detailBuilder.WriteString(line + "\n")

					// Register hitbox if visible
					if lineY >= m.dictionaryDetailScroll && lineY < m.dictionaryDetailScroll+maxResults {
						screenY := layout.Y + listStartLine + (lineY - m.dictionaryDetailScroll)
						relEntry := related
						m.hitboxes = append(m.hitboxes, Hitbox{
							ID:     fmt.Sprintf("dict-related-tc-%d", i),
							View:   ViewDictionary,
							X:      layout.X + listWidth + 3,
							Y:      screenY,
							Width:  lipgloss.Width(line),
							Height: 1,
							Action: func() tea.Cmd {
								m.dictionarySearch = relEntry.Word
								return m.searchDictionary()
							},
						})
					}
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
		writtenDetailLines := maxInt(0, minInt(len(detailLines)-m.dictionaryDetailScroll, maxResults))
		for i := writtenDetailLines; i < maxResults; i++ {
			visibleDetailBuilder.WriteString(strings.Repeat(" ", detailContentWidth) + "\n")
		}

		renderedDetailLines := strings.Split(strings.TrimSuffix(visibleDetailBuilder.String(), "\n"), "\n")
		detailLinesWithScroll := renderScrollbarColumn(renderedDetailLines, maxResults, m.dictionaryDetailTotalLines, m.dictionaryDetailScroll)
		detailPanelContent := strings.Join(detailLinesWithScroll, "\n") + "\n"

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

			starStr := ""
			if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
				starStr = "★ "
			}

			plainLine := fmt.Sprintf("%s%s - %s", starStr, res.Word, res.Translation)
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
						m.inspectDictionaryCursor()
						if idx >= 0 && idx < len(m.dictionaryResults) {
							return m.findRelatedEntries(m.dictionaryResults[idx].Word)
						}
						return nil
					}
					m.dictionaryCursor = idx
					m.dictionaryDetailScroll = 0
					return nil
				},
			})
		}

		// Fill remaining space
		writtenListLines := maxInt(0, minInt(len(m.dictionaryResults)-m.dictionaryScroll, maxResults))
		for i := writtenListLines; i < maxResults; i++ {
			listBuilder.WriteString(strings.Repeat(" ", layout.Width) + "\n")
		}

		// List scrollbar
		listLines := strings.Split(strings.TrimSuffix(listBuilder.String(), "\n"), "\n")
		listLinesWithScroll := renderScrollbarColumn(listLines, maxResults, len(m.dictionaryResults), m.dictionaryScroll)
		listWithScroll := strings.Join(listLinesWithScroll, "\n") + "\n"
		b.WriteString(listWithScroll)

		b.WriteString("\n" + mutedStyle.Render("Press ctrl+d/click selected to view details | Enter to draft | ctrl+a to add | ctrl+e to explain | ctrl+p to play"))
	}

	return b.String()
}

func keyHint(key, action string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(key) + " " + mutedStyle.Render(action)
}

func (m *Model) renderSpotlightDictionary() string {
	boxWidth := 86
	if m.width < 92 {
		boxWidth = m.width - 6
	}
	if boxWidth < 30 {
		boxWidth = 30
	}

	boxHeight := 18
	if m.height < 24 {
		boxHeight = m.height - 6
	}
	if boxHeight < 10 {
		boxHeight = 10
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
	b.WriteString(titleStr + "  " + closeHint + "\n")

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
		b.WriteString(searchBar.Render("🔍 "+searchText) + "\n")
	} else {
		// Clip query if it's too long for the search bar width
		availableWidth := searchBarWidth - 8
		displaySearch := searchText
		if len([]rune(searchText)) > availableWidth && availableWidth > 5 {
			displaySearch = "..." + string([]rune(searchText)[len([]rune(searchText))-availableWidth+3:])
		}
		contentLength := 3 + len([]rune(displaySearch)) + 1
		spaces := searchBarWidth - contentLength - 3
		clearText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true).Render("[x]")
		var clearBtn string
		if spaces > 0 {
			clearBtn = strings.Repeat(" ", spaces) + clearText
		} else {
			clearBtn = " " + clearText
		}
		b.WriteString(searchBar.Render("🔍 "+displaySearch+editStyle.Render("█")+clearBtn) + "\n")

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
	overlayFilterLineY := strings.Count(b.String(), "\n")
	b.WriteString(renderFilterPillsRow(m, startX+2, startY+overlayFilterLineY+1, m.activeView) + "\n")

	headerLines := strings.Count(b.String(), "\n") // Header takes 5 lines (title:1, searchBar:3, filterRow:1)
	interiorWidth := boxWidth - 4
	if interiorWidth < 10 {
		interiorWidth = 10
	}

	// Calculate maximum result lines that fit cleanly before the 1-line footer
	maxResults := boxHeight - 2 - headerLines - 1
	if maxResults < 1 {
		maxResults = 1
	}

	if len(m.dictionaryResults) == 0 {
		if m.dictionarySearch != "" {
			b.WriteString(mutedStyle.Render("No results found.\n"))
		} else {
			b.WriteString(mutedStyle.Render("Type to search local dict.cc dictionary.\n"))
			if len(m.dictionarySearchHistory) > 0 {
				clearHistoryText := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("[Clear]")
				b.WriteString(boldStyle.Render("Recent Searches:") + "  " + clearHistoryText + " " + mutedStyle.Render("(ctrl+x)") + "\n")

				clearHistoryLineY := strings.Count(b.String(), "\n") - 1
				m.hitboxes = append(m.hitboxes, Hitbox{
					ID:     "dict-overlay-history-clear",
					View:   m.activeView,
					X:      startX + 20,
					Y:      startY + clearHistoryLineY + 1,
					Width:  7,
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionarySearchHistory = nil
						m.saveDictionaryHistory()
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
						Y:      startY + lineY + 1,
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

			if len(m.dictionaryRecentlyViewed) > 0 {
				b.WriteString(boldStyle.Render("Recently Inspected Words:") + "\n")
				for i, e := range m.dictionaryRecentlyViewed {
					lineY := strings.Count(b.String(), "\n")
					wordStr := e.Word
					if e.Gender != "" {
						wordStr += " {" + e.Gender + "}"
					}
					if e.Translation != "" {
						wordStr += " - " + e.Translation
					}
					b.WriteString(fmt.Sprintf("  • %s\n", wordStr))
					entry := e
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-overlay-recent-%d", i),
						View:   m.activeView,
						X:      startX + 4,
						Y:      startY + lineY + 1,
						Width:  len([]rune(wordStr)),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = entry.Word
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

			if len(m.dictionaryDiscoverEntries) > 0 {
				discoverIcon := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("✦")
				b.WriteString(discoverIcon + " " + boldStyle.Render("Discover:") + "\n")
				for i, e := range m.dictionaryDiscoverEntries {
					lineY := strings.Count(b.String(), "\n")
					wordStr := e.Word
					if e.Gender != "" {
						wordStr += " " + renderGender(e.Gender)
					}
					if e.Translation != "" {
						trans := e.Translation
						if idx := strings.Index(trans, ";"); idx > 0 {
							trans = strings.TrimSpace(trans[:idx])
						}
						wordStr += " — " + mutedStyle.Render(trans)
					}
					b.WriteString(fmt.Sprintf("  • %s\n", wordStr))
					entry := e
					m.hitboxes = append(m.hitboxes, Hitbox{
						ID:     fmt.Sprintf("dict-overlay-discover-%d", i),
						View:   m.activeView,
						X:      startX + 4,
						Y:      startY + lineY + 1,
						Width:  lipgloss.Width(wordStr),
						Height: 1,
						Action: func() tea.Cmd {
							m.dictionarySearch = entry.Word
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
	} else {
		if m.dictionaryCursor < m.dictionaryScroll {
			m.dictionaryScroll = m.dictionaryCursor
		}
		if m.dictionaryCursor >= m.dictionaryScroll+maxResults {
			m.dictionaryScroll = m.dictionaryCursor - maxResults + 1
		}

		if m.dictionaryDetailView && m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
			res := m.dictionaryResults[m.dictionaryCursor]

			var detailBuilder strings.Builder
			titleText := res.Word
			if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
				titleText = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ") + titleText
			}
			detailBuilder.WriteString(titleStyle.Render(titleText) + "\n")
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
				detailBuilder.WriteString(formatInflectionTable(res.WordClass, res.Gender, res.Forms, m.dictionarySearch))
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

			// Detail content width (leaving room for scrollbar if needed)
			detailContentWidth := interiorWidth - 2

			var visibleDetailBuilder strings.Builder
			for i := m.dictionaryDetailScroll; i < m.dictionaryDetailScroll+maxResults && i < len(detailLines); i++ {
				visibleDetailBuilder.WriteString(padString(detailLines[i], detailContentWidth) + "\n")
			}
			writtenDetailLines := maxInt(0, minInt(len(detailLines)-m.dictionaryDetailScroll, maxResults))
			for i := writtenDetailLines; i < maxResults; i++ {
				visibleDetailBuilder.WriteString(strings.Repeat(" ", detailContentWidth) + "\n")
			}

			renderedDetailLines := strings.Split(strings.TrimSuffix(visibleDetailBuilder.String(), "\n"), "\n")
			detailLinesWithScroll := renderScrollbarColumn(renderedDetailLines, maxResults, m.dictionaryDetailTotalLines, m.dictionaryDetailScroll)
			detailPanelContent := strings.Join(detailLinesWithScroll, "\n") + "\n"

			b.WriteString(detailPanelContent)
		} else if interiorWidth > 70 {
			listWidth := maxInt(25, minInt(40, interiorWidth*4/10))
			detailWidth := interiorWidth - listWidth - 3
			listStartLine := strings.Count(b.String(), "\n")

			var listBuilder strings.Builder
			for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
				res := m.dictionaryResults[i]
				prefix := "  "
				if i == m.dictionaryCursor {
					prefix = "> "
				}
				starStr := ""
				if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
					starStr = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ")
				}
				wordText := res.Word
				if res.Gender != "" {
					wordText += " {" + res.Gender + "}"
				} else if res.WordClass != "" {
					wordText += " [" + res.WordClass + "]"
				}
				if res.Translation != "" {
					wordText = fmt.Sprintf("%s - %s", wordText, res.Translation)
				}
				padded := padString(wordText, maxInt(5, listWidth-2-lipgloss.Width(starStr)))
				highlighted := highlightQuery(padded, m.dictionarySearch, dictHighlightStyle)
				line := prefix + starStr + highlighted
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
					Y:      startY + 1 + listStartLine + (idx - m.dictionaryScroll),
					Width:  listWidth,
					Height: 1,
					Action: func() tea.Cmd {
						m.dictionaryCursor = idx
						m.dictionaryDetailScroll = 0
						m.inspectDictionaryCursor()
						if m.dictionaryDetailVisible() && idx >= 0 && idx < len(m.dictionaryResults) {
							return m.findRelatedEntries(m.dictionaryResults[idx].Word)
						}
						return nil
					},
				})
			}
			writtenListLines := maxInt(0, minInt(len(m.dictionaryResults)-m.dictionaryScroll, maxResults))
			for i := writtenListLines; i < maxResults; i++ {
				listBuilder.WriteString(strings.Repeat(" ", listWidth) + "\n")
			}

			listLines := strings.Split(strings.TrimSuffix(listBuilder.String(), "\n"), "\n")
			listLinesWithScroll := renderScrollbarColumn(listLines, maxResults, len(m.dictionaryResults), m.dictionaryScroll)
			listWithScroll := strings.Join(listLinesWithScroll, "\n") + "\n"

			var detailBuilder strings.Builder
			if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
				res := m.dictionaryResults[m.dictionaryCursor]
				titleText := res.Word
				if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
					titleText = lipgloss.NewStyle().Foreground(colorYellow).Render("★ ") + titleText
				}
				detailBuilder.WriteString(titleStyle.Render(titleText) + "\n")
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
				visibleDetailBuilder.WriteString(padString(detailLines[i], detailWidth-3) + "\n")
			}
			writtenDetailLines := maxInt(0, minInt(len(detailLines)-m.dictionaryDetailScroll, maxResults))
			for i := writtenDetailLines; i < maxResults; i++ {
				visibleDetailBuilder.WriteString(strings.Repeat(" ", detailWidth-3) + "\n")
			}

			renderedDetailLines := strings.Split(strings.TrimSuffix(visibleDetailBuilder.String(), "\n"), "\n")
			detailLinesWithScroll := renderScrollbarColumn(renderedDetailLines, maxResults, m.dictionaryDetailTotalLines, m.dictionaryDetailScroll)
			detailPanelContent := strings.Join(detailLinesWithScroll, "\n") + "\n"

			detailPanel := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(colorPanel).
				Padding(0, 1).
				Width(detailWidth).
				Height(maxResults).
				Render(detailPanelContent)

			joined := lipgloss.JoinHorizontal(lipgloss.Top, listWithScroll, detailPanel)
			b.WriteString(joined)
		} else {
			listStartLine := strings.Count(b.String(), "\n")
			var listBuilder strings.Builder
			for i := m.dictionaryScroll; i < len(m.dictionaryResults) && i < m.dictionaryScroll+maxResults; i++ {
				res := m.dictionaryResults[i]
				prefix := "  "
				if i == m.dictionaryCursor {
					prefix = "> "
				}
				starStr := ""
				if m.dictionaryStarred != nil && m.dictionaryStarred[res.ID] {
					starStr = "★ "
				}
				plainLine := fmt.Sprintf("%s%s - %s", starStr, res.Word, res.Translation)
				padded := padString(plainLine, maxInt(5, interiorWidth-3))
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
					Y:      startY + 1 + listStartLine + (idx - m.dictionaryScroll),
					Width:  interiorWidth,
					Height: 1,
					Action: func() tea.Cmd {
						if m.dictionaryCursor == idx {
							m.dictionaryDetailView = true
							m.dictionaryDetailScroll = 0
							m.inspectDictionaryCursor()
							if idx >= 0 && idx < len(m.dictionaryResults) {
								return m.findRelatedEntries(m.dictionaryResults[idx].Word)
							}
							return nil
						}
						m.dictionaryCursor = idx
						m.dictionaryDetailScroll = 0
						return nil
					},
				})
			}
			writtenListLines := maxInt(0, minInt(len(m.dictionaryResults)-m.dictionaryScroll, maxResults))
			for i := writtenListLines; i < maxResults; i++ {
				listBuilder.WriteString(strings.Repeat(" ", maxInt(5, interiorWidth-1)) + "\n")
			}

			listLines := strings.Split(strings.TrimSuffix(listBuilder.String(), "\n"), "\n")
			listLinesWithScroll := renderScrollbarColumn(listLines, maxResults, len(m.dictionaryResults), m.dictionaryScroll)
			listWithScroll := strings.Join(listLinesWithScroll, "\n") + "\n"
			b.WriteString(listWithScroll)
		}
	}

	targetBodyLines := boxHeight - 3
	currentLines := strings.Count(b.String(), "\n")
	if currentLines < targetBodyLines {
		b.WriteString(strings.Repeat("\n", targetBodyLines-currentLines))
	}

	// Always render key hints in footer so user knows available shortcuts
	var keys []string
	if m.dictionaryDetailView {
		keys = []string{
			keyHint("esc", "back"),
			keyHint("Enter", "draft"),
			keyHint("ctrl+a", "add"),
			keyHint("c", "cloze"),
			keyHint("b", "star"),
			keyHint("ctrl+e", "explain"),
			keyHint("ctrl+p", "play"),
			keyHint("ctrl+g", "deck"),
		}
	} else {
		keys = []string{
			keyHint("Enter", "draft"),
			keyHint("ctrl+a", "add"),
			keyHint("c", "cloze"),
			keyHint("b", "star"),
			keyHint("ctrl+e", "explain"),
			keyHint("ctrl+p", "play"),
			keyHint("ctrl+d", "details"),
			keyHint("ctrl+g", "deck"),
		}
	}
	footerStr := strings.Join(keys, " │ ")

	// If status message is present, show status on top of footer or appended nicely
	if m.status != "" {
		statusText := m.status
		var statusStyle lipgloss.Style
		lowerStatus := strings.ToLower(statusText)
		if strings.Contains(lowerStatus, "added") || strings.Contains(lowerStatus, "found") || strings.Contains(lowerStatus, "ready") || strings.Contains(lowerStatus, "starred") {
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
		} else if strings.Contains(lowerStatus, "no") || strings.Contains(lowerStatus, "failed") || strings.Contains(lowerStatus, "error") {
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
		} else {
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
		}
		// Pad status text nicely next to key hints if width permits
		if lipgloss.Width(footerStr)+lipgloss.Width(statusText)+4 <= interiorWidth {
			footerStr = footerStr + "  " + statusStyle.Render(statusText)
		} else {
			footerStr = statusStyle.Render(statusText)
		}
	}
	b.WriteString(footerStr)

	return boxStyle.Render(b.String())
}

func formatInflectionTable(wordClass, gender, formsStr, query string) string {
	if strings.TrimSpace(formsStr) == "" {
		return ""
	}
	parts := strings.FieldsFunc(formsStr, func(r rune) bool {
		return r == ';' || r == ','
	})
	var cleanParts []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			cleanParts = append(cleanParts, t)
		}
	}

	wc := strings.ToLower(wordClass)
	gen := strings.ToLower(gender)

	var sb strings.Builder

	if (strings.Contains(wc, "verb") || strings.Contains(wc, "v")) && len(cleanParts) == 3 {
		sb.WriteString(boldStyle.Render("Verb Forms:") + "\n")
		p3sg := highlightQuery(cleanParts[0], query, dictHighlightStyle)
		prat := highlightQuery(cleanParts[1], query, dictHighlightStyle)
		perf := highlightQuery(cleanParts[2], query, dictHighlightStyle)
		sb.WriteString(fmt.Sprintf("  • Präsens (3sg): %s\n", p3sg))
		sb.WriteString(fmt.Sprintf("  • Präteritum:    %s\n", prat))
		sb.WriteString(fmt.Sprintf("  • Perfekt:       %s\n\n", perf))
		return sb.String()
	}

	if (strings.Contains(wc, "noun") || gen == "m" || gen == "f" || gen == "n" || gen == "pl") && len(cleanParts) == 2 {
		sb.WriteString(boldStyle.Render("Noun Forms:") + "\n")
		genitiv := highlightQuery(cleanParts[0], query, dictHighlightStyle)
		plural := highlightQuery(cleanParts[1], query, dictHighlightStyle)
		sb.WriteString(fmt.Sprintf("  • Genitiv: %s\n", genitiv))
		sb.WriteString(fmt.Sprintf("  • Plural:  %s\n\n", plural))
		return sb.String()
	}

	if (strings.Contains(wc, "adj") || (len(cleanParts) == 2 && strings.HasPrefix(cleanParts[len(cleanParts)-1], "am "))) && len(cleanParts) == 2 {
		sb.WriteString(boldStyle.Render("Adjective Comparison:") + "\n")
		komp := highlightQuery(cleanParts[0], query, dictHighlightStyle)
		sup := highlightQuery(cleanParts[1], query, dictHighlightStyle)
		sb.WriteString(fmt.Sprintf("  • Komparativ: %s\n", komp))
		sb.WriteString(fmt.Sprintf("  • Superlativ: %s\n\n", sup))
		return sb.String()
	}

	sb.WriteString(boldStyle.Render("Word Forms:") + "\n")
	highlightedForms := highlightQuery(formsStr, query, dictHighlightStyle)
	sb.WriteString("  " + highlightedForms + "\n\n")
	return sb.String()
}
