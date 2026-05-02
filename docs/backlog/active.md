# Active Backlog

Last updated: 2026-05-02

## Current Milestone

Autonomous Feature Pass 18: Final Polish and UI Interactive Elements

## Completed Work

### UI/UX Refinements
- [x] Fixed Statistics scrollbar logic to be more precise and visually representative
- [x] Implemented **interactive scrollbar clicking** for the Statistics view
- [x] Implemented **drag-to-scroll support** for scrollbars in Statistics, Browser, and Cram views
- [x] **Enhanced Decks view** with separate counts for New, Due, and Total cards, featuring color-coded statistics
- [x] **Added Session Statistics** to the status line (Accuracy and Reviewed count)
- [x] **Improved AI view** with interactive [Approve] and [Discard] buttons for each draft
- [x] **Added detailed Preview** to the AI view for the selected draft (showing Extra info, Tags, and Examples)
- [x] **Refactored Browser view** with an explicit **Search Mode** (via `/`) and quick card actions (**b** to bookmark, **x** to suspend)
- [x] **Enhanced Settings view** with interactive hitboxes, section headers, and clickable [+] [-] goal adjustment buttons
- [x] **Enhanced Import view** with interactive buttons and clickable path fields for easier collection management
- [x] **Added interactive Dashboard hitboxes** allowing users to click sections to navigate to Review, Browser, Statistics, etc.
- [x] **Added interactive Cram filters** allowing mouse-based switching between bookmarked, suspended, and other card sets
- [x] Reworked active panel hitbox coordinates to derive content origin from Lip Gloss frame metrics
- [x] Added interactive scrollbar tracks for Browser and Cram views
- [x] Unified all breakpoints (Compact, Medium, Wide) to use a consistent bordered panel for active views
- [x] Added `statsTotalLines` tracking to accurately calculate scroll ratios
- [x] Fixed all Go unit tests for `renderStatistics`

### Bug Fixes
- [x] Added missing `strconv` import in TUI layer
- [x] Resolved regression in `renderCompact` where panels were missing
- [x] Fixed TUI unit tests failing due to signature changes
- [x] **Fixed race condition in Review reveal** by only showing grading hints and hitboxes when fully revealed (synchronizing visual state with interaction)

### Testing
- [x] Added 3 new interactive E2E tests in `e2e_tests/test_interactive_features.py` (mouse navigation, browser movement, scrollbar interaction)
- [x] **Parallelized E2E tests** using `pytest-xdist`, achieving ~8x speedup (from 323s to <45s)

### Current Status (Verification)
- [x] App starts with zero errors
- [x] All views render correctly (Dashboard, Review, Import, AI, Settings, Browser, Cram)
- [x] Core user interactions respond as expected (flashcard reveal, grading, navigation)
- [x] State successfully persisted to SQLite
- [x] All Go unit tests pass (9 test suites)
- [x] All 71 E2E tests pass
- [x] ./scripts/verify.sh executes successfully

## Planned Features

(None at this time)

## Last Verification

- 2026-05-02: `./scripts/verify.sh` passed, including all 71 E2E tests.

## Next Action

Final commit and handoff.

## Summary

The deutsche-tui application is now production-ready with comprehensive German learning content:
- ✅ All 9 Go unit test suites passing
- ✅ All 68 E2E tests passing  
- ✅ App launches without errors
- ✅ All views render correctly
- ✅ Core interactions working
- ✅ State persistence verified
- ✅ 600+ vocabulary items across A1-B2 levels
- ✅ 500+ importable TSV cards
- ✅ Full Anki compatibility

The application successfully meets all requirements for end-to-end functionality.
