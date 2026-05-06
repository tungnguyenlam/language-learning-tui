package sqlite

import (
	"context"
	"strings"
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
	var cardA, cardB core.Card
	found := 0
	for _, n := range deck.Notes {
		for _, c := range n.Cards {
			if found == 0 {
				cardA = c
				found++
			} else if found == 1 {
				cardB = c
				found++
				break
			}
		}
		if found == 2 {
			break
		}
	}
	if found < 2 {
		t.Fatalf("need at least 2 cards for streak test")
	}

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

func TestUndoLastReviewRestoresLeechFlag(t *testing.T) {
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

	// 1. Make it a leech (3 Agains)
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

	stats, err := store.Statistics(ctx)
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if stats.LeechCards != 1 {
		t.Fatalf("expected 1 leech, got %d", stats.LeechCards)
	}

	// 2. Clear leech with a Good grade
	if err := store.RecordReview(ctx, core.ReviewResult{
		CardID:   card.ID,
		Grade:    core.GradeGood,
		Reviewed: now.Add(4 * time.Minute),
		Next:     core.ReviewState{CardID: card.ID, Due: now.Add(24 * time.Hour), Interval: 24 * time.Hour, Reviews: 4, Ease: 2.5},
	}); err != nil {
		t.Fatalf("record good: %v", err)
	}

	stats, err = store.Statistics(ctx)
	if stats.LeechCards != 0 {
		t.Fatalf("expected 0 leeches after Good grade, got %d", stats.LeechCards)
	}

	// 3. Undo the Good grade - should be a leech again!
	if err := store.UndoLastReview(ctx, card.ID); err != nil {
		t.Fatalf("undo review: %v", err)
	}

	stats, err = store.Statistics(ctx)
	if stats.LeechCards != 1 {
		t.Errorf("expected card to be leech again after undoing the Good grade, but got %d leeches", stats.LeechCards)
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

func TestCards_DeckFilter(t *testing.T) {
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

	cards, err := store.Cards(ctx, deck.ID, "")
	if err != nil {
		t.Fatalf("cards: %v", err)
	}
	if len(cards) == 0 {
		t.Fatal("expected cards for deck")
	}
	for _, card := range cards {
		if card.DeckID != deck.ID {
			t.Fatalf("expected deck %s, got %s", deck.ID, card.DeckID)
		}
	}
}

func TestCards_SearchFilter(t *testing.T) {
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
	firstCard := deck.Notes[0].Cards[0]
	searchTerm := firstCard.Prompt[:5]

	cards, err := store.Cards(ctx, "", searchTerm)
	if err != nil {
		t.Fatalf("cards: %v", err)
	}
	if len(cards) == 0 {
		t.Fatal("expected cards containing search term")
	}
	for _, card := range cards {
		promptMatch := strings.Contains(strings.ToLower(card.Prompt), strings.ToLower(searchTerm))
		answerMatch := strings.Contains(strings.ToLower(card.Answer), strings.ToLower(searchTerm))
		if !promptMatch && !answerMatch {
			t.Fatalf("card %s does not contain search term %s", card.ID, searchTerm)
		}
	}
}

func TestReviewHistoryReturnsRecentReviewsForCard(t *testing.T) {
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
	now := time.Date(2026, 4, 30, 8, 0, 0, 0, time.UTC)
	reviews := []core.ReviewResult{
		{
			CardID:   card.ID,
			Grade:    core.GradeAgain,
			Reviewed: now,
			Next:     core.ReviewState{CardID: card.ID, Due: now.Add(time.Hour), Interval: time.Hour, Reviews: 1, Lapses: 1},
		},
		{
			CardID:   card.ID,
			Grade:    core.GradeGood,
			Reviewed: now.Add(time.Minute),
			Next:     core.ReviewState{CardID: card.ID, Due: now.Add(24 * time.Hour), Interval: 24 * time.Hour, Reviews: 2, Lapses: 1},
		},
	}
	for _, review := range reviews {
		if err := store.RecordReview(ctx, review); err != nil {
			t.Fatalf("record review: %v", err)
		}
	}

	history, err := store.ReviewHistory(ctx, card.ID, 5)
	if err != nil {
		t.Fatalf("review history: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("history len = %d, want 2", len(history))
	}
	if history[0].Grade != core.GradeGood || history[0].Reviews != 2 {
		t.Fatalf("newest history = %#v, want good review count 2", history[0])
	}
	if history[1].Grade != core.GradeAgain || history[1].Interval != time.Hour {
		t.Fatalf("oldest history = %#v, want again with 1h interval", history[1])
	}
}

func TestAudioPersistence(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	note := core.Note{
		ID:     "n1",
		DeckID: "d1",
		Front:  "Front",
		Back:   "Back",
		Audio:  "audio.mp3",
	}
	note.Cards = []core.Card{
		{
			ID:     "c1",
			NoteID: "n1",
			DeckID: "d1",
			Kind:   core.CardKindFlashcard,
			Prompt: "Front",
			Answer: "Back",
			Audio:  "audio.mp3",
		},
	}

	if err := store.UpsertDeck(ctx, core.Deck{ID: "d1", Name: "D1", Notes: []core.Note{note}}); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}

	cards, err := store.DueCards(ctx, time.Now(), 10)
	if err != nil {
		t.Fatalf("due cards: %v", err)
	}
	if cards[0].Audio != "audio.mp3" {
		t.Fatalf("audio = %s, want audio.mp3", cards[0].Audio)
	}

	browserCards, err := store.Cards(ctx, "d1", "")
	if err != nil {
		t.Fatalf("cards: %v", err)
	}
	if browserCards[0].Audio != "audio.mp3" {
		t.Fatalf("browser audio = %s, want audio.mp3", browserCards[0].Audio)
	}
}

func TestSetCardsTags(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	deck := core.Deck{ID: "d1", Name: "Test Deck"}
	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}

	note := core.Note{
		ID:     "n1",
		DeckID: "d1",
		Front:  "Front",
		Back:   "Back",
		Tags:   []string{"initial"},
	}
	note.Cards = []core.Card{
		{ID: "c1", NoteID: "n1", DeckID: "d1", Kind: core.CardKindFlashcard, Prompt: "Front", Answer: "Back", Tags: []string{"initial"}},
		{ID: "c2", NoteID: "n1", DeckID: "d1", Kind: core.CardKindFlashcard, Prompt: "Back", Answer: "Front", Tags: []string{"initial"}},
	}

	if err := store.UpsertNote(ctx, note); err != nil {
		t.Fatalf("upsert note: %v", err)
	}

	newTags := []string{"tag1", "tag2"}
	if err := store.SetCardsTags(ctx, []string{"c1"}, newTags); err != nil {
		t.Fatalf("set cards tags: %v", err)
	}

	// Verify both cards and note are updated (since they share noteID)
	c1, err := store.cardsForNote(ctx, "n1")
	if err != nil {
		t.Fatalf("get cards: %v", err)
	}
	for _, c := range c1 {
		if len(c.Tags) != 2 || c.Tags[0] != "tag1" || c.Tags[1] != "tag2" {
			t.Errorf("card %s tags = %v, want %v", c.ID, c.Tags, newTags)
		}
	}

	// Verify through Cards browser query too
	cards, err := store.Cards(ctx, "d1", "")
	if err != nil {
		t.Fatalf("browser cards: %v", err)
	}
	for _, c := range cards {
		if len(c.Tags) != 2 || c.Tags[0] != "tag1" || c.Tags[1] != "tag2" {
			t.Errorf("browser card %s tags = %v, want %v", c.ID, c.Tags, newTags)
		}
	}

	// Verify searching by tag
	searchResult, err := store.Cards(ctx, "d1", "tag1")
	if err != nil {
		t.Fatalf("search browser cards: %v", err)
	}
	if len(searchResult) != 2 {
		t.Errorf("search result length = %d, want 2", len(searchResult))
	}
}

func TestStoreDeckTags(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	deck := core.Deck{
		ID:          "test-deck-tags",
		Name:        "Test Deck With Tags",
		Description: "A test deck with tags",
		Tags:        []string{"german", "beginner", "vocabulary"},
	}

	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck with tags: %v", err)
	}

	// Retrieve the deck and verify tags are preserved
	retrieved, err := store.GetDeck(ctx, deck.ID)
	if err != nil {
		t.Fatalf("get deck with tags: %v", err)
	}

	if retrieved.ID != deck.ID {
		t.Errorf("deck ID mismatch: got %q, want %q", retrieved.ID, deck.ID)
	}
	if retrieved.Name != deck.Name {
		t.Errorf("deck name mismatch: got %q, want %q", retrieved.Name, deck.Name)
	}
	if retrieved.Description != deck.Description {
		t.Errorf("deck description mismatch: got %q, want %q", retrieved.Description, deck.Description)
	}
	if len(retrieved.Tags) != len(deck.Tags) {
		t.Errorf("deck tags length mismatch: got %d, want %d", len(retrieved.Tags), len(deck.Tags))
	}
	for i, tag := range deck.Tags {
		if retrieved.Tags[i] != tag {
			t.Errorf("deck tag mismatch at index %d: got %q, want %q", i, retrieved.Tags[i], tag)
		}
	}

	// Test listing all decks and verify tags are included
	decks, err := store.Decks(ctx)
	if err != nil {
		t.Fatalf("list decks: %v", err)
	}
	if len(decks) != 1 {
		t.Fatalf("expected 1 deck, got %d", len(decks))
	}

	retrievedDeck := decks[0]
	if len(retrievedDeck.Tags) != len(deck.Tags) {
		t.Errorf("listed deck tags length mismatch: got %d, want %d", len(retrievedDeck.Tags), len(deck.Tags))
	}
	for i, tag := range deck.Tags {
		if retrievedDeck.Tags[i] != tag {
			t.Errorf("listed deck tag mismatch at index %d: got %q, want %q", i, retrievedDeck.Tags[i], tag)
		}
	}
}

func TestCleanupTags(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	deckID := "deck1"
	deck := core.Deck{
		ID:   deckID,
		Name: "Deck 1",
		Tags: []string{"unused", "will-be-kept"},
	}
	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}

	note := core.Note{
		ID:     "n1",
		DeckID: deckID,
		Front:  "Front",
		Back:   "Back",
		Tags:   []string{"will-be-kept", "new-tag"},
	}
	note.Cards = []core.Card{
		{
			ID:     "c1",
			NoteID: "n1",
			DeckID: deckID,
			Kind:   core.CardKindFlashcard,
			Prompt: "Front",
			Answer: "Back",
			Tags:   []string{"will-be-kept", "new-tag"},
		},
	}
	if err := store.UpsertNote(ctx, note); err != nil {
		t.Fatalf("upsert note: %v", err)
	}

	// Initial cleanup
	if err := store.CleanupTags(ctx, deckID); err != nil {
		t.Fatalf("cleanup tags: %v", err)
	}

	retrieved, err := store.GetDeck(ctx, deckID)
	if err != nil {
		t.Fatalf("get deck: %v", err)
	}

	// Tags should be "will-be-kept" and "new-tag" (order may vary)
	tagMap := make(map[string]bool)
	for _, tag := range retrieved.Tags {
		tagMap[tag] = true
	}

	if len(tagMap) != 2 {
		t.Errorf("expected 2 tags, got %v", retrieved.Tags)
	}
	if !tagMap["will-be-kept"] {
		t.Errorf("missing tag 'will-be-kept'")
	}
	if !tagMap["new-tag"] {
		t.Errorf("missing tag 'new-tag'")
	}
	if tagMap["unused"] {
		t.Errorf("tag 'unused' should have been removed")
	}
}
