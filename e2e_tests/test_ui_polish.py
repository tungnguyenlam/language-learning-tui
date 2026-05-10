import time
import tempfile
import os
import sys

# Add the project root to sys.path to import tui_tester
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from tui_tester.agent import TUIAgent

def start_agent(tmpdir):
    db_path = os.path.join(tmpdir, "test.db")
    # Use the compiled binary if it exists, otherwise use go run
    bin_path = os.environ.get("DEUTSCH_TUI_BIN", "go run ../cmd/deutsch-tui")
    cmd = f"{bin_path} -data-dir {tmpdir}"
    return TUIAgent(cmd, columns=100, lines=50)

def test_ui_polish_and_content():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # 1. Verify Dashboard (Grammar Tip Example)
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Grammar Tip:")
            # Since tips are daily, we don't know exactly which one, 
            # but we can check for "Example:" if the screen is tall enough
            # Nominalization tip has an example.
            agent.wait_for_text("Example:")

            # 2. Verify New Content (Confusable Words)
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("S") # Seed
            agent.wait_for_text("Seeding standard content...", timeout=5.0)
            # Wait for import to finish
            agent.wait_for_text("Imported", timeout=10.0)
            agent.wait_for_text("Standard Content", timeout=10.0)
            agent.wait_until_stable()

            agent.act("2") # Decks
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.wait_for_text("Search:")
            agent.act("Confusable")
            agent.act("<Enter>")
            agent.wait_for_text("German Confusable Words")
            
            # 3. Verify Scrollbar in Decks
            # By now we have many decks from Seed, so scrollbar should be there.
            # We can check for the scrollbar character '│' or '█'
            agent.wait_for_text("│") 

            # 4. Verify Statistics view Maturity Distribution
            agent.act("4")
            agent.wait_for_text("Statistics:")
            agent.wait_for_text("Maturity Distribution")
            agent.wait_for_text("New:")
            
            # 5. Verify the new deck is actually study-able
            agent.act("2") # Decks
            agent.wait_for_text("DECK LIST")
            # Select Confusable words (it should be filtered)
            agent.act("<Enter>") # Select deck and go to dashboard
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Active Deck: German Confusable Words")
            
            agent.act("3") # Review
            agent.wait_for_text("Review") # Title case
            # Should see one of our confusable words
            # The order might be random or sequential, but let's check for one
            # Note: since they are new, they are in New state.
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_ui_polish_and_content()
