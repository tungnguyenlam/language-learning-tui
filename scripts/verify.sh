#!/usr/bin/env sh
set -eu

export GOCACHE="${GOCACHE:-/tmp/deutsch-tui-gocache}"
export PIP_DISABLE_PIP_VERSION_CHECK="${PIP_DISABLE_PIP_VERSION_CHECK:-1}"

unformatted="$(gofmt -l cmd internal)"
if [ -n "$unformatted" ]; then
	printf 'gofmt required:\n%s\n' "$unformatted" >&2
	exit 1
fi

go test ./...
go vet ./...
./scripts/tui_smoke.sh

# Build binary once for E2E tests to reduce resource contention
printf 'Building binary for E2E tests...\n'
go build -o deutsch-tui-test ./cmd/deutsch-tui
export DEUTSCH_TUI_BIN="$(pwd)/deutsch-tui-test"

if [ -d "tui_tester/venv" ]; then
	printf 'Running E2E tests in parallel...\n'
	# Use a subshell to avoid affecting the current shell's environment
	(
		. tui_tester/venv/bin/activate
		# Install pytest-xdist if missing (for robustness in different environments)
		python3 -m pip install pytest-xdist --quiet
		python3 -m pytest e2e_tests/ -n auto -q
	)
	# Cleanup test binary
	rm -f ./deutsch-tui-test
fi
