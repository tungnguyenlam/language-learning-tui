package content

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"deutsch-tui/internal/core"

	"github.com/klauspost/compress/zstd"
)

// openExported writes notes to an .apkg and opens the collection inside it.
func openExported(t *testing.T, notes []core.Note, deckNames map[string]string) (*sql.DB, map[string][]byte) {
	t.Helper()

	var buf bytes.Buffer
	if err := ExportAnkiAPKGWithDeckNames(&buf, notes, deckNames); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("export is not a valid zip: %v", err)
	}

	entries := make(map[string][]byte)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("opening zip entry %s: %v", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("reading zip entry %s: %v", f.Name, err)
		}
		entries[f.Name] = data
	}

	collection, ok := entries["collection.anki2"]
	if !ok {
		t.Fatalf("export has no collection.anki2, entries: %v", keysOf(entries))
	}

	path := filepath.Join(t.TempDir(), "collection.anki2")
	if err := os.WriteFile(path, collection, 0o600); err != nil {
		t.Fatalf("staging collection: %v", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("opening collection: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db, entries
}

func keysOf[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// An exported package has to satisfy the invariants Anki checks on import.
// Before this was rewritten the exporter emitted notes whose `mid` matched no
// note type, plus a dummy "genesis" note, so real Anki rejected the file even
// though our own importer round-tripped it.
func TestExportProducesAnkiValidCollection(t *testing.T) {
	notes := []core.Note{
		{ID: "n1", DeckID: "a1_food", Front: "das Brot", Back: "bread", Extra: "Ich esse Brot."},
		{ID: "n2", DeckID: "a1_food", Front: "der Apfel", Back: "apple", Type: "Reverse"},
		{ID: "n3", DeckID: "b1_verbs", Front: "Ich {{c1::gehe}} nach Hause.", Back: "I go home."},
	}
	db, entries := openExported(t, notes, map[string]string{
		"a1_food":  "A1 Food & Drink",
		"b1_verbs": "B1 Verbs",
	})
	ctx := context.Background()

	var ver int
	var modelsJSON, decksJSON, confJSON string
	if err := db.QueryRowContext(ctx,
		`SELECT ver, models, decks, conf FROM col WHERE id = 1`).
		Scan(&ver, &modelsJSON, &decksJSON, &confJSON); err != nil {
		t.Fatalf("reading col row: %v", err)
	}
	if ver != 11 {
		t.Errorf("expected schema version 11, got %d", ver)
	}

	var models map[string]struct {
		ID        int64 `json:"id"`
		Templates []any `json:"tmpls"`
		Fields    []any `json:"flds"`
	}
	if err := json.Unmarshal([]byte(modelsJSON), &models); err != nil {
		t.Fatalf("col.models is not valid JSON: %v", err)
	}
	if len(models) == 0 {
		t.Fatal("col.models is empty; Anki cannot render any note")
	}
	for key, model := range models {
		if strconv.FormatInt(model.ID, 10) != key {
			t.Errorf("model key %q does not match its id %d", key, model.ID)
		}
		if len(model.Templates) == 0 || len(model.Fields) == 0 {
			t.Errorf("model %q has no templates or fields", key)
		}
	}

	// Every note must reference a declared note type and carry real content.
	rows, err := db.QueryContext(ctx, `SELECT id, guid, mid, flds FROM notes`)
	if err != nil {
		t.Fatalf("querying notes: %v", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id, mid int64
		var guid, flds string
		if err := rows.Scan(&id, &guid, &mid, &flds); err != nil {
			t.Fatalf("scanning note: %v", err)
		}
		count++
		if id <= 0 {
			t.Errorf("note %q has non-positive id %d", guid, id)
		}
		if _, ok := models[strconv.FormatInt(mid, 10)]; !ok {
			t.Errorf("note %q references unknown note type %d", guid, mid)
		}
		if guid == "genesis" || strings.TrimSpace(strings.ReplaceAll(flds, "\x1f", "")) == "" {
			t.Errorf("export contains a placeholder note (guid %q)", guid)
		}
	}
	if count != len(notes) {
		t.Errorf("expected %d notes, got %d", len(notes), count)
	}

	// Every card must point at a declared deck.
	var decks map[string]struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(decksJSON), &decks); err != nil {
		t.Fatalf("col.decks is not valid JSON: %v", err)
	}
	deckIDs := make(map[int64]string, len(decks))
	for _, d := range decks {
		deckIDs[d.ID] = d.Name
	}

	cardRows, err := db.QueryContext(ctx, `SELECT nid, did FROM cards`)
	if err != nil {
		t.Fatalf("querying cards: %v", err)
	}
	defer cardRows.Close()
	cards := 0
	for cardRows.Next() {
		var nid, did int64
		if err := cardRows.Scan(&nid, &did); err != nil {
			t.Fatalf("scanning card: %v", err)
		}
		cards++
		if _, ok := deckIDs[did]; !ok {
			t.Errorf("card for note %d points at undeclared deck %d", nid, did)
		}
	}
	if cards == 0 {
		t.Fatal("export produced no cards")
	}

	// Anki reads a file literally named "media"; "media.json" is ignored.
	if _, ok := entries["media"]; !ok {
		t.Errorf("export has no media index, entries: %v", keysOf(entries))
	}

	if !strings.Contains(confJSON, "curDeck") {
		t.Errorf("col.conf is missing curDeck: %s", confJSON)
	}
}

func TestExportPreservesDeckNames(t *testing.T) {
	notes := []core.Note{
		{ID: "n1", DeckID: "a1_food", Front: "das Brot", Back: "bread"},
		{ID: "n2", DeckID: "b1_verbs", Front: "gehen", Back: "to go"},
	}
	db, _ := openExported(t, notes, map[string]string{
		"a1_food":  "A1 Food & Drink",
		"b1_verbs": "B1 Verbs",
	})

	var decksJSON string
	if err := db.QueryRow(`SELECT decks FROM col WHERE id = 1`).Scan(&decksJSON); err != nil {
		t.Fatalf("reading decks: %v", err)
	}
	var decks map[string]struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(decksJSON), &decks); err != nil {
		t.Fatalf("col.decks is not valid JSON: %v", err)
	}
	names := make(map[string]bool, len(decks))
	for _, d := range decks {
		names[d.Name] = true
	}
	for _, want := range []string{"A1 Food & Drink", "B1 Verbs"} {
		if !names[want] {
			t.Errorf("expected deck %q, got %v", want, keysOf(names))
		}
	}
}

func TestRoundTripPreservesDecksTypesAndTags(t *testing.T) {
	notes := []core.Note{
		{ID: "n1", DeckID: "a1_food", Front: "das Brot", Back: "bread",
			Extra: "Ich esse Brot.", Tags: []string{"noun", "a1"}},
		{ID: "n2", DeckID: "b1_verbs", Front: "der Apfel", Back: "apple", Type: "Reverse"},
	}

	var buf bytes.Buffer
	if err := ExportAnkiAPKGWithDeckNames(&buf, notes, map[string]string{
		"a1_food":  "A1 Food",
		"b1_verbs": "B1 Verbs",
	}); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	got, err := ImportAnkiAPKG(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(got))
	}

	byID := make(map[string]core.Note, len(got))
	for _, n := range got {
		byID[n.ID] = n
	}

	brot, ok := byID["n1"]
	if !ok {
		t.Fatalf("note n1 missing, got %v", keysOf(byID))
	}
	if brot.DeckID != "A1 Food" {
		t.Errorf("expected deck name to survive the round trip, got %q", brot.DeckID)
	}
	if brot.Front != "das Brot" || brot.Back != "bread" || brot.Extra != "Ich esse Brot." {
		t.Errorf("fields not preserved: %+v", brot)
	}
	if strings.Join(brot.Tags, ",") != "noun,a1" {
		t.Errorf("expected tags to survive, got %v", brot.Tags)
	}

	apfel := byID["n2"]
	if apfel.Type != "Reverse" {
		t.Errorf("expected the reversed note type to survive, got %q", apfel.Type)
	}
	if len(apfel.Cards) != 2 {
		t.Errorf("expected a reversed note to produce 2 cards, got %d", len(apfel.Cards))
	}
}

// Anki fields are HTML. Without stripping, cards render as raw markup in a
// terminal and audio references leak into the prompt text.
func TestImportConvertsHTMLFields(t *testing.T) {
	for _, tc := range []struct {
		name      string
		field     string
		wantText  string
		wantAudio string
	}{
		{"line breaks", "erste Zeile<br>zweite Zeile", "erste Zeile\nzweite Zeile", ""},
		{"div blocks", "<div>eins</div><div>zwei</div>", "eins\nzwei", ""},
		{"entities", "Gr&uuml;&szlig;e &amp; mehr", "Grüße & mehr", ""},
		{"nbsp", "der&nbsp;Hund", "der Hund", ""},
		{"inline markup", "<b>fett</b> und <i>kursiv</i>", "fett und kursiv", ""},
		{"sound reference", "das Brot[sound:brot.mp3]", "das Brot", "brot.mp3"},
		{"image tag", `ein Bild <img src="x.png">`, "ein Bild", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text, audio := stripHTMLFieldWithMedia(tc.field)
			if text != tc.wantText {
				t.Errorf("text: want %q, got %q", tc.wantText, text)
			}
			if audio != tc.wantAudio {
				t.Errorf("audio: want %q, got %q", tc.wantAudio, audio)
			}
		})
	}
}

