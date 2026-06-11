package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

		results, err := dictRepo.Search(context.Background(), query, 50)
		if err != nil {
			return err
		}
		return dictionarySearchResultsMsg{id: id, results: results}
	}
}

func (m *Model) addDictionaryEntryCmd(entry core.DictionaryEntry) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// 1. Ensure "Dictionary" deck exists
		deckID := "dictionary"
		deck, err := m.repo.GetDeck(ctx, deckID)
		if err != nil {
			deck = core.Deck{
				ID:   deckID,
				Name: "Dictionary",
			}
			if err := m.repo.UpsertDeck(ctx, deck); err != nil {
				return err
			}
		}

		// 2. Create note
		note := core.Note{
			ID:     "dict-" + entry.ID,
			DeckID: deckID,
			Type:   "flashcard",
			Front:  entry.Word,
			Back:   entry.Translation,
			Extra:  entry.Forms,
			Tags:   append(entry.Tags, "dictionary"),
		}
		if entry.WordClass != "" {
			note.Extra += "\n[" + entry.WordClass + "]"
		}
		if entry.Gender != "" {
			note.Extra += " {" + entry.Gender + "}"
		}
		note.Cards = content.CardsForNote(note)

		if err := m.repo.UpsertNote(ctx, note); err != nil {
			return err
		}

		return statusMsg{text: fmt.Sprintf("Added '%s' to Dictionary deck", entry.Word)}
	}
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

func (m *Model) openDictionaryOverlay() {
	m.dictionaryOverlayActive = true
	m.resetDictionarySearchState()
	m.status = "Spotlight dictionary open"
}

func (m *Model) closeDictionaryOverlay() {
	m.dictionaryOverlayActive = false
	m.resetDictionarySearchState()
	m.status = "Spotlight dictionary closed"
}
