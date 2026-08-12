package tui

import (
	"fmt"
	"testing"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/ankiweb"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

// Review typing mode must accept 'q' (Qualität, Quelle) instead of quitting.
func TestReviewTypingModeAcceptsQ(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewReview
	m.dueCards = []core.Card{{ID: "c1", Prompt: "Quelle", Answer: "source"}}
	m.typingMode = true
	m.typedAnswer = ""
	m.typingChecked = false

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	m = updated.(*Model)

	if cmd != nil {
		t.Fatal("expected 'q' in review typing mode not to quit the app")
	}
	if m.typedAnswer != "q" {
		t.Fatalf("expected 'q' typed into answer, got %q", m.typedAnswer)
	}
}

// Cram grade keys must not advance before the answer is revealed.
func TestCramGradeKeysRequireReveal(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewCram
	m.cramActive = true
	m.cramRevealed = false
	m.cramCards = []core.Card{{ID: "c1", Prompt: "Haus", Answer: "house"}}
	m.cramCursor = 0
	m.cramReviewed = 0

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "g", Code: 'g'})
	m = updated.(*Model)

	if cmd != nil {
		t.Fatal("expected unrevealed cram grade to produce no command")
	}
	if m.cramReviewed != 0 {
		t.Fatalf("expected no grade before reveal, cramReviewed=%d", m.cramReviewed)
	}
	if !m.cramActive || m.cramCursor != 0 {
		t.Fatal("expected cram session to stay on the same card before reveal")
	}

	m.cramRevealed = true
	updated, _ = m.Update(tea.KeyPressMsg{Text: "g", Code: 'g'})
	m = updated.(*Model)
	if m.cramReviewed != 1 {
		t.Fatalf("expected grade after reveal, cramReviewed=%d", m.cramReviewed)
	}
}

// Relative Clauses trainer must type '=' into the answer, not open Spotlight.
func TestRelativeTrainerEqualsTypesNotDictionary(t *testing.T) {
	m, st := trainerModel(t, PracticeSubViewRelative, 2)

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "=", Code: '='})
	m = updated.(*Model)

	if cmd != nil {
		t.Fatal("expected '=' while typing not to open dictionary overlay cmd")
	}
	if m.dictionaryOverlayActive {
		t.Fatal("expected '=' in Relative trainer not to open dictionary spotlight")
	}
	if st.input != "=" {
		t.Fatalf("expected '=' typed into answer, got %q", st.input)
	}
}

// Gender trainer must not treat unused number keys as global view switches.
func TestGenderTrainerIgnoresGlobalNumberNav(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewGender
	m.practiceItems = []practiceItem{{Word: "Haus", Article: "das"}}
	m.practiceIndex = 0
	m.practiceRevealed = false

	updated, _ := m.Update(tea.KeyPressMsg{Text: "0", Code: '0'})
	m = updated.(*Model)

	if m.practiceSubView != PracticeSubViewGender {
		t.Fatalf("expected to stay in Gender trainer, got %v", m.practiceSubView)
	}
	if m.activeView != ViewPractice {
		t.Fatalf("expected to stay on Practice view, got %v", m.activeView)
	}
}

// Gender "press any key" advance must consume q/? instead of quitting/helping.
func TestGenderRevealedAdvancesOnQAndQuestion(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewGender
	m.practiceItems = []practiceItem{
		{Word: "Haus", Article: "das"},
		{Word: "Katze", Article: "die"},
	}
	m.practiceIndex = 0
	m.practiceRevealed = true

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	m = updated.(*Model)
	if cmd != nil {
		t.Fatal("expected 'q' after gender reveal to advance, not quit")
	}
	if m.practiceRevealed {
		t.Fatal("expected advance to clear practiceRevealed")
	}
	if m.practiceIndex != 1 {
		t.Fatalf("expected advance to next noun, index=%d", m.practiceIndex)
	}

	m.practiceRevealed = true
	updated, cmd = m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(*Model)
	if cmd != nil {
		t.Fatal("expected '?' after gender reveal to advance, not open help")
	}
	if m.showHelp {
		t.Fatal("expected help overlay to stay closed")
	}
}

func TestAIDraftingEscCancels(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewAI
	m.drafting = true
	m.aiInput = "Arztbesuch"

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = updated.(*Model)
	if cmd != nil {
		t.Fatal("expected Esc cancel to produce no command")
	}
	if m.drafting {
		t.Fatal("expected Esc to clear drafting lock")
	}
	if !m.draftCancelled {
		t.Fatal("expected draftCancelled so late results are discarded")
	}
	if m.status != "AI drafting cancelled" {
		t.Fatalf("status = %q, want cancelled message", m.status)
	}

	// Late drafts from the cancelled request must not overwrite the UI.
	updated, _ = m.Update(draftsMsg{{Note: core.Note{ID: "late"}}})
	m = updated.(*Model)
	if len(m.drafts) != 0 {
		t.Fatalf("expected cancelled drafts to be discarded, got %d", len(m.drafts))
	}
	if m.draftCancelled {
		t.Fatal("expected draftCancelled to clear after discarding late drafts")
	}
}

