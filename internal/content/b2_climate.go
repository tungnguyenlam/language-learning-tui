package content

import (
	"deutsch-tui/internal/core"
)

func B2ClimateDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-clm-klima", DeckID: "b2-climate", Front: "das Klima", Back: "climate", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-klimawandel", DeckID: "b2-climate", Front: "der Klimawandel", Back: "climate change", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-erderwaermung", DeckID: "b2-climate", Front: "die Erderwärmung", Back: "global warming", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-treibhausgas", DeckID: "b2-climate", Front: "das Treibhausgas", Back: "greenhouse gas", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-treibhauseffekt", DeckID: "b2-climate", Front: "der Treibhauseffekt", Back: "greenhouse effect", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-kohlendioxid", DeckID: "b2-climate", Front: "das Kohlendioxid", Back: "carbon dioxide", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-emission", DeckID: "b2-climate", Front: "die Emission", Back: "emission", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-co2fussabdr", DeckID: "b2-climate", Front: "der CO2-Fußabdruck", Back: "carbon footprint", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-fossil", DeckID: "b2-climate", Front: "fossile Brennstoffe", Back: "fossil fuels", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-erneuerbar", DeckID: "b2-climate", Front: "erneuerbare Energien", Back: "renewable energy", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-solaren", DeckID: "b2-climate", Front: "die Solarenergie", Back: "solar energy", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-windkraft", DeckID: "b2-climate", Front: "die Windkraft", Back: "wind power", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-wasserkraft", DeckID: "b2-climate", Front: "die Wasserkraft", Back: "hydropower", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-atomkraft", DeckID: "b2-climate", Front: "die Atomkraft", Back: "nuclear power", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-energiewende", DeckID: "b2-climate", Front: "die Energiewende", Back: "energy transition", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-nachhaltig", DeckID: "b2-climate", Front: "nachhaltig", Back: "sustainable", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-nachhaltigk", DeckID: "b2-climate", Front: "die Nachhaltigkeit", Back: "sustainability", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-umweltschutz", DeckID: "b2-climate", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-naturschutz", DeckID: "b2-climate", Front: "der Naturschutz", Back: "nature conservation", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-artenschutz", DeckID: "b2-climate", Front: "der Artenschutz", Back: "species protection", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-aussterben", DeckID: "b2-climate", Front: "aussterben", Back: "to become extinct", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-abholzung", DeckID: "b2-climate", Front: "die Abholzung", Back: "deforestation", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-aufforsten", DeckID: "b2-climate", Front: "aufforsten", Back: "to reforest", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-meeresspieg", DeckID: "b2-climate", Front: "der Meeresspiegel", Back: "sea level", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-gletscher", DeckID: "b2-climate", Front: "der Gletscher", Back: "glacier", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-schmelzen", DeckID: "b2-climate", Front: "schmelzen", Back: "to melt", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-duerre", DeckID: "b2-climate", Front: "die Dürre", Back: "drought", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-flut", DeckID: "b2-climate", Front: "die Flut", Back: "flood", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-hitzewelle", DeckID: "b2-climate", Front: "die Hitzewelle", Back: "heatwave", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-extremwetter", DeckID: "b2-climate", Front: "das Extremwetter", Back: "extreme weather", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-luftverschm", DeckID: "b2-climate", Front: "die Luftverschmutzung", Back: "air pollution", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-feinstaub", DeckID: "b2-climate", Front: "der Feinstaub", Back: "particulate matter", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-recycling", DeckID: "b2-climate", Front: "das Recycling", Back: "recycling", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-muelltren", DeckID: "b2-climate", Front: "die Mülltrennung", Back: "waste separation", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-pfand", DeckID: "b2-climate", Front: "das Pfand", Back: "deposit (on bottles)", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-einweg", DeckID: "b2-climate", Front: "die Einwegverpackung", Back: "single-use packaging", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-plastik", DeckID: "b2-climate", Front: "der Plastikmüll", Back: "plastic waste", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-klimaziel", DeckID: "b2-climate", Front: "das Klimaziel", Back: "climate goal / target", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-pariser", DeckID: "b2-climate", Front: "das Pariser Klimaabkommen", Back: "Paris Climate Agreement", Tags: []string{"b2", "climate"}},
		{ID: "b2-clm-klimaneutral", DeckID: "b2-climate", Front: "klimaneutral", Back: "climate-neutral", Tags: []string{"b2", "climate"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-climate",
		Name:        "B2 German Climate & Sustainability",
		Description: "Intermediate-advanced vocabulary for climate change, sustainability, energy, and environmental policy.",
		Tags:        []string{"german", "b2", "climate", "environment", "sustainability"},
		Notes:       notes,
	}
}
