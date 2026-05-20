package tui

import (
	"testing"
	"time"
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
