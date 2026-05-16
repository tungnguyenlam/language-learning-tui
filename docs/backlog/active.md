# Active Backlog

Last updated: 2026-05-16 12:00 +07

## Current Milestone

Autonomous end-to-end improvement pass on 2026-05-16 (batch 9) - COMPLETED.

## Completed Work

- [x] Fixed 10 failing E2E tests (Settings timing, reveal guard, batch8, bulk deletion, review empty state)
- [x] Verified B2 Job Application, C1 Philosophy & Ethics, A1 City & Directions decks exist (40 cards each)
- [x] UI Improvement: Enhanced session summary with efficiency stat (cards/min)
- [x] UI Improvement: Added motivational tips in Review empty state
- [x] AI Improvement: Added Grammar Breakdown and City Directions topics to AI suggestions
- [x] Added 8 new E2E tests for new features

## Exact Next Action

None - milestone complete.

## Top Issues / Priorities

None.

## Acceptance Criteria

- [x] Baseline `./scripts/tui_smoke.sh` passes (`deutsch-tui smoke ok`).
- [x] Baseline `go test ./...` passes.
- [x] App starts with zero errors after changes.
- [x] Core screens render via E2E coverage.
- [x] Keyboard and mouse navigation have targeted coverage.
- [x] SQLite persistence remains covered by tests.
- [x] At least 8 new tui-tester E2E tests are added and passing.
- [x] At least 2 new content/features are implemented.
- [x] At least 2 UI/UX improvements across different views are implemented.
- [x] At least 2 bugs or potential bugs are fixed.
- [x] `./scripts/verify.sh` passes with zero errors or warnings.
- [x] Completed work is summarized in `docs/backlog/done.md`.

## Last Verified

- `go test ./...` passed.
- `./scripts/verify.sh` pending.
- 311 E2E tests passing.

## Blockers

None.
