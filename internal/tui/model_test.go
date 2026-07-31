package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func executeCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		var msgs []tea.Msg
		for _, c := range batch {
			msgs = append(msgs, executeCmd(c)...)
		}
		return msgs
	}
	return []tea.Msg{msg}
}

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
	model.activateHitboxByID("tab-review")
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

func TestReviewGrammarHintToggle(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "P1", Answer: "A1"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	if model.showGrammarHint {
		t.Fatal("grammar hint should be hidden initially")
	}

	// Press G to toggle hint
	model.Update(tea.KeyPressMsg{Code: 'G'})
	if !model.showGrammarHint {
		t.Fatal("grammar hint should be shown after G")
	}
	if model.grammarHint == nil {
		t.Fatal("grammarHint should not be nil")
	}

	// Press G to toggle off
	model.Update(tea.KeyPressMsg{Code: 'G'})
	if model.showGrammarHint {
		t.Fatal("grammar hint should be hidden after G again")
	}
	if model.grammarHint != nil {
		t.Fatal("grammarHint should be nil when hidden")
	}
}

func TestReviewFlow(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})

	// Load cards
	model.Update(dueCardsMsg(repo.dueCards))

	if len(model.dueCards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(model.dueCards))
	}

	model.activeView = ViewReview
	if model.revealState != RevealIdle {
		t.Fatal("should not be revealed initially")
	}

	// Reveal using space
	model.Update(tea.KeyPressMsg{Code: ' '})
	if model.revealState != RevealFlipping {
		t.Fatal("should be flipping after space")
	}

	// Grade should be blocked while flip animation is in progress.
	if cmd := model.gradeCard(core.GradeGood); cmd != nil {
		t.Fatal("gradeCard should return nil while flip is still animating")
	}

	// Complete flip and reveal
	model.revealState = RevealRevealing
	model.revealProgress = 0
	// Grade after reveal.
	model.revealState = RevealRevealed
	model.revealProgress = 100
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
	if model.revealState != RevealIdle {
		t.Fatal("should be reset to unrevealed")
	}
}

func TestReviewExtraRevealKeyDoesNotBlockGrading(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", Prompt: "P2", Answer: "A2"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview
	model.revealState = RevealRevealed
	model.revealProgress = 100

	updated, cmd := model.Update(tea.KeyPressMsg{Code: ' '})
	model = updated.(*Model)
	if cmd != nil {
		t.Fatal("extra reveal key on revealed card should not grade")
	}
	if model.gradingInProgress {
		t.Fatal("extra reveal key should not block grade keys")
	}

	updated, cmd = model.Update(tea.KeyPressMsg{Code: 'e'})
	model = updated.(*Model)
	if cmd == nil {
		t.Fatal("easy grade should return a command after an extra reveal key")
	}
	if !strings.Contains(model.status, "Grade: Easy") {
		t.Fatalf("status should contain easy grade, got: %s", model.status)
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
	msgs := executeCmd(generate)
	for _, msg := range msgs {
		model.Update(msg)
	}
	if len(model.drafts) != 1 {
		t.Fatalf("drafts = %d, want 1", len(model.drafts))
	}

	approve := model.approveDraft()
	if approve == nil {
		t.Fatal("approveDraft returned nil command")
	}
	msg := approve()
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
	if !strings.Contains(model.status, "Draft saved") {
		t.Fatalf("status = %q, want Draft saved", model.status)
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

	cmd := model.activateHitboxByID("draft-approve-0")
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

	cmd, handled := model.importScreen.HandleKey(model, tea.KeyPressMsg{Code: 'i'})
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
	for _, note := range deck.Notes {
		repo.dueCards = append(repo.dueCards, note.Cards...)
	}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProvider: ai.OfflineProvider{},
		ImportPath: importPath,
		ExportPath: exportPath,
	})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewImport

	cmd, handled := model.importScreen.HandleKey(model, tea.KeyPressMsg{Code: 'x'})
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

func TestExportTSVAppliesStatusFilter(t *testing.T) {
	tmp := t.TempDir()
	exportPath := filepath.Join(tmp, "learning.tsv")
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "n1:front", NoteID: "n1", DeckID: "deck-1", Prompt: "jung", Answer: "young", Mature: false},
			{ID: "n2:front", NoteID: "n2", DeckID: "deck-1", Prompt: "reif", Answer: "mature", Mature: true},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProvider: ai.OfflineProvider{},
		ExportPath: exportPath,
	})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewImport
	model.exportFilter = "Mature"

	cmd := model.exportTSV()
	if cmd == nil {
		t.Fatal("exportTSV should return command")
	}
	model.Update(cmd())

	raw, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	out := string(raw)
	if strings.Contains(out, "jung\tyoung") {
		t.Fatalf("mature-only TSV export included learning card: %s", out)
	}
	if !strings.Contains(out, "reif\tmature") {
		t.Fatalf("mature-only TSV export missing mature card: %s", out)
	}
}

func TestAIGenerateErrorFromProvider(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.FakeProvider{Err: context.Canceled})
	model.Update(decksMsg([]core.Deck{{ID: "deck-1", Name: "Deck One"}}))
	model.activeView = ViewAI

	cmd := model.generateDrafts()
	msgs := executeCmd(cmd)
	for _, msg := range msgs {
		model.Update(msg)
	}
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
	layout := viewportLayout{Width: 82, Height: 24, X: 0, Y: 0}
	view := ansi.Strip(model.renderDecks(layout))
	if !strings.Contains(view, "5D") || !strings.Contains(view, "10T") {
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
	view = ansi.Strip(model.renderDecks(layout))
	if !strings.Contains(view, "Deck Two") || !strings.Contains(view, "D") || !strings.Contains(view, "T") {
		t.Fatalf("decks view rendering missing stats: %s", view)
	}
}

func TestSettingsProviderSwitching(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings
	if model.aiProviderName != "offline" {
		t.Fatalf("initial provider = %q, want offline", model.aiProviderName)
	}

	// New cycle: disabled -> offline -> template -> openai -> anthropic -> disabled
	steps := []struct {
		want     string
		wantType string // "nil", "template", "offline", "openai", "anthropic"
	}{
		{"template", "template"},
		{"openai", "openai"},
		{"anthropic", "anthropic"},
		{"disabled", "nil"},
		{"offline", "offline"},
	}
	for i, step := range steps {
		model.Update(tea.KeyPressMsg{Code: '\r'})
		if model.aiProviderName != step.want {
			t.Fatalf("step %d provider = %q, want %q", i+1, model.aiProviderName, step.want)
		}
		switch step.wantType {
		case "nil":
			if model.aiProvider != nil {
				t.Fatalf("step %d provider should be nil", i+1)
			}
		case "template":
			if _, ok := model.aiProvider.(ai.TemplateProvider); !ok {
				t.Fatalf("step %d provider should be TemplateProvider", i+1)
			}
		case "offline":
			if _, ok := model.aiProvider.(ai.OfflineProvider); !ok {
				t.Fatalf("step %d provider should be OfflineProvider", i+1)
			}
		case "openai":
			if _, ok := model.aiProvider.(ai.OpenAIProvider); !ok {
				t.Fatalf("step %d provider should be OpenAIProvider", i+1)
			}
		case "anthropic":
			if _, ok := model.aiProvider.(ai.AnthropicProvider); !ok {
				t.Fatalf("step %d provider should be AnthropicProvider", i+1)
			}
		}
	}
}

func TestSettingsTemplateEditing(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings
	model.settingsCursor = 2 // Front Template (was 1)

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
	activeSet := model.aiTemplateSets[model.aiTemplateIndex]
	if model.aiTemplates[activeSet]["front"] != want {
		t.Fatalf("template = %q, want %q", model.aiTemplates[activeSet]["front"], want)
	}
}

