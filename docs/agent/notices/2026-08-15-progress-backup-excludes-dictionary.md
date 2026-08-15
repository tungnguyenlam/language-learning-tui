# Progress Backup Excludes Dictionary

Status: active
Scope: storage backup, Import/Export TUI
Related: `internal/core.BackupRepository`, `internal/storage/sqlite/backup.go`, `internal/tui/actions_backup.go`

## Why It Matters

Progress backups copy decks, notes, cards, review history, flags, and app settings.
They deliberately omit the dict.cc tables (~834k rows). Restore must not wipe or
replace dictionary data.

## Required Behavior

- TUI backup/restore must type-assert `core.BackupRepository`. Do not import
  `internal/storage/sqlite` from `internal/tui`.
- Write timestamped files to `{dataDir}/backups/backup-YYYYMMDD-HHMMSS.db`
  using `core.BackupFileName`.
- Restore requires the existing deletion-confirmation prompt (`y/n`).
- If a store does not implement `BackupRepository`, degrade with a status
  message instead of panicking.

## Revisit When

Dictionary data is small enough to include, or a separate dictionary-backup
path is added.
