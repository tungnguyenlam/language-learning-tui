# Active Backlog

Last updated: 2026-05-10 14:30 +07 (Content, UI Polish & SRS Robustness)

## Current Milestone

Completed a multi-area improvement pass covering German content expansion, UI refinements (scrollbars, grammar examples), statistics visualization, and SRS robustness.

## Phase 1 Baseline

- All core views (Dashboard, Decks, Review, Statistics, Import, AI, Settings, Browser, Cram) verified.
- Initial test state: 1 failed (test_new_content_visibility.py), 165 passed.

## Top Issues / Work Plan

- [x] Fix E2E test robustness in `test_new_content_visibility.py`.
- [x] Add "German Confusable Words" deck and expand Grammar Tips with examples.
- [x] Implement visual scrollbar with hitboxes for Decks view.
- [x] Add Maturity Distribution chart to Statistics view.
- [x] Add comprehensive unit tests for SRS Scheduler.
- [x] Verify everything with `./scripts/verify.sh` and fix layout regressions.

## Exact Next Action

Completed. All improvements verified and tests passing.

## Acceptance Criteria

- App starts with zero errors.
- All core screens render correctly.
- New "Confusable Words" deck and Grammar Tip examples are visible.
- Decks view has a functional scrollbar.
- All Go tests and E2E tests (including new ones) pass.
- Backlog updated and changes committed.

## Last Verification

- 2026-05-10 14:25 +07: `./scripts/verify.sh` passed with 163 E2E tests and all Go tests.
- 2026-05-10 14:20 +07: `e2e_tests/test_ui_polish.py` passed individually.
- 2026-05-10 14:15 +07: `go test ./internal/srs/...` passed with new unit tests.
