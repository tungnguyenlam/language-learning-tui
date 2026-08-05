package content

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"deutsch-tui/internal/core"

	"github.com/klauspost/compress/zstd"
	_ "modernc.org/sqlite"
)

// collectionEntries are the collection filenames Anki has used, newest first.
// Modern Anki (2.1.50+) writes collection.anki21b — the same SQLite database
// compressed with zstd. Anki 2.1 wrote collection.anki21 next to a legacy
// collection.anki2 stub, so the stub must be tried last or a shared deck
// imports as empty.
var collectionEntries = []string{
	"collection.anki21b",
	"collection.anki21",
	"collection.anki2",
}

// ImportAnkiAPKG reads notes from an Anki `.apkg` package.
//
// Deck names, note types and tags are preserved. Anki stores fields as HTML, so
// markup is converted to plain text and `[sound:…]` references are lifted onto
// the note's Audio field.
func ImportAnkiAPKG(r io.Reader) ([]core.Note, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading package: %w", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("reading package: not a valid .apkg zip: %w", err)
	}

	entries := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		entries[f.Name] = f
	}

	var collection []byte
	var chosen string
	for _, name := range collectionEntries {
		f, ok := entries[name]
		if !ok {
			continue
		}
		raw, err := readZipEntry(f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", name, err)
		}
		if name == "collection.anki21b" {
			raw, err = decompressZstd(raw)
			if err != nil {
				return nil, fmt.Errorf("decompressing %s: %w", name, err)
			}
		}
		collection, chosen = raw, name
		break
	}
	if collection == nil {
		return nil, fmt.Errorf("no Anki collection found in package (looked for %s)",
			strings.Join(collectionEntries, ", "))
	}

	tmpDir, err := os.MkdirTemp("", "deutsch-tui-apkg-import-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "collection.sqlite")
	if err := os.WriteFile(dbPath, collection, 0o600); err != nil {
		return nil, fmt.Errorf("staging collection: %w", err)
	}

	notes, err := readAnkiCollection(dbPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", chosen, err)
	}
	return notes, nil
}

// ImportAnkiAPKGFromFile imports an APKG from a file path.
func ImportAnkiAPKGFromFile(path string) ([]core.Note, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportAnkiAPKG(f)
}

func readAnkiCollection(dbPath string) ([]core.Note, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening collection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	deckNames := readAnkiDeckNames(ctx, db)
	noteTypes := readAnkiNoteTypes(ctx, db)
	noteDecks := readAnkiNoteDecks(ctx, db)

	rows, err := db.QueryContext(ctx, `SELECT id, guid, mid, flds, tags FROM notes`)
	if err != nil {
		return nil, fmt.Errorf("querying notes: %w", err)
	}
	defer rows.Close()

	var notes []core.Note
	for rows.Next() {
		var id, mid int64
		var guid, flds, tags string
		if err := rows.Scan(&id, &guid, &mid, &flds, &tags); err != nil {
			return nil, fmt.Errorf("scanning note: %w", err)
		}

		raw := strings.Split(flds, "\x1f")
		texts := make([]string, len(raw))
		audio := ""
		for i, f := range raw {
			text, a := stripHTMLFieldWithMedia(f)
			texts[i] = text
			if audio == "" {
				audio = a
			}
		}

		nt := noteTypes[mid]
		frontOrd, backOrd, extraOrds := assignFieldRoles(texts, nt)
		if frontOrd < 0 {
			// Anki collections routinely carry empty scratch notes; one of them
			// must not abort an otherwise good import.
			continue
		}

		note := core.Note{
			ID:     noteIdentity(guid, id),
			DeckID: deckNames[noteDecks[id]],
			Front:  texts[frontOrd],
			Audio:  audio,
			Tags:   strings.Fields(strings.ReplaceAll(tags, "\x1f", " ")),
		}
		if backOrd >= 0 {
			note.Back = texts[backOrd]
		}
		// Anki note types vary wildly in field count; whatever the templates
		// do not use as a side is kept as the reveal-time extra.
		var extras []string
		for _, ord := range extraOrds {
			extras = append(extras, texts[ord])
		}
		note.Extra = strings.Join(extras, "\n")

		if nt.cloze {
			// Keep the {{c1::…}} markers so the cloze parser can see them.
			note.Front, _ = stripHTMLFieldWithMedia(raw[frontOrd])
		} else if nt.reversed {
			note.Type = "Reverse"
		}

		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating notes: %w", err)
	}

	return notes, nil
}

