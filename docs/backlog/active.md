# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Bug-fix and performance pass (per the focus in `AGENTS.md`): audit the remaining
TUI render hot paths and storage/export error paths. Dictionary line-coordinate
lookups are now O(1), card-type statistics failures surface, and APKG export
requires valid deck metadata.

## Exact Next Action

No unfinished executable work remains. Candidate future work:

- Extract remaining view-local state off `Model` (screen files own render +
  keys only today).
- Move history/settings persistence off synchronous TUI input handlers into
  `tea.Cmd`s so a slow repository cannot stall keyboard or mouse handling.

## Completed This Pass

- Performance: dictionary rendering now tracks builder line counts incrementally
  instead of repeatedly allocating `String()` snapshots and rescanning the
  growing buffers for every hitbox coordinate.
- Bug fix: collection/deck statistics now propagate card-type query, scan, and
  row-iteration errors rather than returning plausible but incomplete data.
- Bug fix: APKG export now aborts with context when deck metadata cannot be
  loaded instead of silently exporting internal deck IDs as display names.
- Regression tests cover embedded-newline counting, malformed card-type rows,
  and APKG deck-metadata failures.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- At least 2–3 concrete correctness or hot-path performance improvements land
  with focused regression tests.
- `go test ./...` and `./scripts/verify.sh` pass.

## Last Verification

- Previous pass: `./scripts/verify.sh` passed on 2026-08-23 (35 E2E tests).
- Current pass: `go test ./internal/tui ./internal/storage/sqlite` passed.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.43s).

## Repository State

- All session work is verified and committed; the working tree is clean.
