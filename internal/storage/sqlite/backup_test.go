package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"deutsch-tui/internal/core"
)

func seedBackupStore(t *testing.T, ctx context.Context, store *Store) {
	t.Helper()
	deck := core.Deck{ID: "deck-backup", Name: "Backup Deck", Description: "for tests"}
	note := core.Note{
		ID:     "note-backup",
		DeckID: deck.ID,
		Front:  "die Sicherung",
		Back:   "the backup",
		Extra:  "Sicherungskopie",
		Tags:   []string{"b1", "tech"},
	}
	card := core.Card{
		ID:     "card-backup",
		NoteID: note.ID,
		DeckID: deck.ID,
		Kind:   core.CardKindFlashcard,
		Prompt: "die Sicherung",
		Answer: "the backup",
	}
	if err := store.UpsertDeck(ctx, deck); err != nil {
		t.Fatalf("upsert deck: %v", err)
	}
	if err := store.UpsertNote(ctx, note); err != nil {
		t.Fatalf("upsert note: %v", err)
	}
	if err := store.UpsertCard(ctx, card); err != nil {
		t.Fatalf("upsert card: %v", err)
	}
}

func openTempStore(t *testing.T, name string) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func TestBackupAndRestoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	dest := filepath.Join(t.TempDir(), "backups", BackupFileName(time.Now()))
	info, err := store.Backup(ctx, dest)
	if err != nil {
		t.Fatalf("backup: %v", err)
	}
	if info.Rows["decks"] != 1 || info.Rows["notes"] != 1 || info.Rows["cards"] != 1 {
		t.Fatalf("unexpected backup row counts: %+v", info.Rows)
	}
	if info.SchemaVersion <= 0 {
		t.Fatalf("backup schema version = %d, want > 0", info.SchemaVersion)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("backup file missing: %v", err)
	}

	// Wipe the live data the way the Settings "reset" action would.
	if err := store.Reset(ctx); err != nil {
		t.Fatalf("reset: %v", err)
	}
	decks, err := store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks after reset: %v", err)
	}
	if len(decks) != 0 {
		t.Fatalf("expected empty store after reset, got %d decks", len(decks))
	}

	if _, err := store.Restore(ctx, dest); err != nil {
		t.Fatalf("restore: %v", err)
	}
	decks, err = store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks after restore: %v", err)
	}
	if len(decks) != 1 || decks[0].ID != "deck-backup" || decks[0].Name != "Backup Deck" {
		t.Fatalf("deck not restored: %+v", decks)
	}
	cards, err := store.Cards(ctx, "deck-backup", "", "")
	if err != nil {
		t.Fatalf("cards after restore: %v", err)
	}
	if len(cards) != 1 || cards[0].Prompt != "die Sicherung" {
		t.Fatalf("card not restored: %+v", cards)
	}
}

func TestRestoreIntoDifferentDatabase(t *testing.T) {
	ctx := context.Background()
	source := openTempStore(t, "source.db")
	seedBackupStore(t, ctx, source)

	dest := filepath.Join(t.TempDir(), BackupFileName(time.Now()))
	if _, err := source.Backup(ctx, dest); err != nil {
		t.Fatalf("backup: %v", err)
	}

	// A fresh machine: new data dir, empty database, restore the file.
	target := openTempStore(t, "target.db")
	info, err := target.Restore(ctx, dest)
	if err != nil {
		t.Fatalf("restore into fresh store: %v", err)
	}
	if info.TotalRows == 0 {
		t.Fatal("restore reported zero rows")
	}
	decks, err := target.Decks(ctx)
	if err != nil {
		t.Fatalf("decks: %v", err)
	}
	if len(decks) != 1 || decks[0].ID != "deck-backup" {
		t.Fatalf("deck not restored into fresh store: %+v", decks)
	}
}

func TestRestorePreservesReviewHistory(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	if err := store.RecordReview(ctx, core.ReviewResult{
		CardID:   "card-backup",
		Grade:    core.GradeGood,
		Reviewed: time.Now().UTC(),
		Next: core.ReviewState{
			CardID:   "card-backup",
			Due:      time.Now().Add(48 * time.Hour).UTC(),
			Interval: 48 * time.Hour,
			Ease:     2.6,
			Reviews:  3,
			Lapses:   1,
		},
	}); err != nil {
		t.Fatalf("record review: %v", err)
	}

	dest := filepath.Join(t.TempDir(), BackupFileName(time.Now()))
	info, err := store.Backup(ctx, dest)
	if err != nil {
		t.Fatalf("backup: %v", err)
	}
	if info.Rows["reviews"] == 0 || info.Rows["review_states"] == 0 {
		t.Fatalf("review history missing from backup: %+v", info.Rows)
	}

	if err := store.Reset(ctx); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := store.Restore(ctx, dest); err != nil {
		t.Fatalf("restore: %v", err)
	}

	restored, err := store.GetReviewState(ctx, "card-backup")
	if err != nil {
		t.Fatalf("get review state: %v", err)
	}
	if restored.Reviews != 3 || restored.Lapses != 1 {
		t.Fatalf("review state not restored: %+v", restored)
	}
}

