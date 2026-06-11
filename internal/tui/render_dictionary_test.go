package tui

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"deutsch-tui/internal/core"
)

func TestHighlightQuery(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)

	tests := []struct {
		text     string
		query    string
		expected string
	}{
		{
			text:     "Deutsch",
			query:    "de",
			expected: style.Render("De") + "utsch",
		},
		{
			text:     "Apfelkuchen",
			query:    "kuch",
			expected: "Apfel" + style.Render("kuch") + "en",
		},
		{
			text:     "banana",
			query:    "an",
			expected: "b" + style.Render("an") + style.Render("an") + "a",
		},
		{
			text:     "banana",
			query:    "",
			expected: "banana",
		},
		{
			text:     "Äpfel und Apfel",
			query:    "äpf",
			expected: style.Render("Äpf") + "el und Apfel",
		},
		{
			text:     "Käsekuchen",
			query:    "KÄSE",
			expected: style.Render("Käse") + "kuchen",
		},
		{
			text:     "Straße",
			query:    "ß",
			expected: "Stra" + style.Render("ß") + "e",
		},
	}

	for _, tc := range tests {
		got := highlightQuery(tc.text, tc.query, style)
		if got != tc.expected {
			t.Errorf("highlightQuery(%q, %q) = %q, expected %q", tc.text, tc.query, got, tc.expected)
		}
	}
}

func TestRenderDictionary(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryResults = []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Apfel",
			Translation: "Apple",
			WordClass:   "Noun",
			Gender:      "r",
			Forms:       "Apfel, Äpfel",
			Examples:    []string{"Ein roter Apfel."},
		},
		{
			ID:          "2",
			Word:        "Birne",
			Translation: "Pear",
		},
	}
	m.dictionarySearch = "Ap"
	m.dictionaryCursor = 0

	// Test wide layout (two columns)
	m.width = 100
	m.height = 40
	layout := m.activeViewContentLayout()
	view := m.renderDictionary(layout)
	plainView := stripANSI(view)

	if !strings.Contains(plainView, "Apfel") {
		t.Errorf("expected view to contain Apfel")
	}
	if !strings.Contains(plainView, "Apple") {
		t.Errorf("expected view to contain Apple")
	}
	if !strings.Contains(plainView, "Forms:") {
		t.Errorf("expected view to contain detail panel with Forms. Plain view:\n%s", plainView)
	}
	if !strings.Contains(plainView, "Äpfel") {
		t.Errorf("expected view to contain Äpfel")
	}

	// Test compact layout (single column)
	m.width = 60
	layout = m.activeViewContentLayout()
	view = m.renderDictionary(layout)
	plainView = stripANSI(view)

	if strings.Contains(plainView, "Forms:") {
		t.Errorf("expected compact view NOT to contain detail panel")
	}
}

func TestDictionaryDetailScrollClampsToVisibleRows(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.breakpoint = BreakpointWide
	m.dictionaryDetailTotalLines = 40

	maxScroll := maxInt(0, m.dictionaryDetailTotalLines-dictionaryVisibleRows(m.activeViewContentLayout()))
	for i := 0; i < 100; i++ {
		cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModShift})
		if cmd != nil {
			t.Fatal("shift+down should not return a command")
		}
		if !handled {
			t.Fatal("shift+down should be handled")
		}
	}

	if m.dictionaryDetailScroll != maxScroll {
		t.Fatalf("dictionaryDetailScroll = %d, want %d", m.dictionaryDetailScroll, maxScroll)
	}
}

func TestDictionarySearchHistory(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.recordDictionarySearch("Apfel")
	m.recordDictionarySearch("Birne")
	m.recordDictionarySearch("Apfel") // Duplicate, should be moved to the end (most recent)

	if len(m.dictionarySearchHistory) != 2 {
		t.Fatalf("expected search history to have 2 items, got %d", len(m.dictionarySearchHistory))
	}
	if m.dictionarySearchHistory[0] != "Birne" || m.dictionarySearchHistory[1] != "Apfel" {
		t.Errorf("unexpected history order: %v", m.dictionarySearchHistory)
	}

	// Navigate away to test recording search on view transition
	m.activeView = ViewDictionary
	m.dictionarySearch = "Banane"
	m.updateView(ViewDashboard)

	if len(m.dictionarySearchHistory) != 3 {
		t.Fatalf("expected search history to have 3 items, got %d", len(m.dictionarySearchHistory))
	}
	if m.dictionarySearchHistory[2] != "Banane" {
		t.Errorf("expected Banane at end of search history, got %v", m.dictionarySearchHistory)
	}
}

