package tui

import (
	"context"
	"errors"
	"os"
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

func TestDictionarySearchResultsStatus(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionarySearchID = 1
	m.dictionarySearch = "zzzz"

	updated, cmd := m.Update(dictionarySearchResultsMsg{id: 1, results: nil})
	if cmd != nil {
		t.Fatal("expected no command for dictionary results message")
	}
	got := updated.(*Model)
	if got.status != `No dictionary results for "zzzz"` {
		t.Fatalf("unexpected empty-results status: %q", got.status)
	}

	got.dictionarySearchID = 2
	got.dictionarySearch = ""
	updated, _ = got.Update(dictionarySearchResultsMsg{id: 2, results: nil})
	got = updated.(*Model)
	if got.status != "Dictionary search cleared" {
		t.Fatalf("unexpected cleared-search status: %q", got.status)
	}
}

type captureRepo struct {
	mockRepo
	upsertedNote *core.Note
	notes        map[string]core.Note
}

func (r *captureRepo) UpsertNote(ctx context.Context, note core.Note) error {
	r.upsertedNote = &note
	if r.notes == nil {
		r.notes = make(map[string]core.Note)
	}
	r.notes[note.ID] = note
	return nil
}

func (r *captureRepo) GetNote(ctx context.Context, noteID string) (core.Note, error) {
	note, ok := r.notes[noteID]
	if !ok {
		return core.Note{}, context.Canceled
	}
	return note, nil
}

func (r *captureRepo) UpsertDeck(ctx context.Context, deck core.Deck) error {
	return nil
}

func (r *captureRepo) GetDeck(ctx context.Context, id string) (core.Deck, error) {
	for _, d := range r.decks {
		if d.ID == id {
			return d, nil
		}
	}
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

	msgs := executeCmd(cmd)
	var status statusMsg
	var foundStatus bool
	for _, msg := range msgs {
		if s, ok := msg.(statusMsg); ok {
			status = s
			foundStatus = true
			break
		}
	}
	if !foundStatus {
		t.Fatalf("expected statusMsg in executed msgs, got %v", msgs)
	}

	if !strings.Contains(status.text, "Added") {
		t.Errorf("expected status text to say Added, got %q", status.text)
	}

	if repo.upsertedNote == nil {
		t.Fatal("expected note to be upserted")
	}

	note := repo.upsertedNote
	if note.Front != "das Auto" || note.Back != "Car" {
		t.Errorf("unexpected note content: front=%q, back=%q", note.Front, note.Back)
	}

	if len(note.Cards) == 0 {
		t.Fatal("expected cards to be generated for the quick-added dictionary entry")
	}

	card := note.Cards[0]
	if card.Prompt != "das Auto" || card.Answer != "Car" {
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
	if !handled {
		t.Fatal("expected ctrl+d to be handled")
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
	repo := &mockRepo{}
	m := NewModel(repo, &mockScheduler{})
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
	if saved, _ := repo.GetSetting(context.Background(), "dict_search_history"); saved != "null" {
		t.Errorf("expected click clear to persist an empty history, got %q", saved)
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

func TestSpotlightDictionaryOverlayResultCount(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 120
	m.height = 40
	m.activeView = ViewDashboard
	m.dictionaryOverlayActive = true
	m.dictionarySearch = "Apfel"
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Apfel", Translation: "apple"},
		{ID: "2", Word: "Apfelbaum", Translation: "apple tree"},
	}

	output := stripANSI(m.renderSpotlightDictionary())
	if !strings.Contains(output, "SPOTLIGHT DICTIONARY (2 results)") {
		t.Fatalf("expected overlay title to include result count, got:\n%s", output)
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

func TestDictionaryPersistentSearchHistory(t *testing.T) {
	repo := &mockRepo{}
	m := NewModel(repo, &mockScheduler{})

	// 1. Record search history
	m.recordDictionarySearch("Käse")
	m.recordDictionarySearch("Brot")

	val, err := repo.GetSetting(context.Background(), "dict_search_history")
	if err != nil {
		t.Fatalf("failed to get setting: %v", err)
	}
	if !strings.Contains(val, "Käse") || !strings.Contains(val, "Brot") {
		t.Fatalf("expected saved history to contain Käse and Brot, got %q", val)
	}

	// 2. Load search history
	m2 := NewModel(repo, &mockScheduler{})
	cmd := m2.loadDictionaryHistory()
	if cmd == nil {
		t.Fatal("expected loadDictionaryHistory cmd to be returned")
	}
	msg := cmd()
	loadedMsg, ok := msg.(dictHistoryLoadedMsg)
	if !ok {
		t.Fatalf("expected dictHistoryLoadedMsg, got %T", msg)
	}
	if len(loadedMsg) != 2 || loadedMsg[0] != "Käse" || loadedMsg[1] != "Brot" {
		t.Fatalf("unexpected loaded history: %v", loadedMsg)
	}

	// 3. Clear history
	m.activeView = ViewDictionary
	m.Update(tea.KeyPressMsg{Code: 'x', Mod: tea.ModCtrl}) // ctrl+x

	val, err = repo.GetSetting(context.Background(), "dict_search_history")
	if err != nil {
		t.Fatalf("failed to get setting: %v", err)
	}
	if val != "[]" && val != "" && val != "null" {
		t.Fatalf("expected cleared history to be empty, got %q", val)
	}
}

func TestDictionaryHistoryKeyboardCycling(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionarySearchHistory = []string{"Apfel", "Birne"}
	m.dictionarySearch = ""

	// Down arrow cycles to recent history
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyDown})
	if !handled || cmd == nil {
		t.Fatal("expected Down arrow to be handled with search command")
	}
	if m.dictionarySearch != "Birne" {
		t.Fatalf("expected search query to be 'Birne', got %q", m.dictionarySearch)
	}

	// Up arrow cycles to previous history item
	cmd, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: tea.KeyUp})
	if !handled || cmd == nil {
		t.Fatal("expected Up arrow to be handled with search command")
	}
	if m.dictionarySearch != "Apfel" {
		t.Fatalf("expected search query to be 'Apfel', got %q", m.dictionarySearch)
	}
}

func TestDictionaryCtrlEExplainAndFilterTagHelpers(t *testing.T) {
	// 1. Test filter tag helpers
	q := "Haus :noun"
	if !isFilterActive(q, ":noun") {
		t.Errorf("expected :noun to be active in %q", q)
	}
	if isFilterActive(q, ":verb") {
		t.Errorf("expected :verb to not be active in %q", q)
	}

	toggled := toggleFilterTag(q, ":noun")
	if strings.Contains(toggled, ":noun") {
		t.Errorf("expected toggleFilterTag to remove :noun, got %q", toggled)
	}

	cleared := clearFilterTags("gehen :verb :de")
	if cleared != "gehen" {
		t.Errorf("expected clearFilterTags to return 'gehen', got %q", cleared)
	}

	// 2. Test ctrl+e shortcut switching to ViewAI with explain prompt
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "geheim", Translation: "secret", WordClass: "adj"},
	}
	m.dictionaryCursor = 0

	_, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
	if !handled {
		t.Fatal("expected ctrl+e to be handled")
	}
	if m.activeView != ViewAI {
		t.Fatalf("expected activeView to switch to ViewAI, got %s", m.activeView)
	}
	if !strings.Contains(m.aiInput, "Explain the German word 'geheim'") {
		t.Fatalf("expected aiInput to contain explain prompt, got %q", m.aiInput)
	}
}

