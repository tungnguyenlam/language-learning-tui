# Active Backlog

Last updated: 2026-05-08 (batch B)

## Current Milestone

Making the TUI meaningfully better with 6-8 distinct improvements

## Completed Work

### UI Improvements
- [x] Add visual feedback for keyboard shortcuts in all views
- [x] Add new Travel & Tourism German deck
- [x] Add new Technology & Internet German deck
- [x] Improve the visual hierarchy in the Review view
- [x] Enhance the Statistics view with better data visualization
- [x] Add a loading indicator for AI draft generation

### Bug Fixes
- [x] Fix potential race condition in deck selection
- [x] Improve error handling for file operations

### Developer Experience
- [x] Add more comprehensive logging
- [x] Improve the help documentation

### UI Polish & Developer Experience (DONE)
- [x] Refined UI spacing and layout in Settings: aligned keys, updated formatting.
- [x] Enhanced AI Drafts view with placeholder strings for better context.
- [x] Improved Makefile for easier local developer testing (`make fmt`, `make build`, etc.).
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

## Next Action

- [x] Batch B1: fixed Deck limits left-arrow behavior (now edits limit cursor instead of switching views) and preserved right-arrow global view cycling.
- [x] Batch B1: added AI disabled guard in drafting flow + explicit AI disabled guidance in AI view.
- [x] Batch B1: fixed cram reveal rendering at 100% progress to avoid negative repeat edge cases.
- [x] Batch B1: improved text input safety with rune-aware backspace for deck filter, import/export paths, export tag, browser search, tag input, AI input, and settings templates.
- [x] Batch B2: added new B2 healthcare systems deck (`b2-healthcare-systems.tsv`) and expanded grammar tips coverage.
- [x] Batch B3: added 8 new E2E tests in `e2e_tests/test_batch_safety_and_content.py` covering AI disabled behavior, Unicode backspace flows, deck-limit navigation behavior, arrow-view navigation, and healthcare deck import visibility.
- [x] Phase 4: run full `./scripts/verify.sh`, update done backlog summary, and commit.

## Verification State

- Last verified: `./scripts/verify.sh` passes with no warnings; Go suites pass; E2E now at 124 passing tests (including 8 new batch tests).
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Scalable content Registry with pluggable sources
