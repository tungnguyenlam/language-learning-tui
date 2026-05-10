package content

import "testing"

func TestUrbanMobilityDeckIsEmbedded(t *testing.T) {
	deck, err := DeckByID("b2_urban_mobility")
	if err != nil {
		t.Fatalf("DeckByID failed: %v", err)
	}
	if deck == nil {
		t.Fatal("b2_urban_mobility deck was not loaded")
	}
	if deck.Name != "b2-urban-mobility" {
		t.Fatalf("deck.Name = %q, want b2-urban-mobility", deck.Name)
	}
	if len(deck.Notes) < 25 {
		t.Fatalf("len(deck.Notes) = %d, want at least 25", len(deck.Notes))
	}

	var hasMCQ, hasCloze bool
	for _, note := range deck.Notes {
		if note.Type == "MCQ" {
			hasMCQ = true
		}
		if note.Type == "Cloze" {
			hasCloze = true
		}
	}
	if !hasMCQ || !hasCloze {
		t.Fatalf("deck should include both MCQ and Cloze notes; hasMCQ=%v hasCloze=%v", hasMCQ, hasCloze)
	}
}
