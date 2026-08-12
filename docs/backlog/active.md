# Active Backlog

Last updated: 2026-08-12

## Current Milestone

Decks multi-selection Escape bug is fixed and fully verified.

## Exact Next Action

No unfinished executable work remains. A future pass can migrate Settings, Cram, AI, Practice, or
Dictionary individually from `keys.go` into its registered screen implementation.

## Completed This Pass

- Reproduced the bug in the real TUI: select and deselect a deck, apply the `German` Decks filter,
  then press Escape once; the filter remains visible because the deselected ID remains in the map
  with a `false` value.
- Added a focused screen-contract regression test for the selection-map invariant and one-Escape
  filter clearing; the old behavior fails with `map[string]bool{"deck-1":false}`.
- Changed the Decks selection toggle to delete deselected IDs so the map represents only active
  selections.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- `internal/tui/keys.go` still contains legacy handlers for views not yet fully co-located; migrate
  them individually behind the existing screen contract rather than attempting a broad rewrite.

## Acceptance Criteria

- Deselecting a deck removes its ID from `deckSelected`.
- With no selected decks, one Escape clears an active Decks filter.
- The regression test, `go test ./internal/tui`, and `./scripts/verify.sh` pass.
- Record the completed fix in `docs/backlog/done.md` and commit verified changes.

## Last Verification

- Baseline `go test ./...` passed on 2026-08-12.
- Focused regression and `go test ./internal/tui` passed on 2026-08-12 after the fix.
- Post-fix real-TUI replay passed on 2026-08-12: `Filter: German` was visible before Escape and
  absent after one Escape.
- Final `go test ./...` passed on 2026-08-12.
- Final `./scripts/verify.sh` passed on 2026-08-12: Go tests, vet, offline dict.cc import (834,512
  entries), smoke test, binary build, and core E2E suite (34 passed in 36.86s).
- Failing baseline: `go test ./internal/tui -run
  '^TestDecksScreenEscapeClearsFilterAfterDeselect$' -count=1` failed on 2026-08-12 because the
  selection map retained `"deck-1": false`.
- Real-TUI reproduction passed on 2026-08-12 using a unique data directory and the required
  Observe → Reason → Act → Synchronize loop; after one Escape, `Filter: German` remained visible.
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

- Decks Escape/filter fix is complete and fully verified on `main`.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
