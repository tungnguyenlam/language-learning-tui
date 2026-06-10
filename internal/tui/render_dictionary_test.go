package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"deutsch-tui/internal/core"
)

func TestHighlightQuery(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)

	tests := []struct {
		text     string
		query    string
		expected string
	}{
		{
			text:     "Deutsch",
			query:    "de",
			expected: style.Render("De") + "utsch",
		},
		{
			text:     "Apfelkuchen",
			query:    "kuch",
			expected: "Apfel" + style.Render("kuch") + "en",
		},
		{
			text:     "banana",
			query:    "an",
			expected: "b" + style.Render("an") + style.Render("an") + "a",
		},
		{
			text:     "banana",
			query:    "",
			expected: "banana",
		},
	}

	for _, tc := range tests {
		got := highlightQuery(tc.text, tc.query, style)
		if got != tc.expected {
			t.Errorf("highlightQuery(%q, %q) = %q, expected %q", tc.text, tc.query, got, tc.expected)
		}
	}
}

func TestRenderDictionary(t *testing.T) {
	m := NewModel(&mockRepo{}, &mockScheduler{})
	m.dictionaryResults = []core.DictionaryEntry{
		{
			ID:          "1",
			Word:        "Apfel",
			Translation: "Apple",
			WordClass:   "Noun",
			Gender:      "r",
			Forms:       "Apfel, Äpfel",
			Examples:    []string{"Ein roter Apfel."},
		},
		{
			ID:          "2",
			Word:        "Birne",
			Translation: "Pear",
		},
	}
	m.dictionarySearch = "Ap"
	m.dictionaryCursor = 0

	// Test wide layout (two columns)
	m.width = 100
	m.height = 40
	layout := m.activeViewContentLayout()
	view := m.renderDictionary(layout)
	plainView := stripANSI(view)

	if !strings.Contains(plainView, "Apfel") {
		t.Errorf("expected view to contain Apfel")
	}
	if !strings.Contains(plainView, "Apple") {
		t.Errorf("expected view to contain Apple")
	}
	if !strings.Contains(plainView, "Forms:") {
		t.Errorf("expected view to contain detail panel with Forms. Plain view:\n%s", plainView)
	}
	if !strings.Contains(plainView, "Äpfel") {
		t.Errorf("expected view to contain Äpfel")
	}

	// Test compact layout (single column)
	m.width = 60
	layout = m.activeViewContentLayout()
	view = m.renderDictionary(layout)
	plainView = stripANSI(view)

	if strings.Contains(plainView, "Forms:") {
		t.Errorf("expected compact view NOT to contain detail panel")
	}
}
