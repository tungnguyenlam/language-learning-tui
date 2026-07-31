package tui

import (
	"testing"

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
