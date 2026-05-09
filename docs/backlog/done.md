## 2026-05-10: Advanced Content & UI Polish

- **German Content Expansion:** Added three new decks: "German Prepositions" (mastering Accusative/Dative cases with MCQs), "B1 Nature & Environment" (essential ecological vocabulary), and "Advanced Feelings" (nuanced B2-C1 emotional vocabulary).
- **Dashboard Enhancements:** Added a 14-day "Recent Activity" sparkline chart using Unicode block characters. Reorganized layout to maintain visibility of footer and status line on standard terminals. Added help hints to the dashboard footer.
- **UI/UX Improvements:** Implemented a new Color Theme setting (system, nordic, sunset) accessible via a dedicated 'c' key in Settings. Added a selection counter to the Browser title for better bulk action feedback.
- **Developer Experience:** Integrated structured logging into the SRS Scheduler to track review events and parameter transitions.
- **Robustness & Stability:** Fixed multiple E2E test regressions by implementing intelligent status message preservation (preventing background updates from overwriting help or grading messages).
- **Bug Fix:** Resolved a missing import issue in the dashboard rendering logic.

## 2026-05-10: Multi-Area Enhancements & Cloze Support

- **German Content:** Added new B2-level "Society & Social Issues" deck (30 notes) and expanded Grammar Tips with 10 advanced B2-C1 level tips.
- **Cloze Support:** Fully implemented Cloze deletion support in the TUI, including multi-marker highlighting and correct typing mode verification against the missing word.
- **Statistics View:** Modernized with a 2-column layout, visual progress bars for daily goals and grades, and a new 7-day activity bar chart.
- **Dashboard:** Improved responsiveness by dynamically calculating available height for the Grammar Tip and Quick Actions sections.
- **Decks View:** Added interactive "[Study]" and "[Cram]" buttons for each deck, allowing users to jump directly into a session.
- **Bug Fix:** Resolved a critical bug where typing mode state was not reset between cards, causing input collisions and accidental view switching.
- **E2E Testing:** Added `test_cloze_typing.py` to verify the new Cloze typing functionality.
- **Verification:** Fixed multiple test regressions caused by layout changes, ensuring all 150 tests pass.