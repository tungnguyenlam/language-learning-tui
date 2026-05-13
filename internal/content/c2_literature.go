package content

import (
	"deutsch-tui/internal/core"
)

func C2LiteratureDeck() core.Deck {
	notes := []core.Note{
		{ID: "c2-lit-literatur", DeckID: "c2-literature", Front: "die Literatur", Back: "literature", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-autor", DeckID: "c2-literature", Front: "der Autor / die Autorin", Back: "author", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-schriftsteller", DeckID: "c2-literature", Front: "der Schriftsteller / die Schriftstellerin", Back: "writer", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-roman", DeckID: "c2-literature", Front: "der Roman", Back: "novel", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-gedicht", DeckID: "c2-literature", Front: "das Gedicht", Back: "poem", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-poesie", DeckID: "c2-literature", Front: "die Poesie", Back: "poetry", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-erzaehlung", DeckID: "c2-literature", Front: "die Erzählung", Back: "narrative / story", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-novelle", DeckID: "c2-literature", Front: "die Novelle", Back: "novella", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-drama", DeckID: "c2-literature", Front: "das Drama", Back: "drama / play", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-theaterstueck", DeckID: "c2-literature", Front: "das Theaterstück", Back: "play", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-kapitel", DeckID: "c2-literature", Front: "das Kapitel", Back: "chapter", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-absatz", DeckID: "c2-literature", Front: "der Absatz", Back: "paragraph", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-protagonist", DeckID: "c2-literature", Front: "der Protagonist / die Protagonistin", Back: "protagonist", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-antagonist", DeckID: "c2-literature", Front: "der Antagonist / die Antagonistin", Back: "antagonist", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-handlung", DeckID: "c2-literature", Front: "die Handlung", Back: "plot / storyline", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-metapher", DeckID: "c2-literature", Front: "die Metapher", Back: "metaphor", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-symbol", DeckID: "c2-literature", Front: "das Symbol", Back: "symbol", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-stilmittel", DeckID: "c2-literature", Front: "das Stilmittel", Back: "stylistic device", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-epoche", DeckID: "c2-literature", Front: "die Epoche", Back: "epoch / era", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-rezension", DeckID: "c2-literature", Front: "die Rezension", Back: "review / critique", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-verlag", DeckID: "c2-literature", Front: "der Verlag", Back: "publishing house", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-veroeffentlichen", DeckID: "c2-literature", Front: "veröffentlichen", Back: "to publish", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-herausgeben", DeckID: "c2-literature", Front: "herausgeben", Back: "to publish / to edit", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-lesung", DeckID: "c2-literature", Front: "die Lesung", Back: "reading (event)", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-zitat", DeckID: "c2-literature", Front: "das Zitat", Back: "quote / quotation", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-zitieren", DeckID: "c2-literature", Front: "zitieren", Back: "to quote / to cite", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-meisterwerk", DeckID: "c2-literature", Front: "das Meisterwerk", Back: "masterpiece", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-belletristik", DeckID: "c2-literature", Front: "die Belletristik", Back: "fiction", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-sachbuch", DeckID: "c2-literature", Front: "das Sachbuch", Back: "non-fiction book", Tags: []string{"c2", "literature"}},
		{ID: "c2-lit-biografie", DeckID: "c2-literature", Front: "die Biografie", Back: "biography", Tags: []string{"c2", "literature"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c2-literature",
		Name:        "C2 German Literature",
		Description: "Advanced literary vocabulary for reading and analyzing literature.",
		Tags:        []string{"german", "c2", "literature", "reading"},
		Notes:       notes,
	}
}
