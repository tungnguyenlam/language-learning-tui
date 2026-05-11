package content

import (
	"deutsch-tui/internal/core"
)

func B1TechnologyDeck() core.Deck {
	notes := []core.Note{
		{ID: "tech-001", DeckID: "b1-technology", Front: "der Computer", Back: "computer", Extra: "Plural: die Computer", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-002", DeckID: "b1-technology", Front: "das Laptop", Back: "laptop, notebook", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-003", DeckID: "b1-technology", Front: "das Smartphone", Back: "smartphone", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-004", DeckID: "b1-technology", Front: "das Tablet", Back: "tablet", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-005", DeckID: "b1-technology", Front: "der Bildschirm", Back: "screen, monitor", Extra: "Plural: die Bildschirme", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-006", DeckID: "b1-technology", Front: "die Tastatur", Back: "keyboard", Extra: "Plural: die Tastaturen", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-007", DeckID: "b1-technology", Front: "die Maus", Back: "mouse", Extra: "Plural: die Mäuse", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-008", DeckID: "b1-technology", Front: "der Drucker", Back: "printer", Extra: "Plural: die Drucker", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-009", DeckID: "b1-technology", Front: "der Scanner", Back: "scanner", Extra: "Plural: die Scanner", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-010", DeckID: "b1-technology", Front: "die Webseite", Back: "website", Extra: "Plural: die Webseiten", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-011", DeckID: "b1-technology", Front: "die E-Mail", Back: "email", Extra: "Plural: die E-Mails", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-012", DeckID: "b1-technology", Front: "der Link", Back: "link, hyperlink", Extra: "Plural: die Links", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-013", DeckID: "b1-technology", Front: "das Passwort", Back: "password", Extra: "Plural: die Passwörter", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-014", DeckID: "b1-technology", Front: "der Benutzername", Back: "username", Extra: "Plural: die Benutzernamen", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-015", DeckID: "b1-technology", Front: "der Account", Back: "account", Extra: "Plural: die Accounts", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-016", DeckID: "b1-technology", Front: "hochladen", Back: "to upload", Extra: "Konjugation: ich lade hoch", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-017", DeckID: "b1-technology", Front: "herunterladen", Back: "to download", Extra: "Konjugation: ich lade herunter", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-018", DeckID: "b1-technology", Front: "versenden", Back: "to send", Extra: "Konjugation: ich versende", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-019", DeckID: "b1-technology", Front: "anklicken", Back: "to click", Extra: "Konjugation: ich klicke an", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-020", DeckID: "b1-technology", Front: "ausfüllen", Back: "to fill out", Extra: "Konjugation: ich fülle aus", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-021", DeckID: "b1-technology", Front: "der Browser", Back: "browser", Extra: "Plural: die Browser", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-022", DeckID: "b1-technology", Front: "die Suchmaschine", Back: "search engine", Extra: "Plural: die Suchmaschinen", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-023", DeckID: "b1-technology", Front: "der Cloud-Speicher", Back: "cloud storage", Extra: "Plural: die Cloud-Speicher", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-024", DeckID: "b1-technology", Front: "das WLAN", Back: "WiFi", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-025", DeckID: "b1-technology", Front: "die Internetverbindung", Back: "internet connection", Extra: "Plural: die Internetverbindungen", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-026", DeckID: "b1-technology", Front: "die Datei", Back: "file", Extra: "Plural: die Dateien", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-027", DeckID: "b1-technology", Front: "der Ordner", Back: "folder", Extra: "Plural: die Ordner", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-028", DeckID: "b1-technology", Front: "das Dokument", Back: "document", Extra: "Plural: die Dokumente", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-029", DeckID: "b1-technology", Front: "die Videokonferenz", Back: "video conference", Extra: "Plural: die Videokonferenzen", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-030", DeckID: "b1-technology", Front: "die Software", Back: "software", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-031", DeckID: "b1-technology", Front: "die Hardware", Back: "hardware", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-032", DeckID: "b1-technology", Front: "das Update", Back: "update", Extra: "Plural: die Updates", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-033", DeckID: "b1-technology", Front: "das Backup", Back: "backup", Extra: "Plural: die Backups", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-034", DeckID: "b1-technology", Front: "verschlüsseln", Back: "to encrypt", Extra: "Konjugation: ich verschlüssle", Tags: []string{"b1", "technology", "verb"}},
		{ID: "tech-035", DeckID: "b1-technology", Front: "der Virus", Back: "virus", Extra: "Plural: die Viren", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-036", DeckID: "b1-technology", Front: "der Hacker", Back: "hacker", Extra: "Plural: die Hacker", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-037", DeckID: "b1-technology", Front: "die Firewall", Back: "firewall", Extra: "Plural: die Firewalls", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-038", DeckID: "b1-technology", Front: "sicher", Back: "secure, safe", Extra: "Related: die Sicherheit", Tags: []string{"b1", "technology", "adjective"}},
		{ID: "tech-039", DeckID: "b1-technology", Front: "der Chip", Back: "chip", Extra: "Plural: die Chips", Tags: []string{"b1", "technology", "noun"}},
		{ID: "tech-040", DeckID: "b1-technology", Front: "die Künstliche Intelligenz", Back: "artificial intelligence", Tags: []string{"b1", "technology", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-technology",
		Name:        "B1 Technology & Digital Life",
		Description: "Technology vocabulary, digital communication, and internet terms",
		Tags:        []string{"b1", "technology", "digital", "internet"},
		Notes:       notes,
	}
}
