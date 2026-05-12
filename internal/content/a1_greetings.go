package content

import (
	"deutsch-tui/internal/core"
)

func A1GreetingsDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-greet-hallo", DeckID: "a1-greetings", Front: "Hallo", Back: "Hello", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-morgen", DeckID: "a1-greetings", Front: "Guten Morgen", Back: "Good morning", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-tag", DeckID: "a1-greetings", Front: "Guten Tag", Back: "Good day / Hello", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-abend", DeckID: "a1-greetings", Front: "Guten Abend", Back: "Good evening", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-nacht", DeckID: "a1-greetings", Front: "Gute Nacht", Back: "Good night", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-tschuss", DeckID: "a1-greetings", Front: "Tschüss", Back: "Bye", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-aufwiedersehen", DeckID: "a1-greetings", Front: "Auf Wiedersehen", Back: "Goodbye", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-danke", DeckID: "a1-greetings", Front: "Danke", Back: "Thank you", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-bitte", DeckID: "a1-greetings", Front: "Bitte", Back: "Please / You're welcome", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-entschuldigung", DeckID: "a1-greetings", Front: "Entschuldigung", Back: "Excuse me / Sorry", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-wie-gehts", DeckID: "a1-greetings", Front: "Wie geht's?", Back: "How are you?", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-mir-gehts-gut", DeckID: "a1-greetings", Front: "Mir geht's gut.", Back: "I'm fine.", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-und-dir", DeckID: "a1-greetings", Front: "Und dir?", Back: "And you? (informal)", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-und-ihnen", DeckID: "a1-greetings", Front: "Und Ihnen?", Back: "And you? (formal)", Tags: []string{"a1", "greetings"}},
		{ID: "a1-greet-freut-mich", DeckID: "a1-greetings", Front: "Freut mich", Back: "Nice to meet you", Tags: []string{"a1", "greetings"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-greetings",
		Name:        "A1 Greetings & Farewells",
		Description: "Essential A1 vocabulary for saying hello and goodbye.",
		Tags:        []string{"german", "a1", "greetings"},
		Notes:       notes,
	}
}
