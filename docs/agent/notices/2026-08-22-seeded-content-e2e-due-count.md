# Seeded Content Expansion Coupled to Hard-coded E2E Due-count

Status: active
Scope: StarterDeck auto-seed, e2e_tests, dashboard/review due counts
Related: `internal/content/starter.go`, `cmd/deutsch-tui/main.go`, `e2e_tests` assertions `Due cards:   52` / `Review 1/52` / `51 cards due`

## Why It Matters

On an empty database the app auto-seeds only `StarterDeck()`. `DueCards` then
returns every due card in that collection, so the dashboard `Due cards:   %d`
line and Review `Review 1/%d` / `%d cards due` reflect the **starter** total
(currently 52). The core E2E suite hard-codes that number.

`StandardDecks()` TSV/Go decks are **not** auto-seeded. Expanding those files
does not change `Due cards: 52` on a fresh DB. Tests that press `S` to seed
standard content already wait for deck names, not the starter count.

## Required Behavior

- Changing `StarterDeck()` size: update every hard-coded `52` / `51` / `Review 1/52` assertion, or make those assertions count-agnostic (parse the number and assert it decreases after a grade).
- Expanding Standard Content (embedded TSV, extra Go decks): do **not** bump the starter due-count assertions.
- Do not silently grow `StarterDeck()` expecting the suite to stay green.

## Revisit When

`DueCards` gains a deck filter, the E2E assertions become count-agnostic, or
`StarterDeck()` is expanded.
