package ai

import (
	"context"
	"testing"

	"deutsch-tui/internal/content"
)

func TestValidateDraftsRejectsDuplicateIDs(t *testing.T) {
	deck := content.StarterDeck()
	drafts := []Draft{
		{Note: deck.Notes[0]},
		{Note: deck.Notes[0]},
	}
	if err := ValidateDrafts(drafts); err == nil {
		t.Fatal("expected duplicate id error")
	}
}

func TestFakeProvider(t *testing.T) {
	deck := content.StarterDeck()
	provider := FakeProvider{Drafts: []Draft{{Note: deck.Notes[0]}}}
	drafts, err := provider.GenerateDrafts(context.Background(), DraftRequest{DeckID: deck.ID})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("drafts = %d, want 1", len(drafts))
	}
}

func TestOfflineProviderGeneratesValidDraft(t *testing.T) {
	provider := OfflineProvider{}
	drafts, err := provider.GenerateDrafts(context.Background(), DraftRequest{
		SourceText: "der Kaffee",
		DeckID:     "a1-survival",
		Tags:       []string{"food"},
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if err := ValidateDrafts(drafts); err != nil {
		t.Fatalf("validate drafts: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("drafts = %d, want 1", len(drafts))
	}
	if drafts[0].Note.ID != "ai-der-kaffee" {
		t.Fatalf("note id = %q, want ai-der-kaffee", drafts[0].Note.ID)
	}
	if len(drafts[0].Note.Cards) == 0 {
		t.Fatal("generated draft should include cards")
	}
}

func TestOfflineProviderRequiresSourceText(t *testing.T) {
	provider := OfflineProvider{}
	if _, err := provider.GenerateDrafts(context.Background(), DraftRequest{DeckID: "a1"}); err == nil {
		t.Fatal("expected missing source text error")
	}
}

func TestTemplateProvider(t *testing.T) {
	provider := TemplateProvider{
		Templates: map[string]string{
			"front":   "Word: {{.Topic}}",
			"back":    "German for {{.Topic}} is ...",
			"example": "Ich mag {{.Topic}}.",
		},
	}
	drafts, err := provider.GenerateDrafts(context.Background(), DraftRequest{
		SourceText: "Coffee",
		DeckID:     "a1",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("len(drafts) = %d, want 1", len(drafts))
	}

	note := drafts[0].Note
	if note.Front != "Word: Coffee" {
		t.Errorf("front = %q, want %q", note.Front, "Word: Coffee")
	}
	if note.Back != "German for Coffee is ..." {
		t.Errorf("back = %q, want %q", note.Back, "German for Coffee is ...")
	}
	if note.Examples[0] != "Ich mag Coffee." {
		t.Errorf("example = %q, want %q", note.Examples[0], "Ich mag Coffee.")
	}
}
