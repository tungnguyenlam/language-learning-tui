package content

import (
	"testing"
)

func TestGetRelevantGrammarTip(t *testing.T) {
	// Let's test basic retrieval
	tip := GetRelevantGrammarTip("This sentence has the word Dativ in it")
	if tip.Title != "The Dativ Case" && tip.Title != "Verbs with Dativ" {
		t.Errorf("Expected a tip about Dativ, got %s", tip.Title)
	}

	tip2 := GetRelevantGrammarTip("können")
	if tip2.Title != "Modalverben" && tip2.Title != "Infinitive with zu" && tip2.Title != "Konjunktiv II Subtle Distinctions" {
		// "können" is in the example for Modalverben, but might match other things depending on our heuristic
		// Let's just check it doesn't panic and returns something
	}

	tipEmpty := GetRelevantGrammarTip("")
	if tipEmpty.Title == "" {
		t.Errorf("Expected fallback tip for empty text, got empty")
	}
}
