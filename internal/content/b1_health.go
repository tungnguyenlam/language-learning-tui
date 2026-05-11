package content

import (
	"deutsch-tui/internal/core"
)

func B1HealthDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-hlt-gesundheit", DeckID: "b1-health", Front: "die Gesundheit", Back: "health", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-krankheit", DeckID: "b1-health", Front: "die Krankheit", Back: "illness / disease", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-schmerzen", DeckID: "b1-health", Front: "die Schmerzen (Pl.)", Back: "pain", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-symptom", DeckID: "b1-health", Front: "das Symptom", Back: "symptom", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-behandlung", DeckID: "b1-health", Front: "die Behandlung", Back: "treatment", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-medikament", DeckID: "b1-health", Front: "das Medikament", Back: "medication", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-apotheke", DeckID: "b1-health", Front: "die Apotheke", Back: "pharmacy", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-rezept", DeckID: "b1-health", Front: "das Rezept", Back: "prescription", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-arzt", DeckID: "b1-health", Front: "der Arzt / die Ärztin", Back: "doctor", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-patient", DeckID: "b1-health", Front: "der Patient / die Patientin", Back: "patient", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-krankenhaus", DeckID: "b1-health", Front: "das Krankenhaus", Back: "hospital", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-notaufnahme", DeckID: "b1-health", Front: "die Notaufnahme", Back: "emergency room", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-operation", DeckID: "b1-health", Front: "die Operation", Back: "surgery / operation", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-untersuchung", DeckID: "b1-health", Front: "die Untersuchung", Back: "examination", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-impfung", DeckID: "b1-health", Front: "die Impfung", Back: "vaccination", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-blutdruck", DeckID: "b1-health", Front: "der Blutdruck", Back: "blood pressure", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-fieber", DeckID: "b1-health", Front: "das Fieber", Back: "fever", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-husten", DeckID: "b1-health", Front: "der Husten", Back: "cough", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-schnupfen", DeckID: "b1-health", Front: "der Schnupfen", Back: "runny nose / cold", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-verletzung", DeckID: "b1-health", Front: "die Verletzung", Back: "injury", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-wunde", DeckID: "b1-health", Front: "die Wunde", Back: "wound", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-pflaster", DeckID: "b1-health", Front: "das Pflaster", Back: "plaster / band-aid", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-verband", DeckID: "b1-health", Front: "der Verband", Back: "bandage", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-krankenkasse", DeckID: "b1-health", Front: "die Krankenkasse", Back: "health insurance (provider)", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-versichert", DeckID: "b1-health", Front: "versichert sein", Back: "to be insured", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-krankschreibung", DeckID: "b1-health", Front: "die Krankschreibung", Back: "sick note", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-heilen", DeckID: "b1-health", Front: "heilen", Back: "to heal / cure", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-pflegen", DeckID: "b1-health", Front: "pflegen", Back: "to care for / nurse", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-vorbeugen", DeckID: "b1-health", Front: "vorbeugen", Back: "to prevent", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-erkranken", DeckID: "b1-health", Front: "erkranken", Back: "to fall ill", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-uebelkeit", DeckID: "b1-health", Front: "die Übelkeit", Back: "nausea", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-schwindel", DeckID: "b1-health", Front: "der Schwindel", Back: "dizziness", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-allergie", DeckID: "b1-health", Front: "die Allergie", Back: "allergy", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-diaet", DeckID: "b1-health", Front: "die Diät", Back: "diet", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-ernaehrung", DeckID: "b1-health", Front: "die Ernährung", Back: "nutrition", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-uebergewicht", DeckID: "b1-health", Front: "das Übergewicht", Back: "overweight", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-abnehmen", DeckID: "b1-health", Front: "abnehmen", Back: "to lose weight", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-zunehmen", DeckID: "b1-health", Front: "zunehmen", Back: "to gain weight", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-schwanger", DeckID: "b1-health", Front: "schwanger", Back: "pregnant", Tags: []string{"b1", "health"}},
		{ID: "b1-hlt-geburt", DeckID: "b1-health", Front: "die Geburt", Back: "birth", Tags: []string{"b1", "health"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-health",
		Name:        "German B1 Health & Medicine",
		Description: "Essential vocabulary for discussing health, illness, and medical care.",
		Tags:        []string{"german", "b1", "health"},
		Notes:       notes,
	}
}
