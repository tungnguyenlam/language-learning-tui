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
