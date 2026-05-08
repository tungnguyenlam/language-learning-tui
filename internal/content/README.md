# internal/content

Content management, TSV parsing, and Anki interop.

## Responsibilities

- **TSV Parsing**: `parseTSV` handles the custom flashcard format.
- **Anki Interop**: `apkg.go` handles `.apkg` (SQLite/ZIP) import.
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
