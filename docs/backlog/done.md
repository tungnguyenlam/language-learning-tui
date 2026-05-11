# Done Backlog

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
