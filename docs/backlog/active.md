# Active Backlog

Last updated: 2026-04-29

## Current Milestone

Autonomous product feature pass: polish deutsch-tui with 2-3 vertical learning features.

## Feature Candidates

- Bookmarked cards: high impact / medium effort. Lets learners mark hard or important cards during review and see saved-card counts.
- Undo last review: high impact / medium effort. Fixes accidental grading without manual SQLite edits.
- Daily progress and streak stats: high impact / low effort. Makes the Statistics view actionable by showing today, goal progress, and current streak.
- Custom daily goal setting: medium impact / medium effort. Useful, but config-driven and less central than review correctness.
- Leech detection: medium impact / medium effort. Valuable later once lapse history and bookmark UX are stronger.
- Filtered review mode: high impact / high effort. Needs careful session state and deck/card queries.
- APKG import/export: high impact / high effort. Already on roadmap but larger than a safe 2-3 feature vertical pass.

## Selected Features

### Feature 1: Bookmarked Review Cards

User story: As a learner, I can bookmark a card while reviewing so I can notice and revisit important or troublesome material later.

Acceptance criteria:
- `b` toggles bookmark for the current review card.
- Review view clearly shows bookmark state.
- Dashboard/Statistics show total bookmarked cards.
- Bookmark state persists in SQLite across app restarts.
- Unit coverage includes SQLite persistence and TUI toggling; E2E covers bookmark persistence after restart.

UI/UX and components:
- `internal/tui/model.go`: Review key handling, review rendering, dashboard/statistics counts, status messages.
- `internal/core/core.go`: repository contract and statistics field.
- `internal/storage/sqlite`: safe migration and repository methods.

Data/model changes:
- Add numbered SQLite migration creating `card_flags(card_id PRIMARY KEY, bookmarked, updated_at)`.
- Add repository methods for reading and toggling bookmark state.

Test plan:
- Go unit tests for migration, bookmark toggle/read, stats count, and TUI state/rendering.
- tui-tester E2E: bookmark a review card, restart, confirm bookmark indicator remains.

### Feature 2: Undo Last Review

User story: As a learner, I can undo an accidental review grade immediately and get the card back without corrupting SRS history.

Acceptance criteria:
- `u` after grading restores the most recently graded card to its previous review state.
- Due-card count and review queue refresh immediately.
- The rollback persists in SQLite across restart.
- If no review can be undone, the TUI shows a clear status message.

UI/UX and components:
- `internal/tui/model.go`: track last reviewed card, route `u`, refresh due cards/decks/statistics after undo.
- `internal/storage/sqlite`: delete newest review row for a card and restore `review_states` from the prior review row or remove it for new cards.

Data/model changes:
- No new schema required; use existing `reviews` history and `review_states`.
- Repository contract gets `UndoLastReview(ctx, cardID)`.

Test plan:
- Go unit tests for reverting first review and reverting to prior review state.
- tui-tester E2E: grade one card, undo, restart, confirm due count is restored.

### Feature 3: Daily Progress and Streak Statistics

User story: As a learner, I can see today’s review count and current study streak so the Statistics view supports daily habit formation.

Acceptance criteria:
- Statistics view shows Reviews Today, Daily Goal progress, and Current Streak.
- Counts are derived from persisted review history.
- Values update after grading and survive restart.

UI/UX and components:
- `internal/core/core.go`: extend `Statistics`.
- `internal/storage/sqlite`: aggregate review dates in local persisted history.
- `internal/tui/model.go`: render daily progress in Statistics.

Data/model changes:
- No schema required; derive from `reviews.reviewed_at`.

Test plan:
- Go unit tests for today count and streak derivation.
- tui-tester E2E: grade a card, open Statistics, confirm today/goal progress update.

## Next Action

Autonomous feature pass is complete. Resume the next roadmap milestone by implementing `.apkg` import/export or audio/media support unless a newer user request supersedes it.

## Feature Specs

### Feature 4: Bookmarked Cards Review Mode

User story: As a learner, I can filter my review queue to only bookmarked cards so I can focus on difficult or important material.

