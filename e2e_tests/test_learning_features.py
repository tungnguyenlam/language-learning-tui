import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=110, lines=30):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_review_bookmark_persists_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.assert_text("Bookmark: off")

            agent.act("b")
            agent.wait_for_text("Card bookmarked")
            agent.assert_text("Bookmark: on")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("3")
            restarted.wait_for_text("Review 1/6")
            restarted.assert_text("Bookmark: on")

            restarted.act("4")
            restarted.wait_for_text("Statistics")
            restarted.assert_text("Bookmarked:  1")
        finally:
            restarted.close()


def test_undo_last_review_restores_due_card_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("5 cards due")

            agent.act("u")
            agent.wait_for_text("Last review undone")
            agent.assert_text("Review 1/6")
            agent.assert_text("der Apfel")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text("Due cards:   6")
            restarted.act("3")
            restarted.wait_for_text("Review 1/6")
        finally:
            restarted.close()


def test_statistics_daily_progress_updates_after_review():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Reviews Today: 0/10")
            agent.assert_text("Current Streak: 0 days")

            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("5 cards due")

            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Total Reviews: 1")
            agent.assert_text("Reviews Today: 1/10")
            agent.assert_text("Current Streak: 1 days")
        finally:
            agent.close()
