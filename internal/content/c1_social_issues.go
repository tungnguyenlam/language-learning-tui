package content

import (
	"deutsch-tui/internal/core"
)

func C1SocialIssuesDeck() core.Deck {
	notes := []core.Note{
		{ID: "soc-001", DeckID: "c1-social-issues", Front: "Die {{c1::Auseinandersetzung::debate/clash}} mit der Vergangenheit ist ein wesentlicher Bestandteil der deutschen Identität.", Extra: "Nomen: die Auseinandersetzung, Verb: sich auseinandersetzen mit.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-002", DeckID: "c1-social-issues", Front: "Der {{c1::demografische Wandel::demographic change}} führt zu einer alternden Gesellschaft und belastet die Sozialsysteme.", Tags: []string{"c1", "society", "demographics", "cloze"}},
		{ID: "soc-003", DeckID: "c1-social-issues", Front: "Viele Menschen engagieren sich {{c1::ehrenamtlich::on a voluntary basis}} in Vereinen oder sozialen Organisationen.", Extra: "das Ehrenamt = voluntary work.", Tags: []string{"c1", "society", "work", "cloze"}},
		{ID: "soc-004", DeckID: "c1-social-issues", Front: "Die {{c1::Chancengleichheit::equality of opportunity}} im Bildungssystem ist ein oft diskutiertes politisches Ziel.", Tags: []string{"c1", "society", "education", "cloze"}},
		{ID: "soc-005", DeckID: "c1-social-issues", Front: "Die {{c1::soziale Ungerechtigkeit::social injustice}} hat in den letzten Jahren in vielen Industrieländern zugenommen.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-006", DeckID: "c1-social-issues", Front: "Integration erfordert die {{c1::Bereitschaft::willingness}} beider Seiten, aufeinander zuzugehen.", Tags: []string{"c1", "society", "integration", "cloze"}},
		{ID: "soc-007", DeckID: "c1-social-issues", Front: "Der Schutz der {{c1::Privatsphäre::privacy}} wird im digitalen Zeitalter immer schwieriger.", Tags: []string{"c1", "society", "digital", "cloze"}},
		{ID: "soc-008", DeckID: "c1-social-issues", Front: "Die {{c1::Meinungsfreiheit::freedom of speech}} ist ein Grundpfeiler jeder demokratischen Gesellschaft.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-009", DeckID: "c1-social-issues", Front: "Es ist notwendig, Vorurteile durch Aufklärung und {{c1::Begegnung::encounter}} abzubauen.", Tags: []string{"c1", "society", "psychology", "cloze"}},
		{ID: "soc-010", DeckID: "c1-social-issues", Front: "Die {{c1::Nachhaltigkeit::sustainability}} spielt eine zentrale Rolle in der modernen Umweltpolitik.", Tags: []string{"c1", "society", "environment", "cloze"}},
		{ID: "soc-011", DeckID: "c1-social-issues", Front: "In einer {{c1::pluralistischen::pluralistic}} Gesellschaft müssen verschiedene Lebensentwürfe respektiert werden.", Tags: []string{"c1", "society", "philosophy", "cloze"}},
		{ID: "soc-012", DeckID: "c1-social-issues", Front: "Die {{c1::Zivilgesellschaft::civil society}} leistet einen wichtigen Beitrag zur Stabilität der Demokratie.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-013", DeckID: "c1-social-issues", Front: "Der {{c1::Fachkräftemangel::shortage of skilled workers}} bedroht das Wirtschaftswachstum in vielen Branchen.", Tags: []string{"c1", "society", "economy", "cloze"}},
		{ID: "soc-014", DeckID: "c1-social-issues", Front: "Es bedarf einer {{c1::umfassenden::comprehensive}} Reform des Gesundheitssystems.", Tags: []string{"c1", "society", "health", "cloze"}},
		{ID: "soc-015", DeckID: "c1-social-issues", Front: "Die {{c1::Vereinbarkeit::compatibility}} von Familie und Beruf ist für viele junge Eltern eine Herausforderung.", Tags: []string{"c1", "society", "family", "cloze"}},
		{ID: "soc-016", DeckID: "c1-social-issues", Front: "Die {{c1::Digitalisierung::digitalization}} verändert die Arbeitswelt von Grund auf.", Tags: []string{"c1", "society", "work", "cloze"}},
		{ID: "soc-017", DeckID: "c1-social-issues", Front: "Eine offene Gesellschaft muss auch {{c1::kontroversen::controversial}} Diskussionen standhalten.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-018", DeckID: "c1-social-issues", Front: "Die {{c1::Glaubwürdigkeit::credibility}} von Politikern ist ein wichtiges Gut in der Demokratie.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-019", DeckID: "c1-social-issues", Front: "Man muss die {{c1::Auswirkungen::consequences/effects}} globaler Krisen auf lokaler Ebene betrachten.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-020", DeckID: "c1-social-issues", Front: "Die {{c1::Solidarität::solidarity}} innerhalb der EU wird in Krisenzeiten oft auf die Probe gestellt.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-021", DeckID: "c1-social-issues", Front: "Wir müssen uns gegen jede Form von {{c1::Diskriminierung::discrimination}} zur Wehr setzen.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-022", DeckID: "c1-social-issues", Front: "Die {{c1::Transparenz::transparency}} staatlichen Handelns ist eine Voraussetzung für Vertrauen.", Tags: []string{"c1", "society", "politics", "cloze"}},
		{ID: "soc-023", DeckID: "c1-social-issues", Front: "Medienkompetenz hilft dabei, {{c1::Falschinformationen::false information}} im Internet zu erkennen.", Tags: []string{"c1", "society", "digital", "cloze"}},
		{ID: "soc-024", DeckID: "c1-social-issues", Front: "Der {{c1::Zusammenhalt::cohesion}} der Gesellschaft wird durch wachsende Polarisierung gefährdet.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-025", DeckID: "c1-social-issues", Front: "Staatliche {{c1::Subventionen::subsidies}} können Innovationen in Schlüsselindustrien fördern.", Tags: []string{"c1", "society", "economy", "cloze"}},
		{ID: "soc-026", DeckID: "c1-social-issues", Front: "Die {{c1::Umsetzung::implementation}} neuer Gesetze dauert oft länger als erwartet.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-027", DeckID: "c1-social-issues", Front: "Es ist wichtig, die {{c1::Bedürfnisse::needs}} marginalisierter Gruppen ernst zu nehmen.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-028", DeckID: "c1-social-issues", Front: "Kulturelle {{c1::Vielfalt::diversity}} bereichert das gesellschaftliche Leben.", Tags: []string{"c1", "society", "cloze"}},
		{ID: "soc-029", DeckID: "c1-social-issues", Front: "Die {{c1::Verantwortung::responsibility}} des Einzelnen für das Gemeinwohl sollte gestärkt werden.", Tags: []string{"c1", "society", "philosophy", "cloze"}},
		{ID: "soc-030", DeckID: "c1-social-issues", Front: "Die {{c1::Bürokratie::bureaucracy}} wird oft als Hindernis für wirtschaftliche Dynamik empfunden.", Tags: []string{"c1", "society", "cloze"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-social-issues",
		Name:        "C1 Social Issues & Society",
		Description: "Advanced German vocabulary for social, political, and societal discussions",
		Tags:        []string{"c1", "society", "politics", "cloze"},
		Notes:       notes,
	}
}
