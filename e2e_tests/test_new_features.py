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
