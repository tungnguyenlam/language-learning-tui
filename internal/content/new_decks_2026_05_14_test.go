package content

import (
	"strings"
	"testing"
)

func TestA1AnimalsDeckHasExpectedShape(t *testing.T) {
	deck := A1AnimalsDeck()
	if deck.ID != "a1-animals" {
		t.Fatalf("deck.ID = %q, want a1-animals", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("animals deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasHund, hasKatze := false, false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Hund") {
			hasHund = true
		}
		if strings.Contains(n.Front, "Katze") {
			hasKatze = true
		}
	}
	if !hasHund || !hasKatze {
		t.Fatal("expected animals deck to have Hund and Katze")
	}
}

func TestA2BodyHealthDeckHasExpectedShape(t *testing.T) {
	deck := A2BodyHealthDeck()
	if deck.ID != "a2-body-health" {
		t.Fatalf("deck.ID = %q, want a2-body-health", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("body-health deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasArzt := false
	for _, n := range deck.Notes {
		if strings.Contains(n.Front, "Arzt") {
			hasArzt = true
		}
	}
	if !hasArzt {
		t.Fatal("expected at least one card with 'Arzt'")
	}
}

func TestB1CookingDeckHasExpectedShape(t *testing.T) {
	deck := B1CookingDeck()
	if deck.ID != "b1-cooking" {
		t.Fatalf("deck.ID = %q, want b1-cooking", deck.ID)
	}
	if len(deck.Notes) < 30 {
		t.Fatalf("cooking deck should have >=30 notes, got %d", len(deck.Notes))
	}
	hasKochen := false
	for _, n := range deck.Notes {
		if n.Front == "kochen" {
			hasKochen = true
		}
	}
	if !hasKochen {
		t.Fatal("expected at least one card with 'kochen'")
	}
}

func TestNewDecks20260514AreInStandardDecks(t *testing.T) {
	decks := StandardDecks()
	found := map[string]bool{}
	for _, d := range decks {
		found[d.ID] = true
	}
	for _, id := range []string{"a1-animals", "a2-body-health", "b1-cooking"} {
		if !found[id] {
			t.Errorf("expected deck %q to be in StandardDecks()", id)
		}
	}
}

func TestNewGrammarTips20260514Present(t *testing.T) {
	titles := map[string]bool{}
	for _, tip := range grammarTips {
		titles[tip.Title] = true
	}
	expected := []string{
		"Modal Verbs Word Order",
		"Genitive Case",
		"Reflexive Verbs (Akkusativ)",
		"Future Tense (Futur I)",
		"Konjunktiv II for Politeness",
		"Relative Clauses",
	}
	for _, e := range expected {
		if !titles[e] {
			t.Errorf("expected grammar tip %q to be present", e)
		}
	}
}

func TestNewVerbsOfTheDay20260514Present(t *testing.T) {
	germans := map[string]bool{}
	for _, v := range dailyVerbs {
		germans[v.German] = true
	}
	expected := []string{"schlafen", "tragen", "treffen", "bleiben", "denken", "antworten", "helfen", "warten"}
	for _, e := range expected {
		if !germans[e] {
			t.Errorf("expected verb %q to be present in dailyVerbs", e)
		}
	}
}

func TestAllDeckIDsAreUnique(t *testing.T) {
	decks := StandardDecks()
	seen := map[string][]string{}
	for _, d := range decks {
		seen[d.ID] = append(seen[d.ID], d.Name)
	}
	hasDupes := false
	for id, names := range seen {
		if len(names) > 1 {
			t.Errorf("duplicate deck ID: %q (%d copies): %v", id, len(names), names)
			hasDupes = true
		}
	}
	if hasDupes {
		t.Fatal("duplicate deck IDs found - fix before committing")
	}
}

func TestNewDecksAllNotesHaveCards(t *testing.T) {
	cases := []struct {
		name  string
		notes int
		first bool
	}{}
	an := A1AnimalsDeck()
	bh := A2BodyHealthDeck()
	ck := B1CookingDeck()
	cases = append(cases,
		struct {
			name  string
			notes int
			first bool
		}{"a1-animals", len(an.Notes), len(an.Notes) > 0 && len(an.Notes[0].Cards) > 0},
		struct {
			name  string
			notes int
			first bool
		}{"a2-body-health", len(bh.Notes), len(bh.Notes) > 0 && len(bh.Notes[0].Cards) > 0},
		struct {
			name  string
			notes int
			first bool
		}{"b1-cooking", len(ck.Notes), len(ck.Notes) > 0 && len(ck.Notes[0].Cards) > 0},
	)
	for _, c := range cases {
		if c.notes == 0 {
			t.Errorf("%s: no notes", c.name)
		}
		if !c.first {
			t.Errorf("%s: first note has no cards", c.name)
		}
	}
}
