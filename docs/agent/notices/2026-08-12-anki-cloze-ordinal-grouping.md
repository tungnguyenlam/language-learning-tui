# Anki Cloze Ordinal Grouping

Status: active
Scope: Anki TSV cloze import and review rendering
Related: `internal/content/anki.go`, `internal/tui/render_review.go`

## Why It Matters

Anki creates one card per distinct cloze ordinal, not one card per marker. Repeating an ordinal
means the same card has multiple blanks, and ordinal order—not source-text order—defines card order.

## Required Behavior

Group equal ordinals, sort the generated cards numerically, retain grouped answers in source order,
and render each answer into its matching prompt placeholder. Keep the ordinal in generated card IDs
so imported cards remain stable when source markers appear out of order.

## Revisit When

Revisit if cloze hints, nested markup, or `.apkg` cloze conversion is unified with TSV parsing.
