package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// trainerModel returns a Model sitting in the named trainer with n synthetic
// items loaded and ready for input.
func trainerModel(t *testing.T, kind PracticeSubView, n int) (*Model, *trainerState) {
	t.Helper()
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = kind

	st := m.trainerStateFor(kind)
	st.items = make([]trainerItem, n)
	for i := range st.items {
		st.items[i] = trainerItem{
			Title:  fmt.Sprintf("item-%d", i),
			Answer: fmt.Sprintf("answer-%d", i),
		}
	}
	return m, st
}

// "q" is a global quit shortcut, but German answers contain it (Qualität,
// Quelle). A trainer waiting for input must receive the character instead.
func TestTrainerConsumesGlobalShortcutCharacters(t *testing.T) {
	for _, key := range []rune{'q', '?'} {
		t.Run(string(key), func(t *testing.T) {
			m, st := trainerModel(t, PracticeSubViewPlural, 3)

			updated, cmd := m.Update(tea.KeyPressMsg{Code: key})
			m = updated.(*Model)

			if cmd != nil {
				t.Fatalf("expected %q to produce no command while typing, got one", key)
			}
			if st.input != string(key) {
				t.Fatalf("expected %q to be typed into the answer, input=%q", key, st.input)
			}
			if m.showHelp {
				t.Fatal("expected help overlay to stay closed while typing")
			}
		})
	}
}

// Outside the trainers the same keys must keep their global meaning.
func TestGlobalShortcutsStillWorkOutsideTrainers(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDashboard

	updated, _ := m.Update(tea.KeyPressMsg{Code: '?'})
	if !updated.(*Model).showHelp {
		t.Fatal("expected ? to toggle the help overlay on the Dashboard")
	}
	// Dismiss the modal before checking the Dashboard quit shortcut.
	m.Update(tea.KeyPressMsg{Code: '?'})

	if _, cmd := m.Update(tea.KeyPressMsg{Code: 'q'}); cmd == nil {
		t.Fatal("expected q to quit from the Dashboard")
	}
}

// A revealed trainer still treats any key as "advance", including q and ?.
func TestTrainerRevealedStateAdvancesOnGlobalShortcuts(t *testing.T) {
	for _, key := range []rune{'q', '?'} {
		t.Run(string(key), func(t *testing.T) {
			m, st := trainerModel(t, PracticeSubViewPlural, 3)
			st.revealed = true
			st.index = 0

			if m.trainerInputActive() {
				t.Fatal("expected a revealed trainer not to be treated as accepting typed input")
			}
			if !m.practiceBlocksGlobalShortcut() {
				t.Fatal("expected revealed trainer to block global shortcuts so advance wins")
			}

			updated, cmd := m.Update(tea.KeyPressMsg{Code: key})
			m = updated.(*Model)

			if cmd != nil {
				t.Fatalf("expected %q to advance (no quit/help cmd), got a command", key)
			}
			if st.revealed {
				t.Fatal("expected reveal to clear after advance")
			}
			if st.index != 1 {
				t.Fatalf("expected advance to next item, index=%d", st.index)
			}
			if m.showHelp {
				t.Fatal("expected help overlay to stay closed on revealed advance")
			}
		})
	}
}

func TestTrainerAdvanceReshufflesOnWrap(t *testing.T) {
	_, st := trainerModel(t, PracticeSubViewCase, 12)

	first := st.items[0].Answer
	for i := 0; i < len(st.items)-1; i++ {
		st.advance()
	}
	if st.index != len(st.items)-1 {
		t.Fatalf("expected index at last item, got %d", st.index)
	}
	if st.round != 0 {
		t.Fatalf("expected still on round 0, got %d", st.round)
	}
	if st.items[0].Answer != first {
		t.Fatal("expected the first pass to keep its authored order")
	}

	st.advance()
	if st.index != 0 {
		t.Fatalf("expected wrap to index 0, got %d", st.index)
	}
	if st.round != 1 {
		t.Fatalf("expected round to increment, got %d", st.round)
	}

	// The shuffle must permute, never drop or duplicate, the exercise set.
	seen := make(map[string]int, len(st.items))
	for _, it := range st.items {
		seen[it.Answer]++
	}
	for i := 0; i < 12; i++ {
		if seen[fmt.Sprintf("answer-%d", i)] != 1 {
			t.Fatalf("shuffle did not preserve the exercise set: %v", seen)
		}
	}
}

