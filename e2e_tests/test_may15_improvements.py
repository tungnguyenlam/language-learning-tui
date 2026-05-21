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
    return TUIAgent(cmd, columns=120, lines=60)

def test_may15_improvements():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # 1. Seed content and verify new decks
            agent.wait_for_text("DASHBOARD")
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("S") # Seed
            agent.wait_for_text("Imported", timeout=15.0)
            
            agent.act("2") # Decks
            agent.wait_for_text("DECK LIST")
            
            # Verify A2 Travel & Booking
            agent.act("/")
            agent.act("Travel")
            agent.act("<Enter>")
            agent.wait_for_text("A2 Travel")
            
            # Verify B1 Housing & Apartment
            agent.act("/")
            agent.act("Housing")
            agent.act("<Enter>")
            agent.wait_for_text("B1 Housing & Apartment")
            
            # Verify C1 Environment & Sustainability
            agent.act("/")
            agent.act("C1 Environment")
            agent.act("<Enter>")
            agent.wait_for_text("C1 Environment & Sustai")
            
            # 2. Verify Hint Feature in Review
            # Navigate to Review
            agent.act("3")
            agent.wait_for_text("Review 1/")
            agent.wait_for_text("h hint") # Verify shortcut is in help guide
            
            # Press 'h' to show hint
            agent.act("h")
            agent.wait_for_text("Hint shown")
            agent.wait_for_text("Hint:") # Even if hint value is empty, it shows "Hint: (no hint available)"
            
            # Toggle 'h' to hide hint
            agent.act("h")
            agent.wait_for_text("Hint hidden")
            
            # 3. Verify Dashboard Forecast
            agent.act("1") # Dashboard
            agent.wait_for_text("Review Queue")
            agent.wait_for_text("Next 24h:")
            
            # 4. Verify Statistics "Cards Added" chart
            agent.act("4") # Statistics
            agent.wait_for_text("Statistics:")
            agent.wait_for_text("Cards Added (Last 7 Days)")
            
            # Since we just added decks, today should have many cards added
            # The chart shows day names like Mon, Tue, etc.
            # We can't easily verify the exact bar, but we verify the section exists.
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_may15_improvements()
