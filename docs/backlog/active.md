# Active Backlog

Last updated: 2026-08-23

## Current Milestone

TUI responsiveness: move Dictionary history/recent/star persistence out of
synchronous keyboard and mouse input handlers.

## Exact Next Action

Convert Dictionary search-history record/clear paths to return ordered,
snapshot-based persistence commands and propagate them through every keyboard,
navigation, and hitbox action without dropping the accompanying user command.

## Completed This Pass

- Confirmed `recordDeckSearch` and every Decks history-clear path call
  `Repository.SetSetting` synchronously from input/hitbox handling and discard
  its error.
- Decks history record and clear paths now return snapshot-based persistence
  commands with a two-second timeout and contextual errors.
- Concurrent/stale commands are serialized and only the latest not-yet-started
  snapshot writes, preventing an older save from overwriting newer history.
- Added regression coverage for deferred writes, latest-snapshot behavior,
  surfaced errors, and mouse-clear command propagation.
- Full repository verification passed for the Decks history batch; it is ready
  to commit as a checkpoint before Dictionary work begins.

## Top Issues

- Dictionary search/recent-view/star persistence and Settings config/secrets
  callbacks still perform synchronous writes from TUI handlers.
- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Dictionary search/recent/star persistence performs no repository work until
  returned `tea.Cmd`s execute.
- Compound actions preserve both persistence and their primary command.
- Concurrent saves cannot let an older snapshot overwrite newer state.
- Repository errors reach the established TUI error status path.

## Last Verification

- `go test ./internal/tui -count=1` passed on 2026-08-23.
- `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- `go test ./...` passed on 2026-08-23.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 37.17s).
- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).
- Prompt-only pass: manual consistency review and `git diff --check` passed on
  2026-08-23; no executable code changed.

## Repository State

- Decks history persistence batch is verified and pending commit.
