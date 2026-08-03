# Active Backlog

Last updated: 2026-08-03

## Current Milestone

Dictionary Feature Enhancements & Bug Hardening.

## Exact Next Action

Pick the next dictionary milestone item, or continue bug/perf hardening if regressions appear.

## Top Issues

No active issues remain from the bug fixing pass (Map data race in searchDictionary, UTF-8 string indexing in extractPlural, SQL rows leak in readAnkiDeckNames, unclamped deckCursor in selectDeckByID/selectDeck, Gender/Generic trainer OOB bounds, deck merge unclamped cursor, and SQLite transaction row closure).

## Acceptance Criteria

- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `./scripts/verify.sh` passed on 2026-08-03: Go unit tests with `-race`, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed in 36.50s).

## Completed Work

- [x] Fixed map data race in `searchDictionary()` (`actions_dictionary.go`) by snapshotting `dictionaryStarred` on main UI thread before async command execution.
- [x] Fixed unsafe string indexing in `extractPlural()` (`loaders.go`) by matching case-insensitively directly on `extra` byte boundaries, preventing UTF-8 byte corruption with multi-byte characters like `ẞ`.
- [x] Fixed database `rows` leak in `readAnkiDeckNames()` (`apkg_import.go`) by closing `rows` explicitly when scanning completes.
- [x] Fixed unclamped `deckCursor` panic in `selectDeckByID()` and `selectDeck()` (`handlers.go`) when selecting decks while a filter is active in Decks view.
- [x] Fixed out-of-bounds index panics in Gender Trainer (`render_gender_trainer.go` and `keys.go`) when `practiceItems` is empty or `practiceIndex` is invalid.
- [x] Fixed generic trainer index out of bounds (`trainer.go`) in `renderTrainer` and `updateTrainerKey`.
- [x] Fixed unclamped `deckCursor` panic in `mergeSelectedDecks` (`handlers.go`) and `updateDecksKey` (`keys.go`).
- [x] Fixed unclosed `rows` statement before transaction `ExecContext` in `SaveReview` (`sqlite.go`) and added missing `rows.Err()` checks in `SaveReview` and `Statistics`.
- [x] Fixed AI draft Note ID collision for German words with umlauts/eszett in `draftIDBase` (`ai.go`).
- [x] Fixed missing lower-bound (`>= 0`) cursor checks in `render_browser.go`, `render_cram.go`, and `render_ai.go`.