func TestSettingsBaseURLEditing(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings

	// 1. Edit OpenAI Base URL (index 10)
	model.settingsCursor = 10
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter to edit
	if model.editingSecretKey != "base_url" || model.editingSecretProvider != "openai" {
		t.Fatalf("expected editing openai base_url, got key=%q provider=%q", model.editingSecretKey, model.editingSecretProvider)
	}

	// Type a URL
	for _, char := range "http://localhost:11434" {
		model.Update(tea.KeyPressMsg{Code: char})
	}

	// Save
	model.Update(tea.KeyPressMsg{Code: '\r'})
	if model.editingSecretKey != "" {
		t.Fatal("should not be editing secret after save")
	}

	if model.aiSecrets.OpenAI.BaseURL != "http://localhost:11434" {
		t.Fatalf("expected openai BaseURL to be http://localhost:11434, got %q", model.aiSecrets.OpenAI.BaseURL)
	}

	// 2. Edit Anthropic Base URL (index 13)
	model.settingsCursor = 13
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter to edit
	if model.editingSecretKey != "base_url" || model.editingSecretProvider != "anthropic" {
		t.Fatalf("expected editing anthropic base_url, got key=%q provider=%q", model.editingSecretKey, model.editingSecretProvider)
	}

	// Type a URL
	for _, char := range "http://localhost:9000" {
		model.Update(tea.KeyPressMsg{Code: char})
	}

	// Save
	model.Update(tea.KeyPressMsg{Code: '\r'})
	if model.editingSecretKey != "" {
		t.Fatal("should not be editing secret after save")
	}

	if model.aiSecrets.Anthropic.BaseURL != "http://localhost:9000" {
		t.Fatalf("expected anthropic BaseURL to be http://localhost:9000, got %q", model.aiSecrets.Anthropic.BaseURL)
	}
}

func TestTrimLastRuneHandlesUnicode(t *testing.T) {
	if got := trimLastRune("Käse"); got != "Käs" {
		t.Fatalf("trimLastRune ascii+umlaut = %q, want %q", got, "Käs")
	}
	if got := trimLastRune("🙂a"); got != "🙂" {
		t.Fatalf("trimLastRune emoji tail = %q, want %q", got, "🙂")
	}
	if got := trimLastRune("🙂"); got != "" {
		t.Fatalf("trimLastRune single emoji = %q, want empty", got)
	}
}

func TestSinglePrintableInputAcceptsSingleRune(t *testing.T) {
	if ch, ok := singlePrintableInput("ä"); !ok || ch != "ä" {
		t.Fatalf("singlePrintableInput(ä) = %q,%t; want ä,true", ch, ok)
	}
	if _, ok := singlePrintableInput("\t"); ok {
		t.Fatal("singlePrintableInput(tab) should be rejected")
	}
	if _, ok := singlePrintableInput(""); ok {
		t.Fatal("singlePrintableInput(empty) should be rejected")
	}
}

func TestAIGenerateWhenProviderDisabledSetsStatus(t *testing.T) {
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		AIProviderName: "disabled",
		AIProvider:     nil,
	})
	model.activeView = ViewAI
	model.aiInput = "der Kaffee"

	cmd := model.startDrafting()
	if cmd != nil {
		t.Fatal("startDrafting should return nil when AI provider is disabled")
	}
	if model.drafting {
		t.Fatal("drafting should remain false when provider is disabled")
	}
	if !strings.Contains(model.status, "disabled") {
		t.Fatalf("status = %q, want disabled guidance", model.status)
	}
}

func TestAIGenerateEmptyTopicDoesNotCallProvider(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.FakeProvider{Err: errors.New("provider should not be called")})
	model.activeView = ViewAI
	model.aiInput = "   "

	cmd := model.startDrafting()
	if cmd != nil {
		t.Fatal("startDrafting should return nil for an empty topic")
	}
	if model.drafting {
		t.Fatal("drafting should remain false for an empty topic")
	}
	if !strings.Contains(model.status, "Enter a topic") {
		t.Fatalf("status = %q, want topic guidance", model.status)
	}
}

func TestAIEscapeClearsStarterTopic(t *testing.T) {
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		AIProviderName: "disabled",
		AIProvider:     nil,
	})
	model.activeView = ViewAI
	model.aiInput = "der Kaffee"

	cmd, handled := model.updateAIKey(tea.KeyPressMsg{Code: tea.KeyEsc})
	if !handled {
		t.Fatal("expected escape to be handled in AI view")
	}
	if cmd != nil {
		t.Fatal("expected no command when clearing AI topic")
	}
	if model.aiInput != "" {
		t.Fatalf("expected escape to clear AI topic, got %q", model.aiInput)
	}

	cmd = model.startDrafting()
	if cmd != nil {
		t.Fatal("expected empty topic guard to return no command")
	}
	if !strings.Contains(model.status, "Enter a topic") {
		t.Fatalf("status = %q, want topic guidance", model.status)
	}
}

func TestAIDraftCursorGuardsInvalidSelection(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.OfflineProvider{})
	model.activeView = ViewAI
	model.drafts = []ai.Draft{{Note: core.Note{ID: "n1", DeckID: "d1", Front: "eins", Back: "one"}}}
	model.draftCursor = -1

	if cmd := model.approveDraft(); cmd != nil {
		t.Fatal("approveDraft should not return a command for a negative cursor")
	}
	if !strings.Contains(model.status, "No draft selected") {
		t.Fatalf("status = %q, want invalid selection guidance", model.status)
	}

	model.draftCursor = 4
	if cmd := model.discardDraft(); cmd != nil {
		t.Fatal("discardDraft should not return a command for an out-of-range cursor")
	}
	if len(model.drafts) != 1 {
		t.Fatalf("drafts length changed for invalid cursor: %d", len(model.drafts))
	}
}

func TestDeckLimitEditingClampsFilteredCursor(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDecks
	model.decks = []core.Deck{{ID: "a", Name: "Alpha"}, {ID: "b", Name: "Beta"}}
	model.deckCursor = 99
	model.editingDeckLimits = true

	_, handled := model.updateDecksKey(tea.KeyPressMsg{Code: '+'})
	if !handled {
		t.Fatal("deck limit edit key should be handled")
	}
	if model.deckCursor != 1 {
		t.Fatalf("deckCursor = %d, want clamped last index", model.deckCursor)
	}
}

func TestNewModelWithOptionsTemplateProviderHandlesEmptyTemplates(t *testing.T) {
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		AIProviderName: "template",
		AITemplates:    map[string]map[string]string{},
	})
	if model.aiProvider == nil {
		t.Fatal("aiProvider should be initialized for template mode")
	}
	if model.currentAITemplateSet() == "" {
		t.Fatal("template sets should be initialized with defaults when empty map is provided")
	}
}

func TestDeckLimitLeftMovesCursorBackward(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One", NewCardsPerDay: 10, ReviewLimitPerDay: 50}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewDecks
	model.editingDeckLimits = true
	model.limitCursor = 1
	model.deckCursor = 0

	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	model = updated.(*Model)
	if model.limitCursor != 0 {
		t.Fatalf("limitCursor after left = %d, want 0", model.limitCursor)
	}
}

func TestAIDraftWithTemplateProvider(t *testing.T) {
	repo := &mockRepo{decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}}}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		AIProviderName: "template",
		AITemplates: map[string]map[string]string{
			"vocabulary": {
				"front": "Prefix: {{.Topic}}",
			},
		},
	})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewAI
	model.aiInput = "test"

	generate := model.generateDrafts()
	msgs := executeCmd(generate)
	for _, msg := range msgs {
		model.Update(msg)
	}

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
	model.width = 100
	model.height = 50
	model.Update(statsMsg{stats: core.Statistics{DailyGoal: 10, Grades: map[core.ReviewGrade]int{}}})
	model.activeView = ViewSettings
	model.settingsCursor = 5

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

