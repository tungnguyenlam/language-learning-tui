package content

import (
	"deutsch-tui/internal/core"
)

func A1FamilyDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-family-familie", DeckID: "a1-family", Front: "die Familie", Back: "family", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-eltern", DeckID: "a1-family", Front: "die Eltern", Back: "parents", Extra: "Plural only", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-vater", DeckID: "a1-family", Front: "der Vater", Back: "father", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-mutter", DeckID: "a1-family", Front: "die Mutter", Back: "mother", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-sohn", DeckID: "a1-family", Front: "der Sohn", Back: "son", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-tochter", DeckID: "a1-family", Front: "die Tochter", Back: "daughter", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-bruder", DeckID: "a1-family", Front: "der Bruder", Back: "brother", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-schwester", DeckID: "a1-family", Front: "die Schwester", Back: "sister", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-geschwister", DeckID: "a1-family", Front: "die Geschwister", Back: "siblings", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-grosseltern", DeckID: "a1-family", Front: "die Großeltern", Back: "grandparents", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-grossvater", DeckID: "a1-family", Front: "der Großvater / der Opa", Back: "grandfather", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-grossmutter", DeckID: "a1-family", Front: "die Großmutter / die Oma", Back: "grandmother", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-enkel", DeckID: "a1-family", Front: "der Enkel", Back: "grandson", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-enkelin", DeckID: "a1-family", Front: "die Enkelin", Back: "granddaughter", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-onkel", DeckID: "a1-family", Front: "der Onkel", Back: "uncle", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-tante", DeckID: "a1-family", Front: "die Tante", Back: "aunt", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-cousin", DeckID: "a1-family", Front: "der Cousin", Back: "cousin (male)", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-cousine", DeckID: "a1-family", Front: "die Cousine", Back: "cousin (female)", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-kind", DeckID: "a1-family", Front: "das Kind", Back: "child", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-baby", DeckID: "a1-family", Front: "das Baby", Back: "baby", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-verheiratet", DeckID: "a1-family", Front: "verheiratet", Back: "married", Tags: []string{"a1", "family", "adjective"}},
		{ID: "a1-family-ledig", DeckID: "a1-family", Front: "ledig", Back: "single (unmarried)", Tags: []string{"a1", "family", "adjective"}},
		{ID: "a1-family-geschieden", DeckID: "a1-family", Front: "geschieden", Back: "divorced", Tags: []string{"a1", "family", "adjective"}},
		{ID: "a1-family-freund", DeckID: "a1-family", Front: "der Freund", Back: "friend (male) / boyfriend", Tags: []string{"a1", "family", "noun"}},
		{ID: "a1-family-freundin", DeckID: "a1-family", Front: "die Freundin", Back: "friend (female) / girlfriend", Tags: []string{"a1", "family", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-family",
		Name:        "A1 Family & Friends",
		Description: "Essential A1 vocabulary for family members and relationship statuses.",
		Tags:        []string{"german", "a1", "family"},
		Notes:       notes,
	}
}
