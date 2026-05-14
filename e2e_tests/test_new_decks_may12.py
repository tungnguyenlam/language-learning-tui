import os
import sys
import tempfile
import time

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "../tui_tester")))
from tui_tester import TUIAgent

def start_agent(tmpdir):
    app_cmd = os.getenv('DEUTSCH_TUI_BIN', 'go run ./cmd/deutsch-tui')
    agent = TUIAgent(f'{app_cmd} -data-dir {tmpdir}', columns=100, lines=40)
    agent.wait_for_text("DASHBOARD", timeout=15.0)
    agent.wait_until_stable()
    return agent

def seed_and_go_decks(agent):
    agent.act("5")
    agent.wait_for_text("Import / Export")
    agent.act("S")
    agent.wait_for_text("Seeding standard content...")
    time.sleep(2.0)
    agent.wait_until_stable()
    agent.act("2")
    agent.wait_for_text("DECK LIST", timeout=5.0)
    agent.wait_until_stable()

def test_a2_daily_life_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Daily")
            agent.act("<Enter>")
            agent.wait_for_text("Daily", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Daily", timeout=5.0)
        finally:
            agent.close()

def test_b2_culture_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Culture")
            agent.act("<Enter>")
            agent.wait_for_text("Culture", timeout=10.0)
            agent.act("<Esc>")
            agent.wait_for_text("DECK LIST", timeout=5.0)
        finally:
            agent.close()

def test_c1_academic_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Academic")
            agent.act("<Enter>")
            agent.wait_for_text("Academic", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Academic", timeout=5.0)
        finally:
            agent.close()

def test_a2_shopping_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Shopping")
            agent.act("<Enter>")
            agent.wait_for_text("Shopping", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Shopping", timeout=5.0)
        finally:
            agent.close()

def test_b2_business_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Business")
            agent.act("<Enter>")
            agent.wait_for_text("Business", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Business", timeout=5.0)
        finally:
            agent.close()

def test_c2_legal_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Legal")
            agent.act("<Enter>")
            agent.wait_for_text("Legal", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Legal", timeout=5.0)
        finally:
            agent.close()

def test_daily_life_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("aufstehen")
            agent.act("<Enter>")
            agent.wait_for_text("aufstehen", timeout=10.0)
        finally:
            agent.close()

def test_media_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Nachrichten")
            agent.act("<Enter>")
            agent.wait_for_text("Nachrichten", timeout=10.0)
        finally:
            agent.close()

def test_shopping_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Preis")
            agent.act("<Enter>")
            agent.wait_for_text("Preis", timeout=10.0)
        finally:
            agent.close()

def test_academic_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Forschung")
            agent.act("<Enter>")
            agent.wait_for_text("Forschung", timeout=10.0)
        finally:
            agent.close()

def test_all_new_decks_count():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Shopping")
            agent.act("<Enter>")
            agent.wait_for_text("Shopping", timeout=10.0)
        finally:
            agent.close()

def test_b2_travel_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Travel")
            agent.act("<Enter>")
            agent.wait_for_text("Travel", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Travel", timeout=5.0)
        finally:
            agent.close()

def test_b2_environment_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Environment")
            agent.act("<Enter>")
            agent.wait_for_text("Environment", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Environment", timeout=5.0)
        finally:
            agent.close()

def test_b2_travel_vocabulary_accessible():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("<Esc>")
            agent.wait_until_stable()
            agent.act("/")
            agent.act("Reise")
            agent.act("<Enter>")
            agent.wait_for_text("Reise", timeout=10.0)
        finally:
            agent.close()

def test_c1_science_deck_exists():
    with tempfile.TemporaryDirectory() as tmpdir:
        agent = start_agent(tmpdir)
        try:
            seed_and_go_decks(agent)
            agent.act("/")
            agent.wait_for_text("Search:", timeout=2.0)
            agent.act("Science")
            agent.act("<Enter>")
            agent.wait_for_text("Science", timeout=10.0)
            agent.act("<Enter>")
            agent.wait_for_text("DASHBOARD", timeout=5.0)
            agent.wait_for_text("Science", timeout=5.0)
        finally:
            agent.close()

if __name__ == "__main__":
    test_a2_daily_life_deck_exists()
    test_b2_media_deck_exists()
    test_c1_academic_deck_exists()
    test_a2_shopping_deck_exists()
    test_b2_business_deck_exists()
    test_c2_legal_deck_exists()
    test_daily_life_vocabulary_accessible()
    test_media_vocabulary_accessible()
    test_shopping_vocabulary_accessible()
    test_academic_vocabulary_accessible()
    test_all_new_decks_count()
    test_b2_travel_deck_exists()
    test_b2_environment_deck_exists()
    test_b2_travel_vocabulary_accessible()
    test_c1_science_deck_exists()