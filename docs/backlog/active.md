# Active Backlog

Last updated: 2026-08-01

## Current Milestone

Dictionary Feature Enhancements & Lemmatization Search.

## Exact Next Action

Pick the next dictionary milestone item after this Spotlight/input/browse polish pass.

## Top Issues

- [x] Spotlight `a` no longer steals search typing when prior results remain visible.
- [x] Dictionary detail scroll uses painted viewport height (single-column no longer under-scrolls).
- [x] Filter-only browse (`:verb`, `:pl`, …) queries matching rows in SQL instead of the first 200 alphabetical entries.
- [x] Spotlight empty state respects panel body budget; generic `Ready` status no longer wipes footer shortcuts.

No active issues remain from this pass.

## Acceptance Criteria

- Typing `a` in Spotlight search appends to the query even when prior results exist; focused results still play audio.
- Single-column detail shift/wheel/scrollbar scroll range matches rendered visible rows.
- `:verb` / gender pills return matches from the full dictionary, not only early alphabetical rows.
- Short-terminal Spotlight empty state keeps footer hints and truncates secondary lists.
- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `go test ./internal/tui/ ./internal/storage/sqlite/`
- `./scripts/verify.sh` passed on 2026-08-01: Go unit tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed).

## Completed Work

- [x] **Spotlight Search Types `a`:** Audio shortcut requires results/detail focus; overlay-alone no longer intercepts typing (*Auto*, *haben*).
- [x] **Dictionary Detail Viewport Cache:** Paint caches `dictionaryDetailVisibleRows` / `dictionaryListVisibleRows`; keyboard, wheel, drag, and scrollbar clicks use the painted sizes.
- [x] **Filter-Only Browse Completeness:** Class/gender filter-only search pushes `WHERE` predicates into SQL with `ORDER BY word` + limit.
- [x] **Spotlight Empty-State Budget & Footer:** Empty overlay sections stop at the body budget; mundane `Ready` status no longer replaces shortcut hints when the footer is tight.

## Verification

- `./scripts/verify.sh` passed on 2026-08-01 after Spotlight typing, detail scroll, filter-only browse, and empty-state footer fixes.
