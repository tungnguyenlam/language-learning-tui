package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

func TestBreakpointForWidth(t *testing.T) {
	tests := []struct {
		width int
		want  Breakpoint
	}{
		{50, BreakpointCompact},
		{80, BreakpointMedium},
		{120, BreakpointWide},
	}
	for _, tt := range tests {
		if got := breakpointForWidth(tt.width); got != tt.want {
			t.Fatalf("breakpointForWidth(%d) = %s, want %s", tt.width, got, tt.want)
		}
	}
}

func TestReviewHistoryRenderingAndToggle(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", Prompt: "der Apfel", Answer: "apple"}},
		reviews: []core.ReviewResult{{
			CardID:   "c1",
			Grade:    core.GradeGood,
			Reviewed: time.Date(2026, 4, 30, 8, 0, 0, 0, time.UTC),
			Next:     core.ReviewState{CardID: "c1", Due: time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC), Interval: 24 * time.Hour, Reviews: 1},
		}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewReview
	model.dueCards = repo.dueCards

	cmd := model.toggleReviewHistory()
	if cmd == nil {
		t.Fatal("toggleReviewHistory should load review history")
	}
	msg := cmd()
	model.Update(msg)

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "Review History: der Apfel") {
		t.Fatalf("review history header missing: %s", view)
	}
	if !strings.Contains(view, "good") || !strings.Contains(view, "1 day") {
		t.Fatalf("review history details missing: %s", view)
	}

	model.Update(tea.KeyPressMsg{Code: 'r'})
	if model.showReviewHistory {
		t.Fatal("second r should hide review history")
	}
}

func TestBrowserEnterShowsSelectedCardHistory(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", Prompt: "der Apfel", Answer: "apple"}},
		reviews: []core.ReviewResult{{
			CardID:   "c1",
			Grade:    core.GradeHard,
			Reviewed: time.Date(2026, 4, 30, 8, 0, 0, 0, time.UTC),
			Next:     core.ReviewState{CardID: "c1", Due: time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC), Interval: 2 * time.Hour, Reviews: 1},
		}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserCards = repo.dueCards

	cmd, handled := model.updateBrowserKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !handled || cmd == nil {
		t.Fatal("browser enter should load selected card history")
	}
	model.Update(cmd())

	view := model.renderBrowser()
	if !strings.Contains(view, "Review History: der Apfel") || !strings.Contains(view, "hard") {
		t.Fatalf("browser history missing: %s", view)
	}
}

func TestWindowResizeChangesViewShape(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
	model = updated.(*Model)
	if model.breakpoint != BreakpointCompact {
		t.Fatalf("breakpoint = %s, want compact", model.breakpoint)
	}
	if !strings.Contains(model.View().Content, "Dashboard") {
		t.Fatal("view should render dashboard")
	}
}

func TestHitboxActivation(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.View()
	model.activateHitbox("tab-review")
	if model.activeView != ViewReview {
		t.Fatalf("active view = %s, want review", model.activeView)
	}
}

func TestTabNavigationReturnsViewLoadCommand(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewAI

	updated, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	model = updated.(*Model)
	if model.activeView != ViewSettings {
		t.Fatalf("active view after first tab = %s, want settings", model.activeView)
	}
	if cmd != nil {
		t.Fatal("settings view should not need a load command")
	}

	updated, cmd = model.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	model = updated.(*Model)
	if model.activeView != ViewBrowser {
		t.Fatalf("active view after second tab = %s, want browser", model.activeView)
	}
	if cmd == nil {
		t.Fatal("browser view should return its load command when reached by tab")
	}
}

func TestReviewFlow(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", Prompt: "P2", Answer: "A2"},
		},
	}
	model := NewModel(repo, &mockScheduler{})

	// Load cards
	model.Update(dueCardsMsg(repo.dueCards))

	if len(model.dueCards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(model.dueCards))
	}

	model.activeView = ViewReview
	if model.revealed {
		t.Fatal("should not be revealed initially")
	}

	// Reveal using space
	model.Update(tea.KeyPressMsg{Code: ' '})
	if !model.revealed {
		t.Fatal("should be revealed after space")
	}

	// Grade (in test we don't run the async command, we just check it exists)
	cmd := model.gradeCard(core.GradeGood)
	if cmd == nil {
		t.Fatal("gradeCard should return a command")
	}

	if !strings.Contains(model.status, "Grade: Good") {
		t.Fatalf("status should contain grade, got: %s", model.status)
	}

	// Simulate re-loading with one less card
	repo.dueCards = repo.dueCards[1:]
	model.Update(dueCardsMsg(repo.dueCards))

	if len(model.dueCards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(model.dueCards))
	}
	if model.cursor != 0 {
		t.Fatalf("cursor should be 0, got %d", model.cursor)
	}
	if model.revealed {
		t.Fatal("should be reset to unrevealed")
	}
}