// assignFieldRoles decides which of a note's fields become the prompt, the
// answer, and the extra. It returns (-1, …) when the note has no usable text.
//
// The note type's templates win when they name fields that this note actually
// filled in; otherwise the first two non-empty fields are used, which is what
// makes decks with unconventional field layouts import as usable cards rather
// than as prompts with a blank answer.
func assignFieldRoles(texts []string, nt ankiNoteTypeInfo) (int, int, []int) {
	filled := func(ord int) bool {
		return ord >= 0 && ord < len(texts) && strings.TrimSpace(texts[ord]) != ""
	}

	front, back := -1, -1
	if filled(nt.frontOrd) {
		front = nt.frontOrd
	}
	if filled(nt.backOrd) && nt.backOrd != front {
		back = nt.backOrd
	}

	for ord := range texts {
		if front >= 0 && back >= 0 {
			break
		}
		if !filled(ord) || ord == front || ord == back {
			continue
		}
		if front < 0 {
			front = ord
		} else if back < 0 {
			back = ord
		}
	}
	if front < 0 {
		return -1, -1, nil
	}

	var extras []int
	for ord := range texts {
		if ord != front && ord != back && filled(ord) {
			extras = append(extras, ord)
		}
	}
	return front, back, extras
}

// noteIdentity prefers Anki's stable guid so re-importing an updated package
// updates the same notes instead of duplicating them.
func noteIdentity(guid string, id int64) string {
	if g := strings.TrimSpace(guid); g != "" {
		return g
	}
	return strconv.FormatInt(id, 10)
}

type ankiNoteTypeInfo struct {
	cloze    bool
	reversed bool
	// fields is the note type's field list, in note order. frontOrd/backOrd
	// are the field ordinals the first card template puts on the question and
	// answer side, or -1 when they could not be determined.
	fields   []string
	frontOrd int
	backOrd  int
}

// fieldRefRe matches a template's field references, e.g. {{Front}},
// {{cloze:Text}}, {{type:Back}}, {{#Extra}}.
var fieldRefRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// templateFieldOrds resolves which fields a note type's first template shows on
// the question and answer side.
//
// Anki note types are free-form: a deck may leave field 1 empty and put the
// meaning in field 4. Assuming field 0/1 are the two sides produces cards with
// a blank answer, so follow the template the way Anki does.
func templateFieldOrds(fields []string, qfmt, afmt string) (int, int) {
	ord := make(map[string]int, len(fields))
	for i, name := range fields {
		ord[strings.ToLower(name)] = i
	}

	refs := func(tmpl string) []int {
		var out []int
		for _, m := range fieldRefRe.FindAllStringSubmatch(tmpl, -1) {
			name := strings.TrimSpace(m[1])
			// Drop section markers and filters: {{#Tags}}, {{cloze:Text}}.
			name = strings.TrimLeft(name, "#^/")
			if i := strings.LastIndex(name, ":"); i >= 0 {
				name = name[i+1:]
			}
			name = strings.TrimSpace(name)
			if i, ok := ord[strings.ToLower(name)]; ok {
				out = append(out, i)
			}
		}
		return out
	}

	front, back := -1, -1
	qOrds := refs(qfmt)
	if len(qOrds) > 0 {
		front = qOrds[0]
	}
	onQuestion := make(map[int]bool, len(qOrds))
	for _, o := range qOrds {
		onQuestion[o] = true
	}
	for _, o := range refs(afmt) {
		if !onQuestion[o] {
			back = o
			break
		}
	}
	return front, back
}

// readAnkiNoteTypes reads note types from either schema: the modern `notetypes`
// table, or the legacy `col.models` JSON blob.
func readAnkiNoteTypes(ctx context.Context, db *sql.DB) map[int64]ankiNoteTypeInfo {
	types := make(map[int64]ankiNoteTypeInfo)

	if rows, err := db.QueryContext(ctx, `SELECT id, name FROM notetypes`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int64
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				continue
			}
			types[id] = noteTypeInfoFromName(name)
		}
		// The modern schema keeps template counts in a separate table; a name
		// check is enough to distinguish the stock reversed and cloze types.
		if len(types) > 0 {
			return types
		}
	}

	var modelsJSON string
	if err := db.QueryRowContext(ctx, `SELECT models FROM col LIMIT 1`).Scan(&modelsJSON); err != nil {
		return types
	}
	var models map[string]struct {
		Name      string `json:"name"`
		Type      int    `json:"type"`
		Templates []struct {
			Name string `json:"name"`
			QFmt string `json:"qfmt"`
			AFmt string `json:"afmt"`
		} `json:"tmpls"`
		Fields []struct {
			Name string `json:"name"`
			Ord  int    `json:"ord"`
		} `json:"flds"`
	}
	if err := json.Unmarshal([]byte(modelsJSON), &models); err != nil {
		return types
	}
	for key, model := range models {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			continue
		}
		info := noteTypeInfoFromName(model.Name)
		if model.Type == 1 {
			info.cloze = true
		}
		if len(model.Templates) > 1 {
			info.reversed = true
		}

		info.fields = make([]string, len(model.Fields))
		for _, f := range model.Fields {
			if f.Ord >= 0 && f.Ord < len(info.fields) {
				info.fields[f.Ord] = f.Name
			}
		}
		if len(model.Templates) > 0 {
			info.frontOrd, info.backOrd = templateFieldOrds(
				info.fields, model.Templates[0].QFmt, model.Templates[0].AFmt)
		}
		types[id] = info
	}
	return types
}

