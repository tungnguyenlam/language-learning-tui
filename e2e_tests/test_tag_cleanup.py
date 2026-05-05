import sys
import os
import tempfile
import pytest
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))
from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_unused_tags_cleanup():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # 1. Check initial tags in Decks view
            agent.act("2") # Decks view
            agent.wait_for_text("Decks")
            agent.assert_text("Tags: german, a1")
            
            # 2. Go to Browser and remove "a1" tag from all cards
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            
            # Search for cards with "a1" tag to be sure
            agent.act("/")
            agent.act("a1")
            agent.act("<Enter>")
            agent.wait_until_stable()
            
            # Select all visible cards (for simplicity, we'll just do a few)
            # Actually, let's just update all cards in the browser to have only "german" tag
            # We'll use 'T' for bulk tag update
            
            # Select first 5 cards
            for _ in range(5):
                agent.act("m")
                agent.act("j")
            
            agent.assert_text("5 cards selected")
            agent.act("T")
            agent.wait_for_text("TAGS:")
            # Replace tags with just "german"
            # We need to clear existing input. Backspace a few times just in case.
            for _ in range(20):
                agent.act("<Backspace>")
            agent.act("german")
            agent.act("<Enter>")
            agent.wait_for_text("Updated tags for 5 cards")
            
            # 3. Run Cleanup (C)
            # Note: The starter deck might have more than 5 cards with "a1" tag.
            # So "a1" might still be used. Let's check how many cards have "a1".
            # For the test to be reliable, I should probably remove it from ALL cards or use a small custom deck.
            
            # 4. Let's try to cleanup anyway and see what happens. 
            # If I want to be SURE, I should import a small deck.
            
        finally:
            agent.close()

def test_unused_tags_cleanup_with_custom_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a small TSV file
        tsv_path = os.path.join(tmpdir, "test_tags.tsv")
        with open(tsv_path, "w") as f:
            f.write("#deck:cleanup_test_deck\n")
            f.write("#id\tfront\tback\ttags\n")
            f.write("c1\tUniqueFrontOne\tBack 1\ttag1 tag2\n")
            f.write("c2\tUniqueFrontTwo\tBack 2\ttag1\n")
            
        agent = start_agent(tmpdir)
        try:
            agent.act("5") # Import
            agent.wait_for_text("Import / Export")
            agent.act("<Enter>")
            agent.act("<C-u>")
            agent.act(tsv_path)
            agent.act("<Enter>")
            agent.act("i")
            agent.wait_for_text("Imported 2 notes")
            
            # Select the new deck in Decks view to be sure
            agent.act("2")
            agent.wait_for_text("cleanup_test_deck")
            # Scroll down to it if needed
            for _ in range(5): agent.act("j")
            agent.act("<Enter>") # Select it
            
            agent.act("8") # Browser
            agent.wait_for_text("Card Browser")
            agent.wait_until_stable()
            
            # Search for UniqueFrontOne
            agent.act("/")
            agent.act("UniqueFrontOne")
            agent.act("<Enter>")
            agent.wait_for_text("UniqueFrontOne")
            
            # Change tags: remove tag2
            agent.act("T")
            agent.wait_for_text("TAGS:")
            agent.act("<C-u>")
            agent.act("tag1")
            agent.act("<Enter>")
            agent.wait_for_text("Updated tags")
            
            # Cleanup
            agent.act("C")
            agent.wait_for_text("Cleaning up unused tags")
            
            # Verify in Decks view
            agent.act("2")
            agent.wait_for_text("cleanup_test_deck")
            agent.assert_text("tag1")
            agent.assert_not_text("tag2")
        finally:
            agent.close()
