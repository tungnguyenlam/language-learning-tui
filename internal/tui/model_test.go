package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
			{ID: "deck-1", Name: "Deck One", TotalCards: 10, DueCards: 5},
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
	view := model.renderDecks(0, 0)
	if !strings.Contains(view, "Deck Two (20 total, 2 due)") {
		t.Fatalf("decks view rendering missing stats: %s", view)
	}
}

func TestSettingsProviderSwitching(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewSettings
	if model.aiProviderName != "offline" {
		t.Fatalf("initial provider = %q, want offline", model.aiProviderName)
	}

	// Cursor is at 0 (AI Provider)
	model.Update(tea.KeyPressMsg{Code: '\r'}) // Enter
	if model.aiProviderName != "template" {
		t.Fatalf("provider after toggle = %q, want template", model.aiProviderName)
	}
	if _, ok := model.aiProvider.(ai.TemplateProvider); !ok {
		t.Fatal("aiProvider should be TemplateProvider")
	}

	model.Update(tea.KeyPressMsg{Code: '\r'}) // Toggle back
	if model.aiProviderName != "offline" {
		t.Fatalf("provider after toggle back = %q, want offline", model.aiProviderName)
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
