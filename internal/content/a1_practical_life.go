package content

import (
	"deutsch-tui/internal/core"
)

func A1PracticalLifeDeck() core.Deck {
	notes := []core.Note{
		// Restaurant & Cafe
		{ID: "a1-prac-res-tisch", DeckID: "a1-practical", Front: "Einen Tisch für zwei Personen, bitte.", Back: "A table for two, please.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-karte", DeckID: "a1-practical", Front: "Die Speisekarte, bitte.", Back: "The menu, please.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-bestellen", DeckID: "a1-practical", Front: "Ich möchte bestellen.", Back: "I would like to order.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-wasser", DeckID: "a1-practical", Front: "Ein Wasser, bitte.", Back: "A water, please.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-zahlen", DeckID: "a1-practical", Front: "Zahlen, bitte.", Back: "The bill, please.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-zusammen", DeckID: "a1-practical", Front: "Zusammen oder getrennt?", Back: "Together or separately?", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-stimmt", DeckID: "a1-practical", Front: "Stimmt so.", Back: "Keep the change.", Tags: []string{"a1", "restaurant", "phrase"}},
		{ID: "a1-prac-res-lecker", DeckID: "a1-practical", Front: "Das Essen war sehr lecker.", Back: "The food was very delicious.", Tags: []string{"a1", "restaurant", "phrase"}},

		// Shopping
		{ID: "a1-prac-sho-hilfe", DeckID: "a1-practical", Front: "Kann ich Ihnen helfen?", Back: "Can I help you?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-suche", DeckID: "a1-practical", Front: "Ich suche ein Hemd.", Back: "I am looking for a shirt.", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-kostet", DeckID: "a1-practical", Front: "Was kostet das?", Back: "How much does that cost?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-teuer", DeckID: "a1-practical", Front: "Das ist zu teuer.", Back: "That is too expensive.", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-angebot", DeckID: "a1-practical", Front: "Ist das im Angebot?", Back: "Is this on sale?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-probieren", DeckID: "a1-practical", Front: "Kann ich das anprobieren?", Back: "Can I try this on?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-kabine", DeckID: "a1-practical", Front: "Wo ist die Umkleidekabine?", Back: "Where is the fitting room?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-groesse", DeckID: "a1-practical", Front: "Haben Sie das in Größe M?", Back: "Do you have this in size M?", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-nehme", DeckID: "a1-practical", Front: "Ich nehme es.", Back: "I'll take it.", Tags: []string{"a1", "shopping", "phrase"}},
		{ID: "a1-prac-sho-karte", DeckID: "a1-practical", Front: "Kann ich mit Karte zahlen?", Back: "Can I pay by card?", Tags: []string{"a1", "shopping", "phrase"}},

		// Directions & City
		{ID: "a1-prac-dir-entschuldigung", DeckID: "a1-practical", Front: "Entschuldigung, wo ist der Bahnhof?", Back: "Excuse me, where is the train station?", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-weg", DeckID: "a1-practical", Front: "Ich habe mich verlaufen.", Back: "I am lost.", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-weit", DeckID: "a1-practical", Front: "Ist es weit von hier?", Back: "Is it far from here?", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-links", DeckID: "a1-practical", Front: "Gehen Sie nach links.", Back: "Go to the left.", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-rechts", DeckID: "a1-practical", Front: "Gehen Sie nach rechts.", Back: "Go to the right.", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-geradeaus", DeckID: "a1-practical", Front: "Gehen Sie geradeaus.", Back: "Go straight ahead.", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-naechste", DeckID: "a1-practical", Front: "Die nächste Straße rechts.", Back: "The next street on the right.", Tags: []string{"a1", "directions", "phrase"}},
		{ID: "a1-prac-dir-bus", DeckID: "a1-practical", Front: "Welcher Bus fährt zum Zentrum?", Back: "Which bus goes to the center?", Tags: []string{"a1", "directions", "phrase"}},

		// Social & Small Talk
		{ID: "a1-prac-soc-geht", DeckID: "a1-practical", Front: "Wie geht es dir?", Back: "How are you? (informal)", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-freut", DeckID: "a1-practical", Front: "Freut mich, dich kennenzulernen.", Back: "Nice to meet you.", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-herkunft", DeckID: "a1-practical", Front: "Woher kommst du?", Back: "Where do you come from?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-wohnort", DeckID: "a1-practical", Front: "Wo wohnst du?", Back: "Where do you live?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-beruf", DeckID: "a1-practical", Front: "Was bist du von Beruf?", Back: "What is your profession?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-hobby", DeckID: "a1-practical", Front: "Was sind deine Hobbies?", Back: "What are your hobbies?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-deutsch", DeckID: "a1-practical", Front: "Mein Deutsch ist nicht so gut.", Back: "My German is not so good.", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-langsam", DeckID: "a1-practical", Front: "Können Sie bitte langsamer sprechen?", Back: "Could you please speak slower?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-wiederholen", DeckID: "a1-practical", Front: "Können Sie das bitte wiederholen?", Back: "Could you please repeat that?", Tags: []string{"a1", "social", "phrase"}},
		{ID: "a1-prac-soc-verstanden", DeckID: "a1-practical", Front: "Ich habe das nicht verstanden.", Back: "I didn't understand that.", Tags: []string{"a1", "social", "phrase"}},

		// Hotel & Stay
		{ID: "a1-prac-hot-reservierung", DeckID: "a1-practical", Front: "Ich habe eine Reservierung.", Back: "I have a reservation.", Tags: []string{"a1", "hotel", "phrase"}},
		{ID: "a1-prac-hot-frei", DeckID: "a1-practical", Front: "Haben Sie noch ein Zimmer frei?", Back: "Do you still have a room available?", Tags: []string{"a1-practical", "hotel", "phrase"}},
		{ID: "a1-prac-hot-einzel", DeckID: "a1-practical", Front: "Ein Einzelzimmer für zwei Nächte.", Back: "A single room for two nights.", Tags: []string{"a1", "hotel", "phrase"}},
		{ID: "a1-prac-hot-doppel", DeckID: "a1-practical", Front: "Ein Doppelzimmer, bitte.", Back: "A double room, please.", Tags: []string{"a1", "hotel", "phrase"}},
		{ID: "a1-prac-hot-fruehstueck", DeckID: "a1-practical", Front: "Ist das Frühstück inklusive?", Back: "Is breakfast included?", Tags: []string{"a1", "hotel", "phrase"}},
		{ID: "a1-prac-hot-schluessel", DeckID: "a1-practical", Front: "Hier ist Ihr Schlüssel.", Back: "Here is your key.", Tags: []string{"a1", "hotel", "phrase"}},
		{ID: "a1-prac-hot-wlan", DeckID: "a1-practical", Front: "Wie ist das WLAN-Passwort?", Back: "What is the Wi-Fi password?", Tags: []string{"a1", "hotel", "phrase"}},

		// Emergency & Help
		{ID: "a1-prac-hel-hilfe", DeckID: "a1-practical", Front: "Hilfe!", Back: "Help!", Tags: []string{"a1", "emergency", "phrase"}},
		{ID: "a1-prac-hel-polizei", DeckID: "a1-practical", Front: "Rufen Sie die Polizei!", Back: "Call the police!", Tags: []string{"a1", "emergency", "phrase"}},
		{ID: "a1-prac-hel-arzt", DeckID: "a1-practical", Front: "Ich brauche einen Arzt.", Back: "I need a doctor.", Tags: []string{"a1", "emergency", "phrase"}},
		{ID: "a1-prac-hel-krank", DeckID: "a1-practical", Front: "Ich fühle mich nicht gut.", Back: "I don't feel well.", Tags: []string{"a1", "emergency", "phrase"}},
		{ID: "a1-prac-hel-verloren", DeckID: "a1-practical", Front: "Ich habe meine Tasche verloren.", Back: "I lost my bag.", Tags: []string{"a1", "emergency", "phrase"}},
		{ID: "a1-prac-hel-gestohlen", DeckID: "a1-practical", Front: "Mein Handy wurde gestohlen.", Back: "My phone was stolen.", Tags: []string{"a1", "emergency", "phrase"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-practical",
		Name:        "A1 Practical Phrases",
		Description: "Essential A1 phrases for everyday situations in Germany.",
		Tags:        []string{"german", "a1", "phrases"},
		Notes:       notes,
	}
}
