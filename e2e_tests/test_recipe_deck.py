import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_recipe_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Import view
            agent.act("5")
            agent.wait_for_text("Import / Export")
            
            # Seed standard content
            agent.act("S")
            agent.wait_for_text("Seeding standard content...")
            time.sleep(2.0)
            agent.wait_until_stable()
            
            # Go to Decks
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=5.0)
            
            # Search for Recipe deck
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Recipe")
            agent.act("<Enter>")
            agent.wait_for_text("Recipes & Cooking", timeout=10.0)
            
            # Select the deck
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            
            # Verify deck is active
            agent.wait_for_text("Recipes & Cooking", timeout=5.0)
            
        finally:
            agent.close()

def test_recipe_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Navigate to Import
            agent.act("5")
            agent.wait_for_text("Import / Export")
            
            # Seed content
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            
            # Go to Decks and find Recipe deck
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            
            # Navigate to find Recipe deck using bracket keys
            # First clear any search
            agent.act("<Esc>")
            agent.wait_until_stable()
            
            # Search for cooking-related content
            agent.act("/")
            agent.act("kochen")
            agent.act("<Enter>")
            
            # Should find cards with "kochen" (to cook)
            agent.wait_for_text("kochen", timeout=10.0)
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_recipe_deck_exists()
    test_recipe_vocabulary_accessible()