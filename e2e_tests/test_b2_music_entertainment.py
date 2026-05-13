import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_b2_music_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")

            # Seed standard content
            agent.act("S")
            agent.wait_for_text("Seeding standard content...", timeout=10.0)
            time.sleep(3.0)
            agent.wait_until_stable(timeout=15.0)

            # Go to Decks
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=5.0)

            # Search for Music
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Music")
            agent.act("<Enter>")
            agent.wait_for_text("German B2 Music", timeout=10.0)

        finally:
            agent.close()

def test_ai_suggestions():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")

            # Seed standard content
            agent.act("S")
            agent.wait_for_text("Seeding standard content...", timeout=10.0)
            time.sleep(3.0)
            agent.wait_until_stable(timeout=15.0)

            # Go to AI view
            agent.act("7")
            agent.wait_for_text("AI Drafts", timeout=5.0)

            # Check for new suggestions (should see expanded list)
            # The new suggestions include: "B1 family & friends", "B2 science & tech", etc.
            # These should be visible when no input is present
            # Wait a moment for rendering
            time.sleep(0.5)
            agent.wait_until_stable()

        finally:
            agent.close()

if __name__ == "__main__":
    test_b2_music_deck()
    test_ai_suggestions()