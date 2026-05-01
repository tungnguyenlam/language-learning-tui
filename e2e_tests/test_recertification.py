import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=90, lines=28):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_tab_cycles_through_all_primary_views_and_wraps():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            for text in [
                "Press enter to select deck.",
                "Review 1/6",
                "Statistics",
                "Import / Export",
                "AI Drafts",
                "Settings",
                "Card Browser",
                "Cram Mode",
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
            agent.act("3")
            agent.wait_for_text("Review 1/6")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("h")
            agent.wait_for_text("5 cards due")
            agent.act("1")
            agent.wait_for_text("Due cards:   5")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text("Due cards:   5")
            restarted.act("3")
            restarted.wait_for_text("Review 1/5")
        finally:
            restarted.close()


def test_settings_provider_toggle_persists_after_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("AI Provider: disabled")
            agent.act("<Enter>")
            agent.wait_for_text("AI Provider: offline")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.act("7")
            restarted.wait_for_text("AI Provider: offline")
        finally:
            restarted.close()
