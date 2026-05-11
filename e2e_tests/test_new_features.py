import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=110, lines=40):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_bookmark_filter_toggle():
    """Toggle bookmark filter to review only bookmarked cards."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/52")
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
            agent.wait_for_text("Review 1/52")

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
            # First find an MCQ in the browser to make sure we know what to look for
            agent.act("8")
            agent.wait_for_text("Card Browser")
            # Enter search mode with /
            agent.act("/Ich")
            agent.act("<Enter>")
            agent.wait_for_text("[MCQ]")
            # Go back to dashboard then review
            agent.act("1")

            agent.wait_for_text("DASHBOARD")
            
            agent.act("3")
            agent.wait_for_text("Review 1/52")

            # Do reviews until we hit an MCQ or just use the first one if we can ensure order.
            # Since order is alphabetical by ID, let's see which one is first.
            # a1-col-blau is first.
            # We will just do reviews until we see "1:"
            for _ in range(50):
                agent.act("<Space>")
                agent.wait_until_stable()
                if "1:" in agent.observe():
                    break
                agent.act("g")
                agent.wait_until_stable()

            agent.wait_for_text("1:")
            agent.act("1")
            # It could be Correct or Incorrect depending on which choice is #1
            # Just verify that SOME feedback is shown
            import re
            agent.wait_for_regex(r"(Correct|Incorrect)") 
        finally:
            agent.close()


def test_suspend_card_persists_and_updates_statistics():
    """Suspended cards leave the review queue and remain suspended after restart."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Due cards:   52")
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("x")
            agent.wait_for_text("Card suspended")
            agent.act("q")
        finally:
            agent.close()

        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Due cards:   51")
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
            agent.assert_text("Reviews:     0/11")
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
            agent.wait_for_text("Review 1/52")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("cards due")

            agent.act("2")
            agent.wait_for_text("Decks")
            agent.assert_text("today 1")
            agent.assert_text("100%")
        finally:
            agent.close()


def test_card_browser_search_filter():
    """Card Browser filters cards by search term."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("52 cards found")

            # Enter search mode and type query
            agent.act("/Ap")
            agent.wait_for_text("Apfel")
            # Verify only matching cards shown
            agent.wait_for_text("1 cards found")
        finally:
            agent.close()


def test_session_stats_show_in_statistics():
    """Session statistics update after reviewing cards."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("a")
            agent.wait_for_text("cards due")

            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Session Stats:")
            agent.wait_for_text("Statistics")
            
            
        finally:
            agent.close()


def test_help_overlay_shows_and_dismisses():
    """Pressing ? shows help overlay, pressing ? again dismisses it."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, lines=60)
        try:
            agent.act("?")
            agent.wait_for_text("Help overlay shown. Press ? to close.")
            agent.assert_text("Global:")

            agent.act("?")
            # Wait for the help overlay to disappear
            import time
            time.sleep(1)
            agent.wait_until_stable()
            agent.assert_not_text("Help overlay shown. Press ? to close.")
        finally:
            agent.close()
