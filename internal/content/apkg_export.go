package content

import (
	"archive/zip"
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	_ "modernc.org/sqlite"
)

// ExportAnkiAPKG writes notes as an Anki `.apkg` package.
//
// Deck names default to each note's DeckID. Use ExportAnkiAPKGWithDeckNames to
// give the decks their display names.
func ExportAnkiAPKG(w io.Writer, notes []core.Note) error {
	return ExportAnkiAPKGWithDeckNames(w, notes, nil)
}

// ExportAnkiAPKGWithDeckNames writes notes as an Anki `.apkg` package, mapping
// each note's DeckID through deckNames for the deck title Anki will show.
//
// The package uses Anki's legacy schema-11 collection, which every Anki 2.1+
// release reads. Scheduling state is not exported: cards arrive as new, which
// is what a shared deck should do.
func ExportAnkiAPKGWithDeckNames(w io.Writer, notes []core.Note, deckNames map[string]string) error {
	tmpDir, err := os.MkdirTemp("", "deutsch-tui-apkg-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "collection.anki2")
	media, err := writeAnkiCollection(dbPath, notes, deckNames)
	if err != nil {
		return fmt.Errorf("building collection: %w", err)
	}

	zw := zip.NewWriter(w)

	if err := copyIntoZip(zw, "collection.anki2", dbPath); err != nil {
		return err
	}

	// Anki reads a file literally named "media": a JSON object mapping the
	// numeric name each blob has inside the zip to its real filename.
	index := make(map[string]string, len(media))
	for i, m := range media {
		name := strconv.Itoa(i)
		index[name] = m.name
		if err := copyIntoZip(zw, name, m.path); err != nil {
			return err
		}
	}
	indexJSON, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("encoding media index: %w", err)
	}
	entry, err := zw.Create("media")
	if err != nil {
		return fmt.Errorf("creating media index: %w", err)
	}
	if _, err := entry.Write(indexJSON); err != nil {
		return fmt.Errorf("writing media index: %w", err)
	}

	return zw.Close()
}

type mediaRef struct {
	name string // filename recorded in the media index
	path string // source file on disk
}

