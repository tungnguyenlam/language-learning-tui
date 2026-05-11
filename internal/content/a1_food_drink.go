package content

import (
	"deutsch-tui/internal/core"
)

func A1FoodDrinkDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-food-essen", DeckID: "a1-food-drink", Front: "das Essen", Back: "food", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-trinken", DeckID: "a1-food-drink", Front: "das Getränk", Back: "drink / beverage", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-fruehstueck", DeckID: "a1-food-drink", Front: "das Frühstück", Back: "breakfast", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-mittagessen", DeckID: "a1-food-drink", Front: "das Mittagessen", Back: "lunch", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-abendessen", DeckID: "a1-food-drink", Front: "das Abendessen", Back: "dinner", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-brot", DeckID: "a1-food-drink", Front: "das Brot", Back: "bread", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-broetchen", DeckID: "a1-food-drink", Front: "das Brötchen", Back: "bread roll", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-butter", DeckID: "a1-food-drink", Front: "die Butter", Back: "butter", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-kaese", DeckID: "a1-food-drink", Front: "der Käse", Back: "cheese", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-wurst", DeckID: "a1-food-drink", Front: "die Wurst", Back: "sausage", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-schinken", DeckID: "a1-food-drink", Front: "der Schinken", Back: "ham", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-ei", DeckID: "a1-food-drink", Front: "das Ei", Back: "egg", Extra: "Plural: die Eier", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-marmelade", DeckID: "a1-food-drink", Front: "die Marmelade", Back: "jam", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-honig", DeckID: "a1-food-drink", Front: "der Honig", Back: "honey", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-obst", DeckID: "a1-food-drink", Front: "das Obst", Back: "fruit", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-apfel", DeckID: "a1-food-drink", Front: "der Apfel", Back: "apple", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-banane", DeckID: "a1-food-drink", Front: "die Banane", Back: "banana", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-orange", DeckID: "a1-food-drink", Front: "die Orange", Back: "orange", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-gemuese", DeckID: "a1-food-drink", Front: "das Gemüse", Back: "vegetables", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-tomate", DeckID: "a1-food-drink", Front: "die Tomate", Back: "tomato", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-kartoffel", DeckID: "a1-food-drink", Front: "die Kartoffel", Back: "potato", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-zwiebel", DeckID: "a1-food-drink", Front: "die Zwiebel", Back: "onion", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-fleisch", DeckID: "a1-food-drink", Front: "das Fleisch", Back: "meat", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-hahnchen", DeckID: "a1-food-drink", Front: "das Hähnchen", Back: "chicken", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-rindfleisch", DeckID: "a1-food-drink", Front: "das Rindfleisch", Back: "beef", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-schweinefleisch", DeckID: "a1-food-drink", Front: "das Schweinefleisch", Back: "pork", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-fisch", DeckID: "a1-food-drink", Front: "der Fisch", Back: "fish", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-suppe", DeckID: "a1-food-drink", Front: "die Suppe", Back: "soup", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-salat", DeckID: "a1-food-drink", Front: "der Salat", Back: "salad", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-reis", DeckID: "a1-food-drink", Front: "der Reis", Back: "rice", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-nudeln", DeckID: "a1-food-drink", Front: "die Nudeln", Back: "noodles / pasta", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-kuchen", DeckID: "a1-food-drink", Front: "der Kuchen", Back: "cake", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-eis", DeckID: "a1-food-drink", Front: "das Eis", Back: "ice cream", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-wasser", DeckID: "a1-food-drink", Front: "das Wasser", Back: "water", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-mineralwasser", DeckID: "a1-food-drink", Front: "das Mineralwasser", Back: "mineral water", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-saft", DeckID: "a1-food-drink", Front: "der Saft", Back: "juice", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-kaffee", DeckID: "a1-food-drink", Front: "der Kaffee", Back: "coffee", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-tee", DeckID: "a1-food-drink", Front: "der Tee", Back: "tea", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-milch", DeckID: "a1-food-drink", Front: "die Milch", Back: "milk", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-bier", DeckID: "a1-food-drink", Front: "das Bier", Back: "beer", Tags: []string{"a1", "food", "noun"}},
		{ID: "a1-food-wein", DeckID: "a1-food-drink", Front: "der Wein", Back: "wine", Tags: []string{"a1", "food", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-food-drink",
		Name:        "A1 Food & Drink",
		Description: "Essential A1 vocabulary for food, drinks, and meals.",
		Tags:        []string{"german", "a1", "food"},
		Notes:       notes,
	}
}
