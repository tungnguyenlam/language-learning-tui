package tui

import (
	"testing"
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
