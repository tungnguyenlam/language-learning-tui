package content

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
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

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading TSV input: %w", err)
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
			note.Audio = cleanAudioField(record[7])
		}
		note.Cards = CardsForNote(note)
		notes = append(notes, note)
	}

	return notes, nil
}

func cleanAudioField(s string) string {
	s = strings.TrimSpace(s)
	// Handle [sound:file.mp3]
	if strings.HasPrefix(s, "[sound:") && strings.HasSuffix(s, "]") {
		return s[7 : len(s)-1]
	}
	return s
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
		Audio:  cleanAudioField(field("audio")),
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
		for _, cloze := range clozes {
			choices := append([]string(nil), cloze.Answers...)
			if len(choices) == 0 && cloze.Answer != "" {
				choices = []string{cloze.Answer}
			}
			cards = append(cards, core.Card{
				ID:      fmt.Sprintf("%s:cloze-%d", note.ID, cloze.Ordinal),
				NoteID:  note.ID,
				DeckID:  note.DeckID,
				Kind:    core.CardKindCloze,
				Prompt:  cloze.Prompt,
				Answer:  cloze.Full, // Full sentence for Cloze
				Extra:   note.Extra,
				Choices: choices, // The missing part(s) for this Anki ordinal
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
	Prompt  string
	Answer  string // first missing part; retained for single-answer callers
	Full    string
	Answers []string
	Ordinal int
}

func parseClozes(text string) []clozeInfo {
	markers := parseClozeMarkers(text)
	if len(markers) == 0 {
		return nil
	}

	groups := make(map[int][]clozeMarker)
	ordinals := make([]int, 0, len(markers))
	for _, marker := range markers {
		if _, seen := groups[marker.ordinal]; !seen {
			ordinals = append(ordinals, marker.ordinal)
		}
		groups[marker.ordinal] = append(groups[marker.ordinal], marker)
	}
	sort.Ints(ordinals)

	full := renderClozeText(text, markers, 0)
	result := make([]clozeInfo, 0, len(ordinals))
	for _, ordinal := range ordinals {
		group := groups[ordinal]
		answers := make([]string, 0, len(group))
		for _, marker := range group {
			answers = append(answers, marker.answer)
		}
		answer := ""
		if len(answers) > 0 {
			answer = answers[0]
		}
		result = append(result, clozeInfo{
			Prompt:  renderClozeText(text, markers, ordinal),
			Answer:  answer,
			Full:    full,
			Answers: answers,
			Ordinal: ordinal,
		})
	}
	return result
}

type clozeMarker struct {
	start   int
	end     int
	ordinal int
	answer  string
	hint    string
}

func parseClozeMarkers(text string) []clozeMarker {
	var markers []clozeMarker
	for searchFrom := 0; searchFrom < len(text); {
		relStart := strings.Index(text[searchFrom:], "{{c")
		if relStart == -1 {
			break
		}
		start := searchFrom + relStart
		relEnd := strings.Index(text[start:], "}}")
		if relEnd == -1 {
			break
		}
		end := start + relEnd
		parts := strings.Split(text[start+2:end], "::")
		if len(parts) >= 2 && strings.HasPrefix(parts[0], "c") {
			ordinal, err := strconv.Atoi(strings.TrimPrefix(parts[0], "c"))
			if err == nil && ordinal > 0 && parts[1] != "" {
				hint := ""
				if len(parts) > 2 {
					hint = parts[2]
				}
				markers = append(markers, clozeMarker{
					start:   start,
					end:     end + 2,
					ordinal: ordinal,
					answer:  parts[1],
					hint:    hint,
				})
			}
		}
		searchFrom = end + 2
	}
	return markers
}

// renderClozeText removes all valid cloze markers. When ordinal is non-zero,
// markers with that ordinal become prompts while other markers reveal their
// answer. Grouping by ordinal matches Anki: {{c1::a}} {{c1::b}} creates one
// card with both blanks, while c2 creates the next card.
func renderClozeText(text string, markers []clozeMarker, ordinal int) string {
	var b strings.Builder
	last := 0
	for _, marker := range markers {
		b.WriteString(text[last:marker.start])
		if ordinal != 0 && marker.ordinal == ordinal {
			placeholder := "[...]"
			if marker.hint != "" {
				placeholder = "[" + marker.hint + "]"
			}
			b.WriteString(placeholder)
		} else {
			b.WriteString(marker.answer)
		}
		last = marker.end
	}
	b.WriteString(text[last:])
	return b.String()
}

func splitTags(raw string) []string {
	fields := strings.Fields(strings.ReplaceAll(raw, ",", " "))
	if len(fields) == 0 {
		return nil
	}
	return fields
}
