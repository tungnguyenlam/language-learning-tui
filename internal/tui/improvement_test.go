package tui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"deutsch-tui/internal/core"
)

func TestPracticeHubKeyboardNavigation(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub
	model.practiceHubCursor = 0

	// Move down using 'j'
	cmd, handled := model.updatePracticeKey(tea.KeyPressMsg{Code: 'j'})
	if !handled {
		t.Fatal("Expected 'j' key to be handled in Practice Hub")
	}
	if cmd != nil {
		t.Fatal("Expected cmd to be nil for arrow-key navigation")
	}
	if model.practiceHubCursor != 1 {
		t.Fatalf("Expected cursor to be 1, got %d", model.practiceHubCursor)
	}

	// Move down using Down arrow
	model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyDown})
	if model.practiceHubCursor != 2 {
		t.Fatalf("Expected cursor to be 2, got %d", model.practiceHubCursor)
	}

	// Move up using 'k'
	model.updatePracticeKey(tea.KeyPressMsg{Code: 'k'})
	if model.practiceHubCursor != 1 {
		t.Fatalf("Expected cursor to be 1, got %d", model.practiceHubCursor)
	}

	// Move up using Up arrow
	model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.practiceHubCursor != 0 {
		t.Fatalf("Expected cursor to be 0, got %d", model.practiceHubCursor)
	}

	// Wrap around moving up from 0 — now 12 trainers (indices 0-11)
	model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.practiceHubCursor != 11 {
		t.Fatalf("Expected cursor to be 11 after wrap-around, got %d", model.practiceHubCursor)
	}

	// Wrap around moving down from 11
	model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyDown})
	if model.practiceHubCursor != 0 {
		t.Fatalf("Expected cursor to be 0 after wrap-around, got %d", model.practiceHubCursor)
	}

	// Press Enter to select the gender trainer (index 0)
	cmd, handled = model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled {
		t.Fatal("Expected Enter to be handled")
	}
	if cmd == nil {
		t.Fatal("Expected non-nil command for Enter key entry")
	}
}

func TestReviewUndoAliases(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
		},
		reviews: []core.ReviewResult{{CardID: "c1"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg{cards: repo.dueCards[:1]})
	model.activeView = ViewReview
	model.lastReviewedCardID = "c1"

	// Test 'z' key undo alias
	_, handled := (reviewScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'z'})
	if !handled {
		t.Fatal("Expected 'z' key to be handled in Review")
	}

	// Reset state
	model.lastReviewedCardID = "c1"

	// Test 'ctrl+z' key undo alias
	_, handled = (reviewScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl})
	if !handled {
		t.Fatal("Expected 'ctrl+z' key to be handled in Review")
	}
}

func TestSpotlightDictionaryNarrowScrollbar(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryResults = make([]core.DictionaryEntry, 30)
	for i := 0; i < 30; i++ {
		m.dictionaryResults[i] = core.DictionaryEntry{
			ID:          string(rune('a' + i)),
			Word:        "Word",
			Translation: "Translation",
		}
	}
	m.dictionarySearch = "Word"
	m.dictionaryCursor = 0
	m.width = 60 // Narrow screen triggers single column layout in Spotlight
	m.height = 20

	view := m.renderSpotlightDictionary()
	plainView := stripANSI(view)

	// Since it is narrow and has more results than the interior height, it should render scrollbar characters
	if !strings.Contains(plainView, "┃") && !strings.Contains(plainView, "│") {
		t.Error("Expected narrow spotlight dictionary to render scrollbar characters (┃ or │)")
	}
}

func TestReviewExplanationToggle(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg{cards: repo.dueCards})
	model.activeView = ViewReview

	// First press: m.explanation is empty, so it should trigger explanation retrieval command if provider is enabled
	// Since no provider is enabled, it sets status message.
	model.explanation = "Some explanation"

	// Pressing H again should toggle explanation off
	cmd, handled := (reviewScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'H'})
	if !handled {
		t.Fatal("Expected Shift+H to be handled in Review")
	}
	if cmd != nil {
		t.Fatal("Expected Shift+H toggle-off to return nil command")
	}
	if model.explanation != "" {
		t.Fatalf("Expected explanation to be cleared, got %q", model.explanation)
	}
	if model.status != "Explanation hidden" {
		t.Fatalf("Expected status to be 'Explanation hidden', got %q", model.status)
	}
}

