import os
import sys
import tempfile
import pytest

sys.path.insert(
    0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester"))
)
from tui_tester import TUIAgent


def start_agent(tmpdir):
    app_cmd = os.getenv("DEUTSCH_TUI_BIN", "go run ./cmd/deutsch-tui")
    agent = TUIAgent(f"{app_cmd} -data-dir {tmpdir}", columns=100, lines=50)
    agent.wait_for_text("Dashboard", timeout=15.0)
    agent.wait_until_stable()

    # Seed standard content
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    import time

    time.sleep(2.0)
    agent.wait_until_stable()

    return agent


def test_b1_art_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")

            agent.act("<Esc>")
            agent.wait_until_stable()

            agent.act("/")
            agent.act("literature")
            agent.wait_for_text("B1 Art & Literature")
        finally:
            agent.close()


def test_a2_hobbies_ii_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("hobbies")
            agent.wait_for_text("A2 Hobbies & Free Time II")
        finally:
            agent.close()


def test_b2_science_ii_deck():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("science")
            agent.wait_for_text("B2 Science & Nature II")
        finally:
            agent.close()


def test_verb_of_the_day_presence():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("1")
            agent.wait_for_text("Verb:")
            agent.wait_for_text("ich")
            agent.wait_for_text("du")
        finally:
            agent.close()


def test_cram_mode_visuals():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram Mode")
            agent.act("5")  # Filter all
            agent.act("<Enter>")
            agent.wait_for_text("Cram Review")
            agent.wait_for_text("Type: flashcard")  # should show metadata
            agent.wait_for_text("p play audio")
        finally:
            agent.close()


def test_cram_mode_reveal_and_grade():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("9")
            agent.wait_for_text("Cram Mode")
            agent.act("5")
            agent.act("<Enter>")
            agent.wait_for_text("Press Space or Enter to reveal.")
            agent.act(" ")
            agent.wait_for_text("Grade: a Again | h Hard")
        finally:
            agent.close()


def test_deck_view_visuals():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            # The counts are abbreviated to N D T on terminal widths < 100
            agent.wait_for_text("N ")
            agent.wait_for_text("D ")
            agent.wait_for_text("T ")
            agent.wait_for_text("| today")
        finally:
            agent.close()


def test_browser_empty_search():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("8")
            agent.wait_for_text("Browser")
            agent.act("/")
            agent.act("XYZXYZ_NonExistent")
            agent.act("<Enter>")
            agent.wait_for_text("No cards found")
        finally:
            agent.close()
