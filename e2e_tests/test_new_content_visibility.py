import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_new_decks_visibility():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")
            
            # Seed standard content
            agent.act("S")
            agent.wait_for_text("Seeding standard content...")
            time.sleep(2.0)
            agent.wait_until_stable()

            # Go to Decks
            agent.act("2")
            agent.wait_for_text("Decks", timeout=5.0)
            
            # Check for new decks
            agent.wait_for_text("Science & Technology", timeout=5.0)
            agent.wait_for_text("German Proverbs & Idioms", timeout=5.0)
            agent.wait_for_text("German Comprehensive", timeout=5.0)
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_new_decks_visibility()
