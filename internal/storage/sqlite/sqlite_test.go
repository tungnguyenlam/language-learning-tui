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

func TestGetDeckReturnsNotesAndCards(t *testing.T) {
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
	got, err := store.GetDeck(ctx, deck.ID)
	if err != nil {
		t.Fatalf("get deck: %v", err)
	}
	if len(got.Notes) != len(deck.Notes) {
		t.Fatalf("notes = %d, want %d", len(got.Notes), len(deck.Notes))
	}
	if len(got.Notes[0].Cards) == 0 {
		t.Fatal("expected cards on loaded notes")
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

func TestDecksWithStats(t *testing.T) {
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

	// Count cards
	var totalCards int
	for _, note := range deck.Notes {
		totalCards += len(note.Cards)
	}

	decks, err := store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks: %v", err)
	}
	if len(decks) == 0 {
		t.Fatal("expected at least one deck")
	}
	d := decks[0]
	if d.TotalCards != totalCards {
		t.Fatalf("total cards = %d, want %d", d.TotalCards, totalCards)
	}
	if d.DueCards != totalCards {
		t.Fatalf("initial due cards = %d, want %d", d.DueCards, totalCards)
	}

	// Review one card to make it not due
	card := deck.Notes[0].Cards[0]
	now := time.Now().UTC()
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

	decks, err = store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks: %v", err)
	}
	d = decks[0]
	if d.DueCards != totalCards-1 {
		t.Fatalf("due cards after review = %d, want %d", d.DueCards, totalCards-1)
	}
	if d.ReviewsToday != 1 {
		t.Fatalf("reviews today = %d, want 1", d.ReviewsToday)
	}
	if d.SuccessRate != 1 {
		t.Fatalf("success rate = %.2f, want 1", d.SuccessRate)
	}

	if err := store.RecordReview(ctx, core.ReviewResult{
		CardID:   deck.Notes[1].Cards[0].ID,
		Grade:    core.GradeAgain,
		Reviewed: now.Add(time.Minute),
		Next: core.ReviewState{
			CardID:   deck.Notes[1].Cards[0].ID,
			Due:      now,
			Interval: 0,
			Ease:     2.5,
			Reviews:  1,
		},
	}); err != nil {
		t.Fatalf("record second review: %v", err)
	}
	decks, err = store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks after second review: %v", err)
	}
	d = decks[0]
	if d.ReviewsToday != 2 {
		t.Fatalf("reviews today after second review = %d, want 2", d.ReviewsToday)
	}
	if d.SuccessRate != 0.5 {
		t.Fatalf("success rate after again = %.2f, want 0.50", d.SuccessRate)
	}
}

