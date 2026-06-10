package content

import (
	"bytes"
	"testing"
)

func TestParseDictCCStream_ENDE(t *testing.T) {
	input := `# EN-DE vocabulary database	compiled by dict.cc
# Date and time	2026-06-10 09:13
$500 suit	500-Dollar-Anzug {m}	noun	[cloth.]
% w/w	Gewichtsprozent {n}	noun
apple [fruit]	Apfel {m}	noun	[bot.]
`

	entries, err := ParseDictCCStream(bytes.NewBufferString(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// First entry
	if entries[0].Word != "500-Dollar-Anzug" {
		t.Errorf("expected word '500-Dollar-Anzug', got '%s'", entries[0].Word)
	}
	if entries[0].Translation != "$500 suit" {
		t.Errorf("expected translation '$500 suit', got '%s'", entries[0].Translation)
	}
	if entries[0].Gender != "m" {
		t.Errorf("expected gender 'm', got '%s'", entries[0].Gender)
	}
	if entries[0].WordClass != "noun" {
		t.Errorf("expected word class 'noun', got '%s'", entries[0].WordClass)
	}
	if len(entries[0].Tags) != 1 || entries[0].Tags[0] != "[cloth.]" {
		t.Errorf("expected tags [[cloth.]], got %v", entries[0].Tags)
	}

	// Second entry
	if entries[1].Word != "Gewichtsprozent" {
		t.Errorf("expected word 'Gewichtsprozent', got '%s'", entries[1].Word)
	}
	if entries[1].Gender != "n" {
		t.Errorf("expected gender 'n', got '%s'", entries[1].Gender)
	}

	// Third entry
	if entries[2].Word != "Apfel" {
		t.Errorf("expected word 'Apfel', got '%s'", entries[2].Word)
	}
	if entries[2].Gender != "m" {
		t.Errorf("expected gender 'm', got '%s'", entries[2].Gender)
	}
}

func TestParseDictCCStream_DEEN(t *testing.T) {
	input := `# DE-EN vocabulary database	compiled by dict.cc
# Date and time	2026-06-10 09:11
Apfel {m}	apple [fruit]	noun	[bot.]
Frau {f}	woman	noun
`

	entries, err := ParseDictCCStream(bytes.NewBufferString(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// First entry
	if entries[0].Word != "Apfel" {
		t.Errorf("expected word 'Apfel', got '%s'", entries[0].Word)
	}
	if entries[0].Translation != "apple [fruit]" {
		t.Errorf("expected translation 'apple [fruit]', got '%s'", entries[0].Translation)
	}
	if entries[0].Gender != "m" {
		t.Errorf("expected gender 'm', got '%s'", entries[0].Gender)
	}

	// Second entry
	if entries[1].Word != "Frau" {
		t.Errorf("expected word 'Frau', got '%s'", entries[1].Word)
	}
	if entries[1].Translation != "woman" {
		t.Errorf("expected translation 'woman', got '%s'", entries[1].Translation)
	}
	if entries[1].Gender != "f" {
		t.Errorf("expected gender 'f', got '%s'", entries[1].Gender)
	}
}
