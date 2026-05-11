package content

import (
	"deutsch-tui/internal/core"
)

func B1NatureDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-nat-umwelt", DeckID: "b1-nature", Front: "die Umwelt", Back: "environment", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-klimawandel", DeckID: "b1-nature", Front: "der Klimawandel", Back: "climate change", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-nachhaltigkeit", DeckID: "b1-nature", Front: "die Nachhaltigkeit", Back: "sustainability", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-erneuerbar", DeckID: "b1-nature", Front: "erneuerbare Energien", Back: "renewable energies", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-umweltschutz", DeckID: "b1-nature", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-verschmutzung", DeckID: "b1-nature", Front: "die Verschmutzung", Back: "pollution", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-recycling", DeckID: "b1-nature", Front: "das Recycling", Back: "recycling", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-muelltrennung", DeckID: "b1-nature", Front: "die Mülltrennung", Back: "waste separation", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-oekologisch", DeckID: "b1-nature", Front: "ökologisch", Back: "ecological", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-artenschutz", DeckID: "b1-nature", Front: "der Artenschutz", Back: "species protection", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-treibhauseffekt", DeckID: "b1-nature", Front: "der Treibhauseffekt", Back: "greenhouse effect", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-erwaermung", DeckID: "b1-nature", Front: "die Erderwärmung", Back: "global warming", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-abgase", DeckID: "b1-nature", Front: "die Abgase", Back: "exhaust gases", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-biologisch", DeckID: "b1-nature", Front: "biologisch abbaubar", Back: "biodegradable", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-oekosystem", DeckID: "b1-nature", Front: "das Ökosystem", Back: "ecosystem", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-ressourcen", DeckID: "b1-nature", Front: "die Ressourcen", Back: "resources", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-windkraft", DeckID: "b1-nature", Front: "die Windkraft", Back: "wind power", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-solarenergie", DeckID: "b1-nature", Front: "die Solarenergie", Back: "solar energy", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-wasserkraft", DeckID: "b1-nature", Front: "die Wasserkraft", Back: "hydroelectric power", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-atomausstieg", DeckID: "b1-nature", Front: "der Atomausstieg", Back: "nuclear phase-out", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-oeko-strom", DeckID: "b1-nature", Front: "der Ökostrom", Back: "green electricity", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-naturschutzgebiet", DeckID: "b1-nature", Front: "das Naturschutzgebiet", Back: "nature reserve", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-artenvielfalt", DeckID: "b1-nature", Front: "die Artenvielfalt", Back: "biodiversity", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-aussterben", DeckID: "b1-nature", Front: "aussterben", Back: "to become extinct", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-bedroht", DeckID: "b1-nature", Front: "bedroht sein", Back: "to be threatened/endangered", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-umweltbewusst", DeckID: "b1-nature", Front: "umweltbewusst", Back: "environmentally conscious", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-klimaneutral", DeckID: "b1-nature", Front: "klimaneutral", Back: "climate neutral", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-co2-fussabdruck", DeckID: "b1-nature", Front: "der CO2-Fußabdruck", Back: "carbon footprint", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-einweg", DeckID: "b1-nature", Front: "die Einwegverpackung", Back: "single-use packaging", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-mehrweg", DeckID: "b1-nature", Front: "das Mehrwegsystem", Back: "reusable system", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-verzicht", DeckID: "b1-nature", Front: "der Verzicht auf", Back: "doing without / renunciation of", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-schuetzen", DeckID: "b1-nature", Front: "schützen", Back: "to protect", Tags: []string{"b1", "nature"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-nature",
		Name:        "German B1 Nature & Environment",
		Description: "Essential vocabulary for discussing nature, climate, and sustainability.",
		Tags:        []string{"german", "b1", "nature"},
		Notes:       notes,
	}
}