func TestReviewBookmarkToggleUpdatesCardAndRendering(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"}},
		decks:    []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	if !strings.Contains(model.renderReview(0, 0), "Bookmark: off") {
		t.Fatalf("initial review should show bookmark off: %s", model.renderReview(0, 0))
	}
	cmd := model.toggleBookmark()
	if cmd == nil {
		t.Fatal("toggleBookmark should return command")
	}
	model.Update(cmd())
	if !model.dueCards[0].Bookmarked {
		t.Fatal("card should be bookmarked after toggle")
	}
	if !strings.Contains(model.renderReview(0, 0), "Bookmark: on") {
		t.Fatalf("review should show bookmark on: %s", model.renderReview(0, 0))
	}
}

func TestUndoLastReviewReloadsDueCards(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-1", Prompt: "P2", Answer: "A2"},
		},
		reviews: []core.ReviewResult{{CardID: "c1"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards[:1]))
	model.activeView = ViewReview
	model.lastReviewedCardID = "c1"

	cmd := model.undoLastReview()
	if cmd == nil {
		t.Fatal("undoLastReview should return command")
	}
	model.Update(cmd())
	if model.lastReviewedCardID != "" {
		t.Fatalf("lastReviewedCardID = %q, want empty", model.lastReviewedCardID)
	}
	if len(model.dueCards) != 2 {
		t.Fatalf("due cards after undo = %d, want 2", len(model.dueCards))
	}
	if !strings.Contains(model.status, "Last review undone") {
		t.Fatalf("status = %q, want undo confirmation", model.status)
	}
}

func TestDeckSwitchingFiltersDueCards(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{
			{ID: "deck-1", Name: "Deck One"},
			{ID: "deck-2", Name: "Deck Two"},
		},
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-2", Prompt: "P2", Answer: "A2"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))

	if model.deck.ID != "deck-1" || len(model.dueCards) != 1 || model.dueCards[0].ID != "c1" {
		t.Fatalf("initial deck/filter = %s/%v", model.deck.ID, model.dueCards)
	}

	model.Update(tea.KeyPressMsg{Code: ']'})
	if model.deck.ID != "deck-2" || len(model.dueCards) != 1 || model.dueCards[0].ID != "c2" {
		t.Fatalf("after switch deck/filter = %s/%v", model.deck.ID, model.dueCards)
	}
}

func TestAIDraftApprovalPersistsAndReloadsDueCards(t *testing.T) {
	deck := core.Deck{ID: "deck-1", Name: "Deck One"}
	provider := ai.OfflineProvider{}
	repo := &mockRepo{decks: []core.Deck{deck}}
	model := NewModelWithAI(repo, &mockScheduler{}, provider)
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewAI
	model.aiInput = "der Tee"

	generate := model.generateDrafts()
	msg := generate()
	model.Update(msg)
	if len(model.drafts) != 1 {
		t.Fatalf("drafts = %d, want 1", len(model.drafts))
	}

	approve := model.approveDraft()
	if approve == nil {
		t.Fatal("approveDraft returned nil command")
	}
	msg = approve()
	model.Update(msg)

	if repo.upsertedDeck.ID != "deck-1" {
		t.Fatalf("upserted deck = %q, want deck-1", repo.upsertedDeck.ID)
	}
	if len(repo.upsertedDeck.Notes) != 1 || repo.upsertedDeck.Notes[0].Front != "der Tee" {
		t.Fatalf("upserted notes = %#v", repo.upsertedDeck.Notes)
	}
	if len(model.drafts) != 0 {
		t.Fatalf("drafts after approve = %d, want 0", len(model.drafts))
	}
	if len(model.dueCards) == 0 {
		t.Fatal("approved draft cards should reload into due cards")
	}
	if !strings.Contains(model.status, "Draft approved") {
		t.Fatalf("status = %q, want draft approved", model.status)
	}
}