func TestBrowserDeckSwitchClearsSelection(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewBrowser
	m.browserCards = []core.Card{{ID: "c1", DeckID: "deck-a"}, {ID: "c2", DeckID: "deck-a"}}
	m.browserSelected = map[string]bool{"c1": true, "c2": true}
	m.deck = core.Deck{ID: "deck-b", Name: "Deck B"}

	_ = m.reloadBrowserForSelectedDeck()

	if len(m.browserSelected) != 0 {
		t.Fatalf("expected selection cleared on deck switch, got %v", m.browserSelected)
	}
}

func TestBrowserReentryClearsTransientState(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewBrowser
	m.browserSearch = "Haus"
	m.browserTag = "noun"
	m.searchingBrowser = true
	m.searchingTags = true
	m.taggingCards = true
	m.tagInput = "old-tag"
	m.browserSelected = map[string]bool{"stale-card": true}
	m.browserCards = []core.Card{{ID: "stale-card"}}

	m.updateView(ViewDashboard)
	m.updateView(ViewBrowser)

	if m.browserSearch != "" || m.browserTag != "" || m.searchingBrowser || m.searchingTags || m.taggingCards || m.tagInput != "" {
		t.Fatalf("browser transient state survived re-entry: search=%q tag=%q searching=%v/%v tagging=%v input=%q",
			m.browserSearch, m.browserTag, m.searchingBrowser, m.searchingTags, m.taggingCards, m.tagInput)
	}
	if len(m.browserSelected) != 0 || len(m.browserCards) != 0 {
		t.Fatalf("browser stale data survived re-entry: selected=%v cards=%v", m.browserSelected, m.browserCards)
	}
}

func TestReviewPredictionsIgnoreStaleCardAndRequest(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewReview
	m.dueCards = []core.Card{{ID: "c1"}, {ID: "c2"}}
	m.cursor = 0

	m.loadReviewPredictions("c1") // request 1
	m.cursor = 1
	m.Update(reviewPredictionsMsg{
		id: 1, cardID: "c1", predictions: map[core.ReviewGrade]time.Duration{core.GradeGood: time.Hour},
	})
	if m.reviewPredictions != nil {
		t.Fatal("prediction for the previous card was applied after navigation")
	}

	m.cursor = 0
	m.loadReviewPredictions("c1") // request 2 for the same card
	m.Update(reviewPredictionsMsg{
		id: 1, cardID: "c1", predictions: map[core.ReviewGrade]time.Duration{core.GradeGood: time.Hour},
	})
	if m.reviewPredictions != nil {
		t.Fatal("older prediction request overwrote the current request")
	}

	fresh := map[core.ReviewGrade]time.Duration{core.GradeGood: 2 * time.Hour}
	m.Update(reviewPredictionsMsg{id: 2, cardID: "c1", predictions: fresh})
	if m.reviewPredictions[core.GradeGood] != 2*time.Hour {
		t.Fatalf("current prediction was not applied: %v", m.reviewPredictions)
	}
}

func TestDeckStatisticsIgnoreStaleRequestsAndRefreshOnSwitch(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDashboard
	m.decks = []core.Deck{{ID: "d1", Name: "One"}, {ID: "d2", Name: "Two"}}
	m.deck = m.decks[0]
	m.stats = core.Statistics{TotalCards: 99}
	m.loadStatistics() // request 1 for d1
	m.deck = m.decks[1]
	m.loadStatistics() // request 2 for d2

	m.Update(statsMsg{id: 1, deckID: "d1", stats: core.Statistics{TotalCards: 1}})
	if m.stats.TotalCards != 99 {
		t.Fatalf("stale deck statistics overwrote current stats: %d", m.stats.TotalCards)
	}
	m.Update(statsMsg{id: 2, deckID: "d2", stats: core.Statistics{TotalCards: 2}})
	if m.stats.TotalCards != 2 {
		t.Fatalf("current deck statistics were not applied: %d", m.stats.TotalCards)
	}

	m.deck = m.decks[0]
	m.deckIndex = 0
	updated, cmd := m.Update(tea.KeyPressMsg{Code: ']'})
	m = updated.(*Model)
	if m.deck.ID != "d2" || cmd == nil {
		t.Fatalf("dashboard deck switch did not refresh statistics: deck=%q cmd=%v", m.deck.ID, cmd != nil)
	}
}

func TestGetSelectedCardIDsOnlyVisible(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.browserCards = []core.Card{{ID: "visible"}}
	m.browserSelected = map[string]bool{"visible": true, "stale-other-deck": true}

	ids := m.getSelectedCardIDs()
	if len(ids) != 1 || ids[0] != "visible" {
		t.Fatalf("expected only visible selected IDs, got %v", ids)
	}
}

