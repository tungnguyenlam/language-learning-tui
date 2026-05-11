package content

import (
	"deutsch-tui/internal/core"
)

func B2ScienceIIDeck() core.Deck {
	notes := []core.Note{
		{ID: "sci2-001", DeckID: "b2-science-ii", Front: "die Wissenschaft", Back: "science", Extra: "Plural: die Wissenschaften", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-002", DeckID: "b2-science-ii", Front: "der Wissenschaftler", Back: "scientist", Extra: "Plural: die Wissenschaftler", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-003", DeckID: "b2-science-ii", Front: "forschen", Back: "to research", Extra: "Konjugation: ich forsche", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-004", DeckID: "b2-science-ii", Front: "die Forschung", Back: "research", Extra: "Plural: die Forschungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-005", DeckID: "b2-science-ii", Front: "das Experiment", Back: "experiment", Extra: "Plural: die Experimente", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-006", DeckID: "b2-science-ii", Front: "durchführen", Back: "to carry out, conduct", Extra: "Konjugation: ich führe durch", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-007", DeckID: "b2-science-ii", Front: "die Entdeckung", Back: "discovery", Extra: "Plural: die Entdeckungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-008", DeckID: "b2-science-ii", Front: "entdecken", Back: "to discover", Extra: "Konjugation: ich entdecke", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-009", DeckID: "b2-science-ii", Front: "der Beweis", Back: "proof, evidence", Extra: "Plural: die Beweise", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-010", DeckID: "b2-science-ii", Front: "beweisen", Back: "to prove", Extra: "Konjugation: ich beweise", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-011", DeckID: "b2-science-ii", Front: "die Theorie", Back: "theory", Extra: "Plural: die Theorien", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-012", DeckID: "b2-science-ii", Front: "die Praxis", Back: "practice", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-013", DeckID: "b2-science-ii", Front: "beobachten", Back: "to observe", Extra: "Konjugation: ich beobachte", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-014", DeckID: "b2-science-ii", Front: "die Beobachtung", Back: "observation", Extra: "Plural: die Beobachtungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-015", DeckID: "b2-science-ii", Front: "messen", Back: "to measure", Extra: "Konjugation: ich messe, du misst", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-016", DeckID: "b2-science-ii", Front: "die Messung", Back: "measurement", Extra: "Plural: die Messungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-017", DeckID: "b2-science-ii", Front: "das Ergebnis", Back: "result", Extra: "Plural: die Ergebnisse", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-018", DeckID: "b2-science-ii", Front: "analysieren", Back: "to analyze", Extra: "Konjugation: ich analysiere", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-019", DeckID: "b2-science-ii", Front: "die Analyse", Back: "analysis", Extra: "Plural: die Analysen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-020", DeckID: "b2-science-ii", Front: "die Methode", Back: "method", Extra: "Plural: die Methoden", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-021", DeckID: "b2-science-ii", Front: "entwickeln", Back: "to develop", Extra: "Konjugation: ich entwickle", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-022", DeckID: "b2-science-ii", Front: "die Entwicklung", Back: "development", Extra: "Plural: die Entwicklungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-023", DeckID: "b2-science-ii", Front: "das Phänomen", Back: "phenomenon", Extra: "Plural: die Phänomene", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-024", DeckID: "b2-science-ii", Front: "erklären", Back: "to explain", Extra: "Konjugation: ich erkläre", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-025", DeckID: "b2-science-ii", Front: "die Erklärung", Back: "explanation", Extra: "Plural: die Erklärungen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-026", DeckID: "b2-science-ii", Front: "die Studie", Back: "study", Extra: "Plural: die Studien", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-027", DeckID: "b2-science-ii", Front: "veröffentlichen", Back: "to publish", Extra: "Konjugation: ich veröffentliche", Tags: []string{"b2", "science", "verb"}},
		{ID: "sci2-028", DeckID: "b2-science-ii", Front: "die Physik", Back: "physics", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-029", DeckID: "b2-science-ii", Front: "die Chemie", Back: "chemistry", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-030", DeckID: "b2-science-ii", Front: "die Biologie", Back: "biology", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-031", DeckID: "b2-science-ii", Front: "das Labor", Back: "laboratory", Extra: "Plural: die Labore", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-032", DeckID: "b2-science-ii", Front: "die Energie", Back: "energy", Extra: "Plural: die Energien", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-033", DeckID: "b2-science-ii", Front: "die Materie", Back: "matter", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-034", DeckID: "b2-science-ii", Front: "das Atom", Back: "atom", Extra: "Plural: die Atome", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-035", DeckID: "b2-science-ii", Front: "die Zelle", Back: "cell", Extra: "Plural: die Zellen", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-036", DeckID: "b2-science-ii", Front: "die DNA", Back: "DNA", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-037", DeckID: "b2-science-ii", Front: "die Evolution", Back: "evolution", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-038", DeckID: "b2-science-ii", Front: "das Universum", Back: "universe", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-039", DeckID: "b2-science-ii", Front: "der Planet", Back: "planet", Extra: "Plural: die Planeten", Tags: []string{"b2", "science", "noun"}},
		{ID: "sci2-040", DeckID: "b2-science-ii", Front: "die Schwerkraft", Back: "gravity", Tags: []string{"b2", "science", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-science-ii",
		Name:        "B2 Science & Nature II",
		Description: "Advanced vocabulary for science, research, and natural phenomena",
		Tags:        []string{"b2", "science", "nature"},
		Notes:       notes,
	}
}
