import os
import sys
import tempfile

sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)

from tui_tester import TUIAgent

from e2e_helpers import read_due_count


def start_agent(tmpdir, columns=90, lines=28):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir} -test-mode", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_tab_cycles_through_all_primary_views_and_wraps():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            for text in [
                "Press enter to select deck.",
                f"Review 1/{due}",
                "Statistics",
                "Import / Export",
                "AI Drafts",
                "Settings",
                "Card Browser",
                "Cram Mode",
                "PRACTICE HUB",
                "Use Review (3) to start studying.",
            ]:
                agent.act("<Tab>")
                agent.wait_for_text(text)
                agent.assert_text(text)
        finally:
            agent.close()


def test_hard_grade_persists_review_progress_to_sqlite_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            due = read_due_count(agent)
            agent.act("3")
            agent.wait_for_text(f"Review 1/{due}")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("h")
            agent.wait_for_text(f"{due - 1} cards due")
            agent.act("1")
            agent.wait_for_text(f"Due cards:   {due - 1}")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text(f"Due cards:   {due - 1}")
            restarted.act("3")
            restarted.wait_for_text(f"Review 1/{due - 1}")
        finally:
            restarted.close()


def test_settings_provider_toggle_persists_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("AI Provider:    disabled")
            agent.act("<Enter>")
            agent.wait_for_text("AI Provider:    offline")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("7")
            restarted.wait_for_text("AI Provider:    offline")
        finally:
            restarted.close()
