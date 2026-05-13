# Done Backlog

## 2026-05-13 20:00 +07 (Autonomous Improvement Pass - Additional)

Completed improvement pass with content expansion, AI enhancements, statistics improvements, and test coverage. Total E2E tests: 249 (all passing).

### Content Expansion
- Added **B2 German Music & Entertainment** vocabulary deck (52 cards) - music, concerts, films, artists, instruments, entertainment industry.

### AI Enhancements
- Added 8 new AI suggestions: B1 family & friends, B2 science & tech, B1 daily routine, C2 legal German, A1 colors & shapes, B1 sports & fitness, B2 media & news, C1 academic writing.
- Adjusted suggestion grid from 4 to 5 columns for better screen utilization.

### Statistics Improvements
- Added **Retention Rate** metric to the Statistics view showing mature/total card ratio with color-coded indicators.

### Testing
- Added new E2E test file `test_b2_music_entertainment.py` covering new music deck and AI suggestions.

### Verification
- Total E2E tests: 249 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

---

## 2026-05-13 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with 3 new German content decks, UI enhancements, E2E test expansions, and test regression fixes. Total E2E tests: 247 (all passing).

### Content Expansion
- Added **C2 German Literature** vocabulary deck (30 cards) - literature analysis, book types, stylistic devices.
- Added **C1 German Politics & Society** vocabulary deck (40 cards) - government, elections, diplomacy, societal structures.
- Added **B1 German Jobs & Professions** vocabulary deck (40 cards) - workplace, career, employment contracts, corporate roles.

### UI/UX Improvements
- **Review View:** Centered the text and layout inside the empty state box for better visual alignment.
- **Browser View:** Replaced spacing with pipe separators (`|`) in the Card Preview box for clearer columnar reading.
- **Statistics View:** Fixed a bug where double-rendered ANSI codes were misapplied to the muted scrollbar footer.

### Testing & Developer Experience
- Added unit tests for utility functions `stripANSI`, `trimLastRune`, and `singlePrintableInput`.
- Fixed multiple test regressions caused by search filter state persistence in `test_new_batch3.py` and `test_new_decks_batch_3.py` by pressing Escape between searches.
- Fixed a truncation bug in `test_new_content_visibility.py` where the test failed because a 30-character deck name was visually truncated by the UI layout.
- Added new E2E test files: `test_c_level_decks.py`, `test_b1_jobs_deck.py`, and `test_ui_polish_batch_3.py` for comprehensive test coverage.

### Verification
- Total E2E tests: 247 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-12 16:00 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with content expansion, UX enhancements, bug fixes, and testing expansion. Total E2E tests: 239 (all passing).

### Content Expansion
- Added **A1 Family & Friends** vocabulary deck (25 cards).
- Added **A1 Greetings & Farewells** vocabulary deck (15 cards).
- Expanded **Verb of the Day** with new verbs (bezahlen, bestellen).
- Expanded **Grammar Tips** with 2 advanced topics (Compound Nouns, N-Declension for adjectives).

### UI/UX Improvements
- **Browser:** Made the `SEARCHING` label visually distinct when active.

### Logic & Developer Experience
- **Database:** Improved startup error handling for SQLite locked databases to provide clearer debugging hints to developers.

### Testing
- Added comprehensive E2E test file `test_a1_family_deck.py`.
- Fixed space interpretation issues in search boxes by narrowing search criteria.

### Verification
- Total E2E tests: 239 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-12 14:00 +07 (Autonomous Improvement Pass - Final)

Completed comprehensive improvement pass with content expansion, grading UX improvements, and test updates. Total E2E tests: 238 (all passing).

### Visual Feedback on Grading
- Enhanced grading status to show grade icon (✓/✗/~), remaining cards count, and session accuracy
- Example: "✓ Good | 51 cards due | 100% accuracy"
- Maintains backward compatibility with tests by keeping "cards due" text

