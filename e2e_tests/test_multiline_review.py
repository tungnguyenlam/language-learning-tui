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

def test_multiline_prompt_rendering():
    with tempfile.TemporaryDirectory() as tmpdir:
        tsv_path = os.path.join(tmpdir, "multi.tsv")
        long_prompt = "Line 1\\nLine 2\\nLine 3" # Embedded newlines
        with open(tsv_path, "w") as f:
            f.write(f"id-1\t{long_prompt}\tanswer\t\t\tdeck-1\n")
            
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.act("<Enter>")
            agent.act("<Ctrl-u>")
            for char in tsv_path:
                agent.act(char)
            agent.act("<Enter>")
            agent.act("i")
            time.sleep(1.0)
            
            agent.act("2")
            agent.wait_for_text("deck-1")
            agent.act("j")
            agent.act("<Enter>")
            
            agent.act("3")
            agent.wait_for_text("Review 1/1")
            
            # Verify multi-line prompt
            agent.wait_for_text("Line 1")
            agent.wait_for_text("Line 2")
            agent.wait_for_text("Line 3")
        finally:
            agent.close()
