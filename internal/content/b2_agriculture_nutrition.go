package content

import (
	"deutsch-tui/internal/core"
)

func B2AgricultureNutritionDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-agr-landwirtschaft", DeckID: "b2-agriculture-nutrition", Front: "die Landwirtschaft", Back: "agriculture / farming", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-landwirt", DeckID: "b2-agriculture-nutrition", Front: "der Landwirt / die Landwirtin", Back: "farmer", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-bauernhof", DeckID: "b2-agriculture-nutrition", Front: "der Bauernhof", Back: "farm", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-ackerbau", DeckID: "b2-agriculture-nutrition", Front: "der Ackerbau", Back: "arable farming / crop cultivation", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-viehzucht", DeckID: "b2-agriculture-nutrition", Front: "die Viehzucht", Back: "livestock breeding", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-vieh", DeckID: "b2-agriculture-nutrition", Front: "das Vieh", Back: "livestock / cattle", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-ernte", DeckID: "b2-agriculture-nutrition", Front: "die Ernte", Back: "harvest", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-aussaat", DeckID: "b2-agriculture-nutrition", Front: "die Aussaat", Back: "sowing", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-saatgut", DeckID: "b2-agriculture-nutrition", Front: "das Saatgut", Back: "seed (for sowing)", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-boden", DeckID: "b2-agriculture-nutrition", Front: "der Boden", Back: "soil / ground", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-bodenfruchtbarkeit", DeckID: "b2-agriculture-nutrition", Front: "die Bodenfruchtbarkeit", Back: "soil fertility", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-duenger", DeckID: "b2-agriculture-nutrition", Front: "der Dünger", Back: "fertilizer", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-pestizid", DeckID: "b2-agriculture-nutrition", Front: "das Pestizid", Back: "pesticide", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-unkraut", DeckID: "b2-agriculture-nutrition", Front: "das Unkraut", Back: "weeds", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-bewaesserung", DeckID: "b2-agriculture-nutrition", Front: "die Bewässerung", Back: "irrigation", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-gewaechshaus", DeckID: "b2-agriculture-nutrition", Front: "das Gewächshaus", Back: "greenhouse", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-traktor", DeckID: "b2-agriculture-nutrition", Front: "der Traktor", Back: "tractor", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-acker", DeckID: "b2-agriculture-nutrition", Front: "der Acker", Back: "field (arable land)", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-weide", DeckID: "b2-agriculture-nutrition", Front: "die Weide", Back: "pasture", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-getreide", DeckID: "b2-agriculture-nutrition", Front: "das Getreide", Back: "grain / cereal crops", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-weizen", DeckID: "b2-agriculture-nutrition", Front: "der Weizen", Back: "wheat", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-roggen", DeckID: "b2-agriculture-nutrition", Front: "der Roggen", Back: "rye", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-gerste", DeckID: "b2-agriculture-nutrition", Front: "die Gerste", Back: "barley", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-hafer", DeckID: "b2-agriculture-nutrition", Front: "der Hafer", Back: "oats", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-mais", DeckID: "b2-agriculture-nutrition", Front: "der Mais", Back: "maize / corn", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-raps", DeckID: "b2-agriculture-nutrition", Front: "der Raps", Back: "rapeseed", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-zuckerruebe", DeckID: "b2-agriculture-nutrition", Front: "die Zuckerrübe", Back: "sugar beet", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-ertrag", DeckID: "b2-agriculture-nutrition", Front: "der Ertrag", Back: "yield / return", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-oekolandbau", DeckID: "b2-agriculture-nutrition", Front: "der Ökolandbau", Back: "organic farming", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-massentierhaltung", DeckID: "b2-agriculture-nutrition", Front: "die Massentierhaltung", Back: "factory farming", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-tierschutz", DeckID: "b2-agriculture-nutrition", Front: "der Tierschutz", Back: "animal welfare / animal protection", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-genossenschaft", DeckID: "b2-agriculture-nutrition", Front: "die Genossenschaft", Back: "cooperative", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-subvention", DeckID: "b2-agriculture-nutrition", Front: "die Subvention", Back: "subsidy", Tags: []string{"b2", "agriculture"}},
		{ID: "b2-agr-ernaehrung", DeckID: "b2-agriculture-nutrition", Front: "die Ernährung", Back: "nutrition / diet", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-nahrungsmittel", DeckID: "b2-agriculture-nutrition", Front: "das Nahrungsmittel", Back: "foodstuff", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-lebensmittelsicherheit", DeckID: "b2-agriculture-nutrition", Front: "die Lebensmittelsicherheit", Back: "food safety", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-naehrstoff", DeckID: "b2-agriculture-nutrition", Front: "der Nährstoff", Back: "nutrient", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-eiweiss", DeckID: "b2-agriculture-nutrition", Front: "das Eiweiß / das Protein", Back: "protein", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-kohlenhydrat", DeckID: "b2-agriculture-nutrition", Front: "das Kohlenhydrat", Back: "carbohydrate", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-ballaststoff", DeckID: "b2-agriculture-nutrition", Front: "der Ballaststoff", Back: "dietary fibre", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-vitamin", DeckID: "b2-agriculture-nutrition", Front: "das Vitamin", Back: "vitamin", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-kalorie", DeckID: "b2-agriculture-nutrition", Front: "die Kalorie", Back: "calorie", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-verschwendung", DeckID: "b2-agriculture-nutrition", Front: "die Lebensmittelverschwendung", Back: "food waste", Tags: []string{"b2", "nutrition"}},
		{ID: "b2-agr-anbauen", DeckID: "b2-agriculture-nutrition", Front: "anbauen", Back: "to grow / to cultivate", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-ernten", DeckID: "b2-agriculture-nutrition", Front: "ernten", Back: "to harvest", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-saeen", DeckID: "b2-agriculture-nutrition", Front: "säen", Back: "to sow", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-duengen", DeckID: "b2-agriculture-nutrition", Front: "düngen", Back: "to fertilize", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-bewaessern", DeckID: "b2-agriculture-nutrition", Front: "bewässern", Back: "to irrigate", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-zuechten", DeckID: "b2-agriculture-nutrition", Front: "züchten", Back: "to breed / to rear", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-melken", DeckID: "b2-agriculture-nutrition", Front: "melken", Back: "to milk", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-pfluegen", DeckID: "b2-agriculture-nutrition", Front: "pflügen", Back: "to plough", Tags: []string{"b2", "agriculture", "verb"}},
		{ID: "b2-agr-verarbeiten", DeckID: "b2-agriculture-nutrition", Front: "verarbeiten", Back: "to process", Tags: []string{"b2", "nutrition", "verb"}},
		{ID: "b2-agr-verderben", DeckID: "b2-agriculture-nutrition", Front: "verderben", Back: "to spoil / to go bad", Tags: []string{"b2", "nutrition", "verb"}},
		{ID: "b2-agr-oekologisch", DeckID: "b2-agriculture-nutrition", Front: "ökologisch", Back: "ecological / organic", Tags: []string{"b2", "agriculture", "adjective"}},
		{ID: "b2-agr-ertragreich", DeckID: "b2-agriculture-nutrition", Front: "ertragreich", Back: "high-yielding", Tags: []string{"b2", "agriculture", "adjective"}},
		{ID: "b2-agr-nahrhaft", DeckID: "b2-agriculture-nutrition", Front: "nahrhaft", Back: "nutritious", Tags: []string{"b2", "nutrition", "adjective"}},
		{ID: "b2-agr-saisonal", DeckID: "b2-agriculture-nutrition", Front: "saisonal", Back: "seasonal", Tags: []string{"b2", "nutrition", "adjective"}},
		{ID: "b2-agr-haltbar", DeckID: "b2-agriculture-nutrition", Front: "haltbar", Back: "long-lasting / non-perishable", Tags: []string{"b2", "nutrition", "adjective"}},
		{ID: "b2-agr-gentechnisch", DeckID: "b2-agriculture-nutrition", Front: "gentechnisch verändert", Back: "genetically modified", Tags: []string{"b2", "agriculture", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-agriculture-nutrition",
		Name:        "B2 Landwirtschaft und Ernährung",
		Description: "Farming, crops, livestock and nutrition vocabulary at B2 level.",
		Tags:        []string{"german", "b2", "agriculture", "nutrition"},
		Notes:       notes,
	}
}