func TestImportViewShowsResetAndStatusFilterGuidance(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewImport

	view := ansi.Strip(model.importScreen.Render(model, viewportLayout{}))
	for _, want := range []string{"[R] Reset DB", "Status filters apply to TSV and APKG exports"} {
		if !strings.Contains(view, want) {
			t.Fatalf("import view missing %q:\n%s", want, view)
		}
	}
}

func TestBrowserEmptyStateExplainsActiveFilters(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserSearch = "zzzz"

	view := ansi.Strip(model.renderBrowserAt(viewportLayout{Width: 90, Height: 30}))
	for _, want := range []string{"No cards found in No deck.", "Press Esc to clear search/tag filters"} {
		if !strings.Contains(view, want) {
			t.Fatalf("browser empty state missing %q:\n%s", want, view)
		}
	}
}

func TestCramReviewShowsDeckTagsAndPosition(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.decks = []core.Deck{{ID: "deck-1", Name: "Deck One"}}
	model.cramCards = []core.Card{{
		ID:     "c1",
		DeckID: "deck-1",
		Kind:   core.CardKindFlashcard,
		Prompt: "die Bahn",
		Answer: "train",
		Tags:   []string{"b2", "mobility"},
	}}
	model.cramActive = true

	view := ansi.Strip(model.renderCramAt(viewportLayout{Width: 90, Height: 30}))
	for _, want := range []string{"Deck: Deck One", "1/1", "Tags: #b2 #mobility"} {
		if !strings.Contains(view, want) {
			t.Fatalf("cram review missing %q:\n%s", want, view)
		}
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
	if !strings.Contains(dash, "Streak:") || (!strings.Contains(dash, "🔥") && !strings.Contains(dash, "⚡")) {
		t.Fatalf("dashboard should show current streak with streak emoji: %s", dash)
	}
}

func TestStatisticsShowsLeechAndBookmarkedDue(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"}},
		decks:    []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.height = 40 // Ensure enough height to show all stats without scrolling
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

	view := ansi.Strip(model.renderStatistics(viewportLayout{Width: 100, Height: 40}))
	if !strings.Contains(view, "Bookmarked:") || !strings.Contains(view, "3 (2 due)") {
		t.Fatalf("statistics should show bookmarked due: %s", view)
	}
	if !strings.Contains(view, "Leech:") || !strings.Contains(view, "1") {
		t.Fatalf("statistics should show leech count: %s", view)
	}
}

func TestScrollbarThumbScalesAndReachesTrackEnd(t *testing.T) {
	start, height := scrollbarThumb(40, 10, 0)
	if start != 0 || height != 2 {
		t.Fatalf("top thumb = %d/%d, want 0/2", start, height)
	}

	start, height = scrollbarThumb(40, 10, 30)
	if start != 8 || height != 2 {
		t.Fatalf("bottom thumb = %d/%d, want 8/2", start, height)
	}
}

func TestStatisticsScrollbarHitboxesAlignWithRenderedTrack(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewStatistics
	model.width = 90
	model.height = 30
	model.breakpoint = BreakpointMedium
	model.stats = core.Statistics{
		TotalCards:   20,
		DailyGoal:    10,
		Grades:       map[core.ReviewGrade]int{core.GradeAgain: 3, core.GradeHard: 2, core.GradeGood: 7, core.GradeEasy: 1},
		ReviewsToday: 4,
	}
	model.reviewsPerDay = map[string]int{}
	for i := 0; i < 14; i++ {
		day := time.Now().UTC().AddDate(0, 0, -i)
		dayStr := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		model.reviewsPerDay[dayStr] = i + 1
	}

	view := model.View().Content
	statsHitboxes := hitboxesWithPrefix(model.hitboxes, "stats-scroll-")
	if len(statsHitboxes) == 0 {
		t.Fatal("statistics scrollbar hitboxes were not registered")
	}
	for _, hitbox := range statsHitboxes {
		got := renderedRuneAt(view, hitbox.X, hitbox.Y)
		if got != '│' && got != '█' {
			t.Fatalf("hitbox %s at %d,%d maps to %q, want scrollbar track in:\n%s", hitbox.ID, hitbox.X, hitbox.Y, got, ansi.Strip(view))
		}
	}
}

func TestStatisticsScrollClickAndWheelClampToViewport(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewStatistics
	model.width = 90
	model.height = 30
	model.breakpoint = BreakpointMedium
	model.stats = core.Statistics{DailyGoal: 1, Grades: map[core.ReviewGrade]int{}}
	model.reviewsPerDay = map[string]int{}
	for i := 0; i < 14; i++ {
		day := time.Now().UTC().AddDate(0, 0, -i)
		dayStr := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		model.reviewsPerDay[dayStr] = i + 1
	}

	model.View()
	maxScroll := model.statsMaxScroll()
	if maxScroll <= 0 {
		t.Fatalf("statsMaxScroll = %d, want scrollable content", maxScroll)
	}
	model.activateHitboxByID(fmt.Sprintf("stats-scroll-%d", model.statisticsVisibleLines(model.activeViewContentLayout().Height)-1))
	if model.statsScroll != maxScroll {
		t.Fatalf("statsScroll after bottom click = %d, want %d", model.statsScroll, maxScroll)
	}

	model.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if model.statsScroll != maxScroll {
		t.Fatalf("statsScroll after wheel past bottom = %d, want %d", model.statsScroll, maxScroll)
	}
	model.activateHitboxByID("stats-scroll-0")
	if model.statsScroll != 0 {
		t.Fatalf("statsScroll after top click = %d, want 0", model.statsScroll)
	}
	model.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseWheelUp}))
	if model.statsScroll != 0 {
		t.Fatalf("statsScroll after wheel past top = %d, want 0", model.statsScroll)
	}
}

func TestBrowserAndCramScrollbarHitboxesAlignAndClick(t *testing.T) {
	cards := makeCards(16)

	browser := NewModel(&mockRepo{}, &mockScheduler{})
	browser.activeView = ViewBrowser
	browser.width = 90
	browser.height = 30
	browser.breakpoint = BreakpointMedium
	browser.browserCards = cards
	browser.View()
	assertScrollbarHitboxesAlign(t, browser.View().Content, browser.hitboxes, "browser-scroll-")
	browser.activateHitboxByID("browser-scroll-0")
	if browser.browserCursor != 0 {
		t.Fatalf("browser top click cursor = %d, want 0", browser.browserCursor)
	}
	browser.activateHitboxByID(fmt.Sprintf("browser-scroll-%d", browser.listVisibleLines(browser.activeViewContentLayout().Height)-1))
	if browser.browserCursor != len(cards)-1 {
		t.Fatalf("browser bottom click cursor = %d, want last card", browser.browserCursor)
	}

	cram := NewModel(&mockRepo{}, &mockScheduler{})
	cram.activeView = ViewCram
	cram.width = 90
	cram.height = 30
	cram.breakpoint = BreakpointMedium
	cram.cramCards = cards
	cram.View()
	assertScrollbarHitboxesAlign(t, cram.View().Content, cram.hitboxes, "cram-scroll-")
	cram.activateHitboxByID("cram-scroll-0")
	if cram.cramCursor != 0 {
		t.Fatalf("cram top click cursor = %d, want 0", cram.cramCursor)
	}
	cram.activateHitboxByID(fmt.Sprintf("cram-scroll-%d", cram.listVisibleLines(cram.activeViewContentLayout().Height)-1))
	if cram.cramCursor != len(cards)-1 {
		t.Fatalf("cram bottom click cursor = %d, want last card", cram.cramCursor)
	}
}

