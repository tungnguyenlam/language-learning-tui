# Active Backlog

Last updated: 2026-07-27

## Current Milestone

Dictionary Feature Enhancements & Lemmatization Search.

## Exact Next Action

Continue dictionary UX refinements, vocabulary deck expansion, or AI tutor enhancements.

## Top Issues

None active.

## Completed Work

- [x] **Dictionary One-Key Cloze Flashcard Generation (`c` / `ctrl+c`):** Added `addDictionaryClozeEntryCmd()` in `actions_dictionary.go` and `c` / `ctrl+c` key shortcuts in `keys.go`. Extracts best example sentences from dictionary entries and converts target headwords into `{{c1::<word>}}` Cloze flashcards (or formats fallback Cloze prompts) saved directly to the selected target deck. Added `TestDictionaryClozeCardGeneration`.
- [x] **Recently Inspected Words History Stack & Domain Tag Badges:** Added `dictionaryRecentlyViewed` stack in `Model` and `recordDictionaryView()` in `actions_dictionary.go`. Renders interactive `Recently Inspected Words:` section when dictionary search is cleared. Added styled domain/field tag badges (e.g. `[ZOOL.]`, `[COLLOQUIAL]`, `[MEDICINE]`) in detail panel view (`render_dictionary.go`). Added `TestDictionaryRecentlyViewedAndDomainTags`.
- [x] **Dictionary Entry Star/Bookmark System & `:starred` Filter (`b` / `ctrl+b`):** Added `dictionaryStarred map[string]bool` in `Model`, `toggleStarDictionaryEntry()`, `loadDictionaryStarred()`, and `:starred` filter tag support in `actions_dictionary.go` and `sqlite/dictionary.go`. Renders bright yellow `★` stars in search lists and detail view headers. Starred entries are persisted to SQLite settings (`dict_starred_entries`). Added `[★ Starred]` pill to filter row and `TestDictionaryStarringAndFiltering`.
- [x] **Dictionary Result List Translation Snippets in 2-Column Mode:** Updated `renderDictionary` and `renderSpotlightDictionary` (>80/70 cols) in `render_dictionary.go` so left result items show `Word {gender} - Translation` instead of just headwords, allowing instant scanning of search result meanings. Added `TestDictionaryTwoColumnTranslationDisplay`.
- [x] **Dictionary Target Deck Selection (`ctrl+g`) & Batch Add (`ctrl+s`):** Added `dictionaryTargetDeckID` field in `Model`, `cycleDictionaryTargetDeck()` (`ctrl+g`) to cycle through available decks for dictionary card creation, and `addDictionaryEntriesBatchCmd()` (`ctrl+s`) to batch-add all current dictionary search results directly into the selected target deck. Added `TestDictionaryTargetDeckCyclingAndBatchAdd`.

- [x] **Dictionary Input Fix for 'k' and 'j' Keys:** Updated `updateDictionaryKey()` in `internal/tui/keys.go` so that `'k'` and `'j'` are typed as input characters when the search bar is active (`!m.dictionaryFocusResults && !m.dictionaryDetailView`), preventing search history cycling or list navigation when typing words containing 'k' or 'j'. Added `TestDictionaryKAndJKeyHandling` in `render_dictionary_test.go`.

- [x] **Dictionary Wildcard (* & ?) & Language Scope (de: & en:) Search:** Updated `parseSearchFilters`, `filterEntries`, and `Search` in `internal/storage/sqlite/dictionary.go` to convert wildcard queries into SQLite `LIKE` queries (`*` -> `%`, `?` -> `_`), enabling suffix searches like `*ung`, `*keit`, `*ieren` and language scope filters (`de:` for German headwords/forms, `en:` for English translations). Added `TestDictionarySearchWildcardAndLanguagePrefix`.
- [x] **Interactive Dictionary Filter Pills:** Added `renderFilterPillsRow()` and filter tag helpers (`isFilterActive`, `toggleFilterTag`, `clearFilterTags`) in `internal/tui/render_dictionary.go` and `internal/tui/actions_dictionary.go`. Renders interactive clickable filter pills (`[All]`, `[:verb]`, `[:noun]`, `[:adj]`, `[:adv]`, `[:m]`, `[:f]`, `[:n]`, `[:pl]`) in both standard Dictionary view and Spotlight Dictionary overlay.
- [x] **Dictionary Entry Contextual AI Explainer (`ctrl+e`):** Added `ctrl+e` shortcut in `updateDictionaryKey` (`internal/tui/keys.go`) to automatically pre-populate the AI Tutor (`ViewAI`) with a structured request explaining the selected German dictionary entry (grammar notes, gender usage, collocations, and example sentences). Added unit test `TestDictionaryCtrlEExplainAndFilterTagHelpers`.
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
