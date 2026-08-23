# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Optimize the SQLite statistics streak hot path demonstrated by representative
query-plan and timing evidence.

## Exact Next Action

Commit the verified streak batch, then inspect the `CardsAddedPerDay` consumer
contract and capture focused plan/timing evidence for its full-history
date-grouping query before selecting the next storage change.

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
- Committed Settings persistence as `fc402fa`.
- Reconnaissance found `getCompoundBreakdown` calls
  `DictionaryRepository.Exists` synchronously on cache misses; four Dictionary
  render sites call it, so first render performs SQLite work from `View()`.
- Compound breakdown lookup is now cache-only in render/key paths; validated
  decomposition runs in a two-second command after current search/cursor
  changes and preserves related-entry/recent-view command batches.
- Compound responses and errors carry the search identity, so stale enrichment
  cannot populate cache or overwrite current status.
- Regression coverage proves rendering makes zero `Exists` calls, precomputed
  validated parts still render and drive shortcuts, and stale replies are
  ignored; focused and race tests pass.
- Full repository verification passed for async compound precomputation; the
  batch is ready to commit.
- Follow-up reconnaissance found `refreshLastBackupPath` calls `os.ReadDir`
  synchronously when entering Import and when Restore has no cached path.
- Committed async compound precomputation as `cbee406`.
- Import entry and uncached Restore now return an identity/view-guarded backup
  discovery command; confirmation starts only after the latest result arrives.
- Missing backup directories/files keep the existing no-backup guidance,
  unexpected filesystem errors reach the TUI error path, and stale discovery
  replies are ignored.
- Focused Import/backup tests and the TUI race suite pass.
- Full repository verification passed for async backup discovery; the batch is
  ready to commit.
- Committed async backup discovery as `b9ec6c0`.
- Deck TSV export now creates `exports/` only inside its existing timeout-bound
  command and returns directory failures through the TUI error path; empty
  decks no longer create an unused output directory.
- Regression coverage proves command construction performs no filesystem work
  and directory-creation failures are visible; focused, package, and race tests
  pass.
- Full repository verification passed for the Deck TSV directory batch; it is
  ready to commit.
- Committed Deck TSV directory handling as `9f8c6ed`.
- Storage reconnaissance seeded a temporary 20,000-card/100,000-review database
  and ran `ANALYZE` plus `EXPLAIN QUERY PLAN` for `DueCards` and dashboard
  statistics queries.
- `DueCards` averaged about 18 ms despite its card scan; `Statistics` averaged
  about 709 ms. Its current-streak query accounted for about 524 ms and showed
  `SCAN reviews USING COVERING INDEX idx_reviews_reviewed_at` plus temporary
  B-trees for both `GROUP BY` and `ORDER BY`.
- Global and deck streak queries now stream newest reviews and calculate local
  calendar days in Go, skipping duplicate same-day reviews and stopping at the
  first gap or the existing 365-day cap.
- Permanent regression coverage protects gap, duplicate, deck-scoped, cap, and
  query-plan behavior; focused storage tests pass.
- Representative `Statistics` time fell from about 709 ms to about 243 ms
  (roughly 66%) with all returned fields preserved.
- `go test ./internal/storage/sqlite -count=1`, `go test ./... -count=1`, and the
  full repository verification gate pass for the streak batch.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Global and deck streaks retain local-calendar semantics, ignore duplicate
  reviews on one day, stop at gaps, and remain capped at 365 days.
- The global streak query reads `idx_reviews_reviewed_at` in descending order
  without temporary grouping/sorting, protected by regression/plan coverage.
- Representative `Statistics` timing improves materially without changing its
  returned fields.

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
- Compound batch: focused tests and `go test -race ./internal/tui -count=1`
  passed on 2026-08-23.
- Compound batch: `go test ./...` and `./scripts/verify.sh` passed on
  2026-08-23; core E2E suite 35 passed in 37.77s.
- Backup discovery batch: focused tests, `go test ./internal/tui -count=1`, and
  `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- Backup discovery batch: `go test ./...` and `./scripts/verify.sh` passed on
  2026-08-23; core E2E suite 35 passed in 37.55s.
- Deck TSV directory batch: focused tests, `go test ./internal/tui -count=1`,
  and `go test -race ./internal/tui -count=1` passed on 2026-08-23.
- Deck TSV directory batch: `go test ./...` and `./scripts/verify.sh` passed on
  2026-08-23; core E2E suite 35 passed in 37.76s.
- Streak batch: `go test ./internal/storage/sqlite -count=1` and
  `go test ./... -count=1` passed on 2026-08-23.
- Streak batch: `./scripts/verify.sh` passed on 2026-08-23; Go tests, vet,
  offline dict.cc import (834,512 entries), smoke test, binary build, and core
  E2E suite 35 passed in 37.39s.
- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).
- Prompt-only pass: manual consistency review and `git diff --check` passed on
  2026-08-23; no executable code changed.

## Repository State

- The synthetic reconnaissance harness was removed; only the streak
  implementation, permanent tests, and continuity documentation remain.
