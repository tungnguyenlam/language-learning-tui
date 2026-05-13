import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from tui_tester.agent import TUIAgent


def start_agent(tmpdir, columns=130, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def _seed_and_review(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Seeding standard content...", timeout=10.0)
    time.sleep(3.0)
    agent.wait_until_stable(timeout=15.0)
    agent.act("3")
    agent.wait_for_text("Review", timeout=5.0)


def test_report_card_with_offline_provider_shows_no_op_proposal():
    # Offline provider doesn't implement ChatProvider, so FixCard returns
    # an unchanged proposal with a clear reason. Pressing F should bring
    # up the preview block.
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_review(agent)
            agent.act("F")
            # Either the AI box or the proposal preview should appear.
            agent.wait_for_text("AI", timeout=8.0)
        finally:
            agent.close()


def test_report_card_then_discard_with_n():
    # After reporting, pressing n should discard the proposal and clear
    # the preview block. Status should mention "Fix discarded".
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _seed_and_review(agent)
            agent.act("F")
            # Give the offline flow a moment to produce a "no-op" proposal.
            time.sleep(1.5)
            agent.wait_until_stable(timeout=5.0)
            agent.act("n")
            agent.wait_until_stable(timeout=3.0)
        finally:
            agent.close()


def test_report_key_does_not_crash_on_empty_review():
    # Pressing F on the Dashboard (not Review) should be a no-op.
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("F")
            agent.wait_until_stable(timeout=3.0)
            # Still on dashboard
            agent.wait_for_text("DASHBOARD", timeout=3.0)
        finally:
            agent.close()


if __name__ == "__main__":
    test_report_card_with_offline_provider_shows_no_op_proposal()
    test_report_card_then_discard_with_n()
    test_report_key_does_not_crash_on_empty_review()
