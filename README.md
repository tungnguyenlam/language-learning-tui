# deutsch-tui

A modern, local-first terminal language learning application for German, built with Go and Bubble Tea. It features an advanced Spaced Repetition System (SRS), full mouse support, and AI-powered content generation.

The app is designed for agent-assisted development. Start future sessions by reading [AGENTS.md](AGENTS.md).

## Core Features

- **Spaced Repetition System (SRS):** Powered by the FSRS scheduling algorithm with SQLite persistence.
- **Rich Interactive UI:** A polished Bubble Tea shell with:
    - **Full Mouse Support:** Clickable tabs, interactive buttons, and **drag-to-scroll** scrollbars.
    - **Responsive Layouts:** Dynamic "Wide", "Medium", and "Compact" modes that adapt to your terminal size.
    - **Live Feedback:** Real-time session statistics (Accuracy, reviewed count) and study streaks.
- **Advanced Collection Management:**
    - **Card Browser:** Modal search and quick management actions (Bookmark, Suspend).
    - **Deck Management:** Detailed progress metrics for each deck (New, Due, Total counts).
    - **Cram Mode:** Flexible study sessions for bookmarked, suspended, or leech cards.
- **Powerful Interop & Generation:**
    - **AI Drafting:** Generate high-quality German flashcards using configurable AI providers.
    - **Anki Compatibility:** Support for TSV and `.apkg` import/export.
- **Starter Content:** Includes 600+ vocabulary items across A1-B2 levels.

## Quick Start

### Run the app
```sh
go run ./cmd/deutsch-tui
```

### Run with a custom data directory
```sh
go run ./cmd/deutsch-tui -data-dir ./my-data
```

## Development & Verification

The project maintains a high standard of quality through a comprehensive test suite.

### Full Verification
Runs Go unit tests, code formatting checks, and the full E2E suite:
```sh
./scripts/verify.sh
```

### Testing Infrastructure
- **Go Unit Tests:** 9 test suites covering core logic and storage.
- **E2E Integration:** 76 parallelized tests using the `tui-tester` utility to verify visual state and interaction.

## Architecture

- **Modular Monolith:** Clean separation between `core` logic, `storage` (SQLite), `ai` adapters, and the `tui` layer.
- **Local-First:** All data, progress, and configuration are stored locally in a single directory.
- **Agent Optimized:** Structured for continuous improvement by AI agents with indexed documentation and notices.
