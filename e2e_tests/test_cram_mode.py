import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_cram_review_flow():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create some bookmarked cards
        agent = start_agent(tmpdir, columns=100, lines=30)
        try:
            # Go to review and bookmark a card
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("b")
            agent.wait_for_text("Card bookmarked")
            
            # Go to Cram mode
            agent.act("9")
            agent.wait_for_text("Cram Mode")
            agent.wait_for_text("Filter: bookmarked")
            agent.wait_for_text("> [FC]") # Should show the bookmarked card
            
            # Start Cram Review
            agent.act("<Enter>")
            agent.wait_for_text("Cram Review")
            agent.wait_for_text("Press Space or Enter to reveal.")            
            # Reveal
            agent.act("<Space>")
            agent.wait_for_text("Grade:")
            agent.wait_for_text("cramRevealed")
            
            # Grade
            agent.act("g")
            # After grading the last card, should return to selection
            agent.wait_for_text("Cram Mode")
            agent.assert_text("1 cards loaded")
            agent.wait_for_text("Cram Mode") # Should be back to list
            
        finally:
            agent.close()
