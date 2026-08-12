package sqlite

import (
	"context"
	"database/sql"
)

// withTx runs fn in a transaction, committing only when fn succeeds. The
// deferred rollback is intentionally kept even on success: database/sql makes
// rollback after commit a harmless no-op, while every early return remains
// safe by construction.
func (s *Store) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
