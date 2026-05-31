import sys
import os
import tempfile
import pytest

# Add project root to sys.path so we can import tui_tester package
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=40):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_recently_studied():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Active deck might be anything due to seeding
            # Review one card
            agent.act("3")
            agent.wait_for_text("Review")
            agent.act(" ") # Reveal
            agent.wait_for_text("Again") # Grade predictions show up
            agent.act("g") # Grade
            agent.wait_until_stable()
            
            # Back to Dashboard
            agent.act("1")
            agent.wait_for_text("Recently Studied", timeout=10.0)
            # Just check that there's at least one bullet point with a deck name
            agent.assert_text("•")
        finally:
            agent.close()

def test_business_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Seed standard content first
            agent.act("5") # Import
            agent.wait_for_text("Import")
            agent.act("S") # Seed
            agent.wait_for_text("Standard Content", timeout=30.0)

            # Go to Decks view
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            # Search for Business
            agent.act("/")
            # Use backspace to be safe
            agent.act("<Backspace>" * 20)
            agent.act("business")
            agent.act("<Enter>")
            agent.wait_for_text("Business", timeout=10.0)
        finally:
            agent.close()

def test_card_preview_in_browser():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Card Preview:")
            # Should show preview of first card
            agent.assert_text("Front")
            agent.assert_text("Back")
        finally:
            agent.close()

def test_debug_view_toggle():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Ctrl+D to open debug
            agent.act("<Ctrl-D>")
            agent.wait_for_text("DEBUG LOG")
            agent.assert_text("Active View: debug")
            
            # Ctrl+D to exit
            agent.act("<Ctrl-D>")
            agent.wait_for_text("DASHBOARD")
        finally:
            agent.close()

def test_focus_mode_in_review():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3") # Review
            agent.wait_for_text("Review")
            # Toggle focus mode
            agent.act("f")
            agent.wait_for_text("Focus Mode Active")
            # Check that header info is hidden (or replaced by the banner)
            # 'Session:' string should be gone
            text = agent.screen.get_screen_text()
            assert "Session:" not in text
            
            # Toggle off
            agent.act("f")
            agent.wait_for_text("Session:")
        finally:
            agent.close()

def test_improved_ai_templates():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # First go to Settings to ensure template provider is used
            agent.act("7")
            agent.wait_for_text("SETTINGS")
            # Loop until AI provider is 'template'
            # (disabled -> offline -> template)
            for _ in range(5):
                text = agent.screen.get_screen_text()
                if "AI Provider:    template" in text:
                    break
                agent.act("<Enter>")
                agent.wait_until_stable()

            agent.act("6") # AI view
            agent.wait_for_text("Topic:")
            
            # Type a topic
            agent.act("/")
            agent.act("<Backspace>" * 15) # Clear "der Kaffee"
            agent.act("Apfel")
            agent.act("<Enter>")
            
            # Wait for generation (Template provider)
            agent.wait_for_text("Apfel -> Translation:", timeout=10.0)
            agent.assert_text("Plural:")
            agent.assert_text("Gender:")
        finally:
            agent.close()

def test_grammar_tips_display():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, lines=45) # Need more height
        try:
            agent.wait_for_text("Grammar Tip:")
            # We don't know exactly which tip it is, but it should be there.
            text = agent.screen.get_screen_text()
            assert "Grammar Tip:" in text
        finally:
            agent.close()

def test_browser_bulk_tagging():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            
            # Select first two
            agent.act("m")
            agent.act("j")
            agent.act("m")
            
            # Bulk tags (T)
            agent.act("T")
            agent.wait_for_text("TAGS:")
            agent.act("newtag1,newtag2")
            agent.act("<Enter>")
            
            agent.wait_until_stable()
            # It should show individual tags now
            agent.assert_text("#newtag1")
            agent.assert_text("#newtag2")
        finally:
            agent.close()
