package content

import (
	"deutsch-tui/internal/core"
)

func B2ProgrammingDeck() core.Deck {
	notes := []core.Note{
		{ID: "prog-001", DeckID: "b2-programming", Front: "die Schnittstelle", Back: "interface, API", Extra: "Plural: die Schnittstellen", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-002", DeckID: "b2-programming", Front: "der Quellcode", Back: "source code", Extra: "Auch: der Programmcode", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-003", DeckID: "b2-programming", Front: "bereitstellen", Back: "to deploy, to provide", Extra: "Konjugation: ich stelle bereit. Nomen: die Bereitstellung", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-004", DeckID: "b2-programming", Front: "die Fehlermeldung", Back: "error message", Extra: "Plural: die Fehlermeldungen", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-005", DeckID: "b2-programming", Front: "das Repository", Back: "repository", Extra: "Plural: die Repositories", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-006", DeckID: "b2-programming", Front: "committen", Back: "to commit", Extra: "Konjugation: ich committe", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-007", DeckID: "b2-programming", Front: "mergen", Back: "to merge", Extra: "Konjugation: ich merge", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-008", DeckID: "b2-programming", Front: "der Branch", Back: "branch", Extra: "Plural: die Branches", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-009", DeckID: "b2-programming", Front: "die Versionsverwaltung", Back: "version control", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-010", DeckID: "b2-programming", Front: "die Entwicklungsumgebung", Back: "development environment", Extra: "IDE = Integrierte Entwicklungsumgebung", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-011", DeckID: "b2-programming", Front: "das Framework", Back: "framework", Extra: "Plural: die Frameworks", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-012", DeckID: "b2-programming", Front: "die Bibliothek", Back: "library", Extra: "Plural: die Bibliotheken", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-013", DeckID: "b2-programming", Front: "die Abhängigkeit", Back: "dependency", Extra: "Plural: die Abhängigkeiten", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-014", DeckID: "b2-programming", Front: "kompilieren", Back: "to compile", Extra: "Konjugation: ich kompiliere", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-015", DeckID: "b2-programming", Front: "ausführen", Back: "to execute, to run", Extra: "Konjugation: ich führe aus", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-016", DeckID: "b2-programming", Front: "der Container", Back: "container (e.g. Docker)", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-017", DeckID: "b2-programming", Front: "das Backend", Back: "backend", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-018", DeckID: "b2-programming", Front: "das Frontend", Back: "frontend", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-019", DeckID: "b2-programming", Front: "die Benutzeroberfläche", Back: "user interface", Extra: "UI = User Interface", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-020", DeckID: "b2-programming", Front: "die Abfrage", Back: "query", Extra: "Plural: die Abfragen", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-021", DeckID: "b2-programming", Front: "verschlüsselt", Back: "encrypted", Extra: "Verb: verschlüsseln", Tags: []string{"b2", "programming", "adjective"}},
		{ID: "prog-022", DeckID: "b2-programming", Front: "die Authentifizierung", Back: "authentication", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-023", DeckID: "b2-programming", Front: "die Autorisierung", Back: "authorization", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-024", DeckID: "b2-programming", Front: "skalierbar", Back: "scalable", Extra: "Nomen: die Skalierbarkeit", Tags: []string{"b2", "programming", "adjective"}},
		{ID: "prog-025", DeckID: "b2-programming", Front: "die Wartung", Back: "maintenance", Extra: "Verb: warten", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-026", DeckID: "b2-programming", Front: "der Bugfix", Back: "bug fix", Extra: "Plural: die Bugfixes", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-027", DeckID: "b2-programming", Front: "testen", Back: "to test", Extra: "Konjugation: ich teste", Tags: []string{"b2", "programming", "verb"}},
		{ID: "prog-028", DeckID: "b2-programming", Front: "die Testabdeckung", Back: "test coverage", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-029", DeckID: "b2-programming", Front: "die Dokumentation", Back: "documentation", Tags: []string{"b2", "programming", "noun"}},
		{ID: "prog-030", DeckID: "b2-programming", Front: "das Skript", Back: "script", Extra: "Plural: die Skripte", Tags: []string{"b2", "programming", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-programming",
		Name:        "B2 Programming & Software Engineering",
		Description: "Technical vocabulary for software developers and engineers",
		Tags:        []string{"b2", "programming", "tech", "software"},
		Notes:       notes,
	}
}
