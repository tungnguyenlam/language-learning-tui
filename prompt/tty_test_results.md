# TTY Exploratory Testing Results - 2026-06-14

## Summary
Exploratory testing of `deutsch-tui` revealed a generally stable core UI with some significant performance/reliability issues in the `tui-tester` integration or the app's responsiveness to rapid input.

## Findings

### BUG-001: Tester/App Hangs During Rapid Navigation or Search [High]
- **Description:** The `tui-tester` daemon or the `deutsch-tui` process frequently hangs (timeouts) when performing searches or switching views quickly. This was observed during Browser search and Dictionary lookup.
- **Severity:** High (Impacts automated testing and potentially user experience if the app blocks on UI thread).
- **Reproduction:**
    1. Start app.
    2. Go to Browser (shortcut '8').
    3. Press '/' to search.
    4. Type 'rot' and press Enter.
    5. Occasionally hangs here.

### BUG-002: Practice Mode Shortcut Confusion [Minor]
- **Description:** Shortcut '0' is listed for "Practice" in some views, but it navigates to "Cram Mode" in some contexts or shows "PRACTICE HUB" in others. The sidebar says "Practice [0]" but the navigation behavior was slightly inconsistent during testing.
- **Severity:** Minor.
- **Evidence:** Captured screen showing "DEUTSCH-TUI │ CRAM" when '0' was pressed.

### BUG-003: UI Rendering Artifacts [Minor]
- **Description:** Some screen captures showed overlapping text or partial renders (e.g., `deutsch-tui DASHBOARD │ wide──────────────────────────────────────────────────────────────────────────────╮` appearing inside the Practice Hub).
- **Severity:** Minor (Visual).
- **Evidence:** Multiple `tui-tester observe` outputs showed shifted borders and mixed view content.

### BUG-004: Browser Multi-Select Instruction Error [Trivial]
- **Description:** The footer says "Enter select" but 'm' is used for selection in the list. Enter actually reveals history or edits.
- **Severity:** Trivial.

## Conclusion
The core learning and review loop works well. The app handles initial data import gracefully. The primary focus for "improve.md" should be stabilizing the UI responsiveness and fixing the inconsistent shortcuts.
