package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestWithTxCommitsOnlySuccessfulWork(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer store.Close()

	insertSetting := func(tx *sql.Tx, key string) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO app_settings (key, value) VALUES (?, ?)`, key, "saved")
		return err
	}
	if err := store.withTx(ctx, func(tx *sql.Tx) error {
		return insertSetting(tx, "committed")
	}); err != nil {
		t.Fatalf("commit transaction: %v", err)
	}

	wantRollback := errors.New("force rollback")
	err = store.withTx(ctx, func(tx *sql.Tx) error {
		if err := insertSetting(tx, "rolled-back"); err != nil {
			return err
		}
		return wantRollback
	})
	if !errors.Is(err, wantRollback) {
		t.Fatalf("rollback error = %v, want %v", err, wantRollback)
	}

	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM app_settings WHERE key = 'committed'`).Scan(&count); err != nil {
		t.Fatalf("query committed row: %v", err)
	}
	if count != 1 {
		t.Fatalf("committed rows = %d, want 1", count)
	}
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM app_settings WHERE key = 'rolled-back'`).Scan(&count); err != nil {
		t.Fatalf("query rolled-back row: %v", err)
	}
	if count != 0 {
		t.Fatalf("rolled-back rows = %d, want 0", count)
	}
}
