package content

import (
	"strings"
	"testing"
)

func TestAnalyzeNounsWithArticle(t *testing.T) {
	cases := []struct {
		front        string
		wantArticle  string
		wantGender   string
		wantBase     string
		wantKind     WordKind
		wantInGender bool
	}{
		{"der Kaffee", "der", "masculine", "Kaffee", KindNoun, true},
		{"die Frau", "die", "feminine", "Frau", KindNoun, true},
		{"das Buch", "das", "neuter", "Buch", KindNoun, true},
		{"Die Bahn", "die", "feminine", "Bahn", KindNoun, true},
		{"DAS Telefon", "das", "neuter", "Telefon", KindNoun, true},
	}
	for _, c := range cases {
		got := Analyze(c.front)
		if got.Kind != c.wantKind {
			t.Errorf("%q kind = %v, want %v", c.front, got.Kind, c.wantKind)
		}
		if got.Article != c.wantArticle {
			t.Errorf("%q article = %q, want %q", c.front, got.Article, c.wantArticle)
		}
		if got.Gender != c.wantGender {
			t.Errorf("%q gender = %q, want %q", c.front, got.Gender, c.wantGender)
		}
		if got.Base != c.wantBase {
			t.Errorf("%q base = %q, want %q", c.front, got.Base, c.wantBase)
		}
	}
}

func TestAnalyzeVerbsAndAdjectives(t *testing.T) {
	cases := []struct {
		front string
		want  WordKind
	}{
		{"lernen", KindVerb},
		{"reisen", KindVerb},
		{"feiern", KindVerb},
		{"sammeln", KindVerb},
		{"sich freuen", KindVerb},
		{"groß", KindAdjective},
		{"schnell", KindAdjective},
		{"alt", KindAdjective},
		{"selten", KindAdjective},
		{"trocken", KindAdjective},
		{"offen", KindAdjective},
		{"zufrieden", KindAdjective},
		{"golden", KindAdjective},
		{"modern", KindAdjective},
		{"nüchtern", KindAdjective},
	}
	for _, c := range cases {
		got := Analyze(c.front)
		if got.Kind != c.want {
			t.Errorf("%q kind = %v, want %v", c.front, got.Kind, c.want)
		}
	}
}

// TestAnalyzeIrregularInfinitives guards against the regression where the only
// German infinitives not ending in -en/-n ("sein", "tun") were misclassified as
// adjectives and produced nonsense grammar hints (e.g. "comparative: sein+er").
func TestAnalyzeIrregularInfinitives(t *testing.T) {
	cases := []struct {
		front string
		want  WordKind
	}{
		{"sein", KindVerb},
		{"tun", KindVerb},
	}
	for _, c := range cases {
		got := Analyze(c.front)
		if got.Kind != c.want {
			t.Errorf("%q kind = %v, want %v", c.front, got.Kind, c.want)
		}
	}
}

// TestEnrichIrregularVerbs verifies the two irregular infinitives get a useful
// conjugation from the curated table rather than a guessed/empty stem.
func TestEnrichIrregularVerbs(t *testing.T) {
	for _, verb := range []string{"sein", "tun"} {
		info := Enrich(Analyze(verb))
		if info.Kind != KindVerb {
			t.Fatalf("%q kind = %v, want VERB", verb, info.Kind)
		}
		if len(info.Forms) == 0 {
			t.Errorf("%q: expected conjugation forms, got none", verb)
		}
		if info.Example == "" {
			t.Errorf("%q: expected a curated example", verb)
		}
	}
}

// TestAnalyzeUmlautNounKeepsGender ensures umlaut nouns are still detected as
// nouns (not adjectives) when they appear without an article.
func TestAnalyzeUmlautNounKeepsGender(t *testing.T) {
	got := Analyze("Äpfel")
	if got.Kind != KindNoun {
		t.Errorf("Äpfel kind = %v, want NOUN", got.Kind)
	}
}

// TestAnalyzeCardDetectsEnglishPromptWithToPrefix checks that an English prompt
// beginning with "to " (infinitive marker) is not mistaken for the German side.
func TestAnalyzeCardDetectsEnglishPromptWithToPrefix(t *testing.T) {
	info := AnalyzeCard("to drive", "fahren")
	if info.Kind != KindVerb || info.Base != "fahren" {
		t.Errorf("AnalyzeCard(to drive/fahren) = %+v, want VERB fahren", info)
	}
}

