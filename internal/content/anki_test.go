package content

import (
	"os"
	"strings"
	"testing"
)

func TestImportAnkiTSV(t *testing.T) {
	input := "#separator:tab\n#deck:A1\nid-1\tder Apfel\tapple\tPlural: die Aepfel\ta1 noun\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	if notes[0].DeckID != "A1" || notes[0].Front != "der Apfel" || notes[0].Back != "apple" {
		t.Fatalf("unexpected note: %+v", notes[0])
	}
}

func TestExportAnkiTSVRoundTrip(t *testing.T) {
	deck := StarterDeck()
	out, err := ExportAnkiTSVString(deck.Notes[:1])
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	notes, err := ImportAnkiTSV(strings.NewReader(out), ImportOptions{DefaultDeck: "fallback"})
	if err != nil {
		t.Fatalf("reimport failed: %v", err)
	}
	if notes[0].ID != deck.Notes[0].ID {
		t.Fatalf("id = %q, want %q", notes[0].ID, deck.Notes[0].ID)
	}
}

func TestImportAnkiFixture(t *testing.T) {
	raw, err := os.Open("testdata/anki/a1-basic.tsv")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer raw.Close()

	notes, err := ImportAnkiTSV(raw, ImportOptions{})
	if err != nil {
		t.Fatalf("import fixture: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("notes = %d, want 2", len(notes))
	}
	if notes[0].DeckID != "A1 Survival" {
		t.Fatalf("deck = %q, want A1 Survival", notes[0].DeckID)
	}
}
