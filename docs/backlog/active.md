# Active Backlog

Last updated: 2026-05-17

## Current Milestone

Help overlay & code cleanup - COMPLETED.

## Completed Work

- [x] **Session Summary Greeting:** Changed "No cards due" to "Session complete!" for a more encouraging tone.
- [x] **Help Overlay Keybindings:** Added missing shortcuts (`d` dictionary, `t` typing, `x` suspend, `h` hint, `F` fix, `!`/`@`/`#` recent decks, `B`/`X` bulk actions) to the help overlay.
- [x] **pytest e2e Mark:** Registered `@pytest.mark.e2e` in `pytest.ini` to eliminate warnings across all E2E test files.
- [x] **Code Cleanup:** Removed empty `var ()` block from `utils.go`.

## Exact Next Action

None - milestone complete.

## Top Issues / Priorities

None.

## Last Verified

- `go build ./...` passed.
- `go test ./...` passed.
- `gofmt` and `go vet` passed.
- Smoke test passed.
- Targeted E2E tests passed (14/14, 2/2).

## Blockers

None.