### Content Expansion
- Added **A1 Emergency German** vocabulary deck (20 cards) - emergency phrases and helpful expressions
- Added **A2 Restaurant & Dining** vocabulary deck (52 cards) - restaurant, food, ordering vocabulary
- Expanded **Verb of the Day** with 8 additional verbs (total now 72)
- Added **6 new grammar tips** - Reflexive Verbs, Modal Particles, Word Order in Subordinates, Compound Nouns, etc.

### Testing
- Updated E2E tests to match new grading status format
- Fixed compatibility issues in test_status_and_browser_regressions.py and test_review_history.py

### Verification
- Total E2E tests: 238 (all passing)
- All Go unit tests passing
- `./scripts/verify.sh` executed successfully with zero errors

---

## 2026-05-12 12:00 +07 (Autonomous Improvement Pass - New Content & UI)

Completed comprehensive improvement pass with 3 new German content decks, UI enhancements for Cram mode and Decks view, and 15 new verbs. Total E2E tests: 233 (all passing).

### Content Expansion
- Added **B1 Art & Literature** vocabulary deck (40 cards).
- Added **A2 Hobbies & Free Time II** vocabulary deck (40 cards).
- Added **B2 Science & Nature II** vocabulary deck (40 cards).
- Expanded **Verb of the Day** with 15 additional verbs.

### UI/UX Improvements
- **Cram View:** Added progress bar, formatted the card with a rounded border, and improved alignment.
- **Decks View:** Padded deck names for alignment and colored the New (orange), Due (green), and Total (cyan) counts for better readability.

### Testing
- Added comprehensive E2E test file `test_new_decks_batch_3.py` for new decks.
- Fixed broken UI navigation and layout tests caused by styling upgrades.

### Verification
- Total E2E tests: 233 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-11 14:00 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with new B1 content decks, grammar expansion, and UI/UX refinements. Total E2E tests: 215 (all passing).

### Content Expansion
- Added **B1 Health & Medicine** vocabulary deck (40 cards).
- Added **B1 Culture & Society** vocabulary deck (40 cards).
- Added **B1 Sports & Fitness** vocabulary deck (40 cards).
- Added **B1 Education & Studying** vocabulary deck (40 cards).
- Expanded **Grammar Tips** with 10 advanced topics (C1/C2 level) including participle extensions and nominalization.

### UI/UX Improvements
- **Settings View:** Enhanced visual hierarchy with bold, colored section headers.
- **Review View:** Dynamic card borders that change color based on learning state (New/Learning/Mature).
- **Settings Toggles:** Colorized toggle statuses (e.g., "on" / "off") for better visibility and affordance.

### Verification
- Total E2E tests: 215 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-11 12:30 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with new B1 content decks, expanded vocabulary, enhanced UI, and test fixes. Total E2E tests: 202 (all passing).

### Content Expansion
- Added **B1 Technology & Digital Life** vocabulary deck (40 cards) - technology, computers, digital communication, internet terms.
- Added **B1 Finance & Banking** vocabulary deck (40 cards) - financial vocabulary, banking, money management.
- Added **B1 Environment & Sustainability** vocabulary deck (40 cards) - environmental vocabulary, ecology, sustainability topics.
- Expanded **Verb of the Day** with 15 additional verbs (50 total) - added modal verbs (wollen, können, müssen, sollen, dürfen), erklärären, beginnen, warten, antworten, fragen, zeigen, hören, vergessen, treffen.
- Added **10+ advanced grammar tips** (B1/B2/C1 level) - Konjunktiv II, Passiv, Infinitivgruppen, Konjunktiv I, Genitiv, etc.

### UI/UX Improvements
- **Browser Card Preview:** Enhanced preview panel with State indicator (NEW/LEARNING/MATURE) and Reviews count, replacing Extra field with more useful metadata.
- **AI View Suggestions:** Improved suggestion section with 12 clickable topic suggestions with descriptions, replacing generic "Suggested Topics:" header with "Click a topic or type your own, then press Enter:".

