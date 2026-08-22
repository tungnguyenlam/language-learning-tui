import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

from e2e_helpers import read_due_count


def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_review_history_empty_state_toggles_in_review():
    """Review view shows and hides empty per-card history."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")

            agent.act("r")
            agent.wait_for_text("Review History: blau")
            agent.assert_text("No reviews yet.")

            agent.act("r")
            agent.wait_for_text("Review history hidden")
            agent.wait_until_stable()
            # Wait for the history panel to actually disappear
            start = time.time()
            while time.time() - start < 5.0:
                if "Review History: blau" not in agent.observe():
                    break
                time.sleep(0.1)

        finally:
            agent.close()


def test_browser_history_shows_reviewed_card_attempt():
    """Browser Enter surfaces persisted history for the selected card."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("cards due")

            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("blau")
            agent.act("<Enter>")
            agent.wait_for_text("Review History: blau")
            agent.assert_text("good")
            agent.assert_text("reviews 1")
        finally:
            agent.close()


def test_review_history_persists_after_restart():
    """Review history is loaded from SQLite across app restarts."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("h")
            agent.wait_for_text("cards due")
            agent.act("q")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("8")
            restarted.wait_for_text("Card Browser")
            restarted.wait_for_text("blau")
            restarted.act("<Enter>")
            restarted.wait_for_text("Review History: blau")
            restarted.assert_text("hard")
            restarted.assert_text("reviews 1")
        finally:
            restarted.close()
