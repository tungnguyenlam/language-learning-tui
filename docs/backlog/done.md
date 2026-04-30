# Completed Work

## 2026-04-30 Agent Instruction Maintenance

- Added explicit root guidance that agents may update root/subtree instructions and continuity docs when needed to make future work safer or easier.
- Added durable notices for TUI status-line stability and Browser deck reload behavior.
- Updated `docs/agent/index.md` so future session pickup finds the new notices.

## 2026-04-30 Autonomous Feature Pass 7 - Status Stability and Browser Reload

- Reproduced fresh baseline failure in `./scripts/verify.sh`: 40 E2E tests passed and 7 failed because status messages wrapped across terminal lines.
- Changed TUI rendering to show status on its own stable single-line surface and shortened missing-file import errors.
- Fixed Browser `[`/`]` deck switching so the Browser card list reloads for the newly selected deck.
- Added unit coverage for Browser deck switching reload behavior.
- Added 3 tui-tester E2E regressions for grade status stability, concise missing-file import errors, and Browser deck switching after TSV import.
- Targeted verification passed for all 7 previously failing E2E tests and `go test ./internal/tui`.
- Targeted verification passed for the 3 new E2E tests.
- Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 50 tui-tester E2E tests.

## 2026-04-30 Autonomous Feature Pass 6 - Contextual Help and Additional Tests

- Added contextual help hints to the status bar that show relevant keyboard shortcuts based on the active view
- Created 3 new E2E tests covering uncovered UX behaviors:
  - Cram Mode filter types (suspended, leech, flagged, all)
  - Browser view deck switching functionality using [ and ] keys
  - Settings view daily goal adjustment with +/- keys
- Fixed edge case in ReviewsPerDay SQL query to handle NULL dates
- Extended Statistics view tests with additional coverage for streak indicators and review activity charts
- All Go tests pass, all E2E tests pass, `./scripts/verify.sh` passes with zero errors
- Final verification passed with 44 E2E tests

## 2026-04-30 Autonomous Feature Pass 6 - Tab Navigation Load Commands

- Fixed Tab navigation so it returns the `updateView` load command, matching arrow and mouse navigation behavior.
- Removed unused no-op view navigation wrappers that discarded command-backed loads.
- Added unit coverage proving Tab into Browser returns a load command.
- Added 3 tui-tester E2E tests covering Tab-loaded Browser cards, Tab-loaded Cram bookmarked cards, and Tab-loaded Statistics review progress.
- Added `/deutsch-tui` to `.gitignore` so local root builds are not surfaced as accidental source changes.
- Final verification passed with `./scripts/verify.sh`, including Go tests, go vet, smoke launch, and 41 tui-tester E2E tests.

## 2026-04-30 Autonomous Feature Pass 5 - Review History Visualizations

- Implemented Feature 15: Review History Visualizations (Heatmap/Chart).
- Added `ReviewsPerDay(ctx, days)` method to `core.Repository` interface and SQLite store implementation.
- Added `reviewsPerDay` field to TUI Model and `loadReviewsPerDay()` command.
- Extended `renderStatistics()` to display a 14-day ASCII bar chart of review activity.
- Updated `mockRepo` in `mocks_test.go` to implement the new `ReviewsPerDay` method.
- Statistics view now shows review activity chart with daily counts and visual bars.
- All Go tests pass, all 38 E2E tests pass, `./scripts/verify.sh` passes with zero errors.

## 2026-04-30 Autonomous Feature Pass 4 - Card Browser, Session Stats, Help Overlay

- Implemented Card Browser View: `8` opens browser showing all cards in current deck with type (FC/MCQ), prompt, and status flags (B/L/S).
- Added real-time search filtering by prompt/answer text; `j`/`k` navigation; `[`/`]` deck switching.
- Extended `core.Repository` with `Cards(ctx, deckID, search)` method; implemented in SQLite store with optional deck filter and search.
- Added session statistics tracking: `sessionReviewed`, `sessionCorrect`, accuracy % displayed in Statistics view.
- Implemented Keyboard Shortcut Help Overlay: `?` toggles overlay showing all shortcuts grouped by view.
- Added `showHelp` bool field in TUI model; renders help text with `renderHelp()`.
- Added 3 Go unit tests: `TestCards_DeckFilter`, `TestCards_SearchFilter`, `TestSessionStatsTracking`, `TestHelpOverlayToggle`, `TestBrowserViewNavigation`, `TestBrowserSearchFilter`.
- Added 3 tui-tester E2E tests: `test_card_browser_search_filter`, `test_session_stats_show_in_statistics`, `test_help_overlay_shows_and_dismisses`.
- All 36 E2E tests pass; all Go tests pass; `./scripts/verify.sh` passes with zero errors.

