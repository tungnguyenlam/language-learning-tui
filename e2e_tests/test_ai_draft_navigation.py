import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_ai_draft_navigation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Navigate to Settings and set offline mode
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.act("<Enter>")
            agent.wait_for_text("offline")
            
            # Navigate to AI view
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            
            # Enter edit mode to type topic
            agent.act("/")
            agent.act("<Ctrl-u>")
            
            # Type a simple topic that will generate a single draft
            topic = "Apfel"
            for char in topic:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            
            # Wait for draft to be generated
            agent.wait_for_text("Apfel", timeout=10.0)
            agent.assert_text("Approve")
            
            # Approve the draft with 'A'
            agent.act("A")
            agent.wait_for_text("No drafts", timeout=5.0)
            
            # Now generate another draft with different topic
            agent.act("/")
            agent.act("<Ctrl-u>")
            for _ in range(10):
                agent.act("<Backspace>")
            
            new_topic = "Birne"
            for char in new_topic:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            
            # Wait for new draft
            agent.wait_for_text("Birne", timeout=10.0)
            
            # Approve all with 'A'
            agent.act("A")
            agent.wait_for_text("No drafts", timeout=5.0)
            
            # Navigate to Browser to verify cards were created
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Should be able to find our cards
            agent.act("/")
            agent.act("Apfel")
            agent.wait_for_text("Apfel")
            
        finally:
            agent.close()

def test_ai_draft_approval():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Navigate to Settings and set offline mode
            agent.act("7")
            agent.wait_for_text("Settings")
            agent.act("<Enter>")
            agent.wait_for_text("offline")
            
            # Navigate to AI view
            agent.act("6")
            agent.wait_for_text("AI Drafts")
            
            # Enter edit mode
            agent.act("/")
            agent.act("<Ctrl-u>")
            
            # Type topic 
            topic = "Tisch"
            for char in topic:
                agent.act(char)
                time.sleep(0.01)
            agent.act("<Enter>")
            
            # Should generate a draft
            agent.wait_for_text("Tisch", timeout=10.0)
            
            # Use 'a' to approve single draft
            agent.act("a")
            agent.wait_for_text("No drafts", timeout=5.0)
            
            # Verify by going to browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
        finally:
            agent.close()