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

def test_conjugation_trainer_navigation():
    """Test navigating to and interacting with Verb Conjugation Trainer."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Conjugation view (key K)
            agent.act('K')
            agent.wait_for_text("VERB CONJUGATION TRAINER", timeout=5.0)
            
            # Should show a verb and prompt
            screen = agent.observe()
            assert "Score: 0/0" in screen
            assert "Conjugate for:" in screen
            
            # Type something and press enter
            agent.act('l')
            agent.act('e')
            agent.act('r')
            agent.act('n')
            agent.act('e')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT" in screen or "INCORRECT" in screen
            assert "Score: 0/1" in screen or "Score: 1/1" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should be back to "Conjugate for:" for next verb
            assert "Conjugate for:" in agent.observe()
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
