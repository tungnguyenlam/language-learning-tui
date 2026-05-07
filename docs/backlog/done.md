# Done Backlog

## 2026-05-08: End-to-End TUI Correctness Batch A
- **Bug Fixes**:
  - Clamped Review cursor access across rendering, grading, audio playback, bookmark, suspend, history, and MCQ paths to prevent stale-cursor panics after queue changes.
  - Clamped AI template-set access in AI, Settings, and template-provider switching to prevent panics when configuration indices drift.
  - Added empty import/export path guards for TSV and APKG actions, returning clear status messages instead of attempting invalid file operations.
- **UI/UX Polish**:
  - Added a bordered empty-state panel for Review with concrete next actions.
  - Added Import/Export path warnings, a stronger AI topic placeholder, an AI prompt-quality hint, and clearer Settings daily-goal/help text.
- **Verification**:
  - `go test ./internal/tui ./internal/content` passes.

## 2026-05-08: End-to-End TUI Correctness Batch B
- **Content Addition**:
  - Added `b1-apartment-housing.tsv` with 39 B1 housing, rental, viewing, and move-in cards.
  - Registered the housing deck in `EmbeddedDeckPaths` so deck discovery helpers expose it.
- **E2E Coverage**:
  - Added 8 tui-tester tests for Dashboard progress/digest, Review reveal/grading persistence, Import filters/actions, empty import path status, AI guidance, Settings help, Deck navigation, and all core view tabs.
- **UI/UX Polish**:
  - Added stable AI suggested-topic text and shortened Settings help text to avoid wide-layout wrapping.
- **Verification**:
  - `go test ./internal/content ./internal/tui` passes.
  - `DEUTSCH_TUI_BIN=/tmp/deutsch-tui-e2e pytest e2e_tests/test_end_to_end_core_views.py -q` passes (8/8).

## 2026-05-08: Full-Suite Regression Fixes
- **Content Registry Fix**:
  - Normalized embedded TSV note/card `DeckID` values to the registry deck ID before card generation, fixing foreign-key failures when registry decks are written to SQLite.
- **Startup Compatibility**:
  - Preserved the app's starter-deck default seed so existing E2E flows keep the expected 52-card queue while expanded decks remain available through the content registry.
- **E2E Maintenance**:
  - Updated the agent workflow test to use the current AI navigation key and starter-deck runtime default.
- **Verification**:
  - `go test ./...` passes.
  - Targeted E2E `test_agent_workflow.py test_end_to_end_core_views.py` passes (10/10).

## 2026-05-08: End-to-End TUI Milestone Verified
- **Scope Completed**:
  - App launch, Dashboard, Review, Import, AI, Settings, Decks, Browser, Cram, persistence, and core navigation are covered by the full suite.
  - Completed more than 8 distinct improvements across bug fixes, UI/UX, AI UX, content, tests, and developer verification hygiene.
- **Final Verification**:
  - `./scripts/verify.sh` passes cleanly.
  - Result: all Go tests pass, smoke passes, and 116/116 E2E tests pass with no warning output.

## 2026-05-07: Per-Deck Study Limits
- **Storage & Core**:
  - Extended `Deck` model with `NewCardsPerDay` and `ReviewLimitPerDay`.
  - Added SQLite migration to persist per-deck limits.
  - Implemented `SetDeckLimits` in the repository layer.
- **TUI Enhancements**:
  - Added "Limit Edit" mode to the Decks view (press `L` on a selected deck).
  - Implemented interactive adjustment with `+`/`-` and `Tab` to toggle between New/Review limits.
  - Optimized rendering to only show limits during editing, maintaining compatibility with existing E2E tests.
- **Verification**:
  - All 105 E2E tests pass.
  - Verified storage migrations and persistence via unit tests.

## 2026-05-07: UI Polish & UX Enhancements
- **Dashboard Enhancements**:
  - Implemented a visual progress bar for daily study goals.
  - Added a "Next Card" preview to give users a glimpse of what's coming up.
- **Review View Polish**:
  - Introduced a dedicated "Card" border around flashcard content for better focus.
  - Improved layout separation between deck headers and card content.
  - Refined MCQ feedback with color-coded Correct/Incorrect labels.
- **AI View Improvements**:
  - Added a dedicated, boxed Preview area for AI-generated drafts.
  - Improved draft list alignment and truncation for better readability.
- **Statistics View**:
  - Unified progress bar rendering using a new shared helper function.
- **Test Stability**:
  - Updated E2E tests to match the new UI coordinates and formatting.
  - Verified all 105 tests pass successfully.

## 2026-05-07: Multi-Agent Bug Hunting & Stability
- **Storage Layer Fixes**:
  - Fixed resource leak by adding `defer rows.Close()` in `notesForDeck`.
  - Removed redundant global `PRAGMA foreign_keys = ON` execution calls (now strictly relying on connection DSN).
  - Patched iteration error hiding by adding `rows.Err()` checks across `calculateStreak`, `UndoLastReview`, and `CleanupTags`.
- **TUI Panic Fixes**:
  - Implemented bounds checking for `dueCards` in `selectMCQChoice` to prevent panics on empty queues.
  - Implemented bounds checking for `draftCursor` in `approveDraft` to prevent panics during fast AI draft approvals.
  - Fixed slice bounds panic in `renderCramAt` when applying filters that empty the cram queue.
