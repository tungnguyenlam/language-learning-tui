# Active Backlog

Last updated: 2026-05-08

## Current Milestone

End-to-end TUI correctness and autonomy batch

## Completed Work

### UI Polish & Developer Experience (DONE)
- [x] Refined UI spacing and layout in Settings: aligned keys, updated formatting.
- [x] Enhanced AI Drafts view with placeholder strings for better context.
- [x] Improved Makefile for easier local developer testing (`make fmt`, `make build`, etc.).
- [x] Created and successfully ran new E2E test `e2e_tests/test_settings_ui.py` for Settings UI layout.

### Content Library Expansion (IN PROGRESS)
- [x] Expand Grammar Essentials deck from 14 to 60+ cards
- [x] Expand German Idioms deck from 10 to 40+ cards
- [x] Add new Workplace & Office German deck (`b1-workplace-office.tsv`)
- [x] Add new Shopping & Services deck
- [x] Add new Health & Body deck
- [x] Add new Emotions & Feelings deck
- [x] Add new Transport & Directions deck

### Bug Hunting & Stability Polish
- [x] Fixed E2E test truncation flakiness by configuring `lines=40` for specific tests (`test_audio_autoplay.py`, `test_ui_sanity.py`).
- [x] **TUI Bug Fixes**: Resolved cursor panics, leaky modals, and ignored errors.
- [x] **Data Integrity**: Added transactions, fixed timezone inconsistencies, and enabled foreign keys in connection pool.
- [x] **Performance**: Indexed `reviewed_at` and clamped integer overflow in SRS scheduler.

## Next Action

- [ ] Phase 1: run `cmd/deutsch-tui` through smoke and tui-tester to capture current launch/render state.
- [ ] Phase 2: list top launch, navigation, persistence, and E2E coverage blockers.
- [x] Phase 3 batch A: implemented Review cursor safety, AI/Settings template index safety, empty import/export guards, and cross-view UI hints.
- [x] Phase 3 batch B: added 8 tui-tester E2E tests covering Dashboard, Review persistence, Decks, Import, AI, Settings, and view navigation; added B1 Apartment & Housing content.
- [x] Phase 4: `./scripts/verify.sh` passes with Go tests, smoke, and 116 E2E tests.

## Verification State

- Last verified: `./scripts/verify.sh` passes with no warnings; 116 E2E tests pass. `go run ./cmd/deutsch-tui --help` succeeds.
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Scalable content Registry with pluggable sources
