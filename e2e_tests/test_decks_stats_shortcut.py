import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

from e2e_helpers import read_due_count

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_decks_view_stats_shortcut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("2") # Decks view
            agent.wait_for_text("Decks")

            # Navigate to the first real deck (after All Decks)
            agent.act("j")

            # Press 'v' to view stats
            agent.act("v")

            # Should be in Statistics view for that deck
            agent.wait_for_text("Statistics: German A1 Survival", timeout=10.0)
            # Fresh database: every starter card (including MCQs) is due.
            agent.assert_text(f"Total Cards: {due}")
        finally:
            agent.close()
