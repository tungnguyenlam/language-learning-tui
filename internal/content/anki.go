package content

import (
	"bufio"
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
	deck := opts.DefaultDeck
	deckExplicit := false

	// Read headers and metadata first
	scanner := bufio.NewScanner(r)
	var data strings.Builder
	firstDataRowFound := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#deck:") {
			deck = strings.TrimSpace(strings.TrimPrefix(trimmed, "#deck:"))
			deckExplicit = deck != ""
			continue
		}
		if strings.HasPrefix(trimmed, "#") || trimmed == "" {
			continue
		}
		// Found first non-metadata line
		data.WriteString(line + "\n")
		firstDataRowFound = true
		break
	}

	if firstDataRowFound {
		// Append the rest of the reader to our data builder
		for scanner.Scan() {
			data.WriteString(scanner.Text() + "\n")
		}
	}

	if deck == "" {
		deck = "Imported"
	}

	reader := csv.NewReader(strings.NewReader(data.String()))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = false
	reader.LazyQuotes = true // Helps with some Anki exports

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := []string(nil)
	if len(records) > 0 && isTSVHeader(records[0]) {
		header = normalizeHeader(records[0])
		records = records[1:]
	}

	notes := make([]core.Note, 0, len(records))
	for i, record := range records {
		if len(header) > 0 {
			note, err := noteFromHeaderedRecord(record, header, deck, deckExplicit, i+1)
			if err != nil {
				return nil, err
			}
			notes = append(notes, note)
			continue
		}

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
			Type:   "Basic",
		}
		if len(record) > 3 {
			note.Extra = record[3]
		}
		if len(record) > 4 {
			note.Tags = splitTags(record[4])
		}
		if len(record) > 6 {
			applyNoteType(&note, strings.TrimSpace(record[6]))
		}
		if len(record) > 7 {
			note.Audio = strings.TrimSpace(record[7])
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	return notes, nil
}

func noteFromHeaderedRecord(record []string, header []string, defaultDeck string, deckExplicit bool, row int) (core.Note, error) {
	field := func(name string) string {
		for i, headerName := range headerForRecord(header, record) {
			if headerName == name && i < len(record) {
				return record[i]
			}
		}
		return ""
	}

	front := field("front")
	back := field("back")
	if strings.TrimSpace(front) == "" || strings.TrimSpace(back) == "" {
		return core.Note{}, fmt.Errorf("record %d has empty front/back", row)
	}

	id := strings.TrimSpace(field("id"))
	if id == "" {
		id = generatedNoteID(front, row)
	}

	noteDeck := defaultDeck
	if noteDeck == "" {
		noteDeck = "Imported"
	}
	if deckField := strings.TrimSpace(firstNonEmpty(field("deck"), field("deckid"))); deckField != "" && !deckExplicit {
		noteDeck = deckField
	}

	noteType := strings.TrimSpace(firstNonEmpty(field("notetype"), field("type")))
	if deckField, splitNoteType, ok := splitCombinedDeckNoteType(noteType); ok {
		noteType = splitNoteType
		if !deckExplicit {
			noteDeck = deckField
		}
	}

	note := core.Note{
		ID:     id,
		DeckID: noteDeck,
		Front:  front,
		Back:   back,
		Extra:  field("extra"),
		Tags:   splitTags(field("tags")),
		Type:   "Basic",
		Audio:  strings.TrimSpace(field("audio")),
	}

	applyNoteType(&note, noteType)
	note.Cards = CardsForNote(note)
	return note, nil
}

func applyNoteType(note *core.Note, noteType string) {
	if noteType == "" {
		return
	}
	if strings.HasPrefix(noteType, "MCQ:") {
		choicesStr := strings.TrimPrefix(noteType, "MCQ:")
		choices := strings.Split(choicesStr, "|||")
		if len(choices) == 1 {
			choices = strings.Split(choicesStr, ",")
		}
		note.Choices = note.Choices[:0]
		for _, choice := range choices {
			choice = strings.TrimSpace(choice)
			if choice != "" {
				note.Choices = append(note.Choices, choice)
			}
		}
		if len(note.Choices) < 2 {
			note.Type = "Basic"
			note.Choices = nil
			return
		}
		note.Type = "MCQ"
		return
	}
	if noteType == "Reverse" || noteType == "Basic (and reversed card)" || noteType == "BasicReversed" {
		note.Type = "Reverse"
		return
	}
	note.Type = noteType
}

