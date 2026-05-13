package content

import (
	"deutsch-tui/internal/core"
)

func C1PoliticsDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-pol-regierung", DeckID: "c1-politics", Front: "die Regierung", Back: "government", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-parlament", DeckID: "c1-politics", Front: "das Parlament", Back: "parliament", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-wahl", DeckID: "c1-politics", Front: "die Wahl", Back: "election", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-waehlen", DeckID: "c1-politics", Front: "wählen", Back: "to vote / to elect", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-stimme", DeckID: "c1-politics", Front: "die Stimme", Back: "vote / voice", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-partei", DeckID: "c1-politics", Front: "die Partei", Back: "political party", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-kandidat", DeckID: "c1-politics", Front: "der Kandidat / die Kandidatin", Back: "candidate", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-koalition", DeckID: "c1-politics", Front: "die Koalition", Back: "coalition", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-opposition", DeckID: "c1-politics", Front: "die Opposition", Back: "opposition", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-gesetz", DeckID: "c1-politics", Front: "das Gesetz", Back: "law", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-verfassung", DeckID: "c1-politics", Front: "die Verfassung", Back: "constitution", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-bundestag", DeckID: "c1-politics", Front: "der Bundestag", Back: "German federal parliament", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-kanzler", DeckID: "c1-politics", Front: "der Kanzler / die Kanzlerin", Back: "chancellor", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-minister", DeckID: "c1-politics", Front: "der Minister / die Ministerin", Back: "minister", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-demokratie", DeckID: "c1-politics", Front: "die Demokratie", Back: "democracy", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-diktatur", DeckID: "c1-politics", Front: "die Diktatur", Back: "dictatorship", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-buerger", DeckID: "c1-politics", Front: "der Bürger / die Bürgerin", Back: "citizen", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-rechte", DeckID: "c1-politics", Front: "die Menschenrechte", Back: "human rights", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-mehrheit", DeckID: "c1-politics", Front: "die Mehrheit", Back: "majority", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-minderheit", DeckID: "c1-politics", Front: "die Minderheit", Back: "minority", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-abstimmung", DeckID: "c1-politics", Front: "die Abstimmung", Back: "vote / ballot", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-abgeordneter", DeckID: "c1-politics", Front: "der Abgeordnete / die Abgeordnete", Back: "member of parliament / representative", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-debatte", DeckID: "c1-politics", Front: "die Debatte", Back: "debate", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-reform", DeckID: "c1-politics", Front: "die Reform", Back: "reform", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-steuern", DeckID: "c1-politics", Front: "die Steuern", Back: "taxes", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-wirtschaft", DeckID: "c1-politics", Front: "die Wirtschaft", Back: "economy", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-aussenpolitik", DeckID: "c1-politics", Front: "die Außenpolitik", Back: "foreign policy", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-innenpolitik", DeckID: "c1-politics", Front: "die Innenpolitik", Back: "domestic policy", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-gewerkschaft", DeckID: "c1-politics", Front: "die Gewerkschaft", Back: "trade union", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-streik", DeckID: "c1-politics", Front: "der Streik", Back: "strike", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-demonstration", DeckID: "c1-politics", Front: "die Demonstration", Back: "demonstration / protest", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-forderung", DeckID: "c1-politics", Front: "die Forderung", Back: "demand / claim", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-kompromiss", DeckID: "c1-politics", Front: "der Kompromiss", Back: "compromise", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-verhandeln", DeckID: "c1-politics", Front: "verhandeln", Back: "to negotiate", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-vertrag", DeckID: "c1-politics", Front: "der Vertrag", Back: "treaty / contract", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-beschluss", DeckID: "c1-politics", Front: "der Beschluss", Back: "resolution / decision", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-krise", DeckID: "c1-politics", Front: "die Krise", Back: "crisis", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-frieden", DeckID: "c1-politics", Front: "der Frieden", Back: "peace", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-krieg", DeckID: "c1-politics", Front: "der Krieg", Back: "war", Tags: []string{"c1", "politics"}},
		{ID: "c1-pol-sicherheit", DeckID: "c1-politics", Front: "die Sicherheit", Back: "security / safety", Tags: []string{"c1", "politics"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-politics",
		Name:        "C1 German Politics & Society",
		Description: "Advanced political vocabulary for discussions, news, and societal topics.",
		Tags:        []string{"german", "c1", "politics", "society", "news"},
		Notes:       notes,
	}
}