## 2026-04-29 Autonomous Feature Pass 3 - Review Control, Habit Goal, Deck Insights

- Implemented suspended review cards: `x` in Review suspends the current card, suspended cards are filtered from normal/bookmarked due queues, Dashboard/Statistics show suspended counts, and SQLite persists the flag via migration 13.
- Added targeted Go coverage for suspended-card persistence/filtering/statistics and TUI queue refresh.
- Implemented SQLite-backed daily goal settings: migration 14 adds `app_settings`, Statistics reads the persisted daily goal, and Settings supports `+`/`-` adjustments with a floor of 1.
- Fixed fresh-install Statistics aggregation so empty card tables report zero maturity counts instead of surfacing NULL scan errors.
- Added targeted Go coverage for daily goal persistence/defaults/flooring and Settings rendering/updates.
- Implemented deck progress insights: Decks view now shows per-deck reviews today and all-time success rate derived from persisted review history.
- Added targeted Go coverage for deck progress aggregation and Decks rendering.
- Added 3 tui-tester E2E scenarios covering suspended-card persistence/statistics, daily-goal persistence, and deck progress metrics after review.
- Fixed post-grade aggregate refresh so Decks, Dashboard, and Statistics receive fresh persisted counts immediately after a review.
- Final verification passed with `./scripts/verify.sh`, including all Go tests, go vet, smoke launch, and 33 tui-tester E2E tests.

## 2026-04-29 Autonomous Feature Pass 2 - Bookmark Filter, MCQ Review, Leech Detection

- Added bookmarked-only review mode: `B` (shift+b) in Review toggles a filter showing only bookmarked cards, with a banner indicator and updated Dashboard/Statistics counts.
- Added MCQ card review support: MCQ cards now render numbered choices after Space reveal, accept 1-4 key selection, show Correct/Incorrect feedback, and grade normally. Choices persist in SQLite via new migration 12.
- Added leech detection: cards with 3+ consecutive "Again" grades are auto-flagged as "leech", shown with a LEECH indicator during review, and counted in Dashboard/Statistics. Reset on Hard/Good/Easy grades.
- Added SQLite migration 11 for `card_flags.leech` and `lapse_streak` columns; migration 12 for `cards.choices` column to persist MCQ choices.
- Extended `core.Repository` with `DueCardsBookmarked` method; extended `Statistics` with `BookmarkedDue` and `LeechCards` fields.
- Added 9 Go unit tests covering bookmark filter toggle, MCQ rendering/selection/reset, leech detection/persistence/reset, and Statistics updates.
- Added 3 tui-tester E2E tests covering bookmark filter toggle, leech detection across restart, and MCQ choice selection feedback.
- Final verification passed with `./scripts/verify.sh`, including gofmt, all Go tests, go vet, smoke launch, and 30 E2E tests.

- Added SQLite-backed card bookmarks through a new `card_flags` migration, repository contract support, Review `b` key handling, and Dashboard/Statistics bookmark counts.
- Added undo-last-review support through the repository boundary and Review `u` key handling; SQLite now deletes the newest review and restores the previous review state or returns the card to new.
- Extended Statistics with persisted review-history aggregates for Reviews Today, Daily Goal progress, and Current Streak.
- Added unit coverage for bookmark persistence/counting, undo rollback behavior, daily progress/streak aggregation, and TUI bookmark/undo model flows.
- Added 3 tui-tester E2E tests covering bookmark persistence across restart, undo persistence across restart, and daily Statistics progress after review.
- Final verification passed with `./scripts/verify.sh`, including gofmt, all Go tests, go vet, smoke launch, and 27 E2E tests.

## 2026-04-29 End-to-end stabilization recertification and commit

- Re-ran app launch smoke coverage and `go test ./...`; both passed.
- Re-ran the existing tui-tester E2E suite before changes; all 19 tests passed.
- Added 3 new tui-tester E2E recertification tests covering Tab navigation across all primary views, Hard grading with SQLite persistence across restart, and Settings provider persistence across restart.
- Verified the new E2E tests with `tui_tester/venv/bin/python -m pytest e2e_tests/test_recertification.py -q`.
- Final verification passed with `./scripts/verify.sh`, including gofmt, all Go tests, go vet, smoke launch, and 22 E2E tests.

