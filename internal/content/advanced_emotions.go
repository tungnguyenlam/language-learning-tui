package content

import (
	"deutsch-tui/internal/core"
)

func AdvancedEmotionsDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-emo-sehnsucht", DeckID: "adv-emotions", Front: "die Sehnsucht", Back: "longing / yearning", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-geborgenheit", DeckID: "adv-emotions", Front: "die Geborgenheit", Back: "feeling of security / safety", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-enttauschung", DeckID: "adv-emotions", Front: "die Enttäuschung", Back: "disappointment", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-begeisterung", DeckID: "adv-emotions", Front: "die Begeisterung", Back: "enthusiasm", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-verzweiflung", DeckID: "adv-emotions", Front: "die Verzweiflung", Back: "despair", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-eifersucht", DeckID: "adv-emotions", Front: "die Eifersucht", Back: "jealousy", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-neid", DeckID: "adv-emotions", Front: "der Neid", Back: "envy", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-zuversicht", DeckID: "adv-emotions", Front: "die Zuversicht", Back: "confidence / optimism", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-schadenfreude", DeckID: "adv-emotions", Front: "die Schadenfreude", Back: "joy in others' misfortune", Tags: []string{"c1", "emotion"}},
		{ID: "c1-emo-gleichgultigkeit", DeckID: "adv-emotions", Front: "die Gleichgültigkeit", Back: "indifference", Tags: []string{"c1", "emotion"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "adv-emotions",
		Name:        "German Advanced Feelings",
		Description: "Nuanced vocabulary for emotions (B2-C1).",
		Tags:        []string{"german", "c1", "emotions"},
		Notes:       notes,
	}
}
