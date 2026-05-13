package content

import (
	"deutsch-tui/internal/core"
)

func A1OfficeDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-off-bueroe", DeckID: "a1-office", Front: "das Büro", Back: "office", Tags: []string{"a1", "office"}},
		{ID: "a1-off-schreibtisch", DeckID: "a1-office", Front: "der Schreibtisch", Back: "desk", Tags: []string{"a1", "office"}},
		{ID: "a1-off-stuhl", DeckID: "a1-office", Front: "der Stuhl", Back: "chair", Tags: []string{"a1", "office"}},
		{ID: "a1-off-computer", DeckID: "a1-office", Front: "der Computer", Back: "computer", Tags: []string{"a1", "office"}},
		{ID: "a1-off-laptop", DeckID: "a1-office", Front: "der Laptop", Back: "laptop", Tags: []string{"a1", "office"}},
		{ID: "a1-off-bildschirm", DeckID: "a1-office", Front: "der Bildschirm", Back: "screen / monitor", Tags: []string{"a1", "office"}},
		{ID: "a1-off-tastatur", DeckID: "a1-office", Front: "die Tastatur", Back: "keyboard", Tags: []string{"a1", "office"}},
		{ID: "a1-off-maus", DeckID: "a1-office", Front: "die Maus", Back: "mouse", Tags: []string{"a1", "office"}},
		{ID: "a1-off-drucker", DeckID: "a1-office", Front: "der Drucker", Back: "printer", Tags: []string{"a1", "office"}},
		{ID: "a1-off-papier", DeckID: "a1-office", Front: "das Papier", Back: "paper", Tags: []string{"a1", "office"}},
		{ID: "a1-off-blatt", DeckID: "a1-office", Front: "das Blatt", Back: "sheet (of paper)", Tags: []string{"a1", "office"}},
		{ID: "a1-off-heft", DeckID: "a1-office", Front: "das Heft", Back: "notebook", Tags: []string{"a1", "office"}},
		{ID: "a1-off-buch", DeckID: "a1-office", Front: "das Buch", Back: "book", Tags: []string{"a1", "office"}},
		{ID: "a1-off-ordner", DeckID: "a1-office", Front: "der Ordner", Back: "binder / folder", Tags: []string{"a1", "office"}},
		{ID: "a1-off-mappe", DeckID: "a1-office", Front: "die Mappe", Back: "folder", Tags: []string{"a1", "office"}},
		{ID: "a1-off-kugelschr", DeckID: "a1-office", Front: "der Kugelschreiber", Back: "ballpoint pen", Tags: []string{"a1", "office"}},
		{ID: "a1-off-bleistift", DeckID: "a1-office", Front: "der Bleistift", Back: "pencil", Tags: []string{"a1", "office"}},
		{ID: "a1-off-filzstift", DeckID: "a1-office", Front: "der Filzstift", Back: "felt-tip pen", Tags: []string{"a1", "office"}},
		{ID: "a1-off-marker", DeckID: "a1-office", Front: "der Marker", Back: "marker / highlighter", Tags: []string{"a1", "office"}},
		{ID: "a1-off-radiergummi", DeckID: "a1-office", Front: "der Radiergummi", Back: "eraser / rubber", Tags: []string{"a1", "office"}},
		{ID: "a1-off-lineal", DeckID: "a1-office", Front: "das Lineal", Back: "ruler", Tags: []string{"a1", "office"}},
		{ID: "a1-off-schere", DeckID: "a1-office", Front: "die Schere", Back: "scissors", Tags: []string{"a1", "office"}},
		{ID: "a1-off-kleber", DeckID: "a1-office", Front: "der Kleber", Back: "glue", Tags: []string{"a1", "office"}},
		{ID: "a1-off-tesafilm", DeckID: "a1-office", Front: "der Tesafilm", Back: "sticky tape", Tags: []string{"a1", "office"}},
		{ID: "a1-off-tacker", DeckID: "a1-office", Front: "der Tacker", Back: "stapler", Tags: []string{"a1", "office"}},
		{ID: "a1-off-klammer", DeckID: "a1-office", Front: "die Klammer", Back: "paperclip / staple", Tags: []string{"a1", "office"}},
		{ID: "a1-off-umschlag", DeckID: "a1-office", Front: "der Umschlag", Back: "envelope", Tags: []string{"a1", "office"}},
		{ID: "a1-off-brief", DeckID: "a1-office", Front: "der Brief", Back: "letter", Tags: []string{"a1", "office"}},
		{ID: "a1-off-briefmarke", DeckID: "a1-office", Front: "die Briefmarke", Back: "stamp", Tags: []string{"a1", "office"}},
		{ID: "a1-off-paket", DeckID: "a1-office", Front: "das Paket", Back: "package / parcel", Tags: []string{"a1", "office"}},
		{ID: "a1-off-kalender", DeckID: "a1-office", Front: "der Kalender", Back: "calendar", Tags: []string{"a1", "office"}},
		{ID: "a1-off-notiz", DeckID: "a1-office", Front: "die Notiz", Back: "note", Tags: []string{"a1", "office"}},
		{ID: "a1-off-zettel", DeckID: "a1-office", Front: "der Zettel", Back: "slip of paper / note", Tags: []string{"a1", "office"}},
		{ID: "a1-off-telefon", DeckID: "a1-office", Front: "das Telefon", Back: "telephone", Tags: []string{"a1", "office"}},
		{ID: "a1-off-handy", DeckID: "a1-office", Front: "das Handy", Back: "mobile phone", Tags: []string{"a1", "office"}},
		{ID: "a1-off-email", DeckID: "a1-office", Front: "die E-Mail", Back: "email", Tags: []string{"a1", "office"}},
		{ID: "a1-off-besprechung", DeckID: "a1-office", Front: "die Besprechung", Back: "meeting", Tags: []string{"a1", "office"}},
		{ID: "a1-off-arbeit", DeckID: "a1-office", Front: "die Arbeit", Back: "work", Tags: []string{"a1", "office"}},
		{ID: "a1-off-kollege", DeckID: "a1-office", Front: "der Kollege", Back: "colleague", Tags: []string{"a1", "office"}},
		{ID: "a1-off-chef", DeckID: "a1-office", Front: "der Chef", Back: "boss", Tags: []string{"a1", "office"}},
		{ID: "a1-off-firma", DeckID: "a1-office", Front: "die Firma", Back: "company / firm", Tags: []string{"a1", "office"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-office",
		Name:        "A1 German Office & Stationery",
		Description: "Beginner vocabulary for office items, stationery, and common workplace objects.",
		Tags:        []string{"german", "a1", "office", "stationery", "work"},
		Notes:       notes,
	}
}
