package content

import (
	"strings"
	"testing"
)

func TestDecomposeCompound(t *testing.T) {
	// Simulated dictionary for validation
	dict := map[string]bool{
		"hand": true, "schuh": true, "haus": true, "aufgabe": true,
		"tier": true, "kinder": true, "garten": true, "kind": true,
		"arbeit": true, "tag": true, "schlaf": true, "zimmer": true,
		"tür": true, "flug": true, "zeug": true, "hafen": true,
		"buch": true, "regal": true, "schrank": true, "kühl": true,
		"wort": true, "buch+": true, "wörter": true,
	}
	validate := func(word string) bool {
		return dict[strings.ToLower(word)]
	}

	tests := []struct {
		name      string
		word      string
		wantNil   bool
		wantLeft  string
		wantRight string
	}{
		{
			name:      "Handschuh splits into Hand + Schuh",
			word:      "Handschuh",
			wantLeft:  "Hand",
			wantRight: "Schuh",
		},
		{
			name:      "Hausaufgabe splits",
			word:      "Hausaufgabe",
			wantLeft:  "Haus",
			wantRight: "Aufgabe",
		},
		{
			name:      "Haustier splits",
			word:      "Haustier",
			wantLeft:  "Haus",
			wantRight: "Tier",
		},
		{
			name:      "Arbeitstag splits with Fugen-s",
			word:      "Arbeitstag",
			wantLeft:  "Arbeit",
			wantRight: "Tag",
		},
		{
			name:      "Schlafzimmer splits",
			word:      "Schlafzimmer",
			wantLeft:  "Schlaf",
			wantRight: "Zimmer",
		},
		{
			name:      "Kindergarten splits",
			word:      "Kindergarten",
			wantLeft:  "Kinder",
			wantRight: "Garten",
		},
		{
			name:    "short word returns nil",
			word:    "Haus",
			wantNil: true,
		},
		{
			name:    "empty word returns nil",
			word:    "",
			wantNil: true,
		},
		{
			name:      "word with article strips it",
			word:      "der Handschuh",
			wantLeft:  "Hand",
			wantRight: "Schuh",
		},
		{
			name:      "Kühlschrank splits",
			word:      "Kühlschrank",
			wantLeft:  "Kühl",
			wantRight: "Schrank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := DecomposeCompound(tt.word, validate)
			if tt.wantNil {
				if parts != nil {
					t.Errorf("DecomposeCompound(%q) = %v, want nil", tt.word, parts)
				}
				return
			}
			if parts == nil {
				t.Fatalf("DecomposeCompound(%q) = nil, want parts", tt.word)
			}
			if len(parts) != 2 {
				t.Fatalf("DecomposeCompound(%q) returned %d parts, want 2", tt.word, len(parts))
			}
			if parts[0].Word != tt.wantLeft {
				t.Errorf("left part = %q, want %q", parts[0].Word, tt.wantLeft)
			}
			if parts[1].Word != tt.wantRight {
				t.Errorf("right part = %q, want %q", parts[1].Word, tt.wantRight)
			}
		})
	}
}

func TestDecomposeCompoundWithoutValidator(t *testing.T) {
	// Without a validator, should still return something for long compound words
	parts := DecomposeCompound("Handschuh", nil)
	if parts == nil {
		t.Fatal("DecomposeCompound('Handschuh', nil) = nil, want parts")
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	// Without validation, any split is possible
	t.Logf("Split: %q + %q", parts[0].Word, parts[1].Word)
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hand", "Hand"},
		{"schuh", "Schuh"},
		{"", ""},
		{"A", "A"},
		{"über", "Über"},
	}
	for _, tt := range tests {
		got := capitalizeFirst(tt.input)
		if got != tt.want {
			t.Errorf("capitalizeFirst(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
