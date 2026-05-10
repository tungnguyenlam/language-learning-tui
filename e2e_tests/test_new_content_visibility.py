import os
import sys
import tempfile
import time

sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)
from tui_tester import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
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
            agent.wait_for_text("DECK LIST", timeout=5.0)

            # Check for new decks using search
            for search_term, deck_name_partial in [
                ("Science", "Science & Technology"),
                ("Proverbs", "German Proverbs"),
                ("Comprehensive", "German Comprehensive"),
            ]:
                agent.act("/")
                agent.wait_for_text("Search:", timeout=2.0)
                agent.act(search_term)
                agent.act("<Enter>")
                agent.wait_for_text(deck_name_partial, timeout=5.0)
                # It should have selected the deck (cursor on it)
                # Now press Enter to go to Dashboard and verify
                agent.act("<Enter>")
                agent.wait_for_text("DASHBOARD", timeout=5.0)
                agent.wait_for_text(deck_name_partial, timeout=5.0)
                # Go back to Decks for next search
                agent.act("2")
                agent.wait_for_text("DECK LIST", timeout=5.0)
                agent.wait_until_stable()
                agent.act("<Esc>")  # Clear filter
                agent.wait_until_stable()
        finally:
            agent.close()


if __name__ == "__main__":
    test_new_decks_visibility()
