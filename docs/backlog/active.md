# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 grammar-hint pass classifies numbers and adverbs instead of
faking adjective/verb forms, teaches irregular comparatives, and corrects
weather plus thin unique Standard Content decks. StarterDeck size is unchanged.
Verification is fully green.

## Exact Next Action

No unfinished executable work remains. Candidate future work:
- Make E2E due-count assertions count-agnostic so `StarterDeck()` can grow.
- Extract remaining view-local state off `Model`, or split the large
  async-message switch in `internal/tui/model.go`.

## Completed This Pass

- Grammar hints: cardinal numbers (`sieben`, `dreizehn`, …) and common adverbs
  (`gestern`, `heute`, `oben`) are no longer treated as infinitives or given
  fake comparatives (`ich sieb`, `heuteer`).
- Irregular adjective comparison: `gut → besser`, `hoch → höher`, `alt → älter`,
  `dunkel → dunkler`, and similar high-frequency forms.
- Weather TSV: `schneien` is "to snow"; phrase cards are German-front;
  time-of-day nouns have articles; adjectives/adverbs/verbs are lowercase.
- Split ambiguous `morgen` = Tomorrow/Morning; added teens 13–19; expanded
  `a1-essential`, `b1-workplace-office`, and `b2-urban-mobility` to ~40 notes.
- Weather E2E search now uses `Seasons` so it does not select Time & Weather.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- View-local state still lives on `Model`; screen files own render + keys only.
- StarterDeck expansion is still coupled to hard-coded E2E due-count assertions.

## Acceptance Criteria

- Grammar hints for numbers/adverbs do not show conjugations or naive comparatives.
- Irregular adjectives show the real comparative, not `guter` / `hocher`.
- `schneien` is the weather verb for snow; `StarterDeck()` due count stays 52.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `go test ./...` passed on 2026-08-22.
- `./scripts/verify.sh` passed on 2026-08-22: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 36.66s). Extra: `test_grammar_hint.py` (2 passed) and
  `test_weather_deck.py` (1 passed after unique search).

## Repository State

- Improvement pass is complete and fully verified on this working tree.
