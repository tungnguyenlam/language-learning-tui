package content

import (
	"deutsch-tui/internal/core"
)

func B2BusinessMeetingsDeck() core.Deck {
	notes := []core.Note{
		{ID: "biz-001", DeckID: "b2-business-meetings", Front: "die Tagesordnung", Back: "agenda", Extra: "Plural: die Tagesordnungen. Beispiel: Was steht heute auf der Tagesordnung?", Tags: []string{"b2", "business", "meeting"}},
		{ID: "biz-002", DeckID: "b2-business-meetings", Front: "der Wortbeitrag", Back: "contribution (to a discussion)", Extra: "Plural: die Wortbeiträge", Tags: []string{"b2", "business", "meeting"}},
		{ID: "biz-003", DeckID: "b2-business-meetings", Front: "zustimmen", Back: "to agree", Extra: "Konjugation: ich stimme zu. + Dativ", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-004", DeckID: "b2-business-meetings", Front: "ablehnen", Back: "to reject, to decline", Extra: "Beispiel: Einen Vorschlag ablehnen.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-005", DeckID: "b2-business-meetings", Front: "unterbrechen", Back: "to interrupt", Extra: "Konjugation: du unterbrichst, er unterbricht", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-006", DeckID: "b2-business-meetings", Front: "ausreden lassen", Back: "to let someone finish speaking", Extra: "Lassen Sie mich bitte ausreden.", Tags: []string{"b2", "business", "phrase"}},
		{ID: "biz-007", DeckID: "b2-business-meetings", Front: "das Protokoll", Back: "minutes (of a meeting)", Extra: "Plural: die Protokolle. Protokoll führen.", Tags: []string{"b2", "business", "meeting"}},
		{ID: "biz-008", DeckID: "b2-business-meetings", Front: "der Vorsitzende", Back: "chairman, chairperson", Extra: "Feminin: die Vorsitzende", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-009", DeckID: "b2-business-meetings", Front: "beschlussfähig", Back: "quorate (having a quorum)", Extra: "Adjektiv. Wir sind heute nicht beschlussfähig.", Tags: []string{"b2", "business", "meeting"}},
		{ID: "biz-010", DeckID: "b2-business-meetings", Front: "abstimmen", Back: "to vote", Extra: "Über einen Antrag abstimmen.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-011", DeckID: "b2-business-meetings", Front: "das Ergebnis", Back: "result, outcome", Extra: "Plural: die Ergebnisse", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-012", DeckID: "b2-business-meetings", Front: "zusammenfassen", Back: "to summarize", Extra: "Darf ich die wichtigsten Punkte zusammenfassen?", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-013", DeckID: "b2-business-meetings", Front: "der Kompromiss", Back: "compromise", Extra: "Einen Kompromiss schließen / eingehen.", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-014", DeckID: "b2-business-meetings", Front: "verhandeln", Back: "to negotiate", Extra: "Nomen: die Verhandlung", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-015", DeckID: "b2-business-meetings", Front: "das Gegenargument", Back: "counter-argument", Extra: "Plural: die Gegenargumente", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-016", DeckID: "b2-business-meetings", Front: "überzeugend", Back: "convincing", Extra: "Ein überzeugendes Argument.", Tags: []string{"b2", "business", "adjective"}},
		{ID: "biz-017", DeckID: "b2-business-meetings", Front: "vermitteln", Back: "to mediate, to convey", Extra: "Zwischen den Parteien vermitteln.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-018", DeckID: "b2-business-meetings", Front: "der Standpunkt", Back: "point of view", Extra: "Plural: die Standpunkte", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-019", DeckID: "b2-business-meetings", Front: "berücksichtigen", Back: "to consider, to take into account", Extra: "Alle Faktoren berücksichtigen.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-020", DeckID: "b2-business-meetings", Front: "hervorheben", Back: "to emphasize, to highlight", Extra: "Konjugation: ich hebe hervor", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-021", DeckID: "b2-business-meetings", Front: "das Feedback", Back: "feedback", Extra: "Konstruktives Feedback geben.", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-022", DeckID: "b2-business-meetings", Front: "die Rückfrage", Back: "follow-up question", Extra: "Plural: die Rückfragen", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-023", DeckID: "b2-business-meetings", Front: "ansprechen", Back: "to address (a topic)", Extra: "Ein Problem offen ansprechen.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-024", DeckID: "b2-business-meetings", Front: "verschieben", Back: "to postpone", Extra: "Einen Termin auf nächste Woche verschieben.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-025", DeckID: "b2-business-meetings", Front: "absagen", Back: "to cancel", Extra: "Ein Meeting kurzfristig absagen.", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-026", DeckID: "b2-business-meetings", Front: "die Einladung", Back: "invitation", Extra: "Plural: die Einladungen", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-027", DeckID: "b2-business-meetings", Front: "teilnehmen", Back: "to participate, to attend", Extra: "An einer Sitzung teilnehmen. + an + Dativ", Tags: []string{"b2", "business", "verb"}},
		{ID: "biz-028", DeckID: "b2-business-meetings", Front: "der Teilnehmer", Back: "participant", Extra: "Plural: die Teilnehmer", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-029", DeckID: "b2-business-meetings", Front: "der Diskussionsbedarf", Back: "need for discussion", Extra: "Hier besteht noch Diskussionsbedarf.", Tags: []string{"b2", "business", "noun"}},
		{ID: "biz-030", DeckID: "b2-business-meetings", Front: "zielorientiert", Back: "goal-oriented", Extra: "Wir müssen zielorientiert arbeiten.", Tags: []string{"b2", "business", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "b2-business-meetings",
		Name:        "B2 Business: Meetings & Negotiations",
		Description: "Professional German for meetings, discussions, and negotiations",
		Tags:        []string{"b2", "business", "professional", "career"},
		Notes:       notes,
	}
}