func TestScrollableClippedInputs(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})

	// Test Dictionary long search query clipping
	m.dictionarySearch = "this is an extremely long search query that exceeds normal width limit"
	m.width = 40
	m.height = 20
	layout := m.activeViewContentLayout()
	view := m.renderDictionary(layout)
	if strings.Contains(view, m.dictionarySearch) {
		t.Fatal("Expected long search query to be clipped in Dictionary view")
	}

	// Test Spotlight Dictionary long search query clipping
	viewSpotlight := m.renderSpotlightDictionary()
	if strings.Contains(viewSpotlight, m.dictionarySearch) {
		t.Fatal("Expected long search query to be clipped in Spotlight Dictionary view")
	}

	// Test AI long topic input clipping
	m.activeView = ViewAI
	m.aiInput = "this is an extremely long AI topic input query that exceeds normal width limit"
	viewAI := m.renderAI(0, 0)
	if strings.Contains(viewAI, m.aiInput) {
		t.Fatal("Expected long topic input to be clipped in AI view")
	}
}

func TestReviewMetadataClickHitboxes(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1", Audio: "some-audio.mp3"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg{cards: repo.dueCards})
	model.activeView = ViewReview
	model.focusMode = false

	// Render review view to populate hitboxes
	model.renderReview(0, 0)

	foundBookmark, foundSuspend, foundAudio := false, false, false
	for _, hb := range model.hitboxes {
		if hb.ID == "review-bookmark" {
			foundBookmark = true
		} else if hb.ID == "review-suspend" {
			foundSuspend = true
		} else if hb.ID == "review-audio" {
			foundAudio = true
		}
	}

	if !foundBookmark {
		t.Error("Expected review-bookmark hitbox to be registered")
	}
	if !foundSuspend {
		t.Error("Expected review-suspend hitbox to be registered")
	}
	if !foundAudio {
		t.Error("Expected review-audio hitbox to be registered")
	}
}

func TestPracticeHubVisuals(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	view := model.renderPracticeHub(viewportLayout{Width: 80, Height: 80})

	// Check for icons/labels from all trainers
	icons := []string{"🚻", "🔄", "📐", "🎨", "📍", "👥", "S/", "🔢", "🔗", "Kj", "Pv", "Rl"}
	for _, icon := range icons {
		if !strings.Contains(view, icon) {
			t.Errorf("Expected icon %s in Practice Hub view", icon)
		}
	}
}

func TestDashboardSessionStats(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDashboard

	// Mock a finished session
	model.lastSessionReviewed = 10
	model.lastSessionCorrect = 8
	model.lastSessionDuration = 2 * time.Minute

	view := model.renderDashboard(viewportLayout{Width: 80, Height: 24})

	if !strings.Contains(view, "Last Session: 10 cards") {
		t.Error("Expected last session card count in Dashboard")
	}
	if !strings.Contains(view, "80.0% accuracy") {
		t.Error("Expected accuracy in Dashboard")
	}
	if !strings.Contains(view, "5.0 cards/min") {
		t.Error("Expected cards per minute in Dashboard")
	}
}

func TestRevealSpeedSetting(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})

	// Test instant reveal (0)
	model.revealSpeed = 0
	card := core.Card{ID: "c1", Answer: "test"}
	model.startRevealAnimation(card)

	if model.revealState != RevealRevealed {
		t.Error("Expected instant reveal when speed is 0")
	}
	if model.revealProgress != 100 {
		t.Error("Expected 100% progress when speed is 0")
	}

	// Test animated reveal (5) - now starts with flip animation
	model.revealSpeed = 5
	model.revealState = RevealIdle
	model.revealProgress = 0
	model.flipProgress = 0
	model.flipFrame = 0
	model.startRevealAnimation(card)

	if model.revealState != RevealFlipping {
		t.Error("Expected flipping state when speed is 5")
	}
	if model.flipProgress != 0 {
		t.Error("Expected 0% initial flip progress when speed is 5")
	}
}

func TestMigratedScreensDispatch(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})

	screensToTest := []View{ViewSettings, ViewDebug, ViewStatistics}
	for _, view := range screensToTest {
		if _, registered := model.screens[view]; !registered {
			t.Fatalf("Expected %s screen to be registered in m.screens", view)
		}
		model.activeView = view
		rendered := model.renderActiveViewPlainAt(viewportLayout{Width: 80, Height: 24})
		if rendered == "" || strings.Contains(rendered, "Unknown View") {
			t.Fatalf("Expected valid rendered output for view %s, got empty or unknown", view)
		}
	}
}

