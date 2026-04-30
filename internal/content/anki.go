package content

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"deutsch-tui/internal/core"
)

const exportHeader = "#separator:tab\n#html:false\n#columns:id\tfront\tback\textra\ttags\tdeck\tnotetype\taudio\n"

type ImportOptions struct {
	DefaultDeck string
}

func ImportAnkiTSV(r io.Reader, opts ImportOptions) ([]core.Note, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	deck := opts.DefaultDeck
	lines := strings.Split(string(raw), "\n")
	data := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, "#deck:") {
			deck = strings.TrimSpace(strings.TrimPrefix(line, "#deck:"))
			continue
		}
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		data = append(data, line)
	}
	if deck == "" {
		deck = "Imported"
	}

	reader := csv.NewReader(strings.NewReader(strings.Join(data, "\n")))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = false

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	notes := make([]core.Note, 0, len(records))
	for i, record := range records {
		if len(record) < 3 {
			return nil, fmt.Errorf("record %d has %d fields, want at least 3", i+1, len(record))
		}
		id := strings.TrimSpace(record[0])
		if id == "" {
			return nil, fmt.Errorf("record %d has empty id", i+1)
		}
		noteDeck := deck
		if len(record) > 5 && strings.TrimSpace(record[5]) != "" {
			noteDeck = strings.TrimSpace(record[5])
		}
		note := core.Note{
			ID:     id,
			DeckID: noteDeck,
			Front:  record[1],
			Back:   record[2],
		}
		if len(record) > 3 {
			note.Extra = record[3]
		}
		if len(record) > 4 {
			note.Tags = splitTags(record[4])
		}
		if len(record) > 7 {
			note.Audio = strings.TrimSpace(record[7])
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	return notes, nil
}

func ExportAnkiTSV(w io.Writer, notes []core.Note) error {
	if _, err := io.WriteString(w, exportHeader); err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	writer.Comma = '\t'
	for _, note := range notes {
		if strings.TrimSpace(note.ID) == "" {
			return errors.New("cannot export note with empty id")
		}
		if err := writer.Write([]string{
			note.ID,
			note.Front,
			note.Back,
			note.Extra,
			strings.Join(note.Tags, " "),
			note.DeckID,
			"Basic",
			note.Audio,
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func ExportAnkiTSVString(notes []core.Note) (string, error) {
	var buf bytes.Buffer
	err := ExportAnkiTSV(&buf, notes)
	return buf.String(), err
}

func CardsForNote(note core.Note) []core.Card {
	baseTags := append([]string(nil), note.Tags...)
	cards := []core.Card{
		{
			ID:     note.ID + ":front",
			NoteID: note.ID,
			DeckID: note.DeckID,
			Kind:   core.CardKindFlashcard,
			Prompt: note.Front,
			Answer: note.Back,
			Audio:  note.Audio,
			Tags:   baseTags,
		},
	}
	if len(note.Examples) > 0 {
		cards = append(cards, core.Card{
			ID:      note.ID + ":mcq",
			NoteID:  note.ID,
			DeckID:  note.DeckID,
			Kind:    core.CardKindMCQ,
			Prompt:  note.Examples[0],
			Answer:  note.Back,
			Choices: []string{note.Back, note.Front},
			Audio:   note.Audio,
			Tags:    baseTags,
		})
	}
	return cards
}

func splitTags(raw string) []string {
	fields := strings.Fields(strings.ReplaceAll(raw, ",", " "))
	if len(fields) == 0 {
		return nil
	}
	return fields
}
