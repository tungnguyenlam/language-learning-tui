# Active Backlog

Last updated: 2026-04-29

## Current Milestone

End-to-end deutsch-tui stabilization.

## Next Action

End-to-end stabilization is complete and committed. Resume the next roadmap milestone by implementing the Statistics view unless a newer user request supersedes it.

## Blockers

- None.

## Last Verified

- 2026-04-29: `./scripts/verify.sh` passed with 22 E2E tests, all Go tests, smoke test, gofmt, and go vet.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/test_recertification.py -q` passed with 3 new E2E tests.
- 2026-04-29: Added three recertification E2E tests for Tab view cycling, Hard grade SQLite persistence, and Settings provider persistence.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests -q` passed with 19 E2E tests.
- 2026-04-29: `./scripts/tui_smoke.sh` passed; `go test ./...` passed.
- 2026-04-29: `./scripts/verify.sh` passed with 19 E2E tests, all Go tests, smoke test, gofmt, and go vet.
- 2026-04-29: Cleaned up debug files and added test_data/ to .gitignore.
