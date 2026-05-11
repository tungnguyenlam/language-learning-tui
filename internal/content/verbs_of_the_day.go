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
	{
		German: "sprechen", English: "to speak",
		Ich: "spreche", Du: "sprichst", ErSieEs: "spricht",
		Wir: "sprechen", Ihr: "sprecht", SieSie: "sprechen",
		Example: "Sprechen Sie Deutsch?",
	},
	{
		German: "lesen", English: "to read",
		Ich: "lese", Du: "liest", ErSieEs: "liest",
		Wir: "lesen", Ihr: "lest", SieSie: "lesen",
		Example: "Ich lese ein Buch.",
	},
	{
		German: "schreiben", English: "to write",
		Ich: "schreibe", Du: "schreibst", ErSieEs: "schreibt",
		Wir: "schreiben", Ihr: "schreibt", SieSie: "schreiben",
		Example: "Er schreibt einen Brief.",
	},
	{
		German: "arbeiten", English: "to work",
		Ich: "arbeite", Du: "arbeitest", ErSieEs: "arbeitet",
		Wir: "arbeiten", Ihr: "arbeitet", SieSie: "arbeiten",
		Example: "Wir arbeiten im Büro.",
	},
	{
		German: "wohnen", English: "to live/reside",
		Ich: "wohne", Du: "wohnst", ErSieEs: "wohnt",
		Wir: "wohnen", Ihr: "wohnt", SieSie: "wohnen",
		Example: "Ich wohne in Berlin.",
	},
	{
		German: "verstehen", English: "to understand",
		Ich: "verstehe", Du: "verstehst", ErSieEs: "versteht",
		Wir: "verstehen", Ihr: "versteht", SieSie: "verstehen",
		Example: "Ich verstehe die Frage.",
	},
	{
		German: "kochen", English: "to cook",
		Ich: "koche", Du: "kochst", ErSieEs: "kocht",
		Wir: "kochen", Ihr: "kocht", SieSie: "kochen",
		Example: "Sie kocht gern.",
	},
	{
		German: "fahren", English: "to drive/go",
		Ich: "fahre", Du: "fährst", ErSieEs: "fährt",
		Wir: "fahren", Ihr: "fahrt", SieSie: "fahren",
		Example: "Wir fahren nach München.",
	},
	{
		German: "kommen", English: "to come",
		Ich: "komme", Du: "kommst", ErSieEs: "kommt",
		Wir: "kommen", Ihr: "kommt", SieSie: "kommen",
		Example: "Wann kommst du?",
	},
	{
		German: "machen", English: "to do/make",
		Ich: "mache", Du: "machst", ErSieEs: "macht",
		Wir: "machen", Ihr: "macht", SieSie: "machen",
		Example: "Was machst du?",
	},
	{
		German: "wissen", English: "to know (fact)",
		Ich: "weiß", Du: "weißt", ErSieEs: "weiß",
		Wir: "wissen", Ihr: "wisst", SieSie: "wissen",
		Example: "Ich weiß die Antwort.",
	},
	{
		German: "kennen", English: "to know (person)",
		Ich: "kenne", Du: "kennst", ErSieEs: "kennt",
		Wir: "kennen", Ihr: "kennt", SieSie: "kennen",
		Example: "Ich kenne ihn gut.",
	},
	{
		German: "brauchen", English: "to need",
		Ich: "brauche", Du: "brauchst", ErSieEs: "braucht",
		Wir: "brauchen", Ihr: "braucht", SieSie: "brauchen",
		Example: "Ich brauche Hilfe.",
	},
	{
		German: "heißen", English: "to be called",
		Ich: "heiße", Du: "heißt", ErSieEs: "heißt",
		Wir: "heißen", Ihr: "heißt", SieSie: "heißen",
		Example: "Ich heiße Anna.",
	},
	{
		German: "finden", English: "to find",
		Ich: "finde", Du: "findest", ErSieEs: "findet",
		Wir: "finden", Ihr: "findet", SieSie: "finden",
		Example: "Ich finde das gut.",
	},
	{
		German: "bleiben", English: "to stay",
		Ich: "bleibe", Du: "bleibst", ErSieEs: "bleibt",
		Wir: "bleiben", Ihr: "bleibt", SieSie: "bleiben",
		Example: "Ich bleibe heute zu Hause.",
	},
	{
		German: "denken", English: "to think",
		Ich: "denke", Du: "denkst", ErSieEs: "denkt",
		Wir: "denken", Ihr: "denkt", SieSie: "denken",
		Example: "Was denkst du?",
	},
	{
		German: "glauben", English: "to believe",
		Ich: "glaube", Du: "glaubst", ErSieEs: "glaubt",
		Wir: "glauben", Ihr: "glaubt", SieSie: "glauben",
		Example: "Ich glaube dir.",
	},
	{
		German: "helfen", English: "to help",
		Ich: "helfe", Du: "hilfst", ErSieEs: "hilft",
		Wir: "helfen", Ihr: "helft", SieSie: "helfen",
		Example: "Kannst du mir helfen?",
	},
	{
		German: "lassen", English: "to let",
		Ich: "lasse", Du: "lässt", ErSieEs: "lässt",
		Wir: "lassen", Ihr: "lasst", SieSie: "lassen",
		Example: "Lass uns gehen!",
	},
	{
		German: "nehmen", English: "to take",
		Ich: "nehme", Du: "nimmst", ErSieEs: "nimmt",
		Wir: "nehmen", Ihr: "nehmt", SieSie: "nehmen",
		Example: "Ich nehme einen Kaffee.",
	},
	{
		German: "sagen", English: "to say",
		Ich: "sage", Du: "sagst", ErSieEs: "sagt",
		Wir: "sagen", Ihr: "sagt", SieSie: "sagen",
		Example: "Was sagst du?",
	},
	{
		German: "spielen", English: "to play",
		Ich: "spiele", Du: "spielst", ErSieEs: "spielt",
		Wir: "spielen", Ihr: "spielt", SieSie: "spielen",
		Example: "Wir spielen Fußball.",
	},
	{
		German: "stehen", English: "to stand",
		Ich: "stehe", Du: "stehst", ErSieEs: "steht",
		Wir: "stehen", Ihr: "steht", SieSie: "stehen",
		Example: "Das Auto steht hier.",
	},
	{
		German: "tun", English: "to do",
		Ich: "tue", Du: "tust", ErSieEs: "tut",
		Wir: "tun", Ihr: "tut", SieSie: "tun",
		Example: "Was kann ich für dich tun?",
	},
	{
		German: "verlieren", English: "to lose",
		Ich: "verliere", Du: "verlierst", ErSieEs: "verliert",
		Wir: "verlieren", Ihr: "verliert", SieSie: "verlieren",
		Example: "Ich verliere das Spiel.",
	},
	{
		German: "ziehen", English: "to pull/move",
		Ich: "ziehe", Du: "ziehst", ErSieEs: "zieht",
		Wir: "ziehen", Ihr: "zieht", SieSie: "ziehen",
		Example: "Wir ziehen nach Berlin.",
	},
	{
		German: "bringen", English: "to bring",
		Ich: "bringe", Du: "bringst", ErSieEs: "bringt",
		Wir: "bringen", Ihr: "bringt", SieSie: "bringen",
		Example: "Kannst du mir Wasser bringen?",
	},
	{
		German: "fallen", English: "to fall",
		Ich: "falle", Du: "fällst", ErSieEs: "fällt",
		Wir: "fallen", Ihr: "fallt", SieSie: "fallen",
		Example: "Der Apfel fällt.",
	},
	{
		German: "halten", English: "to hold/stop",
		Ich: "halte", Du: "hältst", ErSieEs: "hält",
		Wir: "halten", Ihr: "haltet", SieSie: "halten",
		Example: "Der Zug hält hier.",
	},
}

func GetVerbOfTheDay() DailyVerb {
	dayOfYear := time.Now().YearDay()
	return dailyVerbs[dayOfYear%len(dailyVerbs)]
}
