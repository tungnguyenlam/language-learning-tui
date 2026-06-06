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
    
    # Insert Katze into the dictionary directly
    db_path = os.path.join(tmpdir, "learning.db")
    import sqlite3
    conn = sqlite3.connect(db_path)
    conn.execute("INSERT INTO dictionary_fts (id, word, translation, word_class, gender, forms, examples, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", ("1", "Katze", "cat", "", "f", "", "", "[\"animal\"]"))
    conn.commit()
    conn.close()
    
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', f'go run {app_path}')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}')
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dictionary_search_results():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Dictionary view
            agent.act("/")
            agent.wait_for_text("Dictionary", timeout=10.0)
            
            # Type "Katze"
            for char in "Katze":
                agent.act(char)
            
            time.sleep(3.0) # Wait for search
            screen = agent.observe()
            
            # Check for results
            assert "Katze" in screen
            # Standard translation for Katze is "cat"
            assert "cat" in screen or "Cat" in screen
            
        finally:
            agent.close()
