package content

import (
	"deutsch-tui/internal/core"
)

func CommonVerbsDeck() core.Deck {
	notes := []core.Note{
		{ID: "v-machen", Front: "machen", Back: "to do / to make", Extra: "ich mache, du machst, er macht", Tags: []string{"a1", "verb"}},
		{ID: "v-sehen", Front: "sehen", Back: "to see", Extra: "ich sehe, du siehst, er sieht", Tags: []string{"a1", "verb"}},
		{ID: "v-essen", Front: "essen", Back: "to eat", Extra: "ich esse, du isst, er isst", Tags: []string{"a1", "verb"}},
		{ID: "v-trinken", Front: "trinken", Back: "to drink", Tags: []string{"a1", "verb"}},
		{ID: "v-schreiben", Front: "schreiben", Back: "to write", Tags: []string{"a1", "verb"}},
		{ID: "v-lesen", Front: "lesen", Back: "to read", Extra: "ich lese, du liest, er liest", Tags: []string{"a1", "verb"}},
		{ID: "v-schlafen", Front: "schlafen", Back: "to sleep", Extra: "ich schlafe, du schläfst, er schläft", Tags: []string{"a1", "verb"}},
		{ID: "v-fahren", Front: "fahren", Back: "to drive / to travel", Extra: "ich fahre, du fährst, er fährt", Tags: []string{"a1", "verb"}},
		{ID: "v-hören", Front: "hören", Back: "to hear / to listen", Tags: []string{"a1", "verb"}},
		{ID: "v-sagen", Front: "sagen", Back: "to say", Tags: []string{"a1", "verb"}},
		{ID: "v-fragen", Front: "fragen", Back: "to ask", Tags: []string{"a1", "verb"}},
		{ID: "v-antworten", Front: "antworten", Back: "to answer", Tags: []string{"a1", "verb"}},
		{ID: "v-arbeiten", Front: "arbeiten", Back: "to work", Tags: []string{"a1", "verb"}},
		{ID: "v-spielen", Front: "spielen", Back: "to play", Tags: []string{"a1", "verb"}},
		{ID: "v-kaufen", Front: "kaufen", Back: "to buy", Tags: []string{"a1", "verb"}},
		{ID: "v-finden", Front: "finden", Back: "to find", Tags: []string{"a1", "verb"}},
		{ID: "v-geben", Front: "geben", Back: "to give", Extra: "ich gebe, du gibst, er gibt", Tags: []string{"a1", "verb"}},
		{ID: "v-denken", Front: "denken", Back: "to think", Tags: []string{"a1", "verb"}},
		{ID: "v-wissen", Front: "wissen", Back: "to know (facts)", Extra: "ich weiß, du weißt, er weiß", Tags: []string{"a1", "verb"}},
		{ID: "v-kennen", Front: "kennen", Back: "to know (people/places)", Tags: []string{"a1", "verb"}},
	}

	for i := range notes {
		notes[i].DeckID = "common-verbs"
		notes[i].Cards = CardsForNote(notes[i])

		// Add some conjugation MCQs
		if notes[i].ID == "v-essen" {
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-conj",
				NoteID:  notes[i].ID,
				DeckID:  "common-verbs",
				Kind:    core.CardKindMCQ,
				Prompt:  "What is the correct form: 'Er ___ einen Apfel.'?",
				Answer:  "isst",
				Choices: []string{"esse", "esst", "isst"},
				Tags:    notes[i].Tags,
			})
		}
		if notes[i].ID == "v-sehen" {
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-conj",
				NoteID:  notes[i].ID,
				DeckID:  "common-verbs",
				Kind:    core.CardKindMCQ,
				Prompt:  "What is the correct form: 'Du ___ den Film.'?",
				Answer:  "siehst",
				Choices: []string{"sehe", "siehst", "seht"},
				Tags:    notes[i].Tags,
			})
		}
	}

	return core.Deck{
		ID:          "common-verbs",
		Name:        "Common German Verbs (A1-A2)",
		Description: "Essential German verbs for daily conversation.",
		Tags:        []string{"german", "a1", "a2", "verbs"},
		Notes:       notes,
	}
}
