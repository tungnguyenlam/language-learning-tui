# Active Backlog

Last updated: 2026-08-12

## Current Milestone

Dashboard, Statistics, and Decks screen-ownership pass is complete and fully verified.

## Exact Next Action

No unfinished executable work remains. A future pass can migrate Settings, Cram, AI, Practice, or
Dictionary individually from `keys.go` into its registered screen implementation.

## Completed This Pass

- Moved Dashboard key handling into `dashboardScreen`, with direct coverage for recent-deck
  selection and unavailable recent-deck shortcuts.
- Moved Statistics scrolling and CSV-export key handling into `statisticsScreen`, with direct
  coverage for bounded viewport movement and export progress.
- Moved Decks search, selection, limit editing, export, study-launch, merge, and delete key handling
  into `decksScreen`; existing limit and export tests now use the registered screen contract.
- Reduced global `keys.go` by 214 lines while preserving the shared routing layer.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- `internal/tui/keys.go` still contains legacy handlers for views not yet fully co-located; migrate
  them individually behind the existing screen contract rather than attempting a broad rewrite.

## Acceptance Criteria

- Dashboard, Statistics, and Decks key behavior lives with each registered screen rather than in
  global `keys.go`.
- Focused tests call the production `screen.HandleKey` contract rather than legacy model helpers.
- Preserve package boundaries and user-visible behavior; add focused tests where logic changes or
  a reusable contract is introduced.
- Run targeted tests after each improvement and `./scripts/verify.sh` before handoff.
- Record completed work in `docs/backlog/done.md` and commit verified changes.

## Last Verification

- After Dashboard screen ownership: `go test ./internal/tui` passed on 2026-08-12.
- After Dashboard screen ownership: `./scripts/verify.sh` passed on 2026-08-12 (34 core E2E tests).
- After Statistics screen ownership: `go test ./internal/tui` passed on 2026-08-12.
- After Statistics screen ownership: `./scripts/verify.sh` passed on 2026-08-12 (34 core E2E tests).
- After Decks screen ownership: `go test ./internal/tui` and `git diff --check` passed on 2026-08-12.
- Final pass: `./scripts/verify.sh` passed on 2026-08-12: Go tests, vet, offline dict.cc import
  (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 36.97s).
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

- Dashboard/Statistics/Decks screen ownership pass is committed on `main`.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
- Pre-existing uncommitted documentation changes in `docs/agent/index.md`,
  `docs/backlog/done.md`, and `prompt/debug.md` belong to the prior debugging-prompt work and must
  be preserved.
