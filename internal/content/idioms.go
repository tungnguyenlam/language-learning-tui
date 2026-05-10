package content

import "deutsch-tui/internal/core"

func IdiomsDeck() core.Deck {
	notes := []core.Note{
		{ID: "idiom-daumen-drucken", DeckID: "idioms", Front: "die Daumen drücken", Back: "to keep one's fingers crossed", Extra: "Example: Ich drücke dir die Daumen!", Tags: []string{"idiom", "common"}},
		{ID: "idiom-nur-bahnhof", DeckID: "idioms", Front: "nur Bahnhof verstehen", Back: "it's all Greek to me / to not understand anything", Extra: "Example: Ich verstehe nur Bahnhof.", Tags: []string{"idiom", "common"}},
		{ID: "idiom-tomaten-augen", DeckID: "idioms", Front: "Tomaten auf den Augen haben", Back: "to be oblivious to what's going on around you", Extra: "Literal: to have tomatoes on one's eyes", Tags: []string{"idiom"}},
		{ID: "idiom-kater-haben", DeckID: "idioms", Front: "einen Kater haben", Back: "to have a hangover", Extra: "Literal: to have a male cat", Tags: []string{"idiom"}},
		{ID: "idiom-schwein-haben", DeckID: "idioms", Front: "Schwein haben", Back: "to be lucky", Extra: "Example: Da hast du aber Schwein gehabt!", Tags: []string{"idiom"}},
		{ID: "idiom-alles-wurst", DeckID: "idioms", Front: "Es ist mir alles Wurst", Back: "I don't care / it's all the same to me", Extra: "Literal: It is all sausage to me", Tags: []string{"idiom"}},
		{ID: "idiom-senf-dazu", DeckID: "idioms", Front: "seinen Senf dazugeben", Back: "to give one's two cents / to put one's oar in", Extra: "Literal: to add one's mustard", Tags: []string{"idiom"}},
		{ID: "idiom-blatt-mund", DeckID: "idioms", Front: "kein Blatt vor den Mund nehmen", Back: "to not mince words / to be outspoken", Extra: "Literal: to not take a leaf before the mouth", Tags: []string{"idiom"}},
		{ID: "idiom-honig-bart", DeckID: "idioms", Front: "jemandem Honig um den Bart schmieren", Back: "to butter someone up", Extra: "Literal: to smear honey around someone's beard", Tags: []string{"idiom"}},
		{ID: "idiom-nase-voll", DeckID: "idioms", Front: "die Nase voll haben", Back: "to be fed up", Extra: "Example: Ich habe die Nase voll davon!", Tags: []string{"idiom"}},
		{ID: "idiom-fisch-wasser", DeckID: "idioms", Front: "wie ein Fisch im Wasser", Back: "like a fish in water / very comfortable", Tags: []string{"idiom"}},
		{ID: "idiom-pech-haben", DeckID: "idioms", Front: "Pech haben", Back: "to have bad luck", Tags: []string{"idiom"}},
		{ID: "idiom-ins-fettnapfchen", DeckID: "idioms", Front: "ins Fettnäpfchen treten", Back: "to put one's foot in it / to make a faux pas", Tags: []string{"idiom"}},
		{ID: "idiom-da-liegt-hund", DeckID: "idioms", Front: "Da liegt der Hund begraben", Back: "That's the heart of the matter / That's the snag", Tags: []string{"idiom"}},
		{ID: "idiom-kalt-warm", DeckID: "idioms", Front: "jemanden kalt oder warm lassen", Back: "to not care / to be indifferent", Tags: []string{"idiom"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "idioms",
		Name:        "German Idioms & Proverbs",
		Description: "Common German expressions that don't translate literally.",
		Tags:        []string{"german", "idioms", "culture"},
		Notes:       notes,
	}
}
