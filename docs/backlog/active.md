# Active Backlog

Last updated: 2026-05-05

## Current Milestone

TUI Refactor Completion & Final Verification

## Completed Work

### Template Editing Cancel Behavior Fix
- [x] **Added revert-on-cancel for template editing**: Esc during template editing now restores the original value instead of keeping changes.
- [x] **Separated Enter/Esc handling**: Enter keeps changes, Esc reverts (matching UI description).
- [x] **Stored original values**: Added `originalTemplateValue` field to preserve pre-edit state.
- [x] **Updated E2E tests**: Fixed test expectations in `test_robustness.py` for cancel behavior.
- [x] **Verified**: All 93 E2E tests passing, all Go unit tests passing.

### Run 2 TUI Refactor Recovery
- [x] Restored stable top-level rendering contracts, status footer text, tab/nav click targets
- [x] Fixed Settings/AI view surfaces, reveal animation state, scrollbar hitbox alignment
- [x] Verified `./scripts/verify.sh` passes with 83 E2E tests

### Decks View Bulk Operations (Pass 21)
- [x] **Implemented Multi-Select**: Added `x` key selection in Decks view.
- [x] **Bulk Deletion**: Added `Backspace` support for deleting selected decks.
- [x] **Deck Merging**: Added `M` key to merge selected decks into current.
- [x] **Repository Methods**: Implemented `DeleteDecks` and `MergeDecks` in SQLite.
- [x] **Verified Suite**: 91 tests passing.

### Export Filtering (Pass 21)
- [x] **Implemented Export Filtering**: Added deck/tag filters to Import/Export view.
- [x] **UI & Key Fixes**: Resolved key prioritization and hitbox alignment.
- [x] **Active Deck Sync**: Default export deck matches current view.
- [x] **Unused Tags Cleanup**: Implemented 'C' key in Browser to remove tags not used by any cards.
- [x] **AI Prompt Templates**: Refined templates to support multiple exercise types.

### Startup Compatibility Fix
- [x] **Legacy AI Template Config Migration**: `make run` accepts flat `ai_templates` config shape.
- [x] **Verified Startup**: `make run` launches TUI successfully.
- [x] **Verified Tests**: `./scripts/verify.sh` passes with 93 E2E tests.

## Next Action

Optimize E2E tests by building the binary once and address unconfirmed deletions.

- [ ] Build binary once in `scripts/verify.sh` and use in E2E tests
- [ ] Add confirmation dialogs for card/deck deletions
- [ ] Align Browser scrollbar hitboxes with Lip Gloss metrics
- [ ] Investigate and fix intermittent parallel test failures (if they persist)
