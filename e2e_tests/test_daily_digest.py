import sys
import os
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_daily_digest_renders_on_dashboard():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("Daily Digest")
            agent.assert_text("M:0 Y:0 N:52")
        finally:
            agent.close()

def test_daily_digest_shows_due_cards_message():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("52 cards waiting.")
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
