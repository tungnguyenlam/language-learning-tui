import sys
import os
import pytest

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def test_dashboard_and_review_flow():
    # Start the Go app using the agent
    agent = TUIAgent('go run ./cmd/deutsch-tui', columns=100, lines=30)
    
    try:
        # Wait for the dashboard to load and become stable
        agent.wait_until_stable(timeout=10.0)
        
        # Verify dashboard content
        agent.assert_text("Dashboard")
        agent.assert_text("Deck: German A1 Survival")
        agent.assert_text("Due cards: 6")
        
        # Switch to Review view
        agent.act('2')
        agent.wait_until_stable()
        
        # Verify review screen for the first card
        agent.assert_text("Review 1/6")
        agent.assert_text("der Apfel")
        agent.assert_text("Press space or enter to reveal.")
        
        # Reveal the card using Enter
        agent.act('<Enter>')
        agent.wait_until_stable()
        
        # Verify revealed content
        agent.assert_text("apple")
        agent.assert_text("Grade: a Again | h Hard | g Good | e")
        
        # Press 'e' for Easy
        agent.act('e')
        agent.wait_until_stable()
        
        # Verify the status bar updated (since grading doesn't advance yet based on current code)
        agent.assert_text("Grade: Easy")
        
        # Try returning to Dashboard
        agent.act('1')
        agent.wait_until_stable()
        agent.assert_text("Use Review to start studying.")
        
    finally:
        # Close the agent
        agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
