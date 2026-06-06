package tui

import (
	tea "charm.land/bubbletea/v2"
	"context"
	"deutsch-tui/internal/core"
	"fmt"
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

		if err := m.repo.UpsertNote(ctx, note); err != nil {
			return err
		}

		return statusMsg{text: fmt.Sprintf("Added '%s' to Dictionary deck", entry.Word)}
	}
}