## Statistics View Implementation (2026-04-29)

Implemented a dedicated Statistics view to provide insights into learning progress and card maturity.

### Implementation
- Added `Statistics` struct and `Repository` method to core domain.
- Implemented `Statistics` fetching in SQLite store with aggregate queries for card maturity (New, Young, Mature) and review success rates.
- Integrated `ViewStatistics` into TUI with dedicated rendering and navigation.
- Refactored TUI navigation to ensure statistics are re-loaded when switching to the view.
- Updated all navigation controls (Tab, Arrows, Number keys 1-7) to support the new 7-view structure.

### Testing
- Created `e2e_tests/test_statistics.py` to verify rendering and real-time updates after reviews.
- Updated all existing E2E tests (`test_tui.py`, `test_wasd_navigation.py`, `test_robustness.py`, `test_recertification.py`) to align with the new view indices and shortcut keys.
- Fixed `mockRepo` in TUI tests to implement the new `Statistics` interface method.

### Verification
- All 24 E2E tests pass.
- `./scripts/verify.sh` passes with zero errors or warnings.

## Navigation Bug Fix and Robustness Enhancement (2026-04-29)

Fixed a conflict between WASD navigation and AI draft discarding, and added new robustness tests.

### Bug Fix
- Removed 'd' key from global view switching to resolve conflict with 'discard' in AI Drafts view.
- Ensured 's' and arrow keys remain available for view navigation.
- Verified fix with targeted E2E tests.

### Testing
- Added `e2e_tests/test_robustness.py` with 3 new tests:
  - `test_settings_template_editing_cancel`: Verifies template edit behavior.
  - `test_import_nonexistent_file_shows_error`: Verifies error handling for missing files.
  - `test_review_grade_updates_dashboard_count`: Verifies end-to-end review and dashboard sync.
- Expanded `e2e_tests/test_wasd_navigation.py` to verify function preservation for 'a' and 'd' keys.
- Updated `scripts/verify.sh` to run the full E2E test suite.

### Verification
- All 19 E2E tests pass.
- `./scripts/verify.sh` passes with zero errors.

## WASD Key Navigation Support (2026-04-29)

Added WASD key support to the deutsch-tui application to improve keyboard navigation:

### Implementation
- Added 'w' key to switch to previous view (like left arrow)
- Added 's' key to switch to next view (like right arrow)
- Preserved 'a' and 'd' keys for existing functions
- Modified `internal/tui/model.go` to handle WASD keys
- Ensured no conflicts with existing functionality

### Testing
- Created comprehensive E2E tests in `e2e_tests/test_wasd_navigation.py`
- Verified all existing tests still pass
- Confirmed WASD keys work correctly for view switching
- Confirmed existing 'a'/'d' functions are preserved

### Verification
- All 16 E2E tests pass including new WASD tests
- `./scripts/verify.sh` passes with zero errors
- No regressions in existing functionality

## 2026-04-29 End-to-end stabilization recertification kickoff

- Reopened the end-to-end deutsch-tui stabilization milestone for the current autonomous pass.
- Confirmed `./scripts/tui_smoke.sh` launches successfully.
- Confirmed direct `TUIAgent` dashboard rendering with `cmd/deutsch-tui --data-dir <tmp>`.
- Verified `go test ./...` and the existing 11 tui-tester E2E tests pass before adding fresh coverage.
- Added three fresh tui-tester E2E recertification tests covering Tab navigation across all primary views, Hard grading with SQLite persistence across restart, and Settings provider persistence across restart.
- Verified the expanded tui-tester suite passes with 14 tests.
- Final verification passed with `./scripts/verify.sh`, including formatting, `go test ./...`, `go vet ./...`, smoke launch, and all 14 E2E tests.

## 2026-04-29 End-to-end deutsch-tui stabilization

- Fixed E2E startup synchronization by waiting for the Dashboard text instead of treating a quiet PTY as a ready screen.
- Added tui-tester mouse click support using xterm SGR mouse events.
- Added E2E coverage for all core views, persisted review grading across restart, and mouse-driven Review tab plus grade interaction.
- Verified with `go test ./...`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 6: AI Provider Configuration and Prompt Templates

