# Done Backlog

## 2026-05-05 - AI View: Prompt Template Refinement
- **Implemented Multi-Set AI Templates**: Refactored the configuration to support multiple template sets (e.g., 'vocabulary', 'grammar', 'articles', 'conjugation').
- **Template Switching**: Added the ability to cycle through template sets in the AI view using the `[` and `]` keys, with immediate visual feedback.
- **Specialized Exercise Templates**: Added default templates for German articles (MCQ) and verb conjugation (Cloze) to facilitate diverse study material generation.
- **Nested Template Editing**: Updated the Settings view to support editing templates within the currently active set.
- **Robustness**: Updated `TemplateProvider` to use the selected set and refined `onConfigChange` to correctly persist nested map structures.
- **Verification**: Updated Go unit tests and adjusted E2E mouse interaction tests to account for layout shifts. All 93 tests verified passing.

## 2026-05-05 - Browser Management: Unused Tags Cleanup
- **Implemented Unused Tags Cleanup**: Added the 'C' key in the Card Browser to automatically remove tags from the current deck that are no longer associated with any cards.
- **Repository Support**: Added `CleanupTags` to the `core.Repository` interface and implemented it in the SQLite store using efficient tag extraction and consolidation.
- **TUI Integration**: Added a dedicated cleanup handler and updated the Browser view help text and global help overlay.
- **Improved Status Management**: Refined status message handling to prevent action-specific messages (like "Cleaning up...") from being prematurely overwritten by general background status updates.
- **Enhanced Verification**: Added `TestCleanupTags` unit test and `test_tag_cleanup.py` E2E test. Verified all 93 project tests passing.

## 2026-05-05 - Decks View Bulk Operations
- **Implemented Multi-Select**: Users can now select multiple decks in the Decks view using the `x` key.
- **Bulk Deletion**: Added support for deleting multiple selected decks (or the current one) using `Backspace`.
- **Deck Merging**: Implemented deck merging via the `M` key, allowing notes from all selected decks to be moved into the highlighted target deck.
- **Repository Extensions**: Added `DeleteDecks` and `MergeDecks` to the core Repository interface and SQLite implementation, using transactional safety for merging.
- **Enhanced E2E Tests**: Added `e2e_tests/test_bulk_decks.py` and updated existing tests to account for UI changes. All 91 tests verified passing.

## 2026-05-05 - Card Export Filtering & UI Robustness
- **Implemented Export Filtering**: Added deck and tag filters to the Import/Export view, allowing users to export specific subsets of their collection.
- **Improved UI Navigation**: Added support for keyboard (`j/k`, `[`/`]`, `t`) and mouse interaction for export filters.
- **Synchronized Active Deck**: The export deck filter now defaults to the currently active deck for a more intuitive workflow.
- **Fixed Key Interception**: Resolved a conflict where global navigation keys stole input from view-specific handlers, ensuring robust control in all views.
- **Enhanced Test Suite**: Added 2 new E2E tests for export filtering and verified all 89 project tests passing.

## 2026-05-05 - Deck Search & Key Handling Refinement
- **Implemented Robust Deck Search**: Added `/` search to Decks view with dedicated `searchingDecks` state and improved UI feedback.
- **Fixed Enter Key Handling**: Resolved a bug where `\r` and `\n` were not recognized as Enter in multiple views (Import, Settings, Decks, Browser), improving cross-terminal and automated testing compatibility.
- **Improved Key Trapping**: Refactored `updateKey` to more robustly trap and delegate keys to active views during text input, ensuring global navigation shortcuts don't interfere with typing.
- **Enhanced E2E Test Coverage**: Added `e2e_tests/test_new_pass_20.py` with 3 new tests covering deck search, deck selection, and settings adjustment.
- **Verified Full Suite**: Successfully ran all 87 tests via `scripts/verify.sh` with zero failures.

## 2026-05-05 - Browser Management: Tag Management
- **Implemented Tag Management**: Users can now add or remove tags from cards directly in the Card Browser using the `T` key.
- **Bulk Tagging**: Added support for bulk tag updates when multiple cards are selected.
- **Note-Card Sync**: Tag updates are applied at the Note level, ensuring consistency across all cards belonging to the same note (e.g., Forward and Reverse cards).
- **Tag Search**: Updated the Browser search (`/`) to include tags in the search query.
- **TUI Tag Input**: Implemented a dedicated tag input mode in the Browser with pre-filling from the current card's tags.
- **UI Feedback**: Added tag display to the Card Browser list (e.g., `#tag1 #tag2`) and updated help text.
- **Input Guard**: Updated `textInputActive` to prevent accidental view switching or global key triggers while typing tags or searching in the browser.
- **Storage Layer**: Added Migration 17 to add `tags` column to `notes` and `cards` tables. Implemented `SetCardsTags` for efficient bulk updates.
- **Verification**: Added `TestSetCardsTags` unit test and `test_tag_management.py` E2E test.

