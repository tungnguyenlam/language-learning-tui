# Scrollbar Hitbox & Panel Layout

Status: active
Scope: `internal/tui/model.go` active panel scrollbars, mouse hitboxes, and panel padding
Related: `contentLayoutForStyle`, `renderActiveView`, `renderStatisticsAt`, `renderBrowserAt`, `renderCramAt`

## Why It Matters

Lip Gloss panel borders and padding shift rendered content relative to the panel origin. Scrollbar hitboxes must use the same frame-derived content coordinates as rendering, otherwise mouse clicks land one or more cells away from the visible track.
Furthermore, the active view content panel padding in `renderActiveView` must pad lines exactly to `layout.Width` and `layout.Height` rather than using hardcoded offsets like `width - 2`, otherwise trailing space wrapping occurs, causing massive scrollbar misalignment and gaps.

## Required Behavior

1. When adding hitboxes inside the active panel, derive the content origin with `contentLayoutForStyle` or pass the existing `viewportLayout` through the render path. Do not reintroduce per-view offsets such as `x+2` or `y+1` for bordered panel content.
2. In `renderActiveView`, always pad line contents and lines slice to `layout.Width` and `layout.Height` respectively to ensure total conformity to style frames without wrapping.

## Revisit When

Revisit if the active view stops using `panelStyle`, Lip Gloss changes style frame sizing semantics, or the TUI adopts a dedicated viewport/list component for scrollable views.