func TestCardBookmarkPersistsAndCountsInStats(t *testing.T) {
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
	if err := store.SetCardBookmark(ctx, card.ID, true); err != nil {
		t.Fatalf("set bookmark: %v", err)
	}
	cards, err := store.DueCards(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	found := false
	for _, due := range cards {
		if due.ID == card.ID {
			found = true
			if !due.Bookmarked {
				t.Fatal("expected due card bookmark to be true")
			}
		}
	}
	if !found {
		t.Fatalf("card %s not found in due cards", card.ID)
	}
	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.BookmarkedCards != 1 {
		t.Fatalf("bookmarked cards = %d, want 1", stats.BookmarkedCards)
	}
}

func TestUndoLastReviewRestoresNewCardToDueQueue(t *testing.T) {
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
	if err := store.RecordReview(ctx, core.ReviewResult{
		CardID:   card.ID,
		Grade:    core.GradeGood,
		Reviewed: now,
		Next: core.ReviewState{
			CardID:   card.ID,
			Due:      now.Add(24 * time.Hour),
			Interval: 24 * time.Hour,
			Reviews:  1,
			Ease:     2.5,
		},
	}); err != nil {
		t.Fatalf("record review: %v", err)
	}
	if err := store.UndoLastReview(ctx, card.ID); err != nil {
		t.Fatalf("undo review: %v", err)
	}
	cards, err := store.DueCards(ctx, now, 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	for _, due := range cards {
		if due.ID == card.ID {
			return
		}
	}
	t.Fatalf("card %s should be due after undoing its first review", card.ID)
}

func TestUndoLastReviewRestoresPreviousReviewState(t *testing.T) {
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
	first := core.ReviewState{CardID: card.ID, Due: now.Add(24 * time.Hour), Interval: 24 * time.Hour, Reviews: 1, Ease: 2.5}
	second := core.ReviewState{CardID: card.ID, Due: now.Add(72 * time.Hour), Interval: 72 * time.Hour, Reviews: 2, Ease: 2.6}
	if err := store.RecordReview(ctx, core.ReviewResult{CardID: card.ID, Grade: core.GradeGood, Reviewed: now, Next: first}); err != nil {
		t.Fatalf("record first review: %v", err)
	}
	if err := store.RecordReview(ctx, core.ReviewResult{CardID: card.ID, Grade: core.GradeEasy, Reviewed: now.Add(time.Minute), Next: second}); err != nil {
		t.Fatalf("record second review: %v", err)
	}
	if err := store.UndoLastReview(ctx, card.ID); err != nil {
		t.Fatalf("undo review: %v", err)
	}
	state, err := store.GetReviewState(ctx, card.ID)
	if err != nil {
		t.Fatalf("review state: %v", err)
	}
	if state.Reviews != 1 || !state.Due.Equal(first.Due) {
		t.Fatalf("state after undo = reviews %d due %s, want reviews 1 due %s", state.Reviews, state.Due, first.Due)
	}
	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.TotalReviews != 1 {
		t.Fatalf("total reviews after undo = %d, want 1", stats.TotalReviews)
	}
}

func TestStatisticsDailyProgressAndStreak(t *testing.T) {
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
	now := time.Now().UTC()
	cardA := deck.Notes[0].Cards[0]
	cardB := deck.Notes[0].Cards[1]
	for _, item := range []struct {
		cardID string
		when   time.Time
	}{
		{cardA.ID, now.AddDate(0, 0, -1)},
		{cardB.ID, now},
	} {
		if err := store.RecordReview(ctx, core.ReviewResult{
			CardID:   item.cardID,
			Grade:    core.GradeGood,
			Reviewed: item.when,
			Next:     core.ReviewState{CardID: item.cardID, Due: item.when.Add(24 * time.Hour), Interval: 24 * time.Hour, Reviews: 1, Ease: 2.5},
		}); err != nil {
			t.Fatalf("record review: %v", err)
		}
	}
	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.ReviewsToday != 1 {
		t.Fatalf("reviews today = %d, want 1", stats.ReviewsToday)
	}
	if stats.DailyGoal != 10 {
		t.Fatalf("daily goal = %d, want 10", stats.DailyGoal)
	}
	if stats.CurrentStreak < 2 {
		t.Fatalf("current streak = %d, want at least 2", stats.CurrentStreak)
	}
}

func TestDailyGoalPersistsInStatistics(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.DailyGoal != 10 {
		t.Fatalf("default daily goal = %d, want 10", stats.DailyGoal)
	}

	if err := store.SetDailyGoal(ctx, 25); err != nil {
		t.Fatalf("set daily goal: %v", err)
	}
	stats, err = store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics after set: %v", err)
	}
	if stats.DailyGoal != 25 {
		t.Fatalf("daily goal = %d, want 25", stats.DailyGoal)
	}

	if err := store.SetDailyGoal(ctx, -3); err != nil {
		t.Fatalf("set daily goal floor: %v", err)
	}
	stats, err = store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics after floor: %v", err)
	}
	if stats.DailyGoal != 1 {
		t.Fatalf("daily goal floor = %d, want 1", stats.DailyGoal)
	}
}

