# Review Reveal Synchronization

Status: active
Scope: internal/tui/model.go, e2e_tests/
Related: `renderReview`, `RevealRevealing`, `RevealRevealed`

## Why It Matters

Previously, the "Grade:" hints were displayed during the `RevealRevealing` animation state. However, the hitboxes for mouse interaction and the full answer text were only added/visible in the `RevealRevealed` state. This caused race conditions in E2E tests that synchronized on the presence of "Grade:" text but then immediately tried to assert the answer or click a grade button before the animation finished.

## Required Behavior

Future agents must ensure that interactive elements (hitboxes) and final answer text are strictly synchronized with the visual state. Do not show grading hints or register grading hitboxes until the `revealState` is `RevealRevealed`.

## Revisit When

Condition: If the animation system is refactored to support interactive elements during the reveal, or if a more robust synchronization primitive is added to the TUI tester.
