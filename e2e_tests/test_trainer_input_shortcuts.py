import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_trainer_typing_does_not_trigger_global_shortcuts():
    """'q' and '?' must be typed into a trainer answer, not quit / open help.

    German answers contain 'q' (Qualität, Quelle); before this was fixed,
    typing one exited the app mid-exercise.
    """
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)

            # Plural Trainer (key 6) takes free-text input.
            agent.act('6')
            agent.wait_for_text("NOUN PLURAL TRAINER", timeout=5.0)

            agent.act('q')
            agent.wait_until_stable()

            screen = agent.observe()
            # Still in the trainer: the app did not quit and help did not open.
            assert "NOUN PLURAL TRAINER" in screen
            assert "Keyboard Shortcuts" not in screen

            agent.act('u')
            agent.act('?')
            agent.wait_until_stable()

            screen = agent.observe()
            assert "NOUN PLURAL TRAINER" in screen
            assert "Keyboard Shortcuts" not in screen

            # Esc clears the typed answer, a second Esc leaves the trainer.
            agent.act('<Esc>')
            agent.wait_until_stable()
            agent.act('<Esc>')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)

            # Outside a trainer, '?' still opens the help overlay.
            agent.act('?')
            agent.wait_for_text("Keyboard Shortcuts", timeout=5.0)

        finally:
            agent.close()


def test_trainer_shows_item_position():
    """The trainer header reports how far through the exercise set you are."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)

            agent.act('3')
            agent.wait_for_text("CASE ENDING TRAINER", timeout=5.0)

            assert "Item 1/" in agent.observe()

            # Answer, then advance: the position moves on.
            agent.act('d')
            agent.act('e')
            agent.act('m')
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.act(' ')
            agent.wait_until_stable()

            assert "Item 2/" in agent.observe()

        finally:
            agent.close()


if __name__ == "__main__":
    pytest.main(["-v", __file__])
