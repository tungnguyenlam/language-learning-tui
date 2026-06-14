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

	// Wrap around moving up from 0
	model.updatePracticeKey(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.practiceHubCursor != 8 {
		t.Fatalf("Expected cursor to be 8 after wrap-around, got %d", model.practiceHubCursor)
	}

	// Wrap around moving down from 8
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
	model.Update(dueCardsMsg(repo.dueCards[:1]))
	model.activeView = ViewReview
	model.lastReviewedCardID = "c1"

	// Test 'z' key undo alias
	_, handled := model.updateReviewKey(tea.KeyPressMsg{Code: 'z'})
	if !handled {
		t.Fatal("Expected 'z' key to be handled in Review")
	}

	// Reset state
	model.lastReviewedCardID = "c1"

	// Test 'ctrl+z' key undo alias
	_, handled = model.updateReviewKey(tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl})
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
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	// First press: m.explanation is empty, so it should trigger explanation retrieval command if provider is enabled
	// Since no provider is enabled, it sets status message.
	model.explanation = "Some explanation"

	// Pressing H again should toggle explanation off
	cmd, handled := model.updateReviewKey(tea.KeyPressMsg{Code: 'H'})
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
	model.Update(dueCardsMsg(repo.dueCards))
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

	view := model.renderPracticeHub(viewportLayout{Width: 80, Height: 24})

	// Check for some icons and colors in the view (using substrings)
	icons := []string{"🚻", "🔄", "📐", "🎨", "📍", "👥", "S/", "🔢", "🔗"}
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
