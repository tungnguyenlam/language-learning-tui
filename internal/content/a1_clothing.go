package content

import (
	"deutsch-tui/internal/core"
)

func A1ClothingDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-clothing-kleidung", DeckID: "a1-clothing", Front: "die Kleidung", Back: "clothes / clothing", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-hemd", DeckID: "a1-clothing", Front: "das Hemd", Back: "shirt", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-t-shirt", DeckID: "a1-clothing", Front: "das T-Shirt", Back: "t-shirt", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-pullover", DeckID: "a1-clothing", Front: "der Pullover", Back: "sweater / pullover", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-hose", DeckID: "a1-clothing", Front: "die Hose", Back: "pants / trousers", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-jeans", DeckID: "a1-clothing", Front: "die Jeans", Back: "jeans", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-kleid", DeckID: "a1-clothing", Front: "das Kleid", Back: "dress", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-rock", DeckID: "a1-clothing", Front: "der Rock", Back: "skirt", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-jacke", DeckID: "a1-clothing", Front: "die Jacke", Back: "jacket", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-mantel", DeckID: "a1-clothing", Front: "der Mantel", Back: "coat", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-schuh", DeckID: "a1-clothing", Front: "der Schuh", Back: "shoe", Extra: "Plural: die Schuhe", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-stiefel", DeckID: "a1-clothing", Front: "der Stiefel", Back: "boot", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-socke", DeckID: "a1-clothing", Front: "die Socke", Back: "sock", Extra: "Plural: die Socken", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-hut", DeckID: "a1-clothing", Front: "der Hut", Back: "hat", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-muetze", DeckID: "a1-clothing", Front: "die Mütze", Back: "beanie / cap", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-schal", DeckID: "a1-clothing", Front: "der Schal", Back: "scarf", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-handschuh", DeckID: "a1-clothing", Front: "der Handschuh", Back: "glove", Extra: "Plural: die Handschuhe", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-guertel", DeckID: "a1-clothing", Front: "der Gürtel", Back: "belt", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-anzug", DeckID: "a1-clothing", Front: "der Anzug", Back: "suit", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-krawatte", DeckID: "a1-clothing", Front: "die Krawatte", Back: "tie", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-tasche", DeckID: "a1-clothing", Front: "die Tasche", Back: "bag", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-brille", DeckID: "a1-clothing", Front: "die Brille", Back: "glasses", Tags: []string{"a1", "clothing", "noun"}},
		{ID: "a1-clothing-anziehen", DeckID: "a1-clothing", Front: "anziehen", Back: "to put on (clothes)", Extra: "separable verb", Tags: []string{"a1", "clothing", "verb"}},
		{ID: "a1-clothing-ausziehen", DeckID: "a1-clothing", Front: "ausziehen", Back: "to take off (clothes)", Extra: "separable verb", Tags: []string{"a1", "clothing", "verb"}},
		{ID: "a1-clothing-tragen", DeckID: "a1-clothing", Front: "tragen", Back: "to wear", Tags: []string{"a1", "clothing", "verb"}},
		{ID: "a1-clothing-passen", DeckID: "a1-clothing", Front: "passen", Back: "to fit", Tags: []string{"a1", "clothing", "verb"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-clothing",
		Name:        "A1 Clothing & Accessories",
		Description: "Essential A1 vocabulary for clothes and getting dressed.",
		Tags:        []string{"german", "a1", "clothing"},
		Notes:       notes,
	}
}
