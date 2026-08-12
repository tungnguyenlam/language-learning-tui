# Active Backlog

Last updated: 2026-08-12

## Current Milestone

Milestone 4 (Hybrid Dictionary) remains complete. The 2026-08-12 autonomous improvement pass is
complete; no new feature milestone is committed yet.

## Exact Next Action

Choose the next milestone with the user from sync/backup, local LLM drafting, or B2/C1 content;
until then continue bug and performance hardening.

## Top Issues

None. The three identified interaction races/state leaks are implemented with tests.

## Acceptance Criteria

- `./scripts/verify.sh` passes and completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-05: Go unit tests with `-race`, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 37.59s).
- Baseline `./scripts/verify.sh` passed on 2026-08-12 before this improvement pass (34 core E2E tests passed).
- Final `./scripts/verify.sh` passed on 2026-08-12: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 36.56s).
- Bug hunting & performance optimization pass complete: `cleanWhitespace` fast-path rune scanner (`dictcc.go`), `addDictionaryClozeEntryCmd` bareWord regex safety guard (`actions_dictionary.go`), and `renderCram` windowed visible cards list rendering (`render_cram.go`).

## Repository State

- `refactor/generic-trainer` was fast-forward merged into `main` on 2026-08-05 and both branches
  are pushed. `main` is at `f56400c`.
- The five `subagent-*` branches are fully merged into `main` but remain checked out in external
  worktrees under `~/.gemini/antigravity-cli/`. Delete those worktrees before deleting the branches.
