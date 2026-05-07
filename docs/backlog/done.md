# Done Backlog

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
