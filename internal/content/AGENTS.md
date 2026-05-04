# Content Agent Rules

- **Type Safety**: When adding new `NoteType` values, ensure they are handled in both the parser (`anki.go`) and the TUI (`render_review.go`).
- **Validation**: Reject or fix malformed TSV lines (e.g., missing columns) gracefully.
- **Embedded Files**: If you update the starter decks, run `go generate` if necessary (though current setup uses `embed` directly).
- **Cloze Syntax**: Maintain compatibility with Anki's `{{c1::text}}` format.
