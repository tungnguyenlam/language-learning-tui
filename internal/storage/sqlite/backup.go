package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"deutsch-tui/internal/core"
)

var _ core.BackupRepository = (*Store)(nil)

// backupFormatVersion identifies the layout of the backup file itself, not
// the schema of the data inside it. Bump it only if the meta table or the
// table-copy strategy changes in a way older readers cannot handle.
const backupFormatVersion = 1

// backupTables lists the tables a backup carries, ordered parent-first so a
// restore can insert them in this order and delete them in reverse without
// tripping foreign keys.
//
// The dictionary tables are deliberately excluded. They hold ~834k rows
// derived from the user's dict.cc download, they are re-importable at any
// time, and including them would turn a small progress backup into a
// multi-hundred-megabyte file. schema_migrations is excluded too: a restore
// targets an already-migrated database, and copying migration bookkeeping
// from an older backup would make the live database lie about its schema.
var backupTables = []string{
	"decks",
	"notes",
	"cards",
	"review_states",
	"reviews",
	"card_flags",
	"app_settings",
}

// ErrNotBackup is returned when a file is readable as SQLite but is not one
// of our backups. The most likely cause is a user pointing the restore at
// their live learning.db, which we must not silently accept.
var ErrNotBackup = errors.New("file is not a deutsch-tui backup")

// Backup writes a portable copy of the user's decks, cards, and review
// history to destPath. It excludes the dictionary, so the result is small
// enough to keep in version control or a cloud folder.
func (s *Store) Backup(ctx context.Context, destPath string) (core.BackupInfo, error) {
	info := core.BackupInfo{Path: destPath, Rows: map[string]int{}}
	if strings.TrimSpace(destPath) == "" {
		return info, errors.New("backup: destination path is required")
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return info, fmt.Errorf("backup: create directory: %w", err)
	}
	// ATTACH will happily open an existing file and then CREATE TABLE will
	// fail on the second run, so clear any previous file (and its WAL
	// sidecars) first.
	for _, p := range []string{destPath, destPath + "-wal", destPath + "-shm"} {
		if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
			return info, fmt.Errorf("backup: replace existing file: %w", err)
		}
	}

	schemaVersion, err := s.currentSchemaVersion(ctx)
	if err != nil {
		return info, err
	}

	// ATTACH is connection-scoped, so every statement below has to run on
	// the same connection rather than on the pool.
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return info, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `ATTACH DATABASE ? AS bkp`, destPath); err != nil {
		return info, fmt.Errorf("backup: attach destination: %w", err)
	}
	defer func() { _, _ = conn.ExecContext(ctx, `DETACH DATABASE bkp`) }()

	createdAt := time.Now().UTC()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE bkp.backup_meta (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`); err != nil {
		return info, fmt.Errorf("backup: create meta table: %w", err)
	}
	meta := [][2]string{
		{"format", fmt.Sprintf("%d", backupFormatVersion)},
		{"schema_version", fmt.Sprintf("%d", schemaVersion)},
		{"created_at", createdAt.Format(time.RFC3339)},
		{"application", "deutsch-tui"},
	}
	for _, kv := range meta {
		if _, err := conn.ExecContext(ctx, `INSERT INTO bkp.backup_meta (key, value) VALUES (?, ?)`, kv[0], kv[1]); err != nil {
			return info, fmt.Errorf("backup: write meta: %w", err)
		}
	}

	for _, table := range backupTables {
		exists, err := connTableExists(ctx, conn, "main", table)
		if err != nil {
			return info, err
		}
		if !exists {
			continue
		}
		// CREATE TABLE ... AS SELECT copies the rows and column names but
		// not the constraints. That is intentional: the backup is a data
		// payload, and restore replays it into the live, migrated schema.
		stmt := fmt.Sprintf(`CREATE TABLE bkp.%q AS SELECT * FROM main.%q`, table, table)
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return info, fmt.Errorf("backup: copy %s: %w", table, err)
		}
		var count int
		if err := conn.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM bkp.%q`, table)).Scan(&count); err != nil {
			return info, fmt.Errorf("backup: count %s: %w", table, err)
		}
		info.Rows[table] = count
		info.TotalRows += count
	}

	info.CreatedAt = createdAt
	info.SchemaVersion = schemaVersion
	return info, nil
}