func TestDictionaryKAndJKeyHandling(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionarySearchHistory = []string{"Apfel", "Birne"}
	m.dictionarySearch = ""
	m.dictionaryFocusResults = false
	m.dictionaryDetailView = false

	// Typing 'k' while in search input must append 'k' to dictionarySearch, NOT cycle history or navigate
	_, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'k'})
	if !handled {
		t.Fatal("expected 'k' to be handled as printable text input")
	}
	if m.dictionarySearch != "k" {
		t.Fatalf("expected search input to be 'k', got %q", m.dictionarySearch)
	}

	// Typing 'j' while in search input must append 'j' to dictionarySearch
	_, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: 'j'})
	if !handled {
		t.Fatal("expected 'j' to be handled as printable text input")
	}
	if m.dictionarySearch != "kj" {
		t.Fatalf("expected search input to be 'kj', got %q", m.dictionarySearch)
	}

	// When results list is focused (dictionaryFocusResults = true), 'k' and 'j' navigate results list
	m.dictionaryFocusResults = true
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Käse"},
		{ID: "2", Word: "Kuchen"},
	}
	m.dictionaryCursor = 0

	_, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: 'j'})
	if !handled {
		t.Fatal("expected 'j' to navigate down when results are focused")
	}
	if m.dictionaryCursor != 1 {
		t.Fatalf("expected cursor to be 1 after 'j', got %d", m.dictionaryCursor)
	}

	_, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: 'k'})
	if !handled {
		t.Fatal("expected 'k' to navigate up when results are focused")
	}
	if m.dictionaryCursor != 0 {
		t.Fatalf("expected cursor to be 0 after 'k', got %d", m.dictionaryCursor)
	}
}

