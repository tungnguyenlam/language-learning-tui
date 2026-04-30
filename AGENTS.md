# Agent Instructions

Use this file at the start of every new chat session.

## Current State

The project is a Go Bubble Tea TUI for German flashcard and MCQ learning. It is local-first, uses SQLite for progress, uses editable deck files for content, and targets Anki-friendly TSV import/export before `.apkg`.

## Startup Routine

1. Read `docs/backlog/active.md` for the current milestone and next action.
2. Read `docs/agent/index.md` for durable notices and decisions.
3. Read `docs/agent/continuity.md` for the required pickup and handoff workflow.
4. Use `rg` to find only the package or feature area you need.
5. Read the package README and nearest subtree `AGENTS.md` before editing that area.
6. **Continuous Documentation:** You MUST update `docs/backlog/active.md` and relevant continuity docs *continuously* as you complete tasks, change plans, or discover new context, not just at the end of the session.
7. Update backlog and notices when leaving unfinished work or changing contracts.

## Instruction Maintenance

- Agents may update these instructions, subtree `AGENTS.md` files, and continuity docs when a change would make future work safer or easier.
- Keep root `AGENTS.md` concise. Put durable traps, workflow notes, and package-specific contracts in `docs/agent/notices/`, subtree `AGENTS.md`, or ADRs, then index them from `docs/agent/index.md`.
- When modifying instructions, make the note concrete and searchable: include the affected files, the required behavior, and when the note should be revisited.
- Do not preserve stale instructions out of inertia. If an instruction becomes wrong, update or remove it in the same change that makes it obsolete.

## Search Hints

- TUI behavior: `rg "type Model|Hitbox|WindowSizeMsg|Mouse" internal/tui`
- Scheduling: `rg "Scheduler|Review" internal/core internal/srs`
- Storage: `rg "Repository|migration|CREATE TABLE" internal/storage`
- Anki import/export: `rg "TSV|Anki|Import|Export" internal/content`
- AI drafting: `rg "Provider|Draft" internal/ai`

## Handoff Rules

- Put unfinished executable work in `docs/backlog/active.md`.
- Put durable warnings or traps in `docs/agent/notices/`.
- Put long-lived architecture decisions in `docs/decisions/`.
- Keep all continuity notes concise and indexed from `docs/agent/index.md`.
- Do not dump chat history into `AGENTS.md`; write searchable, task-shaped records instead.
- Use `./scripts/verify.sh` as the full project verification command.
