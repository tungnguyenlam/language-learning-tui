import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_typing_exercise_mode():
    """Test typing exercise mode in review."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to decks to get some cards
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Select first deck
            agent.act('<Enter>')
            agent.wait_for_text("DASHBOARD")
            
            # Go back to review
            agent.act('3')
            agent.wait_for_text("Review")
            
            # Should have cards now
            assert "Session complete!" not in agent.observe()
            
            # Type a wrong answer
            agent.act('t')  # Toggle typing mode
            agent.wait_for_text("Type your answer")
            
            agent.act('w')
            agent.act('r')
            agent.act('o')
            agent.act('n')
            agent.act('g')
            
            # Submit
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should show incorrect feedback
            screen = agent.observe()
            assert "✗" in screen or "Incorrect" in screen or "Your answer:" in screen
            
        finally:
            agent.close()

def test_typing_cancel():
    """Test cancelling typing mode."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to review with cards
            agent.act('2')  # Decks
            agent.wait_for_text("Decks")
            agent.act('<Enter>')  # Select deck
            agent.wait_for_text("DASHBOARD")
            agent.act('3')  # Review
            agent.wait_for_text("Review")
            
            # Enter typing mode
            agent.act('t')
            agent.wait_for_text("Type your answer")
            
            # Cancel with escape
            agent.act('<Esc>')
            agent.wait_for_text("Typing mode off")
            
            # Should be back to normal review
            screen = agent.observe()
            assert "Press space or enter to reveal" in screen
            
        finally:
            agent.close()

def test_cram_all_shortcut():
    """Test cram all shortcut from decks view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to decks
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Use 'c' to cram all cards from current deck
            agent.act('c')
            agent.wait_for_text("Cram")
            
            # Should be in cram mode
            screen = agent.observe()
            assert "Filter:" in screen or "cards in cram mode" in screen
            
            # Start cram
            agent.act('<Enter>')
            agent.wait_for_text("Press Space or Enter to reveal.")

            # Reveal
            agent.act(' ')
            agent.wait_for_text("Grade:")  # or similar indicator
            
        finally:
            agent.close()

def test_deck_filter_by_tag():
    """Test filtering decks by tag."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to decks
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Open search
            agent.act('/')
            agent.wait_for_text("Search: _")
            
            # Search for a tag that should exist
            agent.act('g')
            agent.act('e')
            agent.act('r')
            agent.act('m')
            agent.act('a')
            agent.act('n')
            agent.act('i')
            agent.act('c')
            
            agent.wait_for_text("Search: german")
            
            # Exit search
            agent.act('<Enter>')
            
            # Should show filtered decks
            screen = agent.observe()
            # Either shows decks or "No decks match search" if none found
            assert "german" in agent.observe().lower() or "No decks match" in agent.observe()
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])