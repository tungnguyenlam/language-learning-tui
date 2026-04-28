# deutsch-tui

A local-first terminal language learning app for German, built with Go and Bubble Tea.

The app is designed for long-running agent-assisted development. Start future sessions by reading [docs/agent/start-here.md](docs/agent/start-here.md), then the active backlog.

## Current MVP

- German A1 starter deck.
- Flashcard and MCQ review primitives.
- SQLite progress storage.
- FSRS-backed scheduling adapter.
- Anki-friendly TSV import/export.
- Responsive Bubble Tea shell with keyboard and mouse input.

## Verify

```sh
./scripts/verify.sh
```

Run:

```sh
go run ./cmd/deutsch-tui
```
