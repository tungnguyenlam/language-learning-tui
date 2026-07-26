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
	# Detect macOS and limit parallelism to avoid PTY exhaustion (BUG-011)
	num_procs="auto"
	if [ "$(uname)" = "Darwin" ]; then
		num_procs=4
	fi

	# Use a subshell to avoid affecting the current shell's environment
	(
		. tui_tester/venv/bin/activate
		# Install pytest-xdist if missing (for robustness in different environments)
		python3 -m pip install pytest-xdist --quiet
		if [ "${FULL_E2E:-0}" = "1" ] || [ "${1:-}" = "--full" ]; then
			printf 'Running FULL E2E test suite (100+ files)...\n'
			python3 -m pytest e2e_tests/ -n "$num_procs" -q
		else
			printf 'Running CORE E2E test suite (pass --full or set FULL_E2E=1 for full 100+ file suite)...\n'
			python3 -m pytest e2e_tests/test_tui.py e2e_tests/test_end_to_end_core_views.py e2e_tests/test_ui_sanity.py e2e_tests/test_interactive_features.py e2e_tests/test_trainer_input_shortcuts.py -n "$num_procs" -q
		fi
	)
	# Cleanup test binary
	rm -f ./deutsch-tui-test
fi
