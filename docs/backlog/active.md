# Active Backlog

Last updated: 2026-07-26

## Current Milestone

Anki Ecosystem Compatibility (follows Practice Trainer Quality, both landed on
`refactor/generic-trainer`).

## Exact Next Action

Continue the screen-interface migration started on `refactor/generic-trainer` —
`internal/tui/screen.go` documents the pattern and only ViewImport /
ViewSessionSummary / ViewAnkiWeb are migrated so far. Alternatively pick the
Backup workflow milestone from `docs/backlog/roadmap.md`.

Possible follow-ups to the AnkiWeb browser, none blocking:

- Media import: downloaded packages' audio blobs are currently dropped;
  `[sound:...]` references are lifted into the note's audio field but the files
  themselves are not stored.
- An E2E scenario for ViewAnkiWeb. It must not touch the network, which means
  the binary needs a way to inject a stub client — the Go tests in
  `internal/tui/screen_ankiweb_test.go` cover the behaviour for now.

## Top Issues

None active.

## Completed Work

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
