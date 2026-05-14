package content

import "deutsch-tui/internal/core"

func B2DigitalPrivacyDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-privacy-datenschutz", DeckID: "b2-digital-privacy", Front: "der Datenschutz", Back: "data protection", Extra: "Datenschutz ist in Deutschland ein wichtiges Thema.", Tags: []string{"b2", "digital", "privacy", "noun"}},
		{ID: "b2-privacy-privatsphaere", DeckID: "b2-digital-privacy", Front: "die Privatsphäre", Back: "privacy", Tags: []string{"b2", "digital", "privacy", "noun"}},
		{ID: "b2-privacy-einwilligung", DeckID: "b2-digital-privacy", Front: "die Einwilligung", Back: "consent", Extra: "Ohne Einwilligung dürfen die Daten nicht genutzt werden.", Tags: []string{"b2", "digital", "law", "noun"}},
		{ID: "b2-privacy-zustimmung", DeckID: "b2-digital-privacy", Front: "die Zustimmung", Back: "approval / consent", Tags: []string{"b2", "digital", "noun"}},
		{ID: "b2-privacy-konto", DeckID: "b2-digital-privacy", Front: "das Benutzerkonto", Back: "user account", Tags: []string{"b1", "digital", "noun"}},
		{ID: "b2-privacy-passwort", DeckID: "b2-digital-privacy", Front: "das Passwort", Back: "password", Tags: []string{"a2", "digital", "noun"}},
		{ID: "b2-privacy-zwei-faktor", DeckID: "b2-digital-privacy", Front: "die Zwei-Faktor-Authentifizierung", Back: "two-factor authentication", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-verschluesselung", DeckID: "b2-digital-privacy", Front: "die Verschlüsselung", Back: "encryption", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-datenleck", DeckID: "b2-digital-privacy", Front: "das Datenleck", Back: "data leak", Extra: "Das Unternehmen meldete ein Datenleck.", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-sicherheitsluecke", DeckID: "b2-digital-privacy", Front: "die Sicherheitslücke", Back: "security vulnerability", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-betrug", DeckID: "b2-digital-privacy", Front: "der Betrug", Back: "fraud", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-phishing", DeckID: "b2-digital-privacy", Front: "das Phishing", Back: "phishing", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-spam", DeckID: "b2-digital-privacy", Front: "der Spam", Back: "spam", Tags: []string{"b1", "digital", "security"}},
		{ID: "b2-privacy-malware", DeckID: "b2-digital-privacy", Front: "die Schadsoftware", Back: "malware", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-firewall", DeckID: "b2-digital-privacy", Front: "die Firewall", Back: "firewall", Tags: []string{"b2", "digital", "security"}},
		{ID: "b2-privacy-geraet", DeckID: "b2-digital-privacy", Front: "das Gerät", Back: "device", Tags: []string{"b1", "digital", "noun"}},
		{ID: "b2-privacy-netzwerk", DeckID: "b2-digital-privacy", Front: "das Netzwerk", Back: "network", Tags: []string{"b1", "digital", "noun"}},
		{ID: "b2-privacy-verlauf", DeckID: "b2-digital-privacy", Front: "der Browserverlauf", Back: "browser history", Tags: []string{"b2", "digital", "noun"}},
		{ID: "b2-privacy-cookie", DeckID: "b2-digital-privacy", Front: "das Cookie", Back: "cookie", Tags: []string{"b1", "digital", "noun"}},
		{ID: "b2-privacy-standort", DeckID: "b2-digital-privacy", Front: "der Standort", Back: "location", Tags: []string{"b1", "digital", "noun"}},
		{ID: "b2-privacy-berechtigung", DeckID: "b2-digital-privacy", Front: "die Berechtigung", Back: "permission / authorization", Tags: []string{"b2", "digital", "noun"}},
		{ID: "b2-privacy-nutzungsbedingungen", DeckID: "b2-digital-privacy", Front: "die Nutzungsbedingungen", Back: "terms of service", Tags: []string{"b2", "digital", "law"}},
		{ID: "b2-privacy-richtlinie", DeckID: "b2-digital-privacy", Front: "die Datenschutzrichtlinie", Back: "privacy policy", Tags: []string{"b2", "digital", "law"}},
		{ID: "b2-privacy-weitergabe", DeckID: "b2-digital-privacy", Front: "die Weitergabe von Daten", Back: "sharing of data", Tags: []string{"b2", "digital", "phrase"}},
		{ID: "b2-privacy-speichern", DeckID: "b2-digital-privacy", Front: "Daten speichern", Back: "to store data", Tags: []string{"b1", "digital", "phrase"}},
		{ID: "b2-privacy-loeschen", DeckID: "b2-digital-privacy", Front: "Daten löschen", Back: "to delete data", Tags: []string{"b1", "digital", "phrase"}},
		{ID: "b2-privacy-sperren", DeckID: "b2-digital-privacy", Front: "ein Konto sperren", Back: "to block an account", Tags: []string{"b2", "digital", "phrase"}},
		{ID: "b2-privacy-melden", DeckID: "b2-digital-privacy", Front: "einen Vorfall melden", Back: "to report an incident", Tags: []string{"b2", "digital", "phrase"}},
		{ID: "b2-privacy-aktualisieren", DeckID: "b2-digital-privacy", Front: "die Software aktualisieren", Back: "to update the software", Tags: []string{"b1", "digital", "phrase"}},
		{ID: "b2-privacy-sichern", DeckID: "b2-digital-privacy", Front: "eine Sicherung erstellen", Back: "to create a backup", Tags: []string{"b2", "digital", "phrase"}},
		{ID: "b2-privacy-verdacht", DeckID: "b2-digital-privacy", Front: "Ich habe den Verdacht, dass mein Konto gehackt wurde.", Back: "I suspect that my account was hacked.", Tags: []string{"b2", "digital", "sentence"}},
		{ID: "b2-privacy-zugriff", DeckID: "b2-digital-privacy", Front: "Wer hat Zugriff auf diese Daten?", Back: "Who has access to this data?", Tags: []string{"b2", "digital", "sentence"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "b2-digital-privacy",
		Name:        "German B2 Digital Privacy & Security",
		Description: "Data protection, account security, and online risk vocabulary for B2 learners.",
		Tags:        []string{"german", "b2", "digital", "privacy", "security"},
		Notes:       notes,
	}
}
