import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=110, lines=30):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=30)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_bookmark_filter_toggle():
    """Toggle bookmark filter to review only bookmarked cards."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.assert_text("Bookmark: off")

            agent.act("b")
            agent.wait_for_text("Card bookmarked")
            agent.assert_text("Bookmark: on")

            agent.act("B")
            agent.wait_for_text("bookmarked cards due")
            agent.assert_text("Bookmarked")
        finally:
            agent.close()


def test_leech_detection_in_statistics():
    """Grade a card 'Again' 3 times and verify leech count in Statistics."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")

            for _ in range(3):
                agent.act("<Space>")
                agent.wait_for_text("cards due")
                agent.act("a")
                agent.wait_for_text("cards due")

            agent.act("q")
        finally:
            agent.close()

        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Leech:")
        finally:
            agent.close()


def test_mcq_review_shows_choices():
    """MCQ cards display numbered choices and accept selection."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")

            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("cards due")

            agent.act("<Space>")
            agent.wait_for_text("1:")

            agent.act("1")
            agent.wait_for_text("Correct")
        finally:
            agent.close()


def test_suspend_card_persists_and_updates_statistics():
    """Suspended cards leave the review queue and remain suspended after restart."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Due cards: 6")
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("x")
            agent.wait_for_text("Card suspended")
            agent.act("q")
        finally:
            agent.close()

        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Due cards: 5")
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Suspended:")
            agent.assert_text("1")
        finally:
            agent.close()


def test_daily_goal_setting_persists_after_restart():
    """Settings +/- updates the SQLite-backed daily review goal."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.assert_text("Daily Goal: 10")
            agent.act("+")
            agent.wait_for_text("Daily goal set to 11")
            agent.act("q")
        finally:
            agent.close()

        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Reviews today: 0/11")
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.assert_text("Daily Goal: 11")
        finally:
            agent.close()


def test_decks_view_shows_progress_metrics_after_review():
    """Decks view surfaces reviews-today and success-rate metrics."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("cards due")

            agent.act("2")
            agent.wait_for_text("Decks")
            agent.assert_text("today 1")
            agent.assert_text("100% success")
        finally:
            agent.close()
