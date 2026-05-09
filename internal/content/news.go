package content

import "deutsch-tui/internal/core"

func NewsDeck() core.Deck {
	return core.Deck{
		ID:          "b2-c1-news-events",
		Name:        "German News & Current Events (B2-C1)",
		Description: "Advanced vocabulary from German news, politics, and society.",
		Notes: []core.Note{
			{
				ID:    "news-1",
				Front: "die Berichterstattung",
				Back:  "reporting, coverage",
				Extra: "Die mediale Berichterstattung war sehr sachlich.",
				Tags:  []string{"news", "media", "B2"},
			},
			{
				ID:    "news-2",
				Front: "der Abgeordnete",
				Back:  "member of parliament, delegate",
				Extra: "Der Abgeordnete hielt eine leidenschaftliche Rede.",
				Tags:  []string{"politics", "C1"},
			},
			{
				ID:    "news-3",
				Front: "die Verhandlung",
				Back:  "negotiation, hearing",
				Extra: "Die Tarifverhandlungen dauerten bis spät in die Nacht.",
				Tags:  []string{"business", "politics", "B2"},
			},
			{
				ID:    "news-4",
				Front: "der Beschluss",
				Back:  "resolution, decision",
				Extra: "Die Regierung fasste einen wichtigen Beschluss zur Klimapolitik.",
				Tags:  []string{"politics", "B2"},
			},
			{
				ID:    "news-5",
				Front: "die Herausforderung",
				Back:  "challenge",
				Extra: "Die Digitalisierung ist eine große Herausforderung für die Gesellschaft.",
				Tags:  []string{"society", "B2"},
			},
			{
				ID:    "news-6",
				Front: "die Schlagzeile",
				Back:  "headline",
				Extra: "Diese Nachricht sorgte für große Schlagzeilen.",
				Tags:  []string{"news", "media", "B2"},
			},
			{
				ID:    "news-7",
				Front: "der Haushalt",
				Back:  "budget (also household)",
				Extra: "Das Parlament debattiert über den neuen Haushalt.",
				Tags:  []string{"politics", "economy", "B2"},
			},
			{
				ID:    "news-8",
				Front: "die Gesetzgebung",
				Back:  "legislation",
				Extra: "Die Gesetzgebung in diesem Bereich ist sehr komplex.",
				Tags:  []string{"politics", "law", "C1"},
			},
			{
				ID:    "news-9",
				Front: "nachhaltig",
				Back:  "sustainable",
				Extra: "Wir brauchen eine nachhaltige Lösung für dieses Problem.",
				Tags:  []string{"environment", "economy", "B2"},
			},
			{
				ID:    "news-10",
				Front: "die Wirtschaftskrise",
				Back:  "economic crisis",
				Extra: "Das Land hat sich von der letzten Wirtschaftskrise erholt.",
				Tags:  []string{"economy", "B2"},
			},
		},
	}
}
