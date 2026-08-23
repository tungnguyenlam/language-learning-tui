# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Bug-fix and performance pass (per the focus in `AGENTS.md`): Cram grade mouse
targets now follow the actual rendered/wrapped controls, database reset reports
post-reset reload failures, and related-dictionary queries no longer suppress
SQLite failures.

## Exact Next Action

No unfinished executable work remains. Candidate future work:

- Extract remaining view-local state off `Model` (screen files own render +
  keys only today).
- The `render_dictionary.go` rescan cleanup (the last render file still using
  `strings.Count(b.String(), "\n")`) is deprioritized: its `b`/`detailBuilder`
  buffers hold only a small header/detail, so the O(n²)→O(n) gain is marginal
  and the two-builder layout makes the refactor higher-risk than the other four
  files.

## Completed This Pass

- Bug fix: Cram's grade hitboxes assumed a centered card and counted source
  newlines, but the card starts at the render context X and Lip Gloss may wrap
  long prompt/answer/context text. Targets are now located from the actual
  rendered card and remain aligned when controls wrap across rows.
- Bug fix: `executeResetDatabase` ignored errors while reloading Decks and Due
  Cards after a successful reset, falsely reporting a clean success with stale
  UI state. Both reload failures now surface with operation context.
- Bug fix: `FindRelatedEntries` discarded errors from its prefix, suffix, and
  exact-match queries, turning database failures into empty related-word lists.
  Query failures now propagate with branch-specific context.
- Regression tests cover all four wrapped Cram grade targets, both reset reload
  paths, and long/short related-entry lookup failures against a closed database.
- `docs/backlog/done.md` updated continuously.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Cram's Again/Hard/Good/Easy mouse targets align with their visible labels even
  when a narrow terminal wraps long card content or the Grade controls.
- Database reset never reports success when its post-reset Decks or Due Cards
  reload fails.
- Related-entry lookup returns SQLite failures rather than an empty result.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `go test ./...` passed on 2026-08-23 after the Cram hitbox, reset error, and
  related-dictionary error fixes.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.51s).

## Repository State

- All session work is verified and committed; the working tree is clean.
