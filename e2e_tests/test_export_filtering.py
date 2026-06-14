import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_export_filtering_by_tag():
    """Verify that export only includes notes with the specified tag."""
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        export_path = os.path.join(tmpdir, "export.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("#deck:TagFilterDeck\n")
            f.write("n1\tF1\tB1\t\ttagX\n")
            f.write("n2\tF2\tB2\t\ttagY\n")

        agent = start_agent(tmpdir)
        try:
            agent.act('5')
            agent.wait_for_text("Import / Export")
            agent.act('i')
            agent.wait_for_text("Imported")
            
            # Set export path
            agent.act('j')
            agent.act('<Enter>')
            agent.act('<C-u>')
            agent.act(export_path)
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Set export tag to 'tagX'
            agent.act('j') # Deck Filter
            agent.act('j') # Tag Filter
            agent.act('t')
            agent.act('tagX')
            agent.act('<Enter>')
            agent.wait_for_text("tagX")
            
            # Export
            agent.act('x')
            agent.wait_for_text("Exported 1 notes")
            
            with open(export_path, 'r') as f:
                content = f.read()
                assert "F1" in content
                assert "F2" not in content
        finally:
            agent.close()

def test_export_filtering_by_deck():
    """Verify that export only includes notes from the specified deck."""
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        export_path = os.path.join(tmpdir, "export.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("#deck:UniqueDeckName\n")
            f.write("u1\tFrontU\tBackU\n")

        agent = start_agent(tmpdir)
        try:
            agent.act('5')
            agent.act('i')
            agent.wait_for_text("Imported")
            
            agent.act('j') # Export Path
            agent.act('<Enter>')
            agent.act('<C-u>')
            agent.act(export_path)
            agent.act('<Enter>')
            
            agent.act('j') # Deck filter
            agent.wait_until_stable()
            
            found = False
            for _ in range(5):
                if "UniqueDeckName" in agent.observe():
                    found = True
                    break
                agent.act(']')
                agent.wait_until_stable()
            
            assert found, "Could not find UniqueDeckName in export filter"
            
            agent.act('x')
            agent.wait_for_text("Exported 1 notes")
            
            with open(export_path, 'r') as f:
                content = f.read()
                assert "FrontU" in content
        finally:
            agent.close()
