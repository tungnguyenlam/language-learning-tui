# TTY Exploratory Test Results

**Date:** 2026-06-12
**Binary:** `deutsch-tui-bin` (built from `./cmd/deutsch-tui`)
**Method:** `tui-tester` CLI with XML-RPC backend, 110x40 terminal

## Test Coverage Summary

All views were explored: Dashboard, Review, Decks, Statistics, Import, AI Drafts, Settings, Browser, Cram, Practice Hub, Session Summary. Interactive features tested: dictionary spotlight search, review card grading, browser search, cram session launch, focus mode, help overlay, deck search, mouse clicks, and rapid view switching.

---

## Bugs Found

### BUG-001 (CRITICAL): Dictionary Spotlight Overlay Keystroke Capture Broken

**Severity:** Critical — Blocks core dictionary functionality
**Reproduction:**
1. From any view, press `=` to open the Spotlight Dictionary overlay
2. Type any characters (e.g., "Haus")
3. Press `<Enter>`

**Observed:** Typed characters are sent to the **underlying view**, not the Spotlight search bar. For example, from AI Drafts, typing "Haus" populates the AI topic field ("Topic: Hausfrau - hausfrau_"). The spotlight search bar remains empty.

**Expected:** Typed characters should be captured by the Spotlight overlay's search input field, and `<Enter>` should execute a dictionary search.

**Evidence:**
```
 DEUTSCH-TUI │ AI │ wide
  AI Drafts
  Deck: All Decks
  Template: vocabulary (use [ / ])
  ┌────────────────────────────────────────────────────────────────────┐
  │ Topic: Hausfrau - hausfrau_                                        │   ← "Haus" typed HERE instead of Spotlight
  └────────────────────────────────────────────────────────────────────┘
```

**Reproduction steps (tui-tester):**
```bash
tui-tester act "="        # Open spotlight
tui-tester act "Haus"     # Type search - goes to underlying view
tui-tester act "<Enter>"  # Triggers underlying view action, not search
```

---

### BUG-002 (MAJOR): ANSI Escape Sequences Leak into Visible UI

**Severity:** Major — Affects status line and content readability
**Reproduction:** Navigate to any view, observe the status line and content areas.

**Observed:** Raw ANSI escape sequences appear in multiple locations:
- Status line: `status: 52 cards found ;5;240m│` (should be colored separator)
- Daily Digest: `│;159\n;4mA1 colors & shapes` (fragment from AI suggestions)
- Quick Actions box: `[5A` appears inline (cursor movement escape)
- Statistics bars: `█�░░░░░ 0  ✗` (truncated/malformed escape sequence)

**Expected:** All ANSI escape sequences should be consumed by the terminal renderer and never appear as visible text.

**Evidence (Dashboard):**
```
 status: 52 cards found ;5;240m│
```
```
                    │  │   M:0 Y:0 N:52  💪 Ready for your     │   │
                    │  ╰───────────────────────────────────────╯ │ daily German practice?                │;159
;4mA1 colors & shapes  │   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    │ │   Next: blau                          │   │
```

**Evidence (Statistics):**
```
│  Success Rate: 0.0%                     █  │
│  Retention: 0.0% (Mature/Total: 0/0)    █�░░░░░ 0  ✗
```

**Probable cause:** The `applyOverlay` function or ANSI parsing in the TUI renderer (`internal/tui/model.go`) is not handling certain multi-character escape sequences correctly. Recent work on `applyOverlay` (see active.md 2026-06-12) may have introduced regressions.

---

### BUG-003 (MINOR): Help Overlay Text Wrapping Artifacts

**Severity:** Minor — Cosmetic only
**Reproduction:** Press `?` from any view to open the Help overlay.

**Observed:** The "Cram" and "Practice" sections in the help overlay have awkward word breaks:
```
│    q/Ctrl+c Quit                   b / B    Bookmark / Unmark    Practice
│                                    u / r    Undo / History     select
│  Dashboard/Decks:                  f / i    Focus / Info         Import
│    [ ]      Prev/next deck         p / d    Play audio / Dict    exp
```

