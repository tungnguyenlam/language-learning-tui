# Active Backlog

Last updated: 2026-06-11

## Current Milestone

Polish and Reliability - IN PROGRESS.

## Completed Work

- [x] **Spotlight Dictionary State Reset:** Centralized dictionary search/detail reset behavior and routed `=`/Dictionary tab overlay open-close through shared helpers so stale search results, detail scroll, and detail mode do not leak between overlay sessions.
- [x] **Spotlight Dictionary Mouse Hitbox Scope:** Gave overlay-only search, history, and result hitboxes distinct IDs and active-view scope, with unit coverage for tab toggling, clear actions, and history mouse targets.
- [x] **E2E Navigation Recertification:** Updated mouse-tab/sidebar coordinates and Browser tab count in E2E tests after Dictionary was removed from the tab cycle and converted to a Spotlight overlay.
- [x] **Dictionary Search History Clear:** Added a clear button `[Clear]` (with mouse hitbox and `ctrl+x` keybinding) to clear recent search queries when the dictionary search bar is empty. Included `TestDictionaryClearSearchHistory` unit test coverage.
- [x] **Practice Hub Session Scores:** Enhanced the Practice Hub view to display the current session score (e.g. `• 4/5 (80%)`) next to each practice trainer button.
- [x] **Practice Hub Reset Scores:** Added an `r` keybinding to reset all practice session scores in the Practice Hub view. Included unit and E2E coverage for scores display and score resetting.

