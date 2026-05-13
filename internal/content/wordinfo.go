package content

import (
	"fmt"
	"strings"
	"unicode"
)

type WordKind int

const (
	KindUnknown WordKind = iota
	KindNoun
	KindVerb
	KindAdjective
	KindPhrase
	KindOther
)

func (k WordKind) String() string {
	switch k {
	case KindNoun:
		return "NOUN"
	case KindVerb:
		return "VERB"
	case KindAdjective:
		return "ADJ"
	case KindPhrase:
		return "PHRASE"
	case KindOther:
		return "WORD"
	default:
		return ""
	}
}

// WordInfo is a compact linguistic profile for a German vocabulary entry,
// derived from the card's German side. Fields are best-effort and may be
// empty when the analyzer cannot infer them confidently.
type WordInfo struct {
	Kind    WordKind
	Article string   // "der" / "die" / "das" for nouns with explicit article
	Base    string   // the headword without article or "sich" prefix
	Display string   // the headword as the user should see it (article + base)
	Gender  string   // "masculine" / "feminine" / "neuter" for nouns
	Forms   []string // related forms: conjugations, comparatives, derived nouns
	Example string   // a usage example sentence
	Note    string   // a short pedagogical note
}

// Analyze inspects a German vocabulary string and returns its linguistic
// profile. The string can be a bare word, an article+noun, a reflexive
// verb ("sich freuen"), or a short phrase. The returned WordInfo always
// has Kind set; downstream code may further enrich the analysis.
func Analyze(front string) WordInfo {
	s := strings.TrimSpace(front)
	if s == "" {
		return WordInfo{}
	}

	if alt, ok := splitAlternates(s); ok {
		// e.g. "der Musiker / die Musikerin" — analyze first alternate
		return Analyze(alt)
	}

	lower := strings.ToLower(s)

	if rest, ok := stripArticle(lower, s, "der "); ok {
		return WordInfo{Kind: KindNoun, Article: "der", Base: rest, Display: "der " + rest, Gender: "masculine"}
	}
	if rest, ok := stripArticle(lower, s, "die "); ok {
		return WordInfo{Kind: KindNoun, Article: "die", Base: rest, Display: "die " + rest, Gender: "feminine"}
	}
	if rest, ok := stripArticle(lower, s, "das "); ok {
		return WordInfo{Kind: KindNoun, Article: "das", Base: rest, Display: "das " + rest, Gender: "neuter"}
	}
	if rest, ok := stripArticle(lower, s, "sich "); ok {
		return WordInfo{Kind: KindVerb, Base: "sich " + rest, Display: "sich " + rest}
	}

	runes := []rune(s)
	if len(runes) == 0 {
		return WordInfo{}
	}

	if strings.Contains(s, " ") {
		// Multi-word entry with no recognized article prefix is treated as a phrase.
		// One special case: capitalized two-word noun compounds like "die Frau Müller"
		// are uncommon in our decks, so phrase is the right default.
		return WordInfo{Kind: KindPhrase, Base: s, Display: s}
	}

	first := runes[0]
	if unicode.IsLower(first) {
		if isInfinitive(s) {
			return WordInfo{Kind: KindVerb, Base: s, Display: s}
		}
		return WordInfo{Kind: KindAdjective, Base: s, Display: s}
	}
	if unicode.IsUpper(first) {
		// Capitalized single word with no article: still a noun in German.
		return WordInfo{Kind: KindNoun, Base: s, Display: s}
	}

	return WordInfo{Kind: KindOther, Base: s, Display: s}
}

// AnalyzeCard picks whichever side of a card looks German and analyzes it.
func AnalyzeCard(prompt, answer string) WordInfo {
	pInfo := Analyze(prompt)
	aInfo := Analyze(answer)
	pScore := germanScore(prompt) + kindScore(pInfo.Kind)
	aScore := germanScore(answer) + kindScore(aInfo.Kind)
	if aScore > pScore {
		return Enrich(aInfo)
	}
	return Enrich(pInfo)
}

// Enrich augments a WordInfo with derived forms and an example sentence.
func Enrich(info WordInfo) WordInfo {
	switch info.Kind {
	case KindNoun:
		enrichNoun(&info)
	case KindVerb:
		enrichVerb(&info)
	case KindAdjective:
		enrichAdjective(&info)
	}
	return info
}

