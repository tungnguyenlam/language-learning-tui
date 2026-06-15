package tui

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\x1b[31mRed Text\x1b[0m", "Red Text"},
		{"Normal Text", "Normal Text"},
		{"\x1b[1mBold\x1b[22m and \x1b[3mItalic\x1b[23m", "Bold and Italic"},
	}

	for _, test := range tests {
		result := stripANSI(test.input)
		if result != test.expected {
			t.Errorf("stripANSI(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestTrimLastRune(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hell"},
		{"", ""},
		{"h", ""},
		{"über", "übe"},
		{"🚀", ""},
	}

	for _, test := range tests {
		result := trimLastRune(test.input)
		if result != test.expected {
			t.Errorf("trimLastRune(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestSinglePrintableInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"a", "a", true},
		{"A", "A", true},
		{" ", " ", true},
		{"ß", "ß", true},
		{"", "", false},
		{"ab", "", false},
		{"\n", "", false},
		{"\t", "", false},
	}

	for _, test := range tests {
		result, ok := singlePrintableInput(test.input)
		if result != test.expected || ok != test.ok {
			t.Errorf("singlePrintableInput(%q) = (%q, %t), expected (%q, %t)", test.input, result, ok, test.expected, test.ok)
		}
	}
}

func TestTruncateLine(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))

	tests := []struct {
		name     string
		input    string
		maxWidth int
		wantVis  int
	}{
		{"ascii exact boundary", "abcdef", 6, 6},
		{"ascii old-overflow case", "12345x", 5, 5},
		{"german umlauts", "Größe und mehr", 10, 10},
		{"wide cjk chars", "中文中文中文", 6, 5},
		{"ansi styled near boundary", style.Render("abc") + "def", 5, 5},
		{"ansi styled old-overflow case", style.Render("abcde") + "f", 5, 5},
		{"small width uses dots", "abcdef", 3, 3},
		{"fits without truncation", "abc", 5, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := truncateLine(tc.input, tc.maxWidth)
			gotVis := lipgloss.Width(got)
			if gotVis != tc.wantVis {
				t.Errorf("truncateLine(%q, %d) visual width = %d, want %d; got %q",
					tc.input, tc.maxWidth, gotVis, tc.wantVis, got)
			}
			if tc.maxWidth > 3 && lipgloss.Width(got) > tc.maxWidth {
				t.Errorf("truncateLine(%q, %d) exceeded maxWidth: got visual width %d",
					tc.input, tc.maxWidth, lipgloss.Width(got))
			}
		})
	}
}

func TestDictionaryHighlightAfterTruncation(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))

	// Simulate a dictionary list item that contains multi-byte German characters
	// and is truncated before highlighting is applied.
	wordText := "Größenordnung und mehr"
	padded := padString(wordText, 12)
	highlighted := highlightQuery(padded, "ße", style)

	// The highlight should be applied to the visible portion of the text.
	if !strings.Contains(highlighted, style.Render("ße")) {
		t.Errorf("expected highlight to contain styled 'ße'; got %q", highlighted)
	}

	// The rendered width must not exceed the requested panel width.
	if lipgloss.Width(highlighted) > 12 {
		t.Errorf("highlighted result width %d exceeds panel width 12: %q",
			lipgloss.Width(highlighted), highlighted)
	}
}

func TestFormatReviewInterval(t *testing.T) {
	tests := []struct {
		interval time.Duration
		expected string
	}{
		{0, "same day"},
		{-1 * time.Hour, "same day"},
		{30 * time.Minute, "1 hour"},
		{1 * time.Hour, "1 hour"},
		{5 * time.Hour, "5 hours"},
		{23 * time.Hour, "23 hours"},
		{24 * time.Hour, "1 day"},
		{7 * 24 * time.Hour, "7 days"},
		{29 * 24 * time.Hour, "29 days"},
		{30 * 24 * time.Hour, "1 month"},
		{45 * 24 * time.Hour, "1 month 15d"},
		{60 * 24 * time.Hour, "2 months"},
		{90 * 24 * time.Hour, "3 months"},
		{120 * 24 * time.Hour, "4 months"},
		{365 * 24 * time.Hour, "1 year"},
		{400 * 24 * time.Hour, "1 year 1 mo"},
		{730 * 24 * time.Hour, "2 years"},
		{800 * 24 * time.Hour, "2 years 2 mo"},
	}

	for _, test := range tests {
		result := formatReviewInterval(test.interval)
		if result != test.expected {
			t.Errorf("formatReviewInterval(%v) = %q, expected %q", test.interval, result, test.expected)
		}
	}
}
