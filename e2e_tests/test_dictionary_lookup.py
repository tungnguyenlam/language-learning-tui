import os
import sys
import tempfile
import time
import subprocess

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    # Pre-seed the database using -smoke flag
    app_path = "./cmd/deutsch-tui/main.go"
    subprocess.run(["go", "run", app_path, "-data-dir", tmpdir, "-smoke"], check=True)
    
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', f'go run {app_path}')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dictionary_lookup_no_crash():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Standard content should be pre-seeded now
            
            # Go to Review mode
            agent.act("3") # Review is 3 based on updateNumberKey
            agent.wait_for_text("REVIEW", timeout=15.0)
            agent.wait_until_stable()
            
            # Press 'd' to trigger dictionary lookup
            agent.act("d")
            
            # It should switch to Dictionary view
            agent.wait_for_text("Dictionary", timeout=10.0)
            
            # The search should be pre-filled with the word from Review
            time.sleep(2.0) # Wait for search to complete
            screen = agent.observe()
            assert "Dictionary" in screen
            # Check for the search bar UI we added
            assert "🔍" in screen
        finally:
            agent.close()
