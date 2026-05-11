package content

import (
	"deutsch-tui/internal/core"
)

func A1HobbiesDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-hob-hobby", DeckID: "a1-hobbies", Front: "das Hobby", Back: "hobby", Extra: "Plural: die Hobbys", Tags: []string{"a1", "hobbies", "noun"}},
		{ID: "a1-hob-lesen", DeckID: "a1-hobbies", Front: "lesen", Back: "to read", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-buch", DeckID: "a1-hobbies", Front: "das Buch", Back: "book", Extra: "Plural: die Bücher", Tags: []string{"a1", "hobbies", "noun"}},
		{ID: "a1-hob-musik", DeckID: "a1-hobbies", Front: "Musik hören", Back: "to listen to music", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-singen", DeckID: "a1-hobbies", Front: "singen", Back: "to sing", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-tanzen", DeckID: "a1-hobbies", Front: "tanzen", Back: "to dance", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-schwimmen", DeckID: "a1-hobbies", Front: "schwimmen", Back: "to swim", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-spielen", DeckID: "a1-hobbies", Front: "spielen", Back: "to play", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-fussball", DeckID: "a1-hobbies", Front: "Fußball spielen", Back: "to play soccer", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-kino", DeckID: "a1-hobbies", Front: "ins Kino gehen", Back: "to go to the cinema", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-kochen", DeckID: "a1-hobbies", Front: "kochen", Back: "to cook", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-backen", DeckID: "a1-hobbies", Front: "backen", Back: "to bake", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-reisen", DeckID: "a1-hobbies", Front: "reisen", Back: "to travel", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-fotografieren", DeckID: "a1-hobbies", Front: "fotografieren", Back: "to take photos", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-malen", DeckID: "a1-hobbies", Front: "malen", Back: "to paint", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-zeichnen", DeckID: "a1-hobbies", Front: "zeichnen", Back: "to draw", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-wandern", DeckID: "a1-hobbies", Front: "wandern", Back: "to hike", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-laufen", DeckID: "a1-hobbies", Front: "laufen / joggen", Back: "to run / to jog", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-radfahren", DeckID: "a1-hobbies", Front: "Rad fahren", Back: "to cycle", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-fernsehen", DeckID: "a1-hobbies", Front: "fernsehen", Back: "to watch TV", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-freunde", DeckID: "a1-hobbies", Front: "Freunde treffen", Back: "to meet friends", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-shoppen", DeckID: "a1-hobbies", Front: "einkaufen / shoppen", Back: "to shop", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-schlafen", DeckID: "a1-hobbies", Front: "schlafen", Back: "to sleep", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-garten", DeckID: "a1-hobbies", Front: "die Gartenarbeit", Back: "gardening", Tags: []string{"a1", "hobbies", "noun"}},
		{ID: "a1-hob-instrument", DeckID: "a1-hobbies", Front: "ein Instrument spielen", Back: "to play an instrument", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-klavier", DeckID: "a1-hobbies", Front: "das Klavier", Back: "piano", Tags: []string{"a1", "hobbies", "noun"}},
		{ID: "a1-hob-gitarre", DeckID: "a1-hobbies", Front: "die Gitarre", Back: "guitar", Tags: []string{"a1", "hobbies", "noun"}},
		{ID: "a1-hob-videospiele", DeckID: "a1-hobbies", Front: "Videospiele spielen", Back: "to play video games", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-faulenzen", DeckID: "a1-hobbies", Front: "faulenzen", Back: "to laze around", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-sport", DeckID: "a1-hobbies", Front: "Sport machen", Back: "to do sports", Tags: []string{"a1", "hobbies", "phrase"}},
		{ID: "a1-hob-reiten", DeckID: "a1-hobbies", Front: "reiten", Back: "to ride (a horse)", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-angeln", DeckID: "a1-hobbies", Front: "angeln", Back: "to fish", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-schreiben", DeckID: "a1-hobbies", Front: "schreiben", Back: "to write", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-sammeln", DeckID: "a1-hobbies", Front: "sammeln", Back: "to collect", Tags: []string{"a1", "hobbies", "verb"}},
		{ID: "a1-hob-basteln", DeckID: "a1-hobbies", Front: "basteln", Back: "to do crafts", Tags: []string{"a1", "hobbies", "verb"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-hobbies",
		Name:        "A1 Hobbies & Free Time",
		Description: "Essential A1 vocabulary for talking about hobbies and free time activities.",
		Tags:        []string{"german", "a1", "hobbies"},
		Notes:       notes,
	}
}
