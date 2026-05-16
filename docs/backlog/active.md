# Active Backlog

Last updated: 2026-05-16 13:45 +07

## Current Milestone

Improve cross-platform audio playback - COMPLETED.

## Completed Work

- [x] UX Improvement: Added interactive "Recently Studied" decks on Dashboard (Clickable + `!/@/#` shortcuts).
- [x] UX Improvement: Enhanced Review view with prominent red "LEECH" badges and bolder state badges.
- [x] UI Improvement: Added celebratory "GOAL MET 🏆" badge on Dashboard when daily goal is reached.
- [x] Content Expansion: Added "B2 Business German: Meetings & Negotiations" deck (30 notes).
- [x] Logic Fix: Ensured statistics reload immediately after recording a review (fixed `0/X` reviews bug).
- [x] Added 3 new E2E tests in `test_batch11_improvements.py`.
- [x] UX Improvement: Added 1-4 grading shortcuts in Review view (verified with E2E test).
- [x] UX Improvement: Added Mouse Wheel support for Decks view scrolling.
- [x] Content Expansion: Added "B2 Programming & Software Engineering" deck (30 technical notes).
- [x] Fixed E2E regression in `test_mcq_navigation_bug.py` caused by new grading shortcuts.
- [x] Fixed 10 failing E2E tests (Batch 9)
- [x] Verified B2 Job Application, C1 Philosophy & Ethics, A1 City & Directions decks exist (40 cards each)
- [x] UI Improvement: Enhanced session summary with efficiency stat (cards/min)
- [x] UI Improvement: Added motivational tips in Review empty state
- [x] AI Improvement: Added Grammar Breakdown and City Directions topics to AI suggestions
- [x] Added 8 new E2E tests for new features

## Exact Next Action

None - terminal audio playback now selects platform-appropriate players.

## Top Issues / Priorities

None.

## Acceptance Criteria

- [x] macOS playback prefers `afplay` with `mpv`/`ffplay` fallbacks.
- [x] Linux playback tries common terminal players (`mpv`, `ffplay`, `play`, etc.).
- [x] Windows playback tries `mpv`/`ffplay` and PowerShell fallback.
- [x] Missing players return a clear actionable error.
- [x] `go test ./...` passes.
- [x] App smoke check passes.

## Last Verified

- `go test ./...` passed.
- `./deutsch-tui-bin -data-dir /tmp/deutsch-tui-audio-player-smoke -smoke` passed.

## Blockers

None.
