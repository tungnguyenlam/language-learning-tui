import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_robust_deck_search():
    """Verify that '/' starts a robust search in Decks view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a second deck via import.tsv
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("#deck:Personal Phrases\n")
            f.write("Hallo\thello\ta1\n")

        agent = start_agent(tmpdir)
        try:
            # Import the deck
            agent.act('5')
            agent.wait_for_text("Import / Export")
            agent.act('i')
            agent.wait_for_text("Imported")
            
            # Go to Decks view
            agent.act('2')
            agent.wait_for_text("Decks")
            agent.assert_text("German A1 Survival")
            agent.assert_text("Personal Phrases")
            
            # Press / to search
            agent.act('/')
            agent.wait_for_text("Search: _")
            
            # Type 'personal'
            agent.act('p')
            agent.act('e')
            agent.act('r')
            agent.wait_for_text("Search: per_")
            agent.assert_text("Personal Phrases")
            
            # Verify other deck is gone
            screen = agent.observe()
            assert "German A1 Survival" not in screen
            
            # Finish search
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("Filter: per (Press / to edit)")
            
            # Clear filter
            agent.act('<Esc>')
            agent.wait_until_stable()
            agent.assert_text("German A1 Survival")
        finally:
            agent.close()

def test_settings_daily_goal_increase():
    """Verify that increasing the daily goal in Settings view works."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=110, lines=40)
        try:
            # Go to Settings view (Tab 7)
            agent.act('7')
            agent.wait_for_text("Settings")
            agent.assert_text("Daily Goal: 10")
            
            # Press + to increase goal
            agent.act('+')
            agent.wait_for_text("Daily Goal: 11")
            
            # Press - to decrease goal
            agent.act('-')
            agent.wait_for_text("Daily Goal: 10")
        finally:
            agent.close()

def test_decks_view_selection():
    """Verify that selecting a deck in Decks view updates active deck and switches to Dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a second deck
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("#deck:Travel\n")
            f.write("Zug\ttrain\ta1\n")

        agent = start_agent(tmpdir)
        try:
            # Import
            agent.act('5')
            agent.act('i')
            agent.wait_for_text("Imported")
            
            # Go to Decks view
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Move down to 'Travel' deck (All Decks -> Survival -> Travel)
            # Or use search to find it and select
            agent.act('/')
            agent.act('T')
            agent.act('r')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Select it
            agent.act('<Enter>')
            agent.wait_for_text("DASHBOARD")
            
            # Active deck should be Travel
            agent.assert_text("Deck: Travel")
        finally:
            agent.close()
