import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_mcq_no_view_switch_after_answer():
    with tempfile.TemporaryDirectory() as tmpdir:
        tsv_path = os.path.join(tmpdir, "mcq_bug.tsv")
        with open(tsv_path, "w") as f:
            # Create an MCQ card
            f.write("id-mcq\tPrompt\tAnswer\t\ttag1\tDeck1\tMCQ:Choice1|||Choice2|||Answer\n")
        
        agent = start_agent(tmpdir)
        try:
            # Import the MCQ card
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.act("<Ctrl-u>")
            for char in tsv_path:
                agent.act(char)
            agent.act("<Enter>")
            agent.act("i")
            time.sleep(2.0)
            
            # Go to Review
            agent.act("3")
            agent.wait_for_text("Review")
            agent.wait_for_text("Prompt")
            
            # The first card is likely the Flashcard (id-mcq:front)
            # Reveal it and grade it to get to the MCQ card
            agent.act(" ")
            agent.wait_for_text("Answer")
            agent.act("g")
            
            # Now we should be on the MCQ card
            agent.wait_for_text("Review 1/1")
            
            # Reveal/Pick choice
            agent.act("3") # Pick "Answer" (3rd choice)
            agent.wait_for_text("Answer")
            agent.wait_for_text("Correct") # Assuming it says Correct
            
            # Now we are in RevealRevealed state and mcqAnswered is true.
            # Pressing '1' should NOT switch to Dashboard because we added it as a grading shortcut.
            agent.act("1")
            
            # If it switched, we'd see "DASHBOARD"
            # Since we fixed it, it should stay in Review (or go to Session Summary if last card)
            time.sleep(1.0)
            screen = agent.observe()
            if "DASHBOARD" in screen:
                 print("BUG REPRODUCED: '1' switched view after MCQ answer")
                 pytest.fail("BUG: '1' switched to dashboard")
            else:
                 print("SUCCESS: '1' did not switch to dashboard")
                 # It might go to Session Summary if it was the last card
                 if "SESSION SUMMARY" not in screen:
                     agent.assert_text("Review")
                     agent.assert_text("Again")
                 else:
                     agent.assert_text("SESSION SUMMARY")
        finally:
            agent.close()

if __name__ == "__main__":
    test_mcq_no_view_switch_after_answer()
