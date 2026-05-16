import time
import tempfile
import os
import sys

# Add the project root to sys.path to import tui_tester
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from tui_tester.agent import TUIAgent

def start_agent(tmpdir):
    # Use the compiled binary if it exists, otherwise use go run
    bin_path = os.environ.get("DEUTSCH_TUI_BIN", "go run ../cmd/deutsch-tui")
    cmd = f"{bin_path} -data-dir {tmpdir}"
    return TUIAgent(cmd, columns=120, lines=45)

def test_may15_batch8():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # 1. Seed content
            agent.wait_for_text("DASHBOARD")
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            time.sleep(1.0)
            agent.act("S") # Seed
            agent.wait_for_text("Imported", timeout=30.0)
            agent.wait_until_stable(timeout=15.0)
            time.sleep(2.0)
            
            # Go back to Dashboard to stabilize
            agent.act("1")
            agent.wait_for_text("DASHBOARD", timeout=10.0)
            agent.wait_until_stable()
            
            # 2. Verify Decks view has many decks
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=10.0)
            agent.act('<Esc>') # Clear any filter
            time.sleep(0.5)
            # Check that we have many decks (should show "of XX" at bottom)
            agent.wait_for_text(" of ", timeout=5.0)
            
            # 3. Verify AI Drafts view
            agent.act("6") # AI
            agent.wait_for_text("AI Drafts")
            agent.wait_for_text("apartment viewing")
            
            # 4. Verify Browser
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Card Preview:")
            
            # 5. Verify Statistics
            agent.act("4") # Statistics
            agent.wait_for_text("Statistics:")
            agent.wait_for_text("Total Cards:")
            
            # 6. Verify Dashboard "Card Mix"
            agent.act("1") # Dashboard
            agent.wait_for_text("Card Mix")
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_may15_batch8()
