#!/usr/bin/env sh
set -eu

export GOCACHE="${GOCACHE:-/tmp/deutsch-tui-gocache}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

DEUTSCH_TUI_BIN="${DEUTSCH_TUI_BIN:-go run ./cmd/deutsch-tui}"
$DEUTSCH_TUI_BIN --data-dir "$tmpdir" --smoke
