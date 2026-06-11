package content

import (
	"testing"
)

func TestA1PracticalLifeDeck(t *testing.T) {
	deck := A1PracticalLifeDeck()
	if deck.ID != "a1-practical" {
		t.Fatalf("deck.ID = %q, want a1-practical", deck.ID)
	}
	if len(deck.Notes) != 49 {
		t.Fatalf("len(deck.Notes) = %d, want 49", len(deck.Notes))
	}

	foundRestaurant := false
	for _, note := range deck.Notes {
		if len(note.Cards) == 0 {
			t.Fatalf("note %q has no generated cards", note.ID)
		}
		for _, tag := range note.Tags {
			if tag == "restaurant" {
				foundRestaurant = true
				break
			}
		}
	}
	if !foundRestaurant {
		t.Error("expected to find restaurant notes in A1PracticalLifeDeck")
	}
}

func TestStandardDecksIncludesA1PracticalLife(t *testing.T) {
	decks := StandardDecks()
	found := false
	for _, d := range decks {
		if d.ID == "a1-practical" {
			found = true
			break
		}
	}
	if !found {
		t.Error("A1PracticalLifeDeck not found in StandardDecks")
	}
}
