# Active Backlog

Last updated: 2026-08-23

## Current Milestone

TUI responsiveness: move Dictionary history/recent/star persistence out of
synchronous keyboard and mouse input handlers.

## Exact Next Action

Convert Dictionary recently-viewed and starred-entry writes/clears to ordered,
snapshot-based commands; propagate commands through cursor inspection,
keyboard actions, and hitboxes without dropping related loads/actions.

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
- Committed the Decks checkpoint as `3ae7568`.
- Dictionary search-history mutation is now in-memory only during handlers;
  screen/navigation boundaries and both full/Spotlight clear hitboxes return an
  ordered snapshot command.
- Navigation batches retain their destination load alongside the history save,
  and persistence errors flow through the TUI error status.
- Extracted the Decks ordering mechanism into a reusable per-value
  `orderedSave` coordinator without changing its verified behavior.
- Full verification passed for Dictionary search-history persistence; the batch
  is ready to commit before recently-viewed/starred work begins.

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
- Dictionary batch: focused Dictionary/Decks persistence tests and
  `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- Dictionary search-history batch: `go test ./...` and `./scripts/verify.sh`
  passed on 2026-08-23; core E2E suite 35 passed in 37.20s.
- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).
- Prompt-only pass: manual consistency review and `git diff --check` passed on
  2026-08-23; no executable code changed.

## Repository State

- Dictionary search-history batch is verified and pending commit.
