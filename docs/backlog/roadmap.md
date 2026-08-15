# Roadmap

## Milestone 1: Bootstrap

- Runnable Bubble Tea shell.
- Core learning domain.
- SQLite progress storage.
- FSRS scheduling adapter.
- Anki-friendly TSV import/export.
- A1 starter content.
- Agent continuity system.

## Milestone 2: Study UX

- Full review session flow.
- Better deck browser.
- Import/export screens.
- TUI smoke test script.

## Milestone 3: AI Drafting

- Provider configuration.
- Prompt templates.
- Draft review and approval workflow.

## Milestone 4: Hybrid Dictionary (Complete)

Completed 2026-08-05 and merged to `main`.

- Transform the app into a unified dictionary-flashcard hybrid.
- **Storage Layer:** FTS5 SQLite table and parser for offline `dict.cc` text file ingestion.
- **Core API:** `DictionaryEntry` domain struct and `SearchDictionary` repository interface.
- **TUI Integration:** A dedicated `ViewDictionary` plus the `=` Spotlight overlay with real-time type-ahead search.
- **Flashcard Loop:** Seamlessly convert dictionary search results into new flashcards via the Drafting flow.
- **Onboarding:** Settings UI and documentation to guide users on importing the `dict.cc` dataset.

## Later

Shipped ahead of schedule and no longer pending: `.apkg` import/export (`internal/content/apkg_*.go`),
audio pronunciation via TTS (`internal/audio`), and the expanded German trainer content
(`internal/tui/trainer_content.go`).

Still open:

- Cloud or multi-device sync (local progress backup/restore is on Import/Export).
- Local LLM provider for offline AI drafting.
- Continued German content expansion (B2/C1 decks).