// InspectBackup reads a backup's metadata without modifying anything, so the
// UI can show the user what they are about to overwrite their data with.
func (s *Store) InspectBackup(ctx context.Context, srcPath string) (core.BackupInfo, error) {
	info := core.BackupInfo{Path: srcPath, Rows: map[string]int{}}
	if _, err := os.Stat(srcPath); err != nil {
		return info, fmt.Errorf("restore: %w", err)
	}
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return info, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `ATTACH DATABASE ? AS bkp`, srcPath); err != nil {
		return info, fmt.Errorf("restore: attach backup: %w", err)
	}
	defer func() { _, _ = conn.ExecContext(ctx, `DETACH DATABASE bkp`) }()

	return readBackupInfo(ctx, conn, srcPath)
}

// Restore replaces the live decks, cards, and review history with the
// contents of a backup file. It is all-or-nothing: the copy runs in a single
// transaction, so a failure part-way leaves the existing data untouched.
//
// The dictionary is not touched, because backups do not carry it.
func (s *Store) Restore(ctx context.Context, srcPath string) (core.BackupInfo, error) {
	info := core.BackupInfo{Path: srcPath, Rows: map[string]int{}}
	if _, err := os.Stat(srcPath); err != nil {
		return info, fmt.Errorf("restore: %w", err)
	}

	currentVersion, err := s.currentSchemaVersion(ctx)
	if err != nil {
		return info, err
	}

	conn, err := s.db.Conn(ctx)
	if err != nil {
		return info, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `ATTACH DATABASE ? AS bkp`, srcPath); err != nil {
		return info, fmt.Errorf("restore: attach backup: %w", err)
	}
	defer func() { _, _ = conn.ExecContext(ctx, `DETACH DATABASE bkp`) }()

	info, err = readBackupInfo(ctx, conn, srcPath)
	if err != nil {
		return info, err
	}
	// A backup from a newer build may contain columns or tables this binary
	// does not know how to place. Refuse rather than restore a subset and
	// leave the user believing their data came back intact.
	if info.SchemaVersion > currentVersion {
		return info, fmt.Errorf("restore: backup was made by a newer version (schema %d, this build supports %d); upgrade first",
			info.SchemaVersion, currentVersion)
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return info, err
	}
	defer tx.Rollback()

	// Children first so foreign keys stay satisfied at every step.
	for i := len(backupTables) - 1; i >= 0; i-- {
		table := backupTables[i]
		exists, err := connTableExists(ctx, conn, "main", table)
		if err != nil {
			return info, err
		}
		if !exists {
			continue
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf(`DELETE FROM main.%q`, table)); err != nil {
			return info, fmt.Errorf("restore: clear %s: %w", table, err)
		}
	}

	restored := map[string]int{}
	total := 0
	for _, table := range backupTables {
		inBackup, err := connTableExists(ctx, conn, "bkp", table)
		if err != nil {
			return info, err
		}
		inMain, err := connTableExists(ctx, conn, "main", table)
		if err != nil {
			return info, err
		}
		if !inBackup || !inMain {
			continue
		}
		// Only copy columns both sides share. An older backup is missing
		// columns added by later migrations (those keep their defaults),
		// and we ignore any column the live schema no longer has.
		cols, err := sharedColumns(ctx, conn, table)
		if err != nil {
			return info, err
		}
		if len(cols) == 0 {
			continue
		}
		quoted := make([]string, len(cols))
		for i, c := range cols {
			quoted[i] = fmt.Sprintf("%q", c)
		}
		list := strings.Join(quoted, ", ")
		stmt := fmt.Sprintf(`INSERT INTO main.%q (%s) SELECT %s FROM bkp.%q`, table, list, list, table)
		res, err := tx.ExecContext(ctx, stmt)
		if err != nil {
			return info, fmt.Errorf("restore: copy %s: %w", table, err)
		}
		if n, err := res.RowsAffected(); err == nil {
			restored[table] = int(n)
			total += int(n)
		}
	}

	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("restore: commit: %w", err)
	}

	info.Rows = restored
	info.TotalRows = total
	return info, nil
}

