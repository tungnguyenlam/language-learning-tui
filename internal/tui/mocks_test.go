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
func (m *mockRepo) RecordReview(ctx context.Context, result core.ReviewResult) error { return nil }
func (m *mockRepo) DueCards(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	return m.dueCards, nil
}
func (m *mockRepo) GetReviewState(ctx context.Context, cardID string) (core.ReviewState, error) {
	return core.ReviewState{CardID: cardID}, nil
}
func (m *mockRepo) Statistics(ctx context.Context) (core.Statistics, error) {
	return core.Statistics{}, nil
}

type mockScheduler struct{}

func (m *mockScheduler) Review(state core.ReviewState, grade core.ReviewGrade, now time.Time) (core.ReviewState, error) {
	return state, nil
}
