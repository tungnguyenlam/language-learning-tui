import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from tui_tester.agent import TUIAgent

def start_agent(tmpdir, columns=110, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', './deutsch-tui-bin')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

@pytest.mark.e2e
def test_b2_business_meetings_deck_exists():
    """Verify that the new Business Meetings deck is available."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("S") # Seed
            agent.wait_for_text("notes from Standard Content")
            
            agent.act("2") # Decks
            agent.wait_for_text("DECK LIST")
            
            # Use search to find the new deck
            agent.act("/")
            for char in "business":
                agent.act(char)
            agent.act("<Enter>")
            
            agent.wait_for_text("B2 Business: Meetings & Negotiations")
            agent.assert_text("Professional German for meetings")
            
        finally:
            agent.close()

@pytest.mark.e2e
def test_dashboard_recent_shortcuts():
    """Verify that !, @, # shortcuts work on the dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5") # Import
            agent.act("S") # Seed
            agent.wait_for_text("notes from Standard Content")
            
            # Go to Decks and pick a specific one
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            
            # Select A1 German Animals
            agent.act("/")
            for char in "German Animals":
                agent.act(char)
            agent.act("<Enter>")
            agent.wait_for_text("A1 German Animals")
            agent.act("<Enter>") # Select it
            
            # Go to Review to make it "recently studied"
            agent.act("3")
            agent.wait_for_text("Review")
            agent.act("<Enter>") # Reveal
            agent.wait_for_text("Grade:")
            agent.act("3") # Grade Good
            agent.wait_for_text("✓ Good")
            time.sleep(2.0)
            
            # Go back to Dashboard
            agent.act("1")
            agent.wait_for_text("Recently Studied")
            agent.wait_for_text("! • A1 German Animals")
            
            # Switch to a different deck first
            agent.act("]")
            
            # Now use '!' to jump back to Animals
            agent.act("!")
            agent.wait_for_text("Review")
            agent.assert_text("A1 German Animals")
        finally:
            agent.close()

@pytest.mark.e2e
def test_goal_met_badge():
    """Verify that the GOAL MET badge appears when goal is reached."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Set daily goal to 1 in settings
            agent.act("7") # Settings
            agent.wait_for_text("Daily Goal")
            
            # Navigate to goal adjustment (it's in the header or specific row)
            # Actually, settings-goal-minus is a thing.
            # But let's just use keys: '-' to decrease until 1.
            for _ in range(20):
                agent.act("-")
            
            agent.wait_for_text("Goal: 1")
            
            # Seed and review 1 card
            agent.act("5")
            agent.act("S")
            agent.wait_for_text("notes from Standard Content")
            
            agent.act("3") # Review
            agent.wait_for_text("Review")
            agent.act("<Enter>") # Reveal
            agent.wait_for_text("Grade:")
            agent.act("3") # Grade Good
            agent.wait_for_text("✓ Good")
            time.sleep(2.0)
            
            # Check Dashboard for badge
            agent.act("1")
            agent.wait_for_text("GOAL MET 🏆")
            agent.assert_text("GOAL MET 🏆")
            
        finally:
            agent.close()
