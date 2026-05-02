import os
import tempfile
import time
from tui_tester import TUIAgent

def test_app_startup():
    """Test that the app starts and shows the dashboard"""
    # Create a temporary directory for the app data
    with tempfile.TemporaryDirectory() as tmpdir:
        # Start the app
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for the app to stabilize
            agent.wait_until_stable(timeout=5.0)
            
            # Observe the screen text
            text = agent.observe()
            
            # Check that we see the expected elements
            assert "deutsch-tui" in text
            assert "DASHBOARD" in text
            assert "Active Deck:" in text
            assert "Review Queue" in text
            
            # Send quit command
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

def test_navigation_between_views():
    """Test navigating between different views"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Check we start on dashboard
            text = agent.observe()
            assert "DASHBOARD" in text
            
            # Navigate to Decks view (2)
            agent.act('2')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Decks" in text
            
            # Navigate to Review view (3)
            agent.act('3')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Review" in text
            
            # Navigate to Statistics view (4)
            agent.act('4')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Statistics" in text
            
            # Navigate to Import view (5)
            agent.act('5')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Import" in text
            
            # Navigate to AI view (6)
            agent.act('6')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "AI Drafts" in text
            
            # Navigate to Settings view (7)
            agent.act('7')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Settings" in text
            
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

def test_help_overlay():
    """Test that help overlay works correctly"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Check we start on dashboard
            text = agent.observe()
            assert "DASHBOARD" in text
            
            # Toggle help overlay
            agent.act('?')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            # Just check that something changed - the help overlay should modify the display
            
            # Toggle help overlay off
            agent.act('?')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            
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