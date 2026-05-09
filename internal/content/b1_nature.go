package content

import (
	"deutsch-tui/internal/core"
)

func B1NatureDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-nat-umwelt", DeckID: "b1-nature", Front: "die Umwelt", Back: "environment", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-klimawandel", DeckID: "b1-nature", Front: "der Klimawandel", Back: "climate change", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-nachhaltigkeit", DeckID: "b1-nature", Front: "die Nachhaltigkeit", Back: "sustainability", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-erneuerbar", DeckID: "b1-nature", Front: "erneuerbare Energien", Back: "renewable energies", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-umweltschutz", DeckID: "b1-nature", Front: "der Umweltschutz", Back: "environmental protection", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-verschmutzung", DeckID: "b1-nature", Front: "die Verschmutzung", Back: "pollution", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-recycling", DeckID: "b1-nature", Front: "das Recycling", Back: "recycling", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-muelltrennung", DeckID: "b1-nature", Front: "die Mülltrennung", Back: "waste separation", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-oekologisch", DeckID: "b1-nature", Front: "ökologisch", Back: "ecological", Tags: []string{"b1", "nature"}},
		{ID: "b1-nat-artenschutz", DeckID: "b1-nature", Front: "der Artenschutz", Back: "species protection", Tags: []string{"b1", "nature"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-nature",
		Name:        "German B1 Nature & Environment",
		Description: "Essential vocabulary for discussing nature, climate, and sustainability.",
		Tags:        []string{"german", "b1", "nature"},
		Notes:       notes,
	}
}
