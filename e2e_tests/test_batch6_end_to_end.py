import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester import TUIAgent


def start_agent(tmpdir, columns=120, lines=46):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def type_text(agent, text):
    for char in text:
        agent.act(char)


def seed_standard(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Imported", timeout=15.0)


def select_deck(agent, query, visible_text):
    agent.act("2")
    agent.wait_for_text("DECK LIST")
    agent.act("/")
    type_text(agent, query)
    agent.wait_for_text(visible_text)
    agent.act("<Enter>")
    agent.wait_for_text("Press enter to select deck.")
    agent.act("<Enter>")
    agent.wait_for_text("DASHBOARD")


def test_bureaucracy_deck_searchable_in_decks_view():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard(agent)
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            type_text(agent, "bureaucracy")
            agent.wait_for_text("German B1 Bureaucracy")
            agent.assert_text("appointments")
        finally:
            agent.close()


def test_digital_privacy_deck_searchable_in_decks_view():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard(agent)
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            type_text(agent, "privacy")
            agent.wait_for_text("B2 Digital Priva")
            agent.assert_text("Data protection")
        finally:
            agent.close()


def test_browser_finds_bureaucracy_vocabulary():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard(agent)
            select_deck(agent, "bureaucracy", "German B1 Bureaucracy")
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.act("/")
            type_text(agent, "Meldebescheinigung")
            agent.wait_for_text("Meldebescheinigung")
            agent.assert_text("registration certificate")
        finally:
            agent.close()


def test_browser_finds_digital_privacy_vocabulary():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard(agent)
            select_deck(agent, "privacy", "German B2 Digital Priva")
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.act("/")
            type_text(agent, "Datenschutz")
            agent.wait_for_text("Datenschutz")
            agent.assert_text("data protection")
        finally:
            agent.close()


def test_ai_empty_topic_guard_is_visible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.act("<Esc>")
            agent.act("<Enter>")
            agent.wait_for_text("Enter a topic before generating AI drafts")
        finally:
            agent.close()


def test_ai_suggestions_include_new_practical_topics():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=140)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.assert_text("B1 bureaucracy appointment")
            agent.assert_text("B2 digital privacy")
        finally:
            agent.close()


def test_statistics_forecast_renders_without_wrapping_status_area():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics")
            agent.assert_text("Forecast:")
            agent.assert_text("10 left")
            agent.assert_text("Use j/k or Mouse Wheel to scroll")
        finally:
            agent.close()


def test_mouse_click_side_navigation_opens_ai_view():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.click(5, 10)
            agent.wait_for_text("AI Drafts")
            agent.assert_text("B2 digital privacy")
        finally:
            agent.close()
