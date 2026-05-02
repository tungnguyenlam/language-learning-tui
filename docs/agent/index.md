# Agent Context Index

Last updated: 2026-05-02

## Active Backlog

- `docs/backlog/active.md`: current milestone, next task, verification state, blockers.
- `docs/agent/continuity.md`: required workflow for session pickup, notices, backlog writing, subtree `AGENTS.md`, and handoff.
- `./scripts/verify.sh`: full project verification command.

## Notices

- `docs/agent/notices/2026-04-29-bubbletea-v2-space-key.md`: Bubble Tea v2 uses "space" string for the space bar.
- `docs/agent/notices/2026-04-29-wasd-navigation.md`: WASD key navigation support implementation.
- `docs/agent/notices/2026-04-30-tui-status-line-stability.md`: keep dynamic TUI status text on a stable single-line `status:` surface for tui-tester.
- `docs/agent/notices/2026-04-30-browser-deck-reload.md`: Browser `[`/`]` deck switching must reload Browser cards for the selected deck.
- `docs/agent/notices/2026-05-02-arrow-key-navigation-fix.md`: Arrow key navigation fix for Review view.
- `docs/agent/notices/2026-05-02-scrollbar-hitbox-layout.md`: derive active-panel hitboxes from Lip Gloss frame metrics instead of per-view offsets.
- `docs/agent/notices/2026-05-02-review-reveal-sync.md`: only show grading options when card is fully revealed to prevent race conditions.
- `docs/agent/notices/2026-05-02-parallel-e2e-tests.md`: E2E tests are parallelized with pytest-xdist for 8x speedup.

## Verification Status

- ✅ All 76 E2E tests passing (as of 2026-05-03, including stability fix for parallel execution)
- ✅ All Go test suites passing (as of 2026-05-03)
- ✅ App launches without errors (smoke test passing)
- ✅ All views render correctly (Dashboard, Review, Import, AI, Settings, Browser, Cram)
- ✅ Core user interactions respond as expected (flashcard reveal, grading, navigation)
- ✅ State successfully persisted to SQLite
- ✅ Code formatting verified (gofmt)
- ✅ Code quality verified (go vet)

## Decisions

- `docs/decisions/ADR-0001-local-first-modular-monolith.md`: local-first modular monolith with indexed continuity docs.
- `docs/dependencies.md`: dependency upgrade and verification policy.
- `docs/storage-migrations.md`: SQLite migration policy.
- `docs/release/checklist.md`: release/checkpoint checklist.

## Package Maps

- `internal/tui/README.md`: Bubble Tea shell, responsive layout, input routing.
- `internal/storage/README.md`: SQLite migrations and repository rules.
- `internal/content/README.md`: deck files and Anki text interop.
- `internal/ai/README.md`: AI provider adapter and draft validation.
