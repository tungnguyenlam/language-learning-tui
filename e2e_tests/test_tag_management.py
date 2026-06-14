import os
import sys
import tempfile

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))

from tui_tester import TUIAgent


def start_agent(tmpdir, columns=120, lines=30):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=columns, lines=lines)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent


def test_browser_tag_management():
    with tempfile.TemporaryDirectory() as tmpdir:
        import_path = os.path.join(tmpdir, "import.tsv")
        with open(import_path, "w", encoding="utf-8") as f:
            f.write("#separator:tab\n")
            f.write("test-card-1\tApple\tApfel\t\tinitial-tag\tTest Deck\tBasic\n")

        agent = start_agent(tmpdir, columns=120, lines=30)
        try:
            # Import the test card
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("i")
            agent.wait_for_text("Imported 1 notes")

            # Go to Browser
            agent.act("8")
            agent.wait_for_text("Card Browser")
            
            # Switch to Test Deck
            agent.act("]")
            agent.wait_for_text("Apple")
            agent.wait_for_text("#initial-tag")

            # Open tag input
            agent.act("T")
            agent.wait_for_text("TAGS: initial-tag_")

            # Change tags
            # Use 11 backspaces to clear "initial-tag"
            for _ in range(11):
                agent.act("<Backspace>")
            agent.act("new-tag")
            agent.act("<Enter>")

            # Verify updated tags
            agent.wait_for_text("#new-tag")
            agent.assert_not_text("#initial-tag")

            # Verify search by tag
            agent.act("/")
            agent.act("new-tag")
            agent.act("<Enter>")
            agent.wait_for_text("Apple")
            agent.wait_for_text("1 cards found")

            agent.act("/")
            agent.act("initial-tag")
            agent.act("<Enter>")
            agent.wait_for_text("No cards found")
            
        finally:
            agent.close()

if __name__ == "__main__":
    test_browser_tag_management()
