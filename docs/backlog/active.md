# Active Backlog

Last updated: 2026-05-07

## Current Milestone

Content Library Expansion & Improvement

## Completed Work

### Content Library Expansion (IN PROGRESS)
- [x] Expand Grammar Essentials deck from 14 to 60+ cards
- [x] Expand German Idioms deck from 10 to 40+ cards
- [x] Add new Shopping & Services deck
- [x] Add new Health & Body deck
- [x] Add new Emotions & Feelings deck
- [x] Add new Transport & Directions deck

### Bug Hunting & Stability (DONE)
- [x] Fixed MCQ key collisions with global navigation.
- [x] Prevented starter deck re-insertion/overwrites in `main.go`.
- [x] Enabled SQLite foreign key constraints for cascading deletes.
- [x] Stabilized deck list cursor during background reloads.
- [x] Implemented stable deletion IDs to prevent race conditions.
- [x] Verified all 105 tests pass successfully.

### E2E Test Optimization & Stabilisation (DONE)
- [x] **Binary pre-build**: Modified `scripts/verify.sh` to build binary once, reducing resource contention.
- [x] **Parallel Test Fixes**: Updated all 28 E2E tests to use `DEUTSCH_TUI_BIN`.
- [x] **Hitbox Alignment**: Aligned hitboxes with Lip Gloss metrics in all views using `contentLayoutForStyle`.
- [x] **Deletion Confirmation**: Added modal confirmation dialogs for card and deck deletions.
- [x] **Transactional Storage**: Wrapped bulk storage updates in transactions for performance and atomicity.

### Template Editing Cancel Behavior Fix
- [x] **Added revert-on-cancel for template editing**: Esc during template editing now restores the original value instead of keeping changes.
- [x] **Stored original values**: Added `originalTemplateValue` field to preserve pre-edit state.

### Bug Hunting & Stability Polish
- [x] **TUI Bug Fixes**: Resolved cursor panics, leaky modals, and ignored errors.
- [x] **Data Integrity**: Added transactions, fixed timezone inconsistencies, and enabled foreign keys in connection pool.
- [x] **Performance**: Indexed `reviewed_at` and clamped integer overflow in SRS scheduler.
- [x] **Multi-Agent Deep Bug Hunt**: Uncovered and fixed out-of-bounds panics, unhandled `rows.Err()` iterations, resource leaks (`rows.Close()`), and incorrect layout height/width propagation.

## Next Action

- [x] Scalable content system with Registry pattern

## Verification State

- Last verified: `./scripts/verify.sh` - 105/105 E2E tests pass, all Go tests pass.
- Environment: Darwin, Go 1.x, Python 3.12.
- Content: 1366+ cards across 14 decks with explanations
- Architecture: Scalable content Registry with pluggable sources
