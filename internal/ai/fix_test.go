package ai

import (
	"context"
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

// fakeChatProvider implements both Provider and ChatProvider so FixCard
// can route through the chat path without touching the network.
type fakeChatProvider struct {
	response string
	err      error
	gotSys   string
	gotUser  string
}

func (f *fakeChatProvider) GenerateDrafts(context.Context, DraftRequest) ([]Draft, error) {
	return nil, nil
}

func (f *fakeChatProvider) SendChat(_ context.Context, system, user string) (string, error) {
	f.gotSys = system
	f.gotUser = user
	return f.response, f.err
}

func TestFixCard_ParsesProposalFromChatJSON(t *testing.T) {
	p := &fakeChatProvider{
		response: `{"front":"der Hund","back":"the dog","extra":"masc, pl: die Hunde","example":"Ich habe einen Hund.","reason":"Added article."}`,
	}
	note := core.Note{ID: "n1", DeckID: "d1", Front: "Hund", Back: "dog"}
	fixed, err := FixCard(context.Background(), p, FixRequest{Note: note, UserReport: "missing article"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fixed.Front != "der Hund" {
		t.Errorf("front = %q, want %q", fixed.Front, "der Hund")
	}
	if fixed.Back != "the dog" {
		t.Errorf("back = %q", fixed.Back)
	}
	if !strings.Contains(fixed.Reason, "article") {
		t.Errorf("reason = %q", fixed.Reason)
	}
	if !strings.Contains(p.gotUser, "missing article") {
		t.Errorf("user prompt should contain the user report; got: %q", p.gotUser)
	}
}

func TestFixCard_RequiresProvider(t *testing.T) {
	_, err := FixCard(context.Background(), nil, FixRequest{Note: core.Note{ID: "n1"}})
	if err == nil {
		t.Fatal("expected error for nil provider")
	}
}

func TestFixCard_NonChatProviderReturnsNoOp(t *testing.T) {
	note := core.Note{
		ID: "n1", DeckID: "d1", Front: "Hund", Back: "dog",
		Extra: "noun", Examples: []string{"Ein Hund läuft."},
	}
	fixed, err := FixCard(context.Background(), OfflineProvider{}, FixRequest{Note: note})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fixed.Front != note.Front || fixed.Back != note.Back {
		t.Errorf("offline provider must leave content unchanged, got %+v", fixed)
	}
	if fixed.Example != "Ein Hund läuft." {
		t.Errorf("example = %q", fixed.Example)
	}
}

func TestParseFixedNoteJSON_RejectsMissingFields(t *testing.T) {
	if _, err := parseFixedNoteJSON(`{"front":"","back":"x"}`); err == nil {
		t.Fatal("expected error for empty front")
	}
	if _, err := parseFixedNoteJSON(`not json`); err == nil {
		t.Fatal("expected error for non-JSON")
	}
}

func TestParseFixedNoteJSON_TrimsWhitespace(t *testing.T) {
	fixed, err := parseFixedNoteJSON(`prose before {"front":"  der Hund  ","back":"the dog ","extra":"","example":"  Ich habe einen Hund. ","reason":""} after`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fixed.Front != "der Hund" {
		t.Errorf("front not trimmed: %q", fixed.Front)
	}
	if fixed.Example != "Ich habe einen Hund." {
		t.Errorf("example not trimmed: %q", fixed.Example)
	}
}