func TestLeechDetectionAfterThreeAgainGrades(t *testing.T) {
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

	for i := 0; i < 3; i++ {
		if err := store.RecordReview(ctx, core.ReviewResult{
			CardID:   card.ID,
			Grade:    core.GradeAgain,
			Reviewed: now.Add(time.Duration(i) * time.Minute),
			Next:     core.ReviewState{CardID: card.ID, Due: now, Interval: 0, Reviews: i + 1, Ease: 2.5},
		}); err != nil {
			t.Fatalf("record review %d: %v", i+1, err)
		}
	}

	cards, err := store.DueCards(ctx, now, 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	found := false
	for _, due := range cards {
		if due.ID == card.ID {
			found = true
			if !due.Leech {
				t.Fatal("card should be flagged as leech after 3 Again grades")
			}
		}
	}
	if !found {
		t.Fatalf("card %s not found in due cards", card.ID)
	}

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.LeechCards != 1 {
		t.Fatalf("leech cards = %d, want 1", stats.LeechCards)
	}
}

func TestLeechResetOnGoodGrade(t *testing.T) {
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

	for i := 0; i < 3; i++ {
		if err := store.RecordReview(ctx, core.ReviewResult{
			CardID:   card.ID,
			Grade:    core.GradeAgain,
			Reviewed: now.Add(time.Duration(i) * time.Minute),
			Next:     core.ReviewState{CardID: card.ID, Due: now, Interval: 0, Reviews: i + 1, Ease: 2.5},
		}); err != nil {
			t.Fatalf("record again %d: %v", i+1, err)
		}
	}

	if err := store.RecordReview(ctx, core.ReviewResult{
		CardID:   card.ID,
		Grade:    core.GradeGood,
		Reviewed: now.Add(4 * time.Minute),
		Next:     core.ReviewState{CardID: card.ID, Due: now.Add(24 * time.Hour), Interval: 24 * time.Hour, Reviews: 4, Ease: 2.5},
	}); err != nil {
		t.Fatalf("record good: %v", err)
	}

	cards, err := store.DueCards(ctx, now.Add(25*time.Hour), 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	for _, due := range cards {
		if due.ID == card.ID && due.Leech {
			t.Fatal("card should not be leech after Good grade")
		}
	}

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.LeechCards != 0 {
		t.Fatalf("leech cards = %d, want 0", stats.LeechCards)
	}
}

func TestDueCardsBookmarkedReturnsOnlyBookmarked(t *testing.T) {
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
	cardA := deck.Notes[0].Cards[0]

	if err := store.SetCardBookmark(ctx, cardA.ID, true); err != nil {
		t.Fatalf("set bookmark: %v", err)
	}

	bookmarked, err := store.DueCardsBookmarked(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("due cards bookmarked: %v", err)
	}
	if len(bookmarked) != 1 {
		t.Fatalf("bookmarked due cards = %d, want 1", len(bookmarked))
	}
	if bookmarked[0].ID != cardA.ID {
		t.Fatalf("bookmarked card = %s, want %s", bookmarked[0].ID, cardA.ID)
	}

	allCards, err := store.DueCards(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	if len(allCards) < 2 {
		t.Fatalf("all due cards = %d, want at least 2", len(allCards))
	}
}

func TestSuspendedCardsAreFilteredAndCounted(t *testing.T) {
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
	initialDue, err := store.DueCards(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("initial due cards: %v", err)
	}
	if err := store.SetCardBookmark(ctx, card.ID, true); err != nil {
		t.Fatalf("set bookmark: %v", err)
	}
	if err := store.SetCardSuspended(ctx, card.ID, true); err != nil {
		t.Fatalf("set suspended: %v", err)
	}

	cards, err := store.DueCards(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	for _, due := range cards {
		if due.ID == card.ID {
			t.Fatal("suspended card should be excluded from normal due cards")
		}
	}

	bookmarked, err := store.DueCardsBookmarked(ctx, time.Now(), 100)
	if err != nil {
		t.Fatalf("bookmarked due cards: %v", err)
	}
	for _, due := range bookmarked {
		if due.ID == card.ID {
			t.Fatal("suspended card should be excluded from bookmarked due cards")
		}
	}

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.SuspendedCards != 1 {
		t.Fatalf("suspended cards = %d, want 1", stats.SuspendedCards)
	}

	decks, err := store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks: %v", err)
	}
	if len(decks) == 0 {
		t.Fatal("expected decks")
	}
	if decks[0].DueCards != len(initialDue)-1 {
		t.Fatalf("deck due cards = %d, want %d", decks[0].DueCards, len(initialDue)-1)
	}
}

func TestStatisticsBookmarkedDueCount(t *testing.T) {
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
	cardA := deck.Notes[0].Cards[0]

	if err := store.SetCardBookmark(ctx, cardA.ID, true); err != nil {
		t.Fatalf("set bookmark: %v", err)
	}

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.BookmarkedDue < 1 {
		t.Fatalf("bookmarked due = %d, want at least 1", stats.BookmarkedDue)
	}
}
