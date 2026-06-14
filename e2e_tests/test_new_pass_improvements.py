import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dashboard_verb_of_the_day():
    """Verify that Verb of the Day is visible on the Dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        agent.wait_for_text("DASHBOARD")
        agent.wait_for_text("Verb:")

def test_review_timer_and_badge():
    """Verify session timer and state badges in Review view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        agent.wait_for_text("DASHBOARD")
        # Go to Review
        agent.act("3")
        agent.wait_for_text("Review")
        
        # Check for timer
        agent.wait_for_text("⏱ 00:")
        
        # Check for badge (should be NEW or LEARNING)
        # Use wait_for_text and check if either exists
        screen = agent.observe()
        assert "NEW" in screen or "LEARNING" in screen or "MATURE" in screen

def test_browser_interval_and_history():
    """Verify intervals and search history in Browser view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        agent.wait_for_text("DASHBOARD")
        
        # 1. Review the first card (should be 'blau' in a fresh DB due to ID order)
        agent.act("3") # Review
        agent.wait_for_text("Review")
        agent.wait_for_text("blau")
        
        agent.act(" ") # Reveal
        agent.wait_for_text("Again") # Wait for reveal to finish
        agent.act("g") # Good
        
        # 2. Go to Browser and search for 'blau'
        agent.act("8") 
        agent.wait_for_text("Card Browser")
        
        agent.act("/")
        agent.wait_for_text("SEARCHING")
        agent.act("blau")
        agent.act("<Enter>")
        
        # 3. Check for history and interval
        agent.act("/")
        agent.wait_for_text("History: blau")
        agent.act("<Esc>")
        
        # Check for interval (1 day or some hours)
        # FSRS might give "1 day" or "4 days" for Good
        agent.wait_for_text("(", timeout=10.0) # Just check if parentheses appear
        screen = agent.observe()
        assert "day" in screen or "hour" in screen
