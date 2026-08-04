# Active Backlog

Last updated: 2026-08-05

## Current Milestone

Milestone 4 (Hybrid Dictionary) is complete and merged to `main`. No new feature milestone is
committed yet — the next milestone must be chosen with the user.

## Exact Next Action

Choose the next milestone with the user from the remaining roadmap items: a sync/backup workflow
for the SQLite progress database, a local LLM provider for offline AI drafting, or continued
B2/C1 German content expansion. Until one is chosen, continue bug and performance hardening.

## Top Issues

None. The 2026-08-03 bug-hardening pass closed all tracked issues.

## Acceptance Criteria

- `./scripts/verify.sh` passes and completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-05: Go unit tests with `-race`, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 35.36s).

## Repository State

- `refactor/generic-trainer` was fast-forward merged into `main` on 2026-08-05 and both branches
  are pushed. `main` is at `f56400c`.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
