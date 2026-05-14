package content

import (
	"deutsch-tui/internal/core"
)

func B2TravelAdventureDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-travel-reise", DeckID: "b2-travel-adventure", Front: "die Reise", Back: "journey / trip", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-urlaub", DeckID: "b2-travel-adventure", Front: "der Urlaub", Back: "vacation / holiday", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseveranstalter", DeckID: "b2-travel-adventure", Front: "der Reiseveranstalter", Back: "tour operator", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseführer", DeckID: "b2-travel-adventure", Front: "der Reiseführer", Back: "travel guide (book/person)", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reisebüro", DeckID: "b2-travel-adventure", Front: "das Reisebüro", Back: "travel agency", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseversicherung", DeckID: "b2-travel-adventure", Front: "die Reiseversicherung", Back: "travel insurance", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseroute", DeckID: "b2-travel-adventure", Front: "die Reiseroute", Back: "itinerary", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reisekoffer", DeckID: "b2-travel-adventure", Front: "der Koffer", Back: "suitcase", Tags: []string{"a2", "travel", "noun"}},
		{ID: "b2-travel-reisegepäck", DeckID: "b2-travel-adventure", Front: "das Gepäck", Back: "luggage", Tags: []string{"a2", "travel", "noun"}},
		{ID: "b2-travel-check-in", DeckID: "b2-travel-adventure", Front: "der Check-in", Back: "check-in", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-einchecken", DeckID: "b2-travel-adventure", Front: "einchecken", Back: "to check in", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-auschecken", DeckID: "b2-travel-adventure", Front: "auschecken", Back: "to check out", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-flugbuchung", DeckID: "b2-travel-adventure", Front: "die Flugbuchung", Back: "flight booking", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-flughafen", DeckID: "b2-travel-adventure", Front: "der Flughafen", Back: "airport", Tags: []string{"a2", "travel", "noun"}},
		{ID: "b2-travel-terminal", DeckID: "b2-travel-adventure", Front: "das Terminal", Back: "terminal", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-abflug", DeckID: "b2-travel-adventure", Front: "der Abflug", Back: "departure", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-ankunft", DeckID: "b2-travel-adventure", Front: "die Ankunft", Back: "arrival", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-passagier", DeckID: "b2-travel-adventure", Front: "der Passagier", Back: "passenger", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-sitzplatz", DeckID: "b2-travel-adventure", Front: "der Sitzplatz", Back: "seat", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-sitzplatzreservierung", DeckID: "b2-travel-adventure", Front: "die Sitzplatzreservierung", Back: "seat reservation", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-notlandung", DeckID: "b2-travel-adventure", Front: "die Notlandung", Back: "emergency landing", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-verspätung", DeckID: "b2-travel-adventure", Front: "die Verspätung", Back: "delay", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-stornierung", DeckID: "b2-travel-adventure", Front: "die Stornierung", Back: "cancellation", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reisepass", DeckID: "b2-travel-adventure", Front: "der Reisepass", Back: "passport", Tags: []string{"a2", "travel", "noun"}},
		{ID: "b2-travel-visum", DeckID: "b2-travel-adventure", Front: "das Visum", Back: "visa", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-impfpass", DeckID: "b2-travel-adventure", Front: "der Impfpass", Back: "vaccination record", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-zoll", DeckID: "b2-travel-adventure", Front: "der Zoll", Back: "customs", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-unterkunft", DeckID: "b2-travel-adventure", Front: "die Unterkunft", Back: "accommodation", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-übernachtung", DeckID: "b2-travel-adventure", Front: "die Übernachtung", Back: "overnight stay", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-frühstück", DeckID: "b2-travel-adventure", Front: "das Frühstück", Back: "breakfast", Tags: []string{"a1", "travel", "noun"}},
		{ID: "b2-travel-halbpension", DeckID: "b2-travel-adventure", Front: "die Halbpension", Back: "half board", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-vollpension", DeckID: "b2-travel-adventure", Front: "die Vollpension", Back: "full board", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-urlaubsort", DeckID: "b2-travel-adventure", Front: "der Urlaubsort", Back: "holiday destination", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-sehenswürdigkeit", DeckID: "b2-travel-adventure", Front: "die Sehenswürdigkeit", Back: "sightseeing attraction", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-touristinfo", DeckID: "b2-travel-adventure", Front: "die Touristeninformation", Back: "tourist information", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-stadtrundfahrt", DeckID: "b2-travel-adventure", Front: "die Stadtrundfahrt", Back: "city tour", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-ausflug", DeckID: "b2-travel-adventure", Front: "der Ausflug", Back: "excursion / day trip", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-wandern", DeckID: "b2-travel-adventure", Front: "wandern", Back: "to hike", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-bergsteigen", DeckID: "b2-travel-adventure", Front: "bergsteigen", Back: "to mountain climb", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-tauchen", DeckID: "b2-travel-adventure", Front: "tauchen", Back: "to dive (scuba)", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-kajaken", DeckID: "b2-travel-adventure", Front: "kajaken", Back: "to kayak", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-skifahren", DeckID: "b2-travel-adventure", Front: "skifahren", Back: "to ski", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-snowboarden", DeckID: "b2-travel-adventure", Front: "snowboarden", Back: "to snowboard", Tags: []string{"b2", "travel", "verb"}},
		{ID: "b2-travel-camping", DeckID: "b2-travel-adventure", Front: "camping", Back: "camping", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-wohnmobil", DeckID: "b2-travel-adventure", Front: "das Wohnmobil", Back: "camper van", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-backpacking", DeckID: "b2-travel-adventure", Front: "das Backpacking", Back: "backpacking", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseblog", DeckID: "b2-travel-adventure", Front: "der Reiseblog", Back: "travel blog", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseerlebnis", DeckID: "b2-travel-adventure", Front: "das Reiseerlebnis", Back: "travel experience", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reisemesse", DeckID: "b2-travel-adventure", Front: "die Reisemesse", Back: "travel fair", Tags: []string{"b2", "travel", "noun"}},
		{ID: "b2-travel-reiseland", DeckID: "b2-travel-adventure", Front: "das Reiseland", Back: "travel destination country", Tags: []string{"b2", "travel", "noun"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "b2-travel-adventure",
		Name:        "German B2 Travel & Adventure",
		Description: "Advanced travel vocabulary for trips, flights, accommodations, and adventure activities.",
		Tags:        []string{"german", "b2", "travel", "adventure"},
		Notes:       notes,
	}
}
