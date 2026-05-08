# Active Backlog

Last updated: 2026-05-08 (batch C1)

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

- [x] Batch C1: fixed startup AI provider wiring so disabled/template config now initializes provider correctly in TUI model options.
- [x] Batch C1: hardened TUI repo command flows by propagating errors in draft approval, bulk bookmark toggle, bulk kind toggle, cleanup tags, and approve-all draft reload path.
- [x] Batch C1: blocked grading while reveal animation is in progress (grade now requires fully revealed state).
- [x] Batch C1: expanded help overlay with AI/deck-limits/browser bulk shortcuts and clarified deck-limit inline hints.
- [x] Batch C1: added unit tests for new error-propagation paths and reveal-state grade guard.
- [x] Batch C2: add a new B2 public-services/civic-life deck and expand grammar tips set for daily rotation variety.
- [x] Batch C3: add 7 new E2E tui-tester tests covering provider persistence on restart, reveal-grade guard, bulk-action failures, and new content import visibility.
- [x] Batch C4: run full `./scripts/verify.sh`, finalize backlog summaries, and commit.
- [x] Batch D1 (UX Polish): fixed browser bulk keys prompt to accurately reflect b/B and x/X toggling, truncated long paths in Import view, added help hint to Dashboard, made AI drafts empty state clearer, grayed out negative goal adjustment in Settings when at zero, and styled Cram/Review shortcuts.
- [x] Batch D2 (Bugfixes & Logic): fixed `browserCursor` out-of-bounds array access panic on bulk delete/conversion, added `ctrl+u` shortcut to safely and quickly clear input across all 5 text input modes (`searchingDecks`, `taggingCards`, `searchingBrowser`, `searchingAI`, `editingTemplate`).
- [ ] Batch D3: finalize and verify all changes.

## Verification State

- Last verified: `./scripts/verify.sh` passes with no warnings; Go suites pass; E2E now at 124 passing tests (including 8 new batch tests).
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Scalable content Registry with pluggable sources
