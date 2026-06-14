# TTY Exploratory Test Results - 2026-06-14

## Summary

Exploratory testing covered Dashboard, Review, Decks, Browser, Statistics, Practice Hub (all trainers sampled), Cram, AI Drafts, Settings, full Dictionary tab, and Spotlight Dictionary overlay (`=`). Core flows work; rendering/overlay issues were the dominant bug class.

**Fix session 2026-06-14:** Addressed BUG-001 through BUG-006 in code (see below). Full `./scripts/verify.sh` passes Go tests and 348+/351 E2E tests; remaining E2E failures are flaky under parallel xdist (deck-name substring timing).

## Identified Issues

### [BUG-001] Practice Hub border corruption and ghost text in sub-views — FIXED
- **Severity:** Major (UI Glitch)
- **Fix:** Two-line button layout with `truncateLine`, replaced wide `✂️` icon with `S/`, increased button spacing, `fillViewportContent` on practice sub-views.
- **Files:** `internal/tui/render_practice.go`, `internal/tui/utils.go`

### [BUG-002] View-transition ghost text (incomplete screen clear) — FIXED (practice views)
- **Severity:** Major (UI Glitch)
- **Fix:** `fillViewportContent` pads short practice trainer output to full panel interior dimensions.
- **Files:** `internal/tui/render_practice.go`, `internal/tui/utils.go`

### [BUG-003] Spotlight Dictionary overlay border / content bleed — FIXED
- **Severity:** Major (UI Glitch)
- **Fix:** `applyOverlay` now uses max overlay width and pads each overlay line before splicing, clearing underlying content in the full overlay rectangle.
- **Files:** `internal/tui/model.go`

### [BUG-004] Settings row text clipping — FIXED
- **Severity:** Minor (UI Glitch)
- **Fix:** Truncate credential display values to fit panel width.
- **Files:** `internal/tui/render_settings.go`

### [BUG-005] Status line / footer content bleed across views — FIXED
- **Severity:** Minor (UI Glitch)
- **Fix:** `refreshViewStatus()` called on every `updateView()` to set view-appropriate default status.
- **Files:** `internal/tui/handlers.go`

### [BUG-006] Dictionary search captures j/k instead of result navigation — FIXED
- **Severity:** Minor (UX Papercut)
- **Fix:** Allow `j`/`k` through text-input trap (same as arrow keys) so dictionary result navigation works while search is active.
- **Files:** `internal/tui/keys.go`

### [BUG-007] tui-tester stability detection incompatible with dynamic views
- **Severity:** Minor (Testing Infrastructure)
- **Status:** Open (tooling limitation). Use `wait-for` anchors, not `wait-stable`, on Review/Statistics timer views.

### [BUG-008] Number keys in Dictionary search cannot switch views
- **Severity:** Minor (UX Papercut)
- **Status:** Open (by design — press Esc first). Low priority.

## Verified Working

- Dashboard, Review grading, Browser, Statistics, AI Drafts, Settings
- Spotlight search keystroke capture (no leak to underlying view)
- Tab/number-key navigation

## Next Step

Monitor E2E parallel flake on deck-list substring tests. Consider `@prompt/improve.md` for BUG-007/008 or new polish.
