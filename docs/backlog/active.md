# Active Backlog

Last updated: 2026-07-27

## Current Milestone

Grammar Trainer Content Expansion & Passive Voice Trainer.

## Exact Next Action

Explore further UI polish, vocabulary deck expansion, or AI tutor enhancements.

## Top Issues

None active.

## Completed Work

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
