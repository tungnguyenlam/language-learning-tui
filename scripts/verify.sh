#!/usr/bin/env sh
set -eu

export GOCACHE="${GOCACHE:-/tmp/deutsch-tui-gocache}"

unformatted="$(gofmt -l cmd internal)"
if [ -n "$unformatted" ]; then
	printf 'gofmt required:\n%s\n' "$unformatted" >&2
	exit 1
fi

go test ./...
go vet ./...
./scripts/tui_smoke.sh
