# internal/content

Content management, TSV parsing, and Anki interop.

## Responsibilities

- **TSV Parsing**: `parseTSV` handles the custom flashcard format.
- **Anki Interop**: `.apkg` (SQLite-in-ZIP) import *and* export — see below.
- **Starter Content**: Go-defined decks plus embedded TSV files under `testdata/german-decks/`.
- **Cloze Deletion**: Parsing and rendering logic for `{{c1::...}}` style cards.

`StandardDecks()` is the production seed list (Go decks + TSV notes keyed by `#deck:`). `AllDecks()` / `DeckByID()` also load embedded TSV files under filename-derived IDs so lookups like `b2_urban_mobility` stay stable.

### Adding New Content

Drop a TSV file into `testdata/german-decks/` or add a Go deck function and register it in `StandardDecks()`.

## Anki Interop (`.apkg`)

`apkg_import.go`, `apkg_export.go`, and `apkg_schema.go` implement a real
round trip with Anki, not just a lossy read.

### Import

- Reads `collection.anki21b` (zstd-compressed), `collection.anki21`, or
  `collection.anki2`, preferring the newest present.
- Handles both the legacy schema-11 layout (note types in the `col.models`
  JSON blob) and the modern one (separate `notetypes`/`fields`/`templates`
  tables).
- **Field roles are derived from the templates, not from position.** Anki note
  types are free-form: field 0 is not reliably the front and field 1 is not
  reliably the back. `templateFieldOrds` reads the `{{Field}}` references out
  of a template's `qfmt`/`afmt` to find which fields the question and answer
  actually use, falling back to the first two non-empty fields. A real shared
  deck that leaves field 1 empty imported with a blank back before this existed.
- Strips HTML to plain text, converting `<br>`/`</div>`/`</p>` to newlines and
  lifting `[sound:...]` references out into the note's audio field.
- Cloze ordinals come from `{{c1::...}}` markers; the template ordinal is the
  cloze number minus one.

### Export

`ExportAnkiAPKG` (and the `…WithDeckNames` / `…ToFile` variants) writes a
schema-11 package Anki will actually open. The invariants Anki checks, all
asserted by `TestExportProducesAnkiValidCollection`:

- `col.ver` is 11.
- Every note's `mid` names a note type present in `col.models` — three are
  declared (`Deutsch-TUI Basic`, `… (and reversed card)`, `… Cloze`).
- Every card's `did` names a deck present in `col.decks`, and `col.conf` sets
  `curDeck`.
- All ids are positive milliseconds, allocated densely.
- The media index is a file literally named `media` (no extension).
- No genesis/placeholder note is emitted.

Fields are joined with `\x1f`, `csum` is the int of the first 8 hex chars of
the SHA1 of the stripped first field, and tags are stored space-padded — all
as Anki expects.

## Flashcard Format (TSV)

The project uses a tab-separated format for easy editing:
- Columns: `Front`, `Back`, `Extra`, `Tags`, `NoteType` (optional).
- `NoteType`: `Basic`, `Reverse`, `Cloze`, `MCQ`.
- `Extra`: Optional explanation shown after reveal.

## Key Symbols

- `StandardDecks`: production seed list (Go decks + TSV by `#deck:`)
- `AllDecks` / `DeckByID`: combined Go + filename-keyed embedded TSV lookup
- `ImportAnkiTSV`: TSV import
- `ExportAnkiAPKG`: Anki package export
