# Active Backlog

Last updated: 2026-05-09 (batch G)

## Current Milestone

Finalizing high-quality TUI experience with expanded content and robust UX.

## Completed Work

### UI Improvements
- [x] Add visual feedback for keyboard navigation
- [x] Responsive layout for Statistics view
- [x] Dashboard: Modernized layout with info-rich sidebar in Wide mode
- [x] Decks View: Integrated mini progress bars for each deck
- [x] AI View: Added interactive, clickable suggested topics
- [x] Statistics: Restored 'due' count for bookmarked cards

### Content & Learning
- [x] Added B2-C1 "German News & Current Events" deck
- [x] Added B2 "Environment & Climate" deck
- [x] Implemented Review Summary screen with accuracy and time tracking
- [x] Improved typing exercise normalization for German characters
- [x] Support space-key for submitting typing answers

### Core Features
- [x] Enhanced Anki APKG export (Cloze support, media files, correct checksums)
- [x] Added hidden test indicators for Cram and Review modes to improve E2E reliability

## Current State

- Verified: All unit tests pass; 137/137 E2E tests pass.
- Environment: Darwin, Go 1.x, Python 3.12.
- Architecture: Modular monolith with pluggable content registry and SQLite persistence.
