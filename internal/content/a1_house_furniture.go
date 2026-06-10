package content

import (
	"deutsch-tui/internal/core"
)

func A1HouseFurnitureDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-hm-haus", DeckID: "a1-house-furniture", Front: "das Haus", Back: "house", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-wohnung", DeckID: "a1-house-furniture", Front: "die Wohnung", Back: "apartment / flat", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-zimmer", DeckID: "a1-house-furniture", Front: "das Zimmer", Back: "room", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-tuer", DeckID: "a1-house-furniture", Front: "die Tür", Back: "door", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-fenster", DeckID: "a1-house-furniture", Front: "das Fenster", Back: "window", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-wand", DeckID: "a1-house-furniture", Front: "die Wand", Back: "wall", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-boden", DeckID: "a1-house-furniture", Front: "der Boden", Back: "floor / ground", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-decke", DeckID: "a1-house-furniture", Front: "die Decke", Back: "ceiling", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-treppe", DeckID: "a1-house-furniture", Front: "die Treppe", Back: "stairs", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-balkon", DeckID: "a1-house-furniture", Front: "der Balkon", Back: "balcony", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-garten", DeckID: "a1-house-furniture", Front: "der Garten", Back: "garden", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-keller", DeckID: "a1-house-furniture", Front: "der Keller", Back: "basement / cellar", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-wohnzimmer", DeckID: "a1-house-furniture", Front: "das Wohnzimmer", Back: "living room", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-schlafzimmer", DeckID: "a1-house-furniture", Front: "das Schlafzimmer", Back: "bedroom", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-kueche", DeckID: "a1-house-furniture", Front: "die Küche", Back: "kitchen", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-bad", DeckID: "a1-house-furniture", Front: "das Bad", Back: "bathroom", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-kinderzimmer", DeckID: "a1-house-furniture", Front: "das Kinderzimmer", Back: "children's room", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-gaestezimmer", DeckID: "a1-house-furniture", Front: "das Gästezimmer", Back: "guest room", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-flur", DeckID: "a1-house-furniture", Front: "der Flur", Back: "hallway / corridor", Tags: []string{"a1", "rooms"}},
		{ID: "a1-hm-tisch", DeckID: "a1-house-furniture", Front: "der Tisch", Back: "table", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-stuhl", DeckID: "a1-house-furniture", Front: "der Stuhl", Back: "chair", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-sessel", DeckID: "a1-house-furniture", Front: "der Sessel", Back: "armchair", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-sofa", DeckID: "a1-house-furniture", Front: "das Sofa", Back: "sofa / couch", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-schrank", DeckID: "a1-house-furniture", Front: "der Schrank", Back: "cupboard / wardrobe", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-bett", DeckID: "a1-house-furniture", Front: "das Bett", Back: "bed", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-regal", DeckID: "a1-house-furniture", Front: "das Regal", Back: "shelf / shelves", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-schreibtisch", DeckID: "a1-house-furniture", Front: "der Schreibtisch", Back: "desk", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-kommode", DeckID: "a1-house-furniture", Front: "die Kommode", Back: "chest of drawers", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-lampe", DeckID: "a1-house-furniture", Front: "die Lampe", Back: "lamp", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-teppich", DeckID: "a1-house-furniture", Front: "der Teppich", Back: "carpet / rug", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-vorhang", DeckID: "a1-house-furniture", Front: "der Vorhang", Back: "curtain", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-kissen", DeckID: "a1-house-furniture", Front: "das Kissen", Back: "cushion / pillow", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-spiegel", DeckID: "a1-house-furniture", Front: "der Spiegel", Back: "mirror", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-bild", DeckID: "a1-house-furniture", Front: "das Bild", Back: "picture / painting", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-herd", DeckID: "a1-house-furniture", Front: "der Herd", Back: "stove / cooker", Tags: []string{"a1", "kitchen"}},
		{ID: "a1-hm-kuehlschrank", DeckID: "a1-house-furniture", Front: "der Kühlschrank", Back: "refrigerator", Tags: []string{"a1", "kitchen"}},
		{ID: "a1-hm-ofen", DeckID: "a1-house-furniture", Front: "der Ofen", Back: "oven", Tags: []string{"a1", "kitchen"}},
		{ID: "a1-hm-spuele", DeckID: "a1-house-furniture", Front: "die Spüle", Back: "sink", Tags: []string{"a1", "kitchen"}},
		{ID: "a1-hm-geschirr", DeckID: "a1-house-furniture", Front: "das Geschirr", Back: "dishes / crockery", Tags: []string{"a1", "kitchen"}},
		{ID: "a1-hm-waschmaschine", DeckID: "a1-house-furniture", Front: "die Waschmaschine", Back: "washing machine", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-muell", DeckID: "a1-house-furniture", Front: "der Müll", Back: "trash / rubbish", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-wohnen", DeckID: "a1-house-furniture", Front: "wohnen", Back: "to live (somewhere)", Tags: []string{"a1", "house", "verb"}},
		{ID: "a1-hm-leben", DeckID: "a1-house-furniture", Front: "leben", Back: "to live", Tags: []string{"a1", "house", "verb"}},
		{ID: "a1-hm-wohnort", DeckID: "a1-house-furniture", Front: "der Wohnort", Back: "place of residence", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-putzen", DeckID: "a1-house-furniture", Front: "putzen", Back: "to clean", Tags: []string{"a1", "house", "verb"}},
		{ID: "a1-hm-wohnungsschluessel", DeckID: "a1-house-furniture", Front: "der Wohnungsschlüssel", Back: "apartment key", Tags: []string{"a1", "house"}},
		{ID: "a1-hm-wohnzimmer-tisch", DeckID: "a1-house-furniture", Front: "der Couchtisch", Back: "coffee table", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-nachttisch", DeckID: "a1-house-furniture", Front: "der Nachttisch", Back: "nightstand", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-uhr", DeckID: "a1-house-furniture", Front: "die Uhr", Back: "clock", Tags: []string{"a1", "furniture"}},
		{ID: "a1-hm-pflanze", DeckID: "a1-house-furniture", Front: "die Pflanze", Back: "plant", Tags: []string{"a1", "house"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-house-furniture",
		Name:        "A1 House & Furniture",
		Description: "Essential A1 vocabulary for rooms, furniture, and household items.",
		Tags:        []string{"german", "a1", "house", "furniture"},
		Notes:       notes,
	}
}
