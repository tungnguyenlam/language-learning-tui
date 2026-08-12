package tui

import "testing"

func TestDebounceSearchTracksIndependentViewGenerations(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})

	if cmd := m.debounceSearch(ViewBrowser); cmd == nil {
		t.Fatal("browser debounce command is nil")
	}
	if cmd := m.debounceSearch(ViewBrowser); cmd == nil {
		t.Fatal("second browser debounce command is nil")
	}
	if cmd := m.debounceSearch(ViewDictionary); cmd == nil {
		t.Fatal("dictionary debounce command is nil")
	}

	if m.browserSearchTimerID != 2 {
		t.Fatalf("browser generation = %d, want 2", m.browserSearchTimerID)
	}
	if m.dictionarySearchTimerID != 1 {
		t.Fatalf("dictionary generation = %d, want 1", m.dictionarySearchTimerID)
	}
	if cmd := m.debounceSearch(ViewDashboard); cmd != nil {
		t.Fatal("unsupported view should not schedule a search")
	}
}

func TestDebouncedSearchDropsStaleGeneration(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewBrowser
	m.debounceSearch(ViewBrowser)
	m.debounceSearch(ViewBrowser)

	_, staleCmd := m.Update(debounceSearchMsg{id: 1, view: ViewBrowser})
	if staleCmd != nil {
		t.Fatal("stale browser generation scheduled a query")
	}
	_, currentCmd := m.Update(debounceSearchMsg{id: 2, view: ViewBrowser})
	if currentCmd == nil {
		t.Fatal("current browser generation did not schedule a query")
	}
}
