# Active Backlog

Last updated: 2026-05-28

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **AI Feature Improvements:** Implemented automatic provider activation on API key entry, added inline tag extraction for AI drafting, and updated Anthropic model defaults.
- [x] **Gender Trainer Practice Mode:** Added a new interactive mode for practicing German noun genders (der/die/das), including automated noun extraction and session scoring.
- [x] **Cloze UI Polish:** Highlight revealed cloze answers in standard review mode.
- [x] **AI Tutor Explanation:** Added `Shift+H` shortcut to get card explanations from AI.
- [x] **Code Deduplication:** Introduced `executeBulkAction` and `executeSingleAction` helpers in `internal/tui/handlers.go`, reducing repetitive boilerplate for browser and review actions.
- [x] **Rendering Consistency:** Centralized `formatReviewInterval` and `renderReviewHistory` in common files (`utils.go` and `render_views.go`) for consistent card metadata display across views.
- [x] **RenderContext Enhancements:** Added `RegisterAction` and `WriteAction` to `RenderContext` to streamline the implementation of interactive UI elements.
- [x] **Dashboard Interactivity:** Added individual hitboxes for recently studied decks and interactive hitboxes for the 'Quick Actions' section.
- [x] **New Content: False Friends Mastery:** Added a new deck with 12 high-quality cards covering German-English false friends.
- [x] **E2E Improvements:** Added `e2e_tests/test_may23_improvements.py` to verify dashboard interactivity and new content.
- [x] **NativeTTS Cache Fix:** Fixed macOS cache mismatch where synthesis output target used `.aiff` while existence checks queried `.wav`.
- [x] **Browser Audio Playback & Select All:** Added `p` (play audio) and `a` (select/deselect all) shortcuts to the Card Browser, with unit tests and help overlay updates.

## Exact Next Action

Implement a "Verb Conjugation Trainer" to further expand interactive learning modes.

## Top Issues / Priorities

None.

## Last Verified

- `./scripts/verify.sh` passed with 0 test failures (335 passed).

## Blockers

None.
selection) or additional learning modes.

## Top Issues / Priorities

None.

## Last Verified

- `./scripts/verify.sh` passed with 0 test failures (332 passed).

## Blockers

None.
