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
	{
		Title: "Imperative Mood",
		Tip:   "For the 'du' form of imperative, usually just remove the -en and -st (Lerne! instead of du lernst).",
	},
	{
		Title: "Reflexive Verbs",
		Tip:   "Reflexive pronouns: mich, dich, sich, uns, euch, sich. Used when the subject and object are the same.",
	},
	{
		Title: "Separable Verbs",
		Tip:   "In simple sentences, the prefix of a separable verb (like 'aufstehen') moves to the very end.",
	},
	{
		Title: "Comparison of Adjectives",
		Tip:   "Positive: schnell; Comparative: schneller; Superlative: am schnellsten.",
	},
	{
		Title: "Konjunktiv II for Politeness",
		Tip:   "Use 'ich hätte gern...' or 'könnten Sie...' to sound polite in requests and service situations.",
	},
	{
		Title: "Passive Voice with werden",
		Tip:   "Form passive with 'werden' + Partizip II: 'Das Paket wird morgen geliefert.'",
	},
	{
		Title: "Noun-Verb Collocations",
		Tip:   "Learn fixed pairs like 'eine Entscheidung treffen' and 'in Betracht ziehen' for natural B1/B2 German.",
	},
	{
		Title: "Subordinate Clauses",
		Tip:   "In clauses with weil, dass, obwohl, and wenn, the conjugated verb goes to the end.",
	},
	{
		Title: "Preposition Choice with Time",
		Tip:   "Use 'seit' for ongoing states, 'ab' for start points, and 'bis' for end points.",
	},
	{
		Title: "Relative Clauses",
		Tip:   "Relative pronouns agree with gender/number of the noun, but case depends on their role in the clause.",
	},
	{
		Title: "Verb Prefix Meaning Shift",
		Tip:   "Prefixes change meaning strongly: 'stehen' (stand) vs. 'verstehen' (understand) vs. 'aufstehen' (get up).",
	},
	{
		Title: "Nominalization",
		Tip:   "German often turns verbs into nouns: 'lesen' -> 'das Lesen'; these nouns are always capitalized.",
	},
	{
		Title: "Adjective Endings After Articles",
		Tip:   "Adjective endings depend on article type (definite/indefinite/none) and case; memorize common tables.",
	},
}

func GetDailyGrammarTip() GrammarTip {
	dayOfYear := time.Now().YearDay()
	return grammarTips[dayOfYear%len(grammarTips)]
}