func TestTrainerAdvanceClearsAnswerState(t *testing.T) {
	_, st := trainerModel(t, PracticeSubViewCase, 3)
	st.revealed = true
	st.showHint = true
	st.input = "dem"

	st.advance()

	if st.revealed || st.showHint || st.input != "" {
		t.Fatalf("expected advance to clear answer state, got revealed=%v hint=%v input=%q",
			st.revealed, st.showHint, st.input)
	}
}

// The fill-in-the-blank sets are hand-authored; this guards the shape every
// item must have so a malformed entry cannot ship a blankless or unanswerable
// exercise.
func TestBlankTrainerContentIsWellFormed(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})

	for _, tc := range []struct {
		name string
		kind PracticeSubView
		load func(*Model) tea.Cmd
		min  int
	}{
		{"case", PracticeSubViewCase, (*Model).loadCaseItems, 25},
		{"adjective", PracticeSubViewAdjective, (*Model).loadAdjectiveItems, 25},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg, ok := tc.load(m)().(trainerItemsMsg)
			if !ok {
				t.Fatal("expected trainerItemsMsg")
			}
			if msg.kind != tc.kind {
				t.Fatalf("expected kind %v, got %v", tc.kind, msg.kind)
			}
			if len(msg.items) < tc.min {
				t.Fatalf("expected at least %d exercises, got %d", tc.min, len(msg.items))
			}

			seen := make(map[string]bool, len(msg.items))
			for _, it := range msg.items {
				if !strings.Contains(it.Title, "_____") {
					t.Errorf("prompt has no blank: %q", it.Title)
				}
				if it.Answer == "" {
					t.Errorf("exercise has no answer: %q", it.Title)
				}
				if it.Context == "" {
					t.Errorf("exercise has no grammar context: %q", it.Title)
				}
				// The reveal must show the answer substituted into the sentence.
				if !strings.Contains(it.RevealTitle, it.Answer) {
					t.Errorf("reveal %q does not contain answer %q", it.RevealTitle, it.Answer)
				}
				if seen[it.Title] {
					t.Errorf("duplicate exercise: %q", it.Title)
				}
				seen[it.Title] = true
			}
		})
	}
}

func TestTrainerRendersItemPosition(t *testing.T) {
	m, st := trainerModel(t, PracticeSubViewCase, 15)
	layout := viewportLayout{Width: 80, Height: 30}

	view := m.renderTrainer(PracticeSubViewCase, layout)
	if !strings.Contains(view, "Item 1/15") {
		t.Fatalf("expected item position in header, got: %s", view)
	}
	// E2E trainer tests match on this exact substring.
	if !strings.Contains(view, "Score: 0/0") {
		t.Fatalf("expected contiguous score text, got: %s", view)
	}
	if strings.Contains(view, "Round") {
		t.Fatal("expected no round counter on the first pass")
	}

	st.round = 1
	view = m.renderTrainer(PracticeSubViewCase, layout)
	if !strings.Contains(view, "Round 2") {
		t.Fatalf("expected round counter after a completed pass, got: %s", view)
	}
}

func TestPracticeHubSearchFilter(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewHub
	m.width = 100
	m.height = 30

	// 1. Press '/' to open filter
	m.updatePracticeKey(tea.KeyPressMsg{Text: "/"})
	if !m.practiceFilterFocus {
		t.Fatalf("expected practiceFilterFocus to be true after '/' key")
	}

	// 2. Type "passive"
	for _, char := range "passive" {
		m.updatePracticeKey(tea.KeyPressMsg{Text: string(char)})
	}
	if m.practiceFilter != "passive" {
		t.Fatalf("expected practiceFilter to be 'passive', got %q", m.practiceFilter)
	}

	// 3. Render Practice Hub and assert filtered title is visible
	out := m.renderPracticeHub(m.activeViewContentLayout())
	if !strings.Contains(out, "Filter: passive") {
		t.Fatalf("expected rendered output to contain 'Filter: passive', got: %s", out)
	}
	if !strings.Contains(out, "Passive Voice Trainer") {
		t.Fatalf("expected rendered output to contain 'Passive Voice Trainer'")
	}
	if strings.Contains(out, "Gender Trainer") {
		t.Fatalf("expected rendered output to NOT contain 'Gender Trainer'")
	}

	// 4. Press Esc to clear filter
	m.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.practiceFilter != "" {
		t.Fatalf("expected practiceFilter to be cleared on Esc, got %q", m.practiceFilter)
	}
}
