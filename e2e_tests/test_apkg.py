import os
import sys
import tempfile
import time
import zipfile
import sqlite3

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=110, lines=30)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_apkg_export_import_full_cycle():
    with tempfile.TemporaryDirectory() as tmpdir:
        # 1. Create a TSV to have some data
        tsv_path = os.path.join(tmpdir, "test_import.tsv")
        with open(tsv_path, "w") as f:
            # Use a unique deck name to avoid confusion
            f.write("id-apkg-1\tfront-1\tback-1\textra-1\ttags-1\tdeck-apkg-test\tBasic\n")
        
        agent = start_agent(tmpdir)
        try:
            # 2. Import TSV
            agent.act("5") # Import view
            agent.wait_for_text("Import / Export")
            
            # Select import field
            agent.act("k") # Ensure import field is selected
            agent.act("<Enter>") # START EDITING
            agent.act("<Ctrl-u>")
            for char in tsv_path:
                agent.act(char)
            agent.act("<Enter>") # FINISH EDITING
            time.sleep(0.5)
            agent.wait_until_stable()
            agent.act("i") # EXECUTE IMPORT TSV
            agent.wait_for_text("Imported", timeout=10.0)
            
            # 3. Set Export path
            export_path = os.path.join(tmpdir, "test_export.apkg")
            agent.act("j") # Select export field
            agent.act("<Enter>") # START EDITING
            agent.act("<Ctrl-u>")
            for char in export_path:
                agent.act(char)
            agent.act("<Enter>") # FINISH EDITING
            time.sleep(0.5)
            agent.wait_until_stable()
            
            # 4. Export as APKG
            agent.act("X") # Shift+X for APKG
            agent.wait_for_text("Exported", timeout=10.0)
            
            # Verify file exists
            assert os.path.exists(export_path)
            
            # 5. Clear database by restarting with a new data dir or just import back
            # For simplicity, let's just import it back and check the message
        
            # Select import field again
            agent.act("k")
            agent.act("<Enter>") # START EDITING
            agent.act("<Ctrl-u>")
            for char in export_path:
                agent.act(char)
            agent.act("<Enter>") # FINISH EDITING
            time.sleep(0.5)
            agent.wait_until_stable()
            agent.act("I") # Shift+I for APKG import
            agent.wait_for_text("Imported", timeout=10.0)
            
            # 6. Verify in Decks view
            agent.act("2")
            agent.wait_for_text("Decks")
            agent.wait_for_text("deck-apkg-test")
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_apkg_export_import_full_cycle()
