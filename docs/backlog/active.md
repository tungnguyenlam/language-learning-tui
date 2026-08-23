# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Bounded storage performance pass complete; rotate to a new TUI hot-path area
without revisiting already-measured storage candidates.

## Exact Next Action

Inspect one unmeasured TUI render path for repeated full-buffer conversion,
sorting, styling of off-screen rows, or other per-message work. Establish a
focused benchmark/allocation or complexity baseline before selecting a change.

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

- The next candidate comes from a previously unmeasured render path and has a
  reproducible benchmark/allocation pattern or clear repeated-work argument.
- Any optimization preserves rendered output, hitboxes, compact layouts, and
  Unicode behavior with focused coverage.

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

## Current Batch Evidence

- The Statistics view consumes only the last seven `CardsAddedPerDay` keys,
  while storage previously grouped the full `notes` table.
- On 100,000 synthetic notes with 90% historical rows, the existing plan showed
  `SCAN n` and took about 42 ms; a seven-day range with
  `idx_notes_created_at` showed an indexed range search and took about 7 ms.
- Regression setup exposed that `date(n.created_at, 'localtime')` cannot parse
  Go-formatted timestamps, leaving the chart empty. The query now follows the
  indexed timestamp notice by grouping `substr(created_at, 1, 19)`, and new
  note creation timestamps are normalized to UTC.
- Migration 27 adds the created-time index. Focused tests protect global/deck
  seven-day counts, exclusion of older notes, migration application, and the
  indexed query plan; they pass.
- `go test ./internal/storage/sqlite -count=1`, `go test ./... -count=1`, and
  `./scripts/verify.sh` pass; the core E2E suite reports 35 passed in 37.69s.
- Review-summary reconnaissance on 100,000 synthetic reviews measured about
  14 ms for the separate success/grouped-grade scans (plus the total pass),
  versus about 2.3 ms for one grouped pass over a covering grade index.
- Migration 28 adds `idx_reviews_grade`; `Statistics` now derives total,
  successful, and per-grade counts from that one grouped result. Focused tests
  protect global/deck scoping, unexpected persisted grades, and a no-temp-B-tree
  covering-index plan; they pass.
- `go test ./internal/storage/sqlite -count=1`, `go test ./... -count=1`, and
  `./scripts/verify.sh` pass; the core E2E suite reports 35 passed in 37.65s.
- On 100,000 synthetic cards, separate maturity and flag passes took about
  58 ms combined; one shared card/state/flag aggregate returned identical
  values in about 43 ms (roughly 26% faster).
- `Statistics` now scans those one-to-one joins once for all eight counters.
  Focused global/deck tests cover absent join rows, due/suspended interactions,
  and the single-card-scan plan; they pass.
- `go test ./internal/storage/sqlite -count=1`, `go test ./... -count=1`, and
  `./scripts/verify.sh` pass; the core E2E suite reports 35 passed in 37.36s.
- Up-to-date `Migrate` benchmarking across three 200-iteration runs measured
  about 59–87 microseconds in memory and 80–82 microseconds file-backed. The
  cost is marginal, so no migration-bookkeeping code was changed.
- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).
- Prompt-only pass: manual consistency review and `git diff --check` passed on
  2026-08-23; no executable code changed.

## Repository State

- Streak optimization is committed as `4e77712`, Cards Added statistics as
  `236c63f`, review-summary consolidation as `2161297`, and shared card
  statistics as `5ea61d3`; the negative startup finding and next-area handoff
  are ready for their documentation checkpoint.
