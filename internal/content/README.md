# internal/content

Content management, TSV parsing, and Anki interop.

## Responsibilities

- **TSV Parsing**: `parseTSV` handles the custom flashcard format.
- **Anki Interop**: `.apkg` (SQLite-in-ZIP) import *and* export — see below.
- **Starter Content**: Embedded TSV files for A1-B2 German levels.
- **Cloze Deletion**: Parsing and rendering logic for `{{c1::...}}` style cards.
- **Content Registry**: Scalable system for loading decks from multiple sources.

## Scalable Content System

The content package uses a scalable Registry pattern that supports multiple content sources:

```go
// Default registry auto-loads from embedded TSV files
decks, err := content.AllDecks()

// Get specific deck by ID
deck, err := content.DeckByID("a1-survival")

// List available deck IDs
ids, err := content.DeckIDs()
```

### Adding New Content

To add new content, simply drop a TSV file into `testdata/german-decks/`:

```
testdata/german-decks/
├── a1-essential.tsv
├── a1-food-drink.tsv
├── a2-grammar-essentials.tsv
├── b1-idioms.tsv
└── ... (auto-discovered)
```

The Registry automatically:
1. Scans for `.tsv` files in the directory
2. Parses each file as a deck
3. Merges notes into deck objects
4. Makes them available via `AllDecks()`

### Content Sources

| Source | Priority | Description |
|--------|----------|-------------|
| EmbeddedSource | 10 | Loads TSV files from `testdata/german-decks/` |
| GoSource | 20 | Loads Go-defined decks (StarterDeck) |

Higher priority = loaded first. Duplicate deck IDs are deduplicated.

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

- `Registry`: Content source registry with pluggable backends
- `ContentSource`: Interface for custom content loaders
- `ImportTSV`: Logic for importing cards from a file or string
- `GermanDecks`: Embedded collection of starter vocabulary (legacy)
