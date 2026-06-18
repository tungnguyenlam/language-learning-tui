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
    bin_path = os.environ.get("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    cmd = f"{bin_path} -data-dir {tmpdir} -test-mode"
    return TUIAgent(cmd, columns=120, lines=50)

def test_ui_polish_and_content():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # 1. Verify Dashboard (Grammar Tip)
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Grammar Tip:")

            # 2. Seed content
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("S") # Seed
            agent.wait_for_text("Seeding standard content...", timeout=5.0)
            agent.wait_for_text("Imported", timeout=30.0)
            agent.wait_until_stable(timeout=15.0)
            time.sleep(2.0)

            # 3. Go back to Dashboard first to ensure state is stable
            agent.act("1")
            agent.wait_for_text("DASHBOARD", timeout=10.0)
            agent.wait_until_stable()
            
            # 4. Verify Review view works
            agent.act("3") # Review
            agent.wait_for_text("Review 1/", timeout=10.0)
            agent.wait_for_text("Deck:", timeout=5.0)
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_ui_polish_and_content()
