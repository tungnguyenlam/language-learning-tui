# Active Backlog

Last updated: 2026-05-15 09:00 +07

## Current Milestone

Autonomous end-to-end improvement pass on 2026-05-15 (batch 7).

## Current Work

- [x] Content Expansion: Add **A2 German Travel & Booking** deck (40 cards).
- [x] Content Expansion: Add **B1 German Housing & Apartment** deck (40 cards).
- [x] Content Expansion: Add **C1 German Environment & Sustainability** deck (40 cards).
- [x] UI Improvement: Add "Hint" support in Review view (Shortcut 'h').
- [x] UI Improvement: Add "Next 24h Forecast" to Dashboard.
- [x] UI Improvement: Add "Cards Added per Day" chart to Statistics.
- [x] AI View: Show incremental drafting progress (via status updates).
- [x] DX/Testing: Add 8+ new E2E tests for these features.

## Exact Next Action

User can give new instructions.

## Top Issues / Priorities

1. Complete at least 6-8 distinct improvements.
2. Ensure high-quality German content.
3. Maintain zero-error launch and passing tests.
4. Enhance visual feedback across views.

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
