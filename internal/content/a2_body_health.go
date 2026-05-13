package content

import (
	"deutsch-tui/internal/core"
)

func A2BodyHealthDeck() core.Deck {
	notes := []core.Note{
		{ID: "a2-bh-kopf", DeckID: "a2-body-health", Front: "der Kopf", Back: "head", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-haar", DeckID: "a2-body-health", Front: "das Haar", Back: "hair", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-gesicht", DeckID: "a2-body-health", Front: "das Gesicht", Back: "face", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-auge", DeckID: "a2-body-health", Front: "das Auge", Back: "eye", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-ohr", DeckID: "a2-body-health", Front: "das Ohr", Back: "ear", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-nase", DeckID: "a2-body-health", Front: "die Nase", Back: "nose", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-mund", DeckID: "a2-body-health", Front: "der Mund", Back: "mouth", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-zahn", DeckID: "a2-body-health", Front: "der Zahn", Back: "tooth", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-zunge", DeckID: "a2-body-health", Front: "die Zunge", Back: "tongue", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-hals", DeckID: "a2-body-health", Front: "der Hals", Back: "neck / throat", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-schulter", DeckID: "a2-body-health", Front: "die Schulter", Back: "shoulder", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-arm", DeckID: "a2-body-health", Front: "der Arm", Back: "arm", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-ellbogen", DeckID: "a2-body-health", Front: "der Ellbogen", Back: "elbow", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-hand", DeckID: "a2-body-health", Front: "die Hand", Back: "hand", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-finger", DeckID: "a2-body-health", Front: "der Finger", Back: "finger", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-daumen", DeckID: "a2-body-health", Front: "der Daumen", Back: "thumb", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-bauch", DeckID: "a2-body-health", Front: "der Bauch", Back: "belly / stomach", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-ruecken", DeckID: "a2-body-health", Front: "der Rücken", Back: "back", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-bein", DeckID: "a2-body-health", Front: "das Bein", Back: "leg", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-knie", DeckID: "a2-body-health", Front: "das Knie", Back: "knee", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-fuss", DeckID: "a2-body-health", Front: "der Fuß", Back: "foot", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-zeh", DeckID: "a2-body-health", Front: "der Zeh", Back: "toe", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-herz", DeckID: "a2-body-health", Front: "das Herz", Back: "heart", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-lunge", DeckID: "a2-body-health", Front: "die Lunge", Back: "lung", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-magen", DeckID: "a2-body-health", Front: "der Magen", Back: "stomach", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-leber", DeckID: "a2-body-health", Front: "die Leber", Back: "liver", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-niere", DeckID: "a2-body-health", Front: "die Niere", Back: "kidney", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-blut", DeckID: "a2-body-health", Front: "das Blut", Back: "blood", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-knochen", DeckID: "a2-body-health", Front: "der Knochen", Back: "bone", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-haut", DeckID: "a2-body-health", Front: "die Haut", Back: "skin", Tags: []string{"a2", "body"}},
		{ID: "a2-bh-arzt", DeckID: "a2-body-health", Front: "der Arzt", Back: "doctor", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-aerztin", DeckID: "a2-body-health", Front: "die Ärztin", Back: "doctor (f.)", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-krankenhaus", DeckID: "a2-body-health", Front: "das Krankenhaus", Back: "hospital", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-praxis", DeckID: "a2-body-health", Front: "die Praxis", Back: "doctor's practice", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-apotheke", DeckID: "a2-body-health", Front: "die Apotheke", Back: "pharmacy", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-medikament", DeckID: "a2-body-health", Front: "das Medikament", Back: "medicine", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-tablette", DeckID: "a2-body-health", Front: "die Tablette", Back: "tablet / pill", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-rezept", DeckID: "a2-body-health", Front: "das Rezept", Back: "prescription / recipe", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-krank", DeckID: "a2-body-health", Front: "krank", Back: "sick / ill", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-gesund", DeckID: "a2-body-health", Front: "gesund", Back: "healthy", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-fieber", DeckID: "a2-body-health", Front: "das Fieber", Back: "fever", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-husten", DeckID: "a2-body-health", Front: "der Husten", Back: "cough", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-schnupfen", DeckID: "a2-body-health", Front: "der Schnupfen", Back: "cold (runny nose)", Tags: []string{"a2", "health"}},
		{ID: "a2-bh-schmerz", DeckID: "a2-body-health", Front: "der Schmerz", Back: "pain", Tags: []string{"a2", "health"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-body-health",
		Name:        "A2 German Body & Health",
		Description: "Elementary vocabulary for body parts, organs, and basic medical situations.",
		Tags:        []string{"german", "a2", "body", "health"},
		Notes:       notes,
	}
}
