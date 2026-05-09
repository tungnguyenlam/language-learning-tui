import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ../cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dictionary_lookup_no_crash():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Review mode
            agent.act("2")
            agent.wait_for_text("Review")
            agent.wait_until_stable()
            
            # Press 'd' to trigger dictionary lookup
            agent.act("d")
            
            # Wait a moment to ensure it didn't crash
            time.sleep(1)
            agent.wait_for_text("Review")
            
            # Verify app is still running and hasn't panicked
            screen = agent.observe()
            assert "panic:" not in screen
            assert "Review" in screen
        finally:
            agent.close()
