import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_practice_hub_mouse_click():
    """Test entering a trainer via mouse click on the Practice Hub."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Decks and select one to ensure we have nouns
            agent.act('2')
            agent.wait_for_text("Decks")
            agent.act('<Enter>')
            agent.wait_for_text("DASHBOARD")

            # Go to Practice Hub (key 0)
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)
            
            # Click the Gender Trainer button (middle column, first row)
            # Center of the button is around X=55, Y=7
            agent.click(55, 7)
            agent.wait_for_text("GENDER TRAINER", timeout=5.0)
            
            # Click one of the gender options (e.g. "die" which is around X=65, Y=15)
            agent.click(65, 15)
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT" in screen or "INCORRECT" in screen
            
            # Click inside content area to continue
            agent.click(50, 10)
            agent.wait_until_stable()
            
            # Should reset to unrevealed state (Which article?)
            assert "Which article?" in agent.observe()
            
        finally:
            agent.close()
