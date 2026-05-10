## 2026-05-10: Autonomous Pass - Batch 1 Mobility & Review Surfaces

- **German Content:** Added an embedded "B2 Urban Mobility" TSV deck with 25 notes covering public transport, commuting, accessibility, MCQs, and cloze cards.
- **Dashboard UI:** Added a tall-layout "Card Mix" panel that visualizes New, Young, and Mature card distribution.
- **Review UI:** Added deck name, card type, and tags to the Review header so learners can keep context while grading.
- **Browser UI:** Expanded Card Preview with deck, kind, extra notes, and tag fallback text.
- **AI/Settings UX:** Reworked AI suggested topics around CEFR levels and added a deterministic provider-cycle hint in Settings.
- **Testing:** Added content and TUI unit tests plus 8 new tui-tester scenarios in `e2e_tests/test_autonomous_pass_mobility_ui.py`.
- **Bug Fix:** Removed unescaped quote characters from mobility MCQ explanations after they caused CSV parsing to collapse rows.
- **Verification:** `go test ./internal/content ./internal/tui` passed on 2026-05-10 13:31 +07.

## 2026-05-10: Autonomous Pass - Batch 2 Export & Empty-State Polish

- **Bug Fix:** TSV export now honors the same Mature/Learning status filter that APKG export already used.
- **Import UX:** Added a visible `[R] Reset DB` action, click handling for Seed Standard and Reset DB, and explicit status-filter guidance.
- **Browser UX:** Empty Browser results now name the active deck and tell learners how to clear search/tag filters.
- **Cram UX:** Cram Review now shows deck name, card position, card type, and tags; all-deck cram lists include deck context.
- **Help Polish:** Help overlay now documents Browser tag filtering and Import database reset.
- **Testing:** Added unit tests for TSV status filtering, Import guidance, Browser empty states, and Cram metadata.
- **Verification:** `go test ./internal/content ./internal/tui` and `pytest -q e2e_tests/test_autonomous_pass_mobility_ui.py` passed on 2026-05-10 13:47 +07.

## 2026-05-10: Autonomous Pass - Final Verification

- **Regression Fixes:** Restored Dashboard "Recently Studied" visibility at standard E2E heights, restored the AI tip text expected by core-view tests, and moved Settings editing/status prompts into visible content for compact terminals.
- **Instruction Hygiene:** Added concise AGENTS.md gotchas for macOS launch checks, embedded TSV quoting, and compact Settings layout constraints.
- **Verification:** `./scripts/verify.sh` passed on 2026-05-10 14:09 +07: all Go tests, `go vet`, smoke startup, and 166 tui-tester E2E tests.
- **Milestone Summary:** Completed a multi-area pass covering content, dashboard/review/browser/cram/settings/import/AI UX, TSV export correctness, click handling, unit tests, E2E coverage, and continuity documentation.

## 2026-05-10: Batch 2 Improvements & UI Polish

- **UI/UX Refinements:**
    - Added "Recently Studied" section to Dashboard, showing the last 3 decks studied.
    - Added "Card Preview" box to Browser view, showing front/back/tags of the selected card.
    - Escaped newlines in AI template previews in Settings to maintain layout integrity.
    - Optimized Dashboard height thresholds for Grammar Tips and Quick Actions.
- **Content Expansion:**
    - Added new "German Business Vocabulary" deck with 15 workplace-relevant notes.
    - Expanded Grammar Tips with advanced topics (Adverbial Genitive, Wissen vs Kennen, etc.).
- **AI Improvements:**
    - Significantly improved AI templates for vocabulary, grammar, and conjugation.
- **Developer Experience:**
    - Added a "Debug Log" view accessible via `Ctrl+D`, showing internal state and active view info.
- **Testing:**
    - Added 8 new E2E tests in `e2e_tests/test_batch_2_improvements.py` verifying all new features.
    - Fixed multiple regressions in existing tests caused by layout and default config changes.
- **Bug Fixes:**
    - Fixed tag splitting logic to handle both commas and whitespace.
    - Fixed special key handling in TUI tester (`Ctrl-D` sequence).
    - Fixed missing imports in `handlers.go`.

## 2026-05-10: Advanced Content & UI Polish

- **German Content Expansion:** Added three new decks: "German Prepositions" (mastering Accusative/Dative cases with MCQs), "B1 Nature & Environment" (essential ecological vocabulary), and "Advanced Feelings" (nuanced B2-C1 emotional vocabulary).
- **Dashboard Enhancements:** Added a 14-day "Recent Activity" sparkline chart using Unicode block characters. Reorganized layout to maintain visibility of footer and status line on standard terminals. Added help hints to the dashboard footer.
- **UI/UX Improvements:** Implemented a new Color Theme setting (system, nordic, sunset) accessible via a dedicated 'c' key in Settings. Added a selection counter to the Browser title for better bulk action feedback.
- **Developer Experience:** Integrated structured logging into the SRS Scheduler to track review events and parameter transitions.
- **Robustness & Stability:** Fixed multiple E2E test regressions by implementing intelligent status message preservation (preventing background updates from overwriting help or grading messages).
- **Bug Fix:** Resolved a missing import issue in the dashboard rendering logic.

## 2026-05-10: Stability, Content & UX Improvements

- **Stability & Bug Fixes:** Fixed a critical panic in the Decks view when applying filters. Resolved E2E test failures in `test_new_content_visibility.py` by improving wait text and robustness.
- **Decks View Navigation:** Renamed Decks view title to "DECK LIST" for better differentiation. Enhanced the search feature to automatically select the first matching deck, streamlining the workflow for large collections.
- **German Content:** Added a new "Common German Verbs (A1-A2)" deck with 20 essential verbs and conjugation MCQ cards. Added advanced grammar tips for C1-C2 level learners.
- **Audio Improvements:** Optimized the `playAudio` command to be non-blocking, preventing TUI hangs during audio playback.
- **Enhanced Statistics:** Added session duration tracking and average review speed (cards/min) to the Statistics view.
- **Developer Tools:** Implemented a "Reset Database" action in the Import view with a confirmation step, making it easier to wipe state during development.
- **UI Polish:** Improved the Browser view's help hint to include the previously hidden Backspace shortcut for deleting cards.
- **Verification:** Successfully passed all 150 project tests, including updated and new E2E scenarios.

## 2026-05-10: Content, UI Polish & SRS Robustness

- **German Content:** Added "German Confusable Words" deck (wissen/kennen, liegen/legen, etc.) with 20+ notes, including MCQs and Cloze cards.
- **Grammar Tips:** Expanded `GrammarTip` struct to include a concrete `Example` field and updated all 50+ tips with German examples.
- **Dashboard UI:** Updated layout to display the new grammar tip examples when space permits, enhancing daily learning value.
- **Decks View:** Implemented a visual scrollbar with mouse-hitbox support, bringing it in line with Browser and Statistics views.
- **Statistics View:** Added a "Maturity Distribution" chart visualizing the percentage of New, Young, and Mature cards in the collection.
- **SRS Robustness:** Added comprehensive unit tests for the scheduler logic, covering various grading sequences and interval growth.
- **Developer Experience:** Enhanced key event logging in `internal/tui/keys.go` to help debug text input trapping and global navigation issues.
- **Bug Fix:** Fixed a race condition in `test_new_content_visibility.py` by adding stability waits and ensuring consistent view state.
- **Verification:** Successfully passed the full suite of 163 tests, including a new E2E test `test_ui_polish.py` and fixes for layout regressions.
