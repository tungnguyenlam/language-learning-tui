package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"
)

func (m *Model) recordDeckSearch(query string) tea.Cmd {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	for i, q := range m.deckSearchHistory {
		if q == query {
			m.deckSearchHistory = append(m.deckSearchHistory[:i], m.deckSearchHistory[i+1:]...)
			m.deckSearchHistory = append(m.deckSearchHistory, query)
			return m.saveDeckHistory()
		}
	}
	m.deckSearchHistory = append(m.deckSearchHistory, query)
	if len(m.deckSearchHistory) > 5 {
		m.deckSearchHistory = m.deckSearchHistory[1:]
	}
	return m.saveDeckHistory()
}

func (m *Model) saveDeckHistory() tea.Cmd {
	history := append([]string(nil), m.deckSearchHistory...)
	requestID := m.deckHistorySaveID.Add(1)
	return func() tea.Msg {
		m.deckHistorySaveMu.Lock()
		defer m.deckHistorySaveMu.Unlock()
		if requestID != m.deckHistorySaveID.Load() {
			return nil
		}
		data, err := json.Marshal(history)
		if err != nil {
			return fmt.Errorf("encode deck search history: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := m.repo.SetSetting(ctx, "deck_search_history", string(data)); err != nil {
			return fmt.Errorf("save deck search history: %w", err)
		}
		return nil
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

func (m *Model) exportDeckTSVCmd(deckID string) tea.Cmd {
	if strings.TrimSpace(deckID) == "" {
		m.status = "No deck selected for export"
		return nil
	}
	dir := "exports"
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.status = fmt.Sprintf("Failed to create export dir: %v", err)
		return nil
	}
	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, deckID)
	filename := fmt.Sprintf("deck_%s_%s.tsv", safeName, time.Now().Format("20060102_150405"))
	exportPath := filepath.Join(dir, filename)
	m.statusSeq++
	m.status = fmt.Sprintf("Exporting %s...", deckID)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cards, err := m.repo.Cards(ctx, deckID, "", "")
		if err != nil {
			return fmt.Errorf("failed to load cards for deck %s: %w", deckID, err)
		}
		if len(cards) == 0 {
			return statusMsg{text: fmt.Sprintf("Deck %q has no cards to export", deckID)}
		}
		seen := make(map[string]bool)
		var notes []core.Note
		for _, c := range cards {
			if seen[c.NoteID] {
				continue
			}
			seen[c.NoteID] = true
			notes = append(notes, core.Note{
				ID:     c.NoteID,
				DeckID: c.DeckID,
				Front:  c.Prompt,
				Back:   c.Answer,
				Extra:  c.Extra,
				Hint:   c.Hint,
				Tags:   c.Tags,
				Audio:  c.Audio,
			})
		}
		file, err := os.Create(exportPath)
		if err != nil {
			return fmt.Errorf("failed to create export file '%s': %w", filepath.Base(exportPath), err)
		}
		defer file.Close()

		if err := content.ExportAnkiTSV(file, notes); err != nil {
			return fmt.Errorf("failed to write TSV data for deck %s: %w", deckID, err)
		}
		return statusMsg{text: fmt.Sprintf("Exported %d notes from deck %q to %s", len(notes), deckID, exportPath)}
	}
}
