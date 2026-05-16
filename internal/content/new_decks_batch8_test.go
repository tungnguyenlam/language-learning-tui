package content

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

func TestB2JobApplicationDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "job") && strings.Contains(d.ID, "application") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected B2 Job Application deck to be in StandardDecks()")
	}
	if len(found.Notes) < 40 {
		t.Fatalf("b2 job application deck should have >=40 notes, got %d", len(found.Notes))
	}
	for _, n := range found.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}

func TestC1PhilosophyEthicsDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "philosophy") && strings.Contains(d.ID, "ethics") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected C1 Philosophy & Ethics deck to be in StandardDecks()")
	}
	if len(found.Notes) < 40 {
		t.Fatalf("c1 philosophy ethics deck should have >=40 notes, got %d", len(found.Notes))
	}
	for _, n := range found.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}

func TestA1CityDirectionsDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "city") && strings.Contains(d.ID, "directions") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected A1 City & Directions deck to be in StandardDecks()")
	}
	if len(found.Notes) < 40 {
		t.Fatalf("a1 city directions deck should have >=40 notes, got %d", len(found.Notes))
	}
	for _, n := range found.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}
