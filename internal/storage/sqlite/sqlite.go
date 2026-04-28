package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (s *Store) UpsertCard(ctx context.Context, card core.Card) error {
	if err := core.ValidateCard(card); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cards (id, note_id, deck_id, kind, prompt, answer)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			note_id = excluded.note_id,
			deck_id = excluded.deck_id,
			kind = excluded.kind,
			prompt = excluded.prompt,
			answer = excluded.answer
	`, card.ID, card.NoteID, card.DeckID, string(card.Kind), card.Prompt, card.Answer)
	return err
}

func (s *Store) RecordReview(ctx context.Context, result core.ReviewResult) error {
	if result.CardID == "" {
		return errors.New("card id is required")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO reviews (card_id, grade, reviewed_at, due_at, interval_seconds, ease, reviews, lapses)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, result.CardID, string(result.Grade), result.Reviewed.UTC(), result.Next.Due.UTC(), int64(result.Next.Interval.Seconds()), result.Next.Ease, result.Next.Reviews, result.Next.Lapses)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO review_states (card_id, due_at, interval_seconds, ease, reviews, lapses)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(card_id) DO UPDATE SET
			due_at = excluded.due_at,
			interval_seconds = excluded.interval_seconds,
			ease = excluded.ease,
			reviews = excluded.reviews,
			lapses = excluded.lapses
	`, result.CardID, result.Next.Due.UTC(), int64(result.Next.Interval.Seconds()), result.Next.Ease, result.Next.Reviews, result.Next.Lapses)
	return err
}

func (s *Store) DueCards(ctx context.Context, now time.Time, limit int) ([]core.Card, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.note_id, c.deck_id, c.kind, c.prompt, c.answer
		FROM cards c
		LEFT JOIN review_states rs ON rs.card_id = c.id
		WHERE rs.due_at IS NULL OR rs.due_at <= ?
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
		if err := rows.Scan(&card.ID, &card.NoteID, &card.DeckID, &kind, &card.Prompt, &card.Answer); err != nil {
			return nil, err
		}
		card.Kind = core.CardKind(kind)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}
