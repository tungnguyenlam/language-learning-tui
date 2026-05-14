package content

import (
	"deutsch-tui/internal/core"
)

func C2LegalDeck() core.Deck {
	notes := []core.Note{
		{ID: "c2-lgx-gesetz", DeckID: "c2-legal", Front: "das Gesetz", Back: "law / statute", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-recht", DeckID: "c2-legal", Front: "das Recht", Back: "law / legal right", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-gericht", DeckID: "c2-legal", Front: "das Gericht", Back: "court", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-richter", DeckID: "c2-legal", Front: "der Richter / die Richterin", Back: "judge", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-anwalt", DeckID: "c2-legal", Front: "der Anwalt / die Anwältin", Back: "lawyer", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-staatsanwalt", DeckID: "c2-legal", Front: "der Staatsanwalt / die Staatsanwältin", Back: "public prosecutor", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-verfahren", DeckID: "c2-legal", Front: "das Verfahren", Back: "procedure / proceedings", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-verhandlung", DeckID: "c2-legal", Front: "die Verhandlung", Back: "trial / hearing", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-urteil", DeckID: "c2-legal", Front: "das Urteil", Back: "verdict / judgment", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-entscheidung", DeckID: "c2-legal", Front: "die Entscheidung", Back: "decision / ruling", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-richterspruch", DeckID: "c2-legal", Front: "der Richterspruch", Back: "judicial ruling / judgment", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-beweis", DeckID: "c2-legal", Front: "der Beweis", Back: "evidence / proof", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-beweislast", DeckID: "c2-legal", Front: "die Beweislast", Back: "burden of proof", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-anklage", DeckID: "c2-legal", Front: "die Anklage", Back: "indictment / charge", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-angeklagt", DeckID: "c2-legal", Front: "der Angeklagte / die Angeklagte", Back: "defendant", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-kläger", DeckID: "c2-legal", Front: "der Kläger / die Klägerin", Back: "plaintiff", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-beklagter", DeckID: "c2-legal", Front: "der Beklagte / die Beklagte", Back: "defendant (civil)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-parteien", DeckID: "c2-legal", Front: "die Parteien", Back: "parties (legal)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-parteivertreter", DeckID: "c2-legal", Front: "der Parteivertreter", Back: "party representative", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-vertreter", DeckID: "c2-legal", Front: "der Vertreter / die Vertreterin", Back: "representative", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-vollmacht", DeckID: "c2-legal", Front: "die Vollmacht", Back: "power of attorney", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtsanwalt", DeckID: "c2-legal", Front: "der Rechtsanwalt / die Rechtsanwältin", Back: "attorney", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-justiz", DeckID: "c2-legal", Front: "die Justiz", Back: "judiciary", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtsprechung", DeckID: "c2-legal", Front: "die Rechtsprechung", Back: "jurisprudence", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-gesetzgebung", DeckID: "c2-legal", Front: "die Gesetzgebung", Back: "legislation", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtsnorm", DeckID: "c2-legal", Front: "die Rechtsnorm", Back: "legal norm", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-paragraph", DeckID: "c2-legal", Front: "der Paragraph / §", Back: "paragraph / section", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-artikel", DeckID: "c2-legal", Front: "der Artikel", Back: "article (law)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-gesetzentwurf", DeckID: "c2-legal", Front: "der Gesetzentwurf", Back: "bill (legislation)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-verfassungs", DeckID: "c2-legal", Front: "verfassungsmäßig", Back: "constitutional", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-grundgesetz", DeckID: "c2-legal", Front: "das Grundgesetz (GG)", Back: "Basic Law (German constitution)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtlich", DeckID: "c2-legal", Front: "rechtlich", Back: "legal / lawful", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtsgültig", DeckID: "c2-legal", Front: "rechtsgültig", Back: "legally valid", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtskräftig", DeckID: "c2-legal", Front: "rechtskräftig", Back: "final / legally binding", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-strafbar", DeckID: "c2-legal", Front: "strafbar", Back: "punishable / criminal", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-verboten", DeckID: "c2-legal", Front: "verboten", Back: "prohibited / forbidden", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-erlaubt", DeckID: "c2-legal", Front: "erlaubt", Back: "permitted / allowed", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-pflicht", DeckID: "c2-legal", Front: "die Pflicht", Back: "obligation / duty", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-haftung", DeckID: "c2-legal", Front: "die Haftung", Back: "liability", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-schadenersatz", DeckID: "c2-legal", Front: "der Schadenersatz", Back: "compensation / damages", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-klage", DeckID: "c2-legal", Front: "die Klage", Back: "lawsuit", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-klagen", DeckID: "c2-legal", Front: "klagen", Back: "to sue / to litigate", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-verklagen", DeckID: "c2-legal", Front: "verklagen", Back: "to sue (someone)", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-strafverfahren", DeckID: "c2-legal", Front: "das Strafverfahren", Back: "criminal proceedings", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-zivilverfahren", DeckID: "c2-legal", Front: "das Zivilverfahren", Back: "civil proceedings", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-berufung", DeckID: "c2-legal", Front: "die Berufung", Back: "appeal", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-revision", DeckID: "c2-legal", Front: "die Revision", Back: "review / appeal", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-justiziabel", DeckID: "c2-legal", Front: "justiziabel", Back: "justiciable", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-juristisch", DeckID: "c2-legal", Front: "juristisch", Back: "juristic / legal", Tags: []string{"c2", "legal"}},
		{ID: "c2-lgx-rechtsstaat", DeckID: "c2-legal", Front: "der Rechtsstaat", Back: "rule of law", Tags: []string{"c2", "legal"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c2-legal",
		Name:        "C2 German Legal & Juridical",
		Description: "Advanced legal vocabulary for juridical and legislative contexts.",
		Tags:        []string{"german", "c2", "legal", "juridical"},
		Notes:       notes,
	}
}
