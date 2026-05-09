package content

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	"database/sql"
	_ "modernc.org/sqlite"
)

// ExportAnkiAPKG exports notes to an APKG file.
// The APKG file contains a SQLite database with the notes, cards, and media.
// Only basic flashcard and MCQ card types are supported for now.
// Progress data (reviews, scheduling) is not exported.
func ExportAnkiAPKG(w io.Writer, notes []core.Note) error {
	// Create a temporary directory for our APKG contents
	tmpDir, err := os.MkdirTemp("", "deutsch-tui-apkg-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the SQLite database
	dbPath := filepath.Join(tmpDir, "collection.anki2")
	if err := createAnkiDatabase(dbPath, notes); err != nil {
		return fmt.Errorf("creating database: %w", err)
	}

	// Create media folder and copy any audio files
	mediaDir := filepath.Join(tmpDir, "media")
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		return fmt.Errorf("creating media dir: %w", err)
	}
	if err := copyMediaFiles(mediaDir, notes); err != nil {
		return fmt.Errorf("copying media: %w", err)
	}

	// Create the media.json file
	if err := createMediaJson(mediaDir, notes); err != nil {
		return fmt.Errorf("creating media json: %w", err)
	}

	// Create the APKG zip file
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Add the database file
	dbFile, err := zipWriter.Create("collection.anki2")
	if err != nil {
		return fmt.Errorf("creating database entry in zip: %w", err)
	}
	if err := addFileToZip(dbPath, dbFile); err != nil {
		return fmt.Errorf("adding database to zip: %w", err)
	}

	// Add media files
	if err := addMediaToZip(mediaDir, zipWriter); err != nil {
		return fmt.Errorf("adding media to zip: %w", err)
	}

	return nil
}

// ImportAnkiAPKG imports notes from an APKG file.
func ImportAnkiAPKG(r io.Reader) ([]core.Note, error) {
	// Create a temporary directory to extract the APKG
	tmpDir, err := os.MkdirTemp("", "deutsch-tui-apkg-import-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Extract the APKG zip file
	if err := extractZip(r, tmpDir); err != nil {
		return nil, fmt.Errorf("extracting apkg: %w", err)
	}

	// Read the SQLite database
	dbPath := filepath.Join(tmpDir, "collection.anki2")
	notes, err := readAnkiDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("reading database: %w", err)
	}

	// Read media files (for now we just note their existence)
	// In a full implementation, we would copy them to our media directory
	// For this implementation, we'll just ensure the notes have the correct audio paths

	return notes, nil
}

// readAnkiDatabase reads notes from an Anki SQLite database.
func readAnkiDatabase(dbPath string) ([]core.Note, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Query notes and their fields
	rows, err := db.QueryContext(context.Background(), `
		SELECT guid, flds, tags
		FROM notes
		WHERE guid != 'genesis'
	`)
	if err != nil {
		return nil, fmt.Errorf("querying notes: %w", err)
	}
	defer rows.Close()

	var notes []core.Note
	for rows.Next() {
		var guid, flds, tags string
		if err := rows.Scan(&guid, &flds, &tags); err != nil {
			return nil, fmt.Errorf("scanning note: %w", err)
		}

		// Split fields by unit separator (\x1f)
		fields := strings.Split(flds, "\x1f")
		if len(fields) < 2 {
			return nil, fmt.Errorf("note %s has insufficient fields", guid)
		}

		// Parse tags (space-separated string)
		var tagList []string
		if tags != "" && tags != "_" {
			tagList = strings.Fields(tags)
		}

		note := core.Note{
			ID:    guid,
			Front: fields[0],
			Back:  fields[1],
			Tags:  tagList,
		}
		if len(fields) > 2 {
			note.Extra = fields[2]
		}

		// Generate cards for the note
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating notes: %w", err)
	}

	return notes, nil
}

// Helper functions for APKG creation