func splitCombinedDeckNoteType(value string) (string, string, bool) {
	for _, marker := range []string{" MCQ:", " Reverse", " Basic", " Cloze"} {
		if idx := strings.LastIndex(value, marker); idx > 0 {
			return strings.TrimSpace(value[:idx]), strings.TrimSpace(value[idx+1:]), true
		}
	}
	return "", "", false
}

func isTSVHeader(record []string) bool {
	if len(record) == 0 {
		return false
	}
	known := 0
	for _, field := range record {
		switch normalizeHeaderName(field) {
		case "id", "deck", "deckid", "front", "back", "extra", "tags", "notetype", "type", "audio":
			known++
		}
	}
	return known >= 2 && known == len(record)
}

func normalizeHeader(record []string) []string {
	header := make([]string, len(record))
	for i, field := range record {
		header[i] = normalizeHeaderName(field)
	}
	return header
}

func normalizeHeaderName(field string) string {
	field = strings.ToLower(strings.TrimSpace(field))
	field = strings.ReplaceAll(field, "_", "")
	field = strings.ReplaceAll(field, " ", "")
	return field
}

func headerForRecord(header []string, record []string) []string {
	if len(record) == len(header)+1 && !containsHeader(header, "deck") && len(header) > 0 && header[len(header)-1] == "notetype" {
		expanded := make([]string, 0, len(header)+1)
		expanded = append(expanded, header[:len(header)-1]...)
		expanded = append(expanded, "deck", "notetype")
		return expanded
	}
	return header
}

func containsHeader(header []string, name string) bool {
	for _, headerName := range header {
		if headerName == name {
			return true
		}
	}
	return false
}

func generatedNoteID(front string, row int) string {
	id := strings.TrimSpace(front)
	if id != "" {
		return id
	}
	return fmt.Sprintf("note-%d", row)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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
			noteType = "MCQ:" + strings.Join(note.Choices, "|||")
		} else if note.Type == "Reverse" {
			noteType = "Reverse"
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
				Answer:  cloze.Full, // Full sentence for Cloze
				Extra:   note.Extra,
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
			Extra:  note.Extra,
			Audio:  note.Audio,
			Tags:   baseTags,
		},
	}

	if note.Type == "Reverse" {
		cards = append(cards, core.Card{
			ID:     note.ID + ":back",
			NoteID: note.ID,
			DeckID: note.DeckID,
			Kind:   core.CardKindFlashcard,
			Prompt: note.Back,
			Answer: note.Front,
			Extra:  note.Extra,
			Audio:  note.Audio,
			Tags:   baseTags,
		})
	}

	if len(note.Choices) > 0 {
		cards = append(cards, core.Card{
			ID:      note.ID + ":mcq",
			NoteID:  note.ID,
			DeckID:  note.DeckID,
			Kind:    core.CardKindMCQ,
			Prompt:  note.Front,
			Answer:  note.Back,
			Extra:   note.Extra,
			Choices: note.Choices,
			Audio:   note.Audio,
			Tags:    baseTags,
		})
	} else if len(note.Examples) > 0 {
		example := note.Examples[0]
		cleanFront := stripArticles(note.Front)
		start, end := findWordInSentence(example, cleanFront)

		var prompt, answer string
		var choices []string

		if start != -1 {
			actualWord := example[start:end]
			prompt = example[:start] + "___" + example[end:]
			prompt = fmt.Sprintf("%s (%s)", prompt, note.Back)
			answer = actualWord
			choices = []string{actualWord, note.Back}
		} else {
			prompt = example
			answer = note.Front
			choices = []string{note.Front, note.Back}
		}

		cards = append(cards, core.Card{
			ID:      note.ID + ":example",
			NoteID:  note.ID,
			DeckID:  note.DeckID,
			Kind:    core.CardKindMCQ,
			Prompt:  prompt,
			Answer:  answer,
			Extra:   note.Extra,
			Choices: choices,
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
