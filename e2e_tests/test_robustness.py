import sys
import os
import pytest
import tempfile
import time

# Add tui_tester to sys.path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_settings_template_editing_cancel():
    """Test that Esc cancels template editing and reverts changes."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Settings
            agent.act('7')
            agent.wait_for_text("Settings")
            
            # Move to Front Template (index 1)
            agent.act('j')
            agent.wait_until_stable()
            
            # Start editing
            agent.act('<Enter>')
            agent.wait_for_text("EDITING")
            
            # Add some text
            agent.act('X')
            agent.act('Y')
            agent.act('Z')
            agent.assert_text("Front Template: {{.Topic}}XYZ")
            
            # Cancel with Esc - should revert changes
            agent.act('<Esc>')
            agent.wait_until_stable()
            agent.assert_not_text("EDITING")
            # Esc should revert the template to its original value
            agent.assert_text("Front Template: {{.Topic}}")
            agent.assert_not_text("Front Template: {{.Topic}}XYZ")


        finally:
            agent.close()

def test_import_nonexistent_file_shows_error():
    """Test that importing a non-existent file shows an error message."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('5')
            agent.wait_for_text("Import / Export")
            
            # workflow: Enter, Ctrl-u, type nonexistent path, Enter, i
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.act('<Ctrl-u>')
            for char in "nonexistent.tsv":
                agent.act(char)
            agent.act('<Enter>')
            agent.wait_until_stable()

            # 'i' triggers import
            agent.act('i')
            agent.wait_for_text("Error:")
            agent.assert_text("no such file or directory")
        finally:
            agent.close()

def test_review_grade_updates_dashboard_count():
    """Test that grading a card updates the due count on the dashboard."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.assert_text("Due cards:   52")
            
            agent.act('3')
            agent.wait_for_text("Review 1/52")
            
            # Reveal and grade Good
            agent.act('<Space>')
            agent.wait_for_text("Grade: a Again")
            agent.act('g')
            agent.wait_for_text("51 cards due")
            
            # Go back to Dashboard
            agent.act('1')
            agent.wait_for_text("DASHBOARD")
            agent.assert_text("Due cards:   51")
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