func TestDecksScrollbarHitboxesAlignWithRenderedTrack(t *testing.T) {
	decks := make([]core.Deck, 15)
	for i := 0; i < 15; i++ {
		decks[i] = core.Deck{
			ID:          fmt.Sprintf("deck-%d", i),
			Name:        fmt.Sprintf("Deck Name %d", i),
			Description: fmt.Sprintf("Description for deck %d that might be long", i),
			Tags:        []string{"tag1", "tag2"},
			TotalCards:  10,
			DueCards:    5,
		}
	}

	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewDecks
	model.width = 90
	model.height = 30
	model.breakpoint = BreakpointMedium
	model.decks = decks
	model.deckSelected = map[string]bool{}

	view := model.View().Content
	assertScrollbarHitboxesAlign(t, view, model.hitboxes, "deck-scroll-")
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

	model.revealState = RevealRevealed
	model.revealProgress = 100
	model.revealState = RevealRevealed
	model.revealProgress = 100
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
	model.revealState = RevealRevealed
	model.revealProgress = 100

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
	model.revealState = RevealRevealed
	model.revealProgress = 100

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
	model.revealState = RevealRevealed
	model.revealProgress = 100

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
	model.revealState = RevealRevealed
	model.revealProgress = 100

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

	if model.revealState == RevealRevealed {
		t.Fatal("flashcard should not be revealed initially")
	}

	view := model.renderReview(0, 0)
	if strings.Contains(view, "1:") {
		t.Fatalf("flashcard view should not show MCQ choices: %s", view)
	}

	model.Update(tea.KeyPressMsg{Code: ' '})
	if model.revealState != RevealFlipping {
		t.Fatalf("flashcard should be flipping after space, revealState=%v", model.revealState)
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
	model.revealState = RevealRevealed
	model.revealProgress = 100

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

func TestSessionGradeDistributionTracking(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
			{ID: "c2", DeckID: "deck-1", Prompt: "P2", Answer: "A2"},
			{ID: "c3", DeckID: "deck-1", Prompt: "P3", Answer: "A3"},
			{ID: "c4", DeckID: "deck-1", Prompt: "P4", Answer: "A4"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))

	model.Update(reviewRecordedMsg{cardID: "c1", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeAgain})
	model.Update(reviewRecordedMsg{cardID: "c2", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeGood})
	model.Update(reviewRecordedMsg{cardID: "c3", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeGood})
	model.Update(reviewRecordedMsg{cardID: "c4", cards: repo.dueCards, decks: repo.decks, stats: core.Statistics{}, grade: core.GradeEasy})

	if model.sessionGrades[core.GradeAgain] != 1 {
		t.Fatalf("GradeAgain = %d, want 1", model.sessionGrades[core.GradeAgain])
	}
	if model.sessionGrades[core.GradeGood] != 2 {
		t.Fatalf("GradeGood = %d, want 2", model.sessionGrades[core.GradeGood])
	}
	if model.sessionGrades[core.GradeEasy] != 1 {
		t.Fatalf("GradeEasy = %d, want 1", model.sessionGrades[core.GradeEasy])
	}
	if model.sessionGrades[core.GradeHard] != 0 {
		t.Fatalf("GradeHard = %d, want 0", model.sessionGrades[core.GradeHard])
	}
}

func TestSessionSummaryShowsGradeDistribution(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "P1", Answer: "A1"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.sessionReviewed = 3
	model.sessionCorrect = 2
	model.sessionGrades = map[core.ReviewGrade]int{
		core.GradeAgain: 1,
		core.GradeGood:  1,
		core.GradeEasy:  1,
	}

	view := summaryScreen{}.Render(model, viewportLayout{Width: 80, Height: 30})
	if !strings.Contains(view, "Grade Distribution") {
		t.Fatalf("session summary should show Grade Distribution: %s", view)
	}
	if !strings.Contains(view, "again") || !strings.Contains(view, "good") || !strings.Contains(view, "easy") {
		t.Fatalf("session summary should show all grade labels: %s", view)
	}
}

func TestHelpOverlayToggle(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})

	if model.showHelp {
		t.Fatal("help should be hidden initially")
	}

	model.Update(tea.KeyPressMsg{Text: "?"})
	if !model.showHelp {
		t.Fatal("help should be shown after pressing ?")
	}

	view := model.View().Content
	stripped := ansi.Strip(view)
	if !strings.Contains(stripped, "Keyboard Shortcuts") {
		t.Fatal("help overlay should be visible in view")
	}

	model.Update(tea.KeyPressMsg{Text: "?"})
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

	model.updateBrowserKey(tea.KeyPressMsg{Code: '/'})
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
	model.revealSpeed = 0
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

func TestCramRevealAtFullProgressRendersWithoutNegativeRepeat(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewCram
	model.cramActive = true
	model.cramCards = []core.Card{{ID: "c1", Prompt: "P1", Answer: "Antwort"}}
	model.cramCursor = 0
	model.revealState = RevealRevealing
	model.revealProgress = 100

	view := model.renderCramAt(viewportLayout{X: 0, Y: 0, Width: 80, Height: 24})
	if !strings.Contains(view, "Antwort") {
		t.Fatalf("expected full answer at 100%% progress, got: %s", view)
	}
}

func TestReviewArrowNavigation(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Test Deck"}},
		dueCards: []core.Card{
			{ID: "c1", Prompt: "Card 1", DeckID: "deck-1"},
			{ID: "c2", Prompt: "Card 2", DeckID: "deck-1"},
			{ID: "c3", Prompt: "Card 3", DeckID: "deck-1"},
		},
	}
	model := NewModel(repo, &mockScheduler{})

	// Load cards
	model.Update(dueCardsMsg(repo.dueCards))

	// Switch to Review view
	model.Update(tea.KeyPressMsg{Code: '3'})
	if model.activeView != ViewReview {
		t.Fatalf("active view = %s, want review", model.activeView)
	}

	// Initial cursor should be at index 0
	if model.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", model.cursor)
	}

	// Move cursor down with arrow key
	model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if model.cursor != 1 {
		t.Fatalf("cursor after down arrow = %d, want 1", model.cursor)
	}

	// Move cursor down again
	model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if model.cursor != 2 {
		t.Fatalf("cursor after second down arrow = %d, want 2", model.cursor)
	}

	// Try to move down when at end (should stay at 2)
	model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if model.cursor != 2 {
		t.Fatalf("cursor after down at end = %d, want 2", model.cursor)
	}

	// Move cursor up
	model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.cursor != 1 {
		t.Fatalf("cursor after up arrow = %d, want 1", model.cursor)
	}

	// Move cursor up again
	model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.cursor != 0 {
		t.Fatalf("cursor after second up arrow = %d, want 0", model.cursor)
	}

	// Try to move up when at start (should stay at 0)
	model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if model.cursor != 0 {
		t.Fatalf("cursor after up at start = %d, want 0", model.cursor)
	}
}

func TestApproveDraftPropagatesDueCardsError(t *testing.T) {
	repo := &mockRepo{
		decks:       []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		errDueCards: errors.New("due cards boom"),
	}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.OfflineProvider{})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewAI
	model.Update(draftsMsg([]ai.Draft{
		{
			Note: core.Note{
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
			},
		},
	}))

	cmd := model.approveDraft()
	if cmd == nil {
		t.Fatal("approveDraft should return command")
	}
	msg := cmd()
	if err, ok := msg.(error); !ok || !strings.Contains(err.Error(), "due cards boom") {
		t.Fatalf("msg = %T (%v), want due cards error", msg, msg)
	}
}

