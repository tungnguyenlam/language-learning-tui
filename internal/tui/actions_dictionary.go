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

		starredOnly := isFilterActive(query, ":starred") || isFilterActive(query, ":star") || isFilterActive(query, ":fav") || isFilterActive(query, ":favorite")

		cleanText := clearFilterTags(query)
		if starredOnly && cleanText == "" {
			var results []core.DictionaryEntry
			for entryID := range m.dictionaryStarred {
				entry, err := dictRepo.GetEntry(ctx, entryID)
				if err == nil {
					results = append(results, entry)
				}
			}
			return dictionarySearchResultsMsg{id: id, results: results}
		}

		results, err := dictRepo.Search(ctx, query, 50)
		if err != nil {
			return err
		}

		if starredOnly {
			var filtered []core.DictionaryEntry
			for _, e := range results {
				if m.dictionaryStarred[e.ID] {
					filtered = append(filtered, e)
				}
			}
			results = filtered
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
			if m.dictionaryTargetDeckID != "" {
				deckID = m.dictionaryTargetDeckID
				for _, d := range m.decks {
					if d.ID == deckID {
						deckName = d.Name
						break
					}
				}
			} else if m.deck.ID != "" && m.deck.ID != "all" {
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

func isFilterActive(query, tag string) bool {
	if tag == "" {
		filterTags := map[string]bool{
			":verb": true, ":v": true, ":noun": true, ":adj": true, ":adv": true,
			":m": true, ":f": true, ":n": true, ":pl": true,
			":starred": true, ":star": true, ":fav": true, ":favorite": true,
		}
		for _, t := range strings.Fields(query) {
			if filterTags[strings.ToLower(t)] {
				return false
			}
		}
		return true
	}
	for _, t := range strings.Fields(query) {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

func toggleFilterTag(query, tag string) string {
	terms := strings.Fields(query)
	found := false
	var newTerms []string
	for _, t := range terms {
		if strings.EqualFold(t, tag) {
			found = true
		} else {
			newTerms = append(newTerms, t)
		}
	}
	if !found && tag != "" {
		newTerms = append(newTerms, tag)
	}
	return strings.Join(newTerms, " ")
}

func clearFilterTags(query string) string {
	filterTags := map[string]bool{
		":verb": true, ":v": true, ":noun": true, ":adj": true, ":adv": true,
		":m": true, ":f": true, ":n": true, ":pl": true,
		":de": true, ":en": true, "lang:de": true, "lang:en": true,
		":starred": true, ":star": true, ":fav": true, ":favorite": true,
	}
	terms := strings.Fields(query)
	var clean []string
	for _, t := range terms {
		if !filterTags[strings.ToLower(t)] && !strings.HasPrefix(strings.ToLower(t), "de:") && !strings.HasPrefix(strings.ToLower(t), "en:") {
			clean = append(clean, t)
		}
	}
	return strings.Join(clean, " ")
}

func (m *Model) toggleStarDictionaryEntry(entry core.DictionaryEntry) tea.Cmd {
	if m.dictionaryStarred == nil {
		m.dictionaryStarred = make(map[string]bool)
	}
	starred := !m.dictionaryStarred[entry.ID]
	if starred {
		m.dictionaryStarred[entry.ID] = true
	} else {
		delete(m.dictionaryStarred, entry.ID)
	}

	var ids []string
	for id := range m.dictionaryStarred {
		ids = append(ids, id)
	}
	data, err := json.Marshal(ids)
	if err == nil {
		ctx := context.Background()
		_ = m.repo.SetSetting(ctx, "dict_starred_entries", string(data))
	}

	if starred {
		m.status = fmt.Sprintf("★ Starred %q", entry.Word)
	} else {
		m.status = fmt.Sprintf("Unstarred %q", entry.Word)
	}

	if isFilterActive(m.dictionarySearch, ":starred") || isFilterActive(m.dictionarySearch, ":star") || isFilterActive(m.dictionarySearch, ":fav") {
		return m.searchDictionary()
	}
	return nil
}

func (m *Model) loadDictionaryStarred() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		raw, err := m.repo.GetSetting(ctx, "dict_starred_entries")
		if err != nil || raw == "" {
			return dictStarredLoadedMsg{}
		}
		var ids []string
		if err := json.Unmarshal([]byte(raw), &ids); err != nil {
			return dictStarredLoadedMsg{}
		}
		starred := make(map[string]bool, len(ids))
		for _, id := range ids {
			starred[id] = true
		}
		return dictStarredLoadedMsg(starred)
	}
}

func (m *Model) cycleDictionaryTargetDeck() tea.Cmd {
	if len(m.decks) == 0 {
		m.dictionaryTargetDeckID = "dictionary"
		m.status = "Quick-add target deck: Dictionary"
		return nil
	}

	currentIndex := -1
	if m.dictionaryTargetDeckID == "dictionary" || m.dictionaryTargetDeckID == "" {
		currentIndex = -1
	} else {
		for i, d := range m.decks {
			if d.ID == m.dictionaryTargetDeckID {
				currentIndex = i
				break
			}
		}
	}

	nextIndex := currentIndex + 1
	if nextIndex >= len(m.decks) {
		m.dictionaryTargetDeckID = "dictionary"
		m.status = "Quick-add target deck: Dictionary"
	} else {
		target := m.decks[nextIndex]
		m.dictionaryTargetDeckID = target.ID
		m.status = fmt.Sprintf("Quick-add target deck: %s", target.Name)
	}
	return nil
}

func (m *Model) addDictionaryEntriesBatchCmd(entries []core.DictionaryEntry) tea.Cmd {
	if len(entries) == 0 {
		m.status = "No dictionary entries to add"
		return nil
	}
	return tea.Batch(
		func() tea.Msg {
			ctx := context.Background()
			deckID := "dictionary"
			deckName := "Dictionary"
			if m.dictionaryTargetDeckID != "" {
				deckID = m.dictionaryTargetDeckID
				for _, d := range m.decks {
					if d.ID == deckID {
						deckName = d.Name
						break
					}
				}
			} else if m.deck.ID != "" && m.deck.ID != "all" {
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

			count := 0
			for _, entry := range entries {
				noteID := "dict-" + entry.ID
				if entry.ID == "" {
					noteID = fmt.Sprintf("dict-%d-%d", time.Now().UnixNano(), count)
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
				if err := m.repo.UpsertNote(ctx, note); err == nil {
					count++
				}
			}

			return statusMsg{text: fmt.Sprintf("Added %d entries to %s deck", count, deckName)}
		},
		m.loadDecks,
		m.loadDueCards,
	)
}

func (m *Model) recordDictionaryView(entry core.DictionaryEntry) {
	if entry.ID == "" && entry.Word == "" {
		return
	}
	var clean []core.DictionaryEntry
	for _, e := range m.dictionaryRecentlyViewed {
		if e.ID != entry.ID && e.Word != entry.Word {
			clean = append(clean, e)
		}
	}
	m.dictionaryRecentlyViewed = append([]core.DictionaryEntry{entry}, clean...)
	if len(m.dictionaryRecentlyViewed) > 10 {
		m.dictionaryRecentlyViewed = m.dictionaryRecentlyViewed[:10]
	}
}

func stripWordArticle(word string) string {
	w := strings.TrimSpace(strings.ToLower(word))
	for _, art := range []string{"der ", "die ", "das ", "den ", "dem ", "des ", "the ", "a ", "an "} {
		if strings.HasPrefix(w, art) {
			return strings.TrimSpace(w[len(art):])
		}
	}
	return w
}

func (m *Model) addDictionaryClozeEntryCmd(entry core.DictionaryEntry) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			ctx := context.Background()

			deckID := "dictionary"
			deckName := "Dictionary"
			if m.dictionaryTargetDeckID != "" {
				deckID = m.dictionaryTargetDeckID
				for _, d := range m.decks {
					if d.ID == deckID {
						deckName = d.Name
						break
					}
				}
			} else if m.deck.ID != "" && m.deck.ID != "all" {
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

			noteID := "dict-cloze-" + entry.ID
			if entry.ID == "" {
				noteID = fmt.Sprintf("dict-cloze-%d", time.Now().UnixNano())
			}

			wordClean := entry.Word
			bareWord := stripWordArticle(wordClean)
			frontText := ""

			for _, ex := range entry.Examples {
				idx := strings.Index(strings.ToLower(ex), strings.ToLower(bareWord))
				if idx != -1 {
					matchStr := ex[idx : idx+len(bareWord)]
					frontText = ex[:idx] + "{{c1::" + matchStr + "}}" + ex[idx+len(bareWord):]
					break
				}
			}

			if frontText == "" {
				frontText = fmt.Sprintf("Das deutsche Wort für '%s' ist {{c1::%s}}.", entry.Translation, formatDictionaryCardFront(entry.Word, entry.Gender))
			}

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
				Type:   "cloze",
				Front:  frontText,
				Back:   entry.Translation,
				Extra:  strings.Join(extraParts, "\n"),
				Tags:   append(entry.Tags, "dictionary", "cloze"),
			}
			note.Cards = content.CardsForNote(note)

			if err := m.repo.UpsertNote(ctx, note); err != nil {
				return err
			}

			return statusMsg{text: fmt.Sprintf("Created Cloze card for '%s' in %s deck", wordClean, deck.Name)}
		},
		m.loadDecks,
		m.loadDueCards,
	)
}