// Digit shortcuts must not abandon an active Cram session.
func TestCramNumberKeysDoNotLeaveSession(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewCram
	m.cramActive = true
	m.cramRevealed = false
	m.cramCards = []core.Card{{ID: "c1", Prompt: "Haus", Answer: "house"}}

	updated, _ := m.Update(tea.KeyPressMsg{Text: "2", Code: '2'})
	m = updated.(*Model)
	if m.activeView != ViewCram {
		t.Fatalf("expected to stay on Cram after '2', got %v", m.activeView)
	}
	if !m.cramActive {
		t.Fatal("expected cram session to remain active")
	}
}

func TestPasteIntoDictionarySearch(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryFocusResults = false
	m.dictionaryDetailView = false

	updated, cmd := m.Update(tea.PasteMsg{Content: "gehen"})
	m = updated.(*Model)
	if m.dictionarySearch != "gehen" {
		t.Fatalf("expected pasted dictionary search, got %q", m.dictionarySearch)
	}
	if cmd == nil {
		t.Fatal("expected debounce search command after paste")
	}
}

func TestPasteIntoReviewTypingAndTrainer(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewReview
	m.typingMode = true
	m.typedAnswer = ""
	m.typingChecked = false

	updated, _ := m.Update(tea.PasteMsg{Content: "Qualität"})
	m = updated.(*Model)
	if m.typedAnswer != "Qualität" {
		t.Fatalf("expected pasted review answer, got %q", m.typedAnswer)
	}

	m, st := trainerModel(t, PracticeSubViewRelative, 2)
	updated, _ = m.Update(tea.PasteMsg{Content: "dessen"})
	m = updated.(*Model)
	if st.input != "dessen" {
		t.Fatalf("expected pasted trainer answer, got %q", st.input)
	}
}

func TestDictionaryMouseWheelScrollsResults(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.width, m.height = 100, 40
	m.dictionaryResults = []core.DictionaryEntry{
		{Word: "a"}, {Word: "b"}, {Word: "c"},
	}
	m.dictionaryCursor = 0
	m.dictionaryFocusResults = true

	updated, _ := m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	m = updated.(*Model)
	if m.dictionaryCursor != 1 {
		t.Fatalf("expected wheel down to advance cursor, got %d", m.dictionaryCursor)
	}

	updated, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelUp}))
	m = updated.(*Model)
	if m.dictionaryCursor != 0 {
		t.Fatalf("expected wheel up to move cursor back, got %d", m.dictionaryCursor)
	}
}

// Typing "d" in the dictionary search bar must not open detail when results exist.
func TestDictionarySearchTypesD(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryFocusResults = false
	m.dictionaryDetailView = false
	m.dictionarySearch = ""
	m.dictionaryResults = []core.DictionaryEntry{{Word: "Haus"}}

	updated, _ := m.Update(tea.KeyPressMsg{Text: "d", Code: 'd'})
	m = updated.(*Model)
	if m.dictionaryDetailView {
		t.Fatal("expected 'd' in search bar not to open detail view")
	}
	if m.dictionarySearch != "d" {
		t.Fatalf("expected 'd' typed into search, got %q", m.dictionarySearch)
	}

	// From results focus, 'd' still toggles detail.
	m.dictionarySearch = "haus"
	m.dictionaryFocusResults = true
	updated, _ = m.Update(tea.KeyPressMsg{Text: "d", Code: 'd'})
	m = updated.(*Model)
	if !m.dictionaryDetailView {
		t.Fatal("expected 'd' on focused results to open detail view")
	}
}

func TestSpotlightDictionarySearchTypesA(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewDashboard
	m.dictionaryOverlayActive = true
	m.dictionaryFocusResults = false
	m.dictionaryDetailView = false
	m.dictionarySearch = ""
	// Prior results remain while typing a new query — 'a' must still type.
	m.dictionaryResults = []core.DictionaryEntry{{ID: "1", Word: "Haus", Translation: "house"}}
	m.dictionaryCursor = 0

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "a", Code: 'a'})
	m = updated.(*Model)
	if m.dictionarySearch != "a" {
		t.Fatalf("expected 'a' typed into Spotlight search, got %q", m.dictionarySearch)
	}
	if cmd == nil {
		t.Fatal("expected debounce search tick after typing 'a'")
	}

	// From results focus, 'a' still plays audio (handled, no search mutation).
	m.dictionarySearch = "haus"
	m.dictionaryFocusResults = true
	updated, _ = m.Update(tea.KeyPressMsg{Text: "a", Code: 'a'})
	m = updated.(*Model)
	if m.dictionarySearch != "haus" {
		t.Fatalf("expected focused 'a' not to type into search, got %q", m.dictionarySearch)
	}
}

func TestCramDigitGradingKeys(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewCram
	m.cramActive = true
	m.cramRevealed = true
	m.cramCards = []core.Card{{ID: "c1", Prompt: "Haus", Answer: "house"}}
	m.cramCursor = 0
	m.cramReviewed = 0

	// 3 (Good) should grade card in Cram mode when revealed
	updated, _ := m.Update(tea.KeyPressMsg{Text: "3", Code: '3'})
	m = updated.(*Model)

	if m.cramReviewed != 1 {
		t.Fatalf("expected digit '3' to grade card after reveal, got cramReviewed=%d", m.cramReviewed)
	}
}

