package content

import (
	"deutsch-tui/internal/core"
)

func C1ScienceDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-sci-wissenschaft", DeckID: "c1-science", Front: "die Wissenschaft", Back: "science", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-forscher", DeckID: "c1-science", Front: "der Forscher / die Forscherin", Back: "researcher / scientist", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-experiment", DeckID: "c1-science", Front: "das Experiment", Back: "experiment", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-hypothese", DeckID: "c1-science", Front: "die Hypothese", Back: "hypothesis", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-theorie", DeckID: "c1-science", Front: "die Theorie", Back: "theory", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-beweis", DeckID: "c1-science", Front: "der Beweis", Back: "proof / evidence", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-physik", DeckID: "c1-science", Front: "die Physik", Back: "physics", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-chemie", DeckID: "c1-science", Front: "die Chemie", Back: "chemistry", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-biologie", DeckID: "c1-science", Front: "die Biologie", Back: "biology", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-mathematik", DeckID: "c1-science", Front: "die Mathematik", Back: "mathematics", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-atom", DeckID: "c1-science", Front: "das Atom", Back: "atom", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-molekül", DeckID: "c1-science", Front: "das Molekül", Back: "molecule", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-zelle", DeckID: "c1-science", Front: "die Zelle", Back: "cell", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-gen", DeckID: "c1-science", Front: "das Gen", Back: "gene", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-dna", DeckID: "c1-science", Front: "die DNS / DNA", Back: "DNA", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-evolution", DeckID: "c1-science", Front: "die Evolution", Back: "evolution", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-mutation", DeckID: "c1-science", Front: "die Mutation", Back: "mutation", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-kraft", DeckID: "c1-science", Front: "die Kraft", Back: "force", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-energie", DeckID: "c1-science", Front: "die Energie", Back: "energy", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-welle", DeckID: "c1-science", Front: "die Welle", Back: "wave", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-licht", DeckID: "c1-science", Front: "das Licht", Back: "light", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-elektron", DeckID: "c1-science", Front: "das Elektron", Back: "electron", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-proton", DeckID: "c1-science", Front: "das Proton", Back: "proton", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-neutron", DeckID: "c1-science", Front: "das Neutron", Back: "neutron", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-quanten", DeckID: "c1-science", Front: "das Quant", Back: "quantum", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-relativität", DeckID: "c1-science", Front: "die Relativitätstheorie", Back: "theory of relativity", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-messung", DeckID: "c1-science", Front: "die Messung", Back: "measurement", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-daten", DeckID: "c1-science", Front: "die Daten (Pl.)", Back: "data", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-analyse", DeckID: "c1-science", Front: "die Analyse", Back: "analysis", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-simulation", DeckID: "c1-science", Front: "die Simulation", Back: "simulation", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-modell", DeckID: "c1-science", Front: "das Modell", Back: "model", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-formel", DeckID: "c1-science", Front: "die Formel", Back: "formula", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-gleichung", DeckID: "c1-science", Front: "die Gleichung", Back: "equation", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-variable", DeckID: "c1-science", Front: "die Variable", Back: "variable", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-konstante", DeckID: "c1-science", Front: "die Konstante", Back: "constant", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-funktion", DeckID: "c1-science", Front: "die Funktion", Back: "function", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-algorithmus", DeckID: "c1-science", Front: "der Algorithmus", Back: "algorithm", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-berechnung", DeckID: "c1-science", Front: "die Berechnung", Back: "calculation", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-statistik", DeckID: "c1-science", Front: "die Statistik", Back: "statistics", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-wahrscheinlichkeit", DeckID: "c1-science", Front: "die Wahrscheinlichkeit", Back: "probability", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-mikroskop", DeckID: "c1-science", Front: "das Mikroskop", Back: "microscope", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-teleskop", DeckID: "c1-science", Front: "das Teleskop", Back: "telescope", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-satelit", DeckID: "c1-science", Front: "der Satellit", Back: "satellite", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-galaxie", DeckID: "c1-science", Front: "die Galaxie", Back: "galaxy", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-universum", DeckID: "c1-science", Front: "das Universum", Back: "universe", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-planet", DeckID: "c1-science", Front: "der Planet", Back: "planet", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-stern", DeckID: "c1-science", Front: "der Stern", Back: "star", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-reaktor", DeckID: "c1-science", Front: "der Reaktor", Back: "reactor", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-partikel", DeckID: "c1-science", Front: "das Teilchen / die Partikel", Back: "particle", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-feld", DeckID: "c1-science", Front: "das Feld", Back: "field (physics)", Tags: []string{"c1", "science"}},
		{ID: "c1-sci-dimension", DeckID: "c1-science", Front: "die Dimension", Back: "dimension", Tags: []string{"c1", "science"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-science",
		Name:        "C1 German Science & Technology",
		Description: "Advanced scientific and technical vocabulary.",
		Tags:        []string{"german", "c1", "science", "technology"},
		Notes:       notes,
	}
}
