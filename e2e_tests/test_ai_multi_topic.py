import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_ai_multi_topic_generation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.act("<Enter>")
            agent.wait_for_text("offline")
            agent.act("6") # AI view
            agent.wait_for_text("Topic:")
            
            # Enter edit mode
            agent.act("/")
            
            # Clear default input
            agent.act("<Ctrl-u>")
            
            # Type multiple topics
            topics = "Apple, Banana, Orange"
            for char in topics:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            
            # Should generate multiple drafts
            agent.wait_for_text("Apple ->", timeout=10.0)
            agent.wait_for_text("Banana ->")
            agent.wait_for_text("Orange ->")
            
            # Approve all with 'A'
            agent.act("A")
            agent.wait_for_text("No drafts yet.", timeout=5.0)
            
            # Verify they are in the browser
            agent.act("8") # Browser view
            agent.wait_for_text("Card Browser")
            
            # Search for Apple
            agent.act("/")
            agent.act("Apple")
            agent.wait_for_text("Apple")
            
            # Search for Banana
            agent.act("<Esc>") # Clear search
            agent.act("/")
            agent.act("Banana")
            agent.wait_for_text("Banana")
            
        finally:
            agent.close()
