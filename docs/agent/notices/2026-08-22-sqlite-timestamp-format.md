# SQLite Timestamp Format Trap

Status: active
Date: 2026-08-22

## Context

`RecordReview` passes `time.Time` values (`result.Reviewed.UTC()`) to the
SQLite driver, which stores them in Go's default string format:

```
2026-08-22 14:52:14.558668 +0000 UTC
```

SQLite's date functions (`date()`, `strftime()`) cannot parse the
` +0000 UTC` suffix and silently return NULL. `ReviewsPerDay` used
`date(reviewed_at, 'localtime')` and therefore always returned an empty map —
the Dashboard "Recent Activity" sparkline and the Statistics per-day chart
never showed real data. TUI tests use `mockRepo` (no SQL), so nothing caught
it.

## Rule

When a SQL query needs a date/datetime from a `time.Time` column, extract the
prefix first:

```sql
date(substr(reviewed_at, 1, 19), 'localtime')  -- "YYYY-MM-DD HH:MM:SS", parsed as UTC
```

- `substr(x, 1, 10)` gives the UTC calendar date (used by the streak and
  recent-decks queries; acceptable for day-granularity grouping).
- `substr(x, 1, 19)` + `'localtime'` gives the true local date and also
  survives RFC3339 (`T` separator) if the driver format ever changes.
- Plain `>=` / `<` comparisons against `time.Time` params work because the
  driver formats both sides identically and the strings align
  lexicographically.

## Regression Coverage

`TestReviewsPerDayCountsRecordedReviews` in
`internal/storage/sqlite/sqlite_test.go` writes through the real
`RecordReview` path and asserts the local-date bucket. It fails (empty map)
if the `date()` form is reintroduced.

## Revisit When

The storage layer changes how timestamps are written (e.g. a migration to
RFC3339 or Unix epoch columns). Then all `substr(...)` date queries should be
revisited together.
