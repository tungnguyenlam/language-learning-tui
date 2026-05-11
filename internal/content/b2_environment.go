package content

import (
	"deutsch-tui/internal/core"
)

func B2EnvironmentDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-env-umwelt", DeckID: "b2-environment", Front: "die Umwelt", Back: "environment", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-natur", DeckID: "b2-environment", Front: "die Natur", Back: "nature", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-ökologie", DeckID: "b2-environment", Front: "die Ökologie", Back: "ecology", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-umweltschutz", DeckID: "b2-environment", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-klima", DeckID: "b2-environment", Front: "das Klima", Back: "climate", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-klimawandel", DeckID: "b2-environment", Front: "der Klimawandel", Back: "climate change", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-treibhauseffekt", DeckID: "b2-environment", Front: "der Treibhauseffekt", Back: "greenhouse effect", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-co2", DeckID: "b2-environment", Front: "die CO2-Emission", Back: "CO2 emission", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-erneuerbar", DeckID: "b2-environment", Front: "erneuerbar", Back: "renewable", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-solar", DeckID: "b2-environment", Front: "solar", Back: "solar", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-windkraft", DeckID: "b2-environment", Front: "die Windkraft", Back: "wind power", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-wasserkraft", DeckID: "b2-environment", Front: "die Wasserkraft", Back: "hydropower", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-geothermie", DeckID: "b2-environment", Front: "die Geothermie", Back: "geothermal energy", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-biomasse", DeckID: "b2-environment", Front: "die Biomasse", Back: "biomass", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-nachhaltig", DeckID: "b2-environment", Front: "nachhaltig", Back: "sustainable", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-ressource", DeckID: "b2-environment", Front: "die Ressource", Back: "resource", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-rohstoff", DeckID: "b2-environment", Front: "der Rohstoff", Back: "raw material", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-energie", DeckID: "b2-environment", Front: "die Energie", Back: "energy", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-strom", DeckID: "b2-environment", Front: "der Strom", Back: "electricity", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-gas", DeckID: "b2-environment", Front: "das Gas", Back: "gas", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-öl", DeckID: "b2-environment", Front: "das Öl", Back: "oil", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-kohle", DeckID: "b2-environment", Front: "die Kohle", Back: "coal", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-atomkraft", DeckID: "b2-environment", Front: "die Atomkraft", Back: "nuclear power", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-atomkraftwerk", DeckID: "b2-environment", Front: "das Atomkraftwerk", Back: "nuclear power plant", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-strahlung", DeckID: "b2-environment", Front: "die Strahlung", Back: "radiation", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-müll", DeckID: "b2-environment", Front: "der Müll", Back: "trash / waste", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-recycling", DeckID: "b2-environment", Front: "das Recycling", Back: "recycling", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-mülldeponie", DeckID: "b2-environment", Front: "die Mülldeponie", Back: "landfill", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-entsorgung", DeckID: "b2-environment", Front: "die Entsorgung", Back: "waste disposal", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-verschmutzung", DeckID: "b2-environment", Front: "die Verschmutzung", Back: "pollution", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-luftverschmutzung", DeckID: "b2-environment", Front: "die Luftverschmutzung", Back: "air pollution", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-wasserverschmutzung", DeckID: "b2-environment", Front: "die Wasserverschmutzung", Back: "water pollution", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-emission", DeckID: "b2-environment", Front: "die Emission", Back: "emission", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-regenwald", DeckID: "b2-environment", Front: "der Regenwald", Back: "rainforest", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-wald", DeckID: "b2-environment", Front: "der Wald", Back: "forest", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-ozean", DeckID: "b2-environment", Front: "der Ozean", Back: "ocean", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-meer", DeckID: "b2-environment", Front: "das Meer", Back: "sea", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-gletscher", DeckID: "b2-environment", Front: "der Gletscher", Back: "glacier", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-artenschutz", DeckID: "b2-environment", Front: "der Artenschutz", Back: "species protection", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-gefährdet", DeckID: "b2-environment", Front: "gefährdet", Back: "endangered", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-art", DeckID: "b2-environment", Front: "die Art", Back: "species", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-ökosystem", DeckID: "b2-environment", Front: "das Ökosystem", Back: "ecosystem", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-biodiversität", DeckID: "b2-environment", Front: "die Biodiversität", Back: "biodiversity", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-lebensraum", DeckID: "b2-environment", Front: "der Lebensraum", Back: "habitat", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-treibhausgas", DeckID: "b2-environment", Front: "das Treibhausgas", Back: "greenhouse gas", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-kohlenstoff", DeckID: "b2-environment", Front: "der Kohlenstoff", Back: "carbon", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-fußabdruck", DeckID: "b2-environment", Front: "der ökologische Fußabdruck", Back: "ecological footprint", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-kreislauf", DeckID: "b2-environment", Front: "der Kreislauf", Back: "cycle", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-reduzieren", DeckID: "b2-environment", Front: "reduzieren", Back: "to reduce", Tags: []string{"b2", "environment"}},
		{ID: "b2-env-vermeiden", DeckID: "b2-environment", Front: "vermeiden", Back: "to avoid", Tags: []string{"b2", "environment"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-environment",
		Name:        "German B2 Environment & Sustainability",
		Description: "Vocabulary for environmental issues, sustainability, and ecology.",
		Tags:        []string{"german", "b2", "environment", "sustainability"},
		Notes:       notes,
	}
}
