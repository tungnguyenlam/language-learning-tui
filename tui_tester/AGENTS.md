# Agent Guide: TUI Tester

You are an AI agent tasked with testing or interacting with a Terminal User Interface (TUI). This utility allows you to treat a terminal application as an observable, interactive state machine.

## Interaction Strategy

Follow the **Observe -> Reason -> Act -> Synchronize** loop.

### 1. Start the Session
Always start by launching the application. If you are in a persistent environment (like a long-running CLI chat), use the background daemon.

```bash
tui-tester start "your-command-here"
```

### 2. Observe
Capture the current state of the screen. The output is a clean text representation of the terminal buffer.

```bash
tui-tester observe
```

### 3. Reason
Analyze the text output. Look for:
- Menu items or buttons.
- Status messages (success, error, loading).
- Active cursor indicators (like `>`).

### 4. Act
Send keystrokes or mouse events.
- **Keystrokes:** Normal characters `"a"`, special keys `"<Enter>"`, `"<Esc>"`, `"<Up>"`, or control chords `"<C-c>"`.
- **Mouse Click:** `tui-tester click X Y` (1-based terminal coordinates).
- **Mouse Drag:** `tui-tester drag StartX StartY EndX EndY` (1-based).

```bash
tui-tester act "<Down>"
tui-tester click 15 3
```

### 5. Synchronize
**Crucial:** TUI applications are asynchronous. After an action, the screen may take time to update. Always wait for a specific visual anchor or for the screen to stop changing.

```bash
# Wait for a specific message
tui-tester wait-for "Successfully saved"

# Wait for the UI to finish rendering/animating
tui-tester wait-stable
```

## Best Practices

1.  **Unique Data Directories:** When testing stateful apps (like SQLite-backed TUIs), always pass a unique temporary directory to the application command (e.g., `-data-dir /tmp/test-1`).
2.  **Explicit Timeouts:** If an application is slow (e.g., AI generation), provide a higher timeout to `wait-for`.
3.  **Clean Up:** Always stop the daemon when finished to release the PTY and system resources.
    ```bash
    tui-tester stop
    ```
4.  **Language Agnostic:** This tool doesn't care if the app is written in Go, Rust, or Python. If it runs in a terminal, you can test it.
