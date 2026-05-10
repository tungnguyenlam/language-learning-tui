package content

import (
	"deutsch-tui/internal/core"
)

func ConfusableWordsDeck() core.Deck {
	notes := []core.Note{
		{
			ID:     "conf-wissen-kennen-1",
			DeckID: "confusable-words",
			Front:  "wissen",
			Back:   "to know (facts, information)",
			Extra:  "ich weiß, du weißt, er weiß, wir wissen... (irregular)",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-wissen-kennen-2",
			DeckID: "confusable-words",
			Front:  "kennen",
			Back:   "to know (people, places, familiarity)",
			Extra:  "Ich kenne ihn. Ich kenne diesen Ort.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-liegen-legen-1",
			DeckID: "confusable-words",
			Front:  "liegen",
			Back:   "to lie / be situated (location, Dativ)",
			Extra:  "Das Buch liegt auf dem Tisch.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-liegen-legen-2",
			DeckID: "confusable-words",
			Front:  "legen",
			Back:   "to lay / put down (movement, Akkusativ)",
			Extra:  "Ich lege das Buch auf den Tisch.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-stehen-stellen-1",
			DeckID: "confusable-words",
			Front:  "stehen",
			Back:   "to stand / be upright (location, Dativ)",
			Extra:  "Die Vase steht im Regal.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-stehen-stellen-2",
			DeckID: "confusable-words",
			Front:  "stellen",
			Back:   "to place / set upright (movement, Akkusativ)",
			Extra:  "Ich stelle die Vase ins Regal.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-sitzen-setzen-1",
			DeckID: "confusable-words",
			Front:  "sitzen",
			Back:   "to sit / be seated (location, Dativ)",
			Extra:  "Er sitzt auf dem Stuhl.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-sitzen-setzen-2",
			DeckID: "confusable-words",
			Front:  "setzen",
			Back:   "to set / seat oneself (movement, Akkusativ)",
			Extra:  "Er setzt sich auf den Stuhl.",
			Tags:   []string{"confusable", "verb"},
		},
		{
			ID:     "conf-viel-viele",
			DeckID: "confusable-words",
			Front:  "viel vs viele",
			Back:   "viel (uncountable) vs viele (countable)",
			Extra:  "viel Wasser, viele Äpfel",
			Tags:   []string{"confusable", "adjective"},
		},
		{
			ID:     "conf-ganz-alle",
			DeckID: "confusable-words",
			Front:  "ganz vs alle",
			Back:   "ganz (whole/entire) vs alle (everyone/all items)",
			Extra:  "der ganze Tag, alle Leute",
			Tags:   []string{"confusable", "adjective"},
		},
		{
			ID:      "conf-mcq-wissen-kennen",
			DeckID:  "confusable-words",
			Front:   "Ich ___ nicht, wo er ist.",
			Back:    "weiß",
			Choices: []string{"weiß", "kenne", "wissen", "kennen"},
			Tags:    []string{"confusable", "mcq"},
		},
		{
			ID:      "conf-mcq-kennen-wissen",
			DeckID:  "confusable-words",
			Front:   "___ du Berlin gut?",
			Back:    "Kennst",
			Choices: []string{"Kennst", "Weißt", "Kennen", "Wissen"},
			Tags:    []string{"confusable", "mcq"},
		},
		{
			ID:     "conf-cloze-liegen",
			DeckID: "confusable-words",
			Front:  "Der Hund {{c1::liegt}} auf dem Sofa.",
			Back:   "liegt",
			Extra:  "Dativ location (auf dem Sofa)",
			Tags:   []string{"confusable", "cloze"},
		},
		{
			ID:     "conf-cloze-legen",
			DeckID: "confusable-words",
			Front:  "Ich {{c1::lege}} die Decke auf das Sofa.",
			Back:   "lege",
			Extra:  "Akkusativ movement (auf das Sofa)",
			Tags:   []string{"confusable", "cloze"},
		},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "confusable-words",
		Name:        "German Confusable Words",
		Description: "Master tricky pairs like wissen/kennen, liegen/legen, and more.",
		Tags:        []string{"german", "vocabulary", "grammar"},
		Notes:       notes,
	}
}
