package content

import "deutsch-tui/internal/core"

func StarterDeck() core.Deck {
	notes := []core.Note{
		{
			ID:       "a1-essen-apfel",
			DeckID:   "a1-survival",
			Front:    "der Apfel",
			Back:     "apple",
			Extra:    "Plural: die Aepfel",
			Tags:     []string{"a1", "noun", "food"},
			Examples: []string{"Ich esse einen Apfel."},
		},
		{
			ID:       "a1-phrase-danke",
			DeckID:   "a1-survival",
			Front:    "Danke",
			Back:     "thank you",
			Extra:    "Polite everyday phrase.",
			Tags:     []string{"a1", "phrase"},
			Examples: []string{"Danke fuer deine Hilfe."},
		},
		{
			ID:       "a1-verb-gehen",
			DeckID:   "a1-survival",
			Front:    "gehen",
			Back:     "to go",
			Extra:    "ich gehe, du gehst, er/sie/es geht",
			Tags:     []string{"a1", "verb"},
			Examples: []string{"Wir gehen nach Hause."},
		},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "a1-survival",
		Name:        "German A1 Survival",
		Description: "Starter vocabulary and phrases for daily German.",
		Tags:        []string{"german", "a1"},
		Notes:       notes,
	}
}