func writeAnkiCollection(dbPath string, notes []core.Note, deckNames map[string]string) ([]mediaRef, error) {
	os.Remove(dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for _, stmt := range strings.Split(ankiSchema11, ";") {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return nil, fmt.Errorf("creating schema: %w", err)
		}
	}

	now := time.Now()
	modSec := now.Unix()
	modMilli := now.UnixMilli()

	decks, deckIDFor := buildAnkiDecks(notes, deckNames, modSec)
	models := ankiModels(modSec, decks[0].ID)

	conf, err := json.Marshal(ankiCollectionConfig(decks[0].ID))
	if err != nil {
		return nil, err
	}
	modelsJSON, err := json.Marshal(models)
	if err != nil {
		return nil, err
	}
	decksMap := make(map[string]ankiDeck, len(decks))
	for _, d := range decks {
		decksMap[itoa(d.ID)] = d
	}
	decksJSON, err := json.Marshal(decksMap)
	if err != nil {
		return nil, err
	}
	dconfJSON, err := json.Marshal(ankiDeckConfig(modSec))
	if err != nil {
		return nil, err
	}

	// crt is the collection's day-rollover anchor: 4am local time today.
	crt := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location()).Unix()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO col (id, crt, mod, scm, ver, dty, usn, ls, conf, models, decks, dconf, tags)
		VALUES (1, ?, ?, ?, 11, 0, 0, 0, ?, ?, ?, ?, '{}')
	`, crt, modMilli, modMilli, string(conf), string(modelsJSON), string(decksJSON), string(dconfJSON)); err != nil {
		return nil, fmt.Errorf("inserting collection row: %w", err)
	}

	// Anki requires unique positive ids and conventionally uses millisecond
	// timestamps, so allocate a dense run of them rather than hashing.
	nextNoteID := modMilli
	nextCardID := modMilli
	var media []mediaRef
	seenMedia := make(map[string]int)
	position := 0

	for _, note := range notes {
		model, fields := ankiNoteFields(note)
		noteID := nextNoteID
		nextNoteID++

		sortField := fields[0]
		sum := sha1.Sum([]byte(stripHTMLField(sortField)))
		csum, _ := strconv.ParseInt(fmt.Sprintf("%x", sum)[:8], 16, 64)

		tags := ""
		if len(note.Tags) > 0 {
			tags = " " + strings.Join(note.Tags, " ") + " "
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO notes (id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data)
			VALUES (?, ?, ?, ?, -1, ?, ?, ?, ?, 0, '')
		`, noteID, ankiGUID(note.ID), model.ID, modSec, tags,
			strings.Join(fields, "\x1f"), sortField, csum); err != nil {
			return nil, fmt.Errorf("inserting note %s: %w", note.ID, err)
		}

		if ref, ok := ankiMediaRef(note.Audio); ok {
			if _, dup := seenMedia[ref.path]; !dup {
				seenMedia[ref.path] = len(media)
				media = append(media, ref)
			}
		}

		did := deckIDFor(note.DeckID)
		for _, ord := range ankiCardOrdinals(note, model) {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO cards (id, nid, did, ord, mod, usn, type, queue, due,
				                   ivl, factor, reps, lapses, left, odue, odid, flags, data)
				VALUES (?, ?, ?, ?, ?, -1, 0, 0, ?, 0, 0, 0, 0, 0, 0, 0, 0, '')
			`, nextCardID, noteID, did, ord, modSec, position); err != nil {
				return nil, fmt.Errorf("inserting card for note %s: %w", note.ID, err)
			}
			nextCardID++
			position++
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing collection: %w", err)
	}
	return media, nil
}

// buildAnkiDecks assigns a stable Anki deck id to every DeckID present in the
// notes. The returned slice always has at least one deck, because a collection
// with no decks has nowhere to put cards.
func buildAnkiDecks(notes []core.Note, deckNames map[string]string, mod int64) ([]ankiDeck, func(string) int64) {
	ids := make(map[string]int64)
	var order []string
	for _, note := range notes {
		key := strings.TrimSpace(note.DeckID)
		if _, ok := ids[key]; ok {
			continue
		}
		ids[key] = 0
		order = append(order, key)
	}
	sort.Strings(order)
	if len(order) == 0 {
		order = []string{""}
		ids[""] = 0
	}

	decks := make([]ankiDeck, 0, len(order))
	next := int64(1)
	for _, key := range order {
		name := deckNames[key]
		if name == "" {
			name = key
		}
		if strings.TrimSpace(name) == "" {
			name = "Deutsch-TUI"
		}
		// "::" is Anki's sub-deck separator; a raw "::" in a generated name
		// would silently create empty parent decks.
		name = strings.ReplaceAll(name, "\x1f", " ")
		ids[key] = next
		decks = append(decks, ankiDeck{
			ID:        next,
			Name:      name,
			Mod:       mod,
			USN:       -1,
			Conf:      1,
			ExtendNew: 10,
			ExtendRev: 50,
			LrnToday:  []int{0, 0},
			RevToday:  []int{0, 0},
			NewToday:  []int{0, 0},
			TimeToday: []int{0, 0},
		})
		next++
	}

	return decks, func(deckID string) int64 {
		if id, ok := ids[strings.TrimSpace(deckID)]; ok {
			return id
		}
		return decks[0].ID
	}
}

// ankiNoteFields picks the note type for a note and lays out its fields in that
// type's order.
func ankiNoteFields(note core.Note) (ankiModel, []string) {
	models := ankiModels(0, 1)
	if len(parseClozes(note.Front)) > 0 {
		return models[itoa(modelIDCloze)], []string{note.Front, note.Extra}
	}
	if note.Type == "Reverse" {
		return models[itoa(modelIDReverse)], []string{note.Front, note.Back, note.Extra}
	}
	return models[itoa(modelIDBasic)], []string{note.Front, note.Back, note.Extra}
}

// ankiCardOrdinals lists the template ordinals Anki should generate cards for.
// For cloze note types the ordinal is the cloze number minus one.
func ankiCardOrdinals(note core.Note, model ankiModel) []int {
	if model.Type == 1 {
		nums := clozeOrdinals(note.Front)
		if len(nums) == 0 {
			return []int{0}
		}
		return nums
	}
	ords := make([]int, len(model.Templates))
	for i := range model.Templates {
		ords[i] = i
	}
	return ords
}

// ankiGUID maps our note id onto a stable per-note string. Anki only requires
// uniqueness, and reusing our own id keeps re-exports idempotent.
func ankiGUID(noteID string) string {
	if strings.TrimSpace(noteID) == "" {
		return fmt.Sprintf("dtui-%d", time.Now().UnixNano())
	}
	return noteID
}

// ankiMediaRef resolves a note's audio value to a file on disk, if it is one.
func ankiMediaRef(audio string) (mediaRef, bool) {
	path := strings.TrimSpace(audio)
	if path == "" {
		return mediaRef{}, false
	}
	if name, ok := soundReference(path); ok {
		path = name
	}
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		return mediaRef{}, false
	}
	return mediaRef{name: filepath.Base(path), path: path}, true
}

func copyIntoZip(zw *zip.Writer, name, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening %s: %w", src, err)
	}
	defer in.Close()

	out, err := zw.Create(name)
	if err != nil {
		return fmt.Errorf("creating zip entry %s: %w", name, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("writing zip entry %s: %w", name, err)
	}
	return nil
}

func itoa(v int64) string { return strconv.FormatInt(v, 10) }