type mockDictRepo struct {
	mockRepo
	entries map[string]core.DictionaryEntry
}

func (r *mockDictRepo) GetEntry(ctx context.Context, id string) (core.DictionaryEntry, error) {
	if entry, ok := r.entries[id]; ok {
		return entry, nil
	}
	return core.DictionaryEntry{}, errors.New("entry not found")
}

func (r *mockDictRepo) Search(ctx context.Context, query string, limit int) ([]core.DictionaryEntry, error) {
	var res []core.DictionaryEntry
	for _, e := range r.entries {
		res = append(res, e)
	}
	return res, nil
}

func (r *mockDictRepo) ImportEntries(ctx context.Context, entries []core.DictionaryEntry) error {
	return nil
}

func (r *mockDictRepo) DictionaryCount(ctx context.Context) (int, error) {
	return len(r.entries), nil
}

func TestDictionaryStarringAndFiltering(t *testing.T) {
	repo := &mockDictRepo{
		entries: map[string]core.DictionaryEntry{
			"1": {ID: "1", Word: "Hund", Translation: "dog", Gender: "m"},
			"2": {ID: "2", Word: "Katze", Translation: "cat", Gender: "f"},
		},
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.dictionaryResults = []core.DictionaryEntry{
		repo.entries["1"],
		repo.entries["2"],
	}
	m.dictionaryCursor = 0
	m.dictionaryFocusResults = true

	// Toggle star on "Hund" using 'b' key
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'b'})
	if !handled {
		t.Fatal("expected 'b' to be handled when results focused")
	}
	_ = cmd
	if !m.dictionaryStarred["1"] {
		t.Fatalf("expected entry '1' (Hund) to be starred")
	}
	if !strings.Contains(m.status, "Starred") || !strings.Contains(m.status, "Hund") {
		t.Fatalf("unexpected status after starring: %q", m.status)
	}

	// Verify rendered 2-column view contains star symbol ★ for Hund
	view := stripANSI(m.renderDictionary(m.activeViewContentLayout()))
	if !strings.Contains(view, "★ Hund") {
		t.Fatalf("expected rendered view to show ★ Hund, got:\n%s", view)
	}

	// Test :starred filter
	m.dictionarySearch = ":starred"
	cmdSearch := m.searchDictionary()
	if cmdSearch == nil {
		t.Fatal("expected searchDictionary command for :starred filter")
	}
	msg := cmdSearch()
	resMsg, ok := msg.(dictionarySearchResultsMsg)
	if !ok {
		t.Fatalf("expected dictionarySearchResultsMsg, got %T", msg)
	}
	if len(resMsg.results) != 1 || resMsg.results[0].Word != "Hund" {
		t.Fatalf("expected :starred filter to return only Hund, got %v", resMsg.results)
	}

	// Unstar "Hund" using ctrl+b
	cmd, handled = m.updateDictionaryKey(tea.KeyPressMsg{Code: 'b', Mod: tea.ModCtrl})
	if !handled {
		t.Fatal("expected ctrl+b to be handled")
	}
	_ = cmd
	if m.dictionaryStarred["1"] {
		t.Fatalf("expected entry '1' (Hund) to be unstarred")
	}
}

