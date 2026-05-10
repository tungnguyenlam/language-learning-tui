import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=118, lines=44):
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


def test_dashboard_card_mix_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("Card Mix")
            agent.assert_text("New")
            agent.assert_text("Young")
            agent.assert_text("Mature")
        finally:
            agent.close()


def test_review_header_shows_deck_and_type():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review")
            agent.wait_for_text("Deck:")
            agent.assert_text("Type:")
        finally:
            agent.close()


def test_browser_preview_shows_kind_and_extra():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Card Preview:")
            agent.assert_text("Kind")
            agent.assert_text("Extra")
        finally:
            agent.close()


def test_ai_suggestions_include_levels_and_mobility():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.assert_text("B2 urban mobility")
            agent.assert_text("C1 business email")
            agent.assert_text("Suggested levels:")
        finally:
            agent.close()


def test_settings_provider_cycle_hint_renders():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("SETTINGS")
            agent.assert_text("Provider cycle:")
            agent.assert_text("disabled -> offline -> template")
        finally:
            agent.close()


def test_seeded_mobility_deck_is_searchable():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard_content(agent)
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("mobility")
            agent.act("<Enter>")
            agent.wait_for_text("B2 Urban Mobility", timeout=10.0)
        finally:
            agent.close()


def test_seeded_mobility_browser_cards_visible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_standard_content(agent)
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("mobility")
            agent.act("<Enter>")
            agent.wait_for_text("B2 Urban Mobility", timeout=10.0)
            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.act("/")
            agent.act("Mobilitätswende")
            agent.act("<Enter>")
            agent.wait_for_text("Mobilitätswende", timeout=10.0)
            agent.assert_text("Card Preview:")
        finally:
            agent.close()


def test_review_reveal_keeps_metadata_visible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review")
            agent.act("<Enter>")
            agent.wait_for_text("Grade: a Again", timeout=10.0)
            agent.assert_text("Deck:")
            agent.assert_text("Type:")
        finally:
            agent.close()
