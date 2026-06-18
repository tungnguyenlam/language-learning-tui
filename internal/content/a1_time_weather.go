package content

import (
	"deutsch-tui/internal/core"
)

func A1TimeWeatherDeck() core.Deck {
	notes := []core.Note{
		{ID: "a1-time-zeit", DeckID: "a1-time-weather", Front: "die Zeit", Back: "time", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-uhr", DeckID: "a1-time-weather", Front: "die Uhr", Back: "clock / watch / o'clock", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-stunde", DeckID: "a1-time-weather", Front: "die Stunde", Back: "hour", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-minute", DeckID: "a1-time-weather", Front: "die Minute", Back: "minute", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-sekunde", DeckID: "a1-time-weather", Front: "die Sekunde", Back: "second", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-tag", DeckID: "a1-time-weather", Front: "der Tag", Back: "day", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-woche", DeckID: "a1-time-weather", Front: "die Woche", Back: "week", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-monat", DeckID: "a1-time-weather", Front: "der Monat", Back: "month", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-jahr", DeckID: "a1-time-weather", Front: "das Jahr", Back: "year", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-morgen", DeckID: "a1-time-weather", Front: "der Morgen", Back: "morning", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-mittag", DeckID: "a1-time-weather", Front: "der Mittag", Back: "noon", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-abend", DeckID: "a1-time-weather", Front: "der Abend", Back: "evening", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-nacht", DeckID: "a1-time-weather", Front: "die Nacht", Back: "night", Tags: []string{"a1", "time", "noun"}},
		{ID: "a1-time-montag", DeckID: "a1-time-weather", Front: "Montag", Back: "Monday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-dienstag", DeckID: "a1-time-weather", Front: "Dienstag", Back: "Tuesday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-mittwoch", DeckID: "a1-time-weather", Front: "Mittwoch", Back: "Wednesday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-donnerstag", DeckID: "a1-time-weather", Front: "Donnerstag", Back: "Thursday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-freitag", DeckID: "a1-time-weather", Front: "Freitag", Back: "Friday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-samstag", DeckID: "a1-time-weather", Front: "Samstag", Back: "Saturday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-sonntag", DeckID: "a1-time-weather", Front: "Sonntag", Back: "Sunday", Tags: []string{"a1", "time", "day"}},
		{ID: "a1-time-gestern", DeckID: "a1-time-weather", Front: "gestern", Back: "yesterday", Tags: []string{"a1", "time", "adverb"}},
		{ID: "a1-time-heute", DeckID: "a1-time-weather", Front: "heute", Back: "today", Tags: []string{"a1", "time", "adverb"}},
		{ID: "a1-time-morgen-adv", DeckID: "a1-time-weather", Front: "morgen", Back: "tomorrow", Tags: []string{"a1", "time", "adverb"}},
		{ID: "a1-weather-wetter", DeckID: "a1-time-weather", Front: "das Wetter", Back: "weather", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-sonne", DeckID: "a1-time-weather", Front: "die Sonne", Back: "sun", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-regen", DeckID: "a1-time-weather", Front: "der Regen", Back: "rain", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-schnee", DeckID: "a1-time-weather", Front: "der Schnee", Back: "snow", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-wind", DeckID: "a1-time-weather", Front: "der Wind", Back: "wind", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-wolke", DeckID: "a1-time-weather", Front: "die Wolke", Back: "cloud", Tags: []string{"a1", "weather", "noun"}},
		{ID: "a1-weather-warm", DeckID: "a1-time-weather", Front: "warm", Back: "warm", Tags: []string{"a1", "weather", "adjective"}},
		{ID: "a1-weather-kalt", DeckID: "a1-time-weather", Front: "kalt", Back: "cold", Tags: []string{"a1", "weather", "adjective"}},
		{ID: "a1-weather-heiss", DeckID: "a1-time-weather", Front: "heiß", Back: "hot", Tags: []string{"a1", "weather", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "a1-time-weather",
		Name:        "A1 Time & Weather",
		Description: "Essential A1 vocabulary for days, time, and basic weather.",
		Tags:        []string{"german", "a1", "time", "weather"},
		Notes:       notes,
	}
}
