package content

import "deutsch-tui/internal/core"

func SlangDeck() core.Deck {
	notes := []core.Note{
		{ID: "slang-geil", DeckID: "slang", Front: "geil", Back: "cool / awesome / (orig. horny)", Extra: "Very common youth slang for 'great'.", Tags: []string{"slang", "informal"}},
		{ID: "slang-krass", DeckID: "slang", Front: "krass", Back: "extreme / intense / sick", Extra: "Can be positive or negative depending on context.", Tags: []string{"slang", "informal"}},
		{ID: "slang-alter", DeckID: "slang", Front: "Alter", Back: "dude / mate", Extra: "Used like 'dude' at the start or end of sentences.", Tags: []string{"slang", "informal"}},
		{ID: "slang-bescheuert", DeckID: "slang", Front: "bescheuert", Back: "stupid / dumb", Extra: "Stronger than 'doof'.", Tags: []string{"slang", "informal"}},
		{ID: "slang-kein-bock", DeckID: "slang", Front: "keinen Bock haben", Back: "to not be in the mood / to not feel like it", Extra: "Example: Ich hab keinen Bock auf Hausaufgaben.", Tags: []string{"slang", "informal"}},
		{ID: "slang-chillig", DeckID: "slang", Front: "chillig", Back: "chilled / relaxed", Extra: "Adopted from English 'chilling'.", Tags: []string{"slang", "informal"}},
		{ID: "slang-lahm", DeckID: "slang", Front: "lahm", Back: "lame / boring", Extra: "Literal: lame/paralyzed.", Tags: []string{"slang", "informal"}},
		{ID: "slang-hammer", DeckID: "slang", Front: "der Hammer", Back: "the bomb / amazing", Extra: "Example: Das ist echt der Hammer!", Tags: []string{"slang", "informal"}},
		{ID: "slang-null-acht-fufzehn", DeckID: "slang", Front: "nullachtfünfzehn (08/15)", Back: "mediocre / run-of-the-mill", Extra: "Originates from a WWI machine gun model.", Tags: []string{"slang", "common"}},
		{ID: "slang-jein", DeckID: "slang", Front: "jein", Back: "yes and no", Extra: "Combination of 'ja' and 'nein'.", Tags: []string{"slang", "common"}},
		{ID: "slang-auf-jeden", DeckID: "slang", Front: "auf jeden Fall / auf jeden", Back: "definitely", Extra: "Often shortened to 'auf jeden'.", Tags: []string{"slang", "informal"}},
		{ID: "slang-passt-schon", DeckID: "slang", Front: "passt schon", Back: "it's okay / it fits", Extra: "Used to say something is 'good enough'.", Tags: []string{"slang", "common"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "slang",
		Name:        "German Slang & Youth Language",
		Description: "Informal expressions and modern slang used in daily life.",
		Tags:        []string{"german", "slang", "informal"},
		Notes:       notes,
	}
}
