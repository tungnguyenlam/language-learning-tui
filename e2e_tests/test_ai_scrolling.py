import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_ai_draft_generation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6") # AI view
            agent.wait_for_text("AI Drafts")
            
            # Clear default input with Ctrl-U and Backspaces
            agent.act("<Ctrl-u>")
            for _ in range(20):
                agent.act("<Backspace>")
            
            topic = "Wunderbar"
            for char in topic:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            
            # Should generate draft
            agent.wait_for_text(f"> {topic}", timeout=10.0)
            agent.assert_text("Approve")
        finally:
            agent.close()
