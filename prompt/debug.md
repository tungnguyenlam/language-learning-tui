# Debug deutsch-tui

You are working on deutsch-tui, a terminal-based German learning app.
Read `GOAL.md` for the intended product experience and `AGENTS.md` for the
repository workflow.

## Your Job

Debug the app autonomously. Find the actual cause of a failure, reproduce it,
fix it at the narrowest responsible layer, and leave behind a regression test
so the same bug does not return.

If a specific bug is described, start there. If no bug is described, inspect
the current branch, recent changes, test failures, and obvious error paths to
choose the highest-signal problem. Do not manufacture a feature just to have
something to change.

Keep going until the reported failure is genuinely handled and closely related
regressions are covered, while keeping each individual fix small and testable.

## Debugging Workflow

1. Start with the startup routine in `AGENTS.md`: read
   `docs/backlog/active.md`, `docs/agent/index.md`, and
   `docs/agent/continuity.md`. Read `GOAL.md` and inspect the working tree
   before editing.

2. Establish a baseline. Run the smallest relevant test or command first and
   capture the exact failure, expected behavior, and observed behavior. Use
   `rg` to trace the failing symbol and read the nearest package `README.md`
   and subtree `AGENTS.md` before touching that package.

3. Reproduce before fixing. For a TUI bug, use `tui-tester` and follow its
   required loop: Observe → Reason → Act → Synchronize. Start the app with a
   unique data directory, wait for a visible anchor after every asynchronous
   action, and read `tui_tester/AGENTS.md` before running E2E or visual tests.
   For storage, scheduling, parsing, or AI bugs, create the smallest focused
   fixture that demonstrates the failure.

4. Identify the root cause rather than masking the symptom. Check state
   ownership, stale async results, bounds and empty-state handling, Unicode
   and user input, persistence errors, and view-specific key routing as
   appropriate. Prefer a fix in the layer that owns the broken invariant.

5. Add or update a regression test that fails for the old behavior and passes
   for the new behavior. Keep the test deterministic: avoid wall-clock races,
   random ordering, network calls, the real user database, and restricted
   local dictionary data.

6. Run verification in layers: the focused test, the affected package tests,
   `go test ./...`, and then `./scripts/verify.sh` when feasible. Use
   `go test -race ./...` for changes involving concurrent commands, maps,
   storage, or asynchronous UI state. If a check fails, fix that failure
   before starting another change.

7. Look for one or two directly related cases exposed by the same root cause
   and cover them if the fix can remain narrow. Stop when the bug family is
   handled; avoid opportunistic refactors and unrelated feature work.

## Debugging Priorities

- Prevent crashes, panics, data loss, and incorrect review grades first.
- Then fix stale async state, persistence failures, input/key conflicts, empty
  states, Unicode corruption, and layout or accessibility regressions.
- Preserve the local-first design and existing keyboard and mouse workflows.
- Keep error messages actionable and keep dynamic TUI status text on the
  stable `status:` surface used by the tester.
- Never commit restricted dictionary archives or other generated/private data.

## Documentation and Handoff

Update documentation continuously, not only at the end:

- Put unfinished executable work and the exact next action in
  `docs/backlog/active.md`.
- Record completed fixes and their verification in `docs/backlog/done.md`.
- Add a dated notice under `docs/agent/notices/` for a durable trap,
  invariant, or cross-package contract, and index it from
  `docs/agent/index.md`.
- Update an ADR only when the debugging fix changes a long-lived architectural
  decision.

Do not leave important reproduction steps, blockers, or warnings only in chat.

## Boundaries

- Preserve unrelated user changes in a dirty worktree.
- Do not use destructive commands such as `git reset --hard` or broad deletes.
- Do not weaken, delete, or skip a test merely to make the suite green.
- Do not add retries, sleeps, ignored errors, or broad recoveries without
  explaining why they address the root cause.
- Do not call external services for a bug that can be reproduced locally.
- If the fix requires new authority, unavailable data, or a materially
  different product decision, document the blocker and ask for direction.

## Done When

- The failure has a reliable reproduction or a documented reason it cannot be
  reproduced locally.
- The root cause is fixed with a narrow implementation change.
- A deterministic regression test covers the old failure.
- Focused tests and `./scripts/verify.sh` pass, or any remaining blocker is
  explicitly recorded in the active backlog.
- Relevant continuity documentation is current.
- The completed changes are committed after verification.
