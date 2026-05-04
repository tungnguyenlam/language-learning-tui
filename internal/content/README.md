# internal/content

Content management, TSV parsing, and Anki interop.

## Responsibilities

- **TSV Parsing**: `parseTSV` handles the custom flashcard format.
- **Anki Interop**: `apkg.go` handles `.apkg` (SQLite/ZIP) import.
- **Starter Content**: Embedded TSV files for A1-B2 German levels.
- **Cloze Deletion**: Parsing and rendering logic for `{{c1::...}}` style cards.

## Flashcard Format (TSV)

The project uses a tab-separated format for easy editing:
- Columns: `Front`, `Back`, `Tags`, `NoteType` (optional).
- `NoteType`: `Basic`, `Reverse`, `Cloze`, `MCQ`.

## Key Symbols

- `Loader`: Interface for loading decks from various sources.
- `ImportTSV`: Logic for importing cards from a file or string.
- `GermanDecks`: Embedded collection of starter vocabulary.
