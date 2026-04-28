package sqlite

import (
	"context"
	"testing"
	"time"

	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"
)

func TestStoreDeckAndFindDueCards(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	deck := content.StarterDeck()
	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}
	cards, err := store.DueCards(ctx, time.Now(), 10)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	if len(cards) == 0 {
		t.Fatal("expected due starter cards")
	}
}

func TestMigrationsAreTracked(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	ids, err := store.AppliedMigrations(ctx)
	if err != nil {
		t.Fatalf("applied migrations: %v", err)
	}
	if len(ids) != len(migrations) {
		t.Fatalf("applied migrations = %d, want %d", len(ids), len(migrations))
	}
	for i, migration := range migrations {
		if ids[i] != migration.ID {
			t.Fatalf("migration id at %d = %d, want %d", i, ids[i], migration.ID)
		}
	}
}

func TestRecordReviewHidesFutureCard(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	deck := content.StarterDeck()
	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}
	card := deck.Notes[0].Cards[0]
	now := time.Date(2026, 4, 29, 10, 0, 0, 0, time.UTC)
	result := core.ReviewResult{
		CardID:   card.ID,
		Grade:    core.GradeGood,
		Reviewed: now,
		Next: core.ReviewState{
			CardID:   card.ID,
			Due:      now.Add(24 * time.Hour),
			Interval: 24 * time.Hour,
			Ease:     2.5,
			Reviews:  1,
		},
	}
	if err := store.RecordReview(ctx, result); err != nil {
		t.Fatalf("record review: %v", err)
	}
	cards, err := store.DueCards(ctx, now, 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	for _, due := range cards {
		if due.ID == card.ID {
			t.Fatalf("reviewed card %s should not be due", card.ID)
		}
	}
}
