# Active Backlog

Last updated: 2026-06-14

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Exact Next Action

Monitor parallel E2E flake on deck-list tests.

## Top Issues

1. **MINOR - BUG-012 [NEW]:** Dictionary result highlighting fails for certain multi-byte character combinations at the boundary of truncation.
2. **MINOR - BUG-013 [NEW]:** Mouse click on dashboard sparkline does not navigate to statistics.

## Completed Work

- [x] **BUG-009 Practice Hub navigation [NEW]:** Added '0' to global navigation, fixed Cram mode blocking number keys.
- [x] **BUG-010 Selection key consistency [NEW]:** Standardized 'm' for multi-select in Decks and Browser; fixed help documentation.
- [x] **applyOverlay ANSI-Aware Rendering Fix:** Rewrote `applyOverlay` (`internal/tui/model.go:1537`) to properly handle ANSI escape sequences in both base and overlay content. Added `spliceVisual`, `visualPrefix`, `findAnsiEnd`, and `runeVisualWidth` helper functions that map visual column positions to rune offsets while preserving ANSI styling codes. This fixes critical rendering corruption in Dictionary Spotlight overlay, Focus Mode, and Dashboard after view navigation. Added comprehensive unit tests `TestSpliceVisual` and `TestApplyOverlay`.
- [x] **Cram Invisible Text Removal:** Removed problematic `\x1b[8mcramRevealed\x1b[0m` invisible text sequence from `render_cram.go:116` that was causing rendering artifacts and potential overlay corruption. The invisible rendition mode was interfering with proper screen clearing during view transitions.
- [x] **Dashboard Dictionary Shortcuts:** Added `Shift+G`, `Shift+V`, and `Shift+W` shortcuts to the Dashboard for instant dictionary lookup of the daily Grammar Tip, Verb of the Day, and Word of the Day. Updated Dashboard labels to `[g/G]`, `[v/V]`, and `[w/W]` for better discoverability.
- [x] **Non-disruptive Dictionary Lookups:** Updated the `d` key in both Review and Browser views to use the Spotlight Dictionary overlay instead of switching to the full Dictionary tab. This allows for fast, context-preserving lookups without breaking the user's flow.
- [x] **A1 Practical Phrases Deck:** Added a new 49-note embedded Go deck `A1PracticalLifeDeck` covering essential phrases for restaurants, shopping, directions, social interactions, hotels, and emergencies. Added unit test coverage in `a1_practical_life_test.go`.
- [x] **Spotlight Open with Query:** Enhanced the Spotlight Dictionary to support opening with a pre-filled query and immediate search execution, powering the new Dashboard and Review shortcuts.
- [x] **Persistent Dictionary Search History:** Implemented persistent search history for Offline Dictionary searches, storing recent searches in the SQLite `app_settings` table using key `'dict_search_history'`. The history is retrieved asynchronously on startup and saved on additions and clearing (`ctrl+x`), maintaining state across application sessions. Added unit coverage `TestDictionaryPersistentSearchHistory`.
- [x] **Spotlight Dictionary Overlay Enhancements:** Improved the Spotlight Dictionary overlay (`=`) with scrollbar support in its detail panel (both in narrow and wide layouts), persistent status/toast messages at the bottom, and colored status/shortcut indicators. Removed `ctrl+a` overlay dismissal, allowing users to add dictionary cards to their collection without closing the overlay search interface.
- [x] **Practical Life Content Decks:** Added three new embedded TSV decks with 40 notes each: A2 Banking Errands, A2 School & Childcare, and B1 Insurance Claims. Added registry tests to verify they load, generate cards, and include key real-life phrases.
- [x] **B1 Household Maintenance Deck:** Added a 40-note embedded TSV deck covering tools, small repairs, plumbing/electrical issues, and repair-service vocabulary. Added registry coverage to verify the deck loads and generates cards.
- [x] **Dictionary Empty Search Feedback:** Updated dictionary search result handling to distinguish cleared searches from real zero-result searches, including the searched query in the status line when no matches are found.
- [x] **Spotlight Dictionary Result Count:** Added result-count feedback to the Spotlight dictionary overlay title, matching the full Dictionary view's count/capped-count behavior. Added unit coverage for the rendered overlay title.
