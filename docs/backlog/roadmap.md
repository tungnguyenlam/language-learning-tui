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

## Later

- `.apkg` import/export.
- Audio/media support.
- More German content.
- Sync or backup workflow.

## Milestone 4: Hybrid Dictionary (Next Focus)
- Transform the app into a unified dictionary-flashcard hybrid.
- **Storage Layer:** FTS5 SQLite table and parser for offline `dict.cc` text file ingestion.
- **Core API:** `DictionaryEntry` domain struct and `SearchDictionary` repository interface.
- **TUI Integration:** A dedicated `ViewDictionary` with real-time, zero-latency type-ahead search.
- **Flashcard Loop:** Seamlessly convert dictionary search results into new flashcards via the Drafting flow.
- **Onboarding:** Settings UI and documentation to guide users on importing the `dict.cc` dataset.