## 2026-05-05 - Documentation & Agent Instruction Pass
- **Comprehensive Documentation Refresh**: Expanded all package-level READMEs in `internal/` to include key symbols, responsibilities, and architectural boundaries.
- **Agent Rule Consolidation**: Created and updated subtree `AGENTS.md` files for `core`, `app`, `srs`, `tui`, `content`, `storage`, and `ai` to codify package-specific constraints and best practices for AI agents.
- **Root Instruction Indexing**: Enhanced the main `AGENTS.md` with a "Common Tasks" section, providing direct pointers to relevant subtree instructions for tasks like adding card types or modifying the UI.
- **Index Alignment**: Synchronized `docs/agent/index.md` with the new documentation structure, ensuring all architectural guides and package maps are correctly indexed.
- **Feature Documentation**: Integrated details about Cloze deletion, Reverse cards, and the responsive layout system into the core documentation to bridge the gap between recent features and long-term guides.


## 2026-05-04 - Browser Management: Note Type Conversion
- **Implemented Card Kind Conversion**: Users can now toggle cards between `Flashcard` and `MCQ` kinds directly in the Card Browser using the `t` key.
- **Bulk Type Conversion**: Added support for bulk kind conversion when multiple cards are selected.
- **Automated MCQ Choice Generation**: When converting a flashcard to MCQ, the system automatically ensures at least two choices exist by using the answer and a placeholder, maintaining data integrity.
- **Repository Support**: Added `SetCardKind` to the `core.Repository` interface and implemented it in the SQLite store.

## 2026-05-04 - Browser Management: Bulk Actions
- **Implemented Multi-Select in Card Browser**: Users can now select multiple cards using the 'm' key. Selected cards are indicated with `[x]` checkboxes and highlighted text.
- **Added Bulk Actions**: Implemented bulk actions for bookmarking (`b`/`B`), suspending (`x`/`X`), and deleting (`backspace`/`delete`) multiple cards at once.
- **Improved Browser UI**: Added a dynamic footer that appears when cards are selected, showing the count of selected cards and available bulk actions.
- **Added Browser E2E Tests**: Created `e2e_tests/test_browser_management.py` to verify multi-select and bulk action workflows.

## 2026-05-04 - Review Intervals & Layout Refinement
- **Implemented Predicted Review Intervals**: The Review view now displays the predicted "Next Review" interval for each grading option (Again, Hard, Good, Easy) based on the FSRS algorithm. This helps users make informed grading decisions.
- **Refactored Wide Mode Layout**: Implemented a dynamic 2/3 column design in Wide mode. The layout now uses 2 columns (Nav + Content) for widths between 100 and 120, and 3 columns (Nav + Content + Detail) for widths above 120. This prevents text wrapping in the content panel and fixed multiple E2E test regressions.
- **Improved TUI Utility**: Added a `formatDuration` helper for user-friendly time interval display.
- **Dynamic Hitboxes**: Updated grading hitboxes in the Review view to be calculated dynamically based on label widths.

## 2026-05-04 - Multi-Deck Studies
- **Implemented Multi-Deck Study Sessions**: Added an "All Decks" virtual deck to the TUI. When selected, the app now interleaves due cards across all available decks, allowing users to perform comprehensive review sessions without manually switching decks.
- **Updated E2E Tests**: Adjusted test assertions in `test_dashboard_features.py` and `test_tui.py` to expect the new "All Decks" default state.

## 2026-05-04 - Advanced Learning & UX Polish Pass (Continued)

### Advanced Learning
- **Implemented Reverse Flashcards**: Notes marked as `Reverse` now automatically generate both Front->Back and Back->Front cards.
- **Improved Example MCQs**: Example sentences are now scanned for the target word, creating fill-in-the-blank style MCQs when possible.
- **B1 Idioms Deck**: Added a new thematic deck covering common German idioms with a mix of Basic and Reverse cards.

### UX & Bug Fixes
- **Dashboard/Stats Counters**: Added "Total Decks" and "Active Decks" to both the Dashboard and Statistics views for better progress tracking.
- **Cram Mode Key Handling**: Fixed a critical bug where global number keys (1-9) were switching views instead of selecting filters in Cram Mode.
- **Test Suite Alignment**: Updated E2E tests to handle shifted layouts and increased content height on the Dashboard.

### Advanced Learning
- **Implemented Reverse Flashcards**: Notes can now be marked as `Reverse` (or `Basic (and reversed card)`) to automatically generate both Forward (Front -> Back) and Reverse (Back -> Front) cards.
- **TSV Import Support**: Enhanced the TSV parser to recognize `notetype:Reverse` and `notetype:Basic (and reversed card)`.
- **Core Model Update**: Added `Type` field to `Note` struct to support extensible note types.

### UX Polish & Dashboard
- **Enhanced Dashboard Stats**: Added "Total Decks" and "Active Decks" counters to the Collection box.
- **Grade Color Coding**: Added visual feedback in Review view with color-coded grade labels (Again=Red, Hard=Orange, Good=Green, Easy=Cyan).
- **Tab Jumping**: Implemented quick view switching using number keys 1-9 (and 0 for Dashboard).
- **Dashboard Interactivity**: Added mouse hitboxes for all major Dashboard boxes (Review Queue, Collection, Progress, Digest).
- **Simplified Help UI**: Moved help overlay to a simple appended text block for better discoverability and test compatibility.
- **Disabled AltScreen**: Temporarily disabled AltScreen to improve terminal output capture in automated testing environments.

