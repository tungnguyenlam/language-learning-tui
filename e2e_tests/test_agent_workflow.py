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
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    yield agent
    agent.close()

def test_decks_view_loads_starter_deck(agent):
    agent.act("2") # 2 is decks
    agent.wait_for_text("German A1 Survival")
    
def test_ai_placeholder(agent):
    # Navigate to AI View
    agent.act("6") # 6 is AI view
        
    agent.wait_for_text("Topic:")
    agent.wait_for_text("A1 survival")