### Testing & Quality
- Fixed E2E tests in `test_autonomous_pass_mobility_ui.py` to match new browser preview format (State field).
- Fixed E2E tests in `test_new_batch3.py` and `test_ui_views_may11.py` to match new AI suggestions format.
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors (202 E2E tests).

### Verification
- Total E2E tests: 202 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-11 10:00 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with new content decks, verbs expansion, UI improvements, and 8 new E2E tests. Total E2E tests: 202 (all passing).

### Content Expansion
- Added **A1 Hobbies & Free Time** vocabulary deck (~35 cards).
- Added **A1 Food & Drink** vocabulary deck (~40 cards).
- Added **A1 Travel & Transport** vocabulary deck (~35 cards).
- Expanded **Verb of the Day** with 15 additional verbs.

### UI/UX Improvements
- Enhanced Dashboard streak indicator with color-coded progression (lightning, fire, double gold fire).
- Improved AI Drafts empty state with a stylized bordered box and clearer call to action.

### Testing
- Added comprehensive E2E test file `test_new_batch3.py` with 8 new scenarios.
- Verified dashboard widgets, UI views, AI draft interactions, and Cram mode navigation.

### Verification
- Total E2E tests: 202 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-11 07:45 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with 7+ new vocabulary decks, grammar expansion, UI polish, and enhanced testing. Total E2E tests: 194 (all passing).

### Content Expansion
- Added **A1 Family & Relationships** vocabulary deck (42 cards) - family members, relationships.
- Added **A1 Animals & Nature** vocabulary deck (47 cards) - animals, plants, natural elements.
- Added **A1 Colors & Shapes** vocabulary deck (35 cards) - colors, geometric shapes, size adjectives.
- Added **A1 Body Parts & Health** vocabulary deck (40 cards) - body parts, health conditions, medical terms.
- Added **A1 Clothing & Fashion** vocabulary deck (53 cards) - clothing items, accessories, shopping.
- Added **A1 School & Education** vocabulary deck (43 cards) - school items, education terms.
- Added **A1 Numbers & Time** vocabulary deck (58 cards) - numbers, time expressions, frequency.
- Expanded **Grammar Tips** with 10+ new A1-level topics (Numbers & Time, Greetings, Word Formation, etc.).
- Expanded **Verb of the Day** with 15+ additional common German verbs.

### UI/UX Improvements
- Added `successStyle` for consistent green success messaging.
- Improved Dashboard motivation messages with new success indicators.
- Enhanced Cram view filter display with success styling.
- Improved AI view empty state messaging and tip visibility.

### Testing
- Added comprehensive E2E tests for new content decks in `test_new_content_may11.py`.
- Added UI view rendering tests in `test_ui_views_may11.py`.

### Verification
- Total E2E tests: 194 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-10 18:00 +07 (Autonomous Improvement Pass)

Completed another improvement pass with 6+ meaningful enhancements. Total E2E tests: 173 (all passing).

### Content Expansion
- Added **A1 Recipes & Cooking** vocabulary deck (86 cards) - food, cooking verbs, meal times.
- Added **A1 Weather & Seasons** vocabulary deck (58 cards) - weather conditions, seasons, time expressions.
- Expanded grammar tips with 8 new topics (relative clauses, informal expressions, wechselpräpositionen, etc.).

### Testing
- Added E2E tests for AI draft navigation (`test_ai_draft_navigation.py`) - verifies AI draft approval workflow.
- Added E2E tests for Recipe deck (`test_recipe_deck.py`) - verifies deck loading and vocabulary access.
- Added E2E test for Weather deck (`test_weather_deck.py`) - verifies deck is accessible.

### Verification
- Total E2E tests: 173 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.


## 2026-05-10 17:30 +07 (Autonomous Improvement Pass)

Completed a comprehensive improvement pass with 11 meaningful enhancements.

### Content Expansion
- Added **German Phrasal Verbs** deck with separable and inseparable verbs.
- Added **Medical German (B2/C1)** deck for advanced learners.
- Implemented **Verb of the Day** on the Dashboard with daily conjugations.