func TestAIDraftMouseApproval(t *testing.T) {
	note := core.Note{
		ID:     "n1",
		DeckID: "deck-1",
		Front:  "der Tee",
		Back:   "tea",
		Cards: []core.Card{{
			ID:     "c1",
			NoteID: "n1",
			DeckID: "deck-1",
			Kind:   core.CardKindFlashcard,
			Prompt: "der Tee",
			Answer: "tea",
		}},
	}
	repo := &mockRepo{decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}}}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.FakeProvider{Drafts: []ai.Draft{{Note: note}}})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewAI
	model.Update(draftsMsg([]ai.Draft{{Note: note}}))
	model.View()

	cmd := model.activateHitbox("draft-approve")
	if cmd == nil {
		t.Fatal("draft approval hitbox should return command")
	}
	model.Update(cmd())
	if repo.upsertedDeck.ID != "deck-1" {
		t.Fatalf("upserted deck = %q, want deck-1", repo.upsertedDeck.ID)
	}
}

func TestImportTSVUpsertsDeckAndRefreshesDueCards(t *testing.T) {
	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "import.tsv")
	exportPath := filepath.Join(tmp, "export.tsv")
	input := "#separator:tab\nid-1\tdie Bahn\ttrain\t\ttransport\tImported A1\tBasic\n"
	if err := os.WriteFile(importPath, []byte(input), 0o644); err != nil {
		t.Fatalf("write import fixture: %v", err)
	}

	repo := &mockRepo{}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProvider: ai.OfflineProvider{},
		ImportPath: importPath,
		ExportPath: exportPath,
	})
	model.Update(decksMsg([]core.Deck{{ID: "a1-survival", Name: "German A1 Survival"}}))
	model.activeView = ViewImport

	cmd, handled := model.updateImportKey(tea.KeyPressMsg{Code: 'i'})
	if !handled || cmd == nil {
		t.Fatal("import key should return command")
	}
	model.Update(cmd())

	if repo.upsertedDeck.ID != "Imported A1" {
		t.Fatalf("upserted deck = %q, want Imported A1", repo.upsertedDeck.ID)
	}
	if len(model.dueCards) != 1 || model.dueCards[0].Prompt != "die Bahn" {
		t.Fatalf("due cards = %#v", model.dueCards)
	}
	if !strings.Contains(model.status, "Imported 1 notes") {
		t.Fatalf("status = %q, want import success", model.status)
	}
}

func TestExportTSVWritesSelectedDeck(t *testing.T) {
	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "import.tsv")
	exportPath := filepath.Join(tmp, "export.tsv")
	deck := core.Deck{
		ID:   "deck-1",
		Name: "Deck One",
		Notes: []core.Note{{
			ID:     "n1",
			DeckID: "deck-1",
			Front:  "die Bahn",
			Back:   "train",
			Cards: []core.Card{{
				ID:     "n1:front",
				NoteID: "n1",
				DeckID: "deck-1",
				Kind:   core.CardKindFlashcard,
				Prompt: "die Bahn",
				Answer: "train",
			}},
		}},
	}
	repo := &mockRepo{decks: []core.Deck{deck}}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProvider: ai.OfflineProvider{},
		ImportPath: importPath,
		ExportPath: exportPath,
	})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewImport

	cmd, handled := model.updateImportKey(tea.KeyPressMsg{Code: 'x'})
	if !handled || cmd == nil {
		t.Fatal("export key should return command")
	}
	model.Update(cmd())

	raw, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if !strings.Contains(string(raw), "die Bahn\ttrain") {
		t.Fatalf("export did not include note: %s", string(raw))
	}
	if !strings.Contains(model.status, "Exported 1 notes") {
		t.Fatalf("status = %q, want export success", model.status)
	}
}

func TestAIGenerateErrorFromProvider(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.FakeProvider{Err: context.Canceled})
	model.Update(decksMsg([]core.Deck{{ID: "deck-1", Name: "Deck One"}}))
	model.activeView = ViewAI

	cmd := model.generateDrafts()
	msg := cmd()
	model.Update(msg)
	if !strings.Contains(model.status, "context canceled") {
		t.Fatalf("status = %q, want provider error", model.status)
	}
}

