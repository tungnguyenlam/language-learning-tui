# Active Backlog

Last updated: 2026-08-22

## Current Milestone

The 2026-08-22 deletion-first cleanup removed unused theme/keymap config, flattened
the content registry into plain functions, and deleted unused wrappers.
Verification is fully green.

## Exact Next Action

No unfinished executable work remains. Candidate future work:
- Expand the thin B2/C1 decks (`b2-c1-news`, `b2-environment`) — but this must
  update the hard-coded `Due cards: 52` / `51 cards due` E2E assertions, since
  `DueCards(ctx, now, limit)` loads all seeded decks with no filter.
- Extract remaining view-local state off `Model`, or split the large
  async-message switch in `internal/tui/model.go`.

## Completed This Pass

- Removed fake Color Theme cycling and the unused `keymap` config field.
- Flattened ContentSource/Registry into `AllDecks`/`DeckByID` functions.
- Deleted unused TTS `ProviderName`, `scrollbarLineWidth`, logger
  `Warn`/`Fatal`/`String`, and one-line wrappers.
- Deduplicated OpenAI/Anthropic HTTP and shared draft-request helpers.

## Top Issues

- `internal/tui/model.go` still contains the large central async-message switch.
- View-local state still lives on `Model`; screen files own render + keys only.
- Seeded-deck content expansion is coupled to hard-coded E2E due-count
  assertions (see Exact Next Action).

## Acceptance Criteria

- Settings no longer advertises a theme that does not apply.
- `content.AllDecks` / `DeckByID` still resolve filename-keyed embedded TSV decks.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `go test ./...` passed on 2026-08-22.
- `./scripts/verify.sh` passed on 2026-08-22: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, and core E2E suite (35
  passed in 36.78s).

## Repository State

- Cleanup pass is complete and fully verified on this working tree.
