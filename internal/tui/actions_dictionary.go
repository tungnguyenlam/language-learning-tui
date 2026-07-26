package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

type dictionarySearchResultsMsg struct {
	id      int
	results []core.DictionaryEntry
}

func (m *Model) searchDictionary() tea.Cmd {
	m.dictionarySearchID++
	id := m.dictionarySearchID
	query := m.dictionarySearch
	return func() tea.Msg {
		if query == "" {
			return dictionarySearchResultsMsg{id: id, results: nil}
		}

		dictRepo, ok := m.repo.(core.DictionaryRepository)
		if !ok {
			return fmt.Errorf("repository does not support dictionary search")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		results, err := dictRepo.Search(ctx, query, 50)
		if err != nil {
			return err
		}
		return dictionarySearchResultsMsg{id: id, results: results}
	}
}

var (
	htmlTagRegex  = regexp.MustCompile(`<[^>]+>`)
	clozeTagRegex = regexp.MustCompile(`\{\{c\d+::([^:}]+)(?:::[^}]*)?\}\}`)
	parenRegex    = regexp.MustCompile(`\([^)]*\)`)
	bracketRegex  = regexp.MustCompile(`\[[^\]]*\]`)
)

func cleanLookupQuery(raw string) string {
	s := raw
	if matches := clozeTagRegex.FindStringSubmatch(s); len(matches) > 1 {
		s = matches[1]
	}

	s = htmlTagRegex.ReplaceAllString(s, "")

	if idx := strings.Index(s, "\n"); idx != -1 {
		s = s[:idx]
	}

	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}

	s = parenRegex.ReplaceAllString(s, "")
	s = bracketRegex.ReplaceAllString(s, "")

	s = strings.TrimLeft(s, "•-–—*# \t")
	s = strings.TrimRight(s, " \t?!.,;:\"'’")

	fields := strings.Fields(s)
	if len(fields) == 0 {
		return strings.TrimSpace(raw)
	}
	return strings.Join(fields, " ")
}

func formatDictionaryCardFront(word, gender string) string {
	w := strings.TrimSpace(word)
	lower := strings.ToLower(w)
	if strings.HasPrefix(lower, "der ") || strings.HasPrefix(lower, "die ") || strings.HasPrefix(lower, "das ") {
		return w
	}
	switch strings.ToLower(strings.TrimSpace(gender)) {
	case "m", "masc", "der":
		return "der " + w
	case "f", "fem", "die":
		return "die " + w
	case "n", "neut", "das":
		return "das " + w
	case "pl", "plural":
		return "die " + w + " (pl.)"
	}
	return w
}

func (m *Model) addDictionaryEntryCmd(entry core.DictionaryEntry) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			ctx := context.Background()

			// 1. Ensure target deck exists
			deckID := "dictionary"
			deckName := "Dictionary"
			if m.deck.ID != "" && m.deck.ID != "all" {
				deckID = m.deck.ID
				deckName = m.deck.Name
			}
			deck, err := m.repo.GetDeck(ctx, deckID)
			if err != nil {
				deck = core.Deck{
					ID:   deckID,
					Name: deckName,
				}
				if err := m.repo.UpsertDeck(ctx, deck); err != nil {
					return err
				}
			}

			// 2. Create note
			noteID := "dict-" + entry.ID
			if entry.ID == "" {
				noteID = fmt.Sprintf("dict-%d", time.Now().UnixNano())
			}

			frontText := formatDictionaryCardFront(entry.Word, entry.Gender)

			extraParts := []string{}
			if entry.Forms != "" {
				extraParts = append(extraParts, "Forms: "+entry.Forms)
			}
			if entry.WordClass != "" {
				extraParts = append(extraParts, "Class: ["+strings.ToUpper(entry.WordClass)+"]")
			}
			if entry.Gender != "" {
				extraParts = append(extraParts, "Gender: {"+entry.Gender+"}")
			}
			if len(entry.Examples) > 0 {
				extraParts = append(extraParts, "Examples:\n• "+strings.Join(entry.Examples, "\n• "))
			}

			note := core.Note{
				ID:     noteID,
				DeckID: deckID,
				Type:   "flashcard",
				Front:  frontText,
				Back:   entry.Translation,
				Extra:  strings.Join(extraParts, "\n"),
				Tags:   append(entry.Tags, "dictionary"),
			}
			note.Cards = content.CardsForNote(note)

			if err := m.repo.UpsertNote(ctx, note); err != nil {
				return err
			}

			return statusMsg{text: fmt.Sprintf("Added '%s' to %s deck", frontText, deck.Name)}
		},
		m.loadDecks,
		m.loadDueCards,
	)
}