func TestDecksViewNavigationAndSelection(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{
			{ID: "deck-1", Name: "Deck One", TotalCards: 10, DueCards: 5, ReviewsToday: 3, SuccessRate: 0.75},
			{ID: "deck-2", Name: "Deck Two", TotalCards: 20, DueCards: 2},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))

	// Switch to Decks view
	model.Update(tea.KeyPressMsg{Code: '2'})
	if model.activeView != ViewDecks {
		t.Fatalf("active view = %s, want decks", model.activeView)
	}

	// Initial cursor should be at index 0 (matching deckIndex)
	if model.deckCursor != 0 {
		t.Fatalf("deckCursor = %d, want 0", model.deckCursor)
	}
	view := model.renderDecks(0, 0)
	if !strings.Contains(view, "today 3") || !strings.Contains(view, "75% success") {
		t.Fatalf("deck view missing progress metrics: %s", view)
	}

	// Move cursor down
	model.Update(tea.KeyPressMsg{Code: 'j'})
	if model.deckCursor != 1 {
		t.Fatalf("deckCursor after j = %d, want 1", model.deckCursor)
	}

	// Select deck-2
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter
	if model.deck.ID != "deck-2" {
		t.Fatalf("active deck = %s, want deck-2", model.deck.ID)
	}
	if model.activeView != ViewDashboard {
		t.Fatalf("view after select = %s, want dashboard", model.activeView)
	}

	// Check rendering contains stats
	view = model.renderDecks(0, 0)
	if !strings.Contains(view, "Deck Two (20 total, 2 due") {
		t.Fatalf("decks view rendering missing stats: %s", view)
	}
}

func TestSettingsProviderSwitching(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings
	if model.aiProviderName != "offline" {
		t.Fatalf("initial provider = %q, want offline", model.aiProviderName)
	}

	// Cycle: offline -> template
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter
	if model.aiProviderName != "template" {
		t.Fatalf("provider after toggle 1 = %q, want template", model.aiProviderName)
	}
	if _, ok := model.aiProvider.(ai.TemplateProvider); !ok {
		t.Fatal("aiProvider should be TemplateProvider")
	}

	// Cycle: template -> disabled
	model.Update(tea.KeyPressMsg{Code: '\r'})
	if model.aiProviderName != "disabled" {
		t.Fatalf("provider after toggle 2 = %q, want disabled", model.aiProviderName)
	}
	if model.aiProvider != nil {
		t.Fatal("aiProvider should be nil when disabled")
	}

	// Cycle: disabled -> offline
	model.Update(tea.KeyPressMsg{Code: '\r'})
	if model.aiProviderName != "offline" {
		t.Fatalf("provider after toggle 3 = %q, want offline", model.aiProviderName)
	}
	if _, ok := model.aiProvider.(ai.OfflineProvider); !ok {
		t.Fatal("aiProvider should be OfflineProvider")
	}
}

func TestSettingsTemplateEditing(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings
	model.settingsCursor = 1 // Front Template

	// Start editing
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter to edit
	if !model.editingTemplate {
		t.Fatal("should be in editing mode")
	}

	// Backspace to clear (partially)
	model.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})

	// Add text
	model.Update(tea.KeyPressMsg{Code: 'X'})

	// Save
	model.Update(tea.KeyPressMsg{Code: '\r'})
	if model.editingTemplate {
		t.Fatal("should not be in editing mode after save")
	}

	want := "{{.Topic}X" // default is {{.Topic}}, minus one char plus X
	if model.aiTemplates["front"] != want {
		t.Fatalf("template = %q, want %q", model.aiTemplates["front"], want)
	}
}

func TestAIDraftWithTemplateProvider(t *testing.T) {
	repo := &mockRepo{decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}}}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProviderName: "template",
		AITemplates: map[string]string{
			"front": "Prefix: {{.Topic}}",
		},
	})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewAI
	model.aiInput = "test"

	generate := model.generateDrafts()
	msg := generate()
	model.Update(msg)

	if len(model.drafts) != 1 {
		t.Fatalf("drafts = %d, want 1", len(model.drafts))
	}
	if model.drafts[0].Note.Front != "Prefix: test" {
		t.Fatalf("draft front = %q, want %q", model.drafts[0].Note.Front, "Prefix: test")
	}
}

func TestBookmarkFilterToggle(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1", Bookmarked: true},
			{ID: "c2", DeckID: "deck-1", Prompt: "P2", Answer: "A2", Bookmarked: false},
			{ID: "c3", DeckID: "deck-1", Prompt: "P3", Answer: "A3", Bookmarked: true},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	if model.bookmarkFilter {
		t.Fatal("bookmark filter should be off initially")
	}

	cmd := model.toggleBookmarkFilter()
	if cmd == nil {
		t.Fatal("toggleBookmarkFilter should return command")
	}
	msg := cmd()
	model.Update(msg)

	if !model.bookmarkFilter {
		t.Fatal("bookmark filter should be on after toggle")
	}
	if len(model.dueCards) != 2 {
		t.Fatalf("due cards = %d, want 2 bookmarked", len(model.dueCards))
	}
	if !strings.Contains(model.status, "bookmarked") {
		t.Fatalf("status = %q, want bookmarked message", model.status)
	}
}

