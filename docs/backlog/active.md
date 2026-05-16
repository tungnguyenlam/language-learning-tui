# Active Backlog

Last updated: 2026-05-16 12:00 +07

## Current Milestone

Autonomous end-to-end improvement pass on 2026-05-16 (batch 11) - COMPLETED.

## Completed Work

- [x] UX Improvement: Added interactive "Recently Studied" decks on Dashboard (Clickable + `!/@/#` shortcuts).
- [x] UX Improvement: Enhanced Review view with prominent red "LEECH" badges and bolder state badges.
- [x] UI Improvement: Added celebratory "GOAL MET 🏆" badge on Dashboard when daily goal is reached.
- [x] Content Expansion: Added "B2 Business German: Meetings & Negotiations" deck (30 notes).
- [x] Logic Fix: Ensured statistics reload immediately after recording a review (fixed `0/X` reviews bug).
- [x] Added 3 new E2E tests in `test_batch11_improvements.py`.
- [x] UX Improvement: Added 1-4 grading shortcuts in Review view (verified with E2E test).
- [x] UX Improvement: Added Mouse Wheel support for Decks view scrolling.
- [x] Content Expansion: Added "B2 Programming & Software Engineering" deck (30 technical notes).
- [x] Fixed E2E regression in `test_mcq_navigation_bug.py` caused by new grading shortcuts.
- [x] Fixed 10 failing E2E tests (Batch 9)
- [x] Verified B2 Job Application, C1 Philosophy & Ethics, A1 City & Directions decks exist (40 cards each)
- [x] UI Improvement: Enhanced session summary with efficiency stat (cards/min)
- [x] UI Improvement: Added motivational tips in Review empty state
- [x] AI Improvement: Added Grammar Breakdown and City Directions topics to AI suggestions
- [x] Added 8 new E2E tests for new features

## Exact Next Action

None - milestone complete. All vision and prompt strategy work committed and pushed.

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
- `./scripts/verify.sh` passed.
- 315 E2E tests passing.

## Blockers

None.
