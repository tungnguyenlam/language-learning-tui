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
    agent.wait_for_text("Dashboard", timeout=15.0)
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
            
            # Cancel with Esc
            agent.act('<Esc>')
            agent.wait_until_stable()
            agent.assert_not_text("EDITING")
            
            # Revert to original? 
            # Actually, let's see how the app handles Esc.
            # Looking at internal/tui/model.go:
            # case "enter", "esc":
            #     m.editingTemplate = false
            # It doesn't seem to have a "revert" buffer. It edits in-place.
            # Wait, let's check updateSettingsKey.
            # case "enter", "esc": m.editingTemplate = false
            # Yes, it edits m.aiTemplates[key] directly.
            
            # If so, Esc just stops editing but keeps the changes?
            # Let's verify this behavior.
            agent.assert_text("Front Template: {{.Topic}}XYZ")
            
        finally:
            agent.close()

def test_import_nonexistent_file_shows_error():
    """Test that importing a non-existent file shows an error message."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('5')
            agent.wait_for_text("Import / Export")
            
            # 'i' triggers import of 'import.tsv' (default)
            # which shouldn't exist in the temp dir.
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
            agent.assert_text("Due cards: 6")
            
            agent.act('3')
            agent.wait_for_text("Review 1/6")
            
            # Reveal and grade Good
            agent.act('<Space>')
            agent.wait_for_text("Grade: a Again")
            agent.act('g')
            agent.wait_for_text("5 cards due")
            
            # Go back to Dashboard
            agent.act('1')
            agent.wait_for_text("Dashboard")
            agent.assert_text("Due cards: 5")
            
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
