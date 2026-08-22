package tui

import (
	"strings"
	"testing"

	"deutsch-tui/internal/core"

	"charm.land/lipgloss/v2"
)

// stripANSI is provided by utils_test.go (package-internal).

func TestNormalizeAnswerNonStrict(t *testing.T) {
	m := &Model{} // strictNormalization defaults to false
	cases := []struct {
		in   string
		want string
	}{
		{"Haus", "haus"},
		{"Über", "ueber"},
		{"Straße", "strasse"},
		{"Der Hund.", "der hund"},
		{"  ich  gehe  ", "ich gehe"},
		{"Apfel!", "apfel"},
	}
	for _, c := range cases {
		if got := m.normalizeAnswer(c.in); got != c.want {
			t.Errorf("normalizeAnswer(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeAnswerStrictKeepsUmlauts(t *testing.T) {
	m := &Model{strictNormalization: true}
	if got := m.normalizeAnswer("Über"); got != "über" {
		t.Errorf("strict normalizeAnswer(Über) = %q, want über", got)
	}
	if got := m.normalizeAnswer("Straße"); got != "straße" {
		t.Errorf("strict normalizeAnswer(Straße) = %q, want straße", got)
	}
}

func TestNormalizeAnswerEqualityAcceptsUmlautVariants(t *testing.T) {
	m := &Model{}
	typed := m.normalizeAnswer("Über")
	target := m.normalizeAnswer("ueber")
	if typed != target {
		t.Errorf("umlaut variants should normalize equal: %q vs %q", typed, target)
	}
}

func TestClozeAnswerText(t *testing.T) {
	// Single choice
	c := core.Card{Kind: core.CardKindCloze, Answer: "antwort", Choices: []string{"antwort"}}
	if got := clozeAnswerText(c); got != "antwort" {
		t.Errorf("single choice = %q, want antwort", got)
	}
	// Multiple choices joined
	c = core.Card{Kind: core.CardKindCloze, Choices: []string{"Hund", "Katze"}}
	if got := clozeAnswerText(c); got != "Hund / Katze" {
		t.Errorf("multi choice = %q, want 'Hund / Katze'", got)
	}
	// No choices falls back to Answer
	c = core.Card{Kind: core.CardKindCloze, Answer: "fallback"}
	if got := clozeAnswerText(c); got != "fallback" {
		t.Errorf("no choice = %q, want fallback", got)
	}
}

func TestRenderClozeAnswersFillsPlaceholders(t *testing.T) {
	prompt := "Ich gehe [hin] und sie kommt [her]."
	out := renderClozeAnswers(prompt, []string{"hin", "her"}, lipgloss.NewStyle())
	plain := stripANSI(out)
	if strings.Contains(plain, "[hin]") || strings.Contains(plain, "[her]") {
		t.Errorf("placeholders were not replaced: %q", plain)
	}
	if !strings.Contains(plain, "hin") || !strings.Contains(plain, "her") {
		t.Errorf("expected choice text in output, got %q", plain)
	}
}

func TestRenderTypingDiffMarksExactMatch(t *testing.T) {
	out := renderTypingDiff("Haus", "Haus")
	plain := stripANSI(out)
	if !strings.Contains(plain, "Haus") {
		t.Errorf("expected typed text in diff output, got %q", plain)
	}
}

func TestRenderTypingDiffShowsExpectedOnMismatch(t *testing.T) {
	// On a mismatch the diff must still render without panicking and surface
	// at least one character of the expected answer (here the shared 'a').
	out := renderTypingDiff("Katze", "Haus")
	plain := stripANSI(out)
	if plain == "" {
		t.Fatal("expected non-empty diff output on mismatch")
	}
	if !strings.Contains(plain, "a") {
		t.Errorf("expected some correct-answer characters in diff output, got %q", plain)
	}
}
