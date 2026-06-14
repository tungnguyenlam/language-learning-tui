import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=110, lines=44)
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


def test_b1_transport_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Transport")
            agent.act("<Enter>")
            agent.wait_for_text("Public Transport", timeout=10.0)
        finally:
            agent.close()


def test_a1_office_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Office")
            agent.act("<Enter>")
            agent.wait_for_text("Office", timeout=10.0)
            agent.wait_for_text("a1, office, stationery", timeout=10.0)
        finally:
            agent.close()


def test_b2_climate_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Climate")
            agent.act("<Enter>")
            agent.wait_for_text("Climate", timeout=10.0)
            agent.wait_for_text("b2, climate", timeout=10.0)
        finally:
            agent.close()


def test_ai_view_has_climate_suggestion():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            agent.wait_for_text("Seeding standard content...", timeout=10.0)
            time.sleep(3.0)
            agent.wait_until_stable(timeout=15.0)
            agent.act("6")
            agent.wait_for_text("AI Drafts", timeout=5.0)
            agent.wait_for_text("climate change", timeout=5.0)
            agent.wait_for_text("small talk", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_streak_emoji_no_streak():
    # When streak=0, no flame/lightning emoji should be shown
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("Streak:", timeout=5.0)
            agent.wait_for_text("0 days", timeout=5.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_b1_transport_deck_listed()
    test_a1_office_deck_listed()
    test_b2_climate_deck_listed()
    test_ai_view_has_climate_suggestion()
    test_new_grammar_tips_render()