func TestDictionaryTwoColumnTranslationDisplay(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 110
	m.height = 40
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Haus", Translation: "house", Gender: "n"},
		{ID: "2", Word: "Baum", Translation: "tree", Gender: "m"},
	}

	view := stripANSI(m.renderDictionary(m.activeViewContentLayout()))
	if !strings.Contains(view, "Haus {n} - house") {
		t.Fatalf("expected 2-column list item to show 'Haus {n} - house', got:\n%s", view)
	}
	if !strings.Contains(view, "Baum {m} - tree") {
		t.Fatalf("expected 2-column list item to show 'Baum {m} - tree', got:\n%s", view)
	}
}

func TestDictionaryTargetDeckCyclingAndBatchAdd(t *testing.T) {
	repo := &captureRepo{}
	repo.decks = []core.Deck{
		{ID: "deck-1", Name: "German A1"},
		{ID: "deck-2", Name: "German B1"},
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewDictionary
	m.decks = repo.decks

	// Cycle deck with ctrl+g
	_, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'g', Mod: tea.ModCtrl})
	if !handled {
		t.Fatal("expected ctrl+g to be handled")
	}
	if m.dictionaryTargetDeckID != "deck-1" {
		t.Fatalf("expected target deck ID 'deck-1', got %q", m.dictionaryTargetDeckID)
	}

	// Add single entry with ctrl+a
	entry := core.DictionaryEntry{ID: "101", Word: "Buch", Translation: "book", Gender: "n"}
	m.dictionaryResults = []core.DictionaryEntry{entry}
	m.dictionaryCursor = 0
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	if !handled || cmd == nil {
		t.Fatal("expected ctrl+a to return batch command")
	}
	msgs := executeCmd(cmd)
	foundStatus := false
	for _, msg := range msgs {
		if s, ok := msg.(statusMsg); ok {
			foundStatus = true
			if !strings.Contains(s.text, "German A1") {
				t.Fatalf("expected status to mention target deck 'German A1', got %q", s.text)
			}
		}
	}
	if !foundStatus {
		t.Fatal("expected statusMsg from add command")
	}

	// Batch add results with ctrl+s
	cmdBatch, handledBatch := m.updateDictionaryKey(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	if !handledBatch || cmdBatch == nil {
		t.Fatal("expected ctrl+s to return batch command")
	}
	batchMsgs := executeCmd(cmdBatch)
	foundBatchStatus := false
	for _, msg := range batchMsgs {
		if s, ok := msg.(statusMsg); ok {
			foundBatchStatus = true
			if !strings.Contains(s.text, "Added 0 entries to German A1 deck; skipped 1 already added") {
				t.Fatalf("unexpected batch status: %q", s.text)
			}
		}
	}
	if !foundBatchStatus {
		t.Fatal("expected statusMsg from batch add command")
	}
}

