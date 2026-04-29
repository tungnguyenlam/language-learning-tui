package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db}
	if err := store.Migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func OpenMemory() (*Store, error) {
	return Open(":memory:")
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		id INTEGER PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := s.migrationApplied(ctx, migration.ID)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %03d: %w", migration.ID, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (id) VALUES (?)`, migration.ID); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) migrationApplied(ctx context.Context, id int) (bool, error) {
	var exists int
	err := s.db.QueryRowContext(ctx, `SELECT 1 FROM schema_migrations WHERE id = ?`, id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) AppliedMigrations(ctx context.Context) ([]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM schema_migrations ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *Store) UpsertDeck(ctx context.Context, deck core.Deck) error {
	if deck.ID == "" {
		return errors.New("deck id is required")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO decks (id, name, description)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET name = excluded.name, description = excluded.description
	`, deck.ID, deck.Name, deck.Description)
	if err != nil {
		return err
	}
	for _, note := range deck.Notes {
		if note.DeckID == "" {
			note.DeckID = deck.ID
		}
		if err := s.UpsertNote(ctx, note); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetDeck(ctx context.Context, id string) (core.Deck, error) {
	var deck core.Deck
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description FROM decks WHERE id = ?`, id).Scan(&deck.ID, &deck.Name, &deck.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return core.Deck{}, fmt.Errorf("deck not found: %s", id)
	}
	if err != nil {
		return core.Deck{}, err
	}
	notes, err := s.notesForDeck(ctx, id)
	if err != nil {
		return core.Deck{}, err
	}
	deck.Notes = notes
	return deck, nil
}

func (s *Store) Decks(ctx context.Context) ([]core.Deck, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrowStart := todayStart.AddDate(0, 0, 1)
	rows, err := s.db.QueryContext(ctx, `
		SELECT d.id, d.name, d.description,
		       COUNT(c.id) as total_cards,
		       SUM(CASE WHEN (rs.due_at IS NULL OR rs.due_at <= ?) AND COALESCE(cf.suspended, 0) = 0 THEN 1 ELSE 0 END) as due_cards,
		       (
		           SELECT COUNT(*)
		           FROM reviews r
		           INNER JOIN cards rc ON rc.id = r.card_id
		           WHERE rc.deck_id = d.id AND r.reviewed_at >= ? AND r.reviewed_at < ?
		       ) as reviews_today,
		       (
		           SELECT COUNT(*)
		           FROM reviews r
		           INNER JOIN cards rc ON rc.id = r.card_id
		           WHERE rc.deck_id = d.id
		       ) as review_count,
		       (
		           SELECT COUNT(*)
		           FROM reviews r
		           INNER JOIN cards rc ON rc.id = r.card_id
		           WHERE rc.deck_id = d.id AND r.grade != ?
		       ) as successful_reviews
		FROM decks d
		LEFT JOIN cards c ON c.deck_id = d.id
		LEFT JOIN review_states rs ON rs.card_id = c.id
		LEFT JOIN card_flags cf ON cf.card_id = c.id
		GROUP BY d.id
		ORDER BY d.name
	`, now, todayStart, tomorrowStart, string(core.GradeAgain))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []core.Deck
	for rows.Next() {
		var deck core.Deck
		var total, due, reviewsToday, reviewCount, successfulReviews sql.NullInt64
		if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description, &total, &due, &reviewsToday, &reviewCount, &successfulReviews); err != nil {
			return nil, err
		}
		deck.TotalCards = int(total.Int64)
		deck.DueCards = int(due.Int64)
		deck.ReviewsToday = int(reviewsToday.Int64)
		if reviewCount.Int64 > 0 {
			deck.SuccessRate = float64(successfulReviews.Int64) / float64(reviewCount.Int64)
		}
		decks = append(decks, deck)
	}
	return decks, rows.Err()
}

