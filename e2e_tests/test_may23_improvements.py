import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=20.0)
    agent.wait_until_stable()
    return agent

@pytest.mark.e2e
def test_false_friends_deck_exists():
    """Verify that the new False Friends deck is available and searchable."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        
        agent.act("5") # Import view
        agent.wait_for_text("Actions:")
        agent.act("S") # Seed standard content
        agent.wait_for_text("Ready", timeout=10.0)
        
        agent.act("2") # Decks view
        agent.wait_for_text("DECK LIST")
        
        # Search for the new deck
        agent.act("/")
        agent.act("mastery")
        agent.wait_for_text("Search: mastery")
        agent.act("<Enter>") # Finish searching
        
        agent.wait_for_text("Filter: mastery")
        agent.act("<Enter>") # Select the deck
        
        agent.wait_for_text("DASHBOARD")
        agent.assert_text("Active Deck: False Friends Mastery")

@pytest.mark.e2e
def test_dashboard_quick_actions_hitboxes():
    """Verify that Quick Actions on the dashboard are displayed."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        
        agent.assert_text("Quick Actions:")
        agent.assert_text("Review")
        agent.assert_text("Cram")
        agent.assert_text("Browser")
        agent.assert_text("Stats")
        agent.assert_text("Import")
        agent.assert_text("AI Draft")

@pytest.mark.e2e
def test_dashboard_recent_decks_navigation():
    """Verify that recent decks on the dashboard are accessible via shortcuts."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        
        # First, study a deck to make it appear in "Recently Studied"
        agent.act("2") # Decks
        agent.act("/")
        agent.act("survival")
        agent.wait_for_text("A1 Survival")
        agent.act("<Enter>") # Finish searching
        agent.act("<Enter>") # Select it
        
        agent.act("3") # Review
        agent.wait_for_text("Survival")
        agent.act("<Enter>") # Reveal
        agent.wait_for_text("Grade:")
        agent.act("g") # Good
        agent.act("1") # Back to dashboard
        
        agent.wait_for_text("Recently Studied")
        agent.assert_text("A1 Survival")
        
        # Now use the shortcut '!' to start review
        agent.act("!")
        agent.wait_for_text("Survival") # Should be back in Review view
