# Active Backlog

Last updated: 2026-08-22

## Current Milestone

Bug-fix and performance pass (per the new focus in `AGENTS.md`): all
"today"/streak statistics now use local calendar days, and two O(n²) render
hot paths were made O(n).

## Exact Next Action

No unfinished executable work remains. Candidate future work:

- Extract remaining view-local state off `Model` (screen files own render +
  keys only today).
- Apply the incremental-line-count pattern from `render_dashboard.go` to the
  other render files that rescan the whole buffer per box
  (`render_dictionary.go`, `render_settings.go`, `render_views.go`,
  `screen_import.go`, `screen_ankiweb.go`).

## Completed This Pass

- Bug fix: `ReviewsToday` (Statistics) and per-deck `reviews_today` (Decks)
  anchored "today" to UTC midnight; daily-goal progress reset hours
  early/late outside UTC. Both now use `localDayStartUTC` (local midnight,
  DST-correct tomorrow). Regression test
  `TestReviewsTodayAnchorsToLocalMidnight` passes under TZ=Asia/Bangkok and
  TZ=America/New_York.
- Streak queries (`currentStreak`, `deckCurrentStreak`) group by
  `date(substr(reviewed_at, 1, 19), 'localtime')` instead of the UTC date
  substring, aligning with `ReviewsPerDay`; `LIMIT 365` now covers 365 local
  days. Guard test `TestCurrentStreakCountsLocalDays`.
- Performance: `render_dashboard.go` replaced 14 full-buffer
  `strings.Count(db.String(), "\n")` scans per frame with an incremental
  counter (output byte-identical).
- Performance: `renderTypingDiff` backtrace appends + writes in reverse
  instead of O(n²) slice prepends per keystroke.
- `prompt/improve.md` lists performance optimization as a direction;
  `AGENTS.md` Current State declares the bug-fix/performance focus.
- New notice `2026-08-22-local-day-statistics.md`.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Statistics, streaks, and per-day activity all agree on the local calendar
  day in any timezone.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-22 after the local-day stats fix and
  render hot-path optimizations: Go tests, vet, offline dict.cc import, smoke
  test, binary build, core E2E suite (35 passed in 37.29s).
- Storage streak/today regression tests pass under TZ=Asia/Bangkok and
  TZ=America/New_York.

## Repository State

- All session work verified on this working tree; ready to commit.

