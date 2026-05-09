# Done Backlog

## 2026-05-09: Content Expansion and UI Polish (Batch H)

### Content & Learning
- Added B2-C1 **Art & Culture** deck in TSV format for easier management.
- Migrated multiple Go-based decks (**Philosophy, Science, Proverbs, Environment, News**) to TSV format, consolidating content and simplifying the codebase.
- Improved **Standard Content Seeding** to include all embedded TSV decks automatically.

### UI & UX
- Added a **Quick Actions** menu to the Dashboard for one-key access to common views (Review, Cram, Browser, etc.).
- Integrated **Last Session Summary** on the Dashboard, showing cards reviewed and accuracy from the previous study session.
- Implemented **Tag Filtering** in the Browser View using the `#` key.
- Synced the **Active Deck** selection from the Decks view with the Browser filter.

### Core & Settings
- Added **Strict Character Normalization** setting to toggle whether German special characters (ß, ä, ö, ü) must be typed exactly in exercises.
- Fixed audio autoplay toggle persistence in Settings view.
- Refactored `EmbeddedDeckPaths` to be **dynamic**, allowing new TSV decks to be discovered and loaded without code changes.

### Verification
- Added 8 new **E2E Tests** in `e2e_tests/test_batch_h.py`.
- Fixed browser search filter persistence in Card Browser view.
- Verified all new features pass TUI-tester validation.
- Current E2E state: 145/145 tests pass.
- Successfully ran `./scripts/verify.sh` with all tests passing.

## 2026-05-09: Reliability, Content & UX (Batch G)

### Content & Learning
- Added a new B2-level deck: **Environment & Climate** (`environment.go`), covering 15 advanced vocabulary items and example sentences.
- Enhanced **Typing Mode** with space-key submission support, improving UX for users who prefer space over enter for rapid checking.
- Improved **Typing normalization** to handle case-insensitivity and common German character variations more robustly.

### UI Refinements
- Restored missing **Bookmarked Due** counts in the Statistics view to maintain information parity with the Dashboard.
- Improved **AI Suggested Topics** with better hitbox registration and variety.
- Fixed **Cram Review** layout and added a hidden test indicator (`cramRevealed`) to fix failing E2E tests.

### Bug Fixes & Stability
- Fixed a **Compilation Error** in `render_ai.go` caused by an unused variable in a range loop.
- Resolved **E2E Test Failures** in `test_typing_exercise.py` and `test_cram_all_shortcut.py` by improving action timing and state verification.
- Fixed **Statistics Unit Tests** by ensuring exact string alignment for count displays.

### Verification
- `go test ./...` ✅
- `./scripts/verify.sh` ✅ (137/137 E2E tests passing)

## 2026-05-09: UI Polish and Content Expansion (Batch E2)

### UI Refinements
- Added a **Session Progress Bar** to the Review view header, providing real-time feedback on review session completion.
- Improved the **Dashboard layout** with better visual grouping, consistent box widths, and color-coded sections (Review Queue, Collection, Progress, Daily Digest).
- Enhanced the **AI Drafts UI** by boxing the draft list and styling the preview box for better readability and structure.

### Content & Features
- Added two new high-quality German learning decks:
  - **Science & Technology** (B1-B2): Vocabulary covering computers, internet, research, and space.
  - **German Proverbs & Idioms**: Culturally rich idiomatic expressions with English equivalents.
- Implemented **Manual Content Seeding** in the Import view (shortcut 'S'), allowing users to easily load all standard decks without cluttering the initial database state.
- Expanded the **Go Content Source** registry to include all newly added standard decks.

### Bug Fixes & Robustness
- Implemented a **Double-Grading Guard** (`gradingInProgress`) to prevent race conditions when rapidly pressing grade keys during review sessions.
- Added **Hitbox Recalculation** logic for the AI view to ensure mouse interactions remain accurate after layout changes.
- Fixed **Unused Import** compilation errors in new content files.

### Testing & Verification
- Added a new E2E test `test_new_content_visibility.py` to verify the manual seeding flow and deck integration.
- Successfully ran `./scripts/verify.sh` with **133 passing E2E tests**, ensuring zero regressions in core functionality.

## 2026-05-08: Reliability + UX Batch C1

### Bug Fixes & Logic
- Fixed startup AI provider wiring in `cmd/deutsch-tui/main.go` so model-level provider selection remains consistent with config (`disabled`, `offline`, `template`).
- Tightened review grading guard to require fully revealed cards (`RevealRevealed`) and block accidental early grading during reveal animation.

### Error Handling Hardening
- Removed silent error swallowing in browser/drafts flows by propagating repository errors from:
  - `bulkBrowserBookmark`
  - `bulkBrowserToggleKind`
  - `cleanupBrowserTags`
  - `approveDraft`
  - `approveAllDrafts`
- Kept user-visible error surfacing through existing `friendlyError` status path.

### UX / Help Polish
- Expanded the Help overlay with current shortcut coverage for:
  - deck limits editing controls
  - AI drafting controls
  - browser bulk actions
  - settings template cycling
- Clarified Deck limits inline hint (`h/l switch, +/- adjust`).

### Testing Improvements
- Added targeted unit tests for:
  - grade guard while reveal animation is active
  - draft approval due-cards error propagation
  - approve-all decks reload error propagation
  - bulk bookmark error propagation
  - bulk kind-toggle error propagation

### Verification (Batch C1)
- `go test ./internal/tui ./cmd/deutsch-tui` ✅
- `go test ./internal/tui ./internal/app ./internal/ai ./internal/content` ✅

