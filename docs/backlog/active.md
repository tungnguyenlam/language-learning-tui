# Active Backlog

Last updated: 2026-08-01

## Current Milestone

Dictionary Feature Enhancements & Lemmatization Search.

## Exact Next Action

Pick the next dictionary milestone item, or continue bug/perf hardening if regressions appear.

## Top Issues

No active issues remain from the Practice Hub filter / due-queue / Cram stale-load pass.

## Acceptance Criteria

- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-01: Go unit tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed).

## Completed Work

- [x] Practice Hub `/` filter traps globals and accepts paste.
- [x] Due-queue / Cram / Gender practice loads ignore stale async results.
- [x] Grade/undo/suspend that finish after a bookmark-filter flip reload the live filter.
