package content

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

func TestA2BankingErrandsDeckLoaded(t *testing.T) {
	deck := findStandardDeck(t, "banking", "errands")

	if len(deck.Notes) < 40 {
		t.Fatalf("a2 banking errands deck should have >=40 notes, got %d", len(deck.Notes))
	}
	assertDeckCardsGenerated(t, deck)
	assertDeckContains(t, deck, "Geld abheben", "Kontostand")
}

func TestB1InsuranceClaimsDeckLoaded(t *testing.T) {
	deck := findStandardDeck(t, "insurance", "claims")

	if len(deck.Notes) < 40 {
		t.Fatalf("b1 insurance claims deck should have >=40 notes, got %d", len(deck.Notes))
	}
	assertDeckCardsGenerated(t, deck)
	assertDeckContains(t, deck, "Schaden", "Widerspruch")
}

func TestA2SchoolChildcareDeckLoaded(t *testing.T) {
	deck := findStandardDeck(t, "school", "childcare")

	if len(deck.Notes) < 40 {
		t.Fatalf("a2 school childcare deck should have >=40 notes, got %d", len(deck.Notes))
	}
	assertDeckCardsGenerated(t, deck)
	assertDeckContains(t, deck, "Elternabend", "krankmelden", "Abholzeit")
}

func findStandardDeck(t *testing.T, parts ...string) core.Deck {
	t.Helper()

	for _, deck := range StandardDecks() {
		id := strings.ToLower(deck.ID)
		matches := true
		for _, part := range parts {
			if !strings.Contains(id, part) {
				matches = false
				break
			}
		}
		if matches {
			return deck
		}
	}

	t.Fatalf("expected deck containing ID parts %v in StandardDecks()", parts)
	return core.Deck{}
}

func assertDeckCardsGenerated(t *testing.T, deck core.Deck) {
	t.Helper()

	for _, note := range deck.Notes {
		if len(note.Cards) == 0 {
			t.Fatalf("note %q has no generated cards", note.Front)
		}
	}
}

func assertDeckContains(t *testing.T, deck core.Deck, needles ...string) {
	t.Helper()

	for _, needle := range needles {
		found := false
		for _, note := range deck.Notes {
			if strings.Contains(note.Front, needle) || strings.Contains(note.Back, needle) || strings.Contains(note.Extra, needle) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected deck %q to contain %q", deck.ID, needle)
		}
	}
}
