# Improve deutsch-tui

You are working on deutsch-tui, a terminal-based German learning app.
Read `GOAL.md` for what the app should ultimately feel like.

## Your Job

Make this app meaningfully better. Don't settle for one tiny tweak — find real problems, fix them properly, and keep going. You have time. Use it.

Some directions worth exploring (pick any, combine them, or find your own):

- Refine the UI — layouts, spacing, colors, visual hierarchy
- Fix bugs — broken state, wrong keybindings, edge cases
- Optimize performance — render hot paths (`View()` runs per message), database queries and indexes, startup latency, avoidable allocations; keep behavior byte-identical unless the change is the fix
- Add learning content — vocabulary decks, grammar exercises
- Improve animations and feedback — transitions, progress indicators
- Enhance the AI view — prompts, responses, UX flow
- Strengthen developer experience — code structure, error messages
- Write or improve tests — unit tests, E2E tui-tester scenarios
- Polish documentation — fix inconsistencies, update stale content

## How to Work

1. Start with the startup routine in `AGENTS.md`.
2. Look around. Understand what's there before changing anything.
3. Pick something concrete, ship it fully, then pick the next thing.
4. Don't stop after one fix if there's more you can clearly improve. Build momentum.
5. Run `./scripts/verify.sh` after each meaningful change. Keep things green as you go — don't accumulate breakage.
6. Update `docs/backlog/active.md` and `docs/backlog/done.md` as you go.
7. If something breaks while you're working, fix it — that's now your task.
8. If an approach isn't working after a few tries, step back and try something different.
9. Commit your work when verification passes.

## Boundaries

- Keep each individual change narrow and testable, but do multiple changes per session.
- Don't refactor for the sake of refactoring.
- Don't add features that conflict with `GOAL.md`.
- Stability comes first. If you're choosing between a risky feature and a solid fix, pick the fix. But don't use "stability" as an excuse to do nothing.

## Tools Available

- Terminal: run Go tests, build, start/stop the app.
- File system: read, write, create files within `internal/` boundaries.
- TUI Testing: `tui-tester` for E2E scenarios (see `tui_tester/README.md`).
- Search: look up Go docs, Bubble Tea examples, etc.

## Done When

- `./scripts/verify.sh` passes.
- You've made at least 2–3 real improvements (not just whitespace or comment changes).
- Your changes are committed.
- `docs/backlog/active.md` reflects the current state.
