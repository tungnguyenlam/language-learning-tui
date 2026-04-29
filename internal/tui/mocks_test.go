package tui

import (
	"context"
	"time"

	"deutsch-tui/internal/core"
)

type mockRepo struct {
	decks        []core.Deck
	dueCards     []core.Card
	upsertedDeck core.Deck
	bookmarks    map[string]bool
	reviews      []core.ReviewResult
}

func (m *mockRepo) UpsertDeck(ctx context.Context, deck core.Deck) error {
	m.upsertedDeck = deck
	replaced := false
	for i := range m.decks {
		if m.decks[i].ID == deck.ID {
			m.decks[i] = deck
			replaced = true
			break
		}
	}
	if !replaced {
		m.decks = append(m.decks, deck)
	}
	for _, note := range deck.Notes {
		m.dueCards = append(m.dueCards, note.Cards...)
	}
	return nil
}
func (m *mockRepo) GetDeck(ctx context.Context, id string) (core.Deck, error) {
	for _, deck := range m.decks {
		if deck.ID == id {
			return deck, nil
		}
	}
	return core.Deck{ID: id, Name: "Mock Deck"}, nil
}
func (m *mockRepo) Decks(ctx context.Context) ([]core.Deck, error) {
	if len(m.decks) > 0 {
		return m.decks, nil
	}
	return []core.Deck{{ID: "mock-1", Name: "Mock Deck"}}, nil
}
func (m *mockRepo) RecordReview(ctx context.Context, result core.ReviewResult) error {
	m.reviews = append(m.reviews, result)
	return nil
}
func (m *mockRepo) UndoLastReview(ctx context.Context, cardID string) error {
	if len(m.reviews) == 0 {
		return core.ErrNoReviewToUndo
	}
	m.reviews = m.reviews[:len(m.reviews)-1]
	return nil
}
func (m *mockRepo) DueCards(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	return m.dueCards, nil
}
func (m *mockRepo) DueCardsBookmarked(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	var bookmarked []core.Card
	for _, card := range m.dueCards {
		if card.Bookmarked {
			bookmarked = append(bookmarked, card)
		}
	}
	return bookmarked, nil
}
func (m *mockRepo) GetReviewState(ctx context.Context, cardID string) (core.ReviewState, error) {
	return core.ReviewState{CardID: cardID}, nil
}
func (m *mockRepo) SetCardBookmark(ctx context.Context, cardID string, bookmarked bool) error {
	if m.bookmarks == nil {
		m.bookmarks = make(map[string]bool)
	}
	m.bookmarks[cardID] = bookmarked
	for i := range m.dueCards {
		if m.dueCards[i].ID == cardID {
			m.dueCards[i].Bookmarked = bookmarked
		}
	}
	return nil
}
func (m *mockRepo) Statistics(ctx context.Context) (core.Statistics, error) {
	stats := core.Statistics{DailyGoal: 10, Grades: map[core.ReviewGrade]int{}}
	for _, bookmarked := range m.bookmarks {
		if bookmarked {
			stats.BookmarkedCards++
		}
	}
	return stats, nil
}

type mockScheduler struct{}

func (m *mockScheduler) Review(state core.ReviewState, grade core.ReviewGrade, now time.Time) (core.ReviewState, error) {
	return state, nil
}
