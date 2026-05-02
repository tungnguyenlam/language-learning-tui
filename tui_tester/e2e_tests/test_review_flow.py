import os
import tempfile
import time
from tui_tester import TUIAgent

def test_review_flow():
    """Test the basic review flow with flashcards"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Navigate to Review view
            agent.act('3')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Review" in text
            
            # Check if there are cards to review
            if "No cards due" not in text:
                # Try to reveal a card
                agent.act(' ')
                agent.wait_until_stable(timeout=2.0)
                text = agent.observe()
                assert "Press space or enter to reveal" not in text  # Card should be revealed
                
                # Try to grade the card (good)
                agent.act('g')
                agent.wait_until_stable(timeout=2.0)
                
            # Quit the app
            agent.act('q')
            
            # Wait for the process to exit
            start = time.time()
            while not agent.done and time.time() - start < 3.0:
                try:
                    agent.observe()
                except Exception:
                    pass
                time.sleep(0.1)
                
            assert agent.done is True
            
        finally:
            agent.close()

def test_bookmark_toggle():
    """Test toggling bookmarks on cards"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Navigate to Review view
            agent.act('3')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            
            # If there are cards to review
            if "No cards due" not in text:
                # Try to toggle bookmark
                agent.act('b')
                agent.wait_until_stable(timeout=2.0)
                # Should see some status change
                
                # Toggle bookmark back
                agent.act('b')
                agent.wait_until_stable(timeout=2.0)
                
            # Quit the app
            agent.act('q')
            
            # Wait for the process to exit
            start = time.time()
            while not agent.done and time.time() - start < 3.0:
                try:
                    agent.observe()
                except Exception:
                    pass
                time.sleep(0.1)
                
            assert agent.done is True
            
        finally:
            agent.close()

def test_deck_switching():
    """Test switching between decks"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Check initial deck
            text = agent.observe()
            initial_deck_present = "Active Deck:" in text
            
            # Try to switch to next deck
            agent.act(']')
            agent.wait_until_stable(timeout=2.0)
            
            # Try to switch to previous deck
            agent.act('[')
            agent.wait_until_stable(timeout=2.0)
            
            # Quit the app
            agent.act('q')
            
            # Wait for the process to exit
            start = time.time()
            while not agent.done and time.time() - start < 3.0:
                try:
                    agent.observe()
                except Exception:
                    pass
                time.sleep(0.1)
                
            assert agent.done is True
            
        finally:
            agent.close()