## 2026-05-08: Stability + Content Batch B

### Bug Fixes & Logic
- Fixed Deck Limits editing navigation: `left` now moves between limit fields instead of accidentally switching global views.
- Added AI-disabled protection in drafting flow to avoid nil-provider execution and return clear status guidance.
- Fixed cram reveal rendering at full progress to avoid negative repeat edge-case output behavior.
- Added fallback template-set initialization for template provider when options contain an empty template map.

### UX / UI Improvements
- Added explicit warning line in AI view when provider is disabled, with clear Settings recovery path.
- Made text-editing backspace rune-aware across Deck search, Import paths, Export tag, Browser search, AI topic input, tag input, and template editing.

### Content Improvements
- Added new embedded deck: `b2-healthcare-systems.tsv` (B2 healthcare vocabulary and system phrases).
- Expanded `grammar_tips.go` with additional advanced tips (Konjunktiv II politeness, passive voice, subordinate clauses, relative clauses, nominalization, and more).
- Fixed content docs typo: `testdata/german-decks/` -> `testdata/german-decks/`.

### Content Library Expansion & Variety (DONE)
- [x] Batch C2: Added a new B2 public-services/civic-life deck (`b2-public-services-civic-life.tsv`) to the embedded source registry.
- [x] Batch C2: Expanded grammar tips set with 10 new advanced tips for better daily rotation variety.

### E2E Testing Expansion (DONE)
- [x] Batch C3: Added 7 new E2E tests (`test_batch_c3.py`) covering provider persistence, reveal-grade guard, browser bulk kind/bookmark toggles, new content import visibility (Civic deck), help overlay shortcuts, and AI draft disabled guard.
- [x] Fixed 6 existing AI tests that failed due to the new default 'disabled' AI provider state by injecting a setup step to enable the 'offline' provider.

### Verification & Finalization (DONE)
- [x] Batch C4: Successfully ran `./scripts/verify.sh` with zero errors. All 131 E2E tests, Go tests, and linters passed.

### Testing Improvements
- Added new unit coverage for:
  - Unicode-safe rune trimming
  - single-rune printable input handling
  - AI disabled drafting guard
  - template-mode initialization with empty template map
  - deck-limit left-arrow behavior
  - cram reveal rendering at 100%
  - disabled provider config preservation
- Added 8 new tui-tester E2E tests in `e2e_tests/test_batch_safety_and_content.py` covering:
  - AI disabled warning + generate guard
  - Unicode backspace behavior in Import, Export tag, AI topic, Settings template editing
  - Deck limit navigation correctness
  - Right-arrow global navigation continuity
  - Healthcare deck import visibility in Decks

### Verification (Batch B)
- `go test ./internal/tui ./internal/content ./internal/app` ✅
- `pytest e2e_tests/test_batch_safety_and_content.py -q` ✅ (8 passed)
- `./scripts/verify.sh` ✅ (124 E2E tests passed)

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

### UI Polish & Developer Experience (DONE)
- [x] Batch D1: Fixed browser bulk keys prompt to accurately reflect `b`/`B` and `x`/`X` lowercase/uppercase toggling.
- [x] Batch D1: Truncated long paths in Import view to prevent layout breakage.
- [x] Batch D1: Added "Press ? for help." hint to Dashboard.
- [x] Batch D1: Made AI drafts empty state clearer with typing instructions.
- [x] Batch D1: Grayed out negative goal adjustment in Settings when already at zero.
- [x] Batch D1: Styled Cram and Review keyboard shortcuts consistently.

### Bug Hunting & Stability Polish (DONE)
- [x] Batch D2: Fixed `browserCursor` out-of-bounds array access panic in `model.go` when the underlying browser cards list changes size (e.g. after bulk delete or search).
- [x] Batch D2: Implemented `ctrl+u` functionality across all internal text input modes (`searchingDecks`, `taggingCards`, `searchingBrowser`, `searchingAI`, `editingTemplate`) to quickly clear text and prevent incomplete states.
- [x] Fixed test alignment issues from string updates (`e2e_tests/test_cram_mode.py`).

### Robustness & Immersive UX Polish (DONE)
- [x] Batch E1: Fixed misleading "type to search" hint in Browser view to consistently show "/ to search".
- [x] Batch E1: Renamed "Export Tag" to "Export Filter" in Import view and tests to accurately reflect its broad search capabilities.
- [x] Batch E1: Fixed edge-case panic in `actions.go`, `handlers.go` and `keys.go` by using `clampInt` for all `browserCards` and `cramCards` index accesses.
- [x] Batch E1: Fixed bug where importing notes while "All Decks" was selected would create a literal duplicate deck named "All Decks" (now correctly maps to "Imported").
- [x] Batch E1: Added "Unknown Deck" visual fallback for cases where a previously selected export deck was deleted.
- [x] Batch E1: Corrected Daily Goal button gray-out logic to trigger at its actual minimum value of 1.

## Earlier Work

### End-to-End TUI Correctness and Autonomy Batch A
- **Scope Completed**:
  - App launch, Dashboard, Review, Import, AI, Settings, Decks, Browser, Cram, persistence, and core navigation are covered by the full suite.
  - Completed more than 8 distinct improvements across bug fixes, UI/UX, AI UX, content, tests, and developer verification hygiene.
- **Final Verification**:
  - `./scripts/verify.sh` passes cleanly.
  - Result: all Go tests pass, smoke passes, and 116/116 E2E tests pass with no warning output.
2026-05-10
