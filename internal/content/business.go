package content

import (
	"deutsch-tui/internal/core"
)

func BusinessDeck() core.Deck {
	notes := []core.Note{
		{ID: "biz-besprechung", DeckID: "business", Front: "die Besprechung", Back: "meeting / discussion", Tags: []string{"business", "b2"}},
		{ID: "biz-verhandlung", DeckID: "business", Front: "die Verhandlung", Back: "negotiation", Tags: []string{"business", "b2"}},
		{ID: "biz-vertrag", DeckID: "business", Front: "der Vertrag", Back: "contract", Tags: []string{"business", "b2"}},
		{ID: "biz-angebot", DeckID: "business", Front: "das Angebot", Back: "offer / quote", Tags: []string{"business", "b2"}},
		{ID: "biz-kundenbetreuung", DeckID: "business", Front: "die Kundenbetreuung", Back: "customer service / support", Tags: []string{"business", "b2"}},
		{ID: "biz-umsatz", DeckID: "business", Front: "der Umsatz", Back: "turnover / revenue", Tags: []string{"business", "b2"}},
		{ID: "biz-bewerbung", DeckID: "business", Front: "die Bewerbung", Back: "application", Tags: []string{"business", "b1"}},
		{ID: "biz-vorstellungsgespraech", DeckID: "business", Front: "das Vorstellungsgespräch", Back: "job interview", Tags: []string{"business", "b1"}},
		{ID: "biz-geschaeftsfuehrer", DeckID: "business", Front: "der Geschäftsführer", Back: "managing director / CEO", Tags: []string{"business", "b2"}},
		{ID: "biz-buchhaltung", DeckID: "business", Front: "die Buchhaltung", Back: "accounting", Tags: []string{"business", "b2"}},
		{ID: "biz-personalabteilung", DeckID: "business", Front: "die Personalabteilung", Back: "HR department", Tags: []string{"business", "b2"}},
		{ID: "biz-marketing", DeckID: "business", Front: "das Marketing", Back: "marketing", Tags: []string{"business", "b2"}},
		{ID: "biz-vertrieb", DeckID: "business", Front: "der Vertrieb", Back: "sales / distribution", Tags: []string{"business", "b2"}},
		{ID: "biz-nische", DeckID: "business", Front: "die Marktnische", Back: "market niche", Tags: []string{"business", "c1"}},
		{ID: "biz-wettbewerb", DeckID: "business", Front: "der Wettbewerb", Back: "competition", Tags: []string{"business", "b2"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "business",
		Name:        "German Business Vocabulary",
		Description: "Essential terms for the German workplace, meetings, and negotiations.",
		Tags:        []string{"german", "business", "b1", "b2"},
		Notes:       notes,
	}
}
