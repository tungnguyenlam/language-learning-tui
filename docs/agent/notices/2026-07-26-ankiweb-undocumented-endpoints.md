# AnkiWeb's shared-deck endpoints are undocumented and may break

Date: 2026-07-26

## What

`internal/ankiweb` talks to three AnkiWeb endpoints that have no published
schema or stability guarantee:

- `GET /svc/shared/list-decks?search=<query>`
- `GET /svc/shared/item-info?sharedId=<id>`
- `GET /svc/shared/download-deck/<id>?t=<downloadKey>`

They return protobuf. The field numbers in `ankiweb.go` were read off live
responses with a hand-written decoder, cross-validated between the two metadata
endpoints (item-info `10.1`/`10.2` notes/cards match list-decks `6`/`7`;
item-info `18`/`19` thumbs match list-decks `3`/`4`). The download token lives
at item-info field `10.5`.

## Why it matters

AnkiWeb can change these at any time without notice. The app must keep working
when that happens.

## How to apply

- **Never let a decode failure become an app failure.** `protobuf.go` returns
  the fields it managed to read plus an error; callers use the partial result
  and drop entries that are unusable (no ID, no title). A new or reordered
  field must degrade the browser, not break it.
- **Refresh `Info` right before `Download`.** The token is short-lived; a stale
  one returns `ErrLinkExpired`.
- **Rate limiting is not a bug.** Anonymous use is capped for both searches and
  downloads (`"Please log in to perform more searches."` /
  `"...to download more decks."`). `ErrAnonymousLimit` is a normal outcome and
  the UI names the manual fallback: download the `.apkg` from ankiweb.net and
  import it with `I`.
- **Keep the browser out of the tab cycle.** `ViewAnkiWeb` is reachable only
  from Import (`A`), so the offline-first flow is never routed through a
  network view. Asserted by `TestAnkiWebIsReachableFromImportButNotInTabCycle`.
- **Tests must not hit the network.** `BaseURL` is a variable; unit tests use
  `httptest`, and the TUI tests use a stub `ankiWebSearcher`.

## Where

- `internal/ankiweb/ankiweb.go`, `internal/ankiweb/protobuf.go`
- `internal/tui/screen_ankiweb.go`
- `internal/ankiweb/README.md`
