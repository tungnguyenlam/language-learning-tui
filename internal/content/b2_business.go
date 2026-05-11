package content

import (
	"deutsch-tui/internal/core"
)

func B2BusinessDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-bus-besprechung", DeckID: "b2-business", Front: "die Besprechung", Back: "meeting", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-konferenz", DeckID: "b2-business", Front: "die Konferenz", Back: "conference", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-präsentation", DeckID: "b2-business", Front: "die Präsentation", Back: "presentation", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-verhandlung", DeckID: "b2-business", Front: "die Verhandlung", Back: "negotiation", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-vertrag", DeckID: "b2-business", Front: "der Vertrag", Back: "contract", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-angebot", DeckID: "b2-business", Front: "das Angebot", Back: "offer / quote", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-preisverhandlungen", DeckID: "b2-business", Front: "die Preisverhandlungen", Back: "price negotiations", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-termin", DeckID: "b2-business", Front: "der Termin", Back: "appointment", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-deadline", DeckID: "b2-business", Front: "die Deadline", Back: "deadline", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-projekt", DeckID: "b2-business", Front: "das Projekt", Back: "project", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-aufgabe", DeckID: "b2-business", Front: "die Aufgabe", Back: "task / assignment", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-verantwortung", DeckID: "b2-business", Front: "die Verantwortung", Back: "responsibility", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-kollege", DeckID: "b2-business", Front: "der Kollege / die Kollegin", Back: "colleague", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-vorgesetzte", DeckID: "b2-business", Front: "der Vorgesetzte / die Vorgesetzte", Back: "superior / boss", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-mitarbeiter", DeckID: "b2-business", Front: "der Mitarbeiter / die Mitarbeiterin", Back: "employee", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-bewerber", DeckID: "b2-business", Front: "der Bewerber / die Bewerberin", Back: "applicant", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-arbeitsplatz", DeckID: "b2-business", Front: "der Arbeitsplatz", Back: "workplace", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-büro", DeckID: "b2-business", Front: "das Büro", Back: "office", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-firmware", DeckID: "b2-business", Front: "die Firma", Back: "company / firm", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-unternehmen", DeckID: "b2-business", Front: "das Unternehmen", Back: "enterprise / company", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-abteilung", DeckID: "b2-business", Front: "die Abteilung", Back: "department", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-gehalt", DeckID: "b2-business", Front: "das Gehalt", Back: "salary", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-boni", DeckID: "b2-business", Front: "die Bonuszahlung", Back: "bonus", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-urlaub", DeckID: "b2-business", Front: "der Urlaub", Back: "vacation / holiday", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-kündigung", DeckID: "b2-business", Front: "die Kündigung", Back: "resignation / termination", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-bewerbung", DeckID: "b2-business", Front: "die Bewerbung", Back: "application (job)", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-lebenslauf", DeckID: "b2-business", Front: "der Lebenslauf", Back: "resume / CV", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-arbeitsvertrag", DeckID: "b2-business", Front: "der Arbeitsvertrag", Back: "employment contract", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-arbeitszeit", DeckID: "b2-business", Front: "die Arbeitszeit", Back: "working hours", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-flexzeit", DeckID: "b2-business", Front: "die Gleitzeit", Back: "flexible working hours", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-homeoffice", DeckID: "b2-business", Front: "das Homeoffice", Back: "home office", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-fortbildung", DeckID: "b2-business", Front: "die Fortbildung", Back: "professional development", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-training", DeckID: "b2-business", Front: "das Training", Back: "training", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-quualifikation", DeckID: "b2-business", Front: "die Qualifikation", Back: "qualification", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-kompetenz", DeckID: "b2-business", Front: "die Kompetenz", Back: "competence", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-erfahrung", DeckID: "b2-business", Front: "die Erfahrung", Back: "experience", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-strategie", DeckID: "b2-business", Front: "die Strategie", Back: "strategy", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-marketing", DeckID: "b2-business", Front: "das Marketing", Back: "marketing", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-umsatz", DeckID: "b2-business", Front: "der Umsatz", Back: "revenue / turnover", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-gewinn", DeckID: "b2-business", Front: "der Gewinn", Back: "profit", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-verlust", DeckID: "b2-business", Front: "der Verlust", Back: "loss", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-investition", DeckID: "b2-business", Front: "die Investition", Back: "investment", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-budget", DeckID: "b2-business", Front: "das Budget", Back: "budget", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-rechnung", DeckID: "b2-business", Front: "die Rechnung", Back: "invoice", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-zahlung", DeckID: "b2-business", Front: "die Zahlung", Back: "payment", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-fällig", DeckID: "b2-business", Front: "fällig", Back: "due / payable", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-zahlungsbedingungen", DeckID: "b2-business", Front: "die Zahlungsbedingungen", Back: "payment terms", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-konto", DeckID: "b2-business", Front: "das Konto", Back: "account", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-rechnungskopie", DeckID: "b2-business", Front: "die Rechnungskopie", Back: "copy of invoice", Tags: []string{"b2", "business"}},
		{ID: "b2-bus-kundendienst", DeckID: "b2-business", Front: "der Kundendienst", Back: "customer service", Tags: []string{"b2", "business"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-business",
		Name:        "German B2 Business & Workplace",
		Description: "Professional vocabulary for business and workplace contexts.",
		Tags:        []string{"german", "b2", "business", "workplace"},
		Notes:       notes,
	}
}
