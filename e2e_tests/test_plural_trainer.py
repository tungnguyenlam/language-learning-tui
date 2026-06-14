import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_plural_trainer_navigation():
    """Test navigating to and interacting with Noun Plural Trainer."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Practice Hub (key 0)
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)
            
            # Select Plural Trainer (key 6)
            agent.act('6')
            agent.wait_for_text("NOUN PLURAL TRAINER", timeout=5.0)
            
            # Should show singular noun
            screen = agent.observe()
            assert "Score: 0/0" in screen
            
            # Type something and press enter
            agent.act('d')
            agent.act('i')
            agent.act('e')
            agent.act(' ')
            agent.act('b')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT" in screen or "INCORRECT" in screen
            assert "Score: 0/1" in screen or "Score: 1/1" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should be back to active question
            assert "Score: 0/1" in agent.observe() or "Score: 1/1" in agent.observe()
            
            # Test Esc to go back to Hub
            agent.act('<Esc>')
            agent.wait_for_text("PRACTICE HUB")
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
