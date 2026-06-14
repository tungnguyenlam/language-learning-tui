import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir} -test-mode', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def test_family_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            agent.wait_for_text("Seeding standard content...")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST", timeout=5.0)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Family")
            agent.act("<Enter>")
            agent.wait_for_text("Family", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Family", timeout=5.0)
        finally:
            agent.close()

def test_animals_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("Animals")
            agent.act("<Enter>")
            agent.wait_for_text("Animals", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Animals", timeout=5.0)
        finally:
            agent.close()

def test_colors_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("Colors")
            agent.act("<Enter>")
            agent.wait_for_text("Colors", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Colors", timeout=5.0)
        finally:
            agent.close()

def test_family_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Mutter")
            agent.act("<Enter>")
            agent.wait_for_text("Mutter", timeout=10.0)
        finally:
            agent.close()

def test_colors_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("blau")
            agent.act("<Enter>")
            agent.wait_for_text("blau", timeout=10.0)
        finally:
            agent.close()

def test_body_health_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("Body")
            agent.act("<Enter>")
            agent.wait_for_text("Body", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Body", timeout=5.0)
        finally:
            agent.close()

def test_clothing_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("Clothing")
            agent.act("<Enter>")
            agent.wait_for_text("Clothing", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Clothing", timeout=5.0)
        finally:
            agent.close()

def test_school_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("School")
            agent.act("<Enter>")
            agent.wait_for_text("School", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("School", timeout=5.0)
        finally:
            agent.close()

def test_numbers_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("/")
            agent.act("Numbers")
            agent.act("<Enter>")
            agent.wait_for_text("Numbers", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD")
            agent.wait_for_text("Numbers", timeout=5.0)
        finally:
            agent.close()

def test_animals_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            agent.act("5")
            agent.wait_for_text("Import / Export")
            agent.act("S")
            time.sleep(2.0)
            agent.wait_until_stable()
            agent.act("2")
            agent.wait_for_text("DECK LIST")
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Hund")
            agent.act("<Enter>")
            agent.wait_for_text("Hund", timeout=10.0)
        finally:
            agent.close()

if __name__ == "__main__":
    test_family_deck_exists()
    test_animals_deck_exists()
    test_colors_deck_exists()
    test_family_vocabulary_accessible()
    test_colors_vocabulary_accessible()
    test_body_health_deck_exists()
    test_clothing_deck_exists()
    test_school_deck_exists()
    test_numbers_deck_exists()
    test_animals_vocabulary_accessible()
