# Trainer Input vs Global Single-Letter Shortcuts

Status: active
Scope: `internal/tui` key routing, practice trainers, Review typing, Cram
Related: `internal/tui/keys.go` (`practiceBlocksGlobalShortcut`, `trainerInputActive`, `textInputActive`, `updatePracticeKey`), `internal/tui/trainer.go`

## Why It Matters

The text-input practice trainers are deliberately **not** part of
`textInputActive()`. That predicate also suppresses Tab / arrow / `w`-`s` view
switching, which trainers must keep. The side effect is that step 1 of
`Model.Update` ("global critical keys") sees trainer keystrokes first, so any
single-letter global shortcut is swallowed as a command instead of being typed
into the answer. `q` used to quit the app in the middle of an exercise whenever
the answer contained one (*Qualität*, *Quelle*).

`Model.practiceBlocksGlobalShortcut()` is the escape hatch used by `q`, `?`,
and `=`: true while a generic trainer has items (typing **or** post-reveal
"press any key"), and while the Gender trainer is in its revealed advance
step. `trainerInputActive()` remains the narrower "typing only" check.

Review typing mode is a separate trap: `q` must reach `updateReviewKey`, while
Cram still uses `q` to exit an active session.

## Required Behavior

- When adding a **new single-letter global shortcut** in step 1 of
  `Model.Update`, guard it with `m.practiceBlocksGlobalShortcut()` (or
  `trainerInputActive()` if only the typing phase matters), or it will steal
  keys from trainers / the reveal-advance step.
- Do **not** "fix" this by adding trainers to `textInputActive()`: that breaks
  Tab / arrow navigation out of a trainer, because `updateTrainerKey` does not
  handle those keys and the fallthrough is gated on `!textInputActive()`.
- Chord shortcuts (`ctrl+…`) need no guard — `singlePrintableInput` rejects
  them, so they fall through on their own.
- Inside a Practice sub-trainer (not Hub), do not let `updateNumberKey` jump
  to global views — Gender/typed answers own those digits.

## Revisit When

The trainers move onto the `screen` interface (`internal/tui/screen.go`) or key
routing stops resolving global shortcuts before view-specific handlers.
