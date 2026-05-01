import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

def start_agent(tmpdir):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}")
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_audio_indicator_in_review():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a TSV with audio field
        tsv_path = "/tmp/audio_test.tsv"
        with open(tsv_path, "w") as f:
            f.write("id-1\tprompt-1\tanswer-1\textra-1\ttags\tdeck-1\tBasic\taudio-1.mp3\n")
        
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")
            
            # Clear path and type new one
            agent.act("<Enter>") # START EDITING
            agent.act("<Ctrl-u>")
            import time
            for char in tsv_path:
                agent.act(char)
                time.sleep(0.05)
            agent.act("<Enter>") # FINISH EDITING
            agent.act("i") # EXECUTE IMPORT TSV
            agent.wait_for_text("Imported 1 notes", timeout=10.0)
            
            # Select the new deck (use left arrows to avoid capturing '2')
            agent.act("<Left>")
            agent.act("<Left>")
            agent.act("<Left>")
            agent.wait_for_text("deck-1")
            agent.act("<Enter>")
            
            # Go to Review
            agent.act("3")
            agent.wait_for_text("Review 1/1")
            agent.wait_for_text("[Audio]")
            agent.wait_for_text("p audio")
            
        finally:
            agent.close()
