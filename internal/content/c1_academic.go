package content

import (
	"deutsch-tui/internal/core"
)

func C1AcademicDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-acd-dissertation", DeckID: "c1-academic", Front: "die Dissertation", Back: "dissertation / thesis", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-forschung", DeckID: "c1-academic", Front: "die Forschung", Back: "research", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-methodik", DeckID: "c1-academic", Front: "die Methodik", Back: "methodology", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-hypothese", DeckID: "c1-academic", Front: "die Hypothese", Back: "hypothesis", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-these", DeckID: "c1-academic", Front: "die These", Back: "thesis / argument", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-argumentation", DeckID: "c1-academic", Front: "die Argumentation", Back: "argumentation", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-folgerung", DeckID: "c1-academic", Front: "die Folgerung", Back: "conclusion / inference", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-korrelation", DeckID: "c1-academic", Front: "die Korrelation", Back: "correlation", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-kausalität", DeckID: "c1-academic", Front: "die Kausalität", Back: "causality", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-empirisch", DeckID: "c1-academic", Front: "empirisch", Back: "empirical", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-qualitativ", DeckID: "c1-academic", Front: "qualitativ", Back: "qualitative", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-quantitativ", DeckID: "c1-academic", Front: "quantitativ", Back: "quantitative", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-analyse", DeckID: "c1-academic", Front: "die Analyse", Back: "analysis", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-synthese", DeckID: "c1-academic", Front: "die Synthese", Back: "synthesis", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-expertise", DeckID: "c1-academic", Front: "die Expertise", Back: "expertise", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-fachgebiet", DeckID: "c1-academic", Front: "das Fachgebiet", Back: "field of study", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-interdisziplinär", DeckID: "c1-academic", Front: "interdisziplinär", Back: "interdisciplinary", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-spezialisierung", DeckID: "c1-academic", Front: "die Spezialisierung", Back: "specialization", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-theoretisch", DeckID: "c1-academic", Front: "theoretisch", Back: "theoretical", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-praktisch", DeckID: "c1-academic", Front: "praktisch", Back: "practical / applied", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-reflexion", DeckID: "c1-academic", Front: "die Reflexion", Back: "reflection", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-kritik", DeckID: "c1-academic", Front: "die Kritik", Back: "criticism / critique", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-review", DeckID: "c1-academic", Front: "das Review / die Begutachtung", Back: "peer review", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-publikation", DeckID: "c1-academic", Front: "die Publikation", Back: "publication", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-fachzeitschrift", DeckID: "c1-academic", Front: "die Fachzeitschrift", Back: "academic journal", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-konferenz", DeckID: "c1-academic", Front: "die Konferenz", Back: "conference", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-vortrag", DeckID: "c1-academic", Front: "der Vortrag", Back: "presentation / lecture", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-referat", DeckID: "c1-academic", Front: "das Referat", Back: "presentation / talk", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-protokoll", DeckID: "c1-academic", Front: "das Protokoll", Back: "minutes / protocol", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-seminar", DeckID: "c1-academic", Front: "das Seminar", Back: "seminar", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-vorlesung", DeckID: "c1-academic", Front: "die Vorlesung", Back: "lecture", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-kolleg", DeckID: "c1-academic", Front: "das Kolleg", Back: "faculty / collegiate", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-professor", DeckID: "c1-academic", Front: "der Professor / die Professorin", Back: "professor", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-dozent", DeckID: "c1-academic", Front: "der Dozent / die Dozentin", Back: "lecturer", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-assistent", DeckID: "c1-academic", Front: "der Assistent / die Assistentin", Back: "research assistant", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-kollege", DeckID: "c1-academic", Front: "der Kollege / die Kollegin", Back: "colleague (academic)", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-kohorte", DeckID: "c1-academic", Front: "die Kohorte", Back: "cohort", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-stichprobe", DeckID: "c1-academic", Front: "die Stichprobe", Back: "sample (research)", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-datenerhebung", DeckID: "c1-academic", Front: "die Datenerhebung", Back: "data collection", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-auswertung", DeckID: "c1-academic", Front: "die Auswertung", Back: "evaluation / analysis", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-statistik", DeckID: "c1-academic", Front: "die Statistik", Back: "statistics", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-diagramm", DeckID: "c1-academic", Front: "das Diagramm", Back: "diagram / chart", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-graph", DeckID: "c1-academic", Front: "die Grafik", Back: "graphic / chart", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-abbildung", DeckID: "c1-academic", Front: "die Abbildung", Back: "figure / illustration", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-literatur", DeckID: "c1-academic", Front: "die Literatur", Back: "literature", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-quellennachweis", DeckID: "c1-academic", Front: "der Quellennachweis", Back: "citation / source reference", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-zitat", DeckID: "c1-academic", Front: "das Zitat", Back: "quotation", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-plagiat", DeckID: "c1-academic", Front: "das Plagiat", Back: "plagiarism", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-fußnote", DeckID: "c1-academic", Front: "die Fußnote", Back: "footnote", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-literaturverzeichnis", DeckID: "c1-academic", Front: "das Literaturverzeichnis", Back: "bibliography", Tags: []string{"c1", "academic"}},
		{ID: "c1-acd-zugang", DeckID: "c1-academic", Front: "der Zugang", Back: "access (to resources)", Tags: []string{"c1", "academic"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-academic",
		Name:        "C1 Academic & Scientific German",
		Description: "Advanced academic vocabulary for research and scientific contexts.",
		Tags:        []string{"german", "c1", "academic", "scientific"},
		Notes:       notes,
	}
}
