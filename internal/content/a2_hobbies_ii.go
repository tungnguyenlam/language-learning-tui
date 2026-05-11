package content

import (
	"deutsch-tui/internal/core"
)

func A2HobbiesIIDeck() core.Deck {
	notes := []core.Note{
		{ID: "hob2-001", DeckID: "a2-hobbies-ii", Front: "die Freizeit", Back: "free time", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-002", DeckID: "a2-hobbies-ii", Front: "das Hobby", Back: "hobby", Extra: "Plural: die Hobbys", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-003", DeckID: "a2-hobbies-ii", Front: "sammeln", Back: "to collect", Extra: "Konjugation: ich sammle", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-004", DeckID: "a2-hobbies-ii", Front: "die Briefmarke", Back: "stamp", Extra: "Plural: die Briefmarken", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-005", DeckID: "a2-hobbies-ii", Front: "das Videospiel", Back: "video game", Extra: "Plural: die Videospiele", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-006", DeckID: "a2-hobbies-ii", Front: "spielen", Back: "to play", Extra: "Konjugation: ich spiele", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-007", DeckID: "a2-hobbies-ii", Front: "das Brettspiel", Back: "board game", Extra: "Plural: die Brettspiele", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-008", DeckID: "a2-hobbies-ii", Front: "die Karte", Back: "card", Extra: "Plural: die Karten", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-009", DeckID: "a2-hobbies-ii", Front: "basteln", Back: "to do crafts, tinker", Extra: "Konjugation: ich bastle", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-010", DeckID: "a2-hobbies-ii", Front: "nähen", Back: "to sew", Extra: "Konjugation: ich nähe", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-011", DeckID: "a2-hobbies-ii", Front: "stricken", Back: "to knit", Extra: "Konjugation: ich stricke", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-012", DeckID: "a2-hobbies-ii", Front: "angeln", Back: "to fish", Extra: "Konjugation: ich angle", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-013", DeckID: "a2-hobbies-ii", Front: "die Natur", Back: "nature", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-014", DeckID: "a2-hobbies-ii", Front: "wandern", Back: "to hike", Extra: "Konjugation: ich wandere", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-015", DeckID: "a2-hobbies-ii", Front: "der Wald", Back: "forest", Extra: "Plural: die Wälder", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-016", DeckID: "a2-hobbies-ii", Front: "das Zelt", Back: "tent", Extra: "Plural: die Zelte", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-017", DeckID: "a2-hobbies-ii", Front: "zelten", Back: "to camp", Extra: "Konjugation: ich zelte", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-018", DeckID: "a2-hobbies-ii", Front: "das Lagerfeuer", Back: "campfire", Extra: "Plural: die Lagerfeuer", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-019", DeckID: "a2-hobbies-ii", Front: "das Instrument", Back: "instrument", Extra: "Plural: die Instrumente", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-020", DeckID: "a2-hobbies-ii", Front: "die Gitarre", Back: "guitar", Extra: "Plural: die Gitarren", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-021", DeckID: "a2-hobbies-ii", Front: "das Klavier", Back: "piano", Extra: "Plural: die Klaviere", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-022", DeckID: "a2-hobbies-ii", Front: "üben", Back: "to practice", Extra: "Konjugation: ich übe", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-023", DeckID: "a2-hobbies-ii", Front: "singen", Back: "to sing", Extra: "Konjugation: ich singe", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-024", DeckID: "a2-hobbies-ii", Front: "das Lied", Back: "song", Extra: "Plural: die Lieder", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-025", DeckID: "a2-hobbies-ii", Front: "tanzen", Back: "to dance", Extra: "Konjugation: ich tanze", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-026", DeckID: "a2-hobbies-ii", Front: "der Verein", Back: "club, association", Extra: "Plural: die Vereine", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-027", DeckID: "a2-hobbies-ii", Front: "das Mitglied", Back: "member", Extra: "Plural: die Mitglieder", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-028", DeckID: "a2-hobbies-ii", Front: "treffen", Back: "to meet", Extra: "Konjugation: ich treffe, du triffst", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-029", DeckID: "a2-hobbies-ii", Front: "die Freunde", Back: "friends", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-030", DeckID: "a2-hobbies-ii", Front: "zusammen", Back: "together", Tags: []string{"a2", "hobbies", "adverb"}},
		{ID: "hob2-031", DeckID: "a2-hobbies-ii", Front: "der Spaß", Back: "fun", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-032", DeckID: "a2-hobbies-ii", Front: "Es macht Spaß", Back: "It is fun", Tags: []string{"a2", "hobbies", "phrase"}},
		{ID: "hob2-033", DeckID: "a2-hobbies-ii", Front: "die Party", Back: "party", Extra: "Plural: die Partys", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-034", DeckID: "a2-hobbies-ii", Front: "feiern", Back: "to celebrate", Extra: "Konjugation: ich feiere", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-035", DeckID: "a2-hobbies-ii", Front: "das Fest", Back: "festival, celebration", Extra: "Plural: die Feste", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-036", DeckID: "a2-hobbies-ii", Front: "einladen", Back: "to invite", Extra: "Konjugation: ich lade ein, du lädst ein", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-037", DeckID: "a2-hobbies-ii", Front: "die Einladung", Back: "invitation", Extra: "Plural: die Einladungen", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-038", DeckID: "a2-hobbies-ii", Front: "das Geschenk", Back: "gift, present", Extra: "Plural: die Geschenke", Tags: []string{"a2", "hobbies", "noun"}},
		{ID: "hob2-039", DeckID: "a2-hobbies-ii", Front: "schenken", Back: "to give (as a gift)", Extra: "Konjugation: ich schenke", Tags: []string{"a2", "hobbies", "verb"}},
		{ID: "hob2-040", DeckID: "a2-hobbies-ii", Front: "genießen", Back: "to enjoy", Extra: "Konjugation: ich genieße", Tags: []string{"a2", "hobbies", "verb"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-hobbies-ii",
		Name:        "A2 Hobbies & Free Time II",
		Description: "More vocabulary for hobbies, crafts, and social activities",
		Tags:        []string{"a2", "hobbies", "free-time"},
		Notes:       notes,
	}
}
