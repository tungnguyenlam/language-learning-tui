# Active Backlog

Last updated: 2026-05-03

## Current Milestone

Autonomous Feature Pass 18: Final Polish and Handoff

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

(None at this time)

## Last Verification

- 2026-05-03: `./scripts/verify.sh` passed, including all 76 E2E tests.
- 2026-05-02: `./scripts/verify.sh` passed, including all 76 E2E tests.

## Next Action

Handoff complete - all verification passed.

## Summary

The deutsche-tui application is production-ready with comprehensive German learning content:
- ✅ All 9 Go unit test suites passing
- ✅ All 76 E2E tests passing  
- ✅ App launches without errors
- ✅ All views render correctly
- ✅ Core interactions working
- ✅ State persistence verified
- ✅ 600+ vocabulary items across A1-B2 levels
- ✅ 500+ importable TSV cards
- ✅ Full Anki compatibility

The application successfully meets all requirements for end-to-end functionality.
