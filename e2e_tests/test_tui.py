import sys
import os
import pytest
import tempfile

# Add tui_tester to sys.path so we can import it
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../tui_tester')))

from tui_tester import TUIAgent

def start_agent(tmpdir, columns=100, lines=30):
    agent = TUIAgent(f'go run ./cmd/deutsch-tui -data-dir {tmpdir}', columns=columns, lines=lines)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_dashboard_and_review_flow():
    # Use a temporary directory for data to ensure a clean state
    with tempfile.TemporaryDirectory() as tmpdir:
        # Start the Go app using the agent
        agent = start_agent(tmpdir)

        try:
            # Verify dashboard content
            agent.assert_text("Dashboard")
            agent.assert_text("Deck: German A1 Survival")
            agent.assert_text("Due cards: 6")

            # Switch to Review view
            agent.act('3')
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

            # Verify it advanced to the next card (the MCQ for the same note)
            # Note: 1/5 because one card is now scheduled in the future
            agent.assert_text("Review 1/5")
            agent.assert_text("Ich esse einen Apfel.")
            agent.assert_text("5 cards due")

            # Try returning to Dashboard
            agent.act('1')
            agent.wait_until_stable()
            agent.assert_text("Use Review to start studying.")

        finally:
            # Close the agent
            agent.close()

def test_ai_draft_approval_persists_across_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)

        try:
            agent.act('5')
            agent.wait_until_stable()
            agent.assert_text("AI Drafts")
            agent.assert_text("Topic: der Kaffee")

            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("> der Kaffee -> German prompt for der")
            agent.assert_text("Kaffee")

            agent.act('a')
            agent.wait_until_stable()
            agent.assert_text("Draft approved")

            agent.act('1')
            agent.wait_until_stable()
            agent.assert_text("Due cards: 8")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text("Dashboard")
            restarted.assert_text("Due cards: 8")
        finally:
            restarted.close()

def test_import_tsv_adds_reviewable_deck_and_export_writes_file():
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        export_path = os.path.join(tmpdir, "export.tsv")
        with open(import_path, "w", encoding="utf-8") as handle:
            handle.write("#separator:tab\n")
            handle.write("import-1\tdie Bahn\ttrain\t\ttransport\tImported A1\tBasic\n")

        agent = start_agent(tmpdir, columns=110, lines=30)

        try:
            agent.act('4')
            agent.wait_until_stable()
            agent.assert_text("Import / Export")
            agent.assert_text("Press i to import TSV.")

            agent.act('i')
            agent.wait_until_stable()
            agent.assert_text("Imported 1 notes")

            agent.act(']')
            agent.wait_until_stable()
            agent.assert_text("Deck: Imported A1")

            agent.act('3')
            agent.wait_until_stable()
            agent.assert_text("die Bahn")

            agent.act('4')
            agent.wait_until_stable()
            agent.act('x')
            agent.wait_until_stable()
            agent.assert_text("Exported 1 notes")
        finally:
            agent.close()

        with open(export_path, "r", encoding="utf-8") as handle:
            exported = handle.read()
        assert "die Bahn\ttrain" in exported
        assert "\tImported A1\tBasic" in exported

def test_settings_and_template_drafting():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            # Go to Settings
            agent.act('6')
            agent.wait_until_stable()
            agent.assert_text("Settings")
            agent.assert_text("AI Provider: disabled")

            # Switch to offline, then template provider
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("AI Provider: offline")
            
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("AI Provider: template")
            agent.assert_text("Switched to template AI provider")

            # Edit Front Template
            agent.act('j') # Move to Front Template
            agent.wait_until_stable()
            agent.act('<Enter>') # Start editing
            agent.wait_until_stable()
            agent.assert_text("EDITING")
            
            # Add suffix to template
            agent.act('P')
            agent.act(':')
            agent.act(' ')
            agent.act('<Enter>') # Save
            agent.wait_until_stable()
            agent.assert_text("Front Template: {{.Topic}}P: ")

            # Go to AI Drafts and verify template is used
            agent.act('5')
            agent.wait_until_stable()
            agent.act('<Enter>') # Generate
            agent.wait_until_stable()
            agent.assert_text("> der KaffeeP: ")
        finally:
            agent.close()

def test_settings_persistence():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('6')
            agent.wait_until_stable()
            
            # Switch to offline, then template provider
            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("AI Provider: offline")

            agent.act('<Enter>')
            agent.wait_until_stable()
            agent.assert_text("AI Provider: template")
        finally:
            agent.close()

        # Restart and verify
        restarted = start_agent(tmpdir)
        try:
            restarted.act('6')
            restarted.wait_until_stable()
            restarted.assert_text("AI Provider: template")
        finally:
            restarted.close()

