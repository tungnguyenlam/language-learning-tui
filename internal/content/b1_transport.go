package content

import (
	"deutsch-tui/internal/core"
)

func B1TransportDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-trn-bahn", DeckID: "b1-transport", Front: "die Bahn", Back: "railway / train (system)", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-zug", DeckID: "b1-transport", Front: "der Zug", Back: "train", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-bahnhof", DeckID: "b1-transport", Front: "der Bahnhof", Back: "railway station", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-haltestelle", DeckID: "b1-transport", Front: "die Haltestelle", Back: "stop (bus / tram)", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-gleis", DeckID: "b1-transport", Front: "das Gleis", Back: "track / platform", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-fahrkarte", DeckID: "b1-transport", Front: "die Fahrkarte", Back: "ticket", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-fahrplan", DeckID: "b1-transport", Front: "der Fahrplan", Back: "timetable / schedule", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-verspaetung", DeckID: "b1-transport", Front: "die Verspätung", Back: "delay", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-anschluss", DeckID: "b1-transport", Front: "der Anschluss", Back: "connection (train)", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-umsteigen", DeckID: "b1-transport", Front: "umsteigen", Back: "to change (trains / buses)", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-aussteigen", DeckID: "b1-transport", Front: "aussteigen", Back: "to get off / disembark", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-einsteigen", DeckID: "b1-transport", Front: "einsteigen", Back: "to get on / board", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-strassenbahn", DeckID: "b1-transport", Front: "die Straßenbahn", Back: "tram / streetcar", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-ubahn", DeckID: "b1-transport", Front: "die U-Bahn", Back: "subway / underground", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-sbahn", DeckID: "b1-transport", Front: "die S-Bahn", Back: "city train / suburban rail", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-bus", DeckID: "b1-transport", Front: "der Bus", Back: "bus", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-taxi", DeckID: "b1-transport", Front: "das Taxi", Back: "taxi", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-fahrrad", DeckID: "b1-transport", Front: "das Fahrrad", Back: "bicycle", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-auto", DeckID: "b1-transport", Front: "das Auto", Back: "car", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-fluss", DeckID: "b1-transport", Front: "die Fähre", Back: "ferry", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-flughafen", DeckID: "b1-transport", Front: "der Flughafen", Back: "airport", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-flugzeug", DeckID: "b1-transport", Front: "das Flugzeug", Back: "aeroplane / aircraft", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-pendler", DeckID: "b1-transport", Front: "der Pendler", Back: "commuter", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-pendeln", DeckID: "b1-transport", Front: "pendeln", Back: "to commute", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-stau", DeckID: "b1-transport", Front: "der Stau", Back: "traffic jam", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-baustelle", DeckID: "b1-transport", Front: "die Baustelle", Back: "construction site / roadworks", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-tankstelle", DeckID: "b1-transport", Front: "die Tankstelle", Back: "petrol / gas station", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-parkplatz", DeckID: "b1-transport", Front: "der Parkplatz", Back: "car park / parking lot", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-autobahn", DeckID: "b1-transport", Front: "die Autobahn", Back: "motorway / highway", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-fahrt", DeckID: "b1-transport", Front: "die Fahrt", Back: "ride / journey", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-reise", DeckID: "b1-transport", Front: "die Reise", Back: "trip / journey", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-monatskarte", DeckID: "b1-transport", Front: "die Monatskarte", Back: "monthly travel pass", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-jahreskarte", DeckID: "b1-transport", Front: "die Jahreskarte", Back: "yearly travel pass", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-rueckfahrt", DeckID: "b1-transport", Front: "die Rückfahrt", Back: "return trip", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-hinfahrt", DeckID: "b1-transport", Front: "die Hinfahrt", Back: "outward journey", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-erstklasse", DeckID: "b1-transport", Front: "die erste Klasse", Back: "first class", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-zweiteklasse", DeckID: "b1-transport", Front: "die zweite Klasse", Back: "second class", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-schaffner", DeckID: "b1-transport", Front: "der Schaffner", Back: "ticket inspector / conductor", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-reisebuero", DeckID: "b1-transport", Front: "das Reisebüro", Back: "travel agency", Tags: []string{"b1", "transport"}},
		{ID: "b1-trn-buchung", DeckID: "b1-transport", Front: "die Buchung", Back: "booking", Tags: []string{"b1", "transport"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-transport",
		Name:        "B1 German Public Transport",
		Description: "Vocabulary for trains, trams, buses, commuting, and getting around in German-speaking cities.",
		Tags:        []string{"german", "b1", "transport", "travel"},
		Notes:       notes,
	}
}
