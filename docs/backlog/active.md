# Active Backlog

Last updated: 2026-04-30

## Current Milestone

Autonomous Feature Pass 5: Advanced Study & Multimedia

## Planned Features

### Feature 13: Full Cram Review Flow

User story: As a learner, I can actually review cards in Cram Mode (reveal/grade) so I can study flagged cards intensely without affecting my SRS schedule.

Acceptance criteria:
- Pressing `Enter` on a card in Cram Mode opens the "Cram Review" view.
- Cram Review works like normal Review: `Space` reveals answer, `a`/`h`/`g`/`e` or `1`-`4` for grading.
- Grading in Cram Mode updates `cramReviewed` and `cramCorrect` session stats but does NOT call SRS scheduler or update `reviews` table.
- Pressing `q` or `Esc` exits Cram Review back to Cram list.
- Accuracy % is displayed after each cram review.

Status: COMPLETED - All E2E tests pass.

### Feature 14: Audio Support Infrastructure

User story: As a learner, I can hear the pronunciation of cards so I can improve my listening comprehension.

Acceptance criteria:
- Add `audio` column to `cards` table via SQLite migration.
- `core.Card` struct includes `Audio` string field.
- TSV import supports a 4th column for audio file paths.
- Review and Cram views show an `[Audio]` indicator if audio is available.
- Pressing `p` in Review/Cram view plays the audio using `afplay` (macOS) or `play` (Linux).
- Auto-play audio on reveal if configured in Settings.

Status: COMPLETED - Audio infrastructure exists, `playAudio` function implemented, `[Audio]` indicator shows in views.

### Feature 15: Review History Visualizations (Heatmap/Chart)

User story: As a learner, I can see my study activity over the last 30 days so I can stay motivated and see my progress trends.

Acceptance criteria:
- Statistics view includes a "Review Activity" section.
- Displays a simple ASCII bar chart or heatmap of reviews per day for the last 14-30 days.
- Derives data from existing `reviews` table.
- Updates in real-time after a review session.

Status: COMPLETED - 14-day ASCII bar chart implemented and showing in Statistics view.

## Next Action

### All planned features for this milestone are complete!

The application is fully functional with:
- All core views rendering correctly
- Complete review flow (Flashcard and MCQ)
- Cram mode with reveal/grade
- Statistics with review activity visualization
- Audio support infrastructure
- 38 E2E tests passing
- `./scripts/verify.sh` passing with zero errors

Future work can focus on:
- Additional UI/UX improvements
- More advanced study features
- APKG import/export (high impact, high effort)

## Current Session Notes

- 2026-04-30: Implemented Feature 15 (Review History Visualizations) with ASCII bar chart in Statistics view.
- Added `ReviewsPerDay` method to repository interface and SQLite implementation.
- Extended TUI model to load and display 14-day review activity chart.
- All tests pass: Go tests, E2E tests, and verify.sh.

## Last Verified

- 2026-04-30: `./scripts/verify.sh` passed: Go tests, smoke launch, and 38 E2E tests.

## Feature Candidates (Future)

- Per-card review history popover: medium impact / medium effort.
- APKG import/export: high impact / high effort.
- Deck tags and filtering: medium impact / medium effort.
- Streak visual indicator: low impact / low effort.
