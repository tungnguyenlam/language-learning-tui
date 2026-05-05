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
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dashboard_polished_layout():
    """Test that the dashboard has the new group headers."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Review Queue")
            agent.assert_text("Collection")
            agent.assert_text("Today's Progress")
            # Verify it's not the old "Collection Stats"
            assert "Collection Stats" not in agent.observe()
        finally:
            agent.close()

def test_review_polished_layout():
    """Test that the review view still works with the new header style."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to review
            agent.act('3')
            agent.wait_for_text("Review")
            
            # Reveal
            agent.act('<Space>')
            agent.wait_until_stable()
            
            agent.assert_text("Grade: a Again")
        finally:
            agent.close()

def test_ai_drafting_status():
    """Test that the AI drafting status and spinner area appear."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to AI
            agent.act('6')
            agent.wait_for_text("AI Drafts")
            
            # Type something and enter
            agent.act('h')
            agent.act('a')
            agent.act('l')
            agent.act('l')
            agent.act('o')
            agent.act('<Enter>')
            
            # We check for the status message
            # It might be "Generating draft..." or already "1 draft ready"
            agent.wait_for_text("draft")
            
            # And for the temporary drafting message in AI view
            # Since offline provider is fast, it might already be done.
            # But let's check if we can see it or the result.
            screen = agent.observe()
            assert "draft ready" in screen or "Generating" in screen or "->" in screen
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