type captureRepo struct {
	mockRepo
	upsertedNote *core.Note
}

func (r *captureRepo) UpsertNote(ctx context.Context, note core.Note) error {
	r.upsertedNote = &note
	return nil
}

func (r *captureRepo) UpsertDeck(ctx context.Context, deck core.Deck) error {
	return nil
}

func (r *captureRepo) GetDeck(ctx context.Context, id string) (core.Deck, error) {
	return core.Deck{ID: id, Name: "Dictionary"}, nil
}

func TestDictionaryQuickAddCardsGenerated(t *testing.T) {
	repo := &captureRepo{}
	m := NewModel(repo, &mockScheduler{})

	entry := core.DictionaryEntry{
		ID:          "test-1",
		Word:        "Auto",
		Translation: "Car",
		Gender:      "n",
	}

	cmd := m.addDictionaryEntryCmd(entry)
	if cmd == nil {
		t.Fatal("expected command from addDictionaryEntryCmd")
	}

	msg := cmd()
	status, ok := msg.(statusMsg)
	if !ok {
		t.Fatalf("expected statusMsg, got %T: %v", msg, msg)
	}

	if !strings.Contains(status.text, "Added") {
		t.Errorf("expected status text to say Added, got %q", status.text)
	}

	if repo.upsertedNote == nil {
		t.Fatal("expected note to be upserted")
	}

	note := repo.upsertedNote
	if note.Front != "Auto" || note.Back != "Car" {
		t.Errorf("unexpected note content: front=%q, back=%q", note.Front, note.Back)
	}

	if len(note.Cards) == 0 {
		t.Fatal("expected cards to be generated for the quick-added dictionary entry")
	}

	card := note.Cards[0]
	if card.Prompt != "Auto" || card.Answer != "Car" {
		t.Errorf("unexpected card content: prompt=%q, answer=%q", card.Prompt, card.Answer)
	}
}

func TestDictionarySearchClearHitbox(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.dictionarySearch = "Apfel"
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Apfel", Translation: "Apple"},
	}

	layout := m.activeViewContentLayout()
	_ = m.renderDictionary(layout)

	// Check if clear button hitbox is registered
	var foundClearHitbox bool
	for _, hb := range m.hitboxes {
		if hb.ID == "dict-search-clear" {
			foundClearHitbox = true
			if hb.Action == nil {
				t.Fatal("expected clear hitbox action to be set")
			}
			hb.Action() // Execute clear
			break
		}
	}

	if !foundClearHitbox {
		t.Fatal("expected dictionary clear button hitbox to be registered")
	}

	if m.dictionarySearch != "" {
		t.Errorf("expected search query to be cleared, got %q", m.dictionarySearch)
	}
	if len(m.dictionaryResults) != 0 {
		t.Errorf("expected results to be cleared, got %d items", len(m.dictionaryResults))
	}
}

func TestDictionaryPreviousView(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryProvider = "Local TUI"
	m.activeView = ViewBrowser

	// Open dictionary via command
	m.openDictionary("Apfel")
	if m.activeView != ViewDictionary {
		t.Fatalf("expected active view to be ViewDictionary, got %v", m.activeView)
	}
	if m.dictionaryPreviousView != ViewBrowser {
		t.Fatalf("expected dictionaryPreviousView to be ViewBrowser, got %v", m.dictionaryPreviousView)
	}

	// Press Esc to go back
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyEsc})
	if !handled {
		t.Fatal("expected Esc to be handled")
	}
	if cmd != nil {
		// updateView returns a Cmd if view requires loading, which ViewBrowser does
		_ = cmd()
	}

	if m.activeView != ViewBrowser {
		t.Fatalf("expected active view to return to ViewBrowser, got %v", m.activeView)
	}
}