- Implemented `TemplateProvider` in `internal/ai` to support customizable prompt templates using `{{.Topic}}` substitution.
- Enhanced `Settings` view to allow toggling between "offline" and "template" AI providers.
- Added support for basic template editing (Front, Back, Example) directly in the `Settings` view.
- Updated `Config` struct to persist `AIProvider` and `AITemplates`.
- Updated `main.go` to initialize the AI provider from the user's configuration.
- Fixed outdated E2E shortcut keys and added a new E2E test for settings and template-based drafting.
- Updated `scripts/verify.sh` to include E2E tests for more comprehensive verification.
- Verified with unit tests, E2E tests, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 5: Deeper Deck Browser

- Implemented a dedicated "Decks" view (`ViewDecks`) to list all available decks with total and due card counts.
- Added keyboard-driven navigation (`j`/`k`, `up`/`down`) for browsing the deck list and `Enter` for selection.
- Updated the SQLite repository to efficiently fetch deck statistics (total/due cards) using a JOIN query.
- Integrated the Decks view into the main TUI navigation, including `Tab` rotation and shortcut key `2`.
- Synchronized deck selection between the new Decks view and existing cycling (`[`/`]`) shortcuts.
- Added unit tests for Decks view rendering, navigation, and storage-layer statistics calculation.
- Verified with `./scripts/verify.sh` and new unit tests.

## 2026-04-29 Milestone 4: TUI Anki TSV Import/Export

- Added Import view actions for importing `import.tsv` and exporting the selected deck to `export.tsv` in the app data directory.
- Preserved exported deck IDs during TSV re-import.
- Loaded full deck notes/cards from SQLite so selected decks can be exported through the repository boundary.
- Added unit coverage for TSV deck-column roundtrip, deck export loading, Import view import/export commands, and mock repository behavior.
- Added E2E coverage for importing a TSV deck, switching to it, reviewing its card, and exporting it back to TSV.
- Verified with `go test ./internal/content ./internal/storage/sqlite ./internal/tui ./cmd/deutsch-tui`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 3: AI Drafting And Deck Selection

- Added an offline AI drafting provider that generates validated draft notes and cards without network access.
- Wired AI draft generation, review, discard, and approval into the TUI.
- Added keyboard deck switching with deck-scoped due-card filtering.
- Added unit coverage for offline draft generation, deck switching, draft approval persistence through the repository boundary, and AI error handling.
- Added E2E coverage for AI draft approval and persistence across restart.
- Verified with `go test ./internal/ai ./internal/tui ./cmd/deutsch-tui`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 2: Review Session Flow

- Implemented full FSRS-backed review session flow.
- Updated SQLite schema with `stability` and `difficulty` fields for FSRS.
- Added `GetDeck` and `Decks` methods to `core.Repository` and `sqlite` store.
- Enhanced TUI Dashboard to load and display deck information.
- Implemented card advancement and status feedback in TUI Review view.
- Added mouse support for grading (Again, Hard, Good, Easy hitboxes).
- Fixed space key binding for `bubbletea/v2`.
- Updated unit and E2E tests to cover the complete review flow.
- Verified with `./scripts/verify.sh` and `pytest e2e_tests/test_tui.py`.

- Bootstrapped Go module, Bubble Tea v2 TUI shell, core domain, content import/export, SQLite storage, SRS adapter, AI draft validation, and agent continuity docs.
- Verified with `GOCACHE=/tmp/deutsch-tui-gocache go test ./...` and `GOCACHE=/tmp/deutsch-tui-gocache go vet ./...`.

## 2026-04-29 Stability Foundation

- Added Git ignore rules, `Makefile`, `./scripts/verify.sh`, and `./scripts/tui_smoke.sh`.
- Added config/logging contracts, numbered SQLite migration tracking, architecture boundary tests, Anki fixture tests, dependency policy, migration policy, release checklist, and fixture documentation.
- Initialized the directory as a Git repository.
- Verified with `./scripts/verify.sh`.

## 2026-04-29 Arrow Key Navigation

- Implemented left and right arrow keys (along with Shift+Tab) to navigate between application views (Dashboard, Decks, Review, Import, AI, Settings).
- Restricted deck switching strictly to `[` and `]` to avoid key mapping conflicts.
- Updated the application footer to clarify arrow key support and account for all 6 views.
- Added 3 new tui-tester E2E tests to verify arrow key navigation across tabs, decks, and settings views.
- Verified with full E2E testing suite and `./scripts/verify.sh` to meet the completion criteria.
