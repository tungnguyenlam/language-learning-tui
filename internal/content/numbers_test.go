package content

import (
	"testing"
)

func TestFormatGermanNumberThousands(t *testing.T) {
	tests := []struct {
		num      int
		expected string
	}{
		{1000, "tausend"},
		{2000, "zweitausend"},
		{2026, "zweitausendsechsundzwanzig"},
		{9999, "neuntausendneunhundertneunundneunzig"},
	}

	for _, tt := range tests {
		actual := formatGermanNumber(tt.num)
		if actual != tt.expected {
			t.Errorf("formatGermanNumber(%d) = %q; want %q", tt.num, actual, tt.expected)
		}
	}
}

func TestFormatGermanOrdinal(t *testing.T) {
	tests := []struct {
		num      int
		expected string
	}{
		{1, "erste"},
		{2, "zweite"},
		{3, "dritte"},
		{4, "vierte"},
		{7, "siebte"},
		{8, "achte"},
		{10, "zehnte"},
		{19, "neunzehnte"},
		{20, "zwanzigste"},
		{21, "einundzwanzigste"},
	}

	for _, tt := range tests {
		actual := formatGermanOrdinal(tt.num)
		if actual != tt.expected {
			t.Errorf("formatGermanOrdinal(%d) = %q; want %q", tt.num, actual, tt.expected)
		}
	}
}

func TestGetNumberExercises(t *testing.T) {
	exercises := GetNumberExercises()
	if len(exercises) == 0 {
		t.Fatalf("expected some exercises, got none")
	}

	hasOrdinal := false
	hasThousands := false

	for _, ex := range exercises {
		if ex.Help == "Ordinal number" {
			hasOrdinal = true
		}
		if ex.Help == "Thousands" {
			hasThousands = true
		}
	}

	if !hasOrdinal {
		t.Errorf("expected exercises to contain ordinal numbers")
	}
	if !hasThousands {
		t.Errorf("expected exercises to contain thousands")
	}

	// Verify we have enough time exercises (at least 30)
	timeCount := 0
	for _, ex := range exercises {
		if ex.Help == "Time expression" {
			timeCount++
		}
	}
	if timeCount != 30 {
		t.Errorf("expected 30 time expressions, got %d", timeCount)
	}
}
