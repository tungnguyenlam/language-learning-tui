# Settings Row Indices

Status: active
Scope: Settings view, E2E tests
Related: `internal/tui/screen_settings.go`, `internal/tui/render_settings.go`, `internal/tui/actions.go`

## Why It Matters

Settings cursor positions are named constants. E2E tests that press `j` a fixed number of times to reach Daily Goal or Auto-play will miss the row if those constants change.

## Required Behavior

Use the constants in `screen_settings.go`. As of 2026-08-22:

- Daily Goal is index 4 (four `j` presses from the AI Provider row)
- Front Template is index 1 (one `j`)
- Auto-play is index 5
- OpenAI Key is index 7
- Reveal Speed is index 15

Dictionary no longer has a Settings row; lookup is the `=` overlay and `/` full tab.

## Revisit When

A Settings row is added or removed above Daily Goal.
