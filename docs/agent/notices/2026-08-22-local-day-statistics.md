# Local-Day Statistics

Status: active
Scope: internal/storage/sqlite (Statistics, Decks, streaks, ReviewsPerDay)
Related: `localDayStartUTC`, `currentStreak`, `deckCurrentStreak`, `ReviewsPerDay`, notice `2026-08-22-sqlite-timestamp-format.md`

## Why It Matters

`reviewed_at` is stored as a UTC Go timestamp string. All user-facing "day"
concepts — Dashboard daily-goal progress (`ReviewsToday`), per-deck
`reviews_today`, current streak, and the activity sparkline — must use the
local calendar day. Anchoring any of them to UTC midnight makes the counters
disagree with each other and reset hours early/late outside UTC.

## Required Behavior

- Day boundaries come from `localDayStartUTC(now)` (local midnight in UTC),
  never `time.Date(..., now UTC fields ..., time.UTC)`.
- Day grouping in SQL uses `date(substr(reviewed_at, 1, 19), 'localtime')`,
  never `substr(reviewed_at, 1, 10)` (that is the UTC date).
- Regression tests `TestReviewsTodayAnchorsToLocalMidnight` and
  `TestCurrentStreakCountsLocalDays` must pass; run them with a non-UTC `TZ`
  (e.g. `TZ=Asia/Bangkok`) when touching this area.

## Revisit When

`reviewed_at` storage changes format, or the app gains an explicit
user-configurable timezone/day-boundary setting.
