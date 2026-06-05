package tui

import (
	"fmt"
	"strings"
)

func (m *Model) renderDictionary(layout viewportLayout) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Dictionary"))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(titleStyle.Render("Search:"))
	b.WriteString(" ")
	searchLine := m.dictionarySearch
	if searchLine == "" {
		searchLine = mutedStyle.Render("Type to search...")
	} else {
		searchLine += editStyle.Render(" ")
	}
	b.WriteString(searchLine)
	b.WriteString("\n\n")

	// Results
	if len(m.dictionaryResults) == 0 {
		if m.dictionarySearch != "" {
			b.WriteString(mutedStyle.Render("No results found."))
		} else {
			b.WriteString(mutedStyle.Render("Search dictionary (local dict.cc)."))
		}
	} else {
		for i, res := range m.dictionaryResults {
			if i >= layout.Height-8 { // basic pagination logic
				b.WriteString(mutedStyle.Render(fmt.Sprintf("+ %d more...", len(m.dictionaryResults)-i)))
				break
			}
			prefix := "  "
			if i == m.dictionaryCursor {
				prefix = "> "
			}
			line := fmt.Sprintf("%s%s - %s", prefix, res.Word, res.Translation)
			if res.WordClass != "" {
				line += mutedStyle.Render(" [" + res.WordClass + "]")
			}
			if res.Gender != "" {
				line += mutedStyle.Render(" {" + res.Gender + "}")
			}
			if i == m.dictionaryCursor {
				b.WriteString(editStyle.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
	}

	return b.String()
}