def test_all_core_views_render_with_keyboard_navigation():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            expected_views = [
                ('1', "Dashboard"),
                ('2', "Decks"),
                ('3', "Review 1/6"),
                ('4', "Import / Export"),
                ('5', "AI Drafts"),
                ('6', "Settings"),
            ]
            for key, text in expected_views:
                agent.act(key)
                agent.wait_for_text(text)
                agent.wait_until_stable()
                agent.assert_text(text)
        finally:
            agent.close()

def test_compact_layout_renders_all_core_views():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=60, lines=24)
        try:
            agent.assert_text("deutsch-tui compact")

            for key, text in [
                ('1', "Dashboard"),
                ('2', "Decks"),
                ('3', "Review 1/6"),
                ('4', "Import / Export"),
                ('5', "AI Drafts"),
                ('6', "Settings"),
            ]:
                agent.act(key)
                agent.wait_for_text(text)
                agent.wait_until_stable()
                agent.assert_text(text)
        finally:
            agent.close()

def test_space_reveal_and_again_grade_keyboard_flow():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.act('3')
            agent.wait_for_text("Review 1/6")

            agent.act('<Space>')
            agent.wait_for_text("Grade: a Again")
            agent.assert_text("apple")

            agent.act('a')
            agent.wait_for_text("Review 1/5")
            agent.assert_text("5 cards due")
            agent.assert_text("Press space or enter to reveal.")
        finally:
            agent.close()

def test_mouse_tabs_open_import_ai_and_settings_views():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.click(31, 3)
            agent.wait_for_text("Import / Export")
            agent.assert_text("Press i to import TSV.")

            agent.click(40, 3)
            agent.wait_for_text("AI Drafts")
            agent.assert_text("Topic: der Kaffee")

            agent.click(45, 3)
            agent.wait_for_text("Settings")
            agent.assert_text("AI Provider: disabled")
        finally:
            agent.close()

def test_review_grade_persists_to_sqlite_across_restart():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act('3')
            agent.wait_for_text("Review 1/6")
            agent.act('<Enter>')
            agent.wait_for_text("Grade: a Again")
            agent.act('e')
            agent.wait_for_text("5 cards due")
            agent.act('1')
            agent.wait_for_text("Due cards: 5")
        finally:
            agent.close()

        restarted = start_agent(tmpdir)
        try:
            restarted.assert_text("Due cards: 5")
            restarted.act('3')
            restarted.wait_for_text("Review 1/5")
        finally:
            restarted.close()

def test_mouse_tab_navigation_and_grade_button():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            # Medium layout tab row starts on terminal row 3; Review tab spans columns 21-28.
            agent.click(22, 3)
            agent.wait_for_text("Review 1/6")
            agent.act('<Enter>')
            agent.wait_for_text("Grade: a Again")

            # In the medium review panel, the Good grade hitbox is on terminal row 11.
            agent.click(30, 11)
            agent.wait_for_text("5 cards due")
            agent.assert_text("Review 1/5")
        finally:
            agent.close()

def test_view_navigation_with_arrow_keys():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.assert_text("Dashboard")
            
            agent.act('<Right>')
            agent.wait_for_text("Press enter to select deck.")
            agent.wait_until_stable()
            
            agent.act('<Right>')
            agent.wait_for_text("Review 1/")
            agent.wait_until_stable()
            
            agent.act('<Left>')
            agent.wait_for_text("Press enter to select deck.")
            agent.wait_until_stable()
            
            agent.act('<Left>')
            agent.wait_for_text("Use Review to start studying.")
            agent.wait_until_stable()
            agent.assert_text("Dashboard")
        finally:
            agent.close()

def test_deck_navigation_with_up_down_arrows():
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, "w", encoding="utf-8") as handle:
            handle.write("#separator:tab\n")
            handle.write("import-1\tA\tB\t\t\tZ-Imported\n")

        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.act('4')
            agent.wait_for_text("Press i to import TSV.")
            agent.act('i')
            agent.wait_for_text("Imported 1 notes")
            
            agent.act('2')
            agent.wait_for_text("Press enter to select deck.")
            
            agent.act('<Down>')
            agent.wait_for_text("> Z-Imported")
            
            agent.act('<Enter>')
            agent.wait_for_text("Use Review to start studying.")
            agent.wait_until_stable()
            agent.assert_text("Deck: Z-Imported")
        finally:
            agent.close()

def test_settings_navigation_with_up_down_arrows():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir, columns=90, lines=28)
        try:
            agent.act('6')
            agent.wait_for_text("AI Provider: disabled")
            
            agent.act('<Down>')
            agent.wait_until_stable()
            
            agent.act('<Enter>')
            agent.wait_for_text("EDITING")
            
            agent.act('<Esc>')
            agent.wait_until_stable()
            
            agent.act('<Up>')
            agent.wait_until_stable()
            
            agent.act('<Enter>')
            agent.wait_for_text("AI Provider: offline")
        finally:
            agent.close()

if __name__ == "__main__":
    pytest.main(["-v", __file__])
