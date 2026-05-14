package content

import (
	"deutsch-tui/internal/core"
)

func B1WeatherDeck() core.Deck {
	notes := []core.Note{
		{ID: "b1-weather-sonnig", DeckID: "b1-weather-seasons", Front: "sonnig", Back: "sunny", Extra: "Das Wetter ist heute sonnig. (The weather is sunny today.)", Tags: []string{"b1", "weather", "a2"}},
		{ID: "b1-weather-regnerisch", DeckID: "b1-weather-seasons", Front: "regnerisch", Back: "rainy", Extra: "Es bleibt regnerisch. (It stays rainy.)", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-bewölkt", DeckID: "b1-weather-seasons", Front: "bewölkt", Back: "cloudy", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-schnee", DeckID: "b1-weather-seasons", Front: "der Schnee", Back: "snow", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-gewitter", DeckID: "b1-weather-seasons", Front: "das Gewitter", Back: "thunderstorm", Extra: "Plural: die Gewitter", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-nebel", DeckID: "b1-weather-seasons", Front: "der Nebel", Back: "fog", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-frost", DeckID: "b1-weather-seasons", Front: "der Frost", Back: "frost", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-windig", DeckID: "b1-weather-seasons", Front: "windig", Back: "windy", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-sturm", DeckID: "b1-weather-seasons", Front: "der Sturm", Back: "storm (gale)", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-eis", DeckID: "b1-weather-seasons", Front: "das Eis", Back: "ice", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-regen", DeckID: "b1-weather-seasons", Front: "der Regen", Back: "rain", Tags: []string{"a2", "weather", "noun"}},
		{ID: "b1-weather-wolke", DeckID: "b1-weather-seasons", Front: "die Wolke", Back: "cloud", Tags: []string{"a2", "weather", "noun"}},
		{ID: "b1-weather-temperatur", DeckID: "b1-weather-seasons", Front: "die Temperatur", Back: "temperature", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-niederschlag", DeckID: "b1-weather-seasons", Front: "der Niederschlag", Back: "precipitation", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-wetter", DeckID: "b1-weather-seasons", Front: "das Wetter", Back: "weather", Tags: []string{"a2", "weather", "noun"}},
		{ID: "b1-weather-prognose", DeckID: "b1-weather-seasons", Front: "die Wetterprognose", Back: "weather forecast", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-feucht", DeckID: "b1-weather-seasons", Front: "feucht", Back: "humid / moist", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-trocken", DeckID: "b1-weather-seasons", Front: "trocken", Back: "dry", Tags: []string{"b1", "weather"}},
		{ID: "b1-weather-hitze", DeckID: "b1-weather-seasons", Front: "die Hitze", Back: "heat", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-kälte", DeckID: "b1-weather-seasons", Front: "die Kälte", Back: "cold", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-blitz", DeckID: "b1-weather-seasons", Front: "der Blitz", Back: "lightning", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-donner", DeckID: "b1-weather-seasons", Front: "der Donner", Back: "thunder", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-sonnenaufgang", DeckID: "b1-weather-seasons", Front: "der Sonnenaufgang", Back: "sunrise", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-sonnenuntergang", DeckID: "b1-weather-seasons", Front: "der Sonnenuntergang", Back: "sunset", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-regenschirm", DeckID: "b1-weather-seasons", Front: "der Regenschirm", Back: "umbrella", Tags: []string{"a2", "weather", "noun"}},
		{ID: "b1-weather-sonnenbrille", DeckID: "b1-weather-seasons", Front: "die Sonnenbrille", Back: "sunglasses", Tags: []string{"a2", "weather", "noun"}},
		{ID: "b1-weather-sonnencreme", DeckID: "b1-weather-seasons", Front: "die Sonnencreme", Back: "sunscreen", Tags: []string{"b1", "weather", "noun"}},
		{ID: "b1-weather-jahreszeit", DeckID: "b1-weather-seasons", Front: "die Jahreszeit", Back: "season", Tags: []string{"a2", "season", "noun"}},
		{ID: "b1-weather-frühling", DeckID: "b1-weather-seasons", Front: "der Frühling", Back: "spring", Tags: []string{"a1", "season", "noun"}},
		{ID: "b1-weather-sommer", DeckID: "b1-weather-seasons", Front: "der Sommer", Back: "summer", Tags: []string{"a1", "season", "noun"}},
		{ID: "b1-weather-herbst", DeckID: "b1-weather-seasons", Front: "der Herbst", Back: "autumn / fall", Tags: []string{"a1", "season", "noun"}},
		{ID: "b1-weather-winter", DeckID: "b1-weather-seasons", Front: "der Winter", Back: "winter", Tags: []string{"a1", "season", "noun"}},
		{ID: "b1-weather-schneien", DeckID: "b1-weather-seasons", Front: "schneien", Back: "to snow", Tags: []string{"a2", "weather", "verb"}},
		{ID: "b1-weather-regnen", DeckID: "b1-weather-seasons", Front: "regnen", Back: "to rain", Tags: []string{"a1", "weather", "verb"}},
		{ID: "b1-weather-strahlen", DeckID: "b1-weather-seasons", Front: "die Sonne scheint", Back: "the sun is shining", Tags: []string{"b1", "weather", "expression"}},
		{ID: "b1-weather-es-gewittert", DeckID: "b1-weather-seasons", Front: "es gewittert", Back: "it's thundering", Tags: []string{"b1", "weather", "expression"}},
		{ID: "b1-weather-draußen", DeckID: "b1-weather-seasons", Front: "Wie ist das Wetter draußen?", Back: "What's the weather like outside?", Tags: []string{"b1", "weather", "phrase"}},
		{ID: "b1-weather-regenwetter", DeckID: "b1-weather-seasons", Front: "Es sieht nach Regen aus.", Back: "It looks like rain.", Tags: []string{"b1", "weather", "phrase"}},
		{ID: "b1-weather-schönwetter", DeckID: "b1-weather-seasons", Front: "Bei dem schönen Wetter...", Back: "In this nice weather...", Tags: []string{"b1", "weather", "phrase"}},
		{ID: "b1-weather-wetterbericht", DeckID: "b1-weather-seasons", Front: "der Wetterbericht", Back: "weather report", Tags: []string{"b1", "weather", "noun"}},
	}
	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}
	return core.Deck{
		ID:          "b1-weather-seasons",
		Name:        "German B1 Weather & Seasons",
		Description: "Comprehensive weather vocabulary, expressions, and seasonal terms for B1 learners.",
		Tags:        []string{"german", "b1", "weather", "seasons"},
		Notes:       notes,
	}
}
