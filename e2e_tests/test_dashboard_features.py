import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_help_overlay_toggle():
    """Test that pressing '?' toggles the help overlay."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Press '?' to toggle help
            agent.act('?')
            agent.wait_for_text("Help overlay shown. Press ? to close.")
            
            # Press '?' again to hide help
            agent.act('?')
            agent.wait_for_text("Help overlay closed.")
        finally:
            agent.close()

def test_dashboard_deck_switching_updates_view():
    """Test that '[' and ']' switch decks on the Dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Check initial deck
            agent.assert_text("Deck: German A1 Survival")
            
            # Switch deck right (only one deck exists initially, so it should stay the same)
            agent.act(']')
            agent.wait_for_text("Deck: German A1 Survival")
            
            # Switch deck left
            agent.act('[')
            agent.wait_for_text("Deck: German A1 Survival")
        finally:
            agent.close()

def test_streak_persists_across_sessions():
    """Test that current streak is saved and reloaded across sessions."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # First session
        agent1 = start_agent(tmpdir)
        try:
            agent1.assert_text("Current Streak: 0 days")

            # Do reviews to build streak
            agent1.act('3')
            agent1.wait_for_text("Review 1/")
            for _ in range(5):
                agent1.act('<Space>')
                agent1.wait_until_stable()
                agent1.act('g')
                agent1.wait_until_stable()
                try:
                    agent1.wait_for_text("Review 1/", timeout=1)
                except:
                    break
            
            # Verify streak updated in session 1
            agent1.act('1')
            agent1.wait_for_text("Dashboard")
            agent1.assert_text("Current Streak: 1 days 🔥")
        finally:
            agent1.close()

        # Second session using same directory
        agent2 = start_agent(tmpdir)
        try:
            # Streak should still be 1 day
            agent2.assert_text("Current Streak: 1 days 🔥")
        finally:
            agent2.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
