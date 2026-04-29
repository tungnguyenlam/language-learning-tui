import os
import tempfile
import time
from tui_tester import TUIAgent

def test_less_agent_workflow():
    # Create a temporary file with known content
    with tempfile.NamedTemporaryFile(mode='w', delete=False) as f:
        f.write("Line 1\nLine 2\nLine 3\nHello TUI Tester!\nLine 5\n")
        temp_path = f.name

    try:
        # Spawn less using the agent
        agent = TUIAgent(f'less {temp_path}')
        
        # Wait for initialization
        agent.wait_until_stable()
        
        # Observe the screen text
        text = agent.observe()
        assert "Hello TUI Tester!" in text
        assert "Line 1" in text
        
        # Send action to quit less
        agent.act('q')
        
        # Wait for the process to exit
        # We can just poll done, wait a maximum of 2 seconds
        start = time.time()
        while not agent.done and time.time() - start < 2.0:
            try:
                agent.observe()
            except Exception:
                pass
            time.sleep(0.1)
            
        assert agent.done is True
        
    finally:
        os.unlink(temp_path)

def test_wait_for_text():
    # Use python to simulate a slow print
    script = "import time, sys; sys.stdout.write('Wait for it...'); sys.stdout.flush(); time.sleep(1); sys.stdout.write(' Boom!'); sys.stdout.flush(); time.sleep(0.5)"
    agent = TUIAgent(f"python3 -c \"{script}\"")
    
    # Wait for the first part
    agent.wait_for_text("Wait for it...", timeout=2.0)
    
    # Let's wait for boom explicitly
    agent.wait_for_text("Boom!", timeout=3.0)
    
    agent.close()