func enrichNoun(info *WordInfo) {
	if info.Base == "" {
		return
	}
	if info.Example == "" {
		switch info.Gender {
		case "masculine":
			info.Example = fmt.Sprintf("Das ist ein %s.", info.Base)
		case "feminine":
			info.Example = fmt.Sprintf("Das ist eine %s.", info.Base)
		case "neuter":
			info.Example = fmt.Sprintf("Das ist ein %s.", info.Base)
		default:
			info.Example = fmt.Sprintf("Das ist %s.", info.Base)
		}
	}
	if info.Note == "" {
		switch info.Gender {
		case "masculine":
			info.Note = "der → den (Akk), dem (Dat)"
		case "feminine":
			info.Note = "die → die (Akk), der (Dat)"
		case "neuter":
			info.Note = "das → das (Akk), dem (Dat)"
		}
	}
}

func enrichVerb(info *WordInfo) {
	if info.Base == "" {
		return
	}
	// Match against curated daily verb conjugations when available.
	for _, v := range dailyVerbs {
		if strings.EqualFold(v.German, info.Base) {
			info.Forms = []string{
				fmt.Sprintf("ich %s", v.Ich),
				fmt.Sprintf("er/sie/es %s", v.ErSieEs),
				fmt.Sprintf("wir %s", v.Wir),
			}
			if info.Example == "" {
				info.Example = v.Example
			}
			info.Note = "Infinitive form (use stem + endings for present tense)"
			return
		}
	}

	stem := guessStem(info.Base)
	if stem != "" {
		info.Forms = []string{
			fmt.Sprintf("ich %se", stem),
			fmt.Sprintf("er/sie/es %st", stem),
		}
	}
	if info.Example == "" && stem != "" {
		info.Example = fmt.Sprintf("Ich %se gern.", stem)
	}
	if info.Note == "" {
		info.Note = "Infinitive — drop -en for the stem"
	}
}

func enrichAdjective(info *WordInfo) {
	if info.Base == "" {
		return
	}
	base := info.Base
	if info.Forms == nil {
		info.Forms = []string{
			"comparative: " + base + "er",
			"superlative: am " + base + "sten",
		}
	}
	if info.Example == "" {
		info.Example = fmt.Sprintf("Es ist sehr %s.", base)
	}
	if info.Note == "" {
		info.Note = "Regular pattern; some adjectives take an umlaut (alt → älter)"
	}
}

func stripArticle(lower, original, prefix string) (string, bool) {
	if strings.HasPrefix(lower, prefix) {
		return strings.TrimSpace(original[len(prefix):]), true
	}
	return "", false
}

func splitAlternates(s string) (string, bool) {
	for _, sep := range []string{" / ", "/", " or ", ", "} {
		if i := strings.Index(s, sep); i > 0 {
			return strings.TrimSpace(s[:i]), true
		}
	}
	return "", false
}

func isInfinitive(s string) bool {
	return strings.HasSuffix(s, "en") || strings.HasSuffix(s, "ern") || strings.HasSuffix(s, "eln")
}

func guessStem(verb string) string {
	switch {
	case strings.HasSuffix(verb, "ern"):
		return strings.TrimSuffix(verb, "n")
	case strings.HasSuffix(verb, "eln"):
		return strings.TrimSuffix(verb, "n")
	case strings.HasSuffix(verb, "en"):
		return strings.TrimSuffix(verb, "en")
	}
	return ""
}

func germanScore(s string) int {
	score := 0
	lower := strings.ToLower(s)
	switch {
	case strings.HasPrefix(lower, "der "),
		strings.HasPrefix(lower, "die "),
		strings.HasPrefix(lower, "das "),
		strings.HasPrefix(lower, "sich "):
		score += 3
	}
	if strings.ContainsAny(s, "äöüÄÖÜß") {
		score += 2
	}
	if strings.HasPrefix(lower, "to ") || strings.HasPrefix(lower, "the ") {
		score -= 3
	}
	return score
}

func kindScore(k WordKind) int {
	switch k {
	case KindNoun, KindVerb, KindAdjective:
		return 1
	case KindPhrase:
		return 0
	}
	return 0
}