func TestTrainerEscKeyOnRevealedCard(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewPreposition
	st := m.trainerStateFor(PracticeSubViewPreposition)
	st.items = []trainerItem{{Answer: "in"}}
	st.index = 0
	st.input = "in"
	st.revealed = true

	updated, _ := m.Update(tea.KeyPressMsg{Text: "esc", Code: 0x1b})
	m = updated.(*Model)

	if m.practiceSubView != PracticeSubViewHub {
		t.Fatalf("expected Esc on revealed card to return to Practice Hub, got subview %v", m.practiceSubView)
	}
}

func TestBrowserLoadIgnoresStaleResults(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "apple", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"},
			{ID: "banana", DeckID: "deck-1", Prompt: "Banana", Answer: "Banane"},
		},
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewBrowser
	m.browserDeckID = "deck-1"

	m.browserSearch = "a"
	stale := m.loadBrowserCards()
	m.browserSearch = "b"
	current := m.loadBrowserCards()

	m.Update(current())
	m.Update(stale())

	if len(m.browserCards) != 1 || m.browserCards[0].ID != "banana" {
		t.Fatalf("stale browser result overwrote current search: %#v", m.browserCards)
	}
}

func TestDictionaryRelatedResultsIgnoreStaleSelection(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryResults = []core.DictionaryEntry{{ID: "new", Word: "neu"}}
	m.dictionaryCursor = 0
	m.dictionaryRelatedID = 2

	m.Update(dictRelatedEntriesMsg{
		id:      1,
		word:    "alt",
		entries: []core.DictionaryEntry{{Word: "altmodisch"}},
	})
	if len(m.dictionaryRelatedEntries) != 0 {
		t.Fatalf("stale related result was applied: %#v", m.dictionaryRelatedEntries)
	}

	m.Update(dictRelatedEntriesMsg{
		id:      2,
		word:    "neu",
		entries: []core.DictionaryEntry{{Word: "neuartig"}},
	})
	if len(m.dictionaryRelatedEntries) != 1 || m.dictionaryRelatedEntries[0].Word != "neuartig" {
		t.Fatalf("current related result was not applied: %#v", m.dictionaryRelatedEntries)
	}
}

func TestDictionaryFindInBrowserKeepsLookupQuery(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "haus", DeckID: "deck-1", Prompt: "Haus", Answer: "house"},
			{ID: "baum", DeckID: "deck-1", Prompt: "Baum", Answer: "tree"},
		},
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewDictionary
	m.dictionaryResults = []core.DictionaryEntry{{ID: "dict-haus", Word: "Haus"}}
	m.dictionaryCursor = 0
	m.dictionaryFocusResults = true

	cmd, handled := m.updateDictionaryKey(tea.KeyPressMsg{Text: "ctrl+f"})
	if !handled {
		t.Fatal("expected ctrl+f to be handled")
	}
	if m.activeView != ViewBrowser {
		t.Fatalf("active view = %s, want browser", m.activeView)
	}
	if m.browserSearch != "Haus" {
		t.Fatalf("browser search = %q, want Haus", m.browserSearch)
	}
	if cmd == nil {
		t.Fatal("expected browser load command")
	}
	m.Update(cmd())
	if len(m.browserCards) != 1 || m.browserCards[0].ID != "haus" {
		t.Fatalf("browser lookup results = %#v, want Haus card", m.browserCards)
	}
}

func TestCramMultiCardProgression(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewCram
	m.cramActive = true
	m.cramCards = []core.Card{
		{ID: "c1", Prompt: "Haus", Answer: "house"},
		{ID: "c2", Prompt: "Baum", Answer: "tree"},
	}
	m.cramCursor = 0
	m.cramRevealed = true

	// Pressing Enter on a revealed card should grade Good, recording review count and advancing cramCursor
	updated, _ := m.Update(tea.KeyPressMsg{Text: "enter", Code: '\r'})
	m = updated.(*Model)

	if m.cramCursor != 1 {
		t.Fatalf("expected cramCursor=1, got %d", m.cramCursor)
	}
	if m.cramReviewed != 1 || m.cramCorrect != 1 {
		t.Fatalf("expected 1/1 reviewed/correct, got %d/%d", m.cramCorrect, m.cramReviewed)
	}

	// Activate card 2 and reveal
	m.cramActive = true
	m.cramRevealed = true
	updated, _ = m.Update(tea.KeyPressMsg{Text: "space", Code: ' '})
	m = updated.(*Model)

	if m.cramReviewed != 2 || m.cramCorrect != 2 {
		t.Fatalf("expected 2/2 reviewed/correct on completion, got %d/%d", m.cramCorrect, m.cramReviewed)
	}
}

