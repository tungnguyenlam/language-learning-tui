# Active Backlog

Last updated: 2026-08-01

## Current Milestone

Dictionary Feature Enhancements & Lemmatization Search.

## Exact Next Action

Pick the next dictionary milestone item after this verified bug/performance pass.

## Top Issues

- [x] `Store.Reset` no longer references the absent `statistics` table; added a regression test.
- [x] `UndoLastReview` and `GetReviewState` now preserve `last_review_at`; added a regression assertion.
- [x] Browser and related-dictionary results carry request identity, and dictionary `ctrl+f` now preserves the lookup query without duplicate loads.
- [x] Normal dictionary form matching uses a companion FTS5 index; quick-add uses loaded deck metadata, and SQLite child-table indexes were added in migration 26.

No active issues remain from this pass.

## Acceptance Criteria

- Reset succeeds on a migrated SQLite store and clears content/progress without referring to absent tables.
- Undo restores both the prior scheduling values and prior last-review timestamp.
- Out-of-order Browser/Dictionary results are ignored.
- Dictionary search keeps inflection matching while avoiding the current unindexed form scan for normal searches.
- `./scripts/verify.sh` passes and the completed work is recorded in `docs/backlog/done.md`.

## Last Verification

- `go test -race ./internal/storage/sqlite`
- `go test -race ./internal/tui`
- `./scripts/verify.sh` passed on 2026-08-01: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed).

## Completed Work

- [x] **Bug Fix: Dictionary Search Mid-Typing Freeze:** Fixed TUI app freeze while typing in the dictionary search bar. `findRelatedEntries` (triggered automatically on every search result when screen width > 80 cols) was executing an unindexed `LIKE '%word%'` query across the 834,512-row `dictionary_fts` SQLite table, locking the database connection and blocking the event loop on keystrokes. Converted `FindRelatedEntries` and `Exists` in `internal/storage/sqlite/dictionary.go` to use fast indexed FTS5 `MATCH` queries (executing in <1ms instead of 1,000ms full table scans).

- [x] **Bug Hunt: Cram 1-4 Digit Grading, FTS Fallback & Trainer Esc:** Added `1`, `2`, `3`, `4` digit grading shortcuts to Cram mode when revealed, with UI hint parity `(1)`, `(2)`, `(3)`, `(4)` matching standard Review mode; ensured SQLite dictionary search gracefully falls back to LIKE queries if FTS parsing fails on special/wildcard tokens; fixed `Esc` key navigation on revealed practice trainer items so it directly exits to the Practice Hub instead of consuming an Esc press on submitted text. Added regression unit tests `TestCramDigitGradingKeys`, `TestTrainerEscKeyOnRevealedCard`, and `TestDictionarySearchFTSFailureFallback`.

- [x] **Bug Hunt: Session Steal, Paste Gaps, Dictionary Input/Wheel:** Fixed digit shortcuts abandoning active Cram mid-session; clipboard paste into Dictionary/Spotlight, Review typing, practice trainers, and AnkiWeb query; mouse-wheel (+ scrollbar drag) for Dictionary results/detail and AnkiWeb; dictionary `d` key stealing search typing (`der`/`das`) after the detail-toggle fix. Regression tests in `bugfix_input_test.go`.

- [x] **Bug Hunt: Input Routing & Browser/AI Safety:** Fixed Review typing `q` quitting; Cram grading before reveal; Practice reveal-advance stealing `q`/`?`; Relative trainer `=` opening Spotlight; Gender trainer `0` jumping views; Browser tag filter substring false positives; AI Esc not cancelling drafting; Browser multi-select surviving deck switches. Added `practiceBlocksGlobalShortcut()`, whole-tag SQL match, `draftCancelled`, and selection clearing/visibility filters. Regression tests in `bugfix_input_test.go`, `trainer_test.go`, `sqlite_test.go`.

