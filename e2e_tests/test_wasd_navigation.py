import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_wasd_view_switching():
    """Test that WASD keys can switch between views like arrow keys"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.assert_text("Dashboard")
            
            # Test 's' key for next view (like right arrow)
            agent.act('s')
            agent.wait_until_stable()
            agent.assert_text("Press enter to select deck.")
            
            agent.act('s')
            agent.wait_until_stable()
            agent.assert_text("Review 1/")
            
            agent.act('s')
            agent.wait_until_stable()
            agent.assert_text("Import / Export")
            
            agent.act('s')
            agent.wait_until_stable()
            agent.assert_text("AI Drafts")
            
            agent.act('s')
            agent.wait_until_stable()
            agent.assert_text("Settings")
            
            # Test 'w' key for previous view (like left arrow)
            agent.act('w')
            agent.wait_until_stable()
            agent.assert_text("AI Drafts")
            
            agent.act('w')
            agent.wait_until_stable()
            agent.assert_text("Import / Export")
            
            agent.act('w')
            agent.wait_until_stable()
            agent.assert_text("Review 1/")
            
            agent.act('w')
            agent.wait_until_stable()
            agent.assert_text("Press enter to select deck.")
            
            agent.act('w')
            agent.wait_until_stable()
            agent.assert_text("Use Review to start studying.")
        finally:
            agent.close()

def test_wasd_navigation_preserves_existing_functions():
    """Test that WASD keys don't interfere with existing functions"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            # Go to Review view
            agent.act('3')
            agent.wait_until_stable()
            agent.assert_text("Review 1/6")
            
            # Reveal card
            agent.act('<Space>')
            agent.wait_until_stable()
            agent.assert_text("Grade: a Again")
            
            # 'a' should grade as "Again", not switch views
            agent.act('a')
            agent.wait_until_stable()
            agent.assert_text("5 cards due")
            
            # Go to AI view
            agent.act('5')
            agent.wait_until_stable()
            agent.assert_text("AI Drafts")
            
            # Generate draft
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("der Kaffee -> German prompt for der")
            
            # 'a' should approve draft, not switch views
            agent.act('a')
            agent.wait_until_stable()
            agent.assert_text("Draft approved")
            
            # Now test 'd' for discard
            # First generate another draft
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("der Kaffee -> German prompt for der")
            
            # 'd' should discard, not switch to Settings (view 6)
            agent.act('d')
            agent.wait_until_stable()
            agent.assert_text("AI Drafts") # Should still be in AI drafts
            agent.assert_text("Topic: der Kaffee") # Should have discarded the generated draft and returned to input
            
            # Verify we didn't switch to Settings
            agent.assert_not_text("AI Provider:")
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])