- **TUI State & UI Logic Fixes**:
  - Fixed "leaky modal" by successfully clearing `m.confirmingDelete` state in `updateView` to prevent overlaps when switching views.
  - Fixed dynamic width/height recalculation using explicit layout dimensions in browser, AI, and deck views instead of root terminal bounds.
  - Resolved `N+1` database query inefficiency in `bulkBrowserToggleKind`.
  - Added robust error handling and returning for all bulk repository operations (`SetCardBookmark`, `SetCardSuspended`, `SetCardKind`).
- **Verification**: All 105 E2E tests and Go tests pass successfully without regressions.

## 2026-05-07: Content Library Expansion
- **Grammar Essentials**: Expanded from 14 to 60 cards (cases, tenses, conjunctions, modal verbs)
- **German Idioms**: Expanded from 10 to 43 cards (common sayings and proverbs)
- **New Decks Added**:
  - A1 Health & Body (97 cards): body parts, health vocabulary, common ailments
  - A2 Shopping & Services (69 cards): retail, banking, postal services
  - A2 Transport & Directions (146 cards): vehicles, directions, travel
  - B1 Emotions & Feelings (110 cards): emotional vocabulary, expressions
  - B1 False Friends (~150 cards): German words that look like English but mean different things
  - B1 Phrasal Verbs (~120 cards): trennbare and untrennbare verbs with examples
- **Total Content**: Expanded from ~600 to 1366 cards across 14 decks (CEFR A1-B2)
- **Verification**: All 105 E2E tests pass, Go tests pass

## 2026-05-06: Critical Stability & Integrity Fixes
- **Deck Management Robustness**:
  - Implemented stable `deleteIDs` in the TUI model to prevent race conditions and ID shifting during deletion confirmation.
  - Added protection to the virtual 'All Decks' item, preventing accidental deletion attempts.
  - Stabilized `deckCursor` in `syncDecks` to prevent the selection from snapping back to the active deck during background reloads.
  - Hardened `main.go` to only insert the starter deck if the database is truly empty, preventing overwrites of user-modified starter decks.
- **Data Integrity & Keybindings**:
  - Explicitly enabled `PRAGMA foreign_keys = ON` immediately after opening the SQLite database to guarantee cascading deletes.
  - Fixed MCQ key collisions in the Review view by surgically trapping 1-4 keys only when a card is present, restoring global navigation for other views.
- **Enhanced Testing**:
  - Added new robust E2E tests for MCQ navigation and UI sanity.
  - Verified the entire project suite (105 tests) with 100% pass rate.

## 2026-05-06: Bug Hunting & Stability Polish
- **TUI Bug Fixes**:
  - Clamped `browserCursor` in `browserCardsMsg` to prevent panics when card list is updated.
  - Properly handled errors returned from repository in bulk browser operations (bookmark, suspend, delete).
  - Fixed leaky modal state by discarding all mouse events when `confirmingDelete` is true.
  - Cleaned up dangling duplicate code block at the end of `model.go`.
- **Data Integrity & SRS Fixes**:
  - Enforced SQLite `PRAGMA foreign_keys=ON` in DSN to ensure correct connection pool behavior.
  - Wrapped `SetCardKind` and `CleanupTags` in transactions for atomicity.
  - Resolved timezone inconsistency in `ReviewsPerDay` by using Go's `time.Local` instead of SQLite `DATE()`.
  - Added clamping to `ScheduledDays` in FSRS state conversions to prevent `time.Duration` integer overflows.
  - Added `idx_reviews_reviewed_at` index in database migrations to improve statistics query performance.

## 2026-05-06: Robustness & Data Integrity (STABILIZED)
- **Bug Fix**: Fixed `UndoLastReview` to correctly recalculate `lapse_streak` and `leech` flags by inspecting remaining history.
- **Bug Fix**: Refactored `ImportAnkiTSV` to robustly handle quoted multiline fields and avoid false-positive comment detection.
- **Improved Testing**: Added regression tests for review undo flag restoration and multiline TSV imports.
- **Security**: Enabled `LazyQuotes` in CSV reader to handle malformed but common Anki export patterns safely.
- **Code Quality**: Cleaned up metadata extraction in importer to use `bufio.Scanner` instead of manual splitting.
- **E2E Stability**: Verified all 102 E2E tests pass with zero regressions.

## 2026-05-05: TUI Stabilisation & UX Polish
- **E2E Performance**: Built binary once in `verify.sh`, reducing test suite time and resource contention.
- **Safety**: Added confirmation dialogs for all destructive actions (card/deck delete).
- **UI Integrity**: Aligned all TUI hitboxes with Lip Gloss metrics for pixel-perfect mouse interaction.
- **Data Integrity**: Implemented transactional bulk updates in the SQLite layer.
- **Test Coverage**: Added 4 new E2E tests and verified 100% pass rate (97/97).

## 2026-05-03: Template Editing & Configuration
- [x] Added revert-on-cancel for template editing.
- [x] Stored original values in `originalTemplateValue`.
- [x] Fixed template set switching logic.
