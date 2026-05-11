import os
import sys
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

@pytest.fixture
def agent():
    app = TUIAgent("go run cmd/deutsch-tui/main.go", columns=100, lines=50)
    yield app

def test_new_a1_decks_loaded(agent: TUIAgent):
    import time
    # Wait for the app to start and dashboard to load
    agent.wait_for_text("DASHBOARD")
    
    # Go to Import / Export view to seed content
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Seeding standard content")
    time.sleep(2.0)
    agent.wait_until_stable()
    
    # Go to Decks view
    agent.act("2")
    agent.wait_for_text("DECK LIST")

    # Search for the new decks by name
    agent.act("/")
    agent.act("Hobbies")
    agent.act("<Enter>")
    agent.wait_for_text("A1 Hobbies & Free Time")
    
    agent.act("/")
    agent.act("Food")
    agent.act("<Enter>")
    agent.wait_for_text("A1 Food & Drink")
    
    agent.act("/")
    agent.act("Travel")
    agent.act("<Enter>")
    agent.wait_for_text("A1 Travel & Transport")

def test_ai_draft_empty_state(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    
    # Go to AI Drafts view
    agent.act("6")
    agent.wait_for_text("AI Drafts")
    
    # Verify the new empty state box is there
    agent.wait_for_text("Ready to create new flashcards")
    agent.wait_for_text("Type a topic and press Enter to generate.")

def test_dashboard_verb_of_the_day(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    # Verify verb section exists
    agent.wait_for_text("Verb:")
    agent.wait_for_text("ich")
    agent.wait_for_text("wir")

def test_dashboard_grammar_tip(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    # Verify grammar tip section exists
    agent.wait_for_text("Grammar Tip:")

def test_statistics_view_renders(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    agent.act("4")
    agent.wait_for_text("Statistics:")
    agent.wait_for_text("Collection")

def test_cram_mode_navigation(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    agent.act("9")
    agent.wait_for_text("Cram Mode")
    agent.wait_for_text("Click a filter to load cards")

def test_ai_draft_suggestions(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    agent.act("6")
    agent.wait_for_text("AI Drafts")
    agent.wait_for_text("Click a topic")
    agent.wait_for_text("A1 survival")
    agent.wait_for_text("B1 doctor visit")

def test_browser_view_loaded(agent: TUIAgent):
    agent.wait_for_text("DASHBOARD")
    agent.act("8")
    agent.wait_for_text("Card Browser")
    agent.wait_for_text("Search:")
