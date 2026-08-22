# Agent Context Index

Last updated: 2026-08-22

## Vision & Prompts

- `GOAL.md`: north-star vision for the app's final state. All work should align with this.
- `prompt/improve.md`: standard autonomous improvement prompt for agent sessions.
- `prompt/debug.md`: autonomous reproduce, diagnose, fix, and regression-test prompt.

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
- `docs/agent/notices/2026-05-12-headered-tsv-decks.md`: headered embedded TSV decks must not leak headers or literal explanations into review answers.
- `docs/agent/notices/2026-05-16-edge-tts-provider.md`: Edge TTS is an optional unofficial online CLI provider and must fail gracefully.
- `docs/agent/notices/2026-06-11-spotlight-dictionary-overlay.md`: Dictionary is a Spotlight-like overlay (`=` key), not a tab. Removed from nav cycle.
- `docs/agent/notices/2026-07-26-trainer-input-vs-global-shortcuts.md`: new single-letter global shortcuts must be guarded with `trainerInputActive()` or they are unreachable as typed characters in the practice trainers.
- `docs/agent/notices/2026-07-26-ankiweb-undocumented-endpoints.md`: AnkiWeb's shared-deck protobuf endpoints have no published schema; decode defensively, refresh download tokens, and keep the browser out of the tab cycle.
- `docs/agent/notices/2026-08-01-dictionary-forms-fts.md`: normal inflection searches use a companion FTS5 table whose rowids must stay aligned with the main dictionary table.
- `docs/agent/notices/2026-08-12-generic-trainer-load-identity.md`: generic trainer item loads must carry a request ID and active-subview guard.
- `docs/agent/notices/2026-08-12-anki-cloze-ordinal-grouping.md`: repeated Anki cloze ordinals form one card and must retain numeric order and grouped answers.
- `docs/agent/notices/2026-08-15-progress-backup-excludes-dictionary.md`: progress backups omit dict.cc tables; TUI must type-assert `BackupRepository`.
- `docs/agent/notices/2026-08-22-cram-session-key-trap.md`: active Cram must consume unhandled keys so global nav cannot abandon the session.
- `docs/agent/notices/2026-08-22-seeded-content-e2e-due-count.md`: seeded-content expansion is coupled to hard-coded `Due cards: 52` / `51 cards due` E2E assertions because `DueCards` loads all decks.

## Verification Status

- ✅ `./scripts/verify.sh` passed on 2026-08-22 after the deletion-first cleanup
  (Go tests, vet, dict.cc import, smoke test, binary build, core E2E suite: 35
  passed in 36.78s)
- ✅ `./scripts/verify.sh` passed on 2026-08-22 after the improvement pass (Go tests, vet, dict.cc
  import, smoke test, binary build, core E2E suite: 35 passed in 38.15s)
- ✅ `./scripts/verify.sh` passed again on 2026-08-22 for the irregular-verb grammar fix + grading/grammar
  test additions (core E2E suite: 35 passed in 37.22s)
- ✅ All Go test suites passing
- ✅ App launches without errors (smoke test passing)
- ✅ Full verification passes with `./scripts/verify.sh`

## Decisions

- `docs/decisions/ADR-0001-local-first-modular-monolith.md`: local-first modular monolith with indexed continuity docs.
- `docs/dependencies.md`: dependency upgrade and verification policy.
- `docs/storage-migrations.md`: SQLite migration policy.
- `docs/release/checklist.md`: release/checkpoint checklist.

## Package Maps

- `internal/core/README.md`: Domain models, interfaces, and architectural boundaries (Pure Go).
- `internal/tui/README.md`: Responsive Bubble Tea shell, View/Model architecture, and mouse routing.
- `internal/storage/README.md`: SQLite migrations, repository pattern, and persistence rules.
- `internal/content/README.md`: TSV format, Anki interop, and embedded starter decks.
- `internal/ankiweb/README.md`: read-only client for AnkiWeb's public shared-deck library (the app's only network surface).
- `internal/ai/README.md`: AI provider adapters, drafting workflow, and validation.
- `internal/srs/README.md`: FSRS scheduling implementation and core mapping.
- `internal/app/README.md`: Config defaults, logging, and data directory setup.
- `tui_tester/README.md`: E2E testing utility and persistent agent mode.
