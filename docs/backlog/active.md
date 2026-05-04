# Active Backlog

Last updated: 2026-05-04

## Current Milestone

Autonomous Feature Pass 20: Multi-Deck Studies & SRS Refinement

## Completed Work

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
- [x] **Improve Navigation**: Allow jumping to specific tabs using number keys (1-6)
- [x] **Fix Cram Mode bug**: Number keys now correctly select filters instead of switching views

### Cloze Deletion Support
- [x] **Implemented Cloze Deletion card kind** in core models
- [x] **Added `parseClozes` logic** to `internal/content/anki.go` for `{{c1::...}}` syntax
- [x] **Updated TUI Review view** to render clozes with `[...]` placeholders
- [x] **Verified with unit tests** in `anki_test.go` and E2E test `test_cloze_deletion.py`

### Dashboard & Content
- [x] **Implemented "Grammar Tip of the Day" system** with rotating daily tips
- [x] **Added responsive Grammar Tip box** to the Dashboard (with height-based visibility check)
- [x] **Created "German Grammar Essentials" deck** featuring both MCQ and Cloze exercises
- [x] **Updated README.md and documentation** to reflect new learning capabilities

### Documentation & Agent Instructions (Pass 20)
- [x] **Expanded all package READMEs** in `internal/` with key symbols and responsibilities.
- [x] **Added subtree `AGENTS.md`** for `core`, `app`, and `srs` packages.
- [x] **Consolidated root `AGENTS.md`** with a "Common Tasks" index.
- [x] **Updated `docs/agent/index.md`** to reflect the new documentation structure.
- [x] **Refined Content and TUI docs** to include Cloze and responsive layout details.

## Next Action

Implement deck search and filtering in the Decks view.
