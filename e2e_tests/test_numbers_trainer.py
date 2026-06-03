import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ../cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_numbers_trainer():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Practice Hub
            agent.act("0")
            agent.wait_for_text("PRACTICE HUB")
            agent.wait_for_text("Numbers & Time")
            
            # Select Numbers Trainer
            agent.act("8")
            agent.wait_for_text("NUMBER & TIME TRAINER")
            
            # Check for a number or time (should have at least one character of question)
            screen = agent.observe()
            assert ":" in screen or any(str(i) in screen for i in range(10))
            
            # Press enter with empty input (should stay on same screen)
            agent.act("<Enter>")
            time.sleep(0.5)
            assert "CORRECT!" not in agent.observe()
            assert "INCORRECT" not in agent.observe()
            
            # Go back to Hub
            agent.act("<Esc>")
            agent.wait_for_text("PRACTICE HUB")
            
            # Go back to Dashboard
            agent.act("<Esc>")
            agent.wait_for_text("DASHBOARD")
            
        finally:
            agent.close()
