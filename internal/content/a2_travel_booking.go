package content

import (
	"deutsch-tui/internal/core"
)

func A2TravelBookingDeck() core.Deck {
	notes := []core.Note{
		{ID: "a2-trb-buchung", DeckID: "a2-travel-booking", Front: "die Buchung", Back: "booking", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-reservieren", DeckID: "a2-travel-booking", Front: "reservieren", Back: "to reserve", Tags: []string{"a2", "travel", "verb"}},
		{ID: "a2-trb-bestaetigung", DeckID: "a2-travel-booking", Front: "die Bestätigung", Back: "confirmation", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-stornieren", DeckID: "a2-travel-booking", Front: "stornieren", Back: "to cancel (a booking)", Tags: []string{"a2", "travel", "verb"}},
		{ID: "a2-trb-umbuchung", DeckID: "a2-travel-booking", Front: "die Umbuchung", Back: "rebooking / change of reservation", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-doppelzimmer", DeckID: "a2-travel-booking", Front: "das Doppelzimmer", Back: "double room", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-einzelzimmer", DeckID: "a2-travel-booking", Front: "das Einzelzimmer", Back: "single room", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-halbpension", DeckID: "a2-travel-booking", Front: "die Halbpension", Back: "half board", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-vollpension", DeckID: "a2-travel-booking", Front: "die Vollpension", Back: "full board", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-uebernachtung", DeckID: "a2-travel-booking", Front: "die Übernachtung", Back: "overnight stay", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-fruehstueck", DeckID: "a2-travel-booking", Front: "inklusive Frühstück", Back: "including breakfast", Tags: []string{"a2", "travel", "phrase"}},
		{ID: "a2-trb-rezeption", DeckID: "a2-travel-booking", Front: "die Rezeption", Back: "reception", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-checkin", DeckID: "a2-travel-booking", Front: "der Check-in", Back: "check-in", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-checkout", DeckID: "a2-travel-booking", Front: "der Check-out", Back: "check-out", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-kurtaxe", DeckID: "a2-travel-booking", Front: "die Kurtaxe", Back: "tourist tax", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-reisebuero", DeckID: "a2-travel-booking", Front: "das Reisebüro", Back: "travel agency", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-pauschalreise", DeckID: "a2-travel-booking", Front: "die Pauschalreise", Back: "package tour", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-reiseleiter", DeckID: "a2-travel-booking", Front: "der Reiseleiter", Back: "tour guide", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-versicherung", DeckID: "a2-travel-booking", Front: "die Versicherung", Back: "insurance", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-reiseruecktritt", DeckID: "a2-travel-booking", Front: "die Reiserücktrittsversicherung", Back: "travel cancellation insurance", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-mietwagen", DeckID: "a2-travel-booking", Front: "der Mietwagen", Back: "rental car", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-tankstelle", DeckID: "a2-travel-booking", Front: "die Tankstelle", Back: "gas station", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-tanken", DeckID: "a2-travel-booking", Front: "tanken", Back: "to refuel / to gas up", Tags: []string{"a2", "travel", "verb"}},
		{ID: "a2-trb-maut", DeckID: "a2-travel-booking", Front: "die Maut", Back: "toll", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-vignette", DeckID: "a2-travel-booking", Front: "die Vignette", Back: "vignette / toll sticker", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-gutschein", DeckID: "a2-travel-booking", Front: "der Reisegutschein", Back: "travel voucher", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-visum", DeckID: "a2-travel-booking", Front: "das Visum", Back: "visa", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-grenze", DeckID: "a2-travel-booking", Front: "der Grenzübergang", Back: "border crossing", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-zoll", DeckID: "a2-travel-booking", Front: "der Zoll", Back: "customs", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-bordkarte", DeckID: "a2-travel-booking", Front: "die Bordkarte", Back: "boarding pass", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-gate", DeckID: "a2-travel-booking", Front: "das Gate", Back: "gate", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-anschlussflug", DeckID: "a2-travel-booking", Front: "der Anschlussflug", Back: "connecting flight", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-terminal", DeckID: "a2-travel-booking", Front: "das Terminal", Back: "terminal", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-sitzplatz", DeckID: "a2-travel-booking", Front: "der Sitzplatz", Back: "seat", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-fensterplatz", DeckID: "a2-travel-booking", Front: "der Fensterplatz", Back: "window seat", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-gangplatz", DeckID: "a2-travel-booking", Front: "der Gangplatz", Back: "aisle seat", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-verspaetung", DeckID: "a2-travel-booking", Front: "die Verspätung", Back: "delay", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-flugausfall", DeckID: "a2-travel-booking", Front: "der Flugausfall", Back: "flight cancellation", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-reklamation", DeckID: "a2-travel-booking", Front: "die Reklamation", Back: "complaint", Tags: []string{"a2", "travel", "noun"}},
		{ID: "a2-trb-fundbuero", DeckID: "a2-travel-booking", Front: "das Fundbüro", Back: "lost and found", Tags: []string{"a2", "travel", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-travel-booking",
		Name:        "A2 Travel & Booking",
		Description: "A2 vocabulary for travel arrangements, hotels, and bookings.",
		Tags:        []string{"german", "a2", "travel"},
		Notes:       notes,
	}
}