func TestDictionaryClozeCardGeneration(t *testing.T) {
	repo := &captureRepo{}
	m := NewModel(repo, &mockScheduler{})

	entry := core.DictionaryEntry{
		ID:          "cloze-1",
		Word:        "Hund",
		Translation: "dog",
		Gender:      "m",
		Examples:    []string{"Mein Hund ist braun."},
	}

	cmd := m.addDictionaryClozeEntryCmd(entry)
	if cmd == nil {
		t.Fatal("expected command from addDictionaryClozeEntryCmd")
	}

	msgs := executeCmd(cmd)
	foundStatus := false
	for _, msg := range msgs {
		if s, ok := msg.(statusMsg); ok {
			foundStatus = true
			if !strings.Contains(s.text, "Created Cloze card for 'Hund'") {
				t.Errorf("unexpected status message: %q", s.text)
			}
		}
	}
	if !foundStatus {
		t.Fatal("expected statusMsg in executed msgs")
	}

	if repo.upsertedNote == nil {
		t.Fatal("expected note to be upserted")
	}
	note := repo.upsertedNote
	if !strings.Contains(note.Front, "{{c1::Hund}}") {
		t.Errorf("expected Cloze prompt to contain {{c1::Hund}}, got %q", note.Front)
	}
	if note.Type != "cloze" {
		t.Errorf("expected note type 'cloze', got %q", note.Type)
	}
	if len(note.Cards) == 0 || note.Cards[0].Kind != core.CardKindCloze {
		t.Fatalf("expected Cloze card in note.Cards, got %v", note.Cards)
	}
}

func TestDictionaryRecentlyViewedAndDomainTags(t *testing.T) {
	repo := &mockRepo{}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40

	entry1 := core.DictionaryEntry{ID: "10", Word: "Käfer", Translation: "beetle", Tags: []string{"[zool.]"}}
	entry2 := core.DictionaryEntry{ID: "20", Word: "Anapher", Translation: "anaphora", Tags: []string{"[lit.]"}}

	m.recordDictionaryView(entry1)
	m.recordDictionaryView(entry2)

	if len(m.dictionaryRecentlyViewed) != 2 || m.dictionaryRecentlyViewed[0].Word != "Anapher" {
		t.Fatalf("unexpected recently viewed stack: %v", m.dictionaryRecentlyViewed)
	}

	// Recent words are part of the learner's lookup workflow and survive a restart.
	m2 := NewModel(repo, &mockScheduler{})
	msg := m2.loadDictionaryRecentlyViewed()()
	loaded, ok := msg.(dictRecentlyViewedLoadedMsg)
	if !ok || len(loaded) != 2 || loaded[0].Word != "Anapher" {
		t.Fatalf("unexpected persisted recently viewed entries: %v", msg)
	}

	// 1. Verify Recently Inspected Words section rendered on empty search
	m.dictionarySearch = ""
	m.dictionaryResults = nil
	viewEmpty := stripANSI(m.renderDictionary(m.activeViewContentLayout()))
	if !strings.Contains(viewEmpty, "Recently Inspected Words:") || !strings.Contains(viewEmpty, "Anapher") {
		t.Fatalf("expected view to contain Recently Inspected Words, got:\n%s", viewEmpty)
	}

	// 2. Verify domain tag badges in detail view
	m.dictionaryResults = []core.DictionaryEntry{entry1}
	m.dictionaryCursor = 0
	m.dictionaryDetailView = true
	viewDetail := stripANSI(m.renderDictionary(m.activeViewContentLayout()))
	if !strings.Contains(viewDetail, "[ZOOL.]") {
		t.Fatalf("expected detail view to contain domain tag badge [ZOOL.], got:\n%s", viewDetail)
	}
}

func TestDictionaryExportTSVAndFilterPills(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width = 100
	m.height = 40
	m.dictionaryResults = []core.DictionaryEntry{
		{ID: "1", Word: "Hund", Translation: "dog", Gender: "m"},
		{ID: "2", Word: "Katze", Translation: "cat", Gender: "f"},
	}

	// Test ctrl+o triggers export command
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl})
	if !handled || cmd == nil {
		t.Fatal("expected ctrl+o to be handled with export command")
	}

	msg := cmd()
	status, ok := msg.(statusMsg)
	if !ok {
		t.Fatalf("expected statusMsg, got %T", msg)
	}
	if !strings.Contains(status.text, "Exported 2 dictionary entries") {
		t.Fatalf("unexpected export status text: %q", status.text)
	}

	// Verify file was created in exports/
	if _, err := os.Stat("exports/dictionary_export.tsv"); os.IsNotExist(err) {
		t.Fatal("expected exports/dictionary_export.tsv to exist")
	}

	// Test DE/EN filter pills hitboxes
	_ = m.renderDictionary(m.activeViewContentLayout())
	var deHitbox, deckPillHitbox bool
	for _, hb := range m.hitboxes {
		if strings.HasPrefix(hb.ID, "dict-filter-pill-") {
			deHitbox = true
		}
		if hb.ID == "dict-target-deck-pill" {
			deckPillHitbox = true
			if hb.Action != nil {
				hb.Action()
			}
		}
	}
	if !deHitbox {
		t.Fatal("expected filter pills hitboxes to be registered")
	}
	if !deckPillHitbox {
		t.Fatal("expected target deck pill hitbox to be registered")
	}
}

