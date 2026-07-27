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

func TestConsolidateDictionaryEntries(t *testing.T) {
	entries := []DictionaryEntry{
		{
			ID:          "1",
			Word:        "stark",
			Translation: "strong",
			WordClass:   "adj",
			Forms:       "stärker; am stärksten",
			Tags:        []string{"[general]"},
			Examples:    []string{"Ein starker Mann."},
		},
		{
			ID:          "2",
			Word:        "stark",
			Translation: "powerful",
			WordClass:   "adj",
			Forms:       "stärker",
			Tags:        []string{"[phys.]"},
			Examples:    []string{"Ein starkes Signal."},
		},
		{
			ID:          "3",
			Word:        "stark",
			Translation: "heavily",
			WordClass:   "adv",
			Forms:       "",
			Tags:        []string{"[coll.]"},
			Examples:    []string{"Es regnet stark."},
		},
		{
			ID:          "4",
			Word:        "Band",
			Gender:      "m",
			Translation: "volume; book",
		},
		{
			ID:          "5",
			Word:        "Band",
			Gender:      "f",
			Translation: "music band",
		},
	}

	merged := ConsolidateDictionaryEntries(entries)

	if len(merged) != 3 {
		t.Fatalf("expected 3 merged entries, got %d", len(merged))
	}

	// 1. "stark" merged entry
	stark := merged[0]
	if stark.Word != "stark" {
		t.Errorf("expected word 'stark', got '%s'", stark.Word)
	}
	expectedTrans := "strong; powerful; heavily"
	if stark.Translation != expectedTrans {
		t.Errorf("expected translation '%s', got '%s'", expectedTrans, stark.Translation)
	}
	if stark.WordClass != "adj, adv" {
		t.Errorf("expected WordClass 'adj, adv', got '%s'", stark.WordClass)
	}
	if stark.Forms != "stärker; am stärksten" {
		t.Errorf("expected Forms 'stärker; am stärksten', got '%s'", stark.Forms)
	}
	if len(stark.Tags) != 3 {
		t.Errorf("expected 3 tags, got %v", stark.Tags)
	}
	if len(stark.Examples) != 3 {
		t.Errorf("expected 3 examples, got %v", stark.Examples)
	}

	// 2. "Band" {m} entry
	bandM := merged[1]
	if bandM.Word != "Band" || bandM.Gender != "m" || bandM.Translation != "volume; book" {
		t.Errorf("unexpected entry for Band {m}: %v", bandM)
	}

	// 3. "Band" {f} entry
	bandF := merged[2]
	if bandF.Word != "Band" || bandF.Gender != "f" || bandF.Translation != "music band" {
		t.Errorf("unexpected entry for Band {f}: %v", bandF)
	}
}
