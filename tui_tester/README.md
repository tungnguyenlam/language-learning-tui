# TUI Tester

A standalone, headless TUI-testing utility designed specifically for agent-based workflows.

## Features
- **Language Agnostic:** Test TUIs written in Go, Rust, Python, C++, etc. Treatments the app as a black box.
- **Full Mouse Support:** Simulate clicks, motion, and **drag-to-scroll** interactions using the xterm SGR protocol.
- **Headless PTY Management:** Uses `pexpect` for cross-platform pseudo-terminal manipulation.
- **ANSI Screen Buffer:** Uses `pyte` to parse escape sequences and maintain a realistic virtual screen.
- **Smart Synchronization:** Wait for text, regex patterns, or screen stability before acting.
- **Agent Loop Ready:** High-level `TUIAgent` wrapper and `AGENTS.md` guide make LLM integration straightforward.
- **CLI Daemon:** Control persistent TUI sessions over multiple shell commands.

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

## Publishing as a Skill or MCP Server

### As a Gemini CLI Skill
You can package this utility as a "Skill" to give Gemini CLI native TUI-testing capabilities.
1.  Add a `SKILL.md` file to the root of the repository.
2.  Define the tools (start, observe, act, wait) using the format expected by the `skill-creator`.
3.  Users can then activate it with `activate_skill("tui-tester")`.

### As an MCP Server (Model Context Protocol)
To use this with other LLM clients (like Claude Desktop), you can wrap the CLI in an MCP server.
1.  Use the `python-mcp-sdk`.
2.  Expose the `TUIAgent` methods as MCP Tools.
3.  This allows any MCP-compatible agent to "see" and "touch" terminal applications.