func TestSpaceKeySelectionParity(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})

	// 1. Decks View Space key toggle ("space")
	m.activeView = ViewDecks
	m.decks = []core.Deck{{ID: "d1", Name: "Deck 1"}}
	m.deckCursor = 0

	m.Update(tea.KeyPressMsg{Text: "space"})
	if !m.deckSelected["d1"] {
		t.Fatal("expected 'space' key in Decks view to select deck d1")
	}

	// 2. Browser View Space key toggle ("space")
	m.activeView = ViewBrowser
	m.browserCards = []core.Card{{ID: "b1", Prompt: "Word"}}
	m.browserCursor = 0

	m.Update(tea.KeyPressMsg{Text: "space"})
	if !m.browserSelected["b1"] {
		t.Fatal("expected 'space' key in Browser view to select card b1")
	}
}

func TestSettingsTemplateEditNilMapSafety(t *testing.T) {
	// Case 1: no active template set
	m1 := NewModel(&mockRepo{}, &mockScheduler{})
	m1.activeView = ViewSettings
	m1.settingsCursor = 2 // Front Template
	m1.aiTemplateSets = nil
	m1.aiTemplates = nil

	updated1, _ := m1.Update(tea.KeyPressMsg{Text: "enter", Code: '\r'})
	m1 = updated1.(*Model)

	if m1.editingTemplate {
		t.Fatal("expected editingTemplate to be false when no template set is active")
	}
	if m1.status != "No template set available to edit" {
		t.Fatalf("expected status message for missing template set, got %q", m1.status)
	}

	// Case 2: active template set present, but aiTemplates map is nil
	m2 := NewModel(&mockRepo{}, &mockScheduler{})
	m2.activeView = ViewSettings
	m2.settingsCursor = 2
	m2.aiTemplateSets = []string{"default"}
	m2.aiTemplates = nil

	updated2, _ := m2.Update(tea.KeyPressMsg{Text: "enter", Code: '\r'})
	m2 = updated2.(*Model)

	if !m2.editingTemplate {
		t.Fatal("expected editingTemplate to be true after initializing nil template map")
	}
	if m2.aiTemplates["default"] == nil {
		t.Fatal("expected aiTemplates['default'] to be safely initialized")
	}
}

func TestUndoRestoresSessionGradeHistogram(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", DeckID: "d1", Prompt: "P", Answer: "A"}},
		decks:    []core.Deck{{ID: "d1", Name: "Deck"}},
	}
	m := NewModel(repo, &mockScheduler{})
	m.Update(reviewRecordedMsg{
		cardID: "c1",
		cards:  repo.dueCards,
		decks:  repo.decks,
		stats:  core.Statistics{},
		grade:  core.GradeGood,
	})
	if m.sessionGrades[core.GradeGood] != 1 {
		t.Fatalf("expected Good=1 after grade, got %d", m.sessionGrades[core.GradeGood])
	}

	m.Update(reviewUndoneMsg{
		cardID: "c1",
		cards:  repo.dueCards,
		decks:  repo.decks,
		stats:  core.Statistics{},
		grade:  core.GradeGood,
	})
	if m.sessionGrades[core.GradeGood] != 0 {
		t.Fatalf("expected Good=0 after undo, got %d", m.sessionGrades[core.GradeGood])
	}
	if m.sessionReviewed != 0 {
		t.Fatalf("expected sessionReviewed=0 after undo, got %d", m.sessionReviewed)
	}
}

func TestExplainAndFixIgnoreStaleOrCancelledResults(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.explainingCard = true
	m.explainCardID = "c-current"
	m.fixingCard = true
	m.fixCardID = "c-fix"

	m.Update(explainMsg{cardID: "c-stale", explanation: "wrong card"})
	if m.explanation != "" {
		t.Fatalf("stale explain applied: %q", m.explanation)
	}
	if !m.explainingCard {
		t.Fatal("explainingCard should remain true while waiting for matching result")
	}

	m.explainingCard = false
	m.explainCardID = ""
	m.Update(explainMsg{cardID: "c-current", explanation: "after cancel"})
	if m.explanation != "" {
		t.Fatalf("cancelled explain applied: %q", m.explanation)
	}

	m.explainingCard = true
	m.explainCardID = "c-current"
	m.Update(explainMsg{cardID: "c-current", explanation: "ok"})
	if m.explanation != "ok" || m.explainingCard {
		t.Fatalf("matching explain not applied: explanation=%q explaining=%v", m.explanation, m.explainingCard)
	}

	m.Update(fixProposalMsg{
		cardID:   "other",
		oldNote:  core.Note{ID: "n1"},
		proposal: ai.FixedNote{Front: "x"},
	})
	if m.fixProposal != nil {
		t.Fatal("mismatched fix proposal should be ignored")
	}

	m.fixingCard = false
	m.Update(fixProposalMsg{
		cardID:   "c-fix",
		oldNote:  core.Note{ID: "n1"},
		proposal: ai.FixedNote{Front: "y"},
	})
	if m.fixProposal != nil {
		t.Fatal("cancelled fix proposal should be ignored")
	}
}

