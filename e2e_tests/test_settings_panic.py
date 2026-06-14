import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def test_settings_panic():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create empty config to force empty template sets?
        # Actually just start agent
        app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ../cmd/deutsch-tui')
        agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
        agent.wait_for_text("DASHBOARD", timeout=15.0)
        agent.wait_until_stable()
        
        try:
            # Go to Settings
            agent.act("7")
            agent.wait_for_text("Settings")
            # Wait to see if it panics
            time.sleep(1)
            screen = agent.observe()
            print("Is panic?", "panic" in screen.lower())
        finally:
            agent.close()

if __name__ == "__main__":
    test_settings_panic()