### UI/UX Improvements
- **Dashboard:** Refactored for extreme compactness to support smaller terminals; restored original wording for E2E compatibility.
- **Browser:** Added **Search History** (up to 5 recent terms) and displayed **Review Intervals** (e.g., "1 day") for each card.
- **Review:** Added **Difficulty Badges** (NEW, LEARNING, MATURE) and a **Session Timer** for better progress tracking.
- **Motivation:** Improved Dashboard messages for daily goals and streaks.

### Logic & Developer Experience
- Enhanced `core.Card` and SQLite storage layer to persistently track and fetch `Interval` and `Reviews` count, enabling richer UI metadata.

### Verification
- Added 3 new E2E tests in `e2e_tests/test_new_pass_improvements.py`.
- Total E2E tests: 168 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully.


## 2026-05-10 15:30 +07 (Autonomous Improvement Pass)

Completed a major improvement pass focusing on Content, UI/UX, and Logic.

### Content Expansion
- Added **German Idioms & Proverbs** deck with 15 common expressions.
- Added **German Slang & Youth Language** deck with 12 modern informal terms.
- Expanded AI Drafting templates with a dedicated **Grammar Explanation** template set.

### UI/UX Refinements
- **Dashboard:** Added a compact, informative header showing "WILLKOMMEN!", daily review progress, and the current date.
- **Statistics:** Added a **"Success Rate per Deck"** visualization with color-coded progress bars.
- **Review View:** Made the **Context/Extra** field significantly more prominent during the reveal phase.
- **AI View:** Fixed the **loading spinner animation** by correctly batching the tick command with the drafting process.

### Logic & Testing
- **Browser:** Improved `Esc` key behavior—it now clears the search/tag filter when in searching mode (while `Enter` preserves it).
- **Testing:** Added a new E2E test file `e2e_tests/test_browser_search.py` covering real-time filtering and filter clearing.
- **Robustness:** Updated existing unit tests to handle batched commands in the TUI model.

### Verification
- All 165 E2E tests passed.
- All Go unit tests passed.
- `./scripts/verify.sh` executed successfully.

## 2026-05-10 14:30 +07 (Content, UI Polish & SRS Robustness)

Completed a multi-area improvement pass covering German content expansion, UI refinements (scrollbars, grammar examples), statistics visualization, and SRS robustness.

### Content Expansion
- Added **German Confusable Words** deck (e.g., *seit* vs *seid*, *wahr* vs *war*).
- Expanded **Grammar Tips** with real-world examples for each tip.

### UI/UX Refinements
- **Decks View:** Implemented a visual **scrollbar** with mouse-clickable hitboxes for long lists.
- **Statistics:** Added a **Maturity Distribution** chart (New vs. Young vs. Mature cards).
- **Polish:** Improved layouts in Settings and Browser views for better spacing on various screen sizes.

### Logic & SRS
- **Scheduler:** Added comprehensive unit tests for the SRS scheduler logic.
- **Robustness:** Fixed a regression where deck selection didn't update the review counts correctly.

### Verification
- All 163 E2E tests passed.
- All Go unit tests passed.
- `./scripts/verify.sh` executed successfully.

## 2026-05-12 10:30 +07 (Autonomous Improvement Pass)

Completed comprehensive improvement pass with 6 new German content decks, grammar expansion, UI polish, and testing. Total E2E tests: 226 (all passing).

### Content Expansion
- Added **A2 Daily Life & Household** vocabulary deck (52 cards) - daily routines, household chores, home objects.
- Added **B2 Media & News** vocabulary deck (50 cards) - journalism, broadcasting, media vocabulary.
- Added **C1 Academic & Scientific German** vocabulary deck (52 cards) - research, methodology, academic contexts.
- Added **A2 Shopping & Services** vocabulary deck (52 cards) - commerce, services, retail vocabulary.
- Added **B2 Business & Workplace** vocabulary deck (52 cards) - professional, corporate, employment.
- Added **C2 German Legal & Juridical** vocabulary deck (52 cards) - law, legislation, judicial terms.
- Expanded **Grammar Tips** with 5 new topics (Präteritum, Konjunktiv II wishes, Imperative forms).

