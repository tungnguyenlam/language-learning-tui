package content

type PrepositionExercise struct {
	Sentence    string
	Answer      string
	Preposition string
	Context     string
}

func GetPrepositionExercises() []PrepositionExercise {
	return []PrepositionExercise{
		{"Wir gehen {{...}} den Park.", "durch", "durch", "Accusative"},
		{"Das Geschenk ist {{...}} dich.", "für", "für", "Accusative"},
		{"Sie kommt {{...}} der Schweiz.", "aus", "aus", "Dative"},
		{"Ich fahre {{...}} dem Zug.", "mit", "mit", "Dative"},
		{"Das Buch liegt auf {{...}} Tisch.", "dem", "auf", "Location (Dative)"},
		{"Ich lege das Buch auf {{...}} Tisch.", "den", "auf", "Motion (Accusative)"},
		{"Wir treffen uns {{...}} 8 Uhr.", "um", "um", "Time (Accusative)"},
		{"Der Hund schläft unter {{...}} Bett.", "dem", "unter", "Location (Dative)"},
		{"Der Hund kriecht unter {{...}} Bett.", "das", "unter", "Motion (Accusative)"},
		{"Wir wohnen in {{...}} Stadt.", "der", "in", "Location (Dative)"},
		{"Wir fahren in {{...}} Stadt.", "die", "in", "Motion (Accusative)"},
	}
}
