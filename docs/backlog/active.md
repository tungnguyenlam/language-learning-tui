# Active Backlog

Last updated: 2026-05-12 (Review grade key stuck-state regression fix complete)

## Current Milestone

Fix review cards getting stuck before `a/h/g/e` grading.

## Current Work

- [x] Traced the stuck path to extra Space/Enter presses on a revealed review card setting `gradingInProgress` without recording a grade.
- [x] Changed revealed-card Space/Enter handling to keep grading keys enabled and show the explicit `a/h/g/e` status hint.
- [x] Added a unit regression test proving `e` still grades after an extra reveal key.

## Exact Next Action

No active executable work. Next agent should follow the user’s newest request.

## Acceptance Criteria

- [x] Extra Space/Enter on a revealed review card does not set `gradingInProgress`.
- [x] `a/h/g/e` still grade normally after that extra reveal key.
- [x] Focused `go test ./internal/tui` passes.
- [x] Full `./scripts/verify.sh` passes with 238 E2E tests.

## Blockers

None.
