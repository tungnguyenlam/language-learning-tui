package content

import (
	"time"
)

type GrammarTip struct {
	Title   string
	Tip     string
	Example string
}

var grammarTips = []GrammarTip{
	{
		Title:   "Articles and Genders",
		Tip:     "Nouns ending in -ung, -heit, -keit, -schaft, -ei, -in are usually feminine (die).",
		Example: "die Erfahrung, die Freiheit, die Möglichkeit",
	},
	{
		Title:   "The Dativ Case",
		Tip:     "Prepositions 'aus, bei, mit, nach, seit, von, zu' always take the Dativ case.",
		Example: "Ich gehe mit meinem Freund (Dativ m.).",
	},
	{
		Title:   "The Akkusativ Case",
		Tip:     "Prepositions 'durch, für, gegen, ohne, um' always take the Akkusativ case.",
		Example: "Das Geschenk ist für dich (Akkusativ).",
	},
	{
		Title:   "Word Order",
		Tip:     "In a main clause, the conjugated verb always takes the second position.",
		Example: "Ich esse (1) heute (2) einen Apfel.",
	},
	{
		Title:   "Plural Forms",
		Tip:     "Nouns ending in -e often add -n to form the plural.",
		Example: "die Lampe -> die Lampen",
	},
	{
		Title:   "Weak Verbs",
		Tip:     "Regular (weak) verbs follow a predictable pattern: -e, -st, -t, -en, -t, -en.",
		Example: "lernen: ich lerne, du lernst, er lernt",
	},
	{
		Title:   "Perfect Tense",
		Tip:     "Most verbs use 'haben' as an auxiliary, but verbs of movement/change use 'sein'.",
		Example: "Ich habe gegessen. Ich bin gegangen.",
	},
	{
		Title:   "Imperative Mood",
		Tip:     "For the 'du' form of imperative, usually just remove the -en and -st.",
		Example: "Lerne! (instead of du lernst)",
	},
	{
		Title:   "Reflexive Verbs",
		Tip:     "Reflexive pronouns: mich, dich, sich, uns, euch, sich.",
		Example: "Ich wasche mich.",
	},
	{
		Title:   "Separable Verbs",
		Tip:     "The prefix of a separable verb moves to the very end in simple sentences.",
		Example: "Ich stehe um 7 Uhr auf. (aufstehen)",
	},
	{
		Title:   "Comparison of Adjectives",
		Tip:     "Positive: schnell; Comparative: schneller; Superlative: am schnellsten.",
		Example: "Berlin ist groß, aber London ist größer.",
	},
	{
		Title:   "Konjunktiv II for Politeness",
		Tip:     "Use 'ich hätte gern...' or 'könnten Sie...' to sound polite.",
		Example: "Ich hätte gern einen Kaffee, bitte.",
	},
	{
		Title:   "Passive Voice with werden",
		Tip:     "Form passive with 'werden' + Partizip II.",
		Example: "Das Haus wird gebaut.",
	},
	{
		Title:   "Noun-Verb Collocations",
		Tip:     "Learn fixed pairs like 'eine Entscheidung treffen'.",
		Example: "Wir müssen endlich eine Entscheidung treffen.",
	},
	{
		Title:   "Subordinate Clauses",
		Tip:     "In clauses with weil, dass, obwohl, the conjugated verb goes to the end.",
		Example: "Ich lerne, weil ich in Berlin wohnen will.",
	},
	{
		Title:   "Preposition Choice with Time",
		Tip:     "Use 'seit' for ongoing states, 'ab' for future start points.",
		Example: "Seit zwei Jahren lerne ich Deutsch.",
	},
	{
		Title:   "Relative Clauses",
		Tip:     "Relative pronouns agree with gender/number, but case depends on its role.",
		Example: "Das ist der Mann, den ich gestern sah.",
	},
	{
		Title:   "Verb Prefix Meaning Shift",
		Tip:     "Prefixes change meaning: 'stehen' (stand) vs. 'verstehen' (understand).",
		Example: "Ich verstehe die Frage nicht.",
	},
	{
		Title:   "Nominalization",
		Tip:     "German often turns verbs into nouns (always capitalized, always neuter).",
		Example: "Das Lesen von Büchern macht Spaß.",
	},
	{
		Title:   "Adjective Endings After Articles",
		Tip:     "Endings depend on article type (definite/indefinite) and case.",
		Example: "der gute Kaffee, ein guter Kaffee",
	},
	{
		Title:   "N-Deklination Nouns",
		Tip:     "Some masculine nouns add -n/-en in all cases except nominative singular.",
		Example: "Ich sehe den Jungen. (der Junge)",
	},
	{
		Title:   "Two-way Prepositions",
		Tip:     "Akkusativ for movement, Dativ for location.",
		Example: "Ich hänge das Bild an die Wand (Akk). Es hängt an der Wand (Dat).",
	},
	{
		Title:   "Verb Position in Questions",
		Tip:     "Yes/no questions start with the verb. W-questions have verb second.",
		Example: "Kommst du? Woher kommst du?",
	},
	{
		Title:   "Da- Compounds",
		Tip:     "Use darauf, damit for things; preposition+pronoun for people.",
		Example: "Ich warte darauf. Ich warte auf ihn.",
	},
	{
		Title:   "Infinitive with zu",
		Tip:     "Many verb pairs use zu + infinitive. Modal verbs do NOT.",
		Example: "Ich versuche zu schlafen. Ich kann schlafen.",
	},
	{
		Title:   "Participle Adjectives",
		Tip:     "Participles often behave like adjectives.",
		Example: "die geschlossene Tür, der schreiende Junge",
	},
	{
		Title:   "Sentence Bracket",
		Tip:     "Finite verb in pos 2 and separable prefix/participle at the end.",
		Example: "Ich habe gestern ein Buch gekauft.",
	},
	{
		Title:   "Reported Speech with Konjunktiv I",
		Tip:     "Formal reported speech often uses Konjunktiv I.",
		Example: "Er sagt, er sei müde.",
	},
	{
		Title:   "Negation with kein vs nicht",
		Tip:     "Use kein before nouns without article, nicht for others.",
		Example: "Ich habe kein Geld. Ich schlafe nicht.",
	},
	{
		Title:   "Temporal Connectors",
		Tip:     "Use nachdem for sequence and während for simultaneity.",
		Example: "Nachdem ich gegessen hatte, ging ich spazieren.",
	},
	{
		Title:   "Plusquamperfekt",
		Tip:     "Action that happened before another past action.",
		Example: "Er war schon weg, als ich ankam.",
	},
	{
		Title:   "Futur I",
		Tip:     "Formed with 'werden' + infinitive. Often expresses assumptions.",
		Example: "Es wird wohl regnen.",
	},
	{
		Title:   "Double Infinitive in Perfect",
		Tip:     "Modal verbs with another verb use a double infinitive.",
		Example: "Ich habe es nicht machen können.",
	},
	{
		Title:   "Um... zu vs. Damit",
		Tip:     "Use 'um... zu' when subjects are the same.",
		Example: "Ich lerne, um den Test zu bestehen.",
	},
	{
		Title:   "Je... desto...",
		Tip:     "The more... the more...",
		Example: "Je mehr ich lerne, desto besser werde ich.",
	},
	{
		Title:   "Als vs. Wenn",
		Tip:     "Use 'als' for a single past event. Use 'wenn' for repeats.",
		Example: "Als ich ein Kind war... Wenn ich Zeit habe...",
	},
	{
		Title:   "Prepositions with Genitive",
		Tip:     "Wegen, während, trotz, and anstatt require Genitive.",
		Example: "Wegen des Regens bleiben wir zu Hause.",
	},
	{
		Title:   "Indefinite Pronouns",
		Tip:     "Pronouns like 'jemand' change endings based on case.",
		Example: "Ich sehe jemanden.",
	},
	{
		Title:   "Modal Particles",
		Tip:     "Words like 'doch', 'mal', 'ja' add flavor/nuance.",
		Example: "Komm mal her! Das ist ja toll!",
	},
	{
		Title:   "N-Deklination Exceptions",
		Tip:     "Exceptions include 'das Herz' (des Herzens).",
		Example: "Ich danke dir von Herzen.",
	},
	{
		Title:   "Konjunktiv II for Unreal Conditions",
		Tip:     "Formed with 'würde' + infinitive or 'wäre', 'hätte'.",
		Example: "Wenn ich reich wäre, würde ich reisen.",
	},
	{
		Title:   "Futur II for Assumptions",
		Tip:     "Assumption about a past action.",
		Example: "Er wird wohl angekommen sein.",
	},
	{
		Title:   "Participle I as Adjective",
		Tip:     "Present participle (-end) describes ongoing action.",
		Example: "das schlafende Baby",
	},
	{
		Title:   "Extended Adjective Phrases",
		Tip:     "German can pack complex info before a noun.",
		Example: "die gestern gekaufte Zeitung",
	},
	{
		Title:   "Alternative Passive with 'lassen'",
		Tip:     "Use 'sich lassen' + infinitive.",
		Example: "Das lässt sich machen.",
	},
	{
		Title:   "Nominalized Adjectives",
		Tip:     "Adjectives as nouns follow adjective declension.",
		Example: "der Alte, ein Alter",
	},
	{
		Title:   "Double Connectors",
		Tip:     "Sowohl... als auch, weder... noch.",
		Example: "Ich mag sowohl Tee als auch Kaffee.",
	},
	{
		Title:   "Modal Verbs with Subjective Meaning",
		Tip:     "Modals can express probability/claims.",
		Example: "Er will es nicht gewusst haben.",
	},
	{
		Title:   "Genitive Prepositions (Advanced)",
		Tip:     "Use 'unweit', 'jenseits' for formal writing.",
		Example: "unweit des Flusses",
	},
	{
		Title:   "Konjunktiv I for Neutral Distance",
		Tip:     "Used in journalism to report statements neutrally.",
		Example: "Der Minister sagte, die Lage sei ernst.",
	},
	{
		Title:   "Nominalstil (Nominal Style)",
		Tip:     "Academic German uses many nouns instead of verbs.",
		Example: "Die Untersuchung ergab... (instead of 'Wir untersuchten...')",
	},
	{
		Title:   "Adverbial Genitive",
		Tip:     "Fixed expressions use genitive for adverbs.",
		Example: "meines Erachtens (in my view)",
	},
	{
		Title:   "Wissen vs Kennen",
		Tip:     "Wissen for facts, Kennen for familiarity/people.",
		Example: "Ich weiß, wo er ist. Ich kenne ihn.",
	},
	{
		Title:   "Position of 'nicht'",
		Tip:     "Usually after the verb, but before what it negates.",
		Example: "Ich esse das nicht. Ich bin nicht müde.",
	},
	{
		Title:   "The Suffix -in",
		Tip:     "Add -in to make professions feminine.",
		Example: "der Arzt -> die Ärztin",
	},
	{
		Title:   "Verb-Preposition Fixed Pairs",
		Tip:     "Many verbs require specific prepositions.",
		Example: "warten auf (+Akk), denken an (+Akk)",
	},
	{
		Title:   "Relative Clauses with 'der/die/das'",
		Tip:     "Use the same gender/number as the noun they describe.",
		Example: "Der Mann, der dort steht, ist mein Lehrer.",
	},
	{
		Title:   "Welch- vs Was für ein",
		Tip:     "welch- asks about specific selection; 'was für ein' asks about type.",
		Example: "Welches Buch? Was für ein Buch?",
	},
	{
		Title:   "Umgangssprachliche Ausdrücke",
		Tip:     "Common informal phrases: 'Alles klar?' (How's it going?), 'Kein Problem!' (No problem!)",
		Example: "Alles klar? Ja, alles gut!",
	},
	{
		Title:   "Wechselpräpositionen",
		Tip:     "Some prepositions can take either Dativ or Akkusativ: 'in, an, auf, über, unter, vor, hinter, neben, zwischen'",
		Example: "Ich bin in dem Haus (wo?). -> Ich gehe in das Haus (wohin?).",
	},
	{
		Title:   "Participial Phrases",
		Tip:     "Partizip I (-nd) for ongoing actions; Partizip II for completed.",
		Example: "Der singende Mann (aktive Handlung). Das gekochte Essen.",
	},
	{
		Title:   "Modal Particles",
		Tip:     "Particles like 'doch', 'mal', 'schon' add nuance.",
		Example: "Kommst du mit? -> Kommst du doch mit? (encouragement)",
	},
	{
		Title:   "Depersonal Phrases",
		Tip:     "Start sentences with 'Es' or 'Man' for impersonal expressions.",
		Example: "Es regnet. Man sagt, dass...",
	},
	{
		Title:   "Numbers & Time",
		Tip:     "Use 'halb' for half hours: halb drei = 2:30.",
		Example: "Es ist halb zwölf.",
	},
	{
		Title:   "Greetings & Farewells",
		Tip:     "Guten Morgen (until 10), Guten Tag (formal), Tschüss (informal).",
		Example: "Guten Morgen! Wie geht es Ihnen?",
	},
	{
		Title:   "Word Formation with Präfixe",
		Tip:     "Verbs can change meaning with prefixes: -ver, -er, -ent-.",
		Example: "stehen vs verstehen (to understand)",
	},
	{
		Title:   "Impersonal Expressions with es",
		Tip:     "Use 'Es gibt' for there is/are, 'Es tut mir leid' for sorry.",
		Example: "Es gibt viele Probleme.",
	},
	{
		Title:   "Adjective Endings After der-words",
		Tip:     "These words carry the adjective ending: der, die, das, ein, dieser, jener.",
		Example: "der große Hund",
	},
	{
		Title:   "Adjective Endings After ein-words",
		Tip:     "With 'ein' words, adjective has weak endings: mein kleines Haus.",
		Example: "ein neues Auto",
	},
	{
		Title:   "Past Tense with haben/sein",
		Tip:     "Narrative German prefers Präteritum for 'sein, haben, werden, modals'.",
		Example: "Ich war müde. Er hatte Hunger.",
	},
	{
		Title:   "Ordinal Numbers",
		Tip:     "Add -te for numbers up to 19, then -ste: der dritte, der zwanzigste.",
		Example: "Am ersten Mai feiern wir.",
	},
	{
		Title:   "Genitive Case Overview",
		Tip:     "Use genitive after 'because of' words: wegen, während, trotz.",
		Example: "Während des Unterrichts.",
	},
}

func GetDailyGrammarTip() GrammarTip {
	dayOfYear := time.Now().YearDay()
	return grammarTips[dayOfYear%len(grammarTips)]
}
