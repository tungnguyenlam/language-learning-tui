# Cram Session Consumes Unhandled Keys

Status: active
Scope: `internal/tui` key routing, Cram
Related: `internal/tui/keys.go` (active-Cram trap), `internal/tui/screen_cram.go`

## Why It Matters

An active Cram session used to dispatch unhandled keys to global navigation.
`s`/`w` switched views and `[`/`]` changed decks while `cramActive` stayed
true, so the learner could leave a live session without `q`.

## Required Behavior

While `activeView == ViewCram && cramActive`, consume every key except
`tab` / `shift+tab` / `ctrl+c` / `q`, even if Cram does not handle it.
`q` still exits the session (then the global `q` handler runs).

## Revisit When

Cram key handling moves behind a stricter modal overlay, or global navigation
is gated on `!cramActive` at the dispatcher.