func TestRestoreRejectsNonBackupFile(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")

	// A live database is valid SQLite but has no backup_meta table. Users
	// will point at learning.db by mistake; that must not wipe their data.
	otherPath := filepath.Join(t.TempDir(), "not-a-backup.db")
	other, err := Open(otherPath)
	if err != nil {
		t.Fatalf("open second store: %v", err)
	}
	seedBackupStore(t, ctx, other)
	if err := other.Close(); err != nil {
		t.Fatalf("close second store: %v", err)
	}

	if _, err := store.Restore(ctx, otherPath); !errors.Is(err, ErrNotBackup) {
		t.Fatalf("Restore(non-backup) error = %v, want ErrNotBackup", err)
	}
}

func TestRestoreRejectsNewerSchema(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	dest := filepath.Join(t.TempDir(), BackupFileName(time.Now()))
	if _, err := store.Backup(ctx, dest); err != nil {
		t.Fatalf("backup: %v", err)
	}

	// Forge a backup that claims to come from a much newer build. Open the
	// file directly rather than through Open(), which would try to migrate
	// the backup as if it were a live database.
	forged, err := sql.Open("sqlite", dest)
	if err != nil {
		t.Fatalf("open backup: %v", err)
	}
	if _, err := forged.ExecContext(ctx, `UPDATE backup_meta SET value = '9999' WHERE key = 'schema_version'`); err != nil {
		t.Fatalf("forge schema version: %v", err)
	}
	if err := forged.Close(); err != nil {
		t.Fatalf("close forged backup: %v", err)
	}

	if _, err := store.Restore(ctx, dest); err == nil {
		t.Fatal("Restore(newer schema) succeeded, want error")
	}
}

func TestRestoreIsAtomicOnFailure(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	missing := filepath.Join(t.TempDir(), "does-not-exist.db")
	if _, err := store.Restore(ctx, missing); err == nil {
		t.Fatal("Restore(missing file) succeeded, want error")
	}

	decks, err := store.Decks(ctx)
	if err != nil {
		t.Fatalf("decks: %v", err)
	}
	if len(decks) != 1 {
		t.Fatalf("failed restore damaged live data: %d decks remain", len(decks))
	}
}

func TestInspectBackupReportsCounts(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	dest := filepath.Join(t.TempDir(), BackupFileName(time.Now()))
	if _, err := store.Backup(ctx, dest); err != nil {
		t.Fatalf("backup: %v", err)
	}

	info, err := store.InspectBackup(ctx, dest)
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	if info.Rows["notes"] != 1 {
		t.Fatalf("inspect note count = %d, want 1", info.Rows["notes"])
	}
	if info.CreatedAt.IsZero() {
		t.Fatal("inspect did not report a creation time")
	}
}

func TestBackupOverwritesExistingFile(t *testing.T) {
	ctx := context.Background()
	store := openTempStore(t, "learning.db")
	seedBackupStore(t, ctx, store)

	dest := filepath.Join(t.TempDir(), "backup-fixed-name.db")
	if _, err := store.Backup(ctx, dest); err != nil {
		t.Fatalf("first backup: %v", err)
	}
	if _, err := store.Backup(ctx, dest); err != nil {
		t.Fatalf("second backup to same path: %v", err)
	}
}

func TestLatestBackupPicksNewest(t *testing.T) {
	dir := t.TempDir()
	names := []string{
		"backup-20260101-000000.db",
		"backup-20260805-120000.db",
		"backup-20260304-093000.db",
		"notes.txt",
	}
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o600); err != nil {
			t.Fatalf("write %s: %v", n, err)
		}
	}
	got, err := LatestBackup(dir)
	if err != nil {
		t.Fatalf("latest backup: %v", err)
	}
	if filepath.Base(got) != "backup-20260805-120000.db" {
		t.Fatalf("LatestBackup = %s, want the 2026-08-05 file", filepath.Base(got))
	}

	if _, err := LatestBackup(t.TempDir()); err == nil {
		t.Fatal("LatestBackup(empty dir) succeeded, want error")
	}
}

func TestBackupFileNameIsChronologicallySortable(t *testing.T) {
	early := BackupFileName(time.Date(2026, 3, 4, 9, 30, 0, 0, time.UTC))
	late := BackupFileName(time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC))
	if !(early < late) {
		t.Fatalf("names do not sort chronologically: %s vs %s", early, late)
	}
}
