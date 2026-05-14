package content

import (
	"deutsch-tui/internal/core"
)

func C1EnvironmentSustainabilityDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-env-nachhaltigkeit", DeckID: "c1-environment", Front: "die Nachhaltigkeit", Back: "sustainability", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-klimawandel", DeckID: "c1-environment", Front: "der Klimawandel", Back: "climate change", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-erderwaermung", DeckID: "c1-environment", Front: "die Erderwärmung", Back: "global warming", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-treibhausgase", DeckID: "c1-environment", Front: "die Treibhausgase", Back: "greenhouse gases", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-fussabdruck", DeckID: "c1-environment", Front: "der ökologische Fußabdruck", Back: "ecological footprint", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-erneuerbar", DeckID: "c1-environment", Front: "die erneuerbaren Energien", Back: "renewable energies", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-energieeffizienz", DeckID: "c1-environment", Front: "die Energieeffizienz", Back: "energy efficiency", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-artenschutz", DeckID: "c1-environment", Front: "der Artenschutz", Back: "species protection", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-biodiversitaet", DeckID: "c1-environment", Front: "die Artenvielfalt / Biodiversität", Back: "biodiversity", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-oekosystem", DeckID: "c1-environment", Front: "das Ökosystem", Back: "ecosystem", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-naturschutz", DeckID: "c1-environment", Front: "der Naturschutz", Back: "nature conservation", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-umweltbelastung", DeckID: "c1-environment", Front: "die Umweltbelastung", Back: "environmental pollution / burden", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-schadstoffausstoss", DeckID: "c1-environment", Front: "der Schadstoffausstoß", Back: "pollutant emissions", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-muellvermeidung", DeckID: "c1-environment", Front: "die Müllvermeidung", Back: "waste avoidance", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-recycling", DeckID: "c1-environment", Front: "das Recycling", Back: "recycling", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-kreislaufwirtschaft", DeckID: "c1-environment", Front: "die Kreislaufwirtschaft", Back: "circular economy", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-ressourcenknappheit", DeckID: "c1-environment", Front: "die Ressourcenknappheit", Back: "resource scarcity", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-umweltschutz", DeckID: "c1-environment", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-umweltpolitik", DeckID: "c1-environment", Front: "die Umweltpolitik", Back: "environmental policy", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-umweltbewusstsein", DeckID: "c1-environment", Front: "das Umweltbewusstsein", Back: "environmental awareness", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-fff", DeckID: "c1-environment", Front: "die Fridays-for-Future-Bewegung", Back: "Fridays for Future movement", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-agrarwende", DeckID: "c1-environment", Front: "die Agrarwende", Back: "agricultural transition", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-mobilitaetswende", DeckID: "c1-environment", Front: "die Mobilitätswende", Back: "mobility transition", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-atomausstieg", DeckID: "c1-environment", Front: "der Atomausstieg", Back: "nuclear phase-out", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-kohleausstieg", DeckID: "c1-environment", Front: "der Kohleausstieg", Back: "coal phase-out", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-verzicht", DeckID: "c1-environment", Front: "der Verzicht", Back: "renunciation / doing without", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-konsumverhalten", DeckID: "c1-environment", Front: "das Konsumverhalten", Back: "consumer behavior", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-wegwerfgesellschaft", DeckID: "c1-environment", Front: "die Wegwerfgesellschaft", Back: "throwaway society", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-massentierhaltung", DeckID: "c1-environment", Front: "die Massentierhaltung", Back: "factory farming", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-ueberfischung", DeckID: "c1-environment", Front: "die Überfischung", Back: "overfishing", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-entwaldung", DeckID: "c1-environment", Front: "die Entwaldung / Abholzung", Back: "deforestation", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-duerre", DeckID: "c1-environment", Front: "die Dürre", Back: "drought", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-ueberschwemmung", DeckID: "c1-environment", Front: "die Überschwemmung", Back: "flooding", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-extremwetter", DeckID: "c1-environment", Front: "die Extremwetterereignisse", Back: "extreme weather events", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-meeresspiegel", DeckID: "c1-environment", Front: "der Meeresspiegelanstieg", Back: "sea level rise", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-co2-steuer", DeckID: "c1-environment", Front: "die CO2-Steuer", Back: "carbon tax", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-emissionshandel", DeckID: "c1-environment", Front: "der Emissionshandel", Back: "emissions trading", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-kompensation", DeckID: "c1-environment", Front: "die Kompensation", Back: "compensation / offsetting", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-gentechnik", DeckID: "c1-environment", Front: "die grüne Gentechnik", Back: "green genetic engineering", Tags: []string{"c1", "environment", "noun"}},
		{ID: "c1-env-plastikverbot", DeckID: "c1-environment", Front: "das Plastikverbot", Back: "plastic ban", Tags: []string{"c1", "environment", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-environment",
		Name:        "C1 Environment & Sustainability",
		Description: "Advanced C1 vocabulary for environmental issues and sustainability.",
		Tags:        []string{"german", "c1", "environment", "may15"},
		Notes:       notes,
	}
}
