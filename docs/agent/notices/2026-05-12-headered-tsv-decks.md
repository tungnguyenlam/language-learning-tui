# Headered TSV Decks

Status: active
Scope: `internal/content/anki.go`, `internal/content/testdata/german-decks/*.tsv`
Related: `ImportAnkiTSV`, `CardsForNote`, `TestEmbeddedDecksDoNotLeakHeadersOrLiteralFieldsIntoCards`

## Why It Matters

Several embedded starter decks use a headered TSV shape without an explicit `id` column, and some rows include a deck label before `notetype` even when the header omits `deck`. If those headers are parsed as data, review cards can show English prompts with `Literal:` explanations as answers.

## Required Behavior

Keep `ImportAnkiTSV` header-aware. Header rows such as `front	back	extra	tags	notetype` must be skipped, generated note IDs for no-id headered decks must continue to use the trimmed front text so reseeding repairs old malformed imports, `#deck:` metadata must win over shorter row deck labels, and `Literal:`/`Literally:` fields must remain `Extra`, never card answers.

## Revisit When

Revisit if embedded decks are migrated to a single canonical TSV schema with explicit `id`, `deck`, `notetype`, and `audio` columns.
