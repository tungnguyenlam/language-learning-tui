package content

import (
	"deutsch-tui/internal/core"
)

func A1AnimalsDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-an-hund", DeckID: "a1-animals", Front: "der Hund", Back: "dog", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-katze", DeckID: "a1-animals", Front: "die Katze", Back: "cat", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-pferd", DeckID: "a1-animals", Front: "das Pferd", Back: "horse", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-kuh", DeckID: "a1-animals", Front: "die Kuh", Back: "cow", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-schwein", DeckID: "a1-animals", Front: "das Schwein", Back: "pig", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-schaf", DeckID: "a1-animals", Front: "das Schaf", Back: "sheep", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-ziege", DeckID: "a1-animals", Front: "die Ziege", Back: "goat", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-huhn", DeckID: "a1-animals", Front: "das Huhn", Back: "chicken", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-hahn", DeckID: "a1-animals", Front: "der Hahn", Back: "rooster", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-ente", DeckID: "a1-animals", Front: "die Ente", Back: "duck", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-gans", DeckID: "a1-animals", Front: "die Gans", Back: "goose", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-vogel", DeckID: "a1-animals", Front: "der Vogel", Back: "bird", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-fisch", DeckID: "a1-animals", Front: "der Fisch", Back: "fish", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-maus", DeckID: "a1-animals", Front: "die Maus", Back: "mouse", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-ratte", DeckID: "a1-animals", Front: "die Ratte", Back: "rat", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-hase", DeckID: "a1-animals", Front: "der Hase", Back: "hare / rabbit", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-kaninchen", DeckID: "a1-animals", Front: "das Kaninchen", Back: "rabbit", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-baer", DeckID: "a1-animals", Front: "der Bär", Back: "bear", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-wolf", DeckID: "a1-animals", Front: "der Wolf", Back: "wolf", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-fuchs", DeckID: "a1-animals", Front: "der Fuchs", Back: "fox", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-hirsch", DeckID: "a1-animals", Front: "der Hirsch", Back: "deer / stag", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-eichhoern", DeckID: "a1-animals", Front: "das Eichhörnchen", Back: "squirrel", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-igel", DeckID: "a1-animals", Front: "der Igel", Back: "hedgehog", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-loewe", DeckID: "a1-animals", Front: "der Löwe", Back: "lion", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-tiger", DeckID: "a1-animals", Front: "der Tiger", Back: "tiger", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-elefant", DeckID: "a1-animals", Front: "der Elefant", Back: "elephant", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-affe", DeckID: "a1-animals", Front: "der Affe", Back: "monkey / ape", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-giraffe", DeckID: "a1-animals", Front: "die Giraffe", Back: "giraffe", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-zebra", DeckID: "a1-animals", Front: "das Zebra", Back: "zebra", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-schlange", DeckID: "a1-animals", Front: "die Schlange", Back: "snake", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-frosch", DeckID: "a1-animals", Front: "der Frosch", Back: "frog", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-spinne", DeckID: "a1-animals", Front: "die Spinne", Back: "spider", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-biene", DeckID: "a1-animals", Front: "die Biene", Back: "bee", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-wespe", DeckID: "a1-animals", Front: "die Wespe", Back: "wasp", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-fliege", DeckID: "a1-animals", Front: "die Fliege", Back: "fly", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-schmetterling", DeckID: "a1-animals", Front: "der Schmetterling", Back: "butterfly", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-ameise", DeckID: "a1-animals", Front: "die Ameise", Back: "ant", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-delfin", DeckID: "a1-animals", Front: "der Delfin", Back: "dolphin", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-wal", DeckID: "a1-animals", Front: "der Wal", Back: "whale", Tags: []string{"a1", "animals"}},
		{ID: "a1-an-haifisch", DeckID: "a1-animals", Front: "der Haifisch", Back: "shark", Tags: []string{"a1", "animals"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-animals",
		Name:        "A1 German Animals",
		Description: "Beginner vocabulary covering pets, farm, wild, and sea animals.",
		Tags:        []string{"german", "a1", "animals"},
		Notes:       notes,
	}
}
