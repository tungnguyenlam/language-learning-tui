package content

import (
	"deutsch-tui/internal/core"
)

func B2CultureLeisureDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-cult-kultur", DeckID: "b2-culture-leisure", Front: "die Kultur", Back: "culture", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-freizeit", DeckID: "b2-culture-leisure", Front: "die Freizeit", Back: "leisure time / free time", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-hobby", DeckID: "b2-culture-leisure", Front: "das Hobby", Back: "hobby", Tags: []string{"a1", "culture", "noun"}},
		{ID: "b2-cult-interesse", DeckID: "b2-culture-leisure", Front: "das Interesse", Back: "interest", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-museum", DeckID: "b2-culture-leisure", Front: "das Museum", Back: "museum", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-ausstellung", DeckID: "b2-culture-leisure", Front: "die Ausstellung", Back: "exhibition", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-galerie", DeckID: "b2-culture-leisure", Front: "die Galerie", Back: "gallery", Tags: []string{"b2", "culture", "noun"}},
		{ID: "b2-cult-theater", DeckID: "b2-culture-leisure", Front: "das Theater", Back: "theatre", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-vorstellung", DeckID: "b2-culture-leisure", Front: "die Vorstellung", Back: "performance / show", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-bühne", DeckID: "b2-culture-leisure", Front: "die Bühne", Back: "stage", Tags: []string{"b2", "culture", "noun"}},
		{ID: "b2-cult-konzert", DeckID: "b2-culture-leisure", Front: "das Konzert", Back: "concert", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-orchester", DeckID: "b2-culture-leisure", Front: "das Orchester", Back: "orchestra", Tags: []string{"b2", "culture", "noun"}},
		{ID: "b2-cult-klassik", DeckID: "b2-culture-leisure", Front: "die klassische Musik", Back: "classical music", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-pop", DeckID: "b2-culture-leisure", Front: "der Pop", Back: "pop music", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-jazz", DeckID: "b2-culture-leisure", Front: "der Jazz", Back: "jazz", Tags: []string{"b2", "culture", "noun"}},
		{ID: "b2-cult-electronic", DeckID: "b2-culture-leisure", Front: "die elektronische Musik", Back: "electronic music", Tags: []string{"b2", "culture", "noun"}},
		{ID: "b2-cult-festival", DeckID: "b2-culture-leisure", Front: "das Festival", Back: "festival", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-feier", DeckID: "b2-culture-leisure", Front: "die Feier", Back: "celebration / party", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-veranstaltung", DeckID: "b2-culture-leisure", Front: "die Veranstaltung", Back: "event", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-treffen", DeckID: "b2-culture-leisure", Front: "das Treffen", Back: "meeting / gathering", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-spaziergang", DeckID: "b2-culture-leisure", Front: "der Spaziergang", Back: "walk", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-wanderung", DeckID: "b2-culture-leisure", Front: "die Wanderung", Back: "hike", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-radfahren", DeckID: "b2-culture-leisure", Front: "das Radfahren", Back: "cycling", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-schwimmen", DeckID: "b2-culture-leisure", Front: "das Schwimmen", Back: "swimming", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-kochen", DeckID: "b2-culture-leisure", Front: "das Kochen", Back: "cooking", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-lesen", DeckID: "b2-culture-leisure", Front: "das Lesen", Back: "reading", Tags: []string{"a1", "culture", "noun"}},
		{ID: "b2-cult-schreiben", DeckID: "b2-culture-leisure", Front: "das Schreiben", Back: "writing", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-film", DeckID: "b2-culture-leisure", Front: "der Film", Back: "film / movie", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-kino", DeckID: "b2-culture-leisure", Front: "das Kino", Back: "cinema", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-podcast", DeckID: "b2-culture-leisure", Front: "der Podcast", Back: "podcast", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-streaming", DeckID: "b2-culture-leisure", Front: "das Streaming", Back: "streaming", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-brett", DeckID: "b2-culture-leisure", Front: "das Brettspiel", Back: "board game", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-karten", DeckID: "b2-culture-leisure", Front: "die Spielkarten", Back: "playing cards", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-urlaub", DeckID: "b2-culture-leisure", Front: "der Urlaub", Back: "vacation / holiday", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-reise", DeckID: "b2-culture-leisure", Front: "die Reise", Back: "trip / journey", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-ausflug", DeckID: "b2-culture-leisure", Front: "der Ausflug", Back: "excursion / day trip", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-park", DeckID: "b2-culture-leisure", Front: "der Park", Back: "park", Tags: []string{"a1", "culture", "noun"}},
		{ID: "b2-cult-garten", DeckID: "b2-culture-leisure", Front: "der Garten", Back: "garden", Tags: []string{"a2", "culture", "noun"}},
		{ID: "b2-cult-fitness", DeckID: "b2-culture-leisure", Front: "das Fitnessstudio", Back: "gym", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-yoga", DeckID: "b2-culture-leisure", Front: "das Yoga", Back: "yoga", Tags: []string{"b1", "culture", "noun"}},
		{ID: "b2-cult-meditation", DeckID: "b2-culture-leisure", Front: "die Meditation", Back: "meditation", Tags: []string{"b2", "culture", "noun"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "b2-culture-leisure",
		Name:        "German B2 Culture & Leisure",
		Description: "Cultural activities, hobbies, entertainment, and leisure vocabulary for B2 learners.",
		Tags:        []string{"german", "b2", "culture", "leisure"},
		Notes:       notes,
	}
}
