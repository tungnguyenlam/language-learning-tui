import os
import pytest
import tempfile
import sys

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_settings_goal_keyboard_interaction():
    """Verify daily goal can be adjusted via +/- keys in Settings."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Settings
            agent.act('7')
            agent.wait_for_text("Settings")
            
            # Navigate to Goal (index 5)
            for _ in range(5):
                agent.act('j')
            agent.wait_for_text("> Daily Goal: 10")
            
            # Increment
            agent.act('+')
            agent.wait_for_text("Daily Goal: 11")
            
            # Decrement twice
            agent.act('-')
            agent.act('-')
            agent.wait_for_text("Daily Goal: 9")
        finally:
            agent.close()

def test_ai_template_set_navigation():
    """Verify switching AI template sets in Settings."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Settings
            agent.act('7')
            agent.wait_for_text("Settings")
            agent.wait_for_text("Template Set: vocabulary")
            
            # Use [ / ] to switch set
            agent.act(']')
            # Should switch to next set alphabetically (articles, conjugation, grammar, vocabulary)
            # If vocabulary was default, next is articles? No, articles is first.
            # Wait! I sorted them. articles, conjugation, grammar, vocabulary.
            # defaultIndex was set to vocabulary (index 3).
            # nextIndex is articles (index 0).
            agent.wait_for_text("Template Set: articles")
            
            agent.act(']')
            agent.wait_for_text("Template Set: conjugation")
        finally:
            agent.close()
