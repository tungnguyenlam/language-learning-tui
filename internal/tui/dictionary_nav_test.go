package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestDictionaryNavigation(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40

	// 1. Full Dictionary tab navigation when empty
	m.activeView = ViewDictionary
	m.dictionarySearch = ""
	m.dictionaryFocusResults = false

	// Pressing '1' should switch to Dashboard
	msg := tea.KeyPressMsg{Code: '1'}
	_, cmd := m.Update(msg)
	if m.activeView != ViewDashboard {
		t.Errorf("expected activeView to be Dashboard, got %s", m.activeView)
	}
	if cmd == nil {
		t.Errorf("expected a command for view update")
	}

	// 2. Dictionary Overlay navigation when empty
	m = NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40
	m.activeView = ViewReview
	m.dictionaryOverlayActive = true
	m.dictionarySearch = ""

	// Pressing '1' should switch to Dashboard AND close overlay
	msg = tea.KeyPressMsg{Code: '1'}
	_, cmd = m.Update(msg)
	if m.activeView != ViewDashboard {
		t.Errorf("expected activeView to be Dashboard from overlay, got %s", m.activeView)
	}
	if m.dictionaryOverlayActive {
		t.Errorf("expected dictionary overlay to be closed")
	}

	// 3. Dictionary search should still work for numbers if not empty
	m = NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 100
	m.height = 40
	m.activeView = ViewDictionary
	m.dictionarySearch = "A"
	m.dictionaryFocusResults = false

	// Pressing '1' should add to search, NOT switch view
	msg = tea.KeyPressMsg{Code: '1'}
	m.Update(msg)
	if m.activeView != ViewDictionary {
		t.Errorf("expected activeView to stay Dictionary, got %s", m.activeView)
	}
	if m.dictionarySearch != "A1" {
		t.Errorf("expected dictionarySearch to be A1, got %s", m.dictionarySearch)
	}
}
