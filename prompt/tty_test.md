# Autonomous TTY Exploratory Testing

You are an expert QA engineer specializing in TUIs. Your mission is to autonomously explore, stress-test, and find bugs in `deutsch-tui` using the `tui-tester` tool.

## The Mission

"See" the app, interact with every corner of it, and try to break it. Look for:
- **Crashes/Panics:** Unexpected exits or hangs.
- **UI Glitches:** Overlapping text, broken borders, misaligned elements, or flickering.
- **Logic Bugs:** Incorrect state transitions, broken keybindings, or data inconsistencies.
- **UX Papercuts:** Confusing navigation, missing help text, or sluggish feedback.

## The Workflow: Observe -> Reason -> Act -> Synchronize

Repeat this loop until you have explored all features or identified critical issues.

1.  **Launch:** Start the app with a unique data directory to avoid state pollution.
    ```bash
    tui-tester start "./deutsch-tui --data-dir /tmp/test-qa-$(date +%s)"
    ```
2.  **Observe:** Capture the screen state.
    ```bash
    tui-tester observe
    ```
3.  **Reason:** Analyze the screen. What is visible? Where can you go? What should happen if you press a certain key?
4.  **Act:** Interact with the UI (keys, mouse).
    ```bash
    tui-tester act "<Key>" # e.g., <Tab>, <Enter>, <Esc>, "j", "k", "/"
    # OR
    tui-tester click X Y
    ```
5.  **Synchronize:** Wait for the UI to stabilize or for a specific anchor.
    ```bash
    tui-tester wait-stable
    # OR
    tui-tester wait-for "DASHBOARD"
    ```

## Testing Strategies

- **Breadth-First:** Visit every screen (Dashboard, Review, Decks, Settings, Dictionary).
- **Stress Test:** Spam keys quickly. Toggle views back and forth.
- **Edge Cases:** Search for empty strings, very long strings, or special characters.
- **Mouse vs. Keyboard:** Ensure both input methods work reliably.
- **Resizing:** (If supported) Simulate different terminal sizes by restarting with different `lines` and `columns`.

## Reporting

For every bug or issue found:
1.  **Capture the evidence:** Save the `tui-tester observe` output that shows the failure.
2.  **Reproduction path:** Document the exact sequence of `act` commands to reproduce it.
3.  **Severity:** Note if it's a blocker, a major glitch, or a minor polish item.

## Cleanup

Always stop the tester when finished.
```bash
tui-tester stop
```

## Next Step: Feed Results into `improve.md`

When testing is complete and the results file is written, the workflow continues with `@prompt/improve.md`. The test results serve as the input bug list for improvement work:

1. Write all findings to `prompt/tty_test_results.md` with clear bug IDs, severity, reproduction steps, and evidence.
2. **Then run `@prompt/improve.md`**, which will read `prompt/tty_test_results.md` and use it as a prioritized fix list alongside its own improvement exploration.

The two prompts form a pipeline: **`tty_test.md` discovers bugs → `improve.md` fixes them.**
