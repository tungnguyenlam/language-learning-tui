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
	{
		Title: "N-Deklination Nouns",
		Tip:   "Some masculine nouns add -n/-en in all cases except nominative singular: 'der Kunde, des Kunden'.",
	},
	{
		Title: "Two-way Prepositions",
		Tip:   "Use Akkusativ for movement and Dativ for location with in, an, auf, über, unter, vor, hinter, neben, zwischen.",
	},
	{
		Title: "Verb Position in Questions",
		Tip:   "Yes/no questions start with the finite verb, while W-questions place the verb second after the question word.",
	},
	{
		Title: "Da- Compounds",
		Tip:   "Use darauf, damit, daran for things and preposition+pronoun for people: 'Ich warte auf ihn', but 'Ich warte darauf'.",
	},
	{
		Title: "Infinitive with zu",
		Tip:   "Many verb pairs use zu + infinitive: 'Ich versuche, Deutsch zu sprechen.' Modal verbs do not use zu.",
	},
	{
		Title: "Participle Adjectives",
		Tip:   "Participles often behave like adjectives: 'ein spannender Film', 'die geschlossene Tür'.",
	},
	{
		Title: "Sentence Bracket",
		Tip:   "In main clauses, verb parts form a bracket: finite verb in position 2 and separable prefix/participle at the end.",
	},
	{
		Title: "Reported Speech with Konjunktiv I",
		Tip:   "Formal reported speech often uses Konjunktiv I: 'Er sagt, er sei krank.'",
	},
	{
		Title: "Negation with kein vs nicht",
		Tip:   "Use kein before nouns without definite article, and nicht for verbs, adjectives, adverbs, and definite nouns.",
	},
	{
		Title: "Temporal Connectors",
		Tip:   "Use nachdem for sequence and während for simultaneity; both send the conjugated verb to the clause end.",
	},
	{
		Title: "Plusquamperfekt",
		Tip:   "Used to express an action that happened before another past action: 'Ich hatte schon gegessen, als er ankam.'",
	},
	{
		Title: "Futur I",
		Tip:   "Formed with 'werden' + infinitive. Often used for assumptions about the present or future.",
	},
	{
		Title: "Double Infinitive in Perfect",
		Tip:   "Modal verbs in the perfect tense with another verb use a double infinitive: 'Ich habe das nicht machen können.'",
	},
	{
		Title: "Um... zu vs. Damit",
		Tip:   "Use 'um... zu' when the subject is the same in both clauses. Use 'damit' when the subjects are different.",
	},
	{
		Title: "Je... desto...",
		Tip:   "The more... the more...: 'Je' + subordinate clause (verb at end), 'desto' + main clause (verb immediately after).",
	},
	{
		Title: "Als vs. Wenn",
		Tip:   "Use 'als' for a single event in the past. Use 'wenn' for repeated events in the past, present, or future.",
	},
	{
		Title: "Prepositions with Genitive",
		Tip:   "Prepositions like 'wegen', 'während', 'trotz', and 'anstatt' require the Genitive case (though Dativ is used colloquially).",
	},
	{
		Title: "Indefinite Pronouns",
		Tip:   "Pronouns like 'jemand', 'niemand', 'manche', 'alle' change endings based on case: 'Ich kenne niemanden.'",
	},
	{
		Title: "Modal Particles",
		Tip:   "Words like 'doch', 'mal', 'ja', 'denn' add flavor but not meaning: 'Komm mal her!' (Just come here!).",
	},
	{
		Title: "N-Deklination Exceptions",
		Tip:   "Exceptions include 'das Herz' (des Herzens) and some Latin words ending in -or/-ent (der Student, des Studenten).",
	},
}

func GetDailyGrammarTip() GrammarTip {
	dayOfYear := time.Now().YearDay()
	return grammarTips[dayOfYear%len(grammarTips)]
}
