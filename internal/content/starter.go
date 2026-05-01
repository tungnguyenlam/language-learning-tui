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
		// MCQ cards for German articles
		{
			ID:     "a1-mcq-articles-1",
			DeckID: "a1-survival",
			Front:  "What is the article for 'Haus' (house)?",
			Back:   "das",
			Extra:  "Haus is neuter in German",
			Tags:   []string{"a1", "mcq", "article"},
			Examples: []string{
				"Das Haus ist gross. (The house is big.)",
			},
		},
		{
			ID:     "a1-mcq-articles-2",
			DeckID: "a1-survival",
			Front:  "What is the article for 'Frau' (woman)?",
			Back:   "die",
			Extra:  "Frau is feminine in German",
			Tags:   []string{"a1", "mcq", "article"},
			Examples: []string{
				"Die Frau liest ein Buch. (The woman reads a book.)",
			},
		},
		{
			ID:     "a1-mcq-articles-3",
			DeckID: "a1-survival",
			Front:  "What is the article for 'Mann' (man)?",
			Back:   "der",
			Extra:  "Mann is masculine in German",
			Tags:   []string{"a1", "mcq", "article"},
			Examples: []string{
				"Der Mann spielt Fussball. (The man plays soccer.)",
			},
		},
		// MCQ cards for German verb conjugations
		{
			ID:     "a1-mcq-conjugation-1",
			DeckID: "a1-survival",
			Front:  "What is the correct form: 'Ich ___ Deutsch.'?",
			Back:   "lerne",
			Extra:  "lernen (to learn) - ich lerne, du lernst, er/sie/es lernt",
			Tags:   []string{"a1", "mcq", "verb"},
			Examples: []string{
				"Ich lerne Deutsch. (I learn German.)",
			},
		},
		{
			ID:     "a1-mcq-conjugation-2",
			DeckID: "a1-survival",
			Front:  "What is the correct form: 'Du ___ Musik.'?",
			Back:   "hoerst",
			Extra:  "hoeren (to listen/hear) - ich hoere, du hoerst, er/sie/es hoert",
			Tags:   []string{"a1", "mcq", "verb"},
			Examples: []string{
				"Du hoerst Musik. (You listen to music.)",
			},
		},
	}
	for i := range notes {
		// Generate standard cards
		notes[i].Cards = CardsForNote(notes[i])

		// Add custom MCQ cards for our specific notes
		switch notes[i].ID {
		case "a1-mcq-articles-1":
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-custom",
				NoteID:  notes[i].ID,
				DeckID:  notes[i].DeckID,
				Kind:    core.CardKindMCQ,
				Prompt:  notes[i].Front,
				Answer:  notes[i].Back,
				Choices: []string{"das", "der", "die"},
				Tags:    notes[i].Tags,
			})
		case "a1-mcq-articles-2":
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-custom",
				NoteID:  notes[i].ID,
				DeckID:  notes[i].DeckID,
				Kind:    core.CardKindMCQ,
				Prompt:  notes[i].Front,
				Answer:  notes[i].Back,
				Choices: []string{"das", "der", "die"},
				Tags:    notes[i].Tags,
			})
		case "a1-mcq-articles-3":
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-custom",
				NoteID:  notes[i].ID,
				DeckID:  notes[i].DeckID,
				Kind:    core.CardKindMCQ,
				Prompt:  notes[i].Front,
				Answer:  notes[i].Back,
				Choices: []string{"das", "der", "die"},
				Tags:    notes[i].Tags,
			})
		case "a1-mcq-conjugation-1":
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-custom",
				NoteID:  notes[i].ID,
				DeckID:  notes[i].DeckID,
				Kind:    core.CardKindMCQ,
				Prompt:  notes[i].Front,
				Answer:  notes[i].Back,
				Choices: []string{"lerne", "lernst", "lernt"},
				Tags:    notes[i].Tags,
			})
		case "a1-mcq-conjugation-2":
			notes[i].Cards = append(notes[i].Cards, core.Card{
				ID:      notes[i].ID + ":mcq-custom",
				NoteID:  notes[i].ID,
				DeckID:  notes[i].DeckID,
				Kind:    core.CardKindMCQ,
				Prompt:  notes[i].Front,
				Answer:  notes[i].Back,
				Choices: []string{"hoere", "hoerst", "hoert"},
				Tags:    notes[i].Tags,
			})
		}
	}
	return core.Deck{
		ID:          "a1-survival",
		Name:        "German A1 Survival",
		Description: "Starter vocabulary and phrases for daily German.",
		Tags:        []string{"german", "a1"},
		Notes:       notes,
	}
}