func TestPracticeHubHitboxSpacing(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	model.hitboxes = nil
	layout := viewportLayout{X: 0, Y: 0, Width: 80, Height: 80}
	model.renderPracticeHub(layout)

	var buttonHitboxes []Hitbox
	for _, hb := range model.hitboxes {
		if !strings.HasPrefix(hb.ID, "practice-scroll-") {
			buttonHitboxes = append(buttonHitboxes, hb)
		}
	}

	if len(buttonHitboxes) != 12 {
		t.Fatalf("Expected 12 button hitboxes in Practice Hub, got %d", len(buttonHitboxes))
	}
	for i, hb := range buttonHitboxes {
		expectedY := layout.Y + 2 + (i * 5)
		if hb.Y != expectedY {
			t.Errorf("Hitbox %d Y coordinate mismatch: expected %d, got %d", i, expectedY, hb.Y)
		}
	}
}

func TestKonjunktivTrainerIsRegistered(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// '0' key should enter Konjunktiv II trainer
	_, handled := model.updatePracticeKey(tea.KeyPressMsg{Code: '0'})
	if !handled {
		t.Fatal("Expected '0' key to be handled in practice hub")
	}
	if model.practiceSubView != PracticeSubViewKonjunktiv {
		t.Fatalf("Expected subview PracticeSubViewKonjunktiv, got %v", model.practiceSubView)
	}
}

func TestPassiveTrainerIsRegistered(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// Cursor at index 10 + enter should enter Passive Voice trainer
	model.practiceHubCursor = 10
	_, handled := model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled {
		t.Fatal("Expected Enter key to be handled in practice hub")
	}
	if model.practiceSubView != PracticeSubViewPassive {
		t.Fatalf("Expected subview PracticeSubViewPassive, got %v", model.practiceSubView)
	}

	st := model.trainerStateFor(PracticeSubViewPassive)
	if st == nil {
		t.Fatal("Expected non-nil trainer state for PracticeSubViewPassive")
	}
}

func TestKonjunktivTrainerRenders(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewKonjunktiv

	// Simulate trainer items being loaded
	st := model.trainerStateFor(PracticeSubViewKonjunktiv)
	st.items = []trainerItem{
		{
			Title:       "Wenn ich mehr Zeit _____, würde ich mehr lesen.",
			Subtitle:    "Meaning: If I had more time, I would read more.",
			Answer:      "hätte",
			RevealTitle: "Wenn ich mehr Zeit hätte, würde ich mehr lesen.",
			Instruction: "Enter the Konjunktiv II form:",
			HintText:    "Strong Konjunktiv II of haben",
			Explanation: "'haben' → Konjunktiv II: 'hätte'.",
		},
	}

	layout := viewportLayout{Width: 80, Height: 30}
	view := stripANSI(model.renderTrainer(PracticeSubViewKonjunktiv, layout))

	if !strings.Contains(view, "KONJUNKTIV II TRAINER") {
		t.Error("Expected 'KONJUNKTIV II TRAINER' in trainer header")
	}
	if !strings.Contains(view, "Wenn ich mehr Zeit") {
		t.Error("Expected trainer item content in rendered view")
	}
}

func TestRelativeTrainerIsRegistered(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// Test '=' shortcut for Relative Clauses trainer
	_, handled := model.updatePracticeKey(tea.KeyPressMsg{Text: "="})
	if !handled {
		t.Fatal("Expected '=' key to be handled in practice hub")
	}
	if model.practiceSubView != PracticeSubViewRelative {
		t.Fatalf("Expected subview PracticeSubViewRelative, got %v", model.practiceSubView)
	}

	st := model.trainerStateFor(PracticeSubViewRelative)
	if st == nil {
		t.Fatal("Expected non-nil trainer state for PracticeSubViewRelative")
	}
}

func TestRelativeTrainerRenders(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewRelative

	st := model.trainerStateFor(PracticeSubViewRelative)
	st.items = []trainerItem{
		{
			Title:       "Der Mann, _____ dort drüben steht, ist mein Lehrer.",
			Subtitle:    "Meaning: The man who is standing over there is my teacher.",
			Answer:      "der",
			RevealTitle: "Der Mann, der dort drüben steht, ist mein Lehrer.",
			Instruction: "Enter the relative pronoun:",
			HintText:    "Nominative Masculine: der",
			Explanation: "'Der Mann' is masculine singular subject → 'der'.",
		},
	}

	layout := viewportLayout{Width: 80, Height: 30}
	view := stripANSI(model.renderTrainer(PracticeSubViewRelative, layout))

	if !strings.Contains(view, "RELATIVE CLAUSES TRAINER") {
		t.Error("Expected 'RELATIVE CLAUSES TRAINER' in trainer header")
	}
	if !strings.Contains(view, "Der Mann") {
		t.Error("Expected trainer item content in rendered view")
	}
}

