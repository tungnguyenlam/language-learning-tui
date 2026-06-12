# TTY Exploratory Test Results

Date: 2026-06-12
App: deutsch-tui (Go Bubble Tea TUI)
Test method: `tui-tester` autonomous exploration

## Critical Bugs

### BUG C1: Focus Mode Completely Broken Rendering
- **View:** Review → Press `f` for Focus Mode
- **Severity:** Critical
- **Description:** Enabling focus mode (`f` key) in the Review view produces a completely garbled render. Borders are shattered (`─╮` at wrong positions, `│` scattered randomly), content (card text, hints) appears at wrong coordinates, and ANSI escape code fragments (`?m│`) leak into the visible output. The entire layout is unreadable.
- **Reproduction:** Start app → Go to Review (key 3) → Press `f`

### BUG C2: Dictionary Spotlight Overlay Rendering Corruption
- **View:** Any view → Press `=` for Spotlight Dictionary
- **Severity:** Critical
- **Description:** The Spotlight Dictionary overlay (`=`) has severe rendering corruption. The overlay panel borders leak ANSI escape codes (`;208m`, `;213m`, `;6m`, `38;5;255m`), underlying view content bleeds through the overlay (sidebar labels, dashboard text), and the right side of the overlay panel shows garbled text from the view beneath it. Occurs on both Dashboard and Review views.
- **Reproduction:** Press `=` from any view

### BUG C3: Dashboard Rendering Corruption After View Navigation
- **View:** Dashboard (after navigating through other views)
- **Severity:** Critical
- **Description:** After navigating between views (especially Practice Hub, Cram, then back to Dashboard via Tab), the Dashboard shows severe rendering corruption: ANSI escape codes in the Quick Actions box (`A8;5;234m`), broken Unicode characters in borders (`�──────`), overlapping box boundaries, and text from previous views bleeding through (`�ctivity` from "Recent Activity"). The entire layout becomes unreadable.
- **Reproduction:** Dashboard → Cram (9) → Practice Hub (Tab) → Dashboard (Tab). Also reproducible by opening then closing Dictionary overlay from Dashboard.

## Major Bugs

### BUG M1: Practice Hub View ANSI Escape Code Leak
- **View:** Practice Hub
- **Severity:** Major
- **Description:** The Practice Hub view shows `[1;16HPRACTICE │ wide` in the card list area, which is a raw ANSI cursor positioning escape sequence leaking into the terminal output. The header bar also incorrectly shows `CRAM` instead of `PRACTICE HUB` after navigating from Cram view.
- **Reproduction:** Dashboard → Cram (9) → Press 5 to load cards → Tab to Practice Hub

### BUG M2: Left Sidebar Content Bleeds Into Main Content
- **Views:** Multiple (Cram, Practice Hub, AI Drafts)
- **Severity:** Major
- **Description:** When navigating between views that have differently-sized content, text from the left sidebar (e.g., "52 cards loaded", "Press enter to s") or other content bleeds through into the main content panel. The rendering doesn't properly clear the previous view's content.
- **Reproduction:** Navigate between Cram and Practice Hub views

### BUG M3: Cram Mode Filter List Shows Duplicate "3: Leeches"
- **View:** Cram (key 9)
- **Severity:** Major
- **Description:** The Cram mode filter list displays "3: Leeches" twice instead of showing "3: Leeches" once and "4: All flagged" in its correct position. The numbering skips from 3 to 4 with two identical Leeches entries.
- **Reproduction:** Go to Cram (key 9) → Select "5: All cards" → Observe filter list

### BUG M4: tui-tester Daemon Hangs on `+` Key in Settings
- **View:** Settings (key 7)
- **Severity:** Major (affects automated testing)
- **Description:** When navigating to Settings, moving to the Daily Goal row, and pressing `+` to increase the value, the `tui-tester` daemon becomes unresponsive. The `+` character may be causing an issue in the key processing pipeline. This was observed multiple times and caused the daemon to require a kill and restart.
- **Reproduction:** Settings → Navigate to "Daily Goal: 10 cards" → Press `+`

## Minor Bugs

### BUG m1: Header Truncation After Focus Mode
- **View:** Review
- **Severity:** Minor
- **Description:** After entering and exiting focus mode (`f` key), the header line becomes misaligned: `   deutsch-tui` has trailing spaces before it and the outer border on the right doesn't close properly (missing `│` on line 1). The `DEUTSCH-TUI │ REVIEW │ wide` header format is disrupted.
- **Reproduction:** Review → Press `f` → Press `f` again to exit

### BUG m2: Bottom-Line Rendering Artifact
- **Views:** Multiple (observed on Decks, Review, Dashboard)
- **Severity:** Minor
- **Description:** A fragment of text "rases for daily German." (truncated from "phrases for daily German") appears persistently on the last line of the terminal output, outside of any border/layout. This text appears to be from deck description content that leaks outside the rendering area.
- **Reproduction:** Navigate to Decks view and observe the bottom-most line of output

### BUG m3: `q` Key in Cram Mode Quits Entire App
- **View:** Cram (key 9)
- **Severity:** Minor (expected but potentially confusing UX)
- **Description:** In Cram mode, pressing `q` quits the entire application rather than just exiting Cram mode back to the Dashboard. The help text says "q to quit" but user expectation in a sub-mode might be that `q` exits the sub-mode. Consider using `Esc` for exiting Cram and reserving `q` for app-wide quit.

## Views Verified (No Critical Issues)

1. **Dashboard** - Clean render on fresh start, all sections display correctly
2. **Review** - Card reveal, grading (a/h/g/e), dictionary lookups, bookmarks all work
3. **Decks** - List view, search/filter work correctly
4. **Settings** - Navigation with j/k works, settings display correctly
5. **Statistics** - Maturity distribution, review stats display correctly
6. **Browser** - Card preview with dictionary details render well
7. **Import/Export** - Form fields and action buttons display correctly
8. **AI Drafts** - Template selection and topic input display correctly (render bugs noted above)
9. **Help Overlay** - Opens/closes correctly, comprehensive shortcut listing

## Root Cause Analysis

The rendering bugs share a common theme: **ANSI escape code leakage** and **incomplete screen clearing** during view transitions. Specifically:

1. The `applyOverlay` function (noted as previously fixed for UTF-8 in the backlog) appears to still have issues with border rendering that uses ANSI color codes - these codes are being included in the visible output rather than being interpreted by the terminal.
2. View transitions don't properly clear the previous view's content before rendering the new view, causing text from the previous view to "bleed through" in the new view.
3. The Focus Mode rendering appears to use a different layout engine that has fundamental issues with box-drawing character placement.

## Recommendations

1. **Audit `applyOverlay` and all border-rendering code** for ANSI escape sequences that are being placed inside visible string content rather than being applied as styling.
2. **Add a full screen clear** (`\033[2J` or Bubble Tea's equivalent) before rendering Focus Mode and Spotlight Dictionary overlays.
3. **Investigate the Cram filter duplication** as a separate data/logic bug in the filter list construction.
4. **Add `+` key handling** to the Settings daily goal adjustment logic and verify it doesn't panic or hang.
5. **Verify `tui-tester` compatibility** with the `+` key to confirm the hanging issue is in the app, not the test tool.