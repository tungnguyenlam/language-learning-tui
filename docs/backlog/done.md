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

## Autonomous Feature Pass 12: UI Refinement and Robust Import/Export ✅
- Fixed import path suffix bug allowing any .tsv/.apkg file to be imported.
- Refined Dashboard and Import/Export views with better styling and structured layouts.
- Implemented 'editing mode' for text fields to resolve key conflicts.
- Protected global navigation keys during text input.
- Added a robust APKG export-import cycle E2E test.
- Updated all existing E2E tests to match the new UI and workflow.
- All 62 E2E tests passing.

## Autonomous Feature Pass 14: UI Polish and AI Feedback ✅
- Refined Dashboard with a professional grouped layout using Lip Gloss borders.
- Implemented an animated spinner for AI drafting to provide better visual feedback.
- Improved Card Browser with a structured title and search bar styling.
- Polished Review and AI views with better casing and title styles.
- Fixed unused `spinnerFrame` field in the Model struct.
- Added 3 new E2E tests for the polished UI and AI drafting status.
- All 65 E2E tests passing.

## Verification Pass: End-to-End ✅
- App starts with zero errors (smoke test passes)
- All views render correctly (Dashboard, Review, Import, AI, Settings, Browser, Cram)
- Core user interactions respond as expected (flashcard reveal, grading, navigation)
- State successfully persisted to SQLite (verified across restarts)
- All 12 unit test suites pass (62 unit tests)
- All 62 E2E tests passing
- ./scripts/verify.sh executes successfully

