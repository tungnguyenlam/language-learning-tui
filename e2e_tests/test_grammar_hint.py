import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=120, lines=44)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def _reveal_first_card(agent):
    agent.act("3")
    agent.wait_until_stable()
    agent.act("<Space>")
    time.sleep(0.5)
    agent.wait_until_stable()
    agent.act("<Space>")
    time.sleep(0.5)
    agent.wait_until_stable()


def test_grammar_hint_appears_for_adjective():
    # The starter deck begins with adjectives (colors: blau, gelb, ...).
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _reveal_first_card(agent)
            agent.wait_for_text("ADJ", timeout=3.0)
            agent.wait_for_text("Forms:", timeout=3.0)
            agent.wait_for_text("comparative", timeout=3.0)
            agent.wait_for_text("Example:", timeout=3.0)
        finally:
            agent.close()


def test_grammar_hint_hidden_before_reveal():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("3")
            agent.wait_until_stable()
            # No reveal yet - grammar hint should NOT be visible
            text = agent.observe()
            assert "Forms:" not in text, "Grammar hint Forms: should not be visible before reveal"
            assert "ADJ" not in text or "ADJ " not in text.replace("\n", " "), \
                "ADJ badge should not be visible before reveal"
        finally:
            agent.close()