func (s *Store) UpsertNote(ctx context.Context, note core.Note) error {
	if note.ID == "" {
		return errors.New("note id is required")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO notes (id, deck_id, front, back, extra)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			deck_id = excluded.deck_id,
			front = excluded.front,
			back = excluded.back,
			extra = excluded.extra
	`, note.ID, note.DeckID, note.Front, note.Back, note.Extra)
	if err != nil {
		return err
	}
	for _, card := range note.Cards {
		if card.NoteID == "" {
			card.NoteID = note.ID
		}
		if card.DeckID == "" {
			card.DeckID = note.DeckID
		}
		if err := s.UpsertCard(ctx, card); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) notesForDeck(ctx context.Context, deckID string) ([]core.Note, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, deck_id, front, back, extra
		FROM notes
		WHERE deck_id = ?
		ORDER BY id
	`, deckID)
	if err != nil {
		return nil, err
	}

	var notes []core.Note
	for rows.Next() {
		var note core.Note
		if err := rows.Scan(&note.ID, &note.DeckID, &note.Front, &note.Back, &note.Extra); err != nil {
			_ = rows.Close()
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range notes {
		cards, err := s.cardsForNote(ctx, notes[i].ID)
		if err != nil {
			return nil, err
		}
		notes[i].Cards = cards
	}
	return notes, nil
}

func (s *Store) cardsForNote(ctx context.Context, noteID string) ([]core.Card, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.note_id, c.deck_id, c.kind, c.prompt, c.answer, c.choices, COALESCE(cf.bookmarked, 0), COALESCE(cf.leech, 0), COALESCE(cf.suspended, 0)
		FROM cards c
		LEFT JOIN card_flags cf ON cf.card_id = c.id
		WHERE c.note_id = ?
		ORDER BY c.id
	`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []core.Card
	for rows.Next() {
		var card core.Card
		var kind string
		var choicesStr string
		var bookmarked, leech, suspended int
		if err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &kind, &card.Prompt, &card.Answer, &choicesStr, &bookmarked, &leech, &suspended); err != nil {
			return nil, err
		}
		card.Kind = core.CardKind(kind)
		card.Bookmarked = bookmarked != 0
		card.Leech = leech != 0
		card.Suspended = suspended != 0
		card.Choices = parseChoices(choicesStr)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (s *Store) UpsertCard(ctx context.Context, card core.Card) error {
	if err := core.ValidateCard(card); err != nil {
		return err
	}
	choices := strings.Join(card.Choices, "|||")
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cards (id, note_id, deck_id, kind, prompt, answer, choices)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			note_id = excluded.note_id,
			deck_id = excluded.deck_id,
			kind = excluded.kind,
			prompt = excluded.prompt,
			answer = excluded.answer,
			choices = excluded.choices
	`, card.ID, card.NoteID, card.DeckID, string(card.Kind), card.Prompt, card.Answer, choices)
	return err
}

func (s *Store) RecordReview(ctx context.Context, result core.ReviewResult) error {
	if result.CardID == "" {
		return errors.New("card id is required")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO reviews (card_id, grade, reviewed_at, due_at, interval_seconds, stability, difficulty, ease, reviews, lapses)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, result.CardID, string(result.Grade), result.Reviewed.UTC(), result.Next.Due.UTC(), int64(result.Next.Interval.Seconds()), result.Next.Stability, result.Next.Difficulty, result.Next.Ease, result.Next.Reviews, result.Next.Lapses)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO review_states (card_id, due_at, interval_seconds, stability, difficulty, ease, reviews, lapses)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET
			due_at = excluded.due_at,
			interval_seconds = excluded.interval_seconds,
			stability = excluded.stability,
			difficulty = excluded.difficulty,
			ease = excluded.ease,
			reviews = excluded.reviews,
			lapses = excluded.lapses
	`, result.CardID, result.Next.Due.UTC(), int64(result.Next.Interval.Seconds()), result.Next.Stability, result.Next.Difficulty, result.Next.Ease, result.Next.Reviews, result.Next.Lapses)
	if err != nil {
		return err
	}

	if result.Grade == core.GradeAgain {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO card_flags (card_id, bookmarked, lapse_streak, leech, updated_at)
			VALUES (?, 0, 1, 0, ?)
			ON CONFLICT(card_id) DO UPDATE SET
				lapse_streak = card_flags.lapse_streak + 1,
				leech = CASE WHEN card_flags.lapse_streak + 1 >= 3 THEN 1 ELSE 0 END,
				updated_at = ?
		`, result.CardID, time.Now().UTC(), time.Now().UTC())
	} else {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO card_flags (card_id, bookmarked, lapse_streak, leech, updated_at)
			VALUES (?, 0, 0, 0, ?)
			ON CONFLICT(card_id) DO UPDATE SET
				lapse_streak = 0,
				leech = 0,
				updated_at = ?
		`, result.CardID, time.Now().UTC(), time.Now().UTC())
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) UndoLastReview(ctx context.Context, cardID string) error {
	if cardID == "" {
		return errors.New("card id is required")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var reviewID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id
		FROM reviews
		WHERE card_id = ?
		ORDER BY reviewed_at DESC, id DESC
		LIMIT 1
	`, cardID).Scan(&reviewID)
	if errors.Is(err, sql.ErrNoRows) {
		return core.ErrNoReviewToUndo
	}
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM reviews WHERE id = ?`, reviewID); err != nil {
		return err
	}

	var state core.ReviewState
	var intervalSec int64
	err = tx.QueryRowContext(ctx, `
		SELECT card_id, due_at, interval_seconds, stability, difficulty, ease, reviews, lapses
		FROM reviews
		WHERE card_id = ?
		ORDER BY reviewed_at DESC, id DESC
		LIMIT 1
	`, cardID).Scan(&state.CardID, &state.Due, &intervalSec, &state.Stability, &state.Difficulty, &state.Ease, &state.Reviews, &state.Lapses)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := tx.ExecContext(ctx, `DELETE FROM review_states WHERE card_id = ?`, cardID); err != nil {
			return err
		}
		return tx.Commit()
	}
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO review_states (card_id, due_at, interval_seconds, stability, difficulty, ease, reviews, lapses)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET
			due_at = excluded.due_at,
			interval_seconds = excluded.interval_seconds,
			stability = excluded.stability,
			difficulty = excluded.difficulty,
			ease = excluded.ease,
			reviews = excluded.reviews,
			lapses = excluded.lapses
	`, state.CardID, state.Due.UTC(), intervalSec, state.Stability, state.Difficulty, state.Ease, state.Reviews, state.Lapses); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) DueCards(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.note_id, c.deck_id, c.kind, c.prompt, c.answer, c.choices, COALESCE(cf.bookmarked, 0), COALESCE(cf.leech, 0), COALESCE(cf.suspended, 0)
		FROM cards c
		LEFT JOIN review_states rs ON rs.card_id = c.id
		LEFT JOIN card_flags cf ON cf.card_id = c.id
		WHERE (rs.due_at IS NULL OR rs.due_at <= ?)
		  AND COALESCE(cf.suspended, 0) = 0
		ORDER BY COALESCE(rs.due_at, '1970-01-01T00:00:00Z'), c.id
		LIMIT ?
	`, now.UTC(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []core.Card
	for rows.Next() {
		var card core.Card
		var kind string
		var choicesStr string
		var bookmarked, leech, suspended int
		if err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &kind, &card.Prompt, &card.Answer, &choicesStr, &bookmarked, &leech, &suspended); err != nil {
			return nil, err
		}
		card.Kind = core.CardKind(kind)
		card.Bookmarked = bookmarked != 0
		card.Leech = leech != 0
		card.Suspended = suspended != 0
		card.Choices = parseChoices(choicesStr)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (s *Store) DueCardsBookmarked(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.note_id, c.deck_id, c.kind, c.prompt, c.answer, c.choices, COALESCE(cf.bookmarked, 0), COALESCE(cf.leech, 0), COALESCE(cf.suspended, 0)
		FROM cards c
		INNER JOIN card_flags cf ON cf.card_id = c.id AND cf.bookmarked = 1
		LEFT JOIN review_states rs ON rs.card_id = c.id
		WHERE (rs.due_at IS NULL OR rs.due_at <= ?)
		  AND COALESCE(cf.suspended, 0) = 0
		ORDER BY COALESCE(rs.due_at, '1970-01-01T00:00:00Z'), c.id
		LIMIT ?
	`, now.UTC(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []core.Card
	for rows.Next() {
		var card core.Card
		var kind string
		var choicesStr string
		var bookmarked, leech, suspended int
		if err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &kind, &card.Prompt, &card.Answer, &choicesStr, &bookmarked, &leech, &suspended); err != nil {
			return nil, err
		}
		card.Kind = core.CardKind(kind)
		card.Bookmarked = bookmarked != 0
		card.Leech = leech != 0
		card.Suspended = suspended != 0
		card.Choices = parseChoices(choicesStr)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (s *Store) SetCardBookmark(ctx context.Context, cardID string, bookmarked bool) error {
	if cardID == "" {
		return errors.New("card id is required")
	}
	value := 0
	if bookmarked {
		value = 1
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO card_flags (card_id, bookmarked, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET
			bookmarked = excluded.bookmarked,
			updated_at = excluded.updated_at
	`, cardID, value, time.Now().UTC())
	return err
}

func (s *Store) SetCardSuspended(ctx context.Context, cardID string, suspended bool) error {
	if cardID == "" {
		return errors.New("card id is required")
	}
	value := 0
	if suspended {
		value = 1
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO card_flags (card_id, suspended, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET
			suspended = excluded.suspended,
			updated_at = excluded.updated_at
	`, cardID, value, time.Now().UTC())
	return err
}

func (s *Store) SetDailyGoal(ctx context.Context, goal int) error {
	if goal < 1 {
		goal = 1
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ('daily_goal', ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, strconv.Itoa(goal), time.Now().UTC())
	return err
}

func (s *Store) dailyGoal(ctx context.Context) (int, error) {
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM app_settings WHERE key = 'daily_goal'`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return 10, nil
	}
	if err != nil {
		return 0, err
	}
	goal, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || goal < 1 {
		return 10, nil
	}
	return goal, nil
}

func (s *Store) Statistics(ctx context.Context) (core.Statistics, error) {
	var stats core.Statistics
	stats.Grades = make(map[core.ReviewGrade]int)

	// Total cards
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cards`).Scan(&stats.TotalCards)
	if err != nil {
		return stats, err
	}

	// Card maturity: New (no reviews), Young (interval < 21 days), Mature (interval >= 21 days)
	// 21 days is a common threshold in Anki.
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN rs.reviews = 0 OR rs.reviews IS NULL THEN 1 ELSE 0 END), 0) as new,
			COALESCE(SUM(CASE WHEN rs.reviews > 0 AND rs.interval_seconds < 1814400 THEN 1 ELSE 0 END), 0) as young,
			COALESCE(SUM(CASE WHEN rs.reviews > 0 AND rs.interval_seconds >= 1814400 THEN 1 ELSE 0 END), 0) as mature
		FROM cards c
		LEFT JOIN review_states rs ON rs.card_id = c.id
	`).Scan(&stats.NewCards, &stats.YoungCards, &stats.MatureCards)
	if err != nil {
		return stats, err
	}

	// Total reviews
	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews`).Scan(&stats.TotalReviews)
	if err != nil {
		return stats, err
	}
	stats.DailyGoal, err = s.dailyGoal(ctx)
	if err != nil {
		return stats, err
	}

	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM reviews
		WHERE reviewed_at >= ? AND reviewed_at < ?
	`, todayStart, todayStart.AddDate(0, 0, 1)).Scan(&stats.ReviewsToday)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM card_flags WHERE bookmarked = 1`).Scan(&stats.BookmarkedCards)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM cards c
		INNER JOIN card_flags cf ON cf.card_id = c.id AND cf.bookmarked = 1
		LEFT JOIN review_states rs ON rs.card_id = c.id
		WHERE rs.due_at IS NULL OR rs.due_at <= ?
	`, now).Scan(&stats.BookmarkedDue)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM card_flags WHERE leech = 1`).Scan(&stats.LeechCards)
	if err != nil {
		return stats, err
	}

	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM card_flags WHERE suspended = 1`).Scan(&stats.SuspendedCards)
	if err != nil {
		return stats, err
	}

	streak, err := s.currentStreak(ctx, now)
	if err != nil {
		return stats, err
	}
	stats.CurrentStreak = streak

	// Success rate (Grade != Again)
	if stats.TotalReviews > 0 {
		var successCount int
		err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews WHERE grade != ?`, string(core.GradeAgain)).Scan(&successCount)
		if err != nil {
			return stats, err
		}
		stats.SuccessRate = float64(successCount) / float64(stats.TotalReviews)
	}

	// Reviews by grade
	rows, err := s.db.QueryContext(ctx, `SELECT grade, COUNT(*) FROM reviews GROUP BY grade`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var grade string
		var count int
		if err := rows.Scan(&grade, &count); err != nil {
			return stats, err
		}
		stats.Grades[core.ReviewGrade(grade)] = count
	}

	return stats, nil
}

