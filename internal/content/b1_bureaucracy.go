package content

import "deutsch-tui/internal/core"

func B1BureaucracyAppointmentsDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-bureaucracy-termin", DeckID: "b1-bureaucracy-appointments", Front: "der Termin", Back: "appointment", Extra: "Ich habe morgen einen Termin beim Amt.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-amt", DeckID: "b1-bureaucracy-appointments", Front: "das Amt", Back: "public office", Extra: "Das Amt öffnet um acht Uhr.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-behoerde", DeckID: "b1-bureaucracy-appointments", Front: "die Behörde", Back: "authority / agency", Extra: "Die Behörde schickt einen Brief.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-formular", DeckID: "b1-bureaucracy-appointments", Front: "das Formular", Back: "form", Extra: "Bitte füllen Sie das Formular aus.", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-antrag", DeckID: "b1-bureaucracy-appointments", Front: "der Antrag", Back: "application / request", Extra: "Der Antrag muss unterschrieben werden.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-bescheinigung", DeckID: "b1-bureaucracy-appointments", Front: "die Bescheinigung", Back: "certificate / confirmation", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-meldebescheinigung", DeckID: "b1-bureaucracy-appointments", Front: "die Meldebescheinigung", Back: "registration certificate", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-ausweis", DeckID: "b1-bureaucracy-appointments", Front: "der Ausweis", Back: "ID card", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-reisepass", DeckID: "b1-bureaucracy-appointments", Front: "der Reisepass", Back: "passport", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-unterschrift", DeckID: "b1-bureaucracy-appointments", Front: "die Unterschrift", Back: "signature", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-unterlage", DeckID: "b1-bureaucracy-appointments", Front: "die Unterlage", Back: "document / paperwork", Extra: "Bringen Sie alle Unterlagen mit.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-frist", DeckID: "b1-bureaucracy-appointments", Front: "die Frist", Back: "deadline", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-gebuehr", DeckID: "b1-bureaucracy-appointments", Front: "die Gebühr", Back: "fee", Extra: "Die Gebühr beträgt zehn Euro.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-wartenummer", DeckID: "b1-bureaucracy-appointments", Front: "die Wartenummer", Back: "queue number", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-wartezimmer", DeckID: "b1-bureaucracy-appointments", Front: "das Wartezimmer", Back: "waiting room", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-sachbearbeiter", DeckID: "b1-bureaucracy-appointments", Front: "der Sachbearbeiter / die Sachbearbeiterin", Back: "case worker", Tags: []string{"b1", "bureaucracy", "person"}},
		{ID: "b1-bureaucracy-versicherung", DeckID: "b1-bureaucracy-appointments", Front: "die Versicherung", Back: "insurance", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-krankenversicherung", DeckID: "b1-bureaucracy-appointments", Front: "die Krankenversicherung", Back: "health insurance", Tags: []string{"b1", "health", "bureaucracy"}},
		{ID: "b1-bureaucracy-steuernummer", DeckID: "b1-bureaucracy-appointments", Front: "die Steuernummer", Back: "tax number", Tags: []string{"b1", "bureaucracy", "finance"}},
		{ID: "b1-bureaucracy-kontoauszug", DeckID: "b1-bureaucracy-appointments", Front: "der Kontoauszug", Back: "bank statement", Tags: []string{"b1", "finance", "bureaucracy"}},
		{ID: "b1-bureaucracy-nachweis", DeckID: "b1-bureaucracy-appointments", Front: "der Nachweis", Back: "proof / evidence", Extra: "Wir brauchen einen Nachweis über Ihr Einkommen.", Tags: []string{"b1", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-kopie", DeckID: "b1-bureaucracy-appointments", Front: "die Kopie", Back: "copy", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-original", DeckID: "b1-bureaucracy-appointments", Front: "das Original", Back: "original document", Tags: []string{"a2", "bureaucracy", "noun"}},
		{ID: "b1-bureaucracy-einreichen", DeckID: "b1-bureaucracy-appointments", Front: "einreichen", Back: "to submit", Extra: "Sie können den Antrag online einreichen.", Tags: []string{"b1", "bureaucracy", "verb"}},
		{ID: "b1-bureaucracy-beantragen", DeckID: "b1-bureaucracy-appointments", Front: "beantragen", Back: "to apply for", Tags: []string{"b1", "bureaucracy", "verb"}},
		{ID: "b1-bureaucracy-ausfuellen", DeckID: "b1-bureaucracy-appointments", Front: "ausfüllen", Back: "to fill out", Tags: []string{"a2", "bureaucracy", "verb"}},
		{ID: "b1-bureaucracy-unterschreiben", DeckID: "b1-bureaucracy-appointments", Front: "unterschreiben", Back: "to sign", Tags: []string{"a2", "bureaucracy", "verb"}},
		{ID: "b1-bureaucracy-bestaetigen", DeckID: "b1-bureaucracy-appointments", Front: "bestätigen", Back: "to confirm", Tags: []string{"b1", "bureaucracy", "verb"}},
		{ID: "b1-bureaucracy-verschieben", DeckID: "b1-bureaucracy-appointments", Front: "einen Termin verschieben", Back: "to reschedule an appointment", Tags: []string{"b1", "bureaucracy", "phrase"}},
		{ID: "b1-bureaucracy-vereinbaren", DeckID: "b1-bureaucracy-appointments", Front: "einen Termin vereinbaren", Back: "to make an appointment", Tags: []string{"b1", "bureaucracy", "phrase"}},
		{ID: "b1-bureaucracy-fehlt", DeckID: "b1-bureaucracy-appointments", Front: "Mir fehlt noch ein Dokument.", Back: "I am still missing one document.", Tags: []string{"b1", "bureaucracy", "phrase"}},
		{ID: "b1-bureaucracy-zustaendig", DeckID: "b1-bureaucracy-appointments", Front: "Wer ist dafür zuständig?", Back: "Who is responsible for that?", Tags: []string{"b1", "bureaucracy", "phrase"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "b1-bureaucracy-appointments",
		Name:        "German B1 Bureaucracy & Appointments",
		Description: "Office, forms, appointments, and document vocabulary for everyday administration in German.",
		Tags:        []string{"german", "b1", "bureaucracy", "appointments"},
		Notes:       notes,
	}
}
