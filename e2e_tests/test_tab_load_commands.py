import os
import sys
import tempfile

sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir} -test-mode", columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent


def tab_to(agent, expected_text, count):
    for _ in range(count):
        agent.act("<Tab>")
    agent.wait_for_text(expected_text)
    agent.wait_until_stable()


def test_tab_to_browser_loads_cards():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            tab_to(agent, "Card Browser", 7)
            agent.wait_for_text("blau")
            agent.assert_text("[FC] blau")
            agent.assert_text("52 cards found")
        finally:
            agent.close()


def test_tab_to_cram_loads_bookmarked_cards():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("b")
            agent.wait_for_text("Card bookmarked")

            tab_to(agent, "Cram Mode", 6)
            agent.wait_until_stable()
            agent.wait_for_text("Filter: bookmarked")
            agent.assert_text("> [FC] blau")
            agent.assert_text("1 cards loaded.")
        finally:
            agent.close()


def test_tab_to_statistics_renders_persisted_progress():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/52")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("51 cards due")

            tab_to(agent, "Statistics", 1)
            agent.assert_text("Total Reviews: 1")
            agent.assert_text("Reviews Today: 1/10")
            agent.wait_for_text("Statistics")
        finally:
            agent.close()
