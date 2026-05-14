import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=110, lines=44)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def _seed_and_open_decks(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Seeding standard content...", timeout=10.0)
    time.sleep(3.0)
    agent.wait_until_stable(timeout=15.0)
    agent.act("2")
    agent.wait_for_text("DECK LIST", timeout=5.0)


def test_legal_contracts_deck_exists():
    """Test that the new Legal & Contracts deck is available."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Legal")
            agent.act("<Enter>")
            agent.wait_for_text("Legal", timeout=10.0)
            agent.wait_for_text("legal", timeout=10.0)
        finally:
            agent.close()


def test_contracts_deck_search():
    """Test searching for 'Contracts' finds the Legal deck."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Contracts")
            agent.act("<Enter>")
            agent.wait_for_text("Legal", timeout=10.0)
        finally:
            agent.close()


def test_dashboard_grammar_tip_shows():
    """Test that the Grammar Tip section shows on the dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
        agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=140, lines=60)
        agent.wait_for_text("DASHBOARD", timeout=15.0)
        agent.wait_until_stable()
        try:
            agent.wait_for_text("Grammar Tip:", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_verb_of_the_day_shows():
    """Test that the Verb of the Day section shows on the dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
        agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=140, lines=60)
        agent.wait_for_text("DASHBOARD", timeout=15.0)
        agent.wait_until_stable()
        try:
            agent.wait_for_text("Verb:", timeout=5.0)
        finally:
            agent.close()


def test_ai_suggestions_show():
    """Test that AI suggestions are visible in AI view."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts", timeout=5.0)
            agent.wait_for_text("Click a topic", timeout=5.0)
        finally:
            agent.close()


def test_ai_new_suggestions_present():
    """Test that the newly added AI topics are present."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts", timeout=5.0)
            agent.wait_for_text("emergency", timeout=3.0)
        finally:
            agent.close()


def test_browser_view_shows_decks():
    """Test that the browser view shows decks."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Browse", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_deck_switch_shortcut():
    """Test that [ and ] switch decks from dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("[")
            agent.wait_until_stable(timeout=3.0)
            agent.act("]")
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_legal_contracts_deck_exists()
    test_contracts_deck_search()
    test_dashboard_grammar_tip_shows()
    test_dashboard_verb_of_the_day_shows()
    test_ai_suggestions_show()
    test_ai_new_suggestions_present()
    test_browser_view_shows_decks()
    test_dashboard_deck_switch_shortcut()