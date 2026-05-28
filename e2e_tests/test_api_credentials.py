import json
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


def _go_to_settings(agent):
    agent.act("7")
    agent.wait_for_text("SETTINGS")
    agent.wait_until_stable()


def test_settings_has_api_credentials_section():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _go_to_settings(agent)
            agent.wait_for_text("API CREDENTIALS", timeout=3.0)
            agent.wait_for_text("API Key:", timeout=3.0)
            agent.wait_for_text("Model:", timeout=3.0)
            agent.wait_for_text("Base URL:", timeout=3.0)
        finally:
            agent.close()


def test_provider_cycle_includes_openai_and_anthropic():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _go_to_settings(agent)
            # Cycle: disabled -> offline -> template -> openai -> anthropic
            expected = ["offline", "template", "openai", "anthropic"]
            for want in expected:
                agent.act("<Enter>")
                time.sleep(0.25)
                agent.wait_until_stable()
                agent.wait_for_text(f"AI Provider:    {want}", timeout=2.0)
        finally:
            agent.close()


def test_typing_api_key_persists_to_secrets_json():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            _go_to_settings(agent)
            # Cycle to openai (disabled -> offline -> template -> openai)
            for _ in range(3):
                agent.act("<Enter>")
                time.sleep(0.2)
                agent.wait_until_stable()
            agent.wait_for_text("AI Provider:    openai", timeout=2.0)

            # Move down to the API Key row (cursor 8) — 8 j keypresses from cursor 0
            for _ in range(8):
                agent.act("j")
                time.sleep(0.05)
            agent.wait_until_stable()
            agent.act("<Enter>")
            agent.wait_until_stable()

            # Type the API key
            agent.act("sk-test-abc123")
            time.sleep(0.3)
            agent.act("<Enter>")
            time.sleep(0.5)
            agent.wait_until_stable()

            # Verify it was masked and the file was saved at 0600
            secrets_path = os.path.join(tmpdir, "secrets.json")
            assert os.path.exists(secrets_path), "secrets.json should be created"
            stat = os.stat(secrets_path)
            assert (stat.st_mode & 0o777) == 0o600, f"mode = {oct(stat.st_mode & 0o777)}"
            with open(secrets_path) as f:
                data = json.load(f)
            try:
                assert data["openai"]["api_key"] == "sk-test-abc123", data
            except KeyError as e:
                from tui_tester.exceptions import TUIAssertionError
                raise TUIAssertionError(f"KeyError: {e} in {data}", agent.observe())

            # Screen should mask it
            text = agent.observe()
            assert "sk-test-abc123" not in text, "API key should not appear in plain text on screen"
            assert "****" in text, "Masked key should appear on screen"
        finally:
            agent.close()
