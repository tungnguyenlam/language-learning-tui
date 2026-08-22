# Seeded Content Expansion Coupled to Hard-coded E2E Due-count

Status: active
Scope: internal/content (embedded TSV/Go decks), e2e_tests, internal/tui dashboard/review
Related: internal/storage/sqlite/sqlite.go `DueCards(ctx, now, limit)`, `internal/tui/loaders.go` `loadDueCards`, e2e_tests assertions `Due cards:   52` / `51 cards due`

## Why It Matters

`Store.DueCards(ctx, time.Time, limit)` ignores deck identity (only `now` and
`limit` are parameters). On a fresh database every seeded card is "due", so the
dashboard `Due cards:   %d` line and the Review status `%d cards due in <deck>`
reflect the TOTAL seeded card count. The core E2E suite hard-codes this number
(`Due cards:   52`, `51 cards due`) and will fail if the seeded total changes.

## Required Behavior

When adding or removing seeded decks/notes, either:
- update every hard-coded due-count assertion in `e2e_tests/` to the new total
  (dashboard shows the new total; the review status shows total−1 after one
  card is revealed), confirm with `./scripts/verify.sh`, and keep the numbers
  in sync; or
- make the assertions count-agnostic (wait for "cards due" substring rather than
  an exact number) so content churn stops breaking the suite.

Do NOT silently add seeded content expecting the suite to stay green — it will
not, and the failure surfaces as a count mismatch deep in the E2E run.

## Revisit When

`DueCards` gains a deck filter parameter, or the E2E assertions are made
count-agnostic, removing the coupling.
