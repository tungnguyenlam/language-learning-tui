package content

import (
	"deutsch-tui/internal/core"
)

func A1TravelDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-tra-reise", DeckID: "a1-travel", Front: "die Reise", Back: "journey / trip", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-urlaub", DeckID: "a1-travel", Front: "der Urlaub", Back: "vacation / holiday", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-ferien", DeckID: "a1-travel", Front: "die Ferien", Back: "school holidays", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-flughafen", DeckID: "a1-travel", Front: "der Flughafen", Back: "airport", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-flugzeug", DeckID: "a1-travel", Front: "das Flugzeug", Back: "airplane", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-flug", DeckID: "a1-travel", Front: "der Flug", Back: "flight", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-bahnhof", DeckID: "a1-travel", Front: "der Bahnhof", Back: "train station", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-zug", DeckID: "a1-travel", Front: "der Zug", Back: "train", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-bahn", DeckID: "a1-travel", Front: "die Bahn", Back: "railway / train", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-gleis", DeckID: "a1-travel", Front: "das Gleis", Back: "platform / track", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-fahrkarte", DeckID: "a1-travel", Front: "die Fahrkarte", Back: "ticket", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-ticket", DeckID: "a1-travel", Front: "das Ticket", Back: "ticket", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-auto", DeckID: "a1-travel", Front: "das Auto", Back: "car", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-bus", DeckID: "a1-travel", Front: "der Bus", Back: "bus", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-strassenbahn", DeckID: "a1-travel", Front: "die Straßenbahn", Back: "tram", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-u-bahn", DeckID: "a1-travel", Front: "die U-Bahn", Back: "subway / underground", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-taxi", DeckID: "a1-travel", Front: "das Taxi", Back: "taxi", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-fahrrad", DeckID: "a1-travel", Front: "das Fahrrad", Back: "bicycle", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-schiff", DeckID: "a1-travel", Front: "das Schiff", Back: "ship", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-hotel", DeckID: "a1-travel", Front: "das Hotel", Back: "hotel", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-zimmer", DeckID: "a1-travel", Front: "das Zimmer", Back: "room", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-gepaeck", DeckID: "a1-travel", Front: "das Gepäck", Back: "luggage", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-koffer", DeckID: "a1-travel", Front: "der Koffer", Back: "suitcase", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-tasche", DeckID: "a1-travel", Front: "die Tasche", Back: "bag", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-reisepass", DeckID: "a1-travel", Front: "der Reisepass", Back: "passport", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-ausweis", DeckID: "a1-travel", Front: "der Ausweis", Back: "ID card", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-tourist", DeckID: "a1-travel", Front: "der Tourist", Back: "tourist", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-stadtplan", DeckID: "a1-travel", Front: "der Stadtplan", Back: "city map", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-sehenswuerdigkeit", DeckID: "a1-travel", Front: "die Sehenswürdigkeit", Back: "sight / attraction", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-ausflug", DeckID: "a1-travel", Front: "der Ausflug", Back: "excursion / trip", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-strand", DeckID: "a1-travel", Front: "der Strand", Back: "beach", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-meer", DeckID: "a1-travel", Front: "das Meer", Back: "sea", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-berge", DeckID: "a1-travel", Front: "die Berge", Back: "mountains", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-see", DeckID: "a1-travel", Front: "der See", Back: "lake", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-insel", DeckID: "a1-travel", Front: "die Insel", Back: "island", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-abfahrt", DeckID: "a1-travel", Front: "die Abfahrt", Back: "departure", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-ankunft", DeckID: "a1-travel", Front: "die Ankunft", Back: "arrival", Tags: []string{"a1", "travel", "noun"}},
		{ID: "a1-tra-verspaetung", DeckID: "a1-travel", Front: "die Verspätung", Back: "delay", Tags: []string{"a1", "travel", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-travel",
		Name:        "A1 Travel & Transport",
		Description: "Essential A1 vocabulary for traveling and transportation.",
		Tags:        []string{"german", "a1", "travel"},
		Notes:       notes,
	}
}
