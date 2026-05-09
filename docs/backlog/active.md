# Active Backlog

Last updated: 2026-05-09 (batch E2)

## Current Milestone

Finalizing high-quality TUI experience with expanded content and robust UX.

## Completed Work

### UI Improvements
- [x] Add visual feedback for keyboard shortcuts in all views
- [x] Enhance Review view with session progress bar and better card styling
- [x] Refine Dashboard with boxed sections and color coding
- [x] Boxed AI Drafts list and improved preview structure

### Content Expansion
- [x] B1-B2 Science & Technology deck
- [x] German Proverbs & Idioms deck
- [x] Manual content seeding (Shortcut 'S')
- [x] Integration of all standard decks into Go source registry

### Bug Fixes & Stability
- [x] Double-grading prevention via state guard
- [x] AI view hitbox recalculation after layout changes
- [x] Unicode-aware backspace and text input hardening

## Next Steps

### Planned Enhancements (Batch F)
- [ ] Implement "Cram All" shortcut from Decks view
- [ ] Add B2-C1 "Philosophy & Literature" deck
- [ ] Implement "Filter by Tag" in Decks view search
- [ ] Add optional typing exercises for Review view
- [ ] Implement deck export to Anki APKG format (already in progress)

## Current State

- Fully verified: `./scripts/verify.sh` passes with 133 E2E tests.
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Modular monolith with pluggable content registry and SQLite persistence.