### Bug Fixes & Refinement
- **Fixed shadowed variables**: Resolved build issue in `internal/storage/sqlite/sqlite.go`.
- **Improved key handling**: Unified number key logic and removed redundant traps in `internal/tui/keys.go`.
- **E2E Test Stability**: Reconciled UI changes with existing test expectations to ensure continued CI/CD reliability.
- **Verified SQLite Persistence**: Confirmed that new deck-level statistics are correctly calculated from the database.

## 2026-05-05 - TUI Refactor Completion & Final Verification

- Completed the comprehensive TUI refactor, solidifying all architectural improvements and rendering contracts.
- Final verification: all 83 E2E tests passing, confirming feature parity and integration stability.
- Verified `./scripts/verify.sh`: all Go tests, smoke tests, and full E2E suite passed with no regressions.

## 2026-05-03 - Run 2 TUI Refactor Recovery

- Completed the in-progress TUI file split recovery from `LEFT_OVER_STATE.md`.
- Restored stable top-level rendering contracts, status footer text, tab/nav click targets, Settings/AI view surfaces, reveal animation state, scrollbar hitbox alignment, APKG export, autoplay status feedback, and import mouse hitboxes.
- Added sequential startup loading with `tea.Sequence` to avoid concurrent init database access.
- Removed temporary recovery/debug handoff files after verification.
- Verified `./scripts/verify.sh`: all Go tests, smoke test, and 77 E2E tests passed.

## 2026-05-02 - Review Reveal Sync and Interactive Testing

- Fixed race condition in Review view where grading hints were displayed before hitboxes and answers were ready.
- Refactored `renderReview` to strictly synchronize visual state with interaction state.
- Implemented **drag-to-scroll support** for scrollbars in Statistics, Browser, and Cram views.
- **Enhanced Decks view** with separate counts for New, Due, and Total cards, providing better progress visibility.
- **Added Session Statistics** to the status line, showing accuracy and cards reviewed in the current session.
- **Improved AI view UX** with interactive [Approve] and [Discard] buttons and a detailed draft preview section.
- **Refactored Card Browser** with a modal **Search Mode** (triggered by `/`) and direct card management via **b** (bookmark) and **x** (suspend).
- **Enhanced Settings view** with interactive hitboxes and clickable [+] and [-] buttons for adjusting the daily goal.
- **Enhanced Import view** with a modern button-like layout and interactive hitboxes for path editing and triggering actions.
- **Added interactive Dashboard hitboxes** linking to Review, Browser, and Statistics views.
- **Added interactive Cram filters** allowing mouse-based switching of card sets.
- **Fixed multiple edge-case bugs** in session statistics (undo support), state management (resetting dragging on view/size change), and MCQ interaction state leakage.
- **Implemented robust ID-based deck synchronization** preventing visual and functional desyncs across all collection-refresh actions.
- **Resolved navigation "traps"** by ensuring global keyboard shortcuts (1-9, ?) take precedence over unhandled local view keys.
- **Fixed deck selection bug** and **secured Cram mode cursor**, improving the reliability of filtered collection management.
- **Fixed audio auto-play and error reporting**, ensuring reliable playback feedback for learners.
- **Documented TUI Tester** as a standalone utility with detailed `README.md` and `AGENTS.md` guides, including instructions for publishing as a Gemini Skill or MCP server.
- **Cleaned up TUI Tester directory** by removing redundant internal `e2e_tests` to ensure a clean standalone repository.
- **Added MIT License** to both the main project and the standalone `tui_tester` library.
- Enhanced `tui-tester` with `move_mouse` and `drag_mouse` simulation capabilities.
- Verified all 76 E2E tests passing.
- **Parallelized E2E tests** with `pytest-xdist`, reducing full verification time from ~5.5 minutes to ~45 seconds.
- Verified all Go unit tests passing.
- Updated project notices and index to reflect architectural requirements for reveal synchronization.

## Animation and Appearance Polish

- **Smother reveal animation**: Increased from 5 to 10 steps (20% → 10% per tick) for smoother text reveal
- **Extended reveal timing**: Increased tick interval from 50ms to 60ms for more comfortable pacing (600ms total)
- **Modern block characters**: Replaced `█` with `▌` in reveal animations for a more elegant look
- **Color palette refresh**: Updated all view titles to use accent color (159) for better consistency
- **Session accuracy colors**: Added dynamic color coding (green/yellow/red) based on performance thresholds
- **Dashboard border colors**: Harmonized with accent color scheme for better visual cohesion
- **Browser search active state**: Enhanced with accent color when search is active
- **Streak/maturity indicators**: Changed from 🔥/⭐ to ✨ for a more modern appearance
- **Refined typography**: Updated muted text from 244 to 248, headers from 229 to 231 for better contrast
