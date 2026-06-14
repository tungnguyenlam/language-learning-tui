import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=90, lines=28):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_review_grade_status_is_single_line():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("cards due")
            agent.assert_text("Review 1/51")
        finally:
            agent.close()


def test_import_missing_file_error_status_is_concise():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("i")
            agent.wait_for_text("status: Error: no such file or directory: import.tsv")
        finally:
            agent.close()


def test_browser_deck_switch_reloads_cards():
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, "w", encoding="utf-8") as f:
            f.write("#separator:tab\n")
            f.write("#deck:Travel A1\n")
            f.write("travel-1\tdie Bahn\ttrain\t\ttransport\tTravel A1\tBasic\n")

        agent = start_agent(tmpdir, columns=110, lines=30)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("i")
            agent.wait_for_text("status: Imported 1 notes from import.tsv")

            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("blau")

            agent.act("]")
            agent.wait_for_text("die Bahn")
            agent.assert_not_text("blau")
        finally:
            agent.close()
