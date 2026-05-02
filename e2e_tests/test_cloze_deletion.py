import sys
import os
import tempfile
import pytest

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_cloze_deletion_flow():
    """Test that Cloze deletion cards work correctly in the TUI."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a TSV with a Cloze card
        tsv_path = os.path.join(tmpdir, "cloze_test.tsv")
        with open(tsv_path, "w") as f:
            f.write("#separator:tab\n#deck:Cloze Test\n")
            f.write("id-1\tIch {{c1::gebe::to give}} dem Mann das Buch.\tgive\t\t\t\t\n")
        
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act('5')
            agent.wait_for_text("Import TSV")
            
            # Type the path
            agent.act('<Enter>') # Edit path
            # Clear existing path (about 100 characters should be enough)
            for _ in range(100):
                agent.act('<Backspace>')
            agent.act(tsv_path)
            agent.act('<Enter>') # Save path
            
            # Import
            agent.act('i')
            agent.wait_for_text("Imported 1 notes")
            
            # Switch to Dashboard
            agent.act('1')
            agent.wait_for_text("DASHBOARD")
            
            # Switch decks until "Cloze Test" is active
            # We'll try up to 10 times
            found_deck = False
            for _ in range(10):
                if "Active Deck: Cloze Test" in agent.observe():
                    found_deck = True
                    break
                agent.act(']')
                agent.wait_until_stable()
            
            if not found_deck:
                raise Exception("Could not find 'Cloze Test' deck")
            
            # Go to Review
            agent.act('3')
            agent.wait_until_stable()
            
            # Verify Cloze prompt
            # It should show [to give] instead of gebe
            # Note: lipgloss styling might be present, but wait_for_text handles plain text
            agent.wait_for_text("Ich [to give] dem Mann das Buch.")
            
            # Reveal
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Verify full sentence is shown
            agent.wait_for_text("Ich gebe dem Mann das Buch.")
            agent.assert_text("Grade: a Again")
            
            # Grade Easy
            agent.act('e')
            agent.wait_until_stable()
            
            # Should be finished
            agent.wait_for_text("No cards due")
            
        finally:
            agent.close()
