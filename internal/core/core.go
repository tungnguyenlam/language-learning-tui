package core

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrNoReviewToUndo = errors.New("no review to undo")

type CardKind string

const (
	CardKindFlashcard CardKind = "flashcard"
	CardKindMCQ       CardKind = "mcq"
	CardKindCloze     CardKind = "cloze"
)

type ReviewGrade string

const (
	GradeAgain ReviewGrade = "again"
	GradeHard  ReviewGrade = "hard"
	GradeGood  ReviewGrade = "good"
	GradeEasy  ReviewGrade = "easy"
)

type Deck struct {
	ID                string
	Name              string
	Description       string
	Tags              []string
	Notes             []Note
	TotalCards        int
	NewCards          int
	DueCards          int
	ReviewsToday      int
	SuccessRate       float64
	NewCardsPerDay    int
	ReviewLimitPerDay int
}

type Note struct {
	ID       string
	DeckID   string
	Type     string
	Front    string
	Back     string
	Extra    string
	Audio    string
	Tags     []string
	Examples []string
	Choices  []string
	Cards    []Card
}

type Card struct {
	ID         string
	NoteID     string
	DeckID     string
	Kind       CardKind
	Prompt     string
	Answer     string
	Extra      string
	Choices    []string
	Audio      string
	Tags       []string
	Bookmarked bool
	Leech      bool
	Suspended  bool
	Mature     bool
	Interval   time.Duration
	Reviews    int
}

type ReviewState struct {
	CardID     string
	Due        time.Time
	LastReview time.Time
	Interval   time.Duration
	Stability  float64
	Difficulty float64
	Ease       float64
	Reviews    int
	Lapses     int
}

type ReviewResult struct {
	CardID   string
	Grade    ReviewGrade
	Reviewed time.Time
	Next     ReviewState
}

type ReviewLog struct {
	CardID   string
	Grade    ReviewGrade
	Reviewed time.Time
	Due      time.Time
	Interval time.Duration
	Reviews  int
	Lapses   int
}

type Statistics struct {
	TotalCards      int
	NewCards        int
	YoungCards      int
	MatureCards     int
	BookmarkedCards int
	BookmarkedDue   int
	LeechCards      int
	SuspendedCards  int
	TotalReviews    int
	ReviewsToday    int
	DailyGoal       int
	CurrentStreak   int
	TotalDecks      int
	ActiveDecks     int
	SuccessRate     float64
	Grades          map[ReviewGrade]int
}

type Repository interface {
	UpsertDeck(ctx context.Context, deck Deck) error
	GetDeck(ctx context.Context, id string) (Deck, error)
	Decks(ctx context.Context) ([]Deck, error)
	RecordReview(ctx context.Context, result ReviewResult) error
	UndoLastReview(ctx context.Context, cardID string) error
	DueCards(ctx context.Context, now time.Time, limit int) ([]Card, error)
	DueCardsBookmarked(ctx context.Context, now time.Time, limit int) ([]Card, error)
	GetReviewState(ctx context.Context, cardID string) (ReviewState, error)
	SetCardBookmark(ctx context.Context, cardID string, bookmarked bool) error
	SetCardSuspended(ctx context.Context, cardID string, suspended bool) error
	SetDailyGoal(ctx context.Context, goal int) error
	SetDeckLimits(ctx context.Context, deckID string, newLimit, reviewLimit int) error
	Statistics(ctx context.Context) (Statistics, error)
	DeckStatistics(ctx context.Context, deckID string) (Statistics, error)
	Cards(ctx context.Context, deckID string, search string, tag string) ([]Card, error)
	ReviewsPerDay(ctx context.Context, days int) (map[string]int, error)
	RecentDecks(ctx context.Context, limit int) ([]string, error)
	ReviewHistory(ctx context.Context, cardID string, limit int) ([]ReviewLog, error)
	DeleteCard(ctx context.Context, cardID string) error
	DeleteDecks(ctx context.Context, ids []string) error
	MergeDecks(ctx context.Context, sourceIDs []string, targetID string) error
	SetCardKind(ctx context.Context, cardID string, kind CardKind) error
	GetNote(ctx context.Context, noteID string) (Note, error)
	UpsertNote(ctx context.Context, note Note) error
	SetCardTags(ctx context.Context, cardID string, tags []string) error
	SetCardsTags(ctx context.Context, cardIDs []string, tags []string) error
	CleanupTags(ctx context.Context, deckID string) error
	Reset(ctx context.Context) error
}

type Scheduler interface {
	Review(state ReviewState, grade ReviewGrade, now time.Time) (ReviewState, error)
	Predict(state ReviewState, now time.Time) map[ReviewGrade]time.Duration
}

func ValidateCard(card Card) error {
	if strings.TrimSpace(card.ID) == "" {
		return errors.New("card id is required")
	}
	if strings.TrimSpace(card.NoteID) == "" {
		return errors.New("note id is required")
	}
	if strings.TrimSpace(card.Prompt) == "" {
		return errors.New("prompt is required")
	}
	if strings.TrimSpace(card.Answer) == "" {
		return errors.New("answer is required")
	}
	if card.Kind == CardKindMCQ && len(card.Choices) < 2 {
		return errors.New("mcq cards require at least two choices")
	}
	if card.Kind != CardKindFlashcard && card.Kind != CardKindMCQ && card.Kind != CardKindCloze {
		return errors.New("unsupported card kind")
	}
	return nil
}

func NewReviewState(cardID string, now time.Time) ReviewState {
	return ReviewState{
		CardID:     cardID,
		Due:        now,
		Interval:   0,
		Difficulty: 0,
		Ease:       2.5,
	}
}
