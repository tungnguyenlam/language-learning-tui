# Done

## Autonomous Feature Pass 10: Deck Tags and Filtering ✅
- Added SQLite migration for tags column on decks table
- Updated repository for deck tag persistence
- Implemented deck filtering in TUI with text input and tag matching
- Added 3 E2E tests for deck tags and filtering
- All tests passing

## Autonomous Feature Pass 11: APKG Import/Export ✅
- Implemented APKG export: generates valid .apkg files with SQLite database
- Implemented APKG import: reads .apkg files, extracts notes/cards
- Added Shift+I (APKG import) and Shift+X (APKG export) in Import view
- Added unit tests in internal/content/apkg_test.go
- Updated architecture boundary test to allow database/sql in content package
- All Go tests pass, E2E tests pass