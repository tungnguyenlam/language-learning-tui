package content

import (
	"deutsch-tui/internal/core"
)

func ProverbsDeck() core.Deck {
	notes := []core.Note{
		{ID: "prov-uebung", DeckID: "german-proverbs", Front: "Übung macht den Meister.", Back: "Practice makes perfect.", Extra: "Literally: Practice makes the master.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-anfang", DeckID: "german-proverbs", Front: "Aller Anfang ist schwer.", Back: "Every beginning is difficult.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-morgenstund", DeckID: "german-proverbs", Front: "Morgenstund hat Gold im Mund.", Back: "The early bird catches the worm.", Extra: "Literally: Morning hour has gold in its mouth.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-geduld", DeckID: "german-proverbs", Front: "Geduld bringt Rosen.", Back: "Patience is a virtue.", Extra: "Literally: Patience brings roses.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-ende", DeckID: "german-proverbs", Front: "Ende gut, alles gut.", Back: "All's well that ends well.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-kleinvieh", DeckID: "german-proverbs", Front: "Kleinvieh macht auch Mist.", Back: "Every little bit helps.", Extra: "Literally: Small cattle also makes manure.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-wer-rastet", DeckID: "german-proverbs", Front: "Wer rastet, der rostet.", Back: "Use it or lose it.", Extra: "Literally: He who rests, rusts.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-schweigen", DeckID: "german-proverbs", Front: "Reden ist Silber, Schweigen ist Gold.", Back: "Speech is silver, silence is golden.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-appetit", DeckID: "german-proverbs", Front: "Appetit kommt beim Essen.", Back: "Appetite comes with eating.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-daumen", DeckID: "german-proverbs", Front: "Daumen drücken!", Back: "Fingers crossed!", Extra: "Literally: Press thumbs!", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-pech", DeckID: "german-proverbs", Front: "Pech im Glück, Glück im Unglück.", Back: "A blessing in disguise.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-hochmut", DeckID: "german-proverbs", Front: "Hochmut kommt vor dem Fall.", Back: "Pride comes before a fall.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-hunger", DeckID: "german-proverbs", Front: "Hunger ist der beste Koch.", Back: "Hunger is the best sauce.", Extra: "Literally: Hunger is the best cook.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-liebe", DeckID: "german-proverbs", Front: "Liebe geht durch den Magen.", Back: "The way to a man's heart is through his stomach.", Extra: "Literally: Love goes through the stomach.", Tags: []string{"proverb", "idiom"}},
		{ID: "prov-not", DeckID: "german-proverbs", Front: "Not macht erfinderisch.", Back: "Necessity is the mother of invention.", Tags: []string{"proverb", "idiom"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "german-proverbs",
		Name:        "German Proverbs & Idioms",
		Description: "Common German proverbs and their English equivalents.",
		Tags:        []string{"german", "proverb", "idiom", "culture"},
		Notes:       notes,
	}
}
