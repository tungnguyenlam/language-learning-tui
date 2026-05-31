package content

import (
	"deutsch-tui/internal/core"
)

func C2FinanceDeck() core.Deck {
	notes := []core.Note{
		{ID: "c2-fin-aktie", DeckID: "c2-finance", Front: "die Aktie", Back: "share / stock", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-börse", DeckID: "c2-finance", Front: "die Börse", Back: "stock exchange", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-zins", DeckID: "c2-finance", Front: "der Zins", Back: "interest (finance)", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-inflation", DeckID: "c2-finance", Front: "die Inflation", Back: "inflation", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-rendite", DeckID: "c2-finance", Front: "die Rendite", Back: "yield / return on investment", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-währung", DeckID: "c2-finance", Front: "die Währung", Back: "currency", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-wechselkurs", DeckID: "c2-finance", Front: "der Wechselkurs", Back: "exchange rate", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-dividende", DeckID: "c2-finance", Front: "die Dividende", Back: "dividend", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-anleihe", DeckID: "c2-finance", Front: "die Anleihe", Back: "bond", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-portfolio", DeckID: "c2-finance", Front: "das Portfolio", Back: "portfolio", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-liquidität", DeckID: "c2-finance", Front: "die Liquidität", Back: "liquidity", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-insolvenz", DeckID: "c2-finance", Front: "die Insolvenz", Back: "insolvency / bankruptcy", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-subvention", DeckID: "c2-finance", Front: "die Subvention", Back: "subsidy", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-kredit", DeckID: "c2-finance", Front: "der Kredit", Back: "credit / loan", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-hypothek", DeckID: "c2-finance", Front: "die Hypothek", Back: "mortgage", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-steuern", DeckID: "c2-finance", Front: "die Steuern", Back: "taxes", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-bruttoinlandsprodukt", DeckID: "c2-finance", Front: "das Bruttoinlandsprodukt (BIP)", Back: "Gross Domestic Product (GDP)", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-konjunktur", DeckID: "c2-finance", Front: "die Konjunktur", Back: "economic cycle / economy", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-rezession", DeckID: "c2-finance", Front: "die Rezession", Back: "recession", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-defizit", DeckID: "c2-finance", Front: "das Defizit", Back: "deficit", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-monopol", DeckID: "c2-finance", Front: "das Monopol", Back: "monopoly", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-oligopol", DeckID: "c2-finance", Front: "das Oligopol", Back: "oligopoly", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-kapital", DeckID: "c2-finance", Front: "das Kapital", Back: "capital", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-gewinn", DeckID: "c2-finance", Front: "der Gewinn", Back: "profit", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-verlust", DeckID: "c2-finance", Front: "der Verlust", Back: "loss", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-umsatz", DeckID: "c2-finance", Front: "der Umsatz", Back: "revenue / turnover", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-investition", DeckID: "c2-finance", Front: "die Investition", Back: "investment", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-spekulation", DeckID: "c2-finance", Front: "die Spekulation", Back: "speculation", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-wachstum", DeckID: "c2-finance", Front: "das Wachstum", Back: "growth", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-wohlstand", DeckID: "c2-finance", Front: "der Wohlstand", Back: "prosperity / wealth", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-bilanz", DeckID: "c2-finance", Front: "die Bilanz", Back: "balance sheet", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-finanzierung", DeckID: "c2-finance", Front: "die Finanzierung", Back: "financing / funding", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-gläubiger", DeckID: "c2-finance", Front: "der Gläubiger", Back: "creditor", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-schuldner", DeckID: "c2-finance", Front: "der Schuldner", Back: "debtor", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-tilgung", DeckID: "c2-finance", Front: "die Tilgung", Back: "repayment / amortization", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-börsencrash", DeckID: "c2-finance", Front: "der Börsencrash", Back: "stock market crash", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-aufschwung", DeckID: "c2-finance", Front: "der Aufschwung", Back: "upswing / boom", Tags: []string{"c2", "finance", "noun"}},
		{ID: "c2-fin-abschwung", DeckID: "c2-finance", Front: "der Abschwung", Back: "downturn", Tags: []string{"c2", "finance", "noun"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "c2-finance",
		Name:        "German C2 Finance & Economics",
		Description: "Advanced financial, economic, and business terminology for C2 learners.",
		Tags:        []string{"german", "c2", "finance", "economics"},
		Notes:       notes,
	}
}
