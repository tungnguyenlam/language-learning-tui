# Skill: TUI Tester

## Description
Expert utility for headless, black-box testing of terminal-based applications (TUIs). It supports full ANSI parsing, smart synchronization, and complex mouse interactions (clicks, motion, dragging).

## Instructions

When this skill is activated, follow the **Observe -> Reason -> Act -> Synchronize** loop for all TUI interactions.

### 1. Initialization
Start the application using the background daemon. Always use a unique data directory for stateful apps to avoid database locks.
```bash
tui-tester start "your-app-command --data-dir /tmp/unique-dir"
```

### 2. Observation
Capture the current screen state. Analyze the text output for visual anchors, menu items, or status changes.
```bash
tui-tester observe
```

### 3. Interaction
Send keystrokes or simulate mouse events.
- **Keys:** Use literal strings for characters or named keys like `"<Enter>"`, `"<Esc>"`, `"<C-c>"`.
- **Mouse:** Use 1-based coordinates for clicks and drags.
```bash
tui-tester act "<Tab>"
tui-tester click 15 3
tui-tester drag 83 6 83 16
```

### 4. Synchronization
**NEVER** skip synchronization. TUI updates are asynchronous.
- Use `wait-for` when expecting specific text (e.g., a success message).
- Use `wait-stable` when waiting for animations to finish or for the initial render.
```bash
tui-tester wait-for "Process Complete" --timeout 10.0
tui-tester wait-stable
```

### 5. Termination
Always stop the daemon when your task is complete to release system resources.
```bash
tui-tester stop
```

## Available Resources

| Command | Usage | Description |
| :--- | :--- | :--- |
| `start` | `tui-tester start "<cmd>"` | Starts the TUI app in a background PTY. |
| `observe` | `tui-tester observe` | Returns the current terminal screen as plain text. |
| `act` | `tui-tester act "<key>"` | Sends a keystroke or control sequence. |
| `click` | `tui-tester click <X> <Y>` | Sends an xterm SGR mouse click (1-based). |
| `drag` | `tui-tester drag <X1> <Y1> <X2> <Y2>` | Simulates a mouse drag from (X1,Y1) to (X2,Y2). |
| `wait-for` | `tui-tester wait-for "<text>"` | Polls the screen until the text appears. |
| `wait-stable` | `tui-tester wait-stable` | Polls the screen until the ANSI buffer stops changing. |
| `stop` | `tui-tester stop` | Terminates the app and stops the background daemon. |
| `AGENTS.md` | `tui_tester/AGENTS.md` | Detailed strategy guide for AI-driven testing. |
