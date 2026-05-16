import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from tui_tester.agent import TUIAgent

def start_agent(tmpdir, columns=100, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', './deutsch-tui-bin')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

@pytest.mark.e2e
def test_grading_shortcuts_1_to_4():
    """Verify that 1, 2, 3, 4 keys work for grading in Review view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Review view (press 3)
            agent.act("3")
            agent.wait_for_text("Review")
            
            if "No cards due" in agent.observe():
                # Seed standard content if no cards due
                agent.act("5") # Import
                agent.wait_for_text("IMPORT")
                agent.act("S") # Seed
                agent.wait_for_text("Standard content seeded")
                agent.act("3") # Back to Review
                agent.wait_for_text("Review")

            # Reveal card
            agent.act(" ")
            agent.wait_for_text("Grade:")
            
            # Test grading with '3' (Good) -> ✓ Good
            agent.act("3")
            agent.wait_for_text("✓ Good")
            
            # Reveal again (next card)
            agent.act(" ")
            agent.wait_for_text("Grade:")
            
            # Test grading with '1' (Again) -> ✗ Again
            agent.act("1")
            agent.wait_for_text("✗ Again")
            
            # Test grading with '2' (Hard) -> ~ Hard
            agent.act(" ")
            agent.wait_for_text("Grade:")
            agent.act("2")
            agent.wait_for_text("~ Hard")

            # Test grading with '4' (Easy) -> ★ Easy
            agent.act(" ")
            agent.wait_for_text("Grade:")
            agent.act("4")
            agent.wait_for_text("★ Easy")
            
        finally:
            agent.close()
