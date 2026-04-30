import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

def start_agent(tmpdir):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}")
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_audio_autoplay_toggle_and_persistence():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Settings
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.wait_for_text("Auto-play audio: off")
            
            # Navigate to Auto-play toggle (index 5)
            for _ in range(5):
                agent.act("j")
            
            # Toggle it on
            agent.act("<Enter>")
            agent.wait_for_text("Auto-play audio: on")
            agent.wait_for_text("Auto-play audio enabled")
            
            # Restart and check persistence
            agent.close()
            agent = start_agent(tmpdir)
            agent.act("7")
            agent.wait_for_text("Auto-play audio: on")
            
        finally:
            agent.close()

def test_audio_autoplay_in_review():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Use a simple path to avoid accidental command triggers from path characters
        tsv_path = "/tmp/autoplay.tsv"
        with open(tsv_path, "w") as f:
            f.write("id-1\tprompt-1\tanswer-1\textra-1\ttags\tdeck-1\tBasic\tbeep.mp3\n")
        
        agent = start_agent(tmpdir)
        try:
            # Enable auto-play first
            agent.act("7")
            agent.wait_for_text("Settings")
            for _ in range(5):
                agent.act("j")
            agent.act("<Enter>")
            agent.wait_for_text("Auto-play audio: on")
            
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")
            
            # Clear path and type new one
            agent.act("<Ctrl-u>")
            time.sleep(0.1)
            for char in tsv_path:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            time.sleep(1.0)
            
            # Select the new deck
            agent.act("2")
            agent.wait_for_text("Decks")
            # Clear any accidental filter
            agent.act("<Esc>")
            agent.wait_for_text("deck-1")
            # Navigate to the second deck (the imported one)
            agent.act("j")
            agent.act("<Enter>")
            
            # Go to Review
            agent.act("3")
            agent.wait_for_text("Review 1/1")
            agent.wait_for_text("[Audio]")
            
            # Reveal card - should trigger auto-play
            agent.act(" ") # Reveal
            agent.wait_for_text("answer-1")
            agent.wait_for_text("Auto-playing audio: beep.mp3")
            
        finally:
            agent.close()
