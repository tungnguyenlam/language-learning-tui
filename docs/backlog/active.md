# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Bounded reconnaissance: identify the next evidence-backed correctness or
hot-path performance batch after removing synchronous TUI persistence writes.

## Exact Next Action

Commit the verified Settings persistence batch, then inspect TUI handlers for
remaining direct repository/filesystem/provider calls outside `tea.Cmd`s. If
none remain, rotate to one measured rendering or SQLite query hot path.

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
- Committed Dictionary search history as `a17284a`.
- Recently inspected words now mutate in memory during input and persist via
  ordered snapshot commands returned by keyboard navigation and full/Spotlight
  hitboxes; related-entry loads remain batched.
- Star toggles now persist a deterministic sorted ID snapshot asynchronously;
  active `:starred` filtering still refreshes in the same command batch.
- Recent/star write failures now reach the established visible error path, with
  focused regression and race coverage.
- Full repository verification passed for the recent/star batch; it is ready to
  commit before Settings callback work begins.
- Committed Dictionary recent/star persistence as `dfa303c`.
- Settings config and secrets callbacks now execute in ordered commands, return
  errors to the visible TUI path, and receive immutable snapshots (including a
  deep copy of nested AI templates).
- Provider/audio/normalization/reveal/template/credential keyboard and hitbox
  paths now return their save commands while updating visible state immediately.
- Focused tests cover deferred callbacks, latest-request suppression, mutable
  snapshot isolation, config/secrets errors, and race safety.
- Full repository verification passed for Settings persistence; the batch is
  ready to commit.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Reconnaissance produces a concrete reproduction/measurement or records that
  the inspected area has no actionable violation.
- The selected next batch is narrow, locally verifiable, and prioritized by the
  bug-fix/performance ordering in `prompt/improve.md`.

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
- Dictionary recent/star batch: focused tests and
  `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- Dictionary recent/star batch: `go test ./...` and `./scripts/verify.sh` passed
  on 2026-08-23; core E2E suite 35 passed in 37.42s.
- Settings batch: `go test ./internal/tui -count=1` and
  `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- Settings batch: `go test ./...` and `./scripts/verify.sh` passed on
  2026-08-23; core E2E suite 35 passed in 37.91s.
- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).
- Prompt-only pass: manual consistency review and `git diff --check` passed on
  2026-08-23; no executable code changed.

## Repository State

- Settings persistence batch is verified and pending commit.