Acceptance criteria:
- `B` (shift+b) in Review view toggles bookmarked-only filter mode.
- When active, only bookmarked cards appear in the review queue.
- A banner shows "Bookmarked mode: ON" in the review view.
- Dashboard shows "X bookmarked due" alongside normal due count.
- Statistics view shows bookmarked mode count if applicable.
- No schema changes (uses existing card_flags table).

UI/UX and components:
- `internal/tui/model.go`: add `bookmarkFilter` bool field, `B` key handler, filtered due cards, banner rendering.
- `internal/core/core.go`: extend `Statistics` with `BookmarkedDue` count.
- `internal/storage/sqlite`: add `DueCardsByBookmark(ctx, now, limit, bookmarked bool)` method.

Test plan:
- Go unit tests for filter toggle, card filtering, and Statistics bookmarked-due count.
- tui-tester E2E: toggle bookmark filter, verify only bookmarked cards shown.

### Feature 5: MCQ Card Review Flow

User story: As a learner, I can answer multiple-choice cards by selecting a numbered option so I can practice recognition alongside recall.

Acceptance criteria:
- MCQ cards render their choices as numbered options (1-4) in the Review view.
- Pressing `1`-`4` selects the corresponding choice and reveals whether it's correct.
- After selection, grading keys work as normal.
- Flashcard flow is unchanged for non-MCQ cards.
- Correct/incorrect feedback shown immediately after choice selection.

UI/UX and components:
- `internal/tui/model.go`: add MCQ rendering branch, choice selection state, correct/incorrect feedback.
- No schema changes (MCQ choices already stored in `cards.kind = 'mcq'` and `core.Card.Choices`).
- Content layer needs to populate `Choices` during TSV import for MCQ notes.

Test plan:
- Go unit tests for MCQ rendering, choice selection, and correctness check.
- tui-tester E2E: review MCQ card, select correct choice, verify feedback.

### Feature 6: Leech Detection

User story: As a learner, I can see which cards I repeatedly fail so I can focus extra effort or suspend them.

Acceptance criteria:
- Cards with 3+ consecutive "Again" grades are automatically flagged as "leech".
- Leech cards show a `LEECH` indicator during review.
- Statistics view shows total leech count.
- Leech flag is persisted in SQLite (new migration for `card_flags.leech` column or separate table).
- Reset leech count when card is graded Hard/Good/Easy.

UI/UX and components:
- `internal/tui/model.go`: leech indicator in Review rendering, leech count in Statistics.
- `internal/storage/sqlite`: migration to add `leech` or `lapse_streak` to `card_flags`, update on review recording.
- `internal/core/core.go`: extend `Statistics` with `LeechCards` count.

Test plan:
- Go unit tests for leech flagging, reset on success, and Statistics count.
- tui-tester E2E: grade card "Again" 3 times, verify leech indicator appears.

## Blockers

- None.

## Last Verified

- 2026-04-29: `./scripts/verify.sh` passed with 27 E2E tests after the autonomous feature pass.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/test_learning_features.py -q` passed with 3 new E2E tests covering bookmark persistence, undo persistence, and daily Statistics progress.
- 2026-04-29: `go test ./...` passed after adding core/SQLite/TUI support for bookmarks, undo-last-review, and daily progress/streak statistics.
- 2026-04-29: Baseline `go test ./...` passed before feature work.
- 2026-04-29: Baseline `./scripts/tui_smoke.sh` passed before feature work.
- 2026-04-29: `./scripts/verify.sh` passed with 24 E2E tests (including new tests for Statistics view and updated navigation).
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests/test_recertification.py -q` passed with 3 new E2E tests.
- 2026-04-29: Added three recertification E2E tests for Tab view cycling, Hard grade SQLite persistence, and Settings provider persistence.
- 2026-04-29: `tui_tester/venv/bin/python -m pytest e2e_tests -q` passed with 19 E2E tests.
- 2026-04-29: `./scripts/tui_smoke.sh` passed; `go test ./...` passed.
- 2026-04-29: `./scripts/verify.sh` passed with 19 E2E tests, all Go tests, smoke test, gofmt, and go vet.
- 2026-04-29: Cleaned up debug files and added test_data/ to .gitignore.
