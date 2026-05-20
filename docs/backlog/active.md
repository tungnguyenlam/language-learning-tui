# Active Backlog

Last updated: 2026-05-20

## Current Milestone

Decks and Browser Scrollbar Alignment - IN PROGRESS.

## Completed Work

- [x] **E2E Test Race Condition Fix:** Made Daily Goal settings adjustments update UI state optimistically to prevent overlapping asynchronous commands from dropping keypresses during rapid test execution.
- [x] **Content Expansion:** Added new **A1 Colors and Shapes** vocabulary deck (23 cards) focusing on foundational adjectives and geometric shapes.
- [x] **UI Polish & Rendering Bug Fix:** Removed redundant "Goal Met" text from the Dashboard progress line, and fixed a severe rendering bug in the Browser Card Preview where shorter lines would not properly pad out their backgrounds, causing overlapping ghost characters on terminal rerenders.
- [x] **Decks View Scrollbar Gaps Fix:** Added dynamic truncation for deck description and tags in Decks view to guarantee lines never exceed viewport width and scrollbars remain perfectly aligned with zero gaps.
- [x] **Browser View Scrollbar Gaps Fix:** Added dynamic truncation for card prompt and tags in Browser view to ensure scrollbar remains aligned and does not get cut off.
- [x] **Settings View Scrollbar Gaps Fix:** Added dynamic truncation for AI templates and API credentials in Settings view to ensure scrollbar remains perfectly aligned.
- [x] **Statistics View Scrollbar Gaps Fix:** Added dynamic truncation for deck names and other stats in Statistics view.
- [x] **E2E Regressions Fix:** Identified and reverted a buggy global padding in `model.go` that was causing severe line-merging glitches in E2E tests.

## Exact Next Action

Ship improvements.

## Top Issues / Priorities

None.

## Last Verified

- `go build ./...` passed.
- `./scripts/verify.sh` passed (317 E2E tests passing).

## Blockers

None.