func TestBookmarkFilterToggleOff(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-1", Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.bookmarkFilter = true

	cmd := model.toggleBookmarkFilter()
	if cmd == nil {
		t.Fatal("toggleBookmarkFilter should return command")
	}
	msg := cmd()
	model.Update(msg)

	if model.bookmarkFilter {
		t.Fatal("bookmark filter should be off after second toggle")
	}
	if len(model.dueCards) != 2 {
		t.Fatalf("due cards = %d, want 2", len(model.dueCards))
	}
}

func TestReviewViewShowsLeechIndicator(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1", Leech: true},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "LEECH") {
		t.Fatalf("review view should show LEECH indicator: %s", view)
	}
}

func TestReviewViewShowsBookmarkFilterBanner(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1", Bookmarked: true},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.bookmarkFilter = true

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "(Bookmarked)") {
		t.Fatalf("review view should show bookmark filter banner: %s", view)
	}
}

func TestSuspendCardRefreshesQueueAndStats(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindFlashcard, Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-1", Kind: core.CardKindFlashcard, Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One", TotalCards: 2, DueCards: 2}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	cmd := model.suspendCard()
	if cmd == nil {
		t.Fatal("suspendCard should return command")
	}
	msg := cmd()
	updated, _ := model.Update(msg)
	model = updated.(*Model)

	if len(model.dueCards) != 1 {
		t.Fatalf("due cards = %d, want 1", len(model.dueCards))
	}
	if model.dueCards[0].ID != "c2" {
		t.Fatalf("remaining card = %s, want c2", model.dueCards[0].ID)
	}
	if model.stats.SuspendedCards != 1 {
		t.Fatalf("suspended cards = %d, want 1", model.stats.SuspendedCards)
	}
	if !strings.Contains(model.status, "suspended") {
		t.Fatalf("status = %q, want suspended", model.status)
	}
}

func TestSettingsDailyGoalAdjustsAndRenders(t *testing.T) {
	repo := &mockRepo{dailyGoal: 10}
	model := NewModel(repo, &mockScheduler{})
	model.Update(statsMsg(core.Statistics{DailyGoal: 10, Grades: map[core.ReviewGrade]int{}}))
	model.activeView = ViewSettings
	model.settingsCursor = 1

	if !strings.Contains(model.renderSettings(0, 0), "Daily Goal: 10") {
		t.Fatalf("settings should show daily goal: %s", model.renderSettings(0, 0))
	}

	updated, cmd := model.Update(tea.KeyPressMsg{Code: '+'})
	model = updated.(*Model)
	if cmd == nil {
		t.Fatal("daily goal increase should return command")
	}
	msg := cmd()
	updated, _ = model.Update(msg)
	model = updated.(*Model)

	if model.stats.DailyGoal != 11 {
		t.Fatalf("daily goal = %d, want 11", model.stats.DailyGoal)
	}
	if !strings.Contains(model.status, "11") {
		t.Fatalf("status = %q, want goal value", model.status)
	}

	model.stats.DailyGoal = 1
	updated, cmd = model.Update(tea.KeyPressMsg{Code: '-'})
	model = updated.(*Model)
	if cmd == nil {
		t.Fatal("daily goal decrease should return command")
	}
	msg = cmd()
	updated, _ = model.Update(msg)
	model = updated.(*Model)
	if model.stats.DailyGoal != 1 {
		t.Fatalf("daily goal floor = %d, want 1", model.stats.DailyGoal)
	}
}

func TestDashboardShowsBookmarkedDueAndLeech(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1", Bookmarked: true, Leech: true},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.stats = core.Statistics{
		BookmarkedCards: 1,
		BookmarkedDue:   1,
		LeechCards:      1,
		ReviewsToday:    3,
		DailyGoal:       10,
		CurrentStreak:   3,
	}

	dash := model.renderActiveViewPlain(0, 0)
	if !strings.Contains(dash, "1 due)") {
		t.Fatalf("dashboard should show bookmarked due count: %s", dash)
	}
	if !strings.Contains(dash, "Leech:") || !strings.Contains(dash, "1") {
		t.Fatalf("dashboard should show leech count: %s", dash)
	}
	if !strings.Contains(dash, "Streak:") || !strings.Contains(dash, "3 days 🔥") {
		t.Fatalf("dashboard should show current streak with fire emoji: %s", dash)
	}
}

