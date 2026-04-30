# TUI Status Line Stability

Status: active
Scope: internal/tui, e2e_tests
Related: internal/tui/model.go, e2e_tests/test_status_and_browser_regressions.py

## Why It Matters

The tui-tester screen model asserts visible text after terminal wrapping. Long dynamic status text mixed into the navigation footer can split phrases such as `5 cards due` or `Draft approved` across lines, causing false E2E failures even when the app behavior is correct.

## Required Behavior

Keep dynamic status text on a stable single-line surface prefixed with `status:`. Normalize embedded whitespace before rendering and truncate the line to terminal width instead of allowing Lip Gloss wrapping to split important assertions. Prefer E2E assertions against stable status text such as `status: 5 cards due`.

## Revisit When

If the TUI test harness gains robust wrapped-text matching or the footer/status layout is redesigned with equivalent single-line guarantees.
