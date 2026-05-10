package content

import (
	"time"
)

type DailyVerb struct {
	German  string
	English string
	Ich     string
	Du      string
	ErSieEs string
	Wir     string
	Ihr     string
	SieSie  string
	Example string
}

var dailyVerbs = []DailyVerb{
	{
		German: "lernen", English: "to learn",
		Ich: "lerne", Du: "lernst", ErSieEs: "lernt",
		Wir: "lernen", Ihr: "lernt", SieSie: "lernen",
		Example: "Ich lerne jeden Tag Deutsch.",
	},
	{
		German: "essen", English: "to eat",
		Ich: "esse", Du: "isst", ErSieEs: "isst",
		Wir: "essen", Ihr: "esst", SieSie: "essen",
		Example: "Wir essen heute im Restaurant.",
	},
	{
		German: "trinken", English: "to drink",
		Ich: "trinke", Du: "trinkst", ErSieEs: "trinkt",
		Wir: "trinken", Ihr: "trinkt", SieSie: "trinken",
		Example: "Trinkst du gern Kaffee?",
	},
	{
		German: "gehen", English: "to go",
		Ich: "gehe", Du: "gehst", ErSieEs: "geht",
		Wir: "gehen", Ihr: "geht", SieSie: "gehen",
		Example: "Wann gehen wir nach Hause?",
	},
	{
		German: "sehen", English: "to see",
		Ich: "sehe", Du: "siehst", ErSieEs: "sieht",
		Wir: "sehen", Ihr: "seht", SieSie: "sehen",
		Example: "Ich sehe den Film.",
	},
}

func GetVerbOfTheDay() DailyVerb {
	dayOfYear := time.Now().YearDay()
	return dailyVerbs[dayOfYear%len(dailyVerbs)]
}