func TestStatisticsShowsLeechAndBookmarkedDue(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"}},
		decks:    []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.stats = core.Statistics{
		TotalCards:      20,
		NewCards:        5,
		YoungCards:      10,
		MatureCards:     5,
		BookmarkedCards: 3,
		BookmarkedDue:   2,
		LeechCards:      1,
		TotalReviews:    50,
		ReviewsToday:    5,
		DailyGoal:       10,
		CurrentStreak:   3,
		SuccessRate:     0.85,
		Grades:          map[core.ReviewGrade]int{core.GradeGood: 40, core.GradeAgain: 10},
	}

	view := model.renderStatistics()
	if !strings.Contains(view, "Bookmarked:  3 (2 due)") {
		t.Fatalf("statistics should show bookmarked due: %s", view)
	}
	if !strings.Contains(view, "Leech:       1") {
		t.Fatalf("statistics should show leech count: %s", view)
	}
}

func TestMCQCardRendersChoicesAfterReveal(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Was ist 'der Apfel'?", Answer: "the apple", Choices: []string{"the apple", "the banana", "the pear", "the grape"}},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "reveal") || !strings.Contains(view, "choices") {
		t.Fatalf("MCQ view should show reveal prompt: %s", view)
	}

	model.revealed = true
	view = model.renderReview(0, 0)
	if !strings.Contains(view, "1:") {
		t.Fatalf("MCQ view should show choice 1 after reveal: %s", view)
	}
	if !strings.Contains(view, "2:") {
		t.Fatalf("MCQ view should show choice 2 after reveal: %s", view)
	}
	if !strings.Contains(view, "3:") {
		t.Fatalf("MCQ view should show choice 3 after reveal: %s", view)
	}
	if !strings.Contains(view, "4:") {
		t.Fatalf("MCQ view should show choice 4 after reveal: %s", view)
	}
}

func TestMCQChoiceSelectionCorrect(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Was ist 'der Apfel'?", Answer: "the apple", Choices: []string{"the apple", "the banana", "the pear", "the grape"}},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealed = true

	model.selectMCQChoice("1")
	if !model.mcqAnswered {
		t.Fatal("MCQ should be answered after choice selection")
	}
	if !model.mcqCorrect {
		t.Fatal("Choice 1 should be correct")
	}
	if model.mcqChoice != 0 {
		t.Fatalf("mcqChoice = %d, want 0", model.mcqChoice)
	}
}

func TestMCQChoiceSelectionIncorrect(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Was ist 'der Apfel'?", Answer: "the apple", Choices: []string{"the apple", "the banana", "the pear", "the grape"}},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealed = true

	model.selectMCQChoice("2")
	if !model.mcqAnswered {
		t.Fatal("MCQ should be answered after choice selection")
	}
	if model.mcqCorrect {
		t.Fatal("Choice 2 should be incorrect")
	}

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "Incorrect") {
		t.Fatalf("MCQ view should show incorrect feedback: %s", view)
	}
}

func TestMCQChoiceSelectionShowsCorrectAnswer(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Was ist 'der Apfel'?", Answer: "the apple", Choices: []string{"the apple", "the banana", "the pear", "the grape"}},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealed = true

	model.selectMCQChoice("3")

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "the apple") {
		t.Fatalf("MCQ view should show correct answer after wrong choice: %s", view)
	}
	if strings.Contains(view, "Correct") {
		t.Fatal("MCQ view should not show Correct for wrong choice")
	}
}

func TestMCQStateResetsOnCardNavigation(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Q1", Answer: "A1", Choices: []string{"A1", "B1"}},
			{ID: "c2", DeckID: "deck-1", Kind: core.CardKindMCQ, Prompt: "Q2", Answer: "A2", Choices: []string{"A2", "B2"}},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealed = true

	model.selectMCQChoice("1")
	if !model.mcqAnswered {
		t.Fatal("MCQ should be answered")
	}

	var updated tea.Model
	updated, _ = model.Update(tea.KeyPressMsg{Code: 'j'})
	model = updated.(*Model)
	if model.mcqAnswered {
		t.Fatal("MCQ state should be reset after navigation")
	}
	if model.mcqChoice != -1 {
		t.Fatalf("mcqChoice = %d, want -1", model.mcqChoice)
	}
}

