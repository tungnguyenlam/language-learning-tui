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
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_bulk_deck_deletion():
    """Verify that multiple decks can be selected and deleted at once."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create two extra decks via import.tsv
        # Format: id\tfront\tback\textra\ttags\tdeck
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("t1\tT1\tB1\t\t\tTrash1\n")
            f.write("t2\tT2\tB2\t\t\tTrash2\n")

        agent = start_agent(tmpdir, columns=110, lines=40)
        try:
            # Import
            agent.act('5')
            agent.wait_for_text("Import / Export")
            agent.act('i')
            agent.wait_for_text("Imported", timeout=10.0)
            agent.wait_until_stable()
            time.sleep(1.0)
            
            # Go to Decks view
            agent.act('2')
            agent.wait_for_text("DECK LIST", timeout=10.0)
            # Clear any existing search filter
            agent.act('<Esc>')
            time.sleep(0.5)
            agent.wait_until_stable()
            
            # Verify Trash decks exist
            agent.wait_for_text("Trash1", timeout=5.0)
            agent.wait_for_text("Trash2", timeout=5.0)
            
            # Select Trash1 using search
            agent.act('/')
            for ch in "Trash1":
                agent.act(ch)
            agent.act('<Enter>')
            time.sleep(0.5)
            agent.wait_until_stable()
            
            agent.act('x') # Select Trash1
            agent.wait_for_text("[x] Trash1")
            
            # Select Trash2 using search
            agent.act('/')
            for ch in "Trash2":
                agent.act(ch)
            agent.act('<Enter>')
            time.sleep(0.5)
            agent.wait_until_stable()
            
            agent.act('x') # Select Trash2
            agent.wait_for_text("2 decks selected")
            
            # Delete both
            agent.act('<Backspace>')
            agent.wait_for_text("CONFIRM DELETION")
            agent.act('y')
            agent.wait_until_stable()
            time.sleep(0.5)
            
            # Verify they are gone
            agent.act('<Esc>')
            agent.wait_until_stable()
            time.sleep(0.5)
            screen = agent.observe()
            assert "Trash1" not in screen
            assert "Trash2" not in screen
        finally:
            agent.close()

def test_deck_merging():
    """Verify that selected decks can be merged into another deck."""
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("s1\tS1\tB1\t\t\tSource1\n")
            f.write("t1\tT1\tB1\t\t\tTarget\n")

        agent = start_agent(tmpdir)
        try:
            agent.act('5')
            agent.act('i')
            agent.wait_for_text("Imported")
            
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Select Source1 using search
            agent.act('/')
            agent.act('S')
            agent.act('o')
            agent.act('u')
            agent.act('r')
            agent.act('c')
            agent.act('e')
            agent.act('1')
            agent.act('<Esc>')
            time.sleep(1.0)
            agent.wait_for_text("multi-select")
            agent.wait_until_stable()
            
            agent.act('x')
            agent.wait_for_text("1 decks selected")
            
            # Highlight Target using search
            agent.act('/')
            agent.act('T')
            agent.act('a')
            agent.act('r')
            agent.act('g')
            agent.act('e')
            agent.act('t')
            agent.act('<Esc>')
            time.sleep(1.0)
            agent.wait_for_text("decks selected")
            agent.wait_until_stable()
            agent.assert_text("Target")

            # Merge M (Source1 -> Target)
            agent.act('M')
            agent.wait_until_stable()

            # Verify Source1 is gone
            agent.act('<Esc>')
            agent.wait_until_stable()
            assert "Source1" not in agent.observe()            
            # Go to Dashboard and verify card counts if possible, 
            # or just verify Target is still there
            agent.assert_text("Target")
        finally:
            agent.close()
