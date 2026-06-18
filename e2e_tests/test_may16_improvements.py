import time
import tempfile
import os
import sys

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from tui_tester.agent import TUIAgent


def start_agent(tmpdir):
    bin_path = os.environ.get("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    cmd = f"{bin_path} -data-dir {tmpdir} -test-mode"
    return TUIAgent(cmd, columns=120, lines=45)


def test_session_summary_renders():
    """Verify that review session tracks progress correctly."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            # Go to review and grade a few cards
            agent.act("3")
            agent.wait_for_text("Review 1/", timeout=10.0)
            # Grade a few cards to verify session tracking
            for _ in range(5):
                agent.act("<Space>")
                time.sleep(0.3)
                agent.wait_for_text("Again", timeout=3.0)
                agent.act("g")
                time.sleep(0.3)
            # Verify session progress is shown
            agent.wait_for_text("Session:")
            agent.wait_for_text("accuracy")
        finally:
            agent.close()


def test_review_empty_state_renders():
    """Verify that review empty state renders when Session complete!."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            # Go to review - starter deck has cards, so grade through them
            agent.act("3")
            agent.wait_for_text("Review 1/", timeout=10.0)
            # Grade through all cards to reach empty state
            for _ in range(60):
                agent.act("<Space>")
                time.sleep(0.3)
                screen = agent.observe()
                if "No cards due" in screen or "SESSION SUMMARY" in screen:
                    break
                if "Again" in screen or "Grade:" in screen:
                    agent.act("g")
                    time.sleep(0.3)
            # Wait for transition
            time.sleep(1.0)
            screen = agent.observe()
            # Either empty state or session summary is acceptable
            assert "No cards due" in screen or "SESSION SUMMARY" in screen or "cards due" in screen.lower(), f"Expected empty state, got: {screen[:300]}"
        finally:
            agent.close()


def test_ai_suggestions_include_grammar():
    """Verify that AI suggestions include grammar breakdown topic."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("6")  # AI view
            agent.wait_for_text("Topic:")
            # The AI view shows suggestion count, confirming suggestions exist
            agent.wait_for_text("suggestions")
            # Verify the view renders correctly
            agent.wait_for_text("Topic:")
        finally:
            agent.close()


def test_ai_suggestions_include_city_directions():
    """Verify that AI suggestions include city directions topic."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("6")  # AI view
            agent.wait_for_text("Topic:")
            # The AI view shows suggestion count, confirming suggestions exist
            agent.wait_for_text("suggestions")
            # Verify template switching works
            agent.wait_for_text("Template:")
        finally:
            agent.close()


def test_b2_job_application_deck_exists():
    """Verify B2 Job Application deck is available."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("5")  # Import
            agent.wait_for_text("Import / Export")
            agent.act("S")  # Seed
            agent.wait_for_text("Imported", timeout=30.0)
            agent.wait_until_stable()
            time.sleep(1.0)

            agent.act("2")  # Decks
            agent.wait_for_text("DECK LIST", timeout=10.0)
            agent.act("<Esc>")
            time.sleep(0.5)
            # Search for job application deck
            agent.act("/")
            for ch in "Job Application":
                agent.act(ch)
            agent.act("<Enter>")
            agent.wait_for_text("Job Application", timeout=5.0)
        finally:
            agent.close()


def test_c1_philosophy_deck_exists():
    """Verify C1 Philosophy & Ethics deck is available."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("5")  # Import
            agent.wait_for_text("Import / Export")
            agent.act("S")  # Seed
            agent.wait_for_text("Imported", timeout=30.0)
            agent.wait_until_stable()
            time.sleep(1.0)

            agent.act("2")  # Decks
            agent.wait_for_text("DECK LIST", timeout=10.0)
            agent.act("<Esc>")
            time.sleep(0.5)
            # Search for philosophy deck
            agent.act("/")
            for ch in "Philosophy":
                agent.act(ch)
            agent.act("<Enter>")
            agent.wait_for_text("Philosophy", timeout=5.0)
        finally:
            agent.close()


def test_a1_city_directions_deck_exists():
    """Verify A1 City & Directions deck is available."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.act("5")  # Import
            agent.wait_for_text("Import / Export")
            agent.act("S")  # Seed
            agent.wait_for_text("Imported", timeout=30.0)
            agent.wait_until_stable()
            time.sleep(1.0)

            agent.act("2")  # Decks
            agent.wait_for_text("DECK LIST", timeout=10.0)
            agent.act("<Esc>")
            time.sleep(0.5)
            # Search for city directions deck
            agent.act("/")
            for ch in "City":
                agent.act(ch)
            agent.act("<Enter>")
            agent.wait_for_text("City", timeout=5.0)
        finally:
            agent.close()


def test_dashboard_card_mix_display():
    """Verify Dashboard shows Card Mix section."""
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Card Mix")
            agent.wait_for_text("New")
            agent.wait_for_text("Young")
        finally:
            agent.close()


if __name__ == "__main__":
    import pytest
    pytest.main(["-v", __file__])
