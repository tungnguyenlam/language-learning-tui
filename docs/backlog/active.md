# Active Backlog

Last updated: 2026-08-15

## Current Milestone

Local progress backup/restore is available on Import/Export. Settings +/- only
change the focused row. Review j/k card navigation is owned by the Review screen.

## Exact Next Action

No unfinished executable work remains. A future pass can migrate Settings, Cram, AI,
Practice, or Dictionary individually from `keys.go` into its registered screen
implementation.

## Completed This Pass

- Aligned `sqlite.Store` with `core.BackupRepository` and added Import/Export `B`
  backup / `U` restore to `{dataDir}/backups/`.
- Settings `+/-` now adjust Daily Goal or Reveal Speed only when that row is
  focused; the provider-cycle hint includes ollama.
- Review j/k card movement lives on `reviewScreen`; the leftover global `p`
  audio fallback was removed.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- `internal/tui/keys.go` still contains legacy handlers for Settings, Cram, AI,
  Practice, and Dictionary; migrate them individually behind the existing screen
  contract rather than attempting a broad rewrite.

## Acceptance Criteria

- `B` on Import writes a timestamped progress backup and reports row count.
- `U` on Import confirms, then restores decks and due cards from the latest backup.
- Settings `+/-` on the AI Provider row does not change the daily goal.
- Review `j`/`k` move the due-card cursor through the screen contract.
- Focused tests, `go test ./internal/tui`, and `./scripts/verify.sh` pass.

## Last Verification

- `go test ./internal/core ./internal/storage/sqlite ./internal/tui` passed on
  2026-08-15.
- `go test ./...` passed on 2026-08-15.
- `./scripts/verify.sh` passed on 2026-08-15: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 37.73s).

## Repository State

- Improvement pass is complete and fully verified on `main`.
- The five `subagent-*` branches remain checked out in external worktrees under
  `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
