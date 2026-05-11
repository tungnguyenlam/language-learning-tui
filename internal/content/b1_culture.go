package content

import (
	"deutsch-tui/internal/core"
)

func B1CultureDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-cul-kultur", DeckID: "b1-culture", Front: "die Kultur", Back: "culture", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-tradition", DeckID: "b1-culture", Front: "die Tradition", Back: "tradition", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-brauch", DeckID: "b1-culture", Front: "der Brauch", Back: "custom / practice", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-fest", DeckID: "b1-culture", Front: "das Fest", Back: "festival / celebration", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-feiertag", DeckID: "b1-culture", Front: "der Feiertag", Back: "public holiday", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-veranstaltung", DeckID: "b1-culture", Front: "die Veranstaltung", Back: "event", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-ausstellung", DeckID: "b1-culture", Front: "die Ausstellung", Back: "exhibition", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-museum", DeckID: "b1-culture", Front: "das Museum", Back: "museum", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-kunst", DeckID: "b1-culture", Front: "die Kunst", Back: "art", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-kuenstler", DeckID: "b1-culture", Front: "der Künstler", Back: "artist", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gemaelde", DeckID: "b1-culture", Front: "das Gemälde", Back: "painting", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-theater", DeckID: "b1-culture", Front: "das Theater", Back: "theatre", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-buehne", DeckID: "b1-culture", Front: "die Bühne", Back: "stage", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-schauspieler", DeckID: "b1-culture", Front: "der Schauspieler", Back: "actor", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-publikum", DeckID: "b1-culture", Front: "das Publikum", Back: "audience", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-konzert", DeckID: "b1-culture", Front: "das Konzert", Back: "concert", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-orchester", DeckID: "b1-culture", Front: "das Orchester", Back: "orchestra", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-literatur", DeckID: "b1-culture", Front: "die Literatur", Back: "literature", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-schriftsteller", DeckID: "b1-culture", Front: "der Schriftsteller", Back: "writer", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-roman", DeckID: "b1-culture", Front: "der Roman", Back: "novel", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gedicht", DeckID: "b1-culture", Front: "das Gedicht", Back: "poem", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-maerchen", DeckID: "b1-culture", Front: "das Märchen", Back: "fairy tale", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-geschichte", DeckID: "b1-culture", Front: "die Geschichte", Back: "history / story", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-vergangenheit", DeckID: "b1-culture", Front: "die Vergangenheit", Back: "past", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gegenwart", DeckID: "b1-culture", Front: "die Gegenwart", Back: "present", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-zukunft", DeckID: "b1-culture", Front: "die Zukunft", Back: "future", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-erbe", DeckID: "b1-culture", Front: "das Erbe", Back: "heritage", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-denkmal", DeckID: "b1-culture", Front: "das Denkmal", Back: "monument / memorial", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gesellschaft", DeckID: "b1-culture", Front: "die Gesellschaft", Back: "society", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gemeinschaft", DeckID: "b1-culture", Front: "die Gemeinschaft", Back: "community", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-buerger", DeckID: "b1-culture", Front: "der Bürger", Back: "citizen", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-staat", DeckID: "b1-culture", Front: "der Staat", Back: "state", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-regierung", DeckID: "b1-culture", Front: "die Regierung", Back: "government", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-wahl", DeckID: "b1-culture", Front: "die Wahl", Back: "election / choice", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-demokratie", DeckID: "b1-culture", Front: "die Demokratie", Back: "democracy", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-recht", DeckID: "b1-culture", Front: "das Recht", Back: "right / law", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-pflicht", DeckID: "b1-culture", Front: "die Pflicht", Back: "duty / obligation", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-gesetz", DeckID: "b1-culture", Front: "das Gesetz", Back: "law", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-religion", DeckID: "b1-culture", Front: "die Religion", Back: "religion", Tags: []string{"b1", "culture"}},
		{ID: "b1-cul-glaube", DeckID: "b1-culture", Front: "der Glaube", Back: "belief / faith", Tags: []string{"b1", "culture"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-culture",
		Name:        "German B1 Culture & Society",
		Description: "Essential vocabulary for discussing culture, arts, history, and society.",
		Tags:        []string{"german", "b1", "culture"},
		Notes:       notes,
	}
}
