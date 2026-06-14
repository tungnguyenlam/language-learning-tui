import time
import os
import sys
import tempfile
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

@pytest.fixture
def agent():
    tmpdir = tempfile.mkdtemp()
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    
    # Seed standard content so we have the new decks
    agent.act("5") # Import view
    agent.wait_for_text("Import / Export")
    agent.act("S") # Seed standard
    # Wait for one of the new decks to appear in the 'Current Deck' line or similar
    agent.wait_for_text("A1 Animals & Nature", timeout=30.0)
    agent.act("1") # Back to dashboard
    agent.wait_for_text("DASHBOARD")
    
    yield agent
    agent.close()

@pytest.mark.e2e
def test_new_medical_deck_exists(agent):
    """Verify that the new A2 Medical Appointment deck is available and searchable."""
    agent.act("2")  # Go to Decks view
    agent.wait_for_text("DECK LIST")
    
    # Search for the new deck specifically
    agent.act("/")
    agent.act("appointment")
    agent.act("<Enter>")
    
    agent.wait_for_text("Medical Appointment")

@pytest.mark.e2e
def test_search_highlights_in_browser(agent):
    """Verify that search terms are found in the Browser view."""
    # First select the medical deck in Decks view
    agent.act("2")
    agent.wait_for_text("DECK LIST")
    agent.act("/")
    agent.act("appointment")
    agent.act("<Enter>") # Finish search
    time.sleep(0.5)
    agent.act("<Enter>") # Select and go to Dashboard
    agent.wait_for_text("DASHBOARD")
    
    # Now go to Browser view
    agent.act("8")  # Go to Browser view
    agent.wait_for_text("Card Browser")
    
    # Search for a common term from the new medical deck
    agent.act("/")
    agent.act("Termin")
    agent.act("<Enter>")
    
    # Verify the text is present.
    agent.wait_for_text("einen Termin vereinbaren")
    
@pytest.mark.e2e
def test_native_tts_fallback_smoke(agent):
    """Smoke test for audio playback initiation."""
    agent.act("3")  # Start review
    
    # Check if we are in review
    agent.wait_for_text("Review")
    
    # Press 'p' to play audio.
    agent.act("p")
    # It might show "Generating..." or "Playing audio..."
    agent.wait_for_text("audio")
