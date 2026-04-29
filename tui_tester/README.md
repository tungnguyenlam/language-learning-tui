# TUI Tester

A standalone, headless TUI-testing utility designed specifically for agent-based workflows.

## Features
- **Headless PTY Management:** Uses `pexpect` for cross-platform pseudo-terminal manipulation.
- **ANSI Screen Buffer:** Uses `pyte` to parse escape sequences and maintain a realistic virtual screen.
- **Smart Synchronization:** Wait for text, regex patterns, or screen stability before acting.
- **Agent Loop Ready:** High-level `TUIAgent` wrapper makes writing LLM agents straightforward.
- **Rich Diagnostics:** Failed assertions automatically attach visual screen dumps.
- **CLI Daemon:** Start applications in the background and interact with them interactively from the command line over a persistent session.

## Installation

You can install dependencies locally:

```bash
pip install .
```

To run tests:
```bash
pip install ".[test]"
pytest
```

## Persistent CLI Agent Mode

You can run `tui-tester` from the command line to control a persistent TUI session over multiple shell commands. This is highly recommended for LLM Agents doing iterative testing.

```bash
# Start the app in a background daemon
tui-tester start "less /etc/hosts"

# Observe the screen (automatically strips trailing blank lines)
tui-tester observe

# Act by sending keystrokes or control sequences
tui-tester act "<Down>"
tui-tester act "<C-c>"

# Stop the daemon
tui-tester stop
```

## Python Agent Loop Example

```python
from tui_tester import TUIAgent

# Initialize the agent wrapper around your target application
tui = TUIAgent('less /etc/hosts')

while not tui.done:
    # 1. Observe
    state = tui.observe()
    
    # 2. Reason (Your LLM logic goes here)
    if "localhost" in state:
        action = "q" # Quit less
    else:
        action = "<Down>" # Scroll down
        
    # 3. Act
    tui.act(action)
    
    # 4. Synchronize
    tui.wait_until_stable()

tui.close()
```

## Extracting to a Separate Repository

This utility is completely self-contained. To extract it:

1. Copy the `tui_tester/` folder to a new location, OR
2. Use `git subtree split`:
   ```bash
   git subtree split --prefix=tui_tester -b tui-tester-branch
   ```
3. Inside the folder, you can run `git init` and push to a new GitHub repository.