func TestDictionaryHeaderResultsCount(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40

	// Case 1: no results
	m.dictionaryResults = nil
	view := m.renderDictionary(m.activeViewContentLayout())
	if !strings.Contains(stripANSI(view), "Dictionary") || strings.Contains(stripANSI(view), "results") {
		t.Errorf("unexpected header for no results: %q", view)
	}

	// Case 2: few results
	m.dictionaryResults = make([]core.DictionaryEntry, 5)
	view = m.renderDictionary(m.activeViewContentLayout())
	if !strings.Contains(stripANSI(view), "Dictionary (5 results)") {
		t.Errorf("expected header to contain results count, got: %q", view)
	}

	// Case 3: 50+ results
	m.dictionaryResults = make([]core.DictionaryEntry, 50)
	view = m.renderDictionary(m.activeViewContentLayout())
	if !strings.Contains(stripANSI(view), "Dictionary (50+ results)") {
		t.Errorf("expected header to contain 50+ results, got: %q", view)
	}
}

func TestDictionarySingleColumnDetailView(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryResults = []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Auto",
			Translation: "Car",
			WordClass:   "Noun",
			Gender:      "n",
			Forms:       "Autos",
			Examples:    []string{"Ein schnelles Auto."},
		},
	}
	m.dictionaryCursor = 0

	// compact layout (width <= 80)
	m.width = 60
	m.height = 40
	m.dictionaryDetailView = false

	// Verify that details are not shown initially
	view := m.renderDictionary(m.activeViewContentLayout())
	plainView := stripANSI(view)
	if strings.Contains(plainView, "Forms:") {
		t.Error("expected compact view not to contain detail panel details initially")
	}

	// Toggle details view using key or directly
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	if !handled || cmd != nil {
		t.Fatal("expected ctrl+d to be handled without cmd")
	}

	if !m.dictionaryDetailView {
		t.Error("expected dictionaryDetailView to be true after ctrl+d")
	}

	// Verify detail panel is now displayed in the single column
	view = m.renderDictionary(m.activeViewContentLayout())
	plainView = stripANSI(view)
	if !strings.Contains(plainView, "Word Forms:") || !strings.Contains(plainView, "Autos") {
		t.Errorf("expected view to display detail view contents, got:\n%s", plainView)
	}

	// Exit detail view with esc
	cmd, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyEsc})
	if !handled || cmd != nil {
		t.Fatal("expected esc to be handled without cmd")
	}
	if m.dictionaryDetailView {
		t.Error("expected dictionaryDetailView to be false after esc")
	}
}

func TestDictionaryClearSearchHistory(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.dictionarySearchHistory = []string{"Apfel", "Birne"}

	layout := m.activeViewContentLayout()
	_ = m.renderDictionary(layout)

	// Check if clear button hitbox is registered
	var foundClearHitbox bool
	for _, hb := range m.hitboxes {
		if hb.ID == "dict-history-clear" {
			foundClearHitbox = true
			if hb.Action == nil {
				t.Fatal("expected clear history hitbox action to be set")
			}
			hb.Action() // Execute clear
			break
		}
	}

	if !foundClearHitbox {
		t.Fatal("expected dictionary clear history hitbox to be registered")
	}

	if len(m.dictionarySearchHistory) != 0 {
		t.Errorf("expected search history to be cleared, got %v", m.dictionarySearchHistory)
	}

	// Re-populate and test via ctrl+x keypress
	m.dictionarySearchHistory = []string{"Banane"}
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl})
	if !handled || cmd != nil {
		t.Fatal("expected ctrl+x to be handled in updateDictionaryKey")
	}

	if len(m.dictionarySearchHistory) != 0 {
		t.Errorf("expected search history to be cleared via ctrl+x, got %v", m.dictionarySearchHistory)
	}
}

