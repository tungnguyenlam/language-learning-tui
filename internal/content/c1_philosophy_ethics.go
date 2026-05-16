package content

import (
	"deutsch-tui/internal/core"
)

func C1PhilosophyEthicsDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-phil-philosophie", DeckID: "c1-philosophy-ethics", Front: "die Philosophie", Back: "philosophy", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-ethik", DeckID: "c1-philosophy-ethics", Front: "die Ethik", Back: "ethics", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-metaphysik", DeckID: "c1-philosophy-ethics", Front: "die Metaphysik", Back: "metaphysics", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-erkenntnistheorie", DeckID: "c1-philosophy-ethics", Front: "die Erkenntnistheorie", Back: "epistemology", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-moral", DeckID: "c1-philosophy-ethics", Front: "die Moral", Back: "morality", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-tugend", DeckID: "c1-philosophy-ethics", Front: "die Tugend", Back: "virtue", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-pflicht", DeckID: "c1-philosophy-ethics", Front: "die Pflicht", Back: "duty / obligation", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-gewissen", DeckID: "c1-philosophy-ethics", Front: "das Gewissen", Back: "conscience", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-vernunft", DeckID: "c1-philosophy-ethics", Front: "die Vernunft", Back: "reason / rationality", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-verstand", DeckID: "c1-philosophy-ethics", Front: "der Verstand", Back: "intellect / mind", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-aufklaerung", DeckID: "c1-philosophy-ethics", Front: "die Aufklärung", Back: "Enlightenment", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-existentialismus", DeckID: "c1-philosophy-ethics", Front: "der Existentialismus", Back: "existentialism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-nihilismus", DeckID: "c1-philosophy-ethics", Front: "der Nihilismus", Back: "nihilism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-utilitarismus", DeckID: "c1-philosophy-ethics", Front: "der Utilitarismus", Back: "utilitarianism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-determinismus", DeckID: "c1-philosophy-ethics", Front: "der Determinismus", Back: "determinism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-willensfreiheit", DeckID: "c1-philosophy-ethics", Front: "die Willensfreiheit", Back: "free will", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-gerechtigkeit", DeckID: "c1-philosophy-ethics", Front: "die Gerechtigkeit", Back: "justice", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-wuerde", DeckID: "c1-philosophy-ethics", Front: "die Würde", Back: "dignity", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-menschenrechte", DeckID: "c1-philosophy-ethics", Front: "die Menschenrechte", Back: "human rights", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-verantwortung", DeckID: "c1-philosophy-ethics", Front: "die Verantwortung", Back: "responsibility", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-werte", DeckID: "c1-philosophy-ethics", Front: "die Werte", Back: "values", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-normen", DeckID: "c1-philosophy-ethics", Front: "die Normen", Back: "norms", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-dilemma", DeckID: "c1-philosophy-ethics", Front: "das Dilemma", Back: "dilemma", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-argumentation", DeckID: "c1-philosophy-ethics", Front: "die Argumentation", Back: "argumentation", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-praemisse", DeckID: "c1-philosophy-ethics", Front: "die Prämisse", Back: "premise", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-schlussfolgerung", DeckID: "c1-philosophy-ethics", Front: "die Schlussfolgerung", Back: "conclusion", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-logik", DeckID: "c1-philosophy-ethics", Front: "die Logik", Back: "logic", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-wahrheit", DeckID: "c1-philosophy-ethics", Front: "die Wahrheit", Back: "truth", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-wirklichkeit", DeckID: "c1-philosophy-ethics", Front: "die Wirklichkeit", Back: "reality", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-schein", DeckID: "c1-philosophy-ethics", Front: "der Schein", Back: "appearance / illusion", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-transzendenz", DeckID: "c1-philosophy-ethics", Front: "die Transzendenz", Back: "transcendence", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-immanenz", DeckID: "c1-philosophy-ethics", Front: "die Immanenz", Back: "immanence", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-subjektivitaet", DeckID: "c1-philosophy-ethics", Front: "die Subjektivität", Back: "subjectivity", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-objektivitaet", DeckID: "c1-philosophy-ethics", Front: "die Objektivität", Back: "objectivity", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-dualismus", DeckID: "c1-philosophy-ethics", Front: "der Dualismus", Back: "dualism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-monismus", DeckID: "c1-philosophy-ethics", Front: "der Monismus", Back: "monism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-idealismus", DeckID: "c1-philosophy-ethics", Front: "der Idealismus", Back: "idealism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-materialismus", DeckID: "c1-philosophy-ethics", Front: "der Materialismus", Back: "materialism", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-phaenomenologie", DeckID: "c1-philosophy-ethics", Front: "die Phänomenologie", Back: "phenomenology", Tags: []string{"c1", "philosophy"}},
		{ID: "c1-phil-hermeneutik", DeckID: "c1-philosophy-ethics", Front: "die Hermeneutik", Back: "hermeneutics", Tags: []string{"c1", "philosophy"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-philosophy-ethics",
		Name:        "C1 German Philosophy & Ethics",
		Description: "Advanced vocabulary for philosophical discourse, ethical theories, and critical thinking.",
		Tags:        []string{"german", "c1", "philosophy", "ethics", "critical-thinking"},
		Notes:       notes,
	}
}
