# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 improvement pass is complete: Cram sessions stay trapped,
Practice Hub filter keys match the visible list, and remaining view handlers
live on their screens.

## Exact Next Action

No unfinished executable work remains. A future pass can start extracting
view-local state off `Model` now that Settings, Cram, AI, Practice, and
Dictionary keys are co-located with their screens, or split the large
async-message switch in `internal/tui/model.go`.

## Completed This Pass

- Active Cram consumes unhandled keys so `s`/`w`/`[`/`]` cannot abandon a
  session; `q` still exits.
- Practice Hub `/` filter accepts Bubble Tea `"space"`, Enter starts the
  visible trainer, and j/k wrap on the filtered list.
- Settings, Cram, AI, Practice, and Dictionary key handlers moved from
  `keys.go` into their screen files. Dictionary detail Space now scrolls.
- Spotlight overlay delegates to `dictionaryScreen.HandleKey`.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- During active Cram, `s`/`w`/`[`/`]` stay on Cram and do not change the deck.
- Practice Hub filter accepts a space character; Enter starts the filtered
  trainer; j/k wrap within the visible list.
- Settings, Cram, AI, Practice, and Dictionary `HandleKey` implementations
  live in `screen_*.go`.
- Focused tests, `go test ./internal/tui`, and `./scripts/verify.sh` pass.

## Last Verification

- `go test ./internal/tui` passed on 2026-08-22.
- `go test ./...` passed on 2026-08-22.
- `./scripts/verify.sh` passed on 2026-08-22: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 38.15s).

## Repository State

- Improvement pass is complete and fully verified on `main`.
