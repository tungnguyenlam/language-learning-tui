# Active Backlog

Last updated: 2026-08-12

## Current Milestone

Developer-experience simplification pass is complete. Persisted card decoding, search debouncing,
routine transaction lifecycle, and Browser/Review screen ownership now have reusable boundaries.

## Exact Next Action

No unfinished executable work remains from this pass. A future pass can continue migrating one
legacy view at a time from `keys.go` into its registered `screen_*.go` implementation.

## Completed This Pass

- Added one canonical SQLite card-row decoder used by deck, note, due, bookmarked-due, and filtered
  card queries, removing five independently maintained field mappings.
- Added one Browser/Dictionary search debounce scheduler with independent per-view generations and
  focused coverage that stale generations cannot trigger queries.
- Moved Browser and Review key handling behind their registered screen implementations, updated
  direct tests to use the production screen contract, and reduced global `keys.go` by 436 lines.
- Added a commit-on-success SQLite transaction helper with direct rollback coverage and adopted it
  across deck/note/card upserts, review recording, card mutations, reset, merge, and tag cleanup.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- `internal/tui/keys.go` still contains legacy handlers for views not yet fully co-located; migrate
  them individually behind the existing screen contract rather than attempting a broad rewrite.

## Acceptance Criteria

- Complete at least 2–3 narrow improvements that make common development paths easier to locate,
  extend, or test.
- Preserve package boundaries and user-visible behavior; add focused tests where logic changes or
  a reusable contract is introduced.
- Run targeted tests after each improvement and `./scripts/verify.sh` before handoff.
- Record completed work in `docs/backlog/done.md` and commit verified changes.

## Last Verification

- Previous pass: `./scripts/verify.sh` and `go test -race ./...` passed on 2026-08-12.
- Current pass baseline: `go test ./...` passed on 2026-08-12.
- After canonical card-row decoding: `go test ./internal/storage/sqlite` passed on 2026-08-12.
- After shared search debouncing: `go test ./internal/tui` passed on 2026-08-12.
- After Browser/Review screen ownership: `go test ./internal/tui` passed on 2026-08-12.
- After shared transaction handling: `go test ./internal/storage/sqlite` passed on 2026-08-12.
- `./scripts/verify.sh` passed on 2026-08-12: Go tests, vet, offline dict.cc import (834,512
  entries), smoke test, binary build, and core E2E suite (34 passed in 36.17s).
- `go test -race ./...` passed on 2026-08-12.

## Repository State

- This developer-experience pass is committed on `main`; the branch remains ahead of `origin/main`
  because it also contains earlier local commits.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
- Pre-existing uncommitted documentation changes in `docs/agent/index.md`,
  `docs/backlog/done.md`, and `prompt/debug.md` belong to the prior debugging-prompt work and must
  be preserved.
