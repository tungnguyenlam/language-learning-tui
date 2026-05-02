# Parallel E2E Testing

Status: active
Scope: e2e_tests/, scripts/verify.sh
Related: `pytest-xdist`, `verify.sh`

## Why It Matters

E2E tests were taking over 5 minutes to run sequentially. They have been parallelized using `pytest-xdist`, reducing the execution time to under 1 minute.

## Required Behavior

- **Maintain Isolation:** Each test must continue to use a unique `-data-dir` (typically via `tempfile.TemporaryDirectory()`) to avoid SQLite database locks or state contamination between concurrent processes.
- **Avoid Global Resources:** Do not introduce dependencies on fixed global resources (like hardcoded network ports or fixed file paths outside the temporary data directory) that would prevent multiple instances of the app from running simultaneously.
- **Verification:** Always run `./scripts/verify.sh` to ensure that new tests are compatible with parallel execution.

## Revisit When

If tests become flaky due to resource contention or if the overhead of `go run` during parallel startup becomes a bottleneck, consider pre-building the binary once before starting the test suite.
