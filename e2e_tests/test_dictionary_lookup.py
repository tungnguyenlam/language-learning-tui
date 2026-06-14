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
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
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
            
            # It should open the Spotlight Dictionary overlay
            agent.wait_for_text("SPOTLIGHT DICTIONARY", timeout=10.0)
            
            # The search should be pre-filled with the word from Review
            time.sleep(2.0) # Wait for search to complete
            screen = agent.observe()
            assert "SPOTLIGHT DICTIONARY" in screen
            # Check for the search bar UI we added
            assert "🔍" in screen

            # Close it with Esc
            agent.act("<Esc>")
            agent.wait_for_text("REVIEW", timeout=5.0)
            assert "SPOTLIGHT DICTIONARY" not in agent.observe()
        finally:
            agent.close()

def test_dictionary_single_column_details():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Pre-seed the database using -smoke flag
        app_path = "./cmd/deutsch-tui/main.go"
        subprocess.run(["go", "run", app_path, "-data-dir", tmpdir, "-smoke"], check=True)
        
        # Start agent with narrow terminal columns=70 to force single-column layout
        app_cmd = os.getenv('DEUTSCH_TUI_BIN', f'go run {app_path}')
        agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=70, lines=30)
        try:
            agent.wait_for_text("DASHBOARD", timeout=15.0)
            agent.wait_until_stable()
            
            # Go to Dictionary view
            agent.act("/")
            agent.wait_for_text("Dictionary", timeout=10.0)
            agent.wait_until_stable()
            
            # Type 'Kaffee' to search
            agent.act("Kaffee")
            time.sleep(2.0)
            
            screen = agent.observe()
            assert "Kaffee" in screen
            # Should have the single column hint at the bottom
            assert "Press ctrl+d/click selected to view details" in screen
            
            # Press ctrl+d to toggle detail view
            agent.act("<C-d>")
            time.sleep(1.0)
            
            screen_details = agent.observe()
            assert "Translations:" in screen_details
            assert "Press esc/ctrl+d to return to list" in screen_details
            
            # Press esc to exit details
            agent.act("<Esc>")
            time.sleep(1.0)
            
            screen_list = agent.observe()
            assert "Press ctrl+d/click selected to view details" in screen_list
        finally:
            agent.close()
