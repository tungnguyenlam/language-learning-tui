package content

import (
	"deutsch-tui/internal/core"
)

func B1EnvironmentDeck() core.Deck {
	notes := []core.Note{
		{ID: "env-001", DeckID: "b1-environment", Front: "die Umwelt", Back: "environment", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-002", DeckID: "b1-environment", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-003", DeckID: "b1-environment", Front: "die Natur", Back: "nature", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-004", DeckID: "b1-environment", Front: "ökologisch", Back: "ecological, eco-friendly", Extra: "Related: die Ökologie", Tags: []string{"b1", "environment", "adjective"}},
		{ID: "env-005", DeckID: "b1-environment", Front: "nachhaltig", Back: "sustainable", Extra: "Related: die Nachhaltigkeit", Tags: []string{"b1", "environment", "adjective"}},
		{ID: "env-006", DeckID: "b1-environment", Front: "der Klimawandel", Back: "climate change", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-007", DeckID: "b1-environment", Front: "die Erderwärmung", Back: "global warming", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-008", DeckID: "b1-environment", Front: "die Treibhausgase", Back: "greenhouse gases", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-009", DeckID: "b1-environment", Front: "der CO2-Ausstoß", Back: "CO2 emissions", Extra: "Plural: die CO2-Ausstoße", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-010", DeckID: "b1-environment", Front: "die Luftverschmutzung", Back: "air pollution", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-011", DeckID: "b1-environment", Front: "die Wasserverschmutzung", Back: "water pollution", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-012", DeckID: "b1-environment", Front: "der Müll", Back: "trash, waste", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-013", DeckID: "b1-environment", Front: "das Recycling", Back: "recycling", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-014", DeckID: "b1-environment", Front: "wiederverwerten", Back: "to recycle, to reuse", Extra: "Konjugation: ich wiederverwerte", Tags: []string{"b1", "environment", "verb"}},
		{ID: "env-015", DeckID: "b1-environment", Front: "die Mülltrennung", Back: "waste separation", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-016", DeckID: "b1-environment", Front: "der Plastikmüll", Back: "plastic waste", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-017", DeckID: "b1-environment", Front: "die Verpackung", Back: "packaging", Extra: "Plural: die Verpackungen", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-018", DeckID: "b1-environment", Front: "reduzieren", Back: "to reduce", Extra: "Konjugation: ich reduziere", Tags: []string{"b1", "environment", "verb"}},
		{ID: "env-019", DeckID: "b1-environment", Front: "sparen", Back: "to save, to conserve", Extra: "Konjugation: ich spare", Tags: []string{"b1", "environment", "verb"}},
		{ID: "env-020", DeckID: "b1-environment", Front: "die erneuerbare Energie", Back: "renewable energy", Extra: "Plural: die erneuerbaren Energien", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-021", DeckID: "b1-environment", Front: "die Solarenergie", Back: "solar energy", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-022", DeckID: "b1-environment", Front: "die Windenergie", Back: "wind energy", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-023", DeckID: "b1-environment", Front: "die Wasserkraft", Back: "hydropower", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-024", DeckID: "b1-environment", Front: "das Kraftwerk", Back: "power plant", Extra: "Plural: die Kraftwerke", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-025", DeckID: "b1-environment", Front: "der Windpark", Back: "wind farm", Extra: "Plural: die Windparks", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-026", DeckID: "b1-environment", Front: "das Solarmodul", Back: "solar panel", Extra: "Plural: die Solarmodule", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-027", DeckID: "b1-environment", Front: "die Energieeffizienz", Back: "energy efficiency", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-028", DeckID: "b1-environment", Front: "der Artenschutz", Back: "species protection", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-029", DeckID: "b1-environment", Front: "die Bedrohung", Back: "threat", Extra: "Plural: die Bedrohungen", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-030", DeckID: "b1-environment", Front: "die Art", Back: "species", Extra: "Plural: die Arten", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-031", DeckID: "b1-environment", Front: "das Ökosystem", Back: "ecosystem", Extra: "Plural: die Ökosysteme", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-032", DeckID: "b1-environment", Front: "der Lebensraum", Back: "habitat", Extra: "Plural: die Lebensräume", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-033", DeckID: "b1-environment", Front: "der Regenwald", Back: "rainforest", Extra: "Plural: die Regenwälder", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-034", DeckID: "b1-environment", Front: "der Wald", Back: "forest", Extra: "Plural: die Wälder", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-035", DeckID: "b1-environment", Front: "das Tier", Back: "animal", Extra: "Plural: die Tiere", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-036", DeckID: "b1-environment", Front: "die Pflanze", Back: "plant", Extra: "Plural: die Pflanzen", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-037", DeckID: "b1-environment", Front: "der Treibstoff", Back: "fuel", Extra: "Plural: die Treibstoffe", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-038", DeckID: "b1-environment", Front: "das Benzin", Back: "gasoline, petrol", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-039", DeckID: "b1-environment", Front: "das Erdöl", Back: "crude oil", Tags: []string{"b1", "environment", "noun"}},
		{ID: "env-040", DeckID: "b1-environment", Front: "die Emission", Back: "emission", Extra: "Plural: die Emissionen", Tags: []string{"b1", "environment", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-environment",
		Name:        "B1 Environment & Sustainability",
		Description: "Environmental vocabulary, sustainability, and ecology topics",
		Tags:        []string{"b1", "environment", "sustainability", "ecology"},
		Notes:       notes,
	}
}
