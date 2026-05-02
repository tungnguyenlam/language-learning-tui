import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=90, lines=30):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_mouse_navigation_tabs():
    """Verify that clicking each tab in the navigation bar switches to the correct view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Dashboard is default
            agent.assert_text("DASHBOARD")
            
            # Click Decks (Tab 2) - In 90 col medium layout
            # Dashboard: X=0..11, Decks: X=12..19, Review: X=20..28
            agent.click(15, 3) 
            agent.wait_for_text("Decks")
            agent.assert_text("German A1 Survival")
            
            # Click Statistics (Tab 4) - X=29..41
            agent.click(35, 3)
            agent.wait_for_text("Statistics")
            agent.assert_text("Total Cards:")
            
            # Click Browser (Tab 8) - Tab 5: 42..50, Tab 6: 51..55, Tab 7: 56..66, Tab 8: 67..76
            agent.click(70, 3)
            agent.wait_for_text("Browser")
            agent.assert_text("52 cards found")
            
            # Click Dashboard again (Tab 1)
            agent.click(5, 3)
            agent.wait_for_text("DASHBOARD")
        finally:
            agent.close()

def test_browser_selection_and_history_toggle():
    """Verify that j/k in the browser moves selection and Enter toggles history."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Initially first card selected
            agent.wait_for_text(">")
            
            # Move down
            agent.act("j")
            agent.wait_until_stable()
            
            # Toggle history (Enter)
            agent.act("<Enter>")
            agent.wait_for_text("Review History:")
            agent.assert_text("No reviews yet.")
            
            # Toggle off
            agent.act("<Enter>")
            agent.assert_not_text("Review History:")
        finally:
            agent.close()

def test_statistics_scrollbar_click():
    """Verify that clicking the scrollbar track in Statistics registers a click."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Statistics
            agent.act("4")
            agent.wait_for_text("Statistics")
            
            # In medium layout (90 wide), panel is 86 wide.
            # Content layout X=3, Y=4. lineWidth = 84.
            # Scrollbar X = 3 + 84 + 1 = 88 (0-based) = 89 (1-based)
            agent.click(89, 10)
            
            # Verify the mouse click is registered in the status line
            # 1-based (89, 10) is 0-based (88, 9)
            agent.wait_for_text("mouse 88,9")
        finally:
            agent.close()
