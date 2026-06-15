# Active Backlog

Last updated: 2026-06-15

## Current Milestone

Polish and Reliability - COMPLETED.

## Exact Next Action

Select next milestone from `docs/backlog/roadmap.md` (e.g., continue dictionary hybrid polish or move to `.apkg` import/export) after committing the current fixes.

## Top Issues

None active.

## Completed Work

- [x] **BUG-013 [FIXED]:** Mouse click on dashboard sparkline now navigates to Statistics. Added `dash-activity` hitbox action in `internal/tui/hitboxes.go`.
- [x] **BUG-012 [FIXED]:** Dictionary result highlighting no longer corrupts or overflows at truncation boundaries for multi-byte/ANSI content. Rewrote `truncateLine` in `internal/tui/utils.go` to correctly reserve space for the ellipsis, handle ANSI escape sequences, and keep output visual width within `maxWidth`. Added `TestTruncateLine` and `TestDictionaryHighlightAfterTruncation`.
- [x] **Parallel E2E flake on deck-list tests [MITIGATED]:** `test_new_a1_decks_loaded` now waits for the visible substring "Travel & Transport" and the deck's unique description instead of the full deck name, which could be truncated or wrapped under parallel PTY load. Full `./scripts/verify.sh` now passes (351 E2E tests).

## Verification

- `./scripts/verify.sh` passed: all Go unit tests, smoke test, binary build, and 351 E2E tests.
