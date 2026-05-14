# Active Backlog

Last updated: 2026-05-14 18:45 +07

## Current Milestone

Autonomous end-to-end improvement pass on 2026-05-14 (batch 6).

## Current Work

No active executable work. Batch 6 is complete and fully verified.

## Exact Next Action

User can give new instructions.

## Top Issues / Priorities

1. Preserve the passing launch/test baseline while improving real end-to-end behavior.
2. Add durable German content that is discoverable through the existing deck/browser flows.
3. Improve visual hierarchy and status feedback without increasing height-sensitive layouts.
4. Harden keyboard/mouse edge cases that could leave stale selections or invalid cursors.
5. Expand E2E coverage with at least 8 new tui-tester tests for new content and core views.

## Acceptance Criteria (2026-05-14, batch 6)

- [x] Baseline `./scripts/tui_smoke.sh` passes (`deutsch-tui smoke ok`).
- [x] Baseline `go test ./...` passes.
- [x] App starts with zero errors after changes.
- [x] Core screens render via E2E coverage.
- [x] Keyboard and mouse navigation have targeted coverage.
- [x] SQLite persistence remains covered by tests.
- [x] At least 8 new tui-tester E2E tests are added and passing.
- [x] At least 2 new content/features are implemented.
- [x] At least 2 UI/UX improvements across different views are implemented.
- [x] At least 2 bugs or potential bugs are fixed.
- [x] `./scripts/verify.sh` passes with zero errors or warnings.
- [x] Completed work is summarized in `docs/backlog/done.md`.

## Last Verified

- `go test ./...` passed.
- `./scripts/verify.sh` passed with 301 E2E tests.

## Blockers

None.
