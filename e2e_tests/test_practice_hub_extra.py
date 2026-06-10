import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=110, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_separable_verb_trainer():
    """Test navigating to and interacting with Separable Verb Trainer."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Practice Hub (key 0)
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)
            
            # Select Separable Verb Trainer (key 7)
            agent.act('7')
            agent.wait_for_text("SEPARABLE VERB TRAINER", timeout=5.0)
            
            # Should show Score: 0/0 and Enter the missing prefix:
            screen = agent.observe()
            assert "Score: 0/0" in screen
            assert "Enter the missing prefix:" in screen
            
            # Type something correct ("auf") and press enter
            agent.act('a')
            agent.act('u')
            agent.act('f')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT!" in screen
            assert "Score: 1/1" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should show Score: 1/1 and back to enter prefix
            screen = agent.observe()
            assert "Score: 1/1" in screen
            
            # Test Esc to go back to Hub
            agent.act('<Esc>')
            agent.wait_for_text("PRACTICE HUB")
            
        finally:
            agent.close()

def test_plural_trainer():
    """Test navigating to and interacting with Noun Plural Trainer."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Practice Hub (key 0)
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)
            
            # Select Plural Trainer (key 6)
            agent.act('6')
            agent.wait_for_text("NOUN PLURAL TRAINER", timeout=5.0)
            
            # Should show singular noun and score 0/0
            screen = agent.observe()
            assert "Score: 0/0" in screen
            
            # Type something correct and press enter. Since the starter deck is seeded,
            # the first noun with a plural is "der Apfel", which has plural "die Äpfel".
            # We can type "aepfel" (a-e-p-f-e-l).
            agent.act('a')
            agent.act('e')
            agent.act('p')
            agent.act('f')
            agent.act('e')
            agent.act('l')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT!" in screen
            assert "Score: 1/1" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should show score 1/1
            screen = agent.observe()
            assert "Score: 1/1" in screen
            
            # Test Esc to go back to Hub
            agent.act('<Esc>')
            agent.wait_for_text("PRACTICE HUB")
            
        finally:
            agent.close()

def test_conjunctions_trainer():
    """Test navigating to and interacting with Conjunctions & Word Order Trainer."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Practice Hub (key 0)
            agent.act('0')
            agent.wait_for_text("PRACTICE HUB", timeout=5.0)
            
            # Select Conjunctions Trainer (key 9)
            agent.act('9')
            agent.wait_for_text("CONJUNCTIONS & WORD ORDER", timeout=5.0)
            
            # Should show Score: 0/0 and Enter the missing word:
            screen = agent.observe()
            assert "Score: 0/0" in screen
            assert "Enter the missing word:" in screen
            
            # Type something correct. The first exercise is:
            # "Ich bleibe heute zu Hause, {{...}} es regnet." -> Answer: "weil"
            agent.act('w')
            agent.act('e')
            agent.act('i')
            agent.act('l')
            agent.act('<Enter>')
            agent.wait_until_stable()
            
            # Should show revealed state
            screen = agent.observe()
            assert "CORRECT!" in screen
            assert "Score: 1/1" in screen
            assert "weil (because) is a subordinating conjunction" in screen
            
            # Press any key for next
            agent.act(' ')
            agent.wait_until_stable()
            
            # Should show Score: 1/1
            screen = agent.observe()
            assert "Score: 1/1" in screen
            
            # Test Esc to go back to Hub
            agent.act('<Esc>')
            agent.wait_for_text("PRACTICE HUB")
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
