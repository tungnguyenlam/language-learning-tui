package content

import (
	"deutsch-tui/internal/core"
)

func B1SportsDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-spt-sport", DeckID: "b1-sports", Front: "der Sport", Back: "sport / physical exercise", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-sportler", DeckID: "b1-sports", Front: "der Sportler / die Sportlerin", Back: "athlete", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-mannschaft", DeckID: "b1-sports", Front: "die Mannschaft", Back: "team", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-trainer", DeckID: "b1-sports", Front: "der Trainer", Back: "coach / trainer", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-schiedsrichter", DeckID: "b1-sports", Front: "der Schiedsrichter", Back: "referee", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-spiel", DeckID: "b1-sports", Front: "das Spiel", Back: "game / match", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-wettkampf", DeckID: "b1-sports", Front: "der Wettkampf", Back: "competition", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-turnier", DeckID: "b1-sports", Front: "das Turnier", Back: "tournament", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-meister", DeckID: "b1-sports", Front: "der Meister", Back: "champion", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-sieg", DeckID: "b1-sports", Front: "der Sieg", Back: "victory", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-niederlage", DeckID: "b1-sports", Front: "die Niederlage", Back: "defeat", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-unentschieden", DeckID: "b1-sports", Front: "das Unentschieden", Back: "draw / tie", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-tor", DeckID: "b1-sports", Front: "das Tor", Back: "goal", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-punkt", DeckID: "b1-sports", Front: "der Punkt", Back: "point", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-rekord", DeckID: "b1-sports", Front: "der Rekord", Back: "record", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-medaille", DeckID: "b1-sports", Front: "die Medaille", Back: "medal", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-pokal", DeckID: "b1-sports", Front: "der Pokal", Back: "cup / trophy", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-stadion", DeckID: "b1-sports", Front: "das Stadion", Back: "stadium", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-halle", DeckID: "b1-sports", Front: "die Sporthalle", Back: "sports hall / gym", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-platz", DeckID: "b1-sports", Front: "der Sportplatz", Back: "sports field / pitch", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-fitnessstudio", DeckID: "b1-sports", Front: "das Fitnessstudio", Back: "gym / fitness center", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-ausruestung", DeckID: "b1-sports", Front: "die Ausrüstung", Back: "equipment", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-trainieren", DeckID: "b1-sports", Front: "trainieren", Back: "to train", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-gewinnen", DeckID: "b1-sports", Front: "gewinnen", Back: "to win", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-verlieren", DeckID: "b1-sports", Front: "verlieren", Back: "to lose", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-teilnehmen", DeckID: "b1-sports", Front: "teilnehmen (an)", Back: "to participate (in)", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-werfen", DeckID: "b1-sports", Front: "werfen", Back: "to throw", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-fangen", DeckID: "b1-sports", Front: "fangen", Back: "to catch", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-laufen", DeckID: "b1-sports", Front: "laufen", Back: "to run", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-rennen", DeckID: "b1-sports", Front: "rennen", Back: "to sprint / run fast", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-springen", DeckID: "b1-sports", Front: "springen", Back: "to jump", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-schwimmen", DeckID: "b1-sports", Front: "schwimmen", Back: "to swim", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-klettern", DeckID: "b1-sports", Front: "klettern", Back: "to climb", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-radfahren", DeckID: "b1-sports", Front: "Rad fahren", Back: "to cycle", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-foulen", DeckID: "b1-sports", Front: "foulen", Back: "to commit a foul", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-verletzen", DeckID: "b1-sports", Front: "sich verletzen", Back: "to get injured", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-muskel", DeckID: "b1-sports", Front: "der Muskel", Back: "muscle", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-ausdauer", DeckID: "b1-sports", Front: "die Ausdauer", Back: "stamina / endurance", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-kraft", DeckID: "b1-sports", Front: "die Kraft", Back: "strength / power", Tags: []string{"b1", "sports"}},
		{ID: "b1-spt-geschwindigkeit", DeckID: "b1-sports", Front: "die Geschwindigkeit", Back: "speed", Tags: []string{"b1", "sports"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-sports",
		Name:        "German B1 Sports & Fitness",
		Description: "Essential vocabulary for discussing sports, games, and fitness.",
		Tags:        []string{"german", "b1", "sports"},
		Notes:       notes,
	}
}
