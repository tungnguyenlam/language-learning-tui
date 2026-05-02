import os
import tempfile
import time
from tui_tester import TUIAgent

def test_ai_view():
    """Test navigating to and interacting with the AI view"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Navigate to AI view
            agent.act('6')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "AI Drafts" in text
            assert "Topic:" in text
            
            # Try typing in the AI input field
            agent.act('H')
            agent.act('a')
            agent.act('l')
            agent.act('l')
            agent.act('o')
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

def test_import_view():
    """Test navigating to and interacting with the Import view"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Navigate to Import view
            agent.act('5')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Import / Export" in text
            assert "Import file:" in text
            assert "Export file:" in text
            
            # Try navigating between fields
            agent.act('j')  # Move down
            agent.wait_until_stable(timeout=2.0)
            
            agent.act('k')  # Move up
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

def test_settings_view():
    """Test navigating to and interacting with the Settings view"""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = TUIAgent(f'go run ../cmd/deutsch-tui --data-dir "{tmpdir}"')
        
        try:
            # Wait for initialization
            agent.wait_until_stable(timeout=5.0)
            
            # Navigate to Settings view
            agent.act('7')
            agent.wait_until_stable(timeout=2.0)
            text = agent.observe()
            assert "Settings" in text
            assert "AI Provider:" in text
            assert "Daily Goal:" in text
            
            # Try navigating between settings
            agent.act('j')  # Move down
            agent.wait_until_stable(timeout=2.0)
            
            agent.act('k')  # Move up
            agent.wait_until_stable(timeout=2.0)
            
            # Try adjusting daily goal
            agent.act('+')
            agent.wait_until_stable(timeout=2.0)
            
            agent.act('-')
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