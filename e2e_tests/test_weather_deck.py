import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_weather_deck_exists():
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
            agent.wait_for_text("DECK LIST", timeout=5.0)
            
            # Search for the unique TSV deck, not "Weather" (that also matches
            # the Go deck "A1 Time & Weather").
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Seasons")
            agent.wait_until_stable()
            agent.wait_for_text("Weather & Seasons", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_until_stable()
            
            # Select the highlighted deck
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("A1 Weather & Seasons", timeout=5.0)
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_weather_deck_exists()