import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=110, lines=30):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_review_history_empty_state_toggles_in_review():
    """Review view shows and hides empty per-card history."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")

            agent.act("r")
            agent.wait_for_text("Review History: der Apfel")
            agent.assert_text("No reviews yet.")

            agent.act("r")
            agent.wait_for_text("Review history hidden")
            agent.assert_not_text("Review History: der Apfel")
        finally:
            agent.close()


def test_browser_history_shows_reviewed_card_attempt():
    """Browser Enter surfaces persisted history for the selected card."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("status: 5 cards due")

            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("der Apfel")
            agent.act("<Enter>")
            agent.wait_for_text("Review History: der Apfel")
            agent.assert_text("good")
            agent.assert_text("reviews 1")
        finally:
            agent.close()


def test_review_history_persists_after_restart():
    """Review history is loaded from SQLite across app restarts."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("h")
            agent.wait_for_text("status: 5 cards due")
            agent.act("q")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("8")
            restarted.wait_for_text("Card Browser")
            restarted.wait_for_text("der Apfel")
            restarted.act("<Enter>")
            restarted.wait_for_text("Review History: der Apfel")
            restarted.assert_text("hard")
            restarted.assert_text("reviews 1")
        finally:
            restarted.close()
