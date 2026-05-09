import sys
import os
import tempfile
import pytest

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', './deutsch-tui-bin')
    # If we're running from tui_tester subdir, adjust path
    if not os.path.exists(app_cmd) and os.path.exists('../deutsch-tui-bin'):
        app_cmd = '../deutsch-tui-bin'
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_cloze_typing_mode():
    """Test that typing mode works correctly for Cloze deletion cards."""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a TSV with two Cloze cards
        tsv_path = os.path.join(tmpdir, "cloze_typing.tsv")
        with open(tsv_path, "w") as f:
            f.write("#separator:tab\n#deck:Cloze Typing\n")
            f.write("id-cloze-1\tDas Wetter ist heute {{c1::schön::beautiful}}.\tnice\t\t\t\t\n")
            f.write("id-cloze-2\tIch {{c1::bin::to be}} müde.\tam\t\t\t\t\n")
        
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act('5')
            agent.wait_for_text("Import TSV")
            
            # Type the path
            agent.act('<Enter>') # Edit path
            for _ in range(100):
                agent.act('<Backspace>')
            agent.act(tsv_path)
            agent.act('<Enter>') # Save path
            
            # Import
            agent.act('i')
            agent.wait_for_text("Imported 2 notes")
            
            # Switch to Dashboard
            agent.act('1')
            agent.wait_for_text("DASHBOARD")
            
            # Switch decks until "Cloze Typing" is active
            found_deck = False
            for _ in range(15):
                screen = agent.observe()
                if "Active Deck: Cloze Typing" in screen:
                    found_deck = True
                    break
                agent.act(']')
                agent.wait_until_stable()
            
            assert found_deck, "Could not find 'Cloze Typing' deck"
            
            # Go to Review
            agent.act('3')
            agent.wait_until_stable()
            
            # Enable Typing Mode
            agent.act('t')
            agent.wait_for_text("Type your answer:")
            
            # Type CORRECT answer for first card
            # First card prompt: Das Wetter ist heute [beautiful].
            # Answer: schön
            agent.act("schön")
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.wait_for_text("✓ schön")
            
            # Grade Good
            agent.act('g')
            agent.wait_until_stable()
            
            # Now on second card: Ich [to be] müde.
            agent.wait_for_text("Ich [to be] müde.")
            
            # Type WRONG answer for second card
            agent.act("bist")
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.wait_for_text("✗ bist")
            agent.wait_for_text("Correct: bin")
            
            # Grade Easy (just to finish)
            agent.act('e')
            agent.wait_until_stable()
            
            # Should be finished
            agent.wait_for_text("No cards due")
            
        finally:
            agent.close()
