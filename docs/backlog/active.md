# Active Backlog

Last updated: 2026-05-05

## Current Milestone

Autonomous Feature Pass 21: Bulk Operations & Advanced Search

## Completed Work

### Deck Search & Key Handling (Pass 20)
- [x] **Implemented Robust Deck Search**: Added `/` search to Decks view with dedicated `searchingDecks` state and improved UI feedback.
- [x] **Fixed Enter Key Handling**: Resolved a bug where `\r` and `\n` were not recognized as Enter in multiple views (Import, Settings, Decks, Browser).
- [x] **Improved Key Trapping**: Refactored `updateKey` to more robustly trap and delegate keys to active views during text input.
- [x] **Added E2E Tests**: Created `e2e_tests/test_new_pass_20.py` with 3 new tests covering deck search, selection, and settings.
- [x] **Verified Full Suite**: 87 tests passed with zero failures.

### Browser Management & Tagging (Pass 21)
- [x] **Implemented Tag Management** in Card Browser (T key to add/remove tags).
- [x] **Bulk Tagging Support**: Apply tags to multiple selected cards at once.
- [x] **Note-Card Sync**: Tags are persisted at the note level and synchronized across all card kinds.
- [x] **Tag Search**: Search by tags in the Browser using the `/` key.
- [x] **Input Guard**: Prevented view switching while typing in Search or Tag prompts.
- [x] **Storage Migrations**: Added `tags` column to `notes` and `cards` tables.

### Advanced Learning & UX Polish (Pass 19)
- [x] **Implemented Reverse Flashcards** (Basic (and reversed card) type)
- [x] **Add "Reverse" support to TSV import** via `notetype:Reverse`
- [x] **Improved "Examples" card generation** to automatically create context-based MCQs
- [x] **Add "B1 - German Idioms" deck** to showcase advanced expressions
- [x] **Add "Total Decks" and "Active Decks" counters** to Dashboard and Statistics
- [x] **Enhance Review Visuals**: Add subtle color coding for Grade options (Again=Red, Good=Green, etc.)
- [x] **Improve Navigation**: Allow jumping to specific tabs using number keys (1-9)
- [x] **Fix Cram Mode bug**: Number keys now correctly select filters instead of switching views

## Next Action

Implement card export filtering (by tag/deck) in the Import/Export view.
