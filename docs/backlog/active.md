# Active Backlog

Last updated: 2026-05-04

## Current Milestone

Autonomous Feature Pass 19: Advanced Learning & UX Polish

## Completed Work

### Cloze Deletion Support
- [x] **Implemented Cloze Deletion card kind** in core models
- [x] **Added robust Cloze parser** supporting Anki-style `{{c1::text::hint}}` syntax
- [x] **Enhanced TUI rendering for Cloze cards** with highlighted placeholders
- [x] **Fixed card generation logic** to prevent duplicate/ugly flashcards for Cloze notes
- [x] **Added comprehensive E2E test** for Cloze deletion flow (`e2e_tests/test_cloze_deletion.py`)

### MCQ Enhancements
- [x] **Expanded MCQ support in TSV imports** using the `notetype` column
- [x] **Updated parser to handle arbitrary choices** (e.g., `MCQ:choice1,choice2,choice3`)
- [x] **Added unit tests for MCQ choices** in `internal/content/anki_test.go`

### Dashboard & Content
- [x] **Implemented "Grammar Tip of the Day" system** with rotating daily tips
- [x] **Added responsive Grammar Tip box** to the Dashboard (with height-based visibility check)
- [x] **Created "German Grammar Essentials" deck** featuring both MCQ and Cloze exercises
- [x] **Updated README.md and documentation** to reflect new learning capabilities

### Bug Fixes & Refinement
- [x] **Fixed Dashboard layout overflows** caused by large boxes on small terminals
- [x] **Synchronized Cloze note card generation** to only produce intended study cards
- [x] **Fixed gofmt violations** across content and test files
- [x] **Hardened E2E test suite** with better environment handling and stabilization

### Current Status (Verification)
- [x] App starts with zero errors
- [x] All views render correctly (including new responsive Dashboard elements)
- [x] Core user interactions respond as expected
- [x] State successfully persisted to SQLite
- [x] All Go unit tests pass (10 test suites)
- [x] All 77 E2E tests pass (including new Cloze test)
- [x] ./scripts/verify.sh executes successfully


## Planned Features

### Advanced Learning
- [ ] **Implement Reverse Flashcards** (Basic (and reversed card) type)
- [ ] **Add "Reverse" support to TSV import** via `notetype:Reverse`
- [ ] **Improve "Examples" card generation** to automatically create context-based MCQs

### UX Polish & Dashboard
- [ ] **Add "Total Decks" and "Active Decks" counters** to Dashboard statistics
- [ ] **Enhance Review Visuals**: Add subtle color coding for Grade options (Again=Red, Good=Green, etc.)
- [ ] **Improve Navigation**: Allow jumping to specific tabs using number keys (1-6)

### Content Expansion
- [ ] **Add "B1 - German Idioms" deck** to showcase advanced expressions

## Last Verification

- 2026-05-04: `./scripts/verify.sh` passed baseline check.
- 2026-05-03: `./scripts/verify.sh` passed, including all 77 E2E tests.

## Next Action

Implement Reverse Flashcards logic in `internal/content/anki.go`.
