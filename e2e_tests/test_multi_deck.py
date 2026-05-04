import sys
import os
import tempfile
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=50):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_multi_deck_shows_all_decks_by_default():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Deck: All Decks")
            agent.assert_text("Due cards:   52")
            agent.assert_text("active)")
        finally:
            agent.close()

def test_multi_deck_review_flow():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3") # Start review
            agent.wait_for_text("Review 1/52")
            agent.act("<Enter>") # Reveal
            agent.wait_for_text("g Good")
            agent.act("g") # Grade Good
            agent.wait_for_text("Review 1/51")
        finally:
            agent.close()

def test_grammar_tip_hides_on_small_screens():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Start with a tall screen to see the tip
        agent = start_agent(tmpdir, columns=100, lines=50)
        try:
            agent.assert_text("Grammar Tip:")
        finally:
            agent.close()

    with tempfile.TemporaryDirectory() as tmpdir:
        # Start with a short screen where the tip shouldn't fit (height <= 20)
        agent = start_agent(tmpdir, columns=100, lines=18)
        try:
            agent.wait_until_stable()
            text = agent.screen.get_screen_text()
            assert "Grammar Tip:" not in text
        finally:
            agent.close()
