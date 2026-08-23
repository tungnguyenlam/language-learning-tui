# Active Backlog

Last updated: 2026-08-23

## Current Milestone

Bug-fix and performance pass (per the focus in `AGENTS.md`): render hot paths
now track line counts incrementally instead of rescanning whole buffers, the
Review typing mode's revealed answer and its grade hitboxes are fixed, and
answer normalization no longer mis-grades multi-space input.

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

- Performance: incremental line counters replaced `strings.Count(b.String(),
  "\n")` buffer rescans in `render_settings.go`, `screen_import.go`,
  `render_views.go`, and `screen_ankiweb.go` (output byte-identical).
- Bug fix: Review typing mode's revealed answer box was wider than its card
  (`Width(maxInt(30, width-6))` > capped `cardWidth`), so its border wrapped
  across lines. Now `Width(cardWidth - 2)`.
- Bug fix: typing-mode grade hitboxes ignored the typing box's border+padding,
  landing 1 row up and 2 columns left of the rendered Grade row; they now offset
  by the style metrics.
- Bug fix: `normalizeAnswer` collapsed spaces with a single `ReplaceAll` pass,
  leaving 3+ space runs intact and mis-grading answers; now collapses all
  whitespace via `strings.Fields`.
- Regression test `TestTypingModeGradeHitboxAlignsWithRenderedGradeLine` plus new
  whitespace cases in `TestNormalizeAnswerNonStrict`.
- `docs/backlog/done.md` updated.

## Top Issues

- View-local state still lives on `Model`; screen files own render + keys only.

## Acceptance Criteria

- Statistics, streaks, and per-day activity all agree on the local calendar
  day in any timezone (unchanged from prior pass).
- Render output is byte-identical after the incremental-line-count refactors.
- Typing-mode Review renders its answer box inside the card and its grade
  click targets align with the displayed Grade row.
- `go test ./...` and `./scripts/verify.sh` stay green.

## Last Verification

- `go test ./...`, `go vet ./...`, and `gofmt` all pass on 2026-08-23 after the
  render rescan cleanup + typing Review fixes.
- `./scripts/verify.sh` passed on 2026-08-23: Go tests, vet, offline dict.cc
  import (834,512 entries), smoke test, binary build, core E2E suite (35 passed
  in 37.62s).

## Repository State

- All session work verified on this working tree; ready to commit.


