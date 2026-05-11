package content

import (
	"deutsch-tui/internal/core"
)

func B1ArtDeck() core.Deck {
	notes := []core.Note{
		{ID: "art-001", DeckID: "b1-art", Front: "die Kunst", Back: "art", Extra: "Plural: die Künste", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-002", DeckID: "b1-art", Front: "der Künstler", Back: "artist (male)", Extra: "Plural: die Künstler", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-003", DeckID: "b1-art", Front: "die Künstlerin", Back: "artist (female)", Extra: "Plural: die Künstlerinnen", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-004", DeckID: "b1-art", Front: "das Gemälde", Back: "painting", Extra: "Plural: die Gemälde", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-005", DeckID: "b1-art", Front: "die Ausstellung", Back: "exhibition", Extra: "Plural: die Ausstellungen", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-006", DeckID: "b1-art", Front: "die Galerie", Back: "gallery", Extra: "Plural: die Galerien", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-007", DeckID: "b1-art", Front: "das Museum", Back: "museum", Extra: "Plural: die Museen", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-008", DeckID: "b1-art", Front: "zeichnen", Back: "to draw", Extra: "Konjugation: ich zeichne", Tags: []string{"b1", "art", "verb"}},
		{ID: "art-009", DeckID: "b1-art", Front: "malen", Back: "to paint", Extra: "Konjugation: ich male", Tags: []string{"b1", "art", "verb"}},
		{ID: "art-010", DeckID: "b1-art", Front: "die Skulptur", Back: "sculpture", Extra: "Plural: die Skulpturen", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-011", DeckID: "b1-art", Front: "die Architektur", Back: "architecture", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-012", DeckID: "b1-art", Front: "der Pinsel", Back: "paintbrush", Extra: "Plural: die Pinsel", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-013", DeckID: "b1-art", Front: "die Leinwand", Back: "canvas", Extra: "Plural: die Leinwände", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-014", DeckID: "b1-art", Front: "die Literatur", Back: "literature", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-015", DeckID: "b1-art", Front: "der Autor", Back: "author (male)", Extra: "Plural: die Autoren", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-016", DeckID: "b1-art", Front: "die Autorin", Back: "author (female)", Extra: "Plural: die Autorinnen", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-017", DeckID: "b1-art", Front: "der Roman", Back: "novel", Extra: "Plural: die Romane", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-018", DeckID: "b1-art", Front: "das Gedicht", Back: "poem", Extra: "Plural: die Gedichte", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-019", DeckID: "b1-art", Front: "der Dichter", Back: "poet", Extra: "Plural: die Dichter", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-020", DeckID: "b1-art", Front: "das Theaterstück", Back: "play (theatre)", Extra: "Plural: die Theaterstücke", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-021", DeckID: "b1-art", Front: "der Verlag", Back: "publisher", Extra: "Plural: die Verlage", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-022", DeckID: "b1-art", Front: "veröffentlichen", Back: "to publish", Extra: "Konjugation: ich veröffentliche", Tags: []string{"b1", "literature", "verb"}},
		{ID: "art-023", DeckID: "b1-art", Front: "das Kapitel", Back: "chapter", Extra: "Plural: die Kapitel", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-024", DeckID: "b1-art", Front: "die Seite", Back: "page", Extra: "Plural: die Seiten", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-025", DeckID: "b1-art", Front: "lesen", Back: "to read", Extra: "Konjugation: ich lese, du liest", Tags: []string{"b1", "literature", "verb"}},
		{ID: "art-026", DeckID: "b1-art", Front: "schreiben", Back: "to write", Extra: "Konjugation: ich schreibe", Tags: []string{"b1", "literature", "verb"}},
		{ID: "art-027", DeckID: "b1-art", Front: "die Hauptfigur", Back: "main character", Extra: "Plural: die Hauptfiguren", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-028", DeckID: "b1-art", Front: "die Handlung", Back: "plot, storyline", Extra: "Plural: die Handlungen", Tags: []string{"b1", "literature", "noun"}},
		{ID: "art-029", DeckID: "b1-art", Front: "spannend", Back: "exciting, thrilling", Tags: []string{"b1", "literature", "adjective"}},
		{ID: "art-030", DeckID: "b1-art", Front: "langweilig", Back: "boring", Tags: []string{"b1", "literature", "adjective"}},
		{ID: "art-031", DeckID: "b1-art", Front: "kreativ", Back: "creative", Tags: []string{"b1", "art", "adjective"}},
		{ID: "art-032", DeckID: "b1-art", Front: "das Meisterwerk", Back: "masterpiece", Extra: "Plural: die Meisterwerke", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-033", DeckID: "b1-art", Front: "der Stil", Back: "style", Extra: "Plural: die Stile", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-034", DeckID: "b1-art", Front: "die Fotografie", Back: "photography", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-035", DeckID: "b1-art", Front: "der Fotograf", Back: "photographer (male)", Extra: "Plural: die Fotografen", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-036", DeckID: "b1-art", Front: "das Foto", Back: "photo", Extra: "Plural: die Fotos", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-037", DeckID: "b1-art", Front: "das Motiv", Back: "motif, subject", Extra: "Plural: die Motive", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-038", DeckID: "b1-art", Front: "die Farbe", Back: "color, paint", Extra: "Plural: die Farben", Tags: []string{"b1", "art", "noun"}},
		{ID: "art-039", DeckID: "b1-art", Front: "bunt", Back: "colorful", Tags: []string{"b1", "art", "adjective"}},
		{ID: "art-040", DeckID: "b1-art", Front: "das Design", Back: "design", Extra: "Plural: die Designs", Tags: []string{"b1", "art", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-art",
		Name:        "B1 Art & Literature",
		Description: "Vocabulary for art, museums, books, and literature",
		Tags:        []string{"b1", "art", "literature"},
		Notes:       notes,
	}
}
