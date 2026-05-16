package content

import (
	"deutsch-tui/internal/core"
)

func B2JobApplicationDeck() core.Deck {
	notes := []core.Note{
		{ID: "b2-job-app-stellenausschreibung", DeckID: "b2-job-application", Front: "die Stellenausschreibung", Back: "job advertisement", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-anforderungsprofil", DeckID: "b2-job-application", Front: "das Anforderungsprofil", Back: "requirements profile", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-anschreiben", DeckID: "b2-job-application", Front: "das Anschreiben", Back: "cover letter", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-bewerbungsunterlagen", DeckID: "b2-job-application", Front: "die Bewerbungsunterlagen", Back: "application documents", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-zeugnis", DeckID: "b2-job-application", Front: "das Zeugnis", Back: "certificate / reference", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-arbeitszeugnis", DeckID: "b2-job-application", Front: "das Arbeitszeugnis", Back: "employment reference", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-qualifikation", DeckID: "b2-job-application", Front: "die Qualifikation", Back: "qualification", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-soft-skills", DeckID: "b2-job-application", Front: "die Soft Skills (soziale Kompetenzen)", Back: "soft skills", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-fachkenntnisse", DeckID: "b2-job-application", Front: "die Fachkenntnisse", Back: "specialist knowledge", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-belastbarkeit", DeckID: "b2-job-application", Front: "die Belastbarkeit", Back: "resilience / ability to work under pressure", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-teamfaehigkeit", DeckID: "b2-job-application", Front: "die Teamfähigkeit", Back: "ability to work in a team", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-zuverlaessigkeit", DeckID: "b2-job-application", Front: "die Zuverlässigkeit", Back: "reliability", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-flexibilitaet", DeckID: "b2-job-application", Front: "die Flexibilität", Back: "flexibility", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-einsatzbereitschaft", DeckID: "b2-job-application", Front: "die Einsatzbereitschaft", Back: "commitment / readiness", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-motivation", DeckID: "b2-job-application", Front: "die Motivation", Back: "motivation", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-vorstellungsgespraech", DeckID: "b2-job-application", Front: "das Vorstellungsgespräch", Back: "job interview", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-selbstpraesentation", DeckID: "b2-job-application", Front: "die Selbstpräsentation", Back: "self-presentation", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-staerken", DeckID: "b2-job-application", Front: "die Stärken", Back: "strengths", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-schwaechen", DeckID: "b2-job-application", Front: "die Schwächen", Back: "weaknesses", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-gehaltsvorstellung", DeckID: "b2-job-application", Front: "die Gehaltsvorstellung", Back: "salary expectation", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-fruehestmoeglicher-eintrittstermin", DeckID: "b2-job-application", Front: "der frühestmögliche Eintrittstermin", Back: "earliest possible start date", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-probezeit", DeckID: "b2-job-application", Front: "die Probezeit", Back: "probationary period", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-befristet", DeckID: "b2-job-application", Front: "befristet", Back: "temporary / fixed-term", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-unbefristet", DeckID: "b2-job-application", Front: "unbefristet", Back: "permanent / open-ended", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-vollzeitstelle", DeckID: "b2-job-application", Front: "die Vollzeitstelle", Back: "full-time position", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-teilzeitstelle", DeckID: "b2-job-application", Front: "die Teilzeitstelle", Back: "part-time position", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-gleitzeit", DeckID: "b2-job-application", Front: "die Gleitzeit", Back: "flextime", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-homeoffice", DeckID: "b2-job-application", Front: "das Homeoffice", Back: "working from home / home office", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-ueberstundenregelung", DeckID: "b2-job-application", Front: "die Überstundenregelung", Back: "overtime policy", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-fortbildung", DeckID: "b2-job-application", Front: "die Fortbildung", Back: "further training / continuing education", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-karrieremoeglichkeiten", DeckID: "b2-job-application", Front: "die Karrieremöglichkeiten", Back: "career opportunities", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-aufstiegsschancen", DeckID: "b2-job-application", Front: "die Aufstiegsschancen", Back: "promotion prospects", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-bewerber", DeckID: "b2-job-application", Front: "der Bewerber / die Bewerberin", Back: "applicant", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-personalabteilung", DeckID: "b2-job-application", Front: "die Personalabteilung (HR)", Back: "human resources department", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-vorabtelefonat", DeckID: "b2-job-application", Front: "das Vorabtelefonat", Back: "preliminary telephone interview", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-assessment-center", DeckID: "b2-job-application", Front: "das Assessment-Center", Back: "assessment center", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-referenzen", DeckID: "b2-job-application", Front: "die Referenzen", Back: "references", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-absage", DeckID: "b2-job-application", Front: "die Absage", Back: "rejection", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-zusage", DeckID: "b2-job-application", Front: "die Zusage", Back: "acceptance / job offer", Tags: []string{"b2", "career"}},
		{ID: "b2-job-app-einarbeitung", DeckID: "b2-job-application", Front: "die Einarbeitung", Back: "induction / orientation", Tags: []string{"b2", "career"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-job-application",
		Name:        "B2 German Job Application & Interviews",
		Description: "Vocabulary for applying for jobs, writing CVs, and mastering interviews in German.",
		Tags:        []string{"german", "b2", "jobs", "career", "interview"},
		Notes:       notes,
	}
}
