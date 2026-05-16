import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from tui_tester.agent import TUIAgent

def start_agent(tmpdir, columns=120, lines=60):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', './deutsch-tui-bin')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

@pytest.mark.e2e
def test_c1_social_issues_deck_and_cloze_visuals():
    """Verify C1 deck exists and Cloze visuals/footer are present."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Seed content
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            agent.wait_for_text("notes from Standard Content")
            time.sleep(2.0)
            
            # Find C1 Social Issues deck
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            for char in "C1 Social Issues":
                agent.act(char)
            agent.act("<Enter>") # Select
            agent.wait_for_text("C1 Social Issues & Society")
            agent.wait_until_stable()
            
            # Start review
            agent.act("3")
            agent.wait_for_text("Review")
            # If still All caught up, wait a bit and try again
            if "All caught up" in agent.observe():
                time.sleep(2.0)
                agent.act("3")
                agent.wait_for_text("Review")
            
            agent.assert_text("C1 Social Issues & Society")
            
            # Check for interactive footer (substrings)
            agent.assert_text("h hint")
            agent.assert_text("b bookmark")
            agent.assert_text("history")
            
            # Check for Cloze visual indicator in prompt
            agent.assert_text("Die")
            agent.assert_text("debate/clash") # This is inside the [hint]
            agent.assert_text("mit der Vergangenheit")
            
        finally:
            agent.close()

@pytest.mark.e2e
def test_global_help_remains_functional():
    """Verify that '?' help still works after UI changes."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("?")
            agent.wait_for_text("Keyboard Shortcuts")
            # New multi-column layout text
            agent.assert_text("Global:")
            agent.assert_text("Review:")
            agent.assert_text("Dashboard/Decks:")
            agent.act("?") # Toggle off
            agent.wait_for_text("DASHBOARD")
        finally:
            agent.close()
