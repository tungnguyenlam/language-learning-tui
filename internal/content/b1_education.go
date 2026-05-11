package content

import (
	"deutsch-tui/internal/core"
)

func B1EducationDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-edu-bildung", DeckID: "b1-education", Front: "die Bildung", Back: "education", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-schule", DeckID: "b1-education", Front: "die Schule", Back: "school", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-schueler", DeckID: "b1-education", Front: "der Schüler / die Schülerin", Back: "pupil / student", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-lehrer", DeckID: "b1-education", Front: "der Lehrer / die Lehrerin", Back: "teacher", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-unterricht", DeckID: "b1-education", Front: "der Unterricht", Back: "lesson / class", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-fach", DeckID: "b1-education", Front: "das Fach", Back: "subject", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-klasse", DeckID: "b1-education", Front: "die Klasse", Back: "class / grade", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-klassenzimmer", DeckID: "b1-education", Front: "das Klassenzimmer", Back: "classroom", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-pause", DeckID: "b1-education", Front: "die Pause", Back: "break", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-hausaufgabe", DeckID: "b1-education", Front: "die Hausaufgabe", Back: "homework", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-pruefung", DeckID: "b1-education", Front: "die Prüfung", Back: "exam", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-test", DeckID: "b1-education", Front: "der Test", Back: "test", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-note", DeckID: "b1-education", Front: "die Note", Back: "grade / mark", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-zeugnis", DeckID: "b1-education", Front: "das Zeugnis", Back: "report card", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-abschluss", DeckID: "b1-education", Front: "der Abschluss", Back: "degree / graduation", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-abitur", DeckID: "b1-education", Front: "das Abitur", Back: "A-levels / high school diploma", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-studium", DeckID: "b1-education", Front: "das Studium", Back: "studies / university study", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-universitaet", DeckID: "b1-education", Front: "die Universität", Back: "university", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-hochschule", DeckID: "b1-education", Front: "die Hochschule", Back: "college / university", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-student", DeckID: "b1-education", Front: "der Student / die Studentin", Back: "university student", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-professor", DeckID: "b1-education", Front: "der Professor / die Professorin", Back: "professor", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-vorlesung", DeckID: "b1-education", Front: "die Vorlesung", Back: "lecture", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-seminar", DeckID: "b1-education", Front: "das Seminar", Back: "seminar", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-semester", DeckID: "b1-education", Front: "das Semester", Back: "semester / term", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-stipendium", DeckID: "b1-education", Front: "das Stipendium", Back: "scholarship", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-bibliothek", DeckID: "b1-education", Front: "die Bibliothek", Back: "library", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-wissen", DeckID: "b1-education", Front: "das Wissen", Back: "knowledge", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-kenntnisse", DeckID: "b1-education", Front: "die Kenntnisse (Pl.)", Back: "skills / knowledge", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-lernen", DeckID: "b1-education", Front: "lernen", Back: "to learn", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-studieren", DeckID: "b1-education", Front: "studieren", Back: "to study (at university)", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-lehren", DeckID: "b1-education", Front: "lehren", Back: "to teach", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-unterrichten", DeckID: "b1-education", Front: "unterrichten", Back: "to teach / give lessons", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-erklaeren", DeckID: "b1-education", Front: "erklären", Back: "to explain", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-verstehen", DeckID: "b1-education", Front: "verstehen", Back: "to understand", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-bestehen", DeckID: "b1-education", Front: "bestehen", Back: "to pass (an exam)", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-durchfallen", DeckID: "b1-education", Front: "durchfallen", Back: "to fail (an exam)", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-wiederholen", DeckID: "b1-education", Front: "wiederholen", Back: "to repeat", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-ueben", DeckID: "b1-education", Front: "üben", Back: "to practice", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-abschreiben", DeckID: "b1-education", Front: "abschreiben", Back: "to copy / cheat", Tags: []string{"b1", "education"}},
		{ID: "b1-edu-aufpassen", DeckID: "b1-education", Front: "aufpassen", Back: "to pay attention", Tags: []string{"b1", "education"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-education",
		Name:        "German B1 Education & Studying",
		Description: "Essential vocabulary for discussing school, university, and learning.",
		Tags:        []string{"german", "b1", "education"},
		Notes:       notes,
	}
}
