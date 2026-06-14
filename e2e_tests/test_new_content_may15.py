import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir, columns=110, lines=44):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_dashboard_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, lines=50)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
        finally:
            agent.close()


def test_ai_view_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("AI Drafts", timeout=5.0)
            agent.wait_until_stable(timeout=2.0)
        finally:
            agent.close()


def test_statistics_view_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_decks_view_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_browser_view_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Browser", timeout=5.0)
            agent.wait_until_stable(timeout=2.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_dashboard_renders()
    test_ai_view_accessible()
    test_statistics_view_accessible()
    test_decks_view_accessible()
    test_browser_view_accessible()