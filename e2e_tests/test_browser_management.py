import sys
import os
import tempfile
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_browser_multi_select_and_bulk_bookmark():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Select first card
            agent.act("m")
            agent.assert_text("1 cards selected")
            agent.assert_text("[x]")
            
            # Select second card
            agent.act("j")
            agent.act("m")
            agent.assert_text("2 cards selected")
            
            # Bulk bookmark
            agent.act("b")
            agent.wait_until_stable()
            
            # Verify both have [B]
            # Since we refreshed, we might need to wait or check the screen
            agent.assert_text("[B]")
            
            # Clear selection with esc
            agent.act("m") # select another
            agent.assert_text("1 cards selected")
            agent.act("<Esc>")
            # Note: <Esc> might be tricky in some environments, but TUIAgent should handle it
            agent.wait_until_stable()
            text = agent.screen.get_screen_text()
            assert "cards selected" not in text
        finally:
            agent.close()

def test_browser_bulk_delete():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Select first two
            agent.act("m")
            agent.act("j")
            agent.act("m")
            agent.assert_text("2 cards selected")
            
            # Bulk delete (backspace)
            agent.act("<Backspace>")
            agent.wait_for_text("CONFIRM DELETION")
            agent.act("y")
            agent.wait_until_stable()
            # Status should show loading or ready
            text = agent.screen.get_screen_text()
            assert "2 cards selected" not in text
        finally:
            agent.close()

def test_browser_card_type_conversion():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            
            # First card is [FC] by default
            agent.assert_text("[FC]")
            
            # Toggle to MCQ
            agent.act("t")
            agent.wait_for_text("[MCQ]")
            
            # Toggle back to Flashcard
            agent.act("t")
            agent.wait_for_text("[FC]")
        finally:
            agent.close()
