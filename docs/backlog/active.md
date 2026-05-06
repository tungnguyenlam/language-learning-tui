# Active Backlog

Last updated: 2026-05-05

## Current Milestone

TUI Refactor Completion & Final Verification (STABILIZED)

## Completed Work

### E2E Test Optimization & Stabilisation
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

## Next Action

- [ ] Add more complex grammar MCQ templates (e.g. passive voice, subjunctive)
- [ ] Implement deck-level statistics export

## Verification State

- Last verified: `./scripts/verify.sh` - 102/102 E2E tests pass, all Go tests pass.
- Environment: Darwin, Go 1.x, Python 3.12.
- Stability: Exceptional. Core data integrity issues (Undo flags, multiline TSV) resolved. Multi-agent code review uncovered and fixed long-standing backend edge cases.
