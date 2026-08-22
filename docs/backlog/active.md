# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 improvement pass fixed irregular-verb grammar-hint classification
(`sein`/`tun`) and added grading-logic + grammar-engine unit tests. Verification
is fully green.

## Exact Next Action

No unfinished executable work remains. Candidate future work:
- Expand the thin B2/C1 decks (`b2-c1-news`, `b2-environment`) — but this must
  update the hard-coded `Due cards: 52` / `51 cards due` E2E assertions, since
  `DueCards(ctx, now, limit)` loads all seeded decks with no filter.
- Extract view-local state off `Model` (Settings, Cram, AI, Practice,
  Dictionary keys now live on screens); or split the large async-message switch
  in `internal/tui/model.go`.

## Completed This Pass

- `isInfinitive` now recognizes the irregular infinitives `sein` and `tun`;
  added them to the curated `dailyVerbs` table with correct conjugations so
  Review grammar hints and Verb-of-the-Day show accurate forms.
- New `internal/tui/grade_logic_test.go` covers `normalizeAnswer`,
  `clozeAnswerText`, `renderClozeAnswers`, and `renderTypingDiff`.
- Extended `internal/content/wordinfo_test.go` with irregular-verb, umlaut,
  reverse-card, and English-prefix edge cases.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- View-local state still lives on `Model`; screen files own render + keys only.
- Seeded-deck content expansion is coupled to hard-coded E2E due-count
  assertions (see Exact Next Action).

## Acceptance Criteria

- `sein`/`tun` classify as verbs and show correct conjugations in Review and
  Verb-of-the-Day.
- `normalizeAnswer` folds umlauts/ß and trims punctuation in non-strict mode.
- Grade-logic and grammar-engine tests pass; `./scripts/verify.sh` stays green.

## Last Verification

- `go test ./...` passed on 2026-08-22.
- `./scripts/verify.sh` passed on 2026-08-22: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 37.22s).

## Repository State

- Improvement pass is complete and fully verified on `main`.
