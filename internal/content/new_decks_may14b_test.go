package content

import (
	"strings"
	"testing"
)

func TestB1WeatherDeckHasExpectedShape(t *testing.T) {
	deck := B1WeatherDeck()
	if deck.ID != "b1-weather-seasons" {
		t.Fatalf("deck.ID = %q, want b1-weather-seasons", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("weather deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasSchnee := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Schnee") || strings.Contains(n.Front, "schneien") {
			hasSchnee = true
		}
	}
	if !hasSchnee {
		t.Fatal("expected weather deck to have snow-related vocabulary")
	}
}

func TestA2ShoppingDeckHasExpectedShape(t *testing.T) {
	deck := A2ShoppingDeck()
	if deck.ID != "a2_purchasing_wear" {
		t.Fatalf("deck.ID = %q, want a2_purchasing_wear", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("shopping deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasKaufen := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "kaufen") {
			hasKaufen = true
		}
	}
	if !hasKaufen {
		t.Fatal("expected shopping deck to have 'kaufen'")
	}
}

func TestB2CultureLeisureDeckHasExpectedShape(t *testing.T) {
	deck := B2CultureLeisureDeck()
	if deck.ID != "b2-culture-leisure" {
		t.Fatalf("deck.ID = %q, want b2-culture-leisure", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("culture deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasMuseum := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Museum") {
			hasMuseum = true
		}
	}
	if !hasMuseum {
		t.Fatal("expected culture deck to have 'Museum'")
	}
}

func TestNewDecks20260514bAreInStandardDecks(t *testing.T) {
	decks := StandardDecks()
	found := map[string]bool{}
	for _, d := range decks {
		found[d.ID] = true
	}
	for _, id := range []string{"b1-weather-seasons", "a2_purchasing_wear", "b2-culture-leisure"} {
		if !found[id] {
			t.Errorf("expected deck %q to be in StandardDecks()", id)
		}
	}
}

func TestGrammarTipsMay14bCount(t *testing.T) {
	want := 128
	got := len(grammarTips)
	if got < want {
		t.Errorf("expected >=%d grammar tips, got %d", want, got)
	}
}

func TestDailyVerbsMay14bCount(t *testing.T) {
	want := 109
	got := len(dailyVerbs)
	if got < want {
		t.Errorf("expected >=%d daily verbs, got %d", want, got)
	}
}
