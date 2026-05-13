package tui

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/core"
)

// fakeFixProvider is a Provider that also satisfies ai.ChatProvider so
// FixCard routes through SendChat instead of the no-op offline path.
type fakeFixProvider struct {
	response string
}

func (f fakeFixProvider) GenerateDrafts(context.Context, ai.DraftRequest) ([]ai.Draft, error) {
	return nil, nil
}

func (f fakeFixProvider) SendChat(context.Context, string, string) (string, error) {
	return f.response, nil
}

func executeFixCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestReportCardWrong_StartsFixFlowAndStoresProposal(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{
			ID: "c1", NoteID: "n1", DeckID: "d1",
			Prompt: "Hund", Answer: "dog",
		}},
	}
	provider := fakeFixProvider{
		response: `{"front":"der Hund","back":"the dog","extra":"masc","example":"Ich habe einen Hund.","reason":"Add article"}`,
	}
	model := NewModelWithAI(repo, &mockScheduler{}, provider)
	model.activeView = ViewReview
	model.dueCards = repo.dueCards

	cmd := model.reportCardWrong()
	if cmd == nil {
		t.Fatal("reportCardWrong should return a command")
	}
	if !model.fixingCard {
		t.Fatal("fixingCard should be true while AI request is in flight")
	}

	msg := executeFixCmd(cmd)
	prop, ok := msg.(fixProposalMsg)
	if !ok {
		t.Fatalf("expected fixProposalMsg, got %T: %+v", msg, msg)
	}
	if prop.proposal.Front != "der Hund" {
		t.Errorf("front in proposal = %q", prop.proposal.Front)
	}

	model.Update(prop)
	if model.fixingCard {
		t.Error("fixingCard should be cleared after proposal arrives")
	}
	if model.fixProposal == nil {
		t.Fatal("fixProposal should be set")
	}
	if model.fixOldNote == nil {
		t.Fatal("fixOldNote should be set")
	}

	view := model.renderReview(0, 0)
	if !strings.Contains(view, "AI-proposed fix") {
		t.Errorf("review render should include proposal block; got:\n%s", view)
	}
	if !strings.Contains(view, "der Hund") {
		t.Errorf("proposal should show new front; got:\n%s", view)
	}
}

func TestReportCardWrong_NoProviderShowsHint(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{ID: "c1", NoteID: "n1", Prompt: "x", Answer: "y"}},
	}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.OfflineProvider{})
	model.aiProvider = nil // simulate "disabled" provider
	model.activeView = ViewReview
	model.dueCards = repo.dueCards

	cmd := model.reportCardWrong()
	if cmd != nil {
		t.Fatal("no provider should produce no command")
	}
	if model.fixingCard {
		t.Error("fixingCard should remain false when provider is nil")
	}
	if !strings.Contains(model.status, "AI provider") {
		t.Errorf("status should hint at enabling provider, got: %q", model.status)
	}
}

func TestDiscardFixProposalResetsState(t *testing.T) {
	model := NewModelWithAI(&mockRepo{}, &mockScheduler{}, ai.OfflineProvider{})
	note := core.Note{ID: "n1", Front: "a", Back: "b"}
	prop := ai.FixedNote{Front: "A", Back: "B"}
	model.fixOldNote = &note
	model.fixProposal = &prop
	model.fixCardID = "c1"

	model.discardFixProposal()
	if model.fixProposal != nil || model.fixOldNote != nil || model.fixingCard {
		t.Fatal("discardFixProposal must clear all fix state")
	}
	if model.status == "" {
		t.Error("status should be updated after discard")
	}
}

func TestApplyFixProposal_PersistsViaUpsertNote(t *testing.T) {
	repo := &mockRepo{
		dueCards: []core.Card{{
			ID: "c1", NoteID: "n1", DeckID: "d1",
			Prompt: "Hund", Answer: "dog",
		}},
	}
	model := NewModelWithAI(repo, &mockScheduler{}, ai.OfflineProvider{})
	model.activeView = ViewReview
	model.dueCards = repo.dueCards

	note := core.Note{ID: "n1", DeckID: "d1", Front: "Hund", Back: "dog"}
	model.fixOldNote = &note
	model.fixProposal = &ai.FixedNote{Front: "der Hund", Back: "the dog", Extra: "masc", Example: "Ich habe einen Hund."}
	model.fixCardID = "c1"

	cmd := model.applyFixProposal()
	if cmd == nil {
		t.Fatal("applyFixProposal should return a command")
	}
	msg := executeFixCmd(cmd)
	if _, ok := msg.(fixAppliedMsg); !ok {
		t.Fatalf("expected fixAppliedMsg, got %T: %+v", msg, msg)
	}
	// mockRepo.UpsertNote should have written the new content back into dueCards
	if repo.dueCards[0].Prompt != "der Hund" || repo.dueCards[0].Answer != "the dog" {
		t.Errorf("mock repo not updated: %+v", repo.dueCards[0])
	}
}
