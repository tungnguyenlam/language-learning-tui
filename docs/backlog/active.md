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

## Next Action

Monitor long-term stability of parallel tests and consider adding more German grammar content.

- [ ] Add more complex grammar MCQ templates (e.g. passive voice, subjunctive)
- [ ] Implement deck-level statistics export

## Verification State

- Last verified: `./scripts/verify.sh` - 100/100 E2E tests pass, all Go tests pass.
- Environment: Darwin, Go 1.x, Python 3.12.
- Stability: Very High. Logic bugs (MCQ, Modal trap, UTF-8 truncation) resolved.
