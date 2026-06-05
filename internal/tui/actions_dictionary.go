package tui

import (
	tea "charm.land/bubbletea/v2"
	"context"
	"deutsch-tui/internal/core"
	"fmt"
)

type dictionarySearchResultsMsg struct {
	results []core.DictionaryEntry
}

func (m *Model) searchDictionary() tea.Cmd {
	query := m.dictionarySearch
	return func() tea.Msg {
		if query == "" {
			return dictionarySearchResultsMsg{results: nil}
		}

		dictRepo, ok := m.repo.(core.DictionaryRepository)
		if !ok {
			return fmt.Errorf("repository does not support dictionary search")
		}

		results, err := dictRepo.Search(context.Background(), query, 50)
		if err != nil {
			return err
		}
		return dictionarySearchResultsMsg{results: results}
	}
}
