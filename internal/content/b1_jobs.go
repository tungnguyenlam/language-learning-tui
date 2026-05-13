package content

import (
	"deutsch-tui/internal/core"
)

func B1JobsDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-job-beruf", DeckID: "b1-jobs", Front: "der Beruf", Back: "profession / occupation", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-arbeitgeber", DeckID: "b1-jobs", Front: "der Arbeitgeber", Back: "employer", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-arbeitnehmer", DeckID: "b1-jobs", Front: "der Arbeitnehmer", Back: "employee", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-bewerbung", DeckID: "b1-jobs", Front: "die Bewerbung", Back: "application", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-lebenslauf", DeckID: "b1-jobs", Front: "der Lebenslauf", Back: "CV / resume", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-vorstellungsgespraech", DeckID: "b1-jobs", Front: "das Vorstellungsgespräch", Back: "job interview", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-gehalt", DeckID: "b1-jobs", Front: "das Gehalt", Back: "salary", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-lohn", DeckID: "b1-jobs", Front: "der Lohn", Back: "wage", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-vertrag", DeckID: "b1-jobs", Front: "der Arbeitsvertrag", Back: "employment contract", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-kuendigung", DeckID: "b1-jobs", Front: "die Kündigung", Back: "dismissal / resignation", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-feierabend", DeckID: "b1-jobs", Front: "der Feierabend", Back: "end of the working day", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-ueberstunde", DeckID: "b1-jobs", Front: "die Überstunde", Back: "overtime", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-kollege", DeckID: "b1-jobs", Front: "der Kollege / die Kollegin", Back: "colleague", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-chef", DeckID: "b1-jobs", Front: "der Chef / die Chefin", Back: "boss", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-buero", DeckID: "b1-jobs", Front: "das Büro", Back: "office", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-besprechung", DeckID: "b1-jobs", Front: "die Besprechung", Back: "meeting", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-erfahrung", DeckID: "b1-jobs", Front: "die Berufserfahrung", Back: "professional experience", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-ausbildung", DeckID: "b1-jobs", Front: "die Ausbildung", Back: "training / apprenticeship", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-praktikum", DeckID: "b1-jobs", Front: "das Praktikum", Back: "internship", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-teilzeit", DeckID: "b1-jobs", Front: "die Teilzeit", Back: "part-time", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-vollzeit", DeckID: "b1-jobs", Front: "die Vollzeit", Back: "full-time", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-selbststaendig", DeckID: "b1-jobs", Front: "selbstständig", Back: "self-employed / independent", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-arbeitslos", DeckID: "b1-jobs", Front: "arbeitslos", Back: "unemployed", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-rente", DeckID: "b1-jobs", Front: "die Rente", Back: "pension", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-steuer", DeckID: "b1-jobs", Front: "die Steuer", Back: "tax", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-team", DeckID: "b1-jobs", Front: "das Team", Back: "team", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-projekt", DeckID: "b1-jobs", Front: "das Projekt", Back: "project", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-ziel", DeckID: "b1-jobs", Front: "das Ziel", Back: "goal / target", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-erfolg", DeckID: "b1-jobs", Front: "der Erfolg", Back: "success", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-verantwortung", DeckID: "b1-jobs", Front: "die Verantwortung", Back: "responsibility", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-aufgaben", DeckID: "b1-jobs", Front: "die Aufgabe", Back: "task / assignment", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-karriere", DeckID: "b1-jobs", Front: "die Karriere", Back: "career", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-bewerben", DeckID: "b1-jobs", Front: "sich bewerben", Back: "to apply", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-einstellen", DeckID: "b1-jobs", Front: "einstellen", Back: "to hire", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-entlassen", DeckID: "b1-jobs", Front: "entlassen", Back: "to dismiss / to fire", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-kuendigen", DeckID: "b1-jobs", Front: "kündigen", Back: "to resign / to quit / to cancel", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-streiken", DeckID: "b1-jobs", Front: "streiken", Back: "to strike", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-verdienen", DeckID: "b1-jobs", Front: "verdienen", Back: "to earn", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-aufsteigen", DeckID: "b1-jobs", Front: "aufsteigen", Back: "to get promoted", Tags: []string{"b1", "jobs"}},
		{ID: "b1-job-zusammenarbeiten", DeckID: "b1-jobs", Front: "zusammenarbeiten", Back: "to collaborate / to work together", Tags: []string{"b1", "jobs"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b1-jobs",
		Name:        "B1 German Jobs & Professions",
		Description: "Vocabulary for the workplace, career, and professional environments.",
		Tags:        []string{"german", "b1", "jobs", "career", "work"},
		Notes:       notes,
	}
}
