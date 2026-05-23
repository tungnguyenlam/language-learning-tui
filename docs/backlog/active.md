# Active Backlog

Last updated: 2026-05-23

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **Code Deduplication:** Introduced `executeBulkAction` and `executeSingleAction` helpers in `internal/tui/handlers.go`, reducing repetitive boilerplate for browser and review actions.
- [x] **Rendering Consistency:** Centralized `formatReviewInterval` and `renderReviewHistory` in common files (`utils.go` and `render_views.go`) for consistent card metadata display across views.
- [x] **RenderContext Enhancements:** Added `RegisterAction` and `WriteAction` to `RenderContext` to streamline the implementation of interactive UI elements.
- [x] **Dashboard Interactivity:** Added individual hitboxes for recently studied decks and interactive hitboxes for the 'Quick Actions' section.
- [x] **New Content: False Friends Mastery:** Added a new deck with 12 high-quality cards covering German-English false friends.
- [x] **E2E Improvements:** Added `e2e_tests/test_may23_improvements.py` to verify dashboard interactivity and new content.

## Exact Next Action

Look for further codebase improvements (refactoring, new features, or content).

## Top Issues / Priorities

None.

## Last Verified

- `./scripts/verify.sh` passed with 0 test failures (320 passed).

## Blockers

None.
