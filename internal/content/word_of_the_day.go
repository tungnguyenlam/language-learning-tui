package content

import "time"

type WordOfDay struct {
	German  string
	English string
	Article string
	Plural  string
	Example string
}

var wordsOfDay = []WordOfDay{
	{German: "der Apfel", English: "apple", Article: "der", Plural: "die Äpfel", Example: "Ich esse einen Apfel."},
	{German: "die Blume", English: "flower", Article: "die", Plural: "die Blumen", Example: "Die Blume ist schön."},
	{German: "das Buch", English: "book", Article: "das", Plural: "die Bücher", Example: "Ich lese ein Buch."},
	{German: "der Hund", English: "dog", Article: "der", Plural: "die Hunde", Example: "Der Hund läuft im Park."},
	{German: "die Katze", English: "cat", Article: "die", Plural: "die Katzen", Example: "Die Katze schläft."},
	{German: "das Haus", English: "house", Article: "das", Plural: "die Häuser", Example: "Das Haus ist groß."},
	{German: "der Tisch", English: "table", Article: "der", Plural: "die Tische", Example: "Das Buch liegt auf dem Tisch."},
	{German: "die Schule", English: "school", Article: "die", Plural: "die Schulen", Example: "Die Kinder gehen zur Schule."},
	{German: "das Wasser", English: "water", Article: "das", Plural: "", Example: "Ich trinke Wasser."},
	{German: "der Freund", English: "friend", Article: "der", Plural: "die Freunde", Example: "Er ist mein bester Freund."},
	{German: "die Zeit", English: "time", Article: "die", Plural: "die Zeiten", Example: "Ich habe keine Zeit."},
	{German: "das Geld", English: "money", Article: "das", Plural: "", Example: "Ich habe kein Geld."},
	{German: "der Weg", English: "way/path", Article: "der", Plural: "die Wege", Example: "Der Weg ist lang."},
	{German: "die Stadt", English: "city", Article: "die", Plural: "die Städte", Example: "Berlin ist eine große Stadt."},
	{German: "das Kind", English: "child", Article: "das", Plural: "die Kinder", Example: "Das Kind spielt im Garten."},
	{German: "der Mann", English: "man", Article: "der", Plural: "die Männer", Example: "Der Mann liest die Zeitung."},
	{German: "die Frau", English: "woman", Article: "die", Plural: "die Frauen", Example: "Die Frau arbeitet im Büro."},
	{German: "das Auto", English: "car", Article: "das", Plural: "die Autos", Example: "Das Auto ist rot."},
	{German: "der Baum", English: "tree", Article: "der", Plural: "die Bäume", Example: "Der Baum ist sehr alt."},
	{German: "die Sonne", English: "sun", Article: "die", Plural: "", Example: "Die Sonne scheint heute."},
	{German: "das Brot", English: "bread", Article: "das", Plural: "die Brote", Example: "Ich kaufe Brot."},
	{German: "der Kaffee", English: "coffee", Article: "der", Plural: "", Example: "Ich trinke gern Kaffee."},
	{German: "die Musik", English: "music", Article: "die", Plural: "", Example: "Ich höre gern Musik."},
	{German: "das Zimmer", English: "room", Article: "das", Plural: "die Zimmer", Example: "Das Zimmer ist sauber."},
	{German: "der Arzt", English: "doctor", Article: "der", Plural: "die Ärzte", Example: "Ich gehe zum Arzt."},
	{German: "die Arbeit", English: "work/job", Article: "die", Plural: "die Arbeiten", Example: "Die Arbeit ist fertig."},
	{German: "das Essen", English: "food/meal", Article: "das", Plural: "die Essen", Example: "Das Essen schmeckt gut."},
	{German: "der Schlüssel", English: "key", Article: "der", Plural: "die Schlüssel", Example: "Wo ist mein Schlüssel?"},
	{German: "die Sprache", English: "language", Article: "die", Plural: "die Sprachen", Example: "Deutsch ist eine schöne Sprache."},
	{German: "das Bild", English: "picture", Article: "das", Plural: "die Bilder", Example: "Das Bild ist wunderschön."},
	{German: "der Stuhl", English: "chair", Article: "der", Plural: "die Stühle", Example: "Setz dich auf den Stuhl."},
	{German: "die Tür", English: "door", Article: "die", Plural: "die Türen", Example: "Bitte schließe die Tür."},
	{German: "das Fenster", English: "window", Article: "das", Plural: "die Fenster", Example: "Öffne das Fenster."},
	{German: "der Bruder", English: "brother", Article: "der", Plural: "die Brüder", Example: "Mein Bruder ist älter."},
	{German: "die Schwester", English: "sister", Article: "die", Plural: "die Schwestern", Example: "Meine Schwester wohnt in München."},
	{German: "das Kleid", English: "dress", Article: "das", Plural: "die Kleider", Example: "Das Kleid ist blau."},
	{German: "der Fluss", English: "river", Article: "der", Plural: "die Flüsse", Example: "Der Fluss ist sehr lang."},
	{German: "die Brücke", English: "bridge", Article: "die", Plural: "die Brücken", Example: "Die Brücke ist alt."},
	{German: "das Meer", English: "sea", Article: "das", Plural: "die Meere", Example: "Wir fahren ans Meer."},
	{German: "der Berg", English: "mountain", Article: "der", Plural: "die Berge", Example: "Der Berg ist sehr hoch."},
}

func GetWordOfTheDay() WordOfDay {
	dayOfYear := time.Now().YearDay()
	return wordsOfDay[dayOfYear%len(wordsOfDay)]
}
