# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 improvement pass fixed `-en` adjective grammar hints and expanded
thin unique Standard Content decks. StarterDeck size is unchanged. Verification
is fully green.

## Exact Next Action

No unfinished executable work remains. Candidate future work:
- Make E2E due-count assertions count-agnostic so `StarterDeck()` can grow.
- Extract remaining view-local state off `Model`, or split the large
  async-message switch in `internal/tui/model.go`.

## Completed This Pass

- Grammar hints: adjectives ending in `-en`/`-ern` (`selten`, `trocken`,
  `offen`, `zufrieden`, `nüchtern`, …) are no longer classified as infinitives.
- Expanded news, emergency, proverbs, preposition, and environment TSV decks;
  fixed `Jetzt oder nie`, `Schweigen ist Gold`, `erneuerbare Energien`, and
  the fake proverb `Pech im Glück, Glück im Unglück`.
- Clarified that only `StarterDeck()` auto-seed is coupled to `Due cards: 52`.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- View-local state still lives on `Model`; screen files own render + keys only.
- StarterDeck expansion is still coupled to hard-coded E2E due-count assertions.

## Acceptance Criteria

- Grammar hints for `-en` adjectives show comparison forms, not fake conjugations.
- Thin unique Standard Content decks have roughly 40 notes each.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `go test ./...` passed on 2026-08-22.
- `./scripts/verify.sh` passed on 2026-08-22: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 38.55s).

## Repository State

- Improvement pass is complete and fully verified on this working tree.