Lines like "Practice\nselect", "Import\nexp" show words split across lines. The "Statistics" shortcut appears as just "Stats".

**Expected:** Help text should wrap cleanly with proper word boundaries.

---

### BUG-004 (MINOR): Dashboard Daily Digest Text Clipping

**Severity:** Minor — Text readability
**Reproduction:** View Dashboard at 110x40 terminal size.

**Observed:** The Daily Digest box has text that bleeds beyond its border. The example sentence text "Ready for your daily German practice?" wraps awkwardly and the box border breaks:
```
│ daily German practice?                │
  ╰───────────────────────────────────────╯
```

Note the left border of the closing edge is at a different indentation than the right, suggesting the box width calculation is off or text padding doesn't match.

---

### BUG-005 (MINOR): Browser Status Line Content Bleed

**Severity:** Minor — Status text readability
**Reproduction:** Navigate to Browser view (key 8).

**Observed:** The status line at the bottom shows truncated help content instead of a proper status message:
```
lect all, p to play audio, t to toggle
```
(This is a fragment of "Select all, p to play audio, t to toggle" from the help overlay.)

**Expected:** The status line should show a clean, complete status message relevant to the current view.

---

### BUG-006 (MINOR): `/` Key Conflict Between Dashboard and Other Views

**Severity:** Minor — UX inconsistency
**Reproduction:**
1. From Dashboard, press `/` → Full Dictionary tab opens (correct)
2. From AI Drafts, press `/` → Edits AI topic instead of opening Dictionary
3. From Decks, press `/` → Opens deck search filter

**Observed:** The `/` key has context-dependent behavior that may confuse users. It does not consistently open Dictionary across views.

**Expected:** Either `/` should always open Dictionary from all views, or a different shortcut should be used for Dictionary (the Spotlight overlay with `=` already provides this).

**Note:** This is likely by design (Dashboard `/` = Dictionary, Decks `/` = search, AI Drafts `/` = edit topic) but the inconsistency could be a UX papercut.

---

### BUG-007 (MINOR): tui-tester Stability Detection Incompatible with Timer Views

**Severity:** Minor — Testing tool limitation
**Reproduction:** Use `tui-tester act` while in Review view.

**Observed:** `tui-tester act` calls `wait_until_stable()` which never completes in timer-based views (Review, Cram active) because the on-screen timer (`⏱ 00:05`) updates the screen content every second.

**Expected:** `wait_until_stable` should either handle timer-based updates gracefully or the tui_tester should provide a non-waiting alternative.

---

## Views Verified Working

The following views rendered and responded correctly to navigation:

| View | Status | Notes |
|------|--------|-------|
| Dashboard | OK | Layout intact, all boxes shown |
| Review (idle) | OK | Card displayed, reveal works |
| Review (revealed) | OK | Grading options shown, grade works |
| Review Focus Mode | OK | `f` toggles correctly |
| Decks | OK | List renders, search works |
| Settings | OK | All items scrollable |
| Statistics | OK | Distribution charts render |
| Import/Export | OK | Paths displayed correctly |
| Browser | OK | Card list + preview shown |
| Cram (idle) | OK | Filter options displayed |
| Cram (active) | OK | 52 cards loaded with All filter |
| Practice Hub | OK | All 9 trainers listed |
| AI Drafts | OK | Suggestions displayed, topic editing works |
| Help Overlay | OK | `?` toggles correctly |
| Spotlight Overlay | OK | `=` opens/closes, result count shown |

---

## Stress Test: Rapid View Switching

Ten rapid view switches (keys 1-9 cycled twice) completed without crashes. The app remained responsive with no state corruption visible in the header.

---

## Summary

- **Total bugs:** 7 (1 Critical, 1 Major, 5 Minor)
- **Critical:** Dictionary Spotlight keystroke capture is non-functional — typed characters leak to the underlying view
- **Major:** ANSI escape sequences leak into visible text across multiple views (status line, dashboard, statistics)
- **Minor:** Help overlay wrapping, dashboard text clipping, browser status truncation, `/` key inconsistency, tui_tester timer incompatibility
- **No crashes** during navigation or stress testing
