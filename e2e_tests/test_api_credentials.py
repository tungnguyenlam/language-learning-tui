import json
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
            agent.wait_for_text("OpenAI Key:", timeout=3.0)
            agent.wait_for_text("OpenAI Model:", timeout=3.0)
            agent.wait_for_text("Anthropic Key:", timeout=3.0)
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
            
            # Now we can edit keys directly without cycling!
            # OpenAI Key is at row 7
            for _ in range(7):
                agent.act("j")
            agent.wait_until_stable()
            agent.act("<Enter>")
            agent.wait_until_stable()

            # Type the API key
            agent.act("sk-test-abc123")
            time.sleep(0.3)
            agent.act("<Enter>")
            time.sleep(0.5)
            agent.wait_until_stable()

            # Verify it was masked
            secrets_path = os.path.join(tmpdir, "secrets.json")
            assert os.path.exists(secrets_path), "secrets.json should be created"
            with open(secrets_path) as f:
                data = json.load(f)
            assert data["openai"]["api_key"] == "sk-test-abc123", data

            # Entering the key should have automatically enabled 'openai'
            agent.wait_for_text("AI Provider:    openai", timeout=2.0)
        finally:
            agent.close()