- [x] **Spotlight Dictionary Detail Parity & Dynamic Height Scaling:** Extended `renderSpotlightDictionary` in `render_dictionary.go` to render complete inflection tables, compound breakdowns (with `dict-overlay-compound-` hitboxes), example sentences, and related words (with `dict-overlay-related-` hitboxes) in both single-column and 2-column Spotlight overlay modes. Made Spotlight height scale dynamically with terminal height (`boxHeight = m.height - 8` on 30+ row screens). Added `TestSpotlightDictionaryDetailParity`.
- [x] **Dictionary Inflection / Grammar Flashcard Quick-Generator (`i` / `ctrl+i`):** Added `addDictionaryInflectionEntryCmd()` in `actions_dictionary.go` and `i` / `ctrl+i` shortcuts in `keys.go` to generate dedicated Grammar Inflection flashcards (`[Grammar] Forms of "..."`) containing word class, gender declensions, and inflection tables. Updated dictionary footer help hints and added `TestDictionaryInflectionCardGeneration`.

- [x] **Dictionary ctrl+e Real Explain Flow:** `ctrl+e` now calls `ai.ExplainDictionaryEntry` (pedagogical tutor explanation) instead of `startDrafting()` flashcard generation. Results render in the AI view with dismiss via `H`/`Esc`. Added `TestExplainDictionaryEntry` and updated `TestDictionaryCtrlEExplainAndFilterTagHelpers`.
- [x] **Two-Column / Spotlight Domain Tag Badges:** Shared `formatDictionaryDomainTags()` so wide two-column dictionary detail and Spotlight detail show the same `[ZOOL.]` / field badges as single-column detail.
- [x] **Clear Recently Inspected Words:** Added `[Clear]` hitboxes (`dict-recent-clear` / `dict-overlay-recent-clear`) and made `ctrl+x` clear recently inspected when search history is already empty. Persists via `clearDictionaryRecentlyViewed()`.
- [x] **AI Draft Row Hitboxes & Context Banner:** Draft text rows are clickable (`draft-row-N`); dictionary context banner shows the headword (`★ Dictionary Context: geheim`) instead of a bare label. Added `TestAIDraftRowHitboxesAndDictionaryContextBanner`.
- [x] **Filter-Only Browse Alphabetical Order:** Filter-pill browse (`:noun`, `:verb`, etc.) uses `ORDER BY word` plus stable in-memory sort. Added `TestDictionaryFilterOnlyBrowseOrdered`.
- [x] **Related Words Compound-Aware Matching:** `FindRelatedEntries` prefers lemma / prefix (≥4 runes) / suffix compounds and drops short-stem middle noise (e.g. `und` → `undurchlässig`). Extended `TestFindRelatedEntries`.

- [x] **Multi-Part German Compound Decomposition & SQLite Dictionary Validation:** Enhanced `DecomposeCompound` in `internal/content/compound.go` to recursively decompose German compound words with 3+ components (e.g. `Krankenhausarzt` -> `Kranken` + `Haus` + `Arzt`). Added `Exists` to `DictionaryRepository` interface and `sqlite/dictionary.go` to validate sub-compounds against the local dictionary DB. Updated single-column and two-column dictionary views in `render_dictionary.go` to display multi-part breakdowns with interactive hitboxes (`dict-compound-sc-` and `dict-compound-tc-`) and expanded number key shortcuts (`1`-`5`) in `keys.go` for all compound parts. Added `TestMultiPartDecomposeCompound` and `TestMultiPartCompoundDecompositionAndHitboxes`.
- [x] **AI Flashcard Draft Batch Operations (`A` / `D` & Hitboxes):** Added `approveAllDrafts()` and `discardAllDrafts()` batch operations for AI flashcard generation. Added `A` / `ctrl+a` (Approve All) and `D` / `ctrl+d` (Discard All) keyboard shortcuts in `updateAIKey` (`keys.go`), rendered `[Approve All]` and `[Discard All]` action buttons with mouse hitboxes (`draft-approve-all` and `draft-discard-all`) in `render_ai.go`, and updated shortcut help footers. Added `TestAIApproveAllAndDiscardAllShortcuts`.

- [x] **Practice Hub `=` vs Dictionary Spotlight Conflict:** Global `=` spotlight handler no longer steals the Practice Hub Relative Clauses trainer shortcut. Outside the hub, `=` still opens Spotlight. Added `TestPracticeHubEqualsOpensRelativeNotDictionary`.
- [x] **Recently Inspected Words Auto-Record & Spotlight Parity:** Navigating/selecting dictionary entries with a visible detail pane now records them via `inspectDictionaryCursor()` / `updateDictionaryKey`; Spotlight empty state renders the same Recently Inspected section with clickable hitboxes.
- [x] **DE/EN Filter-Only Browse:** `Search()` treats bare `de:` / `en:` (and `:de` / `:en`) like other filter pills and returns a browse sample instead of an empty result set. Added `TestDictionarySearchLangFilterOnly`.
- [x] **Stable `:starred` Browse Order:** Starred-only dictionary browse sorts by headword (then ID) instead of map iteration order. Added `TestDictionaryStarredBrowseStableOrder`.
- [x] **Practice Hub Footer/Help Coverage:** Hub status footer documents `1-9,0,-,=`; Relative trainer shares the fill-in hint footer; `TestPracticeHubVisuals` asserts `Pv` and `Rl` icons.

