import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_ui_polish_browser_preview():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()

            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.wait_for_text("Card Preview:")
            agent.wait_for_text("Deck: ")
            agent.wait_for_text(" | ") # Pipe separator should be present
        finally:
            agent.close()

def test_ui_polish_stats_footer():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("4")
            agent.wait_for_text("Statistics:")
            agent.wait_for_text("Lines 1-")
            agent.wait_for_text("Use j/k or Mouse Wheel")
        finally:
            agent.close()

def test_ui_polish_review_empty_box():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            # Should have cards to review, so wait for Review header
            agent.wait_for_text("Review")
            # Exit review
            agent.act("q")
        finally:
            agent.close()

def test_deck_search_clear():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("ZZZZZ")
            agent.act("<Enter>")
            agent.wait_for_text("No decks match search.")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.wait_for_text("All Decks") # should be visible again
        finally:
            agent.close()

def test_browser_search_clear():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()

            agent.act("8")
            agent.wait_for_text("Card Browser")
            agent.act("/")
            agent.act("ZZZZZ")
            agent.act("<Enter>")
            agent.wait_for_text("No cards found")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.wait_for_text("Card Preview:") # should be visible again
        finally:
            agent.close()

def test_cram_filter_navigation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram Mode")
            agent.wait_for_text("Click a filter to load cards")
            agent.act("j")
            agent.act("j")
            agent.act("<Enter>")
            agent.wait_for_text("No cards found for this filter.")
        finally:
            agent.close()

if __name__ == "__main__":
    test_ui_polish_browser_preview()
    test_ui_polish_stats_footer()
    test_ui_polish_review_empty_box()
    test_deck_search_clear()
    test_browser_search_clear()
    test_cram_filter_navigation()
