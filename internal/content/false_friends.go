package content

import (
	"deutsch-tui/internal/core"
)

func FalseFriendsDeck() core.Deck {
	notes := []core.Note{
		{
			ID:     "false-gift",
			DeckID: "false-friends-mastery",
			Front:  "das Gift",
			Back:   "poison",
			Extra:  "NOT 'gift' (Geschenk). Example: Vorsicht, das ist Gift!",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-eventuell",
			DeckID: "false-friends-mastery",
			Front:  "eventuell",
			Back:   "possibly / perhaps",
			Extra:  "NOT 'eventually' (schließlich/letztendlich).",
			Tags:   []string{"false-friend", "adverb"},
		},
		{
			ID:     "false-aktuell",
			DeckID: "false-friends-mastery",
			Front:  "aktuell",
			Back:   "current / up-to-date",
			Extra:  "NOT 'actually' (eigentlich/tatsächlich).",
			Tags:   []string{"false-friend", "adjective"},
		},
		{
			ID:     "false-chef",
			DeckID: "false-friends-mastery",
			Front:  "der Chef",
			Back:   "boss / head of department",
			Extra:  "NOT 'chef' (der Koch).",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-billion",
			DeckID: "false-friends-mastery",
			Front:  "die Billion",
			Back:   "trillion",
			Extra:  "NOT 'billion' (die Milliarde). German uses: Million, Milliarde, Billion, Billiarde.",
			Tags:   []string{"false-friend", "number"},
		},
		{
			ID:     "false-brave",
			DeckID: "false-friends-mastery",
			Front:  "brav",
			Back:   "well-behaved / good (usually for children/pets)",
			Extra:  "NOT 'brave' (mutig/tapfer).",
			Tags:   []string{"false-friend", "adjective"},
		},
		{
			ID:     "false-promotion",
			DeckID: "false-friends-mastery",
			Front:  "die Promotion",
			Back:   "obtaining a doctorate (PhD)",
			Extra:  "NOT 'promotion' at work (Beförderung).",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-rat",
			DeckID: "false-friends-mastery",
			Front:  "der Rat",
			Back:   "advice / counsel",
			Extra:  "NOT 'rat' (die Ratte).",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-see",
			DeckID: "false-friends-mastery",
			Front:  "der See vs die See",
			Back:   "der See (lake) vs die See (sea/ocean)",
			Extra:  "English 'sea' is usually 'das Meer' or 'die See'.",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-mitschreiben",
			DeckID: "false-friends-mastery",
			Front:  "mitschreiben",
			Back:   "to take notes",
			Extra:  "NOT 'to write with' (instrumental).",
			Tags:   []string{"false-friend", "verb"},
		},
		{
			ID:     "false-gymnasium",
			DeckID: "false-friends-mastery",
			Front:  "das Gymnasium",
			Back:   "grammar school / preparatory high school",
			Extra:  "NOT 'gym' (die Turnhalle / das Fitnessstudio).",
			Tags:   []string{"false-friend", "noun"},
		},
		{
			ID:     "false-bekommen",
			DeckID: "false-friends-mastery",
			Front:  "bekommen",
			Back:   "to receive / to get",
			Extra:  "NOT 'to become' (werden).",
			Tags:   []string{"false-friend", "verb"},
		},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "false-friends-mastery",
		Name:        "False Friends Mastery",
		Description: "Watch out for these 'False Friends' (Falsche Freunde) - words that look similar to English but mean something else.",
		Tags:        []string{"german", "vocabulary", "false-friends"},
		Notes:       notes,
	}
}