func TestApproveAllDraftsPropagatesDeckReloadError(t *testing.T) {
	repo := &mockRepo{
		decks:    []core.Deck{{ID: "deck-1", Name: "Deck One"}},
		errDecks: errors.New("decks boom"),
	}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.OfflineProvider{})
	model.Update(decksMsg([]core.Deck{{ID: "deck-1", Name: "Deck One"}}))
	model.activeView = ViewAI
	model.Update(draftsMsg([]ai.Draft{
		{
			Note: core.Note{
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
			},
		},
	}))

	cmd := model.approveAllDrafts()
	if cmd == nil {
		t.Fatal("approveAllDrafts should return command")
	}
	msg := cmd()
	if err, ok := msg.(error); !ok || !strings.Contains(err.Error(), "decks boom") {
		t.Fatalf("msg = %T (%v), want decks error", msg, msg)
	}
}

func TestAIApproveAllAndDiscardAllShortcuts(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.OfflineProvider{})
	model.activeView = ViewAI
	model.drafts = []ai.Draft{
		{Note: core.Note{ID: "n1", DeckID: "deck-1", Front: "Eins", Back: "One"}},
		{Note: core.Note{ID: "n2", DeckID: "deck-1", Front: "Zwei", Back: "Two"}},
	}

	// Test 'D' shortcut to discard all drafts
	cmd, handled := model.updateAIKey(tea.KeyPressMsg{Code: 'D', Text: "D"})
	if !handled || len(model.drafts) != 0 {
		t.Fatalf("expected 'D' shortcut to discard all drafts, got handled=%v, len(drafts)=%d", handled, len(model.drafts))
	}
	_ = cmd

	// Restore drafts and test 'A' shortcut to approve all drafts
	model.drafts = []ai.Draft{
		{Note: core.Note{ID: "n1", DeckID: "deck-1", Front: "Eins", Back: "One"}},
		{Note: core.Note{ID: "n2", DeckID: "deck-1", Front: "Zwei", Back: "Two"}},
	}
	cmd, handled = model.updateAIKey(tea.KeyPressMsg{Code: 'A', Text: "A"})
	if !handled || cmd == nil {
		t.Fatalf("expected 'A' shortcut to return approveAllDrafts command")
	}
	resMsg := cmd()
	impDone, ok := resMsg.(importDoneMsg)
	if !ok || impDone.count != 2 {
		t.Fatalf("expected importDoneMsg with count 2, got %T (%v)", resMsg, resMsg)
	}
}

func TestAIDraftRowHitboxesAndDictionaryContextBanner(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.OfflineProvider{})
	model.activeView = ViewAI
	model.width = 100
	model.height = 40
	model.draftSource = "Word: geheim\nTranslation: secret\n"
	model.drafts = []ai.Draft{
		{Note: core.Note{ID: "n1", DeckID: "deck-1", Front: "Eins", Back: "One"}},
		{Note: core.Note{ID: "n2", DeckID: "deck-1", Front: "Zwei", Back: "Two"}},
	}
	model.draftCursor = 0

	view := stripANSI(model.renderAI(0, 0))
	if !strings.Contains(view, "Dictionary Context: geheim") {
		t.Fatalf("expected dictionary context banner with headword, got:\n%s", view)
	}

	foundRow := false
	for _, hb := range model.hitboxes {
		if hb.ID == "draft-row-1" {
			foundRow = true
			hb.Action()
			break
		}
	}
	if !foundRow {
		t.Fatal("expected draft-row-1 hitbox for mouse selection")
	}
	if model.draftCursor != 1 {
		t.Fatalf("expected clicking draft row to set cursor to 1, got %d", model.draftCursor)
	}

	model.explanation = "geheim is an adjective meaning secret."
	view = stripANSI(model.renderAI(0, 0))
	if !strings.Contains(view, "AI Tutor Explanation") || !strings.Contains(view, "geheim is an adjective") {
		t.Fatalf("expected AI explanation panel, got:\n%s", view)
	}
	_, handled := model.updateAIKey(tea.KeyPressMsg{Code: 'H', Text: "H"})
	if !handled || model.explanation != "" {
		t.Fatalf("expected H to dismiss explanation, handled=%v explanation=%q", handled, model.explanation)
	}
}

func TestBulkBrowserBookmarkPropagatesError(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"},
		},
		errOnCardID: map[string]error{"c1": errors.New("bookmark boom")},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserCards = []core.Card{{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel"}}
	model.browserSelected["c1"] = true

	cmd := model.bulkBrowserBookmark(true)
	if cmd == nil {
		t.Fatal("bulkBrowserBookmark should return command")
	}
	msg := cmd()
	if err, ok := msg.(error); !ok || !strings.Contains(err.Error(), "bookmark boom") {
		t.Fatalf("msg = %T (%v), want bookmark error", msg, msg)
	}
}

func TestBulkBrowserToggleKindPropagatesError(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Kind: core.CardKindFlashcard, Prompt: "Apple", Answer: "Apfel"},
		},
		errOnCardID: map[string]error{"c1": errors.New("kind boom")},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserCards = []core.Card{{ID: "c1", DeckID: "deck-1", Kind: core.CardKindFlashcard, Prompt: "Apple", Answer: "Apfel"}}
	model.browserSelected["c1"] = true

	cmd := model.bulkBrowserToggleKind()
	if cmd == nil {
		t.Fatal("bulkBrowserToggleKind should return command")
	}
	msg := cmd()
	if err, ok := msg.(error); !ok || !strings.Contains(err.Error(), "kind boom") {
		t.Fatalf("msg = %T (%v), want kind error", msg, msg)
	}
}

