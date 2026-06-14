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

def test_gender_trainer_navigation():
    """Test navigating to and interacting with Gender Trainer."""
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
            
            # Select Gender Trainer (key 1)
            agent.act('1')
            agent.wait_for_text("GENDER TRAINER", timeout=5.0)
            
            # Should show a noun and options
            screen = agent.observe()
            assert "Score: 0/0" in screen
            assert "Which article?" in screen
            assert "der" in screen
            assert "die" in screen
            assert "das" in screen

            # Pick an answer (1 for der)
            agent.act('1')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT" in screen or "INCORRECT" in screen
            assert "Score: 0/1" in screen or "Score: 1/1" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should be back to "Which article?" for next noun
            assert "Which article?" in agent.observe()
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
