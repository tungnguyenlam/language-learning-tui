# Seeded Content Expansion Coupled to Hard-coded E2E Due-count

Status: resolved
Scope: StarterDeck auto-seed, e2e_tests, dashboard/review due counts
Related: `internal/content/starter.go`, `cmd/deutsch-tui/main.go`, `e2e_tests/e2e_helpers.py`

## Why It Matters

On an empty database the app auto-seeds only `StarterDeck()`. `DueCards` then
returns every due card in that collection, so the dashboard `Due cards:   %d`
line and Review `Review 1/%d` / `%d cards due` reflect the **starter** total.

`StandardDecks()` TSV/Go decks are **not** auto-seeded. Expanding those files
does not change the starter due count on a fresh DB. Tests that press `S` to
seed standard content already wait for deck names, not the starter count.

## Required Behavior

- E2E assertions are count-agnostic: tests read the live count with
  `read_due_count(agent)` / `read_cards_found(agent)` from
  `e2e_tests/e2e_helpers.py` and compute expectations (e.g. `due - 1` after one
  grade). New tests must use these helpers instead of hard-coding numbers.
- Expanding Standard Content (embedded TSV, extra Go decks): no E2E count
  changes needed.
- `StarterDeck()` may now grow without touching E2E assertions; the helpers
  read the seeded count from the dashboard at runtime.

## Revisit When

`DueCards` gains a deck filter, or the dashboard stops showing a
`Due cards:   N` line (the helpers parse it and fail loudly if absent).
