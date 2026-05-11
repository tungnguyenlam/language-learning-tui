package content

import (
	"deutsch-tui/internal/core"
)

func B1FinanceDeck() core.Deck {
	notes := []core.Note{
		{ID: "fin-001", DeckID: "b1-finance", Front: "das Konto", Back: "account", Extra: "Plural: die Konten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-002", DeckID: "b1-finance", Front: "das Girokonto", Back: "checking account", Extra: "Plural: die Girokonten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-003", DeckID: "b1-finance", Front: "das Sparkonto", Back: "savings account", Extra: "Plural: die Sparkonten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-004", DeckID: "b1-finance", Front: "die Bank", Back: "bank", Extra: "Plural: die Banken", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-005", DeckID: "b1-finance", Front: "der Kontostand", Back: "account balance", Extra: "Plural: die Kontostände", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-006", DeckID: "b1-finance", Front: "die Überweisung", Back: "bank transfer", Extra: "Plural: die Überweisungen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-007", DeckID: "b1-finance", Front: "die Einzahlung", Back: "deposit", Extra: "Plural: die Einzahlungen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-008", DeckID: "b1-finance", Front: "die Auszahlung", Back: "withdrawal", Extra: "Plural: die Auszahlungen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-009", DeckID: "b1-finance", Front: "der Geldautomat", Back: "ATM, cash machine", Extra: "Plural: die Geldautomaten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-010", DeckID: "b1-finance", Front: "die Kreditkarte", Back: "credit card", Extra: "Plural: die Kreditkarten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-011", DeckID: "b1-finance", Front: "die Debitkarte", Back: "debit card", Extra: "Plural: die Debitkarten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-012", DeckID: "b1-finance", Front: "die PIN", Back: "PIN", Extra: "Plural: die PINs", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-013", DeckID: "b1-finance", Front: "das Bargeld", Back: "cash", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-014", DeckID: "b1-finance", Front: "der Wechselkurs", Back: "exchange rate", Extra: "Plural: die Wechselkurse", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-015", DeckID: "b1-finance", Front: "die Währung", Back: "currency", Extra: "Plural: die Währungen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-016", DeckID: "b1-finance", Front: "der Euro", Back: "euro", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-017", DeckID: "b1-finance", Front: "der Cent", Back: "cent", Extra: "Plural: die Cent", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-018", DeckID: "b1-finance", Front: "der Geldschein", Back: "banknote", Extra: "Plural: die Geldscheine", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-019", DeckID: "b1-finance", Front: "die Münze", Back: "coin", Extra: "Plural: die Münzen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-020", DeckID: "b1-finance", Front: "der Kontoauszug", Back: "bank statement", Extra: "Plural: die Kontoauszüge", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-021", DeckID: "b1-finance", Front: "die Rechnung", Back: "invoice, bill", Extra: "Plural: die Rechnungen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-022", DeckID: "b1-finance", Front: "bezahlen", Back: "to pay", Extra: "Konjugation: ich bezahle", Tags: []string{"b1", "finance", "verb"}},
		{ID: "fin-023", DeckID: "b1-finance", Front: "der Preis", Back: "price", Extra: "Plural: die Preise", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-024", DeckID: "b1-finance", Front: "die Kosten", Back: "costs, expenses", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-025", DeckID: "b1-finance", Front: "das Budget", Back: "budget", Extra: "Plural: die Budgets", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-026", DeckID: "b1-finance", Front: "sparen", Back: "to save (money)", Extra: "Konjugation: ich spare", Tags: []string{"b1", "finance", "verb"}},
		{ID: "fin-027", DeckID: "b1-finance", Front: "ausgeben", Back: "to spend", Extra: "Konjugation: ich gebe aus", Tags: []string{"b1", "finance", "verb"}},
		{ID: "fin-028", DeckID: "b1-finance", Front: "einnehmen", Back: "to earn, to receive income", Extra: "Konjugation: ich nehme ein", Tags: []string{"b1", "finance", "verb"}},
		{ID: "fin-029", DeckID: "b1-finance", Front: "der Lohn", Back: "wage, salary", Extra: "Plural: die Löhne", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-030", DeckID: "b1-finance", Front: "das Gehalt", Back: "salary", Extra: "Plural: die Gehälter", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-031", DeckID: "b1-finance", Front: "die Miete", Back: "rent", Extra: "Plural: die Mieten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-032", DeckID: "b1-finance", Front: "die Nebenkosten", Back: "utilities (additional costs)", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-033", DeckID: "b1-finance", Front: "der Vertrag", Back: "contract", Extra: "Plural: die Verträge", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-034", DeckID: "b1-finance", Front: "kündigen", Back: "to cancel, to terminate", Extra: "Konjugation: ich kündige", Tags: []string{"b1", "finance", "verb"}},
		{ID: "fin-035", DeckID: "b1-finance", Front: "die Rate", Back: "installment, rate", Extra: "Plural: die Raten", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-036", DeckID: "b1-finance", Front: "der Kredit", Back: "loan, credit", Extra: "Plural: die Kredite", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-037", DeckID: "b1-finance", Front: "die Schulden", Back: "debts", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-038", DeckID: "b1-finance", Front: "der Zins", Back: "interest", Extra: "Plural: die Zinsen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-039", DeckID: "b1-finance", Front: "die Investition", Back: "investment", Extra: "Plural: die Investitionen", Tags: []string{"b1", "finance", "noun"}},
		{ID: "fin-040", DeckID: "b1-finance", Front: "die Steuer", Back: "tax", Extra: "Plural: die Steuern", Tags: []string{"b1", "finance", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-finance",
		Name:        "B1 Finance & Banking",
		Description: "Financial vocabulary, banking, and money management",
		Tags:        []string{"b1", "finance", "banking", "money"},
		Notes:       notes,
	}
}
