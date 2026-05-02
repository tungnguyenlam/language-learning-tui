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
		if len(record) > 6 {
			noteType := strings.TrimSpace(record[6])
			if strings.HasPrefix(noteType, "MCQ:") {
				choicesStr := strings.TrimPrefix(noteType, "MCQ:")
				note.Choices = strings.Split(choicesStr, ",")
				for i := range note.Choices {
					note.Choices[i] = strings.TrimSpace(note.Choices[i])
				}
			}
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
		noteType := "Basic"
		if len(note.Choices) > 0 {
			noteType = "MCQ:" + strings.Join(note.Choices, ",")
		}
		if err := writer.Write([]string{
			note.ID,
			note.Front,
			note.Back,
			note.Extra,
			strings.Join(note.Tags, " "),
			note.DeckID,
			noteType,
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

	// Cloze Deletion
	clozes := parseClozes(note.Front)
	if len(clozes) > 0 {
		cards := make([]core.Card, 0, len(clozes))
		for i, cloze := range clozes {
			cards = append(cards, core.Card{
				ID:      fmt.Sprintf("%s:cloze-%d", note.ID, i+1),
				NoteID:  note.ID,
				DeckID:  note.DeckID,
				Kind:    core.CardKindCloze,
				Prompt:  cloze.Prompt,
				Answer:  cloze.Full,             // Full sentence for Cloze
				Choices: []string{cloze.Answer}, // The actual missing part
				Audio:   note.Audio,
				Tags:    baseTags,
			})
		}
		return cards
	}

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

	if len(note.Choices) > 0 {
		cards = append(cards, core.Card{
			ID:      note.ID + ":mcq",
			NoteID:  note.ID,
			DeckID:  note.DeckID,
			Kind:    core.CardKindMCQ,
			Prompt:  note.Front,
			Answer:  note.Back,
			Choices: note.Choices,
			Audio:   note.Audio,
			Tags:    baseTags,
		})
	} else if len(note.Examples) > 0 {
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

type clozeInfo struct {
	Prompt string
	Answer string
	Full   string
}

func parseClozes(text string) []clozeInfo {
	var result []clozeInfo
	// Simple regex-like parsing for {{c1::text::hint}}
	// We'll support multiple clozes
	originalText := text
	offset := 0
	for {
		start := strings.Index(text, "{{c")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "}}")
		if end == -1 {
			break
		}
		end += start
		content := text[start+2 : end]
		parts := strings.Split(content, "::")
		if len(parts) < 2 {
			// Not a valid cloze, skip
			text = text[end+2:]
			offset += end + 2
			continue
		}
		// parts[0] is c1, c2, etc.
		answer := parts[1]
		hint := ""
		if len(parts) > 2 {
			hint = parts[2]
		}

		// Replace the cloze with [...] or [hint] for the prompt
		placeholder := "[...]"
		if hint != "" {
			placeholder = "[" + hint + "]"
		}

		// Create the prompt by replacing ONLY this cloze with placeholder
		// and ALL other clozes with their text
		prompt := stripClozeMarkers(originalText, offset+start, offset+end, placeholder)
		full := stripClozeMarkers(originalText, -1, -1, "")

		result = append(result, clozeInfo{
			Prompt: prompt,
			Answer: answer,
			Full:   full,
		})

		// Continue searching in the rest of the text
		text = text[end+2:]
		offset += end + 2
	}
	return result
}

func stripClozeMarkers(text string, activeStart, activeEnd int, placeholder string) string {
	var b strings.Builder
	i := 0
	for i < len(text) {
		if i == activeStart {
			b.WriteString(placeholder)
			i = activeEnd + 2
			continue
		}
		if strings.HasPrefix(text[i:], "{{c") {
			end := strings.Index(text[i:], "}}")
			if end != -1 {
				end += i
				content := text[i+2 : end]
				parts := strings.Split(content, "::")
				if len(parts) >= 2 {
					b.WriteString(parts[1]) // Keep the answer text
					i = end + 2
					continue
				}
			}
		}
		b.WriteByte(text[i])
		i++
	}
	return b.String()
}

func splitTags(raw string) []string {
	fields := strings.Fields(strings.ReplaceAll(raw, ",", " "))
	if len(fields) == 0 {
		return nil
	}
	return fields
}