func TestMCQFlashcardFlowUnchanged(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindFlashcard, Prompt: "P1", Answer: "A1"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	if model.revealed {
		t.Fatal("flashcard should not be revealed initially")
	}

	view := model.renderReview(0, 0)
	if strings.Contains(view, "1:") {
		t.Fatalf("flashcard view should not show MCQ choices: %s", view)
	}

	model.Update(tea.KeyPressMsg{Code: ' '})
	if !model.revealed {
		t.Fatalf("flashcard should be revealed after space, revealed=%v", model.revealed)
	}
}

func TestSessionStatsTracking(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-1", Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealed = true

	if model.sessionReviewed != 0 {
		t.Fatalf("sessionReviewed = %d, want 0", model.sessionReviewed)
	}

	// Simulate grade by sending reviewRecordedMsg
	model.Update(reviewRecordedMsg{cardID: "c1", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeAgain})
	if model.sessionReviewed != 1 {
		t.Fatalf("sessionReviewed = %d, want 1", model.sessionReviewed)
	}
	if model.sessionCorrect != 0 {
		t.Fatalf("sessionCorrect = %d, want 0", model.sessionCorrect)
	}

	model.Update(reviewRecordedMsg{cardID: "c2", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeGood})
	if model.sessionReviewed != 2 {
		t.Fatalf("sessionReviewed = %d, want 2", model.sessionReviewed)
	}
	if model.sessionCorrect != 1 {
		t.Fatalf("sessionCorrect = %d, want 1", model.sessionCorrect)
	}
}

func TestHelpOverlayToggle(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})

	if model.showHelp {
		t.Fatal("help should be hidden initially")
	}

	model.Update(tea.KeyPressMsg{Code: '?'})
	if !model.showHelp {
		t.Fatal("help should be shown after pressing ?")
	}

	view := model.View().Content
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Fatal("help overlay should be visible in view")
	}

	model.Update(tea.KeyPressMsg{Code: '?'})
	if model.showHelp {
		t.Fatal("help should be hidden after pressing ? again")
	}
}

func TestBrowserViewNavigation(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"},
			{ID: "c2", DeckID: "deck-1", Prompt: "Banana", Answer: "Banane"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))

	model.activeView = ViewBrowser
	model.browserDeckID = "deck-1"
	// Execute the cmd to load browser cards
	cmd := model.loadBrowserCards()
	if cmd != nil {
		msg := cmd()
		model.Update(msg)
	}

	if len(model.browserCards) == 0 {
		t.Fatal("browser should have cards")
	}

	model.updateBrowserKey(tea.KeyPressMsg{Code: 'j'})
	if model.browserCursor != 1 {
		t.Fatalf("browserCursor = %d, want 1", model.browserCursor)
	}

	model.updateBrowserKey(tea.KeyPressMsg{Code: 'k'})
	if model.browserCursor != 0 {
		t.Fatalf("browserCursor = %d, want 0", model.browserCursor)
	}
}

