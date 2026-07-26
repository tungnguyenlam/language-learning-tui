package content

// RelativeExercise is a fill-in-the-blank exercise for the Relative Clauses
// (Relativsätze & Relativpronomen) trainer. Answers are the correct relative pronoun.
type RelativeExercise struct {
	Sentence    string // fill-in-the-blank sentence using {{...}}
	Answer      string // correct relative pronoun (e.g. der, die, das, den, dem, deren, dessen, denen, worüber, mit dem)
	Meaning     string // English meaning of the full sentence
	Hint        string // short grammar hint (Gender/Case/Preposition)
	Explanation string // full explanation shown on reveal
}

// GetRelativeExercises returns all Relative Clause exercises.
// Exercises cover:
//   - Nominative relative pronouns (der, die, das, die)
//   - Accusative relative pronouns (den, die, das, die)
//   - Dative relative pronouns (dem, der, dem, denen)
//   - Genitive relative pronouns (dessen, deren, dessen, deren)
//   - Preposition + relative pronoun (mit dem, für die, an denen, worüber, etc.)
func GetRelativeExercises() []RelativeExercise {
	return []RelativeExercise{
		// --- Nominative Relative Pronouns ---
		{
			Sentence:    "Der Mann, {{...}} dort drüben steht, ist mein Deutschlehrer.",
			Answer:      "der",
			Meaning:     "The man who is standing over there is my German teacher.",
			Hint:        "Nominative Masculine: der",
			Explanation: "'Der Mann' is masculine singular and is the subject (Nominative) of the relative clause → 'der'.",
		},
		{
			Sentence:    "Die Frau, {{...}} im Supermarkt arbeitet, ist sehr nett.",
			Answer:      "die",
			Meaning:     "The woman who works in the supermarket is very nice.",
			Hint:        "Nominative Feminine: die",
			Explanation: "'Die Frau' is feminine singular and is the subject (Nominative) of the relative clause → 'die'.",
		},
		{
			Sentence:    "Das Kind, {{...}} auf dem Spielplatz spielt, heißt Lukas.",
			Answer:      "das",
			Meaning:     "The child who is playing on the playground is named Lukas.",
			Hint:        "Nominative Neuter: das",
			Explanation: "'Das Kind' is neuter singular and is the subject (Nominative) of the relative clause → 'das'.",
		},
		{
			Sentence:    "Die Studenten, {{...}} fleißig lernen, bestehen die Prüfung.",
			Answer:      "die",
			Meaning:     "The students who study hard pass the exam.",
			Hint:        "Nominative Plural: die",
			Explanation: "'Die Studenten' is plural and is the subject (Nominative) of the relative clause → 'die'.",
		},
		// --- Accusative Relative Pronouns ---
		{
			Sentence:    "Der Film, {{...}} wir gestern gesehen haben, war unglaublich spannend.",
			Answer:      "den",
			Meaning:     "The movie that we saw yesterday was incredibly exciting.",
			Hint:        "Accusative Masculine: den",
			Explanation: "'Der Film' is masculine singular. In the relative clause 'wir' is the subject and 'den' is the direct object (Accusative) → 'den'.",
		},
		{
			Sentence:    "Die Tasche, {{...}} sie gekauft hat, ist aus echtem Leder.",
			Answer:      "die",
			Meaning:     "The bag that she bought is made of real leather.",
			Hint:        "Accusative Feminine: die",
			Explanation: "'Die Tasche' is feminine singular and serves as the direct object (Accusative) of 'gekauft hat' → 'die'.",
		},
		{
			Sentence:    "Das Buch, {{...}} ich von dir geliehen habe, liegt auf dem Tisch.",
			Answer:      "das",
			Meaning:     "The book that I borrowed from you is lying on the table.",
			Hint:        "Accusative Neuter: das",
			Explanation: "'Das Buch' is neuter singular and serves as the direct object (Accusative) of 'geliehen habe' → 'das'.",
		},
		{
			Sentence:    "Die Vokabeln, {{...}} ich heute gelernt habe, sind sehr nützlich.",
			Answer:      "die",
			Meaning:     "The vocabulary words that I learned today are very useful.",
			Hint:        "Accusative Plural: die",
			Explanation: "'Die Vokabeln' is plural and serves as the direct object (Accusative) in the relative clause → 'die'.",
		},
		// --- Dative Relative Pronouns ---
		{
			Sentence:    "Der Kollege, {{...}} ich geholfen habe, bedankte sich herzlich.",
			Answer:      "dem",
			Meaning:     "The colleague whom I helped thanked me warmly.",
			Hint:        "Dative Masculine: dem (verben mit Dativ: helfen)",
			Explanation: "'helfen' requires the Dative case. 'Der Kollege' is masculine singular → Dative 'dem'.",
		},
		{
			Sentence:    "Die Ärztin, {{...}} ich vertraue, hat viel Erfahrung.",
			Answer:      "der",
			Meaning:     "The doctor whom I trust has a lot of experience.",
			Hint:        "Dative Feminine: der (verben mit Dativ: vertrauen)",
			Explanation: "'vertrauen' requires the Dative case. 'Die Ärztin' is feminine singular → Dative 'der'.",
		},
		{
			Sentence:    "Das Mädchen, {{...}} ich ein Geschenk gegeben habe, lächelte.",
			Answer:      "dem",
			Meaning:     "The girl to whom I gave a gift smiled.",
			Hint:        "Dative Neuter: dem (verben mit Dativ: geben)",
			Explanation: "'geben' takes an indirect object in Dative. 'Das Mädchen' is neuter singular → Dative 'dem'.",
		},
		{
			Sentence:    "Die Freunde, {{...}} wir gratuliert haben, feierten bis spät in die Nacht.",
			Answer:      "denen",
			Meaning:     "The friends whom we congratulated celebrated until late into the night.",
			Hint:        "Dative Plural: denen (verb: gratulieren + Dativ plural)",
			Explanation: "Dative plural relative pronoun is 'denen' (NOT 'die' or 'den'). 'gratulieren' takes Dative.",
		},
		// --- Genitive Relative Pronouns ---
		{
			Sentence:    "Der Autor, {{...}} Buch ein Bestseller wurde, hält eine Vorlesung.",
			Answer:      "dessen",
			Meaning:     "The author whose book became a bestseller is giving a lecture.",
			Hint:        "Genitive Masculine: dessen (whose)",
			Explanation: "Genitive relative pronoun for masculine singular antecedents is 'dessen' ('whose book').",
		},
		{
			Sentence:    "Die Frau, {{...}} Auto gestohlen wurde, hat die Polizei angerufen.",
			Answer:      "deren",
			Meaning:     "The woman whose car was stolen called the police.",
			Hint:        "Genitive Feminine: deren (whose)",
			Explanation: "Genitive relative pronoun for feminine singular antecedents is 'deren' ('whose car').",
		},
		{
			Sentence:    "Das Haus, {{...}} Dach beschädigt ist, wird saniert.",
			Answer:      "dessen",
			Meaning:     "The house whose roof is damaged is being renovated.",
			Hint:        "Genitive Neuter: dessen (whose)",
			Explanation: "Genitive relative pronoun for neuter singular antecedents is 'dessen' ('whose roof').",
		},
		{
			Sentence:    "Eltern, {{...}} Kinder in den Kindergarten gehen, haben oft wenig Zeit.",
			Answer:      "deren",
			Meaning:     "Parents whose children go to kindergarten often have little time.",
			Hint:        "Genitive Plural: deren (whose)",
			Explanation: "Genitive relative pronoun for plural antecedents is 'deren' ('whose children').",
		},
		// --- Preposition + Relative Pronoun ---
		{
			Sentence:    "Der Nachbar, mit {{...}} ich gesprochen habe, wohnt im ersten Stock.",
			Answer:      "dem",
			Meaning:     "The neighbor with whom I spoke lives on the first floor.",
			Hint:        "Preposition + Dative Masculine: mit + dem",
			Explanation: "'mit' always takes Dative. Antecedent 'der Nachbar' is masculine singular → 'mit dem'.",
		},
		{
			Sentence:    "Die Stadt, in {{...}} sie geboren wurde, ist historisch bedeutsam.",
			Answer:      "der",
			Meaning:     "The city in which she was born is historically significant.",
			Hint:        "Preposition + Dative Feminine: in + der (Wo? location)",
			Explanation: "'in' + Dative indicates location (Wo?). 'Die Stadt' is feminine singular → 'in der'.",
		},
		{
			Sentence:    "Das Projekt, an {{...}} wir gearbeitet haben, war erfolgreich.",
			Answer:      "dem",
			Meaning:     "The project on which we worked was successful.",
			Hint:        "Preposition + Dative Neuter: an + dem (arbeiten an + Dat)",
			Explanation: "'arbeiten an' takes Dative. 'Das Projekt' is neuter singular → 'an dem'.",
		},
		{
			Sentence:    "Die Kollegen, ohne {{...}} wir das Ziel nicht erreicht hätten, wurden gelobt.",
			Answer:      "die",
			Meaning:     "The colleagues without whom we would not have reached the goal were praised.",
			Hint:        "Preposition + Accusative Plural: ohne + die",
			Explanation: "'ohne' always takes Accusative. Plural relative pronoun in Accusative is 'die'.",
		},
		{
			Sentence:    "Der Termin, für {{...}} ich mich vorbereite, ist morgen.",
			Answer:      "den",
			Meaning:     "The appointment for which I am preparing is tomorrow.",
			Hint:        "Preposition + Accusative Masculine: für + den",
			Explanation: "'für' always takes Accusative. 'Der Termin' is masculine singular → 'für den'.",
		},
		{
			Sentence:    "Die Freunde, bei {{...}} wir übernachtet haben, wohnen in Berlin.",
			Answer:      "denen",
			Meaning:     "The friends with whom (at whose place) we stayed live in Berlin.",
			Hint:        "Preposition + Dative Plural: bei + denen",
			Explanation: "'bei' always takes Dative. Plural Dative relative pronoun is 'denen' → 'bei denen'.",
		},
		{
			Sentence:    "Das Thema, über {{...}} wir diskutiert haben, war sehr komplex.",
			Answer:      "das",
			Meaning:     "The topic about which we discussed was very complex.",
			Hint:        "Preposition + Accusative Neuter: über + das (diskutieren über + Acc)",
			Explanation: "'diskutieren über' takes Accusative. 'Das Thema' is neuter singular → 'über das'.",
		},
		{
			Sentence:    "Das ist alles, {{...}} ich dir sagen wollte.",
			Answer:      "was",
			Meaning:     "That is everything that I wanted to tell you.",
			Hint:        "Relative pronoun 'was' after indefinites (alles/nichts/etwas)",
			Explanation: "After indefinite pronouns like 'alles', 'nichts', or 'etwas', use 'was' as the relative pronoun instead of 'das'.",
		},
		{
			Sentence:    "Das Beste, {{...}} mir heute passiert ist, war dein Anruf.",
			Answer:      "was",
			Meaning:     "The best thing that happened to me today was your phone call.",
			Hint:        "Relative pronoun 'was' after nominalized superlatives (das Beste)",
			Explanation: "After neuter nominalized superlatives like 'das Beste', 'das Schönste', or 'das Erste', use 'was'.",
		},
	}
}
