# Active Backlog

Last updated: 2026-08-03

## Current Milestone

Dictionary Feature Enhancements & Bug Hardening.

## Exact Next Action

Pick the next dictionary milestone item, or continue bug/perf hardening if regressions appear.

## Top Issues

No active issues remain from the bug fixing pass (Gender trainer empty bounds, generic trainer bounds, deck merge unclamped cursor, SQLite review transaction row closure, AI draft German umlaut ID collision, and lower bound cursor checks).

## Acceptance Criteria

- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-03: Go unit tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 36.97s).

## Completed Work

- [x] Fixed out-of-bounds index panics in Gender Trainer (`render_gender_trainer.go` and `keys.go`) when `practiceItems` is empty or `practiceIndex` is invalid.
- [x] Fixed generic trainer index out of bounds (`trainer.go`) in `renderTrainer` and `updateTrainerKey`.
- [x] Fixed unclamped `deckCursor` panic in `mergeSelectedDecks` (`handlers.go`) and `updateDecksKey` (`keys.go`).
- [x] Fixed unclosed `rows` statement before transaction `ExecContext` in `SaveReview` (`sqlite.go`) and added missing `rows.Err()` checks in `SaveReview` and `Statistics`.
- [x] Fixed AI draft Note ID collision for German words with umlauts/eszett in `draftIDBase` (`ai.go`).
- [x] Fixed missing lower-bound (`>= 0`) cursor checks in `render_browser.go`, `render_cram.go`, and `render_ai.go`.

