package content

import (
	"deutsch-tui/internal/core"
)

func B1CookingDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-co-kochen", DeckID: "b1-cooking", Front: "kochen", Back: "to cook / to boil", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-braten", DeckID: "b1-cooking", Front: "braten", Back: "to fry / to roast", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-backen", DeckID: "b1-cooking", Front: "backen", Back: "to bake", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-duensten", DeckID: "b1-cooking", Front: "dünsten", Back: "to steam / stew", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-grillen", DeckID: "b1-cooking", Front: "grillen", Back: "to grill", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-schneiden", DeckID: "b1-cooking", Front: "schneiden", Back: "to cut", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-hacken", DeckID: "b1-cooking", Front: "hacken", Back: "to chop / mince", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-ruehren", DeckID: "b1-cooking", Front: "rühren", Back: "to stir", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-schlagen", DeckID: "b1-cooking", Front: "schlagen", Back: "to whip / beat", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-mischen", DeckID: "b1-cooking", Front: "mischen", Back: "to mix", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-wuerzen", DeckID: "b1-cooking", Front: "würzen", Back: "to season", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-salzen", DeckID: "b1-cooking", Front: "salzen", Back: "to salt", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-pfeffern", DeckID: "b1-cooking", Front: "pfeffern", Back: "to pepper", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-schaelen", DeckID: "b1-cooking", Front: "schälen", Back: "to peel", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-reiben", DeckID: "b1-cooking", Front: "reiben", Back: "to grate", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-pfanne", DeckID: "b1-cooking", Front: "die Pfanne", Back: "frying pan", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-topf", DeckID: "b1-cooking", Front: "der Topf", Back: "pot", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-deckel", DeckID: "b1-cooking", Front: "der Deckel", Back: "lid", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-messer", DeckID: "b1-cooking", Front: "das Messer", Back: "knife", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-gabel", DeckID: "b1-cooking", Front: "die Gabel", Back: "fork", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-loeffel", DeckID: "b1-cooking", Front: "der Löffel", Back: "spoon", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-schneidebrett", DeckID: "b1-cooking", Front: "das Schneidebrett", Back: "cutting board", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-schuessel", DeckID: "b1-cooking", Front: "die Schüssel", Back: "bowl", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-teller", DeckID: "b1-cooking", Front: "der Teller", Back: "plate", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-tasse", DeckID: "b1-cooking", Front: "die Tasse", Back: "cup", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-glas", DeckID: "b1-cooking", Front: "das Glas", Back: "glass", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-rezept", DeckID: "b1-cooking", Front: "das Rezept", Back: "recipe", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-zutaten", DeckID: "b1-cooking", Front: "die Zutaten", Back: "ingredients", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-mehl", DeckID: "b1-cooking", Front: "das Mehl", Back: "flour", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-zucker", DeckID: "b1-cooking", Front: "der Zucker", Back: "sugar", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-salz", DeckID: "b1-cooking", Front: "das Salz", Back: "salt", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-pfeffer", DeckID: "b1-cooking", Front: "der Pfeffer", Back: "pepper", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-oel", DeckID: "b1-cooking", Front: "das Öl", Back: "oil", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-butter", DeckID: "b1-cooking", Front: "die Butter", Back: "butter", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-essig", DeckID: "b1-cooking", Front: "der Essig", Back: "vinegar", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-hefe", DeckID: "b1-cooking", Front: "die Hefe", Back: "yeast", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-teig", DeckID: "b1-cooking", Front: "der Teig", Back: "dough", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-portion", DeckID: "b1-cooking", Front: "die Portion", Back: "portion / serving", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-vorspeise", DeckID: "b1-cooking", Front: "die Vorspeise", Back: "starter / appetizer", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-hauptgang", DeckID: "b1-cooking", Front: "der Hauptgang", Back: "main course", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-nachtisch", DeckID: "b1-cooking", Front: "der Nachtisch", Back: "dessert", Tags: []string{"b1", "cooking"}},
		{ID: "b1-co-beilage", DeckID: "b1-cooking", Front: "die Beilage", Back: "side dish", Tags: []string{"b1", "cooking"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-cooking",
		Name:        "B1 German Cooking & Kitchen",
		Description: "Intermediate vocabulary for cooking actions, kitchen tools, ingredients, and meal courses.",
		Tags:        []string{"german", "b1", "cooking", "kitchen", "food"},
		Notes:       notes,
	}
}
