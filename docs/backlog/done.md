# Done Backlog

## 2026-05-08: Enhanced TUI Experience Batch

### UI Improvements
- **Keyboard Shortcut Feedback**: Added visual highlighting for keyboard shortcuts in all views to make them more discoverable
- **Visual Hierarchy**: Enhanced the Review view with improved card styling and better prominence for key elements
- **Statistics Visualization**: Improved the Statistics view with better data presentation and visual elements
- **Loading Indicators**: Added enhanced spinner animations for AI draft generation

### Content Additions
- **New Decks**: Added two new German learning decks:
  - B1 Travel & Tourism (`b1-travel-tourism.tsv`) - 97 cards covering travel vocabulary
  - B1 Technology & Internet (`b1-technology-internet.tsv`) - 120 cards covering modern tech terms
- **Expanded Existing Decks**: Both decks follow the established format with CEFR-level tagging

### Bug Fixes & Stability
- **Race Condition Fix**: Addressed potential race condition in deck selection by improving synchronization
- **Enhanced Error Handling**: Improved file operation error messages with more user-friendly feedback
- **Logging System**: Added comprehensive leveled logging throughout the application for better debugging

### Developer Experience
- **Help Documentation**: Enhanced the built-in help system with clearer organization and better formatting
- **Code Quality**: Maintained backward compatibility with existing E2E tests while adding new functionality

### Verification
- `go test ./...` passes with no failures
- `./scripts/tui_smoke.sh` passes successfully
- `./scripts/verify.sh` passes with 116 E2E tests passing
- All existing functionality preserved while adding new features

## Previous Work (from earlier today)

### UI Polish & Developer Experience (DONE)
- [x] Refined UI spacing and layout in Settings: aligned keys, updated formatting.
- [x] Enhanced AI Drafts view with placeholder strings for better context.
- [x] Improved Makefile for easier local developer testing (`make fmt`, `make build`, etc.`).
- [x] Created and successfully ran new E2E test `e2e_tests/test_settings_ui.py` for Settings UI layout.

### Content Library Expansion (IN PROGRESS)
- [x] Expand Grammar Essentials deck from 14 to 60+ cards
- [x] Expand German Idioms deck from 10 to 40+ cards
- [x] Add new Workplace & Office German deck (`b1-workplace-office.tsv`)
- [x] Add new Shopping & Services deck
- [x] Add new Health & Body deck
- [x] Add new Emotions & Feelings deck
- [x] Add new Transport & Directions deck

### Bug Hunting & Stability Polish
- [x] Fixed E2E test truncation flakiness by configuring `lines=40` for specific tests (`test_audio_autoplay.py`, `test_ui_sanity.py`).
- [x] **TUI Bug Fixes**: Resolved cursor panics, leaky modals, and ignored errors.
- [x] **Data Integrity**: Added transactions, fixed timezone inconsistencies, and enabled foreign keys in connection pool.
- [x] **Performance**: Indexed `reviewed_at` and clamped integer overflow in SRS scheduler.

## Earlier Work

### End-to-End TUI Correctness and Autonomy Batch A
- **Scope Completed**:
  - App launch, Dashboard, Review, Import, AI, Settings, Decks, Browser, Cram, persistence, and core navigation are covered by the full suite.
  - Completed more than 8 distinct improvements across bug fixes, UI/UX, AI UX, content, tests, and developer verification hygiene.
- **Final Verification**:
  - `./scripts/verify.sh` passes cleanly.
  - Result: all Go tests pass, smoke passes, and 116/116 E2E tests pass with no warning output.