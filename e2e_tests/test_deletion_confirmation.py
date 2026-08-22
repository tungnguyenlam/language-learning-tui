import os
import pytest
import tempfile
import sys
import time
import re

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

from e2e_helpers import read_cards_found

def start_agent(tmpdir, columns=100, lines=50):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_deck_deletion_confirmation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Decks view
            agent.act('2')
            agent.wait_for_text("Decks")
            agent.assert_text("German A1 Survival")
            
            # Move to German A1 Survival (it's at index 1, All Decks is at 0)
            agent.act('j')
            agent.wait_until_stable()

            # Press Backspace to trigger deletion confirmation
            agent.act('<Backspace>')
            agent.wait_for_text("CONFIRM DELETION")
            agent.assert_text("Delete 1 decks and ALL their cards?")

            # Cancel with 'n'
            agent.act('n')
            agent.wait_until_stable()
            agent.assert_not_text("CONFIRM DELETION")
            agent.assert_text("German A1 Survival")
            agent.assert_text("Deletion cancelled")

            # Trigger again and confirm with 'y'
            agent.act('<Backspace>')
            agent.wait_for_text("CONFIRM DELETION")
            agent.act('y')
            # Wait for deletion to complete and list to refresh
            def deck_gone():
                return "German A1 Survival" not in agent.observe()
            
            start = time.time()
            while time.time() - start < 5.0:
                if deck_gone():
                    break
                time.sleep(0.1)
            
            agent.assert_not_text("German A1 Survival")
        finally:
            agent.close()

def test_card_deletion_confirmation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Browser view
            agent.act('8')
            agent.wait_for_text("Card Browser")
            agent.assert_text("FC")
            
            # Count cards (fails loudly if the status line is missing)
            initial_count = read_cards_found(agent)

            # Press Backspace to trigger deletion confirmation
            agent.act('<Backspace>')
            agent.wait_for_text("CONFIRM DELETION")
            agent.assert_text("Delete card")

            # Cancel with 'esc'
            agent.act('<Esc>')
            agent.wait_until_stable()
            agent.assert_not_text("CONFIRM DELETION")
            agent.assert_text("FC")

            # Trigger again and confirm with 'enter'
            agent.act('<Backspace>')
            agent.wait_for_text("CONFIRM DELETION")
            agent.act('<Enter>')
            
            # Wait for count to decrease
            start = time.time()
            while time.time() - start < 5.0:
                screen = agent.observe()
                match = re.search(r'status: (\d+) cards found', screen)
                if match and int(match.group(1)) < initial_count:
                    break
                time.sleep(0.1)
            
            agent.assert_not_text("CONFIRM DELETION")
            screen = agent.observe()
            match = re.search(r'status: (\d+) cards found', screen)
            new_count = int(match.group(1)) if match else 0
            assert new_count < initial_count
        finally:
            agent.close()
