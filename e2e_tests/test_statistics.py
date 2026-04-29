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

def test_statistics_view_rendering_and_update():
    """Test that Statistics view renders and updates after reviews."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Statistics (key 4)
            agent.act('4')
            agent.wait_for_text("Statistics")
            
            # Initial stats (Survival deck has 6 cards, but let's be generic)
            agent.assert_text("Total Cards:")
            agent.assert_text("Total Reviews: 0")
            agent.assert_text("Success Rate:  0.0%")
            
            # Go to Review and grade one card
            agent.act('3')
            agent.wait_for_text("Review 1/")
            agent.act('<Space>')
            agent.wait_for_text("Grade: a Again")
            agent.act('g') # Good
            agent.wait_until_stable()
            
            # Go back to Statistics
            agent.act('4')
            agent.wait_for_text("Statistics")
            
            # Verify update
            agent.assert_text("Total Reviews: 1")
            agent.assert_text("Success Rate:  100.0%")
            agent.assert_text("good : 1")
            
        finally:
            agent.close()

def test_statistics_view_navigation():
    """Test that Statistics view is reachable via tab and arrows."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Tab navigation: Dashboard -> Decks -> Review -> Statistics
            agent.act('<Tab>') # Decks
            agent.wait_until_stable()
            agent.act('<Tab>') # Review
            agent.wait_until_stable()
            agent.act('<Tab>') # Statistics
            agent.wait_for_text("Statistics")
            
            # Arrow navigation: Statistics -> Import (Right)
            agent.act('<Right>')
            agent.wait_for_text("Import / Export")
            
            # Arrow navigation: Import -> Statistics (Left)
            agent.act('<Left>')
            agent.wait_for_text("Statistics")
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
