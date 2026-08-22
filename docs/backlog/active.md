# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 E2E decoupling pass makes every starter-count assertion
count-agnostic, so `StarterDeck()` can grow without touching tests. A
pre-existing broken Settings daily-goal E2E test was fixed along the way.

## Exact Next Action

No unfinished executable work remains. Candidate future work:

- Extract remaining view-local state off `Model` (screen files own render +
  keys only today).
- Streak queries group by UTC date (`substr(reviewed_at, 1, 10)`); switching
  to local-date grouping would align them with `ReviewsPerDay`.

## Completed This Pass

- E2E due-count decoupling: new `e2e_tests/e2e_helpers.py` with
  `read_due_count` / `read_cards_found`; 16 test files now compute `due`
  from the live dashboard instead of hard-coding 52/51/54.
- Fixed pre-existing broken `test_daily_goal_setting_persists_after_restart`
  (never navigated to the Daily Goal row before pressing `+`; fails on the
  pristine tree, not run by the core suite).
- Resolved notice `2026-08-22-seeded-content-e2e-due-count.md`.
- Split `internal/tui/model.go` (2062 → 1437 lines): the ~45-case async-message
  switch moved to `messages.go` (`handleAsyncMsg`), mouse handling to
  `mouse.go` (`handleMouseMsg`). `Update` is now a small dispatcher.
- User-facing fix: `ReviewsPerDay` used `date(reviewed_at, 'localtime')`,
  which SQLite cannot parse for Go-formatted timestamps, so the Dashboard
  "Recent Activity" sparkline and Statistics per-day chart never showed data.
  Fixed with `date(substr(reviewed_at, 1, 19), 'localtime')`; regression test
  `TestReviewsPerDayCountsRecordedReviews` writes through the real
  `RecordReview` path. See notice `2026-08-22-sqlite-timestamp-format.md`.
- Dashboard empty state: "Recent Activity" shows "No reviews yet — press 3!"
  when the last-14-days total is zero (covered by
  `TestDashboardShowsActivityPlaceholderWithoutReviews`).

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- No E2E test hard-codes the starter due count; helpers fail loudly when the
  dashboard count line is missing.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-22 after the `ReviewsPerDay` fix:
  Go tests, vet, offline dict.cc import, smoke test, binary build, core E2E
  suite (35 passed in 37.56s).
- Live tui-tester check: fresh DB shows the empty-state placeholder; after
  grading one card the sparkline renders a real bar.

## Repository State

- All session work verified on this working tree; ready to commit.