func (m *Model) recordDictionarySearch(query string) {
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}
	for i, q := range m.dictionarySearchHistory {
		if q == query {
			m.dictionarySearchHistory = append(m.dictionarySearchHistory[:i], m.dictionarySearchHistory[i+1:]...)
			m.dictionarySearchHistory = append(m.dictionarySearchHistory, query)
			m.saveDictionaryHistory()
			return
		}
	}
	m.dictionarySearchHistory = append(m.dictionarySearchHistory, query)
	if len(m.dictionarySearchHistory) > 5 {
		m.dictionarySearchHistory = m.dictionarySearchHistory[1:]
	}
	m.saveDictionaryHistory()
}

func (m *Model) cycleDictionaryHistory(direction int) tea.Cmd {
	if len(m.dictionarySearchHistory) == 0 {
		return nil
	}
	currentIndex := -1
	for i, q := range m.dictionarySearchHistory {
		if q == m.dictionarySearch {
			currentIndex = i
			break
		}
	}
	var nextIndex int
	if currentIndex == -1 {
		if direction > 0 {
			nextIndex = len(m.dictionarySearchHistory) - 1
		} else {
			nextIndex = 0
		}
	} else {
		nextIndex = currentIndex + direction
		if nextIndex < 0 {
			nextIndex = 0
		}
		if nextIndex >= len(m.dictionarySearchHistory) {
			nextIndex = len(m.dictionarySearchHistory) - 1
		}
	}
	m.dictionarySearch = m.dictionarySearchHistory[nextIndex]
	m.dictionaryCursor = 0
	m.dictionaryScroll = 0
	m.dictionaryDetailScroll = 0
	m.dictionaryDetailView = false
	return m.searchDictionary()
}

func (m *Model) saveDictionaryHistory() {
	ctx := context.Background()
	data, err := json.Marshal(m.dictionarySearchHistory)
	if err == nil {
		_ = m.repo.SetSetting(ctx, "dict_search_history", string(data))
	}
}

func (m *Model) loadDictionaryHistory() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		raw, err := m.repo.GetSetting(ctx, "dict_search_history")
		if err != nil || raw == "" {
			return dictHistoryLoadedMsg{}
		}
		var history []string
		if err := json.Unmarshal([]byte(raw), &history); err != nil {
			history = strings.Split(raw, "\n")
		}
		var clean []string
		for _, s := range history {
			s = strings.TrimSpace(s)
			if s != "" {
				clean = append(clean, s)
			}
		}
		return dictHistoryLoadedMsg(clean)
	}
}

func (m *Model) resetDictionarySearchState() {
	m.dictionarySearch = ""
	m.dictionaryResults = nil
	m.dictionaryCursor = 0
	m.dictionaryScroll = 0
	m.dictionaryDetailScroll = 0
	m.dictionaryDetailTotalLines = 0
	m.dictionaryDetailView = false
}

func (m *Model) openDictionaryOverlay() tea.Cmd {
	return m.openDictionaryOverlayWithQuery("")
}

func (m *Model) openDictionaryOverlayWithQuery(query string) tea.Cmd {
	m.dictionaryOverlayActive = true
	m.resetDictionarySearchState()
	clean := cleanLookupQuery(query)
	if clean != "" {
		m.dictionarySearch = clean
		m.status = fmt.Sprintf("Spotlight lookup: %s", clean)
		return m.searchDictionary()
	}
	m.status = "Spotlight dictionary open"
	return nil
}

func (m *Model) closeDictionaryOverlay() {
	m.dictionaryOverlayActive = false
	m.resetDictionarySearchState()
	m.status = "Spotlight dictionary closed"
}

func (m *Model) lookupWordOfTheDayInDictionary() tea.Cmd {
	word := content.GetWordOfTheDay()
	return m.openDictionaryOverlayWithQuery(word.German)
}

func (m *Model) lookupVerbOfTheDayInDictionary() tea.Cmd {
	verb := content.GetVerbOfTheDay()
	return m.openDictionaryOverlayWithQuery(verb.German)
}

func (m *Model) lookupGrammarTipInDictionary() tea.Cmd {
	tip := content.GetDailyGrammarTip()
	return m.openDictionaryOverlayWithQuery(tip.Title)
}

func (m *Model) lookupReviewCardInDictionary() tea.Cmd {
	if len(m.dueCards) == 0 {
		return nil
	}
	card := m.dueCards[clampInt(m.cursor, 0, len(m.dueCards)-1)]
	return m.openDictionaryOverlayWithQuery(card.Prompt)
}

func (m *Model) lookupBrowserCardInDictionary() tea.Cmd {
	if len(m.browserCards) == 0 {
		return nil
	}
	card := m.browserCards[clampInt(m.browserCursor, 0, len(m.browserCards)-1)]
	return m.openDictionaryOverlayWithQuery(card.Prompt)
}

func (m *Model) lookupCramCardInDictionary() tea.Cmd {
	if len(m.cramCards) == 0 {
		return nil
	}
	card := m.cramCards[clampInt(m.cramCursor, 0, len(m.cramCards)-1)]
	return m.openDictionaryOverlayWithQuery(card.Prompt)
}
