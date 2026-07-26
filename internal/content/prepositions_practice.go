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
		{"Das Bild hängt an {{...}} Wand.", "der", "an", "Location (Dative)"},
		{"Ich hänge das Bild an {{...}} Wand.", "die", "an", "Motion (Accusative)"},
		{"Das Auto steht hinter {{...}} Haus.", "dem", "hinter", "Location (Dative)"},
		{"Ich parke das Auto hinter {{...}} Haus.", "das", "hinter", "Motion (Accusative)"},
		{"Der Stuhl steht neben {{...}} Tisch.", "dem", "neben", "Location (Dative)"},
		{"Ich stelle den Stuhl neben {{...}} Tisch.", "den", "neben", "Motion (Accusative)"},
		{"Die Lampe hängt über {{...}} Sofa.", "dem", "über", "Location (Dative)"},
		{"Wir hängen die Lampe über {{...}} Sofa.", "das", "über", "Motion (Accusative)"},
		{"Die Kinder spielen vor {{...}} Haus.", "dem", "vor", "Location (Dative)"},
		{"Die Kinder laufen vor {{...}} Haus.", "das", "vor", "Motion (Accusative)"},
		{"Der Stift liegt zwischen {{...}} Büchern.", "den", "zwischen", "Location (Dative)"},
		{"Ich lege den Stift zwischen {{...}} Bücher.", "die", "zwischen", "Motion (Accusative)"},
		// Genitive prepositions
		{Sentence: "Wegen {{...}} Erkältung kann er nicht kommen.", Answer: "seiner", Preposition: "wegen", Context: "Genitive (masculine possessive)"},
		{Sentence: "Trotz {{...}} Regens gehen sie spazieren.", Answer: "des", Preposition: "trotz", Context: "Genitive (masculine/neuter)"},
		{Sentence: "Während {{...}} Mittagspause lese ich.", Answer: "der", Preposition: "während", Context: "Genitive (feminine)"},
		{Sentence: "Statt {{...}} Busses nehme ich die U-Bahn.", Answer: "des", Preposition: "statt", Context: "Genitive (masculine/neuter)"},
		{Sentence: "Innerhalb {{...}} Woche muss das erledigt sein.", Answer: "einer", Preposition: "innerhalb", Context: "Genitive (feminine)"},
		{Sentence: "Außerhalb {{...}} Stadt gibt es viel Natur.", Answer: "der", Preposition: "außerhalb", Context: "Genitive (feminine)"},
		// Temporal prepositions (Dative)
		{Sentence: "Ich lerne Deutsch seit {{...}} Jahr.", Answer: "einem", Preposition: "seit", Context: "Dative (neuter, ein-word)"},
		{Sentence: "Nach {{...}} Arbeit gehe ich ins Fitnessstudio.", Answer: "der", Preposition: "nach", Context: "Dative (feminine)"},
		{Sentence: "Vor {{...}} Essen sollst du Hände waschen.", Answer: "dem", Preposition: "vor", Context: "Dative (neuter, temporal)"},
		{Sentence: "Bis {{...}} nächsten Montag muss der Bericht fertig sein.", Answer: "zum", Preposition: "bis + zu", Context: "Dative (bis zum = bis zu dem)"},
		// Accusative-only prepositions
		{Sentence: "Das Geschenk ist {{...}} seine Mutter.", Answer: "für", Preposition: "für", Context: "Accusative (für + accusative)"},
		{Sentence: "Wir laufen um {{...}} Park.", Answer: "den", Preposition: "um", Context: "Accusative (um + accusative)"},
		{Sentence: "Er kämpft {{...}} die Wahrheit.", Answer: "für", Preposition: "für", Context: "Accusative (für + accusative)"},
		{Sentence: "Ohne {{...}} Schlüssel kann sie nicht rein.", Answer: "den", Preposition: "ohne", Context: "Accusative (ohne + accusative)"},
	}
}
