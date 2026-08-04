package content

import (
	"deutsch-tui/internal/core"
)

func B2MediaJournalismDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-medj-presse", DeckID: "b2-media-journalism", Front: "die Presse", Back: "the press", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-pressefreiheit", DeckID: "b2-media-journalism", Front: "die Pressefreiheit", Back: "freedom of the press", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-journalismus", DeckID: "b2-media-journalism", Front: "der Journalismus", Back: "journalism", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-journalist", DeckID: "b2-media-journalism", Front: "der Journalist / die Journalistin", Back: "journalist", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-redakteur", DeckID: "b2-media-journalism", Front: "der Redakteur / die Redakteurin", Back: "editor", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-redaktion", DeckID: "b2-media-journalism", Front: "die Redaktion", Back: "editorial team / editorial office", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-chefredakteur", DeckID: "b2-media-journalism", Front: "der Chefredakteur / die Chefredakteurin", Back: "editor-in-chief", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-korrespondent", DeckID: "b2-media-journalism", Front: "der Korrespondent / die Korrespondentin", Back: "correspondent", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-moderator", DeckID: "b2-media-journalism", Front: "der Moderator / die Moderatorin", Back: "presenter / host", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-schlagzeile", DeckID: "b2-media-journalism", Front: "die Schlagzeile", Back: "headline", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-leitartikel", DeckID: "b2-media-journalism", Front: "der Leitartikel", Back: "editorial (lead opinion article)", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-reportage", DeckID: "b2-media-journalism", Front: "die Reportage", Back: "feature report", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-interview", DeckID: "b2-media-journalism", Front: "das Interview", Back: "interview", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-recherche", DeckID: "b2-media-journalism", Front: "die Recherche", Back: "research / investigation (journalistic)", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-quelle", DeckID: "b2-media-journalism", Front: "die Quelle", Back: "source", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-quellenschutz", DeckID: "b2-media-journalism", Front: "der Quellenschutz", Back: "protection of sources", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-meldung", DeckID: "b2-media-journalism", Front: "die Meldung", Back: "news report / announcement", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-berichterstattung", DeckID: "b2-media-journalism", Front: "die Berichterstattung", Back: "coverage / reporting", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-kommentar", DeckID: "b2-media-journalism", Front: "der Kommentar", Back: "commentary", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-tageszeitung", DeckID: "b2-media-journalism", Front: "die Tageszeitung", Back: "daily newspaper", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-wochenzeitung", DeckID: "b2-media-journalism", Front: "die Wochenzeitung", Back: "weekly newspaper", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-zeitschrift", DeckID: "b2-media-journalism", Front: "die Zeitschrift", Back: "magazine / periodical", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-boulevardblatt", DeckID: "b2-media-journalism", Front: "das Boulevardblatt", Back: "tabloid", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-auflage", DeckID: "b2-media-journalism", Front: "die Auflage", Back: "circulation / print run", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-abonnement", DeckID: "b2-media-journalism", Front: "das Abonnement", Back: "subscription", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-impressum", DeckID: "b2-media-journalism", Front: "das Impressum", Back: "imprint (legal publisher details)", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-rundfunk", DeckID: "b2-media-journalism", Front: "der öffentlich-rechtliche Rundfunk", Back: "public broadcasting", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-sender", DeckID: "b2-media-journalism", Front: "der Sender", Back: "broadcaster / station", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-sendung", DeckID: "b2-media-journalism", Front: "die Sendung", Back: "programme / broadcast", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-uebertragung", DeckID: "b2-media-journalism", Front: "die Live-Übertragung", Back: "live broadcast", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-pressekonferenz", DeckID: "b2-media-journalism", Front: "die Pressekonferenz", Back: "press conference", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-pressemitteilung", DeckID: "b2-media-journalism", Front: "die Pressemitteilung", Back: "press release", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-meinungsfreiheit", DeckID: "b2-media-journalism", Front: "die Meinungsfreiheit", Back: "freedom of opinion / free speech", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-zensur", DeckID: "b2-media-journalism", Front: "die Zensur", Back: "censorship", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-falschmeldung", DeckID: "b2-media-journalism", Front: "die Falschmeldung", Back: "false report / fake news story", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-desinformation", DeckID: "b2-media-journalism", Front: "die Desinformation", Back: "disinformation", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-glaubwuerdigkeit", DeckID: "b2-media-journalism", Front: "die Glaubwürdigkeit", Back: "credibility", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-objektivitaet", DeckID: "b2-media-journalism", Front: "die Objektivität", Back: "objectivity", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-enthuellung", DeckID: "b2-media-journalism", Front: "die Enthüllung", Back: "revelation / exposé", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-skandal", DeckID: "b2-media-journalism", Front: "der Skandal", Back: "scandal", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-reichweite", DeckID: "b2-media-journalism", Front: "die Reichweite", Back: "reach / audience size", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-leserschaft", DeckID: "b2-media-journalism", Front: "die Leserschaft", Back: "readership", Tags: []string{"b2", "media"}},
		{ID: "b2-medj-berichten", DeckID: "b2-media-journalism", Front: "berichten (über + Akk.)", Back: "to report (on something)", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-recherchieren", DeckID: "b2-media-journalism", Front: "recherchieren", Back: "to research / to investigate", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-veroeffentlichen", DeckID: "b2-media-journalism", Front: "veröffentlichen", Back: "to publish", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-verbreiten", DeckID: "b2-media-journalism", Front: "verbreiten", Back: "to spread / to disseminate", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-zensieren", DeckID: "b2-media-journalism", Front: "zensieren", Back: "to censor", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-enthuellen", DeckID: "b2-media-journalism", Front: "enthüllen", Back: "to reveal / to expose", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-kommentieren", DeckID: "b2-media-journalism", Front: "kommentieren", Back: "to comment on", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-dementieren", DeckID: "b2-media-journalism", Front: "dementieren", Back: "to deny / to issue a denial", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-zitieren", DeckID: "b2-media-journalism", Front: "zitieren", Back: "to quote", Tags: []string{"b2", "media", "verb"}},
		{ID: "b2-medj-glaubwuerdig", DeckID: "b2-media-journalism", Front: "glaubwürdig", Back: "credible", Tags: []string{"b2", "media", "adjective"}},
		{ID: "b2-medj-tendenzioes", DeckID: "b2-media-journalism", Front: "tendenziös", Back: "biased / tendentious", Tags: []string{"b2", "media", "adjective"}},
		{ID: "b2-medj-reisserisch", DeckID: "b2-media-journalism", Front: "reißerisch", Back: "sensationalist / lurid", Tags: []string{"b2", "media", "adjective"}},
		{ID: "b2-medj-unabhaengig", DeckID: "b2-media-journalism", Front: "unabhängig", Back: "independent", Tags: []string{"b2", "media", "adjective"}},
		{ID: "b2-medj-sachlich", DeckID: "b2-media-journalism", Front: "sachlich", Back: "factual / matter-of-fact", Tags: []string{"b2", "media", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-media-journalism",
		Name:        "B2 Medien und Journalismus",
		Description: "Media and journalism vocabulary: press, broadcasting, reporting and media ethics.",
		Tags:        []string{"german", "b2", "media", "journalism"},
		Notes:       notes,
	}
}
