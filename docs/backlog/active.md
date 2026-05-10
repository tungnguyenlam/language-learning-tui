# Active Backlog

Last updated: 2026-05-10 13:31 +07 (Autonomous improvement pass)

## Current Milestone

End-to-end TUI recertification plus multi-area improvement batch.

## Phase 1 Baseline

- Root startup routine, continuity workflow, and TUI/content subtree instructions read.
- `go run ./cmd/deutsch-tui` launched successfully to Dashboard in a TTY on 2026-05-10 13:15 +07.
- Existing uncommitted changes observed and left intact: `README.md`, `demo.tape`, `assets/.gitkeep`, `assets/hero-dark.svg`, `assets/hero-light.svg`.
- `timeout` is not available in this macOS shell; use TTY sessions or project scripts for launch checks.

## Top Issues / Work Plan

- [x] Batch 1: add new German content, improve deck/browser/dashboard review surfaces, and add focused unit/E2E coverage.
- [ ] Batch 2: improve AI/settings/import developer UX, add additional E2E scenarios, and fix regressions found by tests.
- [ ] Batch 3: final polish, AGENTS.md hygiene, full `./scripts/verify.sh`, and commit all changes.

## Exact Next Action

Run the new mobility/UI E2E file, fix any regressions, then continue with Batch 2 polish for import/AI/settings/developer UX.

## Acceptance Criteria

- App starts with zero errors.
- Dashboard, Review, Import, AI, Settings, Browser, Decks, Statistics, and Cram render through Bubble Tea.
- Core flashcard reveal, grading, navigation, and persistence are covered by existing or new tests.
- At least 8 new tui-tester E2E tests are added and passing.
- At least 3 new features/content additions and at least 2 UI/UX improvements are implemented.
- `./scripts/verify.sh` passes with zero errors.
- `docs/backlog/done.md`, continuity docs, and AGENTS.md hygiene are updated.
- Final state is committed to git.

## Last Verification

- 2026-05-10 13:15 +07: `go run ./cmd/deutsch-tui` launched successfully; stopped with `q`.
- 2026-05-10 13:31 +07: `go test ./internal/content ./internal/tui` passed after fixing the mobility TSV quote parsing issue.
