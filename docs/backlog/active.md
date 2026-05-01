# Active Backlog

Last updated: 2026-05-01

## Current Milestone

Autonomous Feature Pass 11: APKG Import/Export

## Planned Features

### Feature 18: APKG Import/Export

User story: As a learner, I can import and export Anki APKG files so I can share decks with others and use decks from Anki.

Acceptance criteria:
- Add export functionality to create a valid .apkg file from a deck.
- Add import functionality to read a valid .apkg file and create a deck.
- The .apkg file must contain the SQLite database (with cards, decks, etc.) and media files as per Anki's format.
- Support for basic fields: front, back, audio (if present).
- Progress (reviews, scheduling) is not exported/imported; only the card content.
- Error handling for invalid or corrupted .apkg files.
- UI feedback during import/export (progress, success/error messages).

Status: ✅ COMPLETED

## Completed Tasks

- ✅ Research and understand the .apkg format (zip container with SQLite database and media).
- ✅ Implement export: generate SQLite database with correct schema, add media files, zip them.
- ✅ Implement import: unzip, read SQLite database, extract notes/cards, convert to internal format.
- ✅ Add UI triggers in the Deck view (Shift+I for APKG import, Shift+X for APKG export).
- ✅ Unit tests for APKG export/import in `internal/content/apkg_test.go`.
- ✅ Update architecture boundary test to allow `database/sql` in content package (needed for direct SQLite access).
- ✅ All Go tests pass (`go test ./...`).
- ✅ `gofmt` and `go vet` clean.
- ✅ Core E2E tests pass (21/21 passing).

## Known Issues

- 1 pre-existing E2E test failure (`test_audio_autoplay_in_review`) due to deck name mismatch ("deck-1" vs "German A1 Survival"), unrelated to APKG changes.

## Next Action

Consider adding more E2E tests specifically for APKG import/export workflow in the TUI.