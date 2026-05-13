package content

import (
	"strings"
	"testing"
)

func TestB1TransportDeckHasExpectedShape(t *testing.T) {
	deck := B1TransportDeck()
	if deck.ID != "b1-transport" {
		t.Fatalf("deck.ID = %q, want b1-transport", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("transport deck should have >=30 notes, got %d", len(deck.Notes))
	}
	if !strings.Contains(deck.Name, "Public Transport") {
		t.Fatalf("deck.Name = %q, want it to mention Public Transport", deck.Name)
	}
	hasZug := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Zug") {
			hasZug = true
		}
	}
	if !hasZug {
		t.Fatal("expected at least one card with 'Zug' (train)")
	}
}

func TestA1OfficeDeckHasExpectedShape(t *testing.T) {
	deck := A1OfficeDeck()
	if deck.ID != "a1-office" {
		t.Fatalf("deck.ID = %q, want a1-office", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("office deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasComputer := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Computer") {
			hasComputer = true
		}
	}
	if !hasComputer {
		t.Fatal("expected at least one card with 'Computer'")
	}
}

func TestB2ClimateDeckHasExpectedShape(t *testing.T) {
	deck := B2ClimateDeck()
	if deck.ID != "b2-climate" {
		t.Fatalf("deck.ID = %q, want b2-climate", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("climate deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasKlimawandel := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Klimawandel") {
			hasKlimawandel = true
		}
	}
	if !hasKlimawandel {
		t.Fatal("expected at least one card with 'Klimawandel' (climate change)")
	}
}

func TestNewDecksAreInStandardDecks(t *testing.T) {
	decks := StandardDecks()
	found := map[string]bool{}
	for _, d := range decks {
		found[d.ID] = true
	}
	for _, id := range []string{"b1-transport", "a1-office", "b2-climate"} {
		if !found[id] {
			t.Errorf("expected deck %q to be in StandardDecks()", id)
		}
	}
}

func TestNewGrammarTipsPresent(t *testing.T) {
	// Spot-check that new grammar tip titles are in the rotation
	titles := map[string]bool{}
	for _, tip := range grammarTips {
		titles[tip.Title] = true
	}
	expected := []string{
		"Two-Way Prepositions (Wechselpräpositionen)",
		"Negation: nicht vs kein",
		"Verbs with Dativ",
		"Comparative & Superlative",
		"Time-Manner-Place (TeKaMoLo)",
		"es gibt + Akkusativ",
	}
	for _, e := range expected {
		if !titles[e] {
			t.Errorf("expected grammar tip %q to be present", e)
		}
	}
}