// Modern Anki (2.1.50+) exports collection.anki21b, a zstd-compressed database.
// Shared decks downloaded from AnkiWeb are increasingly in this format.
func TestImportReadsZstdCollection(t *testing.T) {
	notes := []core.Note{{ID: "n1", DeckID: "deck", Front: "das Brot", Back: "bread"}}

	var legacy bytes.Buffer
	if err := ExportAnkiAPKG(&legacy, notes); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	collection := extractEntry(t, legacy.Bytes(), "collection.anki2")

	enc, err := zstd.NewWriter(nil)
	if err != nil {
		t.Fatalf("creating zstd writer: %v", err)
	}
	compressed := enc.EncodeAll(collection, nil)
	enc.Close()

	var pkg bytes.Buffer
	zw := zip.NewWriter(&pkg)
	writeEntry(t, zw, "collection.anki21b", compressed)
	writeEntry(t, zw, "meta", []byte{0x08, 0x03}) // version marker Anki writes
	if err := zw.Close(); err != nil {
		t.Fatalf("closing zip: %v", err)
	}

	got, err := ImportAnkiAPKG(bytes.NewReader(pkg.Bytes()))
	if err != nil {
		t.Fatalf("import of anki21b package failed: %v", err)
	}
	if len(got) != 1 || got[0].Front != "das Brot" {
		t.Fatalf("unexpected notes from anki21b package: %+v", got)
	}
}

