package content

import (
	"deutsch-tui/internal/core"
)

func A2DailyLifeDeck() core.Deck {
	notes := []core.Note{
		{ID: "a2-dly-aufstehen", DeckID: "a2-daily-life", Front: "aufstehen", Back: "to get up", Extra: "Ich stehe um 7 Uhr auf.", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-zähneputzen", DeckID: "a2-daily-life", Front: "sich die Zähne putzen", Back: "to brush teeth", Extra: "Ich putze mir die Zähne.", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-duschen", DeckID: "a2-daily-life", Front: "duschen", Back: "to shower", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-anziehen", DeckID: "a2-daily-life", Front: "sich anziehen", Back: "to get dressed", Extra: "Er zieht sich an.", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-frühstücken", DeckID: "a2-daily-life", Front: "frühstücken", Back: "to have breakfast", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-kommunikation", DeckID: "a2-daily-life", Front: "telefonieren", Back: "to make a phone call", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-post", DeckID: "a2-daily-life", Front: "die Post", Back: "mail / post office", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-einkaufen", DeckID: "a2-daily-life", Front: "einkaufen", Back: "to shop / do shopping", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-kochen", DeckID: "a2-daily-life", Front: "kochen", Back: "to cook", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-aufräumen", DeckID: "a2-daily-life", Front: "aufräumen", Back: "to tidy up", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-staubsaugen", DeckID: "a2-daily-life", Front: "staubsaugen", Back: "to vacuum", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-wäsche", DeckID: "a2-daily-life", Front: "die Wäsche waschen", Back: "to do laundry", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-bügeln", DeckID: "a2-daily-life", Front: "bügeln", Back: "to iron", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-schlafen", DeckID: "a2-daily-life", Front: "schlafen gehen", Back: "to go to bed", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-wecker", DeckID: "a2-daily-life", Front: "der Wecker", Back: "alarm clock", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-routine", DeckID: "a2-daily-life", Front: "die Routine", Back: "routine", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-morgen", DeckID: "a2-daily-life", Front: "der Morgen", Back: "morning", Extra: "Guten Morgen!", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-abend", DeckID: "a2-daily-life", Front: "der Abend", Back: "evening", Extra: "Guten Abend!", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-nachmittag", DeckID: "a2-daily-life", Front: "der Nachmittag", Back: "afternoon", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-spät", DeckID: "a2-daily-life", Front: "spät", Back: "late", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-früh", DeckID: "a2-daily-life", Front: "früh", Back: "early", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-pause", DeckID: "a2-daily-life", Front: "die Pause", Back: "break / pause", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-frei", DeckID: "a2-daily-life", Front: "frei haben", Back: "to have time off", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-haushalt", DeckID: "a2-daily-life", Front: "der Haushalt", Back: "household", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-ordnung", DeckID: "a2-daily-life", Front: "in Ordnung bringen", Back: "to put in order", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-entsorgen", DeckID: "a2-daily-life", Front: "entsorgen", Back: "to dispose of", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-recyceln", DeckID: "a2-daily-life", Front: "recyceln", Back: "to recycle", Tags: []string{"a2", "environment", "daily"}},
		{ID: "a2-dly-müll", DeckID: "a2-daily-life", Front: "der Müll", Back: "trash / garbage", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-kerzen", DeckID: "a2-daily-life", Front: "die Kerze", Back: "candle", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-lampe", DeckID: "a2-daily-life", Front: "die Lampe", Back: "lamp", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-kissen", DeckID: "a2-daily-life", Front: "das Kissen", Back: "cushion / pillow", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-decke", DeckID: "a2-daily-life", Front: "die Decke", Back: "blanket / cover", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-vorhänge", DeckID: "a2-daily-life", Front: "die Vorhänge", Back: "curtains", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-teppich", DeckID: "a2-daily-life", Front: "der Teppich", Back: "carpet / rug", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-regal", DeckID: "a2-daily-life", Front: "das Regal", Back: "shelf", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-schrank", DeckID: "a2-daily-life", Front: "der Schrank", Back: "cupboard / wardrobe", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-keller", DeckID: "a2-daily-life", Front: "der Keller", Back: "basement", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-dachboden", DeckID: "a2-daily-life", Front: "der Dachboden", Back: "attic", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-treppe", DeckID: "a2-daily-life", Front: "die Treppe", Back: "stairs", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-aufzug", DeckID: "a2-daily-life", Front: "der Aufzug", Back: "elevator", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-klingeln", DeckID: "a2-daily-life", Front: "klingeln", Back: "to ring (doorbell)", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-besuch", DeckID: "a2-daily-life", Front: "der Besuch", Back: "visit / visitor", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-gast", DeckID: "a2-daily-life", Front: "der Gast", Back: "guest", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-einladen", DeckID: "a2-daily-life", Front: "einladen", Back: "to invite", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-klingeln2", DeckID: "a2-daily-life", Front: "an der Tür klingeln", Back: "to ring the doorbell", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-schlüssel", DeckID: "a2-daily-life", Front: "der Schlüssel", Back: "key", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-tür", DeckID: "a2-daily-life", Front: "die Tür", Back: "door", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-fenster", DeckID: "a2-daily-life", Front: "das Fenster", Back: "window", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-klopfen", DeckID: "a2-daily-life", Front: "klopfen", Back: "to knock", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-öffnen", DeckID: "a2-daily-life", Front: "öffnen", Back: "to open", Tags: []string{"a2", "daily"}},
		{ID: "a2-dly-schließen", DeckID: "a2-daily-life", Front: "schließen", Back: "to close", Tags: []string{"a2", "daily"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a2-daily-life",
		Name:        "A2 Daily Life & Household",
		Description: "Vocabulary for daily routines, household chores, and home.",
		Tags:        []string{"german", "a2", "daily", "household"},
		Notes:       notes,
	}
}
