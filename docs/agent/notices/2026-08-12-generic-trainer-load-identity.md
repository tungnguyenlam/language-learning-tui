# Generic Trainer Load Identity

Status: active
Scope: `internal/tui` generic practice trainers
Related: `trainer.go`, `handlers.go`, `model.go`

## Why It Matters

Generic practice trainers load their item pools asynchronously. A learner can leave and re-enter a
trainer before an older loader returns, so a response without request identity can replace the
current session's items.

## Required Behavior

Increment the trainer state's `loadID` for every explicit entry/load, capture that ID in the
returned command, and apply `trainerItemsMsg` only when both its ID and active practice subview
match the current state. New generic trainer configurations must preserve this guard.

## Revisit When

Revisit if generic trainer loading is moved to a shared async service or trainer messages gain a
new transport shape.
