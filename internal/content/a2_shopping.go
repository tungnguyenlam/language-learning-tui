package content

import (
	"deutsch-tui/internal/core"
)

func A2ShoppingDeck() core.Deck {
	deckID := "a2_purchasing_wear"
	notes := []core.Note{
		{ID: "a2-shop-kaufen", DeckID: deckID, Front: "kaufen", Back: "to buy / to purchase", Tags: []string{"a2", "shopping", "verb"}},
		{ID: "a2-shop-laden", DeckID: deckID, Front: "der Laden", Back: "shop / store", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-supermarkt", DeckID: deckID, Front: "der Supermarkt", Back: "supermarket", Tags: []string{"a1", "shopping", "noun"}},
		{ID: "a2-shop-kaufhaus", DeckID: deckID, Front: "das Kaufhaus", Back: "department store", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-boutique", DeckID: deckID, Front: "die Boutique", Back: "boutique", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-markt", DeckID: deckID, Front: "der Markt", Back: "market", Tags: []string{"a1", "shopping", "noun"}},
		{ID: "a2-shop-kasse", DeckID: deckID, Front: "die Kasse", Back: "checkout / cash register", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-preis", DeckID: deckID, Front: "der Preis", Back: "price", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-rabatt", DeckID: deckID, Front: "der Rabatt", Back: "discount", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-angebot", DeckID: deckID, Front: "das Angebot", Back: "offer / special", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-bar", DeckID: deckID, Front: "bar zahlen", Back: "to pay cash", Tags: []string{"a2", "shopping", "phrase"}},
		{ID: "a2-shop-karte", DeckID: deckID, Front: "mit Karte zahlen", Back: "to pay by card", Tags: []string{"a2", "shopping", "phrase"}},
		{ID: "a2-shop-geld", DeckID: deckID, Front: "das Geld", Back: "money", Tags: []string{"a1", "shopping", "noun"}},
		{ID: "a2-shop-münze", DeckID: deckID, Front: "die Münze", Back: "coin", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-schein", DeckID: deckID, Front: "der Schein", Back: "banknote / bill", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-wechselgeld", DeckID: deckID, Front: "das Wechselgeld", Back: "change (money)", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-teuer", DeckID: deckID, Front: "teuer", Back: "expensive", Tags: []string{"a2", "shopping", "adjective"}},
		{ID: "a2-shop-günstig", DeckID: deckID, Front: "günstig", Back: "inexpensive / good value", Tags: []string{"a2", "shopping", "adjective"}},
		{ID: "a2-shop-kostenlos", DeckID: deckID, Front: "kostenlos", Back: "free (of charge)", Tags: []string{"a2", "shopping", "adjective"}},
		{ID: "a2-shop-größe", DeckID: deckID, Front: "die Größe", Back: "size", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-farbe", DeckID: deckID, Front: "die Farbe", Back: "color", Tags: []string{"a1", "shopping", "noun"}},
		{ID: "a2-shop-anprobieren", DeckID: deckID, Front: "anprobieren", Back: "to try on", Tags: []string{"a2", "shopping", "verb"}},
		{ID: "a2-shop-ausverkauf", DeckID: deckID, Front: "der Ausverkauf", Back: "sale (clearance)", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-einkaufstüte", DeckID: deckID, Front: "die Einkaufstüte", Back: "shopping bag", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-einkaufswagen", DeckID: deckID, Front: "der Einkaufswagen", Back: "shopping cart", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-kleidung", DeckID: deckID, Front: "die Kleidung", Back: "clothing / clothes", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-schuh", DeckID: deckID, Front: "der Schuh / die Schuhe", Back: "shoe / shoes", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-hemd", DeckID: deckID, Front: "das Hemd", Back: "shirt", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-hose", DeckID: deckID, Front: "die Hose", Back: "pants / trousers", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-jacke", DeckID: deckID, Front: "die Jacke", Back: "jacket", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-kleid", DeckID: deckID, Front: "das Kleid", Back: "dress", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-rock", DeckID: deckID, Front: "der Rock", Back: "skirt", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-mütze", DeckID: deckID, Front: "die Mütze", Back: "cap / beanie", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-handschuh", DeckID: deckID, Front: "der Handschuh", Back: "glove", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-schal", DeckID: deckID, Front: "der Schal", Back: "scarf", Tags: []string{"a2", "shopping", "noun"}},
		{ID: "a2-shop-können", DeckID: deckID, Front: "Kann ich das anprobieren?", Back: "Can I try this on?", Tags: []string{"a2", "shopping", "phrase"}},
		{ID: "a2-shop-kosten", DeckID: deckID, Front: "Was kostet das?", Back: "How much does this cost?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a2-shop-pay", DeckID: deckID, Front: "Ich möchte das kaufen.", Back: "I would like to buy this.", Tags: []string{"a2", "shopping", "phrase"}},
		{ID: "a2-shop-credit", DeckID: deckID, Front: "Kann ich mit Karte zahlen?", Back: "Can I pay by card?", Tags: []string{"a2", "shopping", "phrase"}},
		{ID: "a2-shop-reduziert", DeckID: deckID, Front: "reduziert", Back: "reduced (on sale)", Tags: []string{"a2", "shopping", "adjective"}},
		{ID: "a2-shop-wie", DeckID: deckID, Front: "Wie viel kostet das?", Back: "How much does that cost?", Tags: []string{"a1", "shopping", "phrase"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          deckID,
		Name:        "German A2 Shopping & Clothing",
		Description: "Shopping vocabulary, clothing items, payment methods, and service phrases.",
		Tags:        []string{"german", "a2", "shopping", "clothing"},
		Notes:       notes,
	}
}