func TestSpotlightDictionaryOverlayActivation(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewDashboard

	// Verify overlay is inactive by default
	if m.dictionaryOverlayActive {
		t.Fatal("expected overlay to be inactive by default")
	}

	// Press '=' to activate overlay
	_, cmd := m.updateKey(tea.KeyPressMsg{Code: '='})
	_ = cmd

	if !m.dictionaryOverlayActive {
		t.Fatal("expected overlay to be active after pressing '='")
	}
	if m.dictionarySearch != "" {
		t.Fatal("expected search to be empty on overlay open")
	}
	if m.activeView != ViewDashboard {
		t.Fatalf("expected active view to remain Dashboard, got %v", m.activeView)
	}
}

func TestSpotlightDictionaryOverlayRendering(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewDashboard
	m.dictionaryOverlayActive = true

	output := m.renderSpotlightDictionary()
	if !strings.Contains(output, "SPOTLIGHT DICTIONARY") {
		t.Fatal("expected spotlight overlay to contain title")
	}
	if !strings.Contains(output, "Search German or English") {
		t.Fatal("expected spotlight overlay to contain placeholder text")
	}
}

func TestSpotlightDictionaryOverlayDeactivation(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewReview
	m.dictionaryOverlayActive = true

	// Press Esc to deactivate overlay
	cmd, handled := m.updateDictionaryOverlayKey(tea.KeyPressMsg{Code: tea.KeyEsc})
	if !handled {
		t.Fatal("expected Esc to be handled")
	}
	_ = cmd
	if m.dictionaryOverlayActive {
		t.Fatal("expected overlay to be deactivated after Esc")
	}
	if m.activeView != ViewReview {
		t.Fatalf("expected active view to remain ViewReview, got %v", m.activeView)
	}
}

func TestSpotlightDictionaryOverlaySearchInput(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewStatistics
	m.dictionaryOverlayActive = true

	// Type a character — the overlay delegates to updateDictionaryKey which handles text input
	// Since dictionaryOverlayActive is true, textInputActive() returns true, blocking number-key navigation
	if !m.textInputActive() {
		t.Fatal("expected textInputActive() to return true when overlay is active")
	}
}

func TestSpotlightDictionaryOverlayToggleOff(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewBrowser
	m.dictionaryOverlayActive = true

	// Press '=' again to toggle off
	cmd, handled := m.updateDictionaryOverlayKey(tea.KeyPressMsg{Code: '='})
	if !handled {
		t.Fatal("expected '=' to be handled in overlay mode")
	}
	_ = cmd
	if m.dictionaryOverlayActive {
		t.Fatal("expected overlay to be toggled off by '='")
	}
}

func TestSpotlightDictionaryNavHitboxResetsState(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewReview
	m.dictionaryOverlayActive = true
	m.dictionarySearch = "Apfel"
	m.dictionaryResults = []core.DictionaryEntry{{ID: "1", Word: "Apfel", Translation: "Apple"}}
	m.dictionaryCursor = 1
	m.dictionaryScroll = 1
	m.dictionaryDetailScroll = 2
	m.dictionaryDetailTotalLines = 10
	m.dictionaryDetailView = true

	cmd := m.activateHitboxByID("tab-dictionary")
	if cmd != nil {
		t.Fatal("dictionary tab overlay toggle should not return a command")
	}
	if m.dictionaryOverlayActive {
		t.Fatal("expected dictionary overlay to close")
	}
	if m.activeView != ViewReview {
		t.Fatalf("active view = %s, want review", m.activeView)
	}
	if m.dictionarySearch != "" || len(m.dictionaryResults) != 0 || m.dictionaryDetailView {
		t.Fatalf("expected dictionary state to reset, search=%q results=%d detail=%v", m.dictionarySearch, len(m.dictionaryResults), m.dictionaryDetailView)
	}
	if m.dictionaryCursor != 0 || m.dictionaryScroll != 0 || m.dictionaryDetailScroll != 0 || m.dictionaryDetailTotalLines != 0 {
		t.Fatalf("expected dictionary indexes to reset, cursor=%d scroll=%d detailScroll=%d detailLines=%d", m.dictionaryCursor, m.dictionaryScroll, m.dictionaryDetailScroll, m.dictionaryDetailTotalLines)
	}

	cmd = m.activateHitboxByID("tab-dictionary")
	if cmd != nil {
		t.Fatal("dictionary tab overlay toggle should not return a command")
	}
	if !m.dictionaryOverlayActive {
		t.Fatal("expected dictionary overlay to reopen")
	}
	if m.activeView != ViewReview {
		t.Fatalf("active view = %s, want review", m.activeView)
	}
}