- [x] **Dictionary Continuity & Safe Batch Add:** Persisted the "Recently Inspected Words" stack in local settings and restored it at startup; made the mouse `[Clear]` recent-search control persist its clearing action; changed `ctrl+s` dictionary batch-add to skip existing notes rather than rewriting them (with accurate added/skipped/failed feedback). Added regression coverage for persisted recent words, persistent mouse clearing, and duplicate-safe batch add.

- [x] **Cram Review Mouse Hitboxes:** Registered interactive mouse hitboxes `cram-grade-again`, `cram-grade-hard`, `cram-grade-good`, and `cram-grade-easy` in `render_cram.go` when `m.cramRevealed` is true. Added `TestCramGradeHitboxes`.
- [x] **Gender Trainer Options Hitbox Alignment:** Fixed `gender-opt-der`, `gender-opt-die`, and `gender-opt-das` mouse hitbox Y coordinates (`layout.Y + 8`) in `render_gender_trainer.go` to align with rendered option buttons. Added `TestGenderTrainerOptionsHitboxYOffset`.
- [x] **AnkiWeb Shared Deck Search Hitboxes:** Registered `ankiweb-result-` hitboxes on search result rows in `screen_ankiweb.go` so clicking any deck in the public library selects it and loads its details. Added `TestAnkiWebSearchResultHitboxes`.
- [x] **Responsive Help Overlay Grid:** Updated `renderHelp` in `render_views.go` to dynamically adapt shortcut column layout (4-column grid on wide screens, 2x2 grid on medium 60-110 col screens, single column stack on <60 col screens) to prevent right-edge clipping. Added `TestResponsiveHelpLayout`.
- [x] **Filter Pills Display Name Fix:** Fixed `renderFilterPillsRow` in `render_dictionary.go` to display nice human-readable pill labels (`★ Starred`, `DE`, `EN`, `Verb`, `Noun`, `Adj`, `Adv`, `Der`, `Die`, `Das`, `Pl`) instead of overriding `label` with raw filter query tags (`:starred`, `de:`, `:verb`, `:m`, etc.). Added `TestFilterPillDisplayNames`.
- [x] **Dashboard Wide-Mode Verb Conjugation Fix:** Fixed wide 3-column dashboard mode (>110 cols) in `render_dashboard.go` to render complete German 3rd person verb conjugation labels (`er/sie/es` and `sie/Sie`) matching 2-column mode. Added `TestWideDashboardVerbConjugationLabels`.
- [x] **Review Grade Hitbox Y Offset Fix:** Fixed `answerYOffset` calculation across standard flashcards, MCQ, and typing cards in `render_review.go` so that `extraDisplay` (card context and `renderGrammarHint`) is included when computing Y coordinates for Grade button hitboxes (`grade-again`, `grade-hard`, `grade-good`, `grade-easy`), preventing mouse hitboxes from misaligning when card context or grammar hints are present. Added `TestReviewGradeHitboxWithExtraContext`.
- [x] **Non-Truncating Footer Shortcut Wrapping:** Implemented `formatWrappedFooter()` in `model.go` to parse shortcut help hints into discrete tokens and wrap them cleanly across multiple lines when terminal width is constrained. Eliminates trailing `...` truncation, ensuring 100% of shortcut commands are visible on any screen width. Added `TestFormatWrappedFooterNoTruncation`.
- [x] **Automatic Viewport Scrollbars & Full Vertical Space Utilization:** Updated `activePanelSize()` in `model.go` (`m.height-4`) to eliminate dead blank space at the bottom of the screen. Implemented `AutoScrollViewport()` in `utils.go` to automatically slice content, attach scrollbars, and register hitboxes whenever any view content exceeds layout height. Applied `AutoScrollViewport` to the Practice Hub (`render_practice.go`), adding `m.practiceScroll`, mouse wheel scrolling, cursor auto-scrolling, and `practice-scroll-` track hitboxes. Added `TestPracticeHubAutoScrollbarAndFullHeight`.
- [x] **Unified Interactive Scrollbar Architecture & Bounds Fix:** Standardized scrollbar rendering across all views (`RenderList` and `renderScrollbarColumn`), using identical track (`│`) and thumb (`█`) characters with `colorPanel`/`colorAccent` styling. Fixed negative fill loop indices (`i := -N`) when scrolling past bounds in `render_dictionary.go`. Registered interactive mouse drag & click hitboxes (`dict-scroll-` and `dict-detail-scroll-`) in `hitboxes.go` and `render_dictionary.go`, making dictionary scrollbars fully interactive alongside Decks, Browser, Cram, Settings, and Statistics. Added `TestDictionaryScrollbarLineCountAndOverscroll`.
- [x] **Dictionary Entry Headword Consolidation:** Added `ConsolidateDictionaryEntries()` in `internal/core/core.go` to group multiple dict.cc dictionary entries for the same German headword and gender into a single entry with combined translations (e.g. `strong; powerful; heavily`), merged word classes (`adj, adv`), deduplicated forms, domain tags, and example sentences. Applied consolidation across `ParseDictCCStream`, SQLite `ImportEntries`, and `Search` result sets. Added `TestConsolidateDictionaryEntries` in `core_test.go`.
- [x] **Duplicate Card Warning on Quick-Add:** Updated `addDictionaryEntryCmd`, `addDictionaryReverseEntryCmd`, and `addDictionaryClozeEntryCmd` in `actions_dictionary.go` to check if a card already exists before upserting, warning the user and requiring an immediate second press to force an update.
- [x] **Related Words & Cross-References:** Added `FindRelatedEntries` to `DictionaryRepository` and `sqlite/dictionary.go`. Displays clickable compound/root related words in detail panel (`render_dictionary.go`) with auto-fetching on selection change.
- [x] **German Compound Word Decomposition:** Added `DecomposeCompound` in `internal/content/compound.go` with Fugenelemente handling and dictionary validation. Renders `Compound Breakdown:` in the detail view.
- [x] **Dictionary Random Discovery (Explore Mode):** Added `RandomEntries` in `sqlite/dictionary.go` and `dictDiscoverEntriesMsg` in `actions_dictionary.go`. Renders interactive `✦ Discover:` section in empty search state for both full view and spotlight overlay.
- [x] **Dictionary One-Key Cloze Flashcard Generation (`c` / `ctrl+c`):** Added `addDictionaryClozeEntryCmd()` in `actions_dictionary.go` and `c` / `ctrl+c` key shortcuts in `keys.go`. Extracts best example sentences from dictionary entries and converts target headwords into `{{c1::<word>}}` Cloze flashcards (or formats fallback Cloze prompts) saved directly to the selected target deck. Added `TestDictionaryClozeCardGeneration`.
- [x] **Recently Inspected Words History Stack & Domain Tag Badges:** Added `dictionaryRecentlyViewed` stack in `Model` and `recordDictionaryView()` in `actions_dictionary.go`. Renders interactive `Recently Inspected Words:` section when dictionary search is cleared. Added styled domain/field tag badges (e.g. `[ZOOL.]`, `[COLLOQUIAL]`, `[MEDICINE]`) in detail panel view (`render_dictionary.go`). Added `TestDictionaryRecentlyViewedAndDomainTags`.
- [x] **Dictionary Entry Star/Bookmark System & `:starred` Filter (`b` / `ctrl+b`):** Added `dictionaryStarred map[string]bool` in `Model`, `toggleStarDictionaryEntry()`, `loadDictionaryStarred()`, and `:starred` filter tag support in `actions_dictionary.go` and `sqlite/dictionary.go`. Renders bright yellow `★` stars in search lists and detail view headers. Starred entries are persisted to SQLite settings (`dict_starred_entries`). Added `[★ Starred]` pill to filter row and `TestDictionaryStarringAndFiltering`.
- [x] **Dictionary Result List Translation Snippets in 2-Column Mode:** Updated `renderDictionary` and `renderSpotlightDictionary` (>80/70 cols) in `render_dictionary.go` so left result items show `Word {gender} - Translation` instead of just headwords, allowing instant scanning of search result meanings. Added `TestDictionaryTwoColumnTranslationDisplay`.
- [x] **Dictionary Target Deck Selection (`ctrl+g`) & Batch Add (`ctrl+s`):** Added `dictionaryTargetDeckID` field in `Model`, `cycleDictionaryTargetDeck()` (`ctrl+g`) to cycle through available decks for dictionary card creation, and `addDictionaryEntriesBatchCmd()` (`ctrl+s`) to batch-add all current dictionary search results directly into the selected target deck. Added `TestDictionaryTargetDeckCyclingAndBatchAdd`.

