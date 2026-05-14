package content

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

func TestBatch6BureaucracyDeck(t *testing.T) {
	deck := B1BureaucracyAppointmentsDeck()
	if deck.ID != "b1-bureaucracy-appointments" {
		t.Fatalf("deck.ID = %q, want b1-bureaucracy-appointments", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("bureaucracy deck should have >=30 notes, got %d", len(deck.Notes))
	}
	if !deckContainsFront(deck, "die Meldebescheinigung") || !deckContainsFront(deck, "einen Termin vereinbaren") {
		t.Fatal("bureaucracy deck should include appointment and document vocabulary")
	}
}

func TestBatch6DigitalPrivacyDeck(t *testing.T) {
	deck := B2DigitalPrivacyDeck()
	if deck.ID != "b2-digital-privacy" {
		t.Fatalf("deck.ID = %q, want b2-digital-privacy", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("digital privacy deck should have >=30 notes, got %d", len(deck.Notes))
	}
	if !deckContainsFront(deck, "der Datenschutz") || !deckContainsFront(deck, "die Zwei-Faktor-Authentifizierung") {
		t.Fatal("digital privacy deck should include Datenschutz and 2FA vocabulary")
	}
}

func TestBatch6NewDecksAreInStandardDecks(t *testing.T) {
	decks := StandardDecks()
	found := map[string]bool{}
	for _, deck := range decks {
		found[deck.ID] = true
	}
	for _, id := range []string{"b1-bureaucracy-appointments", "b2-digital-privacy"} {
		if !found[id] {
			t.Fatalf("expected deck %q in StandardDecks()", id)
		}
	}
}

func TestStandardDeckNoteAndCardIDsAreUnique(t *testing.T) {
	for _, deck := range StandardDecks() {
		noteIDs := map[string]bool{}
		cardIDs := map[string]bool{}
		for _, note := range deck.Notes {
			if strings.TrimSpace(note.ID) == "" {
				t.Fatalf("deck %q has blank note ID for front %q", deck.ID, note.Front)
			}
			if noteIDs[note.ID] {
				t.Fatalf("deck %q has duplicate note ID %q", deck.ID, note.ID)
			}
			noteIDs[note.ID] = true
			for _, card := range note.Cards {
				if strings.TrimSpace(card.ID) == "" {
					t.Fatalf("deck %q note %q has blank card ID", deck.ID, note.ID)
				}
				if cardIDs[card.ID] {
					t.Fatalf("deck %q has duplicate card ID %q", deck.ID, card.ID)
				}
				cardIDs[card.ID] = true
			}
		}
	}
}

func deckContainsFront(deck core.Deck, front string) bool {
	for _, note := range deck.Notes {
		if note.Front == front {
			return true
		}
	}
	return false
}
