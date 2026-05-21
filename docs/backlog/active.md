# Active Backlog

Last updated: 2026-05-21

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **Search Highlights:** Implemented visual feedback for search queries in Browser and Decks views.
- [x] **Native TTS Fallback:** Added support for system-native TTS engines (macOS 'say', Linux 'espeak') as a fallback for Edge TTS.
- [x] **A2 Medical Appointment Deck:** Added a new practical deck for medical situations.
- [x] **Audio Reliability & UX Fix:** Implemented process tracking to prevent overlapping audio playback, fixed hanging "Generating..." status messages, and added support for stripping Anki sound markers.

## Exact Next Action

Monitor search highlight performance in very large lists (though currently optimized).

## Top Issues / Priorities

None.

## Last Verified

- `go build ./...` passed.
- `go test ./internal/audio ./internal/content` passed.
- `pytest e2e_tests/test_may21_improvements.py` passed.

## Blockers

None.