### UI/UX Improvements
- **Statistics View:** Enhanced grade chart with visual icons (✗ for Again, ~ for Hard, ✓ for Good, ★ for Easy).
- **Browser Preview:** Fixed variable scoping issue where `kind` variable was referenced outside its scope.

### Testing
- Added comprehensive E2E test file `test_new_decks_may12.py` with 11 tests for new decks.
- Verified all new decks are accessible and vocabulary is searchable.

### Verification
- Total E2E tests: 226 (all passing).
- All Go unit tests passing.
- `./scripts/verify.sh` executed successfully with zero errors.

## 2026-05-12 11:00 +07 (Autonomous Improvement Pass - Final)

Completed comprehensive improvement pass with 9 new German content decks, UI enhancements, bug fixes, and testing expansion. Total E2E tests: 230 (all passing).

### Content Expansion (468+ new vocabulary cards)
- **A2 Daily Life & Household** (52 cards) - daily routines, household chores, home objects, visitors, rooms
- **B2 Media & News** (50 cards) - journalism, broadcasting, media vocabulary, credibility, digital media
- **C1 Academic & Scientific German** (52 cards) - research methodology, academic writing, citations
- **A2 Shopping & Services** (52 cards) - commerce, retail, services, pricing, measurements
- **B2 Business & Workplace** (52 cards) - corporate, employment, contracts, business operations
- **C2 German Legal & Juridical** (52 cards) - law, legislation, judicial proceedings, rights
- **B2 Travel & Tourism** (52 cards) - travel, accommodation, transportation, attractions
- **B2 Environment & Sustainability** (52 cards) - ecology, climate change, renewable energy, pollution
- **C1 German Science & Technology** (52 cards) - physics, chemistry, biology, astronomy, research
- **B1 Technology expanded** with 16 additional cards (AI, robotics, programming, networks, IT)

### Grammar Expansion
- Added 5 new grammar tips: Präteritum von sein, Präteritum von haben, Konjunktiv II for Wishes, Imperative Forms, Reflexive Pronouns in Dative

### UI/UX Improvements
- **Statistics View:** Enhanced grade chart with visual icons (✗ for Again, ~ for Hard, ✓ for Good, ★ for Easy)
- **Browser Preview:** Fixed variable scoping issue where `kind` variable was referenced outside its definition scope
- **AI Suggestions:** Maintained existing suggestions for E2E test compatibility

### Testing
- Added comprehensive E2E test file `test_new_decks_may12.py` with 15 tests covering all new decks
- All tests verified to pass with `./scripts/verify.sh`

### Verification
- Total E2E tests: 230 (all passing)
- All Go unit tests passing
- `./scripts/verify.sh` executed successfully with zero errors

### Commit Summary
- 17 files changed, 1031 insertions(+), 21 deletions(-)
- 10 new content deck files created
- 1 new E2E test file with 15 tests
- UI fixes in 2 files (render_browser.go, render_statistics.go)

## 2026-05-12 (Review TSV Header Regression Fix)

Fixed headered embedded TSV parsing so review cards no longer treat `front/back/extra/tags/notetype` headers or `Literal:` explanation fields as card content. The B1 idioms Rome reverse card now answers `Alle Wege führen nach Rom`, malformed B1 idioms proverb data was corrected, and a SQLite migration repairs already-seeded header/Rome artifacts.

### Verification
- `go test ./internal/content` passed.
- `./scripts/verify.sh` passed with 238 E2E tests.

## 2026-05-12 (Review Grade Key Stuck-State Regression Fix)

Fixed a review-key regression where pressing Space/Enter after the answer was already revealed set `gradingInProgress` without recording a review, which blocked subsequent `a/h/g/e` grading keys. Extra reveal keys on a revealed card now leave grading available and show the explicit grade-key hint.

### Verification
- `go test ./internal/tui` passed.
- `./scripts/verify.sh` passed with 238 E2E tests.
