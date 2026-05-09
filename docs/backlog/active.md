# Active Backlog

Last updated: 2026-05-10 (Verification Complete)

## Current Milestone

Finalizing high-quality TUI experience with expanded content and robust UX.

## Next Actions

- [ ] Add more C1-C2 advanced content
- [ ] Implement audio recording or external dictionary integration
- [ ] Add deck export filtering (by date or success rate)

## Completed Work (Latest)

### Content & Learning
- [x] Added B2-C1 "Art & Culture" deck (TSV)
- [x] Migrated 5 Go decks to TSV for better editability
- [x] Dynamic TSV deck discovery and loading

### UI & UX
- [x] Dashboard: Quick Actions menu
- [x] Dashboard: Last Session stats summary
- [x] Browser: Tag filtering support (#)
- [x] Browser: Active deck syncing

### Core Features
- [x] Settings: Strict Character Normalization (toggle ss vs ß)
- [x] Fixed audio autoplay toggle persistence in Settings
- [x] Verified with 8 new E2E tests

### Verification
- [x] End-to-end verification: App launches without errors, all views render correctly
- [x] Core user interactions (flashcard reveal, grading, navigation) respond as expected
- [x] State is successfully persisted to SQLite
- [x] End-to-end tests pass using the tui-tester utility (145/145 tests passing)
- [x] Full verification suite (`./scripts/verify.sh`) executes successfully

## Current State

- Verified: All unit tests pass; 145/145 E2E tests pass.
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Modular monolith with dynamic content registry.
- Status: Application is working correctly end-to-end as verified by full test suite.

### UI Improvements
- [x] Add visual feedback for keyboard navigation
- [x] Responsive layout for Statistics view
- [x] Dashboard: Modernized layout with info-rich sidebar in Wide mode
- [x] Decks View: Integrated mini progress bars for each deck
- [x] AI View: Added interactive, clickable suggested topics
- [x] Statistics: Restored 'due' count for bookmarked cards
- [x] Fixed audio autoplay toggle persistence in Settings

### Content & Learning
- [x] Added B2-C1 "German News & Current Events" deck
- [x] Added B2 "Environment & Climate" deck
- [x] Implemented Review Summary screen with accuracy and time tracking
- [x] Improved typing exercise normalization for German characters
- [x] Support space-key for submitting typing answers

### Core Features
- [x] Enhanced Anki APKG export (Cloze support, media files, correct checksums)
- [x] Added hidden test indicators for Cram and Review modes to improve E2E reliability