func TestReviewHeaderShowsDeckTagsAndCardType(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{
			ID:     "c1",
			DeckID: "deck-1",
			Kind:   core.CardKindCloze,
			Prompt: "Die Stadt [baut] das Netz aus.",
			Answer: "baut",
			Tags:   []string{"b2", "mobility"},
		}},
		decks: []core.Deck{{ID: "deck-1", Name: "B2 Urban Mobility"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.Update(dueCardsMsg(repo.dueCards))
	model.activeView = ViewReview

	view := ansi.Strip(model.renderReview(0, 0))
	if !strings.Contains(view, "Deck: B2 Urban Mobility") {
		t.Fatalf("review header missing deck name:\n%s", view)
	}
	if !strings.Contains(view, "Tags: #b2 #mobility") {
		t.Fatalf("review header missing tags:\n%s", view)
	}
	if !strings.Contains(view, "Type: Cloze") {
		t.Fatalf("review header missing card type:\n%s", view)
	}
}

func TestBrowserPreviewShowsDeckKindExtraAndTags(t *testing.T) {
	repo := &mockRepo{
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.Update(decksMsg(repo.decks))
	model.activeView = ViewBrowser
	model.width = 110
	model.height = 40
	model.browserCards = []core.Card{{
		ID:     "c1",
		DeckID: "deck-1",
		Kind:   core.CardKindMCQ,
		Prompt: "Welche Option passt?",
		Answer: "umsteigen",
		Extra:  "Use umsteigen when changing vehicles.",
		Tags:   []string{"b2", "mobility"},
	}}

	view := ansi.Strip(model.renderBrowserAt(viewportLayout{Width: 100, Height: 35}))
	for _, want := range []string{"Card Preview:", "Deck:", "Kind:", "State:", "Front:", "Back:", "Reviews:", "Tags:"} {
		if !strings.Contains(view, want) {
			t.Fatalf("browser preview missing %q:\n%s", want, view)
		}
	}
}

func TestDashboardShowsCardMixForTallLayouts(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.width = 110
	model.height = 42
	model.stats = core.Statistics{NewCards: 4, YoungCards: 8, MatureCards: 12, DailyGoal: 10}

	view := ansi.Strip(model.renderDashboard(viewportLayout{Width: 100, Height: 34}))
	for _, want := range []string{"Card Mix", "New", "Young", "Mature"} {
		if !strings.Contains(view, want) {
			t.Fatalf("dashboard missing %q:\n%s", want, view)
		}
	}
}

func TestSpeechTextForCardPrefersGermanSide(t *testing.T) {
	cases := []struct {
		name string
		card core.Card
		want string
	}{
		{
			name: "front card",
			card: core.Card{ID: "note-1:front", Prompt: "der Kaffee", Answer: "coffee"},
			want: "der Kaffee",
		},
		{
			name: "reverse card",
			card: core.Card{ID: "note-1:back", Prompt: "coffee", Answer: "der Kaffee"},
			want: "der Kaffee",
		},
		{
			name: "cloze card",
			card: core.Card{ID: "note-1:cloze-1", Kind: core.CardKindCloze, Prompt: "Der [drink] ist heiß.", Answer: "Der Kaffee ist heiß."},
			want: "Der Kaffee ist heiß.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := speechTextForCard(tc.card); got != tc.want {
				t.Fatalf("speechTextForCard() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSelectAudioPlayerByPlatform(t *testing.T) {
	lookPath := func(available map[string]bool) func(string) (string, error) {
		return func(name string) (string, error) {
			if available[name] {
				return "/mock/" + name, nil
			}
			return "", errors.New("not found")
		}
	}

	cases := []struct {
		name      string
		goos      string
		available map[string]bool
		wantName  string
		wantArg   string
	}{
		{
			name:      "mac prefers afplay",
			goos:      "darwin",
			available: map[string]bool{"afplay": true, "mpv": true},
			wantName:  "/mock/afplay",
			wantArg:   "card.mp3",
		},
		{
			name:      "linux prefers mpv",
			goos:      "linux",
			available: map[string]bool{"mpv": true, "play": true},
			wantName:  "/mock/mpv",
			wantArg:   "--really-quiet",
		},
		{
			name:      "linux falls back to play",
			goos:      "linux",
			available: map[string]bool{"play": true},
			wantName:  "/mock/play",
			wantArg:   "card.mp3",
		},
		{
			name:      "windows uses powershell fallback",
			goos:      "windows",
			available: map[string]bool{"powershell.exe": true},
			wantName:  "/mock/powershell.exe",
			wantArg:   "-NoProfile",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := selectAudioPlayer("card.mp3", tc.goos, lookPath(tc.available))
			if err != nil {
				t.Fatalf("selectAudioPlayer() error: %v", err)
			}
			if spec.name != tc.wantName {
				t.Fatalf("player = %q, want %q", spec.name, tc.wantName)
			}
			if !containsString(spec.args, tc.wantArg) {
				t.Fatalf("args = %#v, want to contain %q", spec.args, tc.wantArg)
			}
		})
	}
}

func TestSelectAudioPlayerReturnsActionableError(t *testing.T) {
	_, err := selectAudioPlayer("card.mp3", "linux", func(string) (string, error) {
		return "", errors.New("not found")
	})
	if err == nil {
		t.Fatal("expected missing-player error")
	}
	if got := err.Error(); !strings.Contains(got, "install one of") || !strings.Contains(got, "mpv") {
		t.Fatalf("error = %q, want actionable player list", got)
	}
}

func containsString(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func makeCards(count int) []core.Card {
	cards := make([]core.Card, count)
	for i := range cards {
		cards[i] = core.Card{
			ID:     fmt.Sprintf("c-%02d", i+1),
			Prompt: fmt.Sprintf("Card %02d", i+1),
			Answer: fmt.Sprintf("Answer %02d", i+1),
			Kind:   core.CardKindFlashcard,
		}
	}
	return cards
}

func hitboxesWithPrefix(hitboxes []Hitbox, prefix string) []Hitbox {
	matches := make([]Hitbox, 0)
	for _, hitbox := range hitboxes {
		if strings.HasPrefix(hitbox.ID, prefix) {
			matches = append(matches, hitbox)
		}
	}
	return matches
}

func assertScrollbarHitboxesAlign(t *testing.T, view string, hitboxes []Hitbox, prefix string) {
	t.Helper()
	matches := hitboxesWithPrefix(hitboxes, prefix)
	if len(matches) == 0 {
		t.Fatalf("%s hitboxes were not registered", prefix)
	}
	for _, hitbox := range matches {
		got := renderedRuneAt(view, hitbox.X, hitbox.Y)
		if got != '│' && got != '█' {
			t.Fatalf("hitbox %s at %d,%d maps to %q, want scrollbar track in:\n%s", hitbox.ID, hitbox.X, hitbox.Y, got, ansi.Strip(view))
		}
	}
}

func renderedRuneAt(view string, x, y int) rune {
	lines := strings.Split(ansi.Strip(view), "\n")
	if y < 0 || y >= len(lines) {
		return 0
	}
	runes := []rune(lines[y])
	if x < 0 || x >= len(runes) {
		return 0
	}
	return runes[x]
}

func TestBrowserSelectAllAndPlayAudio(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{ID: "c1", DeckID: "deck-1", Prompt: "Apple", Answer: "Apfel", Audio: "apple.mp3"},
			{ID: "c2", DeckID: "deck-1", Prompt: "Banana", Answer: "Banane"},
		},
		decks: []core.Deck{{ID: "deck-1", Name: "Deck One"}},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewBrowser
	model.browserDeckID = "deck-1"
	model.browserCards = repo.dueCards

	// 1. Test Select All
	_, handled := model.updateBrowserKey(tea.KeyPressMsg{Code: 'a'})
	if !handled {
		t.Fatal("expected 'a' to be handled")
	}
	if !model.browserSelected["c1"] || !model.browserSelected["c2"] {
		t.Fatal("expected both cards to be selected after 'a'")
	}

	// 2. Test Deselect All (using toggle 'a' again)
	_, handled = model.updateBrowserKey(tea.KeyPressMsg{Code: 'a'})
	if !handled {
		t.Fatal("expected second 'a' to be handled")
	}
	if model.browserSelected["c1"] || model.browserSelected["c2"] {
		t.Fatal("expected both cards to be deselected after second 'a'")
	}

	// 3. Test Play Audio
	cmd, handled := model.updateBrowserKey(tea.KeyPressMsg{Code: 'p'})
	if !handled {
		t.Fatal("expected 'p' to be handled")
	}
	if cmd == nil {
		t.Fatal("expected command from playing card audio")
	}
	if model.status != "Playing audio..." {
		t.Fatalf("expected status 'Playing audio...', got %q", model.status)
	}
}

func TestAdjectiveEndingTrainer(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// 1. Select the Adjective Ending Trainer (key "4")
	updated, cmd := model.Update(tea.KeyPressMsg{Code: '4'})
	model = updated.(*Model)
	if model.practiceSubView != PracticeSubViewAdjective {
		t.Fatalf("expected practiceSubView to be PracticeSubViewAdjective, got %v", model.practiceSubView)
	}
	if cmd == nil {
		t.Fatal("expected load command")
	}

	// Execute load command
	msgs := executeCmd(cmd)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	updated, _ = model.Update(msgs[0])
	model = updated.(*Model)

	st := model.trainerStateFor(PracticeSubViewAdjective)
	if len(st.items) == 0 {
		t.Fatal("expected loaded adjective items to be non-empty")
	}

	// Test inputs
	// Type correct/incorrect ending
	st.input = "kaltes"
	updated, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model = updated.(*Model)
	if !st.revealed {
		t.Fatal("expected exercise to be revealed after Enter")
	}

	// Press key to move to next item
	updated, _ = model.Update(tea.KeyPressMsg{Code: ' '})
	model = updated.(*Model)
	if st.revealed {
		t.Fatal("expected exercise to reset revealed state on next action")
	}
	if st.index != 1 {
		t.Fatalf("expected adjective index to advance to 1, got %d", st.index)
	}
}

func TestPluralTrainer(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// 1. Select the Plural Trainer (key "6")
	updated, cmd := model.Update(tea.KeyPressMsg{Code: '6'})
	model = updated.(*Model)
	if model.practiceSubView != PracticeSubViewPlural {
		t.Fatalf("expected practiceSubView to be PracticeSubViewPlural, got %v", model.practiceSubView)
	}
	if cmd == nil {
		t.Fatal("expected load command")
	}

	// Execute load command
	msgs := executeCmd(cmd)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	updated, _ = model.Update(msgs[0])
	model = updated.(*Model)

	st := model.trainerStateFor(PracticeSubViewPlural)
	if len(st.items) == 0 {
		t.Fatal("expected loaded plural items to be non-empty")
	}

	// Set singular and plural for testing
	st.items[0] = trainerItem{
		Title:       "das Buch",
		Subtitle:    "(book)",
		Answer:      "die Bücher",
		Instruction: "Enter the plural form (with or without article):",
	}

	// Test inputs
	// Type correct plural form
	st.input = "die Bücher"
	updated, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	model = updated.(*Model)
	if !st.revealed {
		t.Fatal("expected exercise to be revealed after Enter")
	}
	if !st.lastOK {
		t.Fatal("expected correct answer to be recorded")
	}

	// Press key to move to next item
	updated, _ = model.Update(tea.KeyPressMsg{Code: ' '})
	model = updated.(*Model)
	if st.revealed {
		t.Fatal("expected exercise to reset revealed state on next action")
	}
	if st.index != 1 {
		t.Fatalf("expected plural index to advance to 1, got %d", st.index)
	}
}

func TestPracticeMouseClicks(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{
			{
				ID:     "card-1",
				Prompt: "der Tisch",
				Answer: "table",
				DeckID: "mock-1",
			},
			{
				ID:     "card-2",
				Prompt: "die Tür",
				Answer: "door",
				DeckID: "mock-1",
			},
		},
	}
	model := NewModel(repo, &mockScheduler{})
	model.activeView = ViewPractice
	model.practiceSubView = PracticeSubViewHub

	// 1. Click "practice-gender" button on Practice Hub
	cmd := model.activateHitboxByID("practice-gender")
	if cmd == nil {
		t.Fatal("expected load command")
	}
	if model.practiceSubView != PracticeSubViewGender {
		t.Fatalf("expected practiceSubView to transition to PracticeSubViewGender, got %v", model.practiceSubView)
	}

	// 2. Load the items
	msgs := executeCmd(cmd)
	updated, _ := model.Update(msgs[0])
	model = updated.(*Model)

	if len(model.practiceItems) == 0 {
		t.Fatal("expected loaded practice items to be non-empty")
	}

	// 3. Click one of the gender options (e.g. "gender-opt-der")
	// Since we appended the hitboxes inside renderGenderTrainer, let's trigger the render first to populate model.hitboxes.
	model.View() // Render to populate hitboxes

	// Find the hitbox for "gender-opt-der"
	var hitbox Hitbox
	found := false
	for _, h := range model.hitboxes {
		if h.ID == "gender-opt-der" {
			hitbox = h
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected gender-opt-der hitbox to be registered")
	}

	// Activate it
	model.activateHitbox(hitbox)
	if !model.practiceRevealed {
		t.Fatal("expected practice to be revealed after clicking der")
	}

	// 4. Click to proceed to next noun
	model.View() // Render to populate hitboxes in revealed state

	foundNext := false
	for _, h := range model.hitboxes {
		if h.ID == "gender-next" {
			hitbox = h
			foundNext = true
			break
		}
	}
	if !foundNext {
		t.Fatal("expected gender-next hitbox to be registered")
	}

	model.activateHitbox(hitbox)
	if model.practiceRevealed {
		t.Fatal("expected practice to reset revealed state on next action click")
	}
	if model.practiceIndex != 1 {
		t.Fatalf("expected practiceIndex to advance to 1, got %d", model.practiceIndex)
	}
}

func TestConjunctionsTrainer(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewConjunctions

	st := m.trainerStateFor(PracticeSubViewConjunctions)
	st.items = []trainerItem{
		{
			Title:       "Ich bleibe heute zu Hause, {{...}} es regnet.",
			Subtitle:    "Meaning: because",
			Answer:      "weil",
			Instruction: "Enter the missing word:",
			Explanation: "weil is a subordinating conjunction",
		},
		{
			Title:       "Es regnet, {{...}} ich gehe trotzdem spazieren.",
			Subtitle:    "Meaning: but",
			Answer:      "aber",
			Instruction: "Enter the missing word:",
			Explanation: "aber is coordinating",
		},
	}
	st.index = 0

	layout := viewportLayout{Width: 80, Height: 30}
	view := m.renderTrainer(PracticeSubViewConjunctions, layout)
	if !strings.Contains(view, "Ich bleibe heute zu Hause") {
		t.Fatalf("expected view to contain sentence: %s", view)
	}

	for _, char := range "weil" {
		m.updatePracticeKey(tea.KeyPressMsg{Code: char})
	}
	if st.input != "weil" {
		t.Fatalf("expected input to be 'weil', got '%s'", st.input)
	}

	m.updatePracticeKey(tea.KeyPressMsg{Code: '\r'})
	if !st.revealed {
		t.Fatal("expected trainer to be in revealed state")
	}
	if !st.lastOK {
		t.Fatal("expected answer to be correct")
	}
	if st.correct != 1 || st.total != 1 {
		t.Fatalf("expected score to be 1/1, got %d/%d", st.correct, st.total)
	}

	view = m.renderTrainer(PracticeSubViewConjunctions, layout)
	if !strings.Contains(view, "weil is a subordinating conjunction") {
		t.Fatalf("expected explanation to be rendered, got: %s", view)
	}

	m.updatePracticeKey(tea.KeyPressMsg{Code: ' '})
	if st.revealed {
		t.Fatal("expected trainer to reset revealed state on advance")
	}
	if st.index != 1 {
		t.Fatalf("expected trainer to advance to next item, index=%d", st.index)
	}
	if st.input != "" {
		t.Fatalf("expected input to reset, got '%s'", st.input)
	}
}

func TestConjunctionsHint(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewConjunctions

	st := m.trainerStateFor(PracticeSubViewConjunctions)
	st.items = []trainerItem{
		{
			Title:       "Sentence",
			Subtitle:    "Meaning: meaning",
			Answer:      "weil",
			Instruction: "Enter the missing word:",
			HintText:    "THIS IS A HINT",
			Explanation: "explanation",
		},
	}
	st.index = 0

	layout := viewportLayout{Width: 80, Height: 30}
	view := m.renderTrainer(PracticeSubViewConjunctions, layout)
	if strings.Contains(view, "THIS IS A HINT") {
		t.Fatal("hint should not be visible initially")
	}

	m.updatePracticeKey(tea.KeyPressMsg{Code: 'h'})
	if !st.showHint {
		t.Fatal("expected showHint to be true")
	}
	view = m.renderTrainer(PracticeSubViewConjunctions, layout)
	if !strings.Contains(view, "THIS IS A HINT") {
		t.Fatal("hint should be visible after pressing 'h'")
	}

	m.updatePracticeKey(tea.KeyPressMsg{Code: 'h'})
	if st.showHint {
		t.Fatal("expected showHint to be false after toggle")
	}

	// Reveal answer should hide hint input prompt but hint itself?
	// The current implementation hides the hint prompt when revealed.
	m.updatePracticeKey(tea.KeyPressMsg{Code: 'h'})
	for _, char := range "weil" {
		m.updatePracticeKey(tea.KeyPressMsg{Code: char})
	}
	m.updatePracticeKey(tea.KeyPressMsg{Code: '\r'})
	if !st.revealed {
		t.Fatal("expected trainer to be revealed")
	}

	// Advance should reset hint
	m.updatePracticeKey(tea.KeyPressMsg{Code: ' '})
	if st.showHint {
		t.Fatal("expected hint to reset on advance")
	}
}

func TestPracticeHubScoresAndReset(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.activeView = ViewPractice
	m.practiceSubView = PracticeSubViewHub
	m.width = 80
	m.height = 30

	// Mock some practice scores
	m.practiceCorrect, m.practiceTotal = 4, 5
	conj := m.trainerStateFor(PracticeSubViewConjugation)
	conj.correct, conj.total = 3, 3

	layout := m.activeViewContentLayout()
	view := m.renderPracticeHub(layout)
	plainView := stripANSI(view)

	// Verify that score info is rendered on the Practice Hub buttons
	if !strings.Contains(plainView, "4/5 (80%)") {
		t.Errorf("expected score 4/5 (80%%) to be displayed on gender trainer, got:\n%s", plainView)
	}
	if !strings.Contains(plainView, "3/3 (100%)") {
		t.Errorf("expected score 3/3 (100%%) to be displayed on conjugation trainer, got:\n%s", plainView)
	}

	// Press 'r' to reset scores
	cmd, handled := m.updatePracticeKey(tea.KeyPressMsg{Code: 'r'})
	if !handled || cmd != nil {
		t.Fatal("expected 'r' keypress to be handled in Practice Hub")
	}

	// Verify scores are reset to 0
	if m.practiceCorrect != 0 || m.practiceTotal != 0 || conj.correct != 0 || conj.total != 0 {
		t.Errorf("expected practice scores to be reset, got: gender=%d/%d, conj=%d/%d",
			m.practiceCorrect, m.practiceTotal, conj.correct, conj.total)
	}

	// Re-render and check that score info is gone
	view = m.renderPracticeHub(layout)
	plainView = stripANSI(view)
	if strings.Contains(plainView, "4/5") || strings.Contains(plainView, "80%") {
		t.Errorf("expected scores to be cleared from view after reset, got:\n%s", plainView)
	}
}

func TestSpliceVisual(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		startCol    int
		width       int
		replacement string
		want        string
	}{
		{
			name:        "plain text",
			line:        "Hello World",
			startCol:    6,
			width:       5,
			replacement: "DEUTSCH",
			want:        "Hello DEUTSCH",
		},
		{
			name:        "ANSI at start",
			line:        "\x1b[31mRedText\x1b[0m here",
			startCol:    7,
			width:       4,
			replacement: "OVER",
			want:        "\x1b[31mRedText\x1b[0mOVERe",
		},
		{
			name:        "ANSI in prefix preserved",
			line:        "\x1b[38;5;81mBoldLabel\x1b[0m rest",
			startCol:    9,
			width:       5,
			replacement: "REPL",
			want:        "\x1b[38;5;81mBoldLabel\x1b[0mREPL",
		},
		{
			name:        "full replacement with ANSI",
			line:        "\x1b[32mGreen\x1b[0mText",
			startCol:    0,
			width:       10,
			replacement: "ALLNEW",
			want:        "ALLNEW",
		},
		{
			name:        "wide chars",
			line:        "Hallo 世界",
			startCol:    6,
			width:       2,
			replacement: "DE",
			want:        "Hallo DE界",
		},
		{
			name:        "start beyond line length",
			line:        "Short",
			startCol:    20,
			width:       4,
			replacement: "NEW",
			want:        "ShortNEW",
		},
		{
			name:        "zero width replacement",
			line:        "Hello World",
			startCol:    5,
			width:       1,
			replacement: "",
			want:        "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spliceVisual(tt.line, tt.startCol, tt.width, tt.replacement)
			if got != tt.want {
				t.Errorf("spliceVisual(%q, %d, %d, %q)\n= %q\nwant %q", tt.line, tt.startCol, tt.width, tt.replacement, got, tt.want)
			}
		})
	}
}

func TestApplyOverlay(t *testing.T) {
	m := &Model{width: 80, height: 24}
	simpleBase := "Line 1\nLine 2\nLine 3"
	overlay := "\x1b[31mOVERLAY\x1b[0m"
	result := m.applyOverlay(simpleBase, overlay)
	stripped := ansi.Strip(result)
	if !strings.Contains(stripped, "OVERLAY") {
		t.Errorf("applyOverlay should contain OVERLAY, got:\n%s", stripped)
	}
	if !strings.Contains(stripped, "Line 1") {
		t.Errorf("applyOverlay should preserve Line 1, got:\n%s", stripped)
	}

	ansiBase := "\x1b[38;5;81mDashboard\x1b[0m content\n\x1b[32mMore\x1b[0m text"
	result = m.applyOverlay(ansiBase, overlay)
	stripped = ansi.Strip(result)
	if !strings.Contains(stripped, "Dashboard") {
		t.Errorf("applyOverlay should preserve Dashboard, got:\n%s", stripped)
	}
	if !strings.Contains(stripped, "OVERLAY") {
		t.Errorf("applyOverlay should contain OVERLAY, got:\n%s", stripped)
	}

	wideBase := strings.Repeat("X", 40) + "GHOST" + strings.Repeat("Y", 20)
	overlayBox := strings.Repeat(" ", 10) + "BOX" + strings.Repeat(" ", 7)
	result = m.applyOverlay(wideBase, overlayBox)
	stripped = ansi.Strip(result)
	if strings.Contains(stripped, "GHOST") {
		t.Errorf("applyOverlay should clear underlying content in overlay region, got:\n%s", stripped)
	}
	if !strings.Contains(stripped, "BOX") {
		t.Errorf("applyOverlay should contain overlay content, got:\n%s", stripped)
	}
}

func TestFillViewportContent(t *testing.T) {
	layout := viewportLayout{Width: 10, Height: 3}
	filled := fillViewportContent("short\nline", layout)
	lines := strings.Split(filled, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), filled)
	}
	if lipgloss.Width(lines[0]) != 10 {
		t.Fatalf("expected padded width 10, got %d for %q", lipgloss.Width(lines[0]), lines[0])
	}
	if lines[2] != strings.Repeat(" ", 10) {
		t.Fatalf("expected blank padded line, got %q", lines[2])
	}

	tall := strings.Join([]string{strings.Repeat("W", 20), strings.Repeat("X", 20), strings.Repeat("Y", 20)}, "\n")
	filled = fillViewportContent(tall, layout)
	if filled != tall {
		t.Fatalf("expected tall content to pass through unchanged, got %q", filled)
	}
}

func TestTestModeFreezesTimer(t *testing.T) {
	repo := &mockRepo{}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		TestMode: true,
	})
	model.activeView = ViewReview
	model.dueCards = []core.Card{{ID: "c1", Prompt: "P", Answer: "A"}}
	model.sessionStartTime = time.Now().Add(-10 * time.Minute)

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "00:00") {
		t.Errorf("expected timer to be frozen at 00:00 in test mode, got:\n%s", view)
	}

	// Double check without test mode
	model.testMode = false
	view = model.renderReview(0, 0)
	// Since it's 10 minutes ago, it should be 10:00 or similar
	if strings.Contains(view, " 00:00") {
		t.Errorf("expected timer NOT to be 00:00 without test mode, got:\n%s", view)
	}
}
