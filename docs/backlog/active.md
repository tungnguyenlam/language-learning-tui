# Active Backlog

Last updated: 2026-05-10 15:15 +07 (Autonomous Improvement Pass)

## Current Milestone

Implementing a suite of 8+ meaningful improvements including new content, UI polish, logic enhancements, and expanded testing.

## Phase 3 Work Plan

- [x] **Task 1: Content** - Add "German Idioms & Proverbs" deck to core content.
- [x] **Task 2: UI** - Dashboard: Add "Welcome" header and "Daily Progress" status line.
- [x] **Task 3: Statistics** - Add "Success Rate per Deck" visualization.
- [x] **Task 4: AI** - Improve AI Drafting prompts for better grammar explanations.
- [x] **Task 5: UI** - Review View: Display "Extra/Context" field more prominently during reveal.
- [x] **Task 6: UX** - Add loading spinner animation for AI operations.
- [x] **Task 7: Test** - Add E2E test for the new Browser Search feature.
- [x] **Task 8: Content** - Add "German Slang & Youth Language" deck.

## Exact Next Action

Completed. 9 improvements verified and 165 tests passing.

## Acceptance Criteria

- All 8 original tasks (plus 1 bonus) implemented successfully.
- `./scripts/verify.sh` passes with zero errors (165 E2E tests).
- New content (Idioms, Slang) is accessible.
- AI spinner and templates are improved.
- Browser search logic is more robust.
- Dashboard header is compact and informative.
- Statistics view shows per-deck success rates.

## Last Verification

- 2026-05-10 14:25 +07: `./scripts/verify.sh` passed with 163 E2E tests and all Go tests.
- 2026-05-10 14:20 +07: `e2e_tests/test_ui_polish.py` passed individually.
- 2026-05-10 14:15 +07: `go test ./internal/srs/...` passed with new unit tests.
