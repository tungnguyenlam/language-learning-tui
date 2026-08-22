import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

from e2e_helpers import read_due_count


def start_agent(tmpdir, columns=110, lines=42):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir} -test-mode", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_dashboard_digest_and_progress_render():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Today's Progress")
            agent.assert_text("Daily Digest")
            agent.assert_text("Next: blau")
            agent.assert_text("Use Review (3) to start studying.")
        finally:
            agent.close()


def test_review_reveal_grade_and_persist_due_count():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.assert_text("Press space or enter to reveal.")
            agent.act("<Enter>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text(f"Review 1/{due - 1}")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text(f"Due cards:   {due - 1}")
        finally:
            restarted.close()


def test_import_view_renders_filters_and_actions():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.assert_text("Export Deck: All Decks")
            agent.assert_text("Export Filter: (None)")
            agent.assert_text("[i] Import TSV")
            agent.assert_text("[x] Export TSV")
        finally:
            agent.close()


def test_import_empty_path_guard_status():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.act("<Ctrl-u>")
            agent.act("<Enter>")
            agent.act("i")
            agent.wait_for_text("Import path is empty")
        finally:
            agent.close()


def test_ai_view_topic_placeholder_and_tip():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("6")
            agent.wait_for_text("Topic:")
            agent.assert_text("apartment viewing")
            agent.assert_text("Tip: include level and use case")
        finally:
            agent.close()


def test_settings_goal_text_and_template_help():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=110, lines=45)
        try:
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.assert_text("Daily Goal: 10 cards")
            agent.assert_text("AI CONFIGURATION")
        finally:
            agent.close()


def test_deck_navigation_lists_starter_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("Decks")
            agent.assert_text("All Decks")
            agent.assert_text("German A1 Survival")
            agent.assert_text("Press enter to select deck")
        finally:
            agent.close()


def test_import_backup_writes_progress_file():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.assert_text("[B] Backup")
            agent.assert_text("[U] Restore")
            agent.act("B")
            agent.wait_for_text("Backed up", timeout=10.0)
        finally:
            agent.close()

        backups_dir = os.path.join(tmpdir, "backups")
        names = os.listdir(backups_dir)
        assert any(name.startswith("backup-") and name.endswith(".db") for name in names), names


def test_view_tabs_cover_core_screens():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            for key, text in [
                ("4", "Statistics"),
                ("5", "Import / Export"),
                ("6", "AI Drafts"),
                ("7", "Settings"),
                ("8", "Browser"),
                ("9", "Cram"),
            ]:
                agent.act(key)
                agent.wait_for_text(text)
        finally:
            agent.close()
