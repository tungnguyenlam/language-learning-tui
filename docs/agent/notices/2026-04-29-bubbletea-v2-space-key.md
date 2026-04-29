# Bubble Tea v2 Space Key

Status: active
Scope: internal/tui, bubbletea
Related: internal/tui/model.go

## Why It Matters

In `charm.land/bubbletea/v2`, the `KeyPressMsg` for the space key (`Code: ' '`) has a string representation of `"space"` rather than `" "`.

## Required Behavior

When matching key messages in Bubble Tea v2, use `"space"` instead of `" "` if you want to support the space bar.

## Revisit When

If we ever downgrade to Bubble Tea v1 or if v2 changes its string representation of the space key.