func TestPracticeHubEqualsOpensRelativeNotDictionary(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub
	model.width = 100
	model.height = 40

	// '=' must reach the Practice Hub Relative trainer, not the global spotlight overlay.
	updated, _ := model.Update(tea.KeyPressMsg{Text: "="})
	m := updated.(*Model)
	if m.dictionaryOverlayActive {
		t.Fatal("expected '=' on Practice Hub not to open dictionary spotlight")
	}
	if m.practiceSubView != PracticeSubViewRelative {
		t.Fatalf("expected PracticeSubViewRelative, got %v", m.practiceSubView)
	}

	// Outside the hub, '=' still opens the dictionary overlay.
	m.practiceSubView = PracticeSubViewHub
	m.activeView = ViewDashboard
	updated, _ = m.Update(tea.KeyPressMsg{Text: "="})
	m = updated.(*Model)
	if !m.dictionaryOverlayActive {
		t.Fatal("expected '=' on Dashboard to open dictionary spotlight")
	}
}

func TestDecksViewDirectExportShortcut(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{
			{ID: "test_deck", Name: "Test Deck"},
		},
		dueCards: []core.Card{
			{ID: "c1", NoteID: "n1", DeckID: "test_deck", Prompt: "Hund", Answer: "dog"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDecks
	model.decks = repo.decks
	model.deckCursor = 0
	model.width = 100
	model.height = 30

	// Assert footer hint contains 'e export'
	out := stripANSI(model.renderDecks(model.activeViewContentLayout()))
	if !strings.Contains(out, "e export") {
		t.Fatalf("expected Decks view footer to contain 'e export', got: %s", out)
	}

	// Press 'e' key to trigger deck export
	cmd, handled := (decksScreen{}).HandleKey(model, tea.KeyPressMsg{Code: 'e'})
	if !handled || cmd == nil {
		t.Fatalf("expected 'e' key to trigger deck export command")
	}
	msgs := executeCmd(cmd)
	foundExportStatus := false
	for _, m := range msgs {
		if status, ok := m.(statusMsg); ok && strings.Contains(status.text, "Exported") {
			foundExportStatus = true
			break
		}
	}
	if !foundExportStatus {
		t.Fatalf("expected status message confirming export, got msgs: %v", msgs)
	}
}

func TestGenderTrainerEmptyItemsSafety(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewGender
	model.practiceItems = nil
	model.practiceIndex = 5

	// Render should not panic
	out := model.renderGenderTrainer(model.activeViewContentLayout())
	if !strings.Contains(out, "No nouns found for practice") {
		t.Fatalf("expected empty message, got %q", out)
	}

	// Pressing keys should not panic
	for _, k := range []string{"1", "2", "3", "d", "i", "f", "a", "m", "n"} {
		_, _ = model.updatePracticeKey(tea.KeyPressMsg{Text: k})
	}
}

func TestGenericTrainerInvalidIndexSafety(t *testing.T) {
	repo := &mockRepo{}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewCase
	st := model.trainerStateFor(PracticeSubViewCase)
	st.items = []trainerItem{{Title: "Test"}}
	st.index = 99 // out of bounds

	// Render should not panic
	out := model.renderTrainer(PracticeSubViewCase, model.activeViewContentLayout())
	if out == "" {
		t.Fatal("expected non-empty output")
	}

	// Key update should not panic
	_, _ = model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyEnter})
}

func TestDecksViewMergeUnclampedCursorSafety(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{
			{ID: "d1", Name: "Deck 1"},
			{ID: "d2", Name: "Deck 2"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewDecks
	model.decks = repo.decks
	model.deckSelected = map[string]bool{"d1": true}
	model.deckCursor = 50 // out of bounds

	cmd := model.handleDeckMerge()
	if cmd == nil {
		t.Fatal("expected merge command")
	}
	if !strings.Contains(model.status, "Merging 1 decks into Deck 2") {
		t.Fatalf("expected status containing 'Merging 1 decks into Deck 2', got %q", model.status)
	}
}
