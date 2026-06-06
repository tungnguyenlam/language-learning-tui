import os
import sys
import tempfile
import time

sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)
from tui_tester import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=100, lines=50)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_dashboard_header_visible():
    """Verify dashboard header shows WILLKOMMEN and daily progress"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("WILLKOMMEN", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_grammar_tip_visible():
    """Verify grammar tip is displayed on dashboard"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Grammar Tip", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_verb_of_day_visible():
    """Verify Verb of the Day is displayed on dashboard"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Verb:", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_stats_visible():
    """Verify dashboard stats boxes are visible"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Review Queue", timeout=5.0)
            agent.wait_for_text("Collection", timeout=5.0)
        finally:
            agent.close()


def test_ai_view_renders():
    """Verify AI view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts", timeout=5.0)
            agent.wait_for_text("Click a topic", timeout=5.0)
        finally:
            agent.close()


def test_statistics_view_renders():
    """Verify Statistics view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics", timeout=5.0)
        finally:
            agent.close()


def test_browser_view_renders():
    """Verify Browser view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser", timeout=5.0)
            agent.wait_for_text("Search:", timeout=5.0)
        finally:
            agent.close()


def test_settings_view_renders():
    """Verify Settings view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("SETTINGS", timeout=5.0)
        finally:
            agent.close()


def test_cram_view_renders():
    """Verify Cram view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram Mode", timeout=5.0)
        finally:
            agent.close()


def test_import_view_renders():
    """Verify Import view renders correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export", timeout=5.0)
        finally:
            agent.close()


def test_help_overlay():
    """Verify help overlay can be toggled"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.act("?")
            agent.wait_for_text("Keyboard Shortcuts", timeout=5.0)
            agent.act("?")
            agent.wait_until_stable()
        finally:
            agent.close()


if __name__ == "__main__":
    test_dashboard_header_visible()
    test_dashboard_grammar_tip_visible()
    test_dashboard_verb_of_day_visible()
    test_dashboard_stats_visible()
    test_ai_view_renders()
    test_statistics_view_renders()
    test_browser_view_renders()
    test_settings_view_renders()
    test_cram_view_renders()
    test_import_view_renders()
    test_help_overlay()
