# Active Backlog

Last updated: 2026-04-29

## Current Milestone

Milestone 2: build the full review session flow on top of the stable scaffold.

## Next Action

Continue Milestone 2: build the full review session flow on top of the current scaffold.

## Acceptance Criteria

- Review session records FSRS-backed grades into SQLite.
- Due card list updates after grading.
- Review flow is usable with keyboard and mouse.
- Compact, medium, and wide layouts remain covered by tests.

## Blockers

- None recorded.

## Last Verified

- 2026-04-29: `GOCACHE=/tmp/deutsch-tui-gocache go test ./...`
- 2026-04-29: `GOCACHE=/tmp/deutsch-tui-gocache go vet ./...`
- 2026-04-29: `./scripts/verify.sh`
