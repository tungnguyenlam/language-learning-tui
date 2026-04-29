## 2026-04-29 End-to-end stabilization recertification kickoff

- Reopened the end-to-end deutsch-tui stabilization milestone for the current autonomous pass.
- Confirmed `./scripts/tui_smoke.sh` launches successfully.
- Confirmed direct `TUIAgent` dashboard rendering with `cmd/deutsch-tui --data-dir <tmp>`.
- Verified `go test ./...` and the existing 11 tui-tester E2E tests pass before adding fresh coverage.
- Added three fresh tui-tester E2E tests for compact view rendering, space-key reveal plus Again grading, and mouse tab navigation to Import, AI, and Settings.
- Verified the expanded tui-tester suite passes with 14 tests.
- Final verification passed with `./scripts/verify.sh`, including formatting, `go test ./...`, `go vet ./...`, smoke launch, and all 14 E2E tests.

## 2026-04-29 End-to-end deutsch-tui stabilization

- Fixed E2E startup synchronization by waiting for the Dashboard text instead of treating a quiet PTY as a ready screen.
- Added tui-tester mouse click support using xterm SGR mouse events.
- Added E2E coverage for all core views, persisted review grading across restart, and mouse-driven Review tab plus grade interaction.
- Verified with `go test ./...`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 6: AI Provider Configuration and Prompt Templates

- Implemented `TemplateProvider` in `internal/ai` to support customizable prompt templates using `{{.Topic}}` substitution.
- Enhanced `Settings` view to allow toggling between "offline" and "template" AI providers.
- Added support for basic template editing (Front, Back, Example) directly in the `Settings` view.
- Updated `Config` struct to persist `AIProvider` and `AITemplates`.
- Updated `main.go` to initialize the AI provider from the user's configuration.
- Fixed outdated E2E shortcut keys and added a new E2E test for settings and template-based drafting.
- Updated `scripts/verify.sh` to include E2E tests for more comprehensive verification.
- Verified with unit tests, E2E tests, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 5: Deeper Deck Browser

- Implemented a dedicated "Decks" view (`ViewDecks`) to list all available decks with total and due card counts.
- Added keyboard-driven navigation (`j`/`k`, `up`/`down`) for browsing the deck list and `Enter` for selection.
- Updated the SQLite repository to efficiently fetch deck statistics (total/due cards) using a JOIN query.
- Integrated the Decks view into the main TUI navigation, including `Tab` rotation and shortcut key `2`.
- Synchronized deck selection between the new Decks view and existing cycling (`[`/`]`) shortcuts.
- Added unit tests for Decks view rendering, navigation, and storage-layer statistics calculation.
- Verified with `./scripts/verify.sh` and new unit tests.

## 2026-04-29 Milestone 4: TUI Anki TSV Import/Export

- Added Import view actions for importing `import.tsv` and exporting the selected deck to `export.tsv` in the app data directory.
- Preserved exported deck IDs during TSV re-import.
- Loaded full deck notes/cards from SQLite so selected decks can be exported through the repository boundary.
- Added unit coverage for TSV deck-column roundtrip, deck export loading, Import view import/export commands, and mock repository behavior.
- Added E2E coverage for importing a TSV deck, switching to it, reviewing its card, and exporting it back to TSV.
- Verified with `go test ./internal/content ./internal/storage/sqlite ./internal/tui ./cmd/deutsch-tui`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 3: AI Drafting And Deck Selection

- Added an offline AI drafting provider that generates validated draft notes and cards without network access.
- Wired AI draft generation, review, discard, and approval into the TUI.
- Added keyboard deck switching with deck-scoped due-card filtering.
- Added unit coverage for offline draft generation, deck switching, draft approval persistence through the repository boundary, and AI error handling.
- Added E2E coverage for AI draft approval and persistence across restart.
- Verified with `go test ./internal/ai ./internal/tui ./cmd/deutsch-tui`, `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q`, and `./scripts/verify.sh`.

## 2026-04-29 Milestone 2: Review Session Flow

- Implemented full FSRS-backed review session flow.
- Updated SQLite schema with `stability` and `difficulty` fields for FSRS.
- Added `GetDeck` and `Decks` methods to `core.Repository` and `sqlite` store.
- Enhanced TUI Dashboard to load and display deck information.
- Implemented card advancement and status feedback in TUI Review view.
- Added mouse support for grading (Again, Hard, Good, Easy hitboxes).
- Fixed space key binding for `bubbletea/v2`.
- Updated unit and E2E tests to cover the complete review flow.
- Verified with `./scripts/verify.sh` and `pytest e2e_tests/test_tui.py`.

- Bootstrapped Go module, Bubble Tea v2 TUI shell, core domain, content import/export, SQLite storage, SRS adapter, AI draft validation, and agent continuity docs.
- Verified with `GOCACHE=/tmp/deutsch-tui-gocache go test ./...` and `GOCACHE=/tmp/deutsch-tui-gocache go vet ./...`.

## 2026-04-29 Stability Foundation

- Added Git ignore rules, `Makefile`, `./scripts/verify.sh`, and `./scripts/tui_smoke.sh`.
- Added config/logging contracts, numbered SQLite migration tracking, architecture boundary tests, Anki fixture tests, dependency policy, migration policy, release checklist, and fixture documentation.
- Initialized the directory as a Git repository.
- Verified with `./scripts/verify.sh`.

## 2026-04-29 Arrow Key Navigation

- Implemented left and right arrow keys (along with Shift+Tab) to navigate between application views (Dashboard, Decks, Review, Import, AI, Settings).
- Restricted deck switching strictly to `[` and `]` to avoid key mapping conflicts.
- Updated the application footer to clarify arrow key support and account for all 6 views.
- Added 3 new tui-tester E2E tests to verify arrow key navigation across tabs, decks, and settings views.
- Verified with full E2E testing suite and `./scripts/verify.sh` to meet the completion criteria.