func (s *Store) currentStreak(ctx context.Context, now time.Time) (int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT reviewed_at
		FROM reviews
		ORDER BY reviewed_at DESC
	`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	dates := make(map[string]bool)
	for rows.Next() {
		var reviewed time.Time
		if err := rows.Scan(&reviewed); err != nil {
			return 0, err
		}
		day := reviewed.UTC()
		dates[time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")] = true
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if !dates[day.Format("2006-01-02")] {
		yesterday := day.AddDate(0, 0, -1)
		if !dates[yesterday.Format("2006-01-02")] {
			return 0, nil
		}
		day = yesterday
	}
	streak := 0
	for dates[day.Format("2006-01-02")] {
		streak++
		day = day.AddDate(0, 0, -1)
	}
	return streak, nil
}

func (s *Store) GetReviewState(ctx context.Context, cardID string) (core.ReviewState, error) {
	var state core.ReviewState
	var intervalSec int64
	err := s.db.QueryRowContext(ctx, `
		SELECT card_id, due_at, interval_seconds, stability, difficulty, ease, reviews, lapses
		FROM review_states
		WHERE card_id = ?
	`, cardID).Scan(&state.CardID, &state.Due, &intervalSec, &state.Stability, &state.Difficulty, &state.Ease, &state.Reviews, &state.Lapses)
	if errors.Is(err, sql.ErrNoRows) {
		return core.NewReviewState(cardID, time.Now()), nil
	}
	if err != nil {
		return core.ReviewState{}, err
	}
	state.Interval = time.Duration(intervalSec) * time.Second
	return state, nil
}

func parseChoices(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	return strings.Split(raw, "|||")
}

func (s *Store) Cards(ctx context.Context, deckID string, search string) ([]core.Card, error) {
	query := `
		SELECT c.id, c.note_id, c.deck_id, c.kind, c.prompt, c.answer, c.choices, COALESCE(cf.bookmarked, 0), COALESCE(cf.leech, 0), COALESCE(cf.suspended, 0)
		FROM cards c
		LEFT JOIN card_flags cf ON cf.card_id = c.id
		WHERE 1=1
	`
	args := []interface{}{}
	if deckID != "" {
		query += ` AND c.deck_id = ?`
		args = append(args, deckID)
	}
	if search != "" {
		query += ` AND (c.prompt LIKE ? OR c.answer LIKE ?)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}
	query += ` ORDER BY c.id LIMIT 100`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []core.Card
	for rows.Next() {
		var card core.Card
		var kind string
		var choicesStr string
		var bookmarked, leech, suspended int
		if err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &kind, &card.Prompt, &card.Answer, &choicesStr, &bookmarked, &leech, &suspended); err != nil {
			return nil, err
		}
		card.Kind = core.CardKind(kind)
		card.Bookmarked = bookmarked != 0
		card.Leech = leech != 0
		card.Suspended = suspended != 0
		card.Choices = parseChoices(choicesStr)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}
