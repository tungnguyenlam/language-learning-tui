# Done Backlog

## 2026-05-10 15:30 +07 (Autonomous Improvement Pass)

Completed a major improvement pass focusing on Content, UI/UX, and Logic.

### Content Expansion
- Added **German Idioms & Proverbs** deck with 15 common expressions.
- Added **German Slang & Youth Language** deck with 12 modern informal terms.
- Expanded AI Drafting templates with a dedicated **Grammar Explanation** template set.

### UI/UX Refinements
- **Dashboard:** Added a compact, informative header showing "WILLKOMMEN!", daily review progress, and the current date.
- **Statistics:** Added a **"Success Rate per Deck"** visualization with color-coded progress bars.
- **Review View:** Made the **Context/Extra** field significantly more prominent during the reveal phase.
- **AI View:** Fixed the **loading spinner animation** by correctly batching the tick command with the drafting process.

### Logic & Testing
- **Browser:** Improved `Esc` key behavior—it now clears the search/tag filter when in searching mode (while `Enter` preserves it).
- **Testing:** Added a new E2E test file `e2e_tests/test_browser_search.py` covering real-time filtering and filter clearing.
- **Robustness:** Updated existing unit tests to handle batched commands in the TUI model.

### Verification
- All 165 E2E tests passed.
- All Go unit tests passed.
- `./scripts/verify.sh` executed successfully.

## 2026-05-10 14:30 +07 (Content, UI Polish & SRS Robustness)

Completed a multi-area improvement pass covering German content expansion, UI refinements (scrollbars, grammar examples), statistics visualization, and SRS robustness.

### Content Expansion
- Added **German Confusable Words** deck (e.g., *seit* vs *seid*, *wahr* vs *war*).
- Expanded **Grammar Tips** with real-world examples for each tip.

### UI/UX Refinements
- **Decks View:** Implemented a visual **scrollbar** with mouse-clickable hitboxes for long lists.
- **Statistics:** Added a **Maturity Distribution** chart (New vs. Young vs. Mature cards).
- **Polish:** Improved layouts in Settings and Browser views for better spacing on various screen sizes.

### Logic & SRS
- **Scheduler:** Added comprehensive unit tests for the SRS scheduler logic.
- **Robustness:** Fixed a regression where deck selection didn't update the review counts correctly.

### Verification
- All 163 E2E tests passed.
- All Go unit tests passed.
- `./scripts/verify.sh` executed successfully.
