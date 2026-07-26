# Trainer Input vs Global Single-Letter Shortcuts

Status: active
Scope: `internal/tui` key routing, practice trainers
Related: `internal/tui/keys.go` (`trainerInputActive`, `textInputActive`, `updatePracticeKey`), `internal/tui/trainer.go`

## Why It Matters

The text-input practice trainers are deliberately **not** part of
`textInputActive()`. That predicate also suppresses Tab / arrow / `w`-`s` view
switching, which trainers must keep. The side effect is that step 1 of
`Model.Update` ("global critical keys") sees trainer keystrokes first, so any
single-letter global shortcut is swallowed as a command instead of being typed
into the answer. `q` used to quit the app in the middle of an exercise whenever
the answer contained one (*Qualität*, *Quelle*).

`Model.trainerInputActive()` is the narrow escape hatch: true only while a
generic trainer has items loaded and is waiting for an answer (not revealed).
Guarded keys call `updateActiveViewKey` first and fall through to their global
meaning when the trainer does not consume them.

## Required Behavior

- When adding a **new single-letter global shortcut** in step 1 of
  `Model.Update`, guard it with `m.trainerInputActive()` the same way `q` and
  `?` are, or it will be unreachable as a typed character in every trainer.
- Do **not** "fix" this by adding trainers to `textInputActive()`: that breaks
  Tab / arrow navigation out of a trainer, because `updateTrainerKey` does not
  handle those keys and the fallthrough is gated on `!textInputActive()`.
- Chord shortcuts (`ctrl+…`) need no guard — `singlePrintableInput` rejects
  them, so they fall through on their own.

## Revisit When

The trainers move onto the `screen` interface (`internal/tui/screen.go`) or key
routing stops resolving global shortcuts before view-specific handlers.
