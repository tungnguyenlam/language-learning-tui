import os
import sys
import tempfile
import time
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
CIVIC_DECK_PATH = os.path.join(
    REPO_ROOT, "internal", "content", "testdata", "german-decks", "b2-public-services-civic-life.tsv"
)

def start_agent(tmpdir, columns=100, lines=80):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_c3_provider_persistence_on_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7") # Settings
            agent.wait_for_text("AI CONFIGURATION")
            # Toggle provider
            agent.act("<Enter>")
            agent.wait_for_text("AI Provider:    offline")
        finally:
            agent.close()
        
        # Restart and verify
        agent2 = start_agent(tmpdir)
        try:
            agent2.act("7") # Settings
            agent2.wait_for_text("AI Provider:    offline")
        finally:
            agent2.close()

def test_c3_reveal_grade_guard():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=120, lines=40)
        try:
            # Go to review (3)
            agent.act("3")
            agent.wait_for_text("Review 1/", timeout=10.0)
            agent.wait_for_text("Press space or enter to reveal.", timeout=5.0)
            
            # Try to grade prematurely - 'a' should do nothing when not revealed
            agent.act("a")
            time.sleep(0.3)
            # 'h' toggles hint, but should not reveal the answer
            agent.act("h")
            time.sleep(0.3)
            
            # The card should still be in unrevealed state (hint may be shown)
            # Check that grading options are NOT visible
            screen = agent.screen.get_screen_text()
            assert "Again (" not in screen or "Grade:" not in screen, "Grading should not be available before reveal"
            
            # Reveal
            agent.act("<Space>")
            agent.wait_for_text("Again", timeout=5.0)
        finally:
            agent.close()

def test_c3_browser_bulk_kind_toggle():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            
            # Select first card
            agent.act("m")
            agent.assert_text("1 cards selected")
            agent.assert_text("[FC]") # assuming first is FC
            
            # Bulk kind toggle
            agent.act("t")
            agent.wait_for_text("[MCQ]")
        finally:
            agent.close()

def test_c3_browser_bulk_bookmark_toggle():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            
            # Select first two cards
            agent.act("m")
            agent.act("j")
            agent.act("m")
            agent.assert_text("2 cards selected")
            
            # Bulk bookmark
            agent.act("b")
            agent.wait_until_stable()
            agent.assert_text("[B]")
        finally:
            agent.close()

def test_c3_import_visibility_civic():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=120)
        try:
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.wait_until_stable()
            # Clear text and type path
            agent.act("<Ctrl-u>")
            for char in CIVIC_DECK_PATH:
                agent.act(char)
            agent.act("<Enter>")
            agent.wait_until_stable()
            # Press i to import TSV
            agent.act("i")
            agent.wait_for_text("Imported")
            # Verify in Decks view
            agent.act("2")
            agent.wait_for_text("B2 Public Services")
        finally:
            agent.close()

def test_c3_help_overlay_content():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Open help
            agent.act("?")
            agent.wait_for_text("Keyboard Shortcuts")
            agent.wait_for_text("Switch to view")
            # Close help
            agent.act("?")
        finally:
            agent.close()

def test_c3_ai_draft_disabled_guard():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6") # AI View
            agent.wait_for_text("Topic:")
            # Try to generate when disabled
            agent.act("<Enter>")
            agent.wait_until_stable()
            text = agent.screen.get_screen_text()
            assert "disabled" in text.lower() or "not configured" in text.lower() or "error" in text.lower()
        finally:
            agent.close()
