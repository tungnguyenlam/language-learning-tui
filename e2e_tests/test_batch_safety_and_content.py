import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent


REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
HEALTHCARE_DECK_PATH = os.path.join(
    REPO_ROOT, "internal", "content", "testdata", "german-decks", "b2-healthcare-systems.tsv"
)


def start_agent(tmpdir, columns=110, lines=42):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_ai_disabled_warning_and_generate_guard():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("Settings")
            # Cycle through every provider so we land back on disabled.
            for want in ["offline", "template", "openai", "anthropic", "disabled"]:
                agent.act("<Enter>")
                agent.wait_for_text(f"AI Provider:    {want}")
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.assert_text("AI provider is disabled")
            agent.act("<Enter>")
            agent.wait_for_text("AI provider is disabled. Enable it in Settings.")
        finally:
            agent.close()


def test_import_path_backspace_handles_umlaut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.wait_until_stable()
            agent.act("<Ctrl-u>")
            agent.act("k")
            agent.act("ä")
            agent.wait_for_text("Import file: kä_")
            agent.act("<Backspace>")
            agent.wait_for_text("Import file: k_")
        finally:
            agent.close()


def test_export_tag_backspace_handles_umlaut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("t")
            agent.wait_until_stable()
            agent.act("k")
            agent.act("ä")
            agent.wait_for_text("Export Filter: kä_")
            agent.act("<Backspace>")
            agent.wait_for_text("Export Filter: k_")
        finally:
            agent.close()


def test_ai_topic_backspace_handles_umlaut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.wait_until_stable()
            agent.act("k")
            agent.act("ä")
            agent.wait_for_text("Topic: kä_")
            agent.act("<Backspace>")
            agent.wait_for_text("Topic: k_")
        finally:
            agent.close()


def test_settings_template_backspace_handles_umlaut():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.act("<Enter>")
            agent.wait_for_text("AI Provider:    offline")
            agent.act("<Enter>")
            agent.wait_for_text("AI Provider:    template")
            agent.act("j")
            agent.wait_until_stable()
            agent.act("<Enter>")
            agent.wait_for_text("EDITING")
            agent.act("x")
            agent.wait_for_text("{{.Topic}}x")
            agent.act("<Backspace>")
            agent.act("<Enter>")
            agent.wait_for_text("Front Template: {{.Topic}}")
        finally:
            agent.close()


def test_deck_limits_left_arrow_stays_in_decks_view():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("Decks")
            agent.act("j")
            agent.wait_until_stable()
            agent.act("L")
            agent.wait_for_text("Limits: New")
            agent.act("<Left>")
            agent.wait_until_stable()
            agent.assert_text("Decks")
            agent.assert_text("Limits: New")
        finally:
            agent.close()


def test_right_arrow_navigation_still_cycles_views():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("DASHBOARD")
            agent.act("<Right>")
            agent.wait_for_text("Decks")
            agent.act("<Right>")
            agent.wait_for_text("Review")
        finally:
            agent.close()


def test_import_healthcare_systems_deck_is_selectable():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=120, lines=42)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.wait_until_stable()
            agent.act("<Ctrl-u>")
            for char in HEALTHCARE_DECK_PATH:
                agent.act(char)
            agent.act("<Enter>")
            agent.wait_until_stable()
            agent.act("i")
            agent.wait_for_text("Imported")
            agent.act("2")
            agent.wait_for_text("Decks")
            agent.assert_text("B2 Healthcare Systems")
        finally:
            agent.close()
