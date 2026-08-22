import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

from e2e_helpers import read_due_count


def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_review_bookmark_persists_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.assert_text("Bookmark: off")

            agent.act("b")
            agent.wait_for_text("Card bookmarked")
            agent.assert_text("Bookmark: on")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("3")
            restarted.wait_for_text(f"Review 1/{due}")
            restarted.assert_text("Bookmark: on")

            restarted.act("4")
            restarted.wait_for_text("Statistics")
            restarted.assert_text("Bookmarked: 1")
        finally:
            restarted.close()

def test_undo_last_review_restores_due_card_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text(f"{due - 1} cards due")

            agent.act("u")
            agent.wait_for_text("Last review undone")
            agent.assert_text(f"Review 1/{due}")
            agent.assert_text("blau")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text(f"Due cards:   {due}")
            restarted.act("3")
            restarted.wait_for_text(f"Review 1/{due}")
        finally:
            restarted.close()


def test_statistics_daily_progress_updates_after_review():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Reviews Today: 0/10")
            agent.wait_for_text("Statistics")

            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text(f"{due - 1} cards due")

            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Total Reviews: 1")
            agent.assert_text("Reviews Today: 1/10")
            agent.wait_for_text("Statistics")
        finally:
            agent.close()
