# Arrow Key Navigation Fix for Review View

Status: active
Scope: internal/tui
Related: internal/tui/model.go

## Why It Matters

In Bubble Tea v2 applications, arrow keys should navigate within views that support vertical navigation. However, the Review view in deutsch-tui was not properly handling arrow keys due to incorrect key routing logic.

## Issue Description

The Review view was being handled directly in the main `updateKey` function but only for specific action keys (space, enter, grading keys, etc.). The arrow keys ("up", "down") were not included in this handling, causing them to be ignored in the Review view.

Unlike other views that are handled through the `updateActiveViewKey` dispatcher, the Review view had special handling in the main function that bypassed proper arrow key processing.

## Fix Applied

Modified the `updateKey` function in `internal/tui/model.go` to explicitly handle "up" and "down" arrow keys for the Review view, moving the cursor between cards as expected.

## Required Behavior

Arrow keys now properly navigate between cards in the Review view:
- Up arrow (or 'k') moves to the previous card
- Down arrow (or 'j') moves to the next card
- Navigation respects bounds (won't go beyond first or last card)

## Testing

Added `TestReviewArrowNavigation` unit test to verify correct arrow key behavior in the Review view.

## Revisit When

If the key handling architecture is refactored or if additional views need similar fixes.