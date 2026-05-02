import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f"go run ./cmd/deutsch-tui -data-dir {tmpdir}", columns=columns, lines=lines)
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
            agent.wait_for_text("der Apfel")
            agent.assert_text("[FC] der Apfel")
            agent.assert_text("21 cards found")
        finally:
            agent.close()


def test_tab_to_cram_loads_bookmarked_cards():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/21")
            agent.act("b")
            agent.wait_for_text("Card bookmarked")

            tab_to(agent, "Cram Mode", 6)
            agent.wait_for_text("Filter: bookmarked")
            agent.assert_text("> [FC] der Apfel")
            agent.assert_text("1 cards in cram mode")
        finally:
            agent.close()


def test_tab_to_statistics_renders_persisted_progress():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_for_text("Review 1/21")
            agent.act("<Space>")
            agent.wait_for_text("Grade: a Again")
            agent.act("g")
            agent.wait_for_text("20 cards due")

            tab_to(agent, "Statistics", 1)
            agent.assert_text("Total Reviews: 1")
            agent.assert_text("Reviews Today: 1/10")
            agent.assert_text("good: 1")
        finally:
            agent.close()
