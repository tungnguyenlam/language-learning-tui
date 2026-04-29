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

if [ -d "tui_tester/venv" ]; then
	printf 'Running E2E tests...\n'
	# Use a subshell to avoid affecting the current shell's environment
	(
		. tui_tester/venv/bin/activate
		python3 -m pytest e2e_tests/test_tui.py -q
	)
fi
