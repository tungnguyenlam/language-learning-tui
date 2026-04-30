# Browser Deck Reload Contract

Status: active
Scope: internal/tui Browser view
Related: internal/tui/model.go, internal/tui/model_test.go, e2e_tests/test_status_and_browser_regressions.py

## Why It Matters

The Browser view advertises deck switching with `[` and `]`. Changing only the global selected deck leaves the Browser card list stale, so users can believe they are browsing one deck while seeing cards from another.

## Required Behavior

When Browser is active, `[` and `]` must update the selected deck, clear Browser search state, reset the Browser cursor, set `browserDeckID` to the selected deck, and reload Browser cards. Keep unit and tui-tester coverage for this behavior when changing Browser navigation.

## Revisit When

If Browser becomes an all-deck view by design or gets an explicit deck-filter control that replaces global `[` and `]` behavior.
