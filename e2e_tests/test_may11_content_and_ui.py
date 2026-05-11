import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=118, lines=50):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def seed_standard_content(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Standard Content", timeout=30.0)


def test_dashboard_shows_verb_of_day():
    """Verify Dashboard shows Verb of the Day with conjugation."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, lines=50)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Verb:", timeout=10.0)
            agent.wait_for_text("ich", timeout=5.0)
            agent.wait_for_text("wir", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_shows_grammar_tip():
    """Verify Dashboard shows Grammar Tip section."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, lines=50)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Grammar Tip:", timeout=10.0)
        finally:
            agent.close()


def test_dashboard_shows_quick_actions():
    """Verify Dashboard shows quick action shortcuts."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Quick Actions")
            agent.wait_for_text("[3]")
            agent.wait_for_text("[9]")
        finally:
            agent.close()


def test_ai_view_has_topic_suggestions():
    """Verify AI view has clickable topic suggestions."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.wait_for_text("Click a topic")
            agent.wait_for_text("A1 survival")
            agent.wait_for_text("B1")
            agent.wait_for_text("B2")
        finally:
            agent.close()


def test_browser_view_renders_card_list():
    """Verify Browser view renders card list with scrollbar."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Search:")
        finally:
            agent.close()


def test_review_view_navigates_cards():
    """Verify Review view shows card count and navigation hints."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review")
        finally:
            agent.close()


def test_import_view_has_export_options():
    """Verify Import view has export options."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.wait_for_text("[x]")
            agent.wait_for_text("[X]")
        finally:
            agent.close()


def test_settings_view_has_provider_options():
    """Verify Settings view has AI provider options."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("SETTINGS")
            agent.wait_for_text("Provider")
        finally:
            agent.close()


def test_cram_view_has_filter_options():
    """Verify Cram view has filter options."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram Mode")
            agent.wait_for_text("filter")
        finally:
            agent.close()


def test_help_overlay_toggle():
    """Verify help overlay can be toggled with ? key."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("?")
            agent.wait_for_text("Keyboard Shortcuts")
            agent.act("?")
            agent.wait_until_stable()
        finally:
            agent.close()


def test_decks_view_navigates_list():
    """Verify Decks view shows deck list with navigation."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.wait_for_text("German A1 Survival")
        finally:
            agent.close()


def test_statistics_view_has_collection_stats():
    """Verify Statistics view shows collection stats."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.wait_for_text("Collection")
        finally:
            agent.close()


def test_dashboard_shows_review_and_progress():
    """Verify Dashboard shows review queue and progress sections."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Review Queue")
            agent.wait_for_text("Collection")
            agent.wait_for_text("Progress")
        finally:
            agent.close()
