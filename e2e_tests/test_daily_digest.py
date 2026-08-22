import sys
import os
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

from e2e_helpers import read_due_count

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_daily_digest_renders_on_dashboard():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("Daily Digest")
            # Fresh database: every starter card is new.
            agent.assert_text(f"M:0 Y:0 N:{read_due_count(agent)}")
        finally:
            agent.close()

def test_daily_digest_shows_due_cards_message():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text(f"{read_due_count(agent)} cards waiting.")
        finally:
            agent.close()

def test_dashboard_boxes_are_single_line_headers():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("Today's Progress")
            agent.assert_text("Review Queue")
            agent.assert_text("Collection")
        finally:
            agent.close()
