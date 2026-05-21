package content

import (
	"testing"
)

func TestA2MedicalAppointmentDeck(t *testing.T) {
	deck := A2MedicalAppointmentDeck()
	if deck.ID != "a2-medical-appointment" {
		t.Errorf("expected ID 'a2-medical-appointment', got %q", deck.ID)
	}
	if len(deck.Notes) == 0 {
		t.Fatal("expected notes, got none")
	}
	found := false
	for _, n := range deck.Notes {
		if n.Front == "einen Termin vereinbaren" {
			found = true
			if len(n.Cards) == 0 {
				t.Error("note 'einen Termin vereinbaren' has no cards")
			}
			break
		}
	}
	if !found {
		t.Error("expected note 'einen Termin vereinbaren' not found")
	}
}

func TestA2MedicalAppointmentInStandardDecks(t *testing.T) {
	decks := StandardDecks()
	found := false
	for _, d := range decks {
		if d.ID == "a2-medical-appointment" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected deck 'a2-medical-appointment' to be in StandardDecks()")
	}
}
