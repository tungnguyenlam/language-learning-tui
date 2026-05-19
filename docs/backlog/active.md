# Active Backlog

Last updated: 2026-05-20

## Current Milestone

Autonomous Improvement Pass - COMPLETED.

## Completed Work

- [x] **E2E Test Race Condition Fix:** Made Daily Goal settings adjustments update UI state optimistically to prevent overlapping asynchronous commands from dropping keypresses during rapid test execution.
- [x] **Content Expansion:** Added new **A1 Colors and Shapes** vocabulary deck (23 cards) focusing on foundational adjectives and geometric shapes.
- [x] **UI Polish & Rendering Bug Fix:** Removed redundant "Goal Met" text from the Dashboard progress line, and fixed a severe rendering bug in the Browser Card Preview where shorter lines would not properly pad out their backgrounds, causing overlapping ghost characters on terminal rerenders.

## Exact Next Action

None - milestone complete.

## Top Issues / Priorities

None.

## Last Verified

- `go build ./...` passed.
- `./scripts/verify.sh` passed (317 E2E tests passing).

## Blockers

None.