func createAnkiDatabase(dbPath string, notes []core.Note) error {
	// Remove any existing database
	os.Remove(dbPath)

	// Open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	// Create the basic Anki 2.1 schema
	if err := createAnkiSchema(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("creating schema: %w", err)
	}

	// Insert collection metadata
	if err := insertCollectionMetadata(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("inserting collection metadata: %w", err)
	}

	// Insert models (note types)
	modelID, err := insertBasicModel(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("inserting basic model: %w", err)
	}

	// Insert deck
	deckID, err := insertDefaultDeck(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("inserting default deck: %w", err)
	}

	// Insert notes and cards
	noteIDs := make(map[string]int64)
	for _, note := range notes {
		noteID, err := insertNote(tx, modelID, deckID, note)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("inserting note %s: %w", note.ID, err)
		}
		noteIDs[note.ID] = noteID

		// Create cards for this note
		if err := insertCardsForNote(tx, noteID, note); err != nil {
			tx.Rollback()
			return fmt.Errorf("inserting cards for note %s: %w", note.ID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func createAnkiSchema(tx *sql.Tx) error {
	// Collection table
	if _, err := tx.Exec(`
		CREATE TABLE col (
			id              integer primary key,
			crt             integer not null,
			mod             integer not null,
			scm             integer not null,
			ver             integer not null,
			dty             integer not null,
			usn             integer not null,
			ls              integer not null,
			conf            text not null,
			models          text not null,
			decks           text not null,
			dconf           text not null,
			tags            text not null
		);
	`); err != nil {
		return err
	}

	// Notes table
	if _, err := tx.Exec(`
		CREATE TABLE notes (
			id              integer primary key,
			guid            text not null,
			mid             integer not null,
			mod             integer not null,
			usn             integer not null,
			tags            text not null,
			flds            text not null,
			sfld            integer not null,
			csum            integer not null,
			flags           integer not null,
			data            text not null
		);
	`); err != nil {
		return err
	}

	// Cards table
	if _, err := tx.Exec(`
		CREATE TABLE cards (
			id              integer primary key,
			nid             integer not null,
			did             integer not null,
			ord             integer not null,
			mod             integer not null,
			usn             integer not null,
			type            integer not null,
			queue           integer not null,
			due             integer not null,
			ivl             integer not null,
			factor          integer not null,
			reps            integer not null,
			lapses          integer not null,
			left            integer not null,
			odue            integer not null,
			odid            integer not null,
			flags           integer not null,
			data            text not null
		);
	`); err != nil {
		return err
	}

	// Create indexes
	if _, err := tx.Exec(`CREATE INDEX ix_notes_usn on notes (usn);`); err != nil {
		return err
	}
	if _, err := tx.Exec(`CREATE INDEX ix_cards_usn on cards (usn);`); err != nil {
		return err
	}
	if _, err := tx.Exec(`CREATE INDEX ix_cards_nid on cards (nid);`); err != nil {
		return err
	}
	if _, err := tx.Exec(`CREATE INDEX ix_cards_sched on cards (did, queue, due);`); err != nil {
		return err
	}

	return nil
}

func insertCollectionMetadata(tx *sql.Tx) error {
	// Default collection creation time (now)
	now := time.Now().Unix()

	// Empty JSON objects for the various JSON fields
	emptyJSON := "{}"

	_, err := tx.ExecContext(context.Background(), `
		INSERT INTO col (
			id, crt, mod, scm, ver, dty, usn, ls, conf, models, decks, dconf, tags
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`, 1, now, now, now, 1, 0, 0, now, emptyJSON, emptyJSON, emptyJSON, emptyJSON, emptyJSON)
	return err
}

func insertBasicModel(tx *sql.Tx) (int64, error) {
	now := time.Now().UnixMilli()

	// Basic model (note type) that has Front and Back fields
	model := map[string]interface{}{
		"id":    int64(1),
		"name":  "Basic",
		"mod":   now,
		"usn":   int64(-1),
		"sortf": int64(0),
		"req":   []int64{0, 1},
		"vers":  []int64{1},
		"fields": []map[string]string{
			{"name": "Front", "ord": "0"},
			{"name": "Back", "ord": "1"},
		},
		"templates": []map[string]string{
			{
				"name": "Card 1",
				"qfmt": "{{Front}}",
				"afmt": "{{FrontSide}}\n<hr id=\"answer\">{{Back}}",
			},
		},
		"css":     ".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}\n",
		"did":     int64(1), // Default deck
		"type":    int64(0),
		"version": int64(1),
	}

	modelsJSON, err := json.Marshal(map[string]interface{}{
		"1": model,
	})
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO notes (
			id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data
		) VALUES (
			0, 'genesis', 1, ?, ?, '', '', 0, 0, 0, ''
		)
	`, now, now)
	if err != nil {
		return 0, err
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Update the models field in col table
	if _, err := tx.Exec(`UPDATE col SET models = ? WHERE id = 1`, string(modelsJSON)); err != nil {
		return 0, err
	}

	return noteID, nil
}

func insertDefaultDeck(tx *sql.Tx) (int64, error) {
	now := time.Now().UnixMilli()

	// Default deck
	deck := map[string]interface{}{
		"id":               1,
		"name":             "default",
		"mod":              now,
		"usn":              -1,
		"desc":             "",
		"conf":             1,
		"dyn":              0,
		"extendNew":        0,
		"extendRev":        0,
		"collapsed":        false,
		"browserCollapsed": false,
	}

	decksJSON, err := json.Marshal(map[string]interface{}{
		"1": deck,
	})
	if err != nil {
		return 0, err
	}

	// Update the decks field in col table
	if _, err := tx.Exec(`UPDATE col SET decks = ? WHERE id = 1`, string(decksJSON)); err != nil {
		return 0, err
	}

	return 1, nil
}

func insertNote(tx *sql.Tx, modelID, deckID int64, note core.Note) (int64, error) {
	now := time.Now().Unix()
	usn := int64(0)

	// Fields are separated by \x1f (unit separator) in Anki
	fields := strings.Join([]string{
		note.Front,
		note.Back,
		note.Extra,
	}, "\x1f")

	// Tags are stored as a string with spaces
	tags := strings.Join(note.Tags, " ")
	if tags == "" {
		tags = "_"
	}

	// Calculate checksum: first 8 characters of SHA1 hash of the first field
	h := sha1.New()
	h.Write([]byte(note.Front))
	sum := h.Sum(nil)
	csumStr := fmt.Sprintf("%x", sum)[:8]
	csum, _ := strconv.ParseUint(csumStr, 16, 64)

	// Generate a numeric ID for the note (use hash of the string ID for simplicity)
	var noteIDInt int64
	fh := fnv.New64a()
	fh.Write([]byte(note.ID))
	noteIDInt = int64(fh.Sum64())

	result, err := tx.Exec(`
		INSERT INTO notes (
			id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`, noteIDInt, note.ID, modelID, now, usn, tags, fields, 0, int64(csum), 0, "")
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func insertCardsForNote(tx *sql.Tx, noteID int64, note core.Note) error {
	now := time.Now().Unix()
	usn := int64(0)
	due := now      // Due today
	ivl := int64(0) // Interval in days
	factor := 2500  // Factor as per mille (2.5 * 1000)
	reps := int64(0)
	lapses := int64(0)

	for i, card := range note.Cards {
		var cardIDInt int64
		h := fnv.New64a()
		h.Write([]byte(card.ID))
		cardIDInt = int64(h.Sum64())

		// ord is the index of the template or cloze index
		ord := i
		if card.Kind == core.CardKindCloze {
			// Extract cloze number if possible, else use i
			// For now just use i
		}

		if _, err := tx.Exec(`
			INSERT INTO cards (
				id, nid, did, ord, mod, usn, type, queue, due, ivl, factor, reps, lapses, left, odue, odid, flags, data
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`, cardIDInt, noteID, 1, ord, now, usn, 0, 0, due, ivl, factor, reps, lapses, 0, 0, 0, 0, ""); err != nil {
			return err
		}
	}

	return nil
}

func generateGUID() string {
	// Simple GUID generation - in practice this should be more robust
	return fmt.Sprintf("%d%x", time.Now().UnixNano(), time.Now().Unix())
}

func copyMediaFiles(mediaDir string, notes []core.Note) error {
	for i, note := range notes {
		if note.Audio != "" {
			src := note.Audio
			// If it's a path, try to copy it
			if _, err := os.Stat(src); err == nil {
				dest := filepath.Join(mediaDir, fmt.Sprintf("%d", i+1))
				data, err := os.ReadFile(src)
				if err != nil {
					continue
				}
				os.WriteFile(dest, data, 0o644)
			}
		}
	}
	return nil
}

func createMediaJson(mediaDir string, notes []core.Note) error {
	// Create a mapping from media file IDs to filenames
	mediaMap := map[string]string{}

	for i, note := range notes {
		if note.Audio != "" {
			mediaMap[fmt.Sprintf("%d", i+1)] = filepath.Base(note.Audio)
		}
	}

	mediaJSON, err := json.Marshal(mediaMap)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(mediaDir, "media.json"), mediaJSON, 0o644)
}

func extractZip(r io.Reader, dest string) error {
	// Read the entire ZIP file into memory first
	zipData, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading zip data: %w", err)
	}

	// Create a reader from the data
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("creating zip reader: %w", err)
	}

	// Extract each file
	for _, f := range zipReader.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for directory traversal
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(src string, dest io.Writer) error {
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	_, err = io.Copy(dest, inFile)
	return err
}

func addMediaToZip(mediaDir string, zipWriter *zip.Writer) error {
	// Add media.json
	mediaJSONPath := filepath.Join(mediaDir, "media.json")
	if _, err := os.Stat(mediaJSONPath); err == nil {
		mediaFile, err := zipWriter.Create("media.json")
		if err != nil {
			return err
		}
		if err := addFileToZip(mediaJSONPath, mediaFile); err != nil {
			return err
		}
	}

	// Add the media directory contents
	files, err := os.ReadDir(mediaDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || file.Name() == "media.json" {
			continue
		}

		filePath := filepath.Join(mediaDir, file.Name())
		zipPath := filepath.Join(file.Name()) // Media files are at root of APKG

		outFile, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}
		if err := addFileToZip(filePath, outFile); err != nil {
			return err
		}
	}

	return nil
}

// ImportAnkiAPKGFromFile imports an APKG from a file path
func ImportAnkiAPKGFromFile(path string) ([]core.Note, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportAnkiAPKG(f)
}

// ExportAnkiAPKGToFile exports notes to an APKG file at the given path
func ExportAnkiAPKGToFile(path string, notes []core.Note) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExportAnkiAPKG(f, notes)
}
