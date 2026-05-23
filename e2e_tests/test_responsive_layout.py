import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

@pytest.mark.parametrize("width", [40, 60, 100])
def test_dashboard_responsive(width):
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=width, lines=40)
        try:
            agent.assert_text("DASHBOARD")
            agent.assert_text("Review Queue")
            agent.assert_text("Collection")
            agent.assert_text("Today's Progress")
            agent.assert_text("Daily Digest")
        finally:
            agent.close()

@pytest.mark.parametrize("width", [40, 60, 100])
def test_statistics_responsive(width):
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=width, lines=40)
        try:
            agent.act('4') # Switch to Statistics
            agent.wait_for_text("Statistics", timeout=5.0)
            agent.assert_text("Collection")
            agent.assert_text("Today's Performance")
            # "Maturity Distribution" might be scrolled off in narrow views
            if width >= 100:
                agent.assert_text("Maturity Distribution")
        finally:
            agent.close()

@pytest.mark.parametrize("width", [30, 50, 100])
def test_decks_responsive(width):
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=width, lines=40)
        try:
            agent.act('2') # Switch to Decks
            agent.wait_for_text("DECK LIST", timeout=5.0)
            agent.assert_text("All Decks")
        finally:
            agent.close()