func TestBrowserSearchFilter(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"},
			{ID: "c2", DeckID: "deck-1", Prompt: "Banana", Answer: "Banane"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserDeckID = "deck-1"

	model.updateBrowserKey(tea.KeyPressMsg{Code: 'A'})
	model.updateBrowserKey(tea.KeyPressMsg{Code: 'p'})
	// Execute the cmd to load browser cards with search
	cmd := model.loadBrowserCards()
	if cmd != nil {
		msg := cmd()
		model.Update(msg)
	}

	if len(model.browserCards) != 1 {
		t.Fatalf("browserCards = %d, want 1 after search", len(model.browserCards))
	}
	if model.browserCards[0].Prompt != "Apple" {
		t.Fatalf("expected Apple, got %s", model.browserCards[0].Prompt)
	}
}

func TestBrowserDeckSwitchReloadsCards(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"},
			{ID: "c2", DeckID: "deck-2", Prompt: "Coffee", Answer: "Kaffee"},
		},
		decks: []core.Deck{
			{ID: "deck-1", Name: "Deck One"},
			{ID: "deck-2", Name: "Deck Two"},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	load := model.updateView(ViewBrowser)
	model.Update(load())

	if len(model.browserCards) != 1 || model.browserCards[0].ID != "c1" {
		t.Fatalf("initial browser cards = %#v, want deck-1 card", model.browserCards)
	}

	updated, cmd := model.Update(tea.KeyPressMsg{Code: ']'})
	model = updated.(*Model)
	if cmd == nil {
		t.Fatal("browser deck switch should reload cards")
	}
	model.Update(cmd())

	if model.deck.ID != "deck-2" {
		t.Fatalf("selected deck = %s, want deck-2", model.deck.ID)
	}
	if model.browserDeckID != "deck-2" {
		t.Fatalf("browserDeckID = %s, want deck-2", model.browserDeckID)
	}
	if len(model.browserCards) != 1 || model.browserCards[0].ID != "c2" {
		t.Fatalf("browser cards after deck switch = %#v, want deck-2 card", model.browserCards)
	}
}

func TestCramModeFiltering(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "Card 1", Bookmarked: true},
			{ID: "c2", Prompt: "Card 2", Suspended: true},
			{ID: "c3", Prompt: "Card 3", Leech: true},
			{ID: "c4", Prompt: "Card 4"},
		},
	}
	model := NewModel(repo, &mockScheduler{})

	// Mock cram cards loading
	model.Update(cramCardsMsg(repo.dueCards))

	// Test bookmarked filter
	model.cramType = "bookmarked"
	model.cramCards = nil
	model.Update(cramCardsMsg(repo.dueCards))
	if len(model.cramCards) != 1 {
		t.Fatalf("cramCards = %d, want 1 for bookmarked filter", len(model.cramCards))
	}
	if model.cramCards[0].ID != "c1" {
		t.Fatalf("expected c1, got %s", model.cramCards[0].ID)
	}

	// Test suspended filter
	model.cramType = "suspended"
	model.cramCards = nil
	model.Update(cramCardsMsg(repo.dueCards))
	if len(model.cramCards) != 1 {
		t.Fatalf("cramCards = %d, want 1 for suspended filter", len(model.cramCards))
	}
	if model.cramCards[0].ID != "c2" {
		t.Fatalf("expected c2, got %s", model.cramCards[0].ID)
	}

	// Test leech filter
	model.cramType = "leech"
	model.cramCards = nil
	model.Update(cramCardsMsg(repo.dueCards))
	if len(model.cramCards) != 1 {
		t.Fatalf("cramCards = %d, want 1 for leech filter", len(model.cramCards))
	}
	if model.cramCards[0].ID != "c3" {
		t.Fatalf("expected c3, got %s", model.cramCards[0].ID)
	}

	// Test flagged filter (bookmarked, suspended, or leech)
	model.cramType = "flagged"
	model.cramCards = nil
	model.Update(cramCardsMsg(repo.dueCards))
	if len(model.cramCards) != 3 {
		t.Fatalf("cramCards = %d, want 3 for flagged filter", len(model.cramCards))
	}

	// Test all filter
	model.cramType = "all"
	model.cramCards = nil
	model.Update(cramCardsMsg(repo.dueCards))
	if len(model.cramCards) != 4 {
		t.Fatalf("cramCards = %d, want 4 for all filter", len(model.cramCards))
	}
}

func TestCramReviewFlow(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "P1", Answer: "A1", Bookmarked: true},
			{ID: "c2", Prompt: "P2", Answer: "A2", Bookmarked: true},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewCram
	model.cramCards = repo.dueCards
	model.cramType = "bookmarked"

	if model.cramActive {
		t.Fatal("should not be active initially")
	}

	// Start review with enter
	model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !model.cramActive {
		t.Fatal("should be active after enter")
	}
	if model.cramRevealed {
		t.Fatal("should not be revealed initially")
	}

	// Reveal with space
	model.Update(tea.KeyPressMsg{Code: ' '})
	if !model.cramRevealed {
		t.Fatal("should be revealed after space")
	}

	// Grade correct with 'g'
	model.Update(tea.KeyPressMsg{Code: 'g'})
	if model.cramActive {
		t.Fatal("should not be active after grading")
	}
	if model.cramReviewed != 1 || model.cramCorrect != 1 {
		t.Fatalf("expected 1 reviewed/correct, got %d/%d", model.cramReviewed, model.cramCorrect)
	}
	if model.cramCursor != 1 {
		t.Fatalf("expected cursor 1, got %d", model.cramCursor)
	}

	// Start next card
	model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model.Update(tea.KeyPressMsg{Code: ' '})

	// Grade again with 'a'
	model.Update(tea.KeyPressMsg{Code: 'a'})
	if model.cramReviewed != 2 || model.cramCorrect != 1 {
		t.Fatalf("expected 2 reviewed, 1 correct, got %d/%d", model.cramReviewed, model.cramCorrect)
	}
	if model.cramCursor != 0 {
		t.Fatalf("expected cursor 0 (wrapped), got %d", model.cramCursor)
	}
}
