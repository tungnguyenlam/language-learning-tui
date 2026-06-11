# Spotlight Dictionary Overlay

Status: active
Scope: internal/tui, UX direction
Related: internal/tui/keys.go, internal/tui/render_dictionary.go, internal/tui/model.go (dictionaryOverlayActive)

## Why It Matters

The Dictionary was redesigned from a tab-based view to a **Spotlight-like overlay** accessible from anywhere in the app via the `=` key. This follows the macOS Spotlight search UX pattern: a floating search panel that appears on top of the current view without navigating away.

The legacy full Dictionary view (`ViewDictionary`) still exists and is accessible via `/` from the Dashboard (or `d` from Review to look up the current card). However, the future direction is to make the overlay the primary access method and potentially deprecate the full view.

## Required Behavior

- The `=` key opens the spotlight overlay from **any view** (when no text input is active).
- The overlay renders on top of the current view using `applyOverlay` — it does NOT change `m.activeView`.
- `dictionaryOverlayActive` gates key routing: when active, all keys are intercepted by `updateDictionaryOverlayKey` first.
- `textInputActive()` returns `true` when the overlay is active, preventing number-key view switching.
- Pressing `=` or `Esc` closes the overlay and returns to the underlying view.
- Dictionary was **removed from the tab/arrow/WASD navigation cycle** (`nextViewCmd`/`previousViewCmd`).
- Dictionary was **removed from the nav sidebar** and **tab bar** (`renderNav`/`renderTabs`).
- The full `ViewDictionary` tab still exists for `/` and `d` shortcuts.

## Revisit When

- When the full `ViewDictionary` tab is deprecated entirely.
- When the overlay is expanded to support additional features (e.g., Quick Add from overlay, detail view within overlay).
