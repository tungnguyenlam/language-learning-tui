package tui

import (
	"context"
	"encoding/json"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) recordDeckSearch(query string) {
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}
	for i, q := range m.deckSearchHistory {
		if q == query {
			m.deckSearchHistory = append(m.deckSearchHistory[:i], m.deckSearchHistory[i+1:]...)
			m.deckSearchHistory = append(m.deckSearchHistory, query)
			m.saveDeckHistory()
			return
		}
	}
	m.deckSearchHistory = append(m.deckSearchHistory, query)
	if len(m.deckSearchHistory) > 5 {
		m.deckSearchHistory = m.deckSearchHistory[1:]
	}
	m.saveDeckHistory()
}

func (m *Model) saveDeckHistory() {
	ctx := context.Background()
	data, err := json.Marshal(m.deckSearchHistory)
	if err == nil {
		_ = m.repo.SetSetting(ctx, "deck_search_history", string(data))
	}
}

func (m *Model) loadDeckHistory() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		raw, err := m.repo.GetSetting(ctx, "deck_search_history")
		if err != nil || raw == "" {
			return deckHistoryLoadedMsg{}
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
		return deckHistoryLoadedMsg(clean)
	}
}
