package core

import "testing"

func TestValidateCard(t *testing.T) {
	card := Card{
		ID:      "card-1",
		NoteID:  "note-1",
		Kind:    CardKindMCQ,
		Prompt:  "der Apfel",
		Answer:  "apple",
		Choices: []string{"apple", "bread"},
	}
	if err := ValidateCard(card); err != nil {
		t.Fatalf("expected valid card: %v", err)
	}

	card.Choices = []string{"apple"}
	if err := ValidateCard(card); err == nil {
		t.Fatal("expected invalid mcq with one choice")
	}
}