func TestAnkiWebIgnoresStaleSearchAndInfo(t *testing.T) {
	stub := &stubAnkiWeb{decks: []ankiweb.Deck{{ID: 1, Title: "Old Result"}}}
	m, s := ankiWebModel(t, stub)
	s.query = "german"

	staleSearch := s.search(m)
	stub.decks = []ankiweb.Deck{{ID: 2, Title: "Fresh Result"}, {ID: 3, Title: "Also Fresh"}}
	freshSearch := s.search(m)
	for _, msg := range executeCmd(freshSearch) {
		m.Update(msg)
	}
	for _, msg := range executeCmd(staleSearch) {
		m.Update(msg)
	}
	if len(s.results) != 2 || s.results[0].Title != "Fresh Result" {
		t.Fatalf("stale search overwrote fresh results: %+v", s.results)
	}

	// Move cursor, then deliver info for the previous selection.
	stub.details = ankiweb.Details{Title: "Fresh Result"}
	s.cursor = 0
	staleInfo := s.loadDetails(m)
	s.cursor = 1
	s.details = nil
	for _, msg := range executeCmd(staleInfo) {
		m.Update(msg)
	}
	if s.details != nil {
		t.Fatalf("stale info for previous cursor should not apply after moving: %+v", s.details)
	}
}

func TestBrowserSuspendUsesFullDueReloadLimit(t *testing.T) {
	repo := &mockRepo{
		dueCards: make([]core.Card, 600),
		decks:    []core.Deck{{ID: "d1", Name: "Deck"}},
	}
	for i := range repo.dueCards {
		repo.dueCards[i] = core.Card{
			ID:     fmt.Sprintf("c%04d", i),
			DeckID: "d1",
			Prompt: fmt.Sprintf("P%d", i),
			Answer: fmt.Sprintf("A%d", i),
		}
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewBrowser
	m.browserCards = []core.Card{repo.dueCards[0]}
	m.browserCursor = 0
	m.allDue = repo.dueCards[:500] // simulate previously truncated queue

	cmd := m.toggleBrowserSuspension()
	if cmd == nil {
		t.Fatal("expected suspend command")
	}
	msg := cmd()
	suspended, ok := msg.(cardSuspendedMsg)
	if !ok {
		t.Fatalf("expected cardSuspendedMsg, got %T", msg)
	}
	// One card was suspended; the rest of the 600-card collection must remain
	// (previously a hard 500 limit truncated the review queue).
	if len(suspended.cards) != 599 {
		t.Fatalf("browser suspend due reload returned %d cards, want 599 (full queue minus suspended)", len(suspended.cards))
	}
}

func TestBulkBrowserToggleKindUsesVisibleSelection(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "visible", DeckID: "d1", Kind: core.CardKindFlashcard, Prompt: "V", Answer: "1"},
			{ID: "also", DeckID: "d1", Kind: core.CardKindMCQ, Prompt: "A", Answer: "2"},
		},
	}
	m := NewModel(repo, &mockScheduler{})
	m.activeView = ViewBrowser
	m.browserCards = []core.Card{repo.dueCards[0], repo.dueCards[1]}
	m.browserSelected = map[string]bool{"visible": true, "also": true}

	cmd := m.bulkBrowserToggleKind()
	if cmd == nil {
		t.Fatal("expected bulk toggle command")
	}
	_ = cmd()
	if repo.dueCards[0].Kind != core.CardKindMCQ {
		t.Fatalf("visible flashcard should become MCQ, got %s", repo.dueCards[0].Kind)
	}
	if repo.dueCards[1].Kind != core.CardKindFlashcard {
		t.Fatalf("visible MCQ should become flashcard, got %s", repo.dueCards[1].Kind)
	}
}

// Practice Hub `/` filter must trap globals so typing "q" / digits / paste work.
func TestPracticeHubFilterTrapsGlobalsAndPaste(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewHub
	m.practiceFilterFocus = true

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	m = updated.(*Model)
	if cmd != nil {
		t.Fatal("expected 'q' in Practice Hub filter not to quit")
	}
	if m.practiceFilter != "q" {
		t.Fatalf("expected 'q' typed into filter, got %q", m.practiceFilter)
	}
	if m.activeView != ViewPractice {
		t.Fatalf("expected to stay on Practice, got %v", m.activeView)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Text: "2", Code: '2'})
	m = updated.(*Model)
	if m.activeView != ViewPractice || m.practiceSubView != PracticeSubViewHub {
		t.Fatal("expected digit in hub filter not to switch views")
	}
	if m.practiceFilter != "q2" {
		t.Fatalf("expected filter 'q2', got %q", m.practiceFilter)
	}

	updated, _ = m.Update(tea.PasteMsg{Content: "passive"})
	m = updated.(*Model)
	if m.practiceFilter != "q2passive" {
		t.Fatalf("expected paste into hub filter, got %q", m.practiceFilter)
	}
	if !m.practiceFilterFocus {
		t.Fatal("expected filter focus to remain after paste")
	}
}

