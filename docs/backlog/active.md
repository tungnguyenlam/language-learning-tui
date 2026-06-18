# Active Backlog

Last updated: 2026-06-18

## Current Milestone

Content Expansion & Refinement.

## Exact Next Action

Select next milestone from `docs/backlog/roadmap.md` (e.g., Backup workflow or more visual refinements) or continue adding content.

## Top Issues

None active.

## Completed Work

- [x] **A1 Content Expansion:** Added A1 Clothing Deck and A1 Time & Weather Deck to `german_expanded.go` covering essential topics.
- [x] **E2E Stability:** Fixed PTY truncation/rendering glitch in `test_may15_improvements.py` caused by string manipulation affecting UI layout tests.
- [x] **E2E Path Robustness:** Fixed the `app_cmd` binary paths in all E2E pytest files, removing `../cmd/` from `go run ../cmd/deutsch-tui` to prevent `directory not found` exceptions when running individual test files from the project root.

## Verification

- `./scripts/verify.sh` passed: all Go unit tests, smoke test, binary build, and E2E tests.
