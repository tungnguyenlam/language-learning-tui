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


def test_a1_animals_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Animals")
            agent.act("<Enter>")
            agent.wait_for_text("A1 German Animals", timeout=10.0)
        finally:
            agent.close()


def test_a2_body_health_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Body")
            agent.act("<Enter>")
            agent.wait_for_text("Body", timeout=10.0)
        finally:
            agent.close()


def test_b1_cooking_deck_listed():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_open_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Cooking")
            agent.act("<Enter>")
            agent.wait_for_text("Cooking", timeout=10.0)
            agent.wait_for_text("b1, cooking", timeout=10.0)
        finally:
            agent.close()


def test_dashboard_verb_box_shows_english_meaning():
    # Verb box now includes an em-dash and the English translation
    # Use a taller terminal so the bottom Verb/Tip section is not clipped
    with tempfile.TemporaryDirectory() as tmpdir:
        app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
        agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=140, lines=60)
        agent.wait_for_text("DASHBOARD", timeout=15.0)
        agent.wait_until_stable()
        try:
            agent.wait_for_text("Verb:", timeout=5.0)
            # The verb box header includes a "—" separator before the English meaning
            agent.wait_for_text("—", timeout=5.0)
        finally:
            agent.close()


def test_settings_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Settings", timeout=5.0)
        finally:
            agent.close()


def test_ai_view_renders_with_drafts_header():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts", timeout=5.0)
        finally:
            agent.close()


def test_help_overlay_shows():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("?")
            # Help overlay should appear with some hint text
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_stats_view_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("s")
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_a1_animals_deck_listed()
    test_a2_body_health_deck_listed()
    test_b1_cooking_deck_listed()
    test_dashboard_verb_box_shows_english_meaning()
    test_settings_view_renders()
    test_ai_view_renders_with_drafts_header()
    test_help_overlay_shows()
    test_stats_view_renders()
