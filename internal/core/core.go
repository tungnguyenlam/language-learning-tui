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
	ID        string
	DeckID    string
	Type      string
	Front     string
	Back      string
	Extra     string
	Hint      string
	Audio     string
	Tags      []string
	Examples  []string
	Choices   []string
	Cards     []Card
	CreatedAt time.Time
}

type Card struct {
	ID           string
	NoteID       string
	DeckID       string
	Kind         CardKind
	Prompt       string
	Answer       string
	Extra        string
	Hint         string
	Choices      []string
	Audio        string
	Tags         []string
	Bookmarked   bool
	Leech        bool
	Suspended    bool
	Mature       bool
	Interval     time.Duration
	Reviews      int
	Lapses       int
	Ease         float64
	Due          time.Time
	LastReviewed time.Time
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
	Next24hDue      int
	LeechCards      int

	SuspendedCards   int
	TotalReviews     int
	ReviewsToday     int
	DailyGoal        int
	CurrentStreak    int
	TotalDecks       int
	ActiveDecks      int
	SuccessRate      float64
	CardsAddedPerDay map[string]int
	Grades           map[ReviewGrade]int
	CardTypes        map[CardKind]int
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
	// CardsWithFlag returns cards matching a card_flags filter:
	// "bookmarked", "suspended", "leech", or "flagged" (any of the three).
	CardsWithFlag(ctx context.Context, deckID string, flag string) ([]Card, error)
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
	GetSetting(ctx context.Context, key string) (string, error)
	SetSetting(ctx context.Context, key string, value string) error
}

// BackupInfo describes the contents of a backup file, both when one is
// written and when one is inspected before being restored.
type BackupInfo struct {
	Path          string
	CreatedAt     time.Time
	SchemaVersion int
	// Rows counts the rows per table, keyed by table name.
	Rows      map[string]int
	TotalRows int
}

// BackupRepository is an optional capability. It is kept out of Repository
// so in-memory and test doubles do not have to implement it; callers type
// assert and degrade gracefully when a store cannot back itself up.
type BackupRepository interface {
	Backup(ctx context.Context, destPath string) (BackupInfo, error)
	Restore(ctx context.Context, srcPath string) (BackupInfo, error)
	InspectBackup(ctx context.Context, srcPath string) (BackupInfo, error)
}

type Scheduler interface {
	Review(state ReviewState, grade ReviewGrade, now time.Time) (ReviewState, error)
	Predict(state ReviewState, now time.Time) map[ReviewGrade]time.Duration
}

type DictionaryEntry struct {
	ID          string
	Word        string
	Translation string
	WordClass   string
	Gender      string
	Forms       string
	Examples    []string
	Tags        []string
}

type DictionaryRepository interface {
	Search(ctx context.Context, query string, limit int) ([]DictionaryEntry, error)
	GetEntry(ctx context.Context, id string) (DictionaryEntry, error)
	FindRelatedEntries(ctx context.Context, word string, limit int) ([]DictionaryEntry, error)
	RandomEntries(ctx context.Context, limit int) ([]DictionaryEntry, error)
	ImportEntries(ctx context.Context, entries []DictionaryEntry) error
	DictionaryCount(ctx context.Context) (int, error)
	Exists(ctx context.Context, word string) (bool, error)
}

// ConsolidateDictionaryEntries merges dictionary entries sharing the same German headword
// and gender, combining their translations, word classes, forms, tags, and examples.
func ConsolidateDictionaryEntries(entries []DictionaryEntry) []DictionaryEntry {
	if len(entries) <= 1 {
		return entries
	}

	type key struct {
		word   string
		gender string
	}

	order := make([]key, 0, len(entries))
	mergedMap := make(map[key]*DictionaryEntry, len(entries))

	for _, entry := range entries {
		wClean := strings.ToLower(strings.TrimSpace(entry.Word))
		gClean := strings.ToLower(strings.TrimSpace(entry.Gender))
		k := key{word: wClean, gender: gClean}

		existing, found := mergedMap[k]
		if !found {
			entryCopy := entry
			mergedMap[k] = &entryCopy
			order = append(order, k)
		} else {
			existing.Translation = mergeSemicolonJoined(existing.Translation, entry.Translation)

			if existing.WordClass == "" {
				existing.WordClass = entry.WordClass
			} else if entry.WordClass != "" && !strings.Contains(strings.ToLower(existing.WordClass), strings.ToLower(entry.WordClass)) {
				existing.WordClass = existing.WordClass + ", " + entry.WordClass
			}

			existing.Forms = mergeSemicolonJoined(existing.Forms, entry.Forms)
			existing.Tags = appendUniqueStrings(existing.Tags, entry.Tags)
			existing.Examples = appendUniqueStrings(existing.Examples, entry.Examples)
		}
	}

	result := make([]DictionaryEntry, 0, len(order))
	for _, k := range order {
		result = append(result, *mergedMap[k])
	}
	return result
}

func mergeSemicolonJoined(s1, s2 string) string {
	parts1 := splitSemicolon(s1)
	parts2 := splitSemicolon(s2)
	seen := make(map[string]bool)
	var res []string

	for _, p := range append(parts1, parts2...) {
		pTrim := strings.TrimSpace(p)
		if pTrim == "" {
			continue
		}
		pLower := strings.ToLower(pTrim)
		if !seen[pLower] {
			seen[pLower] = true
			res = append(res, pTrim)
		}
	}
	return strings.Join(res, "; ")
}

func splitSemicolon(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return strings.Split(s, ";")
}

func appendUniqueStrings(base []string, add []string) []string {
	seen := make(map[string]bool, len(base)+len(add))
	var res []string
	for _, item := range base {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			itemLower := strings.ToLower(trimmed)
			if !seen[itemLower] {
				seen[itemLower] = true
				res = append(res, trimmed)
			}
		}
	}
	for _, item := range add {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			itemLower := strings.ToLower(trimmed)
			if !seen[itemLower] {
				seen[itemLower] = true
				res = append(res, trimmed)
			}
		}
	}
	return res
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
