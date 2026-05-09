package content

import "deutsch-tui/internal/core"

func EnvironmentDeck() core.Deck {
	return core.Deck{
		ID:          "b2-environment-climate",
		Name:        "Environment & Climate (B2)",
		Description: "Vocabulary related to ecology, climate change, and sustainability.",
		Notes: []core.Note{
			{
				ID:    "env-1",
				Front: "der Umweltschutz",
				Back:  "environmental protection",
				Extra: "Umweltschutz ist eine globale Aufgabe.",
				Tags:  []string{"environment", "B2"},
			},
			{
				ID:    "env-2",
				Front: "der Klimawandel",
				Back:  "climate change",
				Extra: "Der Klimawandel betrifft uns alle.",
				Tags:  []string{"climate", "B2"},
			},
			{
				ID:    "env-3",
				Front: "die Erneuerbare Energien",
				Back:  "renewable energies",
				Extra: "Wir müssen mehr in erneuerbare Energien investieren.",
				Tags:  []string{"energy", "environment", "B2"},
			},
			{
				ID:    "env-4",
				Front: "die Nachhaltigkeit",
				Back:  "sustainability",
				Extra: "Nachhaltigkeit spielt in der modernen Wirtschaft eine große Rolle.",
				Tags:  []string{"sustainability", "B2"},
			},
			{
				ID:    "env-5",
				Front: "der Treibhauseffekt",
				Back:  "greenhouse effect",
				Extra: "CO2 trägt maßgeblich zum Treibhauseffekt bei.",
				Tags:  []string{"climate", "science", "B2"},
			},
			{
				ID:    "env-6",
				Front: "die Artenvielfalt",
				Back:  "biodiversity",
				Extra: "Der Schutz der Artenvielfalt ist extrem wichtig.",
				Tags:  []string{"nature", "B2"},
			},
			{
				ID:    "env-7",
				Front: "die Mülltrennung",
				Back:  "waste separation",
				Extra: "Mülltrennung ist in Deutschland sehr verbreitet.",
				Tags:  []string{"environment", "daily-life", "B2"},
			},
			{
				ID:    "env-8",
				Front: "die Schadstoffemission",
				Back:  "pollutant emission",
				Extra: "Die Schadstoffemissionen müssen drastisch reduziert werden.",
				Tags:  []string{"environment", "science", "C1"},
			},
			{
				ID:    "env-9",
				Front: "das Ökosystem",
				Back:  "ecosystem",
				Extra: "Korallenriffe sind sehr empfindliche Ökosysteme.",
				Tags:  []string{"nature", "B2"},
			},
			{
				ID:    "env-10",
				Front: "die Ressourcenschonung",
				Back:  "resource conservation",
				Extra: "Ressourcenschonung schont auch den Geldbeutel.",
				Tags:  []string{"sustainability", "B2"},
			},
			{
				ID:    "env-11",
				Front: "umweltfreundlich",
				Back:  "environmentally friendly",
				Extra: "Fahrradfahren ist eine umweltfreundliche Art der Fortbewegung.",
				Tags:  []string{"environment", "B2"},
			},
			{
				ID:    "env-12",
				Front: "der Atommüll",
				Back:  "nuclear waste",
				Extra: "Die Entsorgung von Atommüll ist ein ungelöstes Problem.",
				Tags:  []string{"environment", "politics", "B2"},
			},
			{
				ID:    "env-13",
				Front: "die Wiederverwertung",
				Back:  "recycling",
				Extra: "Papier eignet sich hervorragend zur Wiederverwertung.",
				Tags:  []string{"environment", "B2"},
			},
			{
				ID:    "env-14",
				Front: "der Naturschutz",
				Back:  "nature conservation",
				Extra: "Dieses Gebiet steht unter Naturschutz.",
				Tags:  []string{"nature", "B2"},
			},
			{
				ID:    "env-15",
				Front: "die globale Erwärmung",
				Back:  "global warming",
				Extra: "Die globale Erwärmung führt zum Schmelzen der Pole.",
				Tags:  []string{"climate", "B2"},
			},
		},
	}
}
