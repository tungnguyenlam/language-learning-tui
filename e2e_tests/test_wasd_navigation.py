import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path so we can import it
sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_wasd_view_switching():
    """Test that WASD keys can switch between views like arrow keys"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=120, lines=45)
        try:
            agent.assert_text("DASHBOARD")

            # Forward cycle
            # Dashboard -> Decks
            agent.act("s")
            agent.wait_for_text("Decks", timeout=5.0)

            # Decks -> Review
            agent.act("s")
            agent.wait_for_text("Review", timeout=5.0)

            # Review -> Statistics
            agent.act("s")
            agent.wait_for_text("Statistics", timeout=5.0)

            # Statistics -> Import
            agent.act("s")
            agent.wait_for_text("Import / Export", timeout=5.0)

            # Backward cycle
            # Import -> Statistics
            agent.act("w")
            agent.wait_for_text("Statistics", timeout=5.0)

            # Statistics -> Review
            agent.act("w")
            agent.wait_for_text("Review", timeout=5.0)

            # Clean verification based on the confirmed order in handlers.go
            # Cycle: Dashboard(1) -> Decks(2) -> Review(3) -> Statistics(4) -> Import(5) -> AI(6) -> Settings(7) -> Browser(8) -> Cram(9) -> Practice(0)
            
            agent.act("1")
            agent.wait_for_text("DASHBOARD")
            
            agent.act("s") # Dashboard -> Decks
            agent.wait_for_text("Decks")
            
            agent.act("s") # Decks -> Review
            agent.wait_for_text("Review")
            
            agent.act("w") # Review -> Decks
            agent.wait_for_text("Decks")
            
            agent.act("<Left>") # Decks -> Dashboard
            agent.wait_for_text("DASHBOARD")

        finally:
            agent.close()


def test_wasd_navigation_preserves_existing_functions():
    """Test that WASD keys don't interfere with existing functions"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=120, lines=45)
        try:
            # We'll use Review mode to test 'a' (Again) vs view switching
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("S") # Seed (Uppercase S is safer to avoid wasd interference)
            agent.wait_for_regex(r"Imported \d+ notes", timeout=90.0)
            
            agent.act("3") # Review
            agent.wait_for_text("Review 1/", timeout=10.0)
            
            # Space to reveal
            agent.act(" ")
            agent.wait_for_text("Grade: a Again", timeout=5.0)
            
            # Press 'a' to grade as Again. If it was view switching, we'd go to AI view.
            agent.act("a")
            agent.wait_until_stable()
            
            # Verify we are still in Review (should show Review 1/ or similar, or just Review header)
            agent.assert_text("REVIEW")
            # We check the header instead of full screen because AI Drafts is in the sidebar
            screen = agent.observe()
            assert "│ REVIEW │" in screen
            
        finally:
            agent.close()


if __name__ == "__main__":
    pytest.main(["-v", __file__])