// Anki 2.1 packages carry a real collection.anki21 next to a legacy
// collection.anki2 stub kept only for Anki 2.0. Reading the stub first is why
// such decks used to import as empty.
func TestImportPrefersNewerCollectionEntry(t *testing.T) {
	real := []core.Note{{ID: "n1", DeckID: "deck", Front: "die Katze", Back: "cat"}}
	stub := []core.Note{{ID: "old", DeckID: "deck", Front: "STUB", Back: "STUB"}}

	var realBuf, stubBuf bytes.Buffer
	if err := ExportAnkiAPKG(&realBuf, real); err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if err := ExportAnkiAPKG(&stubBuf, stub); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	var pkg bytes.Buffer
	zw := zip.NewWriter(&pkg)
	writeEntry(t, zw, "collection.anki2", extractEntry(t, stubBuf.Bytes(), "collection.anki2"))
	writeEntry(t, zw, "collection.anki21", extractEntry(t, realBuf.Bytes(), "collection.anki2"))
	if err := zw.Close(); err != nil {
		t.Fatalf("closing zip: %v", err)
	}

	got, err := ImportAnkiAPKG(bytes.NewReader(pkg.Bytes()))
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(got) != 1 || got[0].Front != "die Katze" {
		t.Fatalf("expected the anki21 collection to win, got %+v", got)
	}
}

func TestImportRejectsNonPackage(t *testing.T) {
	if _, err := ImportAnkiAPKG(strings.NewReader("this is not a zip")); err == nil {
		t.Fatal("expected an error for a non-package file")
	}

	var pkg bytes.Buffer
	zw := zip.NewWriter(&pkg)
	writeEntry(t, zw, "readme.txt", []byte("nothing here"))
	if err := zw.Close(); err != nil {
		t.Fatalf("closing zip: %v", err)
	}
	_, err := ImportAnkiAPKG(bytes.NewReader(pkg.Bytes()))
	if err == nil || !strings.Contains(err.Error(), "no Anki collection") {
		t.Fatalf("expected a clear 'no collection' error, got %v", err)
	}
}

