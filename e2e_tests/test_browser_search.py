import sys
import os
import tempfile
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("WILLKOMMEN!", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_browser_search_filtering():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Start search
            agent.act("/")
            agent.wait_for_text("SEARCHING")
            
            # Type 'Hallo' (exists in Starter deck)
            agent.act("H")
            agent.act("a")
            agent.act("l")
            agent.act("l")
            agent.act("o")
            
            # It should filter real-time
            agent.wait_for_text("Hallo")
            
            # Clear search with Esc
            agent.act("<Esc>")
            agent.wait_for_text("Card Browser")
            agent.wait_until_stable()
            
            # Should see other cards now
            agent.assert_text("blau")
            
        finally:
            agent.close()

def test_browser_tag_filtering():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Start tag filter
            agent.act("#")
            agent.wait_for_text("FILTER BY TAG")
            
            # Type 'a1'
            agent.act("a")
            agent.act("1")
            
            agent.wait_for_text("#a1")
            
            # Clear with Esc
            agent.act("<Esc>")
            agent.wait_for_text("Card Browser")
            
        finally:
            agent.close()
