package content

import (
	"deutsch-tui/internal/core"
	"fmt"
)

func PrepositionsDeck() core.Deck {
	notes := []core.Note{
		{ID: "prep-acc-durch", DeckID: "grammar-preps", Front: "durch", Back: "through", Tags: []string{"grammar", "preposition", "accusative"}},
		{ID: "prep-acc-fur", DeckID: "grammar-preps", Front: "für", Back: "for", Tags: []string{"grammar", "preposition", "accusative"}},
		{ID: "prep-acc-gegen", DeckID: "grammar-preps", Front: "gegen", Back: "against", Tags: []string{"grammar", "preposition", "accusative"}},
		{ID: "prep-acc-ohne", DeckID: "grammar-preps", Front: "ohne", Back: "without", Tags: []string{"grammar", "preposition", "accusative"}},
		{ID: "prep-acc-um", DeckID: "grammar-preps", Front: "um", Back: "around / at (time)", Tags: []string{"grammar", "preposition", "accusative"}},

		{ID: "prep-dat-aus", DeckID: "grammar-preps", Front: "aus", Back: "out of / from", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-bei", DeckID: "grammar-preps", Front: "bei", Back: "at / near / with", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-mit", DeckID: "grammar-preps", Front: "mit", Back: "with", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-nach", DeckID: "grammar-preps", Front: "nach", Back: "after / to (cities/countries)", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-seit", DeckID: "grammar-preps", Front: "seit", Back: "since / for (time)", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-von", DeckID: "grammar-preps", Front: "von", Back: "from / of", Tags: []string{"grammar", "preposition", "dative"}},
		{ID: "prep-dat-zu", DeckID: "grammar-preps", Front: "zu", Back: "to", Tags: []string{"grammar", "preposition", "dative"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])

		expectedCase := "accusative"
		isDative := false
		for _, tag := range notes[i].Tags {
			if tag == "dative" {
				isDative = true
				break
			}
		}
		if isDative {
			expectedCase = "dative"
		}

		notes[i].Cards = append(notes[i].Cards, core.Card{
			ID:      notes[i].ID + ":mcq-case",
			NoteID:  notes[i].ID,
			DeckID:  notes[i].DeckID,
			Kind:    core.CardKindMCQ,
			Prompt:  fmt.Sprintf("Which case does the preposition '%s' take?", notes[i].Front),
			Answer:  expectedCase,
			Choices: []string{"accusative", "dative", "genitive"},
			Tags:    notes[i].Tags,
		})
	}

	return core.Deck{
		ID:          "grammar-preps",
		Name:        "German Prepositions",
		Description: "Master Accusative and Dative prepositions.",
		Tags:        []string{"grammar", "prepositions"},
		Notes:       notes,
	}
}
