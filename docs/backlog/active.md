# Active Backlog

Last updated: 2026-04-30

## Current Milestone

Autonomous Feature Pass 10: Deck Tags and Filtering

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

### Feature 16: Deck Tags and Filtering

User story: As a learner, I can organize and filter my decks by tags so I can focus on specific topics or categories.

Acceptance criteria:
- Add `tags` column to `decks` table via SQLite migration.
- `core.Deck` struct includes `Tags` string slice field.
- TSV import supports deck tags in a new header directive.
- Decks view shows deck tags below the description.
- Decks view supports text filtering that matches deck names, descriptions, or tags.
- Filtering is case-insensitive and shows live results as the user types.

Status: COMPLETED - Deck tags persisted in database, displayed in UI, and filtering implemented with 3 E2E tests.

### Feature 17: Audio Auto-play on Reveal

User story: As a learner, I can have audio played automatically when I reveal a card, so I can focus on pronunciation without extra steps.

Acceptance criteria:
- Add `autoplay_audio` boolean field to config (default false).
- Settings view includes toggle for "Auto-play audio on reveal".
- Review view automatically plays audio when card is revealed (if audio exists and auto-play enabled).
- Cram view also respects auto-play setting.

Status: PENDING - Implementation in progress.

## Current Pass Plan

Top blockers for this run:

- ✅ Establish a fresh baseline by launching `cmd/deutsch-tui` and running `./scripts/verify.sh`.
- ✅ Add deck tags column to `decks` table and repository support.
- ✅ Add deck filtering functionality to the Decks TUI view.
- ✅ Add at least 3 new tui-tester E2E tests covering deck tags and filtering.
- ✅ Re-run `./scripts/verify.sh` and commit all source changes.
- 🚧 Implement audio auto-play on reveal.

## Next Action

Implement audio auto-play on reveal. Strong next candidates:
- APKG import/export: high impact / high effort.
- Advanced deck organization (nested folders, color coding): medium impact / medium effort.
- Custom study sessions: medium impact / high effort.

## Current Session Notes

- 2026-04-30: Started Autonomous Feature Pass 8. Root and TUI agent instructions read; `go run ./cmd/deutsch-tui --help` launches successfully.
- 2026-04-30: Added visual Streak Indicator (🔥) to the Dashboard view based on the current streak computed from `stats.CurrentStreak`.
- 2026-04-30: Added 3 new E2E tests for Dashboard features and help overlay in `e2e_tests/test_dashboard_features.py`.
- 2026-04-30: Expanded `e2e_tests/test_statistics.py` to check for streak indicator in both Dashboard and Statistics.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 53 tui-tester E2E tests.
- 2026-04-30: Started Autonomous Feature Pass 9. Root continuity docs and TUI instructions read; `./scripts/tui_smoke.sh` passed; `./scripts/verify.sh` passed with 53 E2E tests before source changes.
- 2026-04-30: Implementing per-card review history via `core.Repository`, SQLite `reviews`, and Review/Browser TUI surfaces. Next verification: focused Go tests for storage and TUI rendering.
- 2026-04-30: Added `ReviewHistory` repository support plus Review `r` and Browser `Enter` history panels. `go test ./internal/storage/sqlite ./internal/tui` passed.
- 2026-04-30: Added `e2e_tests/test_review_history.py` with 3 review-history regressions. Targeted E2E passed: `tui_tester/venv/bin/python -m pytest e2e_tests/test_review_history.py -q`.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 56 tui-tester E2E tests.
- 2026-04-30: Started work on Deck tags and filtering feature.
- 2026-04-30: Added SQLite migration 16 to add tags column to decks table.
- 2026-04-30: Updated SQLite Store to persist and retrieve deck tags in UpsertDeck, GetDeck, and Decks functions.
- 2026-04-30: Added unit test coverage for deck tags functionality in SQLite storage.
- 2026-04-30: Added deck filtering functionality in TUI with text input and tag matching.
- 2026-04-30: Added E2E tests for deck tags display and filtering functionality.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 59 tui-tester E2E tests.

## Last Verified

