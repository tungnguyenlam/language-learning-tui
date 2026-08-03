# Active Backlog

Last updated: 2026-08-03

## Current Milestone

Dictionary Feature Enhancements & Lemmatization Search.

## Exact Next Action

Pick the next dictionary milestone item, or continue bug/perf hardening if regressions appear.

## Top Issues

No active issues remain from the bug fixing pass (UTF-8 cloze, SQL aggregates, empty store & intra-day SRS).

## Acceptance Criteria

- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-03: Go unit tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed).

## Completed Work

- [x] Fixed UTF-8 byte offset corruption in cloze card creation (`addDictionaryClozeEntryCmd`).
- [x] Fixed SQL aggregation ordering in `RecentDecks`.
- [x] Fixed empty-store division-by-zero guard in `RandomEntries`.
- [x] Fixed 1000-row streak calculation cap in `currentStreak` / `deckCurrentStreak`.
- [x] Fixed 0-second intra-day learning intervals in `stateFromFSRS`.
