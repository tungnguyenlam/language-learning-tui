package content

import (
	"deutsch-tui/internal/core"
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
	if notes[0].DeckID != deck.Notes[0].DeckID {
		t.Fatalf("deck id = %q, want %q", notes[0].DeckID, deck.Notes[0].DeckID)
	}
}

func TestImportExportMCQ(t *testing.T) {
	input := "#separator:tab\n#deck:Grammar\nid-1\tIch gebe ___ das Buch.\tdem Mann\tDativ\tgrammar\tGrammar\tMCQ:den Mann,der Mann,dem Mann\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]
	if len(note.Choices) != 3 || note.Choices[2] != "dem Mann" {
		t.Fatalf("unexpected choices: %v", note.Choices)
	}

	mcqCard := false
	for _, card := range note.Cards {
		if card.Kind == core.CardKindMCQ {
			mcqCard = true
			if len(card.Choices) != 3 || card.Choices[2] != "dem Mann" {
				t.Fatalf("unexpected card choices: %v", card.Choices)
			}
		}
	}
	if !mcqCard {
		t.Fatal("no MCQ card generated")
	}

	// Test round-trip export
	out, err := ExportAnkiTSVString(notes)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if !strings.Contains(out, "MCQ:den Mann,der Mann,dem Mann") {
		t.Fatalf("exported TSV missing MCQ choices: %q", out)
	}
}

func TestImportCloze(t *testing.T) {
	input := "#separator:tab\n#deck:Cloze\nid-1\tIch {{c1::gebe}} dem Mann das Buch.\tgive\t\t\t\t\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]

	clozeCard := false
	for _, card := range note.Cards {
		if card.Kind == core.CardKindCloze {
			clozeCard = true
			if card.Prompt != "Ich [...] dem Mann das Buch." {
				t.Fatalf("unexpected prompt: %q", card.Prompt)
			}
			if card.Answer != "Ich gebe dem Mann das Buch." {
				t.Fatalf("unexpected answer: %q", card.Answer)
			}
			if len(card.Choices) != 1 || card.Choices[0] != "gebe" {
				t.Fatalf("unexpected choices: %v", card.Choices)
			}
		}
	}
	if !clozeCard {
		t.Fatal("no Cloze card generated")
	}
}

func TestImportMultiCloze(t *testing.T) {
	input := "#separator:tab\n#deck:Cloze\nid-1\t{{c1::Ich}} {{c2::gebe}} dem Mann das Buch.\tgive\t\t\t\t\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	note := notes[0]

	clozeCount := 0
	for _, card := range note.Cards {
		if card.Kind == core.CardKindCloze {
			clozeCount++
			if clozeCount == 1 {
				if card.Prompt != "[...] gebe dem Mann das Buch." {
					t.Fatalf("unexpected prompt 1: %q", card.Prompt)
				}
			} else if clozeCount == 2 {
				if card.Prompt != "Ich [...] dem Mann das Buch." {
					t.Fatalf("unexpected prompt 2: %q", card.Prompt)
				}
			}
		}
	}
	if clozeCount != 2 {
		t.Fatalf("clozeCount = %d, want 2", clozeCount)
	}
}

func TestImportReverse(t *testing.T) {
	input := "#separator:tab\n#deck:Reverse\nid-1\tder Apfel\tapple\t\t\t\tReverse\n"
	notes, err := ImportAnkiTSV(strings.NewReader(input), ImportOptions{})
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(notes))
	}
	note := notes[0]
	if note.Type != "Reverse" {
		t.Fatalf("note.Type = %q, want Reverse", note.Type)
	}

	if len(note.Cards) != 2 {
		t.Fatalf("len(note.Cards) = %d, want 2", len(note.Cards))
	}

	frontCard := note.Cards[0]
	if frontCard.Prompt != "der Apfel" || frontCard.Answer != "apple" {
		t.Fatalf("unexpected front card: %+v", frontCard)
	}

	backCard := note.Cards[1]
	if backCard.Prompt != "apple" || backCard.Answer != "der Apfel" {
		t.Fatalf("unexpected back card: %+v", backCard)
	}
}

func TestExportReverseRoundTrip(t *testing.T) {
	note := core.Note{
		ID:     "rev-1",
		DeckID: "Default",
		Type:   "Reverse",
		Front:  "der Apfel",
		Back:   "apple",
	}
	out, err := ExportAnkiTSVString([]core.Note{note})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if !strings.Contains(out, "Reverse") {
		t.Fatalf("exported TSV missing Reverse type: %q", out)
	}

	notes, err := ImportAnkiTSV(strings.NewReader(out), ImportOptions{})
	if err != nil {
		t.Fatalf("reimport failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len(notes) = %d, want 1", len(notes))
	}
	if notes[0].Type != "Reverse" {
		t.Fatalf("notes[0].Type = %q, want Reverse", notes[0].Type)
	}
	if len(notes[0].Cards) != 2 {
		t.Fatalf("len(notes[0].Cards) = %d, want 2", len(notes[0].Cards))
	}
}

func TestExampleCardGeneration(t *testing.T) {
	note := core.Note{
		ID:       "ex-1",
		Front:    "der Apfel",
		Back:     "apple",
		Examples: []string{"Ich esse einen Apfel."},
	}
	cards := CardsForNote(note)

	// Should have front card and example card
	if len(cards) != 2 {
		t.Fatalf("len(cards) = %d, want 2", len(cards))
	}

	exCard := cards[1]
	if exCard.ID != "ex-1:example" {
		t.Fatalf("id = %q, want ex-1:example", exCard.ID)
	}
	if !strings.Contains(exCard.Prompt, "___") {
		t.Fatalf("prompt missing blank: %q", exCard.Prompt)
	}
	if !strings.Contains(exCard.Prompt, "(apple)") {
		t.Fatalf("prompt missing hint: %q", exCard.Prompt)
	}
	if exCard.Answer != "Apfel" {
		t.Fatalf("answer = %q, want Apfel", exCard.Answer)
	}
}
