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
            # Content layout X=3, Y=4. lineWidth = 78.
            # Scrollbar X = 3 + 78 + 1 = 82 (0-based) = 83 (1-based)
            agent.click(83, 10)
            
            # Verify the mouse click is registered in the status line
            # 1-based (83, 10) is 0-based (82, 9)
            agent.wait_for_text("mouse 82,9")
        finally:
            agent.close()

def test_statistics_drag_to_scroll():
    """Verify that dragging the scrollbar in Statistics scrolls the view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Statistics
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Lines 1-12 of 30")
            time.sleep(0.2) # Ensure hitboxes are stable
            
            # Drag scrollbar from top (approx row 6) to bottom (row 16)
            # In medium layout (90 wide), scrollbar is at X=83.
            # Row 0 is at Y=5, so Row 11 (last) is at Y=16.
            agent.drag_mouse(83, 6, 83, 16, steps=10)
            agent.wait_until_stable()
            
            # The view should have scrolled. We check for the last lines.
            agent.wait_for_text("Lines 19-30 of 30")
        finally:
            agent.close()

def test_ai_draft_interaction():
    """Verify that AI drafts can be approved/discarded via interactive buttons."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to AI View
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            
            # Clear existing topic
            agent.act("<Esc>")
            agent.wait_until_stable()

            # Type a topic
            agent.act("Hallo")
            agent.wait_for_text("Topic: Hallo")
            
            # Generate
            agent.act("<Enter>")
            agent.wait_for_text("Hallo ->", timeout=10.0) # Wait for generation
            
            # Should have at least one draft. 
            # Preview should be visible for the first one
            agent.wait_for_text("Preview:")
            
            # Click Discard on the first draft
            # In 90 col medium layout, AI view starts at Y=4. 
            # From debug logs, draft-discard-0 is at X:47 Y:10 (0-based)
            # So we click (48, 11) (1-based)
            agent.click(48, 11)
            
            # First draft should be gone
            agent.wait_for_text("No drafts yet.")
        finally:
            agent.close()
