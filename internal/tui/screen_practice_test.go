package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestPracticeHubFilterAcceptsSpace(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub
	model.practiceFilterFocus = true

	_, handled := (practiceScreen{}).HandleKey(model, tea.KeyPressMsg{Text: "space"})
	if !handled {
		t.Fatal("space was not handled in the Practice Hub filter")
	}
	if model.practiceFilter != " " {
		t.Fatalf("practice filter = %q, want a single space", model.practiceFilter)
	}
}

func TestPracticeHubFilterEnterStartsVisibleTrainer(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub
	model.practiceFilter = "passive"
	model.practiceHubCursor = 0

	_, handled := (practiceScreen{}).HandleKey(model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled {
		t.Fatal("enter was not handled on the filtered Practice Hub")
	}
	if model.practiceSubView != PracticeSubViewPassive {
		t.Fatalf("enter started %v, want Passive", model.practiceSubView)
	}
}

func TestPracticeHubFilterWrapsVisibleList(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub
	model.practiceFilter = "ending"
	model.practiceHubCursor = 0

	visible := model.visiblePracticeHubModes()
	if len(visible) != 2 {
		t.Fatalf("visible trainers for %q = %d, want 2", model.practiceFilter, len(visible))
	}

	_, handled := (practiceScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'j'})
	if !handled {
		t.Fatal("j was not handled on the filtered Practice Hub")
	}
	if model.practiceHubCursor != 1 {
		t.Fatalf("cursor after j = %d, want 1", model.practiceHubCursor)
	}

	(practiceScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'j'})
	if model.practiceHubCursor != 0 {
		t.Fatalf("cursor after wrap-down = %d, want 0", model.practiceHubCursor)
	}

	(practiceScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'k'})
	if model.practiceHubCursor != 1 {
		t.Fatalf("cursor after wrap-up = %d, want 1", model.practiceHubCursor)
	}
}