- 2026-04-30: `./scripts/verify.sh` passed: Go tests, smoke launch, and 53 E2E tests.
- 2026-04-30: All verification steps completed successfully:
  - gofmt: all files formatted correctly
  - go test ./...: all packages pass
  - go vet ./...: no issues found
  - smoke test: app launches successfully
  - E2E tests: 53/53 pass
- 2026-04-30: Baseline `./scripts/verify.sh` passed at start of Feature Pass 9: Go tests, smoke launch, and 53 E2E tests.
- 2026-04-30: Targeted review-history E2E passed: 3/3.
- 2026-04-30: Final `./scripts/verify.sh` passed: Go tests, smoke launch, and 56 E2E tests.
- 2026-04-30: Final `./scripts/verify.sh` passed for Deck Tags and Filtering: Go tests, smoke launch, and 59 E2E tests.

## Feature Candidates (Future)

- APKG import/export: high impact / high effort.
- Advanced deck organization (nested folders, color coding): medium impact / medium effort.
- Custom study sessions: medium impact / high effort.
- APKG import/export: high impact / high effort.

## Current Session Notes

- 2026-04-30: Started Autonomous Feature Pass 8. Root and TUI agent instructions read; `go run ./cmd/deutsch-tui --help` launches successfully.
- 2026-04-30: Added visual Streak Indicator (🔥) to the Dashboard view based on the current streak computed from `stats.CurrentStreak`.
- 2026-04-30: Added 3 new E2E tests for Dashboard features and help overlay in `e2e_tests/test_dashboard_features.py`.
- 2026-04-30: Expanded `e2e_tests/test_statistics.py` to check for streak indicator in both Dashboard and Statistics.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 53 tui-tester E2E tests.
- 2026-04-30: Started Autonomous Feature Pass 9. Root continuity docs and TUI instructions read; `./scripts/tui_smoke.sh` passed; `./scripts/verify.sh` passed with 53 E2E tests before source changes.
- 2026-04-30: Implementing per-card review history via `core.Repository`, SQLite `reviews`, and Review/Browser TUI surfaces. Next verification: focused Go tests for storage and TUI rendering.
- 2026-04-30: Added `ReviewHistory` repository support plus Review `r` and Browser `Enter` history panels. `go test ./internal/storage/sqlite ./internal/tui` passed.
- 2026-04-30: Added `e2e_tests/test_review_history.py` with 3 review-history regressions. Targeted E2E passed: `tui_tester/venv/bin/python -m pytest e2e_tests/test_review_history.py -q`.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 56 tui-tester E2E tests.
- 2026-04-30: Started work on Deck tags and filtering feature.
- 2026-04-30: Added SQLite migration 16 to add tags column to decks table.
- 2026-04-30: Updated SQLite Store to persist and retrieve deck tags in UpsertDeck, GetDeck, and Decks functions.
- 2026-04-30: Added unit test coverage for deck tags functionality in SQLite storage.
- 2026-04-30: Added deck filtering functionality in TUI with text input and tag matching.
- 2026-04-30: Added E2E tests for deck tags display and filtering functionality.
- 2026-04-30: Final verification passed with `./scripts/verify.sh`: Go tests, smoke launch, and 59 tui-tester E2E tests.

## Last Verified

- 2026-04-30: `./scripts/verify.sh` passed: Go tests, smoke launch, and 53 E2E tests.
- 2026-04-30: All verification steps completed successfully:
  - gofmt: all files formatted correctly
  - go test ./...: all packages pass
  - go vet ./...: no issues found
  - smoke test: app launches successfully
  - E2E tests: 53/53 pass
- 2026-04-30: Baseline `./scripts/verify.sh` passed at start of Feature Pass 9: Go tests, smoke launch, and 53 E2E tests.
- 2026-04-30: Targeted review-history E2E passed: 3/3.
- 2026-04-30: Final `./scripts/verify.sh` passed: Go tests, smoke launch, and 56 E2E tests.

## Feature Candidates (Future)

- APKG import/export: high impact / high effort.
- Deck tags and filtering: medium impact / medium effort.