func TestSpotlightDictionaryClearHitboxIsOverlayScoped(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewStatistics
	m.width = 120
	m.height = 40
	m.dictionaryOverlayActive = true
	m.dictionarySearch = "Apfel"
	m.dictionaryResults = []core.DictionaryEntry{{ID: "1", Word: "Apfel", Translation: "Apple"}}
	m.dictionaryCursor = 1
	m.dictionaryScroll = 1
	m.dictionaryDetailScroll = 1
	m.dictionaryDetailTotalLines = 5
	m.dictionaryDetailView = true

	_ = m.renderSpotlightDictionary()

	var clearHitbox *Hitbox
	for i := range m.hitboxes {
		if m.hitboxes[i].ID == "dict-overlay-search-clear" {
			clearHitbox = &m.hitboxes[i]
			break
		}
	}
	if clearHitbox == nil {
		t.Fatal("expected overlay clear hitbox to be registered")
	}
	if clearHitbox.View != ViewStatistics {
		t.Fatalf("clear hitbox view = %s, want statistics", clearHitbox.View)
	}
	if clearHitbox.Action == nil {
		t.Fatal("expected overlay clear hitbox to use an action")
	}
	if cmd := clearHitbox.Action(); cmd != nil {
		t.Fatal("overlay clear action should not return a command")
	}
	if m.dictionarySearch != "" || len(m.dictionaryResults) != 0 || m.dictionaryDetailView {
		t.Fatalf("expected clear action to reset search/results/detail, search=%q results=%d detail=%v", m.dictionarySearch, len(m.dictionaryResults), m.dictionaryDetailView)
	}
	if m.dictionaryCursor != 0 || m.dictionaryScroll != 0 || m.dictionaryDetailScroll != 0 || m.dictionaryDetailTotalLines != 0 {
		t.Fatalf("expected clear action to reset indexes, cursor=%d scroll=%d detailScroll=%d detailLines=%d", m.dictionaryCursor, m.dictionaryScroll, m.dictionaryDetailScroll, m.dictionaryDetailTotalLines)
	}
}

func TestSpotlightDictionaryHistoryHitboxesAreOverlayScoped(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewBrowser
	m.width = 120
	m.height = 40
	m.dictionaryOverlayActive = true
	m.dictionarySearchHistory = []string{"Apfel", "Birne"}

	_ = m.renderSpotlightDictionary()

	foundClear := false
	foundHistory := false
	for _, hb := range m.hitboxes {
		switch hb.ID {
		case "dict-overlay-history-clear":
			foundClear = true
			if hb.View != ViewBrowser {
				t.Fatalf("history clear hitbox view = %s, want browser", hb.View)
			}
			if hb.Action == nil {
				t.Fatal("expected history clear hitbox to use an action")
			}
			if cmd := hb.Action(); cmd != nil {
				t.Fatal("history clear action should not return a command")
			}
		case "dict-overlay-history-0":
			foundHistory = true
			if hb.View != ViewBrowser {
				t.Fatalf("history hitbox view = %s, want browser", hb.View)
			}
			if hb.Action == nil {
				t.Fatal("expected history hitbox to use an action")
			}
		}
	}
	if !foundClear {
		t.Fatal("expected overlay history clear hitbox")
	}
	if !foundHistory {
		t.Fatal("expected overlay history query hitbox")
	}
	if len(m.dictionarySearchHistory) != 0 {
		t.Fatalf("expected history clear action to empty history, got %v", m.dictionarySearchHistory)
	}
}
