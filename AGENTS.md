# Agent Instructions

Use this file at the start of every new chat session.

## Current State

The project is a Go Bubble Tea TUI for German flashcard and MCQ learning. It is local-first, uses SQLite for progress, uses editable deck files for content, and targets Anki-friendly TSV import/export before `.apkg`.

See `GOAL.md` for the intended final state of the app. All work should move toward that vision.

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

## Common Tasks

- **Adding a new card type**: Read `internal/core/AGENTS.md` (models) and `internal/content/AGENTS.md` (parsing).
- **Modifying the UI**: Read `internal/tui/AGENTS.md` (responsive rules, mouse support).
- **Updating the Database**: Read `internal/storage/AGENTS.md` (migration rules).
- **Adding an AI feature**: Read `internal/ai/AGENTS.md` (mocking, validation).
- **Fixing a scheduling bug**: Read `internal/srs/AGENTS.md` (encapsulation).
- **Running E2E / visual tests**: Read `tui_tester/AGENTS.md` first (see below).

## E2E Testing with tui-tester

For visual/incremental testing of the TUI, use `tui-tester` (Python CLI at `tui_tester/`).

**Before running E2E tests, you MUST read `tui_tester/AGENTS.md`.**

Quick reference:
```bash
# Start app with unique data dir (required for SQLite apps)
tui-tester start "./language-learning-tui --data-dir /tmp/test-$RANDOM"

# Check what's on screen
tui-tester observe

# Send input (keys, mouse clicks, drags)
tui-tester act "<Enter>"
tui-tester act "<Tab>"
tui-tester click 15 3

# Wait for specific text before next action (NEVER skip this!)
tui-tester wait-for "DASHBOARD"
tui-tester wait-stable

# Cleanup when done
tui-tester stop
```

**The loop:** Observe → Reason → Act → Synchronize (repeat). Never skip synchronization — TUI updates are async.

## Search Hints

- TUI behavior: `rg "type Model|Hitbox|WindowSizeMsg|Mouse" internal/tui`
- Scheduling: `rg "Scheduler|Review" internal/core internal/srs`
- Storage: `rg "Repository|migration|CREATE TABLE" internal/storage`
- Anki import/export: `rg "TSV|Anki|Import|Export" internal/content`
- AI drafting: `rg "Provider|Draft" internal/ai`

## Current Gotchas

- macOS shells used in this repo may not provide `timeout`; use a TTY session or `./scripts/tui_smoke.sh` for launch checks.
- Embedded TSV decks are parsed by Go's TSV reader; avoid unescaped quote characters inside fields because they can merge following rows.
- Settings screens are height-sensitive in E2E tests; use 110x40+ terminals for Settings tests to ensure Daily Goal section is visible without scrolling.
- E2E tests often wait for "DASHBOARD" to verify app readiness. When updating the Dashboard header, ensure the word "DASHBOARD" remains visible to maintain test compatibility.
- Dashboard layouts are height-sensitive. Avoid adding vertical space in top sections (header/boxes) as it can push bottom elements (Grammar Tip, help overlays) off-screen, causing E2E test timeouts on standard terminal sizes.
- Bubble Tea unit tests that manually call commands often fail if a command is changed to a `tea.Batch`. Use a `executeCmd` helper in tests to flatten batches and prevent `undefined: msg` errors.
- TSV embedded decks with `#deck:` lines containing `&` characters have different deck IDs than their filename (e.g., `a2-shopping-services.tsv` with `#deck: A2 Shopping & Services` gets ID `a2_shopping_&_services`). Do NOT create Go decks with IDs that would conflict with these. Use `_` underscores in deck IDs (e.g., `a2_purchasing_wear`).
- When adding scrollbar columns to custom render views, pad content with `layout.Width - 2` so the scrollbar stays inside panel bounds on first render.
- **E2E Settings Tests:** The Settings view requires scrolling to reach Daily Goal on smaller terminals. Use `columns=110, lines=40+` and navigate with 'j' keys to reach the Daily Goal row before testing +/- adjustments.
- **E2E Decks View Search:** After searching in Decks view, press `<Esc>` to clear the filter before starting a new search. Search state persists across navigation.
- **E2E Review Empty State:** The starter deck has 52 cards due by default. Tests expecting "No cards due" must grade through all cards first or use a fresh database without seeded content. Expanding Standard Content (Import `S`) does not change this starter count.
- **TSV Import Format:** When creating TSV files for E2E tests, use 6-column format: `id\tfront\tback\textra\ttags\tdeck`. The 6th column sets the deck name.
- **Local Dictionary Data Constraint:** The directory `/local_dict_files/` contains raw zip files of the dict.cc German-English dictionary dataset. This dataset is restricted by license terms and MUST NOT be committed, published, or pushed to any repository. It is gitignored.


## Handoff Rules

- Put unfinished executable work in `docs/backlog/active.md`.
- Put durable warnings or traps in `docs/agent/notices/`.
- Put long-lived architecture decisions in `docs/decisions/`.
- Keep all continuity notes concise and indexed from `docs/agent/index.md`.
- Do not dump chat history into `AGENTS.md`; write searchable, task-shaped records instead.
- Use `./scripts/verify.sh` as the full project verification command.
- When adding UI elements that increase vertical height (e.g. Dashboard boxes), ensure E2E tests using `tui-tester` have sufficient terminal `lines` configured. Elements appended after the main layout (like help overlays) may be clipped and become invisible to `wait_for_text` if the terminal is too short.
- **Bubble Tea Key Strings:** Space character is often received as `"space"` string in `msg.String()`. Always handle both `" "` and `"space"` in single-character input helpers to avoid dropping spaces in search filters.
- **Card Loading Limits:** For large seeded collections (Standard Content), keep the default `DueCards` load cap comfortably above current seeded size (currently 20,000). Lower caps can cause some decks to appear empty in the Review view if their card IDs sort after the limit.
- **E2E Test Selectors:** Prefer unique tags or specific substrings for E2E `wait_for_text` to avoid collisions between similar decks (e.g., "B1 Environment" vs "C1 Environment"; search "Seasons" not "Weather" so "A1 Weather & Seasons" is not confused with "A1 Time & Weather").
- **Dictionary is a Spotlight Overlay:** Dictionary is NOT in the tab/arrow/WASD navigation cycle. It is a Spotlight-like overlay triggered by `=` from any view. The legacy full `ViewDictionary` tab still exists for `/` (Dashboard) and `d` (Review) shortcuts. E2E tests must not assume Dictionary in the tab cycle.

## Shortcuts Reference

- **Global Navigation:** `1-9` for primary views, `0` for Practice Hub. `Tab`/`Arrows` to cycle.
- **Selection:** `m` for multi-select in Decks and Browser. `x` also works in Decks.
- **Dictionary:** `=` for Spotlight (anywhere), `/` for full tab (Dashboard).
- **Review:** `1-4` for grading, `Space`/`Enter` to reveal. `d` for Dictionary lookup.
- **Practice:** `1-9` to select trainer within the Practice Hub.

