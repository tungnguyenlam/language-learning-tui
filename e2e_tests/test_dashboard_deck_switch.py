import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ../cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dashboard_deck_switch():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # We are in DASHBOARD
            print("Initial screen:")
            print(agent.observe())
            
            # Press ] to switch deck
            agent.act("]")
            agent.wait_until_stable()
            
            print("Screen after ]:")
            print(agent.observe())
        finally:
            agent.close()

if __name__ == "__main__":
    test_dashboard_deck_switch()
