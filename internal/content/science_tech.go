package content

import (
	"deutsch-tui/internal/core"
)

func ScienceTechDeck() core.Deck {
	notes := []core.Note{
		// --- Computers & Internet ---
		{ID: "tech-computer", DeckID: "science-tech", Front: "der Computer", Back: "computer", Tags: []string{"tech", "noun"}},
		{ID: "tech-internet", DeckID: "science-tech", Front: "das Internet", Back: "internet", Tags: []string{"tech", "noun"}},
		{ID: "tech-webseite", DeckID: "science-tech", Front: "die Webseite", Back: "website", Tags: []string{"tech", "noun"}},
		{ID: "tech-datei", DeckID: "science-tech", Front: "die Datei", Back: "file", Tags: []string{"tech", "noun"}},
		{ID: "tech-speichern", DeckID: "science-tech", Front: "speichern", Back: "to save", Tags: []string{"tech", "verb"}},
		{ID: "tech-loschen", DeckID: "science-tech", Front: "löschen", Back: "to delete", Tags: []string{"tech", "verb"}},
		{ID: "tech-herunterladen", DeckID: "science-tech", Front: "herunterladen", Back: "to download", Tags: []string{"tech", "verb"}},
		{ID: "tech-hochladen", DeckID: "science-tech", Front: "hochladen", Back: "to upload", Tags: []string{"tech", "verb"}},
		{ID: "tech-passwort", DeckID: "science-tech", Front: "das Passwort", Back: "passwort", Tags: []string{"tech", "noun"}},
		{ID: "tech-benutzername", DeckID: "science-tech", Front: "der Benutzername", Back: "username", Tags: []string{"tech", "noun"}},
		{ID: "tech-kunstliche-intelligenz", DeckID: "science-tech", Front: "die Künstliche Intelligenz", Back: "artificial intelligence (AI)", Tags: []string{"tech", "noun", "ai"}},
		{ID: "tech-algorithmus", DeckID: "science-tech", Front: "der Algorithmus", Back: "algorithm", Tags: []string{"tech", "noun"}},
		{ID: "tech-programmieren", DeckID: "science-tech", Front: "programmieren", Back: "to program", Tags: []string{"tech", "verb"}},
		{ID: "tech-software", DeckID: "science-tech", Front: "die Software", Back: "software", Tags: []string{"tech", "noun"}},
		{ID: "tech-hardware", DeckID: "science-tech", Front: "die Hardware", Back: "hardware", Tags: []string{"tech", "noun"}},

		// --- Science & Research ---
		{ID: "sci-wissenschaft", DeckID: "science-tech", Front: "die Wissenschaft", Back: "science", Tags: []string{"science", "noun"}},
		{ID: "sci-forschung", DeckID: "science-tech", Front: "die Forschung", Back: "research", Tags: []string{"science", "noun"}},
		{ID: "sci-experiment", DeckID: "science-tech", Front: "das Experiment", Back: "experiment", Tags: []string{"science", "noun"}},
		{ID: "sci-labor", DeckID: "science-tech", Front: "das Labor", Back: "laboratory", Tags: []string{"science", "noun"}},
		{ID: "sci-hypothese", DeckID: "science-tech", Front: "die Hypothese", Back: "hypothesis", Tags: []string{"science", "noun"}},
		{ID: "sci-theorie", DeckID: "science-tech", Front: "die Theorie", Back: "theory", Tags: []string{"science", "noun"}},
		{ID: "sci-beweis", DeckID: "science-tech", Front: "der Beweis", Back: "proof/evidence", Tags: []string{"science", "noun"}},
		{ID: "sci-entdeckung", DeckID: "science-tech", Front: "die Entdeckung", Back: "discovery", Tags: []string{"science", "noun"}},
		{ID: "sci-erfindung", DeckID: "science-tech", Front: "die Erfindung", Back: "invention", Tags: []string{"science", "noun"}},
		{ID: "sci-chemie", DeckID: "science-tech", Front: "die Chemie", Back: "chemistry", Tags: []string{"science", "noun"}},
		{ID: "sci-physik", DeckID: "science-tech", Front: "die Physik", Back: "physics", Tags: []string{"science", "noun"}},
		{ID: "sci-biologie", DeckID: "science-tech", Front: "die Biologie", Back: "biology", Tags: []string{"science", "noun"}},
		{ID: "sci-astronomie", DeckID: "science-tech", Front: "die Astronomie", Back: "astronomy", Tags: []string{"science", "noun"}},
		{ID: "sci-mathematik", DeckID: "science-tech", Front: "die Mathematik", Back: "mathematics", Tags: []string{"science", "noun"}},

		// --- Space ---
		{ID: "space-weltraum", DeckID: "science-tech", Front: "der Weltraum", Back: "outer space", Tags: []string{"space", "noun"}},
		{ID: "space-planet", DeckID: "science-tech", Front: "der Planet", Back: "planet", Tags: []string{"space", "noun"}},
		{ID: "space-stern", DeckID: "science-tech", Front: "der Stern", Back: "star", Tags: []string{"space", "noun"}},
		{ID: "space-galaxie", DeckID: "science-tech", Front: "die Galaxie", Back: "galaxy", Tags: []string{"space", "noun"}},
		{ID: "space-rakete", DeckID: "science-tech", Front: "die Rakete", Back: "rocket", Tags: []string{"space", "noun"}},
		{ID: "space-astronaut", DeckID: "science-tech", Front: "der Astronaut", Back: "astronaut", Tags: []string{"space", "noun"}},
		{ID: "space-universum", DeckID: "science-tech", Front: "das Universum", Back: "universum", Tags: []string{"space", "noun"}},
		{ID: "space-schwerkraft", DeckID: "science-tech", Front: "die Schwerkraft", Back: "gravity", Tags: []string{"space", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "science-tech",
		Name:        "Science & Technology",
		Description: "Advanced German vocabulary for science, computers, and space.",
		Tags:        []string{"german", "b1", "b2", "tech", "science"},
		Notes:       notes,
	}
}
