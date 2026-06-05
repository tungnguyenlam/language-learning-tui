package content

import (
	"deutsch-tui/internal/core"
)

func B2LogisticsDeck() core.Deck {
	notes := []core.Note{
		{ID: "log-001", DeckID: "b2-logistics", Front: "die Logistik", Back: "logistics", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-002", DeckID: "b2-logistics", Front: "die Lieferkette", Back: "supply chain", Extra: "Plural: die Lieferketten", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-003", DeckID: "b2-logistics", Front: "das Lager", Back: "warehouse, storage", Extra: "Plural: die Lager", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-004", DeckID: "b2-logistics", Front: "der Zoll", Back: "customs", Extra: "Plural: die Zölle", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-005", DeckID: "b2-logistics", Front: "die Fracht", Back: "freight, cargo", Extra: "Plural: die Frachten", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-006", DeckID: "b2-logistics", Front: "der Spediteur", Back: "freight forwarder", Extra: "Plural: die Spediteure", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-007", DeckID: "b2-logistics", Front: "die Sendung", Back: "shipment", Extra: "Plural: die Sendungen", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-008", DeckID: "b2-logistics", Front: "das Inventar", Back: "inventory", Extra: "Plural: die Inventare", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-009", DeckID: "b2-logistics", Front: "der Lieferant", Back: "supplier", Extra: "Plural: die Lieferanten", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-010", DeckID: "b2-logistics", Front: "die Nachfrage", Back: "demand", Extra: "Angebot und Nachfrage (supply and demand)", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-011", DeckID: "b2-logistics", Front: "das Angebot", Back: "supply, offer", Extra: "Plural: die Angebote", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-012", DeckID: "b2-logistics", Front: "ausliefern", Back: "to deliver", Extra: "Konjugation: ich liefere aus", Tags: []string{"b2", "logistics", "verb"}},
		{ID: "log-013", DeckID: "b2-logistics", Front: "der Versand", Back: "dispatch, shipping", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-014", DeckID: "b2-logistics", Front: "die Lieferung", Back: "delivery", Extra: "Plural: die Lieferungen", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-015", DeckID: "b2-logistics", Front: "der Empfänger", Back: "recipient", Extra: "Plural: die Empfänger", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-016", DeckID: "b2-logistics", Front: "der Gabelstapler", Back: "forklift", Extra: "Plural: die Gabelstapler", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-017", DeckID: "b2-logistics", Front: "die Palette", Back: "pallet", Extra: "Plural: die Paletten", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-018", DeckID: "b2-logistics", Front: "die Route", Back: "route", Extra: "Plural: die Routen", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-019", DeckID: "b2-logistics", Front: "das Transportmittel", Back: "means of transport", Extra: "Plural: die Transportmittel", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-020", DeckID: "b2-logistics", Front: "die Verzögerung", Back: "delay", Extra: "Plural: die Verzögerungen", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-021", DeckID: "b2-logistics", Front: "die Lagerhalle", Back: "warehouse building", Extra: "Plural: die Lagerhallen", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-022", DeckID: "b2-logistics", Front: "die Beschaffung", Back: "procurement", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-023", DeckID: "b2-logistics", Front: "lagern", Back: "to store", Extra: "Konjugation: ich lagere", Tags: []string{"b2", "logistics", "verb"}},
		{ID: "log-024", DeckID: "b2-logistics", Front: "die Bestandskontrolle", Back: "inventory control", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-025", DeckID: "b2-logistics", Front: "verfolgen", Back: "to track, to follow", Extra: "Konjugation: ich verfolge", Tags: []string{"b2", "logistics", "verb"}},
		{ID: "log-026", DeckID: "b2-logistics", Front: "die Sendungsverfolgung", Back: "shipment tracking", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-027", DeckID: "b2-logistics", Front: "der Frachtbrief", Back: "waybill, bill of lading", Extra: "Plural: die Frachtbriefe", Tags: []string{"b2", "logistics", "noun"}},
		{ID: "log-028", DeckID: "b2-logistics", Front: "importieren", Back: "to import", Extra: "Konjugation: ich importiere", Tags: []string{"b2", "logistics", "verb"}},
		{ID: "log-029", DeckID: "b2-logistics", Front: "exportieren", Back: "to export", Extra: "Konjugation: ich exportiere", Tags: []string{"b2", "logistics", "verb"}},
		{ID: "log-030", DeckID: "b2-logistics", Front: "verpacken", Back: "to package", Extra: "Konjugation: ich verpacke", Tags: []string{"b2", "logistics", "verb"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-logistics",
		Name:        "B2 Logistics & Supply Chain",
		Description: "Vocabulary for logistics, transportation, and warehousing",
		Tags:        []string{"b2", "logistics", "business", "transport"},
		Notes:       notes,
	}
}