// readBackupInfo validates that the attached `bkp` database really is one of
// our backups and reads its metadata.
func readBackupInfo(ctx context.Context, conn *sql.Conn, path string) (core.BackupInfo, error) {
	info := core.BackupInfo{Path: path, Rows: map[string]int{}}

	hasMeta, err := connTableExists(ctx, conn, "bkp", "backup_meta")
	if err != nil {
		return info, err
	}
	if !hasMeta {
		return info, fmt.Errorf("restore: %s: %w", filepath.Base(path), ErrNotBackup)
	}

	meta := map[string]string{}
	rows, err := conn.QueryContext(ctx, `SELECT key, value FROM bkp.backup_meta`)
	if err != nil {
		return info, fmt.Errorf("restore: read meta: %w", err)
	}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			rows.Close()
			return info, err
		}
		meta[k] = v
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return info, err
	}
	rows.Close()

	if meta["application"] != "deutsch-tui" {
		return info, fmt.Errorf("restore: %s: %w", filepath.Base(path), ErrNotBackup)
	}
	var format int
	if _, err := fmt.Sscanf(meta["format"], "%d", &format); err != nil || format > backupFormatVersion {
		return info, fmt.Errorf("restore: unsupported backup format %q (this build reads up to %d)",
			meta["format"], backupFormatVersion)
	}
	if _, err := fmt.Sscanf(meta["schema_version"], "%d", &info.SchemaVersion); err != nil {
		return info, fmt.Errorf("restore: backup is missing its schema version")
	}
	if ts, err := time.Parse(time.RFC3339, meta["created_at"]); err == nil {
		info.CreatedAt = ts
	}

	for _, table := range backupTables {
		exists, err := connTableExists(ctx, conn, "bkp", table)
		if err != nil {
			return info, err
		}
		if !exists {
			continue
		}
		var count int
		if err := conn.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM bkp.%q`, table)).Scan(&count); err != nil {
			return info, fmt.Errorf("restore: count %s: %w", table, err)
		}
		info.Rows[table] = count
		info.TotalRows += count
	}
	return info, nil
}

// sharedColumns returns the columns present in both main.table and
// bkp.table, in the live schema's order.
func sharedColumns(ctx context.Context, conn *sql.Conn, table string) ([]string, error) {
	mainCols, err := connTableColumns(ctx, conn, "main", table)
	if err != nil {
		return nil, err
	}
	backupCols, err := connTableColumns(ctx, conn, "bkp", table)
	if err != nil {
		return nil, err
	}
	have := make(map[string]bool, len(backupCols))
	for _, c := range backupCols {
		have[c] = true
	}
	var shared []string
	for _, c := range mainCols {
		if have[c] {
			shared = append(shared, c)
		}
	}
	return shared, nil
}

func connTableColumns(ctx context.Context, conn *sql.Conn, schema, table string) ([]string, error) {
	rows, err := conn.QueryContext(ctx, fmt.Sprintf(`SELECT name FROM %q.pragma_table_info(?)`, schema), table)
	if err != nil {
		return nil, fmt.Errorf("inspect %s.%s: %w", schema, table, err)
	}
	defer rows.Close()
	var cols []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cols, nil
}

func connTableExists(ctx context.Context, conn *sql.Conn, schema, table string) (bool, error) {
	var name string
	query := fmt.Sprintf(`SELECT name FROM %q.sqlite_master WHERE type = 'table' AND name = ?`, schema)
	err := conn.QueryRowContext(ctx, query, table).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect %s.%s: %w", schema, table, err)
	}
	return true, nil
}

func (s *Store) currentSchemaVersion(ctx context.Context) (int, error) {
	applied, err := s.AppliedMigrations(ctx)
	if err != nil {
		return 0, err
	}
	if len(applied) == 0 {
		return 0, nil
	}
	sort.Ints(applied)
	return applied[len(applied)-1], nil
}

// BackupFileName builds the conventional timestamped name used by the UI so
// repeated backups sort chronologically in a directory listing.
func BackupFileName(at time.Time) string {
	return core.BackupFileName(at)
}

// LatestBackup returns the most recent backup file in dir, or an error when
// the directory holds none.
func LatestBackup(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "backup-") && strings.HasSuffix(name, ".db") {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no backups found in %s", dir)
	}
	// The timestamp format is fixed-width, so lexical order is chronological.
	sort.Strings(names)
	return filepath.Join(dir, names[len(names)-1]), nil
}
