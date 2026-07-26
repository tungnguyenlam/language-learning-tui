# internal/ankiweb

Read-only client for AnkiWeb's public shared-deck library. It backs the
`ViewAnkiWeb` browser (`internal/tui/screen_ankiweb.go`), which is the app's
only network-facing screen.

## Responsibilities

- **Search**: `Client.Search` queries `/svc/shared/list-decks?search=`.
- **Details**: `Client.Info` fetches `/svc/shared/item-info?sharedId=`, which
  also carries the short-lived token required to download.
- **Download**: `Client.Download` streams `/svc/shared/download-deck/<id>?t=<key>`
  to an `io.Writer`, reporting progress and capping at `MaxDownloadBytes`.

## AnkiWeb Publishes No API

These endpoints are the ones AnkiWeb's own SvelteKit client calls. They are
protobuf with no published schema, so the field numbers in `ankiweb.go` were
read off live responses. Consequences that must be preserved:

- **Decode defensively.** `protobuf.go` is a hand-written wire reader that
  returns partial results on truncation. An added, reordered, or removed field
  must yield an incomplete deck (dropped if unusable) — never an error and
  never a panic. See `TestSearchToleratesUnknownFields`.
- **Every call has a timeout and a size cap.** A redirect to something
  unexpected must not stream unbounded into memory or onto disk.
- **Download tokens expire.** Refresh via `Info` immediately before
  downloading; a stale token yields `ErrLinkExpired`.

## Error Sentinels

| Sentinel | Meaning | UI response |
|---|---|---|
| `ErrAnonymousLimit` | AnkiWeb allows a few anonymous searches/downloads per address, then asks for an account | Normal outcome, not a fault — point at `ankiweb.net` + manual `.apkg` import |
| `ErrLinkExpired` | Download token went stale | Reload the deck info, then retry |
| `ErrNotAvailable` | Unreachable, or answered with something unreadable | Treat the browser as unavailable; the rest of the app is unaffected |

`classifyError` maps AnkiWeb's plain-text error bodies onto these. The wording
it matches on ("log in", "expired", "link is invalid") is AnkiWeb's, not ours,
so it may drift — the fallback is `ErrNotAvailable` with the server's own text.

## Privacy

No credentials are involved; the shared library is public. Nothing about the
user is sent, and the network is touched only on an explicit search or
download. `BaseURL` is a variable so tests point at `httptest` — the test suite
never leaves the machine.

## Key Symbols

- `Client`, `New`, `NewWithHTTPClient`
- `Deck`, `Details` (its `downloadKey` is deliberately unexported)
- `DeckURL`: the human page, for the manual fallback
- `pbFields`, `pbString`, `pbUint`, `pbSubMessage`: the wire reader
