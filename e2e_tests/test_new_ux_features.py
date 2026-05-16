import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=30)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_cram_mode_filter_types():
    """Test Cram Mode with different filter types: suspended, leech, flagged, all."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # First, create a suspended card
            agent.act("3")  # Go to Review
            agent.wait_for_text("Review 1/52")
            
            # Suspend a card
            agent.act("x")
            agent.wait_for_text("Card suspended")
            
            # Go to Cram Mode
            time.sleep(0.5)
            agent.act("9")
            agent.wait_for_text("Cram Mode")

            agent.wait_for_text("Filter: bookmarked")
            
            # Test suspended filter (press 2)
            agent.act("2")
            agent.wait_for_text("Filter: suspended")
            # Should show suspended card or "No cards found" message
            
            # Test leech filter (press 3)
            agent.act("3")
            agent.wait_for_text("Filter: leech")
            # Should show leech cards or "No cards found" message
            
            # Test flagged filter (press 4)
            agent.act("4")
            agent.wait_for_text("Filter: flagged")
            # Should show flagged cards or "No cards found" message
            
            # Test all filter (press 5)
            agent.act("5")
            agent.wait_for_text("Filter: all")
            agent.wait_for_text("cards loaded")  # Should show card count
        finally:
            agent.close()


def test_browser_deck_switching():
    """Test Browser view deck switching functionality using [ and ] keys."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser view
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Search: _")
            
            # Initially should show cards from current deck
            agent.wait_for_text("blau")
            
            # Try switching to next deck with ]
            agent.act("]")
            agent.wait_until_stable()
            
            # Try switching to previous deck with [
            agent.act("[")
            agent.wait_until_stable()
            
            # Should still show the same cards
            agent.wait_for_text("blau")
            
        finally:
            agent.close()


def test_settings_daily_goal_adjustment():
    """Test Settings view daily goal adjustment with +/- keys."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=110, lines=50)
        try:
            # Go to Settings view
            agent.act("7")
            agent.wait_for_text("Settings")
            # Scroll down to see Daily Goal (it's below AI templates)
            for _ in range(6):
                agent.act("j")
            agent.wait_for_text("Daily Goal:", timeout=5.0)
            
            # Find current daily goal (should be 10 by default)
            import time
            time.sleep(0.5)  # Give time for UI to stabilize
            
            # Increase daily goal with +
            agent.act("+")
            time.sleep(1.0)
            agent.wait_for_text("Daily Goal: 11 cards", timeout=5.0)

            # Decrease daily goal with -
            agent.act("-")
            time.sleep(1.0)
            agent.wait_for_text("Daily Goal: 10 cards", timeout=5.0)

            # Try to decrease below 1 (should stay at 1)
            for _ in range(15):
                agent.act("-")
            time.sleep(1.0)
            agent.wait_for_text("Daily Goal: 1 cards", timeout=5.0)
        finally:
            agent.close()


if __name__ == "__main__":
    import pytest
    pytest.main(["-v", __file__])