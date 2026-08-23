package tui

import (
	"sync"
	"sync/atomic"

	tea "charm.land/bubbletea/v2"
)

// orderedSave serializes writes for one persisted value and suppresses a
// queued snapshot when a newer request exists before that snapshot starts.
type orderedSave struct {
	mu     sync.Mutex
	latest atomic.Uint64
}

func (s *orderedSave) command(save func() error) tea.Cmd {
	requestID := s.latest.Add(1)
	return func() tea.Msg {
		s.mu.Lock()
		defer s.mu.Unlock()
		if requestID != s.latest.Load() {
			return nil
		}
		if err := save(); err != nil {
			return err
		}
		return nil
	}
}
