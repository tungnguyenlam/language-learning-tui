package content

import (
	"deutsch-tui/internal/core"
)

func C1LinguisticsDeck() core.Deck {
	notes := []core.Note{
		{ID: "c1-ling-sprachwissenschaft", DeckID: "c1-linguistics", Front: "die Sprachwissenschaft", Back: "linguistics", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-sprachwissenschaftler", DeckID: "c1-linguistics", Front: "der Sprachwissenschaftler / die Sprachwissenschaftlerin", Back: "linguist", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-phonetik", DeckID: "c1-linguistics", Front: "die Phonetik", Back: "phonetics", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-phonologie", DeckID: "c1-linguistics", Front: "die Phonologie", Back: "phonology", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-phonem", DeckID: "c1-linguistics", Front: "das Phonem", Back: "phoneme", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-laut", DeckID: "c1-linguistics", Front: "der Laut", Back: "speech sound", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-silbe", DeckID: "c1-linguistics", Front: "die Silbe", Back: "syllable", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-betonung", DeckID: "c1-linguistics", Front: "die Betonung", Back: "stress / emphasis", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-aussprache", DeckID: "c1-linguistics", Front: "die Aussprache", Back: "pronunciation", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-vokal", DeckID: "c1-linguistics", Front: "der Vokal", Back: "vowel", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-konsonant", DeckID: "c1-linguistics", Front: "der Konsonant", Back: "consonant", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-umlaut", DeckID: "c1-linguistics", Front: "der Umlaut", Back: "umlaut (vowel mutation)", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-morphologie", DeckID: "c1-linguistics", Front: "die Morphologie", Back: "morphology", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-morphem", DeckID: "c1-linguistics", Front: "das Morphem", Back: "morpheme", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-wortstamm", DeckID: "c1-linguistics", Front: "der Wortstamm", Back: "word stem / root", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-endung", DeckID: "c1-linguistics", Front: "die Endung", Back: "ending / inflectional ending", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-vorsilbe", DeckID: "c1-linguistics", Front: "die Vorsilbe", Back: "prefix", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-nachsilbe", DeckID: "c1-linguistics", Front: "die Nachsilbe", Back: "suffix", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-flexion", DeckID: "c1-linguistics", Front: "die Flexion", Back: "inflection", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-deklination", DeckID: "c1-linguistics", Front: "die Deklination", Back: "declension", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-konjugation", DeckID: "c1-linguistics", Front: "die Konjugation", Back: "conjugation", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-syntax", DeckID: "c1-linguistics", Front: "die Syntax", Back: "syntax", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-satzbau", DeckID: "c1-linguistics", Front: "der Satzbau", Back: "sentence structure", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-satzglied", DeckID: "c1-linguistics", Front: "das Satzglied", Back: "sentence constituent", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-hauptsatz", DeckID: "c1-linguistics", Front: "der Hauptsatz", Back: "main clause", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-nebensatz", DeckID: "c1-linguistics", Front: "der Nebensatz", Back: "subordinate clause", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-wortstellung", DeckID: "c1-linguistics", Front: "die Wortstellung", Back: "word order", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-semantik", DeckID: "c1-linguistics", Front: "die Semantik", Back: "semantics", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-bedeutung", DeckID: "c1-linguistics", Front: "die Bedeutung", Back: "meaning", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-synonym", DeckID: "c1-linguistics", Front: "das Synonym", Back: "synonym", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-antonym", DeckID: "c1-linguistics", Front: "das Antonym", Back: "antonym", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-mehrdeutigkeit", DeckID: "c1-linguistics", Front: "die Mehrdeutigkeit", Back: "ambiguity", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-pragmatik", DeckID: "c1-linguistics", Front: "die Pragmatik", Back: "pragmatics", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-sprechakt", DeckID: "c1-linguistics", Front: "der Sprechakt", Back: "speech act", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-soziolinguistik", DeckID: "c1-linguistics", Front: "die Soziolinguistik", Back: "sociolinguistics", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-dialekt", DeckID: "c1-linguistics", Front: "der Dialekt", Back: "dialect", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-mundart", DeckID: "c1-linguistics", Front: "die Mundart", Back: "local dialect / vernacular", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-umgangssprache", DeckID: "c1-linguistics", Front: "die Umgangssprache", Back: "colloquial language", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-standardsprache", DeckID: "c1-linguistics", Front: "die Standardsprache", Back: "standard language", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-sprachwandel", DeckID: "c1-linguistics", Front: "der Sprachwandel", Back: "language change", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-spracherwerb", DeckID: "c1-linguistics", Front: "der Spracherwerb", Back: "language acquisition", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-muttersprache", DeckID: "c1-linguistics", Front: "die Muttersprache", Back: "native language", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-zweitsprache", DeckID: "c1-linguistics", Front: "die Zweitsprache", Back: "second language", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-lehnwort", DeckID: "c1-linguistics", Front: "das Lehnwort", Back: "loanword", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-fremdwort", DeckID: "c1-linguistics", Front: "das Fremdwort", Back: "foreign word", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-etymologie", DeckID: "c1-linguistics", Front: "die Etymologie", Back: "etymology", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-wortschatz", DeckID: "c1-linguistics", Front: "der Wortschatz", Back: "vocabulary", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-korpus", DeckID: "c1-linguistics", Front: "das Korpus", Back: "corpus (linguistic text collection)", Tags: []string{"c1", "linguistics"}},
		{ID: "c1-ling-ableiten", DeckID: "c1-linguistics", Front: "ableiten", Back: "to derive", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-beugen", DeckID: "c1-linguistics", Front: "beugen", Back: "to inflect (decline / conjugate)", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-aussprechen", DeckID: "c1-linguistics", Front: "aussprechen", Back: "to pronounce", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-bezeichnen", DeckID: "c1-linguistics", Front: "bezeichnen", Back: "to denote / to designate", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-erwerben", DeckID: "c1-linguistics", Front: "erwerben", Back: "to acquire", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-entlehnen", DeckID: "c1-linguistics", Front: "entlehnen", Back: "to borrow (a word from another language)", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-gliedern", DeckID: "c1-linguistics", Front: "gliedern", Back: "to structure / to segment", Tags: []string{"c1", "linguistics", "verb"}},
		{ID: "c1-ling-sprachlich", DeckID: "c1-linguistics", Front: "sprachlich", Back: "linguistic", Tags: []string{"c1", "linguistics", "adjective"}},
		{ID: "c1-ling-muendlich", DeckID: "c1-linguistics", Front: "mündlich", Back: "oral / spoken", Tags: []string{"c1", "linguistics", "adjective"}},
		{ID: "c1-ling-schriftlich", DeckID: "c1-linguistics", Front: "schriftlich", Back: "written", Tags: []string{"c1", "linguistics", "adjective"}},
		{ID: "c1-ling-unregelmaessig", DeckID: "c1-linguistics", Front: "unregelmäßig", Back: "irregular", Tags: []string{"c1", "linguistics", "adjective"}},
		{ID: "c1-ling-bedeutungsgleich", DeckID: "c1-linguistics", Front: "bedeutungsgleich", Back: "synonymous / identical in meaning", Tags: []string{"c1", "linguistics", "adjective"}},
	}

	for i := range notes {
		notes[i].Cards = CardsForNote(notes[i])
	}

	return core.Deck{
		ID:          "c1-linguistics",
		Name:        "C1 Sprachwissenschaft",
		Description: "Linguistics vocabulary: phonetics, morphology, syntax, semantics and sociolinguistics.",
		Tags:        []string{"german", "c1", "linguistics", "language"},
		Notes:       notes,
	}
}
