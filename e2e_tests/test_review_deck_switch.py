import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=120, lines=44)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def _seed(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Seeding standard content...", timeout=10.0)
    time.sleep(3.0)
    agent.wait_until_stable(timeout=15.0)


def test_review_empty_state_shows_current_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed(agent)
            agent.act("3")
            agent.wait_until_stable()
            agent.wait_for_text("Current Deck:", timeout=5.0)
        finally:
            agent.close()


def test_review_deck_switch_updates_label():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed(agent)

            # Get initial deck label on Dashboard
            agent.act("1")
            agent.wait_for_text("DASHBOARD")
            agent.wait_until_stable()
            initial_text = agent.observe()
            initial_active_line = next(
                (ln for ln in initial_text.split("\n") if "Active Deck:" in ln),
                "",
            )

            # Switch deck via Review tab
            agent.act("3")
            agent.wait_until_stable()
            agent.act("]")
            agent.wait_until_stable()
            agent.act("]")
            agent.wait_until_stable()

            # Verify status reflects deck change
            agent.wait_for_text("in ", timeout=2.0)

            # Confirm deck changed on Dashboard
            agent.act("1")
            agent.wait_for_text("DASHBOARD")
            agent.wait_until_stable()
            after_text = agent.observe()
            after_active_line = next(
                (ln for ln in after_text.split("\n") if "Active Deck:" in ln),
                "",
            )
            assert initial_active_line != after_active_line, (
                f"Deck should have changed: before={initial_active_line!r}, after={after_active_line!r}"
            )
        finally:
            agent.close()


def test_review_status_includes_deck_name():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed(agent)
            agent.act("3")
            agent.wait_until_stable()
            # Status line should contain "in <deck name>" not just "All caught up!"
            agent.wait_for_text("status:", timeout=2.0)
            agent.wait_for_text(" in ", timeout=2.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_review_empty_state_shows_current_deck()
    test_review_deck_switch_updates_label()
    test_review_status_includes_deck_name()