func TestAnalyzePhrases(t *testing.T) {
	got := Analyze("guten Morgen")
	if got.Kind != KindPhrase {
		t.Errorf("guten Morgen kind = %v, want PHRASE", got.Kind)
	}
}

func TestAnalyzeAlternateForms(t *testing.T) {
	got := Analyze("der Musiker / die Musikerin")
	if got.Kind != KindNoun {
		t.Errorf("alternate kind = %v, want NOUN", got.Kind)
	}
	if got.Article != "der" {
		t.Errorf("alternate article = %q, want der", got.Article)
	}
	if got.Base != "Musiker" {
		t.Errorf("alternate base = %q, want Musiker", got.Base)
	}
}

func TestEnrichNounAddsExampleAndCaseHint(t *testing.T) {
	info := Enrich(Analyze("der Hund"))
	if info.Example == "" {
		t.Fatal("expected example sentence for noun")
	}
	if !strings.Contains(info.Example, "Hund") {
		t.Errorf("example does not mention base word: %q", info.Example)
	}
	if !strings.Contains(info.Note, "der") {
		t.Errorf("note should mention article cases: %q", info.Note)
	}
}

func TestEnrichVerbUsesCuratedConjugationsWhenAvailable(t *testing.T) {
	// "lernen" is in dailyVerbs
	info := Enrich(Analyze("lernen"))
	if info.Kind != KindVerb {
		t.Fatalf("kind = %v, want VERB", info.Kind)
	}
	if len(info.Forms) < 2 {
		t.Errorf("expected forms from conjugation, got %v", info.Forms)
	}
	if info.Example == "" {
		t.Error("expected curated example for known verb")
	}
}

func TestEnrichVerbFallsBackToStemForUnknownVerb(t *testing.T) {
	info := Enrich(Analyze("xyzzylen"))
	if info.Kind != KindVerb {
		t.Fatalf("kind = %v, want VERB", info.Kind)
	}
	if len(info.Forms) == 0 {
		t.Error("expected synthesized forms for unknown verb")
	}
}

func TestEnrichAdjectiveAddsComparison(t *testing.T) {
	info := Enrich(Analyze("schnell"))
	if info.Kind != KindAdjective {
		t.Fatalf("kind = %v, want ADJ", info.Kind)
	}
	joined := strings.Join(info.Forms, " ")
	if !strings.Contains(joined, "schneller") {
		t.Errorf("expected comparative in forms, got %v", info.Forms)
	}
	if !strings.Contains(joined, "am schnellsten") {
		t.Errorf("expected superlative in forms, got %v", info.Forms)
	}
}

// TestEnrichEnAdjectivesNotFakeConjugations guards the regression where
// adjectives ending in -en/-ern were treated as infinitives and received
// nonsense present-tense forms (e.g. selten → "ich selt").
func TestEnrichEnAdjectivesNotFakeConjugations(t *testing.T) {
	for _, adj := range []string{"selten", "trocken", "offen", "zufrieden", "nüchtern"} {
		info := Enrich(Analyze(adj))
		if info.Kind != KindAdjective {
			t.Errorf("%q kind = %v, want ADJ", adj, info.Kind)
			continue
		}
		joined := strings.Join(info.Forms, " ")
		if strings.Contains(joined, "ich ") {
			t.Errorf("%q: got verb conjugation %v, want adjective comparison", adj, info.Forms)
		}
		if !strings.Contains(joined, "comparative:") {
			t.Errorf("%q: expected comparative in forms, got %v", adj, info.Forms)
		}
	}
}

func TestAnalyzeCardPicksGermanSide(t *testing.T) {
	// English prompt, German answer
	info := AnalyzeCard("the dog", "der Hund")
	if info.Kind != KindNoun || info.Article != "der" {
		t.Errorf("AnalyzeCard expected to pick 'der Hund'; got %+v", info)
	}
	// German prompt, English answer
	info = AnalyzeCard("der Hund", "the dog")
	if info.Kind != KindNoun || info.Article != "der" {
		t.Errorf("AnalyzeCard expected to pick 'der Hund'; got %+v", info)
	}
	// Verb in either direction
	info = AnalyzeCard("lernen", "to learn")
	if info.Kind != KindVerb {
		t.Errorf("AnalyzeCard expected VERB for lernen/to learn; got %+v", info)
	}
}

func TestEmptyInputReturnsEmpty(t *testing.T) {
	info := Analyze("")
	if info.Kind != KindUnknown {
		t.Errorf("empty front kind = %v, want Unknown", info.Kind)
	}
}