- [x] **Dictionary Input Fix for 'k' and 'j' Keys:** Updated `updateDictionaryKey()` in `internal/tui/keys.go` so that `'k'` and `'j'` are typed as input characters when the search bar is active (`!m.dictionaryFocusResults && !m.dictionaryDetailView`), preventing search history cycling or list navigation when typing words containing 'k' or 'j'. Added `TestDictionaryKAndJKeyHandling` in `render_dictionary_test.go`.

- [x] **Dictionary Wildcard (* & ?) & Language Scope (de: & en:) Search:** Updated `parseSearchFilters`, `filterEntries`, and `Search` in `internal/storage/sqlite/dictionary.go` to convert wildcard queries into SQLite `LIKE` queries (`*` -> `%`, `?` -> `_`), enabling suffix searches like `*ung`, `*keit`, `*ieren` and language scope filters (`de:` for German headwords/forms, `en:` for English translations). Added `TestDictionarySearchWildcardAndLanguagePrefix`.
- [x] **Interactive Dictionary Filter Pills:** Added `renderFilterPillsRow()` and filter tag helpers (`isFilterActive`, `toggleFilterTag`, `clearFilterTags`) in `internal/tui/render_dictionary.go` and `internal/tui/actions_dictionary.go`. Renders interactive clickable filter pills (`[All]`, `[:verb]`, `[:noun]`, `[:adj]`, `[:adv]`, `[:m]`, `[:f]`, `[:n]`, `[:pl]`) in both standard Dictionary view and Spotlight Dictionary overlay.
- [x] **Dictionary Entry Contextual AI Explainer (`ctrl+e`):** Added `ctrl+e` shortcut in `updateDictionaryKey` (`internal/tui/keys.go`) to automatically pre-populate the AI Tutor (`ViewAI`) with a structured request explaining the selected German dictionary entry (grammar notes, gender usage, collocations, and example sentences). Added unit test `TestDictionaryCtrlEExplainAndFilterTagHelpers`.
- [x] **Structured German Inflection Table Formatter:** Added `formatInflectionTable()` in `render_dictionary.go` to parse and render labeled grammar inflection tables for Verbs (Präsens 3sg, Präteritum, Perfekt), Nouns (Genitiv, Plural), and Adjectives (Komparativ, Superlativ) in dictionary detail views and Spotlight overlay. Added `TestFormatInflectionTable`.
- [x] **Interactive Compound Breakdown Mouse Hitboxes & Number Key Jump Shortcuts (`1`-`5`):** Registered interactive mouse hitboxes `dict-compound-sc-` and `dict-compound-tc-` in `render_dictionary.go` for compound breakdown parts (e.g. `Hand` + `Schuh`). Added number key shortcuts `1`-`5` in `keys.go` to jump directly to compound sub-words or related words when results or detail panel are focused. Added `TestDictionaryNumberKeyNavigation`.
- [x] **Dictionary Inflected Form Matching & Lemmatization Ranking:** Updated `Search()` and `entryMatchScore()` in `internal/storage/sqlite/dictionary.go` to search `entry.Forms` via `LIKE` alongside FTS matches and score exact form matches at Tier 2. Searching inflected forms (e.g. "ging", "gegangen", "Häuser") now correctly returns and prioritizes their lemma entries ("gehen", "Haus") above prefix/substring matches. Added `TestDictionarySearchInflectedFormsAndFilters`.
- [x] **Dictionary Tag & Part-of-Speech Filters:** Added query filter parsing (`parseSearchFilters`) to `Search()` in `internal/storage/sqlite/dictionary.go`, supporting `:verb` (`:v`), `:noun`, `:adj`, `:adv`, `:m`, `:f`, `:n`, `:pl` tag filters for targeted dictionary searches.
- [x] **Smart Query Sanitization for Contextual Dictionary Lookups:** Added `cleanLookupQuery()` in `internal/tui/actions_dictionary.go` to strip HTML formatting, Cloze tags (`{{c1::...}}`), parenthetical hints `(...)`, bracket tags `[...]`, bullet prefixes, and trailing punctuation when looking up words from Review, Browser, Cram, or Daily Tips. Added `TestCleanLookupQueryAndCardFormat`.
- [x] **Dictionary Quick-Add Flashcard Enrichment:** Updated `addDictionaryEntryCmd()` in `internal/tui/actions_dictionary.go` to automatically attach German gender articles (`der`, `die`, `das`) to noun fronts and format `Extra` card metadata with forms, word class, gender, and example sentences.
- [x] **Dict.cc Data Pipeline & Parsing Enhancements:** Updated `ParseDictCCStream` in `internal/content/dictcc.go` with HTML entity unescaping (`html.UnescapeString`), regex-based gender extraction (`{m}`, `{f}`, `{n}`, `{pl}`, `{m/f}`), word form extraction (`<...>`), and bracketed tag extraction (`[...]`), producing clean headwords and structured forms/tags. Added `TestParseDictCCStream_Enhanced`.
- [x] **Dictionary Search Multi-Tier Ranking & Article Stripping:** Updated `entryMatchScore` and `sortDictionaryEntries` in `internal/storage/sqlite/dictionary.go` with German/English article stripping (`der`, `die`, `das`, `the`, `a`, `an`), semicolon-separated translation matching, and 7-tier score prioritization so headword matches rank above example sentence substrings. Added `TestDictionarySearchArticleAndMultiTranslation`.
- [x] **Cram View Dictionary Lookup & Spotlight Overlay Auto-Close:** Added `lookupCramCardInDictionary()` in `actions_dictionary.go` with `d` key shortcut in Cram mode (`updateCramKey`). Updated Spotlight Dictionary navigation (`ctrl+f` for Browser, `enter` for AI Draft) in `updateDictionaryKey` to automatically close `m.dictionaryOverlayActive` when switching views.
- [x] **Relative Clauses Trainer (trainer #12):** New `PracticeSubViewRelative` mode added to the Practice Hub (icon `Rl`, key `=`). 25 fill-in-the-blank exercises in `internal/content/relative.go` covering Nominative (*der, die, das, die*), Accusative (*den, die, das, die*), Dative (*dem, der, dem, denen*), Genitive (*dessen, deren*), and Prepositional relative clauses (*mit dem, in der, an denen, worüber, für die, was* after indefinites). Added `TestGetRelativeExercises`, `TestRelativeTrainerIsRegistered`, `TestRelativeTrainerRenders`, and updated keybindings, hitboxes, footer, and help overlays.
- [x] **Passive Voice Trainer (trainer #11):** New `PracticeSubViewPassive` mode added to the Practice Hub (icon `Pv`). 25 fill-in-the-blank exercises in `internal/content/passive.go` covering Vorgangspassiv (Präsens, Präteritum, Perfekt), Zustandspassiv, modal passives, agent prepositions (von/durch), passive alternatives (sich lassen), and Futur I passive — each with meaning, hint, and detailed explanation. Added `TestGetPassiveExercises` and `TestPassiveTrainerIsRegistered`.
- [x] **Dashboard Wide-Mode Verb Conjugation Fix:** Updated `render_dashboard.go` so the 3-column wide layout (>110 cols) renders all 6 German verb forms (`ich`, `du`, `er/sie/es`, `wir`, `ihr`, `sie/Sie`) matching the standard 2-column mode.
- [x] **Dashboard Quick Actions Visual Consistency:** Removed opaque `Background("234")` from `dashActionsStyle` in `styles.go` so the Quick Actions block has a transparent background like all other dashboard boxes.
- [x] **Konjunktiv II Trainer (trainer #10):** New `PracticeSubViewKonjunktiv` mode added to the Practice Hub (key `0`, icon `Kj`). 25 fill-in-the-blank exercises in `internal/content/konjunktiv.go` covering würde+Infinitiv, strong Konjunktiv II forms (wäre/hätte/bliebe), modal subjunctives (könnte/sollte/dürfte), and Perfekt conditionals — each with a meaning, grammar hint, and full explanation shown on reveal.
- [x] **Preposition Trainer expanded (24→40 exercises):** Added 16 new exercises to `prepositions_practice.go` covering genitive prepositions (wegen, trotz, während, statt, innerhalb, außerhalb), temporal dative (seit, nach, vor, bis+zu), and accusative-only prepositions (für, um, ohne) to complement the existing two-way preposition pairs.
- [x] **Separable Verb Trainer expanded (15→30 exercises):** Added 15 more exercises in `trainer_content.go` covering abfahren, zurückdenken, anprobieren, teilnehmen, zurückrufen, vorstellen, aufhören, aufmachen, aufgeben, anziehen, aufräumen, mitbringen, ausfüllen, and vorbeikommen.
- [x] **Dictionary Search Exact Match Prioritization:** Updated `Search()` in `internal/storage/sqlite/dictionary.go` with stable tier sorting (`sortDictionaryEntries`), ensuring exact headwords, exact translations, and prefix matches take precedence over substring occurrences in example sentences. Added `TestDictionarySearchExactMatchOrdering`.
- [x] **Dictionary Keyboard History Cycling:** Added `cycleDictionaryHistory()` in `internal/tui/actions_dictionary.go` and updated `updateDictionaryKey()` in `keys.go` so `Up`/`Down` arrow keys navigate recent search history when search input is active/cleared. Added `TestDictionaryHistoryKeyboardCycling`.
- [x] **Active Deck Context for Dictionary Quick-Add:** Updated `addDictionaryEntryCmd()` in `actions_dictionary.go` to target `m.deck.ID` if a specific deck is selected, otherwise defaulting to the `"dictionary"` deck.
- [x] **Full Screen-Interface Migration:** Completed migration of all remaining views (`ViewDashboard`, `ViewDecks`, `ViewReview`, `ViewBrowser`, `ViewAI`, `ViewCram`, `ViewPractice`, `ViewDictionary`) to concrete `screen` types in dedicated files (`screen_dashboard.go`, `screen_decks.go`, `screen_review.go`, `screen_browser.go`, `screen_ai.go`, `screen_cram.go`, `screen_practice.go`, `screen_dictionary.go`). Simplified `renderActiveViewPlainAt` and `updateActiveViewKey` to delegate directly through the `m.screens` map.
- [x] **Dictionary Deck Sync on Card Add:** Updated `addDictionaryEntryCmd` in `actions_dictionary.go` to return a `tea.Batch` including `m.loadDecks` and `m.loadDueCards` so adding cards from the dictionary immediately updates deck state and review queues.
- [x] **Scrollbar variation, smoothness, and UI unification:** Unified scrollbar rendering across all views (`RenderList`, `render_dictionary.go`, `screen_ankiweb.go`) using `renderScrollbarColumn` with consistent track (`│`) and thumb (`█`) characters and Lip Gloss styling. Fixed thumb offset rounding (`+ maxScroll/2`) and track click mapping (`+ (visibleLines-1)/2`) in `utils.go` to eliminate jerky thumb position jumps. Added unit tests in `render_test.go`.
- [x] **Screen-interface migration for ViewSettings, ViewDebug, ViewStatistics:** Created `screen_settings.go`, `screen_debug.go`, `screen_statistics.go` implementing `screen` interface and registered them in `m.registerScreens()`.
- [x] **Practice Hub click alignment fix:** Fixed button hitbox Y spacing in `renderPracticeHub` (`spacing := 5`) so mouse clicks align precisely with rendered button blocks across all terminal heights. Added `TestPracticeHubHitboxSpacing` unit test.
- [x] **Decks view search reset:** Ensured `m.applyDeckFilter()` is called when clearing `m.deckFilter` with `<Esc>` key in Decks view so the deck list immediately refreshes to show all decks.

- [x] **`.apkg` export Anki can actually open:** the old exporter wrote notes
  with `mid=0`, matching no note type in `col.models` — structurally invalid
  for real Anki, and only round-tripping because our own importer ignored
  `mid`. `internal/content/apkg_schema.go` now declares three real note types
  and `apkg_export.go` writes a valid schema-11 collection: positive dense ids,
  declared decks, `col.conf.curDeck`, and a media index named `media`.
  `TestExportProducesAnkiValidCollection` asserts each invariant Anki checks.
- [x] **Import maps fields from templates, not position:** Anki note types are
  free-form, so field 0 is not reliably the front. `templateFieldOrds` reads
  the `{{Field}}` references out of a note type's `qfmt`/`afmt` and falls back
  to the first two non-empty fields. A real 4214-note shared deck imported with
  an empty back before this. Import also now reads zstd `collection.anki21b`,
  the modern `notetypes`/`fields`/`templates` schema, and strips HTML.
- [x] **Deck names survive export:** `exportAPKG` builds real Anki deck names
  from the repository instead of a single default deck, and reconstructs
  `Reverse` note types from swapped card pairs.
- [x] **AnkiWeb shared-deck browser:** press `A` from Import to search
  ankiweb.net's public library, read a deck's description and size, and
  download plus import it without leaving the terminal. `internal/ankiweb` is a
  read-only client over AnkiWeb's undocumented protobuf endpoints, decoded
  defensively by a hand-written wire reader; see the notice
  `docs/agent/notices/2026-07-26-ankiweb-undocumented-endpoints.md`. It is the
  app's only network surface, runs only on explicit user action, sends nothing
  about the user, and is deliberately outside the tab cycle. Rate limits and
  expired tokens are reported as normal outcomes that name the manual fallback.
- [x] **Trainer input no longer quits the app:** `q` and `?` are consumed as
  typed characters while a text-input trainer is waiting for an answer.
  German answers contain `q` (Qualität, Quelle), so typing one used to hit the
  global quit shortcut and exit mid-exercise. See `trainerInputActive()` in
  `internal/tui/keys.go`.
- [x] **Trainers reshuffle between passes:** `trainerState.advance()` and
  `Model.advanceGenderItem()` reshuffle after each completed pass, and the
  Gender Trainer's noun list is shuffled at load. Small exercise sets used to
  cycle in a fixed order forever, so learners recalled the position instead of
  the grammar. The first pass keeps its authored order.
- [x] **Trainer progress indicator:** headers now show `Item n/N` (`Noun n/N`
  for the Gender Trainer) and a `Round n` counter after the first pass.
- [x] **Grammar content fixes and expansion:** corrected three wrong grammar
  labels (`wegen des Wetters` is neuter not masculine; `in einem alten Haus` is
  mixed not weak declension; `warten auf` is a fixed case, not movement) and
  grew the Case and Adjective Ending trainers from 15 to 25 exercises each.
- [x] **Help overlay accuracy:** documented the `=` Dictionary spotlight, which
  was missing entirely, and replaced the stale `Practice 1-3/d-a/m-n` line with
  the real `1-9` trainer keys.

## Verification

- `./scripts/verify.sh` passed on 2026-07-31: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed) after Cram digit-nav, paste, dictionary `d` typing, and dictionary/AnkiWeb mouse-wheel bug fixes.

- `./scripts/verify.sh` passed on 2026-07-31: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed) after input-routing / browser / AI drafting bug fixes.

- `./scripts/verify.sh` passed on 2026-07-31: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed) after dictionary explain flow, domain-tag parity, recently-inspected clear, AI draft row hitboxes, filter-only ordering, and related-word matching.

- `./scripts/verify.sh` passed on 2026-07-27: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite (34 passed) after Practice Hub `=` fix, recently-inspected auto-record/spotlight parity, DE/EN filter-only browse, and starred browse ordering.

- `./scripts/verify.sh` passed on 2026-07-27: Go tests, vet, offline dict.cc import (834,512 entries), smoke test, binary build, and core E2E suite.

- `./scripts/verify.sh` passed: gofmt, all Go unit tests, `go vet`, smoke test,
  binary build, and the E2E suite (353 tests, including the two new
  `e2e_tests/test_trainer_input_shortcuts.py` cases).
- The Anki work is covered by `internal/content/apkg_anki_test.go`,
  `internal/ankiweb/ankiweb_test.go`, and
  `internal/tui/screen_ankiweb_test.go`. None of them touch the network:
  `ankiweb.BaseURL` is pointed at `httptest`, and the TUI tests use a stub
  `ankiWebSearcher`.
- The full chain was also validated once against production by hand — search,
  info, download of a 48 MB deck, and import of its 4214 notes with the correct
  deck name.