func noteTypeInfoFromName(name string) ankiNoteTypeInfo {
	lower := strings.ToLower(name)
	return ankiNoteTypeInfo{
		cloze:    strings.Contains(lower, "cloze"),
		reversed: strings.Contains(lower, "reversed") || strings.Contains(lower, "reverse"),
		frontOrd: -1,
		backOrd:  -1,
	}
}

// readAnkiDeckNames reads deck names from either schema: the modern `decks`
// table, or the legacy `col.decks` JSON blob.
func readAnkiDeckNames(ctx context.Context, db *sql.DB) map[int64]string {
	names := make(map[int64]string)

	if rows, err := db.QueryContext(ctx, `SELECT id, name FROM decks`); err == nil {
		for rows.Next() {
			var id int64
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				continue
			}
			// The modern schema separates sub-deck components with \x1f.
			names[id] = strings.ReplaceAll(name, "\x1f", "::")
		}
		_ = rows.Close()
		if len(names) > 0 {
			return names
		}
	}

	var decksJSON string
	if err := db.QueryRowContext(ctx, `SELECT decks FROM col LIMIT 1`).Scan(&decksJSON); err != nil {
		return names
	}
	var decks map[string]struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(decksJSON), &decks); err != nil {
		return names
	}
	for _, d := range decks {
		names[d.ID] = d.Name
	}
	return names
}

// readAnkiNoteDecks maps each note to a deck. A note's cards can live in
// different decks; the lowest deck id wins so the choice is deterministic.
func readAnkiNoteDecks(ctx context.Context, db *sql.DB) map[int64]int64 {
	assigned := make(map[int64]int64)
	rows, err := db.QueryContext(ctx, `SELECT nid, did FROM cards`)
	if err != nil {
		return assigned
	}
	defer rows.Close()
	for rows.Next() {
		var nid, did int64
		if err := rows.Scan(&nid, &did); err != nil {
			continue
		}
		if cur, ok := assigned[nid]; !ok || did < cur {
			assigned[nid] = did
		}
	}
	return assigned
}

func readZipEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func decompressZstd(data []byte) ([]byte, error) {
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer dec.Close()
	return dec.DecodeAll(data, nil)
}

// --- HTML field handling -----------------------------------------------------

var (
	htmlBreakRe  = regexp.MustCompile(`(?i)<\s*(br|/div|/p|/li|/tr)\s*/?\s*>`)
	htmlTagRe    = regexp.MustCompile(`(?s)<[^>]*>`)
	soundRe      = regexp.MustCompile(`\[sound:([^\]]+)\]`)
	whitespaceRe = regexp.MustCompile(`[ \t]+`)
	clozeIndexRe = regexp.MustCompile(`\{\{c(\d+)::`)
)

// stripHTMLField converts an Anki field to plain text.
func stripHTMLField(field string) string {
	text, _ := stripHTMLFieldWithMedia(field)
	return text
}

// PlainTextFromHTML renders HTML supplied by Anki or AnkiWeb as plain text
// suitable for a terminal. Shared deck descriptions are HTML written by
// strangers, so it strips markup rather than trying to render it.
func PlainTextFromHTML(s string) string {
	return stripHTMLField(s)
}

// stripHTMLFieldWithMedia converts an Anki field to plain text and returns the
// first `[sound:…]` filename it referenced, if any.
func stripHTMLFieldWithMedia(field string) (string, string) {
	audio := ""
	if m := soundRe.FindStringSubmatch(field); len(m) == 2 {
		audio = m[1]
	}
	field = soundRe.ReplaceAllString(field, "")
	return collapseFieldText(htmlTagRe.ReplaceAllString(htmlBreakRe.ReplaceAllString(field, "\n"), "")), audio
}

// stripHTMLKeepingCloze strips markup without disturbing `{{c1::…}}` markers,
// which live in the text rather than in tags.
func stripHTMLKeepingCloze(field string) (string, string) {
	return stripHTMLFieldWithMedia(field)
}

func collapseFieldText(s string) string {
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, " ", " ")
	s = whitespaceRe.ReplaceAllString(s, " ")

	lines := strings.Split(s, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// soundReference reports whether a value is an Anki `[sound:…]` reference and
// returns the filename it points at.
func soundReference(s string) (string, bool) {
	if m := soundRe.FindStringSubmatch(s); len(m) == 2 {
		return m[1], true
	}
	return "", false
}

// clozeOrdinals returns the zero-based template ordinals for the cloze numbers
// present in a text, in ascending order: `{{c1::a}} {{c3::b}}` yields [0 2].
func clozeOrdinals(text string) []int {
	seen := make(map[int]bool)
	for _, m := range clozeIndexRe.FindAllStringSubmatch(text, -1) {
		n, err := strconv.Atoi(m[1])
		if err != nil || n < 1 {
			continue
		}
		seen[n-1] = true
	}
	ords := make([]int, 0, len(seen))
	for ord := range seen {
		ords = append(ords, ord)
	}
	sort.Ints(ords)
	return ords
}
