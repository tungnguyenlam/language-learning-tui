package content

import (
	"deutsch-tui/internal/core"
	"fmt"
)

func StarterDeck() core.Deck {
	notes := []core.Note{
		// --- Greetings & Essentials ---
		{ID: "a1-ess-hallo", DeckID: "a1-survival", Front: "Hallo", Back: "hello", Tags: []string{"a1", "greeting"}},
		{ID: "a1-ess-tschuss", DeckID: "a1-survival", Front: "Tschüss", Back: "bye", Tags: []string{"a1", "greeting"}},
		{ID: "a1-ess-bitte", DeckID: "a1-survival", Front: "Bitte", Back: "please / you're welcome", Tags: []string{"a1", "essential"}},
		{ID: "a1-ess-danke", DeckID: "a1-survival", Front: "Danke", Back: "thank you", Tags: []string{"a1", "essential"}},
		{ID: "a1-ess-ja", DeckID: "a1-survival", Front: "Ja", Back: "yes", Tags: []string{"a1", "essential"}},
		{ID: "a1-ess-nein", DeckID: "a1-survival", Front: "Nein", Back: "no", Tags: []string{"a1", "essential"}},
		{ID: "a1-ess-entschuldigung", DeckID: "a1-survival", Front: "Entschuldigung", Back: "excuse me / sorry", Tags: []string{"a1", "essential"}},

		// --- Nouns: Food & Drink ---
		{ID: "a1-food-apfel", DeckID: "a1-survival", Front: "der Apfel", Back: "apple", Extra: "Plural: die Äpfel", Tags: []string{"a1", "noun", "food"}},
		{ID: "a1-food-brot", DeckID: "a1-survival", Front: "das Brot", Back: "bread", Extra: "Plural: die Brote", Tags: []string{"a1", "noun", "food"}},
		{ID: "a1-food-wasser", DeckID: "a1-survival", Front: "das Wasser", Back: "water", Tags: []string{"a1", "noun", "drink"}},
		{ID: "a1-food-kaffee", DeckID: "a1-survival", Front: "der Kaffee", Back: "coffee", Tags: []string{"a1", "noun", "drink"}},
		{ID: "a1-food-milch", DeckID: "a1-survival", Front: "die Milch", Back: "milk", Tags: []string{"a1", "noun", "drink"}},
		{ID: "a1-food-kase", DeckID: "a1-survival", Front: "der Käse", Back: "cheese", Tags: []string{"a1", "noun", "food"}},
		{ID: "a1-food-fleisch", DeckID: "a1-survival", Front: "das Fleisch", Back: "meat", Tags: []string{"a1", "noun", "food"}},

		// --- Verbs ---
		{ID: "a1-verb-sein", DeckID: "a1-survival", Front: "sein", Back: "to be", Extra: "ich bin, du bist, er ist, wir sind, ihr seid, sie sind", Tags: []string{"a1", "verb"}},
		{ID: "a1-verb-haben", DeckID: "a1-survival", Front: "haben", Back: "to have", Extra: "ich habe, du hast, er hat, wir haben, ihr habt, sie haben", Tags: []string{"a1", "verb"}},
		{ID: "a1-verb-gehen", DeckID: "a1-survival", Front: "gehen", Back: "to go", Tags: []string{"a1", "verb"}},
		{ID: "a1-verb-kommen", DeckID: "a1-survival", Front: "kommen", Back: "to come", Tags: []string{"a1", "verb"}},
		{ID: "a1-verb-lernen", DeckID: "a1-survival", Front: "lernen", Back: "to learn", Tags: []string{"a1", "verb"}},
		{ID: "a1-verb-sprechen", DeckID: "a1-survival", Front: "sprechen", Back: "to speak", Extra: "ich spreche, du sprichst, er spricht", Tags: []string{"a1", "verb"}},

		// --- Numbers ---
		{ID: "a1-num-0", DeckID: "a1-survival", Front: "null", Back: "zero", Tags: []string{"a1", "number"}},
		{ID: "a1-num-1", DeckID: "a1-survival", Front: "eins", Back: "one", Tags: []string{"a1", "number"}},
		{ID: "a1-num-2", DeckID: "a1-survival", Front: "zwei", Back: "two", Tags: []string{"a1", "number"}},
		{ID: "a1-num-3", DeckID: "a1-survival", Front: "drei", Back: "three", Tags: []string{"a1", "number"}},
		{ID: "a1-num-4", DeckID: "a1-survival", Front: "vier", Back: "four", Tags: []string{"a1", "number"}},
		{ID: "a1-num-5", DeckID: "a1-survival", Front: "fünf", Back: "five", Tags: []string{"a1", "number"}},
		{ID: "a1-num-10", DeckID: "a1-survival", Front: "zehn", Back: "ten", Tags: []string{"a1", "number"}},

		// --- Colors ---
		{ID: "a1-col-rot", DeckID: "a1-survival", Front: "rot", Back: "red", Tags: []string{"a1", "color"}},
		{ID: "a1-col-blau", DeckID: "a1-survival", Front: "blau", Back: "blue", Tags: []string{"a1", "color"}},
		{ID: "a1-col-grun", DeckID: "a1-survival", Front: "grün", Back: "green", Tags: []string{"a1", "color"}},
		{ID: "a1-col-gelb", DeckID: "a1-survival", Front: "gelb", Back: "yellow", Tags: []string{"a1", "color"}},
		{ID: "a1-col-schwarz", DeckID: "a1-survival", Front: "schwarz", Back: "black", Tags: []string{"a1", "color"}},
		{ID: "a1-col-weis", DeckID: "a1-survival", Front: "weiß", Back: "white", Tags: []string{"a1", "color"}},

		// --- People & Family ---
		{ID: "a1-fam-mann", DeckID: "a1-survival", Front: "der Mann", Back: "man", Tags: []string{"a1", "noun", "people"}},
		{ID: "a1-fam-frau", DeckID: "a1-survival", Front: "die Frau", Back: "woman", Tags: []string{"a1", "noun", "people"}},
		{ID: "a1-fam-kind", DeckID: "a1-survival", Front: "das Kind", Back: "child", Tags: []string{"a1", "noun", "people"}},
		{ID: "a1-fam-vater", DeckID: "a1-survival", Front: "der Vater", Back: "father", Tags: []string{"a1", "noun", "family"}},
		{ID: "a1-fam-mutter", DeckID: "a1-survival", Front: "die Mutter", Back: "mother", Tags: []string{"a1", "noun", "family"}},

		// --- Places ---
		{ID: "a1-pla-haus", DeckID: "a1-survival", Front: "das Haus", Back: "house", Tags: []string{"a1", "noun", "place"}},
		{ID: "a1-pla-stadt", DeckID: "a1-survival", Front: "die Stadt", Back: "city", Tags: []string{"a1", "noun", "place"}},
		{ID: "a1-pla-schule", DeckID: "a1-survival", Front: "die Schule", Back: "school", Tags: []string{"a1", "noun", "place"}},
		{ID: "a1-pla-supermarkt", DeckID: "a1-survival", Front: "der Supermarkt", Back: "supermarket", Tags: []string{"a1", "noun", "place"}},

		// --- Time ---
		{ID: "a1-tim-heute", DeckID: "a1-survival", Front: "heute", Back: "today", Tags: []string{"a1", "time"}},
		{ID: "a1-tim-morgen", DeckID: "a1-survival", Front: "morgen", Back: "tomorrow", Tags: []string{"a1", "time"}},
		{ID: "a1-tim-jetzt", DeckID: "a1-survival", Front: "jetzt", Back: "now", Tags: []string{"a1", "time"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])

		// Logic for automatically adding MCQ cards based on certain tags
		if contains(notes[i].Tags, "article") || contains(notes[i].Tags, "people") || notes[i].ID == "a1-pla-haus" {
			// Article MCQs
			article := ""
			if len(notes[i].Front) >= 4 && notes[i].Front[0:3] == "der" {
				article = "der"
			}
			if len(notes[i].Front) >= 4 && notes[i].Front[0:3] == "die" {
				article = "die"
			}
			if len(notes[i].Front) >= 4 && notes[i].Front[0:3] == "das" {
				article = "das"
			}

			if article != "" {
				noun := notes[i].Front[4:]
				notes[i].Cards = append(notes[i].Cards, core.Card{
					ID:      notes[i].ID + ":mcq-article",
					NoteID:  notes[i].ID,
					DeckID:  notes[i].DeckID,
					Kind:    core.CardKindMCQ,
					Prompt:  fmt.Sprintf("What is the article for '%s'?", noun),
					Answer:  article,
					Choices: []string{"der", "die", "das"},
					Tags:    notes[i].Tags,
				})
			}
		}

		if contains(notes[i].Tags, "verb") {
			// Basic conjugation MCQ for common verbs
			answer := ""
			prompt := ""
			choices := []string{}

			switch notes[i].ID {
			case "a1-verb-lernen":
				prompt = "What is the correct form: 'Ich ___ Deutsch.'?"
				answer = "lerne"
				choices = []string{"lerne", "lernst", "lernt"}
			case "a1-verb-kommen":
				prompt = "What is the correct form: 'Woher ___ du?'?"
				answer = "kommst"
				choices = []string{"komme", "kommst", "kommt"}
			case "a1-verb-sprechen":
				prompt = "What is the correct form: 'Er ___ gut Deutsch.'?"
				answer = "spricht"
				choices = []string{"spreche", "sprichst", "spricht"}
			}

			if prompt != "" {
				notes[i].Cards = append(notes[i].Cards, core.Card{
					ID:      notes[i].ID + ":mcq-conj",
					NoteID:  notes[i].ID,
					DeckID:  notes[i].DeckID,
					Kind:    core.CardKindMCQ,
					Prompt:  prompt,
					Answer:  answer,
					Choices: choices,
					Tags:    notes[i].Tags,
				})
			}
		}
	}

	return core.Deck{
		ID:          "a1-survival",
		Name:        "German A1 Survival",
		Description: "Comprehensive starter vocabulary and phrases for daily German.",
		Tags:        []string{"german", "a1"},
		Notes:       notes,
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
