# Active Backlog

Last updated: 2026-08-12

## Current Milestone

Milestone 4 (Hybrid Dictionary) remains complete. The autonomous improvement pass for stale
practice loads, Anki cloze fidelity, and the keyboard-help modal is complete.

## Exact Next Action

No unfinished executable work remains from this pass. Select the next roadmap/backlog item when
starting a new improvement session.

## Top Issues

- None identified for this pass.

## Acceptance Criteria

- Generic trainer loads, grouped/out-of-order cloze ordinals, and modal/scrollable help behavior
  each have focused regression coverage.
- `./scripts/verify.sh` passes and completed work is recorded in `docs/backlog/done.md`.
- Changes are committed after verification.

## Last Verification

- Targeted `go test ./internal/content ./internal/tui` passed on 2026-08-12.
- `./scripts/verify.sh` passed on 2026-08-12: Go tests, vet, offline dict.cc import (834,512
  entries), smoke test, binary build, and core E2E suite (34 passed in 36.17s).
- `go test -race ./...` passed on 2026-08-12.

## Repository State

- `refactor/generic-trainer` was fast-forward merged into `main` on 2026-08-05 and both branches
  are pushed. `main` is at `f56400c`.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
