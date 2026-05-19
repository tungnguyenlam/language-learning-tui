package content

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"
)

func TestLegalContractsDeckLoaded(t *testing.T) {
	decks := StandardDecks()
	var found *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "legal") || strings.Contains(d.Name, "Legal") {
			found = &d
			break
		}
	}
	if found == nil {
		t.Fatal("expected Legal & Contracts deck to be in StandardDecks()")
	}
	if len(found.Notes) < 40 {
		t.Fatalf("legal deck should have >=40 notes, got %d", len(found.Notes))
	}
	if len(found.Notes) == 0 {
		t.Fatal("expected legal deck to have notes")
	}
	firstFront := found.Notes[0].Front
	if len(firstFront) < 2 {
		t.Fatalf("expected first note front to have content, got %q", firstFront)
	}
}

func TestGrammarTipsCountMay14c(t *testing.T) {
	if len(grammarTips) < 80 {
		t.Fatalf("expected >=80 grammar tips, got %d", len(grammarTips))
	}
}

func TestNewGrammarTipsMay14cPresent(t *testing.T) {
	titles := map[string]bool{}
	for _, tip := range grammarTips {
		titles[tip.Title] = true
	}
	expected := []string{
		"Weglassen des Artikels",
		"Partizip II mit 'zu'",
		"Futur I",
		"Plusquamperfekt",
		"Konjunktiv II Past",
	}
	for _, e := range expected {
		if !titles[e] {
			t.Errorf("expected grammar tip %q to be present", e)
		}
	}
}

func TestVerbsOfTheDayCountMay14c(t *testing.T) {
	if len(dailyVerbs) < 120 {
		t.Fatalf("expected >=120 verbs of the day, got %d", len(dailyVerbs))
	}
}

func TestNewVerbsOfTheDayMay14cPresent(t *testing.T) {
	germans := map[string]bool{}
	for _, v := range dailyVerbs {
		germans[v.German] = true
	}
	expected := []string{"sagen", "geben", "nehmen", "kommen", "werden", "finden", "stehen", "sitzen"}
	for _, e := range expected {
		if !germans[e] {
			t.Errorf("expected verb %q to be present in dailyVerbs", e)
		}
	}
}

func TestLegalDeckCardsHaveProperStructure(t *testing.T) {
	decks := StandardDecks()
	var legalDeck *core.Deck
	for _, d := range decks {
		if strings.Contains(d.ID, "legal") {
			legalDeck = &d
			break
		}
	}
	if legalDeck == nil {
		t.Fatal("Legal deck not found")
	}
	for _, n := range legalDeck.Notes {
		if len(n.Cards) == 0 {
			t.Errorf("note %s has no cards", n.Front)
		}
	}
}
