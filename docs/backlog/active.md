# Active Backlog

Last updated: 2026-04-29

## Current Milestone

Milestone 7: Statistics and Progress Visualization.

## Next Action

Implement a Statistics view showing learning progress (e.g., retention rate, cards learned over time).

## Acceptance Criteria

- New "Stats" view accessible via Tab or shortcut '7'.
- SQLite repository provides methods to fetch historical review data.
- Stats view displays retention rate and a simple bar chart of reviews per day (text-based).
- Unit tests cover stats calculation logic.

## Blockers

- None.

## Last Verified

- 2026-04-29: `./scripts/verify.sh` (includes Go tests, smoke tests, and E2E tests)
