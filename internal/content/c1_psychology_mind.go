package content

import (
	"deutsch-tui/internal/core"
)

func C1PsychologyMindDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-mind-psyche", DeckID: "c1-psychology-mind", Front: "die Psyche", Back: "psyche / mind", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-bewusstsein", DeckID: "c1-psychology-mind", Front: "das Bewusstsein", Back: "consciousness", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-unbewusstes", DeckID: "c1-psychology-mind", Front: "das Unbewusste", Back: "unconscious (mind)", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-gedächtnis", DeckID: "c1-psychology-mind", Front: "das Gedächtnis", Back: "memory", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-langzeitgedächtnis", DeckID: "c1-psychology-mind", Front: "das Langzeitgedächtnis", Back: "long-term memory", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-kurzzeitgedächtnis", DeckID: "c1-psychology-mind", Front: "das Kurzzeitgedächtnis", Back: "short-term memory", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-wahrnehmung", DeckID: "c1-psychology-mind", Front: "die Wahrnehmung", Back: "perception", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-aufmerksamkeit", DeckID: "c1-psychology-mind", Front: "die Aufmerksamkeit", Back: "attention", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-konzentration", DeckID: "c1-psychology-mind", Front: "die Konzentration", Back: "concentration", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-kognition", DeckID: "c1-psychology-mind", Front: "die Kognition", Back: "cognition", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-denken", DeckID: "c1-psychology-mind", Front: "das Denken", Back: "thinking", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-logik", DeckID: "c1-psychology-mind", Front: "die Logik", Back: "logic", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-abstrakt", DeckID: "c1-psychology-mind", Front: "abstrakt", Back: "abstract", Tags: []string{"b2", "psychology", "adj"}},
		{ID: "c1-mind-konkrekt", DeckID: "c1-psychology-mind", Front: "konkret", Back: "concrete", Tags: []string{"b2", "psychology", "adj"}},
		{ID: "c1-mind-intuition", DeckID: "c1-psychology-mind", Front: "die Intuition", Back: "intuition", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-gefühl", DeckID: "c1-psychology-mind", Front: "das Gefühl", Back: "feeling / emotion", Tags: []string{"a2", "psychology", "noun"}},
		{ID: "c1-mind-emotion", DeckID: "c1-psychology-mind", Front: "die Emotion", Back: "emotion", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-empfindung", DeckID: "c1-psychology-mind", Front: "die Empfindung", Back: "sensation / feeling", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-stimmung", DeckID: "c1-psychology-mind", Front: "die Stimmung", Back: "mood", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-gemüt", DeckID: "c1-psychology-mind", Front: "das Gemüt", Back: "temperament / soul", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-seele", DeckID: "c1-psychology-mind", Front: "die Seele", Back: "soul", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-geist", DeckID: "c1-psychology-mind", Front: "der Geist", Back: "mind / spirit", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-verstand", DeckID: "c1-psychology-mind", Front: "der Verstand", Back: "intellect / reason", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-vernunft", DeckID: "c1-psychology-mind", Front: "die Vernunft", Back: "reason", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-wille", DeckID: "c1-psychology-mind", Front: "der Wille", Back: "will", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-motivation", DeckID: "c1-psychology-mind", Front: "die Motivation", Back: "motivation", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-antrieb", DeckID: "c1-psychology-mind", Front: "der Antrieb", Back: "drive / impulse", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-trieb", DeckID: "c1-psychology-mind", Front: "der Trieb", Back: "urge / instinct", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-instinkt", DeckID: "c1-psychology-mind", Front: "der Instinkt", Back: "instinct", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-reflex", DeckID: "c1-psychology-mind", Front: "der Reflex", Back: "reflex", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-identität", DeckID: "c1-psychology-mind", Front: "die Identität", Back: "identity", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-persönlichkeit", DeckID: "c1-psychology-mind", Front: "die Persönlichkeit", Back: "personality", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-charakter", DeckID: "c1-psychology-mind", Front: "der Charakter", Back: "character", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-temperament", DeckID: "c1-psychology-mind", Front: "das Temperament", Back: "temperament", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-verhalten", DeckID: "c1-psychology-mind", Front: "das Verhalten", Back: "behavior", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-handlung", DeckID: "c1-psychology-mind", Front: "die Handlung", Back: "action / act", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-reaktion", DeckID: "c1-psychology-mind", Front: "die Reaktion", Back: "reaction", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-stimulus", DeckID: "c1-psychology-mind", Front: "der Reiz", Back: "stimulus", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-angst", DeckID: "c1-psychology-mind", Front: "die Angst", Back: "fear / anxiety", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-phobie", DeckID: "c1-psychology-mind", Front: "die Phobie", Back: "phobia", Tags: []string{"c1", "psychology", "noun"}},
		{ID: "c1-mind-stress", DeckID: "c1-psychology-mind", Front: "der Stress", Back: "stress", Tags: []string{"b2", "psychology", "noun"}},
		{ID: "c1-mind-trauma", DeckID: "c1-psychology-mind", Front: "das Trauma", Back: "trauma", Tags: []string{"c1", "psychology", "noun"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "c1-psychology-mind",
		Name:        "German C1 Psychology & Mind",
		Description: "Advanced psychological and philosophical vocabulary for discussing the mind, emotions, and cognition.",
		Tags:        []string{"german", "c1", "psychology", "mind", "philosophy"},
		Notes:       notes,
	}
}
