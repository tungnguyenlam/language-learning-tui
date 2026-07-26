# Active Backlog

Last updated: 2026-07-26

## Current Milestone

Practice Trainer Quality.

## Exact Next Action

Select the next milestone from `docs/backlog/roadmap.md` (Backup workflow, `.apkg`
import/export) or continue the screen-interface migration started on
`refactor/generic-trainer` — `internal/tui/screen.go` documents the pattern and
only ViewImport / ViewSessionSummary are migrated so far.

## Top Issues

None active.

## Completed Work

- [x] **Trainer input no longer quits the app:** `q` and `?` are consumed as
  typed characters while a text-input trainer is waiting for an answer.
  German answers contain `q` (Qualität, Quelle), so typing one used to hit the
  global quit shortcut and exit mid-exercise. See `trainerInputActive()` in
  `internal/tui/keys.go`.
- [x] **Trainers reshuffle between passes:** `trainerState.advance()` and
  `Model.advanceGenderItem()` reshuffle after each completed pass, and the
  Gender Trainer's noun list is shuffled at load. Small exercise sets used to
  cycle in a fixed order forever, so learners recalled the position instead of
  the grammar. The first pass keeps its authored order.
- [x] **Trainer progress indicator:** headers now show `Item n/N` (`Noun n/N`
  for the Gender Trainer) and a `Round n` counter after the first pass.
- [x] **Grammar content fixes and expansion:** corrected three wrong grammar
  labels (`wegen des Wetters` is neuter not masculine; `in einem alten Haus` is
  mixed not weak declension; `warten auf` is a fixed case, not movement) and
  grew the Case and Adjective Ending trainers from 15 to 25 exercises each.
- [x] **Help overlay accuracy:** documented the `=` Dictionary spotlight, which
  was missing entirely, and replaced the stale `Practice 1-3/d-a/m-n` line with
  the real `1-9` trainer keys.

## Verification

- `./scripts/verify.sh` passed: gofmt, all Go unit tests, `go vet`, smoke test,
  binary build, and the E2E suite (353 tests, including the two new
  `e2e_tests/test_trainer_input_shortcuts.py` cases).
