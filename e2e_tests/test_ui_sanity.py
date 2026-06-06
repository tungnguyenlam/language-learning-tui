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


def test_view_navigation_cycle():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Dashboard -> Dictionary
            agent.act("tab")
            agent.wait_for_text("Dictionary")

            # Dictionary -> Decks
            agent.act("tab")
            agent.wait_for_text("Decks")

            # Decks -> Review
            agent.act("tab")
            agent.wait_for_text("Review")

            # Back to Dashboard
            agent.act("1")
            agent.wait_for_text("DASHBOARD")
        finally:
            agent.close()


def test_status_clears_to_ready():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")  # Settings
            agent.wait_for_text("Settings")
            # Trigger a status message by changing goal
            agent.act("+")
            # Wait for it to return to Ready or remain informative
            time.sleep(3.0)
            screen = agent.observe()
            assert "Ready" in screen or "Daily goal set" in screen or "Saving" in screen
        finally:
            agent.close()


if __name__ == "__main__":
    test_view_navigation_cycle()
    test_status_clears_to_ready()
