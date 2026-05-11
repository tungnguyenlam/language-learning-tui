package content

import (
	"deutsch-tui/internal/core"
)

func B2MediaDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-med-nachrichten", DeckID: "b2-media", Front: "die Nachrichten (Pl.)", Back: "news", Tags: []string{"b2", "media"}},
		{ID: "b2-med-zeitung", DeckID: "b2-media", Front: "die Zeitung", Back: "newspaper", Tags: []string{"b2", "media"}},
		{ID: "b2-med-zeitschrift", DeckID: "b2-media", Front: "die Zeitschrift", Back: "magazine", Tags: []string{"b2", "media"}},
		{ID: "b2-med-artikel", DeckID: "b2-media", Front: "der Artikel", Back: "article", Tags: []string{"b2", "media"}},
		{ID: "b2-med-kommentar", DeckID: "b2-media", Front: "der Kommentar", Back: "comment", Tags: []string{"b2", "media"}},
		{ID: "b2-med-sendung", DeckID: "b2-media", Front: "die Sendung", Back: "broadcast / program", Tags: []string{"b2", "media"}},
		{ID: "b2-med-sender", DeckID: "b2-media", Front: "der Sender", Back: "broadcaster / channel", Tags: []string{"b2", "media"}},
		{ID: "b2-med-bericht", DeckID: "b2-media", Front: "der Bericht", Back: "report", Tags: []string{"b2", "media"}},
		{ID: "b2-med-interview", DeckID: "b2-media", Front: "das Interview", Back: "interview", Tags: []string{"b2", "media"}},
		{ID: "b2-med-meldung", DeckID: "b2-media", Front: "die Meldung", Back: "report / news item", Tags: []string{"b2", "media"}},
		{ID: "b2-med-schlagzeile", DeckID: "b2-media", Front: "die Schlagzeile", Back: "headline", Tags: []string{"b2", "media"}},
		{ID: "b2-med-redakteur", DeckID: "b2-media", Front: "der Redakteur / die Redakteurin", Back: "editor", Tags: []string{"b2", "media"}},
		{ID: "b2-med-reporter", DeckID: "b2-media", Front: "der Reporter / die Reporterin", Back: "reporter", Tags: []string{"b2", "media"}},
		{ID: "b2-med-journalist", DeckID: "b2-media", Front: "der Journalist / die Journalistin", Back: "journalist", Tags: []string{"b2", "media"}},
		{ID: "b2-med-quelle", DeckID: "b2-media", Front: "die Quelle", Back: "source", Tags: []string{"b2", "media"}},
		{ID: "b2-med-thema", DeckID: "b2-media", Front: "das Thema", Back: "topic / theme", Tags: []string{"b2", "media"}},
		{ID: "b2-med-berichten", DeckID: "b2-media", Front: "berichten", Back: "to report", Tags: []string{"b2", "media"}},
		{ID: "b2-med-kritisieren", DeckID: "b2-media", Front: "kritisieren", Back: "to criticize", Tags: []string{"b2", "media"}},
		{ID: "b2-med-kontrovers", DeckID: "b2-media", Front: "kontrovers", Back: "controversial", Tags: []string{"b2", "media"}},
		{ID: "b2-med-objektiv", DeckID: "b2-media", Front: "objektiv", Back: "objective", Tags: []string{"b2", "media"}},
		{ID: "b2-med-subjektiv", DeckID: "b2-media", Front: "subjektiv", Back: "subjective", Tags: []string{"b2", "media"}},
		{ID: "b2-med-podcast", DeckID: "b2-media", Front: "der Podcast", Back: "podcast", Tags: []string{"b2", "media"}},
		{ID: "b2-med-blog", DeckID: "b2-media", Front: "der Blog", Back: "blog", Tags: []string{"b2", "media"}},
		{ID: "b2-med-social", DeckID: "b2-media", Front: "die sozialen Medien", Back: "social media", Tags: []string{"b2", "media"}},
		{ID: "b2-med-stream", DeckID: "b2-media", Front: "streamen", Back: "to stream", Tags: []string{"b2", "media"}},
		{ID: "b2-med-videoclip", DeckID: "b2-media", Front: "der Videoclip", Back: "video clip", Tags: []string{"b2", "media"}},
		{ID: "b2-med-werbung", DeckID: "b2-media", Front: "die Werbung", Back: "advertisement", Tags: []string{"b2", "media"}},
		{ID: "b2-med-fernsehgucker", DeckID: "b2-media", Front: "der Fernsehgucker", Back: "TV viewer", Tags: []string{"b2", "media"}},
		{ID: "b2-med-seriös", DeckID: "b2-media", Front: "seriös", Back: "serious / reputable", Tags: []string{"b2", "media"}},
		{ID: "b2-med-sensationell", DeckID: "b2-media", Front: "sensationell", Back: "sensational", Tags: []string{"b2", "media"}},
		{ID: "b2-med-glaubwürdig", DeckID: "b2-media", Front: "glaubwürdig", Back: "credible", Tags: []string{"b2", "media"}},
		{ID: "b2-med-unwahr", DeckID: "b2-media", Front: "unwahr", Back: "untrue", Tags: []string{"b2", "media"}},
		{ID: "b2-med-fake-news", DeckID: "b2-media", Front: "die Falschmeldung", Back: "fake news", Tags: []string{"b2", "media"}},
		{ID: "b2-med-korrespondent", DeckID: "b2-media", Front: "der Korrespondent / die Korrespondentin", Back: "correspondent", Tags: []string{"b2", "media"}},
		{ID: "b2-med-live", DeckID: "b2-media", Front: "live", Back: "live (broadcast)", Tags: []string{"b2", "media"}},
		{ID: "b2-med-ausland", DeckID: "b2-media", Front: "aus dem Ausland", Back: "from abroad", Tags: []string{"b2", "media"}},
		{ID: "b2-med-inland", DeckID: "b2-media", Front: "aus dem Inland", Back: "domestic (news)", Tags: []string{"b2", "media"}},
		{ID: "b2-med-ticker", DeckID: "b2-media", Front: "der Ticker", Back: "news ticker", Tags: []string{"b2", "media"}},
		{ID: "b2-med-digital", DeckID: "b2-media", Front: "digital", Back: "digital", Tags: []string{"b2", "media"}},
		{ID: "b2-med-pressekonferenz", DeckID: "b2-media", Front: "die Pressekonferenz", Back: "press conference", Tags: []string{"b2", "media"}},
		{ID: "b2-med-statement", DeckID: "b2-media", Front: "das Statement", Back: "statement", Tags: []string{"b2", "media"}},
		{ID: "b2-med-veroeffentlichen", DeckID: "b2-media", Front: "veröffentlichen", Back: "to publish", Tags: []string{"b2", "media"}},
		{ID: "b2-med-zensur", DeckID: "b2-media", Front: "die Zensur", Back: "censorship", Tags: []string{"b2", "media"}},
		{ID: "b2-med-meinungsfreiheit", DeckID: "b2-media", Front: "die Meinungsfreiheit", Back: "freedom of speech", Tags: []string{"b2", "media"}},
		{ID: "b2-med-medien", DeckID: "b2-media", Front: "die Medien", Back: "media (plural)", Tags: []string{"b2", "media"}},
		{ID: "b2-med-fernseher", DeckID: "b2-media", Front: "der Fernseher", Back: "TV set", Tags: []string{"b2", "media"}},
		{ID: "b2-med-rundfunk", DeckID: "b2-media", Front: "der Rundfunk", Back: "broadcasting", Tags: []string{"b2", "media"}},
		{ID: "b2-med-print", DeckID: "b2-media", Front: "der Printmedien", Back: "print media", Tags: []string{"b2", "media"}},
		{ID: "b2-med-publikum", DeckID: "b2-media", Front: "das Publikum", Back: "audience", Tags: []string{"b2", "media"}},
		{ID: "b2-med-reichweite", DeckID: "b2-media", Front: "die Reichweite", Back: "reach (audience)", Tags: []string{"b2", "media"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-media",
		Name:        "German B2 Media & News",
		Description: "Vocabulary for discussing media, journalism, and news.",
		Tags:        []string{"german", "b2", "media", "news"},
		Notes:       notes,
	}
}