- [x] **Dictionary Quick Add Card Generation Bugfix:** Resolved a critical bug where quick-adding a dictionary entry to a deck (`ctrl+a`) only created the Note record without generating its associated study/review Card, leaving it permanently inactive/unreviewable. Integrated `content.CardsForNote` and added unit test coverage in `render_dictionary_test.go`.
- [x] **OpenAI & Anthropic Base URLs in Settings TUI:** Exposed both OpenAI and Anthropic Base URLs in the settings TUI list, allowing users to configure them directly (e.g. for offline Ollama or alternative API proxies). Updated settings navigation boundaries, cursor handlers, credential editing cases, Lip Gloss line info builders, and added `TestSettingsBaseURLEditing` unit test coverage.
- [x] **Dictionary Search Bar Clear Button:** Added an interactive `[x]` clear button on the right side of the dictionary search bar when a query is entered. Registered a mouse click hitbox to clear the query, reset results/cursor/scroll, and added `TestDictionarySearchClearHitbox` unit test coverage.
- [x] **Conjunctions & Word Order Trainer:** Implemented a new 9th practice mode in the Practice Hub to master German conjunction word order rules. Includes 15 detailed exercises with structural syntax validation, scoring, and clear grammar explanations.
- [x] **Dictionary Search History:** Added recent search query tracking in the Offline Dictionary view. Keeps up to 5 unique recent searches, rendering them as interactive click-to-search hitboxes when the search input is empty.
- [x] **Dictionary Unicode Highlighting Fix:** Reworked `highlightQuery` in `internal/tui/render_dictionary.go` to match and slice by runes instead of byte offsets, preserving German multi-byte characters like `Ä`, `ä`, and `ß` in styled Dictionary result highlights. Added regression coverage in `render_dictionary_test.go`.
- [x] **Dictionary Detail Scroll Clamp:** Replaced the Dictionary detail panel's hardcoded `shift+down` scroll buffer with the actual visible row count derived from `activeViewContentLayout()`, preventing detail scroll state from drifting beyond the renderable panel. Added a unit test for the clamp.
- [x] **B1 Email Phone Communication Deck:** Added a 40-note embedded TSV deck covering practical email, office message, and phone-call vocabulary/phrases. Added registry coverage to verify the deck loads with email and phone tagged notes.
- [x] **Review State Reset & Deck Switching Fix:** Added `resetReviewState()` helper in `internal/tui/handlers.go` to consolidate review state cleanup (reveal, typing, fix, focus, hint, grammar). Refactored `up`/`down` navigation in Review view to use it. Fixed deck switching (`[`/`]`) on Review view to reset review state and update status, preventing stale card state when switching decks mid-review.
- [x] **Contextual Help Hints for Practice Sub-views:** Added footer help hints for all 8 practice sub-views in `internal/tui/model.go` (Gender, Conjugation, Case, Adjective, Preposition, Plural, Separable, Numbers), plus added missing hints for Statistics (`x` export), Cram active (`q` quit), and Cram idle (`Enter` start).
- [x] **Dictionary in Dashboard Quick Actions:** Added "Dictionary" (`/`) to the Quick Actions row on the Dashboard for one-click access alongside Review, Practice, Cram, etc.
- [x] **SessionSummary State Leak Fix:** Changed the `reviewRecordedMsg` handler (`internal/tui/model.go:890`) to route through `m.updateView(ViewSessionSummary)` instead of directly setting `m.activeView`. Previously, the auto-transition after the last review bypassed `updateView`, leaving stale state (typingMode, fixProposal, showHint, etc.) alive in the summary view.
- [x] **Context-Sensitive helpHint for Missing Views:** Added contextual footer hints for Dictionary, Decks, Statistics, and Practice views in the `View()` method (`internal/tui/model.go:1098-1150`). Previously, these four views showed no contextual guidance in the footer, while all other navigable views did.
- [x] **A1 House & Furniture Deck:** Added a new 50-note A1 vocabulary deck covering rooms (Wohnzimmer, Küche, Bad), furniture (Tisch, Stuhl, Bett, Schrank), kitchen items (Herd, Kühlschrank, Spüle), and household verbs (wohnen, putzen). Registered in `StandardDecks()` and added a Go test in `new_decks_batch8_test.go`.
- [x] **Agent Index Refresh:** Updated `docs/agent/index.md` date from 2026-05-03 to 2026-06-10 and refreshed verification status entries.
- [x] **applyOverlay UTF-8 Corruption Fix:** Fixed `applyOverlay` (`internal/tui/model.go:1417`) which sliced overlay strings by byte offset while content may contain multi-byte UTF-8 characters (umlauts, emojis). Changed to use `[]rune` conversion so character boundaries are respected during overlay placement, preventing garbled characters in confirmation dialogs.
- [x] **State Leak Fixes — typingMode, showCardInfo, fixProposal:** Fixed `updateView` (`internal/tui/handlers.go:36`) to reset `typingMode`, `typedAnswer`, `typingChecked`, `typingCorrect`, `showHint`, `showCardInfo`, `showGrammarHint`, `focusMode`, and all `fix*` fields when switching views. Previously, typing mode would leak across views via number-key or tab navigation, causing printable keys to be swallowed. Also fixed card up/down navigation in Review view (`internal/tui/keys.go:167-186`) to reset `typingMode`/`showCardInfo` state so typed answers don't carry to adjacent cards.
- [x] **Cram Session State Leaks:** Fixed `updateView` Cram entry (`internal/tui/handlers.go:64`) to reset `cramActive`, `cramRevealed`, `revealState`, and `revealProgress` so re-entering Cram via number keys doesn't land in a stale active session. Also fixed cram exit handler (`internal/tui/keys.go:1124`) to reset `cramRevealed`, `revealState`, and `revealProgress` when exiting via `q`/`esc`.
- [x] **Practice Hub Numbers Trainer Mouse Click:** Added missing `case "numbers"` to the `practice-` hitbox handler in `internal/tui/hitboxes.go:93` so clicking the "Numbers & Time" button in the Practice Hub actually activates the trainer.
- [x] **Dictionary UX & UI Enhancements:** Implemented a two-column layout for the Dictionary view (list + details), added support for scrolling through search results with PgUp/PgDn, and introduced a "Quick Add" feature (ctrl+a) to instantly save dictionary entries to a dedicated deck. Improved search robustness with a request ID mechanism to prevent race conditions.
- [x] **Dictionary Audio & Find in Decks:** Added audio pronunciation to dictionary entries (`ctrl+p`), find in decks search (`ctrl+f`), and gender colorization (`{m}` blue, `{f}` pink, `{n}` green) for better visual parsing. Fixes to E2E tests included.
- [x] **Local TUI Dictionary & Seeding:** Implemented a new "Local TUI" dictionary provider that integrates lookups directly into the terminal. Updated standard content seeding to populate the local dictionary from flashcard notes. Refined the Dictionary view with a styled search bar and better interactivity. Added a robust E2E test suite for verification.
- [x] **E2E Practice Trainers Tests:** Created a new E2E test file (`test_practice_hub_extra.py`) verifying both the Separable Verb Trainer (key '7') and the Noun Plural Trainer (key '6') with interactive input and correctness verification.
- [x] **UI Polish & Content Enhancements:** Enhanced header with accent styling and view name, bordered Quick Actions box, styled footer, Dictionary search highlighting, Statistics goal indicators and heatmap legend, Numbers thousands/ordinals, expanded preposition exercises, NativeTTS error messages, dictionary LIKE fallback search, and relocated nav/tabs to render_views.go.
- [x] **E2E Fixes for UI Polish:** Fixed 7 E2E tests broken by UI changes: updated header assertions for new "DEUTSCH-TUI │ VIEW" format, Quick Actions label from colon to bordered box, increased terminal height for Grammar Tip/Verb visibility, and updated statistics visible line count.
- [x] **Dictionary in Tab Cycle:** Added Dictionary view to `nextViewCmd`/`previousViewCmd` navigation cycle so `s`/`tab`/`→` from Dashboard navigates to Dictionary. Updated Dictionary key handler to let `tab`/`right`/`s`/`w`/`left`/`shift+tab` pass through for navigation. Fixed 5 E2E tests for updated navigation order.
- [x] **Milestone 4 (Hybrid Dictionary):** Transformed the app by adding a unified dictionary feature. Built the Core API (`DictionaryEntry`, `DictionaryRepository`), FTS5 SQLite tables for fast indexing, and a dict.cc dataset parser. Added a real-time `ViewDictionary` TUI screen accessible via `/` from the dashboard, seamlessly integrated with the AI Drafting flow.
- [x] **Session Summary Grade Distribution:** Added a visual grade distribution chart (Again/Hard/Good/Easy) to the Session Summary view. Session grades are now tracked during reviews and displayed with colored bars and icons.
- [x] **Dashboard Word of the Day:** Added a "Word of the Day" vocabulary box to the Dashboard, displaying a rotating German noun with its article, plural form, and example sentence. Appears alongside the existing Grammar Tip and Verb of the Day when space allows.
- [x] **Review Session Speed & ETA:** Enhanced the Review header's session progress indicator to show cards-per-minute rate and estimated time remaining based on current review speed.
- [x] **Practice Hub Mouse Interactivity:** Fixed a bug where clicking Practice Hub buttons failed to enter the practice subviews, and added full mouse click interaction to Gender Trainer (der/die/das buttons) and click-to-continue hitboxes for all other trainers.
- [x] **Noun Plural Trainer:** Added a sixth practice mode to the Practice Hub for mastering German noun plural forms with article normalization, automatic extraction from deck notes, and fallback high-quality singular-plural noun lists.
- [x] **Preposition Trainer:** Added a fifth practice mode to the Practice Hub for mastering German two-way prepositions (Wechselpräpositionen) and their correct cases (Dative vs Accusative).
- [x] **Responsive Practice Hub:** Made the vertical spacing in the Practice Hub dynamic, allowing all 5 options to display correctly on smaller terminal sizes without overflow.
- [x] **Adjective Ending Trainer:** Added a fourth practice mode to the Practice Hub for mastering German adjective endings (weak, strong, and mixed declensions) with interactive fill-in-the-blank sentences and dynamic grammar contexts.
- [x] **Leech E2E Test Fix:** Resolved a false-positive E2E timeout in `test_leech_detection_in_statistics` by updating the screen expectation from `"cards due"` (which disappears when revealed) to `"a Again"` (present on the grade selection overlay).
- [x] **Practice Hub & Case Trainer:** Consolidated interactive modes into a central hub and added a new "Case Ending Trainer" for mastering German grammar.
- [x] **Verb Conjugation Trainer:** Added a new interactive mode for practicing German verb conjugations (ich, du, er/sie/es, etc.), including automated person selection and session scoring.
- [x] **Dashboard & UI Polish:** Improved Quick Actions layout, fixed Settings navigation boundaries, and resolved race conditions in daily goal adjustment.
- [x] **AI Feature Improvements:** Implemented automatic provider activation on API key entry, added inline tag extraction for AI drafting, and updated Anthropic model defaults.
- [x] **Gender Trainer Practice Mode:** Added a new interactive mode for practicing German noun genders (der/die/das), including automated noun extraction and session scoring.
- [x] **Cloze UI Polish:** Highlight revealed cloze answers in standard review mode.
- [x] **AI Tutor Explanation:** Added `Shift+H` shortcut to get card explanations from AI.
- [x] **Code Deduplication:** Introduced `executeBulkAction` and `executeSingleAction` helpers in `internal/tui/handlers.go`, reducing repetitive boilerplate for browser and review actions.
- [x] **Rendering Consistency:** Centralized `formatReviewInterval` and `renderReviewHistory` in common files (`utils.go` and `render_views.go`) for consistent card metadata display across views.
- [x] **RenderContext Enhancements:** Added `RegisterAction` and `WriteAction` to `RenderContext` to streamline the implementation of interactive UI elements.
- [x] **Dashboard Interactivity:** Added individual hitboxes for recently studied decks and interactive hitboxes for the 'Quick Actions' section.
- [x] **New Content: False Friends Mastery:** Added a new deck with 12 high-quality cards covering German-English false friends.
- [x] **E2E Improvements:** Added `e2e_tests/test_may23_improvements.py` to verify dashboard interactivity and new content.
- [x] **NativeTTS Cache Fix:** Fixed macOS cache mismatch where synthesis output target used `.aiff` while existence checks queried `.wav`.
- [x] **Browser Audio Playback & Select All:** Added `p` (play audio) and `a` (select/deselect all) shortcuts to the Card Browser, with unit tests and help overlay updates.
- [x] **Grammar Hint Overlay:** Added a `Shift+G` shortcut in the Review view to toggle a contextual grammar hint based on the card contents.
- [x] **List Navigation Shortcuts:** Added `g` (jump to top) and `G` (jump to bottom) shortcuts to Browser, Settings, and Decks views for faster navigation.
- [x] **Browser Contextual Hints:** Integrated the grammar hint analysis directly into the Card Preview box in the Browser view to show contextual word usage information.
- [x] **Interactive Dashboard Elements:** Made Word of the Day, Verb of the Day, and Grammar Tips actionable. Users can now press 'w' to add the word to their collection, 'v' to practice the verb, or 'g' to search for related grammar in the browser.
- [x] **Separable Verb Trainer:** Added a seventh practice mode to the Practice Hub for mastering German separable verbs and their prefixes in sentences.
- [x] **Review Focus Mode Polish:** Enhanced the 'Shift+F' Focus Mode in the Review view with a centered, distraction-free layout and clearer visual feedback, while maintaining compatibility with E2E tests.
- [x] **Practice Hub Navigation & Badge Counts:** Standardized the `esc` key behavior across all 7 practice modes (always clears current input or gracefully returns to the Hub, rather than the dashboard). Added dynamic, parallel-loaded item counts next to each trainer in the Practice Hub to display the number of practice items available.

## Exact Next Action

Await next user feedback or feature request.

## Top Issues / Priorities

None.

## Last Verified

- `PATH="/opt/homebrew/bin:$PATH" ./scripts/verify.sh` passed on 2026-06-11: Go tests, smoke test, binary build, and 349 E2E tests all passed.

## Blockers

None.
