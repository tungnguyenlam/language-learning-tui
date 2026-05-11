package content

import (
	"deutsch-tui/internal/core"
)

func A2ShoppingDeck() core.Deck {
	notes := []core.Note{
		{ID: "a2-shp-einkaufen", DeckID: "a2-shopping", Front: "einkaufen", Back: "to shop / do shopping", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-laden", DeckID: "a2-shopping", Front: "der Laden", Back: "shop / store", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-kaufhaus", DeckID: "a2-shopping", Front: "das Kaufhaus", Back: "department store", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-supermarkt", DeckID: "a2-shopping", Front: "der Supermarkt", Back: "supermarket", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-discounter", DeckID: "a2-shopping", Front: "der Discounter", Back: "discount store", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-markt", DeckID: "a2-shopping", Front: "der Markt", Back: "market", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-basar", DeckID: "a2-shopping", Front: "der Basar", Back: "bazaar", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-handel", DeckID: "a2-shopping", Front: "der Handel", Back: "trade / commerce", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-kasse", DeckID: "a2-shopping", Front: "die Kasse", Back: "cash register / checkout", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-preis", DeckID: "a2-shopping", Front: "der Preis", Back: "price", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-angebot", DeckID: "a2-shopping", Front: "das Angebot", Back: "offer / special", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-rabatt", DeckID: "a2-shopping", Front: "der Rabatt", Back: "discount", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-sale", DeckID: "a2-shopping", Front: "der Ausverkauf", Back: "sale", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-bar", DeckID: "a2-shopping", Front: "bar zahlen", Back: "to pay cash", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-karte", DeckID: "a2-shopping", Front: "mit Karte zahlen", Back: "to pay by card", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-geld", DeckID: "a2-shopping", Front: "das Geld", Back: "money", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-wechselgeld", DeckID: "a2-shopping", Front: "das Wechselgeld", Back: "change (money)", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-quittung", DeckID: "a2-shopping", Front: "die Quittung", Back: "receipt", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-tüte", DeckID: "a2-shopping", Front: "die Tüte", Back: "bag (shopping)", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-waren", DeckID: "a2-shopping", Front: "die Ware", Back: "goods / merchandise", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-ware", DeckID: "a2-shopping", Front: "die Ware", Back: "product", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-größe", DeckID: "a2-shopping", Front: "die Größe", Back: "size", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-farbe", DeckID: "a2-shopping", Front: "die Farbe", Back: "color", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-anprobieren", DeckID: "a2-shopping", Front: "anprobieren", Back: "to try on", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-passen", DeckID: "a2-shopping", Front: "passen", Back: "to fit / to suit", Extra: "Das passt mir gut.", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-kaufen", DeckID: "a2-shopping", Front: "kaufen", Back: "to buy / to purchase", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-verkaufen", DeckID: "a2-shopping", Front: "verkaufen", Back: "to sell", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-verkäufer", DeckID: "a2-shopping", Front: "der Verkäufer / die Verkäuferin", Back: "salesperson", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-kunde", DeckID: "a2-shopping", Front: "der Kunde / die Kundin", Back: "customer", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-wunsch", DeckID: "a2-shopping", Front: "der Wunsch", Back: "wish / desire", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-wünschen", DeckID: "a2-shopping", Front: "wünschen", Back: "to wish / to want", Extra: "Ich wünsche Ihnen einen guten Tag.", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-dienstleistung", DeckID: "a2-shopping", Front: "die Dienstleistung", Back: "service", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-reparatur", DeckID: "a2-shopping", Front: "die Reparatur", Back: "repair", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-friseur", DeckID: "a2-shopping", Front: "der Friseur / die Friseurin", Back: "hairdresser", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-post", DeckID: "a2-shopping", Front: "die Post", Back: "post office / mail", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-bank", DeckID: "a2-shopping", Front: "die Bank", Back: "bank", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-apotheke", DeckID: "a2-shopping", Front: "die Apotheke", Back: "pharmacy", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-arzt", DeckID: "a2-shopping", Front: "der Arzt / die Ärztin", Back: "doctor", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-reinigung", DeckID: "a2-shopping", Front: "die Reinigung", Back: "dry cleaning", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-tankstelle", DeckID: "a2-shopping", Front: "die Tankstelle", Back: "gas station", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-waschanlage", DeckID: "a2-shopping", Front: "die Waschanlage", Back: "car wash", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-reifen", DeckID: "a2-shopping", Front: "der Reifen", Back: "tire", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-behälter", DeckID: "a2-shopping", Front: "der Behälter", Back: "container", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-verpackung", DeckID: "a2-shopping", Front: "die Verpackung", Back: "packaging", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-menge", DeckID: "a2-shopping", Front: "die Menge", Back: "quantity / amount", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-kilo", DeckID: "a2-shopping", Front: "das Kilo", Back: "kilogram", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-gramm", DeckID: "a2-shopping", Front: "das Gramm", Back: "gram", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-liter", DeckID: "a2-shopping", Front: "der Liter", Back: "liter", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-stück", DeckID: "a2-shopping", Front: "das Stück", Back: "piece", Tags: []string{"a2", "shopping"}},
		{ID: "a2-shp-paar", DeckID: "a2-shopping", Front: "das Paar", Back: "pair", Tags: []string{"a2", "shopping"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-shopping",
		Name:        "A2 Shopping & Services",
		Description: "Vocabulary for shopping, services, and commerce.",
		Tags:        []string{"german", "a2", "shopping", "services"},
		Notes:       notes,
	}
}