// Rapid bookmark-filter toggles must not apply a stale due-queue load.
func TestStaleDueQueueLoadIgnored(t *testing.T) {
	allCards := []core.Card{{ID: "a1"}, {ID: "a2"}}
	bookCards := []core.Card{{ID: "b1"}}
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dueLoadID = 2
	m.bookmarkFilter = true

	// Stale unfiltered load from before the filter was enabled.
	m.Update(dueCardsMsg{id: 1, cards: allCards})
	if len(m.allDue) != 0 {
		t.Fatalf("stale dueCardsMsg should be ignored, got %d cards", len(m.allDue))
	}

	m.Update(bookmarkedDueCardsMsg{id: 2, cards: bookCards})
	if len(m.allDue) != 1 || m.allDue[0].ID != "b1" {
		t.Fatalf("expected current bookmarked load, got %#v", m.allDue)
	}

	m.bookmarkFilter = false
	m.dueLoadID = 3
	m.Update(bookmarkedDueCardsMsg{id: 2, cards: bookCards})
	if len(m.allDue) != 1 || m.allDue[0].ID != "b1" {
		t.Fatalf("mismatched bookmarked load should leave prior queue, got %#v", m.allDue)
	}
	m.Update(dueCardsMsg{id: 3, cards: allCards})
	if len(m.allDue) != 2 {
		t.Fatalf("expected fresh unfiltered load, got %d cards", len(m.allDue))
	}
}

// Cram filter/deck changes must drop mid-flight loads for the previous selection.
func TestStaleCramLoadIgnored(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.deck = core.Deck{ID: "d2", Name: "D2"}
	m.cramType = "suspended"
	m.cramLoadID = 2
	m.cramCards = []core.Card{{ID: "keep"}}

	m.Update(cramCardsMsg{
		id:       1,
		cards:    []core.Card{{ID: "old", Bookmarked: true}},
		cramType: "bookmarked",
		deckID:   "d2",
	})
	if len(m.cramCards) != 1 || m.cramCards[0].ID != "keep" {
		t.Fatalf("stale cram filter load applied: %#v", m.cramCards)
	}

	m.Update(cramCardsMsg{
		id:       2,
		cards:    []core.Card{{ID: "wrong-deck", Suspended: true}},
		cramType: "suspended",
		deckID:   "d1",
	})
	if len(m.cramCards) != 1 || m.cramCards[0].ID != "keep" {
		t.Fatalf("stale cram deck load applied: %#v", m.cramCards)
	}

	m.Update(cramCardsMsg{
		id:       2,
		cards:    []core.Card{{ID: "ok", Suspended: true}, {ID: "skip"}},
		cramType: "suspended",
		deckID:   "d2",
	})
	if len(m.cramCards) != 1 || m.cramCards[0].ID != "ok" {
		t.Fatalf("expected current cram load, got %#v", m.cramCards)
	}
}

// Gender practice loads must not apply after leaving the trainer.
func TestStalePracticeItemsIgnored(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.practiceSubView = PracticeSubViewHub
	m.practiceLoadID = 2
	m.practiceItems = []practiceItem{{Word: "Haus", Article: "das"}}

	m.Update(practiceItemsMsg{
		id:    1,
		items: []practiceItem{{Word: "stale", Article: "der"}},
	})
	if len(m.practiceItems) != 1 || m.practiceItems[0].Word != "Haus" {
		t.Fatalf("stale practice items applied off Gender: %#v", m.practiceItems)
	}

	m.practiceSubView = PracticeSubViewGender
	m.Update(practiceItemsMsg{
		id:    2,
		items: []practiceItem{{Word: "Buch", Article: "das"}},
	})
	if len(m.practiceItems) != 1 || m.practiceItems[0].Word != "Buch" {
		t.Fatalf("expected current Gender load, got %#v", m.practiceItems)
	}
}

// Generic trainer loads must not let a previous visit overwrite a newer one.
func TestStaleTrainerItemsIgnored(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewCase
	st := m.trainerStateFor(PracticeSubViewCase)
	st.loadID = 2
	st.items = []trainerItem{{Title: "current", Answer: "dem"}}

	m.Update(trainerItemsMsg{
		kind:  PracticeSubViewCase,
		id:    1,
		items: []trainerItem{{Title: "stale", Answer: "der"}},
	})
	if len(st.items) != 1 || st.items[0].Title != "current" {
		t.Fatalf("stale trainer items applied: %#v", st.items)
	}

	// Leaving the trainer invalidates even the latest request for that view;
	// its result must wait for the next explicit entry/load.
	m.practiceSubView = PracticeSubViewHub
	m.Update(trainerItemsMsg{
		kind:  PracticeSubViewCase,
		id:    2,
		items: []trainerItem{{Title: "off-screen", Answer: "die"}},
	})
	if st.items[0].Title != "current" {
		t.Fatalf("trainer items applied after leaving trainer: %#v", st.items)
	}

	m.practiceSubView = PracticeSubViewCase
	m.Update(trainerItemsMsg{
		kind:  PracticeSubViewCase,
		id:    2,
		items: []trainerItem{{Title: "fresh", Answer: "dem"}},
	})
	if len(st.items) != 1 || st.items[0].Title != "fresh" {
		t.Fatalf("current trainer items were not applied: %#v", st.items)
	}
}

