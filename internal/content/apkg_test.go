package content

import (
	"bytes"
	"testing"

	"deutsch-tui/internal/core"
)

func TestExportImportAPKG(t *testing.T) {
	// Create a test note
	note := core.Note{
		ID:    "test-note-1",
		Front: "Haus",
		Back:  "house",
		Tags:  []string{"noun", "german"},
		Extra: "Das ist ein Haus.",
	}
	notes := []core.Note{note}

	// Export to APKG
	var buf bytes.Buffer
	if err := ExportAnkiAPKG(&buf, notes); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check that we got some data
	if buf.Len() == 0 {
		t.Fatalf("Export produced empty output")
	}

	// Try to import it back
	importedNotes, err := ImportAnkiAPKG(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Check we got the same note back
	if len(importedNotes) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(importedNotes))
	}

	importedNote := importedNotes[0]
	if importedNote.ID != note.ID {
		t.Fatalf("Expected ID %s, got %s", note.ID, importedNote.ID)
	}
	if importedNote.Front != note.Front {
		t.Fatalf("Expected Front %s, got %s", note.Front, importedNote.Front)
	}
	if importedNote.Back != note.Back {
		t.Fatalf("Expected Back %s, got %s", note.Back, importedNote.Back)
	}
	if importedNote.Extra != note.Extra {
		t.Fatalf("Expected Extra %s, got %s", note.Extra, importedNote.Extra)
	}
	// Note: Tags might not be preserved exactly due to Anki format limitations
}
