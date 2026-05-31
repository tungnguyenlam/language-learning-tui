# Active Backlog

Last updated: 2026-05-31

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **Adjective Ending Trainer:** Added a fourth practice mode to the Practice Hub for mastering German adjective endings (weak, strong, and mixed declensions) with interactive fill-in-the-blank sentences and dynamic grammar contexts.
- [x] **Leech E2E Test Fix:** Resolved a false-positive E2E timeout in `test_leech_detection_in_statistics` by updating the screen expectation from `"cards due"` (which disappears when revealed) to `"a Again"` (present on the grade selection overlay).
- [x] **Practice Hub & Case Trainer:** Consolidated interactive modes into a central hub and added a new "Case Ending Trainer" for mastering German grammar.
- [x] **Verb Conjugation Trainer:** Added a new interactive mode for practicing German verb conjugations (ich, du, er/sie/es, etc.), including automated person selection and session scoring.
- [x] **Dashboard & UI Polish:** Improved Quick Actions layout, fixed Settings navigation boundaries, and resolved race conditions in daily goal adjustment.
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
- [x] **Grammar Hint Overlay:** Added a `Shift+G` shortcut in the Review view to toggle a contextual grammar hint based on the card contents.
- [x] **List Navigation Shortcuts:** Added `g` (jump to top) and `G` (jump to bottom) shortcuts to Browser, Settings, and Decks views for faster navigation.
- [x] **Browser Contextual Hints:** Integrated the grammar hint analysis directly into the Card Preview box in the Browser view to show contextual word usage information.

## Exact Next Action

Explore further visual polish or performance optimization, or move on to the next set of user feedback.

## Top Issues / Priorities

None.

## Last Verified

- `./scripts/verify.sh` passed with 0 test failures (338 passed).

## Blockers

None.
