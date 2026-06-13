# TTY Exploratory Test Results - 2026-06-13

## Summary
The application is generally stable and feature-rich, but several UI/UX papercuts and one potential critical hang were identified during exploratory testing.

## Identified Issues

### [BUG-001] Dictionary Search Input Captures Navigation Keys
- **Severity**: Minor (UX Papercut)
- **Description**: In the Dictionary view, the search input field captures `j` and `k` keys, which are also intended for navigating the results list. This makes it impossible to scroll through results using standard TUI navigation keys while the search is active.
- **Reproduction**:
  1. Open the app.
  2. Press `/` to open the Dictionary.
  3. Type a word (e.g., "Haus").
  4. Try to press `j` or `k` to scroll the results.
- **Expected**: `j` and `k` should scroll the results if the search has results, or there should be a clear way to switch focus (e.g., Tab or Down arrow). Currently, `j` and `k` just append characters to the search input.

### [BUG-002] Cram Mode Overlapping UI (Ghost Text)
- **Severity**: Minor (UI Glitch)
- **Description**: When entering Cram review mode, parts of the filter selection screen (e.g., "Click a filter to load cards") occasionally remain visible or overlap with the card review UI.
- **Reproduction**:
  1. Go to Cram view (9).
  2. Select a filter that has cards (e.g., "5: All cards").
  3. Press Enter to start cramming.
  4. Observe the background of the review card.
- **Evidence**: `tui-tester observe` showed "Click a filter to load cards" text behind the card border.

### [BUG-003] Text Clipping in Cram Review
- **Severity**: Minor (UI Glitch)
- **Description**: Usage examples and other long text in the Cram review card are clipped instead of wrapped, making them unreadable.
- **Reproduction**:
  1. Start a Cram session.
  2. Reveal a card with a long example sentence.
  3. Observe that the text is cut off at the card boundary.

### [BUG-004] Practice Hub / AI Drafts Navigation Hang
- **Severity**: Major (Stability)
- **Description**: Navigating to the Practice Hub (0) or AI Drafts (6) frequently causes the application to hang or become unresponsive to TTY input. This happened consistently during automated testing.
- **Reproduction**:
  1. Open the app.
  2. Press `0` or `6`.
  3. The app stops responding to further keypresses or observation requests.
- **Potential Root Cause**: Likely related to heavy initialization logic or deadlocks in the async loaders (e.g., `loadPracticeItems`) when combined with the PTY driver's read loop.

## Next Steps
1. Fix UX focus logic in Dictionary view.
2. Clean up render logic in Cram mode to prevent overlapping text.
3. Implement text wrapping in review cards.
4. Investigate and resolve the hang in Practice/AI views by optimizing loaders or fixing potential deadlocks.
