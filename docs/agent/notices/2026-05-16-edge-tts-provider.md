# Edge TTS Provider

Status: active
Scope: internal/audio, internal/tui audio playback, app config
Related: `internal/audio/edge_tts.go`, `internal/tui/model.go`, `internal/app/config.go`, `github.com/lib-x/edgetts`

## Why It Matters

The Edge TTS integration uses the third-party Go library `github.com/lib-x/edgetts`, which calls Microsoft Edge's online neural TTS service without an Azure key. It is high quality but unofficial and can fail if the network is unavailable or Microsoft changes the endpoint.

## Required Behavior

Keep Edge TTS optional and cached. Provider errors must produce a clear status/error and must not break existing card audio playback. Do not reintroduce a Python/CLI dependency or hardcode Microsoft/Azure credentials for this path.

## Revisit When

Revisit if the app gains an official Azure Speech provider, a bundled offline TTS engine, or Edge TTS changes its API behavior.
