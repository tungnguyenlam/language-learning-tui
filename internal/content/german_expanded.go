package content

import (
	"deutsch-tui/internal/core"
	"fmt"
)

// GermanExpandedDeck provides comprehensive German vocabulary covering A1-B2 levels
// with thematic content for travel, business, food, and daily life
func GermanExpandedDeck() core.Deck {
	notes := []core.Note{
		// --- A1 Level: Essential Survival ---
		// Greetings expanded
		{ID: "a1-greet-guten-tag", DeckID: "german-expanded", Front: "Guten Tag", Back: "good day", Tags: []string{"a1", "greeting", "formal"}},
		{ID: "a1-greet-guten-morgen", DeckID: "german-expanded", Front: "Guten Morgen", Back: "good morning", Tags: []string{"a1", "greeting", "formal"}},
		{ID: "a1-greet-guten-abend", DeckID: "german-expanded", Front: "Guten Abend", Back: "good evening", Tags: []string{"a1", "greeting", "formal"}},
		{ID: "a1-greet-hallo", DeckID: "german-expanded", Front: "Hallo", Back: "hello", Tags: []string{"a1", "greeting", "informal"}},
		{ID: "a1-greet-tschuss", DeckID: "german-expanded", Front: "Tschüss", Back: "bye", Tags: []string{"a1", "greeting", "informal"}},
		{ID: "a1-greet-auf-wiedersehen", DeckID: "german-expanded", Front: "Auf Wiedersehen", Back: "goodbye", Tags: []string{"a1", "greeting", "formal"}},
		{ID: "a1-greet-servus", DeckID: "german-expanded", Front: "Servus", Back: "hi/bye", Tags: []string{"a1", "greeting", "regional"}},
		{ID: "a1-greet-gruss-gott", DeckID: "german-expanded", Front: "Grüß Gott", Back: "greetings", Tags: []string{"a1", "greeting", "regional"}},
		{ID: "a1-greet-moin", DeckID: "german-expanded", Front: "Moin", Back: "hi", Tags: []string{"a1", "greeting", "regional"}},

		// Essential phrases
		{ID: "a1-phrase-bitte", DeckID: "german-expanded", Front: "Bitte", Back: "please / you're welcome", Extra: "Used as both 'please' and 'you're welcome'", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-danke", DeckID: "german-expanded", Front: "Danke", Back: "thank you", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-danke-schoen", DeckID: "german-expanded", Front: "Danke schön", Back: "thank you very much", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-vielen-dank", DeckID: "german-expanded", Front: "Vielen Dank", Back: "many thanks", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-entschuldigung", DeckID: "german-expanded", Front: "Entschuldigung", Back: "excuse me / sorry", Extra: "For apologies and getting attention", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-es-tut-mir-leid", DeckID: "german-expanded", Front: "Es tut mir leid", Back: "I'm sorry", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-ja", DeckID: "german-expanded", Front: "Ja", Back: "yes", Tags: []string{"a1", "essential"}},
		{ID: "a1-phrase-nein", DeckID: "german-expanded", Front: "Nein", Back: "no", Tags: []string{"a1", "essential"}},
		{ID: "a1-phrase-vielleicht", DeckID: "german-expanded", Front: "Vielleicht", Back: "maybe", Tags: []string{"a1", "essential"}},
		{ID: "a1-phrase-bitte-schoen", DeckID: "german-expanded", Front: "Bitte schön", Back: "you're welcome", Tags: []string{"a1", "essential", "polite"}},
		{ID: "a1-phrase-gern-geschehen", DeckID: "german-expanded", Front: "Gern geschehen", Back: "you're welcome", Tags: []string{"a1", "essential", "polite"}},

		// Basic questions
		{ID: "a1-question-wie-gehts", DeckID: "german-expanded", Front: "Wie geht es dir?", Back: "How are you? (informal)", Tags: []string{"a1", "question"}},
		{ID: "a1-question-wie-gehts-formal", DeckID: "german-expanded", Front: "Wie geht es Ihnen?", Back: "How are you? (formal)", Tags: []string{"a1", "question", "formal"}},
		{ID: "a1-question-was-ist", DeckID: "german-expanded", Front: "Was ist...?", Back: "What is...?", Tags: []string{"a1", "question"}},
		{ID: "a1-question-wo-ist", DeckID: "german-expanded", Front: "Wo ist...?", Back: "Where is...?", Tags: []string{"a1", "question"}},
		{ID: "a1-question-wie-spellt", DeckID: "german-expanded", Front: "Wie schreibt man das?", Back: "How do you spell that?", Tags: []string{"a1", "question"}},
		{ID: "a1-question-verstehen", DeckID: "german-expanded", Front: "Verstehen Sie?", Back: "Do you understand? (formal)", Tags: []string{"a1", "question", "formal"}},
		{ID: "a1-question-verstehst", DeckID: "german-expanded", Front: "Verstehst du?", Back: "Do you understand? (informal)", Tags: []string{"a1", "question"}},
		{ID: "a1-question-wie-heisst", DeckID: "german-expanded", Front: "Wie heißt du?", Back: "What's your name? (informal)", Tags: []string{"a1", "question"}},
		{ID: "a1-question-wie-heisst-formal", DeckID: "german-expanded", Front: "Wie heißen Sie?", Back: "What's your name? (formal)", Tags: []string{"a1", "question", "formal"}},
		{ID: "a1-question-kommen-sie", DeckID: "german-expanded", Front: "Kommen Sie aus...?", Back: "Are you from...? (formal)", Tags: []string{"a1", "question", "formal"}},

		// Personal pronouns
		{ID: "a1-grammar-ich", DeckID: "german-expanded", Front: "ich", Back: "I", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-du", DeckID: "german-expanded", Front: "du", Back: "you (informal)", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-er", DeckID: "german-expanded", Front: "er", Back: "he", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-sie-s", DeckID: "german-expanded", Front: "sie", Back: "she", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-es", DeckID: "german-expanded", Front: "es", Back: "it", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-wir", DeckID: "german-expanded", Front: "wir", Back: "we", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-ihr", DeckID: "german-expanded", Front: "ihr", Back: "you all (informal)", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-sie-p", DeckID: "german-expanded", Front: "sie", Back: "they / you (formal)", Tags: []string{"a1", "grammar", "pronoun"}},
		{ID: "a1-grammar-Sie", DeckID: "german-expanded", Front: "Sie", Back: "you (formal)", Tags: []string{"a1", "grammar", "pronoun", "formal"}},
	}

	// Pronoun MCQs - handled in verb conjugation section below

	// Essential verbs with full conjugations
	verbs := []struct {
		infinitive string
		english    string
		singular   []string // ich, du, er/sie/es
		plural     []string // wir, ihr, sie/Sie
	}{
		{"sein", "to be", []string{"bin", "bist", "ist"}, []string{"sind", "seid", "sind"}},
		{"haben", "to have", []string{"habe", "hast", "hat"}, []string{"haben", "habt", "haben"}},
		{"gehen", "to go", []string{"gehe", "gehst", "geht"}, []string{"gehen", "geht", "gehen"}},
		{"kommen", "to come", []string{"komme", "kommst", "kommt"}, []string{"kommen", "kommt", "kommen"}},
		{"machen", "to do/make", []string{"mache", "machst", "macht"}, []string{"machen", "macht", "machen"}},
		{"sprechen", "to speak", []string{"spreche", "sprichst", "spricht"}, []string{"sprechen", "sprecht", "sprechen"}},
		{"essen", "to eat", []string{"esse", "isst", "isst"}, []string{"essen", "esst", "essen"}},
		{"trinken", "to drink", []string{"trinke", "trinkst", "trinkt"}, []string{"trinken", "trinkt", "trinken"}},
		{"lesen", "to read", []string{"lese", "liest", "liest"}, []string{"lesen", "lest", "lesen"}},
		{"schreiben", "to write", []string{"schreibe", "schreibst", "schreibt"}, []string{"schreiben", "schreibt", "schreiben"}},
		{"sehen", "to see", []string{"sehe", "siehst", "sieht"}, []string{"sehen", "seht", "sehen"}},
		{"wollen", "to want", []string{"will", "willst", "will"}, []string{"wollen", "wollt", "wollen"}},
		{"können", "can/to be able to", []string{"kann", "kannst", "kann"}, []string{"können", "könnt", "können"}},
		{"müssen", "must/to have to", []string{"muss", "musst", "muss"}, []string{"müssen", "müsst", "müssen"}},
		{"sollen", "should/to ought to", []string{"soll", "sollst", "soll"}, []string{"sollen", "sollt", "sollen"}},
		{"heißen", "to be called", []string{"heiße", "heißt", "heißt"}, []string{"heißen", "heißt", "heißen"}},
		{"wohnen", "to live/reside", []string{"wohne", "wohnst", "wohnt"}, []string{"wohnen", "wohnt", "wohnen"}},
		{"arbeiten", "to work", []string{"arbeite", "arbeitest", "arbeitet"}, []string{"arbeiten", "arbeitet", "arbeiten"}},
		{"spielen", "to play", []string{"spiele", "spielst", "spielt"}, []string{"spielen", "spielt", "spielen"}},
		{"lernen", "to learn", []string{"lerne", "lernst", "lernt"}, []string{"lernen", "lernt", "lernen"}},
		{"kaufen", "to buy", []string{"kaufe", "kaufst", "kauft"}, []string{"kaufen", "kauft", "kaufen"}},
		{"verkaufen", "to sell", []string{"verkaufe", "verkaufst", "verkauft"}, []string{"verkaufen", "verkauft", "verkaufen"}},
		{"finden", "to find", []string{"finde", "findest", "findet"}, []string{"finden", "findet", "finden"}},
		{"geben", "to give", []string{"gebe", "gibst", "gibt"}, []string{"geben", "gebt", "geben"}},
		{"nehmen", "to take", []string{"nehme", "nimmst", "nimmt"}, []string{"nehmen", "nehmt", "nehmen"}},
		{"lassen", "to let/allow", []string{"lasse", "lässt", "lässt"}, []string{"lassen", "lasst", "lassen"}},
		{"bringen", "to bring", []string{"bringe", "bringst", "bringt"}, []string{"bringen", "bringt", "bringen"}},
		{"bleiben", "to stay", []string{"bleibe", "bleibst", "bleibt"}, []string{"bleiben", "bleibt", "bleiben"}},
		{"beginnen", "to begin", []string{"beginne", "beginnst", "beginnt"}, []string{"beginnen", "beginnt", "beginnen"}},
		{"vergessen", "to forget", []string{"vergesse", "vergisst", "vergisst"}, []string{"vergessen", "vergesst", "vergessen"}},
		{"denken", "to think", []string{"denke", "denkst", "denkt"}, []string{"denken", "denkt", "denken"}},
		{"verstehen", "to understand", []string{"verstehe", "verstehst", "versteht"}, []string{"verstehen", "versteht", "verstehen"}},
		{"brauchen", "to need", []string{"brauche", "brauchst", "braucht"}, []string{"brauchen", "braucht", "brauchen"}},
		{"suchen", "to search/look for", []string{"suche", "suchst", "sucht"}, []string{"suchen", "sucht", "suchen"}},
		{"fahren", "to drive/travel", []string{"fahre", "fährst", "fährt"}, []string{"fahren", "fahrt", "fahren"}},
		{"fliegen", "to fly", []string{"fliege", "fliegst", "fliegt"}, []string{"fliegen", "fliegt", "fliegen"}},
		{"reisen", "to travel", []string{"reise", "reist", "reist"}, []string{"reisen", "reist", "reisen"}},
		{"schwimmen", "to swim", []string{"schwimme", "schwimmst", "schwimmt"}, []string{"schwimmen", "schwimmt", "schwimmen"}},
		{"laufen", "to run/walk", []string{"laufe", "läufst", "läuft"}, []string{"laufen", "lauft", "laufen"}},
		{"tanzen", "to dance", []string{"tanze", "tanzt", "tanzt"}, []string{"tanzen", "tanzt", "tanzen"}},
		{"singen", "to sing", []string{"singe", "singst", "singt"}, []string{"singen", "singt", "singen"}},
		{"hören", "to hear", []string{"höre", "hörst", "hört"}, []string{"hören", "hört", "hören"}},
		{"sehen", "to see", []string{"sehe", "siehst", "sieht"}, []string{"sehen", "seht", "sehen"}},
		{"öffnen", "to open", []string{"öffne", "öffnest", "öffnet"}, []string{"öffnen", "öffnet", "öffnen"}},
		{"schließen", "to close", []string{"schließe", "schließt", "schließt"}, []string{"schließen", "schließt", "schließen"}},
		{"sagen", "to say/tell", []string{"sag", "sagst", "sagt"}, []string{"sagen", "sagt", "sagen"}},
		{"fragen", "to ask", []string{"frage", "fragst", "fragt"}, []string{"fragen", "fragt", "fragen"}},
		{"antworten", "to answer", []string{"antworte", "antwortest", "antwortet"}, []string{"antworten", "antwortet", "antworten"}},
		{"schlafen", "to sleep", []string{"schlafe", "schläfst", "schläft"}, []string{"schlafen", "schläft", "schlafen"}},
		{"wachen", "to wake up", []string{"wache", "wachst", "wacht"}, []string{"wachen", "wacht", "wachen"}},
		{"putzen", "to clean", []string{"putze", "putzt", "putzt"}, []string{"putzen", "putzt", "putzen"}},
		{"waschen", "to wash", []string{"wasche", "wäschst", "wäscht"}, []string{"waschen", "wäscht", "waschen"}},
		{"kämmen", "to comb", []string{"kämme", "kämmst", "kämmt"}, []string{"kämmen", "kämmt", "kämmen"}},
		{"sich-anziehen", "to get dressed", []string{"ziehe mich an", "ziehst du dich an", "zieht sich an"}, []string{"ziehen sich an", "zieht euch an", "ziehen sich an"}},
	}

	for i, v := range verbs {
		note := core.Note{
			ID:     fmt.Sprintf("a1-verb-%d-%s", i, v.infinitive),
			DeckID: "german-expanded",
			Front:  v.infinitive,
			Back:   v.english,
			Extra: fmt.Sprintf("ich %s, du %s, er/sie/es %s | wir %s, ihr %s, sie/Sie %s",
				v.singular[0], v.singular[1], v.singular[2],
				v.plural[0], v.plural[1], v.plural[2]),
			Tags:     []string{"a1", "verb", "conjugation"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Common nouns - Food & Drink
	foodItems := []struct {
		german, english, gender, plural string
	}{
		{"Apfel", "apple", "der", "die Äpfel"},
		{"Birne", "pear", "die", "die Birnen"},
		{"Orange", "orange", "die", "die Orangen"},
		{"Zitrone", "lemon", "die", "die Zitronen"},
		{"Banane", "banana", "die", "die Bananen"},
		{"Traube", "grape", "die", "die Trauben"},
		{"Kirsche", "cherry", "die", "die Kirschen"},
		{"Erdbeere", "strawberry", "die", "die Erdbeeren"},
		{"Brot", "bread", "das", "die Brote"},
		{"Brötchen", "bread roll", "das", "die Brötchen"},
		{"Semmel", "bread roll (Bavaria)", "die", "die Semmeln"},
		{"Käse", "cheese", "der", "die Käse"},
		{"Milch", "milk", "die", "die Milch"},
		{"Butter", "butter", "die", "die Butter"},
		{"Joghurt", "yogurt", "der", "die Joghurts"},
		{"Quark", "quark/curd", "der", "die Quarks"},
		{"Ei", "egg", "das", "die Eier"},
		{"Fleisch", "meat", "das", "die Fleisch"},
		{"Rindfleisch", "beef", "das", "die Rindfleische"},
		{"Schweinefleisch", "pork", "das", "die Schweinefleische"},
		{"Hähnchen", "chicken", "das", "die Hähnchen"},
		{"Wurst", "sausage", "die", "die Würste"},
		{"Schinken", "ham", "der", "die Schinken"},
		{"Salami", "salami", "die", "die Salamis"},
		{"Fisch", "fish", "der", "die Fische"},
		{"Gemüse", "vegetables", "das", "die Gemüse"},
		{"Kartoffel", "potato", "die", "die Kartoffeln"},
		{"Reis", "rice", "der", "die Reis"},
		{"Nudel", "pasta/noodle", "die", "die Nudeln"},
		{"Spätzle", "Swabian noodles", "die", "die Spätzle"},
		{"Salat", "salad", "der", "die Salate"},
		{"Suppe", "soup", "die", "die Suppen"},
		{"Pizza", "pizza", "die", "die Pizzen"},
		{"Kuchen", "cake/tart", "der", "die Kuchen"},
		{"Torte", "cake (layered)", "die", "die Torten"},
		{"Brot", "bread", "das", "die Brote"},
		{"Bier", "beer", "das", "die Biere"},
		{"Wein", "wine", "der", "die Weine"},
		{"Wasser", "water", "das", "die Wasser"},
		{"Saft", "juice", "der", "die Säfte"},
		{"Kaffee", "coffee", "der", "die Kaffees"},
		{"Tee", "tea", "der", "die Tees"},
		{"Cola", "cola", "die", "die Colas"},
		{"Zucker", "sugar", "der", "die Zucker"},
		{"Salz", "salt", "das", "die Salze"},
		{"Pfeffer", "pepper", "der", "die Pfeffer"},
		{"Öl", "oil", "das", "die Öle"},
		{"Essig", "vinegar", "der", "die Essige"},
		{"Senf", "mustard", "der", "die Senf"},
		{"Soße", "sauce", "die", "die Soßen"},
		{"Gewürz", "spice", "das", "die Gewürze"},
	}

	for i, f := range foodItems {
		note := core.Note{
			ID:       fmt.Sprintf("a1-food-%d-%s", i, f.german),
			DeckID:   "german-expanded",
			Front:    fmt.Sprintf("der/die/das %s", f.german),
			Back:     f.english,
			Extra:    fmt.Sprintf("Gender: %s | Plural: %s", f.german[:3], f.plural), // Simplified
			Tags:     []string{"a1", "noun", "food", "drink"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Numbers 0-100
	numberWords := []string{"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn",
		"elf", "zwölf", "dreizehn", "vierzehn", "fünfzehn", "sechzehn", "siebzehn", "achtzehn", "neunzehn",
		"zwanzig", "einundzwanzig", "zweiundzwanzig", "dreiundzwanzig", "vierundzwanzig", "fünfundzwanzig",
		"sechsundzwanzig", "siebenundzwanzig", "achtundzwanzig", "neunundzwanzig", "dreißig",
		"einunddreißig", "zweiunddreißig", "dreiunddreißig", "vierunddreißig", "fünfunddreißig", "sechsunddreißig", "siebenunddreißig", "achtunddreißig", "neununddreißig",
		"vierzig", "einundvierzig", "zweiundvierzig", "dreiundvierzig", "vierundvierzig", "fünfundvierzig",
		"sechsundvierzig", "siebenundvierzig", "achtundvierzig", "neunundvierzig",
		"fünfzig", "einundfünfzig", "zweiundfünfzig", "dreiundfünfzig", "vierundfünfzig", "fünfundfünfzig",
		"sechsundfünfzig", "siebenundfünfzig", "achtundfünfzig", "neunundfünfzig",
		"sechzig", "einundsechzig", "zweiundsechzig", "dreiundsechzig", "vierundsechzig", "fünfundsechzig",
		"sechsundsechzig", "siebenundsechzig", "achtundsechzig", "neunundsechzig",
		"siebzig", "einundsiebzig", "zweiundsiebzig", "dreiundsiebzig", "vierundsiebzig", "fünfundsiebzig",
		"sechsundsiebzig", "siebenundsiebzig", "achtundsiebzig", "neunundsiebzig",
		"achtzig", "einundachtzig", "zweiundachtzig", "dreiundachtzig", "vierundachtzig", "fünfundachtzig",
		"sechsundachtzig", "siebenundachtzig", "achtundachtzig", "neunundachtzig",
		"neunzig", "einundneunzig", "zweiundneunzig", "dreiundneunzig", "vierundneunzig", "fünfundneunzig",
		"sechsundneunzig", "siebenundneunzig", "achtundneunzig", "neunundneunzig", "hundert"}

	for i, num := range numberWords {
		note := core.Note{
			ID:       fmt.Sprintf("a1-num-%d", i),
			DeckID:   "german-expanded",
			Front:    num,
			Back:     fmt.Sprintf("%d", i),
			Tags:     []string{"a1", "number"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Colors
	colors := []struct {
		german, english string
	}{
		{"rot", "red"}, {"blau", "blue"}, {"grün", "green"}, {"gelb", "yellow"},
		{"schwarz", "black"}, {"weiß", "white"}, {"braun", "brown"}, {"grau", "grey/gray"},
		{"rosa", "pink"}, {"lila", "purple"}, {"orange", "orange"}, {"beige", "beige"},
		{"gold", "gold"}, {"silber", "silver"},
	}

	for i, c := range colors {
		note := core.Note{
			ID:       fmt.Sprintf("a1-color-%d-%s", i, c.german),
			DeckID:   "german-expanded",
			Front:    c.german,
			Back:     c.english,
			Tags:     []string{"a1", "color"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Family
	family := []struct {
		german, english, gender, article string
	}{
		{"Mann", "man/husband", "der", "der"}, {"Frau", "woman/wife", "die", "die"}, {"Kind", "child", "das", "das"},
		{"Vater", "father", "der", "der"}, {"Mutter", "mother", "die", "die"}, {"Sohn", "son", "der", "der"},
		{"Tochter", "daughter", "die", "die"}, {"Bruder", "brother", "der", "der"}, {"Schwester", "sister", "die", "die"},
		{"Großvater", "grandfather", "der", "der"}, {"Großmutter", "grandmother", "die", "die"},
		{"Opa", "grandpa", "der", "der"}, {"Oma", "grandma", "die", "die"},
		{"Onkel", "uncle", "der", "der"}, {"Tante", "aunt", "die", "die"},
		{"Cousine", "female cousin", "die", "die"}, {"Cousin", "male cousin", "der", "der"},
		{"Enkel", "grandson", "der", "der"}, {"Enkelin", "granddaughter", "die", "die"},
		{"Ehemann", "husband", "der", "der"}, {"Ehefrau", "wife", "die", "die"},
	}

	for i, f := range family {
		note := core.Note{
			ID:       fmt.Sprintf("a1-fam-%d-%s", i, f.german),
			DeckID:   "german-expanded",
			Front:    fmt.Sprintf("%s %s", f.article, f.german),
			Back:     f.english,
			Extra:    fmt.Sprintf("Gender: %s | Article: %s", f.german[:3], f.article),
			Tags:     []string{"a1", "noun", "family", "people"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Body parts
	body := []struct {
		german, english, article string
	}{
		{"Kopf", "head", "der"}, {"Gesicht", "face", "das"}, {"Haar", "hair", "das"}, {"Auge", "eye", "das"},
		{"Ohr", "ear", "das"}, {"Nase", "nose", "die"}, {"Mund", "mouth", "der"}, {"Zahn", "tooth", "der"},
		{"Zunge", "tongue", "die"}, {"Hals", "neck", "der"}, {"Brust", "chest", "die"}, {"Bauch", "stomach", "der"},
		{"Rücken", "back", "der"}, {"Arm", "arm", "der"}, {"Hand", "hand", "die"}, {"Finger", "finger", "der"},
		{"Bein", "leg", "das"}, {"Fuß", "foot", "der"}, {"Zehe", "toe", "die"}, {"Schulter", "shoulder", "die"},
		{"Ellenbogen", "elbow", "der"}, {"Knie", "knee", "das"}, {"Hüfte", "hip", "die"},
	}

	for i, b := range body {
		note := core.Note{
			ID:       fmt.Sprintf("a1-body-%d-%s", i, b.german),
			DeckID:   "german-expanded",
			Front:    fmt.Sprintf("%s %s", b.article, b.german),
			Back:     b.english,
			Extra:    fmt.Sprintf("Article: %s", b.article),
			Tags:     []string{"a1", "noun", "body"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Clothes - Basic
	clothes := []struct {
		german, english, gender, article string
	}{
		{"Hemd", "shirt", "das", "das"}, {"T-Shirt", "t-shirt", "das", "das"}, {"Pullover", "sweater", "der", "der"},
		{"Jacke", "jacket", "die", "die"}, {"Mantel", "coat", "der", "der"}, {"Hose", "pants", "die", "die"},
		{"Kleid", "dress", "das", "das"}, {"Rock", "skirt", "der", "der"}, {"Shorts", "shorts", "die", "die"},
		{"Socke", "sock", "die", "die"}, {"Schuh", "shoe", "der", "der"}, {"Stiefel", "boot", "der", "der"},
		{"Schal", "scarf", "der", "der"}, {"Mütze", "cap/hat", "die", "die"}, {"Handschuh", "glove", "der", "der"},
		{"Unterwäsche", "underwear", "die", "die"}, {"Badeanzug", "swimsuit", "der", "der"}, {"Bikini", "bikini", "der", "der"},
	}

	for i, c := range clothes {
		note := core.Note{
			ID:       fmt.Sprintf("a1-clothes-%d-%s", i, c.german),
			DeckID:   "german-expanded",
			Front:    fmt.Sprintf("%s %s", c.article, c.german),
			Back:     c.english,
			Extra:    fmt.Sprintf("Gender: %s | Article: %s", c.gender, c.article),
			Tags:     []string{"a1", "noun", "clothes"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// House/Rooms
	rooms := []struct {
		german, english, gender, article string
	}{
		{"Haus", "house", "das", "das"}, {"Wohnung", "apartment/flat", "die", "die"}, {"Zimmer", "room", "das", "das"},
		{"Schlafzimmer", "bedroom", "das", "das"}, {"Wohnzimmer", "living room", "das", "das"}, {"Küche", "kitchen", "die", "die"},
		{"Bad", "bathroom", "das", "das"}, {"Toilette", "toilet", "die", "die"}, {"WC", "toilet", "das", "das"},
		{"Flur", "hallway", "der", "der"}, {"Garten", "garden", "der", "der"}, {"Balkon", "balcony", "der", "der"},
		{"Treppe", "stairs", "die", "die"}, {"Fenster", "window", "das", "das"}, {"Tür", "door", "die", "die"},
		{"Tisch", "table", "der", "der"}, {"Stuhl", "chair", "der", "der"}, {"Sofa", "sofa", "das", "das"},
		{"Bett", "bed", "das", "das"}, {"Lammfell", "fur blanket", "das", "das"},
	}

	for i, r := range rooms {
		note := core.Note{
			ID:       fmt.Sprintf("a1-house-%d-%s", i, r.german),
			DeckID:   "german-expanded",
			Front:    fmt.Sprintf("%s %s", r.article, r.german),
			Back:     r.english,
			Extra:    fmt.Sprintf("Gender: %s | Article: %s", r.gender, r.article),
			Tags:     []string{"a1", "noun", "house", "place"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Places in town
	places := []struct {
		german, english string
	}{
		{"Supermarkt", "supermarket"}, {"Bäckerei", "bakery"}, {"Metzgerei", "butcher shop"}, {"Konditorei", "pastry shop"},
		{"Apotheke", "pharmacy"}, {"Drogerie", "drugstore"}, {"Post", "post office"}, {"Bank", "bank"},
		{"Bibliothek", "library"}, {"Schule", "school"}, {"Universität", "university"}, {"Krankenhaus", "hospital"},
		{"Arztpraxis", "doctor's office"}, {"Kino", "cinema"}, {"Theater", "theater"}, {"Restaurant", "restaurant"},
		{"Café", "cafe"}, {"Bar", "bar"}, {"Hotel", "hotel"}, {"Pension", "guesthouse"},
		{"Bahnhof", "train station"}, {"Flughafen", "airport"}, {"Bushaltestelle", "bus stop"}, {"Tankstelle", "gas station"},
		{"Park", "park"}, {"Platz", "square"}, {"Straße", "street"}, {"U-Bahn", "subway"}, {"S-Bahn", "suburban train"},
		{"Taxi", "taxi"}, {"Bus", "bus"}, {"Auto", "car"},
	}

	for i, p := range places {
		note := core.Note{
			ID:       fmt.Sprintf("a1-place-%d-%s", i, p.german),
			DeckID:   "german-expanded",
			Front:    p.german,
			Back:     p.english,
			Tags:     []string{"a1", "noun", "place"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Time & Calendar
	timeWords := []struct {
		german, english string
	}{
		{"Sekunde", "second"}, {"Minute", "minute"}, {"Stunde", "hour"}, {"Tag", "day"},
		{"Woche", "week"}, {"Monat", "month"}, {"Jahr", "year"}, {"Jahrzehnt", "decade"},
		{"Jahrhundert", "century"}, {"heute", "today"}, {"gestern", "yesterday"}, {"morgen", "tomorrow"},
		{"übermorgen", "day after tomorrow"}, {"vorgestern", "day before yesterday"}, {"nächste Woche", "next week"},
		{"letzte Woche", "last week"}, {"Wochenende", "weekend"}, {"Montag", "Monday"}, {"Dienstag", "Tuesday"},
		{"Mittwoch", "Wednesday"}, {"Donnerstag", "Thursday"}, {"Freitag", "Friday"}, {"Samstag", "Saturday"},
		{"Sonntag", "Sunday"}, {"Januar", "January"}, {"Februar", "February"}, {"März", "March"}, {"April", "April"},
		{"Mai", "May"}, {"Juni", "June"}, {"Juli", "July"}, {"August", "August"}, {"September", "September"},
		{"Oktober", "October"}, {"November", "November"}, {"Dezember", "December"},
	}

	for i, t := range timeWords {
		note := core.Note{
			ID:       fmt.Sprintf("a1-time-%d-%s", i, t.german),
			DeckID:   "german-expanded",
			Front:    t.german,
			Back:     t.english,
			Tags:     []string{"a1", "time"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Weather
	weather := []struct {
		german, english string
	}{
		{"Wetter", "weather"}, {"sonnig", "sunny"}, {"bewölkt", "cloudy"}, {"regnerisch", "rainy"},
		{"Schnee", "snow"}, {"windig", "windy"}, {"stürmisch", "stormy"}, {"Nebel", "fog"},
		{"Hitze", "heat"}, {"Kälte", "cold"}, {"Temperatur", "temperature"}, {"Grad", "degrees"},
		{"regnen", "to rain"}, {"schneien", "to snow"}, {"frieren", "to freeze"}, {"schwitzen", "to sweat"},
		{"heiter", "clear"}, {"warm", "warm"}, {"heiß", "hot"}, {"kalt", "cold"},
		{"kühl", "cool"}, {"mild", "mild"}, {"Hochdruck", "high pressure"}, {"Tiefdruck", "low pressure"},
	}

	for i, w := range weather {
		note := core.Note{
			ID:       fmt.Sprintf("a1-weather-%d-%s", i, w.german),
			DeckID:   "german-expanded",
			Front:    w.german,
			Back:     w.english,
			Tags:     []string{"a1", "weather"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Adjectives (basic)
	adjectives := []struct {
		german, english string
	}{
		{"groß", "big/large"}, {"klein", "small/little"}, {"lang", "long"}, {"kurz", "short"},
		{"hoch", "high/tall"}, {"niedrig", "low"}, {"gut", "good"}, {"schlecht", "bad"},
		{"schön", "beautiful/nice"}, {"häßlich", "ugly"}, {"neu", "new"}, {"alt", "old"},
		{"jung", "young"}, {"alt", "old (person)"}, {"billig", "cheap"}, {"teuer", "expensive"},
		{"leicht", "light/easy"}, {"schwer", "heavy/difficult"}, {"schnell", "fast/quick"}, {"langsam", "slow"},
		{"heiß", "hot"}, {"kalt", "cold"}, {"warm", "warm"}, {"kühl", "cool"},
		{"leer", "empty"}, {"voll", "full"}, {"sauber", "clean"}, {"schmutzig", "dirty"},
		{"stark", "strong"}, {"schwach", "weak"}, {"dick", "thick/fat"}, {"dünn", "thin/slim"},
		{"freundlich", "friendly"}, {"unfreundlich", "unfriendly"}, {"ruhig", "calm/quiet"}, {"laut", "loud"},
		{"sicher", "safe"}, {"gefährlich", "dangerous"}, {"möglich", "possible"}, {"unmöglich", "impossible"},
		{"gesund", "healthy"}, {"krank", "sick"}, {"müde", "tired"}, {"fertig", "ready/finished"},
		{"falsch", "wrong"}, {"richtig", "right/correct"}, {"wichtig", "important"}, {"unwichtig", "unimportant"},
		{"interessant", "interesting"}, {"langweilig", "boring"}, {"schwierig", "difficult"}, {"einfach", "easy"},
		{"glücklich", "happy"}, {"traurig", "sad"}, {"ängstlich", "afraid"}, {"mutig", "brave"},
	}

	for i, a := range adjectives {
		note := core.Note{
			ID:       fmt.Sprintf("a1-adj-%d-%s", i, a.german),
			DeckID:   "german-expanded",
			Front:    a.german,
			Back:     a.english,
			Tags:     []string{"a1", "adjective"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Travel & Transportation (A2 level)
	travel := []struct {
		german, english, extra string
	}{
		{"Reise", "trip/journey", "die Reise"},
		{"Urlaub", "vacation/holiday", "der Urlaub"},
		{"Flug", "flight", "der Flug"},
		{"Zug", "train", "der Zug"},
		{"U-Bahn", "subway", "die U-Bahn"},
		{"S-Bahn", "suburban train", "die S-Bahn"},
		{"Tram", "streetcar", "die Straßenbahn/Tram"},
		{"Fahrkarte", "ticket", "die Fahrkarte"},
		{"Bahnsteig", "platform", "der Bahnsteig"},
		{"Abfahrt", "departure", "die Abfahrt"},
		{"Ankunft", "arrival", "die Ankunft"},
		{"Verspätung", "delay", "die Verspätung"},
		{"Gleis", "track/platform", "das Gleis"},
		{"Pass", "passport", "der Reisepass"},
		{"Gepäck", "luggage", "das Gepäck"},
		{"Koffer", "suitcase", "der Koffer"},
		{"Rucksack", "backpack", "der Rucksack"},
		{"Reiseführer", "guidebook", "der Reiseführer"},
		{"Karte", "map", "die Karte"},
		{"Navi", "GPS navigation", "das Navigationsgerät"},
		{"Grenze", "border", "die Grenze"},
		{"Ausweis", "ID card", "der Ausweis"},
		{"Visum", "visa", "das Visum"},
		{"Flughafen", "airport", "der Flughafen"},
		{"Abflughalle", "departure hall", "die Abflughalle"},
		{"Ankunftshalle", "arrival hall", "die Ankunftshalle"},
		{"Gepäckband", "baggage carousel", "das Gepäckband"},
		{"Kofferraum", "trunk", "der Kofferraum"},
		{"Fahrrad", "bicycle", "das Fahrrad"},
		{"Helm", "helmet", "der Helm"},
	}

	for i, t := range travel {
		note := core.Note{
			ID:       fmt.Sprintf("a2-travel-%d-%s", i, t.german),
			DeckID:   "german-expanded",
			Front:    t.german,
			Back:     t.english,
			Extra:    t.extra,
			Tags:     []string{"a2", "travel", "transportation"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Business German (B1 level)
	business := []struct {
		german, english string
	}{
		{"Arbeit", "work/job"}, {"Beruf", "profession"}, {"Chef", "boss"}, {"Kollege", "colleague"},
		{"Mitarbeiter", "employee"}, {"Angestellte", "employee (female)"}, {"Chefetage", "executive floor"},
		{"Büro", "office"}, {"Schreibtisch", "desk"}, {"Computer", "computer"}, {"Laptop", "laptop"},
		{"Drucker", "printer"}, {"Telefon", "phone"}, {"Handy", "mobile phone"}, {"E-Mail", "email"},
		{"Brief", "letter"}, {"Dokument", "document"}, {"Bericht", "report"}, {"Präsentation", "presentation"},
		{"Meeting", "meeting"}, {"Besprechung", "conference"}, {"Termin", "appointment"}, {"Zeitplan", "schedule"},
		{"Vertrag", "contract"}, {"Angebot", "offer"}, {"Rechnung", "bill/invoice"}, {"Quittung", "receipt"},
		{"Gehalt", "salary"}, {"Lohn", "wage"}, {"Bonus", "bonus"}, {"Urlaubstage", "vacation days"},
		{"Kündigung", "resignation"}, {"Bewerbung", "job application"}, {"Lebenslauf", "CV/resume"},
		{"Vorstellungsgespräch", "job interview"}, {"Referenz", "reference"}, {"Probezeit", "probation period"},
		{"Kunde", "customer"}, {"Käufer", "buyer"}, {"Verkäufer", "seller"}, {"Lieferant", "supplier"},
		{"Produkt", "product"}, {"Dienstleistung", "service"}, {"Qualität", "quality"}, {"Preis", "price"},
		{"Markt", "market"}, {"Wettbewerb", "competition"}, {"Strategie", "strategy"}, {"Ziel", "goal"},
		{"Umsatz", "turnover/revenue"}, {"Gewinn", "profit"}, {"Verlust", "loss"}, {"Kosten", "costs"},
		{"Budget", "budget"}, {"Investition", "investment"}, {"Risiko", "risk"}, {"Chance", "opportunity"},
	}

	for i, b := range business {
		note := core.Note{
			ID:       fmt.Sprintf("b1-business-%d-%s", i, b.german),
			DeckID:   "german-expanded",
			Front:    b.german,
			Back:     b.english,
			Tags:     []string{"b1", "business"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Daily Life & Activities
	activities := []struct {
		german, english string
	}{
		{"aufwachen", "to wake up"}, {"sich anziehen", "to get dressed"}, {"frühstücken", "to have breakfast"},
		{"duschen", "to shower"}, {"Zähne putzen", "to brush teeth"}, {"Haare waschen", "to wash hair"},
		{"rasieren", "to shave"}, {"schminken", "to put on makeup"}, {"kämme", "to comb hair"},
		{"zur Arbeit gehen", "to go to work"}, {"arbeiten", "to work"}, {"Mittagessen", "to have lunch"},
		{"Pause machen", "to take a break"}, {"nach Hause gehen", "to go home"}, {"Abendessen", "to have dinner"},
		{"fernsehen", "to watch TV"}, {"lesen", "to read"}, {"Musik hören", "to listen to music"},
		{"spielen", "to play"}, {"spazieren gehen", "to take a walk"}, {"joggen", "to jog"},
		{"ins Bett gehen", "to go to bed"}, {"schlafen", "to sleep"}, {"Haustier füttern", "to feed pet"},
		{"Gartenarbeit", "gardening"}, {"putzen", "to clean"}, {"Wäsche waschen", "to do laundry"},
		{"einkaufen", "to shop"}, {"kochen", "to cook"}, {"backen", "to bake"},
		{"Abwasch machen", "to wash dishes"}, {"Müll rausbringen", "to take out trash"}, {"Bett machen", "to make bed"},
	}

	for i, a := range activities {
		note := core.Note{
			ID:       fmt.Sprintf("a1-activity-%d-%s", i, a.german),
			DeckID:   "german-expanded",
			Front:    a.german,
			Back:     a.english,
			Tags:     []string{"a1", "activity", "daily"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Common phrases for travel
	travelPhrases := []struct {
		german, english string
	}{
		{"Ich möchte ein Zimmer reservieren.", "I'd like to reserve a room."},
		{"Haben Sie noch ein Zimmer frei?", "Do you have a room available?"},
		{"Wie viel kostet das pro Nacht?", "How much does it cost per night?"},
		{"Ich hätte gerne ein Einzelzimmer.", "I'd like a single room."},
		{"Ich hätte gerne ein Doppelzimmer.", "I'd like a double room."},
		{"Gibt es inklusive Frühstück?", "Is breakfast included?"},
		{"Wo ist die nächste Bushaltestelle?", "Where is the next bus stop?"},
		{"Wie komme ich zum Bahnhof?", "How do I get to the train station?"},
		{"Kann ich mit Karte bezahlen?", "Can I pay by card?"},
		{"Ich hätte gerne die Speisekarte.", "I'd like the menu, please."},
		{"Was empfehlen Sie?", "What do you recommend?"},
		{"Ich bin Vegetarier(in).", "I'm a vegetarian."},
		{"Ich habe eine Allergie gegen Nüsse.", "I have a nut allergy."},
		{"Die Rechnung, bitte.", "The bill, please."},
		{"Können wir getrennt bezahlen?", "Can we pay separately?"},
		{"Entschuldigung, ich habe mich verlaufen.", "Excuse me, I'm lost."},
		{"Ich suche das Museum.", "I'm looking for the museum."},
		{"Können Sie mir helfen?", "Can you help me?"},
		{"Sprechen Sie Englisch?", "Do you speak English?"},
		{"Ich spreche nur ein bisschen Deutsch.", "I only speak a little German."},
	}

	for i, p := range travelPhrases {
		note := core.Note{
			ID:       fmt.Sprintf("phrase-travel-%d", i),
			DeckID:   "german-expanded",
			Front:    p.german,
			Back:     p.english,
			Tags:     []string{"phrase", "travel", "conversation"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Emergency phrases
	emergencyPhrases := []struct {
		german, english string
	}{
		{"Hilfe!", "Help!"},
		{"Rufen Sie einen Arzt!", "Call a doctor!"},
		{"Ich brauche einen Arzt.", "I need a doctor."},
		{"Ich habe Schmerzen.", "I'm in pain."},
		{"Wo ist das Krankenhaus?", "Where is the hospital?"},
		{"Bitte rufen Sie die Polizei.", "Please call the police."},
		{"Ich habe meinen Pass verloren.", "I've lost my passport."},
		{"Mein Geldbeutel wurde gestohlen.", "My wallet was stolen."},
		{"Achtung!", "Watch out!/Caution!"},
		{"Feuer!", "Fire!"},
		{"Es brennt!", "It's on fire!"},
		{"Ich habe einen Unfall gesehen.", "I've seen an accident."},
		{"Notruf", "Emergency call"},
		{"110", "Police (emergency)"},
		{"112", "Fire/medical emergency"},
	}

	for i, e := range emergencyPhrases {
		note := core.Note{
			ID:       fmt.Sprintf("phrase-emergency-%d", i),
			DeckID:   "german-expanded",
			Front:    e.german,
			Back:     e.english,
			Tags:     []string{"phrase", "emergency"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Common connecting words and particles
	connectors := []struct {
		german, english string
	}{
		{"und", "and"},
		{"oder", "or"},
		{"aber", "but"},
		{"denn", "because/for"},
		{"sondern", "but rather"},
		{"weil", "because"},
		{"obwohl", "although"},
		{"wenn", "if/when"},
		{"damit", "so that"},
		{"dass", "that"},
		{"dies", "this"},
		{"das", "that"},
		{"jeder", "every/each"},
		{"alle", "all"},
		{"kein", "no/none"},
		{"ein bisschen", "a little"},
		{"ein wenig", "a little"},
		{"sehr", "very"},
		{"zu", "too"},
		{"auch", "also/too"},
		{"nur", "only"},
		{"noch", "still/yet"},
		{"schon", "already"},
		{"immer", "always"},
		{"nie", "never"},
		{"manchmal", "sometimes"},
		{"oft", "often"},
		{"selten", "seldom"},
		{"heute", "today"},
		{"morgen", "tomorrow"},
		{"gestern", "yesterday"},
		{"jetzt", "now"},
		{"später", "later"},
		{"bald", "soon"},
		{"endlich", "finally"},
		{"leider", "unfortunately"},
		{"glücklicherweise", "fortunately"},
	}

	for i, c := range connectors {
		note := core.Note{
			ID:       fmt.Sprintf("connector-%d-%s", i, c.german),
			DeckID:   "german-expanded",
			Front:    c.german,
			Back:     c.english,
			Tags:     []string{"connector", "grammar"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	// Prepositions
	prepositions := []struct {
		german, english, caseUsed string
	}{
		{"in", "in", "accusative/dative"},
		{"auf", "on/onto", "accusative/dative"},
		{"an", "at/on", "accusative/dative"},
		{"bei", "at/near/with", "dative"},
		{"zu", "to", "dative"},
		{"von", "from/of", "dative"},
		{"für", "for", "accusative"},
		{"mit", "with", "dative"},
		{"durch", "through", "accusative"},
		{"über", "over/about", "accusative/dative"},
		{"unter", "under", "accusative/dative"},
		{"zwischen", "between", "dative"},
		{"gegen", "against/towards", "accusative"},
		{"ohne", "without", "accusative"},
		{"um", "around", "accusative"},
		{"aus", "out of/from", "dative"},
		{"nach", "after/to", "dative"},
		{"seit", "since", "dative"},
		{"gegenüber", "opposite", "dative"},
	}

	for i, p := range prepositions {
		note := core.Note{
			ID:       fmt.Sprintf("prep-%d-%s", i, p.german),
			DeckID:   "german-expanded",
			Front:    p.german,
			Back:     fmt.Sprintf("%s (%s)", p.english, p.caseUsed),
			Tags:     []string{"prep", "grammar"},
			Examples: []string{},
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	return core.Deck{
		ID:          "german-expanded",
		Name:        "German Comprehensive (A1-B2)",
		Description: "Comprehensive German vocabulary covering A1-B2 levels with thematic content for travel, business, food, daily life, and more. Over 500+ flashcards.",
		Tags:        []string{"german", "comprehensive", "a1", "a2", "b1", "b2"},
		Notes:       notes,
	}
}
