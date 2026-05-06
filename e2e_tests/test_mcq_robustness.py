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

def test_mcq_with_commas_in_choices():
    with tempfile.TemporaryDirectory() as tmpdir:
        tsv_path = os.path.join(tmpdir, "mcq_commas.tsv")
        with open(tsv_path, "w") as f:
            # Simple unique string to verify
            f.write("id-1\tprompt\tanswer\textra\ttags\tdeck-1\tMCQ:C1, comma|||C2\n")
        
        agent = start_agent(tmpdir)
        try:
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.act("<Ctrl-u>")
            for char in tsv_path:
                agent.act(char)
            agent.act("<Enter>")
            agent.act("i")
            time.sleep(1.0)
            
            agent.act("2") # Decks
            agent.wait_for_text("deck-1")
            agent.act("j")
            agent.act("<Enter>")
            
            agent.act("3") # Review
            agent.wait_for_text("Review 1/2")
            agent.act(" ") # Reveal
            agent.wait_for_text("answer")
            agent.act("g") # Grade
            agent.wait_for_text("Review 1/1")
            
            agent.act(" ") # Reveal MCQ
            agent.wait_for_text("C1, comma")
            agent.wait_for_text("C2")
        finally:
            agent.close()
