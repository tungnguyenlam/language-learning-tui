package content

import (
	"deutsch-tui/internal/core"
)

func A2MedicalAppointmentDeck() core.Deck {
	notes := []core.Note{
		{ID: "a2-ma-termin", DeckID: "a2-medical-appointment", Front: "einen Termin vereinbaren", Back: "to make an appointment", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-verschieben", DeckID: "a2-medical-appointment", Front: "einen Termin verschieben", Back: "to postpone/reschedule an appointment", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-absagen", DeckID: "a2-medical-appointment", Front: "einen Termin absagen", Back: "to cancel an appointment", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-praxis", DeckID: "a2-medical-appointment", Front: "die Arztpraxis", Back: "doctor's office/practice", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-empfang", DeckID: "a2-medical-appointment", Front: "der Empfang", Back: "reception", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-wartezimmer", DeckID: "a2-medical-appointment", Front: "das Wartezimmer", Back: "waiting room", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-versicherung", DeckID: "a2-medical-appointment", Front: "die Versicherungskarte", Back: "insurance card", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-untersuchung", DeckID: "a2-medical-appointment", Front: "die Untersuchung", Back: "examination", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-untersuchen", DeckID: "a2-medical-appointment", Front: "untersuchen", Back: "to examine", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-wehtun", DeckID: "a2-medical-appointment", Front: "wehtun", Back: "to hurt / to ache", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-beschwerden", DeckID: "a2-medical-appointment", Front: "die Beschwerden", Back: "complaints / symptoms", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-symptome", DeckID: "a2-medical-appointment", Front: "die Symptome", Back: "symptoms", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-fieber", DeckID: "a2-medical-appointment", Front: "Fieber messen", Back: "to take temperature / measure fever", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-blutdruck", DeckID: "a2-medical-appointment", Front: "den Blutdruck messen", Back: "to measure blood pressure", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-einatmen", DeckID: "a2-medical-appointment", Front: "tief einatmen", Back: "to breathe in deeply", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-ausatmen", DeckID: "a2-medical-appointment", Front: "ausatmen", Back: "to breathe out", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-rezept", DeckID: "a2-medical-appointment", Front: "ein Rezept ausstellen", Back: "to issue a prescription", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-krankschreibung", DeckID: "a2-medical-appointment", Front: "die Krankschreibung", Back: "sick note", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-attest", DeckID: "a2-medical-appointment", Front: "das ärztliche Attest", Back: "medical certificate", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-ueberweisung", DeckID: "a2-medical-appointment", Front: "die Überweisung", Back: "referral", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-facharzt", DeckID: "a2-medical-appointment", Front: "der Facharzt", Back: "specialist (doctor)", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-apotheke", DeckID: "a2-medical-appointment", Front: "in die Apotheke gehen", Back: "to go to the pharmacy", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-einnehmen", DeckID: "a2-medical-appointment", Front: "Medikamente einnehmen", Back: "to take medication", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-nuechtern", DeckID: "a2-medical-appointment", Front: "nüchtern sein", Back: "to have an empty stomach (fasting)", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-impfung", DeckID: "a2-medical-appointment", Front: "die Impfung", Back: "vaccination", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-spritze", DeckID: "a2-medical-appointment", Front: "die Spritze", Back: "injection / syringe", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-verband", DeckID: "a2-medical-appointment", Front: "der Verband", Back: "bandage / dressing", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-pflaster", DeckID: "a2-medical-appointment", Front: "das Pflaster", Back: "plaster / band-aid", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-salbe", DeckID: "a2-medical-appointment", Front: "die Salbe", Back: "ointment / cream", Tags: []string{"a2", "health"}},
		{ID: "a2-ma-tropfen", DeckID: "a2-medical-appointment", Front: "die Tropfen", Back: "drops (e.g. eye drops)", Tags: []string{"a2", "health"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-medical-appointment",
		Name:        "A2 Medical Appointment",
		Description: "Essential vocabulary and phrases for visiting a doctor, making appointments, and discussing health issues.",
		Tags:        []string{"german", "a2", "health", "medical"},
		Notes:       notes,
	}
}