func TestTrainerLoadCapturesRequestID(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	first := m.enterPracticeMode(PracticeSubViewCase)
	second := m.enterPracticeMode(PracticeSubViewCase)

	firstResult := first()
	firstMsg, ok := firstResult.(trainerItemsMsg)
	if !ok {
		t.Fatalf("first load returned %T, want trainerItemsMsg", firstResult)
	}
	secondResult := second()
	secondMsg, ok := secondResult.(trainerItemsMsg)
	if !ok {
		t.Fatalf("second load returned %T, want trainerItemsMsg", secondResult)
	}
	if firstMsg.id != 1 || secondMsg.id != 2 {
		t.Fatalf("trainer request IDs = %d, %d; want 1, 2", firstMsg.id, secondMsg.id)
	}
}

func TestHelpOverlayIsModalAndScrollable(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.width = 80
	m.height = 24
	m.breakpoint = BreakpointMedium

	updated, _ := m.Update(tea.KeyPressMsg{Text: "?", Code: '?'})
	m = updated.(*Model)
	if !m.showHelp {
		t.Fatal("expected help overlay to open")
	}
	_ = m.View()
	if m.helpTotalLines <= m.helpViewportLines {
		t.Fatalf("expected help content to require scrolling, total=%d viewport=%d", m.helpTotalLines, m.helpViewportLines)
	}

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "0", Code: '0'})
	m = updated.(*Model)
	if cmd != nil || m.activeView != ViewDashboard || !m.showHelp {
		t.Fatalf("help overlay leaked key to underlying view: view=%s help=%v cmd=%v", m.activeView, m.showHelp, cmd != nil)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(*Model)
	if m.helpScroll != 1 {
		t.Fatalf("help down did not scroll one line, scroll=%d", m.helpScroll)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = updated.(*Model)
	if m.showHelp || m.helpScroll != 0 {
		t.Fatalf("help overlay did not close/reset: help=%v scroll=%d", m.showHelp, m.helpScroll)
	}
}

// Grading that finishes after a bookmark-filter flip must not install the old queue.
func TestGradeRespectsLiveBookmarkFilter(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "d1", Prompt: "P", Answer: "A", Bookmarked: true},
			{ID: "c2", DeckID: "d1", Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "d1", Name: "Deck"}},
	}
	m := NewModel(repo, &mockScheduler{})
	m.bookmarkFilter = true
	m.allDue = []core.Card{repo.dueCards[0]}
	m.applyDeckFilter()

	updated, cmd := m.Update(reviewRecordedMsg{
		cardID:         "c1",
		cards:          []core.Card{repo.dueCards[0]}, // snapshotted bookmarked queue
		decks:          repo.decks,
		stats:          core.Statistics{},
		grade:          core.GradeGood,
		bookmarkFilter: false, // grade started before filter was enabled
	})
	m = updated.(*Model)
	if m.sessionReviewed != 1 {
		t.Fatalf("expected session stats updated, sessionReviewed=%d", m.sessionReviewed)
	}
	if cmd == nil {
		t.Fatal("expected reload cmd when bookmark filter mismatched")
	}
	// Stale unfiltered cards must not replace the live bookmarked queue.
	if len(m.allDue) != 1 || m.allDue[0].ID != "c1" {
		t.Fatalf("mismatched grade should keep prior queue until reload, got %#v", m.allDue)
	}
}

// extractPlural must handle UTF-8 characters whose ToLower changes byte size without slicing out of bounds.
func TestExtractPluralUTF8Safety(t *testing.T) {
	// "ẞ" (capital sharp S, 3 bytes) lowercases to "ß" (2 bytes).
	extra := "ẞ Plural: die Bücher; context"
	plural := extractPlural(extra)
	if plural != "die Bücher" {
		t.Fatalf("expected 'die Bücher', got %q", plural)
	}
}

// selectDeckByID must position deckCursor within filteredDecks when a deck filter is active.
func TestSelectDeckByIDFilteredDeckCursor(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.decks = []core.Deck{
		{ID: "d1", Name: "Alpha"},
		{ID: "d2", Name: "Beta"},
		{ID: "d3", Name: "B1 Environment"},
		{ID: "d4", Name: "B1 Transport"},
	}
	m.deckFilter = "B1"

	m.selectDeckByID("d4") // 4th deck overall, but 2nd deck (index 1) in filtered
	if m.deckCursor != 1 {
		t.Fatalf("expected deckCursor=1 in filtered decks, got %d", m.deckCursor)
	}
}
