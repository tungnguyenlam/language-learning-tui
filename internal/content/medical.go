package content

import (
	"deutsch-tui/internal/core"
)

func MedicalGermanDeck() core.Deck {
	notes := []core.Note{
		{ID: "med-behandlung", DeckID: "medical-german", Front: "die Behandlung", Back: "treatment", Tags: []string{"medical", "noun"}},
		{ID: "med-untersuchung", DeckID: "medical-german", Front: "die Untersuchung", Back: "examination", Tags: []string{"medical", "noun"}},
		{ID: "med-rezept", DeckID: "medical-german", Front: "das Rezept", Back: "prescription", Tags: []string{"medical", "noun"}},
		{ID: "med-krankenhaus", DeckID: "medical-german", Front: "das Krankenhaus", Back: "hospital", Tags: []string{"medical", "noun"}},
		{ID: "med-arzt", DeckID: "medical-german", Front: "der Arzt / die Ärztin", Back: "doctor", Tags: []string{"medical", "noun"}},
		{ID: "med-schmerzen", DeckID: "medical-german", Front: "die Schmerzen", Back: "pain", Tags: []string{"medical", "noun"}},
		{ID: "med-symptom", DeckID: "medical-german", Front: "das Symptom", Back: "symptom", Tags: []string{"medical", "noun"}},
		{ID: "med-diagnose", DeckID: "medical-german", Front: "die Diagnose", Back: "diagnosis", Tags: []string{"medical", "noun"}},
		{ID: "med-patient", DeckID: "medical-german", Front: "der Patient / die Patientin", Back: "patient", Tags: []string{"medical", "noun"}},
		{ID: "med-medikament", DeckID: "medical-german", Front: "das Medikament", Back: "medication", Tags: []string{"medical", "noun"}},
		{ID: "med-impfung", DeckID: "medical-german", Front: "die Impfung", Back: "vaccination", Tags: []string{"medical", "noun"}},
		{ID: "med-operation", DeckID: "medical-german", Front: "die Operation", Back: "operation / surgery", Tags: []string{"medical", "noun"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "medical-german",
		Name:        "Medical German (B2/C1)",
		Description: "Essential vocabulary for medical professionals and patients.",
		Tags:        []string{"german", "medical", "advanced"},
		Notes:       notes,
	}
}
