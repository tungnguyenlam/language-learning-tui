package content

import (
	"deutsch-tui/internal/core"
)

func B1HousingApartmentDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-hou-mietwohnung", DeckID: "b1-housing", Front: "die Mietwohnung", Back: "rental apartment", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-mietvertrag", DeckID: "b1-housing", Front: "der Mietvertrag", Back: "rental contract", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-kaltmiete", DeckID: "b1-housing", Front: "die Kaltmiete", Back: "base rent (without utilities)", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-warmmiete", DeckID: "b1-housing", Front: "die Warmmiete", Back: "rent including utilities", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-nebenkosten", DeckID: "b1-housing", Front: "die Nebenkosten", Back: "additional costs / utilities", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-kaution", DeckID: "b1-housing", Front: "die Kaution", Back: "deposit / security deposit", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-vermieter", DeckID: "b1-housing", Front: "der Vermieter", Back: "landlord", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-mieter", DeckID: "b1-housing", Front: "der Mieter", Back: "tenant", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-wg", DeckID: "b1-housing", Front: "die WG (Wohngemeinschaft)", Back: "shared apartment", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-mitbewohner", DeckID: "b1-housing", Front: "der Mitbewohner", Back: "roommate", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-kuendigen", DeckID: "b1-housing", Front: "kündigen", Back: "to cancel / to terminate (contract)", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-kuendigungsfrist", DeckID: "b1-housing", Front: "die Kündigungsfrist", Back: "notice period", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-einziehen", DeckID: "b1-housing", Front: "einziehen", Back: "to move in", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-ausziehen", DeckID: "b1-housing", Front: "ausziehen", Back: "to move out", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-umziehen", DeckID: "b1-housing", Front: "umziehen", Back: "to move (to another house)", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-umzug", DeckID: "b1-housing", Front: "der Umzug", Back: "the move", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-erdgeschoss", DeckID: "b1-housing", Front: "das Erdgeschoss", Back: "ground floor", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-stockwerk", DeckID: "b1-housing", Front: "das Stockwerk", Back: "floor / story", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-dachboden", DeckID: "b1-housing", Front: "der Dachboden", Back: "attic", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-keller", DeckID: "b1-housing", Front: "der Keller", Back: "basement", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-balkon", DeckID: "b1-housing", Front: "der Balkon", Back: "balcony", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-terrasse", DeckID: "b1-housing", Front: "die Terrasse", Back: "terrace", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-garten", DeckID: "b1-housing", Front: "der Garten", Back: "garden", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-moebliert", DeckID: "b1-housing", Front: "möbliert", Back: "furnished", Tags: []string{"b1", "housing", "adjective"}},
		{ID: "b1-hou-unmoebliert", DeckID: "b1-housing", Front: "unmöbliert", Back: "unfurnished", Tags: []string{"b1", "housing", "adjective"}},
		{ID: "b1-hou-renovieren", DeckID: "b1-housing", Front: "renovieren", Back: "to renovate", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-saniert", DeckID: "b1-housing", Front: "saniert", Back: "refurbished / renovated", Tags: []string{"b1", "housing", "adjective"}},
		{ID: "b1-hou-hausordnung", DeckID: "b1-housing", Front: "die Hausordnung", Back: "house rules", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-muelltrennung", DeckID: "b1-housing", Front: "die Mülltrennung", Back: "waste separation", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-hausmeister", DeckID: "b1-housing", Front: "der Hausmeister", Back: "janitor / caretaker", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-heizung", DeckID: "b1-housing", Front: "die Heizung", Back: "heating", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-lueften", DeckID: "b1-housing", Front: "lüften", Back: "to ventilate / to air out", Tags: []string{"b1", "housing", "verb"}},
		{ID: "b1-hou-anzeige", DeckID: "b1-housing", Front: "die Anzeige", Back: "advertisement / listing", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-inserat", DeckID: "b1-housing", Front: "das Inserat", Back: "advertisement / listing", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-besichtigung", DeckID: "b1-housing", Front: "die Besichtigung", Back: "viewing / inspection", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-quadratmeter", DeckID: "b1-housing", Front: "Quadratmeter (qm/m²)", Back: "square meters", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-wohnflaeche", DeckID: "b1-housing", Front: "die Wohnfläche", Back: "living space", Tags: []string{"b1", "housing", "noun"}},
		{ID: "b1-hou-zentral", DeckID: "b1-housing", Front: "zentral gelegen", Back: "centrally located", Tags: []string{"b1", "housing", "adjective"}},
		{ID: "b1-hou-ruhig", DeckID: "b1-housing", Front: "ruhig gelegen", Back: "quietly located", Tags: []string{"b1", "housing", "adjective"}},
		{ID: "b1-hou-verkehrsguenstig", DeckID: "b1-housing", Front: "verkehrsgünstig", Back: "conveniently located for transport", Tags: []string{"b1", "housing", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-housing",
		Name:        "B1 Housing & Apartment",
		Description: "B1 vocabulary for renting, living, and housing in Germany.",
		Tags:        []string{"german", "b1", "housing"},
		Notes:       notes,
	}
}
