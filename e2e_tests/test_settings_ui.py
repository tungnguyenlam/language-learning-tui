import time
import os
import sys
import tempfile
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

@pytest.fixture
def agent():
    tmpdir = tempfile.mkdtemp()
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=110, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    yield agent
    agent.close()

def test_settings_ui_layout(agent):
    agent.act("7") # Navigate to Settings
    agent.wait_for_text("AI CONFIGURATION")
    agent.wait_for_text("Template Set: vocabulary")
    agent.wait_for_text("AI Provider:    disabled")
    agent.wait_for_text("Front Template:")
    agent.wait_for_text("Back Template: ")
    agent.wait_for_text("Example Tmpl:   ")
    agent.wait_for_text("STUDY PREFERENCES")
    agent.wait_for_text("Daily Goal:")