// Anki note types are free-form. A real AnkiWeb deck ("Deutsch: 4000 German
// Words by Frequency") leaves field 1 empty and puts the meaning in a later
// field, so assuming fields 0/1 are the two sides yields a blank answer.
func TestFieldRolesFollowTemplates(t *testing.T) {
	nt := ankiNoteTypeInfo{fields: []string{"Word", "Unused", "Meaning", "Example"}, frontOrd: -1, backOrd: -1}
	nt.frontOrd, nt.backOrd = templateFieldOrds(nt.fields,
		"{{Meaning}}", "{{FrontSide}}<hr id=answer>{{Word}}<br>{{Example}}")

	if nt.frontOrd != 2 || nt.backOrd != 0 {
		t.Fatalf("template mapping: want front=2 back=0, got front=%d back=%d", nt.frontOrd, nt.backOrd)
	}

	texts := []string{"ein", "", "1) a; 2) one", "Ein Mann."}
	front, back, extras := assignFieldRoles(texts, nt)
	if texts[front] != "1) a; 2) one" {
		t.Errorf("front: got %q", texts[front])
	}
	if back < 0 || texts[back] != "ein" {
		t.Errorf("back: got %d", back)
	}
	if len(extras) != 1 || texts[extras[0]] != "Ein Mann." {
		t.Errorf("extras: got %v", extras)
	}
}

func TestFieldRolesFallBackToNonEmptyFields(t *testing.T) {
	// No usable template information: skip empty fields rather than emitting a
	// card whose answer is blank.
	unknown := ankiNoteTypeInfo{frontOrd: -1, backOrd: -1}

	front, back, extras := assignFieldRoles([]string{"", "das Brot", "", "bread", "Beispiel"}, unknown)
	if front != 1 || back != 3 {
		t.Fatalf("want front=1 back=3, got front=%d back=%d", front, back)
	}
	if len(extras) != 1 || extras[0] != 4 {
		t.Errorf("want extras=[4], got %v", extras)
	}

	// A note with nothing in it is skipped, not turned into a blank card.
	if front, _, _ := assignFieldRoles([]string{"", "  ", ""}, unknown); front != -1 {
		t.Errorf("expected an empty note to be skipped, got front=%d", front)
	}

	// A single-field note is still usable as a one-sided card.
	if front, back, _ := assignFieldRoles([]string{"nur vorne"}, unknown); front != 0 || back != -1 {
		t.Errorf("want front=0 back=-1, got front=%d back=%d", front, back)
	}
}

func TestTemplateFieldOrdsIgnoresFiltersAndSpecials(t *testing.T) {
	fields := []string{"Text", "Extra"}
	front, back := templateFieldOrds(fields, "{{cloze:Text}}", "{{cloze:Text}}<br>{{Extra}}")
	if front != 0 || back != 1 {
		t.Errorf("cloze template: want front=0 back=1, got %d/%d", front, back)
	}

	// {{FrontSide}} and {{#Tags}} are not fields and must not be picked up.
	front, back = templateFieldOrds([]string{"Front", "Back"},
		"{{#Front}}{{type:Front}}{{/Front}}", "{{FrontSide}}{{Back}}")
	if front != 0 || back != 1 {
		t.Errorf("filtered template: want front=0 back=1, got %d/%d", front, back)
	}
}

func TestClozeOrdinals(t *testing.T) {
	for _, tc := range []struct {
		text string
		want []int
	}{
		{"Ich {{c1::gehe}} nach Hause.", []int{0}},
		{"{{c2::A}} und {{c1::B}}", []int{0, 1}},
		{"{{c3::nur eine}}", []int{2}},
		{"kein Cloze", nil},
	} {
		got := clozeOrdinals(tc.text)
		if len(got) != len(tc.want) {
			t.Errorf("%q: want %v, got %v", tc.text, tc.want, got)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("%q: want %v, got %v", tc.text, tc.want, got)
				break
			}
		}
	}
}

func extractEntry(t *testing.T, pkg []byte, name string) []byte {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(pkg), int64(len(pkg)))
	if err != nil {
		t.Fatalf("reading zip: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("opening %s: %v", name, err)
		}
		defer rc.Close()
		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		return data
	}
	t.Fatalf("entry %s not found", name)
	return nil
}

func writeEntry(t *testing.T, zw *zip.Writer, name string, data []byte) {
	t.Helper()
	w, err := zw.Create(name)
	if err != nil {
		t.Fatalf("creating %s: %v", name, err)
	}
	if _, err := w.Write(data); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}
}
