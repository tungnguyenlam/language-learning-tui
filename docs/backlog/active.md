# Active Backlog

Last updated: 2026-04-29

## Current Milestone

End-to-end deutsch-tui stabilization.

## Next Action

End-to-end stabilization is complete. Resume the next roadmap milestone by implementing the Statistics view unless a newer user request supersedes it.

## Acceptance Criteria

- `cmd/deutsch-tui` launches without errors.
- Dashboard, Review, Import, AI, and Settings views render correctly.
- Core keyboard and mouse interactions respond, including flashcard reveal, grading, and navigation.
- Review state persists to SQLite across app restarts.
- Existing tests and at least 3 tui-tester E2E tests pass.
- `./scripts/verify.sh` passes with zero errors or warnings.
- Completed milestone summary is recorded in `docs/backlog/done.md`.

## Blockers

- None.

## Last Verified

- 2026-04-29: `go test ./...` passed.
- 2026-04-29: `./scripts/tui_smoke.sh` passed.
- 2026-04-29: direct `TUIAgent` launch/render of `cmd/deutsch-tui --data-dir <tmp>` showed the Dashboard.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q` passed with 11 tests.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/test_tui.py -q` passed with 14 tests after adding three fresh E2E tests.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/` passed with 16 tests including WASD navigation support.
- 2026-04-29: `./scripts/verify.sh` passed with 16 E2E tests.
- 2026-04-29: port 8765 was free; manual `tui_tester` daemon start/observe/act/stop against `cmd/deutsch-tui` passed.
- 2026-04-29: `./scripts/verify.sh` passed with 19 E2E tests (including 5 new tests for WASD and robustness).
