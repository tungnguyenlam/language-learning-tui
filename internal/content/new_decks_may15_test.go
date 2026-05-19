package content

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

func TestTravelAdventureDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "travel") && strings.Contains(d.ID, "adventure") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected Travel & Adventure deck to be in StandardDecks()")
	}
	if len(found.Notes) < 40 {
		t.Fatalf("travel adventure deck should have >=40 notes, got %d", len(found.Notes))
	}
	firstFront := found.Notes[0].Front
	if len(firstFront) < 2 {
		t.Fatalf("expected first note front to have content, got %q", firstFront)
	}
}

func TestPsychologyMindDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "psychology") && strings.Contains(d.ID, "mind") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected Psychology & Mind deck to be in StandardDecks()")
	}
	if len(found.Notes) < 35 {
		t.Fatalf("psychology deck should have >=35 notes, got %d", len(found.Notes))
	}
	firstFront := found.Notes[0].Front
	if len(firstFront) < 2 {
		t.Fatalf("expected first note front to have content, got %q", firstFront)
	}
}

func TestGrammarTipsCountMay15(t *testing.T) {
	if len(grammarTips) < 80 {
		t.Fatalf("expected >=80 grammar tips, got %d", len(grammarTips))
	}
}

func TestNewGrammarTipsMay15Present(t *testing.T) {
	titles := map[string]bool{}
	for _, tip := range grammarTips {
		titles[tip.Title] = true
	}
	expected := []string{
		"Präpositionale Verben",
		"Wechselpräpositionen",
		"Verben mit Dativ",
		"Modalverben",
	}
	for _, e := range expected {
		if !titles[e] {
			t.Errorf("expected grammar tip %q to be present", e)
		}
	}
}

func TestVerbsOfTheDayCountMay15(t *testing.T) {
	if len(dailyVerbs) < 130 {
		t.Fatalf("expected >=130 verbs of the day, got %d", len(dailyVerbs))
	}
}

func TestNewVerbsOfTheDayMay15Present(t *testing.T) {
	germans := map[string]bool{}
	for _, v := range dailyVerbs {
		germans[v.German] = true
	}
	expected := []string{"arbeiten", "kochen", "reisen", "schlafen", "schreiben", "hören", "verstehen", "machen"}
	for _, e := range expected {
		if !germans[e] {
			t.Errorf("expected verb %q to be present in dailyVerbs", e)
		}
	}
}

func TestTravelDeckCardsHaveProperStructure(t *testing.T) {
	decks := StandardDecks()
	var travelDeck *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "travel-adventure") {
			travelDeck = &d
			break
		}
	}
	if travelDeck == nil {
		t.Fatal("Travel Adventure deck not found")
	}
	for _, n := range travelDeck.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}

func TestPsychologyDeckCardsHaveProperStructure(t *testing.T) {
	decks := StandardDecks()
	var psychDeck *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "psychology-mind") {
			psychDeck = &d
			break
		}
	}
	if psychDeck == nil {
		t.Fatal("Psychology & Mind deck not found")
	}
	for _, n := range psychDeck.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}
