import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_bulk_deck_deletion():
    """Verify that multiple decks can be selected and deleted at once."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create two extra decks via import.tsv
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, 'w') as f:
            f.write("#separator:tab\n")
            f.write("t1\tT1\tB1\t\t\tTrash1\n")
            f.write("t2\tT2\tB2\t\t\tTrash2\n")

        agent = start_agent(tmpdir)
        try:
            # Import
            agent.act('5')
            agent.act('i')
            agent.wait_for_text("Imported")
            
            # Go to Decks view
            agent.act('2')
            agent.wait_for_text("Decks")
            
            # Select Trash1 using search
            agent.act('/')
            agent.act('T')
            agent.act('r')
            agent.act('a')
            agent.act('s')
            agent.act('h')
            agent.act('1')
            agent.act('<Esc>')
            time.sleep(1.0)
            agent.wait_for_text("multi-select")
            agent.wait_until_stable()
            
            agent.act('x') # Select Trash1
            agent.wait_for_text("[x] Trash1")
            
            # Select Trash2 using search
            agent.act('/')
            agent.act('T')
            agent.act('r')
            agent.act('a')
            agent.act('s')
            agent.act('h')
            agent.act('2')
            agent.act('<Esc>')
            time.sleep(1.0)
            agent.wait_for_text("decks selected") # Text changes when one is selected
            agent.wait_until_stable()
            
            agent.act('x') # Select Trash2
            agent.wait_for_text("[x] Trash2")
            
            # Delete both
            agent.act('<Backspace>')
            agent.wait_until_stable()
            
            # Verify they are gone (need to clear filter first or wait for list refresh)
            agent.act('<Esc>')
            agent.wait_until_stable()
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
            agent.wait_for_text("[x] Source1")
            
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
            agent.assert_text("> [ ] Target")
            
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