func TestDictionaryScrollbarLineCountAndOverscroll(t *testing.T) {
	m := NewModel(nil, nil)
	m.width = 120
	m.height = 30
	m.activeView = ViewDictionary

	// Populate results with 5 items
	for i := 0; i < 5; i++ {
		m.dictionaryResults = append(m.dictionaryResults, core.DictionaryEntry{
			ID:          "entry-" + string(rune('a'+i)),
			Word:        "word" + string(rune('a'+i)),
			Translation: "trans" + string(rune('a'+i)),
		})
	}

	// Set scroll to 10 (overscrolled beyond 5 results)
	m.dictionaryScroll = 10
	m.dictionaryDetailScroll = 10

	view := m.renderDictionary(m.activeViewContentLayout())
	lines := strings.Split(view, "\n")

	// Ensure view height does not blow up due to negative fill loop indices
	if len(lines) > 40 {
		t.Fatalf("renderDictionary produced excessive lines (%d) under overscroll conditions", len(lines))
	}
}

func TestFormatInflectionTable(t *testing.T) {
	// Verb forms
	vGot := stripANSI(formatInflectionTable("verb", "", "geht, ging, ist gegangen", "ging"))
	if !strings.Contains(vGot, "Verb Forms:") || !strings.Contains(vGot, "Präsens (3sg): geht") || !strings.Contains(vGot, "Präteritum:    ging") {
		t.Errorf("unexpected verb inflection table output: %s", vGot)
	}

	// Noun forms
	nGot := stripANSI(formatInflectionTable("noun", "m", "des Mannes, die Männer", "Männer"))
	if !strings.Contains(nGot, "Noun Forms:") || !strings.Contains(nGot, "Genitiv: des Mannes") || !strings.Contains(nGot, "Plural:  die Männer") {
		t.Errorf("unexpected noun inflection table output: %s", nGot)
	}

	// Adjective forms
	aGot := stripANSI(formatInflectionTable("adj", "", "größer, am größten", "größer"))
	if !strings.Contains(aGot, "Adjective Comparison:") || !strings.Contains(aGot, "Komparativ: größer") || !strings.Contains(aGot, "Superlativ: am größten") {
		t.Errorf("unexpected adj inflection table output: %s", aGot)
	}
}

func TestDictionaryNumberKeyNavigation(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryResults = []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Handschuh",
			Translation: "Glove",
		},
	}
	m.dictionaryCursor = 0
	m.dictionaryFocusResults = true

	// Pressing '1' should jump search to 1st compound part "Hand"
	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Text: "1", Code: '1'})
	if !handled || cmd == nil {
		t.Fatalf("expected key '1' to be handled with search command")
	}
	if m.dictionarySearch != "Hand" {
		t.Fatalf("expected search to be 'Hand', got %q", m.dictionarySearch)
	}

	// Pressing '2' on Handschuh should jump search to 2nd compound part "Schuh"
	m.dictionarySearch = "Handschuh"
	m.dictionaryFocusResults = true
	cmd, handled = m.updateDictionaryKey(tea.KeyPressMsg{Text: "2", Code: '2'})
	if !handled || cmd == nil {
		t.Fatalf("expected key '2' to be handled with search command")
	}
	if m.dictionarySearch != "Schuh" {
		t.Fatalf("expected search to be 'Schuh', got %q", m.dictionarySearch)
	}
}
