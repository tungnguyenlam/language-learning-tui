#!/usr/bin/env sh
set -eu

export GOCACHE="${GOCACHE:-/tmp/deutsch-tui-gocache}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

go run ./cmd/deutsch-tui --data-dir "$tmpdir" --smoke
