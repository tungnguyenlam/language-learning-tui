package content

import (
	"deutsch-tui/internal/core"
)

func A1CityDirectionsDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-city-stadt", DeckID: "a1-city-directions", Front: "die Stadt", Back: "city / town", Tags: []string{"a1", "city"}},
		{ID: "a1-city-zentrum", DeckID: "a1-city-directions", Front: "das Zentrum / die Stadtmitte", Back: "city center", Tags: []string{"a1", "city"}},
		{ID: "a1-city-strasse", DeckID: "a1-city-directions", Front: "die Straße", Back: "street", Tags: []string{"a1", "city"}},
		{ID: "a1-city-platz", DeckID: "a1-city-directions", Front: "der Platz", Back: "square / place", Tags: []string{"a1", "city"}},
		{ID: "a1-city-weg", DeckID: "a1-city-directions", Front: "der Weg", Back: "way / path", Tags: []string{"a1", "city"}},
		{ID: "a1-city-kreuzung", DeckID: "a1-city-directions", Front: "die Kreuzung", Back: "intersection / crossing", Tags: []string{"a1", "city"}},
		{ID: "a1-city-ampel", DeckID: "a1-city-directions", Front: "die Ampel", Back: "traffic light", Tags: []string{"a1", "city"}},
		{ID: "a1-city-ecke", DeckID: "a1-city-directions", Front: "die Ecke", Back: "corner", Tags: []string{"a1", "city"}},
		{ID: "a1-city-bruecke", DeckID: "a1-city-directions", Front: "die Brücke", Back: "bridge", Tags: []string{"a1", "city"}},
		{ID: "a1-city-bahnhof", DeckID: "a1-city-directions", Front: "der Bahnhof", Back: "train station", Tags: []string{"a1", "city"}},
		{ID: "a1-city-haltestelle", DeckID: "a1-city-directions", Front: "die Haltestelle", Back: "stop (bus/tram)", Tags: []string{"a1", "city"}},
		{ID: "a1-city-flughafen", DeckID: "a1-city-directions", Front: "der Flughafen", Back: "airport", Tags: []string{"a1", "city"}},
		{ID: "a1-city-hotel", DeckID: "a1-city-directions", Front: "das Hotel", Back: "hotel", Tags: []string{"a1", "city"}},
		{ID: "a1-city-restaurant", DeckID: "a1-city-directions", Front: "das Restaurant", Back: "restaurant", Tags: []string{"a1", "city"}},
		{ID: "a1-city-museum", DeckID: "a1-city-directions", Front: "das Museum", Back: "museum", Tags: []string{"a1", "city"}},
		{ID: "a1-city-kirche", DeckID: "a1-city-directions", Front: "die Kirche", Back: "church", Tags: []string{"a1", "city"}},
		{ID: "a1-city-park", DeckID: "a1-city-directions", Front: "der Park", Back: "park", Tags: []string{"a1", "city"}},
		{ID: "a1-city-supermarkt", DeckID: "a1-city-directions", Front: "der Supermarkt", Back: "supermarket", Tags: []string{"a1", "city"}},
		{ID: "a1-city-apotheke", DeckID: "a1-city-directions", Front: "die Apotheke", Back: "pharmacy", Tags: []string{"a1", "city"}},
		{ID: "a1-city-bank", DeckID: "a1-city-directions", Front: "die Bank", Back: "bank", Tags: []string{"a1", "city"}},
		{ID: "a1-city-post", DeckID: "a1-city-directions", Front: "die Post", Back: "post office", Tags: []string{"a1", "city"}},
		{ID: "a1-city-polizei", DeckID: "a1-city-directions", Front: "die Polizei", Back: "police", Tags: []string{"a1", "city"}},
		{ID: "a1-city-krankenhaus", DeckID: "a1-city-directions", Front: "das Krankenhaus", Back: "hospital", Tags: []string{"a1", "city"}},
		{ID: "a1-city-tourist-info", DeckID: "a1-city-directions", Front: "die Touristeninformation", Back: "tourist information", Tags: []string{"a1", "city"}},
		{ID: "a1-city-links", DeckID: "a1-city-directions", Front: "links", Back: "left", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-rechts", DeckID: "a1-city-directions", Front: "rechts", Back: "right", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-geradeaus", DeckID: "a1-city-directions", Front: "geradeaus", Back: "straight ahead", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-zurueck", DeckID: "a1-city-directions", Front: "zurück", Back: "back", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-oben", DeckID: "a1-city-directions", Front: "oben", Back: "up / above", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-unten", DeckID: "a1-city-directions", Front: "unten", Back: "down / below", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-nah", DeckID: "a1-city-directions", Front: "nah", Back: "near", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-weit", DeckID: "a1-city-directions", Front: "weit", Back: "far", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-gegenueber", DeckID: "a1-city-directions", Front: "gegenüber", Back: "opposite", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-neben", DeckID: "a1-city-directions", Front: "neben", Back: "next to", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-zwischen", DeckID: "a1-city-directions", Front: "zwischen", Back: "between", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-abbiegen", DeckID: "a1-city-directions", Front: "abbiegen", Back: "to turn", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-finden", DeckID: "a1-city-directions", Front: "finden", Back: "to find", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-suchen", DeckID: "a1-city-directions", Front: "suchen", Back: "to search / to look for", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-entschuldigung", DeckID: "a1-city-directions", Front: "Entschuldigung, ...", Back: "Excuse me, ...", Tags: []string{"a1", "directions"}},
		{ID: "a1-city-hilfe", DeckID: "a1-city-directions", Front: "Können Sie mir helfen?", Back: "Can you help me?", Tags: []string{"a1", "directions"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-city-directions",
		Name:        "A1 German City & Directions",
		Description: "Essential vocabulary for navigating a city and asking for directions in German.",
		Tags:        []string{"german", "a1", "city", "directions", "travel"},
		Notes:       notes,
	}
}
