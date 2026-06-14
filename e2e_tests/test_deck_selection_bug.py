import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ../cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_deck_selection():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Decks view (key "2")
            agent.act("2")
            agent.wait_for_text("German A1 Survival")
            agent.wait_until_stable()
            
            # Press down to select another deck
            agent.act("j")
            agent.wait_until_stable()
            
            # Press Enter
            agent.act("<Enter>")
            agent.wait_until_stable()
            
            # Should be back at DASHBOARD
            agent.wait_for_text("DASHBOARD")
            print(agent.observe())
        finally:
            agent.close()

if __name__ == "__main__":
    test_deck_selection()
