package content

import (
	"deutsch-tui/internal/core"
)

func B2TravelDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-trv-reisen", DeckID: "b2-travel", Front: "die Reise", Back: "trip / journey", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reisender", DeckID: "b2-travel", Front: "der Reisende / die Reisende", Back: "traveler", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reisebüro", DeckID: "b2-travel", Front: "das Reisebüro", Back: "travel agency", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-urlaub", DeckID: "b2-travel", Front: "der Urlaub", Back: "vacation / holiday", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-destination", DeckID: "b2-travel", Front: "das Reiseziel", Back: "destination", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-route", DeckID: "b2-travel", Front: "die Route", Back: "route", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-pass", DeckID: "b2-travel", Front: "der Reisepass", Back: "passport", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-visum", DeckID: "b2-travel", Front: "das Visum", Back: "visa", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-zoll", DeckID: "b2-travel", Front: "der Zoll", Back: "customs", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-grenzübergang", DeckID: "b2-travel", Front: "der Grenzübergang", Back: "border crossing", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-einreise", DeckID: "b2-travel", Front: "die Einreise", Back: "entry / immigration", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-ausreise", DeckID: "b2-travel", Front: "die Ausreise", Back: "departure / exit", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-aufenthalt", DeckID: "b2-travel", Front: "der Aufenthalt", Back: "stay / residence", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-visumantrag", DeckID: "b2-travel", Front: "der Visumantrag", Back: "visa application", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-unterkunft", DeckID: "b2-travel", Front: "die Unterkunft", Back: "accommodation", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-hotel", DeckID: "b2-travel", Front: "das Hotel", Back: "hotel", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-pension", DeckID: "b2-travel", Front: "die Pension", Back: "guesthouse", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-jugendherberge", DeckID: "b2-travel", Front: "die Jugendherberge", Back: "youth hostel", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-ferienwohnung", DeckID: "b2-travel", Front: "die Ferienwohnung", Back: "holiday apartment", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-campingplatz", DeckID: "b2-travel", Front: "der Campingplatz", Back: "campsite", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-gästezimmer", DeckID: "b2-travel", Front: "das Gästezimmer", Back: "guest room", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reservierung", DeckID: "b2-travel", Front: "die Reservierung", Back: "reservation", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-buchung", DeckID: "b2-travel", Front: "die Buchung", Back: "booking", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-stornieren", DeckID: "b2-travel", Front: "stornieren", Back: "to cancel", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-flug", DeckID: "b2-travel", Front: "der Flug", Back: "flight", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-flughafen", DeckID: "b2-travel", Front: "der Flughafen", Back: "airport", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-gate", DeckID: "b2-travel", Front: "das Gate", Back: "gate", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-boarding", DeckID: "b2-travel", Front: "das Boarding", Back: "boarding", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-terminal", DeckID: "b2-travel", Front: "das Terminal", Back: "terminal", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-gepäck", DeckID: "b2-travel", Front: "das Gepäck", Back: "luggage", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-koffer", DeckID: "b2-travel", Front: "der Koffer", Back: "suitcase", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-rucksack", DeckID: "b2-travel", Front: "der Rucksack", Back: "backpack", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-einchecken", DeckID: "b2-travel", Front: "einchecken", Back: "to check in", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-auschecken", DeckID: "b2-travel", Front: "auschecken", Back: "to check out", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-passagier", DeckID: "b2-travel", Front: "der Passagier / die Passagierin", Back: "passenger", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reiseversicherung", DeckID: "b2-travel", Front: "die Reiseversicherung", Back: "travel insurance", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reisekoordinator", DeckID: "b2-travel", Front: "der Reiseveranstalter", Back: "tour operator", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-rundreise", DeckID: "b2-travel", Front: "die Rundreise", Back: "tour / round trip", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-attraktion", DeckID: "b2-travel", Front: "die Sehenswürdigkeit", Back: "attraction / sight", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-monument", DeckID: "b2-travel", Front: "das Monument", Back: "monument", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-museum", DeckID: "b2-travel", Front: "das Museum", Back: "museum", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-galerie", DeckID: "b2-travel", Front: "die Galerie", Back: "gallery", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-ausstellung", DeckID: "b2-travel", Front: "die Ausstellung", Back: "exhibition", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-denkmal", DeckID: "b2-travel", Front: "das Denkmal", Back: "memorial", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-landkarte", DeckID: "b2-travel", Front: "die Landkarte", Back: "map", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-stadtplan", DeckID: "b2-travel", Front: "der Stadtplan", Back: "city map", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-reiseführer", DeckID: "b2-travel", Front: "der Reiseführer", Back: "travel guide", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-kompass", DeckID: "b2-travel", Front: "der Kompass", Back: "compass", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-wegweiser", DeckID: "b2-travel", Front: "der Wegweiser", Back: "signpost", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-landmark", DeckID: "b2-travel", Front: "das Wahrzeichen", Back: "landmark", Tags: []string{"b2", "travel"}},
		{ID: "b2-trv-tourist", DeckID: "b2-travel", Front: "der Tourist / die Touristin", Back: "tourist", Tags: []string{"b2", "travel"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-travel",
		Name:        "German B2 Travel & Tourism",
		Description: "Advanced vocabulary for travel, tourism, and geography.",
		Tags:        []string{"german", "b2", "travel", "tourism"},
		Notes:       notes,
	}
}
