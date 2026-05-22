# Active Backlog

Last updated: 2026-05-22

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **E2E Test Suite Stabilization:** Resolved 18 test failures caused by recent UI refactoring, fixing responsive breakpoints, help overlay clipping, and missing UI elements.
- [x] **Search Highlights:** Implemented visual feedback for search queries in Browser and Decks views.
- [x] **UI Corruption Fix:** Resolved terminal wrapping and "ghost" character glitches in Decks view on narrow terminals.
- [x] **Native TTS Fallback:** Added support for system-native TTS engines (macOS 'say', Linux 'espeak') as a fallback for Edge TTS.
- [x] **A2 Medical Appointment Deck:** Added a new practical deck for medical situations.
- [x] **Audio Reliability & UX Fix:** Implemented process tracking to prevent overlapping audio playback, fixed hanging "Generating..." status messages, and added support for stripping Anki sound markers.

## Exact Next Action

Look for further codebase improvements (refactoring, new features, or content).

## Top Issues / Priorities

None.

## Last Verified

- `./scripts/verify.sh` passed with 0 test failures (320 passed).

## Blockers

None.
