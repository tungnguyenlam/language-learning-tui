package content

import (
	"deutsch-tui/internal/core"
)

func PhrasalVerbsDeck() core.Deck {
	notes := []core.Note{
		{ID: "pv-aufstehen", DeckID: "phrasal-verbs", Front: "aufstehen", Back: "to get up / stand up", Tags: []string{"separable", "pv"}},
		{ID: "pv-anfangen", DeckID: "phrasal-verbs", Front: "anfangen", Back: "to begin / start", Tags: []string{"separable", "pv"}},
		{ID: "pv-ausgehen", DeckID: "phrasal-verbs", Front: "ausgehen", Back: "to go out", Tags: []string{"separable", "pv"}},
		{ID: "pv-einkaufen", DeckID: "phrasal-verbs", Front: "einkaufen", Back: "to shop", Tags: []string{"separable", "pv"}},
		{ID: "pv-mitkommen", DeckID: "phrasal-verbs", Front: "mitkommen", Back: "to come along", Tags: []string{"separable", "pv"}},
		{ID: "pv-vorbereiten", DeckID: "phrasal-verbs", Front: "vorbereiten", Back: "to prepare", Tags: []string{"separable", "pv"}},
		{ID: "pv-fernsehen", DeckID: "phrasal-verbs", Front: "fernsehen", Back: "to watch TV", Tags: []string{"separable", "pv"}},
		{ID: "pv-anrufen", DeckID: "phrasal-verbs", Front: "anrufen", Back: "to call (phone)", Tags: []string{"separable", "pv"}},
		{ID: "pv-verstehen", DeckID: "phrasal-verbs", Front: "verstehen", Back: "to understand", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-besuchen", DeckID: "phrasal-verbs", Front: "besuchen", Back: "to visit", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-erklären", DeckID: "phrasal-verbs", Front: "erklären", Back: "to explain", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-gefallen", DeckID: "phrasal-verbs", Front: "gefallen", Back: "to like (be pleasing to)", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-bezahlen", DeckID: "phrasal-verbs", Front: "bezahlen", Back: "to pay", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-empfehlen", DeckID: "phrasal-verbs", Front: "empfehlen", Back: "to recommend", Tags: []string{"inseparable", "pv"}},
		{ID: "pv-entdecken", DeckID: "phrasal-verbs", Front: "entdecken", Back: "to discover", Tags: []string{"inseparable", "pv"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "phrasal-verbs",
		Name:        "German Phrasal Verbs",
		Description: "Common separable and inseparable verbs.",
		Tags:        []string{"german", "verbs", "grammar"},
		Notes:       notes,
	}
}
