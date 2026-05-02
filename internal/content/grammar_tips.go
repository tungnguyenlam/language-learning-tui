package content

import (
	"time"
)

type GrammarTip struct {
	Title string
	Tip   string
}

var grammarTips = []GrammarTip{
	{
		Title: "Articles and Genders",
		Tip:   "Nouns ending in -ung, -heit, -keit, -schaft, -ei, -in are usually feminine (die).",
	},
	{
		Title: "The Dativ Case",
		Tip:   "Prepositions 'aus, bei, mit, nach, seit, von, zu' always take the Dativ case.",
	},
	{
		Title: "The Akkusativ Case",
		Tip:   "Prepositions 'durch, für, gegen, ohne, um' always take the Akkusativ case.",
	},
	{
		Title: "Word Order",
		Tip:   "In a main clause, the conjugated verb always takes the second position.",
	},
	{
		Title: "Plural Forms",
		Tip:   "Nouns ending in -e often add -n to form the plural (die Lampe -> die Lampen).",
	},
	{
		Title: "Weak Verbs",
		Tip:   "Regular (weak) verbs follow a predictable pattern: -e, -st, -t, -en, -t, -en.",
	},
	{
		Title: "Perfect Tense",
		Tip:   "Most verbs use 'haben' as an auxiliary verb, but verbs of movement or change of state use 'sein'.",
	},
}

func GetDailyGrammarTip() GrammarTip {
	dayOfYear := time.Now().YearDay()
	return grammarTips[dayOfYear%len(grammarTips)]
}
