# Agent Continuity Workflow

Use this workflow to make development continuous across chat sessions without dumping everything into `AGENTS.md`.

## New Session Pickup

At the start of a new session:

1. Read root `AGENTS.md`.
2. Read `docs/backlog/active.md` and continue the `Next Action` unless the user gives a newer instruction.
3. Read `docs/agent/index.md` and any active notices linked there.
4. Read only the package README and nearest subtree `AGENTS.md` for files you will touch.
5. Use `rg` to inspect the relevant symbols before editing.

If the user switches chats mid-task, treat `docs/backlog/active.md` as the handoff source of truth.

## Status Note

As of 2026-05-03: Project is production-ready. All 77 E2E tests pass, all Go unit test suites pass, verification passes end-to-end. Run 2 TUI refactor recovery is complete and no outstanding issues are tracked.

## Continuous Documentation

Agents MUST NOT wait until the very end of a session to document changes. As soon as a logical sub-task or feature is completed, or a plan changes:

1. Instantly update `docs/backlog/active.md` to reflect the current state and exact next action.
2. If new knowledge was gained, immediately add a notice or update an ADR.
3. Record completed tasks in `docs/backlog/done.md`.

Documenting as you go prevents context loss if a session ends abruptly and ensures the user stays informed of the agent's progress.

## Notices

Write a notice when future agents need durable context that is not executable backlog:

- a known trap or failed approach
- a subtle invariant
- a temporary compatibility rule
- a cross-package contract
- a decision that is too small for an ADR but too important to lose

Create notices in `docs/agent/notices/` named:

```text
YYYY-MM-DD-short-topic.md
```

Use this shape:

```md
# Short Topic

Status: active | resolved | expires YYYY-MM-DD
Scope: package, feature, or file area
Related: paths or symbols

## Why It Matters

One short paragraph.

## Required Behavior

What future agents must do or avoid.

## Revisit When

Condition that makes this notice stale.
```

After adding, resolving, or deleting a notice, update `docs/agent/index.md`.

## Backlog

Use `docs/backlog/active.md` for unfinished executable work. It must answer:

- current milestone
- exact next action
- acceptance criteria
- blockers
- last verified command
- likely files or packages involved when useful

Use `docs/backlog/roadmap.md` for committed future milestones.
Use `docs/backlog/parking-lot.md` for ideas that are not committed work.
Use `docs/backlog/done.md` for completed milestone summaries and verification.

Do not put vague reminders in the active backlog. Convert them into a concrete next action or move them to the parking lot.

## Subtree AGENTS.md

Subtree `AGENTS.md` files are allowed for major areas with local rules. Keep them short and scoped.

Use subtree instructions for:

- package-specific test requirements
- local architecture boundaries
- UI/input rules
- migration or fixture rules
- provider or network restrictions

Do not duplicate root workflow in subtree files. If a tool does not auto-load nested `AGENTS.md`, root `AGENTS.md` and `docs/agent/start-here.md` still require agents to read the nearest subtree file manually.

## Handoff

Before ending work, update continuity state:

1. If work remains, update `docs/backlog/active.md` with the exact next action.
2. If work completed, add a short entry to `docs/backlog/done.md`.
3. If a durable warning was discovered, add or update a notice.
4. If architecture changed, add or update an ADR.
5. If notices or decisions changed, update `docs/agent/index.md`.
6. Run targeted tests and `./scripts/verify.sh` when feasible, then record the latest relevant verification command.

Do not leave important context only in chat.
chat.
