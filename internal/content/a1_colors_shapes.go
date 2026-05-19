package content

import (
	"deutsch-tui/internal/core"
)

func A1ColorsShapesDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-col-rot", DeckID: "a1-colors-shapes", Front: "rot", Back: "red", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-blau", DeckID: "a1-colors-shapes", Front: "blau", Back: "blue", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-gruen", DeckID: "a1-colors-shapes", Front: "grün", Back: "green", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-gelb", DeckID: "a1-colors-shapes", Front: "gelb", Back: "yellow", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-schwarz", DeckID: "a1-colors-shapes", Front: "schwarz", Back: "black", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-weiss", DeckID: "a1-colors-shapes", Front: "weiß", Back: "white", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-grau", DeckID: "a1-colors-shapes", Front: "grau", Back: "grey", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-braun", DeckID: "a1-colors-shapes", Front: "braun", Back: "brown", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-rosa", DeckID: "a1-colors-shapes", Front: "rosa", Back: "pink", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-orange", DeckID: "a1-colors-shapes", Front: "orange", Back: "orange", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-lila", DeckID: "a1-colors-shapes", Front: "lila", Back: "purple", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-hell", DeckID: "a1-colors-shapes", Front: "hell", Back: "light / bright", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-dunkel", DeckID: "a1-colors-shapes", Front: "dunkel", Back: "dark", Tags: []string{"a1", "colors"}},
		{ID: "a1-col-bunt", DeckID: "a1-colors-shapes", Front: "bunt", Back: "colorful", Tags: []string{"a1", "colors"}},
		{ID: "a1-sha-kreis", DeckID: "a1-colors-shapes", Front: "der Kreis", Back: "circle", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-quadrat", DeckID: "a1-colors-shapes", Front: "das Quadrat", Back: "square", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-dreieck", DeckID: "a1-colors-shapes", Front: "das Dreieck", Back: "triangle", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-rechteck", DeckID: "a1-colors-shapes", Front: "das Rechteck", Back: "rectangle", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-stern", DeckID: "a1-colors-shapes", Front: "der Stern", Back: "star", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-rund", DeckID: "a1-colors-shapes", Front: "rund", Back: "round", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-eckig", DeckID: "a1-colors-shapes", Front: "eckig", Back: "angular / square-shaped", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-gerade", DeckID: "a1-colors-shapes", Front: "gerade", Back: "straight", Tags: []string{"a1", "shapes"}},
		{ID: "a1-sha-krumm", DeckID: "a1-colors-shapes", Front: "krumm", Back: "crooked", Tags: []string{"a1", "shapes"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-colors-shapes",
		Name:        "A1 Colors and Shapes",
		Description: "Essential colors and geometric shapes in German.",
		Tags:        []string{"german", "a1", "colors", "shapes"},
		Notes:       notes,
	}
}
