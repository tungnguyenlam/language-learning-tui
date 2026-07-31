package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

type dictionarySearchResultsMsg struct {
	id      int
	results []core.DictionaryEntry
}

type dictRelatedEntriesMsg struct {
	entries []core.DictionaryEntry
}

type dictDiscoverEntriesMsg struct {
	entries []core.DictionaryEntry
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
			sort.Slice(results, func(i, j int) bool {
				if results[i].Word == results[j].Word {
					return results[i].ID < results[j].ID
				}
				return strings.ToLower(results[i].Word) < strings.ToLower(results[j].Word)
			})
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

func (m *Model) findRelatedEntries(word string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dictRepo, ok := m.repo.(core.DictionaryRepository)
		if !ok {
			return nil
		}

		entries, err := dictRepo.FindRelatedEntries(ctx, word, 5)
		if err != nil {
			return nil
		}
		return dictRelatedEntriesMsg{entries: entries}
	}
}

func (m *Model) loadDictionaryDiscoverEntries() tea.Cmd {
	return func() tea.Msg {
		dictRepo, ok := m.repo.(core.DictionaryRepository)
		if !ok {
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		entries, err := dictRepo.RandomEntries(ctx, 5)
		if err != nil {
			return nil
		}
		return dictDiscoverEntriesMsg{entries: entries}
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
	return content.FormatDictionaryCardFront(word, gender)
}

func (m *Model) dictionaryTargetDeck() (deckID, deckName string) {
	deckID = "dictionary"
	deckName = "Dictionary"
	if m.dictionaryTargetDeckID != "" {
		deckID = m.dictionaryTargetDeckID
		for _, d := range m.decks {
			if d.ID == deckID {
				deckName = d.Name
				break
			}
		}
		if deckID == "dictionary" {
			deckName = "Dictionary"
		}
		return deckID, deckName
	}
	if m.deck.ID != "" && m.deck.ID != "all" {
		return m.deck.ID, m.deck.Name
	}
	return deckID, deckName
}

func (m *Model) ensureDictionaryTargetDeck(ctx context.Context) (core.Deck, error) {
	deckID, deckName := m.dictionaryTargetDeck()
	deck, err := m.repo.GetDeck(ctx, deckID)
	if err != nil {
		deck = core.Deck{ID: deckID, Name: deckName}
		if err := m.repo.UpsertDeck(ctx, deck); err != nil {
			return core.Deck{}, err
		}
	}
	return deck, nil
}

func (m *Model) addDictionaryEntryCmd(entry core.DictionaryEntry) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			ctx := context.Background()
			deck, err := m.ensureDictionaryTargetDeck(ctx)
			if err != nil {
				return err
			}

			noteID := "dict-" + entry.ID
			if entry.ID == "" {
				noteID = fmt.Sprintf("dict-%d", time.Now().UnixNano())
			}

			frontText := formatDictionaryCardFront(entry.Word, entry.Gender)

			// Check for duplicates
			if m.dictionaryLastAddAttemptNoteID != noteID {
				if _, err := m.repo.GetNote(ctx, noteID); err == nil {
					m.dictionaryLastAddAttemptNoteID = noteID
					return statusMsg{text: fmt.Sprintf("'%s' already in %s deck (use a/ctrl+a again to update)", frontText, deck.Name)}
				}
			}
			m.dictionaryLastAddAttemptNoteID = ""

			tags := append([]string(nil), entry.Tags...)
			tags = append(tags, "dictionary")
			note := core.Note{
				ID:     noteID,
				DeckID: deck.ID,
				Type:   "flashcard",
				Front:  frontText,
				Back:   entry.Translation,
				Extra:  content.DictionaryEntryExtra(entry),
				Tags:   tags,
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

// addDictionaryReverseEntryCmd creates an EN→DE production flashcard
// (English translation on the front, German headword on the back).
func (m *Model) addDictionaryReverseEntryCmd(entry core.DictionaryEntry) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			ctx := context.Background()
			deck, err := m.ensureDictionaryTargetDeck(ctx)
			if err != nil {
				return err
			}

			noteID := "dict-rev-" + entry.ID
			if entry.ID == "" {
				noteID = fmt.Sprintf("dict-rev-%d", time.Now().UnixNano())
			}

			backText := formatDictionaryCardFront(entry.Word, entry.Gender)
			frontText := strings.TrimSpace(entry.Translation)
			if frontText == "" {
				return statusMsg{text: "Cannot create reverse card without a translation"}
			}
			// Prefer the first sense when multiple translations are semicolon-separated.
			if parts := strings.Split(frontText, ";"); len(parts) > 0 {
				frontText = strings.TrimSpace(parts[0])
			}

			// Check for duplicates
			if m.dictionaryLastAddAttemptNoteID != noteID {
				if _, err := m.repo.GetNote(ctx, noteID); err == nil {
					m.dictionaryLastAddAttemptNoteID = noteID
					return statusMsg{text: fmt.Sprintf("'%s' already in %s deck (use r/ctrl+r again to update)", frontText, deck.Name)}
				}
			}
			m.dictionaryLastAddAttemptNoteID = ""

			tags := append([]string(nil), entry.Tags...)
			tags = append(tags, "dictionary", "reverse")
			note := core.Note{
				ID:     noteID,
				DeckID: deck.ID,
				Type:   "flashcard",
				Front:  frontText,
				Back:   backText,
				Extra:  content.DictionaryEntryExtra(entry),
				Tags:   tags,
			}
			note.Cards = content.CardsForNote(note)

			if err := m.repo.UpsertNote(ctx, note); err != nil {
				return err
			}

			return statusMsg{text: fmt.Sprintf("Added reverse '%s → %s' to %s deck", frontText, backText, deck.Name)}
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

func (m *Model) saveDictionaryRecentlyViewed() {
	data, err := json.Marshal(m.dictionaryRecentlyViewed)
	if err != nil {
		return
	}
	_ = m.repo.SetSetting(context.Background(), "dict_recently_viewed", string(data))
}

func (m *Model) clearDictionaryRecentlyViewed() {
	m.dictionaryRecentlyViewed = nil
	m.saveDictionaryRecentlyViewed()
	m.status = "Cleared recently inspected words"
}

func (m *Model) loadDictionaryRecentlyViewed() tea.Cmd {
	return func() tea.Msg {
		raw, err := m.repo.GetSetting(context.Background(), "dict_recently_viewed")
		if err != nil || raw == "" {
			return dictRecentlyViewedLoadedMsg{}
		}
		var entries []core.DictionaryEntry
		if err := json.Unmarshal([]byte(raw), &entries); err != nil {
			return dictRecentlyViewedLoadedMsg{}
		}
		if len(entries) > 10 {
			entries = entries[:10]
		}
		return dictRecentlyViewedLoadedMsg(entries)
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
	if len(m.dictionaryDiscoverEntries) == 0 {
		return m.loadDictionaryDiscoverEntries()
	}
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

func isLangFilterTerm(term string) bool {
	lower := strings.ToLower(term)
	return lower == ":de" || lower == ":en" || lower == "lang:de" || lower == "lang:en" ||
		strings.HasPrefix(lower, "de:") || strings.HasPrefix(lower, "en:")
}

func isFilterActive(query, tag string) bool {
	if tag == "" {
		filterTags := map[string]bool{
			":verb": true, ":v": true, ":noun": true, ":adj": true, ":adv": true,
			":m": true, ":f": true, ":n": true, ":pl": true,
			":starred": true, ":star": true, ":fav": true, ":favorite": true,
			":de": true, ":en": true, "lang:de": true, "lang:en": true,
		}
		for _, t := range strings.Fields(query) {
			lower := strings.ToLower(t)
			if filterTags[lower] || strings.HasPrefix(lower, "de:") || strings.HasPrefix(lower, "en:") {
				return false
			}
		}
		return true
	}
	tagLower := strings.ToLower(tag)
	for _, t := range strings.Fields(query) {
		lower := strings.ToLower(t)
		if lower == tagLower {
			return true
		}
		// Treat bare "de:" / "en:" language pills as active for any de:/en: scoped term.
		if (tagLower == "de:" && (lower == ":de" || lower == "lang:de" || strings.HasPrefix(lower, "de:"))) ||
			(tagLower == "en:" && (lower == ":en" || lower == "lang:en" || strings.HasPrefix(lower, "en:"))) {
			return true
		}
	}
	return false
}

func toggleFilterTag(query, tag string) string {
	terms := strings.Fields(query)
	found := false
	var newTerms []string
	tagLower := strings.ToLower(tag)
	langToggle := tagLower == "de:" || tagLower == "en:" || tagLower == ":de" || tagLower == ":en"

	for _, t := range terms {
		lower := strings.ToLower(t)
		if lower == tagLower {
			found = true
			continue
		}
		// Language pills are mutually exclusive — drop the other side when switching.
		if langToggle && isLangFilterTerm(t) {
			if (tagLower == "de:" || tagLower == ":de") && (lower == ":de" || lower == "lang:de" || strings.HasPrefix(lower, "de:")) {
				found = true
				continue
			}
			if (tagLower == "en:" || tagLower == ":en") && (lower == ":en" || lower == "lang:en" || strings.HasPrefix(lower, "en:")) {
				found = true
				continue
			}
			if (tagLower == "de:" || tagLower == ":de") && (lower == ":en" || lower == "lang:en" || strings.HasPrefix(lower, "en:")) {
				continue
			}
			if (tagLower == "en:" || tagLower == ":en") && (lower == ":de" || lower == "lang:de" || strings.HasPrefix(lower, "de:")) {
				continue
			}
		}
		newTerms = append(newTerms, t)
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
			deck, err := m.ensureDictionaryTargetDeck(ctx)
			if err != nil {
				return err
			}

			added, skipped, failed := 0, 0, 0
			for _, entry := range entries {
				noteID := "dict-" + entry.ID
				if entry.ID == "" {
					noteID = fmt.Sprintf("dict-%d-%d", time.Now().UnixNano(), added)
				}
				if entry.ID != "" {
					if _, err := m.repo.GetNote(ctx, noteID); err == nil {
						skipped++
						continue
					}
				}
				tags := append([]string(nil), entry.Tags...)
				tags = append(tags, "dictionary")
				note := core.Note{
					ID:     noteID,
					DeckID: deck.ID,
					Type:   "flashcard",
					Front:  formatDictionaryCardFront(entry.Word, entry.Gender),
					Back:   entry.Translation,
					Extra:  content.DictionaryEntryExtra(entry),
					Tags:   tags,
				}
				note.Cards = content.CardsForNote(note)
				if err := m.repo.UpsertNote(ctx, note); err == nil {
					added++
				} else {
					failed++
				}
			}

			status := fmt.Sprintf("Added %d entries to %s deck", added, deck.Name)
			if skipped > 0 {
				status += fmt.Sprintf("; skipped %d already added", skipped)
			}
			if failed > 0 {
				status += fmt.Sprintf("; %d failed", failed)
			}
			return statusMsg{text: status}
		},
		m.loadDecks,
		m.loadDueCards,
	)
}

func (m *Model) dictionaryDetailVisible() bool {
	if m.dictionaryDetailView {
		return true
	}
	if m.dictionaryOverlayActive {
		boxWidth := 86
		if m.width < 92 {
			boxWidth = m.width - 6
		}
		if boxWidth < 30 {
			boxWidth = 30
		}
		return boxWidth-4 > 70
	}
	return m.width > 80
}

// inspectDictionaryCursor records the selected entry when its detail pane is visible.
func (m *Model) inspectDictionaryCursor() {
	if !m.dictionaryDetailVisible() {
		return
	}
	if m.dictionaryCursor >= 0 && m.dictionaryCursor < len(m.dictionaryResults) {
		m.recordDictionaryView(m.dictionaryResults[m.dictionaryCursor])
	}
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
	m.saveDictionaryRecentlyViewed()
}

// explainDictionaryEntry asks the AI tutor for a pedagogical explanation of a
// dictionary entry and surfaces it in the AI view (same explainMsg pipeline as Review).
func (m *Model) explainDictionaryEntry(entry core.DictionaryEntry) tea.Cmd {
	provider := m.aiProvider
	if provider == nil {
		m.status = "Enable an AI provider in Settings to get explanations"
		return nil
	}
	if m.explainingCard {
		return nil
	}
	m.explainingCard = true
	m.explanation = ""
	m.explainError = ""
	m.status = fmt.Sprintf("Asking AI tutor about '%s'…", entry.Word)
	cardID := "dict:" + entry.ID
	if entry.ID == "" {
		cardID = "dict:" + entry.Word
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		explanation, err := ai.ExplainDictionaryEntry(ctx, provider, entry)
		if err != nil {
			return explainErrorMsg{cardID: cardID, err: err}
		}
		return explainMsg{cardID: cardID, explanation: explanation}
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
			deck, err := m.ensureDictionaryTargetDeck(ctx)
			if err != nil {
				return err
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

			// Check for duplicates
			if m.dictionaryLastAddAttemptNoteID != noteID {
				if _, err := m.repo.GetNote(ctx, noteID); err == nil {
					m.dictionaryLastAddAttemptNoteID = noteID
					return statusMsg{text: fmt.Sprintf("'%s' already in %s deck (use c/ctrl+c again to update)", wordClean, deck.Name)}
				}
			}
			m.dictionaryLastAddAttemptNoteID = ""

			tags := append([]string(nil), entry.Tags...)
			tags = append(tags, "dictionary", "cloze")
			note := core.Note{
				ID:     noteID,
				DeckID: deck.ID,
				Type:   "cloze",
				Front:  frontText,
				Back:   entry.Translation,
				Extra:  content.DictionaryEntryExtra(entry),
				Tags:   tags,
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

func (m *Model) exportDictionaryResultsTSVCmd() tea.Cmd {
	entries := m.dictionaryResults
	if len(entries) == 0 {
		m.status = "No dictionary entries to export"
		return nil
	}
	starred := isFilterActive(m.dictionarySearch, ":starred") || isFilterActive(m.dictionarySearch, ":star")
	deckID, _ := m.dictionaryTargetDeck()
	return func() tea.Msg {
		dir := "exports"
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create export dir: %w", err)
		}
		prefix := "dictionary_export"
		if starred {
			prefix = "dictionary_starred"
		}
		filename := fmt.Sprintf("%s_%s.tsv", prefix, time.Now().Format("20060102_150405"))
		exportPath := filepath.Join(dir, filename)

		notes := content.DictionaryEntriesToNotes(entries, deckID)
		file, err := os.Create(exportPath)
		if err != nil {
			return fmt.Errorf("create export file: %w", err)
		}
		defer file.Close()

		if err := content.ExportAnkiTSV(file, notes); err != nil {
			return fmt.Errorf("export TSV: %w", err)
		}

		return statusMsg{text: fmt.Sprintf("Exported %d dictionary entries to %s", len(notes), exportPath)}
	}
}
