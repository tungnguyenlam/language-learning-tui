package tui

import (
	"context"
	"strings"
	"time"

	"deutsch-tui/internal/core"
)

type mockRepo struct {
	decks        []core.Deck
	dueCards     []core.Card
	upsertedDeck core.Deck
	bookmarks    map[string]bool
	suspended    map[string]bool
	reviews      []core.ReviewResult
	dailyGoal    int
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
	var cards []core.Card
	for _, card := range m.dueCards {
		if !card.Suspended {
			cards = append(cards, card)
		}
	}
	return cards, nil
}
func (m *mockRepo) DueCardsBookmarked(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	var bookmarked []core.Card
	for _, card := range m.dueCards {
		if card.Bookmarked && !card.Suspended {
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
func (m *mockRepo) SetCardSuspended(ctx context.Context, cardID string, suspended bool) error {
	if m.suspended == nil {
		m.suspended = make(map[string]bool)
	}
	m.suspended[cardID] = suspended
	for i := range m.dueCards {
		if m.dueCards[i].ID == cardID {
			m.dueCards[i].Suspended = suspended
		}
	}
	return nil
}
func (m *mockRepo) SetDailyGoal(ctx context.Context, goal int) error {
	if goal < 1 {
		goal = 1
	}
	m.dailyGoal = goal
	return nil
}
func (m *mockRepo) Statistics(ctx context.Context) (core.Statistics, error) {
	goal := m.dailyGoal
	if goal == 0 {
		goal = 10
	}
	stats := core.Statistics{DailyGoal: goal, Grades: map[core.ReviewGrade]int{}}
	for _, bookmarked := range m.bookmarks {
		if bookmarked {
			stats.BookmarkedCards++
		}
	}
	for _, suspended := range m.suspended {
		if suspended {
			stats.SuspendedCards++
		}
	}
	return stats, nil
}
func (m *mockRepo) Cards(ctx context.Context, deckID string, search string) ([]core.Card, error) {
	var cards []core.Card
	for _, card := range m.dueCards {
		if deckID != "" && card.DeckID != deckID {
			continue
		}
		if search != "" && !containsIgnoreCase(card.Prompt, search) && !containsIgnoreCase(card.Answer, search) {
			continue
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (m *mockRepo) ReviewsPerDay(ctx context.Context, days int) (map[string]int, error) {
	return map[string]int{}, nil
}

func (m *mockRepo) DeleteCard(ctx context.Context, cardID string) error {
	for i, card := range m.dueCards {
		if card.ID == cardID {
			m.dueCards = append(m.dueCards[:i], m.dueCards[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockRepo) ReviewHistory(ctx context.Context, cardID string, limit int) ([]core.ReviewLog, error) {
	if limit <= 0 {
		limit = 10
	}
	var logs []core.ReviewLog
	for i := len(m.reviews) - 1; i >= 0 && len(logs) < limit; i-- {
		review := m.reviews[i]
		if review.CardID != cardID {
			continue
		}
		logs = append(logs, core.ReviewLog{
			CardID:   review.CardID,
			Grade:    review.Grade,
			Reviewed: review.Reviewed,
			Due:      review.Next.Due,
			Interval: review.Next.Interval,
			Reviews:  review.Next.Reviews,
			Lapses:   review.Next.Lapses,
		})
	}
	return logs, nil
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

type mockScheduler struct{}

func (m *mockScheduler) Review(state core.ReviewState, grade core.ReviewGrade, now time.Time) (core.ReviewState, error) {
	return state, nil
}

func (m *mockScheduler) Predict(state core.ReviewState, now time.Time) map[core.ReviewGrade]time.Duration {
	return map[core.ReviewGrade]time.Duration{
		core.GradeAgain: time.Minute,
		core.GradeHard:  24 * time.Hour,
		core.GradeGood:  3 * 24 * time.Hour,
		core.GradeEasy:  7 * 24 * time.Hour,
	}
}
