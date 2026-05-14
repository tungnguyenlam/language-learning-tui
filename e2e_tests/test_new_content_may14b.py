import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir, columns=110, lines=44):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_dashboard_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_views_accessible_individually():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_until_stable(timeout=2.0)
            agent.act("3")
            agent.wait_until_stable(timeout=2.0)
            agent.act("4")
            agent.wait_until_stable(timeout=2.0)
            agent.act("5")
            agent.wait_until_stable(timeout=2.0)
            agent.act("6")
            agent.wait_until_stable(timeout=2.0)
            agent.act("7")
            agent.wait_until_stable(timeout=2.0)
            agent.act("8")
            agent.wait_until_stable(timeout=2.0)
            agent.act("9")
            agent.wait_until_stable(timeout=2.0)
        finally:
            agent.close()


def test_statistics_view_with_number_key():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=110, lines=30)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_decks_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_help_overlay_shortcut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=140, lines=60)
        try:
            agent.act("?")
            agent.wait_until_stable(timeout=2.0)
            agent.act("?")
            agent.wait_until_stable(timeout=2.0)
        finally:
            agent.close()


def test_browser_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Browser", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_cram_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram", timeout=5.0)
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_dashboard_view_renders()
    test_views_accessible_individually()
    test_statistics_view_with_number_key()
    test_decks_view_renders()
    test_help_overlay_shortcut()
    test_browser_view_renders()
    test_cram_view_